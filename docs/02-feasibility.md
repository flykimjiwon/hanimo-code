# 커스텀 터미널 멀티에이전트 시스템 구축 타당성 분석

> 아키텍처 옵션, LLM 통합, 필수 컴포넌트, 오픈소스 활용, 현실적 타임라인

---

## A. 아키텍처 옵션

### 언어/런타임 비교 (가중 점수)

| 언어 | 가중 점수 | TUI | LLM 에코시스템 | 개발 속도 | 성능 |
|------|----------|-----|---------------|----------|------|
| **TypeScript (Node.js)** | **8.20/10** | Ink (React-in-terminal) | Vercel AI SDK 25+ | 최고 | 보통 |
| **Python** | 7.70/10 | Textual, Rich | LiteLLM 100+ | 빠름 | 보통 |
| **Go** | 7.60/10 | Bubbletea (Elm 아키텍처) | 부족 (프록시 필요) | 보통 | 좋음 |
| **Rust** | 6.80/10 | Ratatui (30~40% 적은 메모리) | 매우 부족 | 느림 | 최고 |

> TypeScript가 LLM 에코시스템 기준 +0.60 포인트 리드. Go는 TUI DX 1위 (OpenCode 검증).

### TUI 프레임워크

- **Ink** (TS) — React 컴포넌트 모델. 프로덕션 CLI 다수 사용. 고빈도 스트리밍에서 약간 래그.
- **Bubbletea** (Go) — Elm Architecture. OpenCode, Charm's Crush가 사용. 최고 TUI DX.
- **Textual** (Python) — OOP 위젯, CSS-like 레이아웃. Python 팀에 최적.
- **Ratatui** (Rust) — Immediate-mode 렌더링. 15% 낮은 CPU, 30~40% 적은 RAM. 구현 복잡.

---

## B. LLM 프로바이더 통합

### 통합 추상화 라이브러리

| 라이브러리 | 언어 | 프로바이더 수 | 추천 대상 |
|-----------|------|-------------|----------|
| **LiteLLM** | Python | 100+ | Python 빌드, `completion()` 한 줄 |
| **Vercel AI SDK** | TypeScript | 25+ | TS 빌드, 스트리밍 우선 |
| **LangChain(.js)** | Both | 다수 | 복잡한 체인/RAG |
| **Mastra** | TypeScript | 다수 | 올인원 TS 에이전트 프레임워크 |

### OpenAI 호환 프로바이더 (코드 변경 없이 연결)

```
Ollama (로컬)     → http://localhost:11434/v1     | Auth: 없음
GLM/Zhipu         → https://open.bigmodel.cn/...  | Auth: Bearer JWT
Azure OpenAI      → https://<resource>.openai...   | Auth: api-key 헤더
OpenRouter        → https://openrouter.ai/api/v1   | Auth: Bearer
vLLM/llama.cpp    → 커스텀 엔드포인트              | Auth: Bearer
```

### 별도 SDK 필요한 프로바이더

| 프로바이더 | SDK | Auth 방식 |
|-----------|-----|----------|
| **Claude** | `@anthropic-ai/sdk` | `x-api-key` + `anthropic-version` 헤더 |
| **Gemini** | `@google/generative-ai` | URL에 `?key=<KEY>` 또는 OAuth2 |

> LiteLLM 사용 시 Claude, Gemini도 별도 코드 불필요.

---

## C. 필수 컴포넌트 8가지

### 1. TUI 렌더링
- 스트리밍 토큰 실시간 렌더링
- 분할 패널: 채팅 + 파일 diff + 에이전트 상태
- 멀티라인 에디터 (vim 키바인딩)
- 상태바: 활성 에이전트, 모델, 토큰/비용 카운터

### 2. 파일 시스템 + 코드 편집
- `.gitignore` 존중하는 읽기/쓰기
- Unified diff 적용 + undo 스택
- 파일 와처 (외부 변경 감지)

### 3. Git 통합
- `status`, `diff`, `add`, `commit`, `log` 서브프로세스 호출
- 선택: libgit2 바인딩 (nodegit, pygit2, go-git)
- aider의 4가지 diff 모드 (whole, diff, udiff, architect) 참고

### 4. LSP 통합 (가장 큰 효과)
- **900x 속도 향상**: LSP 50ms vs grep 45초
- 핵심 작업: `goToDefinition`, `findReferences`, `hover`, `documentSymbol`, `getDiagnostics`
- 라이브러리: `vscode-languageclient` (TS), `multilspy` (Python), 표준 stdio JSON-RPC (Go)

