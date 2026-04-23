package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/flykimjiwon/hanimo/internal/config"
)

// PatchOperation represents a single file operation within an apply_patch call.
type PatchOperation struct {
	Type    string // "add", "update", "delete"
	Path    string
	MoveTo  string // for update with rename
	Hunks   []PatchHunk
	Content string // for add: full file content
}

// PatchHunk represents one change block within an Update operation.
type PatchHunk struct {
	Context     string // @@ anchor (function name or line pattern)
	RemoveLines []string
	AddLines    []string
}

// ParsePatch parses the Codex-style apply_patch format.
//
// Format:
//
//	*** Begin Patch
//	*** Add File: path/to/new.go
//	+line one
//	+line two
//	*** Update File: src/app.go
//	*** Move to: src/main.go
//	@@ func (a *App) Run() {
//	-old line
//	+new line
//	*** Delete File: obsolete.go
//	*** End Patch
func ParsePatch(input string) ([]PatchOperation, error) {
	lines := strings.Split(input, "\n")
	var ops []PatchOperation
	var current *PatchOperation
	var currentHunk *PatchHunk
	inPatch := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "*** Begin Patch" {
			inPatch = true
			continue
		}
		if trimmed == "*** End Patch" {
			if currentHunk != nil && current != nil {
				current.Hunks = append(current.Hunks, *currentHunk)
				currentHunk = nil
			}
			if current != nil {
				ops = append(ops, *current)
				current = nil
			}
			break
		}
		if !inPatch {
			continue
		}

		// New operation headers
		if strings.HasPrefix(trimmed, "*** Add File:") {
			if currentHunk != nil && current != nil {
				current.Hunks = append(current.Hunks, *currentHunk)
				currentHunk = nil
			}
			if current != nil {
				ops = append(ops, *current)
			}
			path := strings.TrimSpace(strings.TrimPrefix(trimmed, "*** Add File:"))
			current = &PatchOperation{Type: "add", Path: path}
			currentHunk = nil
			continue
		}
		if strings.HasPrefix(trimmed, "*** Update File:") {
			if currentHunk != nil && current != nil {
				current.Hunks = append(current.Hunks, *currentHunk)
				currentHunk = nil
			}
			if current != nil {
				ops = append(ops, *current)
			}
			path := strings.TrimSpace(strings.TrimPrefix(trimmed, "*** Update File:"))
			current = &PatchOperation{Type: "update", Path: path}
			currentHunk = nil
			continue
		}
		if strings.HasPrefix(trimmed, "*** Delete File:") {
			if currentHunk != nil && current != nil {
				current.Hunks = append(current.Hunks, *currentHunk)
				currentHunk = nil
			}
			if current != nil {
				ops = append(ops, *current)
			}
			path := strings.TrimSpace(strings.TrimPrefix(trimmed, "*** Delete File:"))
			ops = append(ops, PatchOperation{Type: "delete", Path: path})
			current = nil
			currentHunk = nil
			continue
		}
		if strings.HasPrefix(trimmed, "*** Move to:") && current != nil {
			current.MoveTo = strings.TrimSpace(strings.TrimPrefix(trimmed, "*** Move to:"))
			continue
		}

		if current == nil {
			continue
		}

		// Inside an operation
		switch current.Type {
		case "add":
			// Collect content lines (strip leading +)
			if strings.HasPrefix(line, "+") {
				current.Content += strings.TrimPrefix(line, "+") + "\n"
			} else if trimmed == "" {
				current.Content += "\n"
			}

		case "update":
			// @@ anchor starts a new hunk
			if strings.HasPrefix(trimmed, "@@") {
				if currentHunk != nil {
					current.Hunks = append(current.Hunks, *currentHunk)
				}
				anchor := strings.TrimSpace(strings.TrimPrefix(trimmed, "@@"))
				currentHunk = &PatchHunk{Context: anchor}
				continue
			}
			if currentHunk == nil {
				// Lines before first @@ anchor — create an implicit hunk
				if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "+") {
					currentHunk = &PatchHunk{}
				} else {
					continue
				}
			}
			if strings.HasPrefix(line, "-") {
				currentHunk.RemoveLines = append(currentHunk.RemoveLines, strings.TrimPrefix(line, "-"))
			} else if strings.HasPrefix(line, "+") {
				currentHunk.AddLines = append(currentHunk.AddLines, strings.TrimPrefix(line, "+"))
			}

		default:
			_ = i // suppress unused
		}
	}

	// Handle unterminated patch (no *** End Patch)
	if current != nil {
		if currentHunk != nil {
			current.Hunks = append(current.Hunks, *currentHunk)
		}
		ops = append(ops, *current)
	}

	if len(ops) == 0 {
		return nil, fmt.Errorf("no patch operations found — expected *** Begin Patch / *** End Patch format")
	}
	return ops, nil
}

