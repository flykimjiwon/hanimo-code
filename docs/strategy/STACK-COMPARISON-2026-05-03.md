# Stack Comparison — Go+Bubble Tea / Wails vs React+Ink / Electron / Tauri
**Date:** 2026-05-03
**Context:** hanimo-code (TUI) + hanimo-code-desktop (IDE) 두 surface의 실제 stack vs 오픈소스 진영 표준 비교
**상태:** 실측 + 진영별 레퍼런스 포함

---

## 0. TL;DR

| 질문 | 답 |
|---|---|
| 지금 우리는? | TUI = **Go + Bubble Tea v2 (Charm)** · Desktop = **Wails (Go) + React 18 + TypeScript + CodeMirror + xterm.js** |
| 이게 표준이야? | TUI: AI 에이전트 진영은 React+Ink가 더 표준 / 인프라 CLI 진영은 Go+Charm 표준. Desktop: **Electron이 사실상 표준**, Wails는 IDE급 레퍼런스 거의 없음 |
| 바꿔야 해? | TUI는 안 바꿔도 됨 (charmbracelet/crush와 같은 라인). Desktop은 Wails 유지 시 길을 직접 뚫어야 함 — Electron으로 가면 표준 합류, 그러나 simplicity-first 정책과 충돌 |

---

## 1. 현재 hanimo-code 실측 stack

### 1.1 hanimo-code (TUI)
실측 데이터:

```
module github.com/flykimjiwon/hanimo
go 1.26.1

charm.land/bubbletea/v2  v2.0.2     ← TUI 프레임워크
charm.land/bubbles/v2    v2.1.0     ← 컴포넌트
charm.land/lipgloss/v2   v2.0.2     ← 스타일링
charm.land/glamour/v2    v2.0.0     ← 마크다운 렌더링
sashabaranov/go-openai   v1.41.2    ← LLM SDK
modernc.org/sqlite       v1.48.1    ← 순수 Go SQLite (CGO 없음)
```

| 항목 | 값 |
|---|---|
| Go LOC (cmd + internal) | ~4,660 |
| 단일 바이너리 크기 | 22.7 MB (실측: `hanimo` 바이너리) |
| 외부 런타임 | **없음** — 단일 바이너리 |
| Charm 의존도 | 매우 깊음 (4개 패키지 모두 v2) |
| CGO | **없음** — modernc/sqlite로 우회 |

**핵심 관찰**:
- Charm v2 stack은 2026 기준 가장 최신 라인. 대부분의 OSS는 아직 v1
- modernc/sqlite로 CGO 없이 sqlite를 쓰고 있어 cross-compile 자유도 최상
- **단일 바이너리 22 MB** 는 Node 기반 CLI 평균(50–200 MB 패킹) 대비 절반 이하

### 1.2 hanimo-code-desktop (IDE)
실측 데이터:

```
# Backend (Wails)
github.com/wailsapp/wails/v2     v2.12.0
github.com/sashabaranov/go-openai v1.41.2

# Frontend (Vite + React)
react             ^18.2.0
react-dom         ^18.2.0
typescript        ^4.9.5
vite              ^3.0.7
@vitejs/plugin-react ^2.0.1

# Editor & Terminal
@codemirror/view + 15개 lang-* 패키지   ← Monaco가 아님!
@xterm/addon-fit
@xterm/addon-web-links
xterm             ^5.3.0

# UI
lucide-react      ^1.8.0
react-resizable-panels ^4.10.0
```

| 항목 | 값 |
|---|---|
| Frontend LOC (TS+TSX) | ~6,239 |
| Backend Go 파일 | 약 25개 (app.go, mcp.go, terminal.go, git.go, session.go, providers.go, …) |
| 에디터 엔진 | **CodeMirror 6** (Monaco 아님) |
| 터미널 엔진 | xterm.js |
| 패널 시스템 | react-resizable-panels |
| 아이콘 | lucide-react (메모리 정책: Lucide 단일 source) |

