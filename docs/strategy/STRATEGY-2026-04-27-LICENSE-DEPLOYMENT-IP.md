# Strategy 2026-04-27 — License · Deployment · IP

> 2026-04-27 세션의 전략적 결정사항을 한 문서에 모은 것. 코드/문서 변경의
> WHY는 이 문서에서, WHAT은 `docs/policy/*` + `LICENSE` + `Makefile`에서.

---

## 1. 핵심 결정 5건 (TL;DR)

| # | 결정 | 거부한 대안 | 이유 |
|---|---|---|---|
| 1 | 라이선스 = **Apache 2.0** | MIT (기존 README 잔재), BSL/SSPL | 트레이드마크/특허 grant 명시 + Section 4 attribution이 OSS 정체성과 회사 funnel 모두 만족 |
| 2 | 데스크톱 셸 = **Wails 단독 유지**, Electron 거부 | Electron 추가, Wails+Electron 듀얼 | 외부망 99% Win10/11에 WebView2 자동배포됨. NSIS 인스톨러 + 부재 안내화면으로 1% 케이스도 잡힘. Electron은 16MB 명성 wedge 깨뜨림 |
| 3 | 폐쇄망 대응 = **분기 LTS-Onprem 빌드** (`-tags=onprem`) | Wails+Electron 듀얼, BSL 차단 | 단일 코드베이스 + 빌드 플래그. 택가이가 고통받으며 배운 폐쇄망 패치들이 그대로 onprem variant 자산. 상품은 *서비스*(인증 빌드·SLA·지원)이지 코드가 아님 |
| 4 | 트레이드마크 = **출원 보류 + 선사용 증거 누적** | 즉시 출원 (30~70만원/유사군) | 명성 임계점 전엔 도메인·SNS 핸들·블로그 활동으로 80% 커버. 임계 신호(타사 hanimo 등록, 도메인 사칭, MAU 1만, 첫 엔터프라이즈 계약) 발생 시 즉시 출원 |
| 5 | 텔레메트리 = **컴파일 타임 0**, 향후 도입 시 명시적 옵트인 only | "anonymous installation ID만 보내자" 식 절충 | hanimo의 가장 핵심적 약속. `internal/build` Onprem 분기로 컴파일 시점 보장 |

이 5개가 오늘의 *되돌릴 수 없는* 결정. 나머지는 모두 이 5개의 implementation.

---

## 2. 의사결정 트리

```
"hanimo 명성 폭발하면 어떻게 대응?"
                │
                ▼
   ┌────────── 코드는? ───────────┐
   │ Apache 2.0 (선언적 free)     │
   └──────────────┬───────────────┘
                  │
   ┌──────── 회사가 가져갈 때? ─────────┐
   │ ✅ 자유롭게 (LICENSE/NOTICE 보존)  │
   │ + SSO/감사/SLA 필요 → enterprise  │
   │   funnel = LTS-Onprem 인증 계약   │
   └──────────────┬─────────────────────┘
                  │
   ┌────── 회사가 hanimo 이름 쓰면? ────────┐
   │ ❌ Apache 2.0 §6: 트레이드마크 별도    │
   │   docs/policy/trademark-and-naming.md │
   │   nominative fair use는 OK            │
   │   포크는 별도 이름 + "Based on hanimo" │
   └──────────────┬─────────────────────────┘
                  │
   ┌──── webui 회사가 SaaS 재판매하면? ──────┐
   │ webui는 AGPL-3.0으로 커밋 예정          │
   │ (현재 hanimo-code 자체는 Apache 2.0)    │
   │ → SaaS 재판매 시 수정분 공개 의무       │
   └──────────────┬──────────────────────────┘
                  │
   ┌── 폐쇄망 회사가 사내 배포하면? ────────────┐
   │ ✅ Apache 2.0이라 자유                    │
   │ but 검증·인증·SLA가 필요한 회사들은       │
   │   분기 LTS-Onprem 인증 계약 (서비스 매출) │
   └────────────────────────────────────────────┘
```

