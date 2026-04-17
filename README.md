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

## Internal Process Flow (시각화)

hanimo의 전체 동작을 레이어별로 뜯어놓은 다이어그램입니다. "이 부분 고쳐" 할 때 참고하세요.

### Layer 1: 전체 아키텍처 오버뷰

```
┌─────────────────────────────────────────────────────────┐
│                    사용자 (터미널)                         │
│  키보드 입력 → [textarea] → Enter/Tab/Esc/Ctrl+K        │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│              Bubble Tea 이벤트 루프                       │
│                                                         │
│  Init() → Update(msg) → View() → 화면 렌더링 → 반복     │
│                                                         │
│  msg 종류:                                               │
│  ├─ tea.KeyPressMsg    (키보드)                          │
│  ├─ tea.PasteMsg       (붙여넣기)                        │
│  ├─ tea.WindowSizeMsg  (창 크기)                         │
│  ├─ tea.MouseWheelMsg  (마우스 스크롤)                   │
│  ├─ streamTickMsg      (150ms 스피너)                    │
│  ├─ streamChunkMsg     (LLM 응답 조각)                   │
│  └─ toolResultMsg      (도구 실행 결과)                   │
└──────────────────────┬──────────────────────────────────┘
                       │
          ┌────────────┼────────────┐
          ▼            ▼            ▼
     [키 처리]    [스트림 처리]  [View 렌더링]
```

### Layer 2: 키 입력 → 메시지 전송 파이프라인

```
사용자 키 입력
      │
      ▼
┌─ 오버레이 체크 ──────────────────────────────────┐
│                                                   │
│  m.inSetup?        → viewSetup() (API 키 입력)   │
│  m.showPalette?    → updatePalette() (Ctrl+K)    │
│  m.showMenu?       → updateMenu() (Esc 메뉴)     │
│  m.askQuestion?    → handleAskUserKey() (선택지)  │
│  m.dangerCmd?      → handleDangerKey() (확인)     │
│                                                   │
│  ※ 오버레이가 활성이면 일반 키 처리 스킵           │
└──────────────────────┬───────────────────────────┘
                       │ (오버레이 없음)
                       ▼
┌─ 스트리밍 중? ───────────────────────────────────┐
│                                                   │
│  YES (m.streaming == true):                       │
│  ├─ Ctrl+C/Esc    → cancelStream()               │
│  ├─ PgUp/PgDn     → viewport 스크롤              │
│  ├─ Enter          → pendingQueue에 큐잉          │
│  └─ 그 외          → 무시                         │
│                                                   │
│  NO (idle):                                       │
│  ├─ Enter          → ①로 이동                     │
│  ├─ Tab            → 모드 전환 (Super/Dev/Plan)   │
│  ├─ ↑/↓            → 입력 히스토리 탐색            │
│  ├─ Ctrl+K         → 커맨드 팔레트 열기            │
│  ├─ Ctrl+B         → 마우스 모드 토글              │
│  ├─ Ctrl+L         → 화면 초기화                   │
│  ├─ Shift+Enter    → 줄바꿈                       │
│  └─ Esc            → 메뉴 열기                    │
└──────────────────────┬───────────────────────────┘
                       │ ① Enter 눌림
                       ▼
┌─ 입력 분기 ──────────────────────────────────────┐
│                                                   │
│  input = textarea.Value()                         │
│                                                   │
│  "/" 으로 시작?                                    │
│  ├─ YES → handleSlashCommand(input)               │
│  │        ├─ /auto, /plan, /multi, /git, /copy... │
│  │        ├─ /commands 에 매칭?                    │
│  │        │   → 커스텀 명령 템플릿을 sendMessage() │
│  │        └─ return (m, nil)                      │
│  │                                                │
│  └─ NO  → sendMessage(input)  ← ②로 이동         │
│                                                   │
│  ※ 입력 히스토리에 저장 (최대 100개)               │
│  ※ textarea 초기화                                │
└──────────────────────┬───────────────────────────┘
                       │ ② sendMessage(input)
                       ▼
```

### Layer 3: LLM 스트리밍 + 도구 실행 루프

