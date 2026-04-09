/**
 * 10 TUI theme presets.
 * Each preset defines semantic color mappings used by all components.
 * Switch via /theme command or Esc menu.
 */

export interface ThemeColors {
  userText: string;
  assistantText: string;
  toolCall: string;
  toolResult: string;
  error: string;
  success: string;
  warning: string;
  cost: string;
  provider: string;
  model: string;
  border: string;
  borderFocus: string;
  dimText: string;
  statusIdle: string;
  statusThinking: string;
  statusTool: string;
  hint: string;
  prompt: string;
  promptDisabled: string;
}

export interface ThemePreset {
  id: string;
  name: string;
  description: string;
  colors: ThemeColors;
}

// ── 1. Catppuccin Mocha (default) ──
const catppuccin: ThemePreset = {
  id: 'catppuccin',
  name: 'Catppuccin Mocha',
  description: '부드러운 파스텔 다크 (기본값)',
  colors: {
    userText: '#F5E0DC',
    assistantText: '#CDD6F4',
    toolCall: '#CBA6F7',
    toolResult: '#6C7086',
    error: '#F38BA8',
    success: '#A6E3A1',
    warning: '#F9E2AF',
    cost: '#F9E2AF',
    provider: '#A6E3A1',
    model: '#89B4FA',
    border: '#45475A',
    borderFocus: '#B4BEFE',
    dimText: '#6C7086',
    statusIdle: '#A6E3A1',
    statusThinking: '#94E2D5',
    statusTool: '#CBA6F7',
    hint: '#89DCEB',
    prompt: '#A6E3A1',
    promptDisabled: '#6C7086',
  },
};

// ── 2. Tokyo Night ──
const tokyoNight: ThemePreset = {
  id: 'tokyo-night',
  name: 'Tokyo Night',
  description: '도쿄의 네온 감성',
  colors: {
    userText: '#C0CAF5',
    assistantText: '#A9B1D6',
    toolCall: '#BB9AF7',
    toolResult: '#565F89',
    error: '#F7768E',
    success: '#9ECE6A',
    warning: '#E0AF68',
    cost: '#E0AF68',
    provider: '#9ECE6A',
    model: '#7AA2F7',
    border: '#3B4261',
    borderFocus: '#7AA2F7',
    dimText: '#565F89',
    statusIdle: '#9ECE6A',
    statusThinking: '#7DCFFF',
    statusTool: '#BB9AF7',
    hint: '#7DCFFF',
    prompt: '#9ECE6A',
    promptDisabled: '#565F89',
  },
};

// ── 3. Dracula ──
const dracula: ThemePreset = {
  id: 'dracula',
  name: 'Dracula',
  description: '클래식 뱀파이어 다크',
  colors: {
    userText: '#F8F8F2',
    assistantText: '#F8F8F2',
    toolCall: '#BD93F9',
    toolResult: '#6272A4',
    error: '#FF5555',
    success: '#50FA7B',
    warning: '#F1FA8C',
    cost: '#F1FA8C',
    provider: '#50FA7B',
    model: '#8BE9FD',
    border: '#44475A',
    borderFocus: '#BD93F9',
    dimText: '#6272A4',
    statusIdle: '#50FA7B',
    statusThinking: '#8BE9FD',
    statusTool: '#BD93F9',
    hint: '#8BE9FD',
    prompt: '#50FA7B',
    promptDisabled: '#6272A4',
  },
};

// ── 4. Nord ──
const nord: ThemePreset = {
  id: 'nord',
  name: 'Nord',
  description: '차가운 북유럽 블루',
  colors: {
    userText: '#ECEFF4',
    assistantText: '#D8DEE9',
    toolCall: '#B48EAD',
    toolResult: '#4C566A',
    error: '#BF616A',
    success: '#A3BE8C',
    warning: '#EBCB8B',
    cost: '#EBCB8B',
    provider: '#A3BE8C',
    model: '#81A1C1',
    border: '#3B4252',
    borderFocus: '#88C0D0',
    dimText: '#4C566A',
    statusIdle: '#A3BE8C',
    statusThinking: '#88C0D0',
    statusTool: '#B48EAD',
    hint: '#88C0D0',
    prompt: '#A3BE8C',
    promptDisabled: '#4C566A',
  },
};

// ── 5. Gruvbox Dark ──
const gruvbox: ThemePreset = {
  id: 'gruvbox',
  name: 'Gruvbox',
  description: '레트로 따뜻한 톤',
  colors: {
    userText: '#EBDBB2',
    assistantText: '#D5C4A1',
    toolCall: '#D3869B',
    toolResult: '#665C54',
    error: '#FB4934',
    success: '#B8BB26',
    warning: '#FABD2F',
    cost: '#FABD2F',
    provider: '#B8BB26',
    model: '#83A598',
    border: '#504945',
    borderFocus: '#FE8019',
    dimText: '#665C54',
    statusIdle: '#B8BB26',
    statusThinking: '#8EC07C',
    statusTool: '#D3869B',
    hint: '#83A598',
    prompt: '#B8BB26',
    promptDisabled: '#665C54',
  },
};

