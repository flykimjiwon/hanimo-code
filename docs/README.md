# hanimo Documentation Hub

> hanimo 생태계 모든 문서·HTML 쇼케이스·정책·전략·세션·브랜드 자산을
> 한 페이지에 정리한 진입점. 시각 페이지는 [`index.html`](index.html)에서
> 동일한 내용을 확인.

---

## 처음 오셨다면 — 30초 진입 가이드

1. **hanimo가 뭔가요?** → [`landing/hanimo-products.html`](landing/hanimo-products.html) (5 제품 카드 + 마스코트 + 멀티세션 워크플로우)
2. **기술적으로 어떤 구조?** → [`landing/hanimo-framework.html`](landing/hanimo-framework.html) (12-layer 아키텍처 + 모듈 + 도구)
3. **브랜드·디자인은?** → [`branding/BRAND.md`](branding/BRAND.md) (Modol 마스코트 컨셉) + [`branding/tokens.md`](branding/tokens.md)
4. **회사에서 가져다 쓸 때 라이선스?** → [`policy/README.md`](policy/README.md) (5 정책 인덱스)
5. **최근 변화?** → [`SESSION-2026-05-03-CLOSING.md`](SESSION-2026-05-03-CLOSING.md) (마지막 세션 30초 복구)

---

## 1. HTML Showcases

브라우저에서 바로 열 수 있는 비주얼 페이지. 외부 CDN/폰트 없음, 단일 파일,
에어갭 환경에서도 정상 렌더링 (telemetry-and-privacy 정합).

| Status | Document | 내용 |
|---|---|---|
| **Active · v1** | [`landing/hanimo-products.html`](landing/hanimo-products.html) | 4 제품 카드 + Author Manifest + Multi-Terminal·Multi-Session Workflow + tmux 4-pane + Prompt 템플릿 |
| **Active · v1** | [`landing/hanimo-framework.html`](landing/hanimo-framework.html) | 12-Layer Architecture + 16 Module Map + 23+ Tool Catalog + Modes + Build Profiles + Doc Index |
| **Active** | [`landing/hanimo-dev-landing.html`](landing/hanimo-dev-landing.html) | 개발자 랜딩 (기존 자산) |
| **Legacy · v0** | [`branding/hanimo-brand-design.html`](branding/hanimo-brand-design.html) | Catppuccin + lavender 기반 v0 brand 시스템. 신규 작업은 `BRAND.md` 우선 참조 |
| **Hub** | [`index.html`](index.html) | 본 페이지의 시각 버전 (HTML) |

---

## 2. Brand & Design System

Modol the Honey-Bee Bichon 마스코트 컨셉, Honey 팔레트 토큰, 로고 사양,
README 템플릿, Lucide 단일 아이콘 source 마이그레이션.

| File | 역할 |
|---|---|
| [`branding/README.md`](branding/README.md) | Brand system 인덱스 + 적용 진행 상태 |
| [`branding/BRAND.md`](branding/BRAND.md) | 컨셉·메타포·컬러·타이포·메시징 톤 10절 풀버전 |
| [`branding/tokens.css`](branding/tokens.css) | 5 repo 공유 단일 디자인 토큰 (Honey palette + 타이포 + radius + shadow) |
| [`branding/tokens.md`](branding/tokens.md) | tokens.css의 모든 변수 reference + Tailwind/HTML/React 사용 예시 |
| [`branding/readme-template.md`](branding/readme-template.md) | 11-section 통일 README 템플릿 |
| [`branding/logo-spec.md`](branding/logo-spec.md) | 모돌-꿀벌 SVG 사양 + 3 variant + AI 일러스트 프롬프트 |
| [`branding/lucide-migration.md`](branding/lucide-migration.md) | Phosphor → Lucide 4-phase 점진 이전 체크리스트 |

### Logo Assets (v0 placeholder)

- [`branding/assets/hanimo-logo-square.svg`](branding/assets/hanimo-logo-square.svg) — 256×256, favicon/avatar/앱 아이콘용
- [`branding/assets/hanimo-logo-horizontal.svg`](branding/assets/hanimo-logo-horizontal.svg) — 600×200, README hero/랜딩용
- [`branding/assets/hanimo-logo-compact.svg`](branding/assets/hanimo-logo-compact.svg) — 64×64, 작은 favicon/UI 인라인

---

## 3. Policies

5 정책 문서. 5 repo에 공통 적용되며, hanimo-code가 canonical.

