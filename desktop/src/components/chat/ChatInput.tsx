import { useState, useCallback, KeyboardEvent } from "react";
import { useThemeStore } from "../../stores/theme-store";
import { useChatStore } from "../../stores/chat-store";
import { useOnboardingStore } from "../../stores/onboarding-store";

interface ChatInputProps {
  onSend: (content: string) => void;
  onStop?: () => void;
  disabled?: boolean;
}

export function ChatInput({ onSend, onStop, disabled }: ChatInputProps) {
  const { theme } = useThemeStore();
  const c = theme.colors;
  const [input, setInput] = useState("");
  const [history, setHistory] = useState<string[]>([]);
  const [historyIndex, setHistoryIndex] = useState(-1);
  const { isStreaming, connectionStatus } = useChatStore();

  const canSend = connectionStatus === "connected" && !isStreaming && input.trim().length > 0 && !disabled;

  const handleSend = useCallback(() => {
    if (!canSend) return;
    const trimmed = input.trim();
    onSend(trimmed);
    setHistory((prev) => [...prev, trimmed]);
    setHistoryIndex(-1);
    setInput("");
  }, [canSend, input, onSend]);

  const handleKeyDown = useCallback(
    (e: KeyboardEvent<HTMLTextAreaElement>) => {
      if (e.key === "ArrowUp" && !e.shiftKey) {
        e.preventDefault();
        if (history.length === 0) return;
        const newIndex =
          historyIndex === -1 ? history.length - 1 : Math.max(0, historyIndex - 1);
        setHistoryIndex(newIndex);
        setInput(history[newIndex]);
      } else if (e.key === "ArrowDown" && !e.shiftKey) {
        e.preventDefault();
        if (historyIndex === -1) return;
        const newIndex = historyIndex + 1;
        if (newIndex >= history.length) {
          setHistoryIndex(-1);
          setInput("");
        } else {
          setHistoryIndex(newIndex);
          setInput(history[newIndex]);
        }
      } else if (e.key === "Enter" && !e.shiftKey) {
        e.preventDefault();
        handleSend();
      }
    },
    [handleSend, history, historyIndex]
  );

  return (
    <div
      className="flex items-end gap-2 px-3 py-2 flex-shrink-0"
      style={{ borderTop: `1px solid ${c.border}` }}
    >
      <textarea
        value={input}
        onChange={(e) => setInput(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="Type a message..."
        disabled={disabled && !isStreaming}
        rows={1}
        className="flex-1 rounded-md px-3 py-1.5 text-sm outline-none resize-none"
        style={{
          background: c.inputBg,
          border: `1px solid ${c.inputBorder}`,
          color: c.text,
          maxHeight: 120,
          opacity: disabled && !isStreaming ? 0.5 : 1,
        }}
      />
      {isStreaming ? (
        <button
          onClick={onStop}
          className="px-3 py-1.5 rounded-md text-sm font-medium transition-colors flex-shrink-0"
          style={{
            background: c.error,
            color: "#fff",
            cursor: "pointer",
            border: `1px solid ${c.error}`,
          }}
        >
          Stop
        </button>
      ) : (
        <button
          onClick={handleSend}
          disabled={!canSend}
          className="px-3 py-1.5 rounded-md text-sm font-medium transition-colors flex-shrink-0"
          style={{
            background: canSend ? c.accent : c.bgTertiary,
            color: canSend ? "#ffffff" : c.textMuted,
            cursor: canSend ? "pointer" : "not-allowed",
            border: `1px solid ${canSend ? c.accent : c.border}`,
          }}
        >
          Send
        </button>
      )}
    </div>
  );
}
