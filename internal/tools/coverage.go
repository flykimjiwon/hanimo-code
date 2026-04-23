package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type untested struct {
	path    string
	modTime time.Time
}

// TestCoverageGaps scans source files and finds those without test counterparts.
func TestCoverageGaps(basePath string) string {
	if basePath == "" {
		basePath = "."
	}
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return fmt.Sprintf("Error: invalid path: %v", err)
	}

	gi := LoadGitIgnore(absBase)
	var missing []untested
	const limit = 50

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

		name := d.Name()
		ext := strings.ToLower(filepath.Ext(name))
		base := strings.TrimSuffix(name, filepath.Ext(name))
		dir := filepath.Dir(path)

		switch ext {
		case ".go":
			// skip test files themselves
			if strings.HasSuffix(base, "_test") {
				return nil
			}
			testFile := filepath.Join(dir, base+"_test.go")
			if fileExists(testFile) {
				return nil
			}
		case ".ts", ".tsx", ".js", ".jsx":
			// skip files that are already test files
			if strings.HasSuffix(base, ".test") || strings.HasSuffix(base, ".spec") {
				return nil
			}
			// check foo.test.ts, foo.spec.ts, __tests__/foo.test.ts
			if fileExists(filepath.Join(dir, base+".test"+ext)) ||
				fileExists(filepath.Join(dir, base+".spec"+ext)) ||
				fileExists(filepath.Join(dir, "__tests__", base+".test"+ext)) {
				return nil
			}
		case ".py":
			if strings.HasPrefix(name, "test_") {
				return nil
			}
			if fileExists(filepath.Join(dir, "test_"+name)) ||
				fileExists(filepath.Join(dir, "tests", "test_"+name)) {
				return nil
			}
		default:
			return nil
		}

		info, _ := d.Info()
		mt := time.Time{}
		if info != nil {
			mt = info.ModTime()
		}
		missing = append(missing, untested{path: rel, modTime: mt})
		return nil
	})

	if len(missing) == 0 {
		return "All source files have test counterparts."
	}

	// Sort by mtime descending (recently changed = highest risk)
	sort.Slice(missing, func(i, j int) bool {
		return missing[i].modTime.After(missing[j].modTime)
	})

	if len(missing) > limit {
		missing = missing[:limit]
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d files without tests:\n", len(missing)))
	for _, u := range missing {
		dateStr := "unknown"
		if !u.modTime.IsZero() {
			dateStr = u.modTime.Format("2006-01-02")
		}
		sb.WriteString(fmt.Sprintf("  %s (modified: %s)\n", u.path, dateStr))
	}
	return sb.String()
}

// fileExists is defined in diagnostics.go
