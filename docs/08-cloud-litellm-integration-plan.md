# 08. 클라우드 프로바이더 강화 + 네이티브 고급 기능 계획

> 작성일: 2026-03-27
> 상태: 계획 단계

---

## 1. 현재 상태 분석

### 1.1 클라우드 프로바이더 연결 현황

| 프로바이더 | 지원 방식 | API Key 환경변수 | Base URL | 모델 예시 | 상태 |
|-----------|----------|-----------------|----------|----------|------|
| OpenAI (Codex 포함) | 네이티브 SDK | `OPENAI_API_KEY` | 기본값 | gpt-4o, codex-mini-latest | 완비 |
| Anthropic (Claude) | 네이티브 SDK | `ANTHROPIC_API_KEY` | 기본값 | claude-sonnet-4, claude-opus-4 | 완비 |
| Google (Gemini) | 네이티브 SDK | `GOOGLE_GENERATIVE_AI_KEY` | 기본값 | gemini-2.5-flash, gemini-2.5-pro | 완비 |
| GLM (智谱) | OpenAI-compat | `GLM_API_KEY` | open.bigmodel.cn | glm-4-plus, codegeex-4 | 완비 |
| DeepSeek | OpenAI-compat | `DEEPSEEK_API_KEY` | api.deepseek.com | deepseek-chat, deepseek-coder | 완비 |
| Groq | OpenAI-compat | `GROQ_API_KEY` | api.groq.com | qwen-qwq-32b, llama-3.3-70b | 완비 |
| Together | OpenAI-compat | `TOGETHER_API_KEY` | api.together.xyz | Qwen2.5-Coder-32B 등 | 완비 |
| OpenRouter | OpenAI-compat | `OPENROUTER_API_KEY` | openrouter.ai | 모든 모델 중계 | 완비 |
| Fireworks | OpenAI-compat | `FIREWORKS_API_KEY` | api.fireworks.ai | qwen2.5-coder-32b 등 | 완비 |
| Mistral | OpenAI-compat | `MISTRAL_API_KEY` | api.mistral.ai | codestral, mistral-large | 완비 |
| Ollama | OpenAI-compat | 불필요 | localhost:11434 | qwen3:8b, qwen3:14b 등 | 완비 |
| vLLM | OpenAI-compat | 불필요 | localhost:8000 | 사용자 모델 | 완비 |
| LM Studio | OpenAI-compat | 불필요 | localhost:1234 | 사용자 모델 | 완비 |
| Custom | OpenAI-compat | 선택 | 사용자 지정 | 사용자 모델 | 완비 |

**결론**: 14개 프로바이더 설정 자체는 모두 완비. API 키만 설정하면 즉시 사용 가능.

### 1.2 발견된 미연결/미완성 모듈 (심층 분석 결과)

| 이슈 | 심각도 | 설명 |
|------|--------|------|
| MCP 미연결 | **높음** | `McpBridge` 코드 완성이지만 `cli.ts`에서 호출 안 함 |
| Permission 승인 미연결 | **높음** | 파괴적 작업이 확인 없이 실행됨 |
| Orchestrator 미사용 | **중간** | 멀티에이전트 코드 있지만 호출 경로 없음 |
| TUI `/auto` 미연결 | **중간** | `sendAutoMessage` 미제공 |
| 토큰 사용량 이중카운팅 | **낮음** | `onStepFinish` vs `finalUsage` 충돌 |
| Role 도구 목록 불완전 | **낮음** | dev.json에 15개 중 10개만 등록 |

---

## 2. LiteLLM이란?

### 2.1 개요

