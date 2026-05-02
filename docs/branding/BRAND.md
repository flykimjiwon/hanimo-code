# hanimo Brand — Modol the Honey-Bee Bichon

> hanimo의 정체성·메타포·시각 언어·메시징 톤을 한 문서에 정리.
> 본 문서는 2026-05-03에 정해진 컨셉을 기준으로 하며, 이전의
> Catppuccin/lavender 기반 자산(`hanimo-brand-design.html`)은
> v0 legacy로 보존 후 본 컨셉으로 점진 이전한다.

---

## 1. 이름의 어원

**hanimo (하니모)** = 두 단어의 합성어.

| 음절 | 어원 | 의미 |
|---|---|---|
| **하니** | *harness* + *honey* | AI agent harness(에이전트 골격) + 꿀(따뜻함, 천연, 한국적 정서) |
| **모** | *Modol* (모돌) | 사용자 본인의 미니비숑 반려견 이름 |

따라서 hanimo는 *기능적 메타포*(harness)와 *정서적 메타포*(반려견 모돌)의
이중 정체성을 갖는다. 어느 쪽도 폐기하지 않으며, 마스코트가 두 정체성을
한 캐릭터에 시각화한다.

---

## 2. 마스코트 — Modol the Honey-Bee Bichon

> **모돌이 꿀벌 옷을 입은 미니비숑 강아지.**

### 시각적 핵심
- **종**: 미니비숑(Bichon Frisé) — 흰색·구름같은 털, 검은 눈·코, 작고 둥근 얼굴
- **의상**: 꿀벌 옷 (Honey amber `#f5a623` + 검정 줄무늬), 작은 더듬이, 살짝 들린 날개
- **포즈 기본**: 정면 미소, 한 손(앞발) 살짝 들기. 친근하지만 *진지하게 일하는* 톤
- **표정**: 차분 + 자신감 + 약간의 장난기 (이모지 느낌은 *피하고* 일러스트 느낌)
- **변형**: 노트북 앞에 앉은 모돌 (코딩), 꿀단지 든 모돌 (수익화 X 명성), 망원경 든 모돌 (탐색·검색)

### 마스코트의 역할
- **로고 메인 비주얼** (favicon, GitHub avatar, README hero)
- **빈 상태(empty state)** — "검색 결과 없음", "세션 비어있음" 등에 작은 모돌 일러스트
- **문서·랜딩 일러스트** — 각 섹션 옆에 작은 액션 모돌
- **에러 화면** — 공감 톤으로 "모돌도 막혔어요"

### 차별화 포인트
- 어느 OSS도 미니비숑+꿀벌 마스코트 없음 → *최강 차별화*
- 진정성 — 개인 프로젝트, 반려견 이름에서 직접 따옴
- 한국적 정서 — 한국에서 미니비숑 인기 + 한국식 코스튬 감성
- 시각 메모리 강함 — 강아지 + 노란색 = 즉시 기억
- 굿즈 확장성 — 스티커, 티셔츠, 미니어처, AR 필터, 디스코드 이모지

---

## 3. 컬러 팔레트 — Honey

