package main

import (
	"errors"
	"strings"
	"testing"
)

// Tests for resolveSwitchTarget — postmortem 2026-05-03 Bug #2 root fix.
// All tests use models from the static tier1/tier2 catalogs to avoid the
// Ollama HTTP fallback inside findModelOption.

// noPing stubs the daemon probe to never error. Used when a test wants
// to exercise the Ollama success path without hitting localhost:11434.
func noPing(string) error { return nil }

// failPing stubs the daemon probe to always error. Used to assert the
// pre-flight blocks a config write when Ollama is unreachable.
func failPing(string) error { return errors.New("connection refused (test stub)") }

func TestResolveSwitchTarget_EmptyModelID_Error(t *testing.T) {
	cfg := TGCConfig{}
	_, err := resolveSwitchTarget(cfg, "  ", noPing)
	if err == nil {
		t.Fatal("expected error for empty model id, got nil")
	}
}

func TestResolveSwitchTarget_SameProvider_DoesNotTouchAPI(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-anthropic-key")

	cfg := TGCConfig{Provider: "anthropic"}
	cfg.API.BaseURL = "https://custom.anthropic.example.com/v1"
	cfg.API.APIKey = "user-supplied-key"

	got, err := resolveSwitchTarget(cfg, "claude-sonnet-4-6", noPing)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Same-provider switch must not rewrite BaseURL/APIKey — would clobber
	// the user's working setup. Only Models.Super/Dev should change.
	if got.API.BaseURL != "https://custom.anthropic.example.com/v1" {
		t.Errorf("BaseURL clobbered on same-provider switch: %q", got.API.BaseURL)
	}
	if got.API.APIKey != "user-supplied-key" {
		t.Errorf("APIKey clobbered on same-provider switch: %q", got.API.APIKey)
	}
	if got.Provider != "anthropic" {
		t.Errorf("Provider changed unexpectedly: %q", got.Provider)
	}
	if got.Models.Super != "claude-sonnet-4-6" {
		t.Errorf("Models.Super not updated: %q", got.Models.Super)
	}
	if got.Models.Dev != "claude-sonnet-4-6" {
		t.Errorf("Models.Dev not updated: %q", got.Models.Dev)
	}
}

func TestResolveSwitchTarget_CrossProvider_KeyMissing_NoMutation(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")

	cfg := TGCConfig{Provider: "anthropic"}
	cfg.API.BaseURL = "https://api.anthropic.com/v1"
	cfg.API.APIKey = "sk-anthropic-still-good"

	got, err := resolveSwitchTarget(cfg, "gpt-5", noPing)
	if err == nil {
		t.Fatal("expected errMissingKey, got nil")
	}
	// On error, the returned cfg must equal the input — no half-write.
	if got.Provider != "anthropic" {
		t.Errorf("Provider mutated despite error: %q", got.Provider)
	}
	if got.API.BaseURL != "https://api.anthropic.com/v1" {
		t.Errorf("BaseURL mutated despite error: %q", got.API.BaseURL)
	}
	if got.API.APIKey != "sk-anthropic-still-good" {
		t.Errorf("APIKey mutated despite error: %q", got.API.APIKey)
	}
}

func TestResolveSwitchTarget_CrossProvider_OllamaUnreachable_NoMutation(t *testing.T) {
	cfg := TGCConfig{Provider: "anthropic"}
	cfg.API.BaseURL = "https://api.anthropic.com/v1"
	cfg.API.APIKey = "sk-anthropic-still-good"

	got, err := resolveSwitchTarget(cfg, "qwen3-coder-30b", failPing)
	if err == nil {
		t.Fatal("expected Ollama unreachable error, got nil")
	}
	if !strings.Contains(err.Error(), "Ollama unreachable") {
		t.Errorf("error should mention Ollama unreachable: %v", err)
	}
	// This is the bug #2 acceptance test: picking an Ollama model when the
	// daemon is down must NOT clobber the user's working Anthropic config.
	if got.Provider != "anthropic" {
		t.Errorf("Provider clobbered to %q despite ping failure", got.Provider)
	}
	if got.API.BaseURL != "https://api.anthropic.com/v1" {
		t.Errorf("BaseURL clobbered to %q despite ping failure", got.API.BaseURL)
	}
	if got.API.APIKey != "sk-anthropic-still-good" {
		t.Errorf("APIKey clobbered to %q despite ping failure", got.API.APIKey)
	}
	if got.Models.Super == "qwen3-coder-30b" {
		t.Error("Models.Super updated despite ping failure — should be atomic")
	}
}

func TestResolveSwitchTarget_CrossProvider_OllamaReachable_Switches(t *testing.T) {
	cfg := TGCConfig{Provider: "anthropic"}
	cfg.API.BaseURL = "https://api.anthropic.com/v1"
	cfg.API.APIKey = "sk-anthropic"

	got, err := resolveSwitchTarget(cfg, "qwen3-coder-30b", noPing)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Provider != "ollama" {
		t.Errorf("Provider not switched to ollama: %q", got.Provider)
	}
	if got.API.BaseURL != "http://localhost:11434/v1" {
		t.Errorf("BaseURL not switched to Ollama default: %q", got.API.BaseURL)
	}
	if got.API.APIKey != "" {
		t.Errorf("APIKey should be empty for Ollama, got %q", got.API.APIKey)
	}
	if got.Models.Super != "qwen3-coder-30b" {
		t.Errorf("Models.Super not updated: %q", got.Models.Super)
	}
}

func TestResolveSwitchTarget_UnknownModel_KeepsProvider(t *testing.T) {
	cfg := TGCConfig{Provider: "anthropic"}
	cfg.API.BaseURL = "https://api.anthropic.com/v1"
	cfg.API.APIKey = "sk-anthropic"

	// "my-private-vllm-model" is not in any catalog. findModelOption returns
	// nil, so provider routing is skipped — only Models.Super/Dev change.
	// findModelOption may probe Ollama; pass through with no harm even if
	// the local daemon answers (we only assert provider/baseurl unchanged).
	got, err := resolveSwitchTarget(cfg, "my-private-vllm-model", noPing)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Provider != "anthropic" {
		t.Errorf("Provider changed for unknown model: %q", got.Provider)
	}
	if got.API.BaseURL != "https://api.anthropic.com/v1" {
		t.Errorf("BaseURL changed for unknown model: %q", got.API.BaseURL)
	}
	if got.Models.Super != "my-private-vllm-model" {
		t.Errorf("Models.Super not updated: %q", got.Models.Super)
	}
}

func TestPingOllama_EmptyURL_Error(t *testing.T) {
	if err := pingOllama(""); err == nil {
		t.Fatal("expected error for empty url, got nil")
	}
}
