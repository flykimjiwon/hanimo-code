package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/flykimjiwon/hanimo/internal/config"
)

// readSet is the per-turn set of absolute paths the agent has file_read'd
// in the current conversation turn. Edits and overwrites require a prior
// read so the LLM can't hallucinate an old_string or blow away an unread
// file. Reset by ResetReadSet() at the start of every sendMessage.
var (
	readSetMu sync.RWMutex
	readSet   = map[string]bool{}
)

// ResetReadSet clears the read-before-write tracking set. Call at the
// start of each user turn (sendMessage) so stale reads from previous
// turns don't unlock unrelated edits.
func ResetReadSet() {
	readSetMu.Lock()
	readSet = map[string]bool{}
	readSetMu.Unlock()
}

func markRead(absPath string) {
	readSetMu.Lock()
	readSet[absPath] = true
	readSetMu.Unlock()
}

func wasRead(absPath string) bool {
	readSetMu.RLock()
	defer readSetMu.RUnlock()
	return readSet[absPath]
}

func FileRead(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("read failed: %w", err)
	}
	markRead(absPath)
	return string(data), nil
}

// FileWrite creates a new file or overwrites an existing one. Overwriting
// an existing file requires that it was file_read'd first in the current
// turn — this prevents the LLM from blowing away content it never saw.
// Creating a genuinely new file is always allowed.
func FileWrite(path, content string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	if warn := CheckSensitiveFile(absPath); warn != "" {
		return fmt.Errorf("%s — writing blocked", warn)
	}
	if warn := CheckSecrets(content); warn != "" {
		config.DebugLog("[FILE-WRITE] %s", warn)
	}
	if _, statErr := os.Stat(absPath); statErr == nil {
		// File exists — this is an overwrite.
		if !wasRead(absPath) {
			return fmt.Errorf("cannot overwrite %s without reading it first. Use file_read before file_write", path)
		}
	}
	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir failed: %w", err)
	}
	CreateSnapshot(absPath)
	if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
		return err
	}
	// A fresh write counts as a read for subsequent edits in the same turn.
	markRead(absPath)
	return nil
}

// FileEdit performs a search-and-replace edit on a file with 4-stage fuzzy matching.
// The file must have been file_read'd earlier in the turn — enforcing read-before-write
// eliminates the most common LLM failure mode where old_string is recalled
// from an outdated memory and silently corrupts the file.
//
// Stage 1: ExactMatch — exact string match (fastest)
// Stage 2: LineTrimmed — ignore leading/trailing whitespace per line
// Stage 3: IndentFlex — normalize all indentation differences
// Stage 4: Levenshtein — accept 95%+ similarity match
//
// Returns the number of replacements made, a unified diff preview, and any error.
func FileEdit(path, oldStr, newStr string) (int, string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return 0, "", fmt.Errorf("invalid path: %w", err)
	}
	if !wasRead(absPath) {
		return 0, "", fmt.Errorf("cannot edit %s without reading it first. Use file_read to verify old_string before file_edit", path)
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		return 0, "", fmt.Errorf("read failed: %w", err)
	}
	content := string(data)

	applyEdit := func(original, matched, replacement, stage string) (int, string, error) {
		newContent := strings.Replace(original, matched, replacement, 1)
		diff := GenerateUnifiedDiff(path, original, newContent)
		if err := SnapshotAndWrite(absPath, []byte(newContent)); err != nil {
			return 0, "", fmt.Errorf("write failed: %w", err)
		}
		config.DebugLog("[FILE-EDIT] stage=%s path=%s", stage, path)
		return 1, diff, nil
	}

	// Stage 1: ExactMatch
	if strings.Contains(content, oldStr) {
		return applyEdit(content, oldStr, newStr, "ExactMatch")
	}

	// Stage 2: LineTrimmed — match ignoring leading/trailing whitespace per line
	if _, matchedOld := lineTrimmedFind(content, oldStr); matchedOld != "" {
		return applyEdit(content, matchedOld, newStr, "LineTrimmed")
	}

	// Stage 3: IndentFlex — normalize all indentation
	if _, matchedOld := indentFlexFind(content, oldStr); matchedOld != "" {
		return applyEdit(content, matchedOld, newStr, "IndentFlex")
	}

	// Stage 4: Levenshtein — 95%+ similarity (high threshold to avoid wrong-block matches)
	if matchedOld, similarity := levenshteinFind(content, oldStr, 0.95); matchedOld != "" {
		n, diff, err := applyEdit(content, matchedOld, newStr, fmt.Sprintf("Levenshtein(%.1f%%)", similarity*100))
		return n, diff, err
	}

	// All stages failed
	preview := content
	if pr := []rune(preview); len(pr) > 500 {
		preview = string(pr[:500]) + "..."
	}
	return 0, "", fmt.Errorf("old_string not found in %s (tried 4 fuzzy stages). File preview:\n%s", path, preview)
}

