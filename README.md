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
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-f5a623?style=flat-square" alt="License"></a>
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

### Layer 6: Auto 모드 루프 (자율 실행)

```
/auto "테스트 코드 작성해"
      │
      ├─ m.autoMode = true
      ├─ m.autoTask = "테스트 코드 작성해"
      ├─ 시스템 프롬프트에 AutoPromptSuffix 주입:
      │   "작업 완료 시 [AUTO_COMPLETE], 사용자 입력 필요 시 [AUTO_PAUSE]"
      │
      └─ sendMessage(autoTask) ──→ LLM 스트림 시작
                                        │
      ┌─────────────────────────────────┘
      │
      ▼
┌─ 스트림 완료 후 Auto 체크 ───────────────────────┐
│                                                   │
│  agents.CheckAutoMarkers(streamBuf):              │
│                                                   │
│  [AUTO_COMPLETE] 감지?                            │
│  ├─ YES → autoMode=false, 시스템 프롬프트 복원    │
│  │        "작업 완료" 메시지 표시                  │
│  │        → 종료 (사용자 입력 대기)               │
│  │                                                │
│  [AUTO_PAUSE] 감지?                               │
│  ├─ YES → autoMode=false                          │
│  │        "일시정지 — 사용자 입력 필요" 표시       │
│  │        → 종료 (사용자 입력 대기)               │
│  │                                                │
│  [TASK_COMPLETE] 감지? (Deep Agent)               │
│  ├─ YES → autoMode=false                          │
│  │        → 종료                                  │
│  │                                                │
│  마커 없음?                                       │
│  └─ toolIter < MaxAutoIterations (200)?           │
│     ├─ YES → continueAfterTools()                 │
│     │        → 다시 LLM 스트림 (자동 반복)        │
│     └─ NO  → "tool loop limit" 메시지, 강제 종료  │
│                                                   │
│  ※ pendingQueue에 메시지가 있으면 다음 자동 전송  │
└──────────────────────────────────────────────────┘
```

### Layer 7: Plan 모드 (계획 → 승인 → 단계별 실행)

```
/plan "API 엔드포인트 리팩토링"
      │
      ├─ planAwaitingApproval = true
      ├─ 시스템 프롬프트에 PlanPromptPrefix 주입:
      │   "JSON 형식으로 단계별 계획을 출력하세요"
      │
      └─ sendMessage(planPrompt) ──→ LLM이 JSON 계획 생성
                                        │
      ┌─────────────────────────────────┘
      │
      ▼
┌─ 계획 파싱 ──────────────────────────────────────┐
│                                                   │
│  agents.ParsePlan(streamBuf):                     │
│  ├─ JSON 추출 (```json ... ``` 펜스 제거)         │
│  ├─ Plan { Task, Goal, Steps[] } 구조 파싱        │
│  │                                                │
│  성공?                                            │
│  ├─ YES → m.activePlan = plan                     │
│  │        plan.Render() → 단계 목록 표시          │
│  │        "/approve 로 승인 · /reject 로 거부"    │
│  │        → 사용자 입력 대기                      │
│  └─ NO  → "계획 파싱 실패" 에러                   │
└──────────────────────┬───────────────────────────┘
                       │
      /approve 입력    │    /reject 입력
      ┌────────────────┼────────────────┐
      │                                 │
      ▼                                 ▼
┌─ 실행 시작 ──────┐          계획 폐기, 재작성 가능
│                   │
│  planExecuting    │
│  = true           │
│                   │
│  executeNext      │
│  PlanStep() ──────┼──→ Step N 실행
│                   │         │
│  ┌────────────────┘         ▼
│  │
│  │  step.Status = "in_progress"
│  │  ExecutePromptPrefix로 LLM에 지시:
│  │  "전체 계획 중 Step N/M: {title}"
│  │
│  │  sendMessage(prompt) → LLM 스트림
│  │         │
│  │         ▼
│  │  agents.CheckPlanMarkers(streamBuf):
│  │  ├─ [STEP_COMPLETE] → step.Status="completed"
│  │  │                    current++
│  │  │                    → executeNextPlanStep() (다음 단계)
│  │  ├─ [STEP_FAILED]  → step.Status="failed"
│  │  │                    plan.Status="failed"
│  │  │                    → 실행 중단
│  │  └─ [PLAN_REVISE]  → 계획 재파싱 + 처음부터 재실행
│  │
│  │  모든 단계 완료?
│  │  └─ plan.Status="completed", "모든 단계 완료"
│  │
└──┘
```

