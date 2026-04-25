package main

// Phase 14a — Run / Debug panel data source.
//
// Walks the project root and returns runnable targets from package.json
// scripts, Makefile, and (when go.mod is present) common Go targets. The
// frontend RunPanel renders this list and asks App.RunTarget to spawn the
// actual command in the embedded terminal.

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// RunTarget is a single clickable runnable. Source is one of "npm" |
// "make" | "go" — used by the panel for grouping + an icon hint.
type RunTarget struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Command     string `json:"command"`
	Source      string `json:"source"`
}

// GetRunTargets returns the merged list of runnable scripts in cwd.
// Sorted by source then name so the order is stable across refreshes.
func (a *App) GetRunTargets() []RunTarget {
	var targets []RunTarget
	targets = append(targets, scanNpmScripts(a.cwd)...)
	targets = append(targets, scanMakefile(a.cwd)...)
	targets = append(targets, scanGoTargets(a.cwd)...)
	sort.SliceStable(targets, func(i, j int) bool {
		if targets[i].Source != targets[j].Source {
			return targets[i].Source < targets[j].Source
		}
		return targets[i].Name < targets[j].Name
	})
	return targets
}

func scanNpmScripts(dir string) []RunTarget {
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return nil
	}
	var pkg struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil
	}
	out := make([]RunTarget, 0, len(pkg.Scripts))
	for name, cmd := range pkg.Scripts {
		out = append(out, RunTarget{
			Name:        name,
			Description: cmd,
			Command:     "npm run " + name,
			Source:      "npm",
		})
	}
	return out
}

// makeTargetRE matches `target:` at the start of a line (allowing leading
// non-whitespace identifier chars and an optional dependency list after the
// colon). Skips lines starting with `.` (built-ins like `.PHONY`) and `\t`
// (recipe bodies).
// `target:` (deps optional, allows blank/EOL after colon) — but NOT
// `VAR:=value` which is a make variable assignment. `VAR?=` / `VAR+=`
// are already excluded because they have no `:` for the regex to match.
var makeTargetRE = regexp.MustCompile(`^([A-Za-z][A-Za-z0-9_-]*)\s*:(?:[^=]|$)`)

func scanMakefile(dir string) []RunTarget {
	for _, name := range []string{"Makefile", "makefile", "GNUmakefile"} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			continue
		}
		var out []RunTarget
		seen := map[string]bool{}
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := scanner.Text()
			// Recipe lines (tab-indented) and comments are not targets.
			if strings.HasPrefix(line, "\t") || strings.HasPrefix(strings.TrimSpace(line), "#") {
				continue
			}
			m := makeTargetRE.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			target := m[1]
			if seen[target] {
				continue
			}
			seen[target] = true
			out = append(out, RunTarget{
				Name:        target,
				Description: "make " + target,
				Command:     "make " + target,
				Source:      "make",
			})
		}
		return out // first found Makefile wins
	}
	return nil
}

// scanGoTargets adds canonical Go workflow commands when a go.mod is present.
// These are static — no project introspection — but cover the 90% case.
func scanGoTargets(dir string) []RunTarget {
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err != nil {
		return nil
	}
	return []RunTarget{
		{Name: "build", Description: "Build all packages", Command: "go build ./...", Source: "go"},
		{Name: "test", Description: "Run all tests", Command: "go test ./...", Source: "go"},
		{Name: "vet", Description: "Static analysis", Command: "go vet ./...", Source: "go"},
		{Name: "mod tidy", Description: "Sync dependencies", Command: "go mod tidy", Source: "go"},
	}
}

// RunTarget executes the named target in the embedded terminal. Starts the
// terminal if it isn't running yet so the user doesn't have to open it
// manually first.
func (a *App) RunTarget(command string) error {
	if err := a.StartTerminal(); err != nil {
		return err
	}
	a.WriteTerminal(command + "\n")
	return nil
}