// validatePatchPath checks that a file path doesn't escape the working directory.
func validatePatchPath(path string) error {
	if filepath.IsAbs(path) {
		return fmt.Errorf("absolute paths not allowed: %s", path)
	}
	cwd, _ := os.Getwd()
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %s", path)
	}
	rel, err := filepath.Rel(cwd, absPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("path escapes working directory: %s", path)
	}
	return nil
}

// ApplyPatch parses and applies a multi-file patch. Returns a summary of all operations.
func ApplyPatch(patchText string) (string, error) {
	ops, err := ParsePatch(patchText)
	if err != nil {
		return "", err
	}

	var results []string
	var errors []string

	for _, op := range ops {
		switch op.Type {
		case "add":
			result, err := applyAdd(op)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Add %s: %v", op.Path, err))
			} else {
				results = append(results, result)
			}

		case "update":
			result, err := applyUpdate(op)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Update %s: %v", op.Path, err))
			} else {
				results = append(results, result)
			}

		case "delete":
			result, err := applyDelete(op)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Delete %s: %v", op.Path, err))
			} else {
				results = append(results, result)
			}
		}
	}

	var summary strings.Builder
	if len(results) > 0 {
		summary.WriteString(strings.Join(results, "\n"))
	}
	if len(errors) > 0 {
		if summary.Len() > 0 {
			summary.WriteString("\n\n")
		}
		summary.WriteString("Errors:\n")
		summary.WriteString(strings.Join(errors, "\n"))
	}
	return summary.String(), nil
}

func applyAdd(op PatchOperation) (string, error) {
	if err := validatePatchPath(op.Path); err != nil {
		return "", err
	}
	absPath, err := filepath.Abs(op.Path)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	if _, err := os.Stat(absPath); err == nil {
		return "", fmt.Errorf("file already exists: %s (use Update File instead)", op.Path)
	}

	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("mkdir failed: %w", err)
	}

	content := strings.TrimRight(op.Content, "\n") + "\n"
	if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("write failed: %w", err)
	}

	config.DebugLog("[PATCH-ADD] created %s (%d bytes)", op.Path, len(content))
	return fmt.Sprintf("+ Created %s (%d bytes)", op.Path, len(content)), nil
}

func applyUpdate(op PatchOperation) (string, error) {
	if err := validatePatchPath(op.Path); err != nil {
		return "", err
	}
	if op.MoveTo != "" {
		if err := validatePatchPath(op.MoveTo); err != nil {
			return "", err
		}
	}
	absPath, err := filepath.Abs(op.Path)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("read failed: %w", err)
	}
	content := string(data)
	originalContent := content

	hunksApplied := 0
	var failedHunks []string
	for i, hunk := range op.Hunks {
		newContent, applied := applyHunk(content, hunk)
		if applied {
			content = newContent
			hunksApplied++
		} else {
			config.DebugLog("[PATCH-UPDATE] hunk %d/%d failed in %s (context=%q)", i+1, len(op.Hunks), op.Path, hunk.Context)
			failedHunks = append(failedHunks, fmt.Sprintf("hunk %d (context=%q)", i+1, hunk.Context))
		}
	}

	if hunksApplied == 0 && len(failedHunks) > 0 {
		return "", fmt.Errorf("all hunks failed in %s: %s", op.Path, strings.Join(failedHunks, ", "))
	}

	if content == originalContent {
		return fmt.Sprintf("~ %s: no changes (all hunks matched but content identical)", op.Path), nil
	}

	// Sanity check: detect broken top-of-file insertions
	origFirstLine := strings.TrimSpace(strings.SplitN(originalContent, "\n", 2)[0])
	newFirstLine := strings.TrimSpace(strings.SplitN(content, "\n", 2)[0])
	if origFirstLine != newFirstLine && strings.Contains(content, origFirstLine) {
		isOrigDirective := strings.HasPrefix(origFirstLine, "\"use ") ||
			strings.HasPrefix(origFirstLine, "import ") ||
			strings.HasPrefix(origFirstLine, "package ") ||
			strings.HasPrefix(origFirstLine, "#!")
		isNewDirective := strings.HasPrefix(newFirstLine, "\"use ") ||
			strings.HasPrefix(newFirstLine, "import ") ||
			strings.HasPrefix(newFirstLine, "package ") ||
			strings.HasPrefix(newFirstLine, "#!") ||
			strings.HasPrefix(newFirstLine, "//") ||
			strings.HasPrefix(newFirstLine, "/*")
		if isOrigDirective && !isNewDirective {
			config.DebugLog("[PATCH-SANITY] code inserted above %q directive — rejecting", origFirstLine)
			return "", fmt.Errorf("patch rejected: code was inserted above %q in %s — this would break the file", origFirstLine, op.Path)
		}
	}

	diff := GenerateUnifiedDiff(op.Path, originalContent, content)

	if err := SnapshotAndWrite(absPath, []byte(content)); err != nil {
		return "", fmt.Errorf("write failed: %w", err)
	}

	// Handle move (rename)
	if op.MoveTo != "" {
		newAbsPath, err := filepath.Abs(op.MoveTo)
		if err != nil {
			return "", fmt.Errorf("invalid move target: %w", err)
		}
		newDir := filepath.Dir(newAbsPath)
		if err := os.MkdirAll(newDir, 0755); err != nil {
			return "", fmt.Errorf("mkdir for move failed: %w", err)
		}
		if err := os.WriteFile(newAbsPath, []byte(content), 0644); err != nil {
			return "", fmt.Errorf("write to new path failed: %w", err)
		}
		_ = os.Remove(absPath)
		config.DebugLog("[PATCH-MOVE] %s → %s", op.Path, op.MoveTo)
		result := fmt.Sprintf("~ Updated+Moved %s → %s (%d hunks)", op.Path, op.MoveTo, hunksApplied)
		if diff != "" {
			result += "\n" + diff
		}
		return result, nil
	}

	config.DebugLog("[PATCH-UPDATE] %s (%d/%d hunks applied, %d failed)", op.Path, hunksApplied, len(op.Hunks), len(failedHunks))
	result := fmt.Sprintf("~ Updated %s (%d hunks)", op.Path, hunksApplied)
	if len(failedHunks) > 0 {
		result += fmt.Sprintf("\nWARNING: %d hunk(s) failed: %s. Use file_write to rewrite the full file instead.", len(failedHunks), strings.Join(failedHunks, ", "))
	}
	if diff != "" {
		result += "\n" + diff
	}
	return result, nil
}

