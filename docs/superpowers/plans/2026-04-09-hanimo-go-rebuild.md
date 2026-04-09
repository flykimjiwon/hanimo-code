# hanimo Go Rebuild — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** TECHAI_CODE(Go)를 포크하여 hanimo(오픈소스 멀티프로바이더 AI 코딩 에이전트)를 Go로 리빌드한다.

**Architecture:** TECHAI_CODE의 Bubble Tea TUI + OpenAI-compatible 클라이언트를 기반으로, 14개 프로바이더 레지스트리, SQLite 세션/메모리, Hash-Anchored 편집, Smart Compaction, MCP 클라이언트를 추가한다.

**Tech Stack:** Go 1.26+, Bubble Tea v2, Lipgloss v2, go-openai, anthropic-sdk-go, generative-ai-go, modernc.org/sqlite

---

## File Structure

```
cmd/hanimo/main.go                   ← NEW (tgc 포크 + 리네임)
internal/
├── config/config.go                 ← MODIFY (경로/변수명 변경, 프로바이더 설정 확장)
├── llm/
│   ├── client.go                    ← MODIFY (Provider 인터페이스 사용으로 전환)
│   ├── providers/
│   │   ├── registry.go              ← NEW (Provider 인터페이스 + 팩토리)
│   │   ├── openai_compat.go         ← NEW (8개 클라우드 + custom)
│   │   ├── anthropic.go             ← NEW (anthropic-sdk-go)
│   │   ├── google.go                ← NEW (generative-ai-go)
│   │   └── ollama.go                ← NEW (Ollama REST)
│   ├── compaction.go                ← NEW (Smart Compaction)
│   ├── capabilities.go             ← NEW (모델 능력 매트릭스)
│   ├── context.go                   ← MODIFY (hanimo 브랜딩)
│   ├── models.go                    ← MODIFY (멀티 프로바이더 모델 목록)
│   └── prompt.go                    ← MODIFY (시스템 프롬프트 + Role)
├── tools/
│   ├── registry.go                  ← MODIFY (신규 도구 등록)
│   ├── hashline.go                  ← NEW (Hash-Anchored 편집)
│   ├── git.go                       ← NEW (git 도구)
│   ├── diagnostics.go               ← NEW (LSP 진단)
│   ├── file.go                      ← KEEP
│   ├── shell.go                     ← KEEP
│   └── search.go                    ← KEEP
├── session/
│   ├── db.go                        ← NEW (SQLite 초기화 + 마이그레이션)
│   ├── store.go                     ← NEW (세션 CRUD)
│   ├── memory.go                    ← NEW (장기 메모리)
│   └── usage.go                     ← NEW (토큰/비용 추적)
├── agents/
│   └── auto.go                      ← NEW (자율 모드)
├── mcp/
│   ├── client.go                    ← NEW (MCP JSON-RPC 2.0)
│   ├── transport_stdio.go           ← NEW
│   └── transport_sse.go             ← NEW
├── app/app.go                       ← MODIFY (슬래시 커맨드 확장, 세션 통합)
├── ui/
│   ├── styles.go                    ← MODIFY (hanimo 색상)
│   ├── super.go                     ← MODIFY (hanimo 로고)
│   ├── chat.go                      ← MODIFY (상태바에 프로바이더 표시)
│   ├── tabbar.go                    ← KEEP
│   └── palette.go                   ← NEW (커맨드 팔레트 Ctrl+K)
└── knowledge/                       ← KEEP (store.go, injector.go, extractor.go)

db/schema.sql                        ← NEW
knowledge/                           ← MODIFY (BXM 제거, 오픈소스 문서 추가)
knowledge.go                         ← KEEP (embed.FS)
go.mod                               ← MODIFY (모듈명 + 신규 의존성)
Makefile                             ← MODIFY (hanimo 빌드 타겟)
README.md                            ← REWRITE
```

---

## Phase 1: MVP

### Task 1: TECHAI_CODE 포크 + 리네임

**Files:**
- Copy: TECHAI_CODE → hanimo (Go 소스만)
- Modify: `go.mod`, `cmd/hanimo/main.go`, `internal/config/config.go`, `Makefile`

- [ ] **Step 1: TECHAI_CODE Go 소스를 hanimo에 복사**

```bash
cd /Users/jiwonkim/Desktop/kimjiwon/hanimo

# TECHAI_CODE에서 Go 소스만 복사
cp -r /Users/jiwonkim/Desktop/kimjiwon/TECHAI_CODE/cmd .
cp -r /Users/jiwonkim/Desktop/kimjiwon/TECHAI_CODE/internal .
cp -r /Users/jiwonkim/Desktop/kimjiwon/TECHAI_CODE/knowledge .
cp /Users/jiwonkim/Desktop/kimjiwon/TECHAI_CODE/knowledge.go .
cp /Users/jiwonkim/Desktop/kimjiwon/TECHAI_CODE/go.mod .
cp /Users/jiwonkim/Desktop/kimjiwon/TECHAI_CODE/go.sum .
cp /Users/jiwonkim/Desktop/kimjiwon/TECHAI_CODE/Makefile .
```

- [ ] **Step 2: go.mod 모듈명 변경**

`go.mod`의 모듈명을 변경:
```
module github.com/kimjiwon/tgc
```
→
```
module github.com/flykimjiwon/hanimo
```

- [ ] **Step 3: 전체 import 경로 일괄 변경**

```bash
cd /Users/jiwonkim/Desktop/kimjiwon/hanimo
find . -name "*.go" -exec sed -i '' 's|github.com/kimjiwon/tgc|github.com/flykimjiwon/hanimo|g' {} +
```

- [ ] **Step 4: cmd/tgc → cmd/hanimo 리네임**

```bash
mv cmd/tgc cmd/hanimo
```

`cmd/hanimo/main.go`에서 변경:
```go
// 변경 전
var version = "dev"

// 변경 후
var version = "dev"
var appName = "hanimo"
```

- [ ] **Step 5: config 경로 변경**

`internal/config/config.go`에서:
```go
// 변경 전
ConfigDirName = ".tgc"

// 변경 후
ConfigDirName = ".hanimo"
```

환경변수 프리픽스:
```go
// 변경 전: TGC_API_BASE_URL, TGC_API_KEY, TGC_MODEL_SUPER, TGC_MODEL_DEV
// 변경 후: HANIMO_API_BASE_URL, HANIMO_API_KEY, HANIMO_MODEL_SUPER, HANIMO_MODEL_DEV
```

프로젝트 컨텍스트 파일:
```go
// 변경 전: .techai.md
// 변경 후: .hanimo.md
```

- [ ] **Step 6: Makefile 수정**

```makefile
BINARY_NAME=hanimo
MODULE=github.com/flykimjiwon/hanimo
CMD_DIR=cmd/hanimo

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) ./$(CMD_DIR)

build-all:
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/hanimo-darwin-arm64 ./$(CMD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/hanimo-darwin-amd64 ./$(CMD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/hanimo-linux-amd64 ./$(CMD_DIR)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/hanimo-linux-arm64 ./$(CMD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/hanimo-windows-amd64.exe ./$(CMD_DIR)
```

- [ ] **Step 7: 빌드 확인**

```bash
cd /Users/jiwonkim/Desktop/kimjiwon/hanimo
go build -o hanimo ./cmd/hanimo
./hanimo --version
```
Expected: `hanimo dev`

- [ ] **Step 8: 커밋**

```bash
git add -A
git commit -m "feat: fork TECHAI_CODE as hanimo Go base

Copied Go source from TECHAI_CODE, renamed module to
github.com/flykimjiwon/hanimo, updated all paths and config."
```

---

### Task 2: hanimo 브랜딩 (로고, 색상, UI)

**Files:**
- Modify: `internal/ui/styles.go`, `internal/ui/super.go`, `internal/ui/chat.go`

- [ ] **Step 1: 기존 hanimo 로고 확인**

```bash
grep -A 30 "logo\|Logo\|ASCII" /Users/jiwonkim/Desktop/kimjiwon/hanimo/_legacy_ts/src/tui/*.tsx 2>/dev/null | head -60
```

