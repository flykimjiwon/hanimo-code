# Sprint 4: 킬러 기능 (추후 구현)

> 이 문서는 Sprint 1-3 완료 후 구현할 고급 기능들을 정리한 것입니다.

---

## 1. 체크포인트 + 되돌리기

**참고**: Claude Code (Esc+Esc 3모드), OpenCode (git 기반)

### 기능
- 에이전트 실행 중 자동으로 git snapshot 생성
- `Esc+Esc` 또는 `/undo`로 되돌리기 메뉴 열기
- 3가지 되돌리기 모드:
  - **코드만**: git checkout으로 파일 복원, 대화는 유지
  - **대화만**: 메시지 히스토리 롤백, 코드는 유지
  - **둘 다**: 코드 + 대화 모두 이전 상태로

### 구현 방향
- `src/core/checkpoint.ts` 신규
  - `createCheckpoint()`: git stash 또는 임시 커밋 생성
  - `listCheckpoints()`: 시간순 체크포인트 목록
  - `restoreCheckpoint(id, mode)`: 모드별 복원
- `src/tui/components/checkpoint-dialog.tsx` 신규
  - 체크포인트 목록 + 모드 선택 UI
- `src/tui/app.tsx` 수정: 더블 Esc 감지, 다이얼로그 표시

### 참고 소스
- OpenCode: git snapshot at each agent step
- Claude Code: Esc+Esc rewind with 3 restore modes

---

## 2. `/btw` 사이드 질문

**참고**: Claude Code

### 기능
- 에이전트 작업 중 빠른 질문 가능
- 메인 대화 히스토리에 포함 안 됨 (ephemeral)
- 도구 접근 없음, 순수 Q&A만
- 토큰 절약 (별도 짧은 컨텍스트)

### 구현 방향
- `/btw <질문>` 명령 추가
- `src/core/side-channel.ts` 신규
  - 별도 모델 호출 (메인 대화 컨텍스트 미포함)
  - 결과를 일시적 메시지로 표시 (저장 안 함)
- `src/tui/components/btw-popup.tsx` 신규
  - 팝업으로 질문/답변 표시, 자동 사라짐

### 참고 소스
- Claude Code: /btw command, no tool access, no history pollution
- 토큰 비용 ~50% 절감 효과

---

## 3. `/fork` 대화 분기

**참고**: Codex CLI

### 기능
- 현재 대화를 복제해서 새 분기 생성
- 원본 대화는 보존
- 다른 접근법 탐색 후 비교 가능

### 구현 방향
- `/fork` 명령: 현재 세션 deep copy → 새 세션 ID
- `/fork [message-index]`: 특정 시점부터 분기
- 세션 목록에서 parent-child 관계 표시
- `src/core/session-fork.ts` 신규

### 참고 소스
- Codex CLI: /fork conversation branching, fork from any earlier message

---

## 4. `/review` 빌트인 코드 리뷰

**참고**: Codex CLI

### 기능
- 현재 diff를 분석해서 우선순위별 피드백
- working tree 수정 없음 (읽기만)
- 카테고리: 버그, 보안, 성능, 스타일, 로직

### 구현 방향
- `/review` 명령: `git diff` 수집 → 리뷰 프롬프트로 전달
- `/review --staged`: staged 변경만 리뷰
- `/review <file>`: 특정 파일만 리뷰
- `src/core/code-review.ts` 신규
  - diff 수집 → 구조화된 리뷰 프롬프트 생성
  - 결과: severity (critical/warning/info) + 라인 번호 + 설명

### 참고 소스
- Codex CLI: /review built-in code reviewer with prioritized findings

---

## 5. 추가 고려 기능

| 기능 | 출처 | 우선순위 |
|------|------|---------|
| Watch 모드 (`// AI!` 코멘트) | Aider | 중 |
| Vim 모드 (전체 모션) | Claude Code | 중 |
| 이미지 붙여넣기 (ctrl+v) | Claude Code/Codex | 중 |
| 프롬프트 제안 (git 기반) | Claude Code | 낮 |
| 음성 입력 (push-to-talk) | Claude Code | 낮 |
| 테마 30종 + 라이브 프리뷰 | OpenCode | 낮 |
| Tree-sitter 레포 맵 | Aider | 낮 |
| Mid-turn injection (Enter) | Codex | 중 |
| 자동 린트/테스트 루프 | Aider | 중 |
| 세션 공유 링크 | OpenCode | 낮 |

---

## 타임라인

Sprint 1-3 완료 후, 사용자 피드백 기반으로 우선순위 재조정.
예상: Sprint 4 시작 시점은 Sprint 3 완료 후 1-2주 내.
