package session

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Session represents a chat session row.
type Session struct {
	ID         string
	Name       string
	ProjectDir string
	Provider   string
	Model      string
	Mode       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Message represents a single message in a session.
type Message struct {
	ID         int64
	SessionID  string
	Role       string
	Content    string
	ToolCalls  string
	ToolResult string
	TokensIn   int
	TokensOut  int
	CreatedAt  time.Time
}

// CreateSession inserts a new session and returns its ID.
func CreateSession(projectDir, provider, model, mode string) (string, error) {
	id := uuid.New().String()
	_, err := DB().Exec(
		`INSERT INTO sessions (id, project_dir, provider, model, mode) VALUES (?, ?, ?, ?, ?)`,
		id, projectDir, provider, model, mode,
	)
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}
	return id, nil
}

// SaveMessage appends a message to a session.
func SaveMessage(sessionID, role, content, toolCalls, toolResult string, tokensIn, tokensOut int) error {
	_, err := DB().Exec(
		`INSERT INTO messages (session_id, role, content, tool_calls, tool_result, tokens_in, tokens_out) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		sessionID, role, content, toolCalls, toolResult, tokensIn, tokensOut,
	)
	if err != nil {
		return fmt.Errorf("save message: %w", err)
	}
	// Touch the session's updated_at.
	_, _ = DB().Exec(`UPDATE sessions SET updated_at = CURRENT_TIMESTAMP WHERE id = ?`, sessionID)
	return nil
}

// LoadSession retrieves a session by ID along with all its messages.
func LoadSession(id string) (*Session, []Message, error) {
	s := &Session{}
	err := DB().QueryRow(
		`SELECT id, COALESCE(name,''), project_dir, provider, model, mode, created_at, updated_at FROM sessions WHERE id = ?`, id,
	).Scan(&s.ID, &s.Name, &s.ProjectDir, &s.Provider, &s.Model, &s.Mode, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, fmt.Errorf("session not found: %s", id)
		}
		return nil, nil, fmt.Errorf("load session: %w", err)
	}

	rows, err := DB().Query(
		`SELECT id, session_id, role, content, COALESCE(tool_calls,''), COALESCE(tool_result,''), tokens_in, tokens_out, created_at FROM messages WHERE session_id = ? ORDER BY id`, id,
	)
	if err != nil {
		return s, nil, fmt.Errorf("load messages: %w", err)
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.SessionID, &m.Role, &m.Content, &m.ToolCalls, &m.ToolResult, &m.TokensIn, &m.TokensOut, &m.CreatedAt); err != nil {
			return s, msgs, fmt.Errorf("scan message: %w", err)
		}
		msgs = append(msgs, m)
	}
	return s, msgs, rows.Err()
}

// ListSessions returns sessions ordered by most recently updated, with an optional limit.
func ListSessions(limit int) ([]Session, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := DB().Query(
		`SELECT id, COALESCE(name,''), project_dir, provider, model, mode, created_at, updated_at FROM sessions ORDER BY updated_at DESC LIMIT ?`, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var s Session
		if err := rows.Scan(&s.ID, &s.Name, &s.ProjectDir, &s.Provider, &s.Model, &s.Mode, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return sessions, fmt.Errorf("scan session: %w", err)
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

// SearchSessions finds sessions whose name or messages contain the query string.
func SearchSessions(query string) ([]Session, error) {
	like := "%" + query + "%"
	rows, err := DB().Query(
		`SELECT DISTINCT s.id, COALESCE(s.name,''), s.project_dir, s.provider, s.model, s.mode, s.created_at, s.updated_at
		 FROM sessions s
		 LEFT JOIN messages m ON m.session_id = s.id
		 WHERE s.name LIKE ? OR m.content LIKE ?
		 ORDER BY s.updated_at DESC
		 LIMIT 50`, like, like,
	)
	if err != nil {
		return nil, fmt.Errorf("search sessions: %w", err)
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var s Session
		if err := rows.Scan(&s.ID, &s.Name, &s.ProjectDir, &s.Provider, &s.Model, &s.Mode, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return sessions, fmt.Errorf("scan session: %w", err)
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

// NameSession sets a human-readable name for a session.
func NameSession(id, name string) error {
	res, err := DB().Exec(`UPDATE sessions SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, name, id)
	if err != nil {
		return fmt.Errorf("name session: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("session not found: %s", id)
	}
	return nil
}

// ForkSession duplicates a session and its messages into a new session, returning the new ID.
func ForkSession(srcID string) (string, error) {
	src, msgs, err := LoadSession(srcID)
	if err != nil {
		return "", fmt.Errorf("fork: %w", err)
	}

	newID, err := CreateSession(src.ProjectDir, src.Provider, src.Model, src.Mode)
	if err != nil {
		return "", fmt.Errorf("fork create: %w", err)
	}

	for _, m := range msgs {
		if err := SaveMessage(newID, m.Role, m.Content, m.ToolCalls, m.ToolResult, m.TokensIn, m.TokensOut); err != nil {
			return newID, fmt.Errorf("fork message: %w", err)
		}
	}
	return newID, nil
}
