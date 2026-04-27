# TECHAI → hanimo 지속 동기화 전략

> **작성일**: 2026-04-24
> **목적**: 택가이코드(TECHAI_CODE)가 지속 업데이트되므로, 새 기능을 hanimo로 안전하게 가져오는 반복 가능한 프로세스 정의
> **관련 문서**: `BIDIRECTIONAL-ANALYSIS-2026-04-23.md` (본질 필터), `sync-tracker-2026-04.md` (현황 대시보드)

---

## 0. 핵심 원칙

```
1. 양쪽은 독립 제품이다 — 한쪽의 커밋을 그대로 cherry-pick 하지 않는다
2. 본질 필터를 항상 적용한다 — 폐쇄망 전용 코드는 hanimo에 절대 안 온다
3. 주기적으로 delta를 확인한다 — 2주 Sprint 단위
4. 포팅은 항상 "복사 + 적응" — git subtree/submodule 사용 안 함
5. 포팅 이력은 sync-tracker에 기록한다
```

---

## 1. 동기화 주기 + 워크플로우

### 1.1 2주 Sprint 사이클

```
Sprint 시작 (월)
├── Day 1: Delta 스캔 (30분)
│   └── TECHAI 최근 2주 커밋 분석 → 포팅 후보 선별
├── Day 2~8: Sprint 본작업 (hanimo 자체 개발)
├── Day 9: 포팅 실행 (선별된 항목만)
│   └── 복사 + 치환 + 검증
├── Day 10: sync-tracker 업데이트
│   └── 포팅 완료 항목 ✅ 마킹
└── Sprint 종료
```

### 1.2 Delta 스캔 명령어

```bash
# TECHAI 최근 2주 커밋 중 internal/ 변경분만 확인
cd ~/Desktop/kimjiwon/택가이코드
git log --oneline --since="2 weeks ago" -- internal/ \
  | grep -v "docs:" | grep -v "chore:"

# 변경된 파일 목록
git diff HEAD~20..HEAD --name-only -- internal/ | sort

# hanimo에 없는 새 파일 확인
comm -23 \
  <(git diff HEAD~20..HEAD --name-only -- internal/ | sort) \
  <(cd ~/Desktop/kimjiwon/hanimo && find internal -name "*.go" | sort)
```

### 1.3 본질 필터 체크리스트

매 Delta 스캔 시 아래 질문으로 필터링:

```
[ ] 이 기능이 외부 인터넷 없이 동작하는가? (Yes → 포팅 가능)
[ ] 이 기능이 단일 프로바이더를 전제하는가? (Yes → 포팅 금지 또는 추상화 필요)
[ ] 이 기능이 BXM/Shinhan 전용인가? (Yes → 포팅 금지)
[ ] 이 기능이 한국어 고정을 전제하는가? (Yes → i18n 래핑 필요)
[ ] 이 기능이 audit/감사 전용인가? (Yes → 포팅 금지)
```

---

## 2. 파일별 동기화 분류

### 2.1 상시 동기화 대상 (양쪽 동일하게 유지)

이 파일들은 TECHAI에서 변경이 생기면 hanimo에도 반영해야 한다.

| 영역 | TECHAI 경로 | hanimo 경로 | 치환 규칙 |
|---|---|---|---|
| Agents | `internal/agents/*.go` | 동일 | `tgc` → `hanimo` |
| Compaction | `internal/llm/compaction.go` | 동일 | import만 |
| Context | `internal/llm/context.go` | 동일 | import만 |
| Multi-agent | `internal/multi/*.go` | 동일 | import만 |
| Tools (공통) | `internal/tools/{file,git,shell,search,snapshot,secrets,hashline,cosearch,diff,fuzzy,symbols,gitignore,init,project,commands,diagnostics,registry}.go` | 동일 | import만 |
| Knowledge | `internal/knowledge/*.go` | 동일 | import만 |
| Skills | `internal/skills/*.go` | 동일 | import만 |
| UI (공통) | `internal/ui/{chat,context,dev,menu,palette,persona,plan,styles,super,tabbar}.go` | 동일 | import만 |
| Companion | `internal/companion/*.go` | 동일 | import만 |
| Hooks | `internal/hooks/hooks.go` | 동일 | import만 |

### 2.2 hanimo 전용 (TECHAI에서 절대 덮어쓰지 않음)

