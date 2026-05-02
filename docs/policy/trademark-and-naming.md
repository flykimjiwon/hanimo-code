# 트레이드마크 · 네이밍 정책

## 1. "hanimo" 이름·로고의 법적 상태

- **현재**: 미등록 트레이드마크 (unregistered mark / common-law mark)
- **소유자**: Kim Jiwon
- **사용 시작일**: 2026년 (정확한 일자는 GitHub 첫 공개 릴리스 태그 기준)
- **공식 트레이드마크 출원**: 미실시. 명성 신호가 임계점을 넘으면
  대한민국 특허청에 출원 예정

Apache 2.0 라이선스 Section 6은 명확히 명시한다:

> This License does not grant permission to use the trade names,
> trademarks, service marks, or product names of the Licensor.

즉 **코드는 자유롭지만 이름은 별도다**. 이 문서는 그 "별도"의 규칙을
정의한다.

## 2. 허용되는 사용 (출원 없이도 보호되는 영역)

다음 사용은 *환영*하며, 별도 허락 없이 가능하다:

| 사용 케이스 | OK 여부 | 비고 |
|---|:---:|---|
| "hanimo를 사용하여 X를 만들었다" 같은 *사실 기술* | ✅ | nominative fair use |
| 블로그 포스트·튜토리얼·발표 자료에서 hanimo 언급 | ✅ | 제3자 콘텐츠 환영 |
| 기여자/사용자 페이지에 "Powered by hanimo" 배지 | ✅ | 배지 이미지는 docs/branding/ 참고 |
| 학술 논문에서 hanimo 인용 | ✅ | citation에 GitHub URL 포함 권장 |
| 회사 내부 도구 통합 시 "hanimo 기반" 표기 | ✅ | 라이선스 NOTICE 보존 조건 |
| 비교 기사·리뷰("hanimo vs Cursor") | ✅ | 사실에 기반한 한 |

## 3. 허용되지 않는 사용 (이름·로고 보호 영역)

다음은 명시적 사전 허락 없이 *금지*한다:

| 사용 케이스 | 금지 사유 |
|---|---|
| 포크/파생제품의 *제품명*에 "hanimo" 포함 (예: "hanimo-pro", "hanimo for Korea") | 사용자 혼동 유발 |
| 상용 서비스의 *제품명*으로 "hanimo" 사용 | 트레이드마크 침해 |
| hanimo 로고를 자사 제품 마케팅에 사용 | 후원/제휴 오인 |
| 도메인에 "hanimo" 포함된 사이트 운영 (예: hanimo-clone.com) | 도메인 침해 |
| "공식 hanimo 한국 지부" 같은 자칭 표현 | 권한 사칭 |
| App Store / Play Store에 "hanimo"라는 이름으로 별도 앱 등록 | 사용자 혼동 |

## 4. 포크할 때 권장 네이밍

hanimo를 포크해서 자기 회사·팀의 변형판을 운영하려는 경우:

### ✅ 권장 패턴

- **회사명 + 변경 표시**: `acme-coder` (hanimo 기반이라는 사실은 NOTICE에 기록)
- **목적 표시**: `kimchi-agent` (한국어 친화 변형) — hanimo와 다른 정체성
- **`based-on-hanimo` 명시**: 제품명은 자유, README/메타데이터에
  "Based on hanimo (Apache 2.0)" 표기

### ❌ 회피해야 할 패턴

- `hanimo-pro` / `hanimo-enterprise` / `hanimo-korea`
- `hanimo2` / `hanimo-next`
- `Hanimo Cloud` / `HanimoX`

### 회색 영역 (사전 문의 권장)

- `hanimo-plugin-X` (hanimo의 플러그인을 표방하는 경우): 실제로 hanimo와
  *호환되며* 플러그인으로 동작하면 OK. 사실상 별개 제품이면 회피.
