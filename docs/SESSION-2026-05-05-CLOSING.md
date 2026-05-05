# Session Closing — 2026-05-05 30초 요약

> 다음 세션 진입 시 **30초 안에 맥락 복구**. 본 세션은 05-03 brand v1
> 라운드 + IDE 첫 실행 3-bug fix 직후, 그 *postmortem 후속 작업*을
> 정리한 짧은 정비 세션이다.

---

## 30초 요약

origin/main = `ac72cdc`에서 5 commits 추가. 묶음 셋(A/B/C) 모두 정리:

- **A (청소 3건)**: 5/3 작성됐지만 untracked였던 STACK-COMPARISON 문서 commit · landing-mockups submodule pointer drift(3 commits 흡수) · FEATURES.md/TODO.md 테스트 카운트 실측(30·49 → 44).
- **B (postmortem Bug #2 close)**: `SwitchModel` 부수효과 근본 fix. 결정 로직을 순수 함수 `resolveSwitchTarget`로 분리 → 동일 provider 단축 회로 + Ollama pre-flight ping. picker에서 Ollama 모델 클릭 시 daemon 안 떠 있으면 *config 미수정* + 명시적 에러. 신규 7 단위 테스트.
- **C (postmortem 향후 방지 #1 foundation)**: `internal/build/{default,onprem}.go` 컴파일 타임 분기 도입. default = `openai/gpt-oss-120b`/Novita, onprem = `gpt-oss-120b`/vllm. `GetBuildProfile()` Wails binding 추가. `go test -tags=onprem ./...` 양쪽 53/53 PASS.

테스트: 44 → **53/53 PASS** (default + onprem 양쪽 동일).

---

## 이번 세션 푸시된 커밋 (5 commits, hanimo-code only)

```
6aa92b5  feat(ide): build profile 분기 — default(외부망/Novita) vs onprem(TECHAI/vllm)
c3d5fd7  fix(ide): SwitchModel 부수효과 근본 수정 — 동일 provider 보호 + Ollama pre-flight ping
ab26372  docs(ide): FEATURES.md/TODO.md 테스트 카운트 실측 갱신 (44개)
23ede08  chore(submodule): landing-mockups → e3728f0 (3 commits 흡수)
bdfc0e8  docs(strategy): STACK-COMPARISON-2026-05-03 — TUI/Desktop stack 비교 + 결정 권고
```

(05-03 마라톤 후속이라 본 세션은 hanimo-code 한 레포만 만짐.)

---

## postmortem (`docs/postmortems/2026-05-03-ide-startup-bugs.md`) 진행도

| 항목 | 상태 | 위치 |
|---|---|---|
| Bug #1 (Model ID prefix) | 임시 fix 적용됨 (05-03) | `bindings_phase6.go` Tier 1·2 카탈로그 |
| Bug #2 (SwitchModel 부수효과) | **근본 fix 본 세션** | `bindings_phase6.go::resolveSwitchTarget` + `pingOllama` |
| Bug #3 (가짜 PTY) | 영구 fix 적용됨 (05-03) | `terminal.go` (`creack/pty v1.1.24`) |
| 향후 방지 #1 (build profile 분리) | **foundation 본 세션** · 카탈로그 분리는 후속 PR | `internal/build/{default,onprem}.go` |
| 향후 방지 #2 (외부 endpoint smoke test in CI) | 미시작 | — |
| 향후 방지 #3 (PTY 검증 macOS/Win/Linux) | 미시작 | — |
| 향후 방지 #4 (Picker side effect 명시) | B fix로 자연 해소(error 반환) | — |

---

## 검증 (실측)

```bash
cd hanimo-code-desktop

# default 빌드
go test ./...                     # ok 53/53
go build -o /tmp/h-default .      # 8.0 MB
# onprem 빌드
go test -tags=onprem ./...        # ok 53/53
go build -tags=onprem -o /tmp/h-onprem . # 8.0 MB
```

---

## 다음 세션 진입점 후보

순서는 영향도 높음 → 낮음.

1. **카탈로그 자체를 build profile별로 분리** — 본 세션 C의 *후속 PR*. `tier1Catalog`/`tier2Catalog`를 `catalog_default.go` (//go:build !onprem) + `catalog_onprem.go` (//go:build onprem)로 split. onprem 카탈로그는 사내망에서 도달 불가능한 anthropic/openai/google/novita 항목 제거. 그래야 Bug #1 catalog porting 함정이 *완전히* 차단됨.
2. **외부 endpoint smoke test (CI)** — Novita/Anthropic/OpenAI/Google `/v1/models` ping 하나씩. release 시 한 번. postmortem 향후 방지 #2.
3. **PTY 검증 macOS/Win/Linux** — Wails Linux WebKitGTK가 Stack Comparison 결정 게이트(`docs/strategy/STACK-COMPARISON-2026-05-03.md` §3.5 Scenario A) 안에 들어가는 항목. 6개월 안에 stress test 필요.
4. **TODO.md v0.2.x polish 4건** — title 툴팁 / 플랫폼 감지(Cmd↔Ctrl) / 행 높이 24px+ / QuickOpen 스크롤 따라가기. 모두 "낮음" 난이도.
5. **TEST_CHECKLIST 100+ 수동 UI 항목 회귀 확인** — v0.2.x 빌드에서 한 번도 안 돌렸음.

---

## 정합성 / source-of-truth 메모

- **테스트 카운트**: `go test ./... -v | grep -c "^=== RUN"` 가 정답. 본 세션 종료 시점 53. 향후 갱신 시 README.md / FEATURES.md / TODO.md 셋 다 같은 숫자여야 정합.
- **build profile 상수**: `internal/build/{default,onprem}.go`가 source-of-truth. 새 profile 추가 시 두 파일 모두에 같은 키를 export해야 함 (`Profile`, `RecommendedProvider`, `RecommendedSuperModel`, `RecommendedDevModel`).
- **landing-mockups submodule**: 본 세션 pointer는 `e3728f0`. 안에 `.omc/state/*.json` runtime 캐시가 untracked로 떠다니지만 무시.

---

## 관련 자산
- 핵심 PR (이번 세션): `c3d5fd7` (B) · `6aa92b5` (C)
- 직전 postmortem: [`docs/postmortems/2026-05-03-ide-startup-bugs.md`](postmortems/2026-05-03-ide-startup-bugs.md)
- Stack 결정: [`docs/strategy/STACK-COMPARISON-2026-05-03.md`](strategy/STACK-COMPARISON-2026-05-03.md)
- 직전 closing: [`SESSION-2026-05-03-CLOSING.md`](SESSION-2026-05-03-CLOSING.md)
