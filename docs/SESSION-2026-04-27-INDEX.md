# Session Index — 2026-04-27 시간순 산출

> 이 문서는 2026-04-27 세션의 시간순 흐름·논의 주제·산출 매핑.
> 30초 요약은 [`SESSION-2026-04-27-CLOSING.md`](SESSION-2026-04-27-CLOSING.md),
> 결정 근거는 [`strategy/STRATEGY-2026-04-27-LICENSE-DEPLOYMENT-IP.md`](strategy/STRATEGY-2026-04-27-LICENSE-DEPLOYMENT-IP.md).

---

## 세션 시작 컨텍스트

- **시작 상태**: `origin/main = bb81cf0` · `M landing-mockups` (이전 세션 잔재)
- **최근 작업**: 어제(04-25) v0.2.0 릴리스 + Phase 19(한국 MCP 35종) 마감
- **사용자 첫 질문**: "hanimo code, ide 이쪽이랑 택가이코드쪽 tui, ide 포팅해올거없나?"

---

## 논의 흐름 (시간순)

### Phase 1 — 택가이 → hanimo 포팅 후보 탐색

**질문**: 택가이코드 TUI/IDE에서 hanimo로 가져올만한 게 뭔지

**탐색 대상**:
- `TECHAI_CODE/techai-ide/` (Wails 16MB 자랑하던 풀스택)
- `TECHAI_CODE/techai-ide-electron/` ← *추가 분기 발견*
- `SESSION_SUMMARY.md`, `DEV_GUIDE.md`, `WORKSPACE_GUIDE.md`
- `docs/PORTING_LOG.md`, `ARCHITECTURE_NEXT.md`, `DEBUG_TRANSPORT_FREEZE.md`, `SESSION_LOG_v0.3.0.md`, `MASTER-OVERVIEW-2026-04.md`, `RELEASE_NOTE.md`
- `git log --oneline -40` (TECHAI_CODE)

**발견**: 문서들은 *Wails 시점*이고 Electron 전환은 커밋 로그에만 박혀있음:
```
2533fa4 Windows exe 16MB (WebView2 내장)
db1adbc 폐쇄망 호환 — CDN 제거
bcd9890 Windows 검은 화면 해결           ← 핵심 증상
e01a5b7 WebView2 Fixed Version 번들
811e3ab WebView2 경로 자동 탐색 강화      ← 마지막 발버둥
4d8db69 Electron 버전 초기 구조 완성     ← 결단
ffca695 Electron 완전 포팅 — API 46개
```

**포팅 후보 추림**: toolparse(text tool_call 3포맷) · GitGraph · xterm.js 멀티터미널 · CodeMirror 6 래퍼 · Live Server · 검색 패널 · 세션 영속화 · 마크다운 미리보기 · Cmd+P/Cmd+Shift+P 묶음. (knowledge packs 74개는 hanimo 한국 MCP 35종과 중복이라 제외)

### Phase 2 — Electron 전환 이유 분석 + hanimo도 가야 하나

**질문**: "택가이코드쪽 전면 electron으로 변경했는데 해당이유랑 문서 싹다 참고하고 여기도 어떻게할지 고민좀해보자"

**분석 결론**:
- 택가이의 Electron은 *증상 대응*이지 *철학 변경* 아님
- 핵심 트리거: 사내망 Win + WSUS + WebView2 자동배포 차단 → 검은 화면
- 아키텍처 변환: Wails(Go가 WebView 직접) → Electron(Chromium) + Go-server(REST 46개 + SSE)
- 진짜 자산은 Electron이 아니라 **Go-server 분리 패턴** (`5a59a35` shim)

**hanimo 권장**: 가지 마. Wails 유지 + NSIS 인스톨러 + 부재 안내화면으로 1% 케이스 잡음. 택가이가 배운 Go-server 분리 패턴은 차용.

### Phase 3 — 외부 윈도우는 진짜 문제 없나

**질문**: "그럼 electron 안써도 외부에선 윈도우여도 문제가 없을것이다?"

**솔직한 답변**: 99%는 OK, 1%는 깨짐. 정확한 케이스 표:

| 케이스 | 빈도 | 원인 |
|---|---|---|
| 일반 Win10 1809+ / Win11 | 거의 다 OK | WebView2 자동배포 (2021~) |
| Win10 N/KN 에디션 | 가끔 깨짐 | 미디어팩 미설치 |
| 인터넷 없이 첫 실행 | 깨짐 | 부트스트래퍼 다운로드 불가 |
| 기업 도메인 PC 외부망 사용 | 깨질 수 있음 | WSUS/그룹정책 |
| Win Server, 오래된 빌드 | 깨짐 | Evergreen 미적용 |

**해결**: NSIS 인스톨러 빌드 + 부재 시 안내화면 + Fixed Version 번들은 *옵션* 변형.

### Phase 4 — 듀얼 유지의 비용

**질문**: "그럼 초경량화 wails버전, electron 2가지 버전 둘다 유지하면, 자원낭비야?"

**답변**:
- 자원낭비 *여부는 아키텍처에 달림*. 단일 백엔드 + 다중 셸이면 +10~15%, 분기 코드베이스면 +100%
- 그러나 hanimo한테는 더 나은 선택지 = **브라우저-온리 모드** (Go-server + 사용자 자기 브라우저)
- Electron 대비 패키징 0, 원격/SSH 사용자도 커버, PWA 미래 확장 가능

### Phase 5 — 명성 폭발 시 회사 대응

**질문**: "hanimo 명성이 유명해져서 여러 회사에서 내부에 가져다쓰려고하면 어떻게 대응하지?"

