package tools

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// runGit executes a git command in the given directory with a 10s timeout.
func runGit(dir string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	result := strings.TrimSpace(string(out))
	if err != nil {
		return result, fmt.Errorf("git %s: %s (%w)", args[0], result, err)
	}
	return result, nil
}

// GitStatus returns short status output for the given path.
func GitStatus(path string) (string, error) {
	return runGit(path, "status", "--short")
}

// GitDiff returns diff output. If staged is true, shows staged changes.
func GitDiff(path string, staged bool) (string, error) {
	if staged {
		return runGit(path, "diff", "--staged")
	}
	return runGit(path, "diff")
}

// GitLog returns the last n commits in oneline format.
func GitLog(path string, n int) (string, error) {
	return runGit(path, "log", fmt.Sprintf("-n%d", n), "--oneline")
}

// GitCommit creates a commit with the given message.
func GitCommit(path, message string) (string, error) {
	return runGit(path, "commit", "-m", message)
}

// GitBranch returns the current branch name.
func GitBranch(path string) (string, error) {
	return runGit(path, "branch", "--show-current")
}

// GitBlame returns blame info for a file, grouping consecutive lines by author.
// Format: "L1-L5: Author Name (2026-04-15)\nL6-L20: Other Author (2026-04-10)"
func GitBlame(dir, file string) (string, error) {
	raw, err := runGit(dir, "-C", dir, "blame", "--line-porcelain", file)
	if err != nil {
		return "", err
	}

	type lineInfo struct {
		author string
		date   string
	}

	var lines []lineInfo
	var curAuthor, curDate string

	for _, l := range strings.Split(raw, "\n") {
		if strings.HasPrefix(l, "author ") {
			curAuthor = strings.TrimPrefix(l, "author ")
		} else if strings.HasPrefix(l, "author-time ") {
			// author-time is a unix timestamp; convert to date string
			var ts int64
			fmt.Sscanf(strings.TrimPrefix(l, "author-time "), "%d", &ts)
			curDate = time.Unix(ts, 0).UTC().Format("2006-01-02")
		} else if l == "" || (len(l) > 0 && l[0] == '\t') {
			// tab-prefixed line = actual source line; commit block ends here
			if curAuthor != "" {
				lines = append(lines, lineInfo{author: curAuthor, date: curDate})
				curAuthor = ""
				curDate = ""
			}
		}
	}

	if len(lines) == 0 {
		return "(no blame data)", nil
	}

	// Group consecutive lines by same author+date
	var sb strings.Builder
	start := 1
	for i := 1; i <= len(lines); i++ {
		if i == len(lines) || lines[i].author != lines[i-1].author || lines[i].date != lines[i-1].date {
			end := i
			if start == end {
				sb.WriteString(fmt.Sprintf("L%d: %s (%s)\n", start, lines[start-1].author, lines[start-1].date))
			} else {
				sb.WriteString(fmt.Sprintf("L%d-L%d: %s (%s)\n", start, end, lines[start-1].author, lines[start-1].date))
			}
			if i < len(lines) {
				start = i + 1
			}
		}
	}

	return strings.TrimRight(sb.String(), "\n"), nil
}

// GitHotFiles returns the most frequently changed files in the last N days.
// Returns top 20 sorted by change frequency descending.
func GitHotFiles(dir string, days int) (string, error) {
	if days <= 0 {
		days = 30
	}
	since := fmt.Sprintf("--since=%d.days.ago", days)
	raw, err := runGit(dir, "-C", dir, "log", "--name-only", "--pretty=format:", since)
	if err != nil {
		return "", err
	}

	freq := map[string]int{}
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			freq[line]++
		}
	}

	if len(freq) == 0 {
		return fmt.Sprintf("No file changes in the last %d days.", days), nil
	}

	type entry struct {
		file  string
		count int
	}
	entries := make([]entry, 0, len(freq))
	for f, c := range freq {
		entries = append(entries, entry{f, c})
	}
	// Sort by frequency desc, then name asc for stability
	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].count > entries[i].count ||
				(entries[j].count == entries[i].count && entries[j].file < entries[i].file) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
	if len(entries) > 20 {
		entries = entries[:20]
	}

	var sb strings.Builder
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("  %d changes: %s\n", e.count, e.file))
	}
	return strings.TrimRight(sb.String(), "\n"), nil
}

// GitFileHistory returns recent commits for a specific file.
func GitFileHistory(dir, file string, n int) (string, error) {
	if n <= 0 {
		n = 10
	}
	format := "%h %ad %s"
	return runGit(dir, "-C", dir, "log", fmt.Sprintf("-n%d", n), "--follow",
		fmt.Sprintf("--pretty=format:%s", format), "--date=short", "--", file)
}
