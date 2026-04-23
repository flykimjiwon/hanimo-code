package tools

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Import regex patterns per language.
var (
	// Go
	goImportSingle = regexp.MustCompile(`^\s*import\s+"([^"]+)"`)
	goImportBlock  = regexp.MustCompile(`^\s*"([^"]+)"`)
	goImportStart  = regexp.MustCompile(`^\s*import\s*\(`)
	goImportEnd    = regexp.MustCompile(`^\s*\)`)

	// JS/TS
	jsImportFrom  = regexp.MustCompile(`(?:import|export)\s+.*\s+from\s+['"]([^'"]+)['"]`)
	jsImportPlain = regexp.MustCompile(`import\s+['"]([^'"]+)['"]`)
	jsRequire     = regexp.MustCompile(`require\s*\(\s*['"]([^'"]+)['"]\s*\)`)

	// Python
	pyImport     = regexp.MustCompile(`^import\s+(\S+)`)
	pyFromImport = regexp.MustCompile(`^from\s+(\S+)\s+import`)
)

// sourceExtensions maps file extensions to language identifiers.
var sourceExtensions = map[string]string{
	".go":  "go",
	".js":  "js",
	".ts":  "js",
	".jsx": "js",
	".tsx": "js",
	".py":  "py",
}

// ImportGraph scans source files and builds a map of file -> imported files.
// Supports Go, JS/TS, Python.
func ImportGraph(basePath string) map[string][]string {
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return nil
	}

	gi := LoadGitIgnore(absBase)
	goModule := readGoModule(absBase)

	// Collect all source files first (for resolution lookups).
	sourceFiles := map[string]bool{}
	_ = filepath.Walk(absBase, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(absBase, path)
		name := info.Name()

		if info.IsDir() {
			if shouldSkip(gi, rel, true, name) {
				return filepath.SkipDir
			}
			return nil
		}
		if shouldSkip(gi, rel, false, name) {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(name))
		if _, ok := sourceExtensions[ext]; ok {
			sourceFiles[rel] = true
		}
		return nil
	})

	graph := map[string][]string{}

	for relFile := range sourceFiles {
		absFile := filepath.Join(absBase, relFile)
		ext := strings.ToLower(filepath.Ext(relFile))
		lang := sourceExtensions[ext]

		rawImports := extractImports(absFile, lang)
		resolved := resolveImports(absBase, relFile, rawImports, lang, goModule, sourceFiles)

		if len(resolved) > 0 {
			sort.Strings(resolved)
			graph[relFile] = resolved
		}
	}

	return graph
}

// readGoModule extracts the module path from go.mod if present.
func readGoModule(basePath string) string {
	data, err := os.ReadFile(filepath.Join(basePath, "go.mod"))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}
	return ""
}

// extractImports reads a source file and returns raw import paths.
func extractImports(filePath string, lang string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var imports []string
	scanner := bufio.NewScanner(f)

	switch lang {
	case "go":
		inBlock := false
		for scanner.Scan() {
			line := scanner.Text()
			if goImportStart.MatchString(line) {
				inBlock = true
				continue
			}
			if inBlock {
				if goImportEnd.MatchString(line) {
					inBlock = false
					continue
				}
				if m := goImportBlock.FindStringSubmatch(line); m != nil {
					imports = append(imports, m[1])
				}
				continue
			}
			if m := goImportSingle.FindStringSubmatch(line); m != nil {
				imports = append(imports, m[1])
			}
		}

	case "js":
		for scanner.Scan() {
			line := scanner.Text()
			if m := jsImportFrom.FindStringSubmatch(line); m != nil {
				imports = append(imports, m[1])
			} else if m := jsImportPlain.FindStringSubmatch(line); m != nil {
				imports = append(imports, m[1])
			}
			// require can appear anywhere
			for _, m := range jsRequire.FindAllStringSubmatch(line, -1) {
				imports = append(imports, m[1])
			}
		}

	case "py":
		for scanner.Scan() {
			line := scanner.Text()
			if m := pyFromImport.FindStringSubmatch(line); m != nil {
				imports = append(imports, m[1])
			} else if m := pyImport.FindStringSubmatch(line); m != nil {
				imports = append(imports, m[1])
			}
		}
	}

	return imports
}