**시나리오 분류**:
1. 개발자 1명 사용 — 좋은 일
2. 팀 단위 self-host (webui도 self-host) — webui 수익 모델 위협
3. SI/벤더가 포크+리브랜드 후 납품 — 명성·정체성 강탈

**3가지 포지션 비교**:
- A) Pure Apache 2.0 + 호스팅 SaaS — 명성 ✓ 매출 0
- B) Open Core — 표준 플레이북, GitLab/Sentry/Grafana
- C) BSL/SSPL — OSS 정체성 충돌

**권장 + 사용자 합의 도달**: B 변형 = "코드는 100% OSS, 검증·포팅·지원이 상품" + 트레이드마크 · CLA grant · 텔레메트리 0

### Phase 6 — 저작권 무비용 보호 가능?

**질문**: "특허나 기타 힘 안쓰고 내 저작권을 유지할 방법은 없나?"

**답변**: YES, 거의 다 공짜.
- 저작권 자동 귀속 (베른협약)
- 5가지 무비용 조치: LICENSE / SPDX 헤더 / NOTICE / GPG 서명 / README 저자 명시
- 한국저작권위원회 등록 (수만원, 선택, 명성 신호 후 권장)
- 트레이드마크 차선책: 도메인·SNS 핸들 선점 + 선사용 증거 누적

### Phase 7 — LTS-Onprem 발상

**질문**: "그러면, 내부망같은걸로 버전 중간중간마다 한번씩만 포팅해주는건 어떻게 생각해 only onpremise 이런식으로"

**평가**: 업계 정통 패턴 (Postgres LTS / GitLab CE-EE / Mattermost / Sentry self-host).

**hanimo 매핑**:
- Mainline: 매주~격주, OSS
- LTS-Onprem: 분기 1회, `-tags=onprem` 빌드 플래그
- 단일 코드베이스, 분기 X
- 상품: *서비스* (인증 빌드·SLA·백포트·컨설팅) — 코드는 100% OSS 유지
- 장점: TECHAI에서 검증한 폐쇄망 패치들이 그대로 onprem variant 자산

### Phase 8 — 정책 인프라 일괄 작업

**지시**: "지금까지 정책 라이센스같은거 작업 다 시작하고 문서화 정리까지 싹다해놔"

**산출** (Task 1~5 완료):

| Task | 산출물 |
|---|---|
| 1 | `LICENSE` (Apache 2.0) + `NOTICE` |
| 2 | `docs/policy/{README,copyright,trademark-and-naming,telemetry-and-privacy,lts-onprem}.md` |
| 3 | `CONTRIBUTING.md` (DCO + CLA grant + SPDX 헤더) |
| 4 | `internal/build/{default,onprem}.go` + `Makefile` `build-onprem` 타겟 |
| 5 | `README.md` MIT→Apache 2.0 + Name & Marks/Policies 섹션 + `hanimo-code-desktop/README.md` |

**검증**:
- `go vet ./internal/build/` ✓
- `go vet -tags=onprem ./internal/build/` ✓
- `grep -c MIT` 모든 라이선스 파일·README에서 0건

### Phase 9 — 본 세션 문서화

**지시**: "일단 전부 문서화해서 정리하자~"

**산출**:
- `docs/strategy/STRATEGY-2026-04-27-LICENSE-DEPLOYMENT-IP.md` — 5 결정 + 거부 대안 + 시장 분업 + 시나리오 매트릭스 + 즉시·트리거 액션
- `docs/SESSION-2026-04-27-CLOSING.md` — 30초 복구
- `docs/SESSION-2026-04-27-INDEX.md` (이 파일) — 시간순 산출

---

## 산출 통계

| 항목 | 수 |
|---|---|
| 신규 파일 | 13 |
| 수정 파일 | 3 (README · desktop README · Makefile) |
| 코드 변경 줄 수 | ~50 (build flag 분기 + Makefile 타겟) |
| 문서 변경 줄 수 | ~1500 |
| 거부한 대안 (명문화) | 6 |
| 트리거 기반 후속 액션 (명문화) | 7 |

---

## 본 세션의 *코드* 변경 의의

코드는 거의 안 바꿨지만 *되돌릴 수 없는 시멘틱*을 박았다:

1. **`internal/build` 패키지 신설** — 향후 어떤 코드든 `build.Onprem` 분기 가능. 첫 사용처가 생기면 import 1줄로 끝.
2. **Makefile `build-onprem` 진입로** — 어떤 PR이든 onprem 빌드를 깨뜨리면 CI에서 즉시 발견 가능 (CI 셋업 시 추가 필요).
3. **`-tags=onprem` 컴파일 타임 보장** — 텔레메트리·외부 통신 0의 *증명 가능한* 형태가 처음으로 존재.

이 3개가 향후 PR 리뷰 시 "이 변경이 onprem 약속을 위반하지 않나?"의
체크 포인트가 된다.

---

## 다음 세션 우선순위 (CLOSING 참고)

1. 🥇 Option A — 본 세션 산출 커밋 + 푸시
2. 🥈 Option B — 도메인·핸들·GPG 선점 (사용자 직접)
3. 🥉 Option C — Phase 20 Browser-only 모드 구현
4. 🏅 Option D — Phase 15b2 Anthropic transport
5. 🏵 Option E — TECHAI ↔ Enterprise Edition 리포지셔닝 협의

자세한 명령은 [`SESSION-2026-04-27-CLOSING.md`](SESSION-2026-04-27-CLOSING.md).