[LiteLLM](https://github.com/BerriAI/litellm)은 Python 기반 LLM 통합 라이브러리/프록시 서버로:
- **100+ LLM 프로바이더**를 OpenAI 호환 API로 통합
- **프록시 서버 모드**: `litellm --model gpt-4` → localhost:4000에서 OpenAI 호환 API 제공
- **로드밸런싱, 폴백, 비용 추적, 속도 제한** 내장
- **가상 키 관리**: 팀별 API 키, 예산 설정

### 2.2 devany에 LiteLLM이 필요한가?

| 기능 | devany 현재 | LiteLLM | 비교 |
|------|------------|---------|------|
| 멀티 프로바이더 | 14개 (네이티브) | 100+ | devany도 충분, 추가 필요 시 LiteLLM |
| OpenAI 호환 | 11/14 | 전체 | 동등 |
| 로드밸런싱 | 없음 | 내장 | LiteLLM 우위 |
| 폴백 (모델 실패 시 대체) | 없음 | 내장 | LiteLLM 우위 |
| 비용 추적 | 하드코딩 가격표 | 자동 | LiteLLM 우위 |
| 속도 제한 관리 | 없음 | 내장 | LiteLLM 우위 |
| 가상 키/팀 관리 | 없음 | 내장 | LiteLLM 우위 (팀 사용 시) |
| 스트리밍 | AI SDK 네이티브 | OpenAI 호환 | 동등 |
| 의존성 | Node.js only | Python 필요 | devany 우위 |
| 설치 복잡도 | npm install | pip install + 별도 서버 | devany 우위 |

### 2.3 통합 방식 3가지

#### 방식 A: LiteLLM 프록시 서버 연결 (가장 간단)

```
[devany] ---> [LiteLLM Proxy :4000] ---> [100+ providers]
         OpenAI-compat API
```

- devany의 `custom` 프로바이더로 `http://localhost:4000/v1` 연결
- **이미 가능** — 별도 코드 변경 불필요
- LiteLLM 서버만 별도 실행 필요 (`litellm --model gpt-4 --port 4000`)

```bash
# 사용법 (현재도 가능)
devany -p custom -u http://localhost:4000/v1 -m gpt-4
```

#### 방식 B: LiteLLM을 1등급 프로바이더로 추가

```typescript
// providers/types.ts에 추가
'litellm' // 새 프로바이더

// PROVIDER_BASE_URLS에 추가
litellm: 'http://localhost:4000/v1',
```

- `devany -p litellm -m any-model` 로 사용
- 설정에서 LiteLLM URL/키 관리
- 온보딩에 LiteLLM 옵션 추가

#### 방식 C: LiteLLM의 핵심 기능을 devany에 네이티브 구현

devany가 이미 OpenAI-compat 프로토콜을 사용하므로, LiteLLM의 핵심 가치를 직접 구현:

1. **폴백 체인**: 모델 실패 시 자동 대체 모델로 전환
2. **로드밸런싱**: 같은 모델 여러 엔드포인트 분산
3. **자동 비용 추적**: API 응답의 usage 기반 실시간 계산
4. **속도 제한 관리**: 429 응답 시 자동 재시도 + 백오프

---

## 3. 구현 계획 (우선순위 순)

### Phase 1: 미연결 모듈 수정 (1-2일)

> 기존 코드가 있지만 연결만 안 된 것들. 빠르게 해결 가능.

#### 1-1. MCP 연결

**파일**: `src/cli.ts`, `src/text-mode.ts`

```
변경 사항:
1. cli.ts에서 McpBridge import
2. config 로드 후 McpBridge.loadFromConfig(config.mcp.servers) 호출
3. networkMode 전달하여 onlineOnly 서버 필터링
4. mergeToolSets(tools, mcpBridge.getAvailableTools()) 로 도구 병합
5. text-mode.ts에도 동일 적용
```

#### 1-2. Permission 승인 플로우 연결

**파일**: `src/core/agent-loop.ts`, `src/text-mode.ts`, `src/tui/app.tsx`

```
변경 사항:
1. agent-loop에서 tool 실행 전 getPermissionLevel() 호출
2. 'destructive' 또는 'dangerous' 레벨이면 사용자 확인 요청
3. text-mode: readline으로 Y/N 프롬프트
4. TUI: 확인 다이얼로그 컴포넌트 추가
```

#### 1-3. Role 도구 목록 업데이트

**파일**: `src/roles/built-in/dev.json`, `plan.json`

```
변경 사항:
- dev.json: hashline_read, hashline_edit, webfetch, todo, batch, diagnostics 추가
- plan.json: hashline_read, webfetch, batch, diagnostics 추가 (읽기 전용)
```

#### 1-4. 토큰 사용량 수정

**파일**: `src/core/agent-loop.ts`

```
변경 사항:
- onStepFinish 누적 제거, finalUsage만 사용 (또는 반대)
```

### Phase 2: 네이티브 고급 기능 (3-5일)

> LiteLLM 없이 핵심 가치를 네이티브로 구현. Python 의존성 제로 유지.

#### 3-1. 폴백 체인

**파일**: `src/core/fallback.ts` (신규)

```
설계:
- config.fallback: [{ provider, model }] 배열
- agent-loop에서 API 에러 시 다음 모델로 자동 전환
- 예: anthropic/claude-sonnet-4 실패 → google/gemini-2.5-flash → ollama/qwen3:14b

설정 예시:
{
  "provider": "anthropic",
  "model": "claude-sonnet-4-20250514",
  "fallback": [
    { "provider": "google", "model": "gemini-2.5-flash" },
    { "provider": "ollama", "model": "qwen3:14b" }
  ]
}
```

#### 3-2. 속도 제한 + 자동 재시도

**파일**: `src/core/rate-limiter.ts` (신규)

```
설계:
- 429 응답 감지 → exponential backoff (1s, 2s, 4s, max 30s)
- 프로바이더별 RPM/TPM 한도 관리
- Retry-After 헤더 존중
- 최대 3회 재시도 후 폴백 체인으로 전환
```

#### 3-3. 실시간 비용 추적 개선

**파일**: `src/core/cost-tracker.ts` (신규)

```
설계:
- 하드코딩 가격표 → JSON 외부 파일로 분리 (업데이트 용이)
- API 응답의 usage 필드 기반 실시간 계산
- 세션별 누적 비용 표시
- 일/월 예산 경고 (config.budget.daily, config.budget.monthly)
```

#### 3-4. 멀티 프로바이더 로드밸런싱

**파일**: `src/core/load-balancer.ts` (신규)

```
설계:
- 같은 모델 여러 엔드포인트 등록 가능
- Round-robin 또는 Least-latency 전략
- 헬스체크: 주기적 ping으로 죽은 엔드포인트 제외
- 주 용도: 여러 Ollama 인스턴스 분산, 또는 OpenRouter + 직접 API 병행
```

### Phase 4: Orchestrator 연결 + 멀티에이전트 활성화 (2-3일)

#### 4-1. Orchestrator 호출 경로 추가

**파일**: `src/cli.ts`, `src/text-mode.ts`

```
변경 사항:
1. config.subAgents.enabled === true 일 때 Orchestrator.run() 사용
2. /agent 명령어 또는 자동 판단으로 멀티에이전트 모드 진입
3. FileLockManager를 write 도구에 연결
```

#### 4-2. coordinator.ts / worker-pool.ts 구현

**파일**: `src/agents/coordinator.ts`, `src/agents/worker-pool.ts`

```
설계:
- Coordinator: 작업 큐 관리, 의존성 해결, 진행률 보고
- WorkerPool: 최대 N개 병렬 에이전트 실행, 자원 관리
- FileLockManager 통합으로 파일 충돌 방지
```

#### 4-3. TUI Auto-Loop 연결

**파일**: `src/tui/hooks/use-commands.ts`, `src/tui/app.tsx`

```
변경 사항:
- sendAutoMessage를 CommandContext에 제공
- /auto 명령 시 runAutoLoop 실제 호출
- TUI에서 auto-loop 진행 상태 표시
```

---

## 4. 우선순위 로드맵

```
[Phase 1] 미연결 모듈 수정 ──────────── ✅ 완료 (2026-03-27)
  ├─ ✅ MCP 연결 (cli.ts → McpBridge 초기화 + 도구 병합)
  ├─ ✅ Permission 승인 (permission-gate.ts + text-mode 래핑)
  ├─ ✅ Role 도구 목록 (dev: 10→16, plan: 6→10)
  └─ ✅ 토큰 사용량 수정 (이중 카운팅 제거)

[Phase 2] 네이티브 고급 기능 ─────────── 3-5일
  ├─ 폴백 체인
  ├─ 속도 제한 + 재시도
  ├─ 실시간 비용 추적
  └─ 로드밸런싱

[Phase 3] 멀티에이전트 활성화 ────────── 2-3일
  ├─ Orchestrator 연결
  ├─ coordinator/worker-pool 구현
  └─ TUI Auto-Loop 연결
```

**LiteLLM**: 사용하지 않기로 결정. modol은 이미 14개 프로바이더를 네이티브 지원하며,
OpenAI-compat 프로토콜로 무한 확장 가능. Python 의존성 제로 철학 유지.
필요한 고급 기능(폴백, 재시도, 비용추적)은 Phase 2에서 네이티브 구현.

**총 예상**: 5-8일 (풀타임 기준, Phase 1 완료 기준)

---

## 5. 클라우드 프로바이더 빠른 시작 가이드

현재 modol에서 각 클라우드 프로바이더를 바로 쓰는 방법:

```bash
# Claude (Anthropic)
export ANTHROPIC_API_KEY="sk-ant-..."
modol -p anthropic -m claude-sonnet-4-20250514

# Gemini (Google)
export GOOGLE_GENERATIVE_AI_KEY="AI..."
modol -p google -m gemini-2.5-flash

# Codex (OpenAI)
export OPENAI_API_KEY="sk-..."
modol -p openai -m codex-mini-latest

# GLM (智谱)
export GLM_API_KEY="..."
modol -p glm -m glm-4-plus

# DeepSeek
export DEEPSEEK_API_KEY="sk-..."
modol -p deepseek -m deepseek-chat

# 커스텀 OpenAI 호환 엔드포인트
modol -p custom -u http://my-server:8000/v1 -m my-model
```

---

## 6. 결론

- **클라우드 프로바이더**: 14개 모두 설정 완비. API 키만 넣으면 즉시 사용 가능.
- **Phase 1 완료**: MCP 연결, Permission 승인, Role 도구 목록, 토큰 수정 모두 구현됨.
- **LiteLLM 미채택**: Python 의존성 + 별도 프록시 서버 = modol의 제로 의존성 철학과 충돌. 필요 시 custom 프로바이더로 연결은 이미 가능.
- **다음 단계**: 폴백 체인, 속도 제한, 비용 추적을 네이티브로 구현하면 LiteLLM의 핵심 가치를 Python 없이 달성.
