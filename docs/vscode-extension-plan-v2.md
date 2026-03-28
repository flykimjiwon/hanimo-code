# modol VS Code Extension — 개발 & 배포 플랜

## 개요

modol은 이미 터미널(TUI/Text)에서 완전히 동작합니다.
VS Code 확장은 **부가 채널**로, 기존 core 엔진을 재사용하면서 VS Code UX를 추가합니다.

---

## Phase 1: 터미널 통합 (1일)

비용: $0 | 난이도: 하

### 무엇을 만드는가
VS Code에서 단축키 하나로 modol을 실행하는 최소 확장.

### 구현

```
modol-vscode/
├── package.json          # Extension manifest
├── src/
│   └── extension.ts      # 20줄짜리 진입점
├── tsconfig.json
└── README.md
```

**extension.ts (전체 코드):**
```typescript
import * as vscode from 'vscode';

export function activate(context: vscode.ExtensionContext) {
  // 커맨드 1: 터미널에서 modol 실행
  context.subscriptions.push(
    vscode.commands.registerCommand('modol.open', () => {
      const terminal = vscode.window.createTerminal({
        name: 'modol 🐶',
        cwd: vscode.workspace.workspaceFolders?.[0]?.uri.fsPath,
      });
      terminal.sendText('modol');
      terminal.show();
    })
  );

  // 커맨드 2: 선택 영역을 modol에 전달
  context.subscriptions.push(
    vscode.commands.registerCommand('modol.ask', () => {
      const editor = vscode.window.activeTextEditor;
      const selection = editor?.document.getText(editor.selection);
      const terminal = vscode.window.createTerminal({
        name: 'modol 🐶',
        cwd: vscode.workspace.workspaceFolders?.[0]?.uri.fsPath,
      });
      if (selection) {
        terminal.sendText(`modol "${selection.replace(/"/g, '\\"')}"`);
      } else {
        terminal.sendText('modol');
      }
      terminal.show();
    })
  );
}

export function deactivate() {}
```

**package.json (Extension manifest):**
```json
{
  "name": "modol",
  "displayName": "modol — AI Coding Agent 🐶",
  "description": "Terminal AI coding agent that works with any LLM",
  "version": "0.1.0",
  "publisher": "flykimjiwon",
  "engines": { "vscode": "^1.85.0" },
  "categories": ["AI", "Programming Languages", "Other"],
  "icon": "icon.png",
  "repository": {
    "type": "git",
    "url": "https://github.com/flykimjiwon/dev_anywhere"
  },
  "activationEvents": [],
  "main": "./out/extension.js",
  "contributes": {
    "commands": [
      {
        "command": "modol.open",
        "title": "modol: Open Terminal",
        "icon": "$(terminal)"
      },
      {
        "command": "modol.ask",
        "title": "modol: Ask about selection"
      }
    ],
    "keybindings": [
      {
        "command": "modol.open",
        "key": "ctrl+shift+m",
        "mac": "cmd+shift+m"
      },
      {
        "command": "modol.ask",
        "key": "ctrl+shift+a",
        "mac": "cmd+shift+a"
      }
    ]
  },
  "scripts": {
    "compile": "tsc -p ./",
    "package": "vsce package",
    "publish": "vsce publish"
  },
  "devDependencies": {
    "@types/vscode": "^1.85.0",
    "typescript": "^5.5.0",
    "@vscode/vsce": "^3.0.0"
  }
}
```

### 단축키
| 키 | 동작 |
|---|---|
| `Cmd+Shift+M` | modol 터미널 열기 |
| `Cmd+Shift+A` | 선택 영역을 modol에 질문 |

---

## Phase 2: Sidebar Webview (1주)

비용: $0 | 난이도: 중

### 무엇을 만드는가
VS Code 사이드바에 modol 채팅 UI를 내장. 에디터와 양방향 연동.

### 추가 기능
- 사이드바에 채팅 패널
- 에디터에서 코드 선택 → "modol에 물어보기" 컨텍스트 메뉴
- modol 응답에서 "Apply to editor" 버튼
- 파일 변경 사항 실시간 미리보기 (diff)

### 구조
```
modol-vscode/
├── package.json
├── src/
│   ├── extension.ts          # 진입점 + 커맨드 등록
│   ├── sidebar-provider.ts   # Webview 사이드바
│   ├── bridge.ts             # modol core ↔ VS Code 통신
│   └── diff-viewer.ts        # 파일 변경 diff 미리보기
├── webview/
│   ├── index.html
│   ├── app.tsx               # React 채팅 UI
│   └── styles.css
├── media/
│   └── icon.png              # 🐶 아이콘
└── tsconfig.json
```

### package.json 추가 (Phase 2)
```json
{
  "contributes": {
    "viewsContainers": {
      "activitybar": [{
        "id": "modol",
        "title": "modol",
        "icon": "media/icon.png"
      }]
    },
    "views": {
      "modol": [{
        "type": "webview",
        "id": "modol.chat",
        "name": "Chat"
      }]
    },
    "menus": {
      "editor/context": [{
        "command": "modol.ask",
        "group": "modol",
        "when": "editorHasSelection"
      }]
    }
  }
}
```

---

## Phase 3: 고급 기능 (2주)

### Inline Suggestions (코드 자동완성)
```typescript
// VS Code의 InlineCompletionItemProvider 구현
// modol이 현재 커서 위치의 코드를 읽고 다음 코드를 제안
vscode.languages.registerInlineCompletionItemProvider(
  { pattern: '**' },
  new ModolCompletionProvider()
);
```

### Code Actions (Quick Fix)
```typescript
// TypeScript 에러 위에 "modol: Fix this" 전구 아이콘
vscode.languages.registerCodeActionsProvider(
  { language: 'typescript' },
  new ModolCodeActionProvider()
);
```

### 파일 변경 미리보기
```typescript
// modol이 파일을 수정하기 전에 diff를 보여주고 승인 요청
vscode.commands.executeCommand('vscode.diff', originalUri, modifiedUri, 'modol changes');
```

---

## 배포 & 등록 플랜

### Step 1: VS Code Marketplace 계정 생성

```bash
# 1. Azure DevOps 계정 생성 (무료)
#    https://dev.azure.com 에서 가입

