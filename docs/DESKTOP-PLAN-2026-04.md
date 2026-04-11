# hanimo Desktop — Plan & Design Direction (2026-04)

> **작성일**: 2026-04-11
> **상태**: 플랜 전용. 구현 없음.
> **범위**: Cursor/Zed 스타일 "폴더 탐색 + 코드 에디터 + 내장 hanimo 에이전트"가 **한 개의 데스크탑 앱**으로 결합된 단일 바이너리.
> **디자인 방향**: 흔한 아이콘/테마 금지. 꿀벌·honeycomb 메타포 기반 오리지널 SVG.
> **관계 문서**:
> - `docs/MASTER-OVERVIEW-2026-04.md` — CLI 전체 맥락
> - `docs/PLATFORM-PLAN-2026-04.md` — 5-surface 플랫폼 전략 (WebUI/Community 포함)
> - `docs/ecosystem/ROADMAP.md` — hanimo 생태계 로드맵

---

## 0. TL;DR (30초)

- **무엇**: `hanimo desktop` — 파일 트리 + 에디터 + 내장 AI 패널을 한 앱으로. 기존 `hanimo code` CLI와 동일 엔진 공유.
- **스택**: Wails v3 (Go + webview) + Svelte 5 + Monaco. Electron 회피, 번들 30~60MB 목표.
- **차별점**: VSCode Codicon / Material / Heroicons 일체 미사용. 모든 아이콘·로딩·일러스트 오리지널 SVG. 꿀벌·honeycomb 기반 디자인 언어.
- **관계**: `hanimo code`(CLI)와 **같은 Go 코어** 직접 임포트. CLI를 버리는 게 아니라 GUI 옵션을 추가.
- **타깃**: CLI가 부담스러운 비-터미널 유저 + sealed 빌드를 배포받는 기업 유저.
- **MVP**: 7~8주. 주력 차별화는 **디자인 시스템과 UX**, 기능은 일단 hanimo code와 동등.

---

## 1. 왜 만드는가

### 1.1 기존 Cursor/Zed가 커버하지 못하는 구멍

| 툴 | 한계 |
|---|---|
| Cursor | VSCode fork라 무거움, 디자인 "흔함", 오픈소스 아님 |
| Zed | 빠르지만 Anthropic/OpenAI 중심, 로컬 모델 fastpath 약함 |
| Claude Code (desktop) | Anthropic 종속, 코드 에디터 내장 없음 |
| VSCode + Copilot | Copilot lock-in, 로컬 모델 지원 부족 |

hanimo desktop의 기회:
1. **로컬 모델 1급** — Ollama/vLLM/LM Studio가 1류 시민.
2. **한국어 친화** — intent 감지·clarify-first·UI 한글.
3. **단일 바이너리** — 설치 = 파일 하나. 업데이트 = 바이너리 교체.
4. **baked 빌드** — 관리자가 엔드포인트/키 고정한 빌드를 배포 가능(CLI와 공유).
5. **오리지널 디자인** — 브랜딩이 식별 가능한 차별화.

### 1.2 명시적 비-타깃 (오버스코프 방지)

- ❌ VSCode 확장 생태계 흡수
- ❌ 거대 IDE 기능 (디버거, 리팩터링, IntelliSense 풀스택)
- ❌ 브라우저 협업 (그건 `hanimo webui` 몫)
- ❌ Mobile
- ❌ AI 없는 코드 에디터로서의 경쟁

---

## 2. 아키텍처

### 2.1 3-layer 구조

```
┌────────────────────────────────────────────────────────────┐
│  [Frontend]  Svelte 5 + TypeScript                         │
│              Monaco Editor (lazy-loaded)                   │
│              xterm.js (터미널 패널)                         │
│              custom SVG 아이콘 · 마스코트 · 애니메이션        │
└───────────────────────┬────────────────────────────────────┘
                        │ Wails v3 bindings (struct ↔ TS)
┌───────────────────────▼────────────────────────────────────┐
│  [Bridge]    Go — 얇은 glue                                 │
│              FS watcher · project load · tab state          │
│              LLM stream pump (channel → Wails Events)       │
└───────────────────────┬────────────────────────────────────┘
                        │ Go import (same process)
┌───────────────────────▼────────────────────────────────────┐
│  [Core]      hanimo/internal/* 기존 코어 재사용             │
│              agents · tools · llm · session · mcp · ...     │
│              수정 없이 함수 호출                             │
└────────────────────────────────────────────────────────────┘
```

