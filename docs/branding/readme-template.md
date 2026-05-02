# hanimo Repo README — 11-section Template

> 5 repo가 같은 구조의 README를 갖도록 통일된 템플릿. 신규 repo 또는 기존
> repo의 README 리팩토링 시 본 템플릿을 *복사 후 채워넣기*.

## 11 섹션 (순서 고정)

```
1.  Logo + Name + Korean reading
2.  Tagline (한 줄)
3.  Badges (License · Stack · Status · Stars)
4.  Why hanimo-X? (1 단락)
5.  Features (카테고리 표)
6.  Quick Start (3 옵션 — Binary / Source / Docker 중 해당)
7.  Architecture / Stack
8.  Contributing → CONTRIBUTING.md
9.  Name & Marks → docs/policy/trademark-and-naming.md
10. Policies (인덱스 표)
11. License (Apache 2.0 + Copyright + NOTICE 링크)
```

## 풀 템플릿 (Markdown)

복사해서 `<placeholder>` 부분만 채우세요.

```markdown
<p align="center">
  <img src="docs/branding/assets/hanimo-logo.svg" alt="hanimo" width="120">
</p>

<h1 align="center">hanimo-<NAME></h1>

<p align="center">
  <strong><ENGLISH_TAGLINE_ONE_LINE></strong><br>
  <em><KOREAN_TAGLINE_ONE_LINE></em>
</p>

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-f5a623?style=flat-square" alt="License"></a>
  <a href="<STACK_LINK>"><img src="https://img.shields.io/badge/<STACK>-<VERSION>-f5a623?style=flat-square" alt="Stack"></a>
  <a href="<STATUS_LINK>"><img src="https://img.shields.io/badge/Status-<STATUS>-f5a623?style=flat-square" alt="Status"></a>
  <a href="https://github.com/flykimjiwon/hanimo-<NAME>/stargazers"><img src="https://img.shields.io/github/stars/flykimjiwon/hanimo-<NAME>?style=flat-square&color=f5a623" alt="Stars"></a>
</p>

<p align="center">
  <a href="#features">Features</a> |
  <a href="#quick-start">Quick Start</a> |
  <a href="#architecture">Architecture</a> |
  <a href="#한국어">한국어</a>
</p>

---

## Why hanimo-<NAME>?

<ONE_PARAGRAPH_DIFFERENTIATOR — 무엇이 다른가, 왜 존재하는가, 누가 쓰는가>

---

## Features

<카테고리 표 — repo마다 다름. 예시:>

| Category | Items |
|---|---|
| **<CAT_1>** | <item1> · <item2> · <item3> |
| **<CAT_2>** | <item1> · <item2> |

---

## Quick Start

### Option 1: Binary / Package

\`\`\`bash
<INSTALL_COMMAND>
\`\`\`

### Option 2: Build from Source

\`\`\`bash
git clone https://github.com/flykimjiwon/hanimo-<NAME>.git
cd hanimo-<NAME>
<BUILD_COMMAND>
\`\`\`

### Option 3: Docker (해당 시)

\`\`\`bash
<DOCKER_COMMAND>
\`\`\`

---

## Architecture

<TECH_STACK_LIST>

\`\`\`
<DIRECTORY_LAYOUT>
\`\`\`

---

## Contributing

기여 환영. 라이선스 부여·DCO 서명·SPDX 헤더 컨벤션은
[`CONTRIBUTING.md`](CONTRIBUTING.md) 참고.

---

## Name & Marks

*hanimo* / *hanimo-<NAME>* 및 hanimo 로고는
**Kim Jiwon (김지원)의 미등록 트레이드마크**입니다. Apache 2.0 §6에 의해
라이선스가 트레이드마크 권리를 부여하지 않으며, 포크는 별도 이름 사용을
권장합니다. nominative fair use(블로그·비교·인용·"Powered by hanimo" 배지)는
자유.

자세한 내용: [`docs/policy/trademark-and-naming.md`](https://github.com/flykimjiwon/hanimo/blob/main/docs/policy/trademark-and-naming.md)
(canonical은 hanimo-code 레포)

---

## Policies

| Policy | Document |
|---|---|
| Copyright · 헤더 · GPG | [policy/copyright](https://github.com/flykimjiwon/hanimo/blob/main/docs/policy/copyright.md) |
| Trademark · 네이밍 | [policy/trademark-and-naming](https://github.com/flykimjiwon/hanimo/blob/main/docs/policy/trademark-and-naming.md) |
| Telemetry · 에어갭 | [policy/telemetry-and-privacy](https://github.com/flykimjiwon/hanimo/blob/main/docs/policy/telemetry-and-privacy.md) |
| LTS · 온프레미스 | [policy/lts-onprem](https://github.com/flykimjiwon/hanimo/blob/main/docs/policy/lts-onprem.md) |
| Brand · 디자인 | [branding/BRAND](https://github.com/flykimjiwon/hanimo/blob/main/docs/branding/BRAND.md) |

---

## 한국어

<KOREAN_GUIDE_PARAGRAPH — 한국어 사용자를 위한 한 단락>

\`\`\`bash
<KOREAN_QUICKSTART_COMMAND>
\`\`\`

---

## License

Copyright © 2025-2026 Kim Jiwon (김지원). All rights reserved.

hanimo-<NAME> is licensed under the
[Apache License, Version 2.0](LICENSE). See the [NOTICE](NOTICE) file for
the author declaration and 5-repo lineage.
```

## Placeholder 채우기 가이드

| Placeholder | 예시 (hanimo-rag) | 비고 |
|---|---|---|
| `<NAME>` | `rag` | 레포 이름의 hanimo- 뒤 부분 |
| `<ENGLISH_TAGLINE_ONE_LINE>` | "PostgreSQL-native Hybrid RAG Engine" | 1 줄, 60자 이내 |
| `<KOREAN_TAGLINE_ONE_LINE>` | "PostgreSQL 단독으로 동작하는 RAG 엔진" | 1 줄, 한국어 자연스럽게 |
| `<STACK>` | `Python` / `Next.js` / `Go` | 메인 스택 |
| `<VERSION>` | `3.11+` / `15` / `1.26+` | 메이저 버전 |
| `<STATUS>` | `Active` / `v0.2.0` / `Alpha` | 현재 단계 |
| `<ONE_PARAGRAPH_DIFFERENTIATOR>` | "hanimo-rag is a lightweight, pip-installable RAG..." | 80~150 단어, *다른 RAG와 무엇이 다른가*에 집중 |
| `<CAT_1>`, `<CAT_2>` | `Core`, `Search`, `Deploy` | 3~6개 카테고리 |

## 4 repo 첫 적용 시 작업 분량

| repo | 현재 README | 통일 후 변경 분량 |
|---|---|---|
| **hanimo-code** | 거의 통일됨 (이미 §License/Policies/Name & Marks 박힘) | 로고 SVG 추가 + Quick Start 정리 |
| **hanimo-webui** | 한국어 위주, 배지 통일됨 | 영문 tagline 추가 + Policies 표 추가 + 로고 |
| **hanimo-rag** | 영문 SaaS 스타일, 배지 통일됨 | 한국어 섹션 추가 + Policies 표 + 로고 |
| **hanimo-community** | 한 줄짜리 stub | 풀 템플릿 처음부터 작성 |

## 우선순위

1. **로고 v1 SVG 작성 후** 4 repo 일괄 적용 (logo-spec.md 참고)
2. 또는 *로고 placeholder*(text-only "hanimo")로 먼저 적용 후 v1 합류 시점에 swap
