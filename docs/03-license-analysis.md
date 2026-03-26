# 라이선스/저작권 분석 — 상용 포크 가능 여부

> 20개 도구의 포크, 수정, 상용화, 비공개, 특허, 트레이드마크 리스크 분석.
> **결론: 19/20 상용 포크 가능. Claude Code만 금지.**

---

## 라이선스 유형 비교

| 특성 | MIT | Apache 2.0 | AGPL-3.0 | Proprietary |
|------|:---:|:---:|:---:|:---:|
| 포크 | O | O | O | **X** |
| 클로즈드소스 가능 | O | O | **X (SaaS도)** | N/A |
| SaaS 소스 공개 의무 | X | X | **O** | N/A |
| 저작권 표시 | O | O | O | N/A |
| 변경사항 명시 | X | **O** | O | N/A |
| 명시적 특허 보호 | X | **O** | O | N/A |
| 특허 종료 조항 | X | O (소송 시) | O | N/A |
| 트레이드마크 허가 | X | X | X | N/A |
| 상용 포크 최적? | **최고** | **최고** | **위험** | **금지** |

```
가장 자유로움 ←───────────────────────────────→ 가장 제한적
MIT    Apache 2.0    LGPL    GPL-2.0    GPL-3.0    AGPL-3.0    Proprietary
```

---

## 도구별 분석 요약

### 안전 (LOW 리스크)

| # | 도구 | 라이선스 | 포크 | 상용 | 판매 | 비공개 | 특허 | 비고 |
|---|------|---------|:----:|:----:|:----:|:-----:|:----:|------|
| 1 | **OpenCode** | MIT | O | O | O | O | X | 저작권 표시만 |
| 2 | **aider** | Apache 2.0 | O | O | O | O | **O** | NOTICE 유지 |
| 3 | **Goose** | Apache 2.0 | O | O | O | O | **O** | Block 특허 보호 |
| 4 | **Cline** | Apache 2.0 | O | O | O | O | **O** | Roo Code가 상용 포크 검증 |
| 5 | **Continue** | Apache 2.0 | O | O | O | O | **O** | - |
| 6 | **Codex CLI** | Apache 2.0 | O | O | O | O | **O** | API ToS는 별도 |
| 7 | **Gemini CLI** | Apache 2.0 | O | O | O | O | **O** | Google 특허 보호 |
| 8 | **MetaGPT** | MIT | O | O | O | O | X | - |
| 9 | **CrewAI** | MIT | O | O | O | O | X | - |
| 10 | **SWE-agent** | MIT | O | O | O | O | X | 상용 포크 사례 존재 |
| 11 | **Roo Code** | Apache 2.0 | O | O | O | O | **O** | Cline 포크의 포크 OK |
| 12 | **Bubbletea** | MIT | O | O | O | O | X | 라이브러리로 사용 |
| 13 | **Ink** | MIT | O | O | O | O | X | 라이브러리로 사용 |
| 14 | **Vercel AI SDK** | Apache 2.0 | O | O | O | O | **O** | Vercel 특허 보호 |
| 15 | **Void** | Apache 2.0+MIT | O | O | O | O | **O** | AGPL 아님 (수정됨) |
| 16 | **Smolagents** | Apache 2.0 | O | O | O | O | **O** | HF 제작 |

### 주의 (MEDIUM 리스크)

| # | 도구 | 라이선스 | 리스크 | 대응 |
|---|------|---------|--------|------|
| 17 | **OpenHands** | MIT + Commercial | `enterprise/` 디렉토리 상용 라이선스 | 포크 시 `enterprise/` 제외 |
| 18 | **ChatDev** | Apache 2.0 + **CC BY-NC 4.0** | **데이터셋 상용 불가** | 코드만 포크, 데이터셋 사용 금지 |
| 19 | **LangGraph** | MIT (현재) | LangChain Inc. 향후 라이선스 변경 가능성 | 모니터링 필요 |
| 20 | **LiteLLM** | MIT + Commercial | `enterprise/` 상용 라이선스 | pip install로 코어만 사용 |

### 금지 (CRITICAL)

| # | 도구 | 라이선스 | 사유 |
|---|------|---------|------|
| **7** | **Claude Code** | **Proprietary** | npm 패키지 난독화/컴파일됨. LICENSE.md = "Anthropic Commercial ToS". **포크 = 저작권 침해.** GitHub #22002에서 오픈소스 요청 중이나 2026.03 현재 비공개. |

---

## 핵심 Q&A

