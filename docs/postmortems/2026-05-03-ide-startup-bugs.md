# IDE Startup Bugs — Postmortem (2026-05-03)

> 본 마라톤 세션 brand v1 · 정책 · hub 작업이 *모두 푸시된 후* hanimo IDE
> (hanimo-code-desktop)를 처음 띄워본 시점에 발견된 3개 버그의 사후 기록.
> 모두 해결, *외부 marketplace 환경 차이* 인지 부족이 공통 원인.

---

## TL;DR

| # | 증상 | 진짜 원인 | Fix |
|---|---|---|---|
| 1 | 채팅 시 `404 MODEL_NOT_FOUND` | TECHAI 사내망 vLLM은 바닐라 model id 수용, Novita marketplace는 *vendor-prefix 강제* (`openai/gpt-oss-120b`). TECHAI 카탈로그 그대로 포팅하면서 컨벤션 차이 미인지 | `bindings_phase6.go` Tier 1·2 model id 정리 + `~/.hanimo/config.yaml` 직접 갱신 |
| 2 | 한번 Ollama 모델 클릭 후 *영구* `connection refused` | `SwitchModel()`이 base_url + api_key + provider를 *통째로* 덮어쓰는 부수효과. Ollama 한 번 선택하면 그 후 모든 메시지가 localhost:11434로 감 | 임시: config 수동 갱신. 근본: SwitchModel 로직 분리 (별도 PR 후보) |
| 3 | 터미널 입력 자체가 안 됨 (글자 안 보임, Enter 무시) | `cmd.StdinPipe()` + `cmd.StdoutPipe()`는 *일반 파이프*, 진짜 PTY 아님. zsh/bash가 stdin을 파이프로 받으면 non-interactive 모드로 들어가 echo·prompt·line editing 모두 죽음 | `creack/pty v1.1.24` 도입, `pty.Start(cmd)`로 진짜 PTY master/slave |

---

## Bug #1 — Model ID prefix mismatch (404 MODEL_NOT_FOUND)

### 증상
```
사용자가 IDE 채팅에 "반가워" 입력
→ Error: error, status code: 404, status: 404 Not Found, message: %!s(<nil>)
```

debug.log:
```
{"code":404,"reason":"MODEL_NOT_FOUND","message":"model not found",
 "metadata":{"reason":"model: gpt-oss-120b not found"}}
```

### 진짜 원인
TECHAI 사내망의 vLLM 운영자는 모델을 *바닐라 이름* (`gpt-oss-120b`)으로 등록.
Novita는 multi-tenant marketplace라 *vendor-prefix*가 모델 식별자의 *일부*:

| Endpoint | 모델 ID 컨벤션 |
|---|---|
| TECHAI 사내망 vLLM | `gpt-oss-120b` ✅ 통함 |
| Novita marketplace | `openai/gpt-oss-120b` ✅ / `gpt-oss-120b` ❌ 404 |

hanimo IDE의 Tier 1 카탈로그를 *TECHAI 시절 컨벤션*으로 박았는데, 외부망 빌드는
Novita를 endpoint로 쓰니 정확한 ID 매칭 실패.

### 직접 핑 검증
```bash
# 바닐라
curl -X POST https://api.novita.ai/v3/openai/v1/chat/completions \
  -H "Authorization: Bearer $KEY" \
  -d '{"model":"gpt-oss-120b",...}'
# → 404 MODEL_NOT_FOUND

# Vendor-prefix
curl ... -d '{"model":"openai/gpt-oss-120b",...}'
# → 200 OK
```

### Fix
[`hanimo-code-desktop/bindings_phase6.go`](../../hanimo-code-desktop/bindings_phase6.go)
의 `tier1Catalog` / `tier2Catalog`를 Novita marketplace 호환 ID로 정리:

| 변경 전 | 변경 후 |
|---|---|
| `gpt-oss-120b` (Provider novita) | `openai/gpt-oss-120b` |
| `qwen3-coder-30b` (Provider ollama) | 라벨 명시 "(Ollama)", Hint "Ollama daemon 필요" |
| (신규) | `qwen/qwen3-235b-a22b-instruct-2507` Tier 1 (Qwen3 235B Instruct) |
| (신규) | `qwen/qwen-2.5-72b-instruct` Tier 2 |
| (신규) | `deepseek/deepseek-v3.2` Tier 2 (Novita 직결) |
| `gemma-4-31b-it` | (제거) |

추가: `~/.hanimo/config.yaml`을 *직접 갱신*해서 default model을 Novita로 못박음
(SwitchModel이 ollama로 회귀시키기 *전에* 사용자가 Novita 모델로 메시지 보낼 수
있도록).

