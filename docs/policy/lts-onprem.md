# LTS · 온프레미스 정책

## 1. 두 개의 릴리스 채널

hanimo는 두 개의 평행한 릴리스 채널을 운영한다.

| 채널 | 캐덴스 | 대상 | 빌드 명령 |
|---|---|---|---|
| **Mainline** | 매주~격주 | OSS 사용자 · 개인 개발자 · 빠르게 신기능 받고 싶은 사용자 | `make build` (vanilla) |
| **LTS-Onprem** | 분기 1회 | 사내망 · 폐쇄망 · 한국 금융/공공 · 보안 검증 사이클 필요한 회사 | `make build-onprem` (`-tags=onprem`) |

**핵심 정신**: 코드는 동일한 단일 코드베이스. *빌드 플래그만* 다르다.
별도 디렉토리/포크를 만들지 않는다.

## 2. LTS 캐덴스

- **분기 1회** LTS 태그 발행: `v0.5-LTS`, `v0.6-LTS`, `v0.7-LTS`, ...
- **릴리스 윈도우**: 각 분기 마지막 주 (`-12-25`, `-03-25`, `-06-25`,
  `-09-25` 인근)
- **RC 기간**: LTS 태그 4주 전 RC1 발행 → 사내망 회귀 테스트 → 안정화 후
  최종 LTS 태그
- **Mainline은 영향받지 않음**: LTS 작업이 진행되어도 mainline은 평소
  속도 유지

## 3. 지원 윈도우

**최신 LTS + 직전 LTS, 총 2개**를 동시 지원한다.

```
v0.5-LTS    [지원중 ─────────────────────]
            └─ 2026-09-25 출시 (가상)
            └─ 2027-03-25 EOL

v0.6-LTS               [지원중 ─────────────────────]
                       └─ 2026-12-25 출시
                       └─ 2027-06-25 EOL

v0.7-LTS                          [지원중 ─────────────────────]
                                  └─ 2027-03-25 출시
                                  └─ 2027-09-25 EOL
```

각 LTS는 출시 후 **12개월간 보안 패치 백포트**를 받고, 직전 LTS와
**6개월 중첩**된다 (마이그레이션 여유).

## 4. 백포트 범위

LTS 브랜치에 백포트되는 것:

- ✅ **보안 취약점** (CVE, 시크릿 노출, 권한 우회)
- ✅ **데이터 손상 버그** (file_write 깨짐, undo 실패 등)
- ✅ **빌드 실패 수정** (해당 OS·아키텍처에서 컴파일 안 되는 경우)

LTS 브랜치에 백포트되지 *않는* 것:

- ❌ 신규 기능 (mainline에만)
- ❌ 신규 LLM 프로바이더 추가 (mainline에만)
- ❌ UX 개선 (테마 추가, 단축키 변경 등)
- ❌ 의존성 메이저 업그레이드 (보안 이슈 없으면)
- ❌ 리팩토링·코드 정리

이 분리가 LTS의 *핵심 가치*다: "변하지 않는다"는 약속.

## 5. `-tags=onprem` 빌드의 컴파일 타임 차이

`make build-onprem`은 `-tags=onprem`을 활성화하여 다음을 *컴파일 시점에*
강제한다:

| 영역 | Mainline | Onprem |
|---|---|---|
| 외부 통신 | 사용자 명시 설정 시 | 동일 (정책상 기본 0) |
| 업데이트 확인 코드 | 빌드에 포함 (옵트인) | 빌드에서 제외 |
| 외부 폰트 / CDN 호출 | 없음 (정책상) | 없음 (강제) |
| WebView2 부트스트랩 (Wails 데스크톱) | 활성 | 비활성 (호스트 환경 가정) |
| 기본 LLM 엔드포인트 | `localhost:11434` (Ollama) | 비어 있음 (사용자 명시 필수) |
| 디버그 로그 | 활성 | 사용자 옵트인 후 활성 |

코드에서는 `internal/build`의 `Onprem`/`ProfileName` 상수를 import해서
분기:

```go
import "github.com/flykimjiwon/hanimo/internal/build"

if build.Onprem {
    // 폐쇄망 전용 분기
}
```

## 6. 인증 빌드 vs 셀프 컴파일

