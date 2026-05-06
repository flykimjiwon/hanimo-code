package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const defaultOllamaURL = "http://localhost:11434"

// OllamaProvider wraps OpenAI-compatible for Ollama, adding native model listing.
type OllamaProvider struct {
	*OpenAICompatProvider
	nativeURL string
}

func NewOllama(baseURL, apiKey string) *OllamaProvider {
	if baseURL == "" {
		baseURL = defaultOllamaURL
	}
	// Users (and our default config) commonly pass `http://localhost:11434/v1`
	// because that's what every other OpenAI-compatible endpoint expects.
	// Ollama's native API (/api/tags) lives one level above /v1, so we strip
	// the suffix before computing nativeURL — otherwise ListModels hits
	// /v1/api/tags and gets back a JSON-but-wrong-shape Ollama error page.
	trimmed := strings.TrimRight(baseURL, "/")
	trimmed = strings.TrimSuffix(trimmed, "/v1")
	nativeURL := strings.TrimRight(trimmed, "/")
	compatURL := nativeURL + "/v1"

	return &OllamaProvider{
		OpenAICompatProvider: NewOpenAICompat("ollama", compatURL, apiKey),
		nativeURL:            nativeURL,
	}
}

func (p *OllamaProvider) Name() string { return "ollama" }

// ListModels calls Ollama's native API to list available models.
func (p *OllamaProvider) ListModels() ([]ModelInfo, error) {
	resp, err := http.Get(p.nativeURL + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("ollama list models: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ollama parse models: %w", err)
	}

	models := make([]ModelInfo, 0, len(result.Models))
	for _, m := range result.Models {
		models = append(models, ModelInfo{
			ID:          m.Name,
			DisplayName: m.Name,
			Provider:    "ollama",
		})
	}
	return models, nil
}

func init() {
	Register("ollama", func(baseURL, apiKey string) Provider {
		return NewOllama(baseURL, apiKey)
	})
}