**핵심 관찰 (의외의 발견)**:
- 에디터로 **Monaco가 아니라 CodeMirror 6** 를 쓰고 있음. 이게 큰 의미를 가짐 — 아래 §3.2 참조.
- Wails이라 백엔드는 **Go 그대로 재사용** (LSP/터미널/git/MCP 코드를 hanimo-code TUI와 공유 가능)
- 프론트는 React 18 + TS + Vite로 **VS Code 같은 모듈 구조** 채택 (FileTree, Editor, Terminal, GitPanel, MCPPanel, CommandPalette 등)

---

## 2. TUI 진영: Go+Bubble Tea vs React+Ink (실제 OSS 매핑)

### 2.1 양 진영 대표 OSS

| 진영 | AI 에이전트 | 인프라/생산성 CLI |
|---|---|---|
| **Go + Bubble Tea (Charm)** | charmbracelet/**crush** ⭐ (직접 경쟁자), charmbracelet/mods | gh CLI, **lazygit**, **k9s**, **gum**, **glow**, soft-serve, vhs, freeze |
| **TS + React + Ink** | **claude-code** (Anthropic, 정황상), **gemini-cli** (Google), **opencode**, codex-style 도구들 | Vercel CLI 일부, shopify CLI 일부 |
| **Python + Rich/Textual** | **aider** ⭐ (가장 유명한 AI 코딩 CLI), **continue-cli** | Textual demo apps |
| **Rust + Ratatui** | (AI 사례 적음) | gitui, atuin, bottom |

> AI 코딩 에이전트 진영의 **양대 표준**: Anthropic/Google 라인 = **TS+Ink**, Charm 라인 = **Go+Bubble Tea**.
> 우리(`hanimo`) = Charm 라인. 가장 직접적인 비교 대상은 `charmbracelet/crush`.

### 2.2 항목별 비교

| 항목 | **Go + Bubble Tea (현재 우리)** | **TS + React + Ink** |
|---|---|---|
| **배포 마찰** | 단일 바이너리 (brew/curl, 22 MB) | Node 런타임 필요 (`npx`, `npm i -g`) 또는 pkg/bun으로 50–100 MB |
| **메모리/CPU** | 매우 가벼움 (20–50 MB), 깜빡임 적음 | Node 기준 무거움 (100–300 MB), 큰 출력 시 리렌더 비용 |
| **동시성 모델** | goroutines = LLM stream + file watch + LSP 자연스러움 | Node 이벤트 루프, 무거운 파싱은 worker_threads 필요 |
| **AI/MCP SDK** | Anthropic Go SDK 있음, MCP Go도 있지만 **2nd-class** | **1st-class**: Anthropic SDK, MCP SDK, tree-sitter 모두 npm 표준 |
| **TUI 컴포넌트 풀** | Bubbles + Lip Gloss + Glamour + Huh + Gum + VHS, **폭발적** | Ink 안정적이나 업데이트 느림, 컴포넌트 풀 작음 |
| **컨트리뷰터 풀** | Go 개발자 — 풀 작지만 진중함 | **JS/TS = 압도적 풀**, 프론트 개발자 즉시 합류 |
| **OSS 레퍼런스 (AI agent)** | **crush, mods** | **claude-code, gemini-cli, opencode** |
| **OSS 레퍼런스 (인프라)** | **gh, lazygit, k9s, gum, glow** (압도적) | 거의 없음 |
| **트렌드 (2026)** | 🔥 Charm 진영 떡상, Bubble Tea v2 출시 | 안정 단계, 신선도 떨어짐 |
| **개인 학습 ROI** | Go 한 번 익히면 Wails에도 그대로 활용 | JS 익히면 Ink + Electron 양쪽 활용 |

### 2.3 charmbracelet/crush vs hanimo-code (직접 비교)

`charmbracelet/crush`는 우리의 **가장 직접적인 OSS 경쟁자**. 같은 stack(Go+Bubble Tea), 같은 도메인(AI 코딩 에이전트), Charm이 직접 만든 reference 구현.

| 항목 | crush | hanimo-code |
|---|---|---|
| Stack | Go + Bubble Tea | Go + Bubble Tea v2 ✓ |
| 메이커 | Charmbracelet (TUI 진영 본진) | flykimjiwon |
| 차별화 포인트 | Charm 공식 = 인지도 | **한국어 1st-class + Korea MCP + 신한 폐쇄망(TECHAI)** |
| Lock-in | 없음 | 없음 |

**시사점**: 우리 stack은 정확히 Charm 본진과 같은 라인. 레퍼런스가 이미 검증된 상태. 차별화는 stack이 아니라 **한국어 + 한국 MCP + 폐쇄망 페어** 으로 이미 잡혀 있음.

### 2.4 결론 (TUI)
**유지가 정답**. 다음 근거:
1. 단일 바이너리 22 MB는 명성 공급(install 마찰 0)에 결정적
2. Charm 진영은 떡상 중이고 v2 stack에 이미 진입
3. 직접 비교 대상(crush)이 같은 stack — "왜 Go인가?" 변호 불필요
4. modernc/sqlite + Wails 백엔드 Go 코드 공유 = 두 surface 시너지

**위험**: AI SDK가 Go에서 2nd-class. 새 Anthropic 기능(컴퓨터 사용, prompt caching 옵션 등) 합류가 항상 1–2주 늦음. 직접 unmarshaling이 필요할 때 있음.

---

## 3. Desktop 진영: Wails vs Electron vs Tauri (실제 OSS 매핑)

### 3.1 IDE급 OSS 진영 매핑

| 진영 | IDE/에디터 OSS |
|---|---|
| **Electron + React/TS + Monaco** | **VS Code, Cursor, Void, Trae**, Theia (브라우저 모드도), Atom(EOL), Continue(VS Code 확장) |
| **Electron + 자체 framework** | Atom(killed), Brackets(killed), Light Table(killed) |
| **Tauri + React** | **거의 없음** — Spacedrive(파일매니저), Linear desktop, Pot(번역기), Athens(노트, EOL). **IDE 사례 0** |
| **Wails + React/Vue** | **거의 없음** — Riftshare(파일공유), 일반 데스크톱 앱. **IDE 사례 0** |
| **Rust 네이티브 (Electron 거부)** | **Zed** (GPUI), **Lapce** (Floem), Helix(TUI 전용), neovide |
| **Java/JVM** | JetBrains 전체 라인 (IntelliJ, GoLand, …) |
| **Web (Electron 거부)** | StackBlitz WebContainers, GitHub Codespaces, Gitpod |

> 잔인한 결론: **Electron 외 데스크톱 IDE 사례는 사실상 Rust 네이티브 둘(Zed/Lapce)과 JetBrains 뿐**. Wails/Tauri로 IDE를 만든 OSS는 **0개**.

### 3.2 CodeMirror vs Monaco — 우리가 이미 한 베팅

VS Code/Cursor/Theia가 모두 **Monaco**를 쓰는데 우리는 **CodeMirror 6**을 쓰고 있음. 이건 의식적 선택이든 우연이든 결과적으로 큰 차이를 만듦.

| 항목 | **Monaco (VS Code)** | **CodeMirror 6 (현재 우리)** |
|---|---|---|
| 메이커 | Microsoft | Marijn Haverbeke (ProseMirror 동일) |
| 라이선스 | MIT | MIT |
| 번들 크기 | **5–8 MB gzipped** | **0.5–1 MB gzipped** ✓ |
| 모바일 | 약함 | **강함** ✓ (ProseMirror 라인) |
| LSP 통합 | 풍부함 (vscode-languageserver) | 직접 만들어야 함 (codemirror-languageserver 등) |
| Tree-sitter | 간접 | **직접 통합 가능** (Lezer 또는 web-tree-sitter) |
| 다국어 IME (한국어) | 일부 이슈 알려짐 | **더 안정적** |
| 동시 편집 (CRDT) | 가능하지만 무거움 | **Yjs와 자연스럽게 결합** |
| 진영 | VS Code 표준 | Replit, CodePen, Sourcegraph, Notion-ish |

**시사점**:
- **번들 크기 10배 차이** = Wails+CodeMirror = 30–50 MB 데스크톱 앱 가능
- 만약 Monaco였다면 IDE 호환성은 좋지만 hanimo의 "simplicity-first + 가벼움" 정책과 충돌
- CodeMirror 6은 **모바일 IDE 비전**(메모리: hanimo Vision 5축)과 정확히 정렬 — Monaco는 모바일 거의 불가능
- **이 선택이 우리의 차별화 무기 중 하나** — VS Code 포크 진영(Cursor/Void)이 절대 못 따라오는 영역

### 3.3 항목별 비교

| 항목 | **Wails + Go + React (현재)** | **Electron + React + Monaco** | **Tauri + React + Monaco/CM** |
|---|---|---|---|
| 바이너리 크기 | 30–60 MB (시스템 웹뷰) | **150–300 MB** (Chromium 동봉) | **8–20 MB** (시스템 웹뷰) |
| 메모리 사용 | 100–250 MB | 400 MB – 1.5 GB+ | 100–250 MB |
| 백엔드 언어 | **Go** ✓ (TUI와 코드 공유 가능) | Node.js | **Rust** (학습 비용 큼) |
| 백엔드 ↔ 프론트 IPC | Wails bindings (Go 함수 직접 호출) | IPC + preload script | Tauri commands (Rust→TS) |
| 시스템 웹뷰 | macOS=WKWebView, Win=Edge WebView2, Linux=WebKitGTK | Chromium 동봉 (모든 OS 동일) | macOS=WKWebView, Win=Edge WebView2, Linux=WebKitGTK |
| 웹뷰 일관성 | **3개 엔진 차이 존재** (Linux WebKitGTK가 가장 약점) | 1개 엔진 (예측 가능) | Wails와 동일 약점 |
| Monaco 호환성 | 가능, 그러나 WebKitGTK에서 일부 이슈 | 완벽 | Wails와 동일 |
| CodeMirror 호환성 | **양호** ✓ (현재 채택) | 양호 | 양호 |
| LSP/PTY/Git 풀 | Go에 강함 (lazygit, gopls 라인 재사용) | **node-pty 등 npm 풀 압도적** | Rust 풀 작음, 일부 Node 통합 가능 |
| MCP SDK | Go MCP 2nd-class | **Node MCP 1st-class** | Rust MCP 가능, 풀 작음 |
| 보안 | Go 단일 바이너리 | RCE 사례 다수 (Atom, Discord 등) | **샌드박스 가장 강함** |
| 자동 업데이트 | Wails 자체 기능 빈약 | 표준 (electron-updater) | tauri-updater (양호) |
| 컨트리뷰터 풀 | **작음** (Wails+Go IDE 경험자 거의 없음) | **압도적** | 중간 |
| OSS IDE 레퍼런스 | **0개** | **수십개 + Cursor가 검증** | **0개** |
| 트렌드 | 안정적이나 정체 | 표준이지만 비대화 비판 다수 | 🔥 떡상 중이나 IDE 분야 미증명 |

### 3.4 시나리오별 시사점

#### Scenario A: Wails 유지 (현재)
- ✅ Go 코드 공유 = 두 surface 단일 코드베이스
- ✅ 30–60 MB = simplicity-first 정책 정확히 정렬
- ✅ CodeMirror 6 + Lucide = 이미 모바일 IDE 비전과 정렬
- ❌ IDE OSS 레퍼런스 0 — 길을 직접 뚫어야 함
- ❌ 컨트리뷰터 합류 속도 가장 느림
- ❌ Linux WebKitGTK 호환성 추가 검증 필요
- **명성 가설**: "30 MB IDE를 단일 Go 바이너리로 만들었다" 라는 헤드라인이 먹히면 GitHub trending. 안 먹히면 그냥 unknown

#### Scenario B: Electron으로 전환
- ✅ Cursor/Void/Theia 모든 패턴 그대로 흡수 (LSP, terminal, file watcher, MCP)
- ✅ 컨트리뷰터 풀 5–10배
- ✅ CodeMirror도 Monaco도 자유 선택
- ❌ 200 MB 바이너리 = simplicity-first 정책 위반
- ❌ "또 다른 Cursor 포크" 인식 위험
- ❌ Go 백엔드 코드 → Node로 전부 재작성 = TUI/Desktop 코드 공유 끊김

#### Scenario C: Tauri로 전환
- ✅ 8–20 MB = "10 MB IDE" 헤드라인
- ✅ 보안 모델 가장 강함
- ❌ Rust 학습 비용 + 한국어 컨트리뷰터 풀 작음
- ❌ IDE 분야 미증명 (Wails와 동일 위험)
- ❌ Go 백엔드 코드 → Rust 재작성

### 3.5 결론 (Desktop)
**현재 stack 유지가 일관성 있음**. 근거:
1. simplicity-first 정책 + 모바일 비전과 stack이 정확히 정렬
2. CodeMirror 6 채택은 이미 우연이든 의식적이든 **올바른 베팅** (모바일 + 한국어 + Yjs)
3. Go 백엔드 공유 = TUI와 code reuse라는 큰 자산
4. "30 MB IDE" 가설을 한 번은 시도해볼 가치 있음

**그러나 단서**:
- 6개월 안에 Wails Linux 호환성 stress test 필요
- IDE OSS 레퍼런스 0 = 모든 패턴(LSP wire, terminal PTY, MCP transport, settings sync)을 직접 설계해야 함
- 만약 6개월 시도 후 컨트리뷰터 합류가 부진하면 **Electron 옵션 평가** (Tauri는 IDE 미증명이라 risky)

---

## 4. OSS 명성 관점 종합

### 4.1 명성 = 진입 마찰 ÷ 차별화

| Surface | 진입 마찰 (사용자) | 진입 마찰 (컨트리뷰터) | 차별화 |
|---|---|---|---|
| **hanimo-code TUI (Go)** | 매우 낮음 (단일 바이너리) | 중간 (Go 풀 작음) | 높음 (한국어 + Korea MCP + Charm v2) |
| **hanimo-code-desktop (Wails)** | 중간 (시스템 웹뷰 호환성) | **높음** (Wails IDE 사례 0) | 높음 (30 MB IDE + CodeMirror 모바일 라인) |
| 가상의 Electron 전환 | 중간 (200 MB 다운로드) | 매우 낮음 | 낮음 (또 하나의 Cursor 포크) |
| 가상의 Ink 전환 | 중간 (Node 필요) | 매우 낮음 | 낮음 (또 하나의 claude-code 클론) |

### 4.2 직관 검증
- **TUI 진영**: 우리는 "Charm 본진 + 한국어"로 이미 unique 영역. 바꿀 이유 없음.
- **Desktop 진영**: 우리는 "Wails + CodeMirror + 모바일 라인"으로 unique 영역이지만 검증 필요. 6개월 시도 후 재평가가 합리적.

---

## 5. 권고 결정 (Phase 4 이후 가이드)

| 결정 | 권고 | 근거 |
|---|---|---|
| TUI stack 변경? | **No** | charmbracelet/crush와 같은 라인. 단일 바이너리 22 MB는 명성 자산. AI SDK 격차는 직접 패치로 대응 가능 |
| Desktop stack 변경? | **No (단, 6개월 검증 게이트 설정)** | CodeMirror 6 + Wails는 simplicity + 모바일 비전과 정렬. Electron/Tauri 전환 비용이 차별화 손실보다 큼 |
| Monaco로 전환? | **No** | CodeMirror 6의 가벼움 + 모바일 + 한국어 IME가 우리 비전과 정확히 정렬 |
| Linux WebKitGTK 검증 | **반드시** | Wails 약점. 6개월 안에 stress test |
| 6개월 후 재평가 트리거 | 컨트리뷰터 < 10명 + GitHub stars 정체 | 그때 Electron 옵션 검토 |

---

## 6. 부록 — 추가 OSS 레퍼런스 데이터

### 6.1 AI 코딩 에이전트 stack 매핑 (2026 기준)

| 도구 | Stack | 비고 |
|---|---|---|
| **Anthropic claude-code** | TS + Ink + Node | Anthropic 공식 CLI |
| **Google gemini-cli** | TS + Ink + Node | Apache 2.0 OSS |
| **Charmbracelet crush** | Go + Bubble Tea | **우리와 동일 stack** |
| **Charmbracelet mods** | Go + Bubble Tea | LLM pipe 도구 |
| **opencode** | TS + Ink | OSS Claude Code 클론 |
| **aider** | Python + Rich | 가장 유명한 AI 코딩 CLI |
| **Continue** | TypeScript (VS Code 확장 + JetBrains 플러그인) | IDE 통합 라인 |
| **Cursor** | Electron + React + Monaco (VS Code 포크) | 클로즈드 |
| **Void** | Electron + React + Monaco | OSS Cursor |
| **Zed** | Rust + GPUI | Electron 거부 |
| **hanimo-code** | Go + Bubble Tea v2 | Charm 라인 |

### 6.2 데스크톱 IDE 진영 stack 매핑

| 도구 | Stack |
|---|---|
| **VS Code** | Electron + TS + Monaco |
| **Cursor** | Electron + React + Monaco |
| **Theia** | Electron/Browser + TS + Monaco |
| **Void** | Electron + React + Monaco |
| **Trae** | Electron + React + Monaco |
| **Zed** | Rust 네이티브 (GPUI) |
| **Lapce** | Rust 네이티브 (Floem) |
| **JetBrains** | JVM (IntelliJ Platform) |
| **hanimo-code-desktop** | **Wails (Go) + React + CodeMirror 6 + xterm.js** |

### 6.3 우리 stack의 unique 좌표

```
                  무거움 (Electron)
                       │
                       │  VS Code, Cursor, Void
                       │  Theia, Trae
                       │
   (Web)               │           (Native)
   StackBlitz ─────────┼─────────  Zed, Lapce
                       │           JetBrains
                       │
        Tauri          │   Wails (우리)
        Spacedrive     │   ★ hanimo-code-desktop
                       │
                  가벼움 (시스템 웹뷰)
```

우리는 "**시스템 웹뷰 + Go 백엔드 + CodeMirror**"라는 **빈 사분면**에 있음. IDE 진영에서 OSS 레퍼런스가 없는 이유는 어렵기 때문이지만, 성공하면 **유일한 이름**.

---

## 7. 이전 비교 요약 (2026-05-03 1차 답변 보존)

이 절은 첫 번째 비교 답변 그대로 보존.

### 7.1 hanimo-code (TUI)

| 항목 | Go + Bubble Tea | React + Ink |
|---|---|---|
| 배포 마찰 | 단일 바이너리 (brew/curl, 5–15 MB) | Node 런타임 필요, 또는 50–100 MB 패킹 |
| 메모리/성능 | 매우 가벼움 (20–50 MB) | 무거움 (100–300 MB) |
| 동시성 | goroutines 자연스러움 | Node 이벤트 루프 |
| AI/MCP SDK | 2nd-class | **1st-class** |
| TUI 생태계 | Charm 폭발적 | Ink 안정적 |
| 컨트리뷰터 풀 | Go (작음) | **JS/TS (압도적)** |
| OSS 레퍼런스 | gh, lazygit, k9s, **crush** | claude-code, gemini-cli |
| 트렌드 | 🔥 Charm 떡상 | 안정 단계 |

### 7.2 hanimo-code-desktop (IDE)

| 항목 | Go + Wails | Electron + React | Tauri + React |
|---|---|---|---|
| 바이너리 크기 | 20–50 MB | **100–250 MB** | **5–15 MB** |
| 메모리 | 100–200 MB | 300 MB – 1 GB+ | 100–250 MB |
| 백엔드 언어 | **Go** | Node | **Rust** |
| Monaco/CodeMirror 통합 | 가능, 사례 적음 | **압도적 표준** | 가능, 사례 적음 |
| LSP/터미널/PTY | Go 자체 강함 | npm 풀 압도적 | Rust 풀 작음 |
| 보안 모델 | Go 단일 바이너리 | RCE 사례 많음 | **샌드박스 강함** |
| 컨트리뷰터 풀 | 작음 | **압도적** | 중간 |
| IDE OSS 레퍼런스 | 거의 없음 | VS Code, Cursor, Atom, Theia, Void | 거의 없음 |

### 7.3 1차 답변의 결론
- 터미널 = AI 에이전트 진영은 React+Ink 쏠림, 인프라 CLI는 Go+Charm 압승
- 데스크톱 = Electron 쏠림 절대적, Wails/Tauri는 IDE 미증명
- 명성 목적 + simplicity-first 라면 Go 라인 유지 합리적

---

## 8. 변경 이력
- **2026-05-03**: 1차 작성. 실측 stack + OSS 진영 매핑 + 결정 권고 + 1차 답변 보존.
