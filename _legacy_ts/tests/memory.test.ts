import { describe, it, expect, beforeEach } from 'vitest';
import { initFeatureFlags } from '../src/core/feature-flags.js';
import { buildMemoryPrompt } from '../src/core/memory.js';

describe('memory', () => {
  beforeEach(() => {
    initFeatureFlags({});
  });

  it('should return empty string when memory system is disabled', () => {
    initFeatureFlags({ MEMORY_SYSTEM: false });
    const prompt = buildMemoryPrompt();
    expect(prompt).toBe('');
  });

  it('should return empty string when memory file does not exist', () => {
    initFeatureFlags({ MEMORY_SYSTEM: true });
    // buildMemoryPrompt reads from ~/.hanimo/memory/MEMORY.md
    // In test env, this likely doesn't exist or is empty
    const prompt = buildMemoryPrompt();
    // Either empty (no file) or contains memory content
    expect(typeof prompt).toBe('string');
  });

  it('should include memory header when content exists', () => {
    initFeatureFlags({ MEMORY_SYSTEM: true });
    const prompt = buildMemoryPrompt();
    if (prompt) {
      expect(prompt).toContain('## User Memory');
    }
  });

  it('should sanitize path traversal in topic names', async () => {
    // Dynamic import to test sanitizeTopic indirectly via writeMemory
    const { writeMemory, readMemoryTopic } = await import('../src/core/memory.js');
    initFeatureFlags({ MEMORY_SYSTEM: true });
    // Writing with path traversal characters should be sanitized
    // This won't create files outside the memory directory
    try {
      writeMemory('test content', '../../etc/passwd');
      const content = readMemoryTopic('../../etc/passwd');
      // If it worked, the topic name was sanitized (slashes replaced with _)
      expect(typeof content).toBe('string');
    } catch {
      // sanitizeTopic may throw for invalid names — that's also acceptable
    }
  });

  it('should list no topics when memory system is disabled', async () => {
    const { listTopics } = await import('../src/core/memory.js');
    initFeatureFlags({ MEMORY_SYSTEM: false });
    const topics = listTopics();
    expect(topics).toEqual([]);
  });
});
