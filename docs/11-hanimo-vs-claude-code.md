# Hanimo vs Claude Code — 아키텍처 비교 분석

> **작성일**: 2026-04-01
> **기준**: hanimo v0.1.0 (~10,155 LOC) vs Claude Code v2.1.88 (~512,000 LOC)
> **참고**: `docs/09-claude-code-leak-analysis.md`, `docs/10-modol-improvement-plan.md`

---

## 1. 규모 비교

| 항목 | Hanimo | Claude Code |
|------|--------|-------------|
| **코드량** | ~10,155 LOC / 72파일 | ~512,000 LOC / 2,800+파일 |
| **언어** | TypeScript (ESM) | TypeScript + Rust (WASM) |
| **UI** | React Ink (TUI) | React Ink (TUI) |
| **LLM SDK** | Vercel AI SDK | 자체 구현 |
| **프로바이더** | 10+ (OpenAI, Anthropic, Google 등) | Anthropic 전용 |
| **배포** | npm + SEA Binary | npm + 네이티브 바이너리 |

---

## 2. 핵심 모듈 비교

| 모듈 | Hanimo | Claude Code | 분석 |
|------|--------|-------------|------|
| **에이전트 루프** | `agent-loop.ts` (224 LOC) — streamText 기반 | 자체 구현 — 512K 컨텍스트 관리, KAIROS 데몬 | Hanimo는 Vercel AI SDK 의존으로 깔끔하지만, Claude Code는 컨텍스트 관리가 훨씬 정교 |
| **툴 시스템** | 15개 툴 (file, git, shell, search 등) | 50+ 툴 + 가짜 툴(anti-distillation) | Claude Code는 decoy 툴로 모델 증류 방지. Hanimo는 실용 툴에 집중 |
| **권한 시스템** | `permission.ts` — 경로 차단 + glob 규칙 | 3단계 권한 (allow/ask/deny) + 샌드박스 | 구조적으로 유사. Hanimo도 deny > ask > allow 우선순위 |
| **세션** | JSON 파일 기반 (`~/.hanimo/sessions/`) | SQLite + 암호화 | Hanimo는 가볍지만 검색/포크 지원. Claude Code는 암호화된 영구 저장 |
| **메모리** | `memory.ts` — 마크다운 파일 기반 | autoDream — LLM 기반 자동 정리, KAIROS 데몬 백그라운드 통합 | 가장 큰 격차. Claude Code는 유휴 시에도 메모리 정리 |
| **멀티에이전트** | `orchestrator.ts` — Decompose→Execute→Synthesize | KAIROS 기반 멀티에이전트 + 전문 에이전트 풀 | Hanimo도 병렬 분해/실행/합성 구조. 패턴은 동일하나 규모 차이 |
| **MCP** | `bridge.ts` — 외부 서버 연결 + 툴 병합 | 자체 MCP 구현 + 내장 서버 | 둘 다 MCP 지원. Hanimo는 외부 의존, Claude Code는 자체 서버 포함 |
| **역할** | 4개 빌트인 (chat/dev/plan/super) + 커스텀 | 단일 에이전트 + 모드 전환 | Hanimo의 역할 시스템이 더 유연 |
| **프로바이더** | 10+ 프로바이더 (OpenAI호환 포함) | Anthropic 전용 | Hanimo가 프로바이더 다양성에서 우위 |
| **피처 플래그** | Zod 검증 + 환경변수 오버라이드 | 서버 기반 피처 플래그 + A/B 테스트 | Claude Code는 원격 피처 플래그 관리 |
| **컴팩션** | 40메시지 → LLM 요약 또는 잘라내기 | 자동 다단계 컴팩션 + KAIROS 최적화 | Claude Code가 훨씬 정교한 컨텍스트 관리 |

---

## 3. Hanimo에만 있는 기능

