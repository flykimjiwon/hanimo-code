# Ollama 모델 카탈로그 — hanimo 인증 모델 시스템
# Ollama Models Catalog — hanimo Certified Model System

> **대상 버전 / Target version**: hanimo v0.3.0  
> **조사 기간 / Research period**: 2024-04 ~ 2026-04  
> **최종 갱신 / Last updated**: 2026-04-11  
> **작성자 / Author**: document-specialist agent (oh-my-claudecode)

---

## 빠른 참조 요약 테이블 / Quick Reference Summary Table

| 모델 | 크기 | 컨텍스트 | 도구 호출 | HumanEval | SWE-bench | 라이선스 | Q4_K_M (GB) | Tier |
|------|------|----------|-----------|-----------|-----------|----------|-------------|------|
| **Qwen3-Coder-30B** | 30B (3.3B active) | 256K | ✅ 네이티브 | ~90 | SOTA OSS | Apache 2.0 | ~20 | 🥇 T1 |
| **Qwen2.5-Coder-32B** | 32B | 128K | ⚠️ 제한적 | 92.7 | — | Apache 2.0 | 20 | 🥇 T1 |
| **Devstral Small 2** | 24B | 128K | ✅ 네이티브 | — | 68.0 | Apache 2.0 | 14 | 🥇 T1 |
| **Gemma 4 31B** | 31B | 256K | ✅ 네이티브 | — | — | Apache 2.0 | ~20 | 🥇 T1 |
| **GPT-OSS 20B** | 20B (3.6B active) | 128K | ✅ 네이티브 | — | — | Apache 2.0 | 14 | 🥇 T1 |
| **Llama 3.3 70B** | 70B | 128K | ✅ 네이티브 | 88.4 | — | Llama 3.3 | 43 | 🥇 T1 |
| **DeepSeek-V3 671B** | 671B (37B active) | 160K | ✅ 지원 | 65.2 (base) | 66.0 (V3.1) | MIT | 404 | 🥇 T1* |
| **Qwen3 30B** | 30B (3B active) | 256K | ✅ 네이티브 | ~88 | — | Apache 2.0 | ~20 | 🥇 T1 |
| **Qwen3.5 27B** | 27B | 256K | ✅ 네이티브 | — | 72.4 | Apache 2.0 | ~18 | 🥇 T1 |
| **Devstral 24B (v1)** | 24B | 128K | ✅ 네이티브 | — | 46.8 | Apache 2.0 | 14 | 🥇 T1 |
| **Qwen2.5-Coder-14B** | 14B | 128K | ⚠️ 제한적 | ~80 | — | Apache 2.0 | 9.0 | 🥈 T2 |
| **Qwen2.5-Coder-7B** | 7B | 128K | ⚠️ 제한적 | ~72 | — | Apache 2.0 | 4.7 | 🥈 T2 |
| **Llama 3.1 8B** | 8B | 128K | ✅ 네이티브 | ~72 | — | Llama 3.1 | 5.0 | 🥈 T2 |
| **Llama 3.1 70B** | 70B | 128K | ✅ 네이티브 | ~80 | — | Llama 3.1 | 43 | 🥈 T2 |
| **Mistral Nemo 12B** | 12B | 128K | ✅ 지원 | — | — | Apache 2.0 | 7.1 | 🥈 T2 |
| **Gemma 3 27B** | 27B | 128K | ⚠️ 프롬프트 | 65.2 | — | Apache 2.0 | 16 | 🥈 T2 |
| **Phi-4 14B** | 14B | 16K | ⚠️ 제한적 | 82.6 | — | MIT | 9.1 | 🥈 T2 |
| **Hermes 3 70B** | 70B | 128K | ✅ 네이티브 | — | — | Apache 2.0* | 43 | 🥈 T2 |
| **Qwen3 8B** | 8B | 128K | ✅ 네이티브 | ~82 | — | Apache 2.0 | 5.2 | 🥈 T2 |
| **DeepSeek-R1 32B** | 32B | 128K | ⚠️ 불안정 | — | — | MIT | 20 | 🥈 T2 |
| **Command-R+ 104B** | 104B | 128K | ✅ 네이티브 | — | — | CC-BY-NC | 59 | 🥈 T2 |
| **Qwen3-Coder-480B** | 480B (35B active) | 256K | ✅ 네이티브 | — | SOTA | Apache 2.0 | ~280 | 🥈 T2* |
| **GPT-OSS 120B** | 120B (5.1B active) | 128K | ✅ 네이티브 | — | — | Apache 2.0 | 65 | 🥈 T2 |
| **Codestral 22B** | 22B | 32K | ⚠️ 지원 | — | — | MNPL (비상업) | 13 | 🥉 T3 |
| **StarCoder2 15B** | 15B | 16K | ❌ 없음 | 46.3 | — | OpenRAIL-M | 9.1 | 🥉 T3 |
| **DeepSeek-Coder-V2 16B** | 16B (MoE) | 160K | ❌ 없음 | — | — | 미확인 | 8.9 | 🥉 T3 |
| **Yi-Coder 9B** | 9B | 128K | ❌ 없음 | 85.4 | — | Apache 2.0 | 5.0 | 🥉 T3 |
| **Mixtral 8x7B** | 8x7B (12B active) | 32K | ⚠️ 지원 | — | — | Apache 2.0 | 26 | 🥉 T3 |
| **Gemma 3 12B** | 12B | 128K | ⚠️ 프롬프트 | — | — | Apache 2.0 | 7.4 | 🥉 T3 |
| **QwQ 32B** | 32B | 40K | ❌ 없음 | — | — | 미확인 | 20 | 🥉 T3 |
| **OpenThinker 32B** | 32B | 32K | ❌ 없음 | — | AIME24: 66 | 미확인 | 20 | 🥉 T3 |
| **Phi-4-reasoning 14B** | 14B | 32K | ❌ 없음 | — | — | MIT | 11 | 🥉 T3 |
| **DeepSeek-R1 671B** | 671B | 160K | ⚠️ 불안정 | — | — | MIT | 404 | 🥉 T3* |
| **CodeLlama 34B** | 34B | 16K | ❌ 없음 | — | — | 미확인 | 19 | ⚠️ 비권장 |
| **DeepSeek-Coder v1 33B** | 33B | 16K | ❌ 없음 | — | — | 미확인 | 19 | ⚠️ 비권장 |
| **CodeGemma 7B** | 7B | 8K | ❌ 없음 | — | — | Gemma ToU | 5.0 | ⚠️ 비권장 |
| **Mistral 7B** | 7B | 32K | ⚠️ raw | — | — | Apache 2.0 | 4.5 | ⚠️ 비권장 |

> *T1*: 서버급 하드웨어(400GB+) 필요. hanimo cloud 배포 전용.  
> *T2*: 480B MoE — 클라우드 또는 고성능 워크스테이션 필요.

---

## 목차 / Table of Contents