**핵심**: 모든 FS/LLM 작업은 Go에서. Webview는 read-only view + 사용자 입력만.

### 2.2 프레임워크 비교 & 선택

| 후보 | 번들 크기 | 스타트업 | hanimo 코어 재사용 | 결정 |
|---|---|---|---|---|
| **Wails v3** | 30~60MB | <300ms | Go 직접 import | ✅ **Primary** |
| Tauri v2 | 25~40MB | <300ms | Go sidecar 필요 (FFI/stdio) | 🟡 Fallback |
| Electron | 150~400MB | 1~3s | sidecar | ❌ |
| Neutralino | 2~5MB | <200ms | 런타임 webview 의존 | ❌ |

**결정**: Wails v3.
**Fallback**: v3 alpha가 깨지면 Wails v2 → 마지막 수단 Tauri.

### 2.3 프로세스 모델

- **단일 프로세스**. Frontend webview + Go backend + hanimo core 전부 한 프로세스.
- 긴 작업(LLM stream, tool call)은 goroutine, 진행률은 채널 → `EventsEmit`.
- xterm pty는 별도 goroutine. PTY는 `creack/pty` 사용.

---

## 3. 디자인 방향 (이 문서의 핵심)

> **원칙**: 흔하면 버린다. 꿀벌이 산다. 기하학적 명료함 + 손맛.

### 3.1 메타포

모든 비주얼을 **벌집(hive)** 에서 파생:

| 개념 | 메타포 | 비주얼 |
|---|---|---|
| 파일 | 벌집 셀(cell) | 육각형 아이콘 |
| 폴더 | 셀들의 군집 | 겹친 육각형 덩어리 |
| 프로젝트 | 벌집 전체 (hive) | 큰 단면도 |
| 저장 | 꿀을 채운다 | 셀이 점차 꿀로 차오르는 애니메이션 |
| 에이전트 작업 | 꿀벌 비행 | 2D flight path (점선 궤적) |
| 에러 | 꿀벌 침 | 작은 삼각 sting 아이콘 |
| 로딩 | 꿀 흐름 | 수직 꿀 스트림 (SVG gradient animate) |
| 실행 모드 | 벌 종류 | Super=Worker, Deep Agent=Queen, Plan=Scout |

### 3.2 컬러 팔레트

```css
/* Core */
--honey-gold:     #F9B640;  /* primary brand */
--deep-amber:     #7A4A12;  /* hover, brand dark */
--cream:          #F5EFE3;  /* light surface */
--ink-black:      #161413;  /* dark surface / text */

/* Semantic */
--nectar:         #E8A427;  /* success — 꿀 */
--sting:          #C73E2F;  /* error — 침 색깔 */
--pollen:         #F2D58E;  /* warning — 꽃가루 */
--comb-line:      #C9A970;  /* divider — 벌집 테두리 */

/* Mono text */
--text-strong:    #1B1916;
--text-mid:       #4F4942;
--text-weak:      #8B7F6E;
```

