# dev_anywhere

> **어디서든 개발하는 터미널 멀티에이전트 AI 코딩 시스템**
>
> Windows, macOS, Linux — 내부망, 외부망 — Ollama, OpenAI, Claude, Gemini, GLM 등 모든 LLM 연결

---

## 프로젝트 목표

자체 터미널 기반 멀티에이전트 AI 개발 시스템을 구축하여:

- **크로스 플랫폼**: Windows, macOS, Linux 모두 지원
- **멀티 LLM**: OpenAI, Claude, Gemini, Ollama, GLM/Zhipu, vLLM, 커스텀 엔드포인트 통합
- **내부망/외부망**: 에어갭 사내망(자체 모델 서버) + 일반 인터넷(클라우드 API) 모두 동작
- **멀티에이전트**: 복수 AI 에이전트가 협업하여 코딩, 리뷰, 테스트 수행

---

## 리서치 보고서

| # | 문서 | 내용 |
|---|------|------|
| 1 | [도구 전수 조사](docs/01-tools-survey.md) | 35+ 터미널 AI 코딩 도구, IDE 확장, 멀티에이전트 프레임워크 서베이 |
| 2 | [구축 타당성 분석](docs/02-feasibility.md) | 아키텍처, 언어 선택, 필수 컴포넌트, 타임라인, 난이도 |
| 3 | [라이선스/저작권 분석](docs/03-license-analysis.md) | 20개 도구별 포크 가능 여부, 상용화 조건, 특허, 트레이드마크 |
| 4 | [내부망/외부망 배포](docs/04-deployment.md) | 클라우드 API vs 에어갭 환경, 프로바이더 비용, 로컬 모델 성능 |

---

## 핵심 결론

### 만들 수 있는가? — YES

| 판단 | 결론 |
|------|------|
| 기술적 타당성 | **높음** — OpenCode(120K stars), aider(40K stars) 등 레퍼런스 존재 |
| 다중 LLM 통합 | **쉬움** — LiteLLM/Vercel AI SDK로 100+ 프로바이더 한 줄 통합 |
| 크로스 플랫폼 | **가능** — Go/Node.js 모두 Win+Mac+Linux 지원 |
| 내부망 배포 | **가능** — Ollama/vLLM 100% 오프라인 동작 |
| 포크 저작권 | **19/20 도구 상용 포크 가능** (Claude Code만 금지) |

### 추천 전략

```
OpenCode 포크 (Go, MIT)
├── TUI + LSP + SQLite 세션 — 이미 완성됨
├── 75+ 모델 지원 구조
├── + 커스텀 프로바이더 추가 (GLM, 자체 엔드포인트)
├── + 멀티에이전트 오케스트레이션 레이어
├── + MCP 서버 통합
└── + providers.yaml 환경 전환 (내부망 ↔ 외부망)
```

### MVP 타임라인

| 전략 | 기간 |
|------|------|
| OpenCode 포크 + 확장 | ~5주 |
| 처음부터 구축 (TypeScript) | ~7주 MVP / ~15주 Full |

---

## 기술 스택 후보

| 선택지 | 점수 | TUI | LLM 에코시스템 |
|--------|------|-----|---------------|
| **Go** (Bubbletea) | 7.6/10 | **최고** (OpenCode 검증) | 보통 |
| **TypeScript** (Ink) | 8.2/10 | 좋음 | **최고** (Vercel AI SDK) |
| **Python** (Textual) | 7.7/10 | 보통 | 좋음 (LiteLLM 100+) |
| **Rust** (Ratatui) | 6.8/10 | 좋음 | 부족 |

---

## LLM 프로바이더 맵

```
┌──────────────────────────────────────────────────┐
│           통합 LLM 추상화 레이어                    │
│    (LiteLLM / Vercel AI SDK / 직접 구현)           │
├─────────┬────────┬────────┬────────┬──────┬──────┤
│ OpenAI  │ Claude │ Gemini │  GLM   │Ollama│ vLLM │
│ (직접)  │ (SDK)  │ (SDK)  │(OAI호환)│(로컬)│(로컬) │
└─────────┴────────┴────────┴────────┴──────┴──────┘

외부망: api.openai.com / api.anthropic.com / gemini API
내부망: internal-vllm:8000 / internal-ollama:11434
전환:   OPENAI_BASE_URL 환경변수 하나로
```

---

## 도구 인기순 (GitHub Stars, 2026.03)

| # | 도구 | Stars | 언어 | 라이선스 | 포크 안전 |
|---|------|-------|------|---------|----------|
| 1 | OpenCode | ~120K | Go | MIT | **안전** |
| 2 | Claude Code | ~82K | TS | Proprietary | **금지** |
| 3 | OpenHands | ~69K | Python | MIT* | 주의 |
| 4 | Codex CLI | ~67K | Rust | Apache 2.0 | **안전** |
| 5 | Cline | ~59K | TS | Apache 2.0 | **안전** |
| 6 | Zed | ~55K | Rust | GPL/AGPL | 주의 |
| 7 | AutoGen | ~55K | Python | MIT | **안전** |
| 8 | GPT-Engineer | ~54K | Python | MIT | **안전** |
| 9 | MetaGPT | ~47K | Python | MIT | **안전** |
| 10 | CrewAI | ~45K | Python | MIT | **안전** |

---

## 배포 시나리오

| 시나리오 | 스택 | 비용 | 프라이버시 |
|---------|------|------|----------|
| 소규모 팀, 인터넷 OK | Gemini Flash + 포크 도구 | $0.40/M | API=학습 안 함 |
| 최고 품질 | Claude Sonnet 4 via Bedrock | $15/M | VPC 격리 |
| 한국 엔터프라이즈 | Azure OpenAI (Korea Central) | 엔터프라이즈 | VPC 격리 |
| 에어갭 프로덕션 | vLLM + DeepSeek R1 70B | HW만 | 100% 온프레미스 |
| 에어갭 개인 개발 | Ollama + DeepSeek-Coder | 무료 | 100% 로컬 |

---

## 프로젝트 구조

```
dev_anywhere/
├── README.md                    # 이 파일
├── docs/
│   ├── 01-tools-survey.md       # 35+ 도구 전수 조사
│   ├── 02-feasibility.md        # 구축 타당성 분석
│   ├── 03-license-analysis.md   # 라이선스/저작권 분석
│   └── 04-deployment.md         # 내부망/외부망 배포
└── (향후 소스 코드)
```

---

## 조사 날짜

2026-03-26 기준. 모든 데이터(Stars, 가격, 라이선스)는 이 시점의 스냅샷입니다.

---

## License

MIT