```
sendMessage(input)
      │
      ├─ gitInfo 갱신: gitinfo.Fetch(".")
      ├─ UI 메시지 추가: m.msgs ← user message
      ├─ Intent 감지: agents.DetectIntentLocal(input)
      ├─ Knowledge 주입: knowledgeInj.Inject(mode, input)
      ├─ 시스템 프롬프트 재구성: SystemPrompt + projectCtx + knowledgeCtx
      ├─ history에 user 메시지 추가
      ├─ Compact(history): 40개 넘으면 오래된 도구 결과 스닙
      ├─ 상태 초기화: streaming=true, toolIter=0, tokenCount=0
      │
      └─ startStream() ───────────────────────────────────┐
                                                           │
      ┌────────────────────────────────────────────────────┘
      │
      ▼
┌─ StreamChat (goroutine) ─────────────────────────┐
│                                                   │
│  client.StreamChat(ctx, model, history, tools)    │
│  → OpenAI-compatible SSE stream                   │
│  → chunk별로 streamCh 채널에 전송                  │
│                                                   │
│  동시에: streamTickCmd() → 150ms마다 UI 갱신      │
└──────────────────────┬───────────────────────────┘
                       │ streamChunkMsg 수신
                       ▼
┌─ streamChunkMsg 처리 ────────────────────────────┐
│                                                   │
│  에러?          → 에러 메시지 표시, 스트림 종료    │
│                                                   │
│  done == false? → streamBuf에 텍스트 누적         │
│                   tokenCount 증가                 │
│                   waitForNextChunk() (다음 조각)   │
│                                                   │
│  done == true?  → 분기:                           │
│                                                   │
│  ┌─ toolCalls 있음? ─────────────────────────┐   │
│  │                                            │   │
│  │  history ← assistant msg (with toolCalls)  │   │
│  │                                            │   │
│  │  위험 명령 감지?                            │   │
│  │  ├─ YES → dangerCmd 오버레이 (사용자 확인) │   │
│  │  └─ NO  → 도구 실행 ③으로                  │   │
│  │                                            │   │
│  └────────────────────────────────────────────┘   │
│                                                   │
│  ┌─ toolCalls 없음 (일반 완료) ──────────────┐   │
│  │                                            │   │
│  │  history ← assistant msg                   │   │
│  │  m.msgs ← AI 응답 표시                     │   │
│  │  streaming = false                         │   │
│  │                                            │   │
│  │  autoMode?  → 자동 재전송 (최대 200회)     │   │
│  │  planMode?  → executeNextPlanStep()        │   │
│  │  ASK_USER?  → askQuestion 오버레이         │   │
│  │                                            │   │
│  └────────────────────────────────────────────┘   │
└──────────────────────┬───────────────────────────┘
                       │ ③ 도구 실행
                       ▼
┌─ 도구 실행 루프 ─────────────────────────────────┐
│                                                   │
│  for each toolCall:                               │
│  ├─ 반복 감지: 같은 name+args 3회 이상 → 차단    │
│  ├─ 드리프트 감지: 같은 name 5회 연속 → 차단     │
│  ├─ tools.Execute(name, argsJSON)                 │
│  │   ├─ ripgrep 사용 가능? → RipgrepSearch 우선   │
│  │   ├─ file_write? → CheckSensitiveFile          │
│  │   │               → CheckSecrets               │
│  │   │               → CreateSnapshot             │
│  │   │               → os.WriteFile               │
│  │   ├─ file_edit?  → wasRead 확인                │
│  │   │               → CreateSnapshot             │
│  │   │               → strings.Replace            │
│  │   ├─ shell_exec? → 위험 명령 차단              │
│  │   │               → 30초 타임아웃 실행          │
│  │   └─ 기타 23+ 도구...                          │
│  │                                                │
│  ├─ history ← tool result                         │
│  ├─ m.msgs ← 도구 결과 표시                       │
│  └─ toolIter++ (최대 20회, /auto는 200회)         │
│                                                   │
│  toolIter >= 20?                                  │
│  ├─ YES → "Tool iteration limit" 메시지, 종료     │
│  └─ NO  → continueAfterTools()                    │
│           → startStream() (LLM에게 결과 전달)     │
│           → ④ 다시 streamChunkMsg 루프로           │
└──────────────────────────────────────────────────┘
```

### Layer 4: View() 렌더링 스택