TS hanimo의 로고 ASCII art와 색상을 확인하고 가져온다.

- [ ] **Step 2: styles.go 색상 변경**

`internal/ui/styles.go`에서 hanimo 브랜딩 색상으로 변경:
```go
// hanimo 색상 팔레트 (기존 TS hanimo에서 가져옴)
// 정확한 색상은 Step 1에서 확인 후 적용
// 기본 방향: TECHAI_CODE의 블루 계열 → hanimo 고유 색상
```

- [ ] **Step 3: super.go 로고 교체**

`internal/ui/super.go`의 `RenderLogo()` 함수에서 TECHAI → hanimo 로고로 교체.

- [ ] **Step 4: chat.go 상태바 수정**

`internal/ui/chat.go`의 `RenderStatusBar()`에서:
- "택가이코드" → "hanimo" 텍스트 변경
- 프로바이더명 표시 추가 (예: `ollama/qwen3:8b`)

- [ ] **Step 5: 모드 이름 변경**

`internal/llm/prompt.go`에서:
```go
// 변경 전
var Modes = [ModeCount]ModeInfo{
    {ID: "super", Name: "슈퍼택가이", ...},
    {ID: "dev", Name: "개발", ...},
    {ID: "plan", Name: "플랜", ...},
}

// 변경 후
var Modes = [ModeCount]ModeInfo{
    {ID: "super", Name: "Super", ...},
    {ID: "dev", Name: "Dev", ...},
    {ID: "plan", Name: "Plan", ...},
}
```

- [ ] **Step 6: 빌드 + 시각 확인**

```bash
go build -o hanimo ./cmd/hanimo && ./hanimo
```
로고, 색상, 모드명이 hanimo로 표시되는지 확인.

- [ ] **Step 7: 커밋**

```bash
git add internal/ui/ internal/llm/prompt.go
git commit -m "feat: apply hanimo branding (logo, colors, mode names)"
```

---

### Task 3: 프로바이더 레지스트리 + OpenAI-compatible

**Files:**
- Create: `internal/llm/providers/registry.go`, `internal/llm/providers/openai_compat.go`
- Modify: `internal/llm/client.go`, `internal/config/config.go`

- [ ] **Step 1: Provider 인터페이스 정의**

Create `internal/llm/providers/registry.go`:
```go
package providers

import "context"

type ChatRequest struct {
    Model       string
    Messages    []Message
    Tools       []ToolDef
    Temperature float64
    MaxTokens   int
}

type Message struct {
    Role       string      // system/user/assistant/tool
    Content    string
    ToolCalls  []ToolCall  // assistant 메시지
    ToolCallID string      // tool 응답 메시지
}

type ToolCall struct {
    ID        string
    Name      string
    Arguments string // JSON
}

type ToolDef struct {
    Name        string
    Description string
    Parameters  map[string]interface{} // JSON Schema
}

type ChatChunk struct {
    Content   string
    ToolCalls []ToolCall
    Done      bool
    Error     error
    Usage     *Usage
}

type Usage struct {
    PromptTokens     int
    CompletionTokens int
    TotalTokens      int
}

type ModelInfo struct {
    ID          string
    DisplayName string
    Provider    string
    ContextWindow int
    SupportsTools bool
}

type Provider interface {
    Name() string
    Chat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
    ListModels() ([]ModelInfo, error)
    SupportsTools() bool
}

var registry = map[string]func(baseURL, apiKey string) Provider{}

func Register(name string, factory func(baseURL, apiKey string) Provider) {
    registry[name] = factory
}

func Get(name, baseURL, apiKey string) (Provider, error) {
    factory, ok := registry[name]
    if !ok {
        // fallback to openai-compatible
        factory = registry["openai"]
    }
    return factory(baseURL, apiKey), nil
}

func init() {
    // 기본 프로바이더 등록은 각 파일의 init()에서 수행
}
```

- [ ] **Step 2: OpenAI-compatible 프로바이더 구현**

Create `internal/llm/providers/openai_compat.go`:
```go
package providers

import (
    "context"
    openai "github.com/sashabaranov/go-openai"
)

// 하나의 구현으로 8+ 프로바이더 커버
var openaiCompatProviders = map[string]string{
    "openai":     "https://api.openai.com/v1",
    "novita":     "https://api.novita.ai/v1",
    "openrouter": "https://openrouter.ai/api/v1",
    "deepseek":   "https://api.deepseek.com/v1",
    "groq":       "https://api.groq.com/openai/v1",
    "together":   "https://api.together.xyz/v1",
    "fireworks":  "https://api.fireworks.ai/inference/v1",
    "mistral":    "https://api.mistral.ai/v1",
}

type OpenAICompatProvider struct {
    name    string
    client  *openai.Client
    baseURL string
}

func NewOpenAICompat(name, baseURL, apiKey string) Provider {
    // If baseURL empty, use default for known providers
    if baseURL == "" {
        if defaultURL, ok := openaiCompatProviders[name]; ok {
            baseURL = defaultURL
        }
    }
    cfg := openai.DefaultConfig(apiKey)
    cfg.BaseURL = baseURL
    return &OpenAICompatProvider{
        name:   name,
        client: openai.NewClientWithConfig(cfg),
        baseURL: baseURL,
    }
}

func (p *OpenAICompatProvider) Name() string { return p.name }
func (p *OpenAICompatProvider) SupportsTools() bool { return true }

func (p *OpenAICompatProvider) Chat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error) {
    // Convert ChatRequest → openai.ChatCompletionRequest
    // Stream via CreateChatCompletionStream
    // Yield ChatChunk via channel
    // (기존 client.go의 StreamChat 로직 이동)
}

func (p *OpenAICompatProvider) ListModels() ([]ModelInfo, error) {
    // GET /v1/models API 호출
}

func init() {
    for name := range openaiCompatProviders {
        n := name // capture
        Register(n, func(baseURL, apiKey string) Provider {
            return NewOpenAICompat(n, baseURL, apiKey)
        })
    }
    // custom + vllm + lmstudio도 openai-compatible
    Register("vllm", func(baseURL, apiKey string) Provider {
        return NewOpenAICompat("vllm", baseURL, apiKey)
    })
    Register("lmstudio", func(baseURL, apiKey string) Provider {
        return NewOpenAICompat("lmstudio", baseURL, apiKey)
    })
    Register("custom", func(baseURL, apiKey string) Provider {
        return NewOpenAICompat("custom", baseURL, apiKey)
    })
}
```

- [ ] **Step 3: config.go에 프로바이더 설정 추가**

```go
type ProviderConfig struct {
    APIKey  string `yaml:"api_key,omitempty"`
    BaseURL string `yaml:"base_url,omitempty"`
}

type Config struct {
    Default struct {
        Provider string `yaml:"provider"` // default: "ollama"
        Model    string `yaml:"model"`    // default: "qwen3:8b"
    } `yaml:"default"`
    Providers map[string]ProviderConfig `yaml:"providers"`
    // 기존 API/Models 필드는 하위호환용 유지
}
```

- [ ] **Step 4: client.go를 Provider 인터페이스 사용으로 전환**

기존 `Client` struct의 `StreamChat()`를 Provider의 `Chat()`으로 위임하도록 변경. 기존 직접 호출 코드를 Provider 팩토리 경유로 수정.

- [ ] **Step 5: 빌드 + Ollama 테스트**

```bash
go build -o hanimo ./cmd/hanimo
./hanimo --provider ollama --model qwen3:8b
```
Expected: Ollama에 연결되어 대화 가능

- [ ] **Step 6: 커밋**

```bash
git add internal/llm/providers/ internal/config/ internal/llm/client.go
git commit -m "feat: add provider registry with OpenAI-compatible support

Covers OpenAI, Novita, OpenRouter, DeepSeek, Groq, Together,
Fireworks, Mistral, vLLM, LM Studio, and custom endpoints."
```

---

### Task 4: Anthropic + Google 네이티브 프로바이더

