// Dynamic theme system — colors are mutable, swapped via setTheme()
import { DEFAULT_THEME } from './themes.js';
import type { ThemeColors } from './themes.js';

// Mutable semantic color mapping — all components import this
export const colors: ThemeColors = { ...DEFAULT_THEME.colors };

/** Swap all colors to a new theme. Components pick up changes on next render. */
export function setTheme(newColors: ThemeColors): void {
  for (const key of Object.keys(newColors) as Array<keyof ThemeColors>) {
    colors[key] = newColors[key];
  }
}

// Re-export for backward compat (some files import `theme`)
export const theme = {
  get base() { return '#1E1E2E'; },
  get text() { return colors.assistantText; },
  get blue() { return colors.model; },
  get green() { return colors.success; },
  get red() { return colors.error; },
  get yellow() { return colors.warning; },
  get mauve() { return colors.toolCall; },
  get teal() { return colors.statusThinking; },
  get sky() { return colors.hint; },
};
