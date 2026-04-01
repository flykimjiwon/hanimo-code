import { useThemeStore } from "../../stores/theme-store";

export function FileTree() {
  const { theme } = useThemeStore();
  const c = theme.colors;

  return (
    <div
      className="flex flex-col flex-shrink-0 overflow-hidden"
      style={{ width: 220, background: c.bgSecondary, borderRight: `1px solid ${c.border}` }}
    >
      <div
        className="px-3 py-2 text-xs font-semibold tracking-widest uppercase flex-shrink-0"
        style={{ color: c.textSecondary, borderBottom: `1px solid ${c.border}` }}
      >
        Explorer
      </div>
      <div className="flex-1 flex items-center justify-center px-4">
        <p className="text-center text-sm" style={{ color: c.textMuted }}>
          Open a folder to get started
        </p>
      </div>
    </div>
  );
}