---

## 3. 거부한 대안과 이유

### 거부 1: MIT 유지

**거부 이유**:
- LICENSE 파일이 *부재*했어서 사실상 "All Rights Reserved" 상태였음 — 라이선스 변경이 아니라 *최초 결정*임
- MIT은 트레이드마크 grant 절을 명시하지 않아 §6에 해당하는 안전망 부재
- Apache 2.0의 Section 4 (attribution 의무)가 회사 funnel에 더 강함
- SPDX `Apache-2.0` 식별자가 엔터프라이즈 컴플라이언스 도구에서 더 깔끔히 처리됨

### 거부 2: BSL/SSPL (source-available)

**거부 이유**:
- "OSS 명성" 정체성과 정면 충돌 — Elastic이 AWS와 싸우다 한 선택의 결말은 OSS 커뮤니티 비판
- 한국 시장에서 "진짜 OSS냐 아니냐" 논쟁은 wedge 약화 요인
- 우리는 SaaS 재판매 차단이 *주된* 두려움이 아님 (webui 단계의 문제)

### 거부 3: Electron 추가

**거부 이유** (택가이 커밋 로그가 보여준 패턴):
- 택가이의 Electron 전환은 **사내망 Win + WebView2 부재**라는 매우 구체적 시장 제약 때문
- hanimo 외부망 사용자엔 일반 Win10/11이라 WebView2 자동배포됨 (Microsoft Update Evergreen 2021~)
- 16MB 단일 바이너리는 hanimo 명성 wedge의 일부 — Electron은 80MB+로 늘어나서 정체성 훼손
- 1% 케이스는 NSIS 인스톨러 + 부재 안내화면으로 99.9% 커버 (택가이의 `createErrorWindow` 패턴 차용)

### 거부 4: Wails + Electron 듀얼 유지

**거부 이유**:
- 듀얼 자체는 Go-server 분리 패턴(택가이가 우연히 도달한)으로 비용 +10~15%지만,
- hanimo는 *Electron이 필요한 시장(폐쇄망)을 명시적으로 안 타깃*하므로 비용 정당화 안 됨
- 같은 효과를 **브라우저-온리 모드**로 더 가볍게 얻을 수 있음 (Go-server + 사용자 자기 브라우저)
- 진짜 폐쇄망 회사 들어오면 그제야 Electron variant 추가 (옵션으로 보존)

### 거부 5: 트레이드마크 즉시 출원

**거부 이유**:
- 비용 대비 효과 — 30~70만원/유사군은 명성 임계점 전엔 비효율
- 선사용권(unregistered prior use) + 도메인·SNS·릴리스 기록으로 80% 보호 가능
- 임계 신호(타사 사칭, MAU 1만, 첫 엔터프라이즈 계약) 발생 시 *즉시* 출원하면 충분
- 본 보류 결정은 [`../policy/trademark-and-naming.md`](../policy/trademark-and-naming.md) §6에 트리거 조건 명문화됨

### 거부 6: 즉시 텔레메트리 도입

**거부 이유**:
- "anonymous installation ID 정도면 괜찮지 않나?"라는 식의 절충은 회사들이 hanimo를 *왜* 가져가는지의 핵심 가치 훼손
- 한국 사내망/금융/공공은 *모든* 외부 통신을 의심함. 0이 아니면 검증 통과 못 함
- `-tags=onprem` 빌드의 컴파일 타임 보장으로 "절대 없다"를 *증명* 가능

---

## 4. 양립하는 두 시장의 분업

