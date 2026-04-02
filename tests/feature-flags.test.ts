import { describe, it, expect, afterEach } from 'vitest';
import { initFeatureFlags, isEnabled, getAllFlags } from '../src/core/feature-flags.js';

describe('feature-flags', () => {
  const originalEnv = { ...process.env };

  afterEach(() => {
    process.env = { ...originalEnv };
    initFeatureFlags({});
  });

  it('should default all flags to false', () => {
    initFeatureFlags({});
    expect(isEnabled('HOOK_SYSTEM')).toBe(false);
    expect(isEnabled('MEMORY_SYSTEM')).toBe(false);
    expect(isEnabled('HEADLESS_MODE')).toBe(false);
    expect(isEnabled('MULTI_STAGE_COMPACTION')).toBe(false);
  });

  it('should respect config flags', () => {
    initFeatureFlags({ HOOK_SYSTEM: true, MEMORY_SYSTEM: true });
    expect(isEnabled('HOOK_SYSTEM')).toBe(true);
    expect(isEnabled('MEMORY_SYSTEM')).toBe(true);
    expect(isEnabled('HEADLESS_MODE')).toBe(false);
  });

  it('should allow env override to enable', () => {
    process.env['HANIMO_FF_HEADLESS_MODE'] = '1';
    initFeatureFlags({});
    expect(isEnabled('HEADLESS_MODE')).toBe(true);
  });

  it('should allow env override to disable', () => {
    process.env['HANIMO_FF_HOOK_SYSTEM'] = '0';
    initFeatureFlags({ HOOK_SYSTEM: true });
    expect(isEnabled('HOOK_SYSTEM')).toBe(false);
  });

  it('should return all flags via getAllFlags', () => {
    initFeatureFlags({ PROMPT_CACHE: true });
    const flags = getAllFlags();
    expect(flags.PROMPT_CACHE).toBe(true);
    expect(flags.HOOK_SYSTEM).toBe(false);
    // Verify it returns a copy (not the internal object)
    flags.PROMPT_CACHE = false;
    expect(isEnabled('PROMPT_CACHE')).toBe(true);
  });
});
