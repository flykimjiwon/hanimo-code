import { describe, it, expect } from 'vitest';
import { parseUnifiedDiff, applyHunks } from '../src/tools/patch.js';

describe('patch', () => {
  const sampleDiff = `--- a/src/index.ts
+++ b/src/index.ts
@@ -1,3 +1,4 @@
 import { foo } from './foo';
+import { bar } from './bar';

 const x = 1;
@@ -10,3 +11,3 @@
 function hello() {
-  return 'hello';
+  return 'hello world';
 }`;

  it('should parse unified diff into file patches', () => {
    const patches = parseUnifiedDiff(sampleDiff);
    expect(patches.length).toBe(1);
    expect(patches[0]?.oldFile).toBe('src/index.ts');
    expect(patches[0]?.newFile).toBe('src/index.ts');
    expect(patches[0]?.hunks.length).toBe(2);
  });

  it('should parse hunk headers correctly', () => {
    const patches = parseUnifiedDiff(sampleDiff);
    const hunk1 = patches[0]?.hunks[0];
    expect(hunk1?.oldStart).toBe(1);
    expect(hunk1?.oldCount).toBe(3);
    expect(hunk1?.newStart).toBe(1);
    expect(hunk1?.newCount).toBe(4);
  });

  it('should apply a simple addition hunk', () => {
    const content = `import { foo } from './foo';

const x = 1;`;

    const hunks = [{
      oldStart: 1,
      oldCount: 3,
      newStart: 1,
      newCount: 4,
      lines: [
        " import { foo } from './foo';",
        "+import { bar } from './bar';",
        ' ',
        ' const x = 1;',
      ],
    }];

    const { result, applied, rejected } = applyHunks(content, hunks);
    expect(applied).toBe(1);
    expect(rejected).toEqual([]);
    expect(result).toContain("import { bar } from './bar'");
  });

  it('should apply a replacement hunk', () => {
    const content = `line 1
line 2
old text
line 4`;

    const hunks = [{
      oldStart: 2,
      oldCount: 3,
      newStart: 2,
      newCount: 3,
      lines: [
        ' line 2',
        '-old text',
        '+new text',
        ' line 4',
      ],
    }];

    const { result, applied } = applyHunks(content, hunks);
    expect(applied).toBe(1);
    expect(result).toContain('new text');
    expect(result).not.toContain('old text');
  });

  it('should report rejected hunks when context does not match', () => {
    const content = 'completely different content';
    const hunks = [{
      oldStart: 1,
      oldCount: 1,
      newStart: 1,
      newCount: 1,
      lines: ['-expected line', '+new line'],
    }];

    const { applied, rejected } = applyHunks(content, hunks);
    expect(applied).toBe(0);
    expect(rejected.length).toBe(1);
  });

  it('should return empty patches for invalid diff', () => {
    const patches = parseUnifiedDiff('not a diff');
    expect(patches).toEqual([]);
  });

  it('should handle multi-file diffs', () => {
    const multiDiff = `--- a/file1.ts
+++ b/file1.ts
@@ -1,1 +1,1 @@
-old1
+new1
--- a/file2.ts
+++ b/file2.ts
@@ -1,1 +1,1 @@
-old2
+new2`;

    const patches = parseUnifiedDiff(multiDiff);
    expect(patches.length).toBe(2);
    expect(patches[0]?.newFile).toBe('file1.ts');
    expect(patches[1]?.newFile).toBe('file2.ts');
  });
});