// ── 6. Solarized Dark ──
const solarized: ThemePreset = {
  id: 'solarized',
  name: 'Solarized',
  description: '눈 편한 과학적 배색',
  colors: {
    userText: '#EEE8D5',
    assistantText: '#839496',
    toolCall: '#6C71C4',
    toolResult: '#586E75',
    error: '#DC322F',
    success: '#859900',
    warning: '#B58900',
    cost: '#B58900',
    provider: '#859900',
    model: '#268BD2',
    border: '#073642',
    borderFocus: '#2AA198',
    dimText: '#586E75',
    statusIdle: '#859900',
    statusThinking: '#2AA198',
    statusTool: '#6C71C4',
    hint: '#2AA198',
    prompt: '#859900',
    promptDisabled: '#586E75',
  },
};

// ── 7. One Dark ──
const oneDark: ThemePreset = {
  id: 'one-dark',
  name: 'One Dark',
  description: 'Atom 에디터 감성',
  colors: {
    userText: '#E5C07B',
    assistantText: '#ABB2BF',
    toolCall: '#C678DD',
    toolResult: '#5C6370',
    error: '#E06C75',
    success: '#98C379',
    warning: '#E5C07B',
    cost: '#E5C07B',
    provider: '#98C379',
    model: '#61AFEF',
    border: '#3E4451',
    borderFocus: '#61AFEF',
    dimText: '#5C6370',
    statusIdle: '#98C379',
    statusThinking: '#56B6C2',
    statusTool: '#C678DD',
    hint: '#56B6C2',
    prompt: '#98C379',
    promptDisabled: '#5C6370',
  },
};

// ── 8. Cyberpunk ──
const cyberpunk: ThemePreset = {
  id: 'cyberpunk',
  name: 'Cyberpunk',
  description: '네온 핑크 사이버펑크',
  colors: {
    userText: '#FF2079',
    assistantText: '#E0E0E0',
    toolCall: '#FF6AC1',
    toolResult: '#525252',
    error: '#FF0055',
    success: '#00FF9C',
    warning: '#FFD700',
    cost: '#FFD700',
    provider: '#00FF9C',
    model: '#00D4FF',
    border: '#333333',
    borderFocus: '#FF2079',
    dimText: '#525252',
    statusIdle: '#00FF9C',
    statusThinking: '#00D4FF',
    statusTool: '#FF6AC1',
    hint: '#00D4FF',
    prompt: '#00FF9C',
    promptDisabled: '#525252',
  },
};

// ── 9. Minimal Mono ──
const minimal: ThemePreset = {
  id: 'minimal',
  name: 'Minimal',
  description: '깔끔한 모노크롬',
  colors: {
    userText: '#FFFFFF',
    assistantText: '#C0C0C0',
    toolCall: '#A0A0A0',
    toolResult: '#606060',
    error: '#FF6B6B',
    success: '#C0C0C0',
    warning: '#D0D0D0',
    cost: '#D0D0D0',
    provider: '#C0C0C0',
    model: '#FFFFFF',
    border: '#404040',
    borderFocus: '#808080',
    dimText: '#606060',
    statusIdle: '#C0C0C0',
    statusThinking: '#A0A0A0',
    statusTool: '#A0A0A0',
    hint: '#808080',
    prompt: '#C0C0C0',
    promptDisabled: '#606060',
  },
};

// ── 10. Forest ──
const forest: ThemePreset = {
  id: 'forest',
  name: 'Forest',
  description: '자연 숲속 그린 톤',
  colors: {
    userText: '#D4E7C5',
    assistantText: '#B5C99A',
    toolCall: '#87A96B',
    toolResult: '#4A5D3A',
    error: '#E57373',
    success: '#81C784',
    warning: '#FFD54F',
    cost: '#FFD54F',
    provider: '#81C784',
    model: '#A5D6A7',
    border: '#2E3B2E',
    borderFocus: '#66BB6A',
    dimText: '#4A5D3A',
    statusIdle: '#81C784',
    statusThinking: '#80CBC4',
    statusTool: '#87A96B',
    hint: '#80CBC4',
    prompt: '#81C784',
    promptDisabled: '#4A5D3A',
  },
};

export const THEME_PRESETS: ThemePreset[] = [
  catppuccin,
  tokyoNight,
  dracula,
  nord,
  gruvbox,
  solarized,
  oneDark,
  cyberpunk,
  minimal,
  forest,
];

export function getThemeById(id: string): ThemePreset | undefined {
  return THEME_PRESETS.find((t) => t.id === id);
}

export const DEFAULT_THEME = catppuccin;