// resolveImports converts raw import paths to relative file paths within the project.
func resolveImports(basePath, sourceFile string, rawImports []string, lang, goModule string, sourceFiles map[string]bool) []string {
	sourceDir := filepath.Dir(sourceFile)
	var resolved []string
	seen := map[string]bool{}

	for _, imp := range rawImports {
		var candidates []string

		switch lang {
		case "go":
			candidates = resolveGoImport(imp, goModule, sourceFiles)
		case "js":
			candidates = resolveJSImport(imp, sourceDir, sourceFiles)
		case "py":
			candidates = resolvePyImport(imp, sourceDir, sourceFiles)
		}

		for _, c := range candidates {
			if c != sourceFile && !seen[c] {
				seen[c] = true
				resolved = append(resolved, c)
			}
		}
	}
	return resolved
}

// resolveGoImport resolves a Go import path to local files via go.mod module path.
func resolveGoImport(importPath, goModule string, sourceFiles map[string]bool) []string {
	if goModule == "" {
		return nil
	}
	if !strings.HasPrefix(importPath, goModule) {
		return nil // external package
	}

	// Strip module prefix to get relative package dir
	rel := strings.TrimPrefix(importPath, goModule)
	rel = strings.TrimPrefix(rel, "/")

	var matches []string
	for sf := range sourceFiles {
		if strings.ToLower(filepath.Ext(sf)) != ".go" {
			continue
		}
		dir := filepath.ToSlash(filepath.Dir(sf))
		if dir == rel || (rel == "" && dir == ".") {
			matches = append(matches, sf)
		}
	}
	return matches
}

// resolveJSImport resolves a JS/TS import to local files.
func resolveJSImport(importPath, sourceDir string, sourceFiles map[string]bool) []string {
	// Only resolve relative imports
	if !strings.HasPrefix(importPath, ".") {
		return nil
	}

	target := filepath.ToSlash(filepath.Join(sourceDir, importPath))
	target = filepath.Clean(target)

	// Try exact match, then with extensions, then as index
	jsExts := []string{"", ".js", ".ts", ".jsx", ".tsx", "/index.js", "/index.ts", "/index.jsx", "/index.tsx"}
	for _, ext := range jsExts {
		candidate := target + ext
		if sourceFiles[candidate] {
			return []string{candidate}
		}
	}
	return nil
}

// resolvePyImport resolves a Python import to local files.
func resolvePyImport(importPath, sourceDir string, sourceFiles map[string]bool) []string {
	// Handle relative imports (leading dots)
	cleaned := strings.TrimLeft(importPath, ".")
	dotCount := len(importPath) - len(cleaned)

	var baseDir string
	if dotCount > 0 {
		baseDir = sourceDir
		for i := 1; i < dotCount; i++ {
			baseDir = filepath.Dir(baseDir)
		}
	}

	// Convert module.path to module/path
	parts := strings.Split(cleaned, ".")
	modPath := filepath.Join(parts...)

	var candidates []string
	if baseDir != "" {
		modPath = filepath.Join(baseDir, modPath)
	}

	// Try as module.py or module/__init__.py
	pyFile := filepath.ToSlash(modPath) + ".py"
	initFile := filepath.ToSlash(filepath.Join(modPath, "__init__.py"))

	if sourceFiles[pyFile] {
		candidates = append(candidates, pyFile)
	}
	if sourceFiles[initFile] {
		candidates = append(candidates, initFile)
	}

	return candidates
}

// ReverseImports returns files that import the given file.
func ReverseImports(graph map[string][]string, targetFile string) []string {
	var result []string
	for file, imports := range graph {
		for _, imp := range imports {
			if imp == targetFile {
				result = append(result, file)
				break
			}
		}
	}
	sort.Strings(result)
	return result
}

