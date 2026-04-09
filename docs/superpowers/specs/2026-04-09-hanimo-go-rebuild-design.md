# hanimo Go Rebuild — Design Spec

> Date: 2026-04-09
> Author: 김지원
> Status: Draft → Approved

## 1. Overview

TECHAI_CODE(Go, 폐쇄망 전용 AI 코딩 에이전트)를 포크하여 **hanimo**(오픈소스, 멀티프로바이더 AI 코딩 에이전트)를 Go로 리빌드한다.

기존 hanimo(TypeScript)의 모든 고도화 기능을 Go로 재구현하고, TECHAI_CODE의 배포/성능 장점을 그대로 가져간다.

### 슬로건

> "오픈소스를 위한 AI 코딩 에이전트"

### 핵심 결정

- **언어**: Go (TECHAI_CODE 기반)
- **TUI**: Bubble Tea v2
- **DB**: SQLite (`modernc.org/sqlite`, CGO-free)
- **바이너리**: 단일 파일, 크로스 컴파일 5개 플랫폼
- **브랜딩**: hanimo 로고/색상, `.hanimo.md`, `~/.hanimo/`

---

## 2. Architecture

```
cmd/hanimo/main.go                  ← 엔트리포인트
internal/
├── app/                            ← Bubble Tea TUI (TECHAI_CODE 기반, 브랜딩 교체)
│   ├── model.go                    # Bubble Tea Model/Update/View
│   ├── mode.go                     # 3모드: Super/Dev/Plan (Tab 전환)
│   ├── commands.go                 # 슬래시 커맨드 처리
│   └── queue.go                    # 메시지 큐 (타이핑 중 입력)
├── config/
│   ├── config.go                   # YAML 로드 + 환경변수 오버라이드
│   ├── setup.go                    # 초기 설정 위저드
│   └── defaults.go                 # 빌드타임 기본값 (ldflags 주입)
├── llm/
│   ├── client.go                   # 공통 스트리밍 인터페이스
│   ├── providers/
│   │   ├── registry.go             # Provider 인터페이스 + 팩토리
│   │   ├── openai_compat.go        # OpenAI-compatible (8개 클라우드 + custom)
│   │   ├── anthropic.go            # Anthropic SDK
│   │   ├── google.go               # Google Gemini SDK
│   │   └── ollama.go               # Ollama 네이티브
│   ├── context.go                  # 시스템 컨텍스트 수집
│   ├── compaction.go               # Smart Compaction (snip → micro → LLM 요약)
│   └── capabilities.go            # 모델 능력 매트릭스 (role/coding tier/context window)
├── tools/
│   ├── registry.go                 # 도구 JSON 스키마 정의
│   ├── file.go                     # file_read, file_write, file_edit, file_delete
│   ├── hashline.go                 # Hash-Anchored 편집 (MD5 해시 라인 태깅)
│   ├── shell.go                    # shell_exec
│   ├── search.go                   # grep_search, glob_search, list_files
│   ├── git.go                      # git status/diff/log/commit
│   └── diagnostics.go             # tsc/eslint/go vet 실행 + 결과 파싱
├── session/
│   ├── store.go                    # SQLite CRUD (sessions, messages)
│   ├── memory.go                   # 장기 메모리 (memories 테이블)
│   ├── usage.go                    # 토큰 사용량/비용 추적
│   └── migrate.go                  # DB 스키마 마이그레이션
├── agents/
│   ├── auto.go                     # 자율 모드 (/auto, 최대 20회 루프)
│   └── role.go                     # Role 시스템 (dev/plan/chat, 모델별 자동 할당)
├── mcp/
│   ├── client.go                   # MCP 클라이언트 (JSON-RPC 2.0)
│   ├── transport_stdio.go          # stdio 트랜스포트
│   ├── transport_sse.go            # SSE 트랜스포트
│   └── registry.go                 # MCP 서버 관리 (SQLite)
├── ui/
│   ├── theme.go                    # hanimo 브랜딩 (색상 팔레트, 로고 ASCII)
│   ├── components.go              # 공통 TUI 컴포넌트
│   └── palette.go                  # 커맨드 팔레트 (Ctrl+K 퍼지 검색)
└── knowledge/
    └── store.go                    # embed.FS 기반, Tiered 우선순위, 키워드 매칭

knowledge/                          ← embed.FS로 바이너리에 내장
├── docs/                           # 언어/프레임워크 문서
├── skills/                         # 스킬 문서
└── index.json                      # 메타데이터 인덱스

db/
└── schema.sql                      # SQLite 스키마 정의
```

