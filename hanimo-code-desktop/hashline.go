package main

// Ported from hanimo CLI internal/tools/hashline.go (2026-04-25, Phase 10).
// Same MD5 4-char prefix anchor algorithm so anchors emitted by either
// product line stay interchangeable.

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// hashLine computes a 4-char MD5 hash prefix for a line of text.
func hashLine(line string) string {
	h := md5.Sum([]byte(line))
	return fmt.Sprintf("%x", h)[:4]
}

// parseAnchor parses an anchor like "3#e4d9" into line number and hash.
func parseAnchor(anchor string) (int, string, error) {
	parts := strings.SplitN(anchor, "#", 2)
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid anchor format %q: expected N#hash", anchor)
	}
	line, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", fmt.Errorf("invalid line number in anchor %q: %v", anchor, err)
	}
	if line < 1 {
		return 0, "", fmt.Errorf("line number must be >= 1, got %d", line)
	}
	return line, parts[1], nil
}

// HashlineRead reads a file and returns each line tagged with a 4-char MD5 hash anchor.
// Output format: "1#a3f1| function hello() {"
func HashlineRead(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("cannot resolve path: %v", err)
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("cannot read file: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	var sb strings.Builder
	for i, line := range lines {
		h := hashLine(line)
		fmt.Fprintf(&sb, "%d#%s| %s\n", i+1, h, line)
	}
	return sb.String(), nil
}

// HashlineEdit edits a file using hash anchors to verify lines haven't changed
// since read. startAnchor and endAnchor are in "N#hash" format. Lines from
// startLine to endLine (inclusive) are replaced with newContent.
//
// Caller is responsible for invoking backupBeforeWrite + hash:anchor emit;
// chat.go's hashline_edit tool handler does both.
func HashlineEdit(path, startAnchor, endAnchor, newContent string) (string, int, int, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", 0, 0, fmt.Errorf("cannot resolve path: %v", err)
	}

	startLine, startHash, err := parseAnchor(startAnchor)
	if err != nil {
		return "", 0, 0, err
	}
	endLine, endHash, err := parseAnchor(endAnchor)
	if err != nil {
		return "", 0, 0, err
	}

	if endLine < startLine {
		return "", 0, 0, fmt.Errorf("end line %d is before start line %d", endLine, startLine)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return "", 0, 0, fmt.Errorf("cannot read file: %v", err)
	}

	lines := strings.Split(string(data), "\n")

	if startLine > len(lines) {
		return "", 0, 0, fmt.Errorf("start line %d exceeds file length %d", startLine, len(lines))
	}
	if endLine > len(lines) {
		return "", 0, 0, fmt.Errorf("end line %d exceeds file length %d", endLine, len(lines))
	}

	actualStartHash := hashLine(lines[startLine-1])
	if actualStartHash != startHash {
		return "", 0, 0, fmt.Errorf("hash mismatch at line %d: expected %s, got %s. File changed — re-read with hashline_read", startLine, startHash, actualStartHash)
	}
	actualEndHash := hashLine(lines[endLine-1])
	if actualEndHash != endHash {
		return "", 0, 0, fmt.Errorf("hash mismatch at line %d: expected %s, got %s. File changed — re-read with hashline_read", endLine, endHash, actualEndHash)
	}

	newLines := strings.Split(newContent, "\n")
	oldCount := endLine - startLine + 1

	result := make([]string, 0, len(lines)-oldCount+len(newLines))
	result = append(result, lines[:startLine-1]...)
	result = append(result, newLines...)
	result = append(result, lines[endLine:]...)

	if err := os.WriteFile(absPath, []byte(strings.Join(result, "\n")), 0644); err != nil {
		return "", 0, 0, fmt.Errorf("cannot write file: %v", err)
	}

	return fmt.Sprintf("OK: replaced lines %d-%d (%d lines → %d lines)", startLine, endLine, oldCount, len(newLines)), startLine, len(newLines), nil
}
