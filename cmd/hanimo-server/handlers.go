package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/flykimjiwon/hanimo/internal/config"
	"github.com/flykimjiwon/hanimo/internal/llm"
)

// configDTO is the JSON shape exposed to the extension. The on-disk YAML
// has more fields; we surface only what the UI is allowed to edit.
type configDTO struct {
	BaseURL  string `json:"base_url"`
	APIKey   string `json:"api_key"`
	Super    string `json:"super"`
	Dev      string `json:"dev"`
	Plan     string `json:"plan"`
	HasKey   bool   `json:"has_key"`
	BrandTag string `json:"brand_tag"`
	// IsOnprem mirrors techai's distro flag for UI symmetry. Hanimo is an
	// open-source, multi-provider engine, so this is always false unless the
	// binary was built with a sealed/distro bake.
	IsOnprem bool `json:"is_onprem"`
}

func (s *server) handleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, configToDTO(s.cfg))
	case http.MethodPatch, http.MethodPost:
		var patch struct {
			BaseURL *string `json:"base_url,omitempty"`
			APIKey  *string `json:"api_key,omitempty"`
			Super   *string `json:"super,omitempty"`
			Dev     *string `json:"dev,omitempty"`
			Plan    *string `json:"plan,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		s.chatMu.Lock()
		defer s.chatMu.Unlock()
		// Sealed/distro builds lock BaseURL and model fields to baked values.
		// Silently ignore client-side overrides — the UI marks them readonly.
		if patch.BaseURL != nil && !config.IsDistro() {
			s.cfg.API.BaseURL = strings.TrimSpace(*patch.BaseURL)
		}
		if patch.APIKey != nil {
			s.cfg.API.APIKey = strings.TrimSpace(*patch.APIKey)
		}
		if patch.Super != nil {
			s.cfg.Models.Super = strings.TrimSpace(*patch.Super)
		}
		if patch.Dev != nil {
			s.cfg.Models.Dev = strings.TrimSpace(*patch.Dev)
		}
		// Plan is a synthetic field — hanimo's ModelsConfig has no Plan slot,
		// so plan-mode chat reuses Super. Accept the field for API symmetry
		// with techai but route the value into Super if the user means to
		// change the plan-mode model.
		if patch.Plan != nil && strings.TrimSpace(*patch.Plan) != "" {
			s.cfg.Models.Super = strings.TrimSpace(*patch.Plan)
		}
		if err := config.Save(s.cfg); err != nil {
			http.Error(w, "save failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, configToDTO(s.cfg))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func configToDTO(cfg config.Config) configDTO {
	masked := cfg.API.APIKey
	if len(masked) > 8 {
		masked = masked[:4] + "…" + masked[len(masked)-3:]
	}
	return configDTO{
		BaseURL:  cfg.API.BaseURL,
		APIKey:   masked,
		Super:    cfg.Models.Super,
		Dev:      cfg.Models.Dev,
		Plan:     cfg.Models.Super, // Plan reuses Super in hanimo
		HasKey:   cfg.API.APIKey != "",
		BrandTag: "hanimo",
		IsOnprem: config.IsDistro(),
	}
}

type modelDTO struct {
	ID                string `json:"id"`
	DisplayName       string `json:"display_name"`
	Description       string `json:"description"`
	ContextWindow     int    `json:"context_window"`
	SupportsTools     bool   `json:"supports_tools"`
	SupportsVision    bool   `json:"supports_vision"`
	SupportsReasoning bool   `json:"supports_reasoning"`
}

func (s *server) handleModels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Static registry first — these have curated descriptions and capability
	// flags. Live cache (from /models/refresh) is merged on top so newly
	// shipped models the registry hasn't been updated for still appear.
	seen := make(map[string]int, len(llm.Models))
	out := make([]modelDTO, 0, len(llm.Models))
	for _, m := range llm.Models {
		cap := llm.GetCapability(m.ID)
		out = append(out, modelDTO{
			ID:                m.ID,
			DisplayName:       m.DisplayName,
			Description:       m.Description,
			ContextWindow:     cap.ContextWindow,
			SupportsTools:     cap.SupportsTools,
			SupportsVision:    supportsVision(m.ID),
			SupportsReasoning: supportsReasoning(m.ID),
		})
		seen[m.ID] = len(out) - 1
	}
	provider := detectProviderID(s.cfg.API.BaseURL)
	for _, m := range liveModelsForCurrentProvider(provider, s.cfg.API.BaseURL) {
		if i, ok := seen[m.ID]; ok {
			// Live entry wins for context window / tools when richer.
			if m.ContextWindow > out[i].ContextWindow {
				out[i].ContextWindow = m.ContextWindow
			}
			if m.SupportsTools {
				out[i].SupportsTools = true
			}
			continue
		}
		out = append(out, m)
		seen[m.ID] = len(out) - 1
	}
	sort.Slice(out, func(i, j int) bool { return out[i].DisplayName < out[j].DisplayName })
	writeJSON(w, http.StatusOK, map[string]any{"models": out})
}

// /index/symbols — Phase 2 indexing surface.
//   GET   → cached symbol map
//   POST  → force a rebuild; query ?force=1 also accepted on GET
func (s *server) handleIndexSymbols(w http.ResponseWriter, r *http.Request) {
	force := r.Method == http.MethodPost || r.URL.Query().Get("force") == "1"
	cwd, _ := os.Getwd()
	idx, err := loadOrBuildSymbolIndex(cwd, force)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, idx)
}

// KnowledgeDir is the per-workspace folder holding curated .md docs that
// are auto-prepended to the system prompt. The single-file .hanimo.md is
// also loaded when present.
const KnowledgeDir = ".hanimo-knowledge"

// ProjectContextFile is the single-file pinned project context.
const ProjectContextFile = ".hanimo.md"

// /knowledge — back-compat single-file endpoint for .hanimo.md.
//   GET    → {"path","content"}
//   PUT    → {"content"} replaces file
func (s *server) handleKnowledge(w http.ResponseWriter, r *http.Request) {
	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, ProjectContextFile)
	switch r.Method {
	case http.MethodGet:
		data, err := os.ReadFile(path)
		if err != nil && !os.IsNotExist(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"path": path, "content": string(data)})
	case http.MethodPut, http.MethodPost:
		var body struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if err := os.WriteFile(path, []byte(body.Content), 0644); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"path": path})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

type knowledgeFileMeta struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Modified int64  `json:"modified"`
	Enabled  bool   `json:"enabled"`
}

// /knowledge/files — list + read + write the .hanimo-knowledge/ folder.
//   GET                       → [{name,path,size,modified,enabled}, ...]
//   GET ?path=<rel>           → {path,content}
//   PUT  body {path,content}  → upsert
//   DELETE ?path=<rel>        → remove
func (s *server) handleKnowledgeFiles(w http.ResponseWriter, r *http.Request) {
	cwd, _ := os.Getwd()
	dir := filepath.Join(cwd, KnowledgeDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		http.Error(w, "mkdir: "+err.Error(), http.StatusInternalServerError)
		return
	}
	relParam := r.URL.Query().Get("path")
	switch r.Method {
	case http.MethodGet:
		if relParam != "" {
			abs, err := safeJoin(dir, relParam)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			data, err := os.ReadFile(abs)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, map[string]string{"path": abs, "content": string(data)})
			return
		}
		entries, _ := os.ReadDir(dir)
		out := make([]knowledgeFileMeta, 0, len(entries))
		for _, e := range entries {
			if e.IsDir() || !isMarkdown(e.Name()) {
				continue
			}
			info, err := e.Info()
			if err != nil {
				continue
			}
			out = append(out, knowledgeFileMeta{
				Name:     e.Name(),
				Path:     filepath.Join(dir, e.Name()),
				Size:     info.Size(),
				Modified: info.ModTime().Unix(),
				Enabled:  !strings.HasPrefix(e.Name(), "_"), // _name.md = disabled
			})
		}
		writeJSON(w, http.StatusOK, map[string]any{"dir": dir, "files": out})
	case http.MethodPut, http.MethodPost:
		var body struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		body.Path = strings.TrimSpace(body.Path)
		if body.Path == "" || !isMarkdown(body.Path) {
			http.Error(w, "path must be a .md filename", http.StatusBadRequest)
			return
		}
		abs, err := safeJoin(dir, body.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := os.WriteFile(abs, []byte(body.Content), 0644); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"path": abs})
	case http.MethodDelete:
		if relParam == "" {
			http.Error(w, "?path required", http.StatusBadRequest)
			return
		}
		abs, err := safeJoin(dir, relParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := os.Remove(abs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func isMarkdown(name string) bool {
	n := strings.ToLower(name)
	return strings.HasSuffix(n, ".md") || strings.HasSuffix(n, ".markdown")
}

// safeJoin prevents path traversal — the resolved path must stay within base.
func safeJoin(base, rel string) (string, error) {
	clean := filepath.Clean(rel)
	if filepath.IsAbs(clean) || strings.HasPrefix(clean, "..") {
		return "", fmt.Errorf("invalid path: %q", rel)
	}
	abs := filepath.Join(base, clean)
	absBase, _ := filepath.Abs(base)
	absResolved, _ := filepath.Abs(abs)
	if !strings.HasPrefix(absResolved, absBase+string(filepath.Separator)) && absResolved != absBase {
		return "", fmt.Errorf("path escapes base")
	}
	return abs, nil
}

// SkillsDir holds on-demand skill .md files. Each is shown as a
// (name + when_to_use) line in the system prompt; the LLM loads the body
// via existing file_read when relevant.
const SkillsDir = ".hanimo-skills"

// loadKnowledgeContext reads:
//   - User rules (~/.hanimo/rules.md) — always-on personal instructions
//   - .hanimo.md — project pinned context
//   - .hanimo-knowledge/*.md — enabled curated docs
//   - .hanimo-skills/*.md — listed by name+description only (lazy-loaded)
// Total size capped to avoid blowing the context window.
func loadKnowledgeContext(cwd string, maxBytes int) string {
	var b strings.Builder

	// 1) User rules
	if home, err := os.UserHomeDir(); err == nil {
		if data, err := os.ReadFile(filepath.Join(home, config.ConfigDirName, "rules.md")); err == nil {
			b.WriteString("\n\n## User Rules (~/")
			b.WriteString(config.ConfigDirName)
			b.WriteString("/rules.md)\n")
			b.Write(data)
		}
	}

	// 2) .hanimo.md
	if data, err := os.ReadFile(filepath.Join(cwd, ProjectContextFile)); err == nil {
		b.WriteString("\n\n## Project Context (")
		b.WriteString(ProjectContextFile)
		b.WriteString(")\n")
		b.Write(data)
	}

	// 3) .hanimo-knowledge/
	if entries, err := os.ReadDir(filepath.Join(cwd, KnowledgeDir)); err == nil {
		for _, e := range entries {
			if e.IsDir() || !isMarkdown(e.Name()) || strings.HasPrefix(e.Name(), "_") {
				continue
			}
			if maxBytes > 0 && b.Len() > maxBytes {
				b.WriteString("\n\n[... knowledge truncated to fit context ...]")
				return b.String()
			}
			data, err := os.ReadFile(filepath.Join(cwd, KnowledgeDir, e.Name()))
			if err != nil {
				continue
			}
			b.WriteString("\n\n## Knowledge: ")
			b.WriteString(e.Name())
			b.WriteString("\n")
			b.Write(data)
		}
	}

	// 4) .hanimo-skills/ — index only (name + when_to_use line)
	if entries, err := os.ReadDir(filepath.Join(cwd, SkillsDir)); err == nil && len(entries) > 0 {
		b.WriteString("\n\n## Available Skills (load via file_read when relevant)\n")
		for _, e := range entries {
			if e.IsDir() || !isMarkdown(e.Name()) || strings.HasPrefix(e.Name(), "_") {
				continue
			}
			path := filepath.Join(cwd, SkillsDir, e.Name())
			summary := readSkillSummary(path)
			b.WriteString(fmt.Sprintf("- `%s/%s` — %s\n", SkillsDir, e.Name(), summary))
		}
	}

	return b.String()
}

// readSkillSummary returns the first non-empty line of front-matter
// "description:" or the first heading line, otherwise the first line.
func readSkillSummary(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "(unreadable)"
	}
	lines := strings.Split(string(data), "\n")
	// YAML frontmatter
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		for i := 1; i < len(lines); i++ {
			ln := strings.TrimSpace(lines[i])
			if ln == "---" {
				break
			}
			if strings.HasPrefix(ln, "description:") {
				return strings.TrimSpace(strings.TrimPrefix(ln, "description:"))
			}
		}
	}
	for _, ln := range lines {
		t := strings.TrimSpace(ln)
		if t == "" {
			continue
		}
		t = strings.TrimLeft(t, "# ")
		if len(t) > 100 {
			t = t[:100] + "…"
		}
		return t
	}
	return "(empty)"
}

// /skills — list/read .hanimo-skills/ entries.
//   GET                       → {dir, skills:[{name,path,size,modified,description}]}
//   GET ?path=<rel>           → {path, content}
//   PUT body {path,content}   → upsert
//   DELETE ?path=<rel>        → remove
func (s *server) handleSkills(w http.ResponseWriter, r *http.Request) {
	cwd, _ := os.Getwd()
	dir := filepath.Join(cwd, SkillsDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rel := r.URL.Query().Get("path")
	switch r.Method {
	case http.MethodGet:
		if rel != "" {
			abs, err := safeJoin(dir, rel)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			data, err := os.ReadFile(abs)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, map[string]string{"path": abs, "content": string(data)})
			return
		}
		entries, _ := os.ReadDir(dir)
		out := make([]map[string]any, 0, len(entries))
		for _, e := range entries {
			if e.IsDir() || !isMarkdown(e.Name()) {
				continue
			}
			info, err := e.Info()
			if err != nil {
				continue
			}
			abs := filepath.Join(dir, e.Name())
			out = append(out, map[string]any{
				"name":        e.Name(),
				"path":        abs,
				"size":        info.Size(),
				"modified":    info.ModTime().Unix(),
				"enabled":     !strings.HasPrefix(e.Name(), "_"),
				"description": readSkillSummary(abs),
			})
		}
		writeJSON(w, http.StatusOK, map[string]any{"dir": dir, "skills": out})
	case http.MethodPut, http.MethodPost:
		var body struct{ Path, Content string }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if !isMarkdown(body.Path) {
			http.Error(w, "path must be .md", http.StatusBadRequest)
			return
		}
		abs, err := safeJoin(dir, body.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := os.WriteFile(abs, []byte(body.Content), 0644); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"path": abs})
	case http.MethodDelete:
		if rel == "" {
			http.Error(w, "?path required", http.StatusBadRequest)
			return
		}
		abs, err := safeJoin(dir, rel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = os.Remove(abs)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// /rules — read/write user-level persistent instructions at ~/.hanimo/rules.md.
func (s *server) handleRules(w http.ResponseWriter, r *http.Request) {
	home, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dir := filepath.Join(home, config.ConfigDirName)
	_ = os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, "rules.md")
	switch r.Method {
	case http.MethodGet:
		data, _ := os.ReadFile(path)
		writeJSON(w, http.StatusOK, map[string]string{"path": path, "content": string(data)})
	case http.MethodPut, http.MethodPost:
		var body struct{ Content string }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if err := os.WriteFile(path, []byte(body.Content), 0644); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"path": path})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

// supportsVision and supportsReasoning derive these flags by ID prefix
// because hanimo's ModelCapability struct does not carry them. The WebView
// uses these to render small badges next to each model in the dropdown.
func supportsVision(id string) bool {
	prefixes := []string{"gpt-4o", "claude-", "gemini-", "qwen-vl"}
	for _, p := range prefixes {
		if strings.HasPrefix(id, p) {
			return true
		}
	}
	return false
}

func supportsReasoning(id string) bool {
	prefixes := []string{"o1", "o3", "deepseek-reasoner", "claude-opus", "qwen3-thinking"}
	for _, p := range prefixes {
		if strings.HasPrefix(id, p) {
			return true
		}
	}
	return false
}