**Files:**
- Create: `internal/llm/providers/anthropic.go`, `internal/llm/providers/google.go`
- Modify: `go.mod`

- [ ] **Step 1: 의존성 추가**

```bash
go get github.com/anthropics/anthropic-sdk-go
go get github.com/google/generative-ai-go
```

- [ ] **Step 2: Anthropic 프로바이더 구현**

Create `internal/llm/providers/anthropic.go`:
```go
package providers

import (
    "context"
    "github.com/anthropics/anthropic-sdk-go"
)

type AnthropicProvider struct {
    client *anthropic.Client
}

func NewAnthropic(_, apiKey string) Provider {
    client := anthropic.NewClient(anthropic.WithAPIKey(apiKey))
    return &AnthropicProvider{client: client}
}

func (p *AnthropicProvider) Name() string { return "anthropic" }
func (p *AnthropicProvider) SupportsTools() bool { return true }

func (p *AnthropicProvider) Chat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error) {
    // Convert ChatRequest → anthropic.MessageCreateParams
    // Use Messages.Stream() for SSE streaming
    // Map content blocks + tool_use blocks → ChatChunk
}

func (p *AnthropicProvider) ListModels() ([]ModelInfo, error) {
    return []ModelInfo{
        {ID: "claude-sonnet-4-20250514", DisplayName: "Claude Sonnet 4", Provider: "anthropic", ContextWindow: 200000, SupportsTools: true},
        {ID: "claude-haiku-4-20250414", DisplayName: "Claude Haiku 4", Provider: "anthropic", ContextWindow: 200000, SupportsTools: true},
    }, nil
}

func init() {
    Register("anthropic", func(baseURL, apiKey string) Provider {
        return NewAnthropic(baseURL, apiKey)
    })
}
```

- [ ] **Step 3: Google Gemini 프로바이더 구현**

Create `internal/llm/providers/google.go`:
```go
package providers

import (
    "context"
    "github.com/google/generative-ai-go/genai"
    "google.golang.org/api/option"
)

type GoogleProvider struct {
    apiKey string
}

func NewGoogle(_, apiKey string) Provider {
    return &GoogleProvider{apiKey: apiKey}
}

func (p *GoogleProvider) Name() string { return "google" }
func (p *GoogleProvider) SupportsTools() bool { return true }

func (p *GoogleProvider) Chat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error) {
    // Create genai.NewClient with apiKey
    // Convert ChatRequest → genai.Content messages
    // Use GenerateContentStream for streaming
    // Map genai.GenerateContentResponse → ChatChunk
}

func (p *GoogleProvider) ListModels() ([]ModelInfo, error) {
    return []ModelInfo{
        {ID: "gemini-2.5-flash", DisplayName: "Gemini 2.5 Flash", Provider: "google", ContextWindow: 1000000, SupportsTools: true},
        {ID: "gemini-2.5-pro", DisplayName: "Gemini 2.5 Pro", Provider: "google", ContextWindow: 1000000, SupportsTools: true},
    }, nil
}

func init() {
    Register("google", func(baseURL, apiKey string) Provider {
        return NewGoogle(baseURL, apiKey)
    })
}
```

- [ ] **Step 4: Ollama 네이티브 프로바이더**

Create `internal/llm/providers/ollama.go`:
```go
package providers

// Ollama는 OpenAI-compatible이지만 모델 자동 탐지를 위해 네이티브 구현
type OllamaProvider struct {
    OpenAICompatProvider
}

func NewOllama(baseURL, _ string) Provider {
    if baseURL == "" {
        baseURL = "http://localhost:11434"
    }
    p := NewOpenAICompat("ollama", baseURL+"/v1", "ollama").(*OpenAICompatProvider)
    return &OllamaProvider{*p}
}

func (p *OllamaProvider) ListModels() ([]ModelInfo, error) {
    // GET http://localhost:11434/api/tags
    // Parse Ollama model list response
    // Return as []ModelInfo
}

func init() {
    Register("ollama", func(baseURL, apiKey string) Provider {
        return NewOllama(baseURL, apiKey)
    })
}
```

- [ ] **Step 5: 빌드 확인**

```bash
go build -o hanimo ./cmd/hanimo
```

- [ ] **Step 6: 커밋**

```bash
git add internal/llm/providers/ go.mod go.sum
git commit -m "feat: add Anthropic, Google, Ollama native providers"
```

---

### Task 5: SQLite 통합 (세션, 메모리, 사용량)

**Files:**
- Create: `db/schema.sql`, `internal/session/db.go`, `internal/session/store.go`, `internal/session/memory.go`, `internal/session/usage.go`
- Modify: `go.mod`

- [ ] **Step 1: SQLite 의존성 추가**

```bash
go get modernc.org/sqlite
```

- [ ] **Step 2: 스키마 파일 생성**

Create `db/schema.sql`:
```sql
CREATE TABLE IF NOT EXISTS sessions (
    id          TEXT PRIMARY KEY,
    name        TEXT,
    project_dir TEXT NOT NULL,
    provider    TEXT NOT NULL,
    model       TEXT NOT NULL,
    mode        TEXT DEFAULT 'super',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id  TEXT NOT NULL REFERENCES sessions(id),
    role        TEXT NOT NULL,
    content     TEXT NOT NULL,
    tool_calls  TEXT,
    tool_result TEXT,
    tokens_in   INTEGER DEFAULT 0,
    tokens_out  INTEGER DEFAULT 0,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS memories (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    project_dir TEXT NOT NULL,
    key         TEXT NOT NULL,
    value       TEXT NOT NULL,
    source      TEXT DEFAULT 'auto',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_dir, key)
);

CREATE TABLE IF NOT EXISTS usage_log (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id  TEXT NOT NULL REFERENCES sessions(id),
    provider    TEXT NOT NULL,
    model       TEXT NOT NULL,
    tokens_in   INTEGER NOT NULL,
    tokens_out  INTEGER NOT NULL,
    cost_usd    REAL DEFAULT 0,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS mcp_servers (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT UNIQUE NOT NULL,
    transport   TEXT NOT NULL,
    command     TEXT,
    args        TEXT,
    url         TEXT,
    enabled     BOOLEAN DEFAULT 1,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_messages_session ON messages(session_id);
CREATE INDEX IF NOT EXISTS idx_memories_project ON memories(project_dir);
CREATE INDEX IF NOT EXISTS idx_usage_session ON usage_log(session_id);
```

- [ ] **Step 3: DB 초기화 모듈**

Create `internal/session/db.go`:
```go
package session

import (
    "database/sql"
    "embed"
    "os"
    "path/filepath"
    _ "modernc.org/sqlite"
)

//go:embed ../../db/schema.sql
var schemaSQL string

var db *sql.DB

func InitDB(configDir string) error {
    dbPath := filepath.Join(configDir, "hanimo.db")
    os.MkdirAll(configDir, 0755)

    var err error
    db, err = sql.Open("sqlite", dbPath)
    if err != nil {
        return err
    }

    // WAL 모드 (동시 읽기 성능)
    db.Exec("PRAGMA journal_mode=WAL")
    db.Exec("PRAGMA foreign_keys=ON")

    _, err = db.Exec(schemaSQL)
    return err
}

func CloseDB() {
    if db != nil {
        db.Close()
    }
}

func DB() *sql.DB {
    return db
}
```

- [ ] **Step 4: 세션 CRUD**

