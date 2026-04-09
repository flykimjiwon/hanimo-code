import { describe, it, expect } from 'vitest';
import { getMaxContextMessages } from '../src/providers/context-window.js';

describe('context-window', () => {
  it('should return minimum 20 for small context models', () => {
    // gemma3:1b has 8K context
    const max = getMaxContextMessages('gemma3:1b');
    expect(max).toBe(20);
  });

  it('should return appropriate value for 128K models', () => {
    // gpt-4o has 128K context
    const max = getMaxContextMessages('gpt-4o');
    // (128000 * 0.7 - 2000) / 500 = 175.2 → 175
    expect(max).toBe(175);
  });

  it('should cap at 200 for large context models', () => {
    // gemini-2.5-pro has 1M context
    const max = getMaxContextMessages('gemini-2.5-pro');
    expect(max).toBe(200);
  });

  it('should return 200 for 200K context models', () => {
    // claude-sonnet-4 has 200K context
    const max = getMaxContextMessages('claude-sonnet-4-20250514');
    // (200000 * 0.7 - 2000) / 500 = 276 → capped at 200
    expect(max).toBe(200);
  });

  it('should use safe default for unknown models', () => {
    // Unknown model defaults to 8192 context
    const max = getMaxContextMessages('unknown-model-xyz');
    // (8192 * 0.7 - 2000) / 500 = 7.47 → clamped to 20
    expect(max).toBe(20);
  });

  it('should handle Ollama models with small context', () => {
    // llama3.2:3b has 8K context
    const max = getMaxContextMessages('llama3.2:3b');
    expect(max).toBe(20);
  });
});
