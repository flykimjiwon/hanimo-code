package main

// Phase 15a — Multi-provider catalog.
//
// hanimo-code-desktop currently keeps a single `api.base_url` + `api.api_key`
// in TGCConfig. To make Ollama / Novita / OpenRouter / OpenAI / DeepSeek /
// Mistral all work transparently, this file maps provider name → endpoint +
// env var, and lets SwitchModel rewrite base_url + key automatically when
// the user picks a model.
//
// Anthropic + Google use non-OpenAI APIs (messages, generateContent). They
// are LISTED here as catalogue entries so the dropdown shows them, but the
// actual chat call still goes through the OpenAI SDK — so picking Claude
// Sonnet 4.6 from the dropdown will currently fail at request time.
// A future phase will add a thin transport switch in chat.go.

import (
	"fmt"
	"os"
	"strings"
)

// Provider describes one LLM endpoint family.
type Provider struct {
	Name     string // canonical key — "ollama" | "openai" | "anthropic" | …
	Label    string // shown in UI
	BaseURL  string // OpenAI-compatible /v1 root (or native API root)
	EnvVar   string // env var that holds the API key (empty = no key needed)
	OpenAI   bool   // true if /v1/chat/completions OpenAI-compatible
	KeyHint  string // free-text help shown when key is missing
}

// providerCatalog is the single source of truth. New providers added here
// are immediately routable from the model dropdown without further code.
var providerCatalog = []Provider{
	{
		Name:    "ollama",
		Label:   "Ollama (local)",
		BaseURL: "http://localhost:11434/v1",
		EnvVar:  "",
		OpenAI:  true,
		KeyHint: "no key needed — install from https://ollama.com",
	},
	{
		Name:    "openai",
		Label:   "OpenAI",
		BaseURL: "https://api.openai.com/v1",
		EnvVar:  "OPENAI_API_KEY",
		OpenAI:  true,
		KeyHint: "https://platform.openai.com/api-keys",
	},
	{
		Name:    "anthropic",
		Label:   "Anthropic",
		BaseURL: "https://api.anthropic.com/v1",
		EnvVar:  "ANTHROPIC_API_KEY",
		OpenAI:  false,
		KeyHint: "https://console.anthropic.com/settings/keys",
	},
	{
		Name:    "google",
		Label:   "Google AI",
		BaseURL: "https://generativelanguage.googleapis.com/v1beta",
		EnvVar:  "GOOGLE_API_KEY",
		OpenAI:  false,
		KeyHint: "https://aistudio.google.com/apikey",
	},
	{
		Name:    "novita",
		Label:   "Novita",
		BaseURL: "https://api.novita.ai/v3/openai",
		EnvVar:  "NOVITA_API_KEY",
		OpenAI:  true,
		KeyHint: "https://novita.ai/settings/key-management",
	},
	{
		Name:    "openrouter",
		Label:   "OpenRouter",
		BaseURL: "https://openrouter.ai/api/v1",
		EnvVar:  "OPENROUTER_API_KEY",
		OpenAI:  true,
		KeyHint: "https://openrouter.ai/keys",
	},
	{
		Name:    "deepseek",
		Label:   "DeepSeek",
		BaseURL: "https://api.deepseek.com/v1",
		EnvVar:  "DEEPSEEK_API_KEY",
		OpenAI:  true,
		KeyHint: "https://platform.deepseek.com/api_keys",
	},
	{
		Name:    "mistral",
		Label:   "Mistral",
		BaseURL: "https://api.mistral.ai/v1",
		EnvVar:  "MISTRAL_API_KEY",
		OpenAI:  true,
		KeyHint: "https://console.mistral.ai/api-keys",
	},
	{
		Name:    "groq",
		Label:   "Groq",
		BaseURL: "https://api.groq.com/openai/v1",
		EnvVar:  "GROQ_API_KEY",
		OpenAI:  true,
		KeyHint: "https://console.groq.com/keys",
	},
	{
		Name:    "together",
		Label:   "Together",
		BaseURL: "https://api.together.xyz/v1",
		EnvVar:  "TOGETHER_API_KEY",
		OpenAI:  true,
		KeyHint: "https://api.together.xyz/settings/api-keys",
	},
}

// findProvider returns the catalog entry by name. nil if unknown.
func findProvider(name string) *Provider {
	if name == "" {
		return nil
	}
	for i := range providerCatalog {
		if providerCatalog[i].Name == name {
			return &providerCatalog[i]
		}
	}
	return nil
}

// resolveAPIKey returns the configured key for a provider, falling back
// to the matching env var. Empty string if neither is set.
func resolveAPIKey(p *Provider, cfg TGCConfig) string {
	if p == nil {
		return ""
	}
	if perProvider, ok := cfg.Providers[p.Name]; ok && perProvider.APIKey != "" {
		return perProvider.APIKey
	}
	if p.EnvVar != "" {
		if v := strings.TrimSpace(os.Getenv(p.EnvVar)); v != "" {
			return v
		}
	}
	return ""
}

// resolveBaseURL prefers a per-provider override in config.yaml, falling
// back to the catalog default. Useful for self-hosted vLLM / custom DNS.
func resolveBaseURL(p *Provider, cfg TGCConfig) string {
	if p == nil {
		return ""
	}
	if perProvider, ok := cfg.Providers[p.Name]; ok && perProvider.BaseURL != "" {
		return perProvider.BaseURL
	}
	return p.BaseURL
}

// ProviderListEntry is the frontend-visible projection of one provider
// with its current key/baseUrl status, for a "Providers" Settings tab.
type ProviderListEntry struct {
	Name    string `json:"name"`
	Label   string `json:"label"`
	BaseURL string `json:"baseUrl"`
	EnvVar  string `json:"envVar"`
	OpenAI  bool   `json:"openaiCompatible"`
	HasKey  bool   `json:"hasKey"`
	KeyHint string `json:"keyHint"`
}

// GetProviders is a Wails binding that exposes the catalog + key status
// to the UI. Keys themselves are never returned — only "configured" bools.
func (a *App) GetProviders() []ProviderListEntry {
	cfg := LoadTGCConfig()
	out := make([]ProviderListEntry, 0, len(providerCatalog))
	for _, p := range providerCatalog {
		out = append(out, ProviderListEntry{
			Name:    p.Name,
			Label:   p.Label,
			BaseURL: resolveBaseURL(&p, cfg),
			EnvVar:  p.EnvVar,
			OpenAI:  p.OpenAI,
			HasKey:  resolveAPIKey(&p, cfg) != "" || p.EnvVar == "",
			KeyHint: p.KeyHint,
		})
	}
	return out
}

// errMissingKey is returned by SwitchModel when the user picks a provider
// without a configured key. Frontend surfaces a toast linking to KeyHint.
func errMissingKey(p *Provider) error {
	return fmt.Errorf("API key required for %s — set $%s or providers.%s.api_key (%s)",
		p.Label, p.EnvVar, p.Name, p.KeyHint)
}
