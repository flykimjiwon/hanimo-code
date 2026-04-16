<p align="center">
  <img src="docs/images/hanimo-screenshot.png" alt="hanimo screenshot" width="600">
</p>

<h1 align="center">hanimo</h1>

<p align="center">
  <strong>Open-Source AI Coding Agent for the Terminal</strong><br>
  <em>터미널 기반 오픈소스 AI 코딩 에이전트</em>
</p>

<p align="center">
  <a href="https://github.com/flykimjiwon/hanimo/releases"><img src="https://img.shields.io/github/v/release/flykimjiwon/hanimo?style=flat-square" alt="Release"></a>
  <a href="https://golang.org"><img src="https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat-square&logo=go" alt="Go"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow?style=flat-square" alt="License"></a>
  <a href="https://github.com/flykimjiwon/hanimo/stargazers"><img src="https://img.shields.io/github/stars/flykimjiwon/hanimo?style=flat-square" alt="Stars"></a>
</p>

<p align="center">
  <a href="#features">Features</a> |
  <a href="#quick-start">Quick Start</a> |
  <a href="#commands">Commands</a> |
  <a href="#한국어-가이드">한국어</a>
</p>

---

## Why hanimo?

hanimo ships as a **single ~20MB binary** with zero runtime dependencies. It supports **14+ LLM providers** out of the box, runs entirely in your terminal, and comes with **23+ built-in tools**, **multi-agent orchestration**, a **browser companion dashboard**, and **62 embedded knowledge docs**.

Think of it as an open-source alternative to Claude Code / Cursor Agent that you can run anywhere — from your laptop to an air-gapped server.

---

## Features

### Core
- **14+ LLM Providers** — Ollama, OpenAI, Anthropic, Google, Novita, OpenRouter, DeepSeek, Groq, Together, Fireworks, Mistral, vLLM, LM Studio, Custom
- **3 Modes** — Super (all-purpose), Dev (autonomous coding), Plan (read-only analysis)
- **30+ Slash Commands** — Full CLI experience with `/auto`, `/plan`, `/multi`, `/companion`, and more
- **SQLite Sessions** — Save, load, search, fork conversation sessions
- **5 Themes** — honey (default), ocean, dracula, nord, forest
- **Bilingual UI** — English + Korean (`/lang` to toggle)

### AI Agent
- **Autonomous Mode** (`/auto`) — Self-driving task completion (up to 200 iterations)
- **Planning Mode** (`/plan`) — Step-by-step execution plans with `/approve` / `/reject`
- **Multi-Agent** (`/multi`) — Dual-model orchestration with Review, Consensus, Scan, Auto strategies
- **Smart Compaction** — 3-stage context compression (snip → micro → LLM summary)
- **Context Tracking** — Real-time ctx:XX% in status bar with auto-compact at 90%

### Tools (23+)
| Category | Tools |
|----------|-------|
| **File Operations** | `file_read`, `file_write`, `file_edit`, `hashline_read`, `hashline_edit` |
| **Directory** | `list_files`, `list_tree` |
| **Search** | `grep_search`, `glob_search`, `knowledge_search` |
| **Shell** | `shell_exec` (streaming output, dangerous command blocking) |
| **Git** | `git_status`, `git_diff`, `git_log`, `git_commit` |
| **Project** | `project_detect`, `init_project` |
| **Code Quality** | `diagnostics` (go vet, tsc, eslint, ruff auto-detect) |

### Security
- **Secret Detection** — Blocks API keys, AWS credentials, GitHub tokens, JWTs, private keys from being written
- **Sensitive File Guard** — Prevents writes to `.env`, `.pem`, `credentials.json`, etc.
- **Auto Snapshot** — Every file is backed up before modification (`/undo` to restore)
- **Read-Before-Write** — LLM must read a file before editing (prevents hallucinated overwrites)
- **Dangerous Command Blocking** — `rm -rf /`, `sudo`, credential exfiltration blocked

### Developer Experience
- **Hash-Anchored Editing** — MD5 hash per line prevents stale-edit corruption
- **Git HUD** — Branch name + dirty indicator in status bar
- **Input History** — Arrow `↑/↓` or `Ctrl+P/N` to cycle through previous inputs
- **Command Palette** — `Ctrl+K` fuzzy search across all commands
- **Custom Commands** — Drop `.md` files in `.hanimo/commands/` for reusable prompts
- **Companion Dashboard** — Browser-based real-time SSE dashboard (`/companion`)

