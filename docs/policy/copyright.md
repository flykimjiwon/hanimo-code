# 저작권 정책

> 본 문서는 일반적 정책 안내이며 법률 자문이 아니다. 분쟁이 예상되는
> 상황에서는 한국저작권위원회 무료 상담 또는 변호사 자문을 권장한다.

## 1. 자동 귀속

대한민국은 **베른협약** 가입국이며, 저작권은 코드를 *작성하는 즉시*
저작자에게 자동으로 귀속된다. 별도의 출원·등록·수수료가 필요하지 않다.

- **저작권자**: Kim Jiwon (외부 기여자가 추가되면 해당 기여자도 자기
  기여분의 공동 저작자가 된다)
- **귀속 시점**: 각 파일이 처음 커밋되는 순간
- **증거**: Git 커밋 히스토리(타임스탬프 + 커밋 SHA + 가능한 경우 GPG
  서명)가 1차 증거로 작용한다

## 2. 소스 파일 헤더 컨벤션

신규 또는 의미있게 수정되는 모든 소스 파일에는 아래 한 줄(또는 두 줄)
헤더를 박는다.

### Go / TypeScript / JavaScript / CSS

```go
// Copyright 2025-2026 Kim Jiwon (김지원). All rights reserved.
// SPDX-License-Identifier: Apache-2.0
```

### Python / shell / Makefile

```python
# Copyright 2025-2026 Kim Jiwon (김지원). All rights reserved.
# SPDX-License-Identifier: Apache-2.0
```

### HTML / Markdown

대량의 콘텐츠 파일에는 헤더 생략 가능 (LICENSE + NOTICE가 루트에 있어
저작권은 보존됨). 단 *코드처럼 다뤄지는* HTML(빌드 산출, 데모 페이지)에는
HTML 주석으로 헤더 권장.

```html
<!-- Copyright 2025-2026 Kim Jiwon (김지원). All rights reserved. -->
<!-- SPDX-License-Identifier: Apache-2.0 -->
```

### 왜 SPDX 형식?

`SPDX-License-Identifier:` 한 줄은 자동화 라이선스 스캐너(FOSSA, REUSE,
ScanCode)가 인식하는 표준 형식이다. 회사들이 hanimo를 자기 코드베이스에
편입할 때 이 줄 하나로 컴플라이언스 처리가 끝난다.

### 누락된 파일은 어떡하나

- 신규 PR: 헤더 누락 시 리뷰에서 지적 → 수정 요청
- 기존 파일: 일괄 추가는 `scripts/add-license-headers.sh`(향후 추가 예정)로
  자동 처리. 그 전까지는 *건드리는 김에* 추가하는 점진적 방식

## 3. GPG 서명 커밋 (강력 권장)

GPG 서명은 분쟁 시 "내가 언제 무엇을 작성했다"의 가장 강력한 증거다.

```bash
# 1. GPG 키 생성 (없는 경우)
gpg --full-generate-key

# 2. 키 ID 확인
gpg --list-secret-keys --keyid-format=long

# 3. Git에 서명 활성화
git config --global user.signingkey <KEY_ID>
git config --global commit.gpgsign true
git config --global tag.gpgsign true

# 4. GitHub에 공개키 등록
gpg --armor --export <KEY_ID>
# → GitHub Settings → SSH and GPG keys → New GPG key 에 붙여넣기
```

이후 모든 커밋에 자동으로 서명이 붙고, GitHub에서 "Verified" 배지로
표시된다.

> 참고: 이 정책은 메인테이너에게 *권장*이며 외부 기여자에게 *강제*하지
> 않는다 (DCO 서명만 필수). [`../../CONTRIBUTING.md`](../../CONTRIBUTING.md)
> 참고.

## 4. 한국저작권위원회 등록 (선택)

자동으로 발생한 저작권에 *공식 등록 일자*를 박고 싶다면 한국저작권위원회
(<https://www.cros.or.kr/>)를 통해 1건당 수만 원으로 등록 가능하다.
**이는 특허·상표 출원이 아니다** — 이미 가진 저작권에 등록 일자만 박는
절차이며, 분쟁 시 입증 책임이 상대방으로 넘어가는 *추정력*을 얻는다.

권장 시점:
- 명성이 임계점을 넘어 분쟁 가능성이 보이기 시작할 때
- 메이저 릴리스 (v1.0, v2.0 등)
- 특정 회사가 무단 복제·재배포한 정황을 발견했을 때

지금 단계(v0.x)에서는 **미실시**. 명성 신호가 오면 그제야 등록한다.

## 5. 외부 코드·자산 편입 시 규칙

hanimo에 외부 코드를 가져올 때:

1. **라이선스 호환성 확인**: Apache 2.0과 호환 (MIT, BSD, ISC, Apache 2.0)
   → 가능. 비호환(GPL 계열의 강한 카피레프트) → 분리 모듈로만 가능, 핵심
   라이브러리 통합 금지
2. **NOTICE 갱신**: 외부 저작권자를 `NOTICE` 파일에 추가
3. **`THIRD_PARTY_NOTICES.md` 별도 파일** (대량 의존성 발생 시): 각
   라이브러리·라이선스·저작권자 표 작성

## 6. 침해 발견 시 대응 절차

누군가 hanimo 코드를 라이선스 위반 방식으로 재배포하는 정황 발견 시:

1. **증거 수집** — 위반 페이지/저장소 URL · 스크린샷 · 위반된 코드의
   원본 파일·커밋 SHA · 위반 발견 일자
2. **DMCA / GitHub takedown** — GitHub에 호스팅된 경우
   <https://github.com/contact/dmca-takedown> 양식
3. **이메일 통지** — 가능한 경우 위반자에게 직접 시정 요청 (라이선스 표기
   추가, 또는 제거). 14일 내 회신 없으면 다음 단계
4. **변호사 상담** — 비용 발생, 침해 규모가 크지 않으면 1·2·3까지로
   대부분 해결됨

위 절차는 메인테이너의 재량이며 외부 기여자에게 강제되지 않는다.

## English summary

Copyright in hanimo is automatic under the Berne Convention and vests
in Kim Jiwon and contributors as soon as code is committed. Source files
should carry the SPDX header `// SPDX-License-Identifier: Apache-2.0`
together with the copyright line. GPG-signed commits are encouraged for
maintainers; DCO sign-off is required for all contributors (see
`CONTRIBUTING.md`). Optional registration with the Korea Copyright
Commission is reserved for later milestones.