# 2. Personal Access Token (PAT) 생성
#    Azure DevOps → User Settings → Personal Access Tokens
#    - Name: vscode-marketplace
#    - Organization: All accessible organizations
#    - Scopes: Marketplace (Manage)
#    - 유효기간: 1년

# 3. Publisher 생성
#    https://marketplace.visualstudio.com/manage/createpublisher
#    - Publisher ID: flykimjiwon
#    - Display Name: 김지원
```

### Step 2: 확장 패키징

```bash
# vsce 설치
npm install -g @vscode/vsce

# 패키징 (.vsix 파일 생성)
cd modol-vscode
vsce package
# → modol-0.1.0.vsix 생성됨

# 로컬 테스트
code --install-extension modol-0.1.0.vsix
```

### Step 3: Marketplace 배포

```bash
# 첫 배포
vsce publish
# PAT 입력하면 자동 업로드

# 이후 버전 업데이트
vsce publish patch  # 0.1.0 → 0.1.1
vsce publish minor  # 0.1.1 → 0.2.0
```

### Step 4: 자동 배포 (CI/CD)

```yaml
# .github/workflows/vscode-publish.yml
name: Publish VS Code Extension
on:
  push:
    tags: ['vscode-v*']

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: 20 }
      - run: npm ci
      - run: npm run compile
      - run: npx vsce publish
        env:
          VSCE_PAT: ${{ secrets.VSCE_PAT }}
```

### Step 5: Open VSX 등록 (선택)

VS Code Marketplace 외에 Open VSX (오픈소스 마켓)에도 등록하면 Codium, Gitpod 등에서도 설치 가능.

```bash
npm install -g ovsx
ovsx publish modol-0.1.0.vsix -p <OPEN_VSX_TOKEN>
```

---

## 비용 요약

| 항목 | 비용 |
|---|---|
| Azure DevOps 계정 | 무료 |
| VS Code Marketplace 등록 | 무료 |
| Open VSX 등록 | 무료 |
| Apple Developer (Tauri 앱, 나중에) | $99/yr |
| **총 비용** | **$0** (현 시점) |

---

## Marketplace SEO & 마케팅

### Marketplace 리스팅 최적화

```
이름: modol — AI Coding Agent 🐶
태그: ai, coding-assistant, ollama, gpt, claude, terminal, agent
카테고리: AI, Programming Languages

설명 (첫 3줄이 가장 중요):
"modol은 14개 LLM 프로바이더를 지원하는 터미널 AI 코딩 에이전트입니다.
Ollama 로컬 모델로 무료 사용 가능. Claude Code, Cursor의 오픈소스 대안.
Cmd+Shift+M으로 바로 시작하세요."
```

### 차별화 포인트 (Marketplace 기준)
| 경쟁 | 가격 | 프로바이더 | 로컬 모델 |
|---|---|---|---|
| GitHub Copilot | $10/월 | 1개 | X |
| Cursor | $20/월 | 2개 | X |
| Continue | 무료 | 5개 | O |
| Cody | 무료/유료 | 2개 | X |
| **modol** | **무료** | **14개** | **O (4개)** |

---

## 타임라인

```
Week 1:
├── Day 1: Phase 1 완성 (터미널 통합) + Marketplace 계정
├── Day 2-3: Phase 2 시작 (Sidebar Webview)
├── Day 4-5: Phase 2 완성 + 테스트
├── Day 6: 패키징 + 로컬 테스트
└── Day 7: Marketplace 배포 v0.1.0

Week 2:
├── Phase 3 시작 (Inline Suggestions)
├── 사용자 피드백 수집
└── v0.2.0 배포

Week 3-4:
├── Code Actions (Quick Fix)
├── 파일 변경 diff 미리보기
└── v0.3.0 배포
```

---

## 폴더 구조 제안

```
~/Desktop/kimjiwon/
├── dev_anywhere/          # modol core (현재)
│   ├── src/               # CLI + TUI + 엔진
│   └── package.json       # "modol" npm 패키지
│
├── modol-vscode/          # VS Code 확장 (새로 생성)
│   ├── src/extension.ts
│   ├── webview/
│   └── package.json       # VS Code extension
│
└── modol.app/             # 웹사이트 (나중에)
    ├── landing/           # 랜딩 페이지
    └── docs/              # 문서 사이트
```

---

## 체크리스트

### Phase 1 배포 전
- [ ] Azure DevOps 계정 생성
- [ ] Personal Access Token 발급
- [ ] Publisher ID 생성 (flykimjiwon)
- [ ] extension.ts 작성
- [ ] package.json 작성 (manifest)
- [ ] icon.png 준비 (128x128, 🐶 모돌 아이콘)
- [ ] README.md 작성 (Marketplace에 표시됨)
- [ ] CHANGELOG.md 작성
- [ ] LICENSE 파일 확인 (MIT)
- [ ] `vsce package` → .vsix 생성
- [ ] 로컬 테스트 (`code --install-extension`)
- [ ] `vsce publish` → Marketplace 등록

### Phase 2 배포 전
- [ ] Sidebar Webview 구현
- [ ] 에디터 컨텍스트 메뉴 연동
- [ ] Apply to Editor 기능
- [ ] 다크/라이트 테마 지원
- [ ] v0.2.0 배포
