# hanimo Desktop App — Design Spec

## Overview

hanimo(하니모)의 기존 JS/TS 코드를 그대로 활용하여, Tauri v2 래퍼로 크로스플랫폼 데스크톱 앱을 만든다.
Node.js 설치 없이 누구나 `.dmg` / `.exe` / `.AppImage` 하나만 다운로드하면 바로 사용 가능.

## Goals

- 기존 hanimo core 코드(~7,800줄 TS) 100% 재사용
- Windows, macOS, Linux 크로스플랫폼 지원
- 개발자 + 일반 사용자 모두 타겟
- 경량 배포 (30MB 이하 목표)
- 관리자/고객사별 커스터마이징 가능한 온보딩

## Architecture

```
┌─ Tauri v2 (래퍼) ──────────────────────────┐
│                                             │
│  ┌─ 웹뷰 (시스템 브라우저 엔진) ───────────┐ │
│  │  React + TypeScript UI                  │ │
│  │  - VSCode 스타일 레이아웃               │ │
│  │  - 채팅 패널, 에디터, 파일트리          │ │
│  │  - 다크/라이트 모드                     │ │
│  └─────────────────────────────────────────┘ │
│                                             │
│  ┌─ Node.js Sidecar ──────────────────────┐ │
│  │  기존 hanimo 코드 그대로                │ │
│  │  - core/agent-loop.ts                   │ │
│  │  - tools/*.ts (16개 도구)               │ │
│  │  - providers/*.ts (14개 프로바이더)     │ │
│  │  - roles/, session/, mcp/ 등           │ │
│  └─────────────────────────────────────────┘ │
│                                             │
│  ┌─ Rust (최소한) ─────────────────────────┐ │
│  │  - 윈도우 생성 & 관리                   │ │
│  │  - 시스템 트레이                        │ │
│  │  - 메뉴바                               │ │
│  │  - Sidecar 프로세스 관리                │ │
│  │  (보일러플레이트 수준, 커스텀 로직 없음) │ │
│  └─────────────────────────────────────────┘ │
└─────────────────────────────────────────────┘
```

### Communication Flow

```
React UI  ←──Tauri IPC (invoke)──→  Rust Shell  ←──stdio/HTTP──→  Node.js Sidecar
  │                                                                     │
  │  사용자 입력 (채팅 메시지)                                           │
  │  ──────────────────────────────────────────────────────────────────→ │
  │                                                                     │
  │  에이전트 응답 (스트리밍 토큰, 도구 실행 결과)                       │
  │  ←────────────────────────────────────────────────────────────────── │
```

1. React UI가 사용자 입력을 Tauri IPC로 전달
2. Rust 레이어가 Node.js sidecar 프로세스로 중계
3. Node.js sidecar가 기존 agent-loop 실행, 스트리밍 결과 반환
4. React UI가 실시간 렌더링

## Tech Stack

| Layer | Technology | Role |
|-------|-----------|------|
| Wrapper | Tauri v2 | 네이티브 윈도우, 시스템 통합, 빌드/배포 |
| Frontend | React 18 + TypeScript | UI 컴포넌트, 상태 관리 |
| Styling | TailwindCSS (or CSS Modules) | 다크/라이트 테마, 반응형 |
| Editor | Monaco Editor (@monaco-editor/react) | 코드 편집, 구문 강조, 탭 관리 |
| File Tree | Custom React component | 프로젝트 파일 탐색 |
| Terminal | xterm.js | 내장 터미널 에뮬레이터 (Phase 2) |
| Backend | Node.js sidecar (기존 hanimo) | 에이전트 루프, 도구, 프로바이더 |
| IPC | Tauri invoke API + Event system | 프론트↔백엔드 통신 |

## Target Users

- **개발자**: AI 코딩 에이전트를 GUI 환경에서 사용하고 싶은 사용자
- **일반 사용자**: 코딩 경험 없이 AI로 앱/웹사이트를 만들고 싶은 사용자
- **기업 관리자**: 사내 LLM 엔드포인트에 연결하여 팀에 배포하려는 관리자

