# hanimo Policy Index

이 폴더는 hanimo 프로젝트의 라이선스·저작권·배포·이름 사용에 관한
정책 문서를 모아놓은 곳이다. **법적 효력은 루트의 `LICENSE` + `NOTICE`
파일에서 나오고**, 여기 문서는 그것을 한국어로 풀이하고 정책적 의도를
기록하기 위함이다.

## 정책 문서 목록

| 문서 | 다루는 내용 |
|---|---|
| [`copyright.md`](copyright.md) | 저작권 자동 귀속 · 소스 헤더 · GPG 서명 커밋 · 한국저작권위원회 등록 옵션 |
| [`trademark-and-naming.md`](trademark-and-naming.md) | "hanimo" 이름·로고 사용 정책 · 포크/재배포 시 네이밍 규칙 · 공식 채널 목록 |
| [`telemetry-and-privacy.md`](telemetry-and-privacy.md) | 텔레메트리 0 · 외부 통신 정책 · 에어갭 모드 · 크래시/사용량 로그 |
| [`lts-onprem.md`](lts-onprem.md) | 분기 LTS 캐덴스 · EOL 정책 · 백포트 범위 · `-tags=onprem` 빌드 |
| [`../../CONTRIBUTING.md`](../../CONTRIBUTING.md) | 기여자 라이선스 부여 · DCO 서명 커밋 |

## 변경 절차

이 폴더 안의 정책 문서를 변경하는 PR은 다음 조건을 충족해야 한다:

1. **메인테이너 검토 필수** — 자가 머지 금지
2. **소급 적용 금지** — 정책은 미래 기여에만 적용. 과거 기여물은 기여 시점의 정책을 따름
3. **사용자 영향 분석** — 라이선스 호환성·기여자 권리에 영향이 있으면 RFC 형태로 별도 이슈 발행
4. **이중 언어** — 새 정책 문서는 한국어 본문 + 영문 요약 1단락 권장 (기존 사용자 다수가 한국어 화자)

## 자주 묻는 질문 (개발자 측)

**Q. hanimo 코드를 회사에 가져와서 써도 되나요?**
A. 네. Apache 2.0이라 자유롭게 사용·수정·재배포 가능합니다. 단 `LICENSE`와
`NOTICE` 파일은 보존해야 하며, 제품명에 "hanimo"를 포함하려면
[`trademark-and-naming.md`](trademark-and-naming.md)를 참고하세요.

**Q. 회사에서 SSO/감사 로그/SLA 같은 엔터프라이즈 기능이 필요한데요.**
A. 분기마다 발행되는 LTS 빌드 + 인증·지원 계약을 제공할 예정입니다.
[`lts-onprem.md`](lts-onprem.md)를 참고하세요. 정식 출시 전에는 메일
문의(`enterprise@hanimo.dev` — 도메인 확보 후 활성화)로 연락 바랍니다.

**Q. 텔레메트리는 정말 0인가요?**
A. 네. 빌드 직후 hanimo는 외부와 통신하지 않습니다 (LLM 엔드포인트는
사용자가 명시적으로 설정한 경우에만). [`telemetry-and-privacy.md`](telemetry-and-privacy.md)
참고.

**Q. hanimo를 포크해서 우리 회사 이름으로 재배포해도 되나요?**
A. 코드는 자유롭게. 단 *이름*은 별도. [`trademark-and-naming.md`](trademark-and-naming.md)
의 네이밍 규칙을 따라주세요.
