# hanimo Logo Specification — Modol the Honey-Bee Bichon

> 마스코트 로고를 어떻게 그려야 하는지의 사양. 디자이너에게 직접
> 전달하거나, 일러스트 AI에게 프롬프트로 사용 가능.

## 1. 핵심 메타포

> **모돌(미니비숑 강아지)이 꿀벌 옷을 입은 형태.**

자세한 컨셉 배경: [`BRAND.md`](BRAND.md) §2

## 2. 필수 구성 요소

로고 한 컷에 *반드시* 포함되어야 하는 요소:

| 요소 | 필수도 | 설명 |
|---|---|---|
| 미니비숑 얼굴 | ★★★ | 흰털·둥근 얼굴·검은 눈코·살짝 늘어진 귀 |
| 꿀벌 옷 | ★★★ | Honey amber `#f5a623` 몸통, 검정 줄무늬 (가로 2~3줄) |
| 작은 더듬이 | ★★ | 얼굴 위쪽에 두 가닥, 끝에 작은 점 |
| 작은 날개 | ★ | 옷 등 부분에서 살짝 보이는 투명/반투명 날개 |
| 미소 / 차분 표정 | ★★ | 입꼬리 살짝 올라간 정도, 과도한 웃음 금지 |

## 3. 색상 (정확히)

```
모돌 털:        #f0e4d3   (--hanimo-cream)
꿀벌 옷:        #f5a623   (--hanimo-amber)
검정 줄무늬:    #1a1410   (--hanimo-brown)
눈·코·발끝:     #1a1410   (--hanimo-brown)
하이라이트:     #ffd089   (--hanimo-amber-soft) — 옷의 광택
```

다른 색은 *허용되지 않음*. 4색 + 1 하이라이트만으로 완성되어야 함.

## 4. 변형 (variants)

3 variant 작성 필수:

| Variant | 사이즈 | 용도 |
|---|---|---|
| **Square (정사각)** | 256×256 / 512×512 SVG | favicon · GitHub avatar · 앱 아이콘 |
| **Horizontal (가로형)** | 600×200 SVG | README hero · 랜딩 헤더 |
| **Compact (작은 정사각)** | 64×64 / 32×32 SVG | favicon 16/32 · 작은 UI 아이콘 |

각 variant는 단순화 정도가 다름:
- **Square**: 풀 일러스트 (얼굴 + 옷 + 더듬이 + 날개)
- **Horizontal**: 얼굴+옷만 + "hanimo" 텍스트
- **Compact**: 얼굴 단순 실루엣 (꿀벌 줄무늬는 1~2개로 축소)

## 5. 스타일 가이드

### 권장
- ✅ 플랫 일러스트 (그라디언트 최소, 그림자 없음)
- ✅ 둥글고 부드러운 윤곽선
- ✅ 단순화 — 디테일은 적을수록 좋음
- ✅ 친근함 + 진지함 균형
- ✅ 정면 또는 약간 사선 정면 (3/4 view)

### 거부
- ❌ 사실적 사진풍 일러스트
- ❌ 과도한 카툰화 (눈이 너무 큰 디즈니풍, 챱쌀풍)
- ❌ 그라디언트로 음영 표현 (플랫 유지)
- ❌ 5색 이상 사용
- ❌ 텍스트가 캐릭터를 가리는 구도
- ❌ 과한 표정 (혀 내밀기, 눈 감기 등) — 차분 + 자신감

## 6. 일러스트 AI 프롬프트 (옵션)

Midjourney / DALL·E / Stable Diffusion 등에 그대로 입력 가능:

```
A flat illustration mascot of a Bichon Frisé puppy wearing a bee
costume. The dog has fluffy cream-colored fur (#f0e4d3), round face,
black eyes and nose. The bee costume is honey amber (#f5a623) with
two horizontal black stripes (#1a1410). Two small antennae on top
of the head with tiny dots at the ends. Small translucent wings
peeking from the back. Calm, confident, slightly smiling expression.
Three-quarter front view. Flat design, minimal detail, no gradients,
no shadows, soft rounded outlines. Clean white background. Logo
style, suitable for app icon and avatar. 5 colors only:
cream, amber, brown, black, soft amber highlight.
```

## 7. SVG 작성 시 주의

- **viewBox** 지정 필수 (반응형)
- **Path 단순화** — 중첩 path 5개 이내 권장
- **CSS variable 호환** — 색상은 `currentColor` 또는 직접 hex
- **size-agnostic** — 16×16 viewport에서도 인식 가능
- **단일 파일** — 외부 폰트·이미지 참조 금지 (telemetry-and-privacy 정책)

### SVG 예시 구조

```xml
<!-- hanimo-logo-square.svg -->
<svg xmlns="http://www.w3.org/2000/svg"
     viewBox="0 0 256 256"
     fill="none">

  <!-- background hex (옵션) -->
  <polygon points="..." fill="#1a1410" opacity="0.05"/>

  <!-- bichon body (옷 입은 부분) -->
  <ellipse cx="..." cy="..." rx="..." ry="..."
           fill="#f5a623"/>
  <path d="..." stroke="#1a1410" stroke-width="8"/>  <!-- 줄무늬 -->

  <!-- bichon head -->
  <circle cx="..." cy="..." r="..." fill="#f0e4d3"/>

  <!-- ears -->
  <path d="..." fill="#f0e4d3"/>

  <!-- eyes -->
  <circle cx="..." cy="..." r="3" fill="#1a1410"/>
  <circle cx="..." cy="..." r="3" fill="#1a1410"/>

  <!-- nose -->
  <ellipse cx="..." cy="..." rx="4" ry="3" fill="#1a1410"/>

  <!-- antennae -->
  <line x1="..." y1="..." x2="..." y2="..."
        stroke="#1a1410" stroke-width="2"/>
  <circle cx="..." cy="..." r="3" fill="#1a1410"/>

  <!-- wings (반투명) -->
  <ellipse cx="..." cy="..." rx="..." ry="..."
           fill="#ffd089" opacity="0.6"/>

</svg>
```

## 8. 라이선스

로고 자체의 저작권은 작성자(Kim Jiwon)에게 귀속. AI 일러스트로 시작하더라도
*수정·재해석*을 거쳐 본인이 *명확한 저작자*가 되도록 함.

## 9. 작업 산출물

다음 파일들이 작성되면 본 spec 충족:

```
docs/branding/assets/
├── hanimo-logo-square.svg          (256×256)
├── hanimo-logo-square-512.svg      (512×512, 더 디테일한 버전)
├── hanimo-logo-horizontal.svg      (600×200, 텍스트 포함)
├── hanimo-logo-compact.svg         (64×64, 단순화)
├── hanimo-logo-favicon.svg         (32×32, 가장 단순)
├── hanimo-logo-favicon.png         (32×32 PNG, ICO 변환용)
└── hanimo-logo.ico                 (favicon용 ICO bundle)
```

## 10. 적용 체크리스트

로고 v1이 완성되면:

- [ ] 4 repo의 README hero에 horizontal variant 임베드
- [ ] 각 repo의 favicon (있는 경우) 교체
- [ ] hanimo-code-desktop 앱 아이콘 (.icns / .ico) 갱신
- [ ] GitHub organization avatar 갱신 (확보 후)
- [ ] hanimo.dev 랜딩 페이지 메인 비주얼 (도메인 확보 후)
- [ ] 디스코드 이모지 :modol: (커뮤니티 개설 시)
