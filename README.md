<p align="center">
  <img src="docs/images/hanimo-screenshot.png" alt="hanimo screenshot" width="500">
</p>

<h1 align="center">hanimo</h1>

<p align="center">
  <strong>AI Coding Agent for Open Source</strong><br>
  <em>오픈소스를 위한 AI 코딩 에이전트</em>
</p>

<p align="center">
  <a href="#english">English</a> | <a href="#한국어-가이드">한국어</a>
</p>

---

## English

hanimo is a terminal-based AI coding agent built with Go. Ships as a single binary with zero runtime dependencies, supporting 14+ LLM providers out of the box.

### Features

- **14+ LLM Providers** — Ollama, OpenAI, Anthropic, Google, Novita, OpenRouter, DeepSeek, Groq, Together, Fireworks, Mistral, vLLM, LM Studio, Custom
- **3 Modes** — Super (all-purpose), Dev (coding), Plan (read-only analysis)
- **Hash-Anchored Editing** — MD5 hash per line prevents stale-edit corruption
- **Smart Compaction** — 3-stage context compression (snip → micro → LLM summary)
- **Autonomous Mode** — `/auto` for self-driving task completion (up to 20 iterations)
- **SQLite Sessions** — Save, load, search, fork conversation sessions
- **Project Memory** — Long-term per-project knowledge stored in SQLite
- **MCP Support** — Extend tools via Model Context Protocol (stdio transport)
- **Built-in Tools (14)** — File R/W/E/D, hash-anchored edit, shell, grep, glob, git, diagnostics
- **Command Palette** — Ctrl+K fuzzy search across all commands
- **5 Themes** — honey (default), ocean, dracula, nord, forest
- **Single Binary** — ~20MB, cross-platform (macOS, Linux, Windows)

### Quick Start

#### Option 1: Download Binary

```bash
# macOS (Apple Silicon)
curl -L -o hanimo https://github.com/flykimjiwon/hanimo/releases/latest/download/hanimo-darwin-arm64
chmod +x hanimo
sudo mv hanimo /usr/local/bin/

# macOS (Intel)
curl -L -o hanimo https://github.com/flykimjiwon/hanimo/releases/latest/download/hanimo-darwin-amd64

# Linux (x64)
curl -L -o hanimo https://github.com/flykimjiwon/hanimo/releases/latest/download/hanimo-linux-amd64

# Windows
# Download hanimo-windows-amd64.exe from releases page
```

#### Option 2: Build from Source

```bash
git clone https://github.com/flykimjiwon/hanimo.git
cd hanimo
go build -o hanimo ./cmd/hanimo
./hanimo
```

#### Run

```bash
# Default: connects to Ollama on localhost (no API key needed)
ollama pull qwen3:8b
hanimo

# With specific provider & model
hanimo -p openai -m gpt-4o
hanimo -p anthropic -m claude-sonnet-4
hanimo -p novita -m openai/gpt-oss-120b
hanimo -p deepseek -m deepseek-chat

# With custom endpoint (vLLM, LM Studio, etc.)
hanimo -u http://localhost:8000/v1 -m my-model

# Start in specific mode
hanimo --mode dev
hanimo --mode plan
```

### Configuration

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

Environment variables override config file:

| Variable | Description |
|----------|-------------|
| `HANIMO_API_KEY` | API key |
| `HANIMO_API_BASE_URL` | API endpoint URL |
| `HANIMO_MODEL_SUPER` | Model for Super/Plan mode |
| `HANIMO_MODEL_DEV` | Model for Dev mode |

### Project Context

Create a `.hanimo.md` file in your project root to provide project-specific instructions. hanimo automatically loads it into the system prompt.

```markdown
# My Project
- Framework: Next.js 15 with App Router
- Language: TypeScript strict mode
- Testing: Vitest
- Always use Korean for responses
```

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Enter` | Send message |
| `Shift+Enter` | New line |
| `Tab` | Switch mode (Super → Dev → Plan) |
| `Ctrl+K` | Command palette |
| `Esc` | Cancel streaming |
| `Ctrl+L` | Clear screen |
| `Ctrl+C` | Exit |
| `Alt+↑/↓` | Scroll |

### Slash Commands

| Command | Description |
|---------|-------------|
| `/help` | Show help |
| `/model [name]` | Switch model |
| `/provider [name]` | Switch provider |
| `/auto [task]` | Autonomous mode (up to 20 iterations) |
| `/save [name]` | Save session |
| `/load` | Load session |
| `/search [keyword]` | Search past sessions |
| `/remember [text]` | Save project memory |
| `/memories` | Show saved memories |
| `/diagnostics [file]` | Run code diagnostics (go vet, tsc, eslint, ruff) |
| `/usage` | Token usage & estimated cost |
| `/config` | Show current configuration |
| `/theme [name]` | Change theme (honey, ocean, dracula, nord, forest) |
| `/clear` | Clear conversation |

### Architecture

```
cmd/hanimo/          Entry point
internal/
├── app/             Bubble Tea TUI
├── config/          YAML config + setup wizard
├── llm/
│   ├── providers/   14 LLM providers (OpenAI-compat, Anthropic, Google, Ollama)
│   ├── compaction   Smart context compression
│   └── capabilities Model capability matrix
├── tools/           14 built-in tools (file, shell, search, git, hashline, diagnostics)
├── session/         SQLite (sessions, memory, usage)
├── agents/          Autonomous mode (/auto)
├── mcp/             MCP client (JSON-RPC 2.0, stdio)
├── ui/              TUI components, themes, command palette
└── knowledge/       Embedded knowledge base (Go, JS, TS, React, CSS, Python...)

