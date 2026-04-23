# 양방향 포팅 분석 — 2026-04-23

> **작성일**: 2026-04-23
> **작성자**: 김지원 + Claude Opus 4.7
> **목적**: hanimo-code ↔ TECHAI_CODE 본질(외부망/폐쇄망) 필터링 후 실행 가능한 포팅 리스트 확정
> **상위 문서**: `docs/porting/sync-tracker-2026-04.md` · `docs/porting/hanimo-to-techaicode-plan.md` · `TECHAI_CODE/docs/PORTING_LOG.md`

---

## 0. TL;DR

1. **이번 주 작업 A (hanimo → TECHAI)**: LSP + 룰베이스 도구 5종 + MCP stdio + 테스트 3종. 반나절.
2. **이번 주 작업 B (TECHAI → hanimo)**: memory.go + 테스트 3종 + toolparse 엣지. 1시간.
3. **다음 주**: `hanimo-code-desktop`(IDE) 포팅 (`docs/porting/IDE_PORTING_PLAN.md`).
4. **그 이후 공통 신규**: Repo-map · Subagent · apply_patch AST · JSON headless · AGENTS.md 로더. 양쪽 동시.
5. **각자 전용 영역은 절대 포팅 금지** (아래 4장 참조).

---

## 1. 본질 차이 (필터 기준)

| | **hanimo-code (외부망)** | **TECHAI_CODE (폐쇄망)** |
|---|---|---|
| 네트워크 | 공용 인터넷 OK | 단일 사내 게이트웨이만 |
| 청중 | 전 세계 OSS 개발자 | 신한은행 사내 엔지니어 |
| 프로바이더 | 다중 합법 (Anthropic/OpenAI/Google/OpenRouter/Ollama) | 사내 OpenAI-compat 프록시 1개 |
| 언어 | ko + en (i18n) | ko 고정 |
| 프레임워크 지식 | 범용 (React/Spring/...) | **+ BXM·JEUS·Tibero** (사내 특화) |
| 배포 | 공개 GitHub 바이너리 | USB / 사내 FTP |
| 감사 | 선택 | **필수** (로그/크레덴셜 스크럽, audit trail) |
| 수익 모델 | **없음** (명성 목적 OSS) | 사내 업무 도구 |

---

## 2. hanimo → TECHAI 방향 (본질 필터 통과)

실제 소스 diff 결과 TECHAI에 아직 없는 항목 중 폐쇄망 본질에 맞는 것.

| # | 항목 | 경로 | 판정 근거 |
|:-:|---|---|---|
| A1 | **LSP 클라이언트** | `internal/lsp/` (505 LOC) | gopls/tsserver/pyright 로컬 바이너리, 네트워크 無. 엔터프라이즈 대형 코드베이스에 더 필요 |
| A2 | **Import Graph** | `tools/imports.go` | 파일시스템만, 네트워크 無 |
| A3 | **Test Coverage Gap** | `tools/coverage.go` | 로컬 |
| A4 | **Code Quality Scan** | `tools/quality.go` | 로컬 룰베이스 · 사내 품질 게이트 활용 가능 |
| A5 | **5단계 Edit Replacer** | `tools/replacer.go` | 파일 편집 강건성, 필수 |
| A6 | **Smart Context** | `tools/smartctx.go` | 로컬 |
| A7 | **MCP stdio** | `internal/mcp/` | **폐쇄망에서 오히려 가치 큼** — 사내 Jira/Wiki/Confluence 통합 이미 README에 예시 있음 |
| A8 | **테스트 3종** | `llm/prompt_embed_test.go` · `llm/context_self_test.go` · `skills/loader_test.go` | 회귀 방어 |

**실행**: 복사 + import 경로 치환(`hanimo` → `tgc`) + 기능 검증. 예상 반나절.

## 3. TECHAI → hanimo 방향 (본질 필터 통과)