| 시장 | 채널 | 라이선스 | 빌드 | 수익 모델 |
|---|---|---|---|---|
| **개인 개발자 (외부망)** | mainline | Apache 2.0 | `make build` | 0 (명성 목적) |
| **소규모 팀 (외부망)** | mainline | Apache 2.0 | `make build-distro` (baked endpoint/model) | 0 |
| **사내망 회사 (셀프 컴파일)** | LTS-Onprem | Apache 2.0 | `make build-onprem` | 0 |
| **사내망 회사 (인증 빌드 계약)** | LTS-Onprem | Apache 2.0 + 서비스 계약 | 메인테이너 서명 빌드 | **연 단위 계약** (인증·SLA·백포트·컨설팅) |
| **호스팅 사용자 (미래)** | hanimo-webui SaaS | AGPL-3.0 (webui) | hanimo.dev 호스팅 | **월 구독** |

**핵심**: 단일 코드베이스(hanimo-code, Apache 2.0)에서 4개 채널이 분기, webui만 별도 라이선스(AGPL). TECHAI는 첫 사내망 검증 환경 + first paying customer 후보.

---

## 5. 명성 폭발 시나리오별 대응 매트릭스

| 시나리오 | 빈도 예상 | 대응 |
|---|---|---|
| 개발자 1명이 회사 일에 사용 | 매일 | ✅ 무관, 좋은 일 |
| 팀이 사내에 self-host 배포 | 주 단위 | ✅ Apache 2.0 OK. SSO/감사 필요하면 enterprise@ 안내 |
| 개인 블로그/유튜브에서 hanimo 소개 | 일상 | ✅ nominative fair use, "Powered by hanimo" 배지 환영 |
| SI 회사가 hanimo 포크 → 제품명에 "hanimo" 포함 | 명성 폭발 시 발생 | ❌ 트레이드마크 정책 위반. 시정 요청 → 14일 내 미응답 시 변호사 |
| 경쟁사가 "hanimo는 사실 사내망 안 됨" 식 FUD | 명성 폭발 시 발생 | ✅ LTS-Onprem 정책 + 택가이 검증 사례로 반박 |
| 한국 대기업이 SSO/감사/SLA 요구 | 분기 1건 가능 | ✅ LTS-Onprem 인증 계약 funnel |
| 외국 회사가 라이선스 비호환 의심 | 명성 폭발 시 발생 | ✅ SPDX `Apache-2.0` 명시 + LICENSE 파일 보존 |
| GPL 라이브러리 의존성 추가 PR | 자주 | ❌ CONTRIBUTING.md §7에 거부 명시됨 |
| 텔레메트리 추가 PR (선의의 옵트인 형태로도) | 가끔 | ⚠️ RFC 필수 (telemetry-and-privacy.md §5 강제 규칙) |

---

## 6. 즉시 액션 (코드 외, 1주 내)

| # | 액션 | 비용 | 효과 |
|---|---|---|---|
| 1 | `hanimo.dev` 도메인 등록 | 1만원/년 | 공식 채널 선점, 트레이드마크 사용 증거 |
| 2 | `hanimo.kr` 도메인 등록 | 2~3만원/년 | 한국 시장용 + 사칭 방어 |
| 3 | `@hanimo_dev` X/Twitter 핸들 | 0 | 공식 채널 |
| 4 | `hanimo` GitHub organization (현재 `flykimjiwon` 개인 → org 이전 검토) | 0 | 미래 확장성 + 기여자 추가 용이 |
| 5 | GPG 키 생성 + GitHub 등록 + `git config commit.gpgsign true` | 0 | 모든 커밋 Verified, 분쟁 시 암호학적 증거 |
| 6 | README의 `enterprise@hanimo.dev` 메일 alias 활성화 (Cloudflare Email Routing 무료) | 0 | enterprise funnel 첫 응대 채널 |

위 6개로 *명성이 폭발해도* 즉시 응대 가능한 상태가 됨.

---

## 7. 트리거 기반 후속 액션

다음 액션은 **트리거 발생 시에만** 실행:

