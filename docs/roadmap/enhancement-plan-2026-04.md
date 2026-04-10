# hanimo Enhancement Plan — 2026-04

> **작성일**: 2026-04-10
> **기반 조사**: `docs/research/reference-tools-survey-2026-04.md`
> **목표**: 현재 hanimo(v0.2~) 를 2026 세대 터미널 코딩 에이전트 수준으로 도약.
> **페어**: TECHAI_CODE (`../../TECHAI_CODE/`) 와 쌍방 포팅. hanimo가 upstream, TECHAI_CODE는 simpler downstream.

---

## 0. 현재 상태 점검 (what we already have)

| 영역 | 상태 | 비고 |
|---|---|---|
| Provider 추상화 | ✅ 14+ (ollama, openai, anthropic, google, openai-compat) | `internal/llm/providers/` |
| Hash-anchored edit | ✅ `hashline_read`/`hashline_edit` | `internal/tools/hashline.go` — 업계 드문 자산 |
| MCP (stdio) | ✅ client + transport | `internal/mcp/` |
| MCP (SSE) | ✅ transport 있음 | HTTP streaming은 미완 |
| Compaction | ✅ 3-stage (snip → micro → LLM summary) | `internal/llm/compaction.go` |
| Sessions (SQLite) | ✅ save / load / fork / search | `internal/session/` |
| Knowledge store | ✅ extractor + injector + SQLite | `internal/knowledge/` — 모든 면에서 Hermes memory의 초기형 |
| Agents | ✅ plan / intent / askuser / auto | `internal/agents/` |
| Command palette | ✅ Ctrl+K | `internal/ui/palette.go` |
| Themes / i18n | ✅ 5 themes, ko/en | |
| Deep Agent / 100-iter | ✅ `ModeDev` | |
| **Repo-map** | ❌ 없음 | tree-sitter 미사용 |
| **Skill 시스템** | ❌ 없음 | `docs/skills-and-knowledge-guide.md` 는 문서만 |
| **Hooks 시스템** | ❌ 없음 | |
| **Subagent (context fork)** | ❌ 없음 | 현재 `auto.go`는 main context 공유 |
| **Worktree 병렬화** | ❌ 없음 | |
| **Permission 5-mode** | ⚠️ 부분 | plan은 read-only, 나머지는 모두 허용 |
| **LSP 통합** | ❌ 없음 | Crush 대비 약점 |
| **Prompt caching** | ❌ 없음 | Anthropic cache_control 미활용 |
| **apply_patch (unified diff)** | ❌ 없음 | hashline_edit 하나로 커버 중 |

---

## 1. 설계 원칙

1. **단일 바이너리 유지** — 플러그인/확장은 MCP로 흡수. tree-sitter는 CGO free (`smacker/go-tree-sitter`) 또는 pure-Go 대안.
2. **cost/latency 최적화** — cache breakpoint, lazy skill, repo-map incremental.
3. **Safety by default** — read-before-write 강제, dangerous regex 차단, doom loop 탐지.
4. **TECHAI_CODE와 API 호환성 유지** — 구조체/메시지 이름을 포팅 가능하게.

---

## 2. 로드맵 (3단계)

### Phase A — Foundation Upgrade (2~3주)

#### A1. Repo-map (tree-sitter + PageRank + 영속 캐시)
- **새 패키지**: `internal/repomap/`
- 파일 구조:
  - `repomap/parser.go` — tree-sitter wrapper, 언어별 query.
  - `repomap/graph.go` — symbol definition/reference 그래프.
  - `repomap/rank.go` — PageRank iterative.
  - `repomap/cache.go` — SQLite schema (`symbols(file, lang, name, kind, line, hash)`, `refs(from_sym, to_sym)`, `mtimes(file, mtime)`).
  - `repomap/inject.go` — 토큰 예산 안에서 상위 N개 심볼 마크다운으로 포맷.
- 언어 타깃(우선): go, typescript, python, rust, javascript, java.
- 도구: `repomap_search(query)`, `repomap_top(limit, budget_tokens)`.
- 시스템 프롬프트 inject: 세션 시작 시 상위 1K 토큰 분량을 `<project-map>` 블록으로 주입.
- **성공 기준**: 10K 파일 레포에서 cold start <300ms, warm <80ms.

#### A2. Skill 시스템 (SKILL.md 표준)
- **새 패키지**: `internal/skills/`
- 탐색 경로:
  1. `./.hanimo/skills/*/SKILL.md`
  2. `~/.hanimo/skills/*/SKILL.md`
  3. 모노레포 nested: 현재 파일 기준 상위 디렉토리로 올라가며.
