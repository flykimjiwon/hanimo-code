# Modol Level 2 Skill System Design

## Overview

Modol의 기존 스킬 시스템(마크다운 지식 주입)을 **워크플로우 제어가 가능한 Level 2 스킬 시스템**으로 업그레이드한다. 스킬은 체크리스트, 단계별 실행 흐름, 게이트(승인 구간)를 정의하고 에이전트가 이를 따르도록 강제한다.

### 설계 원칙

- **경량 유지**: Modol의 ~7.8K 라인 철학에 맞게, 스킬 시스템은 500줄 이내로 구현
- **agent-loop 최소 침투**: strict 게이트 처리만 agent-loop에 추가 (~15줄)
- **하위호환**: 기존 `~/.modol/skills/*.md` 순수 마크다운 파일도 계속 동작
- **사용자 확장 가능**: 빌트인 스킬과 동일한 포맷으로 커스텀 스킬 생성 가능

---

## 1. 스킬 파일 포맷

YAML frontmatter + 마크다운 단일 파일 방식.

```markdown
---
name: skill-name
description: 한 줄 설명
triggers: ["키워드1", "keyword2"]
roles: ["dev", "plan"]
type: workflow | knowledge
steps:
  - id: step-id
    title: "단계 제목"
    gate: none | soft | strict
---

## Step: step-id
이 단계에서 에이전트가 따라야 할 지시사항...
```

### 필드 정의

| 필드 | 필수 | 설명 |
|------|------|------|
| `name` | Y | 스킬 고유 식별자 (파일명과 일치 권장) |
| `description` | Y | 스킬 목적 한 줄 설명 |
| `triggers` | N | 자동 트리거 키워드 배열. 없으면 수동 호출만 가능 |
| `roles` | N | 활성화 가능 역할. 없으면 모든 역할에서 사용 가능 |
| `type` | Y | `workflow` (단계별 실행) 또는 `knowledge` (지식 주입, 하위호환) |
| `steps` | workflow만 | 단계 정의 배열 |

### Step 필드

| 필드 | 필수 | 설명 |
|------|------|------|
| `id` | Y | 단계 식별자. 마크다운 `## Step: {id}` 섹션과 매칭 |
| `title` | Y | 사용자에게 표시되는 단계 이름 |
| `gate` | Y | `none`: 자동 진행, `soft`: 프롬프트로 승인 요청 지시, `strict`: agent-loop이 코드 레벨에서 강제 중단 |

### knowledge 타입 (하위호환)

`type: knowledge`이면 steps 없이 마크다운 콘텐츠 전체가 시스템 프롬프트에 주입된다. 기존 `~/.modol/skills/*.md` 파일 중 frontmatter가 없는 파일은 자동으로 `type: knowledge`로 취급한다.

---

## 2. 모듈 구조

```
src/core/
  skills.ts              — 기존 유지 (loadSkills, buildSkillsPrompt 하위호환)
  skill-parser.ts        — frontmatter 파싱 + Step[] 추출 (신규)
  skill-engine.ts        — 트리거 매칭, 활성 스킬 관리, 게이트 판단 (신규)
  skill-prompt.ts        — 활성 스킬 → 시스템 프롬프트 섹션 생성 (신규)

src/skills/              — 빌트인 스킬 마크다운 파일 (신규)
  brainstorming.md
  debugging.md
  git-workflow.md
  planning.md
  tdd.md
  code-review.md
  refactoring.md
  security-check.md

~/.modol/skills/         — 사용자 커스텀 스킬 (기존 경로 유지)
```

### 로딩 우선순위

1. `src/skills/` — 빌트인 (항상 존재)
2. `~/.modol/skills/` — 사용자 커스텀 (빌트인과 동일 name이면 사용자 것이 우선)

---

## 3. 핵심 타입

```typescript
// skill-parser.ts

interface SkillMeta {
  name: string;
  description: string;
  triggers: string[];
  roles: string[];            // 빈 배열이면 모든 역할
  type: 'workflow' | 'knowledge';
  steps: StepDef[];           // knowledge 타입이면 빈 배열
}

interface StepDef {
  id: string;
  title: string;
  gate: 'none' | 'soft' | 'strict';
  content: string;            // 해당 step의 마크다운 지시사항
}

interface ParsedSkill {
  meta: SkillMeta;
  rawContent: string;         // frontmatter 제외 전체 마크다운
  path: string;               // 파일 경로
}
```

