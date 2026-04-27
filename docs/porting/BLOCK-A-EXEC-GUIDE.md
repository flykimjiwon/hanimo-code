# 블록 A — hanimo → TECHAI_CODE 포팅 실행 가이드

> **작성일**: 2026-04-24 (기존 BIDIRECTIONAL-ANALYSIS + sync-tracker 기반 보강)
> **예상 시간**: 반나절 (4~5시간)
> **선행 조건**: 없음 (블록 B와 독립 실행 가능)

---

## 0. TL;DR

hanimo에만 있는 **8개 항목** (LSP + 도구 5종 + MCP stdio + 테스트 3종)을
TECHAI_CODE에 복사 + import 경로 치환(`hanimo` → `tgc`) + 기능 검증.

---

## 1. 포팅 대상 파일 목록

### 1.1 LSP 클라이언트 (505 LOC)

| 파일 (hanimo) | 목적지 (TECHAI) | LOC | 설명 |
|---|---|:---:|---|
| `internal/lsp/client.go` | `internal/lsp/client.go` | ~200 | LSP JSON-RPC 클라이언트 |
| `internal/lsp/protocol.go` | `internal/lsp/protocol.go` | ~150 | LSP 프로토콜 타입 정의 |
| `internal/lsp/servers.go` | `internal/lsp/servers.go` | ~155 | gopls/tsserver/pyright 서버 관리 |

**치환 규칙**:
```bash
# 디렉토리 생성 + 복사
mkdir -p $TECHAI/internal/lsp
cp $HANIMO/internal/lsp/*.go $TECHAI/internal/lsp/

# import 경로 치환
sed -i '' 's|"hanimo/internal/|"tgc/internal/|g' $TECHAI/internal/lsp/*.go
```

**검증**:
```bash
cd $TECHAI && go vet ./internal/lsp/...
```

### 1.2 도구 5종 (룰베이스, 네트워크 불필요)

| # | 파일 (hanimo) | 목적지 (TECHAI) | LOC | 기능 |
|:-:|---|---|:---:|---|
| A2 | `internal/tools/imports.go` | `internal/tools/imports.go` | ~150 | Import graph 분석 |
| A3 | `internal/tools/coverage.go` | `internal/tools/coverage.go` | ~120 | Test coverage gap 탐지 |
| A4 | `internal/tools/quality.go` | `internal/tools/quality.go` | ~130 | Code quality 룰 스캔 |
| A5 | `internal/tools/replacer.go` | `internal/tools/replacer.go` | ~160 | 5단계 안전 편집 교체기 |
| A6 | `internal/tools/smartctx.go` | `internal/tools/smartctx.go` | ~140 | Smart context 선택 |

**치환 규칙**:
```bash
cp $HANIMO/internal/tools/{imports,coverage,quality,replacer,smartctx}.go $TECHAI/internal/tools/

sed -i '' 's|"hanimo/internal/|"tgc/internal/|g' $TECHAI/internal/tools/{imports,coverage,quality,replacer,smartctx}.go
```

**중요**: `internal/tools/registry.go`에 새 도구 5종 등록 필요:
```go
// registry.go에 추가할 항목
"import_graph":   ImportGraphTool,
"test_coverage":  TestCoverageTool,
"code_quality":   CodeQualityTool,
"file_replacer":  FileReplacerTool,  // 기존 file_edit 대체 아님, 보강
"smart_context":  SmartContextTool,
```

### 1.3 MCP stdio 전송 (이미 부분 존재)

| 파일 (hanimo) | 목적지 (TECHAI) | 비고 |
|---|---|---|
| `internal/mcp/transport_stdio.go` | `internal/mcp/transport_stdio.go` | TECHAI는 `types.go`만 있음 — stdio 구현 분리 필요 |

**주의**: TECHAI의 `internal/mcp/types.go`가 인라인으로 stdio를 처리할 수 있음.
복사 전 `types.go`와 `client.go` 비교 후 충돌 해결 필요.

