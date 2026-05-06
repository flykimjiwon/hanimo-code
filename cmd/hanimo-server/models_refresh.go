// Live model discovery for the active provider.
//
// The /models GET handler in handlers.go returns the static llm.Models map —
// a curated, hand-tuned list with display names and capability flags. That
// list lags reality whenever a provider ships a new model.
//
// /models/refresh fixes that: on demand, it asks the active provider's
// ListModels() endpoint for the live list, normalises each entry against the
// curated capability data, and caches the result for the next /models GET.
//
// Cache strategy:
//   - keyed by (provider, baseURL): switching endpoints invalidates
//   - 24h TTL: auto-stale, but never auto-fetches; the client decides when
//     to call /models/refresh (Settings button, periodic from extension)
//   - merged in /models GET so refreshed entries appear without a second call
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/flykimjiwon/hanimo/internal/llm"
	"github.com/flykimjiwon/hanimo/internal/llm/providers"
)

type modelCacheEntry struct {
	Models    []modelDTO
	FetchedAt time.Time
}

var (
	modelCacheMu sync.RWMutex
	modelCache   = map[string]modelCacheEntry{}
)

const modelCacheTTL = 24 * time.Hour

// detectProviderID picks a provider name from the BaseURL. Mirrors the
// extension-side detectProvider() so the server reports the same default
// regardless of which client called it. "openai" is the catch-all.
func detectProviderID(baseURL string) string {
	u := strings.ToLower(baseURL)
	switch {
	case strings.Contains(u, "anthropic.com"):
		return "anthropic"
	case strings.Contains(u, "googleapis.com"):
		return "google"
	case strings.Contains(u, "deepseek.com"):
		return "deepseek"
	case strings.Contains(u, "novita.ai"):
		return "novita"
	case strings.Contains(u, "openrouter.ai"):
		return "openrouter"
	case strings.Contains(u, "127.0.0.1") || strings.Contains(u, "localhost"):
		return "ollama"
	default:
		return "openai"
	}
}

// providerCacheKey is what we hash on; provider+baseURL together so
// switching base_url with the same provider name still busts the cache.
func providerCacheKey(provider, baseURL string) string {
	return provider + "|" + baseURL
}

// /models/refresh — POST [{ "provider": "anthropic" }] (optional)
// On success, returns the same shape as /models GET so the client can swap
// state in one call.
func (s *server) handleModelsRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	provider := r.URL.Query().Get("provider")
	if r.Method == http.MethodPost && provider == "" {
		var body struct {
			Provider string `json:"provider"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		provider = body.Provider
	}
	if provider == "" {
		provider = detectProviderID(s.cfg.API.BaseURL)
	}

	prov, err := providers.Get(provider, s.cfg.API.BaseURL, s.cfg.API.APIKey)
	if err != nil {
		http.Error(w, "unknown provider: "+provider, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	_ = ctx // ListModels has no ctx — kept for symmetry

	live, err := prov.ListModels()
	if err != nil {
		http.Error(w, "list models: "+err.Error(), http.StatusBadGateway)
		return
	}

	out := make([]modelDTO, 0, len(live))
	for _, m := range live {
		cap := llm.GetCapability(m.ID)
		out = append(out, modelDTO{
			ID:                m.ID,
			DisplayName:       fallbackDisplayName(m),
			Description:       providerLabel(m.Provider),
			ContextWindow:     pickContextWindow(m, cap),
			SupportsTools:     m.SupportsTools || cap.SupportsTools,
			SupportsVision:    supportsVision(m.ID),
			SupportsReasoning: supportsReasoning(m.ID),
		})
	}

	modelCacheMu.Lock()
	modelCache[providerCacheKey(provider, s.cfg.API.BaseURL)] = modelCacheEntry{
		Models:    out,
		FetchedAt: time.Now(),
	}
	modelCacheMu.Unlock()

	writeJSON(w, http.StatusOK, map[string]any{
		"provider":   provider,
		"base_url":   s.cfg.API.BaseURL,
		"fetched_at": time.Now().Unix(),
		"count":      len(out),
		"models":     out,
	})
}

// liveModelsForCurrentProvider returns the cached live list if fresh, else
// nil. Read-only — never blocks on a network call.
func liveModelsForCurrentProvider(provider, baseURL string) []modelDTO {
	modelCacheMu.RLock()
	defer modelCacheMu.RUnlock()
	entry, ok := modelCache[providerCacheKey(provider, baseURL)]
	if !ok || time.Since(entry.FetchedAt) > modelCacheTTL {
		return nil
	}
	return entry.Models
}

func fallbackDisplayName(m providers.ModelInfo) string {
	if m.DisplayName != "" {
		return m.DisplayName
	}
	if registered, ok := llm.Models[m.ID]; ok && registered.DisplayName != "" {
		return registered.DisplayName
	}
	return m.ID
}

func providerLabel(name string) string {
	switch name {
	case "anthropic":
		return "Anthropic"
	case "openai":
		return "OpenAI"
	case "google":
		return "Google"
	case "ollama":
		return "Ollama (local)"
	case "deepseek":
		return "DeepSeek"
	case "novita":
		return "Novita"
	case "openrouter":
		return "OpenRouter"
	}
	if name == "" {
		return ""
	}
	return name
}

func pickContextWindow(m providers.ModelInfo, cap llm.ModelCapability) int {
	if m.ContextWindow > 0 {
		return m.ContextWindow
	}
	return cap.ContextWindow
}