## UI/UX Design

### Theme

- **다크 모드** (기본): 미니멀, 검정 배경, 단색 포인트 — 터미널/Claude Code 느낌
- **라이트 모드**: 모던 소프트, 밝은 배경, 둥근 모서리, 부드러운 그림자 — ChatGPT 느낌
- 시스템 설정 연동 (자동 전환) + 수동 토글

### Layout: VSCode Style

```
┌──────────────────────────────────────────────────────────┐
│  Title Bar (hanimo 🐶)                        ─ □ ✕     │
├────┬──────────┬─────────────────────┬────────────────────┤
│ 📁 │ EXPLORER │                     │ 🐶 HANIMO CHAT    │
│ 🔍 │          │  Editor Tabs        │                    │
│ 💬 │ src/     │  ┌─────┬─────────┐  │ User: 로그인 만들어│
│ ⚙️ │  app.tsx │  │code │ preview │  │                    │
│    │  login.  │  │     │  (tab)  │  │ hanimo: 생성      │
│    │  ...     │  │     │         │  │ 완료했습니다...     │
│    │          │  │     │         │  │                    │
│    │          │  │     │         │  │ [메시지 입력...]   │
│    │          │  └─────┴─────────┘  │                    │
├────┴──────────┴─────────────────────┴────────────────────┤
│  Terminal (토글)  $ npm run dev                          │
└──────────────────────────────────────────────────────────┘
```

**Components:**
- **아이콘 사이드바** (왼쪽 40px): 파일 탐색기, 검색, 채팅, 설정 아이콘
- **파일 트리** (왼쪽 200px, 리사이즈 가능): 프로젝트 디렉토리 구조 표시
- **에디터 영역** (중앙): Monaco Editor, 탭으로 파일 전환, 웹 미리보기도 탭으로 표시
- **채팅 패널** (오른쪽 300px, 리사이즈 가능): AI 대화, 도구 실행 결과 인라인 표시
- **터미널** (하단, 토글): xterm.js 기반, 기본 숨김 상태

### Onboarding: Config-Driven Wizard

관리자가 `deploy-config.json`으로 온보딩 동작을 제어:

```jsonc
{
  "onboarding": {
    "mode": "wizard",           // "wizard" | "minimal" | "skip"
    "steps": ["provider", "apiKey", "model"],
    "defaults": {
      "provider": "openai",
      "endpoint": "https://internal-llm.company.com",
      "model": "gpt-4o"
    },
    "locked": ["provider", "endpoint"],
    "branding": {
      "appName": "우리회사 AI",
      "logo": "./assets/logo.png"
    }
  }
}
```

**Scenarios:**

| Scenario | Config | User Experience |
|----------|--------|----------------|
| 오픈소스 일반 배포 | `mode: "wizard"`, full steps | 프로바이더 → API 키 → 모델 (풀 위저드) |
| 사내 배포 (프로바이더 고정) | `locked: ["provider", "endpoint"]` | API 키만 입력 |
| 완전 사전 설정 | `mode: "skip"`, all defaults | 온보딩 없이 바로 시작 |
| 화이트라벨 | `branding` 설정 | 회사 로고 + 이름으로 표시 |

### Web Preview

- 에디터 탭으로 표시 (파일 탭과 동일 레벨)
- hanimo가 HTML/CSS/JS 생성 시 자동 프리뷰 탭 오픈
- 코드 탭과 프리뷰 탭을 좌우 분할 가능
- Tauri 웹뷰를 활용하여 별도 브라우저 불필요

### Terminal

