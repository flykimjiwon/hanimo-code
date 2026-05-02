# 텔레메트리 · 프라이버시 정책

## 1. 핵심 원칙: hanimo는 사용자를 추적하지 않는다

hanimo의 출시 직후 동작:

- **외부 호출 0건** — `~/.hanimo/config.yaml`에 사용자가 명시적으로 LLM
  엔드포인트를 설정하기 전까지 어떤 외부 서버에도 연결하지 않는다
- **고유 식별자 0개** — installation ID, machine ID, anonymous user ID
  등 어떤 형태의 사용자 식별자도 *생성하지 않는다*
- **사용량 보고 0건** — 어떤 명령을 얼마나 자주 사용하는지, 어떤 모델을
  쓰는지에 대한 데이터를 외부로 전송하지 않는다
- **에러 보고 0건** — 크래시·에러 로그는 로컬 파일에만 기록되며 자동
  업로드되지 않는다

이 원칙은 향후에도 *기본값*으로 유지된다. 텔레메트리가 도입되더라도
**명시적 옵트인이 없으면 0이 기본**이다.

## 2. 외부와 통신하는 *유일한* 경로

hanimo가 외부와 통신하는 경우는 다음 사용자가 *직접 트리거*한 것뿐이다:

| 통신 경로 | 트리거 | 데이터 |
|---|---|---|
| LLM API 호출 | 사용자가 메시지 전송 | 사용자 입력 + 컨텍스트 → 사용자가 설정한 엔드포인트 |
| MCP 서버 호출 | 사용자가 MCP 서버 설치/활성화 | 사용자 입력 → 사용자가 설정한 MCP 서버 |
| 업데이트 확인 | (향후) 사용자가 `--check-updates` 명령 실행 | 현재 버전 → GitHub Releases API |
| 패키지 설치 | 사용자가 `go install` / `brew install` 실행 | OS 패키지 매니저 표준 동작 |

**hanimo 자체가 *주도해서* 외부와 통신하는 경로는 없다.**

## 3. 로컬 파일 기록

다음은 *로컬 디스크에만* 기록되며 외부로 전송되지 않는다:

| 파일 | 위치 | 내용 |
|---|---|---|
| 채팅 세션 | `~/.hanimo/sessions.db` (SQLite) | 사용자 메시지 + LLM 응답 + 도구 실행 결과 |
| 설정 | `~/.hanimo/config.yaml` | API 키 · 프로바이더 · 모델 · 테마 |
| 디버그 로그 | `~/.hanimo/debug.log` | 디버그 빌드에서만 생성, 청크 타이밍 등 |
| 스냅샷 | `~/.hanimo/snapshots/` | 파일 수정 전 자동 백업 (undo용) |
| 메모리 | `sessions.db`의 `memories` 테이블 | `/remember`로 저장한 항목 |
| 사용량 통계 | `sessions.db`의 `usage_log` 테이블 | 토큰 수 · 비용 추정 (로컬 표시용) |

이들은 사용자가 자신의 도구 사용을 *스스로* 검토할 수 있도록 기록되며,
hanimo 메인테이너나 제3자에게 공유되지 않는다. 사용자가 직접 삭제할
수 있다(`rm -rf ~/.hanimo/sessions.db`).

## 4. 에어갭 / 폐쇄망 모드

hanimo는 **인터넷이 없는 환경에서도 완전히 동작**하도록 설계됐다:

- LLM 엔드포인트가 사내망 vLLM/LM Studio/Ollama이면 외부 인터넷 불필요
- Knowledge Pack 62개는 바이너리에 `go:embed`로 박혀 있어 다운로드 불요
- 폰트·CSS·아이콘 자산은 모두 번들. CDN 호출 없음
- 의존성 다운로드는 빌드 시점에만 발생 (사용자 실행 시점에는 0)

### `-tags=onprem` 빌드의 추가 보장