### Layer 8: Compaction 3단계 (컨텍스트 압축)

```
llm.Compact(history) — sendMessage()마다 자동 호출
      │
      ▼
┌─ Stage 1: Snip (대량 메시지) ────────────────────┐
│                                                   │
│  조건: len(history) >= 40                         │
│                                                   │
│  최근 10개 메시지 보존 (마지막 대화 보호)         │
│  나머지 중 role=tool && len(content) > 200자:     │
│  → "[snipped: N lines]"로 교체                    │
│                                                   │
│  효과: 오래된 도구 결과 제거, 대화 흐름 유지      │
└──────────────────────┬───────────────────────────┘
                       │
                       ▼
┌─ Stage 2: Micro (개별 메시지 절단) ──────────────┐
│                                                   │
│  조건: 개별 메시지 > 4000자                       │
│                                                   │
│  처리: 앞 2000자 + "[truncated]" + 뒤 2000자     │
│                                                   │
│  효과: 초대형 코드 블록, 로그 출력 축소           │
└──────────────────────┬───────────────────────────┘
                       │
                       ▼
┌─ Stage 3: LLM 요약 (자동 트리거) ────────────────┐
│                                                   │
│  조건: ui.ShouldAutoCompact(ctx%) → ctx >= 90%    │
│  또는: /compact 수동 호출                         │
│                                                   │
│  처리: LLM에게 전체 대화 요약 요청                │
│  → 원래 history를 [system + summary + recent]로   │
│     대폭 축소                                     │
│                                                   │
│  효과: 장기 세션 크래시 방지                      │
│                                                   │
│  Context % 임계값:                                │
│  ├─ < 70%   → Normal (회색)                       │
│  ├─ 70-79%  → Warn (노란색) — 사용자 반응 시간   │
│  ├─ 80-89%  → Critical (빨간색) — 압축 권장      │
│  └─ >= 90%  → Auto-compact 트리거                 │
└──────────────────────────────────────────────────┘
```

### Layer 9: Knowledge 주입 파이프라인

```
sendMessage(input) 중 knowledge 주입
      │
      ▼
┌─ knowledgeInj.Inject(mode, input) ───────────────┐
│                                                   │
│  1단계: 키워드 매칭                               │
│  ├─ 입력에서 키워드 추출                          │
│  ├─ knowledge/index.json의 keywords와 비교        │
│  └─ 매칭 문서 후보 선정                           │
│                                                   │
│  2단계: BM25 스코어링                             │
│  ├─ 후보 문서들의 관련도 점수 계산                │
│  └─ 상위 N개 선택                                 │
│                                                   │
│  3단계: 토큰 예산 내 선택                         │
│  ├─ maxTokenBudget = 8192                         │
│  ├─ 높은 점수 순으로 예산 내에서 포함             │
│  └─ Tier 우선순위: Tier0 > Tier1 > Tier2 > Tier3 │
│                                                   │
│  소스 (3곳 병합):                                 │
│  ├─ knowledge/docs/  — 62개 내장 문서 (embed.FS)  │
│  │   react/, spring/, sql/, terminal/, css/...    │
│  ├─ .hanimo/knowledge/ — 프로젝트 로컬 문서       │
│  └─ ~/.hanimo/knowledge/ — 글로벌 사용자 문서     │
│                                                   │
│  결과: knowledgeCtx 문자열                        │
│  → 시스템 프롬프트에 주입:                        │
│    SystemPrompt(mode) + projectCtx + knowledgeCtx │
└──────────────────────────────────────────────────┘
```

### Layer 10: 위험 명령 확인 플로우