| 파일 | 사유 |
|---|---|
| `internal/llm/providers/anthropic.go` | 다중 프로바이더 |
| `internal/llm/providers/google.go` | 동상 |
| `internal/llm/providers/ollama.go` | 동상 |
| `internal/llm/providers/errors_test.go` | 동상 |
| `internal/lsp/*.go` | hanimo 전용 (TECHAI에 역포팅 예정이지만 hanimo가 원본) |
| `internal/mcp/transport_sse.go` | 외부 SSE |
| `internal/ui/i18n.go` | 다국어 |
| `internal/ui/askuser.go` | hanimo 전용 UX |
| `internal/session/db.go` | SQLite 세션 (구조 다름) |
| `internal/session/memory.go` | 세션 메모리 (구조 다름) |
| `internal/tools/{coverage,imports,quality,replacer,smartctx}.go` | hanimo 전용 도구 (TECHAI에 역포팅 예정) |

### 2.3 TECHAI 전용 (hanimo로 절대 안 오는 것)

| 항목 | 사유 |
|---|---|
| `knowledge-packs/` BXM 13개 | 신한 전용 |
| `cmd/scrape-bxm/` | 동상 |
| `demo-supersol/` | TECHAI 데모 |
| `바보맨/` | TECHAI 전용 |
| 한국어 고정 로직 | i18n 위배 |
| audit trail 강화 | 과투자 |

---

## 3. 포팅 실행 프로토콜

### 3.1 단일 파일 포팅

```bash
# 1. 대상 파일 양쪽 diff
diff $TECHAI/internal/path/file.go $HANIMO/internal/path/file.go

# 2. TECHAI 변경이 hanimo에도 적용 가능한지 판단
#    - 본질 필터 통과?
#    - hanimo 전용 로직과 충돌 없는지?

# 3. 복사 + 치환
cp $TECHAI/internal/path/file.go $HANIMO/internal/path/file.go
sed -i '' 's|"tgc/internal/|"hanimo/internal/|g' $HANIMO/internal/path/file.go

# 4. 검증
cd $HANIMO && go vet ./internal/path/...
cd $HANIMO && go test ./internal/path/...
```

### 3.2 기능 단위 포팅 (여러 파일)

```bash
# 1. TECHAI에서 해당 기능의 관련 커밋 확인
cd $TECHAI
git log --oneline --all -- internal/path/ | head -10

# 2. 변경된 파일 목록 추출
git show <commit-hash> --name-only

# 3. 각 파일에 대해 3.1 반복

# 4. 글로벌 검증
cd $HANIMO && go vet ./...
cd $HANIMO && go test ./...
cd $HANIMO && go build ./cmd/hanimo
```

### 3.3 충돌 해결 전략

양쪽 모두 같은 파일을 수정한 경우:

```
1. hanimo 변경 우선 — hanimo가 upstream
2. TECHAI 변경이 더 나은 경우:
   - hanimo에 TECHAI 버전 적용
   - 단, hanimo 전용 로직(i18n, 다중 프로바이더 등)은 보존
3. 양쪽 변경이 독립적인 경우:
   - 3-way merge (공통 조상 = 마지막 포팅 시점)
```

---

## 4. sync-tracker 업데이트 규칙

`docs/porting/sync-tracker-2026-04.md` 업데이트 시:

```
- 포팅 완료 → ✅ + 커밋 해시 + 날짜
- 포팅 진행 → 🚧 + 담당자
- 포팅 보류 → 📋 + 사유
- 포팅 금지 → ❌ + "본질 위배" 태그
- 불일치 발견 → ⚠️ + 설명
```

---

## 5. TECHAI 모니터링 자동화 (선택)

Sprint마다 수동 스캔이 번거로우면:

### 5.1 Git hook (TECHAI 쪽)

```bash
# TECHAI/.git/hooks/post-commit
#!/bin/bash
# internal/ 변경 시 포팅 후보 로그에 기록
changed=$(git diff HEAD~1 --name-only -- internal/)
if [ -n "$changed" ]; then
  echo "$(date +%Y-%m-%d) $(git log -1 --format='%h %s')" >> .porting-candidates.log
  echo "$changed" >> .porting-candidates.log
  echo "---" >> .porting-candidates.log
fi
```

### 5.2 diff 리포트 스크립트

```bash
#!/bin/bash
# scripts/sync-delta-report.sh
# hanimo 레포에서 실행

TECHAI=~/Desktop/kimjiwon/택가이코드
HANIMO=~/Desktop/kimjiwon/hanimo

echo "=== TECHAI → hanimo Delta Report ==="
echo "Date: $(date)"
echo ""

echo "--- New files in TECHAI (not in hanimo) ---"
comm -23 \
  <(cd $TECHAI && find internal -name "*.go" | sort) \
  <(cd $HANIMO && find internal -name "*.go" | sort)

echo ""
echo "--- Modified files (different content) ---"
for f in $(cd $TECHAI && find internal -name "*.go" | sort); do
  if [ -f "$HANIMO/$f" ]; then
    if ! diff -q "$TECHAI/$f" "$HANIMO/$f" > /dev/null 2>&1; then
      echo "  DIFF: $f"
    fi
  fi
done

echo ""
echo "--- TECHAI recent commits (internal/) ---"
cd $TECHAI && git log --oneline --since="2 weeks ago" -- internal/ | head -20
```