Create `internal/session/store.go`:
```go
package session

import (
    "database/sql"
    "encoding/json"
    "time"
    "github.com/google/uuid"
)

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

type StoredMessage struct {
    ID         int
    SessionID  string
    Role       string
    Content    string
    ToolCalls  string // JSON
    ToolResult string // JSON
    TokensIn   int
    TokensOut  int
    CreatedAt  time.Time
}

func CreateSession(projectDir, provider, model, mode string) (*Session, error) {
    s := &Session{
        ID: uuid.New().String(),
        ProjectDir: projectDir,
        Provider: provider,
        Model: model,
        Mode: mode,
    }
    _, err := db.Exec(
        "INSERT INTO sessions (id, project_dir, provider, model, mode) VALUES (?, ?, ?, ?, ?)",
        s.ID, s.ProjectDir, s.Provider, s.Model, s.Mode,
    )
    return s, err
}

func SaveMessage(sessionID, role, content string, toolCalls interface{}, tokensIn, tokensOut int) error {
    tc, _ := json.Marshal(toolCalls)
    _, err := db.Exec(
        "INSERT INTO messages (session_id, role, content, tool_calls, tokens_in, tokens_out) VALUES (?, ?, ?, ?, ?, ?)",
        sessionID, role, content, string(tc), tokensIn, tokensOut,
    )
    // update session.updated_at
    db.Exec("UPDATE sessions SET updated_at = CURRENT_TIMESTAMP WHERE id = ?", sessionID)
    return err
}

func LoadSession(id string) (*Session, []StoredMessage, error) {
    // SELECT session + messages ORDER BY created_at
}

func ListSessions(projectDir string, limit int) ([]Session, error) {
    // SELECT sessions WHERE project_dir = ? ORDER BY updated_at DESC LIMIT ?
}

func SearchSessions(keyword string) ([]Session, error) {
    // SELECT DISTINCT sessions JOIN messages WHERE content LIKE '%keyword%'
}

func NameSession(id, name string) error {
    _, err := db.Exec("UPDATE sessions SET name = ? WHERE id = ?", name, id)
    return err
}

func ForkSession(sourceID string) (*Session, error) {
    // 1. Load source session + messages
    // 2. Create new session with new ID
    // 3. Copy all messages
}
```

- [ ] **Step 5: 메모리 모듈**

Create `internal/session/memory.go`:
```go
package session

type Memory struct {
    ID         int
    ProjectDir string
    Key        string
    Value      string
    Source     string
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

func SaveMemory(projectDir, key, value, source string) error {
    _, err := db.Exec(
        `INSERT INTO memories (project_dir, key, value, source)
         VALUES (?, ?, ?, ?)
         ON CONFLICT(project_dir, key) DO UPDATE SET value = ?, source = ?, updated_at = CURRENT_TIMESTAMP`,
        projectDir, key, value, source, value, source,
    )
    return err
}

func LoadMemories(projectDir string) ([]Memory, error) {
    // SELECT * FROM memories WHERE project_dir = ? ORDER BY updated_at DESC
}

func DeleteMemory(projectDir, key string) error {
    _, err := db.Exec("DELETE FROM memories WHERE project_dir = ? AND key = ?", projectDir, key)
    return err
}

func FormatMemoriesForPrompt(memories []Memory) string {
    // Format as markdown section for system prompt injection
    // ## Project Memory
    // - key: value
    // - key2: value2
}
```

- [ ] **Step 6: 사용량 추적**

Create `internal/session/usage.go`:
```go
package session

type UsageEntry struct {
    Provider string
    Model    string
    TokensIn  int
    TokensOut int
    CostUSD   float64
}

// 모델별 가격 (per 1M tokens)
var pricing = map[string][2]float64{
    // [input, output]
    "gpt-4o":              {2.50, 10.00},
    "gpt-4o-mini":         {0.15, 0.60},
    "claude-sonnet-4":     {3.00, 15.00},
    "claude-haiku-4":      {0.80, 4.00},
    "gemini-2.5-flash":    {0.15, 0.60},
    "deepseek-chat":       {0.14, 0.28},
    "qwen3:8b":            {0, 0}, // local
}

func LogUsage(sessionID, provider, model string, tokensIn, tokensOut int) error {
    cost := calculateCost(model, tokensIn, tokensOut)
    _, err := db.Exec(
        "INSERT INTO usage_log (session_id, provider, model, tokens_in, tokens_out, cost_usd) VALUES (?, ?, ?, ?, ?, ?)",
        sessionID, provider, model, tokensIn, tokensOut, cost,
    )
    return err
}

func GetSessionUsage(sessionID string) (totalIn, totalOut int, totalCost float64, err error) {
    // SELECT SUM(tokens_in), SUM(tokens_out), SUM(cost_usd) FROM usage_log WHERE session_id = ?
}

func GetTotalUsage() (totalIn, totalOut int, totalCost float64, err error) {
    // SELECT SUM across all sessions
}

func calculateCost(model string, tokensIn, tokensOut int) float64 {
    p, ok := pricing[model]
    if !ok { return 0 }
    return (float64(tokensIn) * p[0] + float64(tokensOut) * p[1]) / 1_000_000
}
```

- [ ] **Step 7: main.go에 DB 초기화 추가**

`cmd/hanimo/main.go`에서:
```go
func main() {
    cfg, _ := config.Load()

    // SQLite 초기화
    session.InitDB(config.ConfigDir())
    defer session.CloseDB()

    // ... 기존 로직
}
```

- [ ] **Step 8: 빌드 확인**

```bash
go build -o hanimo ./cmd/hanimo
./hanimo
# ~/.hanimo/hanimo.db 생성 확인
ls -la ~/.hanimo/
```

- [ ] **Step 9: 커밋**

```bash
git add db/ internal/session/ cmd/hanimo/main.go go.mod go.sum
git commit -m "feat: add SQLite session, memory, and usage tracking"
```

---

### Task 6: Hash-Anchored 편집

**Files:**
- Create: `internal/tools/hashline.go`
- Modify: `internal/tools/registry.go`

- [ ] **Step 1: hashline.go 구현**

Create `internal/tools/hashline.go`:
```go
package tools

import (
    "crypto/md5"
    "fmt"
    "os"
    "strings"
)

// HashlineRead reads a file and tags each line with MD5 hash
func HashlineRead(path string) (string, error) {
    absPath, err := filepath.Abs(path)
    if err != nil { return "", err }

    data, err := os.ReadFile(absPath)
    if err != nil { return "", err }

    lines := strings.Split(string(data), "\n")
    var result strings.Builder
    for i, line := range lines {
        hash := fmt.Sprintf("%x", md5.Sum([]byte(line)))[:4]
        result.WriteString(fmt.Sprintf("%d#%s| %s\n", i+1, hash, line))
    }
    return result.String(), nil
}

// HashlineEdit edits a file using hash anchors for verification
func HashlineEdit(path, startAnchor, endAnchor, newContent string) (string, error) {
    absPath, err := filepath.Abs(path)
    if err != nil { return "", err }

    data, err := os.ReadFile(absPath)
    if err != nil { return "", err }

    lines := strings.Split(string(data), "\n")

    startLine, startHash, err := parseAnchor(startAnchor)
    if err != nil { return "", fmt.Errorf("invalid startAnchor: %v", err) }

    endLine, endHash, err := parseAnchor(endAnchor)
    if err != nil { return "", fmt.Errorf("invalid endAnchor: %v", err) }

    // Verify hashes
    if startLine < 1 || startLine > len(lines) {
        return "", fmt.Errorf("startLine %d out of range (1-%d)", startLine, len(lines))
    }
    actualStartHash := fmt.Sprintf("%x", md5.Sum([]byte(lines[startLine-1])))[:4]
    if actualStartHash != startHash {
        return "", fmt.Errorf("hash mismatch at line %d: expected %s, got %s. File changed — re-read with hashline_read", startLine, startHash, actualStartHash)
    }

    if endLine < 1 || endLine > len(lines) {
        return "", fmt.Errorf("endLine %d out of range (1-%d)", endLine, len(lines))
    }
    actualEndHash := fmt.Sprintf("%x", md5.Sum([]byte(lines[endLine-1])))[:4]
    if actualEndHash != endHash {
        return "", fmt.Errorf("hash mismatch at line %d: expected %s, got %s. File changed — re-read with hashline_read", endLine, endHash, actualEndHash)
    }

    // Replace lines[startLine-1:endLine] with newContent
    newLines := strings.Split(newContent, "\n")
    result := make([]string, 0, len(lines)-endLine+startLine-1+len(newLines))
    result = append(result, lines[:startLine-1]...)
    result = append(result, newLines...)
    result = append(result, lines[endLine:]...)

    err = os.WriteFile(absPath, []byte(strings.Join(result, "\n")), 0644)
    if err != nil { return "", err }

    return fmt.Sprintf("OK: replaced lines %d-%d (%d lines → %d lines)", startLine, endLine, endLine-startLine+1, len(newLines)), nil
}

func parseAnchor(anchor string) (int, string, error) {
    // Parse "3#e4d9" → line=3, hash="e4d9"
    parts := strings.SplitN(anchor, "#", 2)
    if len(parts) != 2 { return 0, "", fmt.Errorf("invalid anchor format: %s", anchor) }
    line, err := strconv.Atoi(parts[0])
    if err != nil { return 0, "", err }
    return line, parts[1], nil
}
```

