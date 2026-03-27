import { describe, it, expect, beforeAll, afterAll } from 'vitest';
import { globSearchTool } from '../src/tools/glob-search.js';
import { writeFileSync, mkdirSync, rmSync } from 'node:fs';
import { join } from 'node:path';

const toolCtx = { toolCallId: 'test', messages: [], abortSignal: undefined as unknown as AbortSignal };
const TEST_DIR = join(process.cwd(), '.test-glob-tmp');

const defaults = { ignore: ['**/node_modules/**', '**/dist/**', '**/.git/**'], limit: 200 };

describe('glob-search tool', () => {
  beforeAll(() => {
    mkdirSync(join(TEST_DIR, 'src'), { recursive: true });
    mkdirSync(join(TEST_DIR, 'lib'), { recursive: true });
    writeFileSync(join(TEST_DIR, 'src', 'app.ts'), 'export {}');
    writeFileSync(join(TEST_DIR, 'src', 'utils.ts'), 'export {}');
    writeFileSync(join(TEST_DIR, 'src', 'style.css'), 'body {}');
    writeFileSync(join(TEST_DIR, 'lib', 'helper.js'), 'module.exports = {}');
    writeFileSync(join(TEST_DIR, 'readme.md'), '# test');
  });

  afterAll(() => {
    rmSync(TEST_DIR, { recursive: true, force: true });
  });

  it('should find files matching pattern', async () => {
    const result = await globSearchTool.execute(
      { ...defaults, pattern: '**/*.ts', cwd: TEST_DIR },
      toolCtx,
    ) as { success: boolean; totalFound: number; files: string[] };

    expect(result.success).toBe(true);
    expect(result.totalFound).toBe(2);
    expect(result.files.some(f => f.endsWith('app.ts'))).toBe(true);
    expect(result.files.some(f => f.endsWith('utils.ts'))).toBe(true);
  });

  it('should respect limit', async () => {
    const result = await globSearchTool.execute(
      { ...defaults, pattern: '**/*', cwd: TEST_DIR, limit: 2 },
      toolCtx,
    ) as { success: boolean; returned: number; truncated: boolean };

    expect(result.success).toBe(true);
    expect(result.returned).toBe(2);
    expect(result.truncated).toBe(true);
  });

  it('should find specific extensions', async () => {
    const result = await globSearchTool.execute(
      { ...defaults, pattern: '**/*.css', cwd: TEST_DIR },
      toolCtx,
    ) as { success: boolean; totalFound: number };

    expect(result.success).toBe(true);
    expect(result.totalFound).toBe(1);
  });

  it('should return empty for no matches', async () => {
    const result = await globSearchTool.execute(
      { ...defaults, pattern: '**/*.xyz', cwd: TEST_DIR },
      toolCtx,
    ) as { success: boolean; totalFound: number };

    expect(result.success).toBe(true);
    expect(result.totalFound).toBe(0);
  });
});