- frontmatter 필드 (YAML):
  ```yaml
  name: deploy-vercel
  description: Vercel에 Next.js 프로젝트 배포
  allowed-tools: [shell_exec, file_read]
  user-invocable: true
  disable-model-invocation: false
  effort: medium    # low|medium|high
  ```
- **Lazy loading**: description만 세션 프롬프트에 주입, 본문은 `/skill deploy-vercel` 호출 시.
- 치환: `$ARGUMENTS`, `$ARGUMENTS[0]`, `${HANIMO_SESSION_ID}`, `${HANIMO_CWD}`.
- Pre-render: `` !`git status -s` `` 같은 inline shell을 렌더 전 실행해 동적 context 삽입.
- **모델-invocable**: 모델이 `/skill <name>` 를 tool call 로 호출 가능하게 한다.
- **추출 자동화** (Hermes 스타일, 선택적):
  - `/auto` 성공 후 knowledge injector에 훅을 달아 "이 문제 해결법을 skill로 저장할까요?" 제안.
  - `.hanimo/skills/_learned/*.md` 에 자동 저장.

#### A3. Hooks 시스템
- **새 패키지**: `internal/hooks/`
- 선언 파일: `.hanimo/hooks.yaml` + `~/.hanimo/hooks.yaml`.
  ```yaml
  hooks:
    - event: PreToolUse
      matcher:
        tool: "shell_exec"
        args_regex: "rm -rf|sudo"
      handler: { type: command, cmd: "./scripts/block-dangerous.sh" }
      on_exit_2: block
    - event: PostToolUse
      matcher: { tool: "file_edit" }
      handler: { type: command, cmd: "gofmt -w {{path}}" }
    - event: PostCompact
      handler: { type: prompt, prompt: "최근 변경 사항 요약해 주세요" }
  ```
- 이벤트 (Phase A 버전, 8개):
  - `SessionStart`, `SessionEnd`
  - `PreToolUse`, `PostToolUse`, `PostToolUseFailure`
  - `PermissionRequest`
  - `PostCompact`
  - `TaskComplete` (deep agent)
- Handler 타입: `command` (shell), `http` (POST JSON), `prompt` (다음 LLM 호출에 context 주입), `skill` (skill 호출).
- exit code 규약: 0=allow, 1=log error, 2=block (PreToolUse only), 3=skip.
- **연동 포인트**: `app.go` 의 tool execution 지점, `compaction.go` 마지막, `session.Save` 전후.

#### A4. Permission 5-mode + Rule 엔진
- 현재 `ModePlan`만 read-only. 이를 직교축으로 분리:
  - Mode (what you can do) — super / dev / plan
  - **Permission** (how strict) — normal / auto-accept / plan / dont-ask / bypass
- `Shift+Tab` 으로 permission 순환 (현재는 mode 순환). mode 순환은 `Tab` 유지.
- Rule 파일: `~/.hanimo/permissions.yaml`
  ```yaml
  allow:
    - Read(./*)
    - Shell(git status)
    - Shell(git diff)
    - Shell(go build ./...)
    - Shell(go test ./...)
  ask:
    - Shell(git push *)
    - Edit(*.env)
  deny:
    - Shell(rm -rf *)
    - Shell(sudo *)
    - Shell(curl * -H 'Authorization*')
  ```
- 세션 중 사용자가 "allow all" 선택 시 자동으로 rules에 append (학습형).
- 기존 `ModePlan` 의 `ReadOnlyTools()` 는 그대로 유지하되, permission engine에 덮어쓰기 됨.

---

### Phase B — Agent Orchestration (2~3주)

#### B1. Subagent with Context Fork
- **새 패키지**: `internal/agents/subagent.go`
- 도구: `spawn_subagent(task, system_prompt?, tools?, model?, max_iterations?)` → 500-token 요약만 반환.
- 구현:
  - 별도 goroutine + 독립 session (`session.Fork`).
  - 부모 context 대신 task + 상위 skill/memory만 전달.
  - 완료 시 `summarizer` 에이전트(작은 모델)로 압축.
- 사용 예:
  - Super 모드에서 탐색이 필요할 때: "코드베이스에서 auth 관련 파일 찾아" → subagent에 위임.
  - 테스트 실행 & 실패 분석.
  - 외부 문서 fetch 후 요약.
- 제한: 토큰 폭발 방지를 위해 동시 3개 상한, depth 2.