```
LLM이 shell_exec("rm -rf build/") 도구 호출
      │
      ▼
┌─ streamChunkMsg (done=true, toolCalls 있음) ─────┐
│                                                   │
│  for each toolCall:                               │
│    if tc.Name == "shell_exec":                    │
│      tools.IsDangerous(cmdStr)?                   │
│                                                   │
│  IsDangerous 패턴:                                │
│  ├─ rm -rf, rm -r (파일 삭제)                     │
│  ├─ sudo (권한 상승)                              │
│  ├─ git push --force (히스토리 파괴)              │
│  ├─ git reset --hard (변경 폐기)                  │
│  ├─ dd if= (디스크 덮어쓰기)                      │
│  ├─ chmod 777 (과도한 권한)                       │
│  ├─ curl.*| sh (원격 스크립트 실행)               │
│  └─ 기타 위험 패턴...                             │
│                                                   │
│  위험 감지됨?                                     │
│  ├─ YES:                                          │
│  │  m.dangerCmd = "rm -rf build/"                 │
│  │  m.dangerReason = "recursive delete"           │
│  │  m.dangerSelected = 1 (기본: Deny)             │
│  │                                                │
│  │  모든 toolCall에 DEFERRED 결과 삽입            │
│  │  → 위험 확인 오버레이 표시                     │
│  │                                                │
│  │  ┌─────────────────────────────┐               │
│  │  │  ⚠ Dangerous command:      │               │
│  │  │  rm -rf build/              │               │
│  │  │  Reason: recursive delete   │               │
│  │  │                             │               │
│  │  │  [Allow]    [Deny]          │               │
│  │  └─────────────────────────────┘               │
│  │                                                │
│  │  사용자 선택:                                  │
│  │  ├─ Allow → resolveDanger(true)                │
│  │  │   실제 도구 실행 → history에 결과 추가      │
│  │  │   → continueAfterTools() (LLM 재개)        │
│  │  └─ Deny  → resolveDanger(false)               │
│  │     "사용자가 거부" 메시지 → LLM에 전달        │
│  │     → continueAfterTools() (LLM이 대안 모색)  │
│  │                                                │
│  └─ NO: 바로 도구 실행 (Layer 3의 ③)             │
└──────────────────────────────────────────────────┘
```

### Layer 11: 세션 영속화 타이밍

```
┌─ SQLite 저장 시점 ───────────────────────────────┐
│                                                   │
│  앱 시작:                                         │
│  ├─ session.InitDB() → ~/.hanimo/sessions.db      │
│  ├─ session.CreateSession() → 새 세션 ID 발급     │
│  └─ system prompt 저장                            │
│                                                   │
│  메시지 전송 시:                                  │
│  ├─ user 메시지 → session.SaveMessage()           │
│  ├─ assistant 응답 → session.SaveMessage()        │
│  └─ tool 결과 → session.SaveMessage()             │
│                                                   │
│  /save [name]:                                    │
│  └─ session.UpdateSessionName() → 이름 지정 저장  │
│                                                   │
│  /load [prefix]:                                  │
│  └─ session.LoadSession() → 전체 대화 복원        │
│                                                   │
│  /search [keyword]:                               │
│  └─ session.SearchSessions() → 전문검색           │
│                                                   │
│  DB 스키마 (db/schema.sql):                       │
│  ├─ sessions: id, name, project_dir, provider,    │
│  │            model, mode, created_at, updated_at │
│  ├─ messages: session_id, role, content,          │
│  │            tool_calls, tool_result, tokens     │
│  ├─ memories: project_dir, key, value, source     │
│  ├─ usage_log: session_id, provider, model,       │
│  │             tokens_in, tokens_out, cost_usd    │
│  └─ mcp_servers: name, transport, command, url    │
└──────────────────────────────────────────────────┘
```

### Layer 12: Multi-Agent 오케스트레이션