// lineTrimmedFind finds oldStr in content by comparing lines with trimmed whitespace.
// Returns the start index and the actual matched substring from content.
func lineTrimmedFind(content, oldStr string) (int, string) {
	contentLines := strings.Split(content, "\n")
	searchLines := strings.Split(oldStr, "\n")
	if len(searchLines) == 0 {
		return -1, ""
	}

	trimmedSearch := make([]string, len(searchLines))
	for i, l := range searchLines {
		trimmedSearch[i] = strings.TrimSpace(l)
	}

	for i := 0; i <= len(contentLines)-len(searchLines); i++ {
		match := true
		for j, ts := range trimmedSearch {
			if strings.TrimSpace(contentLines[i+j]) != ts {
				match = false
				break
			}
		}
		if match {
			matched := strings.Join(contentLines[i:i+len(searchLines)], "\n")
			idx := strings.Index(content, matched)
			return idx, matched
		}
	}
	return -1, ""
}

// indentFlexFind normalizes all whitespace at line starts before comparing.
func indentFlexFind(content, oldStr string) (int, string) {
	contentLines := strings.Split(content, "\n")
	searchLines := strings.Split(oldStr, "\n")
	if len(searchLines) == 0 {
		return -1, ""
	}

	normalizeIndent := func(s string) string {
		trimmed := strings.TrimLeft(s, " \t")
		if len(trimmed) < len(s) {
			return " " + trimmed
		}
		return s
	}

	normSearch := make([]string, len(searchLines))
	for i, l := range searchLines {
		normSearch[i] = normalizeIndent(l)
	}

	for i := 0; i <= len(contentLines)-len(searchLines); i++ {
		match := true
		for j, ns := range normSearch {
			if normalizeIndent(contentLines[i+j]) != ns {
				match = false
				break
			}
		}
		if match {
			matched := strings.Join(contentLines[i:i+len(searchLines)], "\n")
			idx := strings.Index(content, matched)
			return idx, matched
		}
	}
	return -1, ""
}

// levenshteinFind searches content for a block similar to oldStr using Levenshtein distance.
// Only considers blocks of similar line count. Returns matched text and similarity ratio.
func levenshteinFind(content, oldStr string, threshold float64) (string, float64) {
	contentLines := strings.Split(content, "\n")
	searchLines := strings.Split(oldStr, "\n")
	searchLen := len(searchLines)
	if searchLen == 0 || searchLen > len(contentLines) {
		return "", 0
	}

	// For very large files/searches, skip to avoid O(n*m) explosion
	if utf8.RuneCountInString(oldStr) > 2000 {
		return "", 0
	}

	bestSimilarity := 0.0
	bestMatch := ""

	for i := 0; i <= len(contentLines)-searchLen; i++ {
		candidate := strings.Join(contentLines[i:i+searchLen], "\n")
		sim := levenshteinSimilarity(candidate, oldStr)
		if sim > bestSimilarity {
			bestSimilarity = sim
			bestMatch = candidate
		}
		if sim > 0.98 {
			break
		}
	}

	if bestSimilarity >= threshold {
		return bestMatch, bestSimilarity
	}
	return "", 0
}

// levenshteinSimilarity returns similarity ratio (0.0 to 1.0) between two strings.
func levenshteinSimilarity(a, b string) float64 {
	ra := []rune(a)
	rb := []rune(b)
	la, lb := len(ra), len(rb)
	if la == 0 && lb == 0 {
		return 1.0
	}
	maxLen := la
	if lb > maxLen {
		maxLen = lb
	}

	diff := la - lb
	if diff < 0 {
		diff = -diff
	}
	if float64(diff)/float64(maxLen) > 0.3 {
		return 0.0
	}

	prev := make([]int, lb+1)
	curr := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		prev[j] = j
	}
	for i := 1; i <= la; i++ {
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			del := prev[j] + 1
			ins := curr[j-1] + 1
			sub := prev[j-1] + cost
			min := del
			if ins < min {
				min = ins
			}
			if sub < min {
				min = sub
			}
			curr[j] = min
		}
		prev, curr = curr, prev
	}
	dist := prev[lb]
	return 1.0 - float64(dist)/float64(maxLen)
}

