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

### Phase 4 — Central adapter cleanup (별도 PR 권장)

`app/components/icons/index.jsx` 정리. **Phase 1·2·3 대비 위험 계수 큼 —
별도 dedicated PR + 시각 회귀 테스트 확보 후 진행 권장.**

#### 왜 Phase 4가 위험한가

- 100+ 아이콘 매핑이 들어있는 중앙 어댑터 — 한 매핑이라도 틀리면 *수십
  컴포넌트가 동시에 깨짐* (전 화면 회귀)
- Phosphor → Lucide 이름이 *대부분 비슷*하지만 정확히 같지 않은 경우
  많음 (예: Phosphor `Heartbeat` ↔ Lucide `Activity`, `CircleNotch` ↔
  `Loader2`, `MagnifyingGlass` ↔ `Search`, `Prohibit` ↔ `PowerOff`)
- 시각 회귀: Phosphor의 `weight="light"/"duotone"`가 Lucide의 `strokeWidth`
  로 정확히 매칭되지 않음. 일부 아이콘은 시각적으로 *살짝 무거워* 보임
- shadcn/ui 프리미티브 alias가 어댑터에 박혀있어 dialog/select/dropdown 등
  핵심 UI 흐름이 모두 이 어댑터를 의존

#### 권장 Phase 4 진행 절차

1. **사전 작업** (별도 PR 또는 작업 셋업)
   - [ ] Lucide-react 정식 export 이름 매니페스트 추출:
     ```bash
     node -e "console.log(Object.keys(require('lucide-react')).sort().join('\n'))" \
       > lucide-export-list.txt
     ```
   - [ ] 어댑터의 phosphor 사용 100+ 매핑 표 작성:
     `Phosphor 이름 | Lucide 이름 | 시각 차이 노트` 3 컬럼
   - [ ] 매핑 표 메인테이너 검토 (특히 `weight` → `strokeWidth` 변환 정책)
   - [ ] Storybook 또는 *시각 비교 페이지* 1회성 작성: 모든 100+ 아이콘을
     before/after 양쪽 렌더 → screenshot 저장

2. **본 PR**
   - [ ] 어댑터 파일을 lucide-react re-export로 *완전 교체*
   - [ ] `w()` wrapper 함수 제거 (weight prop 무시 의미를 잃음)
   - [ ] 일부 alias 처리: `LucideImage`, `CheckCircle2`, `Code2`, `Table2`,
     `PieChartIcon`, shadcn/ui 프리미티브 11개 (`CheckIcon`, `ChevronDownIcon`
     등)
   - [ ] `package.json`에서 `@phosphor-icons/react` 제거
   - [ ] `package-lock.json` / `yarn.lock` 갱신
   - [ ] **빌드 검증 필수**: `npm run build` (Next.js 프로덕션 빌드)
   - [ ] **타입 검증**: `npm run typecheck` 또는 `tsc --noEmit`
   - [ ] **시각 회귀 검증**: 사전 작업의 Storybook/스크린샷 대비
   - [ ] 첫 PR은 *해결 못한 시각 차이*를 *명시*하는 issue 동반
     (예: "워크플로우 노드 아이콘이 Phosphor 대비 12% 두꺼워 보임 — 디자인
     팀 확인 후 별도 fix")

3. **후속 PR**
   - [ ] 시각 회귀 fix (개별 strokeWidth 조정)
   - [ ] *어댑터 자체 폐기* — 남은 `@/components/icons` import를 직접
     `lucide-react`로 마이그레이션 (long-term)

#### 임시 회피 방법

Phase 4가 진행되기 전까지 어댑터는 *동작*함 (Phosphor 의존). 다만:

- 신규 코드는 *직접 lucide-react* 사용 (어댑터 import 금지)
- `package.json`의 `@phosphor-icons/react`는 *유지* — Phase 4 PR에서 제거
- 디스크 사용량 약 4MB 추가, 번들 사이즈 미세 영향 (tree-shaking 안 되는
  부분 있을 수 있음)

#### 체크리스트 요약

- [ ] Phase 4a: Lucide export 매니페스트 + 100+ 매핑 표 작성
- [ ] Phase 4a: 시각 비교 페이지 / Storybook
- [ ] Phase 4b: 어댑터 lucide-react re-export 교체
- [ ] Phase 4b: 빌드/타입 검증
- [ ] Phase 4b: 시각 회귀 첫 라운드 fix
- [ ] Phase 4c: package.json phosphor dep 제거
- [ ] Phase 4c: lockfile 갱신
- [ ] Phase 4d (long-term): 어댑터 자체 폐기 + 직접 lucide import 마이그레이션

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