### 향후 방지 — 빌드 프로파일별 default model 분리
`-tags=onprem` 빌드는 바닐라 `gpt-oss-120b`, default(외부망) 빌드는
`openai/gpt-oss-120b`로 *컴파일 타임* 분기 권장. `internal/build/onprem.go` +
`default.go`에 `DefaultSuperModel`, `DefaultDevModel` 상수 추가하는 PR이 자연스러움.
docs/branding/lucide-migration.md Phase 4와 별도로 기록되는 후속 작업.

---

## Bug #2 — SwitchModel 부수효과 (config 영구 덮어씀)

### 증상
사용자가 모델 picker에서 *Ollama* 표기 모델(`Qwen3-Coder 30B (Ollama)`) 한 번
클릭 → 그 후 모든 메시지가 `http://localhost:11434/v1/chat/completions`로 감 →
Ollama daemon 안 띄워서 `connection refused`. 다른 모델 다시 선택해야 정상화.

### 진짜 원인
`bindings_phase6.go`의 `SwitchModel(modelID)` 코드:

```go
if opt := findModelOption(modelID); opt != nil {
    p := findProvider(opt.Provider)
    if p != nil {
        ...
        cfg.Provider = p.Name
        cfg.API.BaseURL = resolveBaseURL(p, cfg)  // ← 통째로 덮어씀
        cfg.API.APIKey = key                      // ← 통째로 덮어씀
    }
}
cfg.Models.Super = modelID
cfg.Models.Dev = modelID
saveTGCConfig(cfg)  // ← config.yaml에 영구 박힘
```

즉 모델 선택이 *base_url + api_key + provider까지* 사이드이펙트로 갈아끼움. 이건
*Phase 15a 자동 라우팅 의도* (사용자가 Anthropic 모델 클릭하면 자동으로 Anthropic
endpoint로 전환)였지만 *Ollama 모델 클릭 → localhost로 영구 회귀* 함정 됨.

### Fix (임시)
`~/.hanimo/config.yaml`을 직접 Novita 값으로 다시 박음. SwitchModel 자체는 손대지
않음. 사용자가 *picker 클릭하지 않는 한* 영구 유지됨.

### Fix (근본 — 별도 PR 후보)
1. `SwitchModel`을 *Provider switching* 와 *Model switching* 두 함수로 분리
2. 또는 picker UI에서 "Provider 변경 + Model 변경"을 *명시적 confirm dialog*로
3. 또는 Ollama provider 선택 시 *Ollama daemon ping*하고 실패 시 picker에 비활성화

본 세션은 *임시 fix*만, 근본 수정은 다음 세션 후보.

---

## Bug #3 — 터미널 가짜 PTY (입력 안 됨)

### 증상
- `Cmd+J`로 터미널 패널 열림
- 클릭해도 *프롬프트 안 보임*
- 키 입력해도 *글자 안 보이고* 명령 처리 안 됨
- TECHAI techai-ide에는 *별도 입력란 한 줄*이 있었는데 hanimo는 없어서 더 명확히
  "안 됨"으로 보임

### 진짜 원인
[`hanimo-code-desktop/terminal.go`](../../hanimo-code-desktop/terminal.go) 이전 코드:

```go
stdin, err := cmd.StdinPipe()    // 일반 파이프
stdout, err := cmd.StdoutPipe()  // 일반 파이프
cmd.Stderr = cmd.Stdout
cmd.Start()
```

`cmd.StdinPipe()`는 *진짜 PTY가 아니라* 그냥 *anonymous pipe*. zsh/bash는 `isatty(stdin)`
체크해서 *false*면 *non-interactive 모드*로 진입:

| TTY 모드 (PTY) | non-interactive 모드 (파이프) |
|---|---|
| 프롬프트 표시 | 표시 안 됨 |
| 키 입력 echo | echo 안 됨 (사용자 글자 안 보임) |
| Line editing (←→ Backspace) | 동작 안 함 |
| Signals (Ctrl+C 등) | 무시 |
| 색상 자동 활성 (`isatty(stdout)`) | 비활성 |
| 명령 처리 | line-by-line 일괄 (사용자가 Enter 쳐도 즉시 처리 X) |

즉 *입력 자체는 stdin pipe로 들어가지만* 셸이 *interactive*가 아니라
*batch script 입력*으로 해석. 사용자 관점에선 "아무 반응 없음".

TECHAI techai-ide의 *별도 입력란*은 이 한계 우회 — line-by-line으로 직접
받아서 셸에 전달 (또는 `zsh -c "command"` 형식으로 한 번씩 실행).