```
View() 호출 (매 Update 후)
      │
      ▼
┌─ 화면 구성 (위에서 아래로) ──────────────────────┐
│                                                   │
│  ┌───────────────────────────────────────────┐   │
│  │  viewport (스크롤 가능한 대화 영역)        │   │
│  │  ├─ [SYSTEM] 로고 + 모드 정보              │   │
│  │  ├─ [USER]   사용자 메시지 (파란 블록)     │   │
│  │  ├─ [AI]     AI 응답 (마크다운 렌더링)     │   │
│  │  ├─ [TOOL]   도구 호출 (회색, 접힘)        │   │
│  │  ├─ [SYSTEM] 시스템 메시지 (노란색)        │   │
│  │  └─ ...반복...                             │   │
│  └───────────────────────────────────────────┘   │
│                                                   │
│  ┌───────────────────────────────────────────┐   │
│  │  intent hint (Super 모드, 노란 이탤릭)    │   │
│  │  "Looks like a refactoring request..."     │   │
│  └───────────────────────────────────────────┘   │
│                                                   │
│  ┌───────────────────────────────────────────┐   │
│  │  ASK_USER 블록 (선택지 있을 때만)          │   │
│  │  ├─ 질문 텍스트                            │   │
│  │  └─ [1] 옵션A  [2] 옵션B  [3] 옵션C      │   │
│  └───────────────────────────────────────────┘   │
│                                                   │
│  ┌───────────────────────────────────────────┐   │
│  │  입력 상자 (textarea, 둥근 테두리)         │   │
│  │  │ 여기에 입력하세요...                  │ │   │
│  └───────────────────────────────────────────┘   │
│                                                   │
│  ┌───────────────────────────────────────────┐   │
│  │  상태바 (하단 고정)                        │   │
│  │  Super  gpt-4o  ./myproj  main*           │   │
│  │  ON(16)  [AUTO]  ctx:42%  1234tok  3.2s   │   │
│  │                    Enter  Esc Menu  ^C     │   │
│  └───────────────────────────────────────────┘   │
│                                                   │
│  ┌─ 오버레이 (위에 덮어씌움) ────────────────┐   │
│  │  위험 명령 확인: "rm -rf build/"           │   │
│  │  [Allow]  [Deny]                           │   │
│  ├────────────────────────────────────────────┤   │
│  │  메뉴 오버레이 (Esc):                      │   │
│  │  > Sessions / New / Multi / Git / Help     │   │
│  ├────────────────────────────────────────────┤   │
│  │  커맨드 팔레트 (Ctrl+K):                   │   │
│  │  > /au_  → /auto  /approve                 │   │
│  └───────────────────────────────────────────┘   │
└──────────────────────────────────────────────────┘
```

### Layer 5: 파일별 책임 맵

```
"이 부분 고쳐" 할 때 찾아갈 파일:

입력/이벤트 처리  → internal/app/app.go        Update()
슬래시 명령       → internal/app/app.go        handleSlashCommand()
화면 렌더링       → internal/app/app.go        View()
상태바            → internal/ui/chat.go        RenderStatusBar()
메시지 블록       → internal/ui/chat.go        RenderMessages()
테마/색상         → internal/ui/styles.go
커맨드 팔레트     → internal/ui/palette.go     FuzzyFilter()
메뉴 오버레이     → internal/ui/menu.go
다국어            → internal/ui/i18n.go        T()
컨텍스트 %        → internal/ui/context.go     ContextPercent()

LLM 스트리밍      → internal/llm/client.go     StreamChat()
모델 목록         → internal/llm/models.go
프로바이더 선택   → internal/llm/providers/     registry.go
컨텍스트 압축     → internal/llm/compaction.go  Compact()
시스템 프롬프트   → internal/llm/prompt.go      SystemPrompt()
모델 능력치       → internal/llm/capabilities.go GetCapability()
프롬프트 캐싱     → internal/llm/cache.go

도구 정의         → internal/tools/registry.go  AllTools(), executeInner()
파일 읽기/쓰기    → internal/tools/file.go      FileRead/Write/Edit()
셸 실행           → internal/tools/shell.go     ShellExec()
검색 (grep/glob)  → internal/tools/search.go    GrepSearch(), GlobSearch()
심볼 검색         → internal/tools/symbols.go   SymbolSearch()
ripgrep 폴백      → internal/tools/ripgrep.go   RipgrepSearch()
Git 도구          → internal/tools/git.go       GitStatus/Diff/Log/Commit()
해시라인 편집     → internal/tools/hashline.go  HashlineRead/Edit()
프로젝트 감지     → internal/tools/project.go   DetectProject()
프로젝트 프로파일 → internal/tools/init.go      GenerateProjectProfile()
시크릿 탐지       → internal/tools/secrets.go   CheckSecrets()
스냅샷/undo       → internal/tools/snapshot.go  CreateSnapshot(), UndoLast()
.gitignore 파서   → internal/tools/gitignore.go LoadGitIgnore()
커스텀 명령       → internal/tools/commands.go  LoadCustomCommands()
diff 생성         → internal/tools/diff.go      GenerateUnifiedDiff()
진단              → internal/tools/diagnostics.go RunDiagnostics()

자율 모드         → internal/agents/auto.go     CheckAutoMarkers()
플랜 모드         → internal/agents/plan.go     ParsePlan()
인텐트 감지       → internal/agents/intent.go   DetectIntentLocal()
ASK_USER          → internal/agents/askuser.go

멀티에이전트      → internal/multi/orchestrator.go
전략 선택         → internal/multi/strategy.go
결과 합성         → internal/multi/merge.go

Git HUD           → internal/gitinfo/gitinfo.go Fetch(), Label()
세션 저장         → internal/session/store.go
스킬 로더         → internal/skills/loader.go
MCP 클라이언트    → internal/mcp/client.go
설정              → internal/config/config.go

컴패니언 서버     → internal/companion/server.go
이벤트 브로드캐스트 → internal/companion/hub.go
SSE 이벤트 타입   → internal/companion/events.go
```

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
