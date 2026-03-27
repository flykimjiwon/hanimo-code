import { describe, it, expect } from 'vitest';
import { detectNetworkMode } from '../src/mcp/network.js';
import { createToolSetFromList, mergeToolSets } from '../src/tools/registry.js';

describe('detectNetworkMode', () => {
  it('returns "online" when configured as "online"', async () => {
    const result = await detectNetworkMode('online');
    expect(result).toBe('online');
  });

  it('returns "offline" when configured as "offline"', async () => {
    const result = await detectNetworkMode('offline');
    expect(result).toBe('offline');
  });

  it('returns online or offline for "auto" (DNS probe)', async () => {
    const result = await detectNetworkMode('auto');
    expect(['online', 'offline']).toContain(result);
  });
});

describe('createToolSetFromList', () => {
  it('returns undefined for empty list', () => {
    expect(createToolSetFromList([])).toBeUndefined();
  });

  it('returns correct tools for valid names', () => {
    const toolSet = createToolSetFromList(['read_file', 'glob_search']);
    expect(toolSet).toBeDefined();
    const keys = Object.keys(toolSet!);
    expect(keys).toContain('read_file');
    expect(keys).toContain('glob_search');
    expect(keys.length).toBe(2);
  });

  it('ignores unknown tool names', () => {
    const toolSet = createToolSetFromList(['read_file', 'nonexistent_tool']);
    expect(toolSet).toBeDefined();
    expect(Object.keys(toolSet!)).toEqual(['read_file']);
  });

  it('returns undefined if all names are unknown', () => {
    expect(createToolSetFromList(['fake1', 'fake2'])).toBeUndefined();
  });
});

describe('mergeToolSets', () => {
  it('returns undefined for no inputs', () => {
    expect(mergeToolSets()).toBeUndefined();
  });

  it('returns undefined for all undefined inputs', () => {
    expect(mergeToolSets(undefined, undefined)).toBeUndefined();
  });

  it('merges two tool sets', () => {
    const a = createToolSetFromList(['read_file']);
    const b = createToolSetFromList(['write_file']);
    const merged = mergeToolSets(a, b);
    expect(merged).toBeDefined();
    const keys = Object.keys(merged!);
    expect(keys).toContain('read_file');
    expect(keys).toContain('write_file');
    expect(keys.length).toBe(2);
  });

  it('skips undefined sets in the middle', () => {
    const a = createToolSetFromList(['read_file']);
    const merged = mergeToolSets(a, undefined, createToolSetFromList(['grep_search']));
    expect(merged).toBeDefined();
    expect(Object.keys(merged!).sort()).toEqual(['grep_search', 'read_file']);
  });

  it('later sets override earlier on key conflict', () => {
    const a = createToolSetFromList(['read_file']);
    const b = createToolSetFromList(['read_file', 'write_file']);
    const merged = mergeToolSets(a, b);
    expect(merged).toBeDefined();
    // read_file from b should override a's, but both have the same tool reference
    expect(Object.keys(merged!).sort()).toEqual(['read_file', 'write_file']);
  });
});