```bash
# 먼저 diff 확인
diff $HANIMO/internal/mcp/client.go $TECHAI/internal/mcp/client.go
diff $HANIMO/internal/mcp/manager.go $TECHAI/internal/mcp/manager.go

# types.go ↔ transport_stdio.go 중복 확인 후 복사
```

**SSE 전송은 포팅 금지** — 폐쇄망에서 SSE 외부 서버 연결 불가.

### 1.4 테스트 3종 (회귀 방어)

| 파일 (hanimo) | 목적지 (TECHAI) | LOC |
|---|---|:---:|
| `internal/llm/prompt_embed_test.go` | `internal/llm/prompt_embed_test.go` | ~80 |
| `internal/llm/context_self_test.go` | `internal/llm/context_self_test.go` | ~60 |
| `internal/skills/loader_test.go` | `internal/skills/loader_test.go` | ~90 |

```bash
cp $HANIMO/internal/llm/{prompt_embed_test,context_self_test}.go $TECHAI/internal/llm/
cp $HANIMO/internal/skills/loader_test.go $TECHAI/internal/skills/

sed -i '' 's|"hanimo/internal/|"tgc/internal/|g' $TECHAI/internal/llm/{prompt_embed_test,context_self_test}.go
sed -i '' 's|"hanimo/internal/|"tgc/internal/|g' $TECHAI/internal/skills/loader_test.go
```

---

## 2. 실행 순서 (체크리스트)

```
Phase 1: 복사 + 치환 (30분)
  [ ] LSP 3파일 복사 + import 치환
  [ ] 도구 5파일 복사 + import 치환
  [ ] MCP stdio diff 확인 + 필요시 복사
  [ ] 테스트 3파일 복사 + import 치환

Phase 2: 통합 (1시간)
  [ ] registry.go에 도구 5종 등록
  [ ] go vet ./... 통과 확인
  [ ] go build ./cmd/tgc 성공 확인
  [ ] 컴파일 에러 수정 (import 누락, 타입 불일치)

Phase 3: 기능 검증 (1시간)
  [ ] go test ./internal/lsp/... 통과
  [ ] go test ./internal/tools/... 통과
  [ ] go test ./internal/llm/... 통과
  [ ] go test ./internal/skills/... 통과
  [ ] LSP: 바이너리 실행 → gopls 연동 테스트 (Go 프로젝트에서)

Phase 4: 글로벌 검증 (30분)
  [ ] go vet ./... 클린
  [ ] go test ./... 전체 통과
  [ ] grep -r "hanimo" internal/ → 0건 (치환 누락 없음)
  [ ] grep -r "flykimjiwon" . → 0건 (외부 URL 유입 없음)
  [ ] git diff --stat 확인 + 커밋
```

---

## 3. 포팅 금지 항목 (재확인)

| 파일 | 금지 사유 |
|---|---|
| `internal/llm/providers/anthropic.go` | 폐쇄망 도달 불가 |
| `internal/llm/providers/google.go` | 동상 |
| `internal/llm/providers/ollama.go` | 동상 (사내 프록시 사용) |
| `internal/ui/i18n.go` | 한국어 고정 |
| `internal/mcp/transport_sse.go` | 외부 SSE 서버 연결 불가 |
| `internal/llm/providers/errors_test.go` | 다중 프로바이더 전제 |
| `cmd/spark/` | 외부 제품 |

---

## 4. 커밋 메시지 템플릿

```
feat: port 8 items from hanimo-code (LSP + 5 tools + MCP stdio + 3 tests)

LSP client (gopls/tsserver/pyright) for enterprise codebase navigation.
Rule-based tools: import graph, coverage, quality, replacer, smart context.
MCP stdio transport for internal Jira/Wiki integration.
3 regression test suites for prompt embed, context self-test, skill loader.

Constraint: All items are local-only, no external network required
Constraint: SSE transport excluded (closed-network incompatible)
Rejected: Provider ports (Ollama/Anthropic/Google) | closed-network unreachable
Confidence: high
Scope-risk: moderate
```
