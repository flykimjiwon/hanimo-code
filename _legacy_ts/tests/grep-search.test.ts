import { describe, it, expect, beforeAll, afterAll } from 'vitest';
import { grepSearchTool } from '../src/tools/grep-search.js';
import { writeFileSync, mkdirSync, rmSync } from 'node:fs';
import { join } from 'node:path';

const toolCtx = { toolCallId: 'test', messages: [], abortSignal: undefined as unknown as AbortSignal };
const TEST_DIR = join(process.cwd(), '.test-grep-tmp');

const defaults = { glob: '**/*', ignoreCase: false, maxResults: 100, contextLines: 0 };

describe('grep-search tool', () => {
  beforeAll(() => {
    mkdirSync(TEST_DIR, { recursive: true });
    writeFileSync(
      join(TEST_DIR, 'sample.ts'),
      [
        'import { foo } from "bar";',
        '',
        'function hello() {',
        '  console.log("hello world");',
        '  return 42;',
        '}',
        '',
        'function goodbye() {',
        '  console.log("goodbye");',
        '}',
      ].join('\n'),
    );
    writeFileSync(join(TEST_DIR, 'other.txt'), 'no matches here\njust plain text\n');
  });

  afterAll(() => {
    rmSync(TEST_DIR, { recursive: true, force: true });
  });

  it('should find regex matches in files', async () => {
    const result = await grepSearchTool.execute(
      { ...defaults, pattern: 'console\\.log', path: TEST_DIR, glob: '*.ts' },
      toolCtx,
    ) as { success: boolean; totalMatches: number; matches: Array<{ line: number }> };

    expect(result.success).toBe(true);
    expect(result.totalMatches).toBe(2);
    expect(result.matches[0]!.line).toBe(4);
    expect(result.matches[1]!.line).toBe(9);
  });

  it('should support case-insensitive search', async () => {
    const result = await grepSearchTool.execute(
      { ...defaults, pattern: 'HELLO', path: TEST_DIR, glob: '*.ts', ignoreCase: true },
      toolCtx,
    ) as { success: boolean; totalMatches: number };

    expect(result.success).toBe(true);
    expect(result.totalMatches).toBeGreaterThanOrEqual(1);
  });

  it('should include context lines when requested', async () => {
    const result = await grepSearchTool.execute(
      { ...defaults, pattern: 'return 42', path: TEST_DIR, glob: '*.ts', contextLines: 1 },
      toolCtx,
    ) as { success: boolean; matches: Array<{ content: string; line: number }> };

    expect(result.success).toBe(true);
    expect(result.matches).toHaveLength(1);
    expect(result.matches[0]!.content).toContain('return 42');
  });

  it('should respect maxResults', async () => {
    const result = await grepSearchTool.execute(
      { ...defaults, pattern: 'function', path: TEST_DIR, glob: '*.ts', maxResults: 1 },
      toolCtx,
    ) as { success: boolean; totalMatches: number; truncated: boolean };

    expect(result.success).toBe(true);
    expect(result.totalMatches).toBe(1);
    expect(result.truncated).toBe(true);
  });

  it('should return error for invalid regex', async () => {
    const result = await grepSearchTool.execute(
      { ...defaults, pattern: '[invalid', path: TEST_DIR },
      toolCtx,
    ) as { success: boolean; error: string };

    expect(result.success).toBe(false);
    expect(result.error).toContain('Invalid regex');
  });

  it('should handle no matches gracefully', async () => {
    const result = await grepSearchTool.execute(
      { ...defaults, pattern: 'zzz_nonexistent_pattern_zzz', path: TEST_DIR, glob: '*.ts' },
      toolCtx,
    ) as { success: boolean; totalMatches: number };

    expect(result.success).toBe(true);
    expect(result.totalMatches).toBe(0);
  });
});