| # | 항목 | 경로 | 판정 근거 |
|:-:|---|---|---|
| B1 | **memory.go** (JSON /remember 스토어) | `tools/memory.go` | 범용 메모리 — hanimo엔 `session/memory.go`만 있고 도구 형태 없음 |
| B2 | **hooks_test.go** | `hooks/hooks_test.go` | 회귀 방어 |
| B3 | **llm capabilities/client 테스트** | `llm/capabilities_test.go` · `llm/client_test.go` | 회귀 방어 |
| B4 | **toolparse 엣지케이스 5건** | `toolparse.go` 최근 커밋 (`55e444a`·`f62a1bd` 등) | Qwen3 프록시 호환성 — hanimo도 Ollama Qwen3 사용 |

**실행**: 마찬가지 복사 + import 경로 역치환. 예상 1시간.

---

## 4. 포팅 절대 금지 (본질 위배)

### 4.1 hanimo → TECHAI 금지

| 항목 | 이유 |
|---|---|
| Ollama / Anthropic / Google provider | **폐쇄망 도달 불가** |
| Provider interface 추상화 | 단일 게이트웨이가 본질 |
| Prompt caching (`llm/cache.go`) | 폐쇄망 Novita/사내 프록시 미지원 |
| i18n 다국어 (`ui/i18n.go`) | 한국어 고정 |
| WebFetch / 웹 검색 (존재할 경우) | 외부 HTTP 금지 |
| `cmd/spark/` (DGX Spark 클라이언트) | 외부 제품 |

### 4.2 TECHAI → hanimo 금지

| 항목 | 이유 |
|---|---|
| bxm 지식팩 13개 | **신한 전용 사내 프레임워크** — 오픈소스 법무·브랜딩 리스크 |
| `cmd/scrape-bxm` | 동상 |
| 한국어 고정 · 이모지 제거 · 영어 제거 | 오픈소스는 i18n 유지 |
| 단일 게이트웨이 하드닝 | 의도적 설계 차이 |
| Audit-focused 하드닝 (폐쇄망 전용) | 오픈소스엔 과투자 |

---

## 5. 양쪽 공통 신규 (본질 필터 통과 후 양쪽 동시)

| 항목 | 근거 |
|---|---|
| Repo-map (tree-sitter + PageRank) | 로컬 파싱만, 양쪽 다 필요 |
| Subagent (context fork) | 순수 아키텍처 |
| apply_patch AST 검증 | 양쪽 apply_patch 있음, 검증 미비 |
| Doom loop detector 고도화 | 양쪽 기본만 있음 |
| Session Browser UI + `--resume <id>` | UX |
| JSON headless (`--output json`) | CI/CD |
| DECSET 2026 synchronized output | 터미널 렌더링 |
| AGENTS.md 표준 로드 | 표준 |
| Sandbox (seatbelt/unshare) | 안전망 |

---

## 6. 실행 순서 (확정)

```
Week 17 (이번 주)
├── 월-화: 작업 A (hanimo → TECHAI, 반나절 × 2회 분할)
├── 수: 작업 B (TECHAI → hanimo, 1시간)
└── 목-금: Sprint 1 착수 (Tier 1 완벽주의, REPUTATION-STRATEGY 참조)

Week 18 (다음 주)
└── hanimo-code-desktop (IDE) 포팅 (docs/porting/IDE_PORTING_PLAN.md)
    — 레포명 확정: flykimjiwon/hanimo-code-desktop (obsidian 리브랜드 문서 기준)

Week 19+
└── 공통 신규 (Repo-map → Subagent → Permission 5-mode → ...)
```

---

## 7. 관련 문서

- `docs/porting/sync-tracker-2026-04.md` — 전체 동기화 트래커 (참조용)
- `docs/porting/hanimo-to-techaicode-plan.md` — 상세 실행 가이드 (fresh agent용)
- `docs/porting/IDE_PORTING_PLAN.md` — 다음 주 IDE 포팅
- `docs/SESSION-2026-04-16-17-PORTING-SUMMARY.md` — 지난 포팅 세션 요약
- `docs/strategy/REPUTATION-STRATEGY-2026-04-23.md` — 경쟁 전략 상위 문서
- `TECHAI_CODE/docs/PORTING_LOG.md` — 역방향 포팅 이력 (Tier 1/2 완료)
