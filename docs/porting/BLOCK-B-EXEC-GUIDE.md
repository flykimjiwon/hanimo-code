# 블록 B — TECHAI_CODE → hanimo 포팅 실행 가이드

> **작성일**: 2026-04-24 (기존 BIDIRECTIONAL-ANALYSIS 기반 보강)
> **예상 시간**: 1~2시간
> **선행 조건**: 없음 (블록 A와 독립 실행 가능)

---

## 0. TL;DR

TECHAI에만 있는 **4개 영역** (memory 도구 + hooks 테스트 + LLM 테스트 2종 + toolparse 엣지케이스)을
hanimo에 복사 + import 경로 치환(`tgc` → `hanimo`) + 기능 검증.

**추가**: 2026-04-16 이후 TECHAI에서 진화한 **새 기능 7건** (reasoning_content, tool_call 파싱 강화,
auto-prefetch, status bar 개선 등)도 선별 포팅.

---

## 1. 기존 계획 항목 (4건)

### B1. memory.go — `/remember` 도구 (263 LOC)

| 파일 (TECHAI) | 목적지 (hanimo) | 설명 |
|---|---|---|
| `internal/tools/memory.go` | `internal/tools/memory.go` | JSON 기반 `/remember`, `/memories` 스토어 |

**주의**: hanimo에는 `internal/session/memory.go`가 이미 있음 (세션 메모리).
TECHAI의 `tools/memory.go`는 **도구(tool) 형태** — LLM이 직접 호출하는 `/remember` 커맨드.
두 파일은 역할이 다르므로 충돌 없음.

```bash
cp $TECHAI/internal/tools/memory.go $HANIMO/internal/tools/

sed -i '' 's|"tgc/internal/|"hanimo/internal/|g' $HANIMO/internal/tools/memory.go
```

**통합**: `internal/tools/registry.go`에 등록:
```go
"memory_store":  MemoryStoreTool,   // /remember, /memories
```

### B2. hooks_test.go (151 LOC)

| 파일 (TECHAI) | 목적지 (hanimo) |
|---|---|
| `internal/hooks/hooks_test.go` | `internal/hooks/hooks_test.go` |

```bash
cp $TECHAI/internal/hooks/hooks_test.go $HANIMO/internal/hooks/

sed -i '' 's|"tgc/internal/|"hanimo/internal/|g' $HANIMO/internal/hooks/hooks_test.go
```

### B3. LLM 테스트 2종 (578 LOC)

| 파일 (TECHAI) | 목적지 (hanimo) | LOC |
|---|---|:---:|
| `internal/llm/capabilities_test.go` | `internal/llm/capabilities_test.go` | 80 |
| `internal/llm/client_test.go` | `internal/llm/client_test.go` | 498 |

```bash
cp $TECHAI/internal/llm/{capabilities_test,client_test}.go $HANIMO/internal/llm/

sed -i '' 's|"tgc/internal/|"hanimo/internal/|g' $HANIMO/internal/llm/{capabilities_test,client_test}.go
```

**주의**: `client_test.go`가 498줄로 상당히 큼. TECHAI의 `client.go`와 hanimo의 `client.go`가
구조적으로 달라졌을 수 있으므로, 테스트 실행 전 두 `client.go` diff 확인 필수:
```bash
diff $TECHAI/internal/llm/client.go $HANIMO/internal/llm/client.go | head -50
```

### B4. toolparse 엣지케이스

TECHAI 커밋에서 tool_call 파싱이 대폭 강화됨. 관련 커밋:
- `6a132cf` — text-based tool_calls 파싱 (Qwen3 proxy compat)
- `55e444a` — 엣지케이스 5건 보강
- `5628e06` — `<parameter=key>` value 형식 지원
- `6044126` — 스트리밍 중 부분 태그 노출 방지

**현황**: hanimo의 toolparse는 `internal/exec/exec.go` 또는 `internal/llm/client.go`에 인라인.
TECHAI는 `techai-ide/toolparse.go` (190 LOC) + `toolparse_test.go`로 분리됨.

**포팅 전략**:
```bash
# 1. TECHAI의 toolparse 로직이 어디에 있는지 확인
grep -rn "parseToolCall\|toolparse\|tool_call" $TECHAI/internal/ --include="*.go" | head -20

# 2. hanimo에서 대응 로직 확인
grep -rn "parseToolCall\|tool_call" $HANIMO/internal/ --include="*.go" | head -20

# 3. 차이 비교 후 엣지케이스 패치 적용
```

---

## 2. 신규 포팅 대상 (2026-04-16 이후 TECHAI 진화분)