---

## 3. SQLite Schema

```sql
-- ~/.hanimo/hanimo.db

CREATE TABLE sessions (
    id          TEXT PRIMARY KEY,
    name        TEXT,
    project_dir TEXT NOT NULL,
    provider    TEXT NOT NULL,
    model       TEXT NOT NULL,
    mode        TEXT DEFAULT 'super',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE messages (
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

CREATE TABLE memories (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    project_dir TEXT NOT NULL,
    key         TEXT NOT NULL,
    value       TEXT NOT NULL,
    source      TEXT DEFAULT 'auto',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_dir, key)
);

CREATE TABLE usage_log (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id  TEXT NOT NULL REFERENCES sessions(id),
    provider    TEXT NOT NULL,
    model       TEXT NOT NULL,
    tokens_in   INTEGER NOT NULL,
    tokens_out  INTEGER NOT NULL,
    cost_usd    REAL DEFAULT 0,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE mcp_servers (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT UNIQUE NOT NULL,
    transport   TEXT NOT NULL,
    command     TEXT,
    args        TEXT,
    url         TEXT,
    enabled     BOOLEAN DEFAULT 1,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_messages_session ON messages(session_id);
CREATE INDEX idx_memories_project ON memories(project_dir);
CREATE INDEX idx_usage_session ON usage_log(session_id);
```

---

## 4. Provider Architecture

### Interface

```go
type Provider interface {
    Name() string
    Chat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
    Models() ([]ModelInfo, error)
    SupportsTools() bool
}
```

### Provider Map

| 파일 | 커버하는 프로바이더 | 방식 |
|------|-------------------|------|
| `openai_compat.go` | OpenAI, Novita, OpenRouter, DeepSeek, Groq, Together, Fireworks, Mistral, vLLM, LM Studio, Custom | base URL 교체 |
| `anthropic.go` | Anthropic (Claude) | `anthropic-sdk-go` |
| `google.go` | Google (Gemini) | `generative-ai-go` |
| `ollama.go` | Ollama | REST API 네이티브 |

### Config (`~/.hanimo/config.yaml`)

```yaml
default:
  provider: ollama
  model: qwen3:8b

providers:
  ollama: {}
  novita:
    api_key: "nvt-..."
  openrouter:
    api_key: "sk-or-..."
  anthropic:
    api_key: "sk-ant-..."
  google:
    api_key: "AIza..."
  custom:
    - name: "my-server"
      base_url: "http://192.168.1.100:8000/v1"
      models: ["llama-3.1-70b"]
```

---

## 5. Feature Detail

### 5.1 Hash-Anchored Editing

파일 읽기 시 각 라인에 MD5 해시 태그 부여:
```
1#a3f1| function hello() {
2#7bc2|   return "world"
3#e4d9| }
```

편집 시 해시 앵커로 검증:
```json
{
  "file": "main.go",
  "startAnchor": "1#a3f1",
  "endAnchor": "2#7bc2",
  "newContent": "func hello() string {\n  return \"world\"\n"
}
```

해시 불일치 시 에러 + "re-read" 힌트 반환. 파일 깨짐 방지.

### 5.2 Smart Compaction

3단계 컨텍스트 압축:

1. **Snip** (40+ 메시지): 오래된 도구 결과를 `[snipped: 150 lines]`로 교체
2. **Micro** (여전히 초과): 긴 메시지 앞절반 + 뒷절반만 유지
3. **LLM Summary** (여전히 초과): LLM에게 이전 대화 요약 요청, 요약본으로 교체

보존 대상: 태스크 목표, 완료된 작업, 현재 상태.

### 5.3 Session Management