- [ ] **Step 2: registry.go에 hashline 도구 등록**

```go
// AllTools()에 추가:
{Name: "hashline_read", Description: "Read file with hash-tagged line numbers for safe editing", Parameters: ...},
{Name: "hashline_edit", Description: "Edit file using hash anchors to verify lines haven't changed", Parameters: ...},

// Execute()에 분기 추가:
case "hashline_read": ...
case "hashline_edit": ...
```

- [ ] **Step 3: 빌드 확인**

```bash
go build -o hanimo ./cmd/hanimo
```

- [ ] **Step 4: 커밋**

```bash
git add internal/tools/hashline.go internal/tools/registry.go
git commit -m "feat: add hash-anchored editing (hashline_read, hashline_edit)

MD5 hash per line prevents stale-edit corruption.
Mismatched hash returns error with re-read hint."
```

---

### Task 7: Smart Compaction

**Files:**
- Create: `internal/llm/compaction.go`
- Modify: `internal/app/app.go`

- [ ] **Step 1: compaction.go 구현**

Create `internal/llm/compaction.go`:
```go
package llm

import (
    "context"
    "fmt"
    "strings"
    openai "github.com/sashabaranov/go-openai"
)

const (
    snipThreshold  = 40   // 메시지 수 기준
    microMaxLen    = 4000 // 단일 메시지 최대 길이
    summaryPrompt  = `Summarize the conversation so far in a concise format.
Preserve: task goal, completed work, current state, important decisions.
Discard: tool outputs, intermediate reasoning, duplicate information.
Output in the same language as the conversation.`
)

// Compact reduces conversation history to fit within context window
func Compact(messages []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
    if len(messages) < snipThreshold {
        return messages
    }

    // Stage 1: Snip old tool results
    result := snipToolResults(messages)

    // Stage 2: Micro-truncate long messages
    result = microTruncate(result)

    return result
}

// CompactWithLLM uses LLM to summarize if still too large
func CompactWithLLM(ctx context.Context, client *Client, model string, messages []openai.ChatCompletionMessage, maxTokens int) []openai.ChatCompletionMessage {
    result := Compact(messages)

    // Estimate token count
    totalTokens := estimateTokens(result)
    if totalTokens <= maxTokens {
        return result
    }

    // Stage 3: LLM summary of older messages
    // Keep system prompt (idx 0) + last 10 messages
    // Summarize everything in between
    keepLast := 10
    if len(result) <= keepLast+1 {
        return result
    }

    toSummarize := result[1 : len(result)-keepLast]
    summary, err := client.Chat(ctx, model, []openai.ChatCompletionMessage{
        {Role: "system", Content: summaryPrompt},
        {Role: "user", Content: formatMessagesForSummary(toSummarize)},
    })
    if err != nil {
        return result // 실패 시 원본 유지
    }

    // Reconstruct: system + summary + last N
    compacted := make([]openai.ChatCompletionMessage, 0, keepLast+2)
    compacted = append(compacted, result[0]) // system
    compacted = append(compacted, openai.ChatCompletionMessage{
        Role: "system", Content: "[Previous conversation summary]\n" + summary,
    })
    compacted = append(compacted, result[len(result)-keepLast:]...)

    return compacted
}

func snipToolResults(msgs []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
    // 마지막 10개 제외, tool role 메시지의 content를 축약
    result := make([]openai.ChatCompletionMessage, len(msgs))
    copy(result, msgs)
    cutoff := len(msgs) - 10
    for i := 0; i < cutoff; i++ {
        if result[i].Role == "tool" && len(result[i].Content) > 200 {
            lines := strings.Count(result[i].Content, "\n")
            result[i].Content = fmt.Sprintf("[snipped: %d lines]", lines)
        }
    }
    return result
}

func microTruncate(msgs []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
    result := make([]openai.ChatCompletionMessage, len(msgs))
    copy(result, msgs)
    for i := range result {
        if len(result[i].Content) > microMaxLen {
            half := microMaxLen / 2
            result[i].Content = result[i].Content[:half] + "\n...[truncated]...\n" + result[i].Content[len(result[i].Content)-half:]
        }
    }
    return result
}

func estimateTokens(msgs []openai.ChatCompletionMessage) int {
    total := 0
    for _, m := range msgs {
        total += len(m.Content) / 4
    }
    return total
}

func formatMessagesForSummary(msgs []openai.ChatCompletionMessage) string {
    var sb strings.Builder
    for _, m := range msgs {
        sb.WriteString(fmt.Sprintf("[%s]: %s\n\n", m.Role, m.Content))
    }
    return sb.String()
}
```

- [ ] **Step 2: app.go에 compaction 적용**

`internal/app/app.go`의 `sendMessage()` 함수에서 스트림 시작 전에:
```go
func (m *Model) sendMessage(input string) tea.Cmd {
    // ... 기존 메시지 추가 로직

    // 컴팩션 적용
    m.history = llm.Compact(m.history)

    // ... 기존 스트림 시작 로직
}
```

- [ ] **Step 3: 빌드 확인**

```bash
go build -o hanimo ./cmd/hanimo
```

- [ ] **Step 4: 커밋**

```bash
git add internal/llm/compaction.go internal/app/app.go
git commit -m "feat: add smart compaction (snip → micro → LLM summary)

3-stage context compression for long conversations.
Preserves task goal and recent context."
```

---

### Task 8: MCP 클라이언트

**Files:**
- Create: `internal/mcp/client.go`, `internal/mcp/transport_stdio.go`, `internal/mcp/transport_sse.go`

- [ ] **Step 1: MCP client.go — JSON-RPC 2.0 코어**

Create `internal/mcp/client.go`:
```go
package mcp

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "sync/atomic"
)

type Request struct {
    JSONRPC string      `json:"jsonrpc"`
    ID      int64       `json:"id"`
    Method  string      `json:"method"`
    Params  interface{} `json:"params,omitempty"`
}

type Response struct {
    JSONRPC string          `json:"jsonrpc"`
    ID      int64           `json:"id"`
    Result  json.RawMessage `json:"result,omitempty"`
    Error   *RPCError       `json:"error,omitempty"`
}

type RPCError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

type Tool struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    InputSchema map[string]interface{} `json:"inputSchema"`
}

type Transport interface {
    Send(data []byte) error
    Receive() ([]byte, error)
    Close() error
}

type Client struct {
    transport Transport
    nextID    atomic.Int64
    mu        sync.Mutex
    tools     []Tool
}

func NewClient(transport Transport) *Client {
    return &Client{transport: transport}
}

func (c *Client) Initialize(ctx context.Context) error {
    // Send initialize request
    // Receive capabilities
    // Send initialized notification
    return nil
}

func (c *Client) ListTools(ctx context.Context) ([]Tool, error) {
    resp, err := c.call("tools/list", nil)
    if err != nil { return nil, err }
    var result struct { Tools []Tool `json:"tools"` }
    json.Unmarshal(resp, &result)
    c.tools = result.Tools
    return result.Tools, nil
}

func (c *Client) CallTool(ctx context.Context, name string, args map[string]interface{}) (string, error) {
    params := map[string]interface{}{
        "name":      name,
        "arguments": args,
    }
    resp, err := c.call("tools/call", params)
    if err != nil { return "", err }
    var result struct {
        Content []struct {
            Type string `json:"type"`
            Text string `json:"text"`
        } `json:"content"`
    }
    json.Unmarshal(resp, &result)
    if len(result.Content) > 0 {
        return result.Content[0].Text, nil
    }
    return "", nil
}

func (c *Client) Close() error {
    return c.transport.Close()
}

func (c *Client) call(method string, params interface{}) (json.RawMessage, error) {
    c.mu.Lock()
    defer c.mu.Unlock()
    id := c.nextID.Add(1)
    req := Request{JSONRPC: "2.0", ID: id, Method: method, Params: params}
    data, _ := json.Marshal(req)
    if err := c.transport.Send(data); err != nil { return nil, err }
    respData, err := c.transport.Receive()
    if err != nil { return nil, err }
    var resp Response
    json.Unmarshal(respData, &resp)
    if resp.Error != nil {
        return nil, fmt.Errorf("MCP error %d: %s", resp.Error.Code, resp.Error.Message)
    }
    return resp.Result, nil
}
```