### Q1: OpenCode (MIT) 포크 → 클로즈드소스 + 판매 가능?
**YES.** MIT는 원본 저작권 표시 + MIT 라이선스 텍스트를 어딘가에 포함하면 됨 (LICENSES 파일, About 화면, README). 수정사항 100% 비공개 가능.

### Q2: aider (Apache 2.0) 포크 → 필요한 것?
3가지: (1) 원본 Apache 2.0 LICENSE 파일, (2) NOTICE 파일 유지 + 변경사항 추가, (3) 수정 파일의 저작권 헤더 유지. **소스 공개 불필요.**

### Q3: AGPL-3.0 도구 → 리스크?
AGPL 코드를 수정해서 네트워크로 제공(SaaS/웹앱)하면 **전체 소스 공개 의무** 발생. 상용 제품에 치명적. 단, 이번 조사 대상 20개 중 AGPL 도구는 **없음** (Void는 Apache 2.0으로 확인됨).

### Q4: MIT + Apache 2.0 코드를 한 프로젝트에 합칠 수 있나?
**YES.** 둘 다 permissive 라이선스로 호환됨. Apache 2.0 의무(NOTICE, 변경사항 명시)는 해당 부분에만 적용.

### Q5: LiteLLM을 라이브러리로 사용 (포크 안 함)?
**YES.** `pip install litellm` → MIT 코어 적용. 저작권 표시만 포함. `enterprise/` 폴더는 별도 상용 라이선스.

### Q6: Claude Code 포크 가능?
**ABSOLUTELY NOT.** Proprietary. 난독화된 코드. "Anthropic Commercial Terms of Service" 라이선스. **포크 = 저작권 침해 + 계약 위반.**

### Q7: CLA (Contributor License Agreement) 이슈?
CLA는 **기여자**에게 영향, **사용자/포크자**에게는 무관. 선언된 라이선스가 당신이 받는 라이선스.

### Q8: 트레이드마크?
모든 오픈소스 라이선스(MIT, Apache, GPL)는 트레이드마크 권한을 **부여하지 않음**. "aider", "Cline", "OpenCode", "Claude Code" 등을 제품명으로 사용 불가. "Based on aider" 같은 사실적 설명은 OK.

---

## 상용 포크 추천 등급

### Tier 1 — 최적 (바로 포크 가능)
1. **Cline** — Apache 2.0, Roo Code가 상용 포크로 검증
2. **aider** — Apache 2.0, 명시적 특허 보호
3. **Codex CLI** — Apache 2.0, OpenAI 기여자 특허
4. **Gemini CLI** — Apache 2.0, Google 기여자 특허
5. **MetaGPT / CrewAI / SWE-agent** — MIT, 저작권 표시만

### Tier 2 — 라이브러리로 사용 (포크 불필요)
- Bubbletea, Ink — TUI 프레임워크
- LiteLLM — `pip install litellm`
- Vercel AI SDK — `npm install ai`
- LangGraph — `pip install langgraph` (라이선스 변경 모니터링)

### Tier 3 — 회피
- **Claude Code** — 금지
- **ChatDev** — 데이터셋 CC BY-NC 4.0
- **OpenHands** — enterprise/ 제외 필수

---

## 최소 법적 체크리스트

### MIT 포크 시
```
[ ] 원본 LICENSE 파일 포함
[ ] 원본 프로젝트 이름을 내 제품명으로 사용하지 않음
[ ] 어딘가에 출처 표시 (README, About, /legal)
```

### Apache 2.0 포크 시
```
[ ] 원본 LICENSE 파일 포함
[ ] NOTICE 파일 유지 + 변경사항 추가
[ ] 수정 파일의 저작권 헤더 유지
[ ] 변경사항 문서화
[ ] 원본 트레이드마크 미사용
[ ] 특허 부여 인정 포함
```

### 라이브러리 사용 시 (pip/npm)
```
[ ] 서드파티 귀속 목록 생성 (pip-licenses / license-checker)
[ ] 제품에 포함 (/legal, About, NOTICE)
[ ] FOSSA 등으로 지속적 컴플라이언스 스캔
```

---

## 면책

- 이 분석은 2026-03-26 기준 웹 리서치 결과이며 **법적 조언이 아닙니다**.
- 유의미한 매출이 있는 상용 제품은 **소프트웨어 IP 변호사** 자문을 받으세요.
- 트레이드마크는 저작권과 별도 법률 영역입니다.
- API ToS (OpenAI, Google, Anthropic)는 도구 코드 라이선스와 완전히 별개입니다.
- 항상 최신 LICENSE 파일을 확인하세요.

---

*분석 일시: 2026-03-26*
