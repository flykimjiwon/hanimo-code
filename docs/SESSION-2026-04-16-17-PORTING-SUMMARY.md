# hanimo 대규모 포팅 + 고도화 세션 요약

> 세션: 2026-04-16 ~ 2026-04-17
> 작업자: Claude Opus 4.6 + 김지원

---

## 1. 작업 개요

TECHAI_CODE → hanimo 양방향 포팅 + 룰베이스 기능 고도화 + opencode 분석 + LSP 구현

### 수치 요약

| 항목 | 수치 |
|------|------|
| 커밋 수 | 15개 |
| 추가 코드 | +28,000줄 |
| 삭제 코드 | -50,600줄 (레거시 정리) |
| 새 파일 | 40+ |
| 새 Go 패키지 | 4개 (lsp, companion, gitinfo, multi) |
| 새 도구 | 20개 (16→36) |
| 새 슬래시 명령 | 16개 (26→42) |

---

## 2. 포팅 내역 (TECHAI_CODE → hanimo)

### 덩어리 1: 핵심 도구 + 지식

| 항목 | 파일 |
|------|------|
| diff 생성기 | `tools/diff.go` |
| 시크릿 탐지 | `tools/secrets.go` |
| 스냅샷/undo | `tools/snapshot.go` |
| .gitignore 파서 | `tools/gitignore.go` |
| 프로젝트 감지 | `tools/project.go` |
| /init 프로파일러 | `tools/init.go` |
| 커스텀 명령 | `tools/commands.go` |
| knowledge 62개 문서 | `knowledge/docs/` (16 카테고리) |

### 덩어리 2: UI + 인프라

| 항목 | 파일 |
|------|------|
| Git HUD | `gitinfo/gitinfo.go` + `gitinfo_test.go` |
| Context % 표시 | `ui/context.go` + `context_test.go` |
| 컴패니언 서버 | `companion/` (server, hub, events, browser, embed) |
| 웹 대시보드 | `web/` (index.html, style.css, app.js) |
| Compaction 테스트 | `llm/compaction_test.go` |
| Usage 테스트 | `llm/usage_test.go` |
| 입력 히스토리 | `app.go` (↑↓ + Ctrl+P/N) |

### 덩어리 3: 에이전트 + 명령

| 항목 | 파일 |
|------|------|
| 멀티에이전트 | `multi/` (orchestrator, strategy, merge, types) |
| ripgrep 통합 | `tools/ripgrep.go` |
| 심볼 검색 | `tools/symbols.go` |
| 슬래시 명령 14개 | `app.go` (/new, /git, /copy, /export, /diff, /compact, /init, /forget, /commands, /exit, /quit, /undo, /mcp, /multi) |
| 키보드 단축키 | Ctrl+J (줄바꿈), Ctrl+B (마우스 토글) |

---

## 3. 신규 룰베이스 기능 (LLM 불필요)

| 기능 | 파일 | 설명 |
|------|------|------|
| Import Graph | `tools/imports.go` | Go/JS/TS/Python import 의존성 맵 |
| Reverse Imports | `tools/imports.go` | 역방향 import 추적 |
| Change Impact | `tools/imports.go` | 변경 영향 범위 분석 (BFS) |
| Circular Deps | `tools/imports.go` | 순환 의존성 탐지 (DFS) |
| Find References | `tools/symbols.go` | 심볼 참조 위치 검색 |
| Test Coverage Gap | `tools/coverage.go` | 테스트 파일 누락 탐지 |
| Code Quality Scan | `tools/quality.go` | TODO/대형함수/깊은중첩 |
| Git Blame | `tools/git.go` | 줄별 작성자 |
| Git Hot Files | `tools/git.go` | 최근 N일 변경 빈도 |
| Git File History | `tools/git.go` | 파일별 커밋 히스토리 |
| Smart Context | `tools/smartctx.go` | 관련 파일 자동 탐색 |
| Edit Replacer Chain | `tools/replacer.go` | 5단계 매칭 (exact→trim→whitespace→indent→anchor) |
| Binary Detection | `tools/search.go` | 4KB 바이트 샘플링 |
| Fuzzy File Search | `tools/fuzzy.go` | 퍼지 파일명 검색 |
| 95개 바이너리 확장자 | `tools/search.go` | opencode 참고 확장 |

---

## 4. LSP 클라이언트

| 파일 | 설명 |
|------|------|
| `lsp/protocol.go` | LSP 프로토콜 타입 (Position, Range, Location, Diagnostic) |
| `lsp/client.go` | stdio 클라이언트 (Content-Length 프레이밍, 10초 타임아웃) |
| `lsp/servers.go` | gopls, typescript-language-server, pyright 설정 |

도구 4개: `lsp_definition`, `lsp_references`, `lsp_hover`, `lsp_symbols`

---

