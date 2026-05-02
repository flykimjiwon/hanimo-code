# hanimo Brand System

> hanimo 생태계 5 repo의 *시각·언어·메시징* 통일을 위한 단일 source.
> 변경 시 모든 hanimo 레포에 영향이 있으므로 PR을 통한 명시적 승인 필요.

## 빠른 시작

| 무엇을 찾아왔나 | 어디로 |
|---|---|
| 브랜드 컨셉·정체성 (모돌·꿀벌·하니스) 알고 싶다 | [`BRAND.md`](BRAND.md) |
| 색·폰트·radius 토큰을 *코드에* 박고 싶다 | [`tokens.css`](tokens.css) → import |
| 토큰 reference 표 (어떤 변수가 뭘 의미) | [`tokens.md`](tokens.md) |
| README 통일된 11-section 템플릿 | [`readme-template.md`](readme-template.md) |
| 로고 SVG를 만드는데 사양 필요 | [`logo-spec.md`](logo-spec.md) |
| hanimo-webui의 Phosphor → Lucide 점진 이전 | [`lucide-migration.md`](lucide-migration.md) |

## 5 repo 적용 범위

| repo | 적용 대상 |
|---|---|
| **hanimo-code** | TUI 테마(이미 적용) · 데스크톱 IDE 8 테마 · 본 docs |
| **hanimo-webui** | Next.js Tailwind config · `app/globals.css` · `components/icons` |
| **hanimo-rag** | `dashboard/src/index.css` · README 배지 |
| **hanimo-community** | 랜딩 페이지 (예정) · README 배지 |
| **hanimo-toy-coffee** | 적용 대상 *아님* (내부 샌드박스) |

## 변경 절차

이 폴더 안의 문서·자산을 변경하려면:

1. **PR 필수** — 자가 머지 금지
2. **5-repo 일관성 검토** — 한 곳을 바꾸면 다른 4 repo도 영향 분석
3. **단계별 적용** — 컨셉 변경(BRAND.md) → 토큰 갱신(tokens.css) → 각 repo 개별 PR
4. **legacy 보존** — 기존 자산은 폐기하지 말고 *legacy 표기*로 보존 (예: `hanimo-brand-design.html`은 v0 legacy)

## Legacy 자산 (보존)

| 파일 | 상태 | 비고 |
|---|---|---|
| `hanimo-brand-design.html` | **legacy v0** | Catppuccin 베이스 + lavender accent. Modol 컨셉 도입 전 자산. 기록용 보존, 신규 작업은 본 폴더의 markdown 자산 우선 참조 |

## 메모리 연결

자동 메모리 시스템이 본 컨셉을 자동 로드:
- `project_hanimo_brand_concept.md` (2026-05-03 박힘) — Modol the Honey-Bee Bichon 핵심 요약

향후 모든 디자인 작업 세션에서 본 컨셉이 자동 컨텍스트로 등장.

## 적용 진행 상태 (2026-05-03 기준)

- [x] Honey amber `#f5a623` 라이선스 배지 5 repo 통일 (hanimo-code/rag/webui)
- [x] Lucide 단일 source 컨벤션 박음 (hanimo-webui `app/components/icons/index.jsx`)
- [x] 본 brand docs 7종 신규 작성
- [ ] `tokens.css` 5 repo import 적용
- [ ] 4 repo README hero 통일 (로고 SVG 후 권장)
- [ ] 모돌-꿀벌 SVG 로고 v1 작성
- [ ] hanimo.dev 통합 랜딩 페이지
- [ ] Phosphor → Lucide 도메인 단위 점진 치환 (chat → workflow → admin)