#### B2. Git Worktree 병렬 실행
- **새 패키지**: `internal/worktree/`
- 커맨드: `/parallel "task1" "task2" "task3"`
- 흐름:
  1. 현재 HEAD에서 3개 worktree 생성 (`.hanimo/worktrees/task-{n}-{timestamp}`).
  2. 각 worktree에 subagent + 독립 session.
  3. 완료 후 `git diff main...HEAD` 로 결과 수집.
  4. 충돌 검사 → 없으면 머지 옵션 제시, 있으면 리포트.
- `go-git` 라이브러리 사용 (CGO free).

#### B3. Plan ↔ Build 권한 명시 분리
- 현재 hanimo Plan 모드는 이미 read-only지만 "Plan 결과를 Build에 전달" 흐름이 약함.
- 강화:
  - Plan 결과를 JSON `plan.json` 으로 저장.
  - `/execute plan.json` → Build 모드로 진입 + 각 step을 개별 tool call로 실행.
  - 각 step마다 `PostStep` hook 실행 가능.

#### B4. `apply_patch` (unified diff) 도구
- **새 파일**: `internal/tools/applypatch.go`
- 포맷: Codex CLI의 `*** Begin Patch / *** Update File / *** End Patch`.
- 여러 파일 / 여러 hunk를 한 번에 처리.
- tree-sitter 검증(language 인식): 적용 후 AST 파싱 실패 시 rollback.
- `hashline_edit` 과 공존: 단일 라인 수정엔 hashline, 다중 hunk엔 apply_patch.

---

### Phase C — Polish & Scale (2주)

#### C1. Session Browser UI
- 현재 `session.Store` 있음. Bubble Tea table 컴포넌트로 browser:
  - 최근 세션 리스트 (title, model, updated, tokens).
  - Fuzzy 검색 (`/` 키).
  - `Enter` 로 resume, `Ctrl+D` 로 diff against HEAD.
- CLI: `hanimo --resume`, `hanimo --resume <id|-N>`.

#### C2. Prompt Caching
- Provider별:
  - Anthropic: `cache_control: {"type":"ephemeral"}` at system prompt 끝 + CLAUDE.md 끝.
  - OpenAI: `prompt_cache_key` 자동 생성 (session id + mode).
- 효과: 100회 turn 세션에서 입력 비용 최대 90% 절감.

#### C3. LSP 통합 (선택, 고난이도)
- `gopls`, `tsserver`, `pyright` stdio 연결.
- 도구: `lsp_diagnostics(file)`, `lsp_references(symbol)`, `lsp_rename(old, new)`.
- 우선 `gopls` 만. 이후 YAML 설정으로 확장.

#### C4. Synchronized Output Mode (Bubble Tea v2)
- 터미널 `DECSET 2026` 지원 감지 → streaming 시 flicker 제거.
- Glamour 마크다운 렌더를 일괄 커밋.

#### C5. JSON / Headless 출력 모드
- `hanimo --output json` → stdin으로 프롬프트 받아 JSONL 이벤트 스트림 출력.
- CI/CD 스크립트 친화.

---

## 3. Immediate Wins (오늘 당장 반영 가능, 각 <1h)

1. **System prompt 계층화** (`internal/llm/prompt.go`)
   - `clarifyFirstDirective` 를 파일 분리: `prompts/core.md`, `prompts/super.md`, `prompts/dev.md`, `prompts/plan.md`.
   - `//go:embed` 로 포함.
   - 장점: 수정 rebuild 최소화, diff 가독성 ↑, prompt caching breakpoint 명확.

2. **Tool description에 negative example 추가**
   - `file_edit`: "DO NOT use for new files; use file_write instead."
   - `shell_exec`: "DO NOT use for searching files; use grep_search/glob_search."
   - `hashline_edit`: anchor 포맷 예시 1줄 추가 (`"3#e4d9"`).

3. **Doom loop detector** (`internal/agents/auto.go`)
   - 최근 3 iteration의 tool call 해시가 동일하면 abort + 메시지.
   - 5 lines 안 짤막한 훅.

4. **Dangerous regex 확장** (`internal/tools/shell.go`)
   - 현재 `rm -rf /`, `sudo` 외에:
     - `export\s+(AWS|OPENAI|ANTHROPIC|GITHUB)_.*KEY`
     - `curl\s+.*-H\s+['"]Authorization`
     - `dd\s+if=.*of=/dev/`
     - `:\(\)\{\s*:\|:&\s*\};:` (fork bomb)

5. **`.hanimo.md` 자동 로딩 강화**
   - 현재 프로젝트 루트만 찾음. 부모 디렉토리로 2단계 올라가며 찾기.
   - 발견 시 system prompt에 `<project-context>` 블록으로 포함.

