# hanimo Design Tokens — Reference

> [`tokens.css`](tokens.css)의 모든 변수를 카테고리별로 설명. 변수명·값·
> 사용처·메타포 매핑을 한 표에 정리.

## 색상 — Honey palette

| Token | Hex | 메타포 | 사용처 |
|---|---|---|---|
| `--hanimo-amber` | `#f5a623` | 꿀벌 옷 노란 부분 | primary accent · 라이선스 배지 · CTA 버튼 · 강조 텍스트 |
| `--hanimo-amber-soft` | `#ffd089` | 꿀벌 옷 약한 부분 | hover state · subtle highlight |
| `--hanimo-amber-deep` | `#d99437` | 꿀벌 옷 그림자 부분 | pressed state · 진한 강조 |
| `--hanimo-amber-glow` | `rgba(245,166,35,0.22)` | 꿀의 빛 발산 | radial gradient · selection · focus ring |

## 색상 — Cream / Brown

| Token | Hex | 메타포 | 사용처 |
|---|---|---|---|
| `--hanimo-cream` | `#f0e4d3` | 모돌의 흰털 (따뜻한 크림) | light bg · 카드 |
| `--hanimo-cream-soft` | `#f7eedf` | 모돌 털의 더 밝은 부분 | light elev |
| `--hanimo-brown` | `#1a1410` | 어두운 배경 (벌집 그늘) | dark bg base |
| `--hanimo-brown-soft` | `#1f1813` | 카드/모달 배경 | dark elev |
| `--hanimo-brown-deeper` | `#14100d` | 가장 깊은 배경 | activity bar / sidebar |

## 색상 — Semantic (의미론적, 모두 따뜻한 톤)

| Token | Hex | 의미 | 비고 |
|---|---|---|---|
| `--hanimo-good` | `#6cae75` | 성공·확인 | 차가운 `#10b981` 거부 — Honey 톤과 어울리지 않음 |
| `--hanimo-warn` | `#e8a317` | 경고 | amber와 가까운 톤 |
| `--hanimo-error` | `#d97757` | 에러 | 차가운 `#f87171` 거부 — 따뜻한 적색 |
| `--hanimo-info` | `#6b8db5` | 정보 | 차분한 파랑, 차가움 강도 낮음 |

각 semantic 색상은 `*-soft` variant도 제공 (15% 투명도) — 배경·테두리·glow에 사용.

## 색상 — Foreground

| Token | Hex | 사용처 |
|---|---|---|
| `--hanimo-fg` | `#e8e8ec` | primary text |
| `--hanimo-fg-dim` | `#a0a0a8` | secondary text · 메타데이터 |
| `--hanimo-fg-faint` | `#5a5a64` | placeholder · 매우 약한 힌트 |

## 색상 — Surface

| Token | 사용처 |
|---|---|
| `--hanimo-bg` | 페이지 배경 (다크 모드 기본) |
| `--hanimo-bg-elev` | 카드 배경 (배경보다 약간 밝음) |
| `--hanimo-bg-card` | 강조 카드 |
| `--hanimo-bg-code` | 코드 블록 배경 |
| `--hanimo-border` | 일반 경계 |
| `--hanimo-border-strong` | 강조 경계 (카드 호버 시 등) |
| `--hanimo-border-focus` | 포커스 링 — `--hanimo-amber` 참조 |

## 타이포그래피

| Token | 값 | 용도 |
|---|---|---|
| `--hanimo-font-ui` | `'Pretendard', 'Noto Sans KR', system-ui, ...` | UI / 본문 (한·영 모두 강함) |
| `--hanimo-font-mono` | `'JetBrains Mono', 'SF Mono', 'Menlo', ...` | 코드 / 모노 |

**중요**: 폰트는 **로컬 번들**해야 함. CDN 사용 금지 — `telemetry-and-privacy.md` §4 호환.

## Radius

| Token | px | 용도 |
|---|---|---|
| `--hanimo-radius-xs` | 4 | 매우 작은 chip · pill 일부 |
| `--hanimo-radius-sm` | 6 | 작은 버튼 · input |
| `--hanimo-radius` | 10 | 일반 (기본값) |
| `--hanimo-radius-lg` | 14 | 카드 |
| `--hanimo-radius-xl` | 20 | 모달 · 큰 패널 |
| `--hanimo-radius-pill` | 999 | pill 형태 (badge, chip) |

## Shadow

| Token | 사용처 |
|---|---|
| `--hanimo-shadow-sm` | 약한 elev (작은 카드) |
| `--hanimo-shadow` | 중간 elev (드롭다운, 일반 카드) |
| `--hanimo-shadow-lg` | 강한 elev (모달, 알림) |
| `--hanimo-shadow-glow` | amber 후광 (CTA 강조) |

