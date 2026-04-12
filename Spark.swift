import SwiftUI
import Foundation

// MARK: - Configuration
struct SparkConfig {
    static let appName = "Spark"
    static let endpoint = "https://spark3-share.tech-2030.net/api/v1/chat/completions"
    static let modelsEndpoint = "https://spark3-share.tech-2030.net/api/v1/models"
    static let apiKey = "f0a26c072bd83b635a4daad40a51be068bb80d5a7540adfe"

    static let mainModel = "gemma4:31b"
    static let fastModel = "gemma4:26b"

    static let systemPrompt = "You are Spark, a helpful AI assistant powered by Gemma4. Answer in the user's language. Be concise and accurate."
}

// MARK: - Models
struct ChatMessage: Identifiable, Equatable {
    let id = UUID()
    let role: String
    var content: String
    let model: String?
    let timestamp: Date
    var isStreaming: Bool

    static func == (lhs: ChatMessage, rhs: ChatMessage) -> Bool {
        lhs.id == rhs.id
    }
}

struct APIResponse: Codable {
    struct Choice: Codable {
        struct Message: Codable {
            let role: String
            let content: String
        }
        let message: Message?
        let delta: Message?
    }
    let choices: [Choice]?
    let model: String?
}

struct ModelsResponse: Codable {
    struct Model: Codable {
        let id: String
    }
    let data: [Model]?
}

// MARK: - API Client
class SparkAPI {
    static func sendChat(messages: [(role: String, content: String)], model: String, stream: Bool = true) async throws -> AsyncThrowingStream<String, Error> {
        let url = URL(string: SparkConfig.endpoint)!
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("Bearer \(SparkConfig.apiKey)", forHTTPHeaderField: "Authorization")
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.timeoutInterval = 120

        let body: [String: Any] = [
            "model": model,
            "messages": messages.map { ["role": $0.role, "content": $0.content] },
            "stream": stream,
            "temperature": 0.7,
            "max_tokens": 4096
        ]
        request.httpBody = try JSONSerialization.data(withJSONObject: body)

        if stream {
            return AsyncThrowingStream { continuation in
                let task = Task {
                    do {
                        let (bytes, response) = try await URLSession.shared.bytes(for: request)
                        guard let httpResponse = response as? HTTPURLResponse, httpResponse.statusCode == 200 else {
                            let statusCode = (response as? HTTPURLResponse)?.statusCode ?? 0
                            continuation.finish(throwing: NSError(domain: "SparkAPI", code: statusCode, userInfo: [NSLocalizedDescriptionKey: "HTTP \(statusCode)"]))
                            return
                        }
                        for try await line in bytes.lines {
                            if line.hasPrefix("data: ") {
                                let jsonStr = String(line.dropFirst(6))
                                if jsonStr == "[DONE]" { break }
                                if let data = jsonStr.data(using: .utf8),
                                   let parsed = try? JSONDecoder().decode(APIResponse.self, from: data),
                                   let delta = parsed.choices?.first?.delta?.content {
                                    continuation.yield(delta)
                                }
                            }
                        }
                        continuation.finish()
                    } catch {
                        continuation.finish(throwing: error)
                    }
                }
                continuation.onTermination = { _ in task.cancel() }
            }
        } else {
            return AsyncThrowingStream { continuation in
                Task {
                    do {
                        let (data, response) = try await URLSession.shared.data(for: request)
                        guard let httpResponse = response as? HTTPURLResponse, httpResponse.statusCode == 200 else {
                            let statusCode = (response as? HTTPURLResponse)?.statusCode ?? 0
                            continuation.finish(throwing: NSError(domain: "SparkAPI", code: statusCode, userInfo: [NSLocalizedDescriptionKey: "HTTP \(statusCode)"]))
                            return
                        }
                        if let parsed = try? JSONDecoder().decode(APIResponse.self, from: data),
                           let content = parsed.choices?.first?.message?.content {
                            continuation.yield(content)
                        }
                        continuation.finish()
                    } catch {
                        continuation.finish(throwing: error)
                    }
                }
            }
        }
    }

    static func fetchModels() async -> [String] {
        guard let url = URL(string: SparkConfig.modelsEndpoint) else { return [] }
        var request = URLRequest(url: url)
        request.setValue("Bearer \(SparkConfig.apiKey)", forHTTPHeaderField: "Authorization")
        do {
            let (data, _) = try await URLSession.shared.data(for: request)
            if let response = try? JSONDecoder().decode(ModelsResponse.self, from: data) {
                return response.data?.map(\.id).filter { $0.hasPrefix("gemma4") }.sorted() ?? []
            }
        } catch {}
        return []
    }
}

// MARK: - ViewModel
@MainActor
class ChatViewModel: ObservableObject {
    @Published var messages: [ChatMessage] = []
    @Published var inputText = ""
    @Published var isLoading = false
    @Published var selectedModel: String = SparkConfig.fastModel
    @Published var availableModels: [String] = [SparkConfig.mainModel, SparkConfig.fastModel]
    @Published var errorMessage: String?
    @Published var useDeepThink = false