- [ ] **Step 2: stdio 트랜스포트**

Create `internal/mcp/transport_stdio.go`:
```go
package mcp

import (
    "bufio"
    "io"
    "os/exec"
)

type StdioTransport struct {
    cmd    *exec.Cmd
    stdin  io.WriteCloser
    reader *bufio.Reader
}

func NewStdioTransport(command string, args []string) (*StdioTransport, error) {
    cmd := exec.Command(command, args...)
    stdin, _ := cmd.StdinPipe()
    stdout, _ := cmd.StdoutPipe()
    if err := cmd.Start(); err != nil {
        return nil, err
    }
    return &StdioTransport{
        cmd:    cmd,
        stdin:  stdin,
        reader: bufio.NewReader(stdout),
    }, nil
}

func (t *StdioTransport) Send(data []byte) error {
    // Write Content-Length header + body
    header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
    _, err := t.stdin.Write([]byte(header))
    if err != nil { return err }
    _, err = t.stdin.Write(data)
    return err
}

func (t *StdioTransport) Receive() ([]byte, error) {
    // Read Content-Length header, then body
    // Parse header line for content length
    // Read exact bytes
}

func (t *StdioTransport) Close() error {
    t.stdin.Close()
    return t.cmd.Wait()
}
```

- [ ] **Step 3: SSE 트랜스포트** (스텁)

Create `internal/mcp/transport_sse.go`:
```go
package mcp

// SSETransport implements Transport over HTTP Server-Sent Events
type SSETransport struct {
    url string
    // Phase 2에서 완전 구현
}

// 기본 구조만 정의, 실제 구현은 Phase 2
```

- [ ] **Step 4: MCP 도구를 hanimo 도구 레지스트리에 통합**

`internal/tools/registry.go`에서 MCP 도구를 동적으로 등록:
```go
func RegisterMCPTools(mcpTools []mcp.Tool) {
    for _, t := range mcpTools {
        // openai.Tool 형태로 변환하여 추가
    }
}
```

- [ ] **Step 5: 빌드 확인**

```bash
go build -o hanimo ./cmd/hanimo
```

- [ ] **Step 6: 커밋**

```bash
git add internal/mcp/
git commit -m "feat: add MCP client with stdio transport

JSON-RPC 2.0 protocol, tool discovery and execution.
SSE transport stubbed for Phase 2."
```

---

### Task 9: git 도구 + 슬래시 커맨드 확장

**Files:**
- Create: `internal/tools/git.go`
- Modify: `internal/tools/registry.go`, `internal/app/app.go`

- [ ] **Step 1: git.go 구현**

Create `internal/tools/git.go`:
```go
package tools

import (
    "context"
    "fmt"
    "os/exec"
    "strings"
    "time"
)

func GitStatus(path string) (string, error) {
    return runGit(path, "status", "--short")
}

func GitDiff(path string, staged bool) (string, error) {
    if staged {
        return runGit(path, "diff", "--staged")
    }
    return runGit(path, "diff")
}

func GitLog(path string, n int) (string, error) {
    return runGit(path, "log", fmt.Sprintf("-n%d", n), "--oneline")
}

func GitCommit(path, message string) (string, error) {
    return runGit(path, "commit", "-m", message)
}

func GitBranch(path string) (string, error) {
    return runGit(path, "branch", "--show-current")
}

func runGit(dir string, args ...string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    cmd := exec.CommandContext(ctx, "git", args...)
    cmd.Dir = dir
    out, err := cmd.CombinedOutput()
    return strings.TrimSpace(string(out)), err
}
```

- [ ] **Step 2: registry.go에 git 도구 등록**

```go
// AllTools()에 추가:
{Name: "git_status", ...},
{Name: "git_diff", ...},
{Name: "git_log", ...},
{Name: "git_commit", ...},

// Execute()에 분기 추가
```

- [ ] **Step 3: app.go 슬래시 커맨드 확장**

`handleSlashCommand()`에 추가:
```go
case "/save":    // session.SaveCurrentSession()
case "/load":    // session.ListSessions() → 선택 UI
case "/search":  // session.SearchSessions(keyword)
case "/model":   // 모델 전환
case "/provider": // 프로바이더 전환
case "/usage":   // session.GetSessionUsage() 표시
case "/remember": // session.SaveMemory()
case "/memories": // session.LoadMemories() 표시
case "/config":  // 현재 설정 표시
```

- [ ] **Step 4: 빌드 + 테스트**

```bash
go build -o hanimo ./cmd/hanimo
./hanimo
# /help 입력하여 새 커맨드 표시 확인
```

- [ ] **Step 5: 커밋**

```bash
git add internal/tools/git.go internal/tools/registry.go internal/app/app.go
git commit -m "feat: add git tools and extended slash commands

Git: status, diff, log, commit.
Commands: /save, /load, /search, /model, /provider, /usage,
/remember, /memories, /config."
```

---

### Task 10: CLI 인터페이스 확장

**Files:**
- Modify: `cmd/hanimo/main.go`

- [ ] **Step 1: CLI 플래그 추가**

```go
func main() {
    provider := flag.String("provider", "", "LLM provider (ollama, openai, anthropic, etc.)")
    model    := flag.String("model", "", "Model name")
    apiKey   := flag.String("api-key", "", "API key override")
    baseURL  := flag.String("base-url", "", "Custom endpoint URL")
    mode     := flag.String("mode", "super", "Start mode (super|dev|plan)")
    resume   := flag.String("resume", "", "Resume session by ID or name")
    setup    := flag.Bool("setup", false, "Re-run setup wizard")
    reset    := flag.Bool("reset", false, "Reset config")
    ver      := flag.Bool("version", false, "Print version")
    debug    := flag.Bool("debug", false, "Enable debug mode")

    // Short aliases
    flag.StringVar(provider, "p", "", "")
    flag.StringVar(model, "m", "", "")
    flag.StringVar(apiKey, "k", "", "")
    flag.StringVar(baseURL, "u", "", "")

    flag.Parse()

    // CLI에서 지정한 값은 config보다 우선
    if *provider != "" { cfg.Default.Provider = *provider }
    if *model != "" { cfg.Default.Model = *model }
    // ...
}
```

- [ ] **Step 2: --resume 구현**

```go
if *resume != "" {
    sess, msgs, err := session.LoadSession(*resume)
    // 세션 복원 → Model의 history에 주입
}
```

- [ ] **Step 3: 빌드 + 테스트**

```bash
go build -o hanimo ./cmd/hanimo
./hanimo --version
./hanimo -p ollama -m qwen3:8b
./hanimo --help
```

- [ ] **Step 4: 커밋**

```bash
git add cmd/hanimo/main.go
git commit -m "feat: extend CLI with provider, model, resume flags"
```

---

### Task 11: 지식베이스 정리 + README

**Files:**
- Modify: `knowledge/` (BXM 제거, 범용 문서 유지)
- Rewrite: `README.md`

- [ ] **Step 1: BXM 지식 제거 (신한은행 전용)**

```bash
rm -rf knowledge/docs/bxm/
```

