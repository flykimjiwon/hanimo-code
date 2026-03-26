# 터미널 AI 코딩 도구 전수 조사 (2026.03)

> 35+ 도구를 5개 카테고리로 분류. GitHub Stars, 지원 LLM, 멀티에이전트 여부, MCP 지원 등 비교.

---

## 1. 공식 CLI 도구 (빅테크)

### Claude Code (Anthropic)

| 항목 | 내용 |
|------|------|
| GitHub | [anthropics/claude-code](https://github.com/anthropics/claude-code) |
| 언어 | TypeScript / Node.js |
| 플랫폼 | macOS, Linux, Windows |
| LLM | Claude 전용 |
| Stars | ~82,000 |
| 라이선스 | **Proprietary** (source-available) |
| 멀티에이전트 | O (워크트리 병렬) |
| MCP | O |

- 코드베이스 전체 이해 + 자연어 코딩
- VS Code 확장, GitHub @claude 태깅
- 플러그인 시스템, 커스텀 슬래시 커맨드
- GA v1.0.0 (2025.05), 2026.03 현재 활발히 개발 중

---

### OpenAI Codex CLI

| 항목 | 내용 |
|------|------|
| GitHub | [openai/codex](https://github.com/openai/codex) |
| 언어 | Rust |
| 플랫폼 | macOS, Linux, Windows |
| LLM | OpenAI (GPT-4o, o-series) |
| Stars | ~67,000 |
| 라이선스 | Apache 2.0 |
| 멀티에이전트 | O (워크트리 병렬) |
| MCP | O |

- Rust 기반 고성능 TUI
- 슬래시 커맨드, 이미지 첨부, 엔터프라이즈 기능
- v0.116.0 (2026.03.19): ChatGPT 디바이스 코드 로그인

---

### Google Gemini CLI

| 항목 | 내용 |
|------|------|
| GitHub | [google-gemini/gemini-cli](https://github.com/google-gemini/gemini-cli) |
| 언어 | TypeScript / Node.js |
| 플랫폼 | macOS, Linux, Windows |
| LLM | Gemini 3 (1M 토큰 컨텍스트) |
| 라이선스 | Apache 2.0 |
| 멀티에이전트 | 신생 (GitHub Actions 베타) |
| MCP | O |

- Google Search 그라운딩, 파일/셸/웹 도구
- GitHub Actions 연동 (2026.03.19 베타)

---

### GitHub Copilot CLI

| 항목 | 내용 |
|------|------|
| GitHub | [github/copilot-cli](https://github.com/github/copilot-cli) |
| 언어 | TypeScript / Node.js |
| 플랫폼 | macOS, Linux, Windows |
| LLM | **멀티모델** — Claude Opus 4.6, Sonnet 4.6, GPT-5.3-Codex, Gemini 3 Pro |
| 라이선스 | Proprietary (구독 필요) |
| 멀티에이전트 | O (클라우드 위임, 전문 에이전트) |
| MCP | O |

- GA 2026.02.25
- Plan 모드 (사전 계획) + Autopilot 모드
- `&` 접두사로 백그라운드 클라우드 위임

---

## 2. 오픈소스 터미널 도구

| 도구 | Stars | 언어 | LLM 다중 | 멀티에이전트 | MCP | 라이선스 | 핵심 |
|------|-------|------|---------|-----------|-----|---------|------|
| **OpenCode** | ~120K | Go | 75+ | O | O | MIT | 최고 TUI, LSP, SQLite |
| **aider** | ~40K | Python | 100+ (LiteLLM) | X | 제한 | Apache 2.0 | Git-native 페어 프로그래밍 |
| **OpenHands** | ~69K | Python | 다수 | O | O | MIT | Docker 샌드박스, Devin 대안 |
| **Cline** | ~59K | TS | BYOK | O | O | Apache 2.0 | VS Code 확장 |
| **Continue** | ~31K | TS | 50+ | O | O | Apache 2.0 | VS Code + JetBrains |
| **Goose** | ~29K | Rust | MCP 기반 | O | O | Apache 2.0 | MCP-native, Block 제작 |
| **GPT-Engineer** | ~54K | Python | OpenAI | X | X | MIT | 전체 프로젝트 생성 |
| **SWE-agent** | ~15K | Python | 다수 | X | X | MIT | 학술 벤치마크 최고 |
| **Void** | ~15K | TS | 다수 | O | O | Apache 2.0 | 오픈소스 Cursor 대안 |
| **Roo Code** | ~15K | TS | OpenAI호환 | O | O | Apache 2.0 | Cline 포크, 상용 |
| **Sweep** | ~7.5K | Python | OpenAI | X | X | Elastic 2.0 | GitHub 이슈→PR 자동화 |
| **Mentat** | ~2.5K | Python | GPT-4 | X | X | MIT | 멀티파일 편집 |
| **auto-coder** | ~4K | Python | 다수 | X | X | Apache 2.0 | 인터랙티브 채팅+실행 |

### 주요 도구 상세

**OpenCode** — Go 기반, Bubbletea TUI, 75+ 모델, LSP 통합, SQLite 세션, Vim 에디터. 5M+ 월간 유저. 프라이버시 우선 (코드 외부 미전송).

**aider** — Paul Gauthier 제작. Git-native 자동 커밋, 음성→코드, 이미지/웹 컨텍스트. o3-pro로 84.9% 정확도 (polyglot 벤치마크).

**OpenHands** — 구 OpenDevin. Docker 샌드박스에서 코드 작성, 명령 실행, 웹 브라우징, API 호출. v1.5.0 (2026.03).

**Goose** — Block(Square) 제작. Rust 기반, MCP-native. CLI + 데스크탑 앱. Jira, GitHub, 시스템 명령 통합.

---

## 3. 멀티에이전트 프레임워크

| 프레임워크 | Stars | 언어 | 핵심 아이디어 | MCP |
|-----------|-------|------|------------|-----|
| **AutoGen** (→ MS Agent Framework) | ~55K | Python/.NET | 대화형 멀티에이전트 | O |
| **MetaGPT** | ~47K | Python | 소프트웨어 회사 시뮬레이션 | 부분 |
| **CrewAI** | ~45K | Python | 역할 기반 크루 오케스트레이션 | O |
| **LangGraph** | ~25K | Python/TS | 그래프 기반 상태 에이전트 | O |
| **ChatDev** | ~26K | Python | 가상 SW 회사 대화 | X |
| **Smolagents** (HF) | ~18K | Python | 초경량 코드 실행 | O |
| **PydanticAI** | ~9K | Python | 타입-세이프 에이전트 | O |
| **AgentScope** | ~6K | Python | 프로덕션급 MsgHub | O |
| **Magentic-One** | (AutoGen 내) | Python | 오케스트레이터+4 전문가 | X |

**AutoGen**: Microsoft가 Semantic Kernel과 합쳐 Microsoft Agent Framework로 통합 (2025.10). 유지보수 모드.

**MetaGPT**: PM→아키텍트→개발→QA 역할 시뮬레이션. ICLR 2025 oral (상위 1.8%).

**CrewAI**: 역할 기반 자율 에이전트 오케스트레이션. Crews(자율) + Flows(파이프라인). 100K+ 인증 개발자.

---

## 4. 오케스트레이션/확장 레이어

### MCP (Model Context Protocol)

- 2026년 **사실상 표준** — 1,000+ MCP 서버, 34,700+ npm 종속 프로젝트
- MCP Apps (2026.01 GA): 대화 내 인터랙티브 UI 컴포넌트
- 2026 로드맵: 트리거, 스트리밍 결과, 엔터프라이즈 확장 (SSO, 감사 로그)
- ChatGPT, Claude, VS Code, Goose 등 모든 주요 도구에서 지원

### Oh My Claude Code (OMC)

- Claude Code 위에 올리는 멀티에이전트 오케스트레이션 레이어
- 32개 전문 에이전트 (아키텍처, 리서치, 디자인, 테스트, 데이터 사이언스)
- Ultrapilot: 5개 동시 워커, 3~5배 빠른 실행
- 스킬 시스템, 딥 인터뷰, 상태 관리

---

## 5. 상용 IDE 도구

| 도구 | 제조사 | Stars | 핵심 | 가격 |
|------|--------|-------|------|------|
| **Cursor** | Anysphere | N/A | #1 AI IDE, 100만+ 유저 | $20/월 |
| **Windsurf** | Cognition AI | N/A | Cascade 에이전트, LogRocket #1 | Free+ |
| **Devin** | Cognition AI | N/A | 자율 AI 개발자 | $20~$500/월 |
| **Kiro** | Amazon/AWS | N/A | AWS 네이티브, IAM Autopilot | TBD |
| **Amp** | Sourcegraph | N/A | 멀티모델 라우팅, 서브에이전트 | 상용 |
| **Zed** | Zed Industries | ~55K | Rust 네이티브, ACP | 무료 |

---

## 6. 인기순 랭킹 (GitHub Stars)

| # | 도구 | Stars |
|---|------|-------|
| 1 | OpenCode | ~120,000 |
| 2 | Claude Code | ~82,000 |
| 3 | OpenHands | ~68,600 |
| 4 | Codex CLI | ~67,000 |
| 5 | Zed | ~55,000 |
| 6 | AutoGen | ~54,600 |
| 7 | GPT-Engineer | ~54,000 |
| 8 | Cline | ~59,400 |
| 9 | MetaGPT | ~47,000 |
| 10 | CrewAI | ~44,589 |

---

## 7. 2026년 주요 관찰

1. **MCP가 범용 표준이 됨** — 모든 주요 도구가 MCP 지원
2. **공식 CLI가 마인드셰어 지배** — Claude Code, Codex CLI, Gemini CLI, Copilot CLI 모두 GA
3. **OpenCode가 깜짝 1위** — 120K stars, 5M+ MAU, Go TUI + 프라이버시
4. **멀티에이전트가 표준** — 단일 에이전트 도구(aider, Mentat) 상대적 하락
5. **IDE 시장 통합** — Windsurf→Cognition 인수 ($250M), Kiro는 Amazon Q 리브랜딩
6. **AutoGen 유지보수 모드** — Microsoft Agent Framework로 통합

---

*조사 일시: 2026-03-26 | 출처: GitHub, 공식 문서, 벤치마크 사이트*
