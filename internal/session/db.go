package session

import (
	"database/sql"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"
)

const schemaSQL = `
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    name TEXT,
    project_dir TEXT NOT NULL,
    provider TEXT NOT NULL,
    model TEXT NOT NULL,
    mode TEXT DEFAULT 'super',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL REFERENCES sessions(id),
    role TEXT NOT NULL,
    content TEXT NOT NULL,
    tool_calls TEXT,
    tool_result TEXT,
    tokens_in INTEGER DEFAULT 0,
    tokens_out INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS memories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_dir TEXT NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    source TEXT DEFAULT 'auto',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_dir, key)
);
CREATE TABLE IF NOT EXISTS usage_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL REFERENCES sessions(id),
    provider TEXT NOT NULL,
    model TEXT NOT NULL,
    tokens_in INTEGER NOT NULL,
    tokens_out INTEGER NOT NULL,
    cost_usd REAL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS mcp_servers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    transport TEXT NOT NULL,
    command TEXT,
    args TEXT,
    url TEXT,
    enabled BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_messages_session ON messages(session_id);
CREATE INDEX IF NOT EXISTS idx_memories_project ON memories(project_dir);
CREATE INDEX IF NOT EXISTS idx_usage_session ON usage_log(session_id);
`

var (
	globalDB *sql.DB
	dbMu     sync.Mutex
)

// InitDB opens (or creates) the SQLite database at configDir/hanimo.db,
// enables WAL mode, and runs the schema migration.
func InitDB(configDir string) error {
	dbMu.Lock()
	defer dbMu.Unlock()

	dbPath := fmt.Sprintf("%s/hanimo.db", configDir)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	// Enable WAL mode for better concurrent read performance.
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return fmt.Errorf("set WAL mode: %w", err)
	}

	// Run schema migration.
	if _, err := db.Exec(schemaSQL); err != nil {
		db.Close()
		return fmt.Errorf("migrate schema: %w", err)
	}

	globalDB = db
	return nil
}

// DB returns the global database handle. Panics if InitDB has not been called.
func DB() *sql.DB {
	dbMu.Lock()
	defer dbMu.Unlock()
	if globalDB == nil {
		panic("session.DB() called before InitDB")
	}
	return globalDB
}

// CloseDB closes the global database connection.
func CloseDB() {
	dbMu.Lock()
	defer dbMu.Unlock()
	if globalDB != nil {
		globalDB.Close()
		globalDB = nil
	}
}