TECHAI가 지속 업데이트되면서 hanimo에 아직 없는 기능들:

### N1. reasoning_content 지원 (TECHAI 커밋 `ef79db5`, `fbc744b`, `905199d`)

| 기능 | 설명 | hanimo 영향 |
|---|---|---|
| `reasoning_content` 필드 파싱 | GPT-OSS-120B 등 reasoning 모델의 사고 과정 분리 표시 | `internal/llm/client.go` 스트리밍 로직 수정 |
| thinking/content 스트림 분리 | 사고 과정(회색) vs 최종 답변(흰색) 색상 구분 | `internal/ui/chat.go` 렌더링 수정 |
| status bar reasoning/writing 표시 | 현재 단계를 상태바에 표시 | `internal/ui/tabbar.go` 수정 |

**판정**: **포팅 가능** — Ollama Qwen3도 reasoning_content 지원. 범용 기능.

### N2. auto-prefetch (TECHAI 커밋 `47557fa`, `0de24f5`, `eab9817`)

| 기능 | 설명 | 판정 |
|---|---|---|
| 한국어 조사 제거 후 파일명 검색 | "파일을" → "파일" fallback | **포팅 가능** (i18n 확장) |
| prefetch 지속성 | tool iteration 간 prefetch 상태 유지 | **포팅 가능** |
| prefetch 시 file_write 강제 | apply_patch 대신 file_write 사용 | **검토 필요** (hanimo는 hashline 우선) |

### N3. 스트리밍 tool_call 부분 태그 방지 (TECHAI 커밋 `6044126`)

**판정**: **포팅 필수** — 모든 모델에서 발생 가능한 UX 버그.

### N4. config 자동 마이그레이션 (TECHAI 커밋 `9884959`)

**판정**: **선택적 포팅** — hanimo는 다중 프로바이더라 config 구조가 다름.

---

## 3. 실행 순서 (체크리스트)

```
Phase 1: 기존 4건 복사 + 치환 (20분)
  [ ] memory.go 복사 + import 치환
  [ ] hooks_test.go 복사 + import 치환
  [ ] capabilities_test.go + client_test.go 복사 + import 치환
  [ ] client.go diff 확인 (테스트 호환성)

Phase 2: toolparse 엣지케이스 (30분)
  [ ] TECHAI의 toolparse 변경분 확인 (4 커밋)
  [ ] hanimo 대응 코드 위치 확인
  [ ] 엣지케이스 패치 수동 적용 (cherry-pick 불가 — 구조 다름)

Phase 3: 신규 기능 선별 포팅 (40분)
  [ ] reasoning_content 필드 파싱 (client.go 스트리밍)
  [ ] thinking/content 색상 분리 (chat.go)
  [ ] 스트리밍 부분 태그 방지 (client.go)
  [ ] 한국어 조사 제거 fallback (선택)

Phase 4: 검증 (20분)
  [ ] go vet ./... 클린
  [ ] go test ./internal/hooks/... 통과
  [ ] go test ./internal/llm/... 통과
  [ ] go test ./internal/tools/... 통과
  [ ] grep -r "tgc/" internal/ → 0건 (치환 누락 없음)
  [ ] 바이너리 빌드 확인: go build ./cmd/hanimo
```

---

## 4. 포팅 금지 항목 (재확인)

| 항목 | 금지 사유 |
|---|---|
| BXM 지식팩 13개 (`knowledge-packs/`) | 신한 전용 — 법무/브랜딩 리스크 |
| `cmd/scrape-bxm/` | 동상 |
| 한국어 고정 로직 (이모지 제거, 영어 제거) | hanimo는 i18n 유지 |
| 단일 게이트웨이 하드닝 | hanimo는 다중 프로바이더 |
| Audit trail 강화 | 오픈소스엔 과투자 |
| `demo-supersol/` | TECHAI 전용 데모 |
| `바보맨/` | TECHAI 전용 모듈 |

---

## 5. 커밋 메시지 템플릿

```
feat: port from TECHAI_CODE — memory tool + tests + toolparse + reasoning_content

Block B items:
- memory.go: /remember, /memories JSON store (263 LOC)
- hooks_test.go: lifecycle hooks regression tests (151 LOC)
- LLM capabilities + client tests (578 LOC combined)
- toolparse edge cases: Qwen3 proxy compat, partial tag prevention

New from TECHAI (2026-04-16~24):
- reasoning_content field parsing for thinking models
- thinking/content stream color separation
- streaming partial tool_call tag prevention

Constraint: No BXM knowledge packs or closed-network-specific code
Rejected: auto-prefetch file_write force | conflicts with hashline-first strategy
Confidence: high
Scope-risk: narrow
```