| Token | Hex | 사용처 |
|---|---|---|
| `--hanimo-amber` | `#f5a623` | **primary accent** — 로고 옷, 라이선스 배지, 강조 텍스트, 버튼 |
| `--hanimo-amber-soft` | `#ffd089` | hover/glow, 부드러운 강조 |
| `--hanimo-amber-deep` | `#d99437` | pressed state, 진한 강조 |
| `--hanimo-amber-glow` | `rgba(245,166,35,0.22)` | radial gradient, 부드러운 후광 |
| `--hanimo-cream` | `#f0e4d3` | light bg, 모돌 털 색 |
| `--hanimo-brown` | `#1a1410` | dark bg base (메인) |
| `--hanimo-brown-soft` | `#1f1813` | dark elev (카드, 모달) |
| `--hanimo-good` | `#6cae75` | success — 따뜻한 녹색 (차가운 #10b981 X) |
| `--hanimo-warn` | `#e8a317` | warning |
| `--hanimo-error` | `#d97757` | error — 따뜻한 적색 (차가운 #f87171 X) |
| `--hanimo-info` | `#6b8db5` | info — 차분한 파랑 |

### 메타포의 컬러 매핑
- **모돌 털** = cream `#f0e4d3` (강아지 흰털을 따뜻한 크림으로 표현)
- **꿀벌 옷 노란 부분** = amber `#f5a623`
- **꿀벌 옷 검정 줄무늬** = brown `#1a1410`
- **눈·코·발끝** = brown `#1a1410`

이 4색만으로 마스코트가 완성되어야 하며, 다른 색이 들어가면 그건 *다른 캐릭터*가 된다.

### 거부된 톤
- **차가운 파랑** (`#3b82f6`, `#0ea5e9`) — Apache/Pure-tech 인상이 너무 강함
- **lavender** (`#CBA6F7`) — v0 legacy에서 사용. 모돌 정체성과 어울리지 않음
- **네온 핑크/그린** — 한국 정서와 충돌

---

## 4. 타이포그래피

| 용도 | Stack | 비고 |
|---|---|---|
| UI / 본문 | `'Pretendard', 'Noto Sans KR', system-ui, -apple-system, sans-serif` | 한국어 + 영문 모두 강함, 가변 폰트 |
| 코드 / 모노 | `'JetBrains Mono', 'SF Mono', 'Menlo', 'Consolas', monospace` | hanimo 코드 톤 |
| 디스플레이 | UI stack에 `font-weight: 800` + `letter-spacing: -2.5px` | 헤로 헤더 |

### 폰트 호스팅
- **CDN 금지** — `telemetry-and-privacy.md` §4 호환 (외부 호출 0)
- **로컬 번들** — `public/fonts/Pretendard-*.woff2` 형태로 각 repo에 직접 번들
- **시스템 폰트 폴백** — 폰트 다운로드 실패 시에도 안정적 렌더링

### Pretendard 라이선스
- SIL Open Font License 1.1 — Apache 2.0과 호환
- 번들 시 `LICENSE-Pretendard.txt` 함께 배포

---

## 5. 모티프 — Hexagon (배경 패턴)

미니비숑·꿀벌이 *메인*이고, hexagon은 *보조 배경*으로 강등.

### 사용 규칙
- **메인 비주얼에는 hexagon 단독 사용 금지** — 항상 마스코트와 함께
- **배경 패턴**: 매우 낮은 opacity (5~10%), 큰 사이즈, 화면 가득 채우는 배경
- **카드 그리드**: hexagon 모양의 카드는 *3개 이상 모일 때만* 허용 (1개 단독 금지)
- **divider**: 작은 hexagon 점 3개 (· · ·)로 변형 가능

### 의미
- 벌집 셀 = hanimo 5 repo가 한 hive
- 모돌이 그 벌집의 마스터 — 꿀벌 옷 입고 *벌집을 지키는* 강아지

---

## 6. 메시징 톤

### 4 원칙

1. **차분 + 자신감** — "혁명적", "최고의", "유일한" 같은 마케팅 fluff 금지. 사실로만 말함.
2. **한국어 친화** — 한국어 본문 + 정확한 영문 병기. "오픈소스" "에이전트" 같은 외래어 자연스럽게 사용.
3. **솔직** — 한계와 거부된 대안을 *명시*. STRATEGY 문서의 §"거부된 대안" 패턴 유지.
4. **이모지 절제** — 코딩 컨벤션과 일치. 마스코트가 이미 친근하므로 텍스트는 차분하게.

### 좋은 예
> "hanimo는 단일 ~20MB 바이너리로 배포됩니다. 14+ LLM 프로바이더를 지원하고 외부 통신 없이 동작합니다."

### 나쁜 예
> "🚀 hanimo는 AI 코딩의 미래를 혁신합니다! ✨ 강력하고 놀라운 기능을 경험해보세요! 🎉"

### 마스코트가 직접 말할 때
- 1인칭 "저", "모돌이", 친근하지만 정확
- "모돌이 도와드릴게요" ✅ / "모돌이가 도와줍니다요" ❌

---

## 7. 적용 단위 — 5 repo 매핑

| repo | 컬러 | 타입 | 마스코트 노출 | 비고 |
|---|---|---|---|---|
| **hanimo-code** | TUI 5 테마 + 데스크톱 8 테마 모두 amber | UI: ko 사용 시 Pretendard 권장 / Mono: JetBrains Mono | favicon, README hero | 가장 성숙 |
| **hanimo-webui** | Tailwind config의 primary를 amber로 | Pretendard + JetBrains Mono | hero, empty state, error | DESIGN_SYSTEM.md와 정합 |
| **hanimo-rag** | dashboard CSS 변수 amber | Pretendard + JetBrains Mono | dashboard hero, README | 별도 dashboard 있음 |
| **hanimo-community** | 랜딩 페이지 amber 메인 | Pretendard 기반 | 페이지 메인 비주얼 | 가장 마스코트 노출 큼 |
| **hanimo-toy-coffee** | (적용 안 함) | (적용 안 함) | (적용 안 함) | 내부 샌드박스 |

---

## 8. 거부된 컨셉 (이력 보존)

| 폐기 컨셉 | 시점 | 이유 |
|---|---|---|
| Catppuccin Mocha 베이스 + lavender accent | v0 (2026-04 이전) | hanimo의 honey 정체성과 충돌, lavender는 hanimo 어원과 무관 |
| 단순 꿀벌 마스코트 | 2026-05-03 직전 | 차별화 약함, hanimo의 'modol' 어원이 시각화되지 않음 |
| Hexagon 단독 메인 비주얼 | 2026-05-03 | 인지·기억성 약함, 다른 OSS와 차별화 부족 |
| 다크 블루 톤 | 2026-04 검토 | "쿨하지만 차가움", honey의 따뜻함과 충돌 |

---

## 9. 다음 단계 로드맵

| 단계 | 작업 | 예상 시간 | 의존 |
|---|---|---|---|
| **Phase 1** ✅ | Honey 라이선스 배지 5 repo 통일 + Lucide 컨벤션 박기 | (완료) | - |
| **Phase 1** ✅ | Brand docs 7종 신규 작성 | (완료) | - |
| **Phase 2** | `tokens.css` 5 repo import 적용 | 반나절 | - |
| **Phase 2** | 4 repo README hero 통일 (텍스트만, 로고 placeholder) | repo당 30분 | - |
| **Phase 3** | 모돌-꿀벌 SVG 로고 v1 (정사각 + 가로형) | 2~3h | 디자이너 또는 AI 일러스트 |
| **Phase 3** | favicon · GitHub org avatar 적용 | 1h | 로고 v1 |
| **Phase 4** | `@hanimo/ui` 컴포넌트 라이브러리 (shadcn 기반) | 며칠 | - |
| **Phase 4** | hanimo-webui Phosphor → Lucide 도메인 단위 치환 | 도메인당 1~2h | - |
| **Phase 5** | hanimo.dev 통합 랜딩 페이지 (모돌 메인) | 1~2일 | 도메인 + 로고 |
| **Phase 5** | Pretendard 폰트 5 repo 번들 | repo당 15분 | - |

---

## 10. 의사결정 기록

본 컨셉이 합의된 시점: **2026-05-03**.

채택 근거:
1. *기능적 어원*(harness)과 *정서적 어원*(modol)을 한 마스코트에 통합
2. Honey amber `#f5a623`는 hanimo-code-desktop 디자인 v1에 *이미* 박혀있던 자산 — 폐기 대신 확장
3. 차별화 — 어느 OSS도 미니비숑+꿀벌 마스코트 없음
4. 한국 정서 + OSS 명성 wedge 동시 강화
5. 5-repo 시각 통일을 위한 *구체적 anchor* 제공 (추상 컨셉이 아닌 구체적 캐릭터)

거부 근거 (위 §8 참고):
- 단순 꿀벌, 단순 hexagon, lavender, 다크 블루 — 모두 정체성 명료성 부족

본 결정은 *되돌릴 수 없는* 정체성 결정이며, 변경 시 5 repo 전체 리브랜딩 발생.
변경하려면 새 ADR 문서 + 5 repo 동시 조율 필요.

---

_Last updated: 2026-05-03 · brand v1 — Modol the Honey-Bee Bichon_