- 학습용/연구용 비공개 포크: 제품화하지 않으면 자유. 공개 시 위 규칙 적용.

## 5. 공식 채널 목록 (사칭 방지)

다음만이 공식 hanimo 채널이다. 이 외의 "hanimo 공식 X"는 사칭이다.

| 채널 | URL / 핸들 | 상태 |
|---|---|---|
| 공식 GitHub 조직 | <https://github.com/flykimjiwon> | 활성 |
| 공식 hanimo 레포 | <https://github.com/flykimjiwon/hanimo> | 활성 |
| 공식 도메인 | `hanimo.dev` (확보 예정) | 미확보 — 사용자가 직접 등록 권장 |
| 공식 X/Twitter | `@hanimo_dev` (확보 예정) | 미확보 |
| 공식 Discord | TBD | 미개설 |
| 공식 이메일 | `enterprise@hanimo.dev` (도메인 확보 후) | 미활성 |

> **메인테이너 액션 아이템 (이 정책 발효 후 1주 내)**: `hanimo.dev` /
> `hanimo.kr` 도메인, `@hanimo_dev` 핸들 선점. 도메인 비용 연 1~2만원
> 수준이며 향후 트레이드마크 출원 시 *사용 증거*로도 작용한다.

## 6. 향후 정식 출원 트리거

다음 중 하나라도 발생하면 대한민국 특허청 트레이드마크 출원을
*적극 검토*한다:

- 외부 회사·팀이 "hanimo" 이름으로 자기 제품을 등록·홍보
- 누군가 `hanimo.com` / `hanimo.kr` / 유사 도메인을 선점하여 사칭
- 월간 활성 사용자 1만 명 돌파
- 첫 엔터프라이즈 계약 체결 (방어용)

출원 비용은 유사군당 약 30~70만원 수준이며 메인테이너 개인 부담.

## 7. 침해 발견 시 대응 절차

위 §3의 금지 사용을 발견한 경우:

1. **증거 수집** — 침해 페이지/제품 URL · 스크린샷 · 발견 일자 · 침해
   범위(사용자 수, 매출 추정 등 가능한 한)
2. **공개 통지** (선택) — `docs/policy/known-trademark-issues.md`(향후
   생성)에 기록하여 커뮤니티가 인지하도록
3. **시정 요청** — 침해자에게 이메일/이슈로 정중히 시정 요청. 14일 기한
4. **공식 절차** — 사이즈가 큰 경우 변호사 자문 후 내용증명 / 가처분 /
   특허심판원 이의신청

선사용권(unregistered prior use)을 입증하려면 *공개적이고 지속적인
사용 기록*이 핵심이다. 이를 위해:

- 모든 GitHub 릴리스를 공개로 유지 (날짜 증거)
- 블로그/SNS에서 hanimo 이름으로 활동 (사용 증거)
- 회의·발표 자료 보관 (인지도 증거)

## 8. 라이선스 변경 vs 트레이드마크 변경

이 두 가지를 혼동하지 말 것:

- **라이선스 변경** (Apache 2.0 → BSL 등): 미래 기여에만 적용 가능,
  과거 코드는 영원히 Apache 2.0
- **트레이드마크 정책 변경**: 미래 사용자에게 즉시 적용 가능. 과거에
  허락받은 사용은 그 시점 정책에 따름

이 문서가 변경될 때 변경 이력은 Git 히스토리로 보존되며, 주요 변경은
`docs/policy/changelog.md`(향후 생성)에 별도 기록한다.

## English summary

The names "hanimo" and the hanimo logo are unregistered marks of Kim Jiwon.
Apache 2.0 explicitly does not grant trademark rights (Section 6). Forks
must use a distinct product name and may not place "hanimo" in their
product title; "Based on hanimo" attribution in README/metadata is
welcomed. Nominative fair use (blog posts, comparisons, citations) is
explicitly permitted. Formal trademark registration is deferred until
adoption thresholds are reached.
