import { useState } from "react";
import { useThemeStore } from "../../stores/theme-store";

type ActivePanel = "explorer" | "search" | "chat" | "settings";

export function Sidebar() {
  const { theme, toggle, mode } = useThemeStore();
  const c = theme.colors;
  const [active, setActive] = useState<ActivePanel>("explorer");

  const buttons: { id: ActivePanel; icon: string; label: string }[] = [
    { id: "explorer", icon: "📁", label: "Explorer" },
    { id: "search", icon: "🔍", label: "Search" },
    { id: "chat", icon: "💬", label: "Chat" },
    { id: "settings", icon: "⚙️", label: "Settings" },
  ];

  return (
    <div
      className="flex flex-col items-center py-2 flex-shrink-0"
      style={{ width: 48, background: c.sidebarBg, borderRight: `1px solid ${c.border}` }}
    >
      <div className="flex flex-col gap-1 flex-1">
        {buttons.map((btn) => (
          <button
            key={btn.id}
            onClick={() => setActive(btn.id)}
            title={btn.label}
            className="flex items-center justify-center rounded-md text-lg transition-colors"
            style={{
              width: 40,
              height: 40,
              color: active === btn.id ? c.sidebarIconActive : c.sidebarIcon,
              background: active === btn.id ? c.bgSecondary : "transparent",
            }}
          >
            {btn.icon}
          </button>
        ))}
      </div>
      <button
        onClick={toggle}
        title="Toggle theme"
        className="flex items-center justify-center rounded-md text-lg transition-colors"
        style={{
          width: 40,
          height: 40,
          color: c.sidebarIcon,
          background: "transparent",
        }}
      >
        {mode === "dark" ? "☀️" : "🌙"}
      </button>
    </div>
  );
}
