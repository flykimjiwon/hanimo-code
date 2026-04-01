import { useState } from "react";
import { useThemeStore } from "../../stores/theme-store";

export function ChatPanel() {
  const { theme } = useThemeStore();
  const c = theme.colors;
  const [input, setInput] = useState("");

  return (
    <div
      className="flex flex-col flex-shrink-0"
      style={{ width: 340, background: c.bg, borderLeft: `1px solid ${c.border}` }}
    >
      {/* Header */}
      <div
        className="flex items-center gap-2 px-3 flex-shrink-0"
        style={{ height: 36, borderBottom: `1px solid ${c.border}`, background: c.bgSecondary }}
      >
        <span className="text-base">🐶</span>
        <span className="text-sm font-medium" style={{ color: c.text }}>
          hanimo chat
        </span>
      </div>

      {/* Messages area */}
      <div className="flex-1 overflow-y-auto px-3 py-4">
        <p className="text-sm text-center" style={{ color: c.textMuted }}>
          No messages yet. Start a conversation!
        </p>
      </div>

      {/* Input area */}
      <div
        className="flex items-center gap-2 px-3 py-2 flex-shrink-0"
        style={{ borderTop: `1px solid ${c.border}` }}
      >
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Type a message..."
          className="flex-1 rounded-md px-3 py-1.5 text-sm outline-none"
          style={{
            background: c.inputBg,
            border: `1px solid ${c.inputBorder}`,
            color: c.text,
          }}
        />
        <button
          className="px-3 py-1.5 rounded-md text-sm font-medium transition-colors"
          style={{ background: c.accent, color: "#ffffff" }}
          onClick={() => setInput("")}
        >
          Send
        </button>
      </div>
    </div>
  );
}
