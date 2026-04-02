import { describe, it, expect, beforeEach } from 'vitest';
import { executeHooks } from '../src/core/hooks.js';
import { initFeatureFlags } from '../src/core/feature-flags.js';
import type { HookConfig } from '../src/config/schema.js';

const emptyHookConfig: HookConfig = {
  PreToolUse: [],
  PostToolUse: [],
  SessionStart: [],
  SessionStop: [],
  UserPromptSubmit: [],
};

describe('hooks', () => {
  beforeEach(() => {
    initFeatureFlags({});
  });

  it('should return non-blocked when hook system is disabled', async () => {
    initFeatureFlags({ HOOK_SYSTEM: false });
    const result = await executeHooks('PreToolUse', { toolName: 'shell_exec' }, emptyHookConfig);
    expect(result.blocked).toBe(false);
    expect(result.systemReminders).toEqual([]);
  });

  it('should return non-blocked when no hooks are configured', async () => {
    initFeatureFlags({ HOOK_SYSTEM: true });
    // executeHooks reads from ~/.hanimo/config.json which won't exist in test
    const result = await executeHooks('SessionStart', {}, emptyHookConfig);
    expect(result.blocked).toBe(false);
    expect(result.systemReminders).toEqual([]);
  });

  it('should handle PreToolUse event type', async () => {
    initFeatureFlags({ HOOK_SYSTEM: true });
    const result = await executeHooks('PreToolUse', { toolName: 'write_file' }, emptyHookConfig);
    expect(result).toHaveProperty('blocked');
    expect(result).toHaveProperty('systemReminders');
    expect(Array.isArray(result.systemReminders)).toBe(true);
  });

  it('should handle PostToolUse event type', async () => {
    initFeatureFlags({ HOOK_SYSTEM: true });
    const result = await executeHooks('PostToolUse', { toolName: 'read_file' }, emptyHookConfig);
    expect(result.blocked).toBe(false);
  });

  it('should handle UserPromptSubmit with user prompt context', async () => {
    initFeatureFlags({ HOOK_SYSTEM: true });
    const result = await executeHooks('UserPromptSubmit', { userPrompt: 'hello' }, emptyHookConfig);
    expect(result.blocked).toBe(false);
  });

  it('should handle SessionStop event', async () => {
    initFeatureFlags({ HOOK_SYSTEM: true });
    const result = await executeHooks('SessionStop', {}, emptyHookConfig);
    expect(result.blocked).toBe(false);
  });
});