knowledge/           Docs embedded into binary via embed.FS
```

### Supported Providers

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

- **14개 이상 LLM 프로바이더** — Ollama, OpenAI, Anthropic, Google, Novita, OpenRouter, DeepSeek, Groq, Together, Fireworks, Mistral, vLLM, LM Studio, Custom
- **3가지 모드** — Super (만능), Dev (코딩 전용), Plan (분석 전용, 읽기만)
- **Hash-Anchored 편집** — 라인별 MD5 해시로 파일 손상 방지
- **Smart Compaction** — 3단계 컨텍스트 압축 (스닙 → 마이크로 → LLM 요약)
- **자율 모드** — `/auto`로 최대 20회 반복 자동 실행
- **SQLite 세션** — 대화 저장, 불러오기, 검색, 분기(fork)
- **프로젝트 메모리** — SQLite에 프로젝트별 장기 기억 저장
- **MCP 지원** — Model Context Protocol로 외부 도구 확장
- **내장 도구 14종** — 파일 읽기/쓰기/편집/삭제, 셸, grep, glob, git, 진단
- **커맨드 팔레트** — Ctrl+K로 모든 명령어 퍼지 검색
- **5가지 테마** — honey (기본), ocean, dracula, nord, forest
- **단일 바이너리** — ~20MB, 크로스 플랫폼 (macOS, Linux, Windows)

### 빠른 시작

#### 방법 1: 바이너리 다운로드

```bash
# macOS (Apple Silicon)
curl -L -o hanimo https://github.com/flykimjiwon/hanimo/releases/latest/download/hanimo-darwin-arm64
chmod +x hanimo
sudo mv hanimo /usr/local/bin/
```

#### 방법 2: 소스에서 빌드

```bash
git clone https://github.com/flykimjiwon/hanimo.git
cd hanimo
go build -o hanimo ./cmd/hanimo
./hanimo
```

#### 실행

```bash
# 기본: Ollama에 자동 연결 (API 키 불필요)
ollama pull qwen3:8b
hanimo

# 프로바이더 지정
hanimo -p openai -m gpt-4o
hanimo -p anthropic -m claude-sonnet-4
hanimo -p novita -m openai/gpt-oss-120b
hanimo -p deepseek -m deepseek-chat

# 커스텀 엔드포인트 (vLLM, LM Studio 등)
hanimo -u http://localhost:8000/v1 -m my-model

# 모드 지정
hanimo --mode dev    # 코딩 모드
hanimo --mode plan   # 분석 모드
```

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
  openrouter:
    api_key: "sk-or-..."
  anthropic:
    api_key: "sk-ant-..."
  google:
    api_key: "AIza..."
```

환경변수로도 설정 가능:

| 환경변수 | 설명 |
|---------|------|
| `HANIMO_API_KEY` | API 키 |
| `HANIMO_API_BASE_URL` | API 엔드포인트 URL |
| `HANIMO_MODEL_SUPER` | Super/Plan 모드 모델 |
| `HANIMO_MODEL_DEV` | Dev 모드 모델 |

### 프로젝트 컨텍스트

프로젝트 루트에 `.hanimo.md` 파일을 만들면 프로젝트 맞춤 지시를 줄 수 있습니다. hanimo가 자동으로 시스템 프롬프트에 주입합니다.

```markdown
# 내 프로젝트
- 프레임워크: Next.js 15 App Router
- 언어: TypeScript strict
- 테스트: Vitest
- 한국어로 응답
```

### 단축키

| 키 | 동작 |
|----|------|
| `Enter` | 메시지 전송 |
| `Shift+Enter` | 줄바꿈 |
| `Tab` | 모드 전환 (Super → Dev → Plan) |
| `Ctrl+K` | 커맨드 팔레트 |
| `Esc` | 스트리밍 취소 |
| `Ctrl+L` | 화면 초기화 |
| `Ctrl+C` | 종료 |
| `Alt+↑/↓` | 스크롤 |

### 슬래시 커맨드

| 커맨드 | 설명 |
|--------|------|
| `/help` | 도움말 |
| `/model [이름]` | 모델 전환 |
| `/provider [이름]` | 프로바이더 전환 |
| `/auto [작업]` | 자율 모드 (최대 20회 반복) |
| `/save [이름]` | 세션 저장 |
| `/load` | 세션 불러오기 |
| `/search [키워드]` | 과거 세션 검색 |
| `/remember [텍스트]` | 프로젝트 메모리 저장 |
| `/memories` | 저장된 메모리 조회 |
| `/diagnostics [파일]` | 코드 진단 (go vet, tsc, eslint, ruff) |
| `/usage` | 토큰 사용량 및 예상 비용 |
| `/config` | 현재 설정 표시 |
| `/theme [이름]` | 테마 변경 (honey, ocean, dracula, nord, forest) |
| `/clear` | 대화 초기화 |

### 프로바이더 목록

| 프로바이더 | 유형 | 기본 엔드포인트 |
|-----------|------|----------------|
| Ollama | 로컬 | `localhost:11434` |
| OpenAI | 클라우드 | `api.openai.com` |
| Anthropic | 클라우드 | `api.anthropic.com` |
| Google | 클라우드 | `generativelanguage.googleapis.com` |
| Novita | 클라우드 | `api.novita.ai` |
| OpenRouter | 클라우드 | `openrouter.ai` |
| DeepSeek | 클라우드 | `api.deepseek.com` |
| Groq | 클라우드 | `api.groq.com` |
| Together | 클라우드 | `api.together.xyz` |
| Fireworks | 클라우드 | `api.fireworks.ai` |
| Mistral | 클라우드 | `api.mistral.ai` |
| vLLM | 로컬 | 직접 지정 |
| LM Studio | 로컬 | 직접 지정 |
| Custom | 자유 | 직접 지정 |

---

## License

MIT