6. **Read-before-write 강제**
   - `internal/tools/file.go` 에서 session-local set 유지.
   - `file_edit` / `hashline_edit` / `file_write(overwrite)` 시 해당 set에 없으면 경고 반환 (강제 block은 permission mode에 따라).

---

## 4. 새 파일/패키지 매핑 (구현 시 참고)

```
internal/
├── repomap/          [A1] tree-sitter + PageRank symbol map
│   ├── parser.go
│   ├── graph.go
│   ├── rank.go
│   ├── cache.go      (SQLite)
│   └── inject.go
├── skills/           [A2] SKILL.md loader + registry
│   ├── loader.go
│   ├── frontmatter.go
│   ├── registry.go
│   ├── invoke.go
│   └── extractor.go  (learned skills from /auto)
├── hooks/            [A3] event-driven hook system
│   ├── events.go
│   ├── loader.go
│   ├── executor.go
│   └── handlers/
│       ├── command.go
│       ├── http.go
│       ├── prompt.go
│       └── skill.go
├── permission/       [A4] permission engine + rule matcher
│   ├── engine.go
│   ├── rules.go
│   └── modes.go
├── agents/
│   ├── subagent.go   [B1] context-fork subagent
│   └── summarizer.go (small model for subagent return)
├── worktree/         [B2] git worktree parallel
│   └── manager.go
├── tools/
│   └── applypatch.go [B4] unified-diff tool
└── lsp/              [C3] optional LSP bridge
    ├── client.go
    └── gopls.go
```

---

## 5. 기존 파일 수정 지점

| 파일 | 수정 |
|---|---|
| `internal/app/app.go` | hook 실행 지점 5곳, subagent spawn 라우팅, permission engine 주입 |
| `internal/llm/prompt.go` | 계층화, embed 분리, repo-map/skill description 주입 |
| `internal/llm/compaction.go` | PostCompact hook 호출, cache breakpoint 삽입 |
| `internal/tools/registry.go` | 새 도구 등록 (repomap_*, spawn_subagent, apply_patch) |
| `internal/tools/shell.go` | dangerous regex 확장, credential scrubbing |
| `internal/tools/file.go` | read-before-write 강제 |
| `internal/session/store.go` | fork 메서드, worktree 링크 |
| `internal/config/config.go` | `hooks`, `permissions`, `skills_dir`, `repomap` 섹션 추가 |
| `internal/ui/palette.go` | skill / hook / session browser 엔트리 |
| `cmd/hanimo/main.go` | `--resume`, `--output json`, `--permission` 플래그 |

---

## 6. 성공 지표 (Phase A 완료 시점)

- [ ] 10K 파일 레포에서 `hanimo` 시작 ≤ 300ms (cold), ≤ 80ms (warm).
- [ ] Skill 5개 기본 제공 (`review-pr`, `write-tests`, `refactor`, `deploy`, `debug`).
- [ ] Hooks 예제 3개 동작 (gofmt on save, block dangerous, summarize on compact).
- [ ] Permission rule 5개 학습 경험 가능 (사용자가 allow 시 yaml append).
- [ ] hashline_edit + apply_patch 둘 다 정상 동작.
- [ ] Doom loop detector가 실제로 1회 이상 무한루프 방지.

---

## 7. TECHAI_CODE와의 포팅 관계

이 계획의 각 항목은 TECHAI_CODE `docs/roadmap/enhancement-plan-2026-04.md` 의 해당 항목과 동기화. 포팅 순서:

1. **hanimo → TECHAI_CODE (catch-up)**: hash-anchored edit, MCP, multi-provider, 3-stage compaction.
2. **hanimo가 먼저**: repo-map, skills, hooks, subagent, worktree.
3. **TECHAI_CODE 실험 → hanimo 흡수**: TECHAI_CODE가 단순해서 실험 친화적. 새 아이디어는 여기서 먼저 검증 후 hanimo에 머지.

---

## 8. 다음 액션 (구체 순서)

1. 이 문서 리뷰 + 우선순위 사인오프.
2. `docs/research/reference-tools-survey-2026-04.md` 스킵 없이 정독.
3. Phase A1 (Repo-map) 스파이크: `smacker/go-tree-sitter` 로 Go 심볼 추출 PoC.
4. Phase A2 (Skill) 스파이크: `.hanimo/skills/hello/SKILL.md` 로딩 PoC.
5. Immediate Wins 6개를 1일 스프린트로 마무리.
6. 본 격 Phase A 작업 진입.
