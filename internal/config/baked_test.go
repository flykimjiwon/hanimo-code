package config

import "testing"

func withBaked(mode, url, provider, model, key string, fn func()) {
	oldMode, oldURL, oldProvider, oldModel, oldKey := BakedMode, BakedBaseURL, BakedProvider, BakedModel, BakedAPIKey
	defer func() {
		BakedMode, BakedBaseURL, BakedProvider, BakedModel, BakedAPIKey = oldMode, oldURL, oldProvider, oldModel, oldKey
	}()
	BakedMode, BakedBaseURL, BakedProvider, BakedModel, BakedAPIKey = mode, url, provider, model, key
	fn()
}

func TestApplyBaked_Vanilla(t *testing.T) {
	withBaked("", "https://baked.example/v1", "novita", "qwen", "sk-baked", func() {
		cfg := applyBaked(DefaultConfig())
		if cfg.API.BaseURL != DefaultBaseURL {
			t.Errorf("vanilla should not apply baked URL, got %q", cfg.API.BaseURL)
		}
		if cfg.API.APIKey != "" {
			t.Errorf("vanilla should not apply baked key")
		}
	})
}

func TestApplyBaked_Distro(t *testing.T) {
	// Simulate a fresh install (no user config.yaml on disk) so the
	// model fallback branch of applyBaked actually fires.
	prev := distroUserHasConfig
	defer func() { distroUserHasConfig = prev }()
	distroUserHasConfig = false

	withBaked("distro", "https://baked.example/v1", "novita", "qwen/qwen3-coder-30b", "sk-should-ignore", func() {
		cfg := applyBaked(DefaultConfig())
		if cfg.API.BaseURL != "https://baked.example/v1" {
			t.Errorf("distro should apply baked URL, got %q", cfg.API.BaseURL)
		}
		if cfg.API.APIKey != "" {
			t.Errorf("distro must NOT apply baked key, got %q", cfg.API.APIKey)
		}
		if cfg.Models.Super != "qwen/qwen3-coder-30b" {
			t.Errorf("distro should apply baked model, got %q", cfg.Models.Super)
		}
		if cfg.Default.Provider != "novita" {
			t.Errorf("distro should apply baked provider, got %q", cfg.Default.Provider)
		}
		if IsSealed() {
			t.Errorf("distro should not report IsSealed")
		}
		if !IsDistro() {
			t.Errorf("distro should report IsDistro")
		}
	})
}

// Regression: when a user has a config.yaml on disk, distro must NOT
// rewrite their model choice even if the chosen model happens to match
// the built-in default. The old heuristic (compare to DefaultModel)
// silently swapped user-chosen models under them.
func TestDistroDoesNotClobberUserModel(t *testing.T) {
	prev := distroUserHasConfig
	defer func() { distroUserHasConfig = prev }()
	distroUserHasConfig = true

	withBaked("distro", "https://baked.example/v1", "novita", "qwen/qwen3-coder-30b", "", func() {
		cfg := DefaultConfig()
		cfg.Models.Super = DefaultModel
		cfg.Models.Dev = DefaultDevModel
		got := applyBaked(cfg)
		if got.Models.Super != DefaultModel {
			t.Errorf("distro must not rewrite user Super model when config.yaml exists, got %q", got.Models.Super)
		}
		if got.Models.Dev != DefaultDevModel {
			t.Errorf("distro must not rewrite user Dev model when config.yaml exists, got %q", got.Models.Dev)
		}
	})
}

func TestApplyBaked_Sealed(t *testing.T) {
	withBaked("sealed", "https://baked.example/v1", "novita", "qwen", "sk-sealed", func() {
		cfg := applyBaked(DefaultConfig())
		if cfg.API.APIKey != "sk-sealed" {
			t.Errorf("sealed should apply baked key, got %q", cfg.API.APIKey)
		}
		if !IsSealed() {
			t.Errorf("sealed should report IsSealed")
		}
	})
}

func TestDistroPreservesUserKey(t *testing.T) {
	withBaked("distro", "https://baked.example/v1", "novita", "qwen", "", func() {
		cfg := DefaultConfig()
		cfg.API.APIKey = "sk-user" // user already configured
		got := applyBaked(cfg)
		if got.API.APIKey != "sk-user" {
			t.Errorf("distro should not clobber user key, got %q", got.API.APIKey)
		}
	})
}
