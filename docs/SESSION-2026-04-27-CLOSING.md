# Session Closing — 2026-04-27 30초 요약

> 이 문서는 다음 세션 진입 시 **30초 안에 모든 맥락 복구**를 위해 클로징
> 시점에 박은 단일 요약. 모든 세부사항은 이 파일이 가리키는 STRATEGY/
> policy 문서에서 확장.

---

## 30초 요약

`origin/main = bb81cf0` (이번 세션 미푸시 상태). 이번 세션은 **코드 변경 0,
정책·라이선스·문서 인프라 100%**.

배경: 택가이코드(TECHAI)가 사내망 Windows의 WebView2 부재 때문에 Wails
→ Electron 전면 전환한 사실 발견 → hanimo도 같은 길 가야 하나 고민 →
*가지 마, 그 대신 LTS-Onprem 빌드 채널과 Apache 2.0 + 트레이드마크 정책
인프라부터 박자*로 결론. 명성 폭발해서 회사들이 사내 도입 시도하기 *전에*
응답할 정책 무기 다 갖춤.

산출:
- LICENSE (Apache 2.0) + NOTICE 신설
- CONTRIBUTING.md (DCO + CLA grant)
- 정책 5종: copyright/trademark/telemetry/lts-onprem/index
- internal/build 빌드 태그 골격 + Makefile `build-onprem` 타겟
- README 라이선스 섹션 전면 갱신 (MIT → Apache 2.0)
- STRATEGY 문서 1종 (오늘의 모든 결정 + 거부 대안 + 트리거 액션)

---

## 이번 세션 변경 파일 (커밋 전)

```
M  Makefile                              ← .PHONY + build-onprem 타겟 2개
M  README.md                             ← MIT → Apache 2.0 + Name&Marks/Policies/License
M  hanimo-code-desktop/README.md         ← 라이선스 포인터
M  landing-mockups                       ← 이번 세션 무관 (이전 세션 잔재)

??  LICENSE                              ← Apache 2.0 풀텍스트
??  NOTICE                               ← 트레이드마크 안내 포함
??  CONTRIBUTING.md                      ← DCO + CLA grant + SPDX 헤더 컨벤션
??  internal/build/                      ← default.go + onprem.go (//go:build 분기)
??  docs/policy/                         ← README + 4 정책 문서
??  docs/strategy/STRATEGY-2026-04-27-LICENSE-DEPLOYMENT-IP.md
??  docs/SESSION-2026-04-27-CLOSING.md   ← 이 파일
??  docs/SESSION-2026-04-27-INDEX.md     ← 시간순 산출
```

검증 통과:
- `go vet ./internal/build/` (default 프로파일) ✓
- `go vet -tags=onprem ./internal/build/` (onprem 프로파일) ✓
- `grep -c MIT` README/NOTICE/LICENSE 모두 0건

---

## 박힌 결정 5건 (되돌릴 수 없음)

전체 근거는 `docs/strategy/STRATEGY-2026-04-27-LICENSE-DEPLOYMENT-IP.md`.

| # | 결정 |
|---|---|
| 1 | **라이선스 = Apache 2.0** (MIT/BSL 거부) |
| 2 | **데스크톱 셸 = Wails 단독** (Electron 추가 거부) |
| 3 | **폐쇄망 대응 = 분기 LTS-Onprem 빌드** (`-tags=onprem`) — 듀얼 셸 거부 |
| 4 | **트레이드마크 = 출원 보류 + 선사용 증거 누적** — 즉시 출원 거부 |
| 5 | **텔레메트리 = 컴파일 타임 0** — 향후 도입 시 명시적 옵트인 only |

---

## 다음 세션 진입점

### 🥇 Option A — 이번 세션 산출 커밋 + 푸시 (즉시, 5분)

```bash
git add LICENSE NOTICE CONTRIBUTING.md \
        README.md hanimo-code-desktop/README.md Makefile \
        internal/build/ docs/policy/ \
        docs/strategy/STRATEGY-2026-04-27-LICENSE-DEPLOYMENT-IP.md \
        docs/SESSION-2026-04-27-CLOSING.md \
        docs/SESSION-2026-04-27-INDEX.md

git commit -s -m "$(cat <<'EOF'
docs: 라이선스·정책 인프라 박음 (Apache 2.0 + LTS-Onprem)

- LICENSE/NOTICE 신설 (Apache 2.0 + 트레이드마크 안내)
- docs/policy/ 5종: copyright/trademark/telemetry/lts-onprem/index
- CONTRIBUTING.md + DCO 서명 + SPDX 헤더 컨벤션
- internal/build 빌드 태그 골격 (default/onprem)
- Makefile: build-onprem / build-onprem-all
- README: MIT → Apache 2.0 + Name & Marks + Policies 섹션
- docs/strategy/STRATEGY-2026-04-27-LICENSE-DEPLOYMENT-IP.md
- docs/SESSION-2026-04-27-{CLOSING,INDEX}.md

Constraint: 명성 폭발 전에 라이선스·트레이드마크 입장 명문화 필요
Rejected: 트레이드마크 즉시 출원 | 비용·시점 부적절, 선사용 증거로 대체
Rejected: BSL/SSPL | OSS 명성 정체성과 충돌
Rejected: Wails+Electron 듀얼 | 16MB 정체성 훼손, 정당화 시장 없음
Rejected: 즉시 텔레메트리 | 사내망 검증 차단 + OSS 신뢰 훼손
Confidence: high
Scope-risk: narrow
Directive: 정책 변경은 소급 적용 금지, 미래 기여에만 적용

Signed-off-by: Kim Jiwon <shbtechgroup5@gmail.com>
EOF
)"

git push origin main
```