| Policy | Document | Scope |
|---|---|---|
| Index + FAQ | [`policy/README.md`](policy/README.md) | 5 정책 진입점 |
| Copyright | [`policy/copyright.md`](policy/copyright.md) | 자동 귀속 · SPDX 헤더 · GPG 서명 · 한국저작권위원회 등록 |
| Trademark | [`policy/trademark-and-naming.md`](policy/trademark-and-naming.md) | 미등록 트레이드마크 · 포크 네이밍 · 공식 채널 |
| Telemetry & Privacy | [`policy/telemetry-and-privacy.md`](policy/telemetry-and-privacy.md) | 외부 통신 0 · 에어갭 · 향후 도입 시 강제 규칙 |
| LTS & On-Premises | [`policy/lts-onprem.md`](policy/lts-onprem.md) | 분기 LTS 캐덴스 · 12개월 백포트 · 인증 빌드 |

---

## 4. Strategy

의사결정·시장 분석·비전·인증 모델. 거부된 대안 + 트리거 액션이 함께 박힘.

| File | 역할 |
|---|---|
| [`strategy/STRATEGY-2026-04-27-LICENSE-DEPLOYMENT-IP.md`](strategy/STRATEGY-2026-04-27-LICENSE-DEPLOYMENT-IP.md) | 5 결정 + 거부 대안 6 + 시장 분업 + 시나리오 매트릭스 + 트리거 액션 |
| [`strategy/VISION-2026-04-25-MULTI-MODEL-MULTI-DEVICE.md`](strategy/VISION-2026-04-25-MULTI-MODEL-MULTI-DEVICE.md) | 6축 비전 + 시장 매트릭스 13 제품 + hanimo 자리 |
| [`strategy/MARKET-ANALYSIS-2026-04-24.md`](strategy/MARKET-ANALYSIS-2026-04-24.md) | 시장 진단 · IDE wedge · 플러그인 흡수 |
| [`strategy/REPUTATION-STRATEGY-2026-04-23.md`](strategy/REPUTATION-STRATEGY-2026-04-23.md) | 명성 wedge — 한국어 + OSS |
| [`strategy/hanimo-certified-models-v0.2.4.md`](strategy/hanimo-certified-models-v0.2.4.md) | Tier 1/2/3 모델 정책 |
| [`strategy/project-strategy-2026-04.md`](strategy/project-strategy-2026-04.md) | 종합 프로젝트 전략 |

---

## 5. Roadmap

| File | 역할 |
|---|---|
| [`roadmap/enhancement-plan-2026-04.md`](roadmap/enhancement-plan-2026-04.md) | Phase A/B/C 풀 플랜 (Repo-map · Skill · Hooks · Subagent · LSP · Caching) |
| [`roadmap/v0.2-v0.4-improvement-plan.md`](roadmap/v0.2-v0.4-improvement-plan.md) | 버전별 누적 작업 목록 |
| [`roadmap/2026-04-16-hanimo-enhancement-plan.md`](roadmap/2026-04-16-hanimo-enhancement-plan.md) | 중간 갱신본 |

---

## 6. Recent Sessions

시간순 세션 마감 기록. 새 세션 진입 시 *가장 최근 CLOSING*부터 읽으면 30초 안에 맥락 복구.

| Date | Document | Highlight |
|---|---|---|
| 2026-05-03 | [`SESSION-2026-05-03-CLOSING.md`](SESSION-2026-05-03-CLOSING.md) | **Brand v1** — Modol 마스코트, 5 repo 통일, Lucide Phase 1·2·3 |
| 2026-04-27 | [`SESSION-2026-04-27-CLOSING.md`](SESSION-2026-04-27-CLOSING.md) | 정책 인프라 — Apache 2.0, LTS-Onprem 빌드 골격, 정책 5종 |
| 2026-04-27 | [`SESSION-2026-04-27-INDEX.md`](SESSION-2026-04-27-INDEX.md) | 시간순 산출 기록 |
| 2026-04-25 | [`SESSION-2026-04-25-CLOSING.md`](SESSION-2026-04-25-CLOSING.md) | Phase 19 + v0.2.0 출시 — 한국 MCP 35종, Multi-OS 빌드 |
| 2026-04-25 | [`SESSION-2026-04-25-INDEX.md`](SESSION-2026-04-25-INDEX.md) | 같은 세션 시간순 |
| 2026-04-25 | [`SESSION-2026-04-25-RESUME.md`](SESSION-2026-04-25-RESUME.md) | 정책 + 큰 작업 후보 + 핵심 파일 |
| 2026-04-23 | [`SESSION-2026-04-23-INDEX.md`](SESSION-2026-04-23-INDEX.md) | IDE 디자인 v1 — Honey 팔레트, 14 Activity Bar 아이콘 |
| 2026-04-16/17 | [`SESSION-2026-04-16-17-PORTING-SUMMARY.md`](SESSION-2026-04-16-17-PORTING-SUMMARY.md) | 포팅 라운드 — hashline · auto · git tools · diagnostics |
| 2026-04-11 | [`SESSION-2026-04-11-INDEX.md`](SESSION-2026-04-11-INDEX.md) | Master Overview — 14+ provider, Tier 시스템 피벗 |