## 5. 문서화

| 문서 | 내용 |
|------|------|
| `README.md` | 전면 개편 — 36도구, 42명령, 12레이어 시각화 |
| `docs/roadmap/2026-04-16-hanimo-enhancement-plan.md` | 고도화 로드맵 (v0.7~v1.0, 1133줄) |
| `docs/research/2026-04-17-opencode-analysis.md` | opencode 분석 — 12개 룰베이스 기능 비교 |
| `docs/superpowers/plans/2026-04-16-techai-to-hanimo-batch1.md` | 포팅 플랜 (Batch 1 상세) |

---

## 6. 정리 작업

| 항목 | 상태 |
|------|:----:|
| `_legacy_ts/` 삭제 (3.1GB) | ✅ 완료 |
| models.go TECHAI 모델 → 범용 모델 | ✅ 완료 |
| 한국어 문자열 영어화 (tools, companion, gitinfo, web) | ✅ 완료 |
| BXM/TECHAI 참조 제거 | ✅ 완료 |
| 테스트 assertions 영어화 | ✅ 완료 |

### 미완료 (수동 삭제 필요)

```bash
cd ~/Desktop/kimjiwon/hanimo
rm -f Spark.swift welcome.html llm-intro.html "스크린샷.png" "스크린샷 2026-04-02 오후 1.04.56.png"
rm -f knowledge/docs/java/jeus.md knowledge/docs/database/tibero.md knowledge/docs/database/oracle.md
git add -A && git commit -m "chore: remove enterprise docs + stale files" && git push
```

---

## 7. 전체 도구 현황 (36개)

### LLM 연동 도구 (16개)
```
file_read, file_write, file_edit, hashline_read, hashline_edit,
list_files, list_tree, shell_exec, grep_search, glob_search,
knowledge_search, git_status, git_diff, git_log, git_commit,
diagnostics
```

### 룰베이스 도구 (16개, LLM 불필요)
```
project_detect, init_project, symbol_search, find_references,
import_graph, reverse_imports, change_impact, test_coverage_gaps,
code_quality_scan, smart_context, fuzzy_search,
git_blame, git_hot_files, git_file_history,
lsp_definition, lsp_references, lsp_hover, lsp_symbols
```

### 내부 전용 (도구로 등록되지 않음)
```
snapshot (CreateSnapshot/UndoLast), secrets (CheckSecrets),
gitignore (LoadGitIgnore), diff (GenerateUnifiedDiff),
replacer (SmartReplace), binary detect (isBinaryFile),
commands (LoadCustomCommands)
```

---

## 8. hanimo vs TECHAI_CODE 차별화 전략

| 구분 | TECHAI_CODE (폐쇄망) | hanimo (오픈소스) |
|------|---------------------|------------------|
| 네트워크 | 인터넷 불가 | 인터넷 가능 |
| 지식 | 62개 내장 필수 | 사용자 docs 폴더 + LLM 학습 데이터 |
| 모델 | 사내 전용 (GPT-OSS, Qwen-Coder) | 14+ 프로바이더 자유 선택 |
| 빌드 | baked (endpoint/key 고정) | vanilla (사용자 설정) |
| 확장 | 제한적 | MCP + Skills + LSP |
| LSP | 없음 | gopls, tsserver, pyright |
| 멀티에이전트 | 사내 모델 2개 | 다중 프로바이더 라우팅 가능 |
| 커뮤니티 | 없음 | GitHub 오픈소스 |
| 바이너리 | BXM 지식 포함 | 경량 (지식 외부화) |
| 브랜딩 | 택가이코드 | hanimo |

### 추후 hanimo_onprem 출시 시
hanimo 기반 + 내장 지식 62개 복원 + baked 모드 + 사내 모델 → hanimo_onprem으로 리네이밍

---

## 9. 다음 세션 TODO

### 높은 우선순위
- [ ] 불필요 파일 삭제 (Spark.swift, 스크린샷, JEUS/Tibero/Oracle docs)
- [ ] 웹 검색/브라우징 도구 추가
- [ ] 파일 변경 승인 흐름 (diff 미리보기 → approve/deny)
- [ ] LSP diagnostics (편집 후 자동 에러 피드백)
- [ ] 내장 knowledge 62개 제거 또는 축소 검토

### 중간 우선순위
- [ ] 이미지 입력 지원 (멀티모달)
- [ ] Anthropic 네이티브 프로바이더 (캐싱, thinking)
- [ ] 비용 제한/자동 중단
- [ ] 도구 수 최적화 (36 → 25 목표)
- [ ] frontend/ vs web/ 통합

### 낮은 우선순위
- [ ] VSCode 확장
- [ ] diff 기반 편집 모드 (unified diff)
- [ ] Git 자동 커밋 워크플로우
- [ ] 프롬프트 캐싱 실제 활용
