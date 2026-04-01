# Claude Code 소스코드 유출 분석서

> **작성일**: 2026-04-01
> **대상**: Claude Code npm 패키지 v2.1.88 (2026-03-31 유출)
> **출처**: choi.openai 스레드, WikiDocs 별첨 91, instructkr/claw-code GitHub, 웹 리서치

---

## 1. 유출 경위

2026년 3월 31일, 보안 연구자 **Chaofan Shou** (@shoucccc)가 `@anthropic-ai/claude-code` npm 패키지 v2.1.88에서 **59.8MB 소스맵 파일** (`cli.js.map`)을 발견. 이 `.map` 파일에 `sourceMappingURL` 디렉티브가 Anthropic의 **Cloudflare R2** 버킷을 가리키고 있어, 비난독화된 전체 TypeScript 소스를 다운로드할 수 있었음.

**근본 원인**: Anthropic이 2025년 말 Bun을 인수하여 빌드에 사용. Bun의 알려진 이슈(oven-sh/bun#28001)로 프로덕션에서도 소스맵이 생성됨. `.npmignore`에 `*.map`을 추가하지 않은 실수.

- **총 파일 수**: ~1,900개 TypeScript 파일
- **총 코드량**: 512,000+ 라인
- **빌드 시스템**: Bun (esbuild 기반)
- **UI**: React + Ink (TUI), Zustand (상태관리)
- **보존 리포**: `instructkr/claw-code` (GitHub)

### 코드 규모 참고

| 파일 | 라인 수 |
|------|--------|
| `QueryEngine.ts` | ~46,000 |
| `Tool.ts` | ~29,000 |
| `commands.ts` | ~25,000 |
| `print.ts` | 5,594 (단일 함수 3,167줄, 12단계 중첩) |
| `bashSecurity.ts` | 2,592 |

---

## 2. 전체 아키텍처

### 2.1 4단계 실행 파이프라인 (TAOR)

```
Think → Act → Observe → Repeat
```

1. **Startup**: 설정 로드, 권한 모드 결정, 컨텍스트 수집, 피처 플래그 평가
2. **Query Loop**: 사용자 입력 → LLM 호출 → 응답 스트리밍 (반복)
3. **Tool Execution**: 10단계 파이프라인 (스키마 검증 → 권한 → 실행 → 결과)
4. **Display**: React/Ink TUI 또는 JSON 스트림 (headless)

### 2.2 6가지 실행 모드

| 모드 | 설명 |
|------|------|
| **REPL** | 기본 대화형 CLI (React/Ink TUI) |
| **Headless** | JSON 스트림 I/O, CI/CD 통합용 |
| **Coordinator** | Leader + Worker 멀티 에이전트 오케스트레이션 |
| **Bridge** | 로컬 CLI ↔ 클라우드(CCR) WebSocket 연결 |
| **Kairos** | 데몬 모드 — 백그라운드 상주, 트리거 기반 실행 |
| **Viewer** | 세션 재생 뷰어 (읽기 전용) |

### 2.3 핵심 클래스 구조

```
QueryEngine (46,000 lines)
  ├── ToolManager (40+ tools, 85+ slash commands)
  ├── PermissionManager (5 modes)
  ├── CompactionManager (5 strategies)
  ├── HookManager (13 events, 5 hook types)
  ├── MemoryManager (MEMORY.md + autoDream)
  ├── SessionManager (JSONL, fork, replay)
  ├── FeatureFlagManager (44 flags)
  └── BashSecurityManager (23+ security checks)
```

### 2.4 시스템 프롬프트 어셈블리

- **~60개 시스템 프롬프트 컴포넌트** + **~40개 시스템 리마인더** 조건부 조립
- `SYSTEM_PROMPT_DYNAMIC_BOUNDARY` 마커: 전역 캐시 가능 영역 / 세션별 동적 영역 분리
- CLAUDE.md는 **매 쿼리 반복마다 재로딩** (세션 중 수정 즉시 반영)
- CLAUDE.md 문자 제한: ~40,000자
- 도구 정의: 알파벳 순 정렬 (캐시 히트율 극대화)

### 2.5 멀티모델 전략

| 모델 | 용도 |
|------|------|
| **Sonnet 3.5/4.0** | 메인 추론 루프 |
| **Haiku 3.5/4.5** | 경량 작업 (제목 생성, 분류, sub-agent) |
| **Opus 4.6** | ULTRAPLAN 심층 계획 |

---

## 3. 도구 시스템 (Tools)

### 3.1 40+ 내장 도구 + 85+ 슬래시 명령

**파일 조작**: Read, Write, Edit (MultiEdit), Glob, Grep
**실행**: Bash, BashBackground (백그라운드)
**에이전트**: Task (subagent 생성), SendMessage
**네비게이션**: WebFetch, WebSearch, NotebookEdit
**상호작용**: AskUserQuestion, EnterPlanMode, ExitPlanMode, TodoWrite, TodoRead
**고급**: EnterWorktreeTool, ExitWorktreeTool (git worktree 격리)
**MCP 발견**: 8가지 스코프 (config, CLI, env, org policies, .claude, Claude.ai 등)

### 3.2 10단계 도구 실행 파이프라인

```
1. 스키마 검증 (Zod)
2. PreToolUse 훅 실행
3. 권한 확인 (PermissionManager)
4. 경로 샌드박스 검증
5. 동시성 분류 (safe/unsafe)
6. 도구 실행 (execute)
7. 결과 크기 제한 (truncation)
8. PostToolUse 훅 실행
9. 결과 캐싱
10. 결과 반환
```

### 3.3 동시성 모델

- **Safe 도구** (Read, Glob, Grep, TodoWrite, ToolSearch 등): 최대 10개 병렬 실행
- **Unsafe 도구** (Write, Bash 등): 순차 실행
- 분류 기준: `toolDefinition.isReadOnly` 또는 화이트리스트
- Write/Edit 도구: 프로젝트 내부면 fast-path, 외부면 분류기 실행

### 3.4 MCP 통합

- Model Context Protocol SDK 내장
- stdio/SSE/streamable-http 3가지 전송 방식
- 외부 MCP 서버 도구를 네이티브 도구와 동등하게 취급
- 도구 스키마 자동 변환 (MCP → AI SDK)
- **주의**: MCP 도구 결과는 microcompaction 대상에서 제외 (autocompact까지 유지)

---

## 4. 권한 시스템 (Permissions)

### 4.1 5가지 권한 모드

| 모드 | 설명 |
|------|------|
| **plan** | 읽기 전용 모드 |
| **default** | 모든 동작에 대해 대화형 확인 |
| **acceptEdits** | 파일 편집 자동 승인, 기타 확인 |
| **auto** | ML 기반 LLM 분류기로 자동 승인 (리서치 프리뷰) |
| **bypassPermissions** | 전체 허용 (원격 킬스위치 가능) |

### 4.2 Auto 모드 2단계 LLM 분류기

```
Stage 1: 빠른 단일 토큰 필터
  → allow / soft_deny / hard_deny + 신뢰도 + 근거

Stage 2: Chain-of-Thought 추론 (Stage 1에서 플래그 시)
  → 컨텍스트 기반 정밀 판단
  → 여러 resolver를 병렬 레이싱
```

- 화이트리스트 safe 도구 (Read, Grep, Glob, TodoWrite, ToolSearch): **분류 스킵**
- 프로젝트 내부 Write/Edit: fast-path 처리
- 외부 경로 Write/Edit: 분류기 실행

### 4.3 glob 패턴 규칙 시스템

`settings.json`에서 3가지 규칙 유형:
- **Allow**: 수동 승인 없이 실행
- **Ask**: 확인 프롬프트
- **Deny**: 실행 차단

```json
{
  "permissions": {
    "allow": ["Bash(npm test*)", "Bash(git commit *)"],
    "deny": ["Bash(rm -rf *)", "Bash(curl *)"],
    "ask": ["Edit(*.config.*)", "Bash(git push *)"]
  }
}
```

**평가 순서**: deny → ask → allow (첫 매칭 규칙 적용)
**5단계 설정 캐스케이드**: CLI 인자 > 프로젝트 `.claude/settings.local.json` > 글로벌 `~/.claude/settings.json` > LLM 판단 > 기본값

### 4.4 Bash 보안 (bashSecurity.ts — 2,592줄)

**23+ 보안 검사**:
- 18개 Zsh 빌트인 차단
- Zsh equals 확장 방어 (`=curl` 우회 방지)
- 제로 폭 공백 주입 감지
- IFS null-byte 주입 방어
- HackerOne 발견 토큰 우회 방어
- **3개 독립 파서**: `splitCommand_DEPRECATED`, `tryParseShellCommand`, `ParsedCommand.parse`
- **Tree-sitter AST** 분석으로 구조적 명령 파싱
- Heredoc 내 숨겨진 명령 추출
- 보조 Haiku 모델 호출로 명령 접두사 검증
- Git 작업은 trust 대화 승인까지 지연

**알려진 약점**: `validateGitCommit`이 `allow` 반환 시 → `bashCommandIsSafe` short-circuit → `validateRedirections` **실행 안 됨** (소스 주석에 기록됨)

### 4.5 위험도 분류

모든 동작에 **LOW / MEDIUM / HIGH** 위험도 지정.
보호 파일: `.gitconfig`, `.bashrc`, `.zshrc`, `.mcp.json`, `.claude.json`

### 4.6 관련 CVE

| CVE | 설명 |
|-----|------|
| CVE-2025-52882 | IDE origin/websocket 검증 |
| CVE-2025-59828 | Pre-trust Yarn config 실행 |
| CVE-2025-58764 | 명령 파싱 프롬프트 우회 |
| CVE-2025-64755 | sed 파싱으로 읽기 전용 우회 |
| CVE-2026-21852 | 악성 repo 설정으로 API 키 유출 |

---

## 5. 컨텍스트 압축 (Compaction)

`query.ts`의 4단계 파이프라인에서 트리거되는 **5가지 전략**:

### 5.1 Snip Compact

- 오래된 메시지를 임계값 이하로 제거
- 가장 빠르지만 손실적
- 도구 결과 본문을 `[snipped: 200 lines]`로 대체

### 5.2 Microcompact

- 턴 내 도구 결과를 **도구 호출 구조는 보존**하면서 압축
- API 호출 없이 로컬 캐시 편집
- `COMPACTABLE_TOOLS`에 속한 도구만 대상
- **MCP 도구 결과, Agent 결과, 커스텀 도구는 제외** (autocompact까지 유지)
- `maxResultSizeChars: Infinity`인 파일 읽기도 제외
- `CACHED_MICROCOMPACT` 플래그: 동일 microcompaction 재계산 방지

### 5.3 Auto-Compact (LLM 요약)

- 컨텍스트 윈도우의 설정된 % 초과 시 발동
- **포크된 서브프로세스**에서 실행
- 최대 **20,000 토큰** 요약 생성 (13,000 토큰 버퍼)
- 프롬프트: "사용자 피드백에 특별히 주의", "도구 결과가 아닌 모든 사용자 메시지 보존"
- 압축 후 모델에게 "사용자에게 추가 질문 없이 계속하라" 지시

**발견된 버그**: 1,279개 세션에서 **50+ 연속 실패** — 전역적으로 일일 ~250,000 API 호출 낭비.
수정: `MAX_CONSECUTIVE_AUTOCOMPACT_FAILURES = 3` (3회 연속 실패 후 세션 내 압축 비활성화)

### 5.4 Reactive Compact

- API가 "prompt too long" 에러 반환 시 긴급 폴백
- `REACTIVE_COMPACT` 피처 플래그 뒤에 게이트

### 5.5 Context Collapse

- 오래된 메시지의 단계적 공격적 정리
- 지능적 메시지 선택 제거
- `CONTEXT_COLLAPSE` 피처 플래그 뒤에 게이트

### 5.6 추가 메커니즘: Content Replacement

- 큰 도구 결과를 **참조 포인터**로 대체, 원본은 별도 저장
- 메시지별 예산 적용
- `seenIds`로 캐시된 결정은 세션 내 잠금

### 5.7 수동 압축

- `/compact <focus>` 명령으로 사용자 주도 요약 + 토픽 우선순위 지정

---

## 6. 훅 시스템 (Hooks)

### 6.1 13개 확인된 이벤트

| 이벤트 | 설명 |
|--------|------|
| **SessionStart** | 세션 시작/재개 시 |
| **SessionEnd** | 세션 종료 시 (차단 불가) |
| **UserPromptSubmit** | 사용자 입력 제출 시, 처리 전 |
| **PreToolUse** | 도구 실행 전 (차단/수정 가능) |
| **PostToolUse** | 도구 실행 후 |
| **PostToolUseFailure** | 도구 실행 실패 시 |
| **PermissionRequest** | 권한 대화 표시 시 |
| **PreCompact** | 대화 압축 직전 |
| **SubagentStart** | 서브에이전트 생성 시 |
| **SubagentStop** | 서브에이전트 종료 시 |
| **Notification** | 알림 이벤트 시 |
| **Stop** | Claude 생성 중단 시 |

### 6.2 5가지 훅 유형

| 유형 | 설명 |
|------|------|
| **Shell Command** (`command`) | 외부 쉘 명령 실행; exit code 2 = 차단 |
| **LLM 주입** (`prompt`) | LLM으로 허용/차단 평가 |
| **Agent 검증** (`agent`) | 전체 도구 접근 가능한 에이전트 검증기 spawn |
| **HTTP Webhook** | 외부 엔드포인트에 POST |
| **JavaScript 함수** | 직접 JS 실행 |

### 6.3 훅 실행 모델

- exit code 0 = 통과, exit code 2 = **차단**
- stdout은 `<system-reminder>` 태그로 LLM 컨텍스트에 주입
- 타임아웃 설정 가능 (기본 10초)
- 훅은 **사용자/관리자만** 구성 가능 (모델이 자기 수정 불가)
- SessionStart는 command 유형만 지원

### 6.4 UserPromptSubmit 보안 주의

- 사용자 요청 처리 **전**에 트리거 → 시스템 프롬프트에 임의 내용 주입 가능
- `.claude/settings.local.json`의 훅이 주입한 내용 → 모델이 처리 → 어시스턴트 메시지에 반영 → autocompact 시 "사용자 피드백"으로 보존 → **지시문 세탁(laundering)** 가능

---

## 7. 메모리 시스템 (Memory)

### 7.1 6계층 메모리 계층구조

| 계층 | 설명 | 압축 우선순위 |
|------|------|-------------|
| **Managed Policies** | 조직 수준 규칙 | 최고 (절대 삭제 안 함) |
| **Project Config** | CLAUDE.md 파일들 | 매우 높음 |
| **User Preferences** | `~/.claude/` 설정 | 높음 |
| **Session History** | 현재 대화 기록 | 중간 |
| **Auto-Learned** | MEMORY.md (25KB 제한) | 낮음 |
| **Real-Time** | 실시간 상호작용 스트림 | 최저 (가장 먼저 압축) |

### 7.2 MEMORY.md

- 위치: `~/.claude/projects/<project-hash>/memory/MEMORY.md`
- **200라인 또는 25KB** 중 먼저 도달하는 제한
- 매 세션 시작 시 자동 로드 → 시스템 프롬프트에 주입
- **인덱스/포인터 시스템** 역할: 라인당 ~150자, 전체 데이터가 아닌 위치 참조
- YAML frontmatter로 4가지 메모리 유형: `user`, `feedback`, `project`, `reference`
- **Sonnet 기반 관련성 선택기**: 턴당 최대 5개 관련 메모리 검색

### 7.3 autoDream 메모리 통합

UC Berkeley의 "Sleep-time Compute" 연구에서 영감. **포크된 서브에이전트**로 실행.

**4가지 필수 트리거 게이트**:
1. 마지막 통합으로부터 >= 24시간 경과
2. >= 5개 새 세션 축적
3. 활성 통합 프로세스 없음 (락 획득)
4. 마지막 스캔으로부터 >= 10분 경과

**4단계 통합 프로세스**:
1. **Orient**: 디렉토리 스캔 + 기존 메모리 리뷰
2. **Gather Signal**: 최근 로그/트랜스크립트에서 시그널 추출
3. **Consolidate**: 메모리 파일 작성/갱신, 절대 날짜 기록, 이산 관찰 병합, 논리적 모순 제거, 모호한 인사이트 → 구조화된 사실로 변환
4. **Prune**: MEMORY.md를 **200라인/25KB** 이하로 유지

### 7.4 에이전트별 영속 메모리

- 옵션: `~/.claude/agent-memory/<name>/`
- `AGENT_MEMORY_SNAPSHOT` 피처 플래그 뒤에 게이트
- 세션 간 에이전트별 지식 축적

---

## 8. 세션 관리 (Sessions)

### 8.1 JSONL 형식

- **Append-only JSONL** (JSON Lines): 새 메시지는 파일 끝에 추가, 기존 데이터 재작성 없음
- 저장 위치: `~/.claude/projects/<encoded-cwd>/*.jsonl`
- `<encoded-cwd>`: 절대 경로에서 비영숫자를 `-`로 대체
- 알려진 이슈: 큰 JSONL 파일이 메모리 전체 소비 (GitHub #22365)

### 8.2 세션 제어

| 명령 | 설명 |
|------|------|
| `--continue` | 현재 디렉토리의 가장 최근 세션 재개 |
| `--resume <id>` | 특정 세션 ID로 재개 |
| `--fork-session` | 원본 히스토리 복사 후 분기, 두 세션 독립 유지 |

### 8.3 Attribution Tracking (기여 추적)

`commitAttribution.ts`로 파일별 수정 이력 추적:
- 프롬프트 횟수, 권한 상호작용 횟수
- 사용 환경 식별 (CLI, VS Code, Web)
- 변경당 모델/버전 연결
- AI 기여 비율 (%) + 문자 수 (Claude vs 사람)
- git notes에 저장 (기본적으로 공개 커밋에 포함 안 됨)

---

## 9. 서브에이전트 시스템 (3가지 모델)

### 9.1 Fork 모델

- 부모 컨텍스트의 **바이트 동일 복사** 생성
- 초기 컨텍스트가 동일하므로 **프롬프트 캐시 히트** → spawn 비용 저렴
- 독립 컨텍스트 윈도우, 별도 토큰 예산, 독립 압축 트리거
- 내장 에이전트 타입: **ExploreAgent** (읽기 전용, Haiku), **VerificationAgent** (작업 검증)
- VerificationAgent: 3+ 작업 완료 후 검증 없으면 자동 트리거

### 9.2 Teammate 모델

- **파일 기반 메일박스** 시스템으로 터미널 패인 간 통신
- `AsyncLocalStorage`로 프로세스 격리
- 공유 스크래치패드 디렉토리 (`tengu_scratch` 피처 플래그) → 내구적 크로스 워커 지식 공유

### 9.3 Worktree 모델

- 각 에이전트가 독립 **git 브랜치** 할당 (`EnterWorktreeTool` / `ExitWorktreeTool`)
- 병렬 작업 시 머지 충돌 방지

### 9.4 서브에이전트 격리 속성

- 독립 컨텍스트 윈도우 (fork 후 부모와 공유 안 함)
- 별도 토큰 예산 (비용 폭주 방지)
- 제한된 도구 접근 (대부분 모드에서 sub-agent 재spawn 불가)
- 독립 압축 트리거
- 새 MCP 서버 연결, 종료 시 정리

---

## 10. Coordinator 모드 (멀티 에이전트)

`CLAUDE_CODE_COORDINATOR_MODE=1`로 활성화. 오케스트레이션 알고리즘이 **시스템 프롬프트 지시**로 구현 (하드코딩 아님).

### 10.1 4단계 워크플로우

```
1. Research    — Worker들이 병렬 조사
2. Synthesis   — Coordinator가 결과 이해, 구현 스펙 작성
3. Implementation — Worker들이 스펙에 따라 타겟 변경
4. Verification — 독립 검증 에이전트가 PASS/FAIL 판정
```

### 10.2 핵심 원칙

- **Coordinator는 자체 검증 불가** — 독립 검증자만 판정
- 핵심 지시: "약한 작업을 고무도장 찍지 말라", "후속 작업 지시 전에 결과를 이해하라"
- Worker는 공유 스크래치패드에 권한 프롬프트 없이 접근
- `atomicClaim` 메커니즘 (`createResolveOnce`): 여러 Worker가 같은 권한 요청 중복 처리 방지

### 10.3 vs hanimo Orchestrator

| 항목 | Claude Code Coordinator | hanimo Orchestrator |
|------|------------------------|--------------------|
| Leader 역할 | 계획/감독만, 코드 수정 불가 | LLM이 분해+합성 |
| Worker 격리 | 독립 QueryEngine | 독립 SubAgent |
| 통신 | SendMessage 양방향 | 결과 합성만 (단방향) |
| 공유 상태 | TaskList, 스크래치패드 | 없음 |
| 동적 확장 | Worker 추가 spawn 가능 | 고정 subAgentCount |
| 검증 | 독립 검증 에이전트 | 없음 |

---

## 11. 피처 플래그 (44개)

### 11.1 주요 플래그 상세

| 플래그 | 상세 |
|--------|------|
| **KAIROS** | 데몬 모드. 15초 블로킹 예산, append-only 일일 로그, `<tick>` 프롬프트 라이프사이클, 독점 도구 (PushNotification, SendUserFile, SubscribePRTool), 크로스세션 연속성. 터미널 포커스 감지: unfocused="자율 행동 강화", focused="선택지 표면화, 커밋 전 확인" |
| **BUDDY** | 타마고치 스타일 펫. **18종** (duck, goose, blob, cat, dragon, octopus, owl, penguin, turtle, snail, ghost, axolotl, capybara, cactus, robot, rabbit, mushroom, chonk). **5 희귀도** (Common 60%, Uncommon 25%, Rare 10%, Epic 4%, Legendary 1%). 1% shiny 확률. **5 스탯** (DEBUGGING, PATIENCE, CHAOS, WISDOM, SNARK). 6 눈 스타일, 8 모자. Mulberry32 PRNG seed: `userId + "friend-2026-401"`. 종 이름은 내부 유출 감지 스캐너 회피 위해 **hex 코드로 인코딩** |
| **ULTRAPLAN** | CCR(Cloud Container Runtime)에서 **Opus 4.6으로 최대 30분** 계획 오프로드. 클라이언트 3초 간격 폴링, 5회 실패 허용. `__ULTRAPLAN_TELEPORT_LOCAL__` 센티널로 결과 반환. ULTRAREVIEW 변형 (코드 리뷰용) |
| **COORDINATOR_MODE** | 멀티 에이전트 스웜 오케스트레이션 |
| **VOICE_MODE** | Push-to-talk (Space 누르고 녹음, 놓으면 전송). Anthropic 독점 `voice_stream` STT. 현재 ant 전용 |
| **WEB_BROWSER_TOOL** | Playwright 기반 인프로세스 웹 브라우저 |
| **ANTI_DISTILLATION_CC** | API 요청에 가짜 도구 정의 주입 → 경쟁사 학습 데이터 오염 |
| **PROACTIVE** | 사용자 입력 없이 자율 실행 |
| **TEAMMEM** | 팀 간 메모리 공유 |
| **TOKEN_BUDGET** | 토큰 예산 관리 |
| **REACTIVE_COMPACT** | API 에러 시 긴급 압축 |
| **CACHED_MICROCOMPACT** | microcompaction 캐싱 |
| **CONTEXT_COLLAPSE** | 단계적 메시지 정리 |
| **TERMINAL_PANEL** | 터미널 캡처 |
| **NATIVE_CLIENT_ATTESTATION** | Bun/Zig 수준 해시 계산 → 바이너리 검증 |
| **AGENT_MEMORY_SNAPSHOT** | 에이전트별 영속 메모리 |

### 11.2 플래그 제어 계층

```
서버 사이드 (Anthropic API 응답, GrowthBook)
  → 조직/계정 설정
    → 로컬 설정 파일
      → 환경 변수
        → CLI 인자 (최우선)
```

---

## 12. 프롬프트 캐시 경제학

### 12.1 캐시 전략

- `SYSTEM_PROMPT_DYNAMIC_BOUNDARY`: 경계 이전 = 전역 캐시, 이후 = 동적
- 시스템 프롬프트: `cache_control: { type: 'ephemeral' }` → TTL 5분
- 도구 정의: 알파벳 순 정렬 → 결정적 캐시 히트
- `DANGEROUS_uncachedSystemPromptSection()`: 휘발성 데이터용 캐시 명시적 파괴

### 12.2 캐시 브레이크 감지

`promptCacheBreakDetection.ts`가 **14가지 캐시 파괴 벡터** 추적.
**"Sticky latches"**: 모드 토글이 캐시를 무효화하지 않도록 방지.

### 12.3 비용 구조 (Sonnet 기준)

```
일반 입력:   $3/M tokens
캐시 히트:   $0.3/M tokens (90% 할인)
캐시 쓰기:   $3.75/M tokens (25% 추가)
```

→ 3+ 턴 대화에서 캐시가 순비용 절감 효과 발생

---

## 13. 보안 메커니즘

### 13.1 Anti-Distillation (지식 증류 방지)

`ANTI_DISTILLATION_CC` 활성 시, API 요청에 `anti_distillation: ['fake_tools']` 전송.
**4가지 조건 모두 충족 시만 활성**: 컴파일타임 플래그, CLI 엔트리포인트, 1st-party API, GrowthBook 플래그.
**우회**: `CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS` 설정, 3rd-party API/SDK는 자동 우회.

### 13.2 Connector-Text Summarization

- 서버 사이드 메커니즘: 도구 호출 사이의 어시스턴트 텍스트를 버퍼링 → **암호화 서명**으로 요약 → 후속 턴에서 복원
- 전체 추론 체인을 기록된 트래픽에서 숨김
- Anthropic 내부 사용자 (`USER_TYPE === 'ant'`) 전용

### 13.3 네이티브 클라이언트 인증

- API 요청에 `cch=00000` 플레이스홀더 포함
- Bun의 Zig 네이티브 HTTP 스택이 전송 전 해시 값으로 대체
- 서버가 해시 검증 → 진짜 Claude Code 바이너리 확인
- 플레이스홀더 길이 = 해시 길이 → Content-Length 변경 방지

### 13.4 Undercover 모드 (`undercover.ts`)

- Anthropic 직원 (`USER_TYPE === 'ant'`)이 공개/오픈소스 리포에서 자동 활성
- 차단 대상: 내부 코드명 (Capybara, Tengu), Slack 채널, 미공개 모델 버전 (opus-4-7, sonnet-4-8), "Claude Code" 자체
- **강제 비활성 없음**: `CLAUDE_CODE_UNDERCOVER=1`은 강제 활성만 가능
- "Co-Authored-By" 라인 및 모든 귀속 푸터 제거
- 허용된 내부 리포: `claude-cli-internal`, `casino`, `feldspar-testing`, `forge-web`, `mobile-apps`

---

## 14. Bridge 모드

```
Local CLI (Bridge Client)
  ↕ WebSocket (JWT 인증)
Cloud (CCR - Claude Code Remote)
  ↕
AI Model (Anthropic API)
```

- 최대 32개 병렬 세션
- WebSocket 기반 실시간 양방향 통신
- JWT 토큰 인증 (갱신 자동화)
- 로컬 파일 시스템 프록시: 클라우드에서 로컬 파일 접근 가능
- 연결 끊김 시 자동 재연결 + 세션 복원
- VS Code / JetBrains IDE 통합

---

## 15. 숨겨진 기능 + 내부 특이사항

### 15.1 Frustration Detection (좌절 감지)

`userPromptKeywords.ts`에서 정규식 패턴 매칭:
- 감지 대상: "wtf", "shit(ty|tiest)?", "dumbass", "fucking? (broken|useless)"
- 대소문자 무시, 단어 경계 매칭
- Datadog 분석에 로깅
- 별도 함수로 연속 키워드 ("continue", "keep going") 감지하여 좌절과 구분

### 15.2 내부 vs 외부 사용자 차이 (`USER_TYPE === 'ant'`)

| 내부 (Anthropic 직원) | 외부 (일반 사용자) |
|----------------------|------------------|
| 완전한 문장, 결과 충실 보고, 오해 반박 | "요점으로 바로", 판단 없이 요청 실행 |
| 내부 메모: Capybara v8 **29-30% 거짓 주장율** (이전 16.7%) | - |
| `/share` → 내부 Slack `#claude-code-feedback`에 세션 게시 | - |
| `@[MODEL LAUNCH]` 태그로 외부 롤아웃 패치 표시 | - |

### 15.3 모델 코드명

| 코드명 | 대상 |
|--------|------|
| **Tengu** | Claude Code 자체 (수백 회 등장) |
| **Capybara** | Claude 4.6 변형 |
| **Fennec** | Opus 4.6 릴리스 |
| **Numbat** | 출시 전 테스트 버전 |

### 15.4 Fast Mode ("Penguin Mode")

- 엔드포인트: `/api/claude_code_penguin_mode`
- 설정 키: `penguinModeOrgEnabled`
- 킬스위치: `tengu_penguins_off`

### 15.5 TUI 렌더링 최적화

게임 엔진 기법 적용:
- `Int32Array` 기반 ASCII 풀
- 비트마스크 인코딩 스타일링
- 커서 이동 병합
- 자기 축출(self-evicting) 라인 폭 캐시 → `stringWidth` 호출 **~50배 감소**
- React + Ink + Meta Yoga flexbox
- 더블 버퍼 스크린 렌더링 + 셀 수준 블리팅

### 15.6 분석 래퍼

의도적 마찰 패턴: `AnalyticsMetadata_I_VERIFIED_THIS_IS_NOT_CODE_OR_FILEPATHS`

---

## 16. hanimo에 적용 가능한 핵심 패턴 (우선순위 순)

| # | 패턴 | 영향도 | 난이도 |
|---|------|--------|--------|
| 1 | **훅 시스템** (13 이벤트, 5 유형) — 확장성의 핵심 | 매우 높음 | 중간 |
| 2 | **메모리 시스템** (MEMORY.md + autoDream) — 세션 간 학습 | 매우 높음 | 중간 |
| 3 | **다단계 압축** (snip/micro/auto/reactive/collapse) | 높음 | 중간 |
| 4 | **도구 동시성** (safe 병렬 10개, unsafe 순차) | 높음 | 높음 |
| 5 | **glob 기반 권한 규칙** + Auto 모드 LLM 분류기 | 높음 | 중간 |
| 6 | **Coordinator 패턴** (Leader 계획/감독, Worker 실행, 독립 검증) | 매우 높음 | 매우 높음 |
| 7 | **프롬프트 캐시** (DYNAMIC_BOUNDARY, 도구 정렬) | 중간 | 낮음 |
| 8 | **세션 포크** (대화 분기 → 탐색적 개발) | 중간 | 낮음 |
| 9 | **Headless 모드** (JSON I/O → CI/CD 통합) | 중간 | 낮음 |
| 10 | **피처 플래그** (런타임 기능 토글) | 높음 | 낮음 |

---

## 참고 자료

- [DEV Community — Claude Code Leaked via npm Source Maps](https://dev.to/gabrielanhaia/claude-codes-entire-source-code-was-just-leaked-via-npm-source-maps-heres-whats-inside-cjo)
- [Alex Kim — The Claude Code Source Leak](https://alex000kim.com/posts/2026-03-31-claude-code-source-leak/)
- [o-mega.ai — Inside Claude Code: Leaked Source Analysis](https://o-mega.ai/articles/inside-claude-code-the-leaked-source-analysis)
- [Straiker — With Great Agency Comes Great Responsibility](https://www.straiker.ai/blog/claude-code-source-leak-with-great-agency-comes-great-responsibility)
- [Gizmodo — Source Code Leaks at the Exact Wrong Time](https://gizmodo.com/source-code-for-anthropics-claude-code-leaks-at-the-exact-wrong-time-2000740379)
- [VentureBeat — Claude Code's source code appears to have leaked](https://venturebeat.com/technology/claude-codes-source-code-appears-to-have-leaked-heres-what-we-know)
- [GitHub — awesome-claude-code-postleak-insights](https://github.com/nblintao/awesome-claude-code-postleak-insights)
- [MindStudio — Claude Code AutoDream](https://www.mindstudio.ai/blog/what-is-claude-code-autodream-memory-consolidation-2)
- [ccleaks.com — Hidden Features](https://www.ccleaks.com/)

*이 문서는 유출된 코드의 아키텍처 분석 목적으로 작성되었으며, 코드 자체를 포함하지 않습니다.*
