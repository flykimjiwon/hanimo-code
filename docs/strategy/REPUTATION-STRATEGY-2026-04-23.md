# hanimo Reputation Strategy — 2026-04-23

> **작성일**: 2026-04-23
> **작성자**: 김지원 (flykimjiwon) + Claude Opus 4.7
> **리브랜드 기준**: 오늘 `Modol*` → `hanimo-*` 4개 레포 전면 리브랜드 완료.
> 원본 기록: `/Users/jiwonkim/Desktop/kimjiwon/obsidian/projects/hanimo_리브랜드_2026-04-23.md`
>
> **전제 피벗**:
> 1. `hanimo-code` + `hanimo-code-desktop`(IDE) = **완전 무료 오픈소스 + 누구나 기여 가능**. 유료화 금지.
> 2. 목적 = **김지원 개인 명성(reputation) · 영향력 · 커뮤니티 매력도**. 매출 아님.
> 3. `hanimo-webui`는 이미 관리자 기능이 풍부하게 자체 운영 중 — 본 문서 **범위 제외**.
> 4. `hanimo-rag`(pip 설치형)·`hanimo-community`(AgentRank 프로토타입)는 레포 존재 — 본 문서 **범위 제외**.
> 5. 별도 Desktop 앱·VSCode Extension·Spark Native App 계획은 **전부 보류** (IDE가 Desktop 수요 흡수).
> 6. TECHAI_CODE(신한 폐쇄망)는 downstream. 본 문서는 외부망 오픈소스 hanimo만 다룸.

## 0.5. hanimo 생태계 정식 네이밍 (2026-04-23 확정)

