# TODO — 착수 전 결정 필요 (2026-04-23)

> Sprint 1 및 포팅 실제 착수 전 김지원이 확정해야 할 항목.
> 참조: `docs/strategy/REPUTATION-STRATEGY-2026-04-23.md` · `docs/porting/BIDIRECTIONAL-ANALYSIS-2026-04-23.md`

---

## 1. 전략 결정 대기 (5건)

| # | 질문 | 옵션 | 영향 |
|:-:|---|---|---|
| 1 | 주당 투입 시간 | ① 10h (취미) ② 20h ③ 30h ④ 풀타임 | Sprint 2주 주기 유지 가능 여부 |
| 2 | 런치 순서 | A 한국 먼저 (권장) → 3개월 후 영어권 · B 동시 · C 영어권 먼저 | 마케팅 노력 배분 |
| 3 | 슬로건 1문장 | A *"It asks. Others assume."* · B *"The agent that can't silently overwrite your edits."* · C 자체 제안 | 랜딩 · HN 타이틀 · README 상단 |
| 4 | 보류 문서 `status: deferred` frontmatter 추가 여부 | Yes / No | `DESKTOP-PLAN`·`vscode-extension-*`·`AGENTRANK-INTEGRATION`·`SPARK-NATIVE-APP`·`EMBEDDED-BROWSER` 5건 대상 |
| 5 | 브랜드 허브 레포 | A 신규 `flykimjiwon/hanimo` · B `hanimo-community` 겸임 · C `hanimo-code/landing-mockups/` 승격 | 랜딩 hosting 경로 |

---

## 2. 포팅 착수 대기 (3블록)

### 블록 A — hanimo → TECHAI_CODE (반나절)

| 파일/패키지 | 경로 (hanimo) | 목적지 (TECHAI) |
|---|---|---|
| LSP 클라이언트 | `internal/lsp/` (505 LOC) | `internal/lsp/` |
| Import Graph | `tools/imports.go` | `tools/imports.go` |
| Test Coverage | `tools/coverage.go` | `tools/coverage.go` |
| Code Quality | `tools/quality.go` | `tools/quality.go` |
| Edit Replacer | `tools/replacer.go` | `tools/replacer.go` |
| Smart Context | `tools/smartctx.go` | `tools/smartctx.go` |
| MCP stdio | `internal/mcp/` | `internal/mcp/` |
| 테스트 3종 | `llm/prompt_embed_test.go` · `llm/context_self_test.go` · `skills/loader_test.go` | 동일 경로 |

### 블록 B — TECHAI → hanimo (1시간)

| 파일 | 경로 (TECHAI) | 목적지 (hanimo) |
|---|---|---|
| memory.go | `tools/memory.go` | `tools/memory.go` |
| hooks_test | `hooks/hooks_test.go` | `hooks/hooks_test.go` |
| llm 테스트 2종 | `llm/capabilities_test.go` · `llm/client_test.go` | 동일 경로 |
| toolparse 엣지 | 최근 커밋 `55e444a`·`f62a1bd` 변경분 | hanimo `toolparse.go` |

### 블록 C — hanimo-code-desktop (IDE) 포팅 (2026-04-24~ 다음 주)

상세 플랜: `docs/porting/IDE_PORTING_PLAN.md`
레포명 확정: `flykimjiwon/hanimo-code-desktop`

---

## 3. Sprint 1 착수 대기 (Tier 1 완벽주의)

`docs/strategy/REPUTATION-STRATEGY-2026-04-23.md` §4 Sprint 1 참조.

| 태스크 | 상태 |
|---|:-:|
| `hanimo models` 명령 + Tier 시각화 표 | 대기 |
| Tier 1 10개 SWE-bench-Lite 회귀 러너 | 대기 |
| 3-Layer Defense 고도화 | 대기 |
| Prompt caching (Anthropic `cache_control` + OpenAI `prompt_cache_key`) | 대기 |
| **체크포인트**: Tier 1 5개 SWE-bench 태스크 성공률 ≥80% | 대기 |

---

## 4. 디자인 결정 대기 (5건) — `hanimo-code-desktop`

관련 문서: `docs/porting/HANIMO-DESKTOP-DESIGN-PLAN-2026-04-23.md` · mock: `designs/hanimo-desktop-v1.html`

| # | 질문 | 옵션 | 영향 |
|:-:|---|---|---|
| D1 | 아이콘 라이브러리 | A) Lucide (현재 v1 mock) · B) Phosphor (꿀벌 메타포 친화) · C) 자체 SVG (오리지널) | 브랜드 독창성 vs 개발 속도 |
| D2 | 코드 에디터 엔진 | A) CodeMirror 6 (techai-ide 기본) · B) Monaco (VSCode 엔진) · C) Shiki 정적 | 번들 크기 vs 언어 지원 |
| D3 | 폰트 | A) Geist + JetBrains Mono (현재) · B) Pretendard + D2 Coding (한국어 최적) | 한국어 가독성 |
| D4 | Theme 개수 | A) 8개 (현재) · B) 5개만 (Honey/Slate/Ocean/Claude/Paper) · C) 사용자 커스텀 JSON | 번들 크기 |
| D5 | 꿀벌 로고 아이콘 | 현재 `hexagon` lucide 사용 · 오리지널 꿀벌 SVG 제작 필요? | 브랜드 기억도 |

## 5. 세컨더리 TODO

- [ ] obsidian 내부 폴더명 rename: `obsidian/projects/modolai/` → `hanimo-webui/`, `modolrag/` → `hanimo-rag/` (리브랜드 문서 §6)
- [ ] hanimo-webui DB 식별자(`modol`, `modol_dev`) 변경 여부 검토 (마이그레이션 필수)
- [ ] hanimo-rag Docker 파일 유지/제거 결정 (본인 미사용)
- [ ] CONTRIBUTING.md · good-first-issue 라벨 30개 작성 (Sprint 5)
- [ ] `hanimo.dev` 도메인 인수 및 DNS 설정

## 6. Tooling Policy 편입 (2026-04-24 추가)

관련 문서: `docs/superpowers/TOOLING-POLICY-2026-04-24.md`

- [ ] Sprint 1 착수 전 `REPUTATION-STRATEGY`를 `superpowers:writing-plans` 포맷으로 재구성 (Phase별 verification 기준 추가)
- [ ] `designs/hanimo-desktop-v1.html` `web-design-guidelines` 스킬로 A11y·contrast·keyboard 검증
- [ ] Phase 0~5 각 진입 시 `writing-plans` → `executing-plans` 체인 적용
- [ ] 2개 이상 독립 변경 발생 시 gstack 스택 분리 규칙 적용
- [ ] 완료 선언 직전 `superpowers:verification-before-completion` 항상 invoke
- [ ] UI 작업 시 `frontend-design` 또는 `tui-designer` 명시적 invoke
