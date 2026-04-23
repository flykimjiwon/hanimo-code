# hanimo-code-desktop — Design v1 + 확장 기능 플랜

> 작성일: 2026-04-23
> 모체 디자인: `TECHAI_CODE/designs-v3.html` (487 LOC, 11 테마, Bubble Tea 스타일)
> 출력물: `hanimo-code/designs/hanimo-desktop-v1.html` (678 LOC, 8 테마 + 브랜드)
> 상위 문서: `docs/porting/IDE_PORTING_PLAN.md` · `docs/strategy/REPUTATION-STRATEGY-2026-04-23.md`

---

## 0. TL;DR

1. **기본 테마 Honey** — 꿀벌/앰버 팔레트를 브랜드 기본값으로. 앰버 `#f5a623` + 크림 `#f4ecd8` + 깊은 갈색 `#1a1410`.
2. **Top Ribbon 재설계** — Brand + Mode Switcher + Theme Bar + Provider Chip을 한 줄로 통합. techai-v3의 2줄 header → 1줄로 축소해 **더 가볍게**.
3. **Activity Bar 14개 아이콘** — techai 8개에서 확장: LSP Problems, Knowledge Packs, Skills, MCP, Subagents, Sessions, Web Preview, Permissions 추가.
4. **Right Panel Metrics Row** — Context % / Cache hit% / Iter / Provider tier를 한 줄에 압축. techai엔 없음.
5. **Problems Strip (LSP)** — 에디터 하단 4px 띠로 문제·경고·힌트 집계 + hash-anchor 상태. techai엔 없음.
6. **Command Palette 오버레이** — Ctrl+K 전역 단축키. v3엔 단축키 힌트만, 여기선 실제 팔레트 UI mock.
7. **Hash-anchored Edit 시각화** — gutter에 🔒 아이콘 · tools log `↺ undo` · `a3f9` 해시 표기.

---

## 1. Design v1 요소 매핑

### 1.1 techai-v3 → hanimo-desktop-v1 변환표

| 영역 | techai-v3 | hanimo-desktop-v1 | 이유 |
|---|---|---|---|
| Title | "TECHAI CODE IDE — Design v3" | "hanimo Desktop — Design v1 (Honey)" | 브랜드 |
| 상단 | Theme bar 1줄 (30px, 11 테마) | **Ribbon 1줄 (38px)** = Brand + Mode + Theme + Provider | 정보밀도 ↑ |
| 기본 테마 | Slate (blue) | **Honey (amber)** | 브랜드 |
| 테마 수 | 11 (Slate/Cursor/Linear/GitHub/Dracula/Nord/OneDark/Monokai/Solarized/Claude/Vercel) | 8 (Honey/Slate/Ocean/Dracula/Nord/Forest/GitHub/Claude/Paper) | 단순화 + 오리지널 |
| Mode switcher | 없음 (TUI에만) | **Super / Deep / Plan 3-segment** | 3모드 핵심 기능 노출 |
| Provider | `Qwen3-Coder` 고정 badge | **드롭다운 chip** `Claude Sonnet 4.6 [T1] ▼` | 14+ provider 선택 |
| Activity Bar | 5 (files/search/git/play/brain) + user/settings | **14개** (추가: problems/knowledge/skills/MCP/subagents/sessions/web/permissions) | 외부망 확장 |
| Editor gutter | 48px · 줄번호만 | **54px · 줄번호 + hash anchor 🔒** | hashline_edit 가시화 |
| Problems Strip | 없음 | 에디터-터미널 사이 4px 띠 (LSP 집계 + hash 상태) | LSP 통합 |
| Right Panel | 360px · header + chat + input + tools | **380px · header + metrics row + chat + input + tools** | 지표 가시화 |
| Metrics | 우측 상단 `2,481 tok`만 | **4열 그리드** Context·Cache·Iter·Provider | 비용/상태 투명성 |
| Chat input | placeholder 영문 고정 | **한국어 + chips** `@context /skill mcp` | 한국어 우선 + 컨텍스트 주입 UI |
| Tool log | tname + result | + `↺ undo` 버튼 per entry | 스냅샷 가시화 |
| Status bar | 좌: git/cloud/model · 우: UTF-8/Go/Tool/v | + **ctx:23% 진행바 + cache:87%** | 비용/맥락 상태 |
| Command Palette | 힌트 텍스트만 | **전역 오버레이 mock** (Ctrl+K 토글) | 접근성 |

### 1.2 Honey 팔레트 (브랜드 정식)

```
--bg-base:       #1a1410  /* 깊은 갈색 */
--bg-sidebar:    #1f1813
--bg-panel:      #221a14
--accent:        #f5a623  /* 꿀 앰버 */
--accent-glow:   rgba(245,166,35,0.22)
--success:       #6cae75  /* 꿀벌 녹 */
--warning:       #e8a317
--fg-primary:    #f4ecd8  /* 크림 */
--fg-secondary:  #c9b890
```

