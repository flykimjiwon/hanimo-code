package tools

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var todoPattern = regexp.MustCompile(`(?i)(TODO|FIXME|HACK)[\s:!]`)

// CodeQualityScan runs rule-based quality checks on source files.
// checks is a comma-separated list of: "todo", "large_functions", "deep_nesting", or "all".
func CodeQualityScan(basePath string, checks string) string {
	if basePath == "" {
		basePath = "."
	}
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return fmt.Sprintf("Error: invalid path: %v", err)
	}

	if checks == "" {
		checks = "all"
	}
	active := map[string]bool{}
	if checks == "all" {
		active["todo"] = true
		active["large_functions"] = true
		active["deep_nesting"] = true
	} else {
		for _, c := range strings.Split(checks, ",") {
			active[strings.TrimSpace(c)] = true
		}
	}

	gi := LoadGitIgnore(absBase)

	var todoFindings []string
	var largeFnFindings []string
	var nestingFindings []string

	sourceExts := map[string]bool{
		".go": true, ".js": true, ".ts": true, ".tsx": true,
		".jsx": true, ".py": true, ".java": true, ".rs": true,
	}

	_ = filepath.WalkDir(absBase, func(path string, d os.DirEntry, err error) error {
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

		ext := strings.ToLower(filepath.Ext(d.Name()))
		if !sourceExts[ext] {
			return nil
		}

		if active["todo"] {
			todoFindings = append(todoFindings, scanTODO(path, rel)...)
		}
		if active["large_functions"] {
			largeFnFindings = append(largeFnFindings, scanLargeFunctions(path, rel, ext)...)
		}
		if active["deep_nesting"] {
			nestingFindings = append(nestingFindings, scanDeepNesting(path, rel)...)
		}
		return nil
	})

	var sb strings.Builder

	if active["todo"] {
		sb.WriteString("=== TODO/FIXME/HACK ===\n")
		if len(todoFindings) == 0 {
			sb.WriteString("  (none found)\n")
		} else {
			for _, f := range todoFindings {
				sb.WriteString("  " + f + "\n")
			}
		}
		sb.WriteString("\n")
	}

	if active["large_functions"] {
		sb.WriteString("=== Large Functions (>50 lines) ===\n")
		if len(largeFnFindings) == 0 {
			sb.WriteString("  (none found)\n")
		} else {
			for _, f := range largeFnFindings {
				sb.WriteString("  " + f + "\n")
			}
		}
		sb.WriteString("\n")
	}

	if active["deep_nesting"] {
		sb.WriteString("=== Deep Nesting (>5 levels) ===\n")
		if len(nestingFindings) == 0 {
			sb.WriteString("  (none found)\n")
		} else {
			for _, f := range nestingFindings {
				sb.WriteString("  " + f + "\n")
			}
		}
		sb.WriteString("\n")
	}

	result := strings.TrimRight(sb.String(), "\n")
	if result == "" {
		return "No quality issues found."
	}
	return result
}

func scanTODO(absPath, relPath string) []string {
	f, err := os.Open(absPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var findings []string
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if todoPattern.MatchString(line) {
			findings = append(findings, fmt.Sprintf("%s:%d: %s", relPath, lineNum, strings.TrimSpace(line)))
		}
	}
	return findings
}

func scanLargeFunctions(absPath, relPath, ext string) []string {
	patterns, ok := symbolPatterns[ext]
	if !ok || len(patterns) == 0 {
		return nil
	}

	f, err := os.Open(absPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	type fnStart struct {
		name string
		line int
	}

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var starts []fnStart
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		for _, pat := range patterns {
			m := pat.FindStringSubmatch(trimmed)
			if len(m) >= 2 {
				starts = append(starts, fnStart{name: m[1], line: i + 1})
				break
			}
		}
	}

	var findings []string
	for i, fn := range starts {
		endLine := len(lines)
		if i+1 < len(starts) {
			endLine = starts[i+1].line - 1
		}
		size := endLine - fn.line + 1
		if size > 50 {
			findings = append(findings, fmt.Sprintf("%s:%d: %s (%d lines)", relPath, fn.line, fn.name, size))
		}
	}
	return findings
}

func scanDeepNesting(absPath, relPath string) []string {
	f, err := os.Open(absPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var findings []string
	scanner := bufio.NewScanner(f)
	lineNum := 0
	inDeep := false
	deepStart := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Count indentation in spaces (tabs count as 4)
		indent := 0
		for _, ch := range line {
			if ch == ' ' {
				indent++
			} else if ch == '\t' {
				indent += 4
			} else {
				break
			}
		}
		level := indent / 4

		if level > 5 {
			if !inDeep {
				inDeep = true
				deepStart = lineNum
			}
		} else {
			if inDeep {
				findings = append(findings, fmt.Sprintf("%s:%d-%d: deep nesting (>5 levels)", relPath, deepStart, lineNum-1))
				inDeep = false
			}
		}
	}
	if inDeep {
		findings = append(findings, fmt.Sprintf("%s:%d-%d: deep nesting (>5 levels)", relPath, deepStart, lineNum))
	}
	return findings
}
