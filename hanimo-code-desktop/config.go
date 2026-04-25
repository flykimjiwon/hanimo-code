package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// TGCConfig mirrors .hanimo/config.yaml so we share settings with the TUI.
type TGCConfig struct {
	API struct {
		BaseURL string `yaml:"base_url"`
		APIKey  string `yaml:"api_key"`
	} `yaml:"api"`
	Models struct {
		Super string `yaml:"super"`
		Dev   string `yaml:"dev"`
	} `yaml:"models"`
	// Phase 15a — per-provider override map. Keyed by provider name from
	// providers.go (ollama / openai / anthropic / novita / openrouter / …).
	// Each entry can override base_url and supply api_key. If left empty we
	// fall back to the catalog default + env var.
	Providers map[string]ProviderEntry `yaml:"providers,omitempty"`
	// Active provider name — set by SwitchModel based on the chosen model's
	// provider field. Used by newChatEngine to pick base_url + api_key.
	Provider string `yaml:"provider,omitempty"`
}

// ProviderEntry is one row in the providers map.
type ProviderEntry struct {
	BaseURL string `yaml:"base_url,omitempty"`
	APIKey  string `yaml:"api_key,omitempty"`
}

const (
	defaultBaseURL = "http://localhost:11434/v1"
	defaultModel   = "qwen3:8b"
)

// LoadTGCConfig reads the shared .hanimo/config.yaml (or .hanimo/).
func LoadTGCConfig() TGCConfig {
	cfg := TGCConfig{}
	cfg.API.BaseURL = defaultBaseURL
	cfg.Models.Super = defaultModel
	cfg.Models.Dev = defaultModel

	home, _ := os.UserHomeDir()

	// Try config dirs in priority order
	dirs := []string{".hanimo", ".hanimo"}
	for _, dir := range dirs {
		path := filepath.Join(home, dir, "config.yaml")
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		_ = yaml.Unmarshal(data, &cfg)
		break
	}

	// Env overrides
	if v := os.Getenv("HANIMO_API_BASE_URL"); v != "" {
		cfg.API.BaseURL = v
	}
	if v := os.Getenv("HANIMO_API_KEY"); v != "" {
		cfg.API.APIKey = v
	}
	if v := os.Getenv("HANIMO_MODEL_SUPER"); v != "" {
		cfg.Models.Super = v
	}

	return cfg
}
