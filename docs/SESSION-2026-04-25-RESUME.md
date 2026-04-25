# 다음 세션 재개 가이드 — 2026-04-25

> 이 문서를 새 세션에서 첫 번째로 읽으면 이전 맥락 100% 복구 가능.
> 짝꿍 문서: `docs/SESSION-2026-04-25-INDEX.md` (이번 세션 산출 요약)
> 상위 전략: `docs/strategy/REPUTATION-STRATEGY-2026-04-23.md`
> 디자인 mock: `designs/hanimo-desktop-v1.html` (모든 Phase의 시각 기준)

---

## 1. 30초 컨텍스트 복구

- **프로젝트**: `hanimo-code-desktop` (Wails IDE, hanimo-code 서브디렉토리)
- **레포 상태**: `origin/main = d67e52a` (Phase 11 완료)
- **누적 진척**: Phase 0~11 완성 (스켈레톤 → 브랜드 → UI → 데이터 → Hash anchor 종단 종단 → Skills)
- **빌드 검증**: vite 1527 KiB · go test 49/49 · TS clean
- **다음 자연스러운 단계**: Phase 12 — MCP 클라이언트 또는 cache metrics 실값
- **사용자 정책 (반드시 준수)**:
  - hanimo-code/desktop = **완전 무료 OSS** · 유료화 금지 (목적: 명성)
  - 5-surface 중 Code + Desktop만 집중 (WebUI/RAG/Community/Spark는 보류)
  - Honey 팔레트가 브랜드 기본 테마

---

## 2. 즉시 실행 가능한 스모크 테스트

새 세션 시작 시 5분 안에 모든 게 정상인지 확인:

```bash
cd /Users/jiwonkim/Desktop/kimjiwon/hanimo-code

# 1. 동기화 상태
git status -sb
# → "## main...origin/main" 만 보여야 정상

# 2. CLI 헬스체크
go test ./internal/... -count=1 2>&1 | tail -5

# 3. Desktop Go 헬스체크
cd hanimo-code-desktop
go build ./... && go test ./... -count=1 | tail -3

# 4. Frontend 헬스체크
cd frontend && npm run build 2>&1 | tail -5
```

각 단계 실패 시 즉시 중단하고 디버깅. 모두 OK면 그대로 진행.

---

## 3. Phase 12 후보 (우선도 순)

### 🥇 Option A: MCP 클라이언트 통합 (권장 · 中복잡도 · 4~6시간)

**왜**: 사용자에게 즉각적 가치 (Jira/Slack/Figma 같은 외부 도구 채팅에서 바로 호출), placeholder 패널 1개 실제로 살림, hanimo CLI 자산 활용 패턴 (Skills 처럼) 재사용.

**입구 파일**:
- `internal/mcp/client.go` — hanimo CLI 자산
- `hanimo-code-desktop/frontend/src/components/PlaceholderPanel.tsx` (현재 mcp 분기는 여기)

**작업 단위** (5단계):
1. `hanimo-code-desktop/mcp.go` 신규 — Client struct + ListServers + CallTool
2. App.GetMCPServers() · App.CallMCPTool(server, tool, args) Wails 메서드
3. `MCPPanel.tsx` 신규 (Skills 패턴 복제)
4. App.tsx 라우팅: PlaceholderPanel → MCPPanel
5. 빌드 + 커밋 + 푸시

**리스크**: stdio transport 프로세스 라이프사이클 관리 (한 번 spawn → 세션 동안 유지)

### 🥈 Option B: Cache hit% saved $ 실값 (작은 · 1~2시간)

**왜**: MetricsRow의 마지막 placeholder를 채워 디자인 v1 mock과 100% 일치. chat.go에서 토큰 카운트만 추적하면 됨.

**입구**:
- `hanimo-code-desktop/chat.go` — usage 추적 추가
- `hanimo-code-desktop/bindings_phase3.go` — GetMetrics 실값 채움

**작업 단위**:
1. chat.go에 sessionTokens · cacheHits 카운터 추가
2. Stream 응답 시 usage 정보 누적
3. GetMetrics에서 누적값 반환
4. 빌드 + 커밋

### 🥉 Option C: LSP 서버 통합 (큰 · 1~2일)

**왜**: ProblemsStrip이 진짜 진단 표시, hanimo의 또 다른 차별화 (CLI에 internal/lsp 자산 있음).

**복잡도 ↑**: gopls/tsserver 프로세스 관리 + LSP 프로토콜 + 라이프사이클. 별도 세션 권장.

### 🏅 Option D: Subagents / Permissions / Run

**큰 작업들**. 별도 세션마다 하나씩.

---

## 4. 절대 잊지 말 것 (Constraints)