### Knowledge & Skills
- **62 Embedded Docs** — React, Next.js, Spring, SQL, PostgreSQL, Tailwind, Vite, and more
- **User Knowledge** — Drop `.md`/`.txt` files in `.hanimo/knowledge/` for project-specific context
- **SKILL.md Loader** — agentskills.io compatible skill system
- **Project Profiles** — `/init` generates `.hanimo.md` with auto-detected structure

### Extensibility
- **MCP Support** — Extend tools via Model Context Protocol (stdio + SSE transports)
- **Multi-Provider** — Route different modes to different models/providers
- **Prompt Caching** — Infrastructure for cache-aware prompt construction
- **Baked Builds** — Compile with frozen endpoint/model/key for distribution

---

## Quick Start

### Option 1: Download Binary

```bash
# macOS (Apple Silicon)
curl -L -o hanimo https://github.com/flykimjiwon/hanimo/releases/latest/download/hanimo-darwin-arm64
chmod +x hanimo
sudo mv hanimo /usr/local/bin/

# macOS (Intel)
curl -L -o hanimo https://github.com/flykimjiwon/hanimo/releases/latest/download/hanimo-darwin-amd64

# Linux (x64)
curl -L -o hanimo https://github.com/flykimjiwon/hanimo/releases/latest/download/hanimo-linux-amd64
```

### Option 2: Build from Source

```bash
git clone https://github.com/flykimjiwon/hanimo.git
cd hanimo
go build -o hanimo ./cmd/hanimo
./hanimo
```

### Run

```bash
# Default: connects to Ollama on localhost (no API key needed)
ollama pull qwen3:8b
hanimo

# With specific provider
hanimo -p openai -m gpt-4o
hanimo -p anthropic -m claude-sonnet-4
hanimo -p deepseek -m deepseek-chat

# Custom endpoint (vLLM, LM Studio, etc.)
hanimo -u http://localhost:8000/v1 -m my-model

# Start in specific mode
hanimo --mode dev     # Autonomous coding
hanimo --mode plan    # Read-only analysis

# Resume a previous session
hanimo --resume my-project
```

---

## Commands

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Enter` | Send message |
| `Shift+Enter` | New line |
| `Tab` | Switch mode (Super → Dev → Plan) |
| `↑/↓` | Input history |
| `Ctrl+K` | Command palette |
| `Ctrl+P/N` | Input history (alt) |
| `Ctrl+L` | Clear screen |
| `Esc` | Menu / Cancel streaming |
| `Ctrl+C` | Exit |
| `Alt+↑/↓` | Scroll viewport |

### Slash Commands

| Command | Description |
|---------|-------------|
| **Session** | |
| `/new` | Start fresh session |
| `/save [name]` | Save session to SQLite |
| `/load [prefix]` | Load saved session |
| `/search [keyword]` | Search past sessions |
| `/export [filename]` | Export session to markdown |
| **AI Control** | |
| `/auto [task]` | Autonomous mode (up to 200 iterations) |
| `/plan [task]` | Create execution plan |
| `/approve` | Approve plan |
| `/reject` | Reject plan |
| `/multi [strategy]` | Multi-agent toggle (on/off/review/consensus/scan/auto) |
| `/compact` | Manual context compression |
| **Tools** | |
| `/git` | Show git branch and status |
| `/diff` | Show git diff |
| `/init` | Generate .hanimo.md project profile |
| `/diagnostics [file]` | Run code diagnostics |
| `/undo [N\|list]` | Undo file changes (snapshot restore) |
| `/copy` | Copy last AI response to clipboard |
| `/commands` | List custom commands |
| **Configuration** | |
| `/model [name]` | Switch model |
| `/provider [name]` | Switch provider |
| `/config` | Show current configuration |
| `/theme [name]` | Change theme (honey/ocean/dracula/nord/forest) |
| `/lang` | Toggle language (English/Korean) |
| `/setup` | Re-run API key configuration |
| `/companion` | Launch browser dashboard |
| `/mcp` | MCP server status |
| **Info** | |
| `/help` | Show help |
| `/version` | Show version |
| `/usage` | Token usage and cost |
| `/remember [text]` | Save project memory |
| `/memories` | Show saved memories |
| `/forget [key]` | Remove memory entry |
| `/clear` | Clear conversation |
| `/exit`, `/quit` | Exit hanimo |

---

## Multi-Agent Mode

hanimo supports dual-model orchestration where two AI agents work together on your task.

```bash
/multi review     # Agent1 generates → Agent2 reviews (read-only)
/multi consensus  # Both agents generate independently → LLM synthesizes
/multi scan       # Partition work by file → each agent handles a section
/multi auto       # Auto-detect best strategy based on query
/multi off        # Disable multi-agent
```

The orchestrator runs both agents in parallel with real-time progress reporting. Results are merged via LLM synthesis.

---

## Companion Dashboard

Launch a browser-based real-time dashboard:

```bash
/companion
```

Opens `http://localhost:8420` with:
- Live streaming of AI responses
- Tool call monitoring
- Mode/session status
- SSE-based auto-reconnect

