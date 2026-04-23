# Session 2026-04-23 — Strategy Pivot + Bidirectional Port Analysis

> 작업자: 김지원 (flykimjiwon) + Claude Opus 4.7 (1M context)
> 세션 목적: hanimo 생태계 리브랜드 이후 전략 재정렬 + hanimo ↔ TECHAI_CODE 양방향 포팅 본질 필터링 + IDE 네이밍 통일
> 선행 세션: `docs/SESSION-2026-04-11-INDEX.md`, `docs/SESSION-2026-04-16-17-PORTING-SUMMARY.md`
> 모체 리브랜드 로그: `~/Desktop/kimjiwon/obsidian/projects/hanimo_리브랜드_2026-04-23.md`

---

## 1. 결정 사항 요약

### 1.1 피벗 결정

| # | 결정 | 근거 |
|:-:|---|---|
| 1 | hanimo-code + hanimo-code-desktop = **완전 무료 OSS · 유료화 금지** | 목적은 김지원 개인 명성 |
| 2 | hanimo-webui는 자체 운영 중 · 본 문서 범위 제외 | 이미 관리자 기능 풍부 |
| 3 | 5-surface 전략 → **Code + Desktop(IDE) 올인** | 1인 리소스 한계 |
| 4 | Desktop 별도 앱 · VSCode Ext · Spark Native · AgentRank 통합 = **전부 보류** | IDE가 Desktop 수요 흡수 |
| 5 | IDE 공식 레포명 = **`hanimo-code-desktop`** | 리브랜드 문서 §1 기준 |
| 6 | 포팅 시 본질(외부망/폐쇄망) 필터 적용 · 한쪽 전용은 절대 포팅 금지 | |

### 1.2 이기는 각도 3개 (REPUTATION-STRATEGY §2)

1. **Korea-First Coding Agent** — 한국어 intent/clarify + 한국 스택 지식팩
2. **Hash-Anchored Safety** — `hashline_edit` 브랜드화, *"agent that can't silently overwrite"*
3. **Local-First + Air-gap Ready** — Ollama/vLLM 1급 + baked 빌드 (무료)

### 1.3 버리는 것 (REPUTATION-STRATEGY §3)

- 유료·엔터프라이즈 세일즈 · 14+ provider 자랑 · 5-surface 병렬 · Desktop 별도 앱 · VSCode Extension · Spark Native · AgentRank 통합 · 영어권 우선 런치 · 기능 합산 자랑

---

## 2. 이번 세션 산출물

### 2.1 신규 문서 (3개)

| 경로 | 내용 |
|---|---|
| `docs/strategy/REPUTATION-STRATEGY-2026-04-23.md` | 경쟁 전략 마스터 · 90일 6-sprint · 12개월 이정표 · 킬러 무브 · 결정 대기 5건 |
| `docs/porting/BIDIRECTIONAL-ANALYSIS-2026-04-23.md` | hanimo ↔ TECHAI 본질 필터링 포팅 리스트 (A 8건 · B 4건 · 양쪽 9건 · 금지 12건) |
| `docs/TODO-decisions-2026-04-23.md` | 실제 착수 전 확정 필요 5건 |

### 2.2 수정 문서 (1개)

| 경로 | 변경 |
|---|---|
| `docs/porting/IDE_PORTING_PLAN.md` | `hanimo-ide` → `hanimo-code-desktop` 전역 치환 (25건) |

### 2.3 상위 문서 재정의 (기존 유효 / 본 세션으로 갱신)

| 기존 문서 | 본 세션 이후 위상 |
|---|---|
| `docs/PLATFORM-PLAN-2026-04.md` | ⚠️ 5-surface 축소됨 — REPUTATION-STRATEGY가 상위 |
| `docs/COMPETITIVE-LANDSCAPE-2026-04.md` | ✅ 분석 유효 · 결론은 REPUTATION-STRATEGY로 갱신 |
| `docs/DESKTOP-PLAN-2026-04.md` | ❌ 보류 (IDE로 흡수) — `status: deferred` 추가 필요 |
| `docs/vscode-extension-detailed-plan.md` | ❌ 보류 — `status: deferred` 추가 필요 |
| `docs/AGENTRANK-INTEGRATION-2026-04.md` | ❌ hanimo-community 레포에서 독립 진행 · code 쪽 통합은 보류 |
| `docs/SPARK-NATIVE-APP-ROADMAP.md` | ❌ 보류 — Go CLI v1만 유지 |
| `docs/EMBEDDED-BROWSER-PLAN.md` | ❌ 보류 — IDE 기본 Webview로 충분 |

---

## 3. 메모리 반영 (~/.claude/projects/.../memory/)

| 파일 | 역할 |
|---|---|
| `project_hanimo_techai.md` | hanimo ↔ TECHAI 페어 개요 |
| `project_hanimo_ecosystem_naming.md` | 4-repo 공식 네이밍 + IDE = hanimo-code-desktop |
| `project_scope_pivot_2026-04-23.md` | 5-surface → Code+IDE 축소 |
| `feedback_monetization_policy.md` | 유료화 금지, 명성 목적 (code+desktop 한정) |
| `MEMORY.md` | 인덱스 갱신 |

---

## 4. 실행 안 함 (의도적 보류)

사용자 지시: **"전부 문서화만 해놔"** (2026-04-23).

본 세션은 계획·문서화만 진행. 아래는 **별도 세션에서 실행 예정**:

- Sprint 1 (Tier 1 완벽주의) 실제 코드 작업
- 작업 A (hanimo → TECHAI 포팅)
- 작업 B (TECHAI → hanimo 포팅)
- hanimo-code-desktop 포팅 (2026-04-24~)
- 보류 문서들에 `status: deferred` frontmatter 실제 추가

각 항목은 `docs/TODO-decisions-2026-04-23.md`와 관련 플랜 문서 참조하여 착수.