| # | 제약 |
|:-:|---|
| 1 | hanimo-code-desktop은 **완전 무료 OSS** · 유료화 제안 금지 |
| 2 | bxm·scrape-bxm·shinhan·사내 전용 코드는 desktop 포팅 절대 금지 |
| 3 | TECHAI 기준 작업 시 본질 필터 (외부망 vs 폐쇄망) 적용 |
| 4 | landing-mockups submodule은 **건드리지 않음** (다른 저장소) |
| 5 | main 브랜치 직접 푸시는 사용자 기존 패턴 — explicit consent 가정 |
| 6 | 커밋 메시지에 Constraint/Rejected/Confidence 트레일러 포함 (CLAUDE.md commit_protocol) |
| 7 | 한국어로 응답 · 사용자가 영어 강제하지 않는 한 |

---

## 5. 알려진 미해결 (의도적 보류)

| 항목 | 이유 |
|---|---|
| `navigator.platform` deprecated 경고 | Wails 환경 종속, 대체 API 비표준 |
| Wails App.d.ts 자동 재생성 | `wails dev` 한 번 돌리면 해소, 현재 `(mod as any)` 캐스트로 우회 |
| `tools/memory.go` registry 등록 | 사용자 판단 대기 (session/memory.go와 공존 여부) |
| ChatPanel 양방향 sync (Knowledge 토글이 ChatPanel header 카운트와) | Context provider 도입 필요, 우선순위 낮음 |
| split editor의 hash anchor 표시 | main editor만 ref 보관 — 단일 정책 유지 |

---

## 6. 핵심 파일 위치 빠른 참조

```
hanimo-code-desktop/
├── app.go                    — Wails App struct, file ops, lifecycle
├── chat.go                   — chatEngine + tool 핸들러 + SendMessage
├── config.go                 — TGCConfig (yaml) + LoadTGCConfig
├── bindings_phase3.go        — UndoLastEdit, GetMetrics(stub), GetProblems(stub)
├── bindings_phase6.go        — GetAvailableModels (Tier1+Ollama), SwitchModel
├── hash_anchor_emit.go       — emitHashAnchorsFor / emitHashAnchorsForLines
├── hashline.go               — hashLine + parseAnchor + HashlineRead/Edit (Phase 10)
├── skills.go                 — GetSkills + parseSkillMeta (Phase 11)
└── frontend/src/
    ├── App.tsx               — 모든 라우팅 + Top Ribbon
    └── components/
        ├── ActivityBar.tsx   — 14 icons (Phase 2)
        ├── ModeSwitcher.tsx  — Super/Deep/Plan
        ├── ProviderChip.tsx  — 드롭다운 + Tier 추론 (Phase 6)
        ├── MetricsRow.tsx    — 4열 · stub 값 fallback
        ├── ProblemsStrip.tsx — LSP 진단 strip + hash-anchor 상태
        ├── KnowledgePanel.tsx · SessionsPanel.tsx · SkillsPanel.tsx
        ├── PlaceholderPanel.tsx — MCP/Subagents/Permissions/Run/Problems/_
        ├── CommandPalette.tsx — 24 commands · multi-token fuzzy
        └── hashAnchorGutter.ts — CodeMirror StateField + GutterMarker
```

---

## 7. 커밋 메시지 템플릿 (CLAUDE.md commit_protocol 준수)

```
feat(desktop): Phase NN — <한 줄 요약>

<2-3 문장 본문>

## 신규/변경
- 파일별 핵심 변경

## 검증
- vite build: OK
- go build/test: OK 49/49

## 자가 리뷰
- 트레이드오프 / 결정 근거

Constraint: <지킨 제약>
Rejected: <대안> | <기각 이유>
Directive: <후속 작업자 주의>
Confidence: high
Scope-risk: narrow
Not-tested: <엣지케이스>

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
```

---

## 8. 새 세션에서 첫 메시지 권장

> "어제 d67e52a까지 Phase 0-11 완성. `docs/SESSION-2026-04-25-RESUME.md` 보고 Phase 12 (MCP) 진행해."

또는 옵션 명시:
> "Option A (MCP) 가자" / "Option B (cache metrics) 먼저" / "다른 거"

---

## 9. 메모리 인덱스 (참고)

`/Users/jiwonkim/.claude/projects/-Users-jiwonkim-Desktop-kimjiwon-hanimo-code/memory/MEMORY.md` 의 핵심 항목:

- `feedback_monetization_policy.md` — 유료화 금지
- `feedback_tooling_policy.md` — superpowers/gstack 활용 규칙
- `project_hanimo_techai.md` — hanimo ↔ TECHAI 페어
- `project_hanimo_ecosystem_naming.md` — 4-repo 네이밍
- `project_hanimo_desktop_design.md` — Honey 팔레트 + 8 테마
- `project_market_decisions_2026-04-24.md` — wedge 두 축 결정
- `project_scope_pivot_2026-04-23.md` — Code+IDE 올인

새 세션이 자동으로 이들을 로드하므로 별도 호출 불필요. 단 충돌 시 SESSION-RESUME 우선.
