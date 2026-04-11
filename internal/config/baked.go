package config

// Build-time baked configuration.
//
// hanimo supports three distribution modes, selected at `go build` time via
// `-ldflags -X`:
//
//   1. vanilla  — nothing baked. Users run `hanimo` and go through the
//      setup wizard to enter endpoint + key. Default for the public repo.
//
//   2. distro   — admin bakes the endpoint/provider/model into the binary
//      but leaves the API key empty. End users drop in HANIMO_API_KEY (or
//      fill it at the setup wizard) and everything else is fixed. Typical
//      for teams distributing a pre-configured hanimo to employees.
//
//   3. sealed   — admin bakes EVERYTHING including the API key. No user
//      input required; launching the binary goes straight to the TUI.
//      Typical for demos, kiosks, or fully managed environments. The
//      sealed binary MUST NOT be redistributed because it contains a
//      secret.
//
// Build mode is selected by setting BakedMode to "distro" or "sealed" at
// link time. In vanilla (the default "") all Baked* values are ignored.
//
// Example (distro):
//
//   go build -ldflags "\
//     -X 'github.com/flykimjiwon/hanimo/internal/config.BakedMode=distro' \
//     -X 'github.com/flykimjiwon/hanimo/internal/config.BakedBaseURL=https://api.novita.ai/v3/openai' \
//     -X 'github.com/flykimjiwon/hanimo/internal/config.BakedProvider=novita' \
//     -X 'github.com/flykimjiwon/hanimo/internal/config.BakedModel=qwen/qwen3-coder-30b'" \
//     -o hanimo ./cmd/hanimo
//
// Example (sealed):
//
//   go build -ldflags "\
//     -X 'github.com/flykimjiwon/hanimo/internal/config.BakedMode=sealed' \
//     -X 'github.com/flykimjiwon/hanimo/internal/config.BakedBaseURL=https://api.novita.ai/v3/openai' \
//     -X 'github.com/flykimjiwon/hanimo/internal/config.BakedProvider=novita' \
//     -X 'github.com/flykimjiwon/hanimo/internal/config.BakedModel=qwen/qwen3-coder-30b' \
//     -X 'github.com/flykimjiwon/hanimo/internal/config.BakedAPIKey=sk-…'" \
//     -o hanimo ./cmd/hanimo
//
// Convenience Makefile targets (build-distro / build-sealed) wrap these
// invocations so the admin doesn't have to memorize them.
var (
	// BakedMode selects the distribution mode.
	//   ""       — vanilla
	//   "distro" — endpoint+provider+model baked, API key user-provided
	//   "sealed" — everything baked including API key
	BakedMode = ""

	// BakedBaseURL is the API endpoint frozen into the binary. Ignored
	// unless BakedMode is "distro" or "sealed".
	BakedBaseURL = ""

	// BakedProvider is the logical provider name (openai / anthropic /
	// novita / ollama / ...). Used for display and for provider-specific
	// behaviour (e.g. Anthropic native path).
	BakedProvider = ""

	// BakedModel is the default Super/Plan model.
	BakedModel = ""

	// BakedDevModel is the default Deep Agent model. If empty, BakedModel
	// is used for both.
	BakedDevModel = ""

	// BakedAPIKey is the API key frozen into a "sealed" build. MUST NOT
	// be set for "distro" builds — leave empty and let the user supply
	// HANIMO_API_KEY at runtime.
	BakedAPIKey = ""

	// BakedBrand overrides the displayed product name (for white-label
	// enterprise distributions). Empty keeps the default "hanimo".
	BakedBrand = ""
)

// BakeInfo returns a human-readable summary of the baked configuration,
// suitable for /config output and the startup banner.
func BakeInfo() string {
	switch BakedMode {
	case "distro":
		return "distro — endpoint/provider/model baked; API key user-provided"
	case "sealed":
		return "sealed — everything baked (do not redistribute this binary)"
	default:
		return "vanilla — user-configurable"
	}
}

// IsSealed reports whether this binary was built with a baked API key.
// Sealed builds skip the setup wizard entirely.
func IsSealed() bool { return BakedMode == "sealed" && BakedAPIKey != "" }

// IsDistro reports whether this binary was built with a baked endpoint
// but no baked key (typical team distribution).
func IsDistro() bool { return BakedMode == "distro" }

// BakedBrandName returns the product name, respecting a white-label
// override from BakedBrand.
func BakedBrandName() string {
	if BakedBrand != "" {
		return BakedBrand
	}
	return "hanimo"
}

// applyBaked returns cfg with any baked values merged in. Baked values
// only override empty/default fields — user-supplied config.yaml settings
// still win over the distro bake, except for sealed mode which forces
// everything.
func applyBaked(cfg Config) Config {
	switch BakedMode {
	case "sealed":
		// Sealed: force every baked field, ignore user config.
		if BakedBaseURL != "" {
			cfg.API.BaseURL = BakedBaseURL
		}
		if BakedAPIKey != "" {
			cfg.API.APIKey = BakedAPIKey
		}
		if BakedModel != "" {
			cfg.Models.Super = BakedModel
			if BakedDevModel == "" {
				cfg.Models.Dev = BakedModel
			}
		}
		if BakedDevModel != "" {
			cfg.Models.Dev = BakedDevModel
		}
		if BakedProvider != "" {
			cfg.Default.Provider = BakedProvider
		}
	case "distro":
		// Distro: force endpoint/provider/model, leave API key to user.
		if BakedBaseURL != "" {
			cfg.API.BaseURL = BakedBaseURL
		}
		if BakedModel != "" && cfg.Models.Super == DefaultModel {
			cfg.Models.Super = BakedModel
		}
		if BakedDevModel != "" && cfg.Models.Dev == DefaultDevModel {
			cfg.Models.Dev = BakedDevModel
		} else if BakedModel != "" && cfg.Models.Dev == DefaultDevModel {
			cfg.Models.Dev = BakedModel
		}
		if BakedProvider != "" && cfg.Default.Provider == "" {
			cfg.Default.Provider = BakedProvider
		}
	}
	return cfg
}