**금지**: 고채도 파랑(#0066FF 류), 고채도 녹색 성공(#00C853 류), 네온 계열. "Material / Bootstrap / Electron 기본값" 냄새 제거.

### 3.3 타이포 스케일

```css
--font-display: 'Instrument Serif', 'Noto Serif KR', serif;   /* 제목 전용 */
--font-body:    'Pretendard Variable', -apple-system, sans-serif;
--font-mono:    'JetBrains Mono', 'D2 Coding', monospace;

--text-xs:  11px / 1.45;
--text-sm:  13px / 1.5;
--text-md:  15px / 1.55;  /* 본문 */
--text-lg:  18px / 1.4;   /* display H3 */
--text-xl:  24px / 1.25;  /* display H2 */
--text-2xl: 32px / 1.15;  /* display H1 */
```

**왜 세리프 + 모노**: 일반 SaaS는 거의 전부 Inter/SF Pro/Geist sans. Instrument Serif + 한글 세리프 조합은 즉시 "다른 것"처럼 읽힌다.

### 3.4 레이아웃 초안

```
┌──────────────────────────────────────────────────────────────────┐
│ ≡  hanimo   /Users/kim/my-app                  🐝  ▢  ▭  ✕      │  custom titlebar
├─────────┬─────────────────────────────────────┬──────────────────┤
│ HIVE    │                                     │  HANIMO PANEL    │
│         │             EDITOR                  │                  │
│ ⬡ cmd   │                                     │  > 리팩터 해줘    │
│ ⬡ src   │   // main.go                        │                  │
│   ⬡ api │   func main() {                     │  >> file_read    │
│   ⬡ ui  │     fmt.Println("hi")               │  >> file_edit    │
│ ⬡ docs  │   }                                 │  ✓ 적용           │
│         │                                     │                  │
│ 🐝 봄    │  [ main.go ●] [api.go]              │  [Super] [Deep]  │
│ 🐝 채집  │                                     │  [Plan]          │
│ 🐝 설정  │  ───── TERMINAL ─────               │                  │
│         │   $ go test ./...                   │                  │
│         │   PASS                              │                  │
├─────────┴─────────────────────────────────────┴──────────────────┤
│ main • modified    🐝 수신중 2.3s · 45tok · 19t/s                 │  status bar
└──────────────────────────────────────────────────────────────────┘
```

**차별점**:
- 파일 트리 아이콘 = **육각형 셀** (접힌 셀/열린 셀로 폴더/파일 구분)
- 들여쓰기 라인 = 벌집 셀 테두리 스타일 (얇은 꺾인 선)
- 탭 active 마커 = 점이 아닌 꿀방울(•) 형태
- 상태바의 spinner = 꿀 흐름 바

### 3.5 SVG 아이콘 세트 (사이드바)

표준 VSCode는 Files/Search/Git/Debug/Extensions. 우리는 **전혀 다른 명명 + 전혀 다른 그래픽**:

| 슬롯 | 이름 | 역할 | SVG 컨셉 | 제작 노트 |
|---|---|---|---|---|
| 1 | **Hive** | 탐색기 | 벌집 단면, 한 셀 강조 | 24×24, 1.5px stroke, honey-gold fill |
| 2 | **Scent** | 검색 | 꿀벌이 남긴 향기의 나선 | 나선 path + 끝점 dot |
| 3 | **Waggle** | git | 꿀벌의 8자 춤 궤적 | lemniscate + 미세 점선 |
| 4 | **Pollen** | memory (remember/memories) | 3점 꽃가루 + 연결 호 | 3개 원 + arc |
| 5 | **Flightpath** | auto/deep history | 위에서 본 비행 궤적 | dashed curved path |
| 6 | **Nectar** | sessions | 꿀방울 + 미세 반사광 | drop + highlight ellipse |
| 7 | **Drone** | MCP servers | 기하학 드론 벌 | hex body + 2 wings |
| 8 | **Queen** | 설정 | 왕관 쓴 벌 | bee silhouette + crown |

**제작 지침**:
- 모두 24×24 viewBox
- 1.5px uniform stroke
- `currentColor` 사용 → 테마 전환 시 자동
- 호버 애니메이션 = 0.2s cubic-bezier(.34,1.56,.64,1) — 약간의 bounce
- 활성 아이콘은 0.5° 기울기 (벌이 날아가는 느낌)

### 3.6 로고·마스코트

#### 3.6.1 앱 로고

```
파일: docs/images/hanimo-desktop-logo.svg
컨셉: 정육각형 안에 앉은 꿀벌 실루엣.
     벌의 가로 줄무늬 3개가 미묘하게 'h' 모양으로 배열.
컬러: honey-gold + ink-black 2-tone
사이즈: 16 · 22 · 32 · 64 · 128 · 512 (macOS ICNS용)
변형: outline / filled / mono (타이틀바용)
```

#### 3.6.2 마스코트 일러스트 3종

| 이름 | 쓰임 | 컨셉 |
|---|---|---|
| **hani idle** | 온보딩 | 셀 위에 앉아 한쪽 날개 접은 벌 |
| **hani work** | 로딩 | 벌이 여러 셀 사이를 점선 궤적으로 날아가는 모습 |
| **hani sleep** | 빈 상태(empty state) | 셀 하나 위에 몸을 웅크린 벌, 미세 Z |

스타일: **핸드드로잉 + 기하학 혼합**. 완전 카툰 아님, 완전 기하학 아님. 참고는 Excalidraw 정밀도에 Saul Bass 미니멀 기하학.

#### 3.6.3 텍스처

- 전 화면 배경에 **SVG turbulence grain** 오버레이 (`<filter id="grain"><feTurbulence baseFrequency="0.9" numOctaves="2"/></filter>`)
- 투명도 3~5% — 평면이 완전히 평평하지 않게
- 다크 모드에서 조금 더 강하게 (5~8%)

### 3.7 모션 가이드

```css
--ease-bounce:  cubic-bezier(0.34, 1.56, 0.64, 1);  /* UI 전환 */
--ease-soft:    cubic-bezier(0.16, 1, 0.3, 1);      /* 패널 슬라이드 */
--ease-drop:    cubic-bezier(0.4, 0, 0.2, 1);       /* 모달 드롭 */

--duration-snap: 120ms;
--duration-ui:   180ms;
--duration-panel: 260ms;
```

**금지**: `linear`, `ease`, `ease-in-out` 디폴트. 즉시 "흔한" 느낌 나옴.

**Sign-off 시그니처 애니메이션**: 파일 저장 성공 시 "꿀방울이 파일 탭 위로 떨어지고 tab이 미세하게 bounce" — 이 0.4초 애니메이션이 제품의 시그니처.

### 3.8 흔한 → hanimo 대체표

| 흔한 UI | hanimo desktop |
|---|---|
| 파란 체크마크 ✓ | 꿀방울이 떨어지는 0.4s 애니메이션 |
| 빨간 에러 뱃지 | sting 작은 삼각형 + 미세 진동 |
| 녹색 체크 | 육각 셀 하나가 꿀로 채워짐 |
| 원형 spinner | 세로 꿀 스트림 (gradient animate) |
| 알림 뱃지 | 꽃잎 한 장 (ellipse rotate) |
| Progress bar 가로 | 세로 꿀 스트림 + 높이 채움 |
| Modal backdrop | 벌집 패턴 희미한 오버레이 |
| 포커스 링 파란색 | honey-gold 2px + 3px outer glow |
| 컨텍스트 메뉴 | 둥근 육각형 프레임 |
| Dropdown ▼ | 벌집 셀이 접히며 열림 |

### 3.9 다크 모드

- `ink-black` 배경 + `cream` 텍스트
- honey-gold는 그대로 유지(대비 충분)
- grain 강도 5→8%
- 탭 활성 = honey-gold 2px 아래 라인 + 미세 bloom

---

## 4. 기능 범위

### 4.1 MVP (v0.1) — 7~8주

**목표**: "hanimo code의 GUI 버전" + 디자인 시그니처.

- [ ] 폴더 열기 / 최근 프로젝트
- [ ] Hive 파일 트리 (honeycomb 아이콘, collapse)
- [ ] Monaco Editor 내장 (lazy load, 필요 언어만)
- [ ] 파일 CRUD + Ctrl+S (read-before-write 규칙 CLI와 동일)
- [ ] 탭 멀티 파일 + Cmd+W 닫기
- [ ] Hanimo Panel — Super/Deep/Plan 모드 전환
- [ ] 스트리밍 응답 live render (CLI와 동일 엔진)
- [ ] 터미널 패널 (xterm.js + pty) — `Ctrl+J` 토글
- [ ] baked distro / sealed 빌드 (CLI와 Makefile 공유)
- [ ] 커스텀 타이틀바 (darwin/windows)
- [ ] 시그니처 꿀방울 저장 애니메이션
- [ ] 8개 사이드바 아이콘 전부 완성
- [ ] 마스코트 3종 완성

### 4.2 v0.2

- [ ] Git 통합 (status, diff 뷰, stage/commit — CLI의 git 도구 GUI 래핑)
- [ ] Command palette (Cmd+K) — CLI palette 항목 재사용
- [ ] Checkpoint/Rollback 버튼 (CLI `/checkpoint` 시각화)
- [ ] MCP 서버 관리 UI (목록/추가/상태)
- [ ] 세션 브라우저 (Nectar 탭)

### 4.3 v0.3

- [ ] Plan 모드 시각화 — 단계별 plan을 flowchart 보드로
- [ ] Tool call timeline (flightpath 시각화)
- [ ] Diff 승인 UX — GitHub 리뷰 스타일
- [ ] Queen 설정 페이지 (YAML 아닌 form)

### 4.4 v0.4+ (수요 생기면)

- 원격 프로젝트(ssh, devcontainer)
- 경량 자동완성(gopls 통합)
- 테마 마켓

---

## 5. 폴더 구조 (별도 레포)

```
hanimo-desktop/
├── cmd/hanimo-desktop/main.go   # Wails entry
├── internal/
│   ├── bridge/                   # Wails bindings
│   │   ├── fs.go
│   │   ├── project.go
│   │   ├── editor.go
│   │   └── agent.go              # hanimo core로 위임
│   ├── fswatcher/                # fsnotify wrapper
│   ├── pty/                      # creack/pty wrapper
│   └── baked/                    # CLI와 동일 ldflags 시스템
├── frontend/
│   ├── src/
│   │   ├── lib/
│   │   │   ├── hive/             # 파일 트리
│   │   │   ├── editor/           # Monaco wrapper
│   │   │   ├── panel/            # hanimo AI 패널
│   │   │   ├── terminal/         # xterm
│   │   │   ├── mascot/           # SVG mascot components
│   │   │   └── icons/            # 커스텀 아이콘 8개
│   │   ├── routes/
│   │   └── styles/
│   │       ├── tokens.css        # 디자인 토큰 (본 문서 기반)
│   │       └── global.css
│   ├── static/svg/               # 모든 오리지널 에셋
│   └── package.json
├── build/
│   ├── darwin/icons.icns
│   ├── linux/
│   └── windows/
├── docs/
│   └── design-system.md
├── wails.json
├── go.mod                         # replace github.com/flykimjiwon/hanimo → ../hanimo
└── Makefile                       # build / build-distro / build-sealed
```

**replace directive**로 hanimo 코어를 가져옴. CI에서는 태그 pin.

---

## 6. 핵심 기술 결정

### 6.1 에디터 — Monaco vs CodeMirror 6

- **Monaco** 선택: VSCode 사용자 근육 기억 + 신택스 완비.
- 단점: 번들 3MB → lazy-load + 필요 언어만 등록 (Go, TS, Python, Markdown, YAML) → ~1MB.
- CodeMirror 6는 더 가볍지만 언어 설정 수동.

### 6.2 Frontend — Svelte 5 vs Solid vs React

- **Svelte 5** 선택: runes API + 작은 번들 + `<svg>` 친화.
- React 제외: "흔함" + 번들 부담.

### 6.3 터미널 — xterm.js + creack/pty

- Linux/macOS: creack/pty.
- Windows: ConPTY (Go 어댑터 있음).
- fallback: shell exec — PTY 못 쓸 때 dumb mode.

### 6.4 아이콘 렌더링

- 모두 `<svg>` 인라인 (이미지 태그 X) — `currentColor` 활용.
- 번들 시 Svelte 컴포넌트로 import: `import HiveIcon from '$lib/icons/hive.svelte'`.

### 6.5 보안

- Webview CSP 엄격 설정, `--disable-web-security` 절대 금지.
- FS 접근 전부 Go bridge.
- 파괴적 도구 (shell_exec dangerous pattern)는 GUI 모달 확인 (CLI danger confirm 대응).
- Sealed 빌드 ldflags 주입 (CLI와 공유).

### 6.6 업데이트

- `minio/selfupdate` 또는 Wails 자체 업데이트 훅.
- `distro`/`sealed` 빌드는 업데이트 비활성화 옵션.

---

## 7. baked 빌드 일관성

CLI의 `make build-distro` / `make build-sealed` 와 **동일 파라미터 체계**:

```bash
# vanilla
make build-desktop

# distro (엔드포인트만 고정, 사용자는 키만 입력)
make build-desktop-distro \
  ENDPOINT=https://api.novita.ai/v3/openai \
  PROVIDER=novita \
  MODEL=qwen/qwen3-coder-30b

# sealed (키까지 포함, 재배포 금지)
make build-desktop-sealed \
  ENDPOINT=https://api.novita.ai/v3/openai \
  PROVIDER=novita \
  MODEL=qwen/qwen3-coder-30b \
  API_KEY=sk-...
```

내부적으로 `internal/config` (hanimo core)의 `Baked*` ldflags를 그대로 주입 — 별도 시스템 만들지 않음.

---

## 8. 개발 단계 (타임라인)

### Stage 1 — PoC (1주)
- Wails v3 hello world
- hanimo core `replace` 동작
- Svelte 5 + Monaco 통합 1파일 렌더
- `bridge.StartChat("hello")` → frontend에 스트림 출력
- **산출**: 동작하는 단일 화면

### Stage 2 — MVP skeleton (2주)
- Hive 파일 트리 (아이콘은 임시)
- Monaco 멀티탭
- Hanimo Panel (스트리밍 + 모드 전환)
- 터미널 패널
- **산출**: CLI 기능 GUI로 이식 완료

### Stage 3 — 디자인 시스템 (2주)
- 8 사이드바 아이콘 SVG 완성
- 마스코트 3종 완성
- 타이포 / 컬러 토큰 적용
- 시그니처 꿀방울 저장 애니메이션
- 다크/라이트 테마
- **산출**: "다른 앱" 느낌 확보

### Stage 4 — Polish & Release (2주)
- Windows/Linux 빌드
- macOS 공증 (Apple Developer $99)
- baked 빌드 CI 파이프라인
- 랜딩 페이지 섹션 (hanimo.dev/desktop)
- **산출**: v0.1 공개 릴리즈

**합계**: 7주 (혼자 10주 버퍼).

---

## 9. 리스크

| 리스크 | 영향 | 완화 |
|---|---|---|
| Wails v3 alpha 불안정 | 중 | Wails v2 fallback → 최후 Tauri+sidecar |
| Monaco 번들 부담 | 저 | lazy-load + minimal languages |
| **SVG 직접 제작 시간** | 높음 | 1차 스케치는 기하학만 → 2차 iteration에서 일러스트 디테일 |
| macOS 공증 | 저 | Apple 개인 계정 |
| Windows ConPTY 호환 | 중 | Windows 빌드 늦게 (MVP+1) |
| hanimo core 변경 시 desktop 깨짐 | 중 | `go.work` + CI smoke test |
| 사용자가 Cursor/VSCode 근육 기억 요구 | 중 | 단축키는 최대한 호환 (Cmd+P, Cmd+B 등) |

### 미결정
- 코어 레포에 `desktop/` 서브디렉토리 vs 별도 `hanimo-desktop` 레포 → **별도 레포** 권장
- 패널 탈착 가능 여부 → v0.2 결정
- 테마 시스템 외부 공개 여부 → 판단 보류

---

## 10. 디자인 제작 플로우 (SVG 직접 만들기)

유저 요구: "흔한 아이콘·흔한 디자인은 쓰지 않도록 세련되게."

### 10.1 제작 툴 옵션

| 옵션 | 장점 | 단점 |
|---|---|---|
| **Figma** | 빠름, 벡터 정확 | 세리프 한글 미리보기 빈약 |
| **Affinity Designer** | 1회 구매, 강력 | 협업 약함 |
| **Illustrator** | 표준 | 비쌈 |
| **Excalidraw → SVG** | 손맛 유지 | 세밀 조정 어려움 |
| **tldraw → SVG** | 마찬가지 | |
| **Rive** | 애니메이션 | 러닝커브 |

**추천 조합**: Figma (아이콘 8개 + 토큰) + Affinity Designer (마스코트 3종) + Rive (시그니처 애니메이션).

### 10.2 AI 보조 옵션

- **Recraft v3**: 벡터 일러스트 direct 생성 → 수작업 정리
- **Midjourney** → trace: 텍스처/컨셉 → Figma에서 벡터화
- **ControlNet Canny**: 기존 스케치 기반 변형

**주의**: AI가 낸 결과는 **반드시 수작업 정리** — 자잘한 불규칙성이 "흔함"을 만든다.

### 10.3 제작 순서

1. **그리드 정의** — 24×24, 16 주변, 1.5 stroke
2. **Hive 아이콘 먼저** (복잡도 낮음, 기준)
3. **나머지 7개** 같은 언어로
4. **로고**
5. **마스코트 idle**
6. **work** / **sleep**
7. **그레인 패턴** — SVG turbulence 튜닝
8. **애니메이션** — 저장 시그니처부터

---

## 11. 브랜딩 일관성

| 속성 | 값 |
|---|---|
| 제품명 | `hanimo desktop` (lowercase) |
| 워드마크 | `hanimo` + bee 아이콘 조합 |
| 기본 컬러 | `#F9B640` |
| 보조 | `#161413`, `#F5EFE3` |
| 마스코트 | hani (꿀벌), modol (강아지) — 항상 오리지널, 재사용 금지 |
| 슬로건 KR | "둥지에서 코딩하는 작은 에이전트" |
| 슬로건 EN | "A tiny coding agent that lives in your hive" |
| 앱 아이콘 | honeycomb-inset bee, rounded square 64/128/512 |

---

## 12. 배포 & 라이선스

- **라이선스**: MIT (CLI와 동일)
- **채널**:
  - GitHub Releases (공식 바이너리)
  - Homebrew Cask: `brew install --cask hanimo-desktop`
  - Winget: `winget install hanimo`
  - AppImage: Linux
- **서명/공증**:
  - macOS: Apple 공증 v0.2부터 (Developer Program $99/년)
  - Windows: code signing cert (v0.3부터, 비용 시 보류)
- **Sealed 빌드 재배포 금지 경고**: README / build output / 바이너리 기동 시 로그.

---

## 13. 다음 액션 (구체)

- [ ] `hanimo-desktop` 별도 레포 생성 (초기 private)
- [ ] Wails v3 alpha PoC — hanimo core replace 동작 확인
- [ ] Figma 디자인 토큰 파일 초안 (`tokens.css` 대응)
- [ ] 아이콘 8개 1차 스케치 (종이 OK)
- [ ] 마스코트 스타일 결정 — 핸드드로잉 vs 기하학 vs 혼합
- [ ] `docs/DESIGN-SYSTEM.md` 별도 문서 분리 (본 문서에서 §3 섹션 발췌 확장)
- [ ] 1주 PoC 스프린트 착수

---

## 14. 참고 자료

- Wails v3 docs — https://v3alpha.wails.io
- Monaco embedding — https://microsoft.github.io/monaco-editor/
- Svelte 5 runes — https://svelte.dev/blog/runes
- xterm.js — https://xtermjs.org
- creack/pty — https://github.com/creack/pty
- 핸드드로잉+기하학 레퍼런스: Excalidraw, tldraw, Saul Bass posters
- 꿀벌 메타포: Honeycomb.io 로고 시스템, Bees & Bombs (gif 아티스트)
- 시그니처 애니메이션 레퍼런스: Notion의 저장 체크, Linear의 단축키 힌트

---

## 15. 핵심 원칙 (한 번 더)

1. **흔하면 버린다** — 기본 라이브러리 아이콘/팔레트는 즉시 제거 대상
2. **꿀벌이 산다** — 모든 비주얼 메타포는 hive에서 파생
3. **단일 바이너리** — 설치 복잡도는 타협 불가
4. **hanimo 코어 그대로** — 데스크탑은 CLI의 대안이 아닌 GUI 옵션
5. **사용자 키만** — 기업/개인 배포 모두 baked distro 모드로 키 입력만 받음

이 다섯 줄에서 벗어나는 제안은 일단 거절.

---

_Last updated: 2026-04-11 · hanimo desktop 0.0 (plan-only)_
