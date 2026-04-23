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

// Directories to skip during search traversal.
var skipDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	"dist":         true,
	"__pycache__":  true,
	".next":        true,
	"vendor":       true,
	".omc":         true,
}

// Binary file extensions to skip.
var binaryExts = map[string]bool{
	// Executables & Libraries
	".exe": true, ".dll": true, ".so": true, ".dylib": true, ".a": true,
	".lib": true, ".obj": true, ".o": true, ".ko": true, ".elf": true,
	".bin": true, ".out": true, ".app": true, ".deb": true, ".rpm": true,
	".msi": true, ".dmg": true, ".pkg": true, ".snap": true, ".flatpak": true,
	// Archives
	".zip": true, ".tar": true, ".gz": true, ".bz2": true, ".xz": true,
	".7z": true, ".rar": true, ".zst": true, ".lz4": true, ".lzma": true,
	".tgz": true, ".tbz2": true, ".war": true, ".ear": true, ".jar": true,
	// Images
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".bmp": true,
	".ico": true, ".webp": true, ".avif": true, ".tiff": true, ".tif": true,
	// .svg is intentionally excluded — SVG is text/XML
	".psd": true, ".ai": true, ".eps": true, ".raw": true, ".cr2": true,
	".nef": true, ".heic": true, ".heif": true,
	// Audio
	".mp3": true, ".wav": true, ".flac": true, ".aac": true, ".ogg": true,
	".wma": true, ".m4a": true, ".opus": true,
	// Video
	".mp4": true, ".avi": true, ".mkv": true, ".mov": true, ".wmv": true,
	".flv": true, ".webm": true, ".m4v": true, ".3gp": true,
	// Fonts
	".ttf": true, ".otf": true, ".woff": true, ".woff2": true, ".eot": true,
	// Documents
	".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
	".ppt": true, ".pptx": true, ".odt": true, ".ods": true, ".odp": true,
	// Data
	".db": true, ".sqlite": true, ".sqlite3": true, ".mdb": true,
	".dat": true, ".sav": true, ".bak": true,
	// Compiled/Bytecode
	".pyc": true, ".pyo": true, ".class": true, ".wasm": true,
	// Other
	".iso": true, ".img": true, ".vmdk": true, ".qcow2": true,
	".DS_Store": true, ".lock": true,
}

const (
	maxGrepMatches = 300
	maxGrepBytes   = 60000
	maxGlobFiles   = 5000
	maxLineChars   = 2000
)

// isBinaryFile checks if a file is binary by sampling the first 4096 bytes.
// Returns true if null bytes found or >30% non-printable characters.
func isBinaryFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 4096)
	n, err := f.Read(buf)
	if n == 0 {
		return false
	}
	buf = buf[:n]

	nonPrintable := 0
	for _, b := range buf {
		if b == 0x00 {
			return true // null byte → binary
		}
		if b < 0x20 && b != '\n' && b != '\r' && b != '\t' {
			nonPrintable++
		}
	}

	return float64(nonPrintable)/float64(n) > 0.3
}

func truncateLine(line string) string {
	runes := []rune(line)
	if len(runes) > maxLineChars {
		return string(runes[:maxLineChars]) + "..."
	}
	return line
}

// GrepSearch searches file contents by regex pattern.
// Returns matches in "file:line:content" format.
// Uses ripgrep (rg) when available for 100x performance, falls back to Go.
func GrepSearch(pattern, basePath, glob string, ignoreCase bool, contextLines int) (string, error) {
	if pattern == "" {
		return "", fmt.Errorf("pattern is required")
	}

	// Try ripgrep first for performance
	if IsRipgrepAvailable() {
		result, err := RipgrepSearch(pattern, basePath, glob, ignoreCase, contextLines)
		if err == nil {
			return result, nil
		}
		config.DebugLog("[GREP] ripgrep failed, falling back to Go: %v", err)
	}

	flags := ""
	if ignoreCase {
		flags = "(?i)"
	}
	re, err := regexp.Compile(flags + pattern)
	if err != nil {
		return "", fmt.Errorf("invalid regex: %w", err)
	}

	if basePath == "" {
		basePath = "."
	}
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	var globRe *regexp.Regexp
	if glob != "" {
		globPattern := globToRegex(glob)
		globRe, err = regexp.Compile(globPattern)
		if err != nil {
			return "", fmt.Errorf("invalid glob filter %q: %w", glob, err)
		}
	}

	var results strings.Builder
	matchCount := 0

	walkErr := filepath.WalkDir(absBase, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip binary files (extension check + byte sampling)
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if binaryExts[ext] {
			return nil
		}
		if isBinaryFile(path) {
			return nil
		}

		// Apply glob filter on relative path (use forward slashes for cross-platform regex)
		rel, _ := filepath.Rel(absBase, path)
		rel = filepath.ToSlash(rel)
		if globRe != nil && !globRe.MatchString(rel) && !globRe.MatchString(d.Name()) {
			return nil
		}

		// Skip large files (>1MB)
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if info.Size() > 5*1024*1024 {
			return nil
		}

		matches, err := searchFile(path, rel, re, contextLines)
		if err != nil {
			return nil
		}

		for _, m := range matches {
			if matchCount >= maxGrepMatches || results.Len() >= maxGrepBytes {
				results.WriteString(fmt.Sprintf("\n... (truncated, %d+ matches)\n", matchCount))
				return fmt.Errorf("limit reached")
			}
			results.WriteString(m)
			results.WriteString("\n")
			matchCount++
		}

		return nil
	})

	if walkErr != nil && walkErr.Error() != "limit reached" {
		config.DebugLog("[GREP] walk error: %v", walkErr)
	}

	if matchCount == 0 {
		return "No matches found.", nil
	}

	config.DebugLog("[GREP] pattern=%q path=%s matches=%d bytes=%d", pattern, basePath, matchCount, results.Len())
	return results.String(), nil
}