---

## Configuration

Config file: `~/.hanimo/config.yaml`

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
```

Environment variables override config:

| Variable | Description |
|----------|-------------|
| `HANIMO_API_KEY` | API key |
| `HANIMO_API_BASE_URL` | API endpoint URL |
| `HANIMO_MODEL_SUPER` | Model for Super/Plan mode |
| `HANIMO_MODEL_DEV` | Model for Dev mode |

---

## Project Context

Create `.hanimo.md` in your project root (or run `/init` to auto-generate):

```markdown
# My Project
- Framework: Next.js 15 with App Router
- Language: TypeScript strict mode
- Testing: Vitest
- Always respond in Korean
```

hanimo automatically loads this into the system prompt every session.

### Knowledge Base

Drop `.md` or `.txt` files into `.hanimo/knowledge/` for project-specific docs:

```
.hanimo/knowledge/
├── api-conventions.md
├── database-schema.md
└── deployment-guide.txt
```

hanimo indexes and searches these docs when relevant to your query.

---

## Architecture

```
cmd/hanimo/            CLI entry point
internal/
├── app/               Bubble Tea TUI (3000+ lines)
├── agents/            Autonomous mode + planning + intent
├── companion/         Browser SSE dashboard (server, hub, events)
├── config/            YAML config + baked builds (vanilla/distro/sealed)
├── gitinfo/           Git branch/dirty HUD
├── knowledge/         Embedded knowledge base (62 docs, 16 categories)
├── llm/
│   ├── providers/     14 LLM providers (OpenAI, Anthropic, Google, Ollama)
│   ├── compaction     3-stage context compression
│   ├── capabilities   Model capability matrix
│   └── cache          Prompt caching infrastructure
├── mcp/               MCP client (JSON-RPC 2.0, stdio/SSE)
├── multi/             Multi-agent orchestration (review/consensus/scan)
├── session/           SQLite (sessions, messages, memories, usage)
├── skills/            SKILL.md loader (agentskills.io compatible)
├── tools/             23+ built-in tools
└── ui/                TUI components, themes, i18n, command palette

knowledge/docs/        62 embedded docs (React, Spring, SQL, Terminal, CSS...)
web/                   Companion dashboard (HTML/CSS/JS + SSE)
```

---

## Supported Providers

| Provider | Type | Default Endpoint |
|----------|------|-----------------|
| Ollama | Local | `localhost:11434` |
| OpenAI | Cloud | `api.openai.com` |
| Anthropic | Cloud | `api.anthropic.com` |
| Google | Cloud | `generativelanguage.googleapis.com` |
| Novita | Cloud | `api.novita.ai` |
| OpenRouter | Cloud | `openrouter.ai` |
| DeepSeek | Cloud | `api.deepseek.com` |
| Groq | Cloud | `api.groq.com` |
| Together | Cloud | `api.together.xyz` |
| Fireworks | Cloud | `api.fireworks.ai` |
| Mistral | Cloud | `api.mistral.ai` |
| vLLM | Local | Custom URL |
| LM Studio | Local | Custom URL |
| Custom | Any | Custom URL |

---

## 한국어 가이드

hanimo(하니모)는 Go로 만든 터미널 기반 AI 코딩 에이전트입니다. 단일 바이너리로 배포되며, 런타임 의존성 없이 14개 이상의 LLM 프로바이더를 지원합니다.

### 주요 기능

- **14개 이상 LLM 프로바이더** — Ollama, OpenAI, Anthropic, Google, Novita, OpenRouter, DeepSeek, Groq 등
- **3가지 모드** — Super (만능), Dev (자율 코딩), Plan (분석 전용)
- **30개 이상 슬래시 명령** — `/auto`, `/plan`, `/multi`, `/companion`, `/copy`, `/undo` 등
- **멀티에이전트** — 듀얼 모델 오케스트레이션 (Review, Consensus, Scan, Auto 전략)
- **23개 이상 내장 도구** — 파일 편집, 셸, git, 검색, 프로젝트 감지, 진단 등
- **보안** — 시크릿 탐지, 민감 파일 차단, 자동 스냅샷, read-before-write 강제
- **Hash-Anchored 편집** — 라인별 MD5 해시로 파일 손상 방지
- **Smart Compaction** — 3단계 컨텍스트 압축 (스닙 → 마이크로 → LLM 요약)
- **컨텍스트 추적** — 상태바에 ctx:XX% 실시간 표시, 90%에서 자동 압축
- **Git HUD** — 브랜치명 + dirty 상태 표시
- **입력 히스토리** — ↑/↓ 화살표로 이전 입력 탐색
- **컴패니언 대시보드** — 브라우저 기반 실시간 SSE 대시보드
- **62개 내장 지식 문서** — React, Spring, SQL, Terminal, CSS, Testing 등
- **스킬 시스템** — SKILL.md 로더 (agentskills.io 호환)
- **SQLite 세션** — 대화 저장/불러오기/검색/분기
- **5가지 테마** — honey, ocean, dracula, nord, forest
- **단일 바이너리** — ~20MB, 크로스 플랫폼 (macOS, Linux, Windows)

### 빠른 시작

```bash
# 소스에서 빌드
git clone https://github.com/flykimjiwon/hanimo.git
cd hanimo
go build -o hanimo ./cmd/hanimo
./hanimo