`internal/knowledge/extractor.go`에서 BXM 관련 키워드 제거.

- [ ] **Step 2: 지식 인덱스 업데이트**

`knowledge/index.json` 재생성 (BXM 항목 제거).

- [ ] **Step 3: README.md 작성**

hanimo 오픈소스 프로젝트용 README:
- 프로젝트 소개 + 슬로건
- 주요 기능
- 설치 방법 (바이너리 다운로드 / 소스 빌드)
- Quick Start (Ollama, OpenAI, Anthropic 예시)
- 설정 (config.yaml)
- 슬래시 커맨드 목록
- 키바인딩
- 지원 프로바이더 목록
- 라이선스

- [ ] **Step 4: 커밋**

```bash
git add knowledge/ README.md internal/knowledge/
git commit -m "feat: clean knowledge base and write README

Remove BXM (Shinhan-specific) docs.
Add comprehensive README for open-source release."
```

---

## Phase 2: 고도화

### Task 12: 자율 모드 (/auto)

**Files:**
- Create: `internal/agents/auto.go`
- Modify: `internal/app/app.go`

- [ ] **Step 1: auto.go 구현**

Create `internal/agents/auto.go`:
```go
package agents

import (
    "context"
    "fmt"
    "strings"

    "github.com/flykimjiwon/hanimo/internal/llm"
    "github.com/flykimjiwon/hanimo/internal/tools"
)

const (
    MaxAutoIterations = 20
    AutoCompleteMarker = "[AUTO_COMPLETE]"
    AutoPauseMarker    = "[AUTO_PAUSE]"
)

type AutoLoop struct {
    client     *llm.Client
    model      string
    iteration  int
    maxIter    int
}

func NewAutoLoop(client *llm.Client, model string) *AutoLoop {
    return &AutoLoop{client: client, model: model, maxIter: MaxAutoIterations}
}

// RunStep executes one iteration of the auto loop
// Returns: (response string, toolCalls, shouldContinue bool, error)
func (a *AutoLoop) RunStep(ctx context.Context, history []openai.ChatCompletionMessage) (string, []llm.ToolCallInfo, bool, error) {
    a.iteration++
    if a.iteration > a.maxIter {
        return "", nil, false, fmt.Errorf("auto loop reached max iterations (%d)", a.maxIter)
    }

    // Inject auto-mode system prompt addition
    // "Complete the task autonomously. When done, output [AUTO_COMPLETE]. If blocked, output [AUTO_PAUSE]."

    // Stream response
    // Check for markers
    // Return shouldContinue based on markers and tool calls
}

func (a *AutoLoop) Iteration() int { return a.iteration }
func (a *AutoLoop) MaxIterations() int { return a.maxIter }
```

- [ ] **Step 2: app.go에 /auto 커맨드 통합**

```go
case "/auto":
    // 1. Parse remaining text as task description
    // 2. Create AutoLoop
    // 3. Inject task + "자율 실행 모드" system message
    // 4. Start loop: stream → tool calls → execute → repeat
    // 5. Display iteration count in status bar
    // 6. Stop on AUTO_COMPLETE, AUTO_PAUSE, or max iterations
```

- [ ] **Step 3: 빌드 + 테스트**

```bash
go build -o hanimo ./cmd/hanimo
./hanimo
# /auto List all Go files in current directory
```

- [ ] **Step 4: 커밋**

```bash
git add internal/agents/ internal/app/app.go
git commit -m "feat: add autonomous mode (/auto)

Agent loops up to 20 iterations until [AUTO_COMPLETE].
Executes tools, reads diagnostics, and self-corrects."
```

---

### Task 13: LSP Diagnostics

**Files:**
- Create: `internal/tools/diagnostics.go`
- Modify: `internal/tools/registry.go`

- [ ] **Step 1: diagnostics.go 구현**

Create `internal/tools/diagnostics.go`:
```go
package tools

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
)

type Diagnostic struct {
    File     string `json:"file"`
    Line     int    `json:"line"`
    Column   int    `json:"column"`
    Severity string `json:"severity"` // error/warning
    Message  string `json:"message"`
    Source   string `json:"source"`   // tsc/eslint/go
}

func RunDiagnostics(dir string, targetFile string) ([]Diagnostic, error) {
    projectType := detectProjectType(dir)

    var diagnostics []Diagnostic

    switch projectType {
    case "go":
        d, _ := runGoVet(dir, targetFile)
        diagnostics = append(diagnostics, d...)
    case "typescript", "javascript":
        d, _ := runTsc(dir)
        diagnostics = append(diagnostics, d...)
        d2, _ := runEslint(dir, targetFile)
        diagnostics = append(diagnostics, d2...)
    case "python":
        d, _ := runRuff(dir, targetFile)
        diagnostics = append(diagnostics, d...)
    }

    return diagnostics, nil
}

func detectProjectType(dir string) string {
    if fileExists(filepath.Join(dir, "go.mod")) { return "go" }
    if fileExists(filepath.Join(dir, "tsconfig.json")) { return "typescript" }
    if fileExists(filepath.Join(dir, "package.json")) { return "javascript" }
    if fileExists(filepath.Join(dir, "pyproject.toml")) || fileExists(filepath.Join(dir, "requirements.txt")) { return "python" }
    return "unknown"
}

func runGoVet(dir, file string) ([]Diagnostic, error) {
    args := []string{"vet"}
    if file != "" { args = append(args, file) } else { args = append(args, "./...") }
    // Execute and parse output
}

func runTsc(dir string) ([]Diagnostic, error) {
    // npx tsc --noEmit --pretty false
    // Parse "file(line,col): error TS1234: message" format
}

func runEslint(dir, file string) ([]Diagnostic, error) {
    // npx eslint --format json file
    // Parse JSON output
}

func runRuff(dir, file string) ([]Diagnostic, error) {
    // ruff check --output-format json file
}

func FormatDiagnostics(diags []Diagnostic) string {
    if len(diags) == 0 { return "No issues found." }
    var sb strings.Builder
    for _, d := range diags {
        sb.WriteString(fmt.Sprintf("%s:%d:%d [%s] %s (%s)\n", d.File, d.Line, d.Column, d.Severity, d.Message, d.Source))
    }
    sb.WriteString(fmt.Sprintf("\nTotal: %d issues", len(diags)))
    return sb.String()
}
```

- [ ] **Step 2: registry.go에 등록**

```go
{Name: "diagnostics", Description: "Run code diagnostics (go vet, tsc, eslint, ruff)", Parameters: ...},
```

- [ ] **Step 3: /diagnostics 슬래시 커맨드 연결**

app.go에서 `/diagnostics [file]` → `tools.RunDiagnostics()` → 결과 표시

- [ ] **Step 4: 빌드 + 테스트**

```bash
go build -o hanimo ./cmd/hanimo
./hanimo
# /diagnostics
```

- [ ] **Step 5: 커밋**

```bash
git add internal/tools/diagnostics.go internal/tools/registry.go internal/app/app.go
git commit -m "feat: add LSP diagnostics (go vet, tsc, eslint, ruff)

Auto-detects project type, runs appropriate linters,
returns structured error/warning list."
```

---

### Task 14: Role 시스템 + 모델 능력 매트릭스

**Files:**
- Create: `internal/llm/capabilities.go`
- Modify: `internal/llm/prompt.go`

- [ ] **Step 1: capabilities.go 구현**

