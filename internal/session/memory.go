package session

import (
	"fmt"
	"strings"
	"time"
)

// Memory represents a project-scoped key-value memory entry.
type Memory struct {
	ID         int64
	ProjectDir string
	Key        string
	Value      string
	Source     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// SaveMemory upserts a memory entry for the given project.
func SaveMemory(projectDir, key, value, source string) error {
	if source == "" {
		source = "auto"
	}
	_, err := DB().Exec(
		`INSERT INTO memories (project_dir, key, value, source)
		 VALUES (?, ?, ?, ?)
		 ON CONFLICT(project_dir, key) DO UPDATE SET
		   value = excluded.value,
		   source = excluded.source,
		   updated_at = CURRENT_TIMESTAMP`,
		projectDir, key, value, source,
	)
	if err != nil {
		return fmt.Errorf("save memory: %w", err)
	}
	return nil
}

// LoadMemories returns all memory entries for a project directory.
func LoadMemories(projectDir string) ([]Memory, error) {
	rows, err := DB().Query(
		`SELECT id, project_dir, key, value, source, created_at, updated_at
		 FROM memories WHERE project_dir = ? ORDER BY key`, projectDir,
	)
	if err != nil {
		return nil, fmt.Errorf("load memories: %w", err)
	}
	defer rows.Close()

	var mems []Memory
	for rows.Next() {
		var m Memory
		if err := rows.Scan(&m.ID, &m.ProjectDir, &m.Key, &m.Value, &m.Source, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return mems, fmt.Errorf("scan memory: %w", err)
		}
		mems = append(mems, m)
	}
	return mems, rows.Err()
}

// DeleteMemory removes a memory entry by project and key.
func DeleteMemory(projectDir, key string) error {
	res, err := DB().Exec(`DELETE FROM memories WHERE project_dir = ? AND key = ?`, projectDir, key)
	if err != nil {
		return fmt.Errorf("delete memory: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("memory not found: %s/%s", projectDir, key)
	}
	return nil
}

// FormatMemoriesForPrompt formats memories as a markdown block suitable for
// injection into a system prompt.
func FormatMemoriesForPrompt(memories []Memory) string {
	if len(memories) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("## Project Memory\n\n")
	for _, m := range memories {
		b.WriteString(fmt.Sprintf("- **%s**: %s\n", m.Key, m.Value))
	}
	return b.String()
}
