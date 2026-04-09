import { describe, it, expect, beforeAll, afterAll } from 'vitest';
import { batchTool } from '../src/tools/batch.js';
import { writeFileSync, mkdirSync, rmSync } from 'node:fs';
import { join } from 'node:path';

const toolCtx = { toolCallId: 'test', messages: [], abortSignal: undefined as unknown as AbortSignal };
const TEST_DIR = join(process.cwd(), '.test-batch-tmp');

describe('batch tool', () => {
  beforeAll(() => {
    mkdirSync(TEST_DIR, { recursive: true });
    writeFileSync(join(TEST_DIR, 'a.txt'), 'hello from a');
    writeFileSync(join(TEST_DIR, 'b.txt'), 'hello from b');
    writeFileSync(join(TEST_DIR, 'c.json'), '{"key": "value"}');
  });

  afterAll(() => {
    rmSync(TEST_DIR, { recursive: true, force: true });
  });

  it('should read multiple files in parallel', async () => {
    const result = await batchTool.execute(
      {
        reads: [
          { path: join(TEST_DIR, 'a.txt') },
          { path: join(TEST_DIR, 'b.txt') },
        ],
      },
      toolCtx,
    ) as { success: boolean; results: { reads: Array<{ path: string; content?: string; error?: string }> } };

    expect(result.success).toBe(true);
    expect(result.results.reads).toHaveLength(2);
    const readA = result.results.reads.find(r => r.path.endsWith('a.txt'))!;
    const readB = result.results.reads.find(r => r.path.endsWith('b.txt'))!;
    expect(readA.content).toContain('hello from a');
    expect(readB.content).toContain('hello from b');
  });

  it('should handle missing files gracefully', async () => {
    const result = await batchTool.execute(
      {
        reads: [
          { path: join(TEST_DIR, 'nonexistent.txt') },
          { path: join(TEST_DIR, 'a.txt') },
        ],
      },
      toolCtx,
    ) as { success: boolean; results: { reads: Array<{ path: string; content?: string; error?: string }> } };

    expect(result.success).toBe(true);
    expect(result.results.reads).toHaveLength(2);
    const missing = result.results.reads.find(r => r.path.endsWith('nonexistent.txt'))!;
    const found = result.results.reads.find(r => r.path.endsWith('a.txt'))!;
    expect(missing.error).toBeDefined();
    expect(found.content).toContain('hello from a');
  });

  it('should run glob patterns in parallel', async () => {
    const result = await batchTool.execute(
      {
        globs: [
          { pattern: '*.txt', cwd: TEST_DIR },
          { pattern: '*.json', cwd: TEST_DIR },
        ],
      },
      toolCtx,
    ) as { success: boolean; results: { globs: Array<{ pattern: string; files: string[] }> } };

    expect(result.success).toBe(true);
    expect(result.results.globs).toHaveLength(2);
    const txtGlob = result.results.globs.find(g => g.pattern === '*.txt')!;
    const jsonGlob = result.results.globs.find(g => g.pattern === '*.json')!;
    expect(txtGlob.files).toContain('a.txt');
    expect(txtGlob.files).toContain('b.txt');
    expect(jsonGlob.files).toContain('c.json');
  });

  it('should handle empty input gracefully', async () => {
    const result = await batchTool.execute({}, toolCtx) as { success: boolean };
    expect(result.success).toBe(true);
  });

  it('should combine reads and globs', async () => {
    const result = await batchTool.execute(
      {
        reads: [{ path: join(TEST_DIR, 'a.txt') }],
        globs: [{ pattern: '*.json', cwd: TEST_DIR }],
      },
      toolCtx,
    ) as { success: boolean; results: { reads: unknown[]; globs: unknown[] } };

    expect(result.success).toBe(true);
    expect(result.results.reads).toHaveLength(1);
    expect(result.results.globs).toHaveLength(1);
  });
});
