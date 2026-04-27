# hanimo ↔ TECHAI_CODE Sync Tracker — 2026-04

> **목적**: 두 프로젝트 간 기능 포팅 상태를 한 눈에. 양쪽 동일 파일로 유지(쌍방 복사).
> **업데이트**: 항목 이동 시 즉시. PR/커밋 단위로 tick.

---

## Legend

- ✅ 구현 완료
- 🚧 작업 중
- 📋 계획됨
- ❌ 미정
- ⚠️ 취약점/불일치

---

## 1. Core LLM & Tools

| 항목 | hanimo | TECHAI_CODE | 방향 | 비고 |
|---|:---:|:---:|:---:|---|
| OpenAI-compatible provider | ✅ | ✅ | — | 양쪽 기본 |
| Ollama provider | ✅ | 📋 | →TECHAI | T1-5 |
| Anthropic provider | ✅ | 📋 | →TECHAI | T1-5 |
| Google provider | ✅ | 📋 | →TECHAI | T1-5 |
| Provider interface 추상화 | ✅ | 📋 | →TECHAI | T1-5 |
| 3-stage compaction | ✅ | ✅ | ⇄ | 둘 다 유사 |
| Context window tracker | ✅ | ✅ | ⇄ | |
| Prompt caching (anthropic) | 📋 | ❌ | hanimo 먼저 | C2 |
| Tool: file_read/write/edit | ✅ | ✅ | ⇄ | |
| Tool: grep/glob | ✅ | ✅ | ⇄ | |
| Tool: shell_exec | ✅ | ✅ | ⇄ | dangerous regex 강화 필요 |
| Tool: hashline_read/edit | ✅ | 📋 | →TECHAI | T1-1 ⭐ |
| Tool: git_status/diff/log/commit | ✅ | 📋 | →TECHAI | T1-3 |
| Tool: diagnostics | ✅ | 📋 | →TECHAI | T1-4 |
| Tool: apply_patch (unified diff) | 📋 | 📋 | TECHAI 먼저 PoC | T2-4 |
| Read-before-write enforcement | 📋 | 📋 | 양쪽 동시 | IW |

---

## 2. Agent System

| 항목 | hanimo | TECHAI_CODE | 방향 | 비고 |
|---|:---:|:---:|:---:|---|
| Super / Dev / Plan 3-mode | ✅ | ✅ | ⇄ | |
| Plan 실제 read-only | ✅ | ⚠️ | →TECHAI | T1-2 🔥 |
| Agents 패키지 분리 | ✅ | 📋 | →TECHAI | T1-7 |
| `/auto` deep agent | ✅ | ⚠️ 부분 | →TECHAI | T1-7 |
| `intent.go` 의도 감지 | ✅ | ❌ | →TECHAI | T1-7 |
| `askuser.go` ASK_USER | ✅ | ⚠️ 미구조화 | →TECHAI | |
| Doom loop detector | 📋 | 📋 | 양쪽 동시 | IW |
| Iteration cap 환경변수 | ⚠️ 하드코딩 | ⚠️ 하드코딩 | 양쪽 동시 | |
| Subagent (context fork) | 📋 | 📋 | hanimo 먼저 | B1 |
| Git worktree 병렬화 | 📋 | ❌ | hanimo 먼저 | B2 |

---

## 3. Context Engineering

| 항목 | hanimo | TECHAI_CODE | 방향 | 비고 |
|---|:---:|:---:|:---:|---|
| `.hanimo.md` / `.techai.md` 로드 | ✅ | ⚠️ | →TECHAI | |
| clarifyFirstDirective | ✅ | 📋 | →TECHAI | IW-3 |
| 계층형 프롬프트 (embed 분리) | 📋 | 📋 | 양쪽 동시 | IW-2 |
| Repo-map (tree-sitter) | 📋 | 📋 | TECHAI PoC → hanimo | A1 / T2-3 |
| Knowledge store (SQLite) | ✅ | ✅ | ⇄ | |
| Skill 시스템 | 📋 | 📋 | TECHAI PoC → hanimo | A2 / T2-1 |
| Learned-skill 자동 추출 | 📋 | ❌ | hanimo 먼저 | A2 |
| AGENTS.md 표준 지원 | 📋 | 📋 | 양쪽 동시 | |

---

## 4. Extensibility