## Spacing — 4 단위 스케일

| Token | px | 일반 사용처 |
|---|---|---|
| `--hanimo-sp-1` | 4 | tight grid gap |
| `--hanimo-sp-2` | 8 | small padding |
| `--hanimo-sp-3` | 12 | input padding |
| `--hanimo-sp-4` | 16 | 일반 padding (기본) |
| `--hanimo-sp-5` | 20 | 카드 내부 |
| `--hanimo-sp-6` | 24 | 섹션 간 여백 |
| `--hanimo-sp-8` | 32 | 큰 섹션 |
| `--hanimo-sp-10` | 40 | hero padding |
| `--hanimo-sp-12` | 48 | 큰 hero |
| `--hanimo-sp-16` | 64 | 페이지 상하 여백 |

## Transition

| Token | 값 | 용도 |
|---|---|---|
| `--hanimo-transition-fast` | `100ms ease-out` | hover 색상 변화 등 빠른 피드백 |
| `--hanimo-transition` | `180ms ease-out` | 일반 (기본값) |
| `--hanimo-transition-slow` | `300ms ease-out` | 모달 열림, 큰 영역 변화 |

## Z-index 스케일

| Token | 값 | 용도 |
|---|---|---|
| `--hanimo-z-dropdown` | 100 | 드롭다운 메뉴 |
| `--hanimo-z-sticky` | 200 | sticky 헤더 |
| `--hanimo-z-overlay` | 300 | 백드롭 |
| `--hanimo-z-modal` | 400 | 모달 |
| `--hanimo-z-popover` | 500 | 팝오버 |
| `--hanimo-z-tooltip` | 600 | 툴팁 |
| `--hanimo-z-toast` | 700 | 토스트 (최상단) |

## 사용 예시

### 1. 일반 HTML/CSS

```html
<link rel="stylesheet" href="path/to/tokens.css">

<style>
  .button {
    background: var(--hanimo-amber);
    color: var(--hanimo-brown);
    padding: var(--hanimo-sp-3) var(--hanimo-sp-5);
    border-radius: var(--hanimo-radius);
    transition: background var(--hanimo-transition);
  }
  .button:hover { background: var(--hanimo-amber-soft); }
  .button:active { background: var(--hanimo-amber-deep); }
</style>
```

### 2. Tailwind 매핑

`tailwind.config.js`:

```js
module.exports = {
  theme: {
    extend: {
      colors: {
        'hanimo-amber':         '#f5a623',
        'hanimo-amber-soft':    '#ffd089',
        'hanimo-amber-deep':    '#d99437',
        'hanimo-cream':         '#f0e4d3',
        'hanimo-brown':         '#1a1410',
        'hanimo-good':          '#6cae75',
        'hanimo-warn':          '#e8a317',
        'hanimo-error':         '#d97757',
        'hanimo-info':          '#6b8db5',
      },
      fontFamily: {
        sans: ['Pretendard', 'Noto Sans KR', 'system-ui'],
        mono: ['JetBrains Mono', 'SF Mono', 'Menlo'],
      },
      borderRadius: {
        'hanimo': '10px',
        'hanimo-lg': '14px',
      },
    },
  },
}
```

### 3. React inline (변수 직접 참조)

```jsx
<button style={{
  background: 'var(--hanimo-amber)',
  color: 'var(--hanimo-brown)',
  padding: 'var(--hanimo-sp-3) var(--hanimo-sp-5)',
  borderRadius: 'var(--hanimo-radius)',
}}>
  click me
</button>
```

### 4. 다크 → 라이트 토글

```html
<!-- 기본 (다크) -->
<html>...</html>

<!-- 라이트 모드 -->
<html data-theme="light">...</html>
```

## 거부된 토큰

| 토큰 후보 | 거부 사유 |
|---|---|
| `--hanimo-blue: #3b82f6` | 차가운 톤, honey 정체성과 충돌 |
| `--hanimo-purple: #cba6f7` | v0 legacy의 lavender 잔재 |
| `--hanimo-pink: #ec4899` | 한국 정서와 충돌, 마스코트 정체성 흐림 |
| 채도 100% 네온 색상군 | "쿨하지만 차가움" — 따뜻함이 hanimo의 핵심 |

## 변경 절차

토큰 값을 바꾸려면:
1. PR 필수 (자가 머지 금지)
2. 5 repo 영향 분석 (BRAND.md §7 참고)
3. `BRAND.md`의 §3·§4와 정합 유지 (메타포 깨지면 안 됨)
4. CHANGELOG에 *어떤 토큰이 어떻게 변경됐는지* 명시