// defaultSkipDirs is the set of directory names that recursive ListFiles
// should skip entirely. These are dependency caches, build artifacts, VCS
// metadata, IDE folders, and hanimo-specific legacy paths that pollute
// recursive listings and hit the 500-file cap without yielding useful info.
var defaultSkipDirs = map[string]bool{
	".git":          true,
	".svn":          true,
	".hg":           true,
	"node_modules":  true,
	"dist":          true,
	"build":         true,
	"__pycache__":   true,
	".next":         true,
	".nuxt":         true,
	".svelte-kit":   true,
	"vendor":        true,
	".venv":         true,
	"venv":          true,
	"target":        true, // Rust
	".gradle":       true,
	".idea":         true,
	".vscode":       true,
	".omc":          true,
	".superpowers":  true,
	"_legacy_ts":    true, // hanimo-specific
	".DS_Store":     true,
	".pytest_cache": true,
	".mypy_cache":   true,
	".ruff_cache":   true,
	"coverage":      true,
}

// ListFiles lists files in a directory, optionally recursive.
func ListFiles(dir string, recursive bool) ([]string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}
	info, err := os.Stat(absDir)
	if err != nil {
		return nil, fmt.Errorf("stat failed: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dir)
	}

	var files []string
	truncated := false
	if recursive {
		err = filepath.WalkDir(absDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil // skip errors
			}
			name := d.Name()
			// Always skip known-noise directories entirely (don't walk into
			// them). Also skip any other hidden dir not explicitly listed.
			if d.IsDir() && path != absDir {
				if defaultSkipDirs[name] || strings.HasPrefix(name, ".") {
					return filepath.SkipDir
				}
			}
			rel, _ := filepath.Rel(absDir, path)
			if d.IsDir() {
				files = append(files, rel+"/")
			} else {
				files = append(files, rel)
			}
			if len(files) > 500 {
				truncated = true
				return filepath.SkipAll
			}
			return nil
		})
		if truncated {
			files = append(files, "", "[Hint] Too many files. Retry with recursive=false, or specify a subdirectory.")
		}
	} else {
		entries, err2 := os.ReadDir(absDir)
		if err2 != nil {
			return nil, fmt.Errorf("readdir failed: %w", err2)
		}
		for _, e := range entries {
			if e.IsDir() {
				files = append(files, e.Name()+"/")
			} else {
				files = append(files, e.Name())
			}
		}
	}
	return files, err
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ListTree renders a directory-only tree (no files) up to maxDepth. Skips the
// same noise directories as ListFiles. Gives the LLM a fast, token-cheap view
// of project structure without enumerating every file.
func ListTree(dir string, maxDepth int) (string, error) {
	if maxDepth <= 0 {
		maxDepth = 3
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}
	info, err := os.Stat(absDir)
	if err != nil {
		return "", fmt.Errorf("stat failed: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", dir)
	}

	var sb strings.Builder
	sb.WriteString(filepath.Base(absDir) + "/\n")

	var walk func(path string, prefix string, depth int) error
	walk = func(path string, prefix string, depth int) error {
		if depth > maxDepth {
			return nil
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil
		}
		dirs := make([]os.DirEntry, 0, len(entries))
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			name := e.Name()
			if defaultSkipDirs[name] || strings.HasPrefix(name, ".") {
				continue
			}
			dirs = append(dirs, e)
		}
		for i, e := range dirs {
			last := i == len(dirs)-1
			branch := "├── "
			nextPrefix := prefix + "│   "
			if last {
				branch = "└── "
				nextPrefix = prefix + "    "
			}
			sb.WriteString(prefix + branch + e.Name() + "/\n")
			_ = walk(filepath.Join(path, e.Name()), nextPrefix, depth+1)
		}
		return nil
	}
	_ = walk(absDir, "", 1)
	return sb.String(), nil
}