| 커맨드 | 동작 |
|--------|------|
| `/save [name]` | 현재 세션 SQLite에 저장 |
| `/load` | 세션 목록 표시 + 선택 |
| `/search <keyword>` | 과거 세션 검색 |
| `--resume [id\|name]` | CLI에서 세션 복원 |
| `--fork <id>` | 세션 분기 (새 ID, 히스토리 복사) |

### 5.4 Memory System

```
/remember "이 프로젝트는 PostgreSQL 16 사용"
  → memories 테이블에 저장 (project_dir + key 기준 UPSERT)

다음 세션 시작 시:
  → 해당 project_dir의 memories 자동 로드 → 시스템 프롬프트에 주입
```

### 5.5 MCP Client

- JSON-RPC 2.0 프로토콜 구현
- stdio 트랜스포트: 외부 프로세스 spawn + stdin/stdout 통신
- SSE 트랜스포트: HTTP 연결
- MCP 서버 설정은 SQLite `mcp_servers` 테이블 또는 config.yaml
- 도구 자동 등록: MCP 서버의 tools를 hanimo 도구로 노출

### 5.6 Autonomous Mode (/auto)

```
/auto Fix all TypeScript errors in src/

루프:
  1. LLM에게 현재 상태 전달
  2. LLM이 도구 호출 (파일 수정, 진단 실행 등)
  3. 결과 확인
  4. [AUTO_COMPLETE] 마커 → 종료
  5. 아니면 → 1로 돌아감 (최대 20회)
```

### 5.7 LSP Diagnostics

```
/diagnostics              → 프로젝트 전체 진단
/diagnostics src/main.go  → 특정 파일 진단

실행: go vet, tsc --noEmit, eslint (프로젝트 타입 자동 감지)
결과를 구조화해서 에이전트에게 전달 → 자동 수정 가능
```

---

## 6. Branding

### 6.1 이름/경로

| 항목 | TECHAI_CODE | hanimo |
|------|------------|--------|
| 바이너리 | `techai` | `hanimo` |
| 설정 디렉토리 | `~/.tgc/` | `~/.hanimo/` |
| 프로젝트 설정 | `.techai.md` | `.hanimo.md` |
| DB 파일 | — | `~/.hanimo/hanimo.db` |
| 디버그 로그 | `~/.tgc/debug.log` | `~/.hanimo/debug.log` |

### 6.2 TUI 브랜딩

- 기존 hanimo의 로고 ASCII art 그대로 사용
- 색상 팔레트: 기존 hanimo 테마 유지
- 모드 표시, 프롬프트 스타일: TECHAI_CODE 레이아웃 기반 + hanimo 색상

---

## 7. CLI Interface

```
hanimo [prompt...]
  -p, --provider <name>      프로바이더 선택
  -m, --model <name>         모델 선택
  -k, --api-key <key>        API 키 오버라이드
  -u, --base-url <url>       커스텀 엔드포인트
  --mode <super|dev|plan>    시작 모드
  --resume [id|name]         세션 복원
  --fork <id>                세션 분기
  --setup                    설정 위저드 재실행
  --reset                    설정 초기화
  --version                  버전 출력
  --debug                    디버그 모드
```

### 키바인딩

| 키 | 동작 |
|----|------|
| Enter | 메시지 전송 |
| Shift+Enter | 줄바꿈 |
| Tab | 모드 전환 (Super → Dev → Plan) |
| Esc | 스트리밍 취소 |
| Ctrl+K | 커맨드 팔레트 |
| Ctrl+L | 화면 클리어 |
| Ctrl+C | 종료 |
| Alt+↑/↓ | 스크롤 |

### 슬래시 커맨드

| 커맨드 | 동작 |
|--------|------|
| `/help` | 도움말 |
| `/model [name]` | 모델 전환 |
| `/provider [name]` | 프로바이더 전환 |
| `/mode [super\|dev\|plan]` | 모드 전환 |
| `/tools [on\|off]` | 도구 토글 |
| `/config` | 현재 설정 표시 |
| `/usage` | 토큰 사용량 + 비용 |
| `/theme [id]` | 테마 변경 |
| `/auto [msg]` | 자율 모드 |
| `/save [name]` | 세션 저장 |
| `/load` | 세션 로드 |
| `/search [keyword]` | 세션 검색 |
| `/remember [text]` | 메모리 저장 |
| `/memories` | 저장된 메모리 조회 |
| `/diagnostics [file]` | 코드 진단 |
| `/clear` | 대화 초기화 |
| `/exit` | 종료 |

