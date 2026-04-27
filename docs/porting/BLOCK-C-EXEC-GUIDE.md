# 블록 C — techai-ide → hanimo-code-desktop IDE 포팅 실행 가이드

> **작성일**: 2026-04-24 (기존 IDE_PORTING_PLAN + DESKTOP-DESIGN-PLAN 기반 보강)
> **예상 시간**: 1~2일
> **선행 조건**: 블록 A/B 완료 권장 (TUI 코드 동기화 후 IDE 작업)
> **레포명**: `flykimjiwon/hanimo-code-desktop`

---

## 0. TL;DR

TECHAI의 `techai-ide/` (Wails GUI, 68 features, Go 2,033 LOC + React 2,911 LOC)를
`hanimo-code-desktop`으로 리브랜딩 + 다중 프로바이더 지원 추가.

---

## 1. techai-ide 현재 구조 (2026-04-24 기준)

### 1.1 Go 백엔드 (10파일, 2,033 LOC)

| 파일 | LOC | 역할 | 포팅 시 변경 |
|---|:---:|---|---|
| `main.go` | 142 | Wails 앱 부트스트랩 | 브랜딩 치환 |
| `app.go` | 356 | 앱 라이프사이클, 파일 I/O, 프로젝트 관리 | `.techai.md` → `.hanimo.md` |
| `chat.go` | 567 | **LLM 통신 핵심** — 10 도구, 스트리밍 | **대공사**: 단일 엔드포인트 → 14+ 프로바이더 |
| `config.go` | 60 | 설정 로드/저장 | `~/.tgc/` → `~/.hanimo/`, `TGC_*` → `HANIMO_*` |
| `git.go` | 200 | Git 연동 (branch, diff, log, commit) | 변경 없음 |
| `knowledge.go` | 190 | 74 knowledge packs 로딩 | 74 → 62 (BXM 13 제거) |
| `session.go` | 141 | 채팅 세션 저장/로드/자동저장 | 경로 치환만 |
| `settings.go` | 57 | 사용자 설정 UI 바인딩 | 경로 치환만 |
| `terminal.go` | 130 | PTY 터미널 | 변경 없음 |
| `toolparse.go` | 190 | tool_call 파싱 (Qwen3 호환) | 변경 없음 |

### 1.2 React 프론트엔드 (18 컴포넌트, 2,911 LOC)

| 컴포넌트 | LOC | 역할 | 포팅 시 변경 |
|---|:---:|---|---|
| `App.tsx` | 178 | 메인 레이아웃 | 브랜딩 |
| `ChatPanel.tsx` | 365 | **채팅 UI** | 프로바이더 선택 UI 추가 |
| `Editor.tsx` | 408 | 파일 에디터 (CodeMirror 6) | 변경 없음 |
| `CodeEditor.tsx` | 177 | 코드 하이라이팅 | 변경 없음 |
| `FileTree.tsx` | 275 | 파일 트리 | 변경 없음 |
| `SearchPanel.tsx` | 223 | 검색 패널 | 변경 없음 |
| `Terminal.tsx` | 202 | 터미널 패널 | 변경 없음 |
| `GitGraph.tsx` | 178 | Git 그래프 시각화 | 변경 없음 |
| `GitPanel.tsx` | 159 | Git 패널 (branch UI, 알림 뱃지) | 변경 없음 |
| `QuickOpen.tsx` | 130 | 파일 퀵 오픈 | 변경 없음 |
| `CommandPalette.tsx` | 116 | 명령 팔레트 | 변경 없음 |
| `SettingsPanel.tsx` | 104 | 설정 패널 | **14+ 프로바이더 드롭다운** 추가 |
| `ThemePicker.tsx` | 93 | 테마 선택 | Honey 테마 추가 |
| `ActivityBar.tsx` | 76 | 좌측 활동 바 (8 아이콘) | 14 아이콘으로 확장 |
| `Toast.tsx` | 61 | 토스트 알림 | 변경 없음 |
| `StatusBar.tsx` | 54 | 하단 상태바 | 브랜딩 + 프로바이더 표시 |
| `ResizeHandle.tsx` | 54 | 패널 리사이즈 핸들 | 변경 없음 |
| `AboutDialog.tsx` | 44 | 정보 다이얼로그 | 브랜딩 전체 교체 |
| `main.tsx` | 14 | React 엔트리 | 변경 없음 |

### 1.3 최근 추가 기능 (2026-04-24 커밋)

techai-ide에 마지막 3 커밋으로 추가된 기능:
- Git branch UI + 알림 뱃지
- 앱 아이콘 + 최근 프로젝트
- 자동저장 + 채팅 세션 저장 + 마크다운 미리보기

이 기능들은 **모두 포팅 대상** — 범용 IDE 기능.

---

## 2. 포팅 실행 계획 (4 Phase)

### Phase 0: 프로젝트 스캐폴딩 (30분)