| 레포 | GitHub | 역할 | 본 문서 범위 |
|---|---|---|:-:|
| **hanimo-code** | [flykimjiwon/hanimo-code](https://github.com/flykimjiwon/hanimo-code) | Go 터미널 AI 코딩 에이전트 | ✅ |
| **hanimo-code-desktop** | (예정) | IDE 변형 (Wails 기반, 2026-04-24~ 포팅) | ✅ |
| hanimo-webui | [flykimjiwon/hanimo-webui](https://github.com/flykimjiwon/hanimo-webui) | Next.js 웹 AI 챗 플랫폼 (자체 운영 중) | ❌ |
| hanimo-rag | [flykimjiwon/hanimo-rag](https://github.com/flykimjiwon/hanimo-rag) | PostgreSQL RAG 엔진 (pip 설치형) | ❌ |
| hanimo-community | [flykimjiwon/hanimo-community](https://github.com/flykimjiwon/hanimo-community) | 커뮤니티 + AgentRank 프로토타입 | ❌ |

**브랜드 허브 결정 대기** (obsidian 리브랜드 문서 §6):
- A) 새 레포 `flykimjiwon/hanimo` 만들어 랜딩·문서 허브로 사용
- B) `hanimo-community`가 허브 겸임
- C) `hanimo-code`의 `landing-mockups/` 승격

**식별자 규칙** (혼동 방지):
- GitHub 레포명·CLI 커맨드·npm/PyPI 배포명: **hyphen** (`hanimo-code`)
- Python 모듈·import 경로: **underscore** (`hanimo_rag`)
- 환경변수 prefix: **UPPER underscore** (`HANIMO_CODE_*`)
>
> **상위 문서 재정의**:
> - `docs/PLATFORM-PLAN-2026-04.md` → 본 문서로 대체 (5-surface 축소)
> - `docs/COMPETITIVE-LANDSCAPE-2026-04.md` → 분석은 유효, 결론 본 문서로 갱신
> - `docs/DESKTOP-PLAN-2026-04.md` → **보류** (IDE에 흡수)
> - `docs/AGENTRANK-INTEGRATION-2026-04.md` → **Phase ∞**
> - `docs/vscode-extension-detailed-plan.md` → **보류**

---

## 0. TL;DR (30초)

- **집중 축** = 오직 **hanimo-code (Go TUI)** + **hanimo-code-desktop (Wails IDE, 다음 주 포팅)**. 다른 전부 ❌.
- **성공 지표** = GitHub ⭐, 컨트리뷰터 수, HN 프론트, 미디어 커버리지, "hash-anchored edit" 용어 인용.
- **이기는 3 각도**:
  1. **Korea-First** — 한국어 intent/clarify, 한국 스택 지식팩
  2. **Hash-Anchored Safety** — "AI가 네 파일을 덮어쓰지 않음"의 브랜드화
  3. **Local-First + Air-gap Ready** — Ollama/vLLM 1급 + baked 빌드 (수요는 있지만 **무료**로 제공)
- **버리는 것들** — 유료화, 엔터프라이즈 세일즈, 5-surface 병렬, 14+ provider 자랑, Desktop/VSCode Ext/Spark Native 동시 추진.
- **6개 스프린트 (12주)** — Tier-1 완벽주의 → Repo-map/Skill/Hooks → Subagent/Permission → Korea Pack → Launch 준비 → Public β.
- **12개월 목표** — ⭐ 5,000+, 활성 컨트리뷰터 20명+, Show HN front, 한국 개발자 컨퍼런스 발표 1회.

---

## 1. 냉정한 시장 진단

### 1.1 레드오션 현실

| 경쟁자 | 규모 | 후발자가 못 따라잡는 것 |
|---|---|---|
| Claude Code | 공식 (Anthropic) · SWE-bench 80.9% · Agent Teams · 1M ctx | **모델 자체 품질 + 브랜드** |
| Codex CLI | 공식 (OpenAI) · Terminal-Bench 77.3% | **모델 자체** |
| opencode | **⭐ 95K · 월 2.5M 유저 · 75+ provider** | **네트워크 효과 + 커뮤니티** |
| Aider | **⭐ 39K · 주 15B 토큰** · Git 일체화 | "Aider = 페어프로그래밍" 브랜드 각인 |
| Cursor Agent | IDE 딥 통합 | **IDE 실시간 편집 UX** |
| Crush (Charm) | Bubble Tea 계열 · LSP | Charm.sh 브랜드 파워 |
| Continue / Goose / Kilo / Cline | 다양한 오픈소스 | — |

**결론:** 정면 승부(provider 수·모델 품질·벤치 점수) = **패배 확정**.
opencode가 이미 "75+ provider / Bubble Tea TUI / 무료"를 차지 → hanimo가 14+ provider로 붙으면 "opencode 하위호환" 포지션이 될 뿐.

### 1.2 후발자 명성 전략의 공식

> **"남보다 먼저 한 것 하나 + 남이 안 하는 것 하나 + 한국어 우위 하나"** 세 줄을 영어와 한국어로 일관되게 반복.

명성은 "기능 합산"으로 쌓이지 않음. **문장 하나**로 쌓임. 예:
- Aider → "git-native AI pair programmer"
- opencode → "provider-agnostic TUI agent"
- Claude Code → "agent teams that work while you sleep"
- **hanimo** → **"The agent that asks before it breaks your code."** (제안)

---

## 2. 이길 수 있는 3 각도 (정직하게)

### 🎯 각도 1 — **Korea-First Coding Agent**

**근거**:
- opencode / Aider / Claude Code 한국어 UX 전부 열악 (영어 plan, 영어 에러, 영어 commit tone).
- Cursor 한국 사용자 수십만 명이 영어 UX로 고생 중.
- 한국 스타트업·주니어 개발자 시장 = **번역되지 않은 허기**.

**실행**:
- Intent detection (한글 모호성) ← 이미 있음 (`agents/intent.go`) → **브랜드화**
- clarify-first ASK_USER 구조화 ← 이미 있음 → **영상화**
- 한국어 error message / 한국어 plan / 한국어 commit message tone
- **한국 스택 지식팩**: Spring + JEUS/Tomcat + Oracle/Tibero + NiceDB + 네이버 오픈소스 (LINE FOSS, NAVER Labs) + 카카오(toasted-bread 등) + 토스(slash) + 당근마켓(Daangn OSS)
- **한국 개발 문화 지식**: `README.md` 한국어 템플릿, 커밋 컨벤션 한국어 가이드, 팀 온보딩 문서 자동 생성 스킬

**왜 못 따라오나**: opencode 메인테이너 한국어 모름. Anthropic·OpenAI는 한국을 tier-2 취급. **언어 장벽 해자**.

### 🎯 각도 2 — **Hash-Anchored Safety** (브랜드 무기)

**근거**:
- `hashline_edit`는 **업계 유일**. "AI가 내 파일을 덮어쓰는" 공포는 모든 사용자 공통.
- 이 공포는 **Aider를 쓰게 만든 주 동기**와 동일 (`git commit -a && aider`).
- hanimo는 git 의존 없이 파일 단위로 이 안전성을 제공 → 더 강함.

**실행**:
- `hashline_*` 기능을 **1문장 슬로건**으로 고유명사화: "Hash-Anchored Editing™"
- 모든 홈페이지·데모·Show HN에 이 한 줄 반복
- **논문 스타일 블로그 1편 작성**: *"Why AI Coding Agents Need Hash-Anchored Edits"*
  - 경쟁자 failure case 3개 재현 (Cursor/Aider가 stale overwrite로 커밋 날린 실제 사례)
  - MD5 4자리 앵커 방식 수학적 증명 (충돌 확률, 윈도우 크기 tradeoff)
  - hanimo / Aider / Cursor 비교 표
- **HN 타이틀 제안**: *"Show HN: hanimo — an AI coding agent that can't silently overwrite your edits"*

**왜 못 따라오나**: 구조적으로 복제 가능하지만 **용어를 먼저 찜한 쪽이 이김**. Aider가 "repo-map"을, Cursor가 "Composer"를 찜한 방식.

### 🎯 각도 3 — **Local-First + Air-gap Ready** (무료로)

**근거**:
- 한국 금융/공공/제조 = 외부 클라우드 금지. Claude Code·Cursor는 본질적으로 클라우드 종속 → 진입 불가.
- 하지만 **돈 안 받으므로** 이 시장을 **유료 매출 근거로 쓰지 않음**. 대신:
  - "hanimo로 폐쇄망에서 Claude Code를 대체했다"는 **블로그 후기** 자체가 명성.
  - TECHAI_CODE가 이미 증명한 패턴을 OSS로 공개 → 한국 오픈소스 계에서 고유 포지션.

**실행**:
- Baked 빌드 (엔드포인트·키·모델 고정 배포) **문서화 + 자동화** (`make bake`)
- Ollama `/api/tags` 자동 감지 + 로컬 모델 retry/fallback 1급 시민화
- vLLM endpoint 예시 템플릿 내장
- SBOM + 서명 바이너리
- **Air-gapped deploy playbook 공개**: *"Run hanimo on a network with zero outbound access, in 10 minutes"*

**왜 못 따라오나**: Claude Code/Cursor는 비즈니스 모델상 on-prem 무료 제공 불가. opencode는 다중 provider지만 "감사·보안 하드닝" 문서 없음.

---

## 3. 버릴 것들 (Negative Scope) — **가장 중요**

| 버려라 | 이유 |
|---|---|
| ❌ Desktop (Wails v3 별도 앱) | **hanimo-code-desktop(IDE)가 같은 수요 커버**. `docs/DESKTOP-PLAN-2026-04.md` 폐기/흡수 |
| ❌ VSCode Extension | monorepo 리팩터 비용 대비 효과 작음. `docs/vscode-extension-detailed-plan.md` 보류 |
| ❌ Spark Native App | DGX Spark 전용 — 수요 작음. 현 Go CLI 유지, native app 계획 보류 |
| ❌ `hanimo-code`에서 webui/rag/community 기능 통합 | 각 레포가 별도 제품. 크로스 링크 수준에서만 참조 |
| ❌ AgentRank 통합 작업 (code 쪽에서) | `hanimo-community` 레포에서 독립 진행 |
| ❌ 14+ provider 자랑 | opencode 75+가 압도. **Tier 1 10개**만 쥐고 "완벽 지원"으로 각인 |
| ❌ 영어권 우선 런치 | 한국 10K 유저 > 북미 100 유저 (니치 장악 후 확장) |
| ❌ 유료·엔터프라이즈 세일즈 | 본 전략 전제 — 명성 목적이지 매출 아님 |
| ❌ 기능 합산 자랑 | "36 tools / 42 commands" 식 README 금지. **슬로건 3줄로 축약** |

---

## 4. 90일 실행 플랜 (2주 스프린트 × 6)

### Sprint 1 (Week 1-2) — "Tier 1 완벽주의"

| 태스크 | 산출 |
|---|---|
| Certified Model System 코드화 | `hanimo models` 명령 + Tier 시각화 표 |
| Tier 1 10개 회귀 러너 | SWE-bench-Lite 일부 태스크 자체 러너 (cron 주 1회) |
| 3-Layer Defense 고도화 | pre-send regex 확장 · post-receive secret leak · loop hash abort |
| Prompt caching 도입 | Anthropic `cache_control` + OpenAI `prompt_cache_key` + Novita 대응 |
| **체크포인트** | Tier 1 5개 SWE-bench 태스크 성공률 ≥80% |

### Sprint 2 (Week 3-4) — "Repo-map + Skill + Hooks"

| 태스크 | 산출 |
|---|---|
| `internal/repomap/` 신규 | tree-sitter + PageRank + SQLite cache (Aider 포팅, 라이선스 MIT 확인) |
| `internal/skills/` 고도화 | `$ARGUMENTS` · inline `!shell` · learned skill 자동 추출 |
| `internal/hooks/` 고도화 | 8 events + `.hanimo/hooks.yaml` |
| `internal/permission/` 5-mode | Shift+Tab 순환 + 학습형 `permissions.yaml` |
| **체크포인트** | 10K 파일 레포 repo-map cold <300ms · warm <80ms |

### Sprint 3 (Week 5-6) — "Subagent + Doom-loop + Git Worktree"

| 태스크 | 산출 |
|---|---|
| `agents/subagent.go` | context fork + 요약 반환 |
| Doom-loop detector 고도화 | 최근 3회 tool hash → abort + iteration cap env |
| Git worktree 병렬 실행 | `hanimo worktree` CLI |
| Session Browser UI + `--resume <id>` | C1 |
| **체크포인트** | 동시 3 subagent 성공, loop abort 재현 테스트, worktree 3개 병렬 실험 |

### Sprint 4 (Week 7-8) — "Korea Pack + Local-First Polish"

| 태스크 | 산출 |
|---|---|
| 한국 스택 지식팩 | Spring / JEUS / Tibero / 네이버 · 카카오 · 토스 · 당근 OSS 패턴 |
| Ollama `/api/tags` 자동 감지 | 로컬 모델 retry/fallback |
| Baked 빌드 파이프라인 | `make bake` + 서명 + SBOM |
| JSON headless mode | `--output json` |
| DECSET 2026 synchronized output | 스트리밍 깜빡임 제거 |
| **체크포인트** | 한국 네이티브 5명 베타 NPS ≥40 |

### Sprint 5 (Week 9-10) — "Launch 준비"

| 태스크 | 산출 |
|---|---|
| `hanimo.dev` 랜딩 | comparison-first · 3 각도 메시지 · 데모 GIF |
| README 전면 개편 | 한국어/영어 병행, 슬로건 3줄 상단 |
| 논문 스타일 블로그 1편 | *"Why AI Coding Agents Need Hash-Anchored Edits"* |
| 데모 영상 3편 × 3분 | 한국어 / 영어 / 폐쇄망 시나리오 |
| CONTRIBUTING.md · good-first-issue 30개 | 기여자 매력도 |
| IDE 포팅 완료 (다른 스프린트에서 병행) | `hanimo-code-desktop` MVP |
| **체크포인트** | 한국 OSS 기여자 5명 pre-launch 리뷰 완료 |

### Sprint 6 (Week 11-12) — "Public β"

| 태스크 | 산출 |
|---|---|
| Show HN · r/LocalLLaMA · r/golang 동시 런치 | 3 창구 |
| 한국 커뮤니티 공지 | okky · 클리앙 개발프로그래밍 · Disquiet · GeekNews |
| 두 번째 블로그 | *"How we run a Claude Code alternative on air-gapped networks, for free"* |
| Hacker News front 입장 시 대응 문서 | FAQ · 벤치 표 · 로드맵 |
| **목표** | ⭐ 500 · 1주 DAU 200 · 컨트리뷰터 5명 확보 |

---

## 5. 6·12개월 이정표 (명성 지표만)

| 시점 | ⭐ | 활성 컨트리뷰터 | 미디어 | 이벤트 |
|:-:|:-:|:-:|:-:|---|
| +3개월 | 1,500 | 8명 | 한국 기술 블로그 3개 언급 | — |
| +6개월 | 3,000 | 15명 | HN front 1회 · Hacker News 월간 top 3 | FEConf / AWSKRUG / GDG 1회 발표 |
| +9개월 | 4,000 | 18명 | Charm.sh · Anthropic 블로그 피인용 | — |
| +12개월 | **5,000+** | **20명+** | 한국 개발자 컨퍼런스 keynote 1회 | 한국 오픈소스 시상 후보 |

**지표 의식적 제외**: 매출, DAU/MAU(절대 수치보다 비율), 엔터프라이즈 로고.

---

## 6. 킬러 무브 (1~2개로 충분)

### 🔑 Move 1 — **"Clarify-first as default"**

경쟁자들은 즉시 실행. hanimo는 **모호한 요청이면 항상 한 번 되묻는다**. 이미 있는 기능 → **브랜드화만** 필요.

- 데모 영상 1편 = **경쟁자 5개가 잘못된 코드를 쓸 때 hanimo만 되묻는 split-screen**.
- 슬로건: *"It asks. Others assume."*

### 🔑 Move 2 — **"Undo Everything"**

`hashline_edit` + `snapshot` + `git worktree` 조합으로 **어떤 AI 작업도 1초 롤백**.

- 영상 1편 = hanimo가 대규모 리팩터링 후 `/undo all` 한 방으로 복원.
- 슬로건: *"The only AI agent with a true undo button."*

---

## 7. 기여자 매력도 체크리스트

명성 = 별 + 컨트리뷰터. 컨트리뷰터 유입을 위한 필수 정비.

- [ ] `CONTRIBUTING.md` — 30분 안에 첫 PR 날리는 방법
- [ ] `ARCHITECTURE.md` — 12-layer 시각화 + 각 패키지 한 문단 설명
- [ ] `good-first-issue` 라벨 30개 + **실제 해결 가능한 난이도 검증**
- [ ] `help-wanted` 라벨 10개 (medium 난이도)
- [ ] PR 응답 **24시간 내 첫 코멘트** 약속
- [ ] DCO/CLA **없음** (진입 장벽 최소화) — Apache-2.0로 충분
- [ ] Discord/Slack 대신 **GitHub Discussions** (권위 집중)
- [ ] 한국 기여자용 한국어 이슈 템플릿 병행

---

## 8. 당신이 지금 결정해야 할 5가지

Sprint 1 실제 착수 전에 확정 필요.

| # | 질문 | 옵션 |
|:-:|---|---|
| 1 | 주당 투입 시간 | ① 10h (취미) ② 20h ③ 30h ④ 풀타임 |
| 2 | 한국 vs 글로벌 런치 순서 | **권장 = 한국 먼저 → 3개월 후 영어권** |
| 3 | 슬로건 1문장 확정 | A *"It asks. Others assume."* / B *"The agent that can't silently overwrite your edits."* / C (자체 제안) |
| 4 | **보류 문서 공식 폐기** | `docs/DESKTOP-PLAN-2026-04.md`·`docs/vscode-extension-detailed-plan.md`·`docs/AGENTRANK-INTEGRATION-2026-04.md` frontmatter에 **`status: deferred`** 추가 여부 |
| 5 | **브랜드 허브 레포** | A) 새 레포 `flykimjiwon/hanimo` / B) `hanimo-community` 겸임 / C) `hanimo-code/landing-mockups/` 승격 |

---

## 9. 관련 문서

- `docs/COMPETITIVE-LANDSCAPE-2026-04.md` — 경쟁 분석 원본 (분석은 유효)
- `docs/MASTER-OVERVIEW-2026-04.md` — Phase A/B/C 기술 로드맵
- `docs/DASHBOARD.md` — Certified Model System Tier 구조
- `docs/strategy/hanimo-certified-models-v0.2.4.md` — Tier 1 10개 모델 명세
- `docs/strategy/project-strategy-2026-04.md` — 이전 전략 문서
- `docs/porting/BIDIRECTIONAL-ANALYSIS-2026-04-23.md` — 본 세션 양방향 포팅 분석
- `docs/porting/IDE_PORTING_PLAN.md` — `hanimo-code-desktop` 포팅 플랜 (2026-04-24~)
- `~/Desktop/kimjiwon/obsidian/projects/hanimo_리브랜드_2026-04-23.md` — 4-repo 리브랜드 작업 로그 (원본)
