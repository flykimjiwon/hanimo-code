export interface Theme {
  name: "dark" | "light";
  colors: {
    bg: string;
    bgSecondary: string;
    bgTertiary: string;
    border: string;
    text: string;
    textSecondary: string;
    textMuted: string;
    accent: string;
    accentHover: string;
    userBubble: string;
    assistantBubble: string;
    success: string;
    error: string;
    sidebarBg: string;
    sidebarIcon: string;
    sidebarIconActive: string;
    editorBg: string;
    tabBg: string;
    tabActiveBg: string;
    inputBg: string;
    inputBorder: string;
  };
}

export const darkTheme: Theme = {
  name: "dark",
  colors: {
    bg: "#0d0d0d",
    bgSecondary: "#141414",
    bgTertiary: "#1a1a1a",
    border: "#2a2a2a",
    text: "#e8e8e8",
    textSecondary: "#a0a0a0",
    textMuted: "#606060",
    accent: "#6366f1",
    accentHover: "#4f52d9",
    userBubble: "#1e1e2e",
    assistantBubble: "#141414",
    success: "#7ee787",
    error: "#ff7b72",
    sidebarBg: "#0a0a0a",
    sidebarIcon: "#606060",
    sidebarIconActive: "#6366f1",
    editorBg: "#111111",
    tabBg: "#141414",
    tabActiveBg: "#1a1a1a",
    inputBg: "#1a1a1a",
    inputBorder: "#2a2a2a",
  },
};

export const lightTheme: Theme = {
  name: "light",
  colors: {
    bg: "#ffffff",
    bgSecondary: "#f8f9fa",
    bgTertiary: "#f1f3f5",
    border: "#e0e0e0",
    text: "#1a1a1a",
    textSecondary: "#555555",
    textMuted: "#999999",
    accent: "#1a73e8",
    accentHover: "#1558b0",
    userBubble: "#e8f0fe",
    assistantBubble: "#f8f9fa",
    success: "#28a745",
    error: "#dc3545",
    sidebarBg: "#f1f3f5",
    sidebarIcon: "#999999",
    sidebarIconActive: "#1a73e8",
    editorBg: "#ffffff",
    tabBg: "#f8f9fa",
    tabActiveBg: "#ffffff",
    inputBg: "#ffffff",
    inputBorder: "#e0e0e0",
  },
};

export function getTheme(name: "dark" | "light"): Theme {
  return name === "dark" ? darkTheme : lightTheme;
}