```bash
# 1. 디렉토리 복사 (node_modules 제외)
mkdir -p ~/Desktop/kimjiwon/hanimo-code-desktop
rsync -av --exclude='node_modules' --exclude='dist' --exclude='.git' \
  $TECHAI/techai-ide/ ~/Desktop/kimjiwon/hanimo-code-desktop/

# 2. git init + 첫 커밋
cd ~/Desktop/kimjiwon/hanimo-code-desktop
git init
git add -A
git commit -m "feat: initial scaffold from techai-ide (pre-rebranding)"

# 3. 벌크 브랜딩 치환
find . -type f \( -name "*.go" -o -name "*.tsx" -o -name "*.ts" -o -name "*.json" -o -name "*.html" -o -name "*.md" \) \
  -exec sed -i '' \
    -e 's/techai-ide/hanimo-code-desktop/g' \
    -e 's/TECHAI IDE/hanimo-code Desktop/g' \
    -e 's/TECHAI_CODE/hanimo-code/g' \
    -e 's/techai/hanimo/g' \
    -e 's/TECHAI/HANIMO/g' \
    -e 's|tgc/internal/|hanimo/internal/|g' \
    -e 's|~/.tgc/|~/.hanimo/|g' \
    -e 's/TGC_/HANIMO_/g' \
    -e 's/.techai.md/.hanimo.md/g' \
    {} \;

# 4. go.mod 업데이트
sed -i '' 's|module tgc|module hanimo-desktop|g' go.mod
# 또는 새 모듈명 결정 후 치환
```

### Phase 1: Go 백엔드 적응 (1~2시간)

#### 1-1. chat.go — 대공사 (핵심)

현재 TECHAI의 `chat.go`는 단일 OpenAI-compat 엔드포인트.
hanimo-code-desktop은 **14+ 프로바이더** 지원 필요.

**전략**: hanimo TUI의 `internal/llm/providers/` 추상화 계층을 가져와서
`chat.go`의 LLM 호출부를 교체.

```go
// 현재 (TECHAI)
client := openai.NewClient(apiKey)
resp, err := client.CreateChatCompletion(ctx, req)

// 목표 (hanimo-desktop)
provider := providers.Get(cfg.Provider)
resp, err := provider.Chat(ctx, messages, tools)
```

**필요 작업**:
1. `internal/llm/providers/` 패키지를 desktop에도 포함 (또는 shared 패키지 참조)
2. `config.go`에 프로바이더 선택 필드 추가
3. `chat.go`에서 프로바이더 추상화 계층 사용
4. 프론트엔드 `SettingsPanel.tsx`에 프로바이더 드롭다운 추가

#### 1-2. knowledge.go — BXM 13 제거

```bash
# knowledge-packs/ 에서 BXM 관련 13개 제거
# 정확한 파일명은 TECHAI의 knowledge-packs/ 디렉토리에서 확인
ls knowledge-packs/ | grep -i bxm
# 해당 파일 삭제
```

#### 1-3. config.go / settings.go — 경로 치환

Phase 0의 벌크 치환으로 대부분 완료. 수동 확인:
```bash
grep -rn "tgc\|techai\|TGC_" *.go
```

### Phase 2: React 프론트엔드 적응 (1시간)

#### 2-1. ChatPanel.tsx — 프로바이더 선택 UI

```tsx
// 상단에 프로바이더/모델 선택 드롭다운 추가
<div className="provider-selector">
  <select value={provider} onChange={setProvider}>
    <option value="ollama">Ollama</option>
    <option value="openai">OpenAI</option>
    <option value="anthropic">Anthropic</option>
    {/* ... */}
  </select>
  <input value={model} onChange={setModel} placeholder="Model name" />
</div>
```

#### 2-2. SettingsPanel.tsx — 14+ 프로바이더 설정

프로바이더별 API 키, 엔드포인트 URL 입력 폼 추가.

#### 2-3. ThemePicker.tsx — Honey 테마 추가

```tsx
const themes = [
  // 기존 TECHAI 테마
  { name: 'dark', label: 'Dark' },
  { name: 'light', label: 'Light' },
  // hanimo 테마 추가
  { name: 'honey', label: 'Honey', colors: { primary: '#f5a623', bg: '#f4ecd8', text: '#1a1410' } },
  { name: 'paper', label: 'Paper', colors: { primary: '#b8860b', bg: '#fffdf7' } },
  { name: 'ocean', label: 'Ocean' },
  { name: 'dracula', label: 'Dracula' },
  { name: 'nord', label: 'Nord' },
  { name: 'forest', label: 'Forest' },
];
```

#### 2-4. ActivityBar.tsx — 14 아이콘 확장

기존 8 → 14:
- 추가: LSP Problems, Knowledge Packs, Skills, MCP, Subagents, Sessions

#### 2-5. StatusBar.tsx / AboutDialog.tsx — 브랜딩