---

## 6. 현재 동기화 상태 스냅샷 (2026-04-24)

### TECHAI에만 있는 파일 (hanimo로 아직 안 온 것)

| 파일 | 상태 | 포팅 대상? |
|---|---|---|
| `internal/hooks/hooks_test.go` | 미포팅 | **Yes** (블록 B) |
| `internal/llm/capabilities_test.go` | 미포팅 | **Yes** (블록 B) |
| `internal/llm/client_test.go` | 미포팅 | **Yes** (블록 B) |
| `internal/tools/memory.go` | 미포팅 | **Yes** (블록 B) |
| `internal/tools/apply_patch_test.go` | 미포팅 | **Yes** (범용 테스트) |
| `internal/tools/grep_patterns_test.go` | 미포팅 | **Yes** (범용 테스트) |
| `internal/tools/search_test.go` | 미포팅 | **Yes** (범용 테스트) |
| `internal/mcp/types.go` | 구조 다름 | 비교 후 결정 |

### hanimo에만 있는 파일 (TECHAI로 아직 안 간 것)

| 파일 | 상태 | 역포팅 대상? |
|---|---|---|
| `internal/lsp/` (3파일) | 미포팅 | **Yes** (블록 A) |
| `internal/tools/{imports,coverage,quality,replacer,smartctx}.go` | 미포팅 | **Yes** (블록 A) |
| `internal/llm/providers/{anthropic,google,ollama}.go` | hanimo 전용 | **No** (본질 위배) |
| `internal/mcp/transport_sse.go` | hanimo 전용 | **No** |
| `internal/ui/{i18n,askuser}.go` | hanimo 전용 | **No** |
| `internal/session/{db,memory}.go` | 구조 다름 | 비교 후 결정 |
| `internal/config/{baked_test,validate_test}.go` | 테스트 | **Yes** (블록 A) |
| `internal/knowledge/userdocs_e2e_test.go` | 테스트 | **Yes** (블록 A) |
| `internal/llm/{context_self_test,environment_test,prompt_embed_test}.go` | 테스트 | **Yes** (블록 A) |
| `internal/llm/providers/errors_test.go` | 다중 프로바이더 전제 | **No** |
| `internal/skills/loader_test.go` | 테스트 | **Yes** (블록 A) |

### TECHAI 최근 진화분 (2026-04-16 이후, 포팅 가능)

| 커밋 | 기능 | 포팅 우선순위 |
|---|---|---|
| `ef79db5` + `fbc744b` + `905199d` | reasoning_content 지원 | **높음** |
| `6044126` | 스트리밍 부분 태그 방지 | **높음** |
| `55e444a` + `5628e06` | toolparse 엣지케이스 | **중간** |
| `9884959` | config 자동 마이그레이션 | **낮음** |
| `b598651` | status bar reasoning/writing 표시 | **중간** |
| `0de24f5` | 한국어 조사 제거 fallback | **낮음** (hanimo는 i18n) |
| `47557fa` | prefetch 지속성 | **낮음** |
| `d891bcf` + `93a6e01` + `faf6712` | IDE: Git branch UI + 자동저장 + 세션 | **블록 C** |

---

## 7. 관련 문서 인덱스

| 문서 | 역할 |
|---|---|
| `BLOCK-A-EXEC-GUIDE.md` | hanimo → TECHAI 실행 가이드 |
| `BLOCK-B-EXEC-GUIDE.md` | TECHAI → hanimo 실행 가이드 |
| `BLOCK-C-EXEC-GUIDE.md` | techai-ide → hanimo-desktop 실행 가이드 |
| `BIDIRECTIONAL-ANALYSIS-2026-04-23.md` | 본질 필터 분석 (마스터) |
| `sync-tracker-2026-04.md` | 전체 동기화 대시보드 |
| `hanimo-to-techaicode-plan.md` | 상세 포팅 가이드 (fresh agent용) |
| `IDE_PORTING_PLAN.md` | 기존 IDE 포팅 플랜 (블록 C 전신) |
| `HANIMO-DESKTOP-DESIGN-PLAN-2026-04-23.md` | 데스크탑 디자인 스펙 |