# Ollama 사용 (API 키 불필요)
ollama pull qwen3:8b
hanimo

# 프로바이더 지정
hanimo -p openai -m gpt-4o
hanimo -p anthropic -m claude-sonnet-4
hanimo -p deepseek -m deepseek-chat
```

### 단축키

| 키 | 동작 |
|----|------|
| `Enter` | 메시지 전송 |
| `Shift+Enter` | 줄바꿈 |
| `Tab` | 모드 전환 (Super → Dev → Plan) |
| `↑/↓` | 입력 히스토리 |
| `Ctrl+K` | 커맨드 팔레트 |
| `Esc` | 메뉴 / 스트리밍 취소 |
| `Ctrl+C` | 종료 |

### 슬래시 커맨드

| 커맨드 | 설명 |
|--------|------|
| `/new` | 새 세션 시작 |
| `/auto [작업]` | 자율 모드 (최대 200회 반복) |
| `/plan [작업]` | 실행 계획 생성 |
| `/multi [전략]` | 멀티에이전트 (on/off/review/consensus/scan/auto) |
| `/git` | Git 브랜치/상태 표시 |
| `/diff` | Git diff 표시 |
| `/init` | .hanimo.md 프로젝트 프로파일 생성 |
| `/copy` | 마지막 AI 응답 클립보드 복사 |
| `/undo [N\|list]` | 파일 변경 되돌리기 |
| `/compact` | 수동 컨텍스트 압축 |
| `/companion` | 브라우저 대시보드 실행 |
| `/save [이름]` | 세션 저장 |
| `/load` | 세션 불러오기 |
| `/search [키워드]` | 세션 검색 |
| `/export` | 세션 마크다운 내보내기 |
| `/model [이름]` | 모델 전환 |
| `/provider [이름]` | 프로바이더 전환 |
| `/theme [이름]` | 테마 변경 |
| `/lang` | 언어 전환 (한국어/영어) |
| `/help` | 도움말 |
| `/version` | 버전 정보 |
| `/exit` | 종료 |

### 설정

설정 파일: `~/.hanimo/config.yaml`

```yaml
default:
  provider: ollama
  model: qwen3:8b

providers:
  ollama: {}
  novita:
    api_key: "nvt-..."
  anthropic:
    api_key: "sk-ant-..."
```

### 프로젝트 컨텍스트

프로젝트 루트에 `.hanimo.md` 파일을 만들거나 `/init`으로 자동 생성:

```markdown
# 내 프로젝트
- 프레임워크: Next.js 15 App Router
- 언어: TypeScript strict
- 테스트: Vitest
```

---

## Contributing

Contributions welcome! Please read our contributing guide before submitting PRs.

```bash
# Development
git clone https://github.com/flykimjiwon/hanimo.git
cd hanimo
make build      # Build binary
make test       # Run tests
make lint       # Run go vet
```

---

## License

MIT