```typescript
// skill-engine.ts

interface ActiveSkill {
  skill: ParsedSkill;
  currentStepIndex: number;
  completedSteps: string[];   // 완료된 step id 배열
  startedAt: number;          // Date.now()
  gated: boolean;             // strict 게이트에서 대기 중인지
}

interface SkillEngine {
  // 초기화: 빌트인 + 사용자 스킬 로드
  loadAll(): ParsedSkill[];

  // 메시지에서 트리거 매칭 (현재 역할 고려)
  matchTriggers(message: string, role: string): ParsedSkill | null;

  // 스킬 활성화/비활성화
  activate(skillName: string): ActiveSkill;
  deactivate(): void;
  getActive(): ActiveSkill | null;

  // step 진행
  advanceStep(): StepDef | null;    // 다음 step 반환, 없으면 null (완료)
  isGated(): boolean;               // 현재 strict 게이트에서 대기 중?
  approveGate(): void;              // 사용자가 게이트 승인

  // 현재 활성 스킬의 프롬프트 데이터 반환 (프롬프트 생성은 skill-prompt.ts가 담당)
  getPromptData(): { skill: ParsedSkill; currentStepIndex: number; completedSteps: string[] } | null;
}
```

---

## 4. 트리거 매칭 로직

```
사용자 메시지 입력
  │
  ├─ 슬래시 커맨드인가? (/skill activate brainstorming)
  │   └─ Yes → 해당 스킬 직접 활성화
  │
  ├─ 이미 활성 스킬이 있는가?
  │   └─ Yes → 현재 스킬 계속 진행 (새 트리거 무시)
  │
  └─ triggers 키워드 매칭
      └─ 매칭 발견 → 역할 체크
          ├─ 역할 허용 → 스킬 활성화 + 알림 표시
          └─ 역할 불가 → 무시 (일반 실행)
```

### 매칭 규칙

- **대소문자 무시**: "Brainstorm" == "brainstorm"
- **부분 매칭**: 메시지에 트리거 단어가 포함되면 매칭 (단어 경계 체크)
- **복수 매칭 시**: 가장 구체적인(triggers 배열이 짧은) 스킬 우선
- **활성 스킬 중복 방지**: 이미 스킬이 활성화되어 있으면 새 트리거 무시. `/skill deactivate` 후 재활성화 가능

---

## 5. agent-loop 통합

### 수정 범위: ~15줄

agent-loop.ts의 `runAgentLoop()` 함수에 다음 로직 추가:

```typescript
// 1. 시스템 프롬프트에 활성 스킬 지시사항 주입
const promptData = skillEngine.getPromptData();
const skillPrompt = promptData ? buildSkillPrompt(promptData) : '';
const systemPrompt = baseSystemPrompt + skillPrompt;

// 2. 응답 완료 후 strict 게이트 체크
if (skillEngine.isGated()) {
  // 스트림 중단, 게이트 이벤트 발생
  onEvent?.('skill:gate', {
    skill: activeSkill.skill.meta.name,
    step: activeSkill.skill.meta.steps[activeSkill.currentStepIndex],
  });
  return { gated: true, response: fullResponse };
}

// 3. step 자동 진행 (soft 게이트는 프롬프트에서 처리)
skillEngine.advanceStep();
```

### TUI/Text-mode 통합

게이트 이벤트를 받으면:
- TUI: 상태바에 `[brainstorming: propose 단계 승인 대기]` 표시, Enter로 승인
- Text-mode: `"[GATE] brainstorming → propose 단계입니다. 계속하려면 Enter..."` 출력

---

## 6. 시스템 프롬프트 생성

### workflow 스킬 활성 시

```
## Active Skill: brainstorming (Step 2/4: clarify)

### Current Step Instructions
한 번에 하나씩 질문하여 요구사항을 명확히 하세요.
가능하면 객관식으로 제시하세요.

### Workflow Progress
- [x] explore — 프로젝트 컨텍스트 파악
- [>] clarify — 요구사항 명확화 (현재)
- [ ] propose — 접근 방식 제안 [GATE: 사용자 승인 필요]
- [ ] implement — 승인 후 구현

### Rules
- 현재 단계의 지시사항을 따르세요.
- 단계가 완료되면 "[STEP_COMPLETE: clarify]"를 응답에 포함하세요.
- [GATE] 표시된 단계 전에는 반드시 사용자 확인을 받으세요.
```

### knowledge 스킬 활성 시

기존 방식과 동일: 전체 마크다운 콘텐츠를 시스템 프롬프트에 주입.

### Step 완료 감지

에이전트 응답에서 `[STEP_COMPLETE: {step-id}]` 패턴을 감지하면 `skillEngine.advanceStep()` 호출. 감지 못 하면 수동으로 `/skill next`로 진행 가능.

---

## 7. 슬래시 커맨드

