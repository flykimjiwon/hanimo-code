# modol 다음 세션 TODO

## 우선순위 1: 멀티 엔드포인트 TUI 연결

### 현재 상태
- EndpointManager 엔진 구현됨 (`src/providers/endpoint-manager.ts`)
- config.json에 endpoints 배열 지원됨
- `/endpoint add/list/remove` TUI 커맨드 구현됨
- **아직 TUI 모델 메뉴에 실제 연결 안 됨**

### 구현할 것

1. **모델 메뉴에 엔드포인트 표시**
   - 모델명 옆에 `@엔드포인트명` 표시
   - 예: `qwen3:8b @local`, `qwen3:8b @remote`, `gpt-oss:20b @dgx`

2. **같은 모델 여러 엔드포인트 → 서브메뉴**
   ```
   모델 선택: qwen3:8b (2 endpoints)
   → 1. @local (http://localhost:11434) — priority:10
     2. @remote (http://192.168.1.100:11434) — priority:3
     3. Auto (round-robin)
   ```

3. **프로바이더 전환 시 연결 검증**
   - API 키 없는 클라우드 프로바이더 → 경고 메시지
   - 엔드포인트 연결 실패 → 에러 메시지 + 다른 엔드포인트 제안

4. **app.tsx에서 EndpointManager 초기화**
   ```typescript
   // startup에서:
   const epManager = new EndpointManager();
   await epManager.loadEndpoints(config.endpoints);

   // 모델 메뉴에서:
   const models = epManager.getAllModels(); // 엔드포인트별 모델

   // 모델 선택 시:
   const ep = epManager.getEndpoint(modelId); // 라운드로빈 or 선택
   const modelInstance = getModel(ep.provider, modelId, { baseURL: ep.endpointURL, apiKey: ep.apiKey });
   ```

5. **StatusBar에 현재 엔드포인트 표시**
   ```
   modol │ ollama/qwen3:8b @local │ 🐶 Super Modol │ tools:ON
   ```

## 우선순위 2: /auto 커맨드 실제 동작 연결

### 현재 상태
- auto-loop.ts 엔진 구현됨
- /auto 커맨드는 메시지만 표시하고 실제 동작 안 함
- `sendAutoMessage`가 CommandContext에 없음

### 구현할 것
- CommandContext에 `sendAutoMessage` 추가
- app.tsx에서 auto-loop 실행 연결
- 자율 루프 중 TUI에 진행 상황 표시

## 우선순위 3: 스킬 시스템

### 구현할 것
- `~/.modol/skills/*.md` 파일 자동 로드
- 시스템 프롬프트에 주입
- `/skill list` — 로드된 스킬 목록
- `/skill add <url>` — URL에서 스킬 파일 생성 (webfetch 활용)

## 우선순위 4: text-mode i18n

### 구현할 것
- text-mode.ts에 한국어/영어 분기 추가
- onboarding.ts에 언어 선택 추가

## 우선순위 5: TUI 추가 개선

- 프로바이더 전환 시 연결 검증 (API 키 확인)
- StatusBar에 openai라고 뜨는데 실제론 ollama 치는 버그 수정
- 테마 설명 i18n (현재 한국어만)
- leader key 힌트 i18n
