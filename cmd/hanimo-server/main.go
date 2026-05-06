// Command hanimo-server exposes the hanimo engine over a local HTTP+SSE API.
// It is intended to be spawned as a subprocess by the VS Code extension
// (or any other UI) and is bound to a single workspace directory.
//
// Ported from techai/cmd/tgc-server, adapted to hanimo's richer engine
// (multi-provider LLM, LSP integration, larger tool set).
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"github.com/flykimjiwon/hanimo/internal/config"
	"github.com/flykimjiwon/hanimo/internal/hooks"
	"github.com/flykimjiwon/hanimo/internal/llm"
	"github.com/flykimjiwon/hanimo/internal/tools"
)

var version = "dev"

type chatRequest struct {
	Prompt    string   `json:"prompt"`
	Stdin     string   `json:"stdin,omitempty"`
	Model     string   `json:"model,omitempty"`
	Mode      string   `json:"mode,omitempty"`
	MaxTurns  int      `json:"max_turns,omitempty"`
	Ephemeral bool     `json:"ephemeral,omitempty"`
	Images    []string `json:"images,omitempty"` // data URLs, e.g. "data:image/png;base64,..."
	Client    string   `json:"client,omitempty"` // "vscode" → extension-specific system prompt
}

type sseEvent struct {
	Event string
	Data  any
}

func main() {
	host := flag.String("host", "127.0.0.1", "bind host (default 127.0.0.1, never expose to LAN)")
	port := flag.Int("port", 0, "bind port (0 = OS-assigned)")
	cwd := flag.String("cwd", "", "workspace directory (default: process cwd)")
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("hanimo-server %s\n", version)
		return
	}

	if *cwd != "" {
		if err := os.Chdir(*cwd); err != nil {
			fmt.Fprintf(os.Stderr, "hanimo-server: cannot chdir to %q: %v\n", *cwd, err)
			os.Exit(1)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		cfg = config.DefaultConfig()
	}

	mux := http.NewServeMux()
	srv := &server{cfg: cfg}
	mux.HandleFunc("/health", srv.handleHealth)
	mux.HandleFunc("/version", srv.handleVersion)
	mux.HandleFunc("/chat", srv.handleChat)
	mux.HandleFunc("/shutdown", srv.handleShutdown)
	mux.HandleFunc("/config", srv.handleConfig)
	mux.HandleFunc("/models", srv.handleModels)
	mux.HandleFunc("/models/refresh", srv.handleModelsRefresh)
	mux.HandleFunc("/index/symbols", srv.handleIndexSymbols)
	mux.HandleFunc("/knowledge", srv.handleKnowledge)
	mux.HandleFunc("/knowledge/files", srv.handleKnowledgeFiles)
	mux.HandleFunc("/skills", srv.handleSkills)
	mux.HandleFunc("/rules", srv.handleRules)
	mux.HandleFunc("/permissions", srv.handlePermissions)
	mux.HandleFunc("/confirm", srv.handleConfirm)

	addr := fmt.Sprintf("%s:%d", *host, *port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "hanimo-server: listen %s: %v\n", addr, err)
		os.Exit(1)
	}
	actualAddr := listener.Addr().(*net.TCPAddr)

	srv.httpServer = &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Print addr line that the parent (VS Code extension) parses.
	// Format MUST stay stable: "hanimo-server listening on http://HOST:PORT"
	fmt.Printf("hanimo-server listening on http://%s:%d\n", *host, actualAddr.Port)
	os.Stdout.Sync()

	if err := srv.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "hanimo-server: serve: %v\n", err)
		os.Exit(1)
	}
}

type server struct {
	cfg        config.Config
	httpServer *http.Server
	// chatMu serializes /chat requests because tools mutate process cwd
	// and global state. One workspace = one in-flight chat.
	chatMu sync.Mutex
}

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cwd, _ := os.Getwd()
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"version": version,
		"cwd":     cwd,
		"model":   s.cfg.Models.Super,
	})
}

func (s *server) handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"version": version})
}

func (s *server) handleShutdown(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	go func() {
		time.Sleep(100 * time.Millisecond)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = s.httpServer.Shutdown(ctx)
	}()
}