벌크 치환 후 수동 확인.

### Phase 3: 빌드 + 검증 (30분)

```bash
# Go 백엔드 빌드
cd ~/Desktop/kimjiwon/hanimo-code-desktop
go build -o hanimo-desktop .

# 프론트엔드 빌드
cd frontend
npm install
npm run build

# Wails 빌드 (전체)
wails build

# 검증
[ ] go vet ./... 클린
[ ] 프론트엔드 TypeScript 에러 0건
[ ] Wails 바이너리 실행 확인
[ ] 기본 채팅 동작 확인 (Ollama 연결)
[ ] Git 패널 동작 확인
[ ] 파일 에디터 동작 확인
[ ] 터미널 동작 확인
[ ] 설정에서 프로바이더 전환 확인
```

### Phase 4: 디자인 적용 (선택, 별도 PR)

`docs/porting/HANIMO-DESKTOP-DESIGN-PLAN-2026-04-23.md` 기반:
- Top Ribbon (1 line): 브랜드 + Mode Switcher + Theme Bar + Provider Chip
- Right Panel Metrics: Context % / Cache hit% / Iter / Provider tier
- Problems Strip: 4px band below editor
- Hash-anchor 시각화: 🔒 gutter icon + undo button

---

## 3. 파일 매핑 전체 테이블

### Go 파일

| techai-ide 파일 | hanimo-desktop 파일 | 변경 수준 |
|---|---|---|
| `main.go` | `main.go` | 치환만 |
| `app.go` | `app.go` | 경로 치환 + `.hanimo.md` |
| `chat.go` | `chat.go` | **대공사** (프로바이더 추상화) |
| `config.go` | `config.go` | 경로 + 프로바이더 필드 |
| `git.go` | `git.go` | 변경 없음 |
| `knowledge.go` | `knowledge.go` | BXM 13 제거 |
| `session.go` | `session.go` | 경로 치환만 |
| `settings.go` | `settings.go` | 경로 치환만 |
| `terminal.go` | `terminal.go` | 변경 없음 |
| `toolparse.go` | `toolparse.go` | 변경 없음 |
| `app_test.go` | `app_test.go` | 치환만 |
| `config_test.go` | `config_test.go` | 치환만 |
| `toolparse_test.go` | `toolparse_test.go` | 변경 없음 |

### React 컴포넌트

| 컴포넌트 | 변경 수준 |
|---|---|
| `ChatPanel.tsx` | **중** (프로바이더 UI) |
| `SettingsPanel.tsx` | **중** (프로바이더 설정) |
| `ThemePicker.tsx` | **소** (Honey 추가) |
| `ActivityBar.tsx` | **소** (6 아이콘 추가) |
| `StatusBar.tsx` | **소** (브랜딩) |
| `AboutDialog.tsx` | **소** (브랜딩) |
| 나머지 12개 | 변경 없음 |

---

## 4. 리스크 + 완화

| 리스크 | 확률 | 완화 |
|---|:---:|---|
| Wails 빌드 환경 미설치 | 높음 | `brew install wails` 또는 `go install github.com/wailsapp/wails/v2/cmd/wails@latest` |
| npm 의존성 충돌 | 중간 | `rm -rf node_modules package-lock.json && npm install` |
| PTY 플랫폼 차이 (Windows) | 낮음 | macOS 먼저, Windows는 후순위 |
| go-openai 버전 충돌 | 중간 | `go.mod` 비교 후 최신 사용 |
| TypeScript 타입 에러 | 중간 | Wails 바인딩 재생성: `wails generate module` |
| chat.go 프로바이더 통합 복잡도 | **높음** | hanimo TUI의 `internal/llm/providers/` 그대로 가져오기 |

---

## 5. 미결정 사항 (TODO-decisions 연동)

| # | 질문 | 추천 |
|:-:|---|---|
| D1 | 아이콘 라이브러리 | Lucide (현재 사용 중, 변경 불필요) |
| D2 | 코드 에디터 엔진 | CodeMirror 6 (현재 사용 중, 유지) |
| D3 | 폰트 | Pretendard + D2 Coding (한국어 최적) |
| D4 | 테마 개수 | 8개 유지 (Honey 기본) |
| D5 | 꿀벌 로고 | 오리지널 SVG 필요 (Phase 4) |

---

## 6. 커밋 메시지 템플릿

```
feat: scaffold hanimo-code-desktop from techai-ide

Wails GUI IDE with 68 features, rebrand from TECHAI to hanimo-code.
Go backend (2,033 LOC) + React frontend (2,911 LOC).
14+ LLM provider support via hanimo provider abstraction layer.
BXM knowledge packs removed (13), Honey theme added.

Constraint: PTY terminal macOS-first, Windows deferred
Constraint: Config path ~/.hanimo/ shared with TUI
Rejected: Monaco editor | bundle size too large for Wails
Confidence: high
Scope-risk: broad
```
