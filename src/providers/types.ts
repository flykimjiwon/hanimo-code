export const PROVIDER_NAMES = [
  // Cloud APIs (native SDK)
  'openai',
  'anthropic',
  'google',
  // Cloud APIs (OpenAI-compatible)
  'deepseek',
  'groq',
  'together',
  'openrouter',
  'fireworks',
  'mistral',
  'glm',
  // Local / self-hosted
  'ollama',
  'vllm',
  'lmstudio',
  'custom',
] as const;

export type ProviderName = typeof PROVIDER_NAMES[number];

export interface ProviderConfig {
  apiKey?: string;
  baseURL?: string;
  defaultModel?: string;
}

// Providers that run locally (no API key needed, tools may not work)
export const LOCAL_PROVIDERS = new Set<ProviderName>([
  'ollama',
  'vllm',
  'lmstudio',
  'custom',
]);

// Default base URLs for OpenAI-compatible providers
export const PROVIDER_BASE_URLS: Partial<Record<ProviderName, string>> = {
  deepseek: 'https://api.deepseek.com/v1',
  groq: 'https://api.groq.com/openai/v1',
  together: 'https://api.together.xyz/v1',
  openrouter: 'https://openrouter.ai/api/v1',
  fireworks: 'https://api.fireworks.ai/inference/v1',
  mistral: 'https://api.mistral.ai/v1',
  glm: 'https://open.bigmodel.cn/api/paas/v4',
  ollama: 'http://localhost:11434/v1',
  vllm: 'http://localhost:8000/v1',
  lmstudio: 'http://localhost:1234/v1',
};

// Recommended coding models per provider (20-40B sweet spot highlighted)
export const KNOWN_MODELS: Record<string, string[]> = {
  openai: ['gpt-4o', 'gpt-4o-mini', 'gpt-4.1', 'gpt-4.1-mini', 'gpt-4.1-nano', 'o3-mini', 'codex-mini-latest'],
  anthropic: ['claude-sonnet-4-20250514', 'claude-haiku-4-20250414', 'claude-opus-4-20250514'],
  google: ['gemini-2.5-flash', 'gemini-2.5-pro', 'gemini-2.0-flash'],
  deepseek: ['deepseek-chat', 'deepseek-coder', 'deepseek-reasoner'],
  groq: ['qwen-qwq-32b', 'llama-3.3-70b-versatile', 'deepseek-r1-distill-llama-70b-specdec', 'llama-3.1-8b-instant'],
  together: ['Qwen/Qwen2.5-Coder-32B-Instruct', 'deepseek-ai/DeepSeek-V3', 'meta-llama/Llama-3.3-70B-Instruct-Turbo', 'Qwen/Qwen2.5-72B-Instruct-Turbo'],
  openrouter: ['deepseek/deepseek-chat-v3-0324:free', 'qwen/qwen3-coder', 'google/gemini-2.5-flash-preview', 'anthropic/claude-sonnet-4'],
  fireworks: ['accounts/fireworks/models/qwen2p5-coder-32b-instruct', 'accounts/fireworks/models/deepseek-v3', 'accounts/fireworks/models/llama-v3p3-70b-instruct'],
  mistral: ['codestral-latest', 'mistral-large-latest', 'mistral-small-latest'],
  glm: ['glm-4-plus', 'glm-4-flash', 'codegeex-4'],
  // Local models — 20-40B coding sweet spot
  ollama: ['qwen3:32b', 'qwen3-coder:30b', 'deepseek-coder-v2:16b', 'codellama:34b', 'codegemma:7b', 'starcoder2:15b', 'devstral:24b', 'llama3.2'],
  vllm: [],
  lmstudio: [],
};