| 기능 | 설명 |
|------|------|
| **멀티 프로바이더** | OpenAI, Anthropic, Google, Ollama, vLLM, DeepSeek 등 10+ 프로바이더 |
| **역할 시스템** | chat/dev/plan/super + 커스텀 역할 (툴셋/프롬프트 분리) |
| **해시앵커 에디팅** | `hashline_edit` — 해시 기반 검증으로 stale edit 방지 |
| **테마 시스템** | TUI 테마 프리뷰 + 실시간 전환 |
| **리더키** | Ctrl+X 프리픽스 (vim-like 키바인딩) |
| **i18n** | 한국어/영어 시스템 메시지 전환 |
| **SEA Binary** | 단일 실행파일 빌드 (esbuild + Node SEA) |
| **비용 추적** | 실시간 토큰/비용 추정 (프로바이더별 가격표) |
| **오프라인 모드** | MCP 서버 네트워크 모드 분리 |

---

## 4. Claude Code에만 있는 기능

| 기능 | 설명 |
|------|------|
| **KAIROS 데몬** | 백그라운드 프로세스 — 컨텍스트 관리, 메모리 정리, 사전 로딩 |
| **autoDream** | 유휴 시 LLM이 메모리를 자동 정리/통합 |
| **Anti-distillation** | 가짜 툴/응답으로 모델 증류 방지 |
| **Undercover 모드** | 감지 회피 모드 (봇 탐지 우회) |
| **좌절 감지** | 사용자 감정 분석 regex (대응 전략 변경) |
| **Rust WASM** | 성능 크리티컬 부분 Rust로 구현 |
| **암호화 세션** | 세션 데이터 암호화 저장 |
| **원격 피처 플래그** | 서버 기반 A/B 테스트 + 롤아웃 |
| **텔레메트리** | 상세 사용량 분석 + 크래시 리포트 |

---

## 5. 아키텍처 패턴 일치도

### 공통 패턴 (둘 다 사용)

- React Ink TUI
- 스트리밍 텍스트 생성
- 툴 레지스트리 패턴
- MCP 통합
- 권한 게이트 (deny > ask > allow)
- 세션 영속성
- 시스템 프롬프트 빌더 (계층적)
- 프로젝트 지침 파일 (`.claude.md` / `.hanimo.md`)
- 슬래시 커맨드
- 메시지 컴팩션
- 피처 플래그
- 헤드리스 모드 (JSON I/O)

### Hanimo만의 패턴

- 멀티 프로바이더 추상화
- 역할 기반 툴 격리
- 해시앵커 편집 검증
- 비용 추적/추정

### Claude Code만의 패턴

- 백그라운드 데몬 (KAIROS)
- 메모리 자동 정리 (autoDream)
- 보안 위장 (anti-distillation, undercover)
- 감정 분석 (좌절 감지)
- Rust/WASM 하이브리드

---

## 6. 격차 분석 및 로드맵 제안

| 우선순위 | 영역 | 현재 격차 | 구현 난이도 | 기대 효과 |
|---------|------|----------|-----------|----------|
| **P0** | 컨텍스트 관리 | 단순 잘라내기 vs 다단계 요약 | 중 | 장시간 대화 품질 극대화 |
| **P0** | 메모리 자동화 | 수동 저장 vs autoDream | 중 | 에이전트 학습 능력 |
| **P1** | 세션 암호화 | 평문 JSON vs 암호화 | 하 | 보안 강화 |
| **P1** | 텔레메트리 | 없음 vs 상세 분석 | 중 | 사용자 행동 이해 |
| **P2** | 성능 최적화 | JS only vs Rust WASM | 상 | 대규모 파일 처리 속도 |
| **P2** | 원격 피처 플래그 | 로컬 only vs 서버 기반 | 중 | A/B 테스트 지원 |

---

## 7. 결론

**Hanimo의 강점**: 멀티 프로바이더, 역할 시스템, 오픈 아키텍처, 가벼운 코드베이스 (10K LOC로 핵심 기능 커버), 비용 추적

**Claude Code의 강점**: 메모리 자동화 (KAIROS/autoDream), 컨텍스트 최적화, 보안 계층 (anti-distillation), 성능 (Rust WASM)

**핵심 인사이트**: Hanimo는 Claude Code 대비 **1/50 규모**이지만 **핵심 아키텍처 패턴의 80%를 공유**. 프로바이더 다양성과 역할 시스템은 오히려 우위. 가장 큰 격차는 **메모리 자동화**와 **컨텍스트 관리 정교함** — 이 두 영역이 다음 진화 포인트.
