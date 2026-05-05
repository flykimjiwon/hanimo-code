//go:build !onprem

// Package build holds compile-time profile constants. Two profiles exist:
//
//   - default (this file, no build tag): hanimo-code-desktop for the public
//     OSS distribution. Recommended provider is Novita (vendor-prefix
//     marketplace), recommended super model is `openai/gpt-oss-120b`.
//   - onprem (`-tags=onprem`): TECHAI 사내망 / Shinhan internal vLLM.
//     Vanilla model ids without vendor-prefix.
//
// The split exists because both endpoints serve the same gpt-oss-120b model
// under DIFFERENT identifiers (Novita: "openai/gpt-oss-120b"; vanilla vLLM:
// "gpt-oss-120b"). Hardcoding either convention as a single default caused
// postmortem 2026-05-03 Bug #1 (404 MODEL_NOT_FOUND on first message after
// porting the TECHAI catalog to the public build).
//
// LoadTGCConfig keeps the zero-config Ollama localhost fallback for cold
// starts. These constants are the *recommended* defaults exposed to the
// frontend (BuildProfile / RecommendedDefaults bindings) for first-run
// wizards or a "Reset to Recommended" action — they do not change the
// behavior of an already-configured ~/.hanimo/config.yaml.
package build

const (
	Profile              = "default"
	RecommendedProvider  = "novita"
	RecommendedSuperModel = "openai/gpt-oss-120b"
	RecommendedDevModel   = "openai/gpt-oss-120b"
)
