package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/flykimjiwon/hanimo/internal/config"
)

// fuzzyScore computes a score for how well query matches candidate.
// Returns -1 if query is not a subsequence of candidate.
func fuzzyScore(query, candidate string) int {
	lowerQuery := strings.ToLower(query)
	lowerCand := strings.ToLower(candidate)

	// Subsequence check
	qi := 0
	for ci := 0; ci < len(lowerCand) && qi < len(lowerQuery); ci++ {
		if lowerCand[ci] == lowerQuery[qi] {
			qi++
		}
	}
	if qi < len(lowerQuery) {
		return -1 // not a subsequence
	}

	score := 0

	// Consecutive match bonus
	qi = 0
	consecutive := 0
	for ci := 0; ci < len(lowerCand) && qi < len(lowerQuery); ci++ {
		if lowerCand[ci] == lowerQuery[qi] {
			consecutive++
			if consecutive > 1 {
				score += 3 // consecutive bonus
			}
			// Exact case match bonus
			if qi < len(query) && ci < len(candidate) && candidate[ci] == query[qi] {
				score++
			}
			qi++
		} else {
			consecutive = 0
		}
	}

	// Basename match bonus
	base := filepath.Base(candidate)
	lowerBase := strings.ToLower(base)
	bqi := 0
	for ci := 0; ci < len(lowerBase) && bqi < len(lowerQuery); ci++ {
		if lowerBase[ci] == lowerQuery[bqi] {
			bqi++
		}
	}
	if bqi == len(lowerQuery) {
		score += 5
	}

	return score
}

// FuzzyFileSearch finds files matching a fuzzy query string.
// Uses character-by-character subsequence matching with scoring.
// Returns results formatted as "score\tpath\n".
func FuzzyFileSearch(query, basePath string, maxResults int) (string, error) {
	if query == "" {
		return "", fmt.Errorf("query is required")
	}
	if basePath == "" {
		basePath = "."
	}
	if maxResults <= 0 {
		maxResults = 20
	}

	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	gi := LoadGitIgnore(absBase)
	queryStartsDot := strings.HasPrefix(query, ".")

	type scored struct {
		path  string
		score int
	}
	var candidates []scored

	walkErr := filepath.WalkDir(absBase, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		rel, _ := filepath.Rel(absBase, path)
		rel = filepath.ToSlash(rel)

		if d.IsDir() {
			if shouldSkip(gi, rel, true, d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if shouldSkip(gi, rel, false, d.Name()) {
			return nil
		}

		s := fuzzyScore(query, rel)
		if s < 0 {
			return nil
		}

		// Hidden files rank lower unless query starts with '.'
		if !queryStartsDot && strings.HasPrefix(d.Name(), ".") {
			s -= 10
		}

		candidates = append(candidates, scored{path: rel, score: s})
		return nil
	})

	if walkErr != nil {
		config.DebugLog("[FUZZY] walk error: %v", walkErr)
	}

	// Sort by score descending, then path ascending for ties
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].score != candidates[j].score {
			return candidates[i].score > candidates[j].score
		}
		return candidates[i].path < candidates[j].path
	})

	if len(candidates) > maxResults {
		candidates = candidates[:maxResults]
	}

	if len(candidates) == 0 {
		return "No files matched.", nil
	}

	var sb strings.Builder
	for _, c := range candidates {
		sb.WriteString(fmt.Sprintf("%d\t%s\n", c.score, c.path))
	}

	config.DebugLog("[FUZZY] query=%q path=%s results=%d", query, basePath, len(candidates))
	return sb.String(), nil
}
