import type { LanguageModelV1 } from 'ai';
import { createOpenAI } from '@ai-sdk/openai';
import { createAnthropic } from '@ai-sdk/anthropic';
import { createGoogleGenerativeAI } from '@ai-sdk/google';
import type { ProviderName, ProviderConfig } from './types.js';
import { PROVIDER_BASE_URLS, LOCAL_PROVIDERS } from './types.js';

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