    private var currentTask: Task<Void, Never>?

    var modelLabel: String {
        selectedModel == SparkConfig.mainModel ? "31B Deep" : "26B Fast"
    }

    func loadModels() {
        Task {
            let models = await SparkAPI.fetchModels()
            if !models.isEmpty {
                availableModels = models
            }
        }
    }

    func send() {
        let text = inputText.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !text.isEmpty, !isLoading else { return }

        inputText = ""
        errorMessage = nil

        let userMsg = ChatMessage(role: "user", content: text, model: nil, timestamp: Date(), isStreaming: false)
        messages.append(userMsg)

        let model = useDeepThink ? SparkConfig.mainModel : selectedModel

        let assistantMsg = ChatMessage(role: "assistant", content: "", model: model, timestamp: Date(), isStreaming: true)
        messages.append(assistantMsg)
        isLoading = true

        let history: [(role: String, content: String)] = [
            (role: "system", content: SparkConfig.systemPrompt)
        ] + messages.dropLast().suffix(20).map { (role: $0.role, content: $0.content) }

        currentTask = Task {
            do {
                let stream = try await SparkAPI.sendChat(messages: history, model: model)
                for try await chunk in stream {
                    if let idx = messages.lastIndex(where: { $0.id == assistantMsg.id }) {
                        messages[idx].content += chunk
                    }
                }
                if let idx = messages.lastIndex(where: { $0.id == assistantMsg.id }) {
                    messages[idx].isStreaming = false
                }
            } catch {
                errorMessage = "Error: \(error.localizedDescription)"
                if let idx = messages.lastIndex(where: { $0.id == assistantMsg.id }) {
                    if messages[idx].content.isEmpty {
                        messages.remove(at: idx)
                    } else {
                        messages[idx].isStreaming = false
                    }
                }
            }
            isLoading = false
        }
    }

    func stop() {
        currentTask?.cancel()
        currentTask = nil
        isLoading = false
        if let idx = messages.indices.last, messages[idx].isStreaming {
            messages[idx].isStreaming = false
        }
    }

    func clearChat() {
        messages.removeAll()
        errorMessage = nil
    }
}

// MARK: - Views
struct MessageBubble: View {
    let message: ChatMessage

    var body: some View {
        HStack(alignment: .top, spacing: 10) {
            if message.role == "user" { Spacer(minLength: 60) }

            VStack(alignment: message.role == "user" ? .trailing : .leading, spacing: 4) {
                if message.role == "assistant", let model = message.model {
                    HStack(spacing: 4) {
                        Image(systemName: "bolt.fill")
                            .font(.system(size: 9))
                        Text(model)
                            .font(.system(size: 10, weight: .medium, design: .monospaced))
                    }
                    .foregroundColor(.secondary)
                }

                Text(message.content + (message.isStreaming ? " \u{258C}" : ""))
                    .textSelection(.enabled)
                    .padding(.horizontal, 14)
                    .padding(.vertical, 10)
                    .background(
                        message.role == "user"
                        ? Color.accentColor.opacity(0.15)
                        : Color(nsColor: .controlBackgroundColor)
                    )
                    .cornerRadius(16)
                    .overlay(
                        RoundedRectangle(cornerRadius: 16)
                            .stroke(Color.primary.opacity(0.06), lineWidth: 1)
                    )

                Text(message.timestamp, style: .time)
                    .font(.system(size: 9))
                    .foregroundColor(.secondary.opacity(0.6))
            }

            if message.role == "assistant" { Spacer(minLength: 60) }
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 2)
    }
}

struct ChatView: View {
    @StateObject private var vm = ChatViewModel()
    @FocusState private var inputFocused: Bool
    @State private var showSettings = false