---

## 7. Code Infrastructure (루트)

| File | Purpose |
|---|---|
| [`../LICENSE`](../LICENSE) | Apache License 2.0 풀텍스트 |
| [`../NOTICE`](../NOTICE) | Kim Jiwon 단독 저작자 6조항 + 5-repo lineage + 트레이드마크 |
| [`../CONTRIBUTING.md`](../CONTRIBUTING.md) | 라이선스 부여 + DCO 서명 + SPDX 헤더 + PR 체크리스트 |
| [`../README.md`](../README.md) | hanimo-code 메인 README |
| [`../Makefile`](../Makefile) | 빌드 — vanilla / distro / sealed / **onprem** 4 프로파일 |
| [`../internal/build/`](../internal/build/) | `//go:build onprem` 컴파일 타임 분기 |

---

## 8. 기타 문서 (히스토리·연구·세부)

`docs/` 루트의 기타 markdown 문서:

- [`MASTER-OVERVIEW-2026-04.md`](MASTER-OVERVIEW-2026-04.md) — 전체 마스터 개요
- [`OPEN-SOURCE-LAUNCH-PLAN-2026-04.md`](OPEN-SOURCE-LAUNCH-PLAN-2026-04.md) — OSS 런치 플랜
- [`PLATFORM-PLAN-2026-04.md`](PLATFORM-PLAN-2026-04.md) — 플랫폼 전략
- [`COMPETITIVE-LANDSCAPE-2026-04.md`](COMPETITIVE-LANDSCAPE-2026-04.md) — 경쟁 환경
- [`AGENTRANK-INTEGRATION-2026-04.md`](AGENTRANK-INTEGRATION-2026-04.md) — 에이전트 랭킹 통합
- [`DESIGN-BRIEF-2026-04.md`](DESIGN-BRIEF-2026-04.md) — 디자인 브리프
- [`DESKTOP-PLAN-2026-04.md`](DESKTOP-PLAN-2026-04.md) — 데스크톱 IDE 계획
- [`EMBEDDED-BROWSER-PLAN.md`](EMBEDDED-BROWSER-PLAN.md) — 임베디드 브라우저 계획
- [`SPARK-NATIVE-APP-ROADMAP.md`](SPARK-NATIVE-APP-ROADMAP.md) — Spark 네이티브 앱
- [`skills-and-knowledge-guide.md`](skills-and-knowledge-guide.md) — Skill·Knowledge 가이드
- [`vscode-extension-detailed-plan.md`](vscode-extension-detailed-plan.md) — VS Code 확장 상세
- [`vscode-extension-plan-v2.md`](vscode-extension-plan-v2.md) — VS Code 확장 v2
- [`vscode-extension-open-questions.md`](vscode-extension-open-questions.md) — 미해결 질문

서브 폴더:
- [`branding/`](branding/) — 브랜드 시스템 (위 §2 참고)
- [`landing/`](landing/) — HTML 쇼케이스 (위 §1 참고)
- [`policy/`](policy/) — 5 정책 (위 §3 참고)
- [`strategy/`](strategy/) — 전략 (위 §4 참고)
- [`roadmap/`](roadmap/) — 로드맵 (위 §5 참고)
- [`research/`](research/) — 리서치 노트
- [`porting/`](porting/) — 포팅 기록 (hanimo ↔ TECHAI 등)
- [`superpowers/`](superpowers/) — superpowers 관련
- [`ecosystem/`](ecosystem/) — 생태계 정의 (별도 폴더, 본 hub와 별도 운영)
- [`html/`](html/) · [`images/`](images/) — 자산

---

## 변경 절차

본 hub(README.md, index.html)를 변경하려면:

1. *카테고리 추가* — 신규 정책·전략·세션 등 → 해당 섹션에 한 줄 추가
2. *문서 이동* — 링크 갱신 + git history 보존 (rename 사용)
3. *legacy 표기* — 폐기하지 말고 "Legacy" 배지 유지
4. *5 repo 영향 분석* — 본 hub는 hanimo-code canonical. 다른 repo가 참조하는 링크가 깨지지 않는지 확인

---

_Last updated: 2026-05-03 · brand v1 · Modol the Honey-Bee Bichon_
