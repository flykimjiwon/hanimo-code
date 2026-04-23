package tools

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/flykimjiwon/hanimo/internal/config"
)

// Symbol represents a code symbol (function, class, method, etc.)
type Symbol struct {
	Name     string
	Kind     string // "function", "class", "method", "interface", "type", "const"
	File     string
	Line     int
	Preview  string
	ModTime  time.Time
}

// Language-specific symbol patterns
var symbolPatterns = map[string][]*regexp.Regexp{
	".go": {
		regexp.MustCompile(`^func\s+(\w+)`),                           // func Foo
		regexp.MustCompile(`^func\s+\([^)]+\)\s+(\w+)`),              // func (r *R) Foo
		regexp.MustCompile(`^type\s+(\w+)\s+(struct|interface)`),       // type Foo struct/interface
		regexp.MustCompile(`^const\s+(\w+)`),                           // const Foo
		regexp.MustCompile(`^var\s+(\w+)`),                             // var Foo
	},
	".js": {
		regexp.MustCompile(`^(?:export\s+)?function\s+(\w+)`),         // function foo
		regexp.MustCompile(`^(?:export\s+)?(?:const|let|var)\s+(\w+)\s*=\s*(?:async\s+)?\(`), // const foo = (
		regexp.MustCompile(`^(?:export\s+)?class\s+(\w+)`),            // class Foo
	},
	".ts": {
		regexp.MustCompile(`^(?:export\s+)?function\s+(\w+)`),
		regexp.MustCompile(`^(?:export\s+)?(?:const|let|var)\s+(\w+)\s*=\s*(?:async\s+)?\(`),
		regexp.MustCompile(`^(?:export\s+)?class\s+(\w+)`),
		regexp.MustCompile(`^(?:export\s+)?interface\s+(\w+)`),
		regexp.MustCompile(`^(?:export\s+)?type\s+(\w+)`),
	},
	".tsx": nil, // filled below
	".jsx": nil,
	".py": {
		regexp.MustCompile(`^(?:async\s+)?def\s+(\w+)`),               // def foo
		regexp.MustCompile(`^class\s+(\w+)`),                           // class Foo
	},
	".java": {
		regexp.MustCompile(`(?:public|private|protected)\s+(?:static\s+)?(?:\w+\s+)+(\w+)\s*\(`), // public void foo(
		regexp.MustCompile(`(?:public|private|protected)?\s*class\s+(\w+)`),                       // class Foo
		regexp.MustCompile(`(?:public|private|protected)?\s*interface\s+(\w+)`),                    // interface Foo
	},
	".rs": {
		regexp.MustCompile(`^(?:pub\s+)?fn\s+(\w+)`),                  // fn foo / pub fn foo
		regexp.MustCompile(`^(?:pub\s+)?struct\s+(\w+)`),              // struct Foo
		regexp.MustCompile(`^(?:pub\s+)?enum\s+(\w+)`),                // enum Foo
		regexp.MustCompile(`^(?:pub\s+)?trait\s+(\w+)`),               // trait Foo
		regexp.MustCompile(`^impl\s+(\w+)`),                            // impl Foo
	},
	".sh": {
		regexp.MustCompile(`^(\w+)\s*\(\)`),                            // foo()
		regexp.MustCompile(`^function\s+(\w+)`),                        // function foo
	},
}

func init() {
	// tsx/jsx share ts/js patterns
	symbolPatterns[".tsx"] = symbolPatterns[".ts"]
	symbolPatterns[".jsx"] = symbolPatterns[".js"]
}

// SymbolSearch finds code symbols (functions, classes, etc.) matching a query.
func SymbolSearch(query, basePath string) (string, error) {
	if query == "" {
		return "", fmt.Errorf("query is required")
	}
	if basePath == "" {
		basePath = "."
	}
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	gi := LoadGitIgnore(absBase)
	queryLower := strings.ToLower(query)

	var symbols []Symbol
	maxSymbols := 50

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

		ext := strings.ToLower(filepath.Ext(d.Name()))
		patterns, ok := symbolPatterns[ext]
		if !ok || len(patterns) == 0 {
			return nil
		}

		info, _ := d.Info()
		mtime := time.Time{}
		if info != nil {
			mtime = info.ModTime()
		}

		found := extractSymbols(path, rel, patterns, queryLower, mtime)
		symbols = append(symbols, found...)

		if len(symbols) >= maxSymbols*3 {
			return fmt.Errorf("limit reached")
		}
		return nil
	})

	if walkErr != nil && walkErr.Error() != "limit reached" {
		config.DebugLog("[SYMBOL] walk error: %v", walkErr)
	}

	if len(symbols) == 0 {
		return "No symbols found.", nil
	}

	// Sort: exact name match first, then mtime descending
	sort.Slice(symbols, func(i, j int) bool {
		iExact := strings.EqualFold(symbols[i].Name, query)
		jExact := strings.EqualFold(symbols[j].Name, query)
		if iExact != jExact {
			return iExact
		}
		return symbols[i].ModTime.After(symbols[j].ModTime)
	})

	if len(symbols) > maxSymbols {
		symbols = symbols[:maxSymbols]
	}

	var sb strings.Builder
	for _, s := range symbols {
		sb.WriteString(fmt.Sprintf("%s:%d [%s] %s\n", s.File, s.Line, s.Kind, s.Preview))
	}

	config.DebugLog("[SYMBOL] query=%q matches=%d", query, len(symbols))
	return sb.String(), nil
}