- 하단 패널, 기본 숨김
- 토글로 열기/닫기 (Ctrl+` 또는 메뉴)
- hanimo의 shell_exec 실행 결과를 채팅 인라인으로도 표시
- 사용자가 직접 명령어 입력 가능

## Phase Plan

### Phase 1 — MVP

최소 동작하는 데스크톱 앱. 핵심 기능만 포함.

| Component | Description |
|-----------|-------------|
| Tauri 셸 | 윈도우 생성, Node.js sidecar 실행, IPC 통신 |
| 온보딩 위저드 | 프로바이더 선택 → API 키 입력 → 모델 선택 (기본 wizard 모드) |
| VSCode 레이아웃 | 사이드바 + 파일트리 + 에디터(Monaco) + 채팅 패널 |
| AI 채팅 | 스트리밍 응답, 도구 실행 결과 인라인 표시, 역할 전환 |
| 다크/라이트 모드 | 두 테마 + 수동 토글 + 시스템 연동 |
| 파일 CRUD | 파일트리에서 파일 열기/생성/삭제, 에디터에서 편집/저장 |
| 빌드 & 배포 | macOS(.dmg), Windows(.exe), Linux(.AppImage) 빌드 파이프라인 |

### Phase 2

| Component | Description |
|-----------|-------------|
| 내장 터미널 | xterm.js 기반 하단 터미널 패널 |
| 웹 미리보기 | 에디터 탭으로 HTML/CSS/JS 실시간 렌더링 |
| deploy-config.json | 온보딩 커스터마이징 (locked, defaults, branding) |
| 세션 관리 | 세션 저장/불러오기/검색 |

### Phase 3

| Component | Description |
|-----------|-------------|
| 자율 모드 (/auto) | 에이전트 자동 반복 실행 |
| 멀티 에이전트 | 작업 분해 → 병렬 실행 → 결과 합성 |
| 브라우저 비주얼 프리뷰 | 앱 내 웹뷰로 목업/다이어그램 실시간 표시 |
| MCP 관리 UI | MCP 서버 연결/해제 GUI |

## Build & Distribution

| Platform | Format | Signing |
|----------|--------|---------|
| macOS | `.dmg`, `.app` | Apple Developer ID (선택) |
| Windows | `.exe`, `.msi` | Code signing cert (선택) |
| Linux | `.AppImage`, `.deb` | N/A |

- GitHub Actions CI/CD로 자동 빌드
- GitHub Releases로 배포
- 자동 업데이트: Tauri updater plugin

## File Structure (예상)

```
modol/
├── src/                        # 기존 hanimo core (변경 없음)
│   ├── core/
│   ├── tools/
│   ├── providers/
│   └── ...
├── desktop/                    # 새로 추가 — Tauri 앱
│   ├── src-tauri/              # Rust (Tauri 설정)
│   │   ├── src/
│   │   │   └── main.rs         # 윈도우 생성, sidecar 관리
│   │   ├── Cargo.toml
│   │   └── tauri.conf.json     # Tauri 빌드 설정
│   ├── src/                    # React UI
│   │   ├── App.tsx
│   │   ├── components/
│   │   │   ├── Sidebar.tsx
│   │   │   ├── FileTree.tsx
│   │   │   ├── Editor.tsx
│   │   │   ├── ChatPanel.tsx
│   │   │   ├── Terminal.tsx
│   │   │   └── Onboarding.tsx
│   │   ├── hooks/
│   │   ├── stores/
│   │   └── themes/
│   ├── package.json
│   └── vite.config.ts
├── deploy-config.json          # 배포 설정 (Phase 2)
└── package.json                # 루트
```

## Constraints & Decisions

- Rust 코드는 Tauri 보일러플레이트 수준만 작성. 비즈니스 로직은 모두 Node.js sidecar에서 처리.
- 기존 `src/` 디렉토리는 변경하지 않는다. 데스크톱 앱은 `desktop/` 하위에 격리.
- Monaco Editor를 사용하여 코드 편집 기능 구현. CodeMirror 대비 VSCode 호환성이 높음.
- 프론트엔드 상태 관리는 Zustand (경량, React 친화적).
- IPC는 Tauri의 invoke + event system 사용. WebSocket이나 HTTP 서버 불필요.