| 트리거 | 액션 |
|---|---|
| MAU 1만 돌파 OR 외부 회사 hanimo 등록 시도 OR 도메인 사칭 발견 | 트레이드마크 출원 즉시 (한국 + 미국 우선) |
| 외부 회사 3곳이 LTS-Onprem 명시적 요청 | 첫 LTS 빌드 (`v1.0-LTS`) 발행 |
| 첫 엔터프라이즈 계약 협의 | 변호사 자문 + 표준 계약서 템플릿 작성 |
| 외부 PR 누적 30건 | DCO bot (probot/dco) 도입 검토 |
| 의존성 누적 50개 이상 | `THIRD_PARTY_NOTICES.md` 자동 생성 + SBOM (cyclonedx) 도입 |
| webui 첫 PoC | webui repo를 AGPL-3.0으로 라이선스 |
| GPL 의존성이 *진짜로* 필요한 경우 | 분리 모듈로만 가능, 핵심 라이브러리 통합 금지 (CONTRIBUTING.md §7) |

이 트리거 표는 *예언*이 아니라 *준비된 응답*. 트리거 안 와도 무관, 오면 즉시 대응.

---

## 8. 미해결 / 보류 결정

명시적으로 *나중에* 결정하기로 한 것:

1. **TECHAI ↔ hanimo Enterprise Edition 리포지셔닝**: 신한과의 계약 관계상
   "TECHAI는 hanimo-onprem의 first paying customer"로 발화 가능한지는 별도 협의 필요. 당장은 TECHAI를 거론하지 않고 hanimo-onprem만 단독 상품으로 포지셔닝.
2. **Browser-only 모드 구현 시점**: Go HTTP 서버 분리 리팩토링 1~2주 작업. Phase 20 후보. 우선순위는 사용자 결정.
3. **코드사이닝 인증서 구매**: macOS Developer ID($99/년) + Windows EV($300~/년)는 *첫 엔터프라이즈 계약 협의 시* 비용 회수 가능 시점에 구매.
4. **문서 자동 영문 번역**: 정책 문서들이 한국어 본문 + 영문 요약 1단락 형식인데, 명성이 글로벌로 확산되면 본문 영문화 필요. 현재 보류.

---

## 9. 핵심 파일 위치 (변경 시 함께 검토)

```
hanimo-code/
├── LICENSE                                          ← Apache 2.0 풀텍스트
├── NOTICE                                           ← 저작권 + 트레이드마크 안내
├── CONTRIBUTING.md                                  ← DCO 서명 + CLA grant 절
├── README.md (§Name & Marks, §Policies, §License)   ← 사용자 진입점
├── Makefile (build-onprem, build-onprem-all)        ← 빌드 채널
├── docs/policy/
│   ├── README.md                                    ← 정책 인덱스
│   ├── copyright.md
│   ├── trademark-and-naming.md
│   ├── telemetry-and-privacy.md
│   └── lts-onprem.md
├── docs/strategy/
│   └── STRATEGY-2026-04-27-LICENSE-DEPLOYMENT-IP.md ← 본 문서
├── internal/build/
│   ├── default.go                                   ← Onprem=false
│   └── onprem.go                                    ← //go:build onprem · Onprem=true
└── hanimo-code-desktop/README.md                    ← 라이선스 포인터
```

위 파일군이 *전부* 동일한 정책을 가리키도록 일관성 유지. 한 곳을 바꾸면
나머지 검토 필수.

---

## 10. 영문 한 단락 요약 (외부 인용용)

> hanimo is licensed under Apache 2.0 with telemetry zero by default,
> compile-time enforced via the `onprem` build tag for closed-network
> deployments. Trademark protection is deferred until adoption thresholds
> are reached, replaced by the standard public-prior-use evidence stack
> (domains, social handles, dated releases). The product channel is
> bipartite: a continuous mainline for individuals and small teams, and
> quarterly LTS-Onprem snapshots for enterprises that require certified
> builds, backported security patches, and SLA support. The code is and
> will remain 100% Apache 2.0; revenue stems from services (certification,
> backport SLAs, consulting), not from license tiers.

---

_Last updated: 2026-04-27 · 라이선스/정책 인프라 첫 박음 + STRATEGY 문서_
