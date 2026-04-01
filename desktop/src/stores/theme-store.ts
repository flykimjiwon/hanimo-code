import { create } from "zustand";
import { persist } from "zustand/middleware";
import { type Theme, getTheme } from "../lib/theme";

interface ThemeState {
  mode: "dark" | "light";
  theme: Theme;
  toggle: () => void;
  setMode: (mode: "dark" | "light") => void;
}

export const useThemeStore = create<ThemeState>()(
  persist(
    (set) => ({
      mode: "dark",
      theme: getTheme("dark"),
      toggle: () =>
        set((state) => {
          const newMode = state.mode === "dark" ? "light" : "dark";
          return { mode: newMode, theme: getTheme(newMode) };
        }),
      setMode: (mode) => set({ mode, theme: getTheme(mode) }),
    }),
    {
      name: "hanimo-theme",
      partialize: (state) => ({ mode: state.mode }),
      onRehydrateStorage: () => (state) => {
        if (state) {
          state.theme = getTheme(state.mode);
        }
      },
    }
  )
);
