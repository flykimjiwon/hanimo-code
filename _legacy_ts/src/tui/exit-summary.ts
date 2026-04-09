/**
 * Print a session summary box when exiting hanimo.
 * Rendered directly to stdout after Ink unmounts.
 */

export interface SessionStats {
  sessionId: string;
  toolCalls: number;
  toolSuccesses: number;
  toolErrors: number;
  wallTimeMs: number;
  agentActiveMs: number;
  apiTimeMs: number;
  promptTokens: number;
  completionTokens: number;
  totalCost: number;
}

function formatDuration(ms: number): string {
  if (ms < 1000) return `${ms}ms`;
  const seconds = Math.floor(ms / 1000);
  if (seconds < 60) return `${seconds}s`;
  const minutes = Math.floor(seconds / 60);
  const secs = seconds % 60;
  if (minutes < 60) return `${minutes}m ${secs}s`;
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  return `${hours}h ${mins}m ${secs}s`;
}

function padLeft(str: string, len: number): string {
  return ' '.repeat(Math.max(0, len - str.length)) + str;
}

export function printExitSummary(stats: SessionStats): void {
  const BOX_WIDTH = 65;
  const INNER = BOX_WIDTH - 4; // 2 border + 2 padding

  const successRate = stats.toolCalls > 0
    ? ((stats.toolSuccesses / stats.toolCalls) * 100).toFixed(1)
    : '0.0';

  const apiPercent = stats.agentActiveMs > 0
    ? ((stats.apiTimeMs / stats.agentActiveMs) * 100).toFixed(1)
    : '0.0';

  const costStr = stats.totalCost > 0
    ? `$${stats.totalCost.toFixed(4)}`
    : '$0.00';

  const lines = [
    '',
    '  Agent powering down. Goodbye!',
    '',
    '  Interaction Summary',
    `  Session ID:${padLeft(stats.sessionId.slice(0, 36), INNER - 14)}`,
    `  Tool Calls:${padLeft(`${stats.toolCalls} ( \u2713 ${stats.toolSuccesses} \u2717 ${stats.toolErrors} )`, INNER - 14)}`,
    `  Success Rate:${padLeft(`${successRate}%`, INNER - 16)}`,
    `  Tokens:${padLeft(`${stats.promptTokens} in / ${stats.completionTokens} out`, INNER - 10)}`,
    `  Cost:${padLeft(costStr, INNER - 8)}`,
    '',
    '  Performance',
    `  Wall Time:${padLeft(formatDuration(stats.wallTimeMs), INNER - 13)}`,
    `  Agent Active:${padLeft(formatDuration(stats.agentActiveMs), INNER - 16)}`,
    `    \u00BB API Time:${padLeft(`${formatDuration(stats.apiTimeMs)} (${apiPercent}%)`, INNER - 16)}`,
    '',
  ];

  const top = '\u256D' + '\u2500'.repeat(BOX_WIDTH - 2) + '\u256E';
  const bottom = '\u2570' + '\u2500'.repeat(BOX_WIDTH - 2) + '\u256F';

  const output = [
    top,
    ...lines.map(line => {
      const visible = line;
      const pad = Math.max(0, BOX_WIDTH - 4 - visible.length);
      return `\u2502 ${visible}${' '.repeat(pad)} \u2502`;
    }),
    bottom,
  ].join('\n');

  process.stdout.write('\n' + output + '\n\n');
}