| 커맨드 | 설명 |
|--------|------|
| `/skill list` | 설치된 스킬 목록 (빌트인 + 커스텀) |
| `/skill activate <name>` | 스킬 수동 활성화 |
| `/skill deactivate` | 현재 활성 스킬 비활성화 |
| `/skill next` | 현재 step 수동 완료 + 다음 step 진행 |
| `/skill status` | 활성 스킬 진행 상황 표시 |
| `/skill create <name>` | 새 스킬 템플릿 생성 (`~/.modol/skills/`) |
| `/skill info <name>` | 스킬 상세 정보 (steps, triggers, roles) |

### TUI 통합

- `Esc` 메뉴에 `Skills` 항목 추가
- 스킬 선택 시 `SelectMenu`로 목록 표시 → 활성화
- 상태바에 활성 스킬 + 현재 step 표시

---

## 8. 빌트인 스킬 — 1차 (3개)

### 8.1 brainstorming

```yaml
name: brainstorming
description: 구현 전 설계/요구사항 탐색
triggers: ["설계해", "brainstorm", "디자인해", "기획해", "어떻게 만들"]
roles: ["dev", "plan"]
type: workflow
steps:
  - id: explore
    title: "프로젝트 컨텍스트 파악"
    gate: none
  - id: clarify
    title: "요구사항 명확화"
    gate: none
  - id: propose
    title: "접근 방식 제안"
    gate: strict
  - id: design
    title: "상세 설계"
    gate: soft
```

- **explore**: 관련 파일, 최근 커밋, 문서를 읽어 현재 상태 파악
- **clarify**: 한 번에 하나씩 질문. 객관식 선호. 목적/제약/성공 기준 확인
- **propose**: 2-3가지 접근 방식 + 트레이드오프 + 추천안 제시. **strict 게이트: 사용자가 방식을 승인해야 다음 진행**
- **design**: 승인된 방식으로 상세 설계. 아키텍처, 컴포넌트, 데이터 흐름, 에러 처리

### 8.2 debugging

```yaml
name: debugging
description: 체계적 버그 분석 및 수정
triggers: ["버그", "에러", "debug", "오류", "안돼", "안됨", "왜 안"]
roles: ["dev"]
type: workflow
steps:
  - id: reproduce
    title: "재현 확인"
    gate: none
  - id: hypothesize
    title: "원인 가설 수립"
    gate: none
  - id: isolate
    title: "원인 격리"
    gate: none
  - id: fix
    title: "수정 적용"
    gate: strict
  - id: verify
    title: "수정 검증"
    gate: none
```

- **reproduce**: 에러 메시지, 스택 트레이스, 재현 조건 확인. 재현 불가면 로그/조건 수집
- **hypothesize**: 가능한 원인 2-3개 가설. 각 가설의 검증 방법 제시
- **isolate**: 가설을 하나씩 검증. git blame, grep, 관련 코드 읽기. 근본 원인 특정
- **fix**: 수정안 제시. **strict 게이트: 사용자가 수정 방향을 승인해야 코드 변경**
- **verify**: 수정 후 테스트 실행, 회귀 확인, 관련 코드 영향 범위 체크

### 8.3 git-workflow

```yaml
name: git-workflow
description: 안전한 Git 커밋/브랜치 워크플로우
triggers: ["커밋", "commit", "브랜치", "branch", "PR", "풀리퀘"]
roles: ["dev"]
type: workflow
steps:
  - id: status
    title: "변경사항 확인"
    gate: none
  - id: stage
    title: "스테이징 선별"
    gate: soft
  - id: message
    title: "커밋 메시지 작성"
    gate: soft
  - id: commit
    title: "커밋 실행"
    gate: strict
```

- **status**: `git status`, `git diff`로 변경사항 전체 파악. 민감 파일(.env 등) 경고
- **stage**: 논리적 단위로 스테이징 제안. `git add -A` 대신 파일별 선택 권장
- **message**: conventional commit 형식으로 메시지 초안 작성. why 중심
- **commit**: **strict 게이트: 최종 메시지 + 스테이지 파일 목록을 보여주고 사용자 승인 후 커밋**

---

## 9. 빌트인 스킬 — 2차 (5개)

### 9.1 planning

```yaml
name: planning
description: 멀티스텝 구현 계획 수립
triggers: ["계획", "plan", "플랜", "로드맵", "어떤 순서"]
roles: ["dev", "plan"]
type: workflow
steps:
  - id: analyze
    title: "요구사항 분석"
    gate: none
  - id: decompose
    title: "태스크 분해"
    gate: none
  - id: sequence
    title: "실행 순서 결정"
    gate: soft
  - id: write-plan
    title: "계획 문서화"
    gate: strict
```