| 항목 | hanimo | TECHAI_CODE | 방향 | 비고 |
|---|:---:|:---:|:---:|---|
| MCP stdio transport | ✅ | 📋 | →TECHAI | T1-6 |
| MCP SSE transport | ✅ | ❌ | →TECHAI | |
| MCP HTTP streaming | 📋 | ❌ | hanimo 먼저 | |
| Hooks 시스템 | 📋 | 📋 | TECHAI PoC → hanimo | A3 / T2-2 |
| Permission rule 엔진 | 📋 | 📋 | hanimo 먼저 | A4 |
| Permission 5-mode 순환 | 📋 | 📋 | hanimo 먼저 | A4 |
| permissions.yaml 학습형 | 📋 | ❌ | hanimo 먼저 | A4 |

---

## 5. UX / TUI

| 항목 | hanimo | TECHAI_CODE | 방향 | 비고 |
|---|:---:|:---:|:---:|---|
| Bubble Tea 기반 | ✅ | ✅ | ⇄ | |
| Command palette (Ctrl+K) | ✅ | 📋 | →TECHAI | T3-1 |
| i18n (ko/en) | ✅ | 📋 | →TECHAI | T3-2 |
| Themes (5개) | ✅ | 📋 | →TECHAI | T3-3 |
| Session 저장 (SQLite) | ✅ | ✅ | ⇄ | |
| Session Browser UI | 📋 | ❌ | hanimo 먼저 | C1 |
| `--resume <id>` CLI | 📋 | 📋 | 양쪽 동시 | |
| JSON / headless 출력 | 📋 | 📋 | 양쪽 동시 | C5 |
| Synchronized Output (DECSET 2026) | 📋 | 📋 | 양쪽 동시 | C4 |
| LSP 통합 (gopls 등) | 📋 | ❌ | hanimo 먼저 | C3 |
| Tab / Shift+Tab 분리 | ⚠️ | ⚠️ | 양쪽 동시 | A4 |

---

## 6. Safety

| 항목 | hanimo | TECHAI_CODE | 방향 | 비고 |
|---|:---:|:---:|:---:|---|
| Dangerous regex 기본 | ✅ | ⚠️ 부분 | →TECHAI | T1-8 |
| 확장 dangerous regex (env leak 등) | 📋 | 📋 | 양쪽 동시 | IW-8 / T1-8 |
| Credential scrubbing (output) | 📋 | 📋 | hanimo 먼저 | |
| Sandbox (seatbelt/unshare) | ❌ | ❌ | hanimo 먼저 | C-future |
| Shell timeout env 가변 | ⚠️ 하드 | ⚠️ 하드 | 양쪽 동시 | IW-7 |

---

## 7. 포팅 규칙

1. **단방향 기본**: hanimo → TECHAI_CODE. 코드 베이스 스타일 일치를 위해 파일 이름/패키지명 동일 유지.
2. **PoC 역방향**: Skill, Repo-map, Hooks, apply_patch는 TECHAI_CODE에서 먼저 PoC → 검증 후 hanimo 포팅. 포팅 시 generic화 + 에러 핸들링 + 설정화 추가.
3. **양쪽 커밋 메시지 컨벤션**:
   - `port(hashline): from hanimo@<sha>` — 포팅 시
   - `poc(skill): initial loader` — PoC 시
   - `sync(docs): update sync-tracker` — 본 문서 업데이트
4. **테스트 동반**: 포팅 시 원본 프로젝트의 테스트 파일도 함께 복사·수정.
5. **문서 동기화**: 본 문서는 양쪽 `docs/porting/` 에 동일 내용으로 복사.

---

## 8. 다음 밀수 (Next Drop) — 우선 처리

### 8.1 블록 실행 (즉시)

| 블록 | 내용 | 예상 | 실행 가이드 |
|---|---|:---:|---|
| **A** | hanimo → TECHAI: LSP + 도구 5종 + MCP stdio + 테스트 3종 | 반나절 | `BLOCK-A-EXEC-GUIDE.md` |
| **B** | TECHAI → hanimo: memory.go + tests + toolparse + reasoning_content | 1~2h | `BLOCK-B-EXEC-GUIDE.md` |
| **C** | techai-ide → hanimo-code-desktop: Wails IDE 리브랜딩 | 1~2d | `BLOCK-C-EXEC-GUIDE.md` |

### 8.2 블록 이후 (Sprint 순서)