```
/multi review  (또는 consensus, scan, auto)
      │
      ▼
┌─ 전략 선택 ──────────────────────────────────────┐
│                                                   │
│  /multi auto 일 때:                               │
│  ├─ strategy.go: DetectStrategy(input)            │
│  ├─ 리뷰 키워드? → StrategyReview                 │
│  ├─ 비교 키워드? → StrategyConsensus              │
│  ├─ 대량 파일?   → StrategyScan                   │
│  └─ 기본        → StrategyReview                  │
│                                                   │
│  사용 모델:                                       │
│  ├─ Agent1 (Super model) — 주 실행, 쓰기 권한     │
│  └─ Agent2 (Dev model)   — 보조, 읽기 전용        │
└──────────────────────┬───────────────────────────┘
                       │
          ┌────────────┼────────────┐
          ▼            ▼            ▼
     [Review]    [Consensus]    [Scan]

┌─ Review 전략 ────────────────────────────────────┐
│  Agent1: 코드 생성/수정 (전체 도구 접근)          │
│       ↓ 결과                                      │
│  Agent2: 읽기 전용 리뷰 (검토 의견 생성)          │
│       ↓                                           │
│  merge.go: 두 결과 합성 (LLM 기반)               │
└──────────────────────────────────────────────────┘

┌─ Consensus 전략 ─────────────────────────────────┐
│  Agent1: 독립 실행 ──┐                            │
│                       ├─→ merge.go: LLM 합성      │
│  Agent2: 독립 실행 ──┘   "두 응답을 비교 분석"    │
│                                                   │
│  ※ 병렬 실행 (goroutine)                         │
│  ※ 실시간 AgentProgress 이벤트 전송              │
└──────────────────────────────────────────────────┘

┌─ Scan 전략 ──────────────────────────────────────┐
│  파일/함수 단위 파티셔닝:                         │
│  ├─ Agent1: src/api/ + src/models/ 담당           │
│  └─ Agent2: src/handlers/ + src/utils/ 담당       │
│                                                   │
│  ※ 각 에이전트에 "당신은 AgentN, 이 영역 담당"   │
│     힌트 주입                                     │
│  ※ 결과 합산 후 LLM 합성                         │
└──────────────────────────────────────────────────┘

공통: 각 에이전트 최대 20 tool iterations
      실시간 progress → companion SSE로 전송
      합성 실패 시 Agent1 결과만 사용 (graceful)
```

---

## Desktop IDE — `hanimo-code-desktop/`

Wails(Go + React/TS) 기반 데스크톱 IDE. CLI와 같은 `~/.hanimo/config.yaml`,
같은 SKILL/MCP 자산을 공유하며 Honey 팔레트 + 8 테마 + 14 Activity 아이콘
+ hash-anchor gutter로 *"Agent can't silently overwrite your edits"* 라는
브랜드 약속을 시각적으로 보여준다.

- **현재 진척**: Phase 0~13 + 리뷰픽스 (24 커밋)
- **빌드 상태**: `vite 1532 KiB` · `go test 49/49` · TS clean
- **자세한 가이드**: [`hanimo-code-desktop/README.md`](hanimo-code-desktop/README.md)
- **다음 세션 진입점**: [`docs/SESSION-2026-04-25-RESUME.md`](docs/SESSION-2026-04-25-RESUME.md)

```bash
cd hanimo-code-desktop
go install github.com/wailsapp/wails/v2/cmd/wails@latest
go mod download && (cd frontend && npm ci)
wails dev    # hot reload · localhost:34115에 devtools
wails build  # macOS .app / Linux binary / Windows .exe
```

데스크톱 IDE는 CLI와 마찬가지로 **완전 무료 OSS** — 유료화 계획 없음.

---

## Architecture

```
cmd/hanimo/            CLI entry point
hanimo-code-desktop/   Wails IDE (Go + React/TS) — Phase 0-13 완성
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

## Name & Marks

The names *hanimo* / *Hanimo Code* (하니모 코드) and the hanimo logo are
**unregistered marks of Kim Jiwon (김지원)**. Apache 2.0 (Section 6) does
not grant trademark rights — forks are welcome but must use a distinct
product name. Nominative fair use (blog posts, comparisons, citations,
"Powered by hanimo" badges) is explicitly permitted.

See [`docs/policy/trademark-and-naming.md`](docs/policy/trademark-and-naming.md)
for the full policy.

---

## Policies

| Policy | Document |
|---|---|
| Copyright · source headers · GPG signing | [`docs/policy/copyright.md`](docs/policy/copyright.md) |
| Trademark · forking · naming | [`docs/policy/trademark-and-naming.md`](docs/policy/trademark-and-naming.md) |
| Telemetry · privacy · air-gap | [`docs/policy/telemetry-and-privacy.md`](docs/policy/telemetry-and-privacy.md) |
| LTS · on-premises · enterprise | [`docs/policy/lts-onprem.md`](docs/policy/lts-onprem.md) |
| Contributing · DCO · CLA grant | [`CONTRIBUTING.md`](CONTRIBUTING.md) |
| **Brand · 디자인 · 마스코트** | [`docs/branding/BRAND.md`](docs/branding/BRAND.md) — Modol the Honey-Bee Bichon |

Index: [`docs/policy/README.md`](docs/policy/README.md) · [`docs/branding/README.md`](docs/branding/README.md)

---

## License

Apache License 2.0 — see [LICENSE](LICENSE) and [NOTICE](NOTICE).

Copyright 2025-2026 Kim Jiwon (김지원). All rights reserved.