### Fix
[`creack/pty v1.1.24`](https://github.com/creack/pty) 도입:

```go
import "github.com/creack/pty"

ptmx, err := pty.Start(cmd)  // ← 진짜 PTY master/slave 생성
                              //   cmd.Stdin/Stdout/Stderr는 자동으로 slave에 연결
```

PTY master 한 개의 `io.ReadWriteCloser`로 input/output 양방향 통신. PTY slave는
셸이 *진짜 TTY*로 인식 → echo·prompt·line editing·signals 정상 동작.

추가 fix:
- `ResizeTerminal(rows, cols)` 실제 동작 (`pty.Setsize` 호출) — vim/htop 등이 정확한
  창 크기로 렌더
- 환경변수: `TERM=xterm-256color` + `COLORTERM=truecolor` 추가

### 검증
```
$ echo $TERM
xterm-256color
$ ls
[색상 출력 정상, 사용자 입력 echo 정상, Enter 즉시 처리]
```

### 영구 영향
- `hanimo-code-desktop/go.mod`에 `github.com/creack/pty v1.1.24` 의존 추가
- 터미널 입력은 *영구 fix*. 별도 입력란 UI는 불필요해짐 (TECHAI 패턴 모방 거부)

---

## 공통 교훈

세 버그 모두 *외부 환경* (marketplace API / desktop OS PTY 표준)과 *TECHAI 시절
컨벤션*의 차이를 인지 못 한 데서 옴. *포팅* 작업의 함정.

### 향후 방지 체크리스트

- [ ] **빌드 프로파일별 default 분리**: `-tags=onprem` (사내망 컨벤션) vs default
      (marketplace 컨벤션) — model id, base_url, 텔레메트리 등 컴파일 타임 분기
- [ ] **외부 의존 fixture 테스트**: Novita / Anthropic / OpenAI 등 *실제 endpoint
      ping*을 CI에 한 번씩 (이번 release smoke test) — 404 케이스 즉시 발견
- [ ] **PTY/TTY 검증**: macOS arm64/Intel + Windows + Linux에서 *실제 터미널 입력*
      smoke test (스크립트로 `ls\n` 보내고 echo 받는지)
- [ ] **Picker side effect 명시**: SwitchModel 같은 부수효과가 있는 함수는 confirm
      dialog 또는 별도 명시적 동의

### 메모리 박음 (자동 컨텍스트)
미래 세션이 같은 실수 반복 안 하도록 메모리에 박음:
- *외부 marketplace는 vendor-prefix 강제* — 단일 컨벤션 가정 금지
- *Wails desktop 터미널은 진짜 PTY 필수* — `cmd.StdinPipe()` 사용 금지

---

## 시간선

| 시각 | 이벤트 |
|---|---|
| 22:22 | 사용자 IDE에 "반가워" → 404 |
| 22:30 | 진단 시작 — config.yaml + Novita /v1/models 핑 |
| 22:32 | Bug #1 원인 확정 — vendor-prefix 누락 |
| 22:35 | config.yaml 임시 갱신 + IDE 재시작 → 채팅 정상 |
| 22:38 | bindings_phase6.go 카탈로그 정리 |
| 22:40 | wails 미설치 → `go install github.com/wailsapp/wails/v2/cmd/wails@latest` |
| 22:43 | wails build → IDE 재시작 → 새 카탈로그 노출 |
| 22:45 | 사용자가 picker에서 `Qwen3-Coder 30B (Ollama)` 클릭 → Bug #2 발현 (config 회귀) |
| 22:48 | terminal.go 분석 → Bug #3 발견 (가짜 PTY) |
| 22:52 | `creack/pty v1.1.24` 도입 + terminal.go 재작성 |
| 22:55 | wails build → IDE 재시작 → 채팅·터미널 모두 정상 |
| 22:57 | 사용자 "잘된다" 확인 |
| 23:00 | 본 사후 기록 작성 |

---

## 관련 자산
- 코드: [`bindings_phase6.go`](../../hanimo-code-desktop/bindings_phase6.go) ·
  [`terminal.go`](../../hanimo-code-desktop/terminal.go) ·
  [`go.mod`](../../hanimo-code-desktop/go.mod)
- 정책: [`telemetry-and-privacy.md`](../policy/telemetry-and-privacy.md)
- 전략: [`STRATEGY-2026-04-27-LICENSE-DEPLOYMENT-IP.md`](../strategy/STRATEGY-2026-04-27-LICENSE-DEPLOYMENT-IP.md)
- 세션: [`SESSION-2026-05-03-CLOSING.md`](../SESSION-2026-05-03-CLOSING.md)