// CircularDeps finds circular dependency chains in the import graph.
func CircularDeps(graph map[string][]string) [][]string {
	var cycles [][]string
	visited := map[string]bool{}
	recStack := map[string]bool{}
	path := []string{}

	var dfs func(node string)
	dfs = func(node string) {
		visited[node] = true
		recStack[node] = true
		path = append(path, node)

		for _, dep := range graph[node] {
			if !visited[dep] {
				dfs(dep)
			} else if recStack[dep] {
				// Found cycle — extract it from path
				cycleStart := -1
				for i, p := range path {
					if p == dep {
						cycleStart = i
						break
					}
				}
				if cycleStart >= 0 {
					cycle := make([]string, len(path)-cycleStart)
					copy(cycle, path[cycleStart:])
					cycle = append(cycle, dep) // close the cycle
					cycles = append(cycles, cycle)
				}
			}
		}

		path = path[:len(path)-1]
		recStack[node] = false
	}

	// Sort keys for deterministic output
	keys := make([]string, 0, len(graph))
	for k := range graph {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, node := range keys {
		if !visited[node] {
			dfs(node)
		}
	}
	return cycles
}

// ChangeImpact returns all files transitively depending on the target file.
func ChangeImpact(graph map[string][]string, targetFile string) []string {
	// Build reverse graph
	reverse := map[string][]string{}
	for file, imports := range graph {
		for _, imp := range imports {
			reverse[imp] = append(reverse[imp], file)
		}
	}

	// BFS on reverse graph
	visited := map[string]bool{targetFile: true}
	queue := []string{targetFile}
	var result []string

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, dep := range reverse[current] {
			if !visited[dep] {
				visited[dep] = true
				result = append(result, dep)
				queue = append(queue, dep)
			}
		}
	}

	sort.Strings(result)
	return result
}

// FormatImportTree prints the dependency tree for a file.
func FormatImportTree(graph map[string][]string, rootFile string, maxDepth int) string {
	if maxDepth <= 0 {
		maxDepth = 4
	}

	var b strings.Builder
	b.WriteString(rootFile + "\n")
	visited := map[string]bool{}
	formatTreeRecurse(&b, graph, rootFile, "", maxDepth, 0, visited)
	return b.String()
}

func formatTreeRecurse(b *strings.Builder, graph map[string][]string, node, prefix string, maxDepth, depth int, visited map[string]bool) {
	if depth >= maxDepth {
		return
	}
	visited[node] = true

	deps := graph[node]
	for i, dep := range deps {
		isLast := i == len(deps)-1
		connector := "├── "
		childPrefix := "│   "
		if isLast {
			connector = "└── "
			childPrefix = "    "
		}

		suffix := ""
		if visited[dep] {
			suffix = " (circular)"
		}
		b.WriteString(prefix + connector + dep + suffix + "\n")

		if !visited[dep] {
			formatTreeRecurse(b, graph, dep, prefix+childPrefix, maxDepth, depth+1, visited)
		}
	}

	visited[node] = false // allow node to appear in other branches
}

// FormatImportGraph formats the entire import graph as a readable string.
func FormatImportGraph(graph map[string][]string) string {
	if len(graph) == 0 {
		return "No imports found."
	}

	keys := make([]string, 0, len(graph))
	for k := range graph {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Import graph (%d files with imports):\n\n", len(keys)))
	for _, file := range keys {
		deps := graph[file]
		b.WriteString(fmt.Sprintf("%s\n", file))
		for _, dep := range deps {
			b.WriteString(fmt.Sprintf("  → %s\n", dep))
		}
		b.WriteString("\n")
	}

	// Report circular deps
	cycles := CircularDeps(graph)
	if len(cycles) > 0 {
		b.WriteString(fmt.Sprintf("⚠ %d circular dependency chain(s) detected:\n", len(cycles)))
		for _, cycle := range cycles {
			b.WriteString("  " + strings.Join(cycle, " → ") + "\n")
		}
	}

	return b.String()
}