func (s *server) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Prompt) == "" {
		http.Error(w, "prompt is required", http.StatusBadRequest)
		return
	}
	if s.cfg.API.APIKey == "" && !strings.Contains(s.cfg.API.BaseURL, "localhost") && !strings.Contains(s.cfg.API.BaseURL, "127.0.0.1") {
		http.Error(w, "API key not configured. Run 'hanimo --setup' first.", http.StatusPreconditionFailed)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	emit := func(ev sseEvent) {
		payload, err := json.Marshal(ev.Data)
		if err != nil {
			return
		}
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ev.Event, payload)
		flusher.Flush()
	}

	s.chatMu.Lock()
	defer s.chatMu.Unlock()

	if err := s.runChat(r.Context(), req, emit); err != nil {
		emit(sseEvent{Event: "error", Data: map[string]string{"message": err.Error()}})
	}
	emit(sseEvent{Event: "done", Data: map[string]any{}})
}

// runChat is a server-side mirror of internal/exec.Run that emits SSE
// events instead of writing to stdout. Keep behavior aligned with exec.Run.
func (s *server) runChat(ctx context.Context, req chatRequest, emit func(sseEvent)) error {
	maxTurns := req.MaxTurns
	if maxTurns <= 0 {
		maxTurns = 20
	}

	model := s.cfg.Models.Super
	if req.Model != "" {
		model = req.Model
	}

	mode := parseMode(req.Mode)

	client := llm.NewClient(s.cfg.API.BaseURL, s.cfg.API.APIKey)
	systemPrompt := llm.SystemPrompt(mode)
	cwd, _ := os.Getwd()
	// Loads ~/.hanimo/rules.md + .hanimo.md + .hanimo-knowledge/*.md + skills index.
	// Cap at ~100KB to stay clear of even a 32K-token context budget.
	systemPrompt += loadKnowledgeContext(cwd, 100_000)

	// VS Code extension renders the assistant text through a Markdown engine
	// in the WebView. Tell the model to format final answers with proper
	// headings/bullets/code blocks so the rendered output looks polished.
	if req.Client == "vscode" {
		systemPrompt += "\n\n## VS Code Display Hint\n" +
			"Your response is rendered through a Markdown engine in a VS Code sidebar WebView.\n" +
			"After running tools, synthesize a clean final answer with:\n" +
			"- Short headings (## / ###) for sections when the answer has multiple parts\n" +
			"- Bullet points for lists, **bold** for key terms, `inline code` for paths/identifiers\n" +
			"- Triple-backtick code blocks with language tags for code snippets\n" +
			"- Tables for structured comparisons\n" +
			"- Concise prose — avoid restating tool output verbatim; summarize what matters\n" +
			"- For file proposals, prefix the code block's first line with a comment like `// path/to/file.ext` so the user can apply with one click.\n" +
			"Skip formatting fluff for trivially short answers (one or two sentences).\n"
	}

	userMsg := req.Prompt
	if req.Stdin != "" {
		userMsg += "\n\n<stdin>\n" + req.Stdin + "\n</stdin>"
	}

	// Multimodal: build a multipart message when images are attached. Falls
	// back to the legacy single-string Content when no images. The provider
	// must be vision-capable for this to render properly.
	userContent := userMsg
	var multipartUser []openai.ChatMessagePart
	if len(req.Images) > 0 {
		multipartUser = append(multipartUser, openai.ChatMessagePart{
			Type: openai.ChatMessagePartTypeText,
			Text: userMsg,
		})
		for _, dataURL := range req.Images {
			multipartUser = append(multipartUser, openai.ChatMessagePart{
				Type:     openai.ChatMessagePartTypeImageURL,
				ImageURL: &openai.ChatMessageImageURL{URL: dataURL},
			})
		}
		userContent = "" // ignored when MultiContent is set
	}

	hookMgr := hooks.NewManager()
	hookMgr.RunSessionStart("vscode-server")
	defer hookMgr.RunStop("vscode-server")

	userMessage := openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser}
	if multipartUser != nil {
		userMessage.MultiContent = multipartUser
	} else {
		userMessage.Content = userContent
	}
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
		userMessage,
	}
	toolDefs := tools.AllTools()

	for turn := 0; turn < maxTurns; turn++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		turnCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		ch := client.StreamChat(turnCtx, model, messages, toolDefs)

		var contentBuf strings.Builder
		var toolCalls []llm.ToolCallInfo

		for chunk := range ch {
			if chunk.Err != nil {
				cancel()
				return fmt.Errorf("LLM error: %w", chunk.Err)
			}
			if chunk.Content != "" {
				contentBuf.WriteString(chunk.Content)
				emit(sseEvent{Event: "content", Data: map[string]string{"text": chunk.Content}})
			}
			if chunk.Done {
				toolCalls = chunk.ToolCalls
				if chunk.Usage != nil {
					emit(sseEvent{Event: "usage", Data: map[string]int{
						"prompt_tokens":     chunk.Usage.PromptTokens,
						"completion_tokens": chunk.Usage.CompletionTokens,
						"total_tokens":      chunk.Usage.TotalTokens,
					}})
				}
			}
		}
		cancel()

		assistantMsg := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: contentBuf.String(),
		}
		if len(toolCalls) > 0 {
			assistantMsg.ToolCalls = toOpenAIToolCalls(toolCalls)
		}
		messages = append(messages, assistantMsg)

		if len(toolCalls) == 0 {
			emit(sseEvent{Event: "turn_end", Data: map[string]any{"turn": turn}})
			return nil
		}

		for _, tc := range toolCalls {
			emit(sseEvent{Event: "tool_call", Data: map[string]string{
				"id":        tc.ID,
				"name":      tc.Name,
				"arguments": tc.Arguments,
			}})

			// Permission gate — supports allow/deny/ask. "ask" pauses the loop,
			// emits a tool_confirm event over SSE, and waits up to 5min for
			// the user's decision before proceeding.
			decision, reason := CheckTool(tc.Name, tc.Arguments)
			if decision == "deny" {
				result := "[permission denied] " + reason
				emit(sseEvent{Event: "tool_result", Data: map[string]string{
					"id": tc.ID, "result": result, "aborted": "true",
				}})
				messages = append(messages, openai.ChatCompletionMessage{
					Role: openai.ChatMessageRoleTool, Content: result, ToolCallID: tc.ID,
				})
				continue
			}
			if decision == "ask" {
				_ = registerConfirm(tc.ID)
				emit(sseEvent{Event: "tool_confirm", Data: map[string]string{
					"id": tc.ID, "name": tc.Name, "arguments": tc.Arguments,
				}})
				dec := waitConfirm(tc.ID, 5*time.Minute)
				if !dec.Approve {
					reason := dec.Reason
					if reason == "" {
						reason = "user denied"
					}
					result := "[user denied] " + reason
					emit(sseEvent{Event: "tool_result", Data: map[string]string{
						"id": tc.ID, "result": result, "aborted": "true",
					}})
					messages = append(messages, openai.ChatCompletionMessage{
						Role: openai.ChatMessageRoleTool, Content: result, ToolCallID: tc.ID,
					})
					continue
				}
			}

			if hookMgr.RunPreToolUse(tc.Name, tc.Arguments) == hooks.HookFailAbort {
				result := "Aborted by pre_tool_use hook"
				emit(sseEvent{Event: "tool_result", Data: map[string]string{
					"id": tc.ID, "result": result, "aborted": "true",
				}})
				messages = append(messages, openai.ChatCompletionMessage{
					Role: openai.ChatMessageRoleTool, Content: result, ToolCallID: tc.ID,
				})
				continue
			}

			result := tools.Execute(tc.Name, tc.Arguments)
			hookMgr.RunPostToolUse(tc.Name, tc.Arguments, result)

			emit(sseEvent{Event: "tool_result", Data: map[string]string{
				"id": tc.ID, "result": result,
			}})
			messages = append(messages, openai.ChatCompletionMessage{
				Role: openai.ChatMessageRoleTool, Content: result, ToolCallID: tc.ID,
			})
		}
	}

	return fmt.Errorf("max turns (%d) reached without completion", maxTurns)
}

func parseMode(s string) llm.Mode {
	switch strings.ToLower(s) {
	case "dev", "deep":
		return llm.ModeDev
	case "plan":
		return llm.ModePlan
	default:
		return llm.ModeSuper
	}
}

func toOpenAIToolCalls(calls []llm.ToolCallInfo) []openai.ToolCall {
	out := make([]openai.ToolCall, len(calls))
	for i, tc := range calls {
		out[i] = openai.ToolCall{
			ID:       tc.ID,
			Type:     openai.ToolTypeFunction,
			Function: openai.FunctionCall{Name: tc.Name, Arguments: tc.Arguments},
		}
	}
	return out
}