// searchFile scans a single file for regex matches.
func searchFile(absPath, relPath string, re *regexp.Regexp, contextLines int) ([]string, error) {
	f, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var allLines []string
	var matchLineNums []int

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 256*1024)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		allLines = append(allLines, line)
		if re.MatchString(line) {
			matchLineNums = append(matchLineNums, lineNum)
		}
	}

	if len(matchLineNums) == 0 {
		return nil, nil
	}

	var results []string
	if contextLines <= 0 {
		for _, ln := range matchLineNums {
			results = append(results, fmt.Sprintf("%s:%d:%s", relPath, ln, truncateLine(allLines[ln-1])))
		}
	} else {
		// With context lines, group nearby matches
		shown := make(map[int]bool)
		for _, ln := range matchLineNums {
			start := ln - contextLines
			if start < 1 {
				start = 1
			}
			end := ln + contextLines
			if end > len(allLines) {
				end = len(allLines)
			}
			for i := start; i <= end; i++ {
				if !shown[i] {
					prefix := " "
					if i == ln {
						prefix = ">"
					}
					results = append(results, fmt.Sprintf("%s:%d:%s %s", relPath, i, prefix, truncateLine(allLines[i-1])))
					shown[i] = true
				}
			}
		}
	}

	return results, nil
}

// GlobSearch finds files matching a glob pattern (supports **).
func GlobSearch(pattern, basePath string) (string, error) {
	if pattern == "" {
		return "", fmt.Errorf("pattern is required")
	}

	if basePath == "" {
		basePath = "."
	}
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	// Convert glob pattern to regex for ** support
	globRe, err := regexp.Compile(globToRegex(pattern))
	if err != nil {
		return "", fmt.Errorf("invalid glob pattern: %w", err)
	}

	type globEntry struct {
		rel   string
		mtime time.Time
	}
	var matches []globEntry

	walkErr := filepath.WalkDir(absBase, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		rel, _ := filepath.Rel(absBase, path)
		rel = filepath.ToSlash(rel)
		if globRe.MatchString(rel) {
			mtime := time.Time{}
			if info, err := d.Info(); err == nil {
				mtime = info.ModTime()
			}
			matches = append(matches, globEntry{rel: rel, mtime: mtime})
			if len(matches) >= maxGlobFiles {
				return fmt.Errorf("limit reached")
			}
		}
		return nil
	})

	if walkErr != nil && walkErr.Error() != "limit reached" {
		config.DebugLog("[GLOB] walk error: %v", walkErr)
	}

	if len(matches) == 0 {
		return "No files matched.", nil
	}

	// Sort by mtime descending (most recently modified first)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].mtime.After(matches[j].mtime)
	})

	var names []string
	for _, m := range matches {
		names = append(names, m.rel)
	}

	result := strings.Join(names, "\n")
	if len(matches) >= maxGlobFiles {
		result += fmt.Sprintf("\n... (truncated at %d files)", maxGlobFiles)
	}

	config.DebugLog("[GLOB] pattern=%q path=%s matches=%d", pattern, basePath, len(matches))
	return result, nil
}

// globToRegex converts a glob pattern to a regex string.
// Supports: ** (any path), * (any non-separator), ? (single char), {a,b} (alternation)
func globToRegex(pattern string) string {
	var b strings.Builder
	b.WriteString("^")

	braceDepth := 0
	i := 0
	for i < len(pattern) {
		ch := pattern[i]
		switch ch {
		case '*':
			if i+1 < len(pattern) && pattern[i+1] == '*' {
				// ** matches any number of path segments
				if i+2 < len(pattern) && pattern[i+2] == '/' {
					b.WriteString("(.*/)?")
					i += 3
				} else {
					b.WriteString(".*")
					i += 2
				}
			} else {
				// * matches anything except /
				b.WriteString("[^/]*")
				i++
			}
		case '?':
			b.WriteString("[^/]")
			i++
		case '.':
			b.WriteString("\\.")
			i++
		case '\\':
			b.WriteString("\\\\")
			i++
		case '{':
			// {a,b,c} → (a|b|c)
			b.WriteString("(")
			braceDepth++
			i++
		case '}':
			b.WriteString(")")
			if braceDepth > 0 {
				braceDepth--
			}
			i++
		case ',':
			if braceDepth > 0 {
				b.WriteString("|")
			} else {
				b.WriteByte(',') // literal comma outside braces
			}
			i++
		default:
			b.WriteByte(ch)
			i++
		}
	}

	b.WriteString("$")
	return b.String()
}
