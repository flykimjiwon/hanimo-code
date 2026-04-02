import { describe, it, expect, beforeEach } from 'vitest';
import { getPermissionLevel, setPermissionRules, checkPathSandbox } from '../src/core/permission.js';
import { initFeatureFlags } from '../src/core/feature-flags.js';

describe('permission-glob', () => {
  beforeEach(() => {
    setPermissionRules([]);
    initFeatureFlags({});
  });

  it('should allow read-only tools without glob rules', () => {
    expect(getPermissionLevel('read_file', true)).toBe('allow');
    expect(getPermissionLevel('glob_search', true)).toBe('allow');
    expect(getPermissionLevel('grep_search', true)).toBe('allow');
    expect(getPermissionLevel('git_status', true)).toBe('allow');
  });

  it('should ask for destructive tools when approval required', () => {
    expect(getPermissionLevel('write_file', true)).toBe('ask');
    expect(getPermissionLevel('edit_file', true)).toBe('ask');
    expect(getPermissionLevel('shell_exec', true)).toBe('ask');
    expect(getPermissionLevel('git_commit', true)).toBe('ask');
  });

  it('should allow destructive tools when approval not required', () => {
    expect(getPermissionLevel('write_file', false)).toBe('allow');
    expect(getPermissionLevel('shell_exec', false)).toBe('allow');
  });

  it('should apply deny rules with glob matching', () => {
    initFeatureFlags({ PERMISSION_GLOB_RULES: true });
    setPermissionRules([
      { tool: 'shell_exec', action: 'deny' },
    ]);
    expect(getPermissionLevel('shell_exec', true)).toBe('deny');
  });

  it('should apply allow rules via glob pattern', () => {
    initFeatureFlags({ PERMISSION_GLOB_RULES: true });
    setPermissionRules([
      { tool: 'write_*', action: 'allow' },
    ]);
    expect(getPermissionLevel('write_file', true)).toBe('allow');
  });

  it('should prioritize deny over allow rules', () => {
    initFeatureFlags({ PERMISSION_GLOB_RULES: true });
    setPermissionRules([
      { tool: 'shell_exec', action: 'allow' },
      { tool: 'shell_exec', action: 'deny' },
    ]);
    // deny rules are sorted first
    expect(getPermissionLevel('shell_exec', true)).toBe('deny');
  });

  it('should block sensitive file paths', () => {
    const blocked = checkPathSandbox('/home/user/.ssh/id_rsa', '/home/user/project');
    expect(blocked).not.toBeNull();
    expect(blocked).toContain('sensitive');
  });

  it('should block path traversal', () => {
    const blocked = checkPathSandbox('../../etc/passwd', '/home/user/project');
    expect(blocked).not.toBeNull();
    expect(blocked).toContain('outside');
  });

  it('should allow paths within sandbox', () => {
    const result = checkPathSandbox('src/index.ts', '/home/user/project');
    expect(result).toBeNull();
  });
});
