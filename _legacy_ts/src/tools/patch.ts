import { tool } from 'ai';
import { z } from 'zod';
import { readFile, writeFile } from 'node:fs/promises';
import { resolve } from 'node:path';
import { checkPathSandbox } from '../core/permission.js';

interface Hunk {
  oldStart: number;
  oldCount: number;
  newStart: number;
  newCount: number;
  lines: string[];
}

interface FilePatch {
  oldFile: string;
  newFile: string;
  hunks: Hunk[];
}

/**
 * Parse unified diff text into structured file patches.
 */
export function parseUnifiedDiff(diff: string): FilePatch[] {
  const patches: FilePatch[] = [];
  const lines = diff.split('\n');
  let i = 0;

  while (i < lines.length) {
    const line = lines[i] ?? '';

    // Look for --- header
    if (line.startsWith('--- ')) {
      const oldFile = (line.slice(4).split('\t')[0] ?? '').replace(/^[ab]\//, '');
      i++;
      const plusLine = lines[i] ?? '';
      if (!plusLine.startsWith('+++ ')) {
        i++;
        continue;
      }
      const newFile = (plusLine.slice(4).split('\t')[0] ?? '').replace(/^[ab]\//, '');
      i++;

      const hunks: Hunk[] = [];

      // Parse hunks
      while (i < lines.length) {
        const hunkLine = lines[i] ?? '';
        if (hunkLine.startsWith('--- ')) break; // Next file

        const hunkMatch = hunkLine.match(/^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@/);
        if (!hunkMatch) {
          i++;
          continue;
        }

        const hunk: Hunk = {
          oldStart: parseInt(hunkMatch[1] ?? '1', 10),
          oldCount: parseInt(hunkMatch[2] ?? '1', 10),
          newStart: parseInt(hunkMatch[3] ?? '1', 10),
          newCount: parseInt(hunkMatch[4] ?? '1', 10),
          lines: [],
        };
        i++;

        // Collect hunk lines
        while (i < lines.length) {
          const cl = lines[i] ?? '';
          if (cl.startsWith('@@') || cl.startsWith('--- ') || cl.startsWith('diff ')) break;
          if (cl.startsWith('+') || cl.startsWith('-') || cl.startsWith(' ') || cl === '') {
            hunk.lines.push(cl);
          } else if (cl.startsWith('\\')) {
            // "\ No newline at end of file" — skip
          } else {
            break;
          }
          i++;
        }

        hunks.push(hunk);
      }

      if (hunks.length > 0) {
        patches.push({ oldFile, newFile, hunks });
      }
    } else {
      i++;
    }
  }

  return patches;
}

/**
 * Apply hunks to file content. Applies from bottom to top to preserve line offsets.
 */
export function applyHunks(content: string, hunks: Hunk[]): { result: string; applied: number; rejected: string[] } {
  const fileLines = content.split('\n');
  let applied = 0;
  const rejected: string[] = [];

  // Sort hunks by oldStart descending (apply from bottom)
  const sorted = [...hunks].sort((a, b) => b.oldStart - a.oldStart);

  for (const hunk of sorted) {
    const oldLines: string[] = [];
    const newLines: string[] = [];

    for (const line of hunk.lines) {
      if (line.startsWith('-')) {
        oldLines.push(line.slice(1));
      } else if (line.startsWith('+')) {
        newLines.push(line.slice(1));
      } else {
        // Context line (starts with ' ' or is empty)
        const ctx = line.startsWith(' ') ? line.slice(1) : line;
        oldLines.push(ctx);
        newLines.push(ctx);
      }
    }

    // Find the hunk location — try exact match at expected line, then fuzzy
    const startIdx = hunk.oldStart - 1;
    let matchIdx = -1;

    // Try exact position first
    if (matchesAt(fileLines, oldLines, startIdx)) {
      matchIdx = startIdx;
    } else {
      // Fuzzy search: try nearby lines (within 20 lines)
      for (let offset = 1; offset <= 20; offset++) {
        if (matchesAt(fileLines, oldLines, startIdx - offset)) {
          matchIdx = startIdx - offset;
          break;
        }
        if (matchesAt(fileLines, oldLines, startIdx + offset)) {
          matchIdx = startIdx + offset;
          break;
        }
      }
    }

    if (matchIdx >= 0) {
      fileLines.splice(matchIdx, oldLines.length, ...newLines);
      applied++;
    } else {
      rejected.push(`@@ -${hunk.oldStart},${hunk.oldCount} +${hunk.newStart},${hunk.newCount} @@`);
    }
  }

  return { result: fileLines.join('\n'), applied, rejected };
}

function matchesAt(fileLines: string[], oldLines: string[], startIdx: number): boolean {
  if (startIdx < 0 || startIdx + oldLines.length > fileLines.length) return false;
  for (let j = 0; j < oldLines.length; j++) {
    if (fileLines[startIdx + j] !== oldLines[j]) return false;
  }
  return true;
}

export const patchTool = tool({
  description:
    'Apply a unified diff patch to one or more files. Useful for making multiple ' +
    'related changes across a file in a single operation. Accepts standard unified diff format (git diff output).',
  parameters: z.object({
    patch: z.string().describe('Unified diff content (--- a/file ... +++ b/file ...)'),
    cwd: z.string().optional().describe('Working directory for relative paths'),
    dryRun: z
      .boolean()
      .default(false)
      .describe('If true, validate patch without applying'),
  }),
  execute: async ({ patch, cwd, dryRun }) => {
    try {
      const workDir = cwd ?? process.cwd();
      const patches = parseUnifiedDiff(patch);

      if (patches.length === 0) {
        return {
          success: false,
          error: 'No valid patches found. Ensure the input is in unified diff format (--- a/file ... +++ b/file ... @@ ... @@)',
        };
      }

      const filesModified: string[] = [];
      let totalHunks = 0;
      const allRejected: string[] = [];

      for (const fp of patches) {
        const filePath = resolve(workDir, fp.newFile);

        // Path sandbox check
        const blocked = checkPathSandbox(filePath, workDir);
        if (blocked) {
          return { success: false, error: blocked };
        }

        let content: string;
        try {
          content = await readFile(filePath, 'utf-8');
        } catch {
          return {
            success: false,
            error: `Cannot read file: ${fp.newFile}`,
          };
        }

        const { result, applied, rejected } = applyHunks(content, fp.hunks);

        if (rejected.length > 0) {
          allRejected.push(...rejected.map(r => `${fp.newFile}: ${r}`));
        }

        if (applied > 0 && !dryRun) {
          await writeFile(filePath, result, 'utf-8');
        }

        if (applied > 0) {
          filesModified.push(fp.newFile);
          totalHunks += applied;
        }
      }

      return {
        success: allRejected.length === 0,
        dryRun,
        filesModified,
        hunksApplied: totalHunks,
        ...(allRejected.length > 0 ? { rejected: allRejected } : {}),
      };
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : String(err);
      return { success: false, error: `Patch failed: ${message}` };
    }
  },
});
