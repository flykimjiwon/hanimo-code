import type { LanguageModelV1 } from 'ai';
import { createOpenAI } from '@ai-sdk/openai';
import { createAnthropic } from '@ai-sdk/anthropic';
import { createGoogleGenerativeAI } from '@ai-sdk/google';
import type { ProviderName, ProviderConfig } from './types.js';
import { PROVIDER_BASE_URLS, LOCAL_PROVIDERS, KNOWN_MODELS } from './types.js';
import type { CustomProvider } from '../config/schema.js';

interface ProviderInstance {
  getModel(modelId: string): LanguageModelV1;
}

function createOpenAICompatible(
  name: string,
  config: ProviderConfig,
): ProviderInstance {
  const baseURL = config.baseURL ?? PROVIDER_BASE_URLS[name as ProviderName];
  if (!baseURL) {
    throw new Error(`Provider "${name}" requires a baseURL in config`);
  }
  const isLocal = LOCAL_PROVIDERS.has(name as ProviderName);
  const provider = createOpenAI({
    apiKey: config.apiKey ?? (isLocal ? 'not-needed' : undefined),
    baseURL,
    compatibility: 'compatible',
  });
  return {
    getModel(modelId: string) {
      return provider(modelId);
    },
  };
}

function createProviderInstance(
  name: ProviderName,
  config: ProviderConfig,
): ProviderInstance {
  // Native SDK providers
  switch (name) {
    case 'openai': {
      const provider = createOpenAI({
        apiKey: config.apiKey,
        ...(config.baseURL ? { baseURL: config.baseURL } : {}),
      });
      return {
        getModel(modelId: string) {
          return provider(modelId);
        },
      };
    }

    case 'anthropic': {
      const provider = createAnthropic({
        apiKey: config.apiKey,
        ...(config.baseURL ? { baseURL: config.baseURL } : {}),
      });
      return {
        getModel(modelId: string) {
          return provider(modelId);
        },
      };
    }

    case 'google': {
      const provider = createGoogleGenerativeAI({
        apiKey: config.apiKey,
        ...(config.baseURL ? { baseURL: config.baseURL } : {}),
      });
      return {
        getModel(modelId: string) {
          return provider(modelId);
        },
      };
    }

    // All other providers use OpenAI-compatible protocol
    default:
      return createOpenAICompatible(name, config);
  }
}

const providerCache = new Map<string, ProviderInstance>();

function getCacheKey(name: ProviderName, config: ProviderConfig): string {
  return `${name}:${config.baseURL ?? ''}:${config.apiKey ?? ''}`;
}

export function createProvider(
  name: ProviderName,
  config: ProviderConfig = {},
): ProviderInstance {
  const key = getCacheKey(name, config);
  let instance = providerCache.get(key);
  if (!instance) {
    instance = createProviderInstance(name, config);
    providerCache.set(key, instance);
  }
  return instance;
}

export function getModel(
  providerName: ProviderName,
  modelId: string,
  config: ProviderConfig = {},
): LanguageModelV1 {
  const key = getCacheKey(providerName, config);
  let instance = providerCache.get(key);
  if (!instance) {
    instance = createProviderInstance(providerName, config);
    providerCache.set(key, instance);
  }
  return instance.getModel(modelId);
}

export function clearProviderCache(): void {
  providerCache.clear();
}

// Custom provider registry — loaded from config.customProviders
const customProviderRegistry = new Map<string, CustomProvider>();

export function registerCustomProviders(providers: CustomProvider[]): void {
  customProviderRegistry.clear();
  for (const cp of providers) {
    customProviderRegistry.set(cp.name, cp);
    // Merge models into KNOWN_MODELS for discovery
    if (cp.models.length > 0) {
      const existing = KNOWN_MODELS[cp.name] ?? [];
      KNOWN_MODELS[cp.name] = [...new Set([...existing, ...cp.models])];
    }
  }
}

export function getCustomProvider(name: string): CustomProvider | undefined {
  return customProviderRegistry.get(name);
}

export function getCustomProviderNames(): string[] {
  return Array.from(customProviderRegistry.keys());
}

export function getModelForCustomProvider(
  providerName: string,
  modelId: string,
): LanguageModelV1 | null {
  const cp = customProviderRegistry.get(providerName);
  if (!cp) return null;

  const config: ProviderConfig = {
    baseURL: cp.baseURL,
    apiKey: cp.apiKey,
  };

  if (cp.protocol === 'anthropic') {
    const instance = createAnthropic({
      apiKey: config.apiKey ?? 'not-needed',
      baseURL: config.baseURL,
    });
    return instance(modelId);
  }

  // Default: OpenAI-compatible
  const instance = createOpenAI({
    apiKey: config.apiKey ?? 'not-needed',
    baseURL: config.baseURL,
    compatibility: 'compatible',
  });
  return instance(modelId);
}