`make build-onprem`으로 컴파일된 바이너리는 다음을 *컴파일 타임에 강제*한다:

- 업데이트 확인 명령 비활성화 (코드 자체에 없음)
- 외부 도메인 default 값 제거 (사용자가 명시적으로 설정해야만 동작)
- 데스크톱 변형의 경우 WebView2 부트스트랩 호출 제거

자세한 내용은 [`lts-onprem.md`](lts-onprem.md) 참고.

## 5. 향후 텔레메트리 도입 시 규칙

만약 hanimo의 미래 버전에 텔레메트리가 추가된다면, 다음 규칙을 *반드시*
따른다:

1. **명시적 옵트인 only** — 첫 실행 시 또는 설정에서 사용자가 `true`로
   바꾸지 않는 한 0 데이터
2. **데이터 최소화** — 무엇이 전송되는지를 한 페이지에 모두 열거. 익명
   사용자 ID조차 *별도 옵트인*
3. **로컬 미리보기** — `hanimo telemetry preview` 같은 명령으로 *지금
   전송될* 데이터를 사용자가 볼 수 있어야 함
4. **종료 명령** — `hanimo telemetry off`로 *영구적으로* 비활성화 가능 +
   설정 파일에 기록되어 업데이트 후에도 유지
5. **이 문서 갱신** — 어떤 데이터가 어디로 가는지 명시
6. **PR 리뷰 필수** — 텔레메트리 추가 PR은 메인테이너 외 최소 1인의
   외부 리뷰 필요

위 규칙을 위반하는 텔레메트리는 *디자인부터 거부*된다.

## 6. 데스크톱 IDE의 추가 고려사항

`hanimo-code-desktop` (Wails 기반 IDE)에는 다음이 *없다*:

- Sentry / Bugsnag / 기타 에러 추적 SDK
- Google Analytics / 사용량 분석 SDK
- 광고 / 추적 픽셀
- 외부 폰트 CDN (Geist 등은 번들)
- 외부 아이콘 CDN (Lucide 등은 번들)
- 자동 업데이트 (사용자가 `wails build` 또는 GitHub Releases에서 직접
  다운로드)

택가이코드(TECHAI) 폐쇄망 작업에서 검증된 패턴이 그대로 적용되어 있다.

## 7. 사용자가 추가로 보호할 수 있는 방법

특별히 민감한 환경에서 사용 시:

1. **방화벽 규칙** — hanimo 프로세스의 outbound를 LLM 엔드포인트 도메인
   하나로 제한 (예: `localhost:11434`만 허용)
2. **API 키 분리** — hanimo용 API 키를 별도 발급해서 *그 키의 사용량만*
   추적 (혹시라도 의심되는 트래픽이 있으면 그 키만 폐기)
3. **세션 자동 삭제** — `crontab`으로 `~/.hanimo/sessions.db` 주기적 삭제
4. **시크릿 검증 강화** — `internal/tools/secrets.go`의 패턴 확장
5. **소스 빌드** — 바이너리 다운로드 대신 `go build`로 직접 컴파일

## 8. 침해 의심 시

만약 hanimo가 본 정책에 반하여 외부와 통신하는 정황을 발견하면:

1. **재현 절차 + 패킷 캡처(가능한 경우)** 첨부
2. **GitHub Issue로 공개 보고** (`type: privacy-violation` 라벨)
3. **메인테이너 응답 시한**: 72시간

본 정책은 hanimo의 가장 핵심적인 약속 중 하나이며, 위반 시 즉시
릴리스 차단·핫픽스·공개 사과의 대상이다.

## English summary

hanimo performs zero outbound calls at startup and never collects user
identifiers, usage metrics, or crash reports. The only outbound traffic
is what the user explicitly configures (LLM endpoints, MCP servers).
All session data, logs, and snapshots stay on the local disk. The
`onprem` build tag enforces these defaults at compile time. Any future
telemetry must be opt-in by default and documented here before merge.