### 🥈 Option B — 도메인·핸들·GPG 선점 (30분, 코드 외)

`STRATEGY-2026-04-27...md` §6 즉시 액션 6개:
1. `hanimo.dev` 도메인 (1만원/년)
2. `hanimo.kr` 도메인 (2~3만원/년)
3. `@hanimo_dev` X/Twitter 핸들
4. `hanimo` GitHub organization 검토
5. GPG 키 생성 + GitHub 등록 + `git config commit.gpgsign true`
6. `enterprise@hanimo.dev` Cloudflare Email Routing

Option A 푸시 *후* 진행 권장 (도메인 등록은 사용자 본인 결제 필요).

### 🥉 Option C — Phase 20: Browser-only 모드 구현 (1~2일)

`internal/` 백엔드를 standalone HTTP 서버로 분리. Wails 데스크톱과 같은
백엔드 공유, 사용자가 `hanimo serve --port 8080` 으로 띄우면 자기 브라우저
(Chrome/Edge/Safari)로 접속 가능. WebView2 부재 1% / Linux / 원격 SSH
사용자 *전부 커버* + Electron 추가 안 함.

전제: Option A 먼저 푸시.

### 🏅 Option D — Phase 15b2 Anthropic transport (4~6h)

이전 세션 `SESSION-2026-04-25-CLOSING.md` Option B 그대로. 라이선스 인프라와
독립이라 언제든 가능.

### 🏵 Option E — TECHAI ↔ hanimo Enterprise Edition 리포지셔닝 협의

신한과의 계약 관계 검토 + TECHAI를 *first paying customer*로 발화 가능한지
판단. 비코드, 사용자 본인 협의 필요.

---

## 새 세션에서 첫 메시지 권장

> "어제 라이선스·정책 인프라 다 박았어. `docs/SESSION-2026-04-27-CLOSING.md`
> 부터 봐. Option A (커밋 푸시) 부터 가자."

또는 명시:
> "Option B (도메인) 했어, 다음 가자" / "Option C (Browser-only) ㄱㄱ"

---

## 진입 문서 우선순위

1. **`docs/SESSION-2026-04-27-CLOSING.md`** (이 파일) — 30초 복구
2. `docs/strategy/STRATEGY-2026-04-27-LICENSE-DEPLOYMENT-IP.md` — 5 결정 + 거부 대안 + 트리거 액션
3. `docs/policy/README.md` — 정책 5종 인덱스
4. `docs/SESSION-2026-04-27-INDEX.md` — 시간순 산출
5. 자동 로드: `~/.claude/.../memory/MEMORY.md` (project + feedback)

---

## 빌드 헬스체크 (다음 세션 시작 시 5분 내)

```bash
cd /Users/jiwonkim/Desktop/kimjiwon/hanimo-code
git status -sb

# 빌드 태그 양쪽 검증
go vet ./internal/build/
go vet -tags=onprem ./internal/build/

# Onprem 바이너리 첫 빌드 확인 (이번 세션 이후 처음 시도)
make build-onprem 2>&1 | tail -20
ls -la hanimo-onprem
```

`build-onprem` 빌드가 처음 돌 때 다음 가능성을 살핌:
- `internal/build` import 경로 오타 (있으면 수정)
- 다른 패키지에서 `build.Onprem` 사용 안 하면 빌드는 OK지만 *분기 동작*은 미검증

상기 모두 OK 면 그대로 다음 작업 진입.

---

_링크 시 함께 보면 좋은 문서_:
- 어제: `docs/SESSION-2026-04-25-CLOSING.md` (v0.2.0 릴리스 + Phase 19)
- 비전: `docs/strategy/VISION-2026-04-25-MULTI-MODEL-MULTI-DEVICE.md`
- 시장 매트릭스: `docs/strategy/MARKET-ANALYSIS-2026-04-24.md`
