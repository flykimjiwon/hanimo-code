package build

import "testing"

// Each profile must export non-empty Profile/Recommended* constants so that
// frontend BuildProfile() / RecommendedDefaults bindings always return
// something usable. This single test runs in whichever profile is active
// at compile time (default by default, onprem with -tags=onprem) — no
// way to test both in one invocation, so CI must run both build tags.
func TestProfileConstants_NonEmpty(t *testing.T) {
	if Profile == "" {
		t.Error("Profile must be set")
	}
	if RecommendedProvider == "" {
		t.Error("RecommendedProvider must be set")
	}
	if RecommendedSuperModel == "" {
		t.Error("RecommendedSuperModel must be set")
	}
	if RecommendedDevModel == "" {
		t.Error("RecommendedDevModel must be set")
	}
}

func TestProfileConstants_KnownProfile(t *testing.T) {
	switch Profile {
	case "default", "onprem":
		// ok
	default:
		t.Errorf("unknown profile %q — expected default or onprem", Profile)
	}
}
