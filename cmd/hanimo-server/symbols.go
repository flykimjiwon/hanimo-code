// Lightweight symbol indexer for the workspace. Regex-only — no tree-sitter,
// no embeddings. Trades 100% accuracy for speed (<1s on a 50k-line repo) and
// zero install footprint. The LLM uses find_symbol to jump to definitions
// quickly; for anything beyond exact-name lookups, it falls back to grep_search.
package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Symbol struct {
	Name string `json:"name"`
	Kind string `json:"kind"` // function | class | method | const | var | export
	File string `json:"file"`
	Line int    `json:"line"`
}

type SymbolIndex struct {
	BuiltAt   time.Time `json:"built_at"`
	Root      string    `json:"root"`
	FileCount int       `json:"file_count"`
	Symbols   []Symbol  `json:"symbols"`
}

var (
	indexMu    sync.Mutex
	cachedIdx  *SymbolIndex
	cachedRoot string
)

const cacheDir = ".hanimo-cache"
const cacheFile = "symbols.json"

// Languages we scan and the regex used per language. Each capture group #1
// is the symbol name, the leading keyword classifies the kind.
var langPatterns = map[string][]struct {
	kind string
	re   *regexp.Regexp
}{
	".go": {
		{"function", regexp.MustCompile(`^func\s+(?:\([^)]*\)\s+)?([A-Z_a-z]\w*)\s*\(`)},
		{"struct", regexp.MustCompile(`^type\s+([A-Z_a-z]\w*)\s+struct\b`)},
		{"interface", regexp.MustCompile(`^type\s+([A-Z_a-z]\w*)\s+interface\b`)},
		{"const", regexp.MustCompile(`^(?:var|const)\s+([A-Z][A-Z_a-z0-9]*)\s`)},
	},
	".ts":  tsLikePatterns,
	".tsx": tsLikePatterns,
	".js":  tsLikePatterns,
	".jsx": tsLikePatterns,
	".py": {
		{"function", regexp.MustCompile(`^def\s+([A-Z_a-z]\w*)\s*\(`)},
		{"class", regexp.MustCompile(`^class\s+([A-Z_a-z]\w*)\b`)},
	},
	".rs": {
		{"function", regexp.MustCompile(`^(?:pub(?:\([^)]*\))?\s+)?fn\s+([A-Z_a-z]\w*)\s*[(<]`)},
		{"struct", regexp.MustCompile(`^(?:pub(?:\([^)]*\))?\s+)?struct\s+([A-Z_a-z]\w*)\b`)},
		{"enum", regexp.MustCompile(`^(?:pub(?:\([^)]*\))?\s+)?enum\s+([A-Z_a-z]\w*)\b`)},
		{"trait", regexp.MustCompile(`^(?:pub(?:\([^)]*\))?\s+)?trait\s+([A-Z_a-z]\w*)\b`)},
	},
}

var tsLikePatterns = []struct {
	kind string
	re   *regexp.Regexp
}{
	{"function", regexp.MustCompile(`^(?:export\s+)?(?:async\s+)?function\s+([A-Z_a-z]\w*)\s*[(<]`)},
	{"class", regexp.MustCompile(`^(?:export\s+)?(?:abstract\s+)?class\s+([A-Z_a-z]\w*)\b`)},
	{"interface", regexp.MustCompile(`^(?:export\s+)?interface\s+([A-Z_a-z]\w*)\b`)},
	{"type", regexp.MustCompile(`^(?:export\s+)?type\s+([A-Z_a-z]\w*)\s*=`)},
	{"const", regexp.MustCompile(`^(?:export\s+)?const\s+([A-Z_a-z]\w*)\s*[=:]`)},
}

var skipDirs = map[string]bool{
	"node_modules": true, ".git": true, "dist": true, "build": true,
	".hanimo-cache": true, "vendor": true, "__pycache__": true, ".next": true,
	"target": true, ".venv": true, "venv": true,
}

// loadOrBuildSymbolIndex returns a cached index when available and fresh,
// otherwise rebuilds from disk. force=true bypasses the cache.
func loadOrBuildSymbolIndex(root string, force bool) (*SymbolIndex, error) {
	indexMu.Lock()
	defer indexMu.Unlock()

	if !force && cachedIdx != nil && cachedRoot == root && time.Since(cachedIdx.BuiltAt) < 5*time.Minute {
		return cachedIdx, nil
	}

	cachePath := filepath.Join(root, cacheDir, cacheFile)
	if !force {
		if data, err := os.ReadFile(cachePath); err == nil {
			var idx SymbolIndex
			if json.Unmarshal(data, &idx) == nil && time.Since(idx.BuiltAt) < 30*time.Minute {
				cachedIdx = &idx
				cachedRoot = root
				return &idx, nil
			}
		}
	}

	idx, err := buildSymbolIndex(root)
	if err != nil {
		return nil, err
	}
	cachedIdx = idx
	cachedRoot = root
	_ = os.MkdirAll(filepath.Dir(cachePath), 0755)
	if data, err := json.Marshal(idx); err == nil {
		_ = os.WriteFile(cachePath, data, 0644)
	}
	return idx, nil
}

func buildSymbolIndex(root string) (*SymbolIndex, error) {
	idx := &SymbolIndex{
		BuiltAt: time.Now(),
		Root:    root,
		Symbols: make([]Symbol, 0, 1024),
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if skipDirs[d.Name()] || strings.HasPrefix(d.Name(), ".") && d.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}
		ext := filepath.Ext(path)
		patterns, ok := langPatterns[ext]
		if !ok {
			return nil
		}
		// Skip very large files (>1MB) — likely generated/minified
		info, err := d.Info()
		if err != nil || info.Size() > 1<<20 {
			return nil
		}
		idx.FileCount++
		rel, _ := filepath.Rel(root, path)
		scanFile(path, rel, patterns, &idx.Symbols)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return idx, nil
}

func scanFile(absPath, relPath string, patterns []struct {
	kind string
	re   *regexp.Regexp
}, out *[]Symbol) {
	f, err := os.Open(absPath)
	if err != nil {
		return
	}
	defer f.Close()
	data, err := os.ReadFile(absPath)
	if err != nil {
		return
	}
	for i, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimLeft(line, " \t")
		for _, p := range patterns {
			if m := p.re.FindStringSubmatch(trimmed); m != nil {
				*out = append(*out, Symbol{
					Name: m[1], Kind: p.kind, File: relPath, Line: i + 1,
				})
				break
			}
		}
	}
}

// FindSymbols returns symbols whose name matches the query (case-insensitive
// substring). Used by the find_symbol tool.
func FindSymbols(root, query string, limit int) ([]Symbol, error) {
	idx, err := loadOrBuildSymbolIndex(root, false)
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 30
	}
	q := strings.ToLower(query)
	out := make([]Symbol, 0, limit)
	// Exact matches first.
	for _, s := range idx.Symbols {
		if strings.EqualFold(s.Name, query) {
			out = append(out, s)
			if len(out) >= limit {
				return out, nil
			}
		}
	}
	for _, s := range idx.Symbols {
		if strings.EqualFold(s.Name, query) {
			continue
		}
		if strings.Contains(strings.ToLower(s.Name), q) {
			out = append(out, s)
			if len(out) >= limit {
				return out, nil
			}
		}
	}
	if len(out) == 0 {
		return out, fmt.Errorf("no symbol matched %q", query)
	}
	return out, nil
}