    var body: some View {
        VStack(spacing: 0) {
            // Header
            HStack {
                HStack(spacing: 6) {
                    Image(systemName: "sparkle")
                        .font(.system(size: 16, weight: .semibold))
                        .foregroundColor(.orange)
                    Text("Spark")
                        .font(.system(size: 18, weight: .bold, design: .rounded))
                }

                Spacer()

                if let err = vm.errorMessage {
                    Text(err)
                        .font(.system(size: 10))
                        .foregroundColor(.red)
                        .lineLimit(1)
                }

                Spacer()

                HStack(spacing: 12) {
                    // Model selector
                    Menu {
                        ForEach(vm.availableModels, id: \.self) { model in
                            Button(action: { vm.selectedModel = model }) {
                                HStack {
                                    Text(model)
                                    if model == vm.selectedModel {
                                        Image(systemName: "checkmark")
                                    }
                                }
                            }
                        }
                    } label: {
                        HStack(spacing: 4) {
                            Circle()
                                .fill(vm.selectedModel == SparkConfig.mainModel ? Color.purple : Color.green)
                                .frame(width: 7, height: 7)
                            Text(vm.modelLabel)
                                .font(.system(size: 11, weight: .medium, design: .monospaced))
                        }
                        .padding(.horizontal, 8)
                        .padding(.vertical, 4)
                        .background(Color.primary.opacity(0.05))
                        .cornerRadius(6)
                    }
                    .buttonStyle(.plain)

                    // Deep think toggle
                    Button(action: { vm.useDeepThink.toggle() }) {
                        HStack(spacing: 3) {
                            Image(systemName: vm.useDeepThink ? "brain.fill" : "brain")
                                .font(.system(size: 11))
                            Text("Deep")
                                .font(.system(size: 10, weight: .medium))
                        }
                        .padding(.horizontal, 6)
                        .padding(.vertical, 3)
                        .background(vm.useDeepThink ? Color.purple.opacity(0.15) : Color.clear)
                        .cornerRadius(4)
                    }
                    .buttonStyle(.plain)
                    .help("Deep Think: Use 31B model for complex reasoning")

                    Button(action: { vm.clearChat() }) {
                        Image(systemName: "trash")
                            .font(.system(size: 12))
                    }
                    .buttonStyle(.plain)
                    .help("Clear chat")
                }
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 10)
            .background(.ultraThinMaterial)

            Divider()

            // Messages
            ScrollViewReader { proxy in
                ScrollView {
                    LazyVStack(spacing: 4) {
                        if vm.messages.isEmpty {
                            VStack(spacing: 12) {
                                Image(systemName: "sparkle")
                                    .font(.system(size: 40))
                                    .foregroundColor(.orange.opacity(0.4))
                                Text("Spark")
                                    .font(.system(size: 28, weight: .bold, design: .rounded))
                                    .foregroundColor(.primary.opacity(0.3))
                                Text("DGX Spark \u{00B7} Gemma4")
                                    .font(.system(size: 13, design: .monospaced))
                                    .foregroundColor(.secondary.opacity(0.5))

                                HStack(spacing: 8) {
                                    QuickButton(text: "Explain quantum computing", icon: "atom") {
                                        vm.inputText = "Explain quantum computing simply"
                                    }
                                    QuickButton(text: "Write Python code", icon: "chevron.left.forwardslash.chevron.right") {
                                        vm.inputText = "Write a Python function to find prime numbers"
                                    }
                                    QuickButton(text: "Translate Korean", icon: "globe") {
                                        vm.inputText = "Translate this to English: "
                                        inputFocused = true
                                    }
                                }
                                .padding(.top, 8)
                            }
                            .frame(maxWidth: .infinity, maxHeight: .infinity)
                            .padding(.top, 80)
                        }

                        ForEach(vm.messages) { msg in
                            MessageBubble(message: msg)
                                .id(msg.id)
                        }

                        Color.clear.frame(height: 1).id("bottom")
                    }
                    .padding(.vertical, 8)
                }
                .onChange(of: vm.messages.count) { _ in
                    withAnimation(.easeOut(duration: 0.2)) {
                        proxy.scrollTo("bottom", anchor: .bottom)
                    }
                }
                .onChange(of: vm.messages.last?.content) { _ in
                    proxy.scrollTo("bottom", anchor: .bottom)
                }
            }

            Divider()

            // Input
            HStack(spacing: 10) {
                TextField("Message Spark...", text: $vm.inputText, axis: .vertical)
                    .textFieldStyle(.plain)
                    .lineLimit(1...6)
                    .focused($inputFocused)
                    .onSubmit {
                        if !NSEvent.modifierFlags.contains(.shift) {
                            vm.send()
                        }
                    }
                    .font(.system(size: 14))

                if vm.isLoading {
                    Button(action: { vm.stop() }) {
                        Image(systemName: "stop.circle.fill")
                            .font(.system(size: 22))
                            .foregroundColor(.red.opacity(0.8))
                    }
                    .buttonStyle(.plain)
                } else {
                    Button(action: { vm.send() }) {
                        Image(systemName: "arrow.up.circle.fill")
                            .font(.system(size: 22))
                            .foregroundColor(vm.inputText.trimmingCharacters(in: .whitespaces).isEmpty ? .secondary : .accentColor)
                    }
                    .buttonStyle(.plain)
                    .disabled(vm.inputText.trimmingCharacters(in: .whitespaces).isEmpty)
                }
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 12)
            .background(.ultraThinMaterial)
        }
        .frame(minWidth: 500, minHeight: 400)
        .onAppear {
            vm.loadModels()
            inputFocused = true
        }
    }
}

struct QuickButton: View {
    let text: String
    let icon: String
    let action: () -> Void

    var body: some View {
        Button(action: action) {
            HStack(spacing: 4) {
                Image(systemName: icon)
                    .font(.system(size: 10))
                Text(text)
                    .font(.system(size: 11))
            }
            .padding(.horizontal, 10)
            .padding(.vertical, 6)
            .background(Color.primary.opacity(0.05))
            .cornerRadius(8)
        }
        .buttonStyle(.plain)
    }
}

// MARK: - App Entry
@main
struct SparkApp: App {
    var body: some Scene {
        WindowGroup {
            ChatView()
        }
        .windowStyle(.titleBar)
        .defaultSize(width: 700, height: 600)
    }
}