Claude Code·Cursor·opencode 모두 **blue/purple 주류**. Honey는 **비슷한 도구 어디에도 없는 색**이라 스크린샷 한 장에도 식별 가능.

### 1.3 라이트 테마 Paper (신규)

Claude 테마는 이미 있고, 더 중립적인 라이트 옵션 `Paper`도 추가. 종이질감 `#fffdf7` + 골드 accent `#b8860b`. 한국어 문서·학술 환경에 적합.

---

## 2. 외부망이라서 추가한 기능 (hanimo-code-desktop 전용)

폐쇄망 TECHAI IDE에는 들어갈 수 없지만 외부망 hanimo라서 넣을 수 있는 것들.

### 2.1 UI로 노출된 것 (v1 mock 포함)

| # | 기능 | 위치 | 구현 우선도 |
|:-:|---|---|:-:|
| 1 | **Provider 드롭다운** (14+ 14 공식) | Top Ribbon · `.provider-chip` | 🔴 P0 |
| 2 | **Mode Switcher** (Super/Deep/Plan) | Top Ribbon · `.modes` | 🔴 P0 |
| 3 | **Prompt Cache Indicator** (hit%, saved $) | Right Panel metrics · `Cache hit` 열 | 🔴 P0 |
| 4 | **Context % 진행바** | Status Bar · `.ctx-bar` | 🔴 P0 |
| 5 | **LSP Problems Strip** | Editor 하단 4px | 🟠 P1 |
| 6 | **Hash-anchor 가시화** | Gutter 🔒 + Tools log `a3f9` | 🟠 P1 |
| 7 | **Undo 버튼** per tool entry | Tools log · `↺ undo` | 🟠 P1 |
| 8 | **MCP · Skills · Subagents Activity 아이콘** | Activity Bar (14 아이콘) | 🟠 P1 |
| 9 | **Knowledge Packs 토글** | Activity Bar book 아이콘 · 패널 | 🟡 P2 |
| 10 | **Web Preview 패널** | Activity Bar monitor 아이콘 | 🟡 P2 |
| 11 | **Sessions / History 브라우저** | Activity Bar history 아이콘 | 🟡 P2 |
| 12 | **Permissions 5-mode indicator** | Status Bar `Ask (Shift+Tab)` | 🟡 P2 |
| 13 | **Command Palette 오버레이** | Ctrl+K 전역 | 🟢 P3 |
| 14 | **한국어 토글 배지** `한` | Right Panel header | 🟢 P3 |

### 2.2 아직 UI에 없지만 플랜에 넣어야 할 것

| 기능 | 설명 | 구현 힌트 |
|---|---|---|
| **Ollama /api/tags 자동 감지** | 로컬 모델 목록 실시간 풀 | Provider 드롭다운에 "Discovered" 섹션 |
| **HuggingFace model discovery** | 최신 OSS 모델 풀 | Provider 드롭다운 "Community" 섹션 |
| **Hot Reload 라이브 프리뷰** | 파일 수정 → 내장 브라우저 자동 새로고침 | Web Preview 패널 |
| **Gist/세션 공유 버튼** | 공개 링크로 세션 업로드 | Sessions 패널 상단 Share |
| **Git Worktree Manager** | 병렬 실험 브랜치 관리 | Activity Bar git 아이콘 확장 |
| **Subagent Dashboard** | 서브에이전트 현황/로그 | Activity Bar share-2 아이콘 |
| **Doom-loop warning toast** | 반복 tool call 감지 시 | Toast 컴포넌트 (없음) |
| **Cost estimator** | 현재 세션 토큰 × 모델 단가 | Metrics row 확장 |
| **`.hanimo/hooks.yaml` visual editor** | Hook 8 events 드래그-드롭 설정 | Settings 패널 |
| **Baked build 마법사** | 엔드포인트·키·모델 고정 배포 wizard | Settings `Build & Share` 탭 |

---

## 3. techai-ide에서 **빼는** 것 (폐쇄망 전용 기능)

| 항목 | 이유 |
|---|---|
| 단일 API endpoint 하드코딩 | 외부망은 14+ provider |
| 한국어 고정 | i18n 유지 |
| 이모지 제거 | 이모지 OK |
| BXM·JEUS·Tibero 지식팩 | 신한 전용 → 오픈소스 금지 |
| `scrape-bxm` 명령 | 동상 |
| 사내 Jira/Wiki MCP 기본 번들 | 개인 설정 |
| Audit log persistent store | 오픈소스엔 과투자 |
| Credential scrubbing 룰셋 (사내 특화) | 기본 scrubbing만 유지 |

---

## 4. 구현 로드맵 (다음 주 IDE 포팅 후)

### Phase 0 — Wails 스켈레톤 (2026-04-24, 반나절)

`docs/porting/IDE_PORTING_PLAN.md` 기존 플랜 그대로:
1. `hanimo-code/hanimo-code-desktop/` 디렉토리 생성 (또는 독립 레포)
2. techai-ide 복사 + sed 치환
3. `chat.go` 멀티프로바이더 어댑터 삽입
4. `wails build` 성공