1. [연구 방법론 / Research Methodology](#연구-방법론)
2. [🥇 Tier 1 — 인증 후보 모델 (Certified Candidates)](#tier-1)
3. [🥈 Tier 2 — 주의사항 포함 지원 (Supported with Caveats)](#tier-2)
4. [🥉 Tier 3 — 실험적 / 로컬 전용 (Experimental / Local-only)](#tier-3)
5. [⚠️ 비권장 모델 (Avoid / Not Recommended)](#avoid)
6. [멀티모달 모델 개요 (Multimodal Brief)](#multimodal)
7. [추론 모델 개요 (Reasoning Models Brief)](#reasoning)
8. [hanimo v0.3.0 인증 세트 최종 권고 (Final Recommendation)](#recommendation)

---

## 연구 방법론 / Research Methodology {#연구-방법론}

### 정보 출처 / Sources

- **Ollama 라이브러리**: `ollama.com/library/{model}` — 각 모델 페이지 직접 WebFetch
- **Qwen 공식 블로그**: `qwenlm.github.io/blog/`
- **DeepSeek API 문서**: `api-docs.deepseek.com`
- **Google Gemma 블로그**: `blog.google/innovation-and-ai/technology/developers-tools/gemma-4/`
- **OpenAI GPT-OSS 발표**: `openai.com/index/introducing-gpt-oss/`
- **Mistral AI 뉴스**: `mistral.ai/news/devstral`
- **HuggingFace 모델 카드**: 각 조직의 HF 저장소
- **BenchLM.ai**: `benchlm.ai/coding` — 최신 코딩 벤치마크 리더보드

### 분류 기준 / Classification Criteria

**Tier 1 필수 조건**:
- 코딩 성능 강함 (HumanEval ≥80 또는 SWE-bench ≥40)
- 안정적인 도구 호출 (네이티브 JSON/함수 호출 지원)
- 활발히 유지보수됨 (2024년 이후 출시 또는 갱신)
- Apache 2.0 / MIT / Llama Community 라이선스
- Ollama에서 `ollama pull`로 즉시 사용 가능

**메모리 추정 공식 (Q4_K_M)**:
- 7B: ~4.5GB | 8B: ~5.0GB | 14B: ~9.0GB | 22B: ~13GB
- 24B: ~14GB | 32B: ~20GB | 70B: ~43GB | 72B: ~44GB

---

## 🥇 Tier 1 — 인증 후보 모델 {#tier-1}

> **정의**: hanimo 기본 지원 대상. 강력한 코딩 능력, 신뢰할 수 있는 도구 호출, 허용적 라이선스, 활발한 유지보수. 8~10개 목표.

---

### 1. Qwen3-Coder-30B-A3B

**핵심 요약**: Qwen 시리즈 최고 코딩 에이전트 모델. SWE-bench RL 훈련, 256K 컨텍스트, Apache 2.0.

| 항목 | 내용 |
|------|------|
| **파라미터** | 30B 총 / 3.3B 활성 (MoE) |
| **출시일** | 2025년 7월 (Qwen3-Coder 시리즈) |
| **컨텍스트** | 256K 네이티브 / 1M 외삽 |
| **도구 호출** | ✅ 네이티브 (JSON/함수 호출) |
| **HumanEval** | ~90+ (EvalPlus 기준) |
| **SWE-bench** | OSS 최고 수준 (정확한 수치 미공개) |
| **라이선스** | Apache 2.0 |
| **Q4_K_M** | ~20GB |
| **Ollama 태그** | `qwen3-coder:30b` |

**특징 / Features**:
- SWE-bench 및 유사 벤치마크에서 장기 강화학습(Long-horizon RL) 훈련
- 7.5조 토큰 훈련 (70% 코드 구성)
- 30B 파라미터 중 3.3B만 활성화 → 24GB VRAM으로 실행 가능
- 에이전트 태스크에 최적화된 아키텍처

**에이전틱 사용 시 주의사항**:
- SGLang/vLLM에서 새 도구 파서 필요 (Qwen3-Coder 전용)
- Qwen-Agent 프레임워크 사용 시 성능 극대화

**출처**: [Qwen3-Coder Ollama](https://ollama.com/library/qwen3-coder) | [QwenLM GitHub](https://github.com/QwenLM/Qwen3-Coder) | [Qwen Blog](https://qwenlm.github.io/blog/qwen3-coder/)

---

### 2. Qwen2.5-Coder-32B-Instruct

**핵심 요약**: GPT-4o 수준 코딩 능력. 128K 컨텍스트, Apache 2.0, 24GB VRAM 구동. 현재 가장 검증된 로컬 코딩 모델 중 하나.

| 항목 | 내용 |
|------|------|
| **파라미터** | 32B |
| **출시일** | 2024년 11월 |
| **컨텍스트** | 128K |
| **도구 호출** | ⚠️ 제한적 (Ollama 이슈 #7051 참조) |
| **HumanEval** | **92.7%** |
| **Aider 벤치마크** | 73.7 (GPT-4o와 유사) |
| **라이선스** | Apache 2.0 |
| **Q4_K_M** | 20GB |
| **Ollama 태그** | `qwen2.5-coder:32b` |

**특징 / Features**:
- 5.5조 토큰 훈련 (소스 코드 + 합성 데이터)
- 92개 이상 프로그래밍 언어 지원
- EvalPlus, LiveCodeBench, BigCodeBench OSS 최고 수준
- McEval 다언어 점수 65.9, MdEval 수리 벤치마크 75.2 (OSS 1위)

**에이전틱 사용 시 주의사항**:
- Qwen2.5-Coder는 **네이티브 함수 호출이 불완전**함 — `finish_reason`이 `stop`으로 반환되고 `tool_calls` 필드가 비어있는 문제 보고됨 (GitHub issue #180)
- Qwen-Agent 프레임워크나 커스텀 도구 파싱으로 우회 가능
- hanimo에서는 프롬프트 기반 도구 호출 템플릿으로 처리 권장

**출처**: [Ollama 페이지](https://ollama.com/library/qwen2.5-coder) | [Qwen 블로그](https://qwenlm.github.io/blog/qwen2.5-coder-family/) | [NVIDIA NIM 모델카드](https://build.nvidia.com/qwen/qwen2_5-coder-32b-instruct/modelcard)

---

### 3. Devstral Small 2 (24B)

**핵심 요약**: SWE-bench 68.0% — 소비자 하드웨어에서 실행 가능한 최고의 에이전트 코딩 모델.

| 항목 | 내용 |
|------|------|
| **파라미터** | 24B |
| **출시일** | 2025년 12월 (Devstral 2 시리즈) |
| **컨텍스트** | 128K (Devstral 2: 256K) |
| **도구 호출** | ✅ 네이티브 (도구 호출 성공률 클로즈드 모델 수준) |
| **SWE-bench Verified** | **68.0%** |
| **라이선스** | Apache 2.0 |
| **Q4_K_M** | 14GB |
| **Ollama 태그** | `devstral-small-2` (또는 `devstral`) |

**특징 / Features**:
- Mistral AI + All Hands AI (OpenHands) 공동 개발
- Mistral Small 3.1 기반 파인튜닝
- RTX 4090 또는 32GB Mac에서 실행 가능
- OpenHands 프레임워크와 완벽 통합
- XML 구조화 출력 + 네이티브 함수 호출 지원

**에이전틱 사용 시 주의사항**:
- 소프트웨어 엔지니어링 특화 — 일반 대화는 다른 모델에 비해 약함
- hanimo의 `--agentic` 모드에 최적

**출처**: [Ollama 페이지](https://ollama.com/library/devstral) | [Mistral 발표](https://mistral.ai/news/devstral) | [Devstral 2 발표](https://mistral.ai/news/devstral-2-vibe-cli)

---

### 4. Gemma 4 31B

**핵심 요약**: Google의 2026년 4월 최신 오픈 모델. Apache 2.0으로 라이선스 변경. LiveCodeBench 80.0%, 네이티브 함수 호출, τ²-bench 86.4%.

| 항목 | 내용 |
|------|------|
| **파라미터** | 31B (Dense) |
| **출시일** | **2026년 4월 2일** |
| **컨텍스트** | 256K |
| **도구 호출** | ✅ 네이티브 (모델 가중치에 내장) |
| **LiveCodeBench v6** | **80.0%** |
| **MMLU Pro** | 85.2 |
| **AIME 2026** | 89.2 |
| **Codeforces ELO** | 2150 |
| **τ²-bench (Retail)** | 86.4% |
| **라이선스** | **Apache 2.0** (이전 Gemma ToU에서 변경) |
| **Q4_K_M** | ~20GB |
| **Ollama 태그** | `gemma4:31b` |

**특징 / Features**:
- Arena AI 텍스트 리더보드 OSS 3위 (2026년 4월 기준)
- Gemma 3 27B (LiveCodeBench 29.1%) 대비 175% 향상
- 구조화 JSON 출력 + 네이티브 시스템 명령 지원
- 멀티모달 (텍스트 + 이미지) 지원
- 함수 호출 Gemma 3 대비 대폭 강화 (프롬프트 기반 → 가중치 내장)

**에이전틱 사용 시 주의사항**:
- 최신 모델 (2026-04-02) — Ollama 지원 안정화 필요 시 v0.3.1에서 추가 검토
- 이전 Gemma 모델들과 달리 Apache 2.0이므로 상업적 제약 없음

**출처**: [Google Gemma 4 블로그](https://blog.google/innovation-and-ai/technology/developers-tools/gemma-4/) | [Google OSS 블로그](https://opensource.googleblog.com/2026/03/gemma-4-expanding-the-gemmaverse-with-apache-20.html) | [VentureBeat](https://venturebeat.com/technology/google-releases-gemma-4-under-apache-2-0-and-that-license-change-may-matter)

---

### 5. GPT-OSS 20B

**핵심 요약**: OpenAI의 첫 오픈웨이트 모델 (GPT-2 이후). Apache 2.0, 16GB RAM으로 실행, 네이티브 도구 호출.

| 항목 | 내용 |
|------|------|
| **파라미터** | 21B 총 / 3.6B 활성 (MoE) |
| **출시일** | **2025년 8월 5일** |
| **컨텍스트** | 128K |
| **도구 호출** | ✅ 네이티브 (브라우저, Python, 함수 호출) |
| **LiveCodeBench** | **70 pass@1** (LCB v6, 높은 추론 모드) |
| **TauBench** | o3-mini 수준 |
| **라이선스** | Apache 2.0 |
| **Q4_K_M** | 14GB |
| **Ollama 태그** | `gpt-oss:20b` |

**특징 / Features**:
- MXFP4 네이티브 양자화 (4.25비트/파라미터) — 하나의 포맷으로 압축
- 단일 16GB RAM 기기에서 실행 가능
- 웹 브라우징, Python 도구 호출, 구조화 출력 내장
- Azure, HuggingFace, vLLM, Ollama 등 주요 플랫폼 동시 지원

**에이전틱 사용 시 주의사항**:
- MoE 특성상 일부 추론 프레임워크에서 배치 처리 비효율
- MXFP4 양자화 — 일부 llama.cpp 구버전에서 호환성 문제 가능

**출처**: [OpenAI GPT-OSS 발표](https://openai.com/index/introducing-gpt-oss/) | [GPT-OSS 모델 카드](https://openai.com/index/gpt-oss-model-card/) | [Ollama 페이지](https://ollama.com/library/gpt-oss)

---

### 6. Llama 3.3 70B

**핵심 요약**: Meta의 70B 인스트럭트 모델. Llama 3.1 405B와 유사한 성능, 강력한 도구 호출, HumanEval 88.4%.

| 항목 | 내용 |
|------|------|
| **파라미터** | 70B |
| **출시일** | 2024년 12월 |
| **컨텍스트** | 128K |
| **도구 호출** | ✅ 네이티브 |
| **HumanEval** | **88.4%** |
| **MBPP EvalPlus** | 87.6 |
| **BFCL v2** | 77.3 |
| **라이선스** | Llama 3.3 Community License |
| **Q4_K_M** | 43GB |
| **Ollama 태그** | `llama3.3:70b` |

**특징 / Features**:
- Llama 3.1 405B와 유사한 성능을 70B에서 달성
- 8개 언어 공식 지원 (영어, 독어, 불어, 이탈리아어, 포르투갈어, 힌디어, 스페인어, 태국어)
- 합성 데이터 생성 및 모델 증류 허용
- 도구 호출 검증 완료 (BFCL v2 기준)

**에이전틱 사용 시 주의사항**:
- 48GB+ VRAM 필요 (또는 2×24GB)
- Llama 라이선스는 700만+ MAU 서비스는 별도 계약 필요
- 한국어 공식 지원 없음 (한국어 성능은 Qwen 계열에 비해 약함)

**출처**: [Meta Llama 3.1 블로그](https://ai.meta.com/blog/meta-llama-3-1/) | [Ollama 페이지](https://ollama.com/library/llama3.3)

---

### 7. DeepSeek-V3 / V3.1 (671B)

**핵심 요약**: MIT 라이선스 OSS 최강 MoE. SWE-bench 66.0% (V3.1), 에이전트 코딩 벤치마크 40%+ 향상. 단, 400GB+ 서버 필요.

| 항목 | 내용 |
|------|------|
| **파라미터** | 671B 총 / 37B 활성 (MoE) |
| **출시일** | 2024년 12월 (V3), 2025년 8월 (V3.1) |
| **컨텍스트** | 160K |
| **도구 호출** | ✅ 지원 (V3.1: 구조화 API 지원 강화) |
| **HumanEval** | 65.2 (base) |
| **SWE-bench** | **66.0** (V3.1) |
| **라이선스** | **MIT** |
| **Q4_K_M** | 404GB |
| **Ollama 태그** | `deepseek-v3:671b` |

**특징 / Features**:
- MIT 라이선스 — 상업적 사용, 수정, 재배포 완전 허용
- V3.1: R1의 추론 능력 + V3의 범용 능력 통합
- SWE-bench V3.1 기준 Claude Opus 4.1 (74.5%)에 근접
- 에이전트 코딩 벤치마크에서 V3 대비 40%+ 향상

**에이전틱 사용 시 주의사항**:
- **로컬 실행 불가** — 서버급 하드웨어 필수 (400GB+)
- hanimo의 경우 클라우드 API (DeepSeek API) 연동 또는 전용 서버 배포 전용
- Ollama 0.5.5+ 필요

**출처**: [DeepSeek V3.1 발표](https://felloai.com/deepseek-v3-1-is-here-chinese-most-advanced-open-source-ai-yet/) | [HuggingFace](https://huggingface.co/deepseek-ai/DeepSeek-V3) | [Ollama 페이지](https://ollama.com/library/deepseek-v3)

---

### 8. Qwen3 30B (Dense)

**핵심 요약**: Qwen3 시리즈 Dense 플래그십. 256K 컨텍스트, 네이티브 도구 호출, Apache 2.0. QwQ-32B 대비 10배 효율.

| 항목 | 내용 |
|------|------|
| **파라미터** | 30B (Dense) |
| **출시일** | 2025년 4-5월 |
| **컨텍스트** | 256K |
| **도구 호출** | ✅ 네이티브 (thinking/non-thinking 모드 모두) |
| **HumanEval** | ~88 (EvalPlus 추정) |
| **라이선스** | Apache 2.0 |
| **Q4_K_M** | ~20GB |
| **Ollama 태그** | `qwen3:30b` |

**특징 / Features**:
- QwQ-32B를 10배 적은 활성 파라미터로 능가
- 사고(thinking) 모드와 비사고(non-thinking) 모드 전환 가능
- Qwen3 전체 시리즈 Apache 2.0
- 100개 이상 언어 지원

**에이전틱 사용 시 주의사항**:
- 도구 호출 시 Hermes-style 템플릿 사용 시 최고 성능
- Qwen-Agent 프레임워크와 함께 사용 권장
- thinking 모드는 컨텍스트를 빠르게 소비 — num_ctx 설정 주의

**출처**: [Qwen3 블로그](https://qwenlm.github.io/blog/qwen3/) | [Ollama 페이지](https://ollama.com/library/qwen3) | [Apache 2.0 라이선스 확인](https://ollama.com/library/qwen3:latest/blobs/d18a5cc71b84)

---

### 9. Qwen3.5 27B

**핵심 요약**: 2026년 2월 출시. SWE-bench 72.4%, BFCL-V4 72.2%, Apache 2.0. GPT-5 Mini 수준을 소비자 GPU에서 달성.

| 항목 | 내용 |
|------|------|
| **파라미터** | 27B |
| **출시일** | **2026년 2월 16일** |
| **컨텍스트** | 256K |
| **도구 호출** | ✅ 네이티브 |
| **SWE-bench Verified** | **72.4%** |
| **BFCL-V4** | 72.2 |
| **MMLU-Pro** | 87.8 |
| **라이선스** | Apache 2.0 |
| **Q4_K_M** | ~18GB |
| **Ollama 태그** | `qwen3.5:27b` |

**특징 / Features**:
- 멀티모달 네이티브 (텍스트 + 이미지)
- 201개 언어 지원
- Gated Delta Networks + sparse MoE 하이브리드 아키텍처
- GPT-5 Mini 대비 함수 호출에서 30% 우위

**에이전틱 사용 시 주의사항**:
- Qwen3.5 도구 호출 관련 Ollama 이슈 (#14493) 보고됨 — 반복 페널티 설정 필요
- 최신 모델 (2026-02) — Ollama 통합 성숙도 모니터링 필요

**출처**: [Qwen3.5 블로그](https://qwen.ai/blog?id=qwen3.5) | [HuggingFace Qwen3.5-27B](https://huggingface.co/Qwen/Qwen3.5-27B)

---

### 10. Devstral 24B (v1, Original)

**핵심 요약**: 소비자 하드웨어에서 실행 가능한 최초의 고성능 에이전트 코딩 모델. SWE-bench 46.8%.

| 항목 | 내용 |
|------|------|
| **파라미터** | 24B |
| **출시일** | 2025년 5월 |
| **컨텍스트** | 128K |
| **도구 호출** | ✅ 네이티브 |
| **SWE-bench Verified** | **46.8%** (당시 OSS 1위) |
| **라이선스** | Apache 2.0 |
| **Q4_K_M** | 14GB |
| **Ollama 태그** | `devstral:24b` |

**특징 / Features**:
- Devstral Small 2로 업그레이드 시 68.0%로 향상
- RTX 4090 단일 GPU 실행 가능
- Mistral Small 3.1 기반

**에이전틱 사용 시 주의사항**:
- Devstral Small 2 (v2)가 있으면 v2를 우선 사용
- 일반 대화나 창작 태스크에는 적합하지 않음

**출처**: [Ollama 페이지](https://ollama.com/library/devstral) | [Mistral Devstral 발표](https://mistral.ai/news/devstral)

---

## 🥈 Tier 2 — 주의사항 포함 지원 {#tier-2}

> **정의**: 작동하지만 모델별 튜닝이 필요하거나 도구 호출이 불완전한 모델. hanimo에서 `-model` 플래그로 선택 가능하며, 성능 보증은 Tier 1보다 낮음.

---

### Qwen2.5-Coder-14B / 7B

| 항목 | 14B | 7B |
|------|-----|-----|
| **컨텍스트** | 128K | 128K |
| **HumanEval** | ~80 | ~72 |
| **도구 호출** | ⚠️ 제한적 | ⚠️ 제한적 |
| **라이선스** | Apache 2.0 | Apache 2.0 |
| **Q4_K_M** | 9.0GB | 4.7GB |
| **Ollama 태그** | `qwen2.5-coder:14b` | `qwen2.5-coder:7b` |

**용도**: 메모리 제한 환경 (16GB 이하). 코딩 성능은 우수하나 도구 호출 신뢰성이 낮아 에이전트 워크플로에서 주의 필요.

**출처**: [Ollama qwen2.5-coder](https://ollama.com/library/qwen2.5-coder)

---

### Llama 3.1 8B / 70B

| 항목 | 8B | 70B |
|------|-----|-----|
| **컨텍스트** | 128K | 128K |
| **HumanEval** | ~72 | ~80 |
| **도구 호출** | ✅ 네이티브 | ✅ 네이티브 |
| **라이선스** | Llama 3.1 | Llama 3.1 |
| **Q4_K_M** | 5.0GB | 43GB |
| **Ollama 태그** | `llama3.1:8b` | `llama3.1:70b` |

**용도**: 도구 호출은 안정적이나 Llama 3.3 70B 출시 후 70B는 사실상 구버전. 8B는 경량 에이전트에 여전히 유용. Llama 라이선스: 700만+ MAU 서비스에 별도 계약 필요.

**출처**: [Ollama llama3.1](https://ollama.com/library/llama3.1) | [Meta Llama 3.1 블로그](https://ai.meta.com/blog/meta-llama-3-1/)

---

### Mistral Nemo 12B

| 항목 | 내용 |
|------|------|
| **파라미터** | 12B |
| **출시일** | 2024년 7월 18일 |
| **컨텍스트** | 128K |
| **도구 호출** | ✅ 지원 (함수 호출 훈련됨) |
| **라이선스** | Apache 2.0 |
| **Q4_K_M** | 7.1GB |
| **Ollama 태그** | `mistral-nemo:12b` |

**특징**: Mistral 7B의 드롭인 대체재. NVIDIA 협력. 11개 언어 강점. 한국어 포함. Gemma 2 9B, Llama 3 8B 대비 최대 10% 우위.

**출처**: [Ollama mistral-nemo](https://ollama.com/library/mistral-nemo) | [Mistral NeMo 발표](https://mistral.ai/news/mistral-nemo)

---

### Gemma 3 27B

| 항목 | 내용 |
|------|------|
| **파라미터** | 27B |
| **출시일** | 2025년 초 |
| **컨텍스트** | 128K |
| **도구 호출** | ⚠️ 프롬프트 엔지니어링 기반 (가중치 내장 아님) |
| **HumanEval** | 65.2 |
| **MMLU** | 78.6 |
| **라이선스** | Apache 2.0 |
| **Q4_K_M** | 16GB |
| **Ollama 태그** | `gemma3:27b` |

**특징**: Gemma 4 출시 전까지 Google 오픈 모델 최강. 멀티모달 (4B 이상). 도구 호출은 Gemma 4와 달리 시스템 프롬프트 기반.

**출처**: [Ollama gemma3](https://ollama.com/library/gemma3) | [Google AI 함수 호출 가이드](https://ai.google.dev/gemma/docs/capabilities/function-calling)

---

### Phi-4 14B (Microsoft)

| 항목 | 내용 |
|------|------|
| **파라미터** | 14B |
| **출시일** | 2024년 12월 12일 |
| **컨텍스트** | 16K ⚠️ (짧음) |
| **도구 호출** | ⚠️ Phi-4 14B는 불완전. Phi-4-mini는 지원 |
| **HumanEval** | **82.6%** |
| **라이선스** | MIT |
| **Q4_K_M** | 9.1GB |
| **Ollama 태그** | `phi4:14b` |

**특징**: MIT 라이선스, 수학/STEM 강점. 단, 컨텍스트 16K는 에이전트 워크플로에 매우 제한적. 함수 호출은 Phi-4-mini (3.8B)가 지원, 14B는 커뮤니티 구현에 의존.

**출처**: [Ollama phi4](https://ollama.com/library/phi4) | [Microsoft Phi-4 발표](https://techcommunity.microsoft.com/blog/azure-ai-foundry-blog/introducing-phi-4-microsoft%E2%80%99s-newest-small-language-model-specializing-in-comple/4357090)

---

### Hermes 3 (8B / 70B / 405B)

| 항목 | 내용 |
|------|------|
| **파라미터** | 3B, 8B, 70B, 405B |
| **출시일** | 2024년 8월 |
| **컨텍스트** | 128K |
| **도구 호출** | ✅ 네이티브 (90% 정확도) |
| **라이선스** | Apache 2.0 (데이터셋) / Llama 3.1 (모델 기반) |
| **Q4_K_M** | 8B: 5GB, 70B: 43GB |
| **Ollama 태그** | `hermes3:8b`, `hermes3:70b` |

**특징**: NousResearch가 Llama 3.1 위에 파인튜닝. 함수 호출 90% 정확도, 구조화 JSON 출력 84%. 에이전트 롤플레이, 멀티턴 대화 강점.

**출처**: [Ollama hermes3](https://ollama.com/library/hermes3) | [NousResearch Hermes 3](https://nousresearch.com/hermes3/) | [Hermes 기술 보고서](https://arxiv.org/pdf/2408.11857)

---

### Qwen3 8B

| 항목 | 내용 |
|------|------|
| **파라미터** | 8B |
| **컨텍스트** | 128K (기본 40K, 확장 128K) |
| **도구 호출** | ✅ 네이티브 |
| **라이선스** | Apache 2.0 |
| **Q4_K_M** | 5.2GB |
| **Ollama 태그** | `qwen3:8b` |

**특징**: Qwen2.5-72B-Instruct 수준의 성능을 8B에서 달성. 경량 에이전트에 탁월. thinking/non-thinking 전환 가능.

**출처**: [Ollama qwen3](https://ollama.com/library/qwen3) | [HuggingFace Qwen3-8B](https://huggingface.co/Qwen/Qwen3-8B)

---

### DeepSeek-R1 32B (Distill-Qwen)

| 항목 | 내용 |
|------|------|
| **파라미터** | 32B (Qwen2.5 기반 증류) |
| **컨텍스트** | 128K |
| **도구 호출** | ⚠️ 불안정 (Ollama 이슈 #8517) |
| **라이선스** | MIT (일부 Qwen 기반 → Apache 2.0) |
| **Q4_K_M** | 20GB |
| **Ollama 태그** | `deepseek-r1:32b` |

**특징**: OpenAI o1-mini를 여러 벤치마크에서 능가. 수학/추론 탁월. 그러나 도구 호출 템플릿 누락 문제 보고됨 — 커뮤니티 버전 (`MFDoom/deepseek-r1-tool-calling`) 사용 권장.

**출처**: [Ollama deepseek-r1](https://ollama.com/library/deepseek-r1) | [Ollama 이슈 #8517](https://github.com/ollama/ollama/issues/8517)

---

### Command-R+ 104B

| 항목 | 내용 |
|------|------|
| **파라미터** | 104B |
| **출시일** | 2024년 4월 |
| **컨텍스트** | 128K |
| **도구 호출** | ✅ 네이티브 (멀티스텝 도구 사용, RAG) |
| **라이선스** | CC-BY-NC ⚠️ (비상업적) |
| **Q4_K_M** | 59GB |
| **Ollama 태그** | `command-r-plus:104b` |

**특징**: RAG 및 엔터프라이즈 도구 사용에 최적화. GPT-4-turbo 대비 공개 도구 사용 벤치마크에서 우위. 단, CC-BY-NC 라이선스로 상업적 사용 제한.

**출처**: [Ollama command-r-plus](https://ollama.com/library/command-r-plus) | [Cohere Command R+](https://docs.cohere.com/docs/command-r-plus)

---

### Qwen3-Coder-480B-A35B (클라우드/서버 전용)

| 항목 | 내용 |
|------|------|
| **파라미터** | 480B 총 / 35B 활성 (MoE) |
| **출시일** | 2025년 7월 31일 |
| **컨텍스트** | 256K 네이티브 / 1M 외삽 |
| **도구 호출** | ✅ 네이티브 |
| **SWE-bench** | OSS 최고 (정확한 수치 미공개) |
| **라이선스** | Apache 2.0 |
| **Q4_K_M** | ~280GB |
| **Ollama 태그** | `qwen3-coder:480b` |

**특징**: Qwen3-Coder의 플래그십. 로컬 실행은 불가능하며 클라우드/전용 서버 배포 전용. 30B 버전이 로컬 대안.

**출처**: [Ollama qwen3-coder:480b](https://ollama.com/library/qwen3-coder:480b) | [HuggingFace 480B](https://huggingface.co/Qwen/Qwen3-Coder-480B-A35B-Instruct)

---

### GPT-OSS 120B

| 항목 | 내용 |
|------|------|
| **파라미터** | 117B 총 / 5.1B 활성 (MoE) |
| **출시일** | 2025년 8월 5일 |
| **컨텍스트** | 128K |
| **도구 호출** | ✅ 네이티브 |
| **TauBench** | o4-mini 수준 |
| **라이선스** | Apache 2.0 |
| **Q4_K_M** | 65GB |
| **Ollama 태그** | `gpt-oss:120b` |

**특징**: 단일 80GB GPU 실행 가능. o4-mini 수준 추론. 코딩, 수학, 건강 쿼리에서 o3-mini 능가. 도구 호출 훈련 내장.

**출처**: [Ollama gpt-oss:120b](https://ollama.com/library/gpt-oss:120b) | [OpenAI GPT-OSS](https://openai.com/index/introducing-gpt-oss/)

---

## 🥉 Tier 3 — 실험적 / 로컬 전용 {#tier-3}

> **정의**: 소규모, 틈새, 도구 호출 미약, 또는 라이선스 제한 모델. 채팅 전용이나 특수 목적으로는 유용.

---

### Codestral 22B (Mistral)

| 항목 | 내용 |
|------|------|
| **파라미터** | 22B |
| **출시일** | 2024년 5월 |
| **컨텍스트** | 32K |
| **도구 호출** | ⚠️ 지원 (신뢰성 낮음) |
| **라이선스** | ⚠️ MNPL (비상업적, 비프로덕션) |
| **Q4_K_M** | 13GB |

**특징**: 80개 이상 프로그래밍 언어, FIM(Fill-in-the-Middle) 지원, RepoBench 강점. 단 MNPL 라이선스로 상업적 프로덕션 사용 금지. hanimo가 상업 제품이라면 사용 불가.

**출처**: [Ollama codestral](https://ollama.com/library/codestral)

---

### StarCoder2 15B

| 항목 | 내용 |
|------|------|
| **파라미터** | 3B, 7B, 15B |
| **출시일** | 2024년 2월 |
| **컨텍스트** | 16K |
| **도구 호출** | ❌ 없음 |
| **HumanEval** | 46.3 (15B) |
| **라이선스** | BigCode OpenRAIL-M |
| **Q4_K_M** | 9.1GB (15B) |

**특징**: Hugging Face + ServiceNow + NVIDIA 협업. 600개 이상 언어 훈련. 그러나 컨텍스트 16K 제한, 도구 호출 없음, HumanEval 46.3%로 현재 세대 모델에 비해 크게 뒤떨어짐.

**출처**: [Ollama starcoder2](https://ollama.com/library/starcoder2) | [StarCoder2 논문](https://arxiv.org/pdf/2402.19173)

---

### DeepSeek-Coder-V2 16B (MoE)

| 항목 | 내용 |
|------|------|
| **파라미터** | 16B / 236B |
| **출시일** | 2024년 5월 |
| **컨텍스트** | 160K (16B) |
| **도구 호출** | ❌ 공식 지원 없음 |
| **라이선스** | 미확인 (연구 전용 가능성) |
| **Q4_K_M** | 8.9GB (16B) |

**특징**: 2024년 당시 GPT-4-Turbo 수준 코드 성능. 6조 토큰 훈련, 338개 프로그래밍 언어. 그러나 도구 호출 부재, 라이선스 불명확으로 Tier 3.

**출처**: [Ollama deepseek-coder-v2](https://ollama.com/library/deepseek-coder-v2)

---

### Yi-Coder 9B (01.AI)

| 항목 | 내용 |
|------|------|
| **파라미터** | 1.5B, 9B |
| **출시일** | 2024년 9월 |
| **컨텍스트** | 128K |
| **도구 호출** | ❌ 공식 미지원 |
| **HumanEval** | **85.4%** |
| **MBPP** | 73.8 |
| **라이선스** | Apache 2.0 |
| **Q4_K_M** | 5.0GB |

**특징**: 10B 미만 최고 HumanEval 중 하나 (85.4%). 52개 언어 지원. 단, 도구 호출 미지원으로 에이전트 워크플로 불가.

**출처**: [Ollama yi-coder](https://ollama.com/library/yi-coder) | [Yi-Coder GitHub](https://github.com/01-ai/Yi-Coder)

---

### Mixtral 8x7B

| 항목 | 내용 |
|------|------|
| **파라미터** | 8x7B (12B 활성) |
| **출시일** | 2023년 12월 |
| **컨텍스트** | 32K |
| **도구 호출** | ⚠️ 제한적 |
| **라이선스** | Apache 2.0 |
| **Q4_K_M** | 26GB |

**특징**: 2024년 초까지 오픈 MoE 모델 최강. 현재는 더 나은 대안이 많음. 8x22B는 네이티브 함수 호출 지원하나 64K 컨텍스트, 85GB 필요.

**출처**: [Ollama mixtral](https://ollama.com/library/mixtral)

---

### Gemma 3 12B

| 항목 | 내용 |
|------|------|
| **파라미터** | 12B |
| **컨텍스트** | 128K |
| **도구 호출** | ⚠️ 프롬프트 기반 |
| **라이선스** | Apache 2.0 |
| **Q4_K_M** | 7.4GB |

**특징**: Gemma 4 출시로 동일 메모리 예산에서 Gemma 4 4B/12B가 우위. 경과적 지원.

**출처**: [Ollama gemma3](https://ollama.com/library/gemma3)

---

### QwQ 32B

| 항목 | 내용 |
|------|------|
| **파라미터** | 32B |
| **컨텍스트** | 40K |
| **도구 호출** | ❌ 없음 |
| **라이선스** | 미확인 |
| **Q4_K_M** | 20GB |

**특징**: Qwen 추론 모델. 수학/논리 강점. DeepSeek-R1, o1-mini와 경쟁. 그러나 도구 호출 부재, 40K 컨텍스트 제한으로 에이전트 사용 불가.

**출처**: [Ollama qwq](https://ollama.com/library/qwq)

---

### OpenThinker 32B

| 항목 | 내용 |
|------|------|
| **파라미터** | 7B, 32B |
| **컨텍스트** | 32K |
| **도구 호출** | ❌ 없음 |
| **AIME24** | 66.0 (32B) |
| **MATH500** | 90.6 (32B) |
| **라이선스** | 미확인 |
| **Q4_K_M** | 20GB |

**특징**: DeepSeek-R1 증류 + OpenThoughts-114k 데이터셋. 수학 추론 탁월. 코딩 에이전트보다는 STEM 추론 전용.

**출처**: [Ollama openthinker](https://ollama.com/library/openthinker)

---

### Phi-4-reasoning 14B

| 항목 | 내용 |
|------|------|
| **파라미터** | 14B |
| **컨텍스트** | 32K |
| **도구 호출** | ❌ 없음 |
| **라이선스** | MIT |
| **Q4_K_M** | 11GB |

**특징**: SFT + RL 기반 추론 강화. DeepSeek-R1 671B와 경쟁 주장. 그러나 도구 호출 부재로 에이전트 사용 불가.

**출처**: [Ollama phi4-reasoning](https://ollama.com/library/phi4-reasoning)

---

### DeepSeek-R1 671B

| 항목 | 내용 |
|------|------|
| **파라미터** | 671B |
| **컨텍스트** | 160K |
| **도구 호출** | ⚠️ 불안정 |
| **라이선스** | MIT |
| **Q4_K_M** | 404GB |

**특징**: 수학/추론 벤치마크에서 o3/Gemini 2.5 Pro에 근접. MIT 라이선스. 단 로컬 실행 불가 (400GB+), 도구 호출 불안정. 서버 배포 시 DeepSeek-V3.1 우선 권장.

**출처**: [Ollama deepseek-r1](https://ollama.com/library/deepseek-r1)

---

## ⚠️ 비권장 모델 / Avoid / Not Recommended {#avoid}

| 모델 | 이유 |
|------|------|
| **CodeLlama (7B/13B/34B/70B)** | 2023년 모델, 도구 호출 없음, 16K 이하 컨텍스트, 현세대 모델에 완전히 뒤떨어짐 |
| **DeepSeek-Coder v1 (1.3B/6.7B/33B)** | 2023년 모델, 16K 컨텍스트, 도구 호출 없음, V2/V3로 대체됨 |
| **CodeGemma (2B/7B)** | 8K 컨텍스트, 도구 호출 없음, Gemma ToU 라이선스 제한, Gemma 3/4로 대체됨 |
| **Mistral 7B v0.1/v0.2/v0.3** | 32K 컨텍스트, 도구 호출은 raw 모드만 (v0.3), 현세대에 비해 성능 낮음 |
| **LLaVA 7B/13B/34B** | 2023년 멀티모달, 컨텍스트 2-32K, 도구 호출 없음, Llama 3.2 Vision으로 대체됨 |
| **Phi-3 Mini/Medium** | 4K-128K 컨텍스트 (설정 필요), 도구 호출 미약, Phi-4로 대체됨 |
| **Phi-3.5 Mini** | 128K 지원하나 도구 호출 없음, 3.8B로 에이전트 능력 제한 |
| **Gemma 2 (2B/9B/27B)** | 8K 컨텍스트, 도구 호출 없음, Gemma 3/4로 대체됨 |
| **Command-R 35B** | CC-BY-NC 라이선스 (비상업적), Command-R+ 104B 대비 성능 낮음 |
| **Llama 3.2 1B/3B** | 경량 모델, 코딩/에이전트에는 너무 작음 |

---

## 멀티모달 모델 개요 / Multimodal Models Brief {#multimodal}

hanimo는 현재 텍스트 전용 에이전트이지만, 향후 스크린샷 분석이나 UI 이해 기능 추가 시 참고.

| 모델 | 크기 | 컨텍스트 | 도구 호출 | 비전 강점 | 라이선스 |
|------|------|----------|-----------|-----------|----------|
| **Gemma 4 31B** | 31B | 256K | ✅ | 텍스트+이미지 | Apache 2.0 |
| **Qwen3.5 (전 크기)** | 0.8B~122B | 256K | ✅ | 텍스트+이미지+오디오 | Apache 2.0 |
| **Qwen3-VL 8B** | 8B | 256K | ✅ GUI 에이전트 | OS World 최고 | 미확인 |
| **Llama 3.2 Vision 11B** | 11B | 128K | ❌ | 이미지 추론 | Llama 3.2 |
| **Llama 4 Scout/Maverick** | 109B/400B | 10M/10M | 미확인 | 텍스트+이미지 | Llama 4 |
| **Pixtral 12B** | 12B | 128K | 미확인 | MMMU 62.5%, ChartQA 83.7% | Apache 2.0 |
| **Gemma 3 4B/12B/27B** | 다양 | 128K | ⚠️ | 텍스트+이미지 | Apache 2.0 |

**주목 모델**: Qwen3-VL은 GUI 에이전트 능력으로 OS World 벤치마크 1위. 미래 hanimo 스크린 이해 기능에 최적 후보.

**출처**: [Ollama qwen3-vl](https://ollama.com/library/qwen3-vl) | [Ollama llama3.2-vision](https://ollama.com/library/llama3.2-vision) | [Pixtral 12B 논문](https://arxiv.org/abs/2410.07073)

---

## 추론 모델 개요 / Reasoning Models Brief {#reasoning}

hanimo의 `--think` 모드나 복잡한 계획 수립에 활용 가능. 단, 대부분 도구 호출 미지원.

| 모델 | 크기 | 특징 | 도구 호출 | 권장 용도 |
|------|------|------|-----------|-----------|
| **DeepSeek-R1 671B** | 671B | MIT, o3 수준 | ⚠️ 불안정 | 서버 추론 |
| **DeepSeek-R1 32B distill** | 32B | Qwen2.5 기반, o1-mini 능가 | ⚠️ 불안정 | 로컬 추론 |
| **DeepSeek-R1 14B distill** | 14B | 경량 추론 | ⚠️ 불안정 | 16GB 추론 |
| **QwQ 32B** | 32B | Qwen 추론 전용 | ❌ | 수학/논리 |
| **OpenThinker 32B** | 32B | OSS 완전공개 데이터 | ❌ | STEM 추론 |
| **Phi-4-reasoning 14B** | 14B | MIT, SFT+RL | ❌ | 엣지 추론 |
| **Qwen3 (thinking 모드)** | 8B~235B | 도구 호출 + 추론 통합 | ✅ | **최우선 권장** |

**핵심 인사이트**: Qwen3 시리즈는 thinking 모드와 도구 호출을 동시에 지원하는 유일한 계열. 추론이 필요한 에이전트 태스크에는 Qwen3 30B 또는 Qwen3-Coder 30B가 최적.

**출처**: [DeepSeek-R1 발표](https://api-docs.deepseek.com/news/news250120) | [Ollama openthinker](https://ollama.com/library/openthinker)

---

## hanimo v0.3.0 인증 세트 최종 권고 / Final Recommendation for Certified Set {#recommendation}

### 설계 원칙 / Design Principles

hanimo의 "인증 모델 시스템"은 다음 원칙으로 운영:

1. **인증 = 테스트 완료**: 각 인증 모델은 hanimo 도구 호출 형식, 프롬프트 템플릿, 반복 제한과 함께 실제 검증됨
2. **하드웨어 계층화**: 8GB, 16GB, 24GB, 48GB+ VRAM별 최적 모델 명확히 분리
3. **라이선스 청결도**: Apache 2.0 / MIT 우선. Llama 라이선스는 수용(단 MAU 제한 주의). CC-BY-NC 제외.
4. **도구 호출 신뢰성**: Tier 1 인증 모델은 모두 hanimo의 도구 호출 워크플로와 명시적으로 테스트되어야 함

---

### 권고 인증 세트 (8~10개) / Recommended Certified Set

```
hanimo v0.3.0 공식 지원 모델 (Certified Models)
```

#### 하드웨어별 분류 / By Hardware Tier

**A. 엔트리 레벨 (8~16GB VRAM) — 로컬 기본**

| 순위 | 모델 | 태그 | 크기 | 이유 |
|------|------|------|------|------|
| 1 | Qwen3 8B | `qwen3:8b` | 5.2GB | Apache 2.0, 128K ctx, 네이티브 도구 호출, thinking 모드 |
| 2 | Devstral 24B (v1) | `devstral:24b` | 14GB | SWE-bench 46.8%, 에이전트 특화, Apache 2.0 |
| 3 | GPT-OSS 20B | `gpt-oss:20b` | 14GB | Apache 2.0, 16GB 실행, 네이티브 도구 호출 |

**B. 메인스트림 (24~32GB VRAM) — 권장 사용 환경**

| 순위 | 모델 | 태그 | 크기 | 이유 |
|------|------|------|------|------|
| 4 | Qwen3-Coder 30B | `qwen3-coder:30b` | 20GB | 에이전트 코딩 특화, 256K ctx, RL on SWE-bench |
| 5 | Qwen2.5-Coder 32B | `qwen2.5-coder:32b` | 20GB | HumanEval 92.7%, 가장 검증된 코딩 모델 |
| 6 | Devstral Small 2 | `devstral-small-2` | 14GB | SWE-bench 68.0%, 24GB GPU 실행 가능 |
| 7 | Gemma 4 31B | `gemma4:31b` | ~20GB | Apache 2.0, LiveCodeBench 80%, 네이티브 fn calling |

**C. 고성능 (48GB+ VRAM) — 서버/워크스테이션**

| 순위 | 모델 | 태그 | 크기 | 이유 |
|------|------|------|------|------|
| 8 | Llama 3.3 70B | `llama3.3:70b` | 43GB | HumanEval 88.4%, 광범위 커뮤니티 지원 |
| 9 | Qwen3.5 27B | `qwen3.5:27b` | ~18GB | SWE-bench 72.4%, 2026 최신, 멀티모달 |

**D. 클라우드 API 전용 (로컬 실행 불가)**

| 순위 | 모델 | 태그 | 크기 | 이유 |
|------|------|------|------|------|
| 10 | DeepSeek-V3.1 | API | 404GB | MIT, SWE-bench 66%, 에이전트 도구 호출 강화 |

---

### 기본 모델 (Default Model) 권고

```
hanimo 기본값: qwen3-coder:30b
```

**이유**:
- 20GB Q4_K_M → 24GB VRAM (RTX 3090/4090, Mac M2 Pro 32GB) 대응
- SWE-bench RL 훈련으로 에이전트 코딩에 최적화
- Apache 2.0 — 상업적 제약 없음
- 256K 컨텍스트로 대형 레포지토리 처리 가능
- 네이티브 도구 호출

**폴백 체인 (Fallback Chain)**:
```
qwen3-coder:30b → qwen2.5-coder:32b → devstral-small-2 → qwen3:8b → gpt-oss:20b
```

---

### 모델 태그 구성 예시 / Example Config

```go
// hanimo 인증 모델 목록 (v0.3.0)
var CertifiedModels = []Model{
    {Name: "qwen3-coder:30b",    Tier: 1, ToolCalling: Native,   VRAM: "24GB", Default: true},
    {Name: "qwen2.5-coder:32b",  Tier: 1, ToolCalling: Limited,  VRAM: "24GB"},
    {Name: "devstral-small-2",   Tier: 1, ToolCalling: Native,   VRAM: "16GB"},
    {Name: "gemma4:31b",         Tier: 1, ToolCalling: Native,   VRAM: "24GB"},
    {Name: "gpt-oss:20b",        Tier: 1, ToolCalling: Native,   VRAM: "16GB"},
    {Name: "llama3.3:70b",       Tier: 1, ToolCalling: Native,   VRAM: "48GB"},
    {Name: "qwen3.5:27b",        Tier: 1, ToolCalling: Native,   VRAM: "24GB"},
    {Name: "qwen3:8b",           Tier: 1, ToolCalling: Native,   VRAM: "8GB"},
    {Name: "deepseek-v3:671b",   Tier: 1, ToolCalling: Supported, VRAM: "Cloud"},
}
```

---

### v0.3.0 로드맵 제안 / Roadmap Suggestions

1. **즉시 지원**: `qwen3-coder:30b`, `qwen2.5-coder:32b`, `qwen3:8b`, `gpt-oss:20b`
2. **v0.3.1 추가**: `gemma4:31b` (2026-04-02 출시, Ollama 안정화 후), `devstral-small-2`
3. **v0.4.0 검토**: `qwen3.5:27b` (2026-02 출시, SWE-bench 72.4%), `Devstral 2 123B` (서버)
4. **제거 예정**: `codellama`, `deepseek-coder` v1, `codegemma` — 더 이상 경쟁력 없음

---

## 참고 문헌 / References

### Ollama 공식 페이지
- [Ollama Library](https://ollama.com/library)
- [qwen3-coder](https://ollama.com/library/qwen3-coder)
- [qwen2.5-coder](https://ollama.com/library/qwen2.5-coder)
- [qwen3](https://ollama.com/library/qwen3)
- [qwen3.5](https://ollama.com/library/qwen3.5)
- [deepseek-r1](https://ollama.com/library/deepseek-r1)
- [deepseek-v3](https://ollama.com/library/deepseek-v3)
- [devstral](https://ollama.com/library/devstral)
- [devstral-small-2](https://ollama.com/library/devstral-small-2)
- [gemma3](https://ollama.com/library/gemma3)
- [gemma4](https://ollama.com/library/gemma4) — via Ollama
- [gpt-oss](https://ollama.com/library/gpt-oss)
- [llama3.1](https://ollama.com/library/llama3.1)
- [llama3.3](https://ollama.com/library/llama3.3)
- [phi4](https://ollama.com/library/phi4)
- [hermes3](https://ollama.com/library/hermes3)
- [mistral-nemo](https://ollama.com/library/mistral-nemo)
- [command-r-plus](https://ollama.com/library/command-r-plus)
- [starcoder2](https://ollama.com/library/starcoder2)
- [yi-coder](https://ollama.com/library/yi-coder)

### 공식 발표 / Official Announcements
- [Qwen3-Coder 블로그](https://qwenlm.github.io/blog/qwen3-coder/)
- [Qwen3 블로그](https://qwenlm.github.io/blog/qwen3/)
- [Qwen2.5-Coder 블로그](https://qwenlm.github.io/blog/qwen2.5-coder-family/)
- [Qwen3.5 블로그](https://qwen.ai/blog?id=qwen3.5)
- [Google Gemma 4 블로그](https://blog.google/innovation-and-ai/technology/developers-tools/gemma-4/)
- [Google Gemma 4 Apache 2.0 발표](https://opensource.googleblog.com/2026/03/gemma-4-expanding-the-gemmaverse-with-apache-20.html)
- [OpenAI GPT-OSS 발표](https://openai.com/index/introducing-gpt-oss/)
- [Mistral Devstral 발표](https://mistral.ai/news/devstral)
- [Mistral Devstral 2 발표](https://mistral.ai/news/devstral-2-vibe-cli)
- [Mistral NeMo 발표](https://mistral.ai/news/mistral-nemo)
- [Meta Llama 3.1 블로그](https://ai.meta.com/blog/meta-llama-3-1/)
- [DeepSeek V3.1 발표](https://api-docs.deepseek.com/updates)
- [DeepSeek-R1 발표](https://api-docs.deepseek.com/news/news250120)
- [NousResearch Hermes 3](https://nousresearch.com/hermes3/)

### 벤치마크 및 분석 / Benchmarks & Analysis
- [BenchLM 코딩 리더보드](https://benchlm.ai/coding)
- [LLM Stats Qwen2.5-Coder 32B](https://llm-stats.com/models/qwen-2.5-coder-32b-instruct)
- [Qwen2.5-Coder NVIDIA NIM](https://build.nvidia.com/qwen/qwen2_5-coder-32b-instruct/modelcard)
- [Gemma 4 벤치마크](https://gemma4all.com/blog/gemma-4-benchmarks-performance)
- [Morph Ollama 모델 순위](https://www.morphllm.com/best-ollama-models)
- [StarCoder2 논문](https://arxiv.org/pdf/2402.19173)
- [Hermes 3 기술 보고서](https://arxiv.org/pdf/2408.11857)
- [Qwen3 기술 보고서](https://arxiv.org/pdf/2505.09388)

### 라이선스 / Licenses
- [Meta Llama 3.1 라이선스](https://www.llama.com/llama3_1/license/)
- [Gemma 4 Apache 2.0 전환](https://venturebeat.com/technology/google-releases-gemma-4-under-apache-2-0-and-that-license-change-may-matter)
- [DeepSeek MIT 라이선스](https://github.com/deepseek-ai/DeepSeek-R1/blob/main/LICENSE)
- [Ollama 도구 호출 모델 목록](https://ollama.com/search?c=tools)

---

*이 문서는 oh-my-claudecode document-specialist 에이전트가 WebFetch + WebSearch를 통해 ollama.com/library, 공식 모델 블로그, HuggingFace 모델 카드, 벤치마크 리더보드를 조사하여 작성했습니다.*  
*This document was researched and written by the oh-my-claudecode document-specialist agent via direct WebFetch of ollama.com/library pages, official model blogs, HuggingFace model cards, and benchmark leaderboards.*