Create `internal/llm/capabilities.go`:
```go
package llm

type CodingTier int
const (
    CodingTierStrong   CodingTier = iota // GPT-4o, Claude Sonnet, Qwen3-Coder
    CodingTierModerate                    // GPT-4o-mini, Gemini Flash
    CodingTierWeak                        // Small Ollama models
    CodingTierNone                        // Chat-only models
)

type RoleAssignment int
const (
    RoleAgent     RoleAssignment = iota // Full tools, file writes
    RoleAssistant                        // Read-only tools
    RoleChat                             // No tools
)

type ModelCapability struct {
    ContextWindow  int
    CodingTier     CodingTier
    DefaultRole    RoleAssignment
    SupportsTools  bool
    SupportsVision bool
}

var knownModels = map[string]ModelCapability{
    // Cloud models
    "gpt-4o":           {128000, CodingTierStrong, RoleAgent, true, true},
    "gpt-4o-mini":      {128000, CodingTierModerate, RoleAgent, true, true},
    "claude-sonnet-4":  {200000, CodingTierStrong, RoleAgent, true, true},
    "claude-haiku-4":   {200000, CodingTierModerate, RoleAgent, true, false},
    "gemini-2.5-flash": {1000000, CodingTierModerate, RoleAgent, true, true},
    "deepseek-chat":    {128000, CodingTierStrong, RoleAgent, true, false},
    // Ollama models
    "qwen3:8b":         {32768, CodingTierModerate, RoleAssistant, true, false},
    "qwen3:32b":        {32768, CodingTierStrong, RoleAgent, true, false},
    "llama3.1:8b":      {128000, CodingTierWeak, RoleChat, false, false},
    "codellama:13b":    {16384, CodingTierModerate, RoleAssistant, true, false},
}

func GetCapability(model string) ModelCapability {
    if cap, ok := knownModels[model]; ok { return cap }
    // Unknown model: assume moderate capability
    return ModelCapability{32768, CodingTierModerate, RoleAssistant, true, false}
}

func AutoAssignRole(model string) RoleAssignment {
    return GetCapability(model).DefaultRole
}
```

- [ ] **Step 2: prompt.go에 Role 기반 시스템 프롬프트 적용**

```go
func SystemPromptForRole(role RoleAssignment, mode Mode) string {
    // RoleAgent: 기존 Super/Dev 프롬프트 (전체 도구)
    // RoleAssistant: Plan 프롬프트 (읽기 전용)
    // RoleChat: 도구 없음, 대화만
}
```

- [ ] **Step 3: 빌드 + 테스트**

```bash
go build -o hanimo ./cmd/hanimo
./hanimo -p ollama -m llama3.1:8b
# 자동으로 Chat role 할당, 도구 비활성화 확인
```

- [ ] **Step 4: 커밋**

```bash
git add internal/llm/capabilities.go internal/llm/prompt.go
git commit -m "feat: add role system with model capability matrix

Auto-assigns Agent/Assistant/Chat role based on model.
160+ models with context window, coding tier metadata."
```

---

### Task 15: 커맨드 팔레트 + 테마 + 최종 정리

**Files:**
- Create: `internal/ui/palette.go`
- Modify: `internal/app/app.go`, `internal/ui/styles.go`

- [ ] **Step 1: palette.go 구현**

Create `internal/ui/palette.go`:
```go
package ui

import (
    "strings"
    "sort"
)

type PaletteItem struct {
    Label       string
    Description string
    Action      string // slash command to execute
}

var paletteItems = []PaletteItem{
    {Label: "Save Session", Description: "Save current session", Action: "/save"},
    {Label: "Load Session", Description: "Load a previous session", Action: "/load"},
    {Label: "Search Sessions", Description: "Search past sessions", Action: "/search"},
    {Label: "Switch Model", Description: "Change LLM model", Action: "/model"},
    {Label: "Switch Provider", Description: "Change LLM provider", Action: "/provider"},
    {Label: "Usage Stats", Description: "Show token usage and cost", Action: "/usage"},
    {Label: "Auto Mode", Description: "Start autonomous execution", Action: "/auto"},
    {Label: "Diagnostics", Description: "Run code diagnostics", Action: "/diagnostics"},
    {Label: "Remember", Description: "Save a memory", Action: "/remember"},
    {Label: "Memories", Description: "Show saved memories", Action: "/memories"},
    {Label: "Config", Description: "Show current configuration", Action: "/config"},
    {Label: "Clear", Description: "Clear conversation", Action: "/clear"},
    {Label: "Help", Description: "Show help", Action: "/help"},
}

func FuzzySearch(query string, items []PaletteItem) []PaletteItem {
    if query == "" { return items }
    query = strings.ToLower(query)
    var matched []PaletteItem
    for _, item := range items {
        label := strings.ToLower(item.Label)
        desc := strings.ToLower(item.Description)
        if strings.Contains(label, query) || strings.Contains(desc, query) {
            matched = append(matched, item)
        }
    }
    return matched
}

func RenderPalette(items []PaletteItem, selected int, query string, width int) string {
    // Render floating overlay with search input + item list
    // Highlight selected item
}
```

- [ ] **Step 2: app.go에 Ctrl+K 팔레트 연결**

```go
// KeyMsg 핸들링에서:
case "ctrl+k":
    m.showPalette = !m.showPalette
    m.paletteQuery = ""
    m.paletteSelected = 0
```

- [ ] **Step 3: 테마 시스템 기본 구현**

`internal/ui/styles.go`에 테마 프리셋 추가:
```go
type Theme struct {
    Name    string
    Primary lipgloss.Color
    // ...
}

var themes = map[string]Theme{
    "default": {...},
    "dark":    {...},
    "ocean":   {...},
}
```

- [ ] **Step 4: 전체 빌드 + 통합 테스트**

```bash
go build -o hanimo ./cmd/hanimo
./hanimo --version
./hanimo -p ollama -m qwen3:8b
# Ctrl+K 팔레트, /auto, /diagnostics, /save, /load 전부 확인
```

- [ ] **Step 5: 커밋**

```bash
git add internal/ui/palette.go internal/app/app.go internal/ui/styles.go
git commit -m "feat: add command palette (Ctrl+K) and theme system

Fuzzy search across all slash commands.
Theme presets: default, dark, ocean."
```

---

### Task 16: 최종 빌드 + 크로스 컴파일 + 태그

- [ ] **Step 1: 전체 린트**

```bash
go vet ./...
```

- [ ] **Step 2: 크로스 컴파일**

```bash
make build-all
ls -la dist/
```

Expected:
```
hanimo-darwin-arm64
hanimo-darwin-amd64
hanimo-linux-amd64
hanimo-linux-arm64
hanimo-windows-amd64.exe
```

- [ ] **Step 3: 로컬 실행 테스트**

```bash
./dist/hanimo-darwin-arm64 --version
./dist/hanimo-darwin-arm64 -p ollama -m qwen3:8b
```

- [ ] **Step 4: .gitignore 업데이트**

```
dist/
*.exe
hanimo
_legacy_ts/node_modules/
.DS_Store
```

- [ ] **Step 5: 최종 커밋 + 태그**

```bash
git add -A
git commit -m "feat: hanimo v0.1.0 — Go rebuild complete

TECHAI_CODE base + 14 providers + SQLite sessions/memory +
Hash-Anchored editing + Smart Compaction + MCP client +
autonomous mode + LSP diagnostics + command palette."

git tag v0.1.0
```

---

## Phase 3 (문서만 — 미구현)

### Future: 멀티 에이전트 (SubAgent)
- Orchestrator가 태스크를 독립 서브태스크로 분해
- 각 SubAgent를 goroutine으로 병렬 실행
- 결과 합성하여 사용자에게 반환
- `internal/agents/orchestrator.go`, `internal/agents/subagent.go`

### Future: 데스크톱 앱 (Wails)
- Go 백엔드 (hanimo-core 직접 임포트)
- React/Svelte 프론트엔드 (웹뷰)
- Monaco 에디터 통합
- 단일 `.app` / `.exe` 빌드

### Future: VS Code Extension
- Go 바이너리를 `--headless` 모드로 사이드카 실행
- JSON-RPC stdin/stdout 통신
- TS Extension: 웹뷰 채팅 UI + 코드렌즈 + 인라인 제안
- `.vsix` + 바이너리 번들 배포

### Future: headless 모드
- `hanimo --headless`: JSON-RPC 서버로 동작
- stdin으로 요청 수신, stdout으로 응답 전송
- 데스크톱/Extension/외부 도구 연동 기반 인프라

### Future: 플러그인 시스템
- Go 플러그인 또는 WASM 기반 확장
- 사용자 정의 도구/프로바이더 동적 로딩