### Phase 1 — Honey 테마 + Top Ribbon (0.5일)

1. `frontend/src/styles/themes.ts` — Honey 기본값으로 등록
2. `App.tsx` header 재구성 — Ribbon 1줄 레이아웃
3. `ModeSwitcher.tsx` 신규 컴포넌트 (Super/Deep/Plan)
4. `ProviderChip.tsx` 신규 컴포넌트 (드롭다운)
5. `ThemePicker.tsx` — 8 테마로 축소·재정렬

### Phase 2 — Activity Bar 확장 + Right Panel Metrics (1일)

1. `ActivityBar.tsx` — 14 아이콘 + dot/num 배지
2. `MetricsRow.tsx` 신규 — context/cache/iter/provider 4열
3. Go 백엔드 `app.go` — `GetMetrics()` 바인딩 추가
4. LSP Problems 집계 API — Go에서 hanimo 기존 `internal/lsp/` 호출

### Phase 3 — Hash-anchor 시각화 + Undo (1일)

1. `CodeEditor.tsx` gutter — `hashline_*` 이벤트 구독
2. 수정 중인 줄에 🔒 + 앵커 해시 `a3f9`
3. `ToolsLog.tsx` — 각 entry에 `↺ undo` 버튼
4. Go `app.go` — `UndoLastEdit()` 바인딩 (`internal/tools/snapshot.go` 활용)

### Phase 4 — Command Palette + LSP Problems Strip (0.5일)

1. `CommandPalette.tsx` — Ctrl+K 전역 + fuzzy search
2. `ProblemsStrip.tsx` — 에디터-터미널 사이 4px strip
3. Go `app.go` — `ListCommands()`, `LSPDiagnostics(file)` 바인딩

### Phase 5 — 확장 Activity 패널 (1~2일, 선택)

- `KnowledgePanel.tsx` · `MCPPanel.tsx` · `SkillsPanel.tsx` · `SubagentsPanel.tsx` · `WebPreview.tsx` · `SessionsBrowser.tsx`
- 각 패널 Go 바인딩 · 스키마 · CRUD

**총 추산**: Phase 0~4 = **3일** / Phase 5 포함 = **5일**

---

## 5. 품질 체크리스트 (Launch 전)

- [ ] Honey 테마 3개 스크린샷 (에디터 / 채팅 / 다크톤)
- [ ] Paper 라이트 테마 스크린샷 1개
- [ ] Provider 드롭다운 14개 확인
- [ ] Ctrl+K 팔레트 fuzzy 검색 동작
- [ ] Mode 전환 시 badge/placeholder/색 바뀜
- [ ] hash-anchor 🔒 gutter 표시 재현
- [ ] Cache hit% metric 실시간 갱신
- [ ] 한국어 · 영어 i18n 토글
- [ ] baked build 마법사 Wizard
- [ ] macOS Gatekeeper 우회 문서
- [ ] Windows 설치 경로 `C:\hanimo-desktop` 자동화
- [ ] Linux AppImage / Flatpak 어느 쪽?

---

## 6. 결정 대기 (디자인 관련)

| # | 질문 | 옵션 | 영향 |
|:-:|---|---|---|
| 1 | 아이콘 라이브러리 | A) Lucide (현재 v1 mock) · B) Phosphor (꿀벌 메타포 더 잘 맞음) · C) 자체 SVG (오리지널) | 브랜드 독창성 vs 개발 속도 |
| 2 | 코드 에디터 | A) CodeMirror 6 (techai-ide 기본) · B) Monaco (VSCode 엔진) · C) Shiki 정적 | 번들 크기 vs 언어 지원 |
| 3 | 폰트 | A) Geist + JetBrains Mono (현재) · B) Pretendard + D2 Coding (한국어 최적) | 한국어 가독성 |
| 4 | Theme 개수 | A) 8개 (현재) · B) 5개만 (Honey/Slate/Ocean/Claude/Paper — 더 가볍게) · C) 사용자 커스텀 JSON | 번들 크기 |
| 5 | 꿀벌 로고 아이콘 | 현재 `hexagon` lucide 사용 · 오리지널 꿀벌 SVG 제작 필요? | 브랜드 기억도 |

---

## 7. 관련 문서

- `designs/hanimo-desktop-v1.html` — 이번 세션 산출 mock
- `docs/porting/IDE_PORTING_PLAN.md` — Wails 포팅 실행 플랜
- `docs/strategy/REPUTATION-STRATEGY-2026-04-23.md` — 상위 경쟁 전략
- `docs/DESIGN-BRIEF-2026-04.md` — 브랜드 톤 브리프 (꿀벌/honeycomb 기반)
- `TECHAI_CODE/designs-v3.html` — 모체 디자인 (폐쇄망 IDE)
- `TECHAI_CODE/prototype.html` — 동작 프로토타입 (891 LOC, 참고용)
