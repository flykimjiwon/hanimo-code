package main

import "hanimo-code-desktop/internal/build"

// BuildProfileInfo is the frontend projection of compile-time profile
// constants from internal/build. Used by status-bar badges and a future
// first-run wizard or "Reset to Recommended" action.
type BuildProfileInfo struct {
	Profile              string `json:"profile"`              // "default" | "onprem"
	RecommendedProvider  string `json:"recommendedProvider"`  // catalog provider name
	RecommendedSuperModel string `json:"recommendedSuperModel"`
	RecommendedDevModel   string `json:"recommendedDevModel"`
}

// GetBuildProfile returns the compile-time profile + recommended defaults.
// Frontend can use Profile to render a 사내망/외부망 badge and the model
// fields to seed a config wizard. Note: this does NOT reflect the user's
// current ~/.hanimo/config.yaml — it is the recommended seed for fresh
// installs only.
func (a *App) GetBuildProfile() BuildProfileInfo {
	return BuildProfileInfo{
		Profile:              build.Profile,
		RecommendedProvider:  build.RecommendedProvider,
		RecommendedSuperModel: build.RecommendedSuperModel,
		RecommendedDevModel:   build.RecommendedDevModel,
	}
}