- **analyze**: 요구사항, 제약 조건, 의존성 파악. 기존 코드/아키텍처와의 관계 확인
- **decompose**: 작업을 독립적이고 테스트 가능한 단위로 분해. 각 단위의 예상 복잡도 표시
- **sequence**: 의존성 기반 실행 순서 결정. 병렬 가능한 태스크 식별. 크리티컬 패스 명시
- **write-plan**: **strict 게이트**: 계획서를 마크다운으로 정리하여 사용자 승인 후 저장

### 9.2 tdd

```yaml
name: tdd
description: 테스트 주도 개발 사이클
triggers: ["tdd", "테스트 먼저", "test first", "테스트부터"]
roles: ["dev"]
type: workflow
steps:
  - id: write-test
    title: "실패하는 테스트 작성"
    gate: none
  - id: run-fail
    title: "테스트 실패 확인"
    gate: none
  - id: implement
    title: "최소 구현"
    gate: strict
  - id: run-pass
    title: "테스트 통과 확인"
    gate: none
  - id: refactor
    title: "리팩토링"
    gate: soft
```

- **write-test**: 구현할 기능의 기대 동작을 테스트로 먼저 작성. 엣지 케이스 포함
- **run-fail**: 테스트 실행, 올바른 이유로 실패하는지 확인. 엉뚱한 이유면 테스트 수정
- **implement**: **strict 게이트**: 테스트를 통과할 최소한의 구현만 작성. 과도한 구현 방지
- **run-pass**: 모든 테스트 통과 확인. 기존 테스트 회귀 없음 확인
- **refactor**: 중복 제거, 네이밍 개선. 테스트가 계속 통과하는지 확인

### 9.3 code-review

```yaml
name: code-review
description: 변경사항 체계적 리뷰
triggers: ["리뷰", "review", "검토", "코드 봐줘", "확인해줘"]
roles: ["dev", "plan"]
type: workflow
steps:
  - id: diff-scan
    title: "변경 범위 파악"
    gate: none
  - id: logic-check
    title: "로직 검증"
    gate: none
  - id: security-check
    title: "보안 체크"
    gate: none
  - id: style-check
    title: "스타일/일관성"
    gate: none
  - id: summary
    title: "리뷰 요약"
    gate: strict
```

- **diff-scan**: `git diff`로 변경 범위 파악. 파일별 변경 규모, 영향 범위 정리
- **logic-check**: 비즈니스 로직 오류, 엣지 케이스 누락, 잘못된 조건문 검사
- **security-check**: 인젝션, 하드코딩된 비밀, 안전하지 않은 패턴 검사
- **style-check**: 기존 코드 컨벤션과의 일관성, 네이밍, 파일 구조
- **summary**: **strict 게이트**: 심각도별(critical/warning/info) 분류된 리뷰 결과 제시. 사용자 확인 후 수정 진행 가능

### 9.4 refactoring

```yaml
name: refactoring
description: 안전한 코드 리팩토링 가이드
triggers: ["리팩토링", "refactor", "정리", "개선해", "클린업"]
roles: ["dev"]
type: workflow
steps:
  - id: identify
    title: "개선 대상 식별"
    gate: none
  - id: safety-check
    title: "안전성 확인"
    gate: soft
  - id: extract
    title: "리팩토링 실행"
    gate: strict
  - id: verify
    title: "동작 검증"
    gate: none
```

- **identify**: 코드 스멜 식별 — 긴 함수, 중복 코드, 복잡한 조건문, 과도한 결합
- **safety-check**: 테스트 커버리지 확인. 테스트 없는 코드면 먼저 테스트 추가 제안. 영향 범위 분석
- **extract**: **strict 게이트**: 리팩토링 계획(어떤 패턴 적용, 어떤 파일 변경) 승인 후 실행. 한 번에 하나의 리팩토링만
- **verify**: 기존 테스트 전체 통과 확인. 동작 변경 없음 확인

### 9.5 security-check

```yaml
name: security-check
description: 보안 취약점 스캔
triggers: ["보안", "security", "취약점", "vulnerability", "보안 검사"]
roles: ["dev", "plan"]
type: workflow
steps:
  - id: scan-deps
    title: "의존성 취약점 스캔"
    gate: none
  - id: check-inputs
    title: "입력 검증 확인"
    gate: none
  - id: check-secrets
    title: "비밀 노출 검사"
    gate: none
  - id: report
    title: "보안 리포트"
    gate: strict
```

