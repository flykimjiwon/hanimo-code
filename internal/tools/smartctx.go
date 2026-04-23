package tools

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// SmartContext returns the top N most relevant files for a given target file.
// Uses import graph + git hot files to rank relevance.
func SmartContext(targetFile, basePath string, maxFiles int) (string, error) {
	if targetFile == "" {
		return "", fmt.Errorf("target file is required")
	}
	if basePath == "" {
		basePath = "."
	}
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}
	if maxFiles <= 0 {
		maxFiles = 5
	}

	// Normalize target to forward-slash relative path
	target := filepath.ToSlash(targetFile)

	// 1. Build import graph
	graph := ImportGraph(absBase)

	// 2. Score each file
	scores := map[string]int{}

	// Files that targetFile imports (direct dependencies)
	for _, imp := range graph[target] {
		scores[imp] += 10
	}

	// Files that import targetFile (reverse dependencies)
	importers := ReverseImports(graph, target)
	for _, imp := range importers {
		scores[imp] += 8
	}

	// 3. Git hot files for recency signal
	hotRaw, err := GitHotFiles(absBase, 14)
	if err == nil {
		for _, line := range strings.Split(hotRaw, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// Format from GitHotFiles: "  N  path/to/file"
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				hotFile := parts[len(parts)-1]
				if hotFile != target {
					scores[hotFile] += 5
				}
			}
		}
	}

	// 4. Shared directory bonus
	targetDir := filepath.Dir(target)
	for file := range graph {
		if file == target {
			continue
		}
		if filepath.Dir(file) == targetDir {
			scores[file] += 2
		}
	}

	// Remove the target file itself
	delete(scores, target)

	if len(scores) == 0 {
		return "No related files found.", nil
	}

	// 5. Sort by score descending
	type scored struct {
		file  string
		score int
	}
	var ranked []scored
	for f, s := range scores {
		ranked = append(ranked, scored{f, s})
	}
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].score != ranked[j].score {
			return ranked[i].score > ranked[j].score
		}
		return ranked[i].file < ranked[j].file
	})

	if len(ranked) > maxFiles {
		ranked = ranked[:maxFiles]
	}

	// 6. Format output with first 30 lines of each file
	var sb strings.Builder
	sb.WriteString("## Related Files\n")
	for _, r := range ranked {
		sb.WriteString(fmt.Sprintf("### %s (score: %d)\n", r.file, r.score))
		absFile := filepath.Join(absBase, r.file)
		lines := readFirstLines(absFile, 30)
		if lines != "" {
			sb.WriteString("```\n")
			sb.WriteString(lines)
			sb.WriteString("```\n")
		}
	}
	return sb.String(), nil
}

// readFirstLines reads the first n lines of a file.
func readFirstLines(path string, n int) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	var sb strings.Builder
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() && count < n {
		sb.WriteString(scanner.Text())
		sb.WriteString("\n")
		count++
	}
	return sb.String()
}