---

## 8. Build & Distribution

### Makefile

```makefile
build:          # 현재 플랫폼 → ./hanimo
build-all:      # 5개 플랫폼 크로스 컴파일
                # hanimo-darwin-arm64, hanimo-darwin-amd64
                # hanimo-linux-amd64, hanimo-linux-arm64
                # hanimo-windows-amd64.exe
install:        # go install → $GOPATH/bin/hanimo
test:           # go test ./...
lint:           # go vet + staticcheck
clean:          # rm -rf dist/
```

### ldflags

```makefile
LDFLAGS = -ldflags "-s -w \
  -X main.version=$(VERSION) \
  -X main.buildDate=$(DATE)"
```

---

## 9. Dependencies (go.mod)

```
github.com/charmbracelet/bubbletea   v2  # TUI 프레임워크
github.com/charmbracelet/lipgloss    v2  # 스타일링
github.com/charmbracelet/bubbles     v2  # TUI 컴포넌트
github.com/charmbracelet/glamour     v2  # 마크다운 렌더링
github.com/sashabaranov/go-openai        # OpenAI-compatible 클라이언트
github.com/anthropics/anthropic-sdk-go   # Anthropic 클라이언트
github.com/google/generative-ai-go      # Google Gemini 클라이언트
modernc.org/sqlite                       # SQLite (CGO-free)
gopkg.in/yaml.v3                         # 설정 파일 파싱
github.com/google/uuid                   # 세션 ID 생성
```

**총 10개 의존성.** (node_modules 329개 → 10개)

---

## 10. Phases

### Phase 1 — MVP

1. TECHAI_CODE 포크 + 리네임 (cmd/hanimo, 설정 경로, 브랜딩)
2. 14개 프로바이더 레지스트리 구현
3. SQLite 통합 (세션, 메모리, 사용량)
4. Hash-Anchored 편집
5. Smart Compaction
6. MCP 클라이언트
7. git 도구 추가
8. 슬래시 커맨드 확장

### Phase 2 — 고도화

9. 자율 모드 (/auto)
10. LSP Diagnostics
11. Role 시스템 (모델 능력 기반 자동 할당)
12. 멀티 엔드포인트 + 로드밸런싱
13. 토큰 비용 추적 (/usage)
14. 테마 시스템 (/theme)
15. 커맨드 팔레트 (Ctrl+K)

### Phase 3 — 확장 (문서만, 이번 범위 아님)

16. 멀티 에이전트 (SubAgent 병렬 실행)
17. 데스크톱 앱 (Wails: Go 백엔드 + 웹 프론트)
18. VS Code Extension (Go 사이드카 + TS Extension)
19. headless 모드 (JSON-RPC stdin/stdout)
20. 플러그인 시스템

---

## 11. TECHAI_CODE 역방향 포팅 (다음 작업)

이번 작업 완료 후, hanimo에서 검증된 기능을 TECHAI_CODE로 포팅:

| 기능 | 우선순위 |
|------|---------|
| SQLite 세션/메모리 | 높음 (TECHAI_CODE에도 예정) |
| Hash-Anchored 편집 | 높음 |
| Smart Compaction | 중간 |
| 자율 모드 | 중간 |
| 추가 도구 (git, diagnostics) | 낮음 |

프로바이더/MCP는 폐쇄망이라 불필요.

---

## 12. Risks & Mitigations

| 리스크 | 대응 |
|--------|------|
| MCP Go 구현 복잡도 | JSON-RPC 2.0은 단순한 프로토콜. stdio 먼저, SSE는 후순위 |
| anthropic-sdk-go 안정성 | 공식 SDK, 2025 출시. 불안정하면 OpenAI-compatible 폴백 |
| modernc.org/sqlite 바이너리 크기 증가 | ~5MB 추가 예상 (17→22MB). 허용 범위 |
| 기존 hanimo TS 사용자 마이그레이션 | config 포맷 유지, 세션 JSON→SQLite 변환 스크립트 제공 |