- **scan-deps**: `npm audit` 또는 의존성 목록에서 알려진 취약점 확인
- **check-inputs**: 사용자 입력이 들어오는 경로 추적. SQL 인젝션, XSS, 커맨드 인젝션 패턴 검사
- **check-secrets**: `.env`, API 키, 토큰이 코드에 하드코딩되지 않았는지 검사. `.gitignore` 누락 체크
- **report**: **strict 게이트**: 심각도별(critical/high/medium/low) 분류된 보안 리포트 제시. 수정 권장사항 포함

---

## 10. Step 완료 감지 전략

에이전트가 step을 완료했는지 판단하는 방법:

### 1차: 응답 내 시그널 감지

시스템 프롬프트에서 에이전트에게 `[STEP_COMPLETE: {step-id}]` 시그널을 응답 끝에 포함하도록 지시. skill-engine이 응답을 파싱하여 감지.

### 2차: 수동 진행

시그널이 감지되지 않으면 사용자가 `/skill next`로 수동 진행 가능.

### 자동 진행 금지 규칙

- **strict 게이트가 있는 step 자체는 자동 완료 불가**. `[STEP_COMPLETE]` 시그널이 감지되어도 strict step은 사용자 `approveGate()` 호출이 있어야 진행
- gate가 `none`인 step은 `[STEP_COMPLETE]` 시그널로 자동 진행
- gate가 `soft`인 step은 `[STEP_COMPLETE]` 시그널로 자동 진행 (프롬프트에서 승인 요청을 지시하지만 코드로 강제하지 않음)

---

## 11. 에러 처리

| 상황 | 처리 |
|------|------|
| frontmatter 파싱 실패 | `type: knowledge`로 폴백, 경고 로그 |
| step id와 마크다운 섹션 불일치 | 로드 시 경고, 콘텐츠 없는 step은 빈 지시사항 |
| 트리거 충돌 (2개 스킬 동시 매칭) | triggers 배열 길이가 짧은(더 구체적인) 스킬 우선 |
| strict 게이트에서 사용자가 거부 | 스킬 비활성화, 일반 모드로 복귀 |
| 스킬 활성 중 역할 변경 | 새 역할에서 허용되지 않으면 스킬 자동 비활성화 |

---

## 12. 구현 범위 추정

| 모듈 | 예상 라인 | 설명 |
|------|----------|------|
| `skill-parser.ts` | ~80줄 | frontmatter 파싱, step 추출 |
| `skill-engine.ts` | ~150줄 | 트리거 매칭, 상태 관리, 게이트 |
| `skill-prompt.ts` | ~60줄 | 시스템 프롬프트 생성 |
| `skills.ts` 수정 | ~20줄 | 새 모듈 통합, 하위호환 유지 |
| `agent-loop.ts` 수정 | ~15줄 | strict 게이트 처리 |
| `text-mode.ts` 수정 | ~30줄 | /skill 커맨드 + 게이트 UI |
| `tui/app.tsx` 수정 | ~40줄 | 스킬 메뉴 + 상태바 + 게이트 UI |
| 빌트인 스킬 8개 | ~800줄 | 마크다운 파일 (각 ~100줄) |
| 테스트 | ~200줄 | 파서, 엔진, 트리거 매칭 테스트 |
| **합계** | **~1,395줄** | |

---

## 13. 테스트 전략

### 단위 테스트

- `skill-parser.test.ts`: frontmatter 파싱, step 추출, 하위호환(frontmatter 없는 md), 잘못된 frontmatter 폴백
- `skill-engine.test.ts`: 트리거 매칭(대소문자, 부분매칭, 충돌), 스킬 활성화/비활성화, step 진행, 게이트 상태, 역할 필터링
- `skill-prompt.test.ts`: workflow 프롬프트 생성, knowledge 프롬프트 생성, 진행 상황 마크다운

### 통합 테스트

- 빌트인 스킬 파일이 모두 파싱 가능한지 검증
- agent-loop에 스킬 엔진을 연결한 상태에서 strict 게이트 동작 확인

---

## 14. 마이그레이션

기존 사용자의 `~/.modol/skills/*.md` 파일:
- frontmatter가 없으면 → `type: knowledge`로 자동 취급 (기존과 동일 동작)
- frontmatter가 있으면 → 새 스킬 시스템으로 파싱

**Breaking change 없음.**

---

## 15. 향후 확장 (Level 3 경로)

Level 2가 안정화된 후:
- 스킬이 MCP 서버를 참조할 수 있는 `mcp` frontmatter 필드 추가
- 스킬 간 체이닝: `brainstorming` 완료 시 자동으로 `planning` 활성화하는 `next-skill` 필드
- 스킬 마켓플레이스: `modol skill install <url>` 커맨드