Apache 2.0 라이선스 하에 누구나 `make build-onprem`으로 직접 컴파일할
수 있다. 그러나 회사들이 *진짜로 원하는 것*은 코드가 아니라:

| 상품 | 내용 |
|---|---|
| **인증 빌드** | 메인테이너가 직접 컴파일·서명한 바이너리. 해시·서명 검증 가능 |
| **보증서** | "이 빌드는 X 일자 기준 OWASP/CVE 검증 통과" PDF |
| **백포트 SLA** | 12개월간 발견되는 보안 이슈에 대한 백포트 보장 |
| **한국어 지원** | 평일 4시간 응답, 사내망 트러블슈팅 |
| **온사이트 도입** | 모델 endpoint 연동, knowledge pack 커스텀 |

이는 모두 *서비스*이며 코드는 그대로 OSS다. 라이선스/특허/상표를
건드리지 않는다.

> **상품 가격 가이드는 이 문서에 명시하지 않는다**. 첫 계약 시점에
> 시장 가격을 보고 책정. 현재로서는 enterprise 문의가 들어오면 개별
> 견적.

## 7. 호환성 약속 / 비약속

LTS 브랜치 *내에서*:

- ✅ `~/.hanimo/config.yaml` 형식 호환 유지
- ✅ `~/.hanimo/sessions.db` 스키마 호환 (마이그레이션 도구 제공)
- ✅ `.hanimo.md` 프로젝트 프로파일 호환
- ✅ MCP 프로토콜 호환
- ✅ 슬래시 명령 이름·동작 호환

LTS 브랜치 *간*:

- ⚠️ 위 항목들은 호환되지 않을 수 있음 (마이그레이션 가이드 동봉)
- ⚠️ CLI 플래그·환경변수 변경 가능 (changelog 명시)

## 8. EOL 절차

LTS의 EOL 시점에:

1. EOL 30일 전 GitHub Discussion에 공지
2. EOL 당일: README의 LTS 표에서 해당 버전을 "EOL" 표시로 이동
3. EOL 후: 보안 패치 중단. 단 *심각한 (CVE-Critical)* 이슈는 예외적
   백포트 가능 (메인테이너 재량)
4. 인증 계약 고객에게는 별도 EOL 통지 + 마이그레이션 권유

## 9. 신규 LTS 출시 체크리스트

각 LTS 발행 시 반드시 통과해야 하는 체크:

- [ ] `make build-onprem-all` 모든 OS×Arch 빌드 성공
- [ ] 회귀 테스트 (`go test ./...`) 100% 통과
- [ ] 에어갭 환경 시뮬레이션: 외부 통신 0건 검증
  (`docker run --network=none ... hanimo-onprem-...`)
- [ ] CDN/외부 도메인 grep 0건 (`grep -rn "https://" --include="*.go"
  --include="*.html" --include="*.css"`)
- [ ] CHANGELOG.md 갱신
- [ ] `docs/policy/lts-onprem-history.md` (향후 생성)에 새 LTS 항목 추가
- [ ] GitHub Release: signed binaries + checksums.txt + SBOM (향후)
- [ ] 직전 LTS의 EOL 일자 README 표 갱신

## 10. 첫 LTS 발행 트리거 조건

현재 hanimo는 **0.x 버전대**이며, 첫 LTS는 다음 조건이 모두 충족될 때
발행한다:

- mainline이 `v1.0` 도달 (API 안정화)
- 외부 회사 3곳 이상이 LTS 빌드를 명시적으로 요청
- 보안 회귀 테스트 자동화 완료 (CI에서 매 커밋 실행)

이전까지는 **mainline만 운영**한다. 본 문서는 그 시점이 왔을 때 즉시
적용 가능하도록 *미리 정의*해둔 것이다.

## English summary

hanimo runs two parallel release channels from a single codebase
controlled by Go build tags. **Mainline** ships continuously for OSS
users; **LTS-Onprem** ships quarterly snapshots tailored for closed
networks (Korean enterprises, banks, public sector). LTS branches receive
12 months of backported security patches with a 6-month overlap between
adjacent LTS versions. The product is a *service* (certified builds, SLA,
support) rather than the code itself — the code remains 100% Apache 2.0.
The first LTS will be cut after `v1.0` and three enterprise requests.