func extractSymbols(absPath, relPath string, patterns []*regexp.Regexp, queryLower string, mtime time.Time) []Symbol {
	f, err := os.Open(absPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var symbols []Symbol
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 256*1024)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		for _, pat := range patterns {
			matches := pat.FindStringSubmatch(trimmed)
			if len(matches) < 2 {
				continue
			}
			name := matches[1]
			if !strings.Contains(strings.ToLower(name), queryLower) {
				continue
			}

			kind := detectKind(trimmed)
			symbols = append(symbols, Symbol{
				Name:    name,
				Kind:    kind,
				File:    relPath,
				Line:    lineNum,
				Preview: truncateLine(trimmed),
				ModTime: mtime,
			})
			break
		}
	}
	return symbols
}

// FindReferences finds files that reference (call/use) a given symbol name.
// Uses SymbolSearch to find the definition, then GrepSearch for all usages,
// excluding the definition line itself.
func FindReferences(symbolName, basePath string) (string, error) {
	if symbolName == "" {
		return "", fmt.Errorf("symbolName is required")
	}
	if basePath == "" {
		basePath = "."
	}

	// Find definition location first
	defResult, err := SymbolSearch(symbolName, basePath)
	if err != nil {
		return "", fmt.Errorf("symbol search failed: %w", err)
	}

	// Parse definition file:line from first exact-match result
	defFile := ""
	defLine := 0
	if defResult != "No symbols found." {
		for _, line := range strings.Split(defResult, "\n") {
			if line == "" {
				continue
			}
			// Format: "file:line [kind] preview"
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				defFile = parts[0]
				fmt.Sscanf(parts[1], "%d", &defLine)
				break
			}
		}
	}

	// Grep for all occurrences using word-boundary pattern
	pattern := `\b` + regexp.QuoteMeta(symbolName) + `\b`
	raw, err := GrepSearch(pattern, basePath, "", false, 0)
	if err != nil {
		return "", fmt.Errorf("grep failed: %w", err)
	}

	if raw == "" || strings.HasPrefix(raw, "No matches") {
		return fmt.Sprintf("No references found for %q", symbolName), nil
	}

	// Group by file, exclude definition line
	fileLines := map[string][]string{}
	var fileOrder []string
	seen := map[string]bool{}

	for _, line := range strings.Split(raw, "\n") {
		if line == "" {
			continue
		}
		// format: file:line:content
		parts := strings.SplitN(line, ":", 3)
		if len(parts) < 3 {
			continue
		}
		f, lnStr := parts[0], parts[1]
		var ln int
		fmt.Sscanf(lnStr, "%d", &ln)

		// Exclude definition line
		if f == defFile && ln == defLine {
			continue
		}

		if !seen[f] {
			seen[f] = true
			fileOrder = append(fileOrder, f)
		}
		fileLines[f] = append(fileLines[f], fmt.Sprintf("  L%s: %s", lnStr, strings.TrimSpace(parts[2])))
	}

	sort.Strings(fileOrder)

	if len(fileOrder) == 0 {
		return fmt.Sprintf("No references found for %q (only definition)", symbolName), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("References to %q:\n", symbolName))
	if defFile != "" {
		sb.WriteString(fmt.Sprintf("  Definition: %s:%d\n\n", defFile, defLine))
	}
	for _, f := range fileOrder {
		sb.WriteString(fmt.Sprintf("%s:\n", f))
		for _, l := range fileLines[f] {
			sb.WriteString(l + "\n")
		}
	}
	return strings.TrimRight(sb.String(), "\n"), nil
}

func detectKind(line string) string {
	lower := strings.TrimSpace(strings.ToLower(line))
	switch {
	case strings.HasPrefix(lower, "class "), strings.Contains(lower, " class "):
		return "class"
	case strings.HasPrefix(lower, "interface "), strings.Contains(lower, " interface "):
		return "interface"
	case strings.HasPrefix(lower, "type ") && (strings.Contains(lower, "struct") || strings.Contains(lower, "interface")):
		return "type"
	case strings.HasPrefix(lower, "trait "):
		return "trait"
	case strings.HasPrefix(lower, "enum "), strings.Contains(lower, " enum "):
		return "enum"
	case strings.HasPrefix(lower, "const "):
		return "const"
	case strings.HasPrefix(lower, "impl "):
		return "impl"
	default:
		return "function"
	}
}
