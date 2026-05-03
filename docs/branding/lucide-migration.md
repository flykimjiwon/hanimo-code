# Phosphor → Lucide 점진 이전 — hanimo-webui

> 2026-05-03 컨벤션: hanimo 생태계 4 repo는 **Lucide React 단일 아이콘
> source**로 통일. hanimo-webui의 Phosphor Icons는 도메인 단위로 점진
> 치환한다. 무리한 일괄 치환 금지.

## 현재 상태 (2026-05-03)

`@phosphor-icons/react` 사용처:

| 파일 | Phosphor 사용 | 우선순위 |
|---|---|---|
| `app/components/icons/index.jsx` | 중앙 어댑터 (다수) | **마지막** — 모든 사용처 끝난 후 |
| `app/components/chat/ChatInput.js` | `PaintBrush, UserGear` | 🥇 1순위 |
| `app/components/chat/MessageList.js` | `PaintBrush` | 🥇 1순위 |
| `app/admin/settings/page.js` | `PenNib` | 🥈 2순위 |
| `app/components/workflow/WorkflowNode.js` | 다수 | 🥉 3순위 |
| `app/components/workflow/NodePalette.js` | 다수 | 🥉 3순위 |
| `app/components/workflow/WorkflowCanvas.js` | 다수 | 🥉 3순위 |

## 도메인 단위 단계

### Phase 1 — Chat domain (chat → ChatInput, MessageList)

가장 작고 격리된 영역. 1~2시간.

체크리스트:
- [ ] `lucide-react`에서 등가 아이콘 찾기
  - `PaintBrush` → `Paintbrush2` 또는 `Brush`
  - `UserGear` → `UserCog` (조합) 또는 `Settings`
- [ ] `app/components/chat/ChatInput.js` import 교체
- [ ] `app/components/chat/MessageList.js` import 교체
- [ ] 시각 회귀 테스트 — 두 아이콘이 비슷한 크기·무게로 보이는지
- [ ] PR + 커밋 메시지: `refactor(chat): Phosphor → Lucide 이전`

### Phase 2 — Admin / Settings domain

`app/admin/settings/page.js` 등.

체크리스트:
- [ ] `PenNib` → `PenLine` 또는 `Pen`
- [ ] 다른 admin 페이지에서 phosphor 사용 검색 (`grep -r "@phosphor-icons" app/admin/`)
- [ ] 일괄 교체
- [ ] PR + 커밋 메시지: `refactor(admin): Phosphor → Lucide 이전`

### Phase 3 — Workflow domain (가장 무거움)

`app/components/workflow/` 3 파일에 phosphor 다수 사용. **가장 마지막에**
하는 게 안전. 4~8시간.

체크리스트:
- [ ] 각 파일의 phosphor import 목록 추출
  ```bash
  grep -A 20 "@phosphor-icons" app/components/workflow/*.js
  ```
- [ ] 각 phosphor 아이콘별 lucide 등가 매핑 표 작성
- [ ] *작은* 단위로 PR (한 파일씩) — 회귀 위험 최소화
- [ ] 워크플로우 노드 시각 회귀 테스트 — 노드 아이콘 인지 가능한지
- [ ] PR + 커밋 메시지: `refactor(workflow): Phosphor → Lucide 이전 (N/3)`

### Phase 4 — Central adapter cleanup

`app/components/icons/index.jsx` 정리.

체크리스트:
- [ ] phosphor import 제거
- [ ] 중앙 adapter를 lucide-react re-export로 단순화 또는 *완전 제거*
  (직접 lucide-react import 권장)
- [ ] `import` 일괄 검색으로 잔존 phosphor 사용 0건 확인
- [ ] `package.json`에서 `@phosphor-icons/react` 제거
- [ ] `npm install` → lockfile 갱신
- [ ] 빌드 확인: `npm run build`
- [ ] 시각 회귀 테스트 (전체)
- [ ] PR + 커밋 메시지: `refactor(icons): Phosphor 의존 완전 제거`

## Phosphor → Lucide 매핑 표 (시작점)

자주 사용되는 매핑. 정확한 매칭이 없으면 *유사 의미*로 대체.

| Phosphor | Lucide | 비고 |
|---|---|---|
| `Heartbeat` | `Activity` | 둘 다 펄스 라인 |
| `PaintBrush` | `Paintbrush2` | 정확한 매칭 |
| `UserGear` | `UserCog` | 사용자 + 톱니 조합 |
| `PenNib` | `PenLine` 또는 `Pen` | 펜촉 강조 시 PenLine |
| `Sparkle` | `Sparkles` | 정확한 매칭 |
| `Robot` | `Bot` | 정확한 매칭 |
| `ArrowsClockwise` | `RefreshCw` | 정확한 매칭 |
| `MagnifyingGlass` | `Search` | 정확한 매칭 |
| `ChatCircleDots` | `MessageCircle` | 점 없는 버전 |
| `Lightning` | `Zap` | 정확한 매칭 |
| `Compass` | `Compass` | 동일 이름 |
| `BookOpen` | `BookOpen` | 동일 이름 |

찾을 때: <https://lucide.dev/icons/>

## 시각 회귀 테스트

각 단계 완료 후:

1. **개발 서버 실행**: `npm run dev`
2. **각 도메인 페이지 방문**:
   - chat: `/chat` 또는 채팅창 진입
   - admin: `/admin/settings`
   - workflow: `/workflow` (또는 워크플로우 빌더 페이지)
3. **시각 체크**:
   - 아이콘 크기가 비슷한지 (Lucide는 보통 약간 작음, `size={20}` 명시 권장)
   - 색상이 정상 (currentColor 상속 잘 되는지)
   - 클릭 가능한 영역 유지
4. **screenshot diff** (옵션, playwright 사용 시):
   ```bash
   npx playwright test --update-snapshots  # 새 baseline
   ```

## 거부된 접근

- ❌ **일괄 sed 치환** — 의미 매핑이 틀릴 위험. 도메인별 수동 검토 필수
- ❌ **Phosphor 즉시 완전 제거** — 워크플로우 노드 다수 회귀 위험
- ❌ **Lucide 외 다른 아이콘 라이브러리 도입** (heroicons, react-icons 등) — 단일 source 원칙 위배
- ❌ **자체 SVG 아이콘 라이브러리 작성** — 휠 재발명, lucide가 1500+ 아이콘 제공

## 진행 보고

각 Phase 완료 시 본 문서의 체크박스 업데이트 + 진행 상태:

```
- [x] Phase 1 — Chat domain (commit c2d6627, 2026-05-03)
- [x] Phase 2 — Admin / Settings domain (commit 20a674a, 2026-05-03)
- [x] Phase 3 — Workflow domain 3 파일 일괄 (commit e9d1d8f, 2026-05-03)
       · 스코프 예상보다 작음 — 3 phosphor 아이콘만 + wrapper 패턴 동일
       · CircleArrowRight/Left + MessageSquareText, strokeWidth=1.5
- [ ] Phase 4 — Central adapter cleanup
       · app/components/icons/index.jsx 중앙 어댑터 phosphor 제거
       · package.json @phosphor-icons/react 의존성 제거
       · npm install lockfile 갱신
```

## 첫 시도 시 추천 시점

다음 두 조건 중 하나일 때:
1. **chat 컴포넌트를 만질 일이 생긴 PR** — 김에 Lucide로 교체 (Phase 1만)
2. **별도 1시간 슬롯 확보 시** — Phase 1 단독 PR

Phase 3은 *최소 반나절* 확보 후 진행. 무리하지 말 것.