### 5. Tool/Function Calling
- OpenAI `tools[]` JSON Schema가 사실상 표준
- Anthropic: `tools[]` + `input_schema` / Gemini: `functionDeclarations`
- LiteLLM이 모두 정규화

### 6. 에이전트 간 통신
- **MVP**: 내부 async 이벤트 큐 / 메시지 버스
- **풀**: A2A 프로토콜 (Google, Linux Foundation) — HTTP JSON-RPC, P2P

### 7. MCP 지원
- SDK: `@modelcontextprotocol/sdk` (TS), `mcp` (Python)
- 2025.11 스펙: Tasks primitive, OAuth auth, 구조화된 도구 출력
- 내 도구를 MCP 서버로 노출 → Claude Desktop, Cursor, Zed에서 호출 가능

### 8. 샌드박싱/보안
- 셸 실행 전 사용자 승인 필수
- 파일시스템 화이트리스트/블랙리스트
- Docker 컨테이너 격리
- MCP 보안 주의: 공개 서버 중 인증 있는 서버 0% (Knostic, 2025.07)

---

## D. 기존 오픈소스 활용

### 포크/라이브러리 후보

| 프로젝트 | Stars | 언어 | 라이선스 | 활용 방법 |
|---------|-------|------|---------|----------|
| **OpenCode** | ~120K | Go | MIT | **포크** — TUI+LSP+세션 기반 |
| **aider** | ~40K | Python | Apache 2.0 | **포크** — Git+LiteLLM 기반 |
| **LiteLLM** | 50K+ | Python | MIT | **라이브러리** — pip install |
| **Bubbletea** | 30K+ | Go | MIT | **라이브러리** — go.mod |
| **Ink** | 25K+ | TS | MIT | **라이브러리** — npm |
| **MCP SDK** | - | TS/Python | MIT | **라이브러리** |

### 카테고리별 최적 라이브러리

```
LLM 추상화:  LiteLLM (Python) / Vercel AI SDK (TS)
TUI:         Bubbletea (Go) / Ink (TS)
LSP:         multilspy (Python) / vscode-languageclient (TS)
Git:         gitpython / go-git / subprocess
MCP:         @modelcontextprotocol/sdk (TS) / mcp (Python)
세션 저장:    SQLite (better-sqlite3 / go-sqlite)
설정:        TOML/YAML + 환경변수 오버라이드
```

---

## E. 현실적 평가

### 타임라인

| 범위 | 기간 | 산출물 |
|------|------|--------|
| **MVP** | ~7주 | 코어 루프 + TUI + 파일/Git + 멀티프로바이더 |
| **풀 v1** | ~15주 | + LSP + 멀티에이전트 + MCP + 보안 |
| **프로덕션** | 22+주 | + 테스트, 문서, 패키징, CI |

### 포크 시 단축

| 전략 | MVP 기간 | 절약 |
|------|---------|------|
| **OpenCode 포크** (Go) | **2~3주** | ~4주 절약 |
| **aider 포크** (Python) | 3~4주 | ~3주 절약 |
| 처음부터 구축 (TS) | ~7주 | 기준 |

### 가장 어려운 부분 TOP 5

1. **스트리밍 TUI 동기화** — 멀티 에이전트 토큰 스트리밍 + 키보드 + 상태바 동시 렌더링
2. **멀티에이전트 파일 충돌** — 3개 에이전트가 같은 파일 수정 시 락킹/이벤트 소싱
3. **LSP 라이프사이클** — 언어 서버 시작/중지, 비동기 진단, Windows `\` vs `/`
4. **크로스 플랫폼 터미널** — Windows ConPTY, ANSI 이스케이프, 마우스, tmux
5. **프로바이더 API 변경** — Anthropic, Google, OpenAI 도구 스키마 빈번한 변경

### Build vs Fork 판단

| 조건 | 추천 |
|------|------|
| 빠른 결과물이 중요 | **OpenCode 포크** |
| Python 선호 + LiteLLM 100+ | **aider 포크** |
| TS 네이티브 + 완전한 소유권 | **처음부터 구축** |
| 학습 자체가 목적 | **처음부터 구축** |

---

*분석 일시: 2026-03-26*
