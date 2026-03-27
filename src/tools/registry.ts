import type { ToolSet } from 'ai';
import { readFileTool, writeFileTool, editFileTool } from './file-ops.js';
import { gitStatusTool, gitDiffTool, gitCommitTool, gitLogTool } from './git-tools.js';
import { shellExecTool } from './shell-exec.js';
import { globSearchTool } from './glob-search.js';
import { grepSearchTool } from './grep-search.js';

export function createToolRegistry() {
  return {
    read_file: readFileTool,
    write_file: writeFileTool,
    edit_file: editFileTool,
    git_status: gitStatusTool,
    git_diff: gitDiffTool,
    git_commit: gitCommitTool,
    git_log: gitLogTool,
    shell_exec: shellExecTool,
    glob_search: globSearchTool,
    grep_search: grepSearchTool,
  };
}

// Read-only tools for Assistant role (no write/edit/shell/git-commit)
export function createReadOnlyTools() {
  return {
    read_file: readFileTool,
    glob_search: globSearchTool,
    grep_search: grepSearchTool,
  };
}

/**
 * Create a ToolSet from a list of tool names.
 * Used by the role system to build tool sets from role JSON definitions.
 */
export function createToolSetFromList(toolNames: string[]): ToolSet | undefined {
  if (toolNames.length === 0) return undefined;
  const allTools = createToolRegistry();
  const filtered: Record<string, unknown> = {};
  for (const name of toolNames) {
    const tool = allTools[name as keyof ReturnType<typeof createToolRegistry>];
    if (tool) {
      filtered[name] = tool;
    }
  }
  return Object.keys(filtered).length > 0 ? (filtered as ToolSet) : undefined;
}

/**
 * Merge multiple ToolSets into one. Later sets override earlier ones on conflict.
 */
export function mergeToolSets(...sets: Array<ToolSet | undefined>): ToolSet | undefined {
  const merged: Record<string, unknown> = {};
  for (const set of sets) {
    if (!set) continue;
    Object.assign(merged, set);
  }
  return Object.keys(merged).length > 0 ? (merged as ToolSet) : undefined;
}

export {
  readFileTool,
  writeFileTool,
  editFileTool,
  gitStatusTool,
  gitDiffTool,
  gitCommitTool,
  gitLogTool,
  shellExecTool,
  globSearchTool,
  grepSearchTool,
};
