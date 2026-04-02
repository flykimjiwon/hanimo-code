import { getModelCapability } from './model-capabilities.js';

/**
 * Calculate the maximum number of context messages based on model's context window.
 * Uses conservative estimates: 1 message ≈ 500 tokens avg, system prompt ≈ 2000 tokens.
 * Only uses 70% of context window to leave room for response generation.
 */
export function getMaxContextMessages(modelName: string, provider?: string): number {
  const cap = getModelCapability(modelName, provider);
  const ctx = cap.contextWindow ?? 8192; // Safe default for unknown models

  // Usable tokens: 70% of context minus system prompt overhead
  const usableTokens = Math.floor(ctx * 0.7) - 2000;
  const maxMessages = Math.floor(usableTokens / 500);

  // Clamp: minimum 20 (usable conversation), maximum 200 (diminishing returns)
  return Math.max(20, Math.min(200, maxMessages));
}
