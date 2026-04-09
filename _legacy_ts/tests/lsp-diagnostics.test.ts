import { describe, it, expect } from 'vitest';

// Test the parsing functions by importing the module and checking tool structure
import { lspDiagnosticsTool } from '../src/tools/lsp-diagnostics.js';

describe('lsp-diagnostics tool', () => {
  it('should be a valid tool with correct parameters', () => {
    expect(lspDiagnosticsTool).toBeDefined();
    expect(lspDiagnosticsTool.description).toContain('TypeScript');
    expect(lspDiagnosticsTool.description).toContain('ESLint');
  });

  it('should run diagnostics on this project (has tsconfig.json)', async () => {
    const result = await lspDiagnosticsTool.execute(
      { checker: 'tsc' },
      { toolCallId: 'test', messages: [], abortSignal: undefined as unknown as AbortSignal },
    ) as { success: boolean; checkers: string[]; summary: { total: number; errors: number; warnings: number } };

    expect(result.success).toBe(true);
    expect(result.checkers).toContain('tsc');
    expect(result.summary).toBeDefined();
    expect(typeof result.summary.total).toBe('number');
    expect(typeof result.summary.errors).toBe('number');
  }, 45000);

  it('should handle file-specific diagnostics', async () => {
    const result = await lspDiagnosticsTool.execute(
      { file: 'src/tools/todo.ts', checker: 'tsc' },
      { toolCallId: 'test2', messages: [], abortSignal: undefined as unknown as AbortSignal },
    ) as { success: boolean };

    expect(result.success).toBe(true);
  }, 45000);
});