1. **Sprint 1**: Tier 1 완벽주의 (SWE-bench, 3-Layer Defense)
2. **Anthropic 네이티브 provider** (#5, 1일)
3. **Rate limit tracker** (#6, 4시간)
4. **Repo-map PoC** — TECHAI에서
5. **Subagent (context fork)** — hanimo 먼저
6. **Permission 5-mode** — hanimo 먼저

### 8.3 지속 동기화

2주 Sprint마다 TECHAI delta 스캔 → 선별 포팅.
상세: `CONTINUOUS-SYNC-STRATEGY.md`

---

## 9. 파일 대조표 (2026-04-24 최신)

### 9.1 TECHAI에만 있는 파일 → hanimo로 가져올 것 (블록 B)

| TECHAI 경로 | hanimo 목적지 | LOC | 상태 |
|---|---|:---:|---|
| `internal/tools/memory.go` | `internal/tools/memory.go` | 263 | 📋 블록 B |
| `internal/hooks/hooks_test.go` | `internal/hooks/hooks_test.go` | 151 | 📋 블록 B |
| `internal/llm/capabilities_test.go` | `internal/llm/capabilities_test.go` | 80 | 📋 블록 B |
| `internal/llm/client_test.go` | `internal/llm/client_test.go` | 498 | 📋 블록 B |
| `internal/tools/apply_patch_test.go` | `internal/tools/apply_patch_test.go` | 276 | 📋 블록 B |
| `internal/tools/grep_patterns_test.go` | `internal/tools/grep_patterns_test.go` | — | 📋 블록 B |
| `internal/tools/search_test.go` | `internal/tools/search_test.go` | — | 📋 블록 B |

### 9.2 hanimo에만 있는 파일 → TECHAI로 보낼 것 (블록 A)

| hanimo 경로 | TECHAI 목적지 | LOC | 상태 |
|---|---|:---:|---|
| `internal/lsp/client.go` | `internal/lsp/client.go` | ~200 | 📋 블록 A |
| `internal/lsp/protocol.go` | `internal/lsp/protocol.go` | ~150 | 📋 블록 A |
| `internal/lsp/servers.go` | `internal/lsp/servers.go` | ~155 | 📋 블록 A |
| `internal/tools/imports.go` | `internal/tools/imports.go` | ~150 | 📋 블록 A |
| `internal/tools/coverage.go` | `internal/tools/coverage.go` | ~120 | 📋 블록 A |
| `internal/tools/quality.go` | `internal/tools/quality.go` | ~130 | 📋 블록 A |
| `internal/tools/replacer.go` | `internal/tools/replacer.go` | ~160 | 📋 블록 A |
| `internal/tools/smartctx.go` | `internal/tools/smartctx.go` | ~140 | 📋 블록 A |
| `internal/mcp/transport_stdio.go` | `internal/mcp/transport_stdio.go` | — | 📋 블록 A |

### 9.3 양쪽 공통 (동기화 유지)

| 경로 | 동기화 상태 |
|---|---|
| `internal/agents/{auto,intent,plan,askuser}.go` | ⇄ 동기 |
| `internal/tools/{file,git,shell,search,hashline,...}.go` | ⇄ 동기 |
| `internal/llm/{client,compaction,context,prompt,...}.go` | ⇄ 동기 |
| `internal/ui/{chat,menu,palette,...}.go` | ⇄ 동기 |

### 9.4 포팅 금지 (본질 위배)

| 방향 | 항목 |
|---|---|
| TECHAI → hanimo 금지 | BXM 지식팩, `cmd/scrape-bxm/`, 한국어 고정, audit 강화, `demo-supersol/`, `바보맨/` |
| hanimo → TECHAI 금지 | Ollama/Anthropic/Google provider, i18n, SSE transport, Spark |

---

## 10. TECHAI 최근 진화 (2026-04-16~24) — 포팅 후보

| 커밋 | 기능 | 파일 | hanimo 포팅 | 우선순위 |
|---|---|---|:---:|---|
| `ef79db5`+`fbc744b`+`905199d` | reasoning_content 파싱 + 색상 분리 | `llm/client.go`, `ui/chat.go` | **Yes** | 높음 |
| `6044126` | 스트리밍 부분 tool_call 태그 방지 | `llm/client.go` | **Yes** | 높음 |
| `55e444a`+`5628e06` | toolparse 엣지케이스 7건 | `exec/exec.go` 또는 인라인 | **Yes** | 중간 |
| `b598651` | status bar reasoning/writing 표시 | `ui/tabbar.go` | **Yes** | 중간 |
| `9884959` | config 자동 마이그레이션 | `config/config.go` | 검토 | 낮음 |
| `0de24f5` | 한국어 조사 제거 fallback | 검색 로직 | 검토 | 낮음 |
| `d891bcf`+`93a6e01`+`faf6712` | IDE: Git branch UI + 자동저장 + 세션 | `techai-ide/` | **블록 C** | 별도 |

---

_Last updated: 2026-04-24_
