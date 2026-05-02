# Contributing to hanimo

hanimo is open source under the [Apache License 2.0](LICENSE) and welcomes
contributions. This document covers the legal & process pieces — for the
*how* of building and testing, see the project [`README.md`](README.md)
and the desktop subproject's [`hanimo-code-desktop/README.md`](hanimo-code-desktop/README.md).

한국어 안내는 본 문서 하단 §8을 참고.

---

## 1. License grant

By submitting a contribution to this repository (pull request, patch,
issue attachment containing code, or any other form), you agree that:

1. **Apache 2.0 licensing** — Your contribution is licensed under the
   Apache License, Version 2.0, the same license that governs the rest of
   the project.
2. **Authorship** — You are the author of the contribution, or you have
   sufficient rights from the actual author(s) to submit it under
   Apache 2.0.
3. **Future relicensing flexibility** — You grant the project maintainer
   the right to relicense the project (including your contribution) under
   any OSI-approved open-source license in the future, provided that the
   destination license is no more restrictive in *receiver* terms than
   Apache 2.0 + AGPL-3.0. This clause exists so that the project can adapt
   to ecosystem shifts (e.g., a future "OSI-approved successor to AGPL")
   without requiring per-contributor consent again.
4. **No warranty** — You provide the contribution "as is", with no
   warranty as to fitness for any purpose, in line with Apache 2.0 §7.

This is a lightweight contributor license grant in lieu of a separate
formal CLA infrastructure. It is not legal advice; if your employer
requires a corporate CLA, please open an issue and we will arrange one.

## 2. Sign-off (DCO)

We use the [Developer Certificate of Origin](https://developercertificate.org/)
(DCO). Every commit must carry a `Signed-off-by:` line:

```
Signed-off-by: Random J Developer <random@developer.example.org>
```

The easiest way is to sign-off automatically:

```bash
git commit -s -m "fix: handle nil context in chat stream"
```

`git commit -s` appends the `Signed-off-by:` trailer for you using your
configured `user.name` / `user.email`.

By signing off, you affirm the four DCO clauses (you wrote the patch, or
you have the right to submit it under the project's license, etc.). Full
DCO text: <https://developercertificate.org/>.

PRs without a sign-off on every commit will be asked to amend before
merge. There is no DCO bot configured (yet) — sign-off is checked
manually in review.

## 3. GPG-signed commits (optional, recommended for maintainers)

For maintainers and frequent contributors, GPG signing on top of DCO
sign-off provides cryptographic authorship evidence. Setup is in
[`docs/policy/copyright.md`](docs/policy/copyright.md) §3.

External one-off contributors are not required to GPG-sign.

## 4. Source file headers

New source files should carry the SPDX header:

```go
// Copyright 2025-2026 Kim Jiwon (김지원). All rights reserved.
// SPDX-License-Identifier: Apache-2.0
```

(Adjust the comment syntax for the file format — see
[`docs/policy/copyright.md`](docs/policy/copyright.md) §2 for examples.)

When modifying an existing file that lacks a header, you may add one as
part of your change. Mass-adding headers across the repo is reserved for
a separate dedicated PR.

## 5. Commit message conventions

We follow Conventional Commits with a few project-specific trailers
(see `~/.claude/CLAUDE.md` if you use Claude Code, or the existing
`git log` for examples):

```
feat(desktop): add hash-anchor gutter for hashline_edit

Visualizes which line ranges the agent must specify by hash.
Backend emit happens via runtime.EventsEmit("hash:anchor", ...).

Constraint: anchor flash must complete within one frame
Rejected: animation library | adds runtime cost
Confidence: high
Scope-risk: narrow
Signed-off-by: Random J Developer <random@dev.example.org>
```

Useful trailers (apply when relevant):

- `Constraint:` active constraint shaping the decision
- `Rejected:` alternative considered | reason
- `Confidence:` high / medium / low
- `Scope-risk:` narrow / moderate / broad
- `Directive:` warning for future modifiers
- `Not-tested:` edge case not covered

These are *encouraged*, not enforced. Trivial commits (typos, formatting)
can skip them.

## 6. Pull request checklist

Before opening a PR:

- [ ] Every commit has `Signed-off-by:` (run `git rebase -i` to fix
      missing ones with `--signoff`)
- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes (or new tests are added for new behavior)
- [ ] `go vet ./...` clean
- [ ] Frontend: `cd frontend && npm run build` (or
      `cd hanimo-code-desktop/frontend && npm run build`) succeeds if
      frontend was touched
- [ ] New source files carry the SPDX header
- [ ] If touching policy/license docs, the change is *additive* or
      properly RFC'd (see [`docs/policy/README.md`](docs/policy/README.md))
- [ ] If adding network calls, the [`docs/policy/telemetry-and-privacy.md`](docs/policy/telemetry-and-privacy.md)
      principles still hold

## 7. What we don't accept

- **License-incompatible code** — GPL-only / proprietary code cannot be
  merged. MIT, BSD, ISC, Apache 2.0 imports are fine; flag them in NOTICE
- **Telemetry without explicit RFC** — see telemetry-and-privacy policy
- **Vendoring entire third-party trees** without a NOTICE update
- **Logo/name uses violating trademark policy** — see
  [`docs/policy/trademark-and-naming.md`](docs/policy/trademark-and-naming.md)
- **Changes that silently break onprem build** — `make build-onprem`
  must still compile after your change

## 8. 한국어 가이드

### 라이선스 부여
PR을 제출한다는 것은 다음에 동의하는 것을 의미합니다:

1. 본 기여는 **Apache License 2.0**으로 제공됩니다 (프로젝트와 동일)
2. 본인이 저작자이거나, 저작자로부터 충분한 권리를 받았습니다
3. 메인테이너가 향후 (Apache 2.0 / AGPL-3.0 보다 *수신자에게 더 제약적이지
   않은*) OSI 승인 라이선스로 재라이선싱할 수 있는 권리를 부여합니다
4. 무보증으로 제공됩니다 (Apache 2.0 §7)

회사 차원의 CLA가 필요하면 이슈로 알려주세요.

### DCO 서명
모든 커밋에 `Signed-off-by:` 줄이 있어야 합니다. 가장 쉬운 방법:

```bash
git commit -s -m "메시지"
```

DCO 본문: <https://developercertificate.org/>

### 소스 헤더
신규 파일에 SPDX 헤더를 박아주세요:

```go
// Copyright 2025-2026 Kim Jiwon (김지원). All rights reserved.
// SPDX-License-Identifier: Apache-2.0
```

### 커밋 메시지
한국어 또는 영어 모두 환영. Conventional Commits 형식 권장
(`feat:`, `fix:`, `docs:` 등).

### 거부되는 변경
- 라이선스 비호환 코드 (GPL-only 등)
- RFC 없이 텔레메트리 추가
- 트레이드마크 정책 위반하는 로고/이름 사용
- onprem 빌드를 깨뜨리는 변경
- 정책 문서의 소급 변경

자세한 정책: [`docs/policy/README.md`](docs/policy/README.md)