// applyHunk applies a single hunk to the content string.
func applyHunk(content string, hunk PatchHunk) (string, bool) {
	lines := strings.Split(content, "\n")

	startIdx := 0
	if hunk.Context != "" {
		found := false
		for i, line := range lines {
			if strings.Contains(line, hunk.Context) {
				startIdx = i
				found = true
				break
			}
		}
		if !found {
			for i, line := range lines {
				if strings.Contains(strings.TrimSpace(line), strings.TrimSpace(hunk.Context)) {
					startIdx = i
					found = true
					break
				}
			}
			if !found {
				return "", false
			}
		}
	}

	// Pure insertion (no remove lines)
	if len(hunk.RemoveLines) == 0 {
		if hunk.Context == "" && len(lines) > 1 {
			config.DebugLog("[PATCH-HUNK] pure insertion without @@ anchor on non-empty file — rejecting")
			return "", false
		}
		insertIdx := startIdx + 1
		newLines := make([]string, 0, len(lines)+len(hunk.AddLines))
		newLines = append(newLines, lines[:insertIdx]...)
		newLines = append(newLines, hunk.AddLines...)
		newLines = append(newLines, lines[insertIdx:]...)
		return strings.Join(newLines, "\n"), true
	}

	// Find remove lines starting from startIdx
	matchIdx := -1
	for i := startIdx; i <= len(lines)-len(hunk.RemoveLines); i++ {
		matched := true
		for j, removeLine := range hunk.RemoveLines {
			if strings.TrimRight(lines[i+j], " \t") != strings.TrimRight(removeLine, " \t") {
				matched = false
				break
			}
		}
		if matched {
			matchIdx = i
			break
		}
	}

	// Fallback: fuzzy match with trimmed whitespace
	if matchIdx == -1 {
		for i := startIdx; i <= len(lines)-len(hunk.RemoveLines); i++ {
			matched := true
			for j, removeLine := range hunk.RemoveLines {
				if strings.TrimSpace(lines[i+j]) != strings.TrimSpace(removeLine) {
					matched = false
					break
				}
			}
			if matched {
				matchIdx = i
				break
			}
		}
	}

	if matchIdx == -1 {
		return "", false
	}

	newLines := make([]string, 0, len(lines)-len(hunk.RemoveLines)+len(hunk.AddLines))
	newLines = append(newLines, lines[:matchIdx]...)
	newLines = append(newLines, hunk.AddLines...)
	newLines = append(newLines, lines[matchIdx+len(hunk.RemoveLines):]...)

	return strings.Join(newLines, "\n"), true
}

func applyDelete(op PatchOperation) (string, error) {
	if err := validatePatchPath(op.Path); err != nil {
		return "", err
	}
	absPath, err := filepath.Abs(op.Path)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	CreateSnapshot(absPath)

	if err := os.Remove(absPath); err != nil {
		return "", fmt.Errorf("delete failed: %w", err)
	}

	config.DebugLog("[PATCH-DELETE] removed %s", op.Path)
	return fmt.Sprintf("- Deleted %s", op.Path), nil
}
