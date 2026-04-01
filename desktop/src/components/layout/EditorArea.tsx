import { useThemeStore } from "../../stores/theme-store";

export function EditorArea() {
  const { theme } = useThemeStore();
  const c = theme.colors;

  return (
    <div
      className="flex flex-col flex-1 overflow-hidden"
      style={{ background: c.editorBg }}
    >
      {/* Tab bar */}
      <div
        className="flex items-center flex-shrink-0 px-2 gap-1"
        style={{ height: 36, background: c.tabBg, borderBottom: `1px solid ${c.border}` }}
      >
        <div
          className="flex items-center px-3 text-sm rounded-t-sm"
          style={{
            height: 36,
            background: c.tabActiveBg,
            color: c.text,
            borderTop: `1px solid ${c.accent}`,
          }}
        >
          Welcome
        </div>
      </div>

      {/* Center content */}
      <div className="flex-1 flex flex-col items-center justify-center gap-3">
        <span className="text-5xl">🐶</span>
        <p className="text-sm" style={{ color: c.textMuted }}>
          Start a conversation to begin
        </p>
      </div>
    </div>
  );
}
