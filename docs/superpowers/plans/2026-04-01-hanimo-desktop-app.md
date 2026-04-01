# hanimo Desktop App (Phase 1 MVP) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Tauri v2 래퍼로 hanimo를 크로스플랫폼 데스크톱 앱으로 만든다. 기존 코드는 Node.js sidecar로 그대로 재사용.

**Architecture:** React 프론트엔드가 Tauri IPC를 통해 Rust 셸과 통신하고, Rust 셸은 Node.js sidecar(기존 `headless.ts`)를 stdio로 제어한다. 프론트엔드는 VSCode 스타일 레이아웃(사이드바 + 파일트리 + Monaco 에디터 + 채팅 패널)이며, 다크/라이트 테마를 지원한다.

**Tech Stack:** Tauri v2, React 18, TypeScript, Vite, Monaco Editor, Zustand, TailwindCSS

**Spec:** `docs/superpowers/specs/2026-04-01-hanimo-desktop-app-design.md`

---

## File Structure

```
desktop/
├── src-tauri/
│   ├── src/
│   │   ├── main.rs              # Tauri entry, window creation
│   │   ├── sidecar.rs           # Node.js sidecar process management
│   │   └── commands.rs          # Tauri IPC commands (invoke handlers)
│   ├── Cargo.toml               # Rust dependencies
│   └── tauri.conf.json          # Tauri build config, sidecar declaration
├── src/
│   ├── main.tsx                 # React entry point
│   ├── App.tsx                  # Root layout component
│   ├── stores/
│   │   ├── chat-store.ts        # Chat messages state (Zustand)
│   │   ├── editor-store.ts      # Open files, active tab state
│   │   ├── theme-store.ts       # Dark/light theme state
│   │   └── onboarding-store.ts  # Onboarding wizard state
│   ├── components/
│   │   ├── layout/
│   │   │   ├── Sidebar.tsx      # Icon sidebar (left 40px)
│   │   │   ├── FileTree.tsx     # File explorer panel
│   │   │   ├── EditorArea.tsx   # Tab bar + Monaco editor container
│   │   │   └── ChatPanel.tsx    # Right chat panel
│   │   ├── chat/
│   │   │   ├── MessageList.tsx  # Chat message list renderer
│   │   │   ├── MessageBubble.tsx # Single message (user/assistant/tool)
│   │   │   └── ChatInput.tsx    # Message input with send button
│   │   ├── editor/
│   │   │   ├── TabBar.tsx       # Editor tab bar
│   │   │   └── MonacoWrapper.tsx # Monaco editor wrapper
│   │   └── onboarding/
│   │       └── OnboardingWizard.tsx # Step-by-step setup wizard
│   ├── hooks/
│   │   ├── use-sidecar.ts       # Hook for sidecar IPC communication
│   │   └── use-file-system.ts   # Hook for file system operations via Tauri
│   ├── lib/
│   │   ├── ipc.ts               # Tauri invoke/listen wrappers
│   │   └── theme.ts             # Theme definitions (dark/light)
│   └── styles/
│       └── globals.css          # TailwindCSS imports + base styles
├── index.html                   # Vite HTML entry
├── package.json                 # Frontend dependencies
├── tsconfig.json                # TypeScript config
├── vite.config.ts               # Vite config for Tauri
├── tailwind.config.ts           # TailwindCSS config
└── postcss.config.js            # PostCSS for Tailwind
```

---

## Task 1: Tauri 프로젝트 초기화

**Files:**
- Create: `desktop/src-tauri/src/main.rs`
- Create: `desktop/src-tauri/Cargo.toml`
- Create: `desktop/src-tauri/tauri.conf.json`
- Create: `desktop/src-tauri/capabilities/default.json`
- Create: `desktop/package.json`
- Create: `desktop/index.html`
- Create: `desktop/vite.config.ts`
- Create: `desktop/tsconfig.json`
- Create: `desktop/src/main.tsx`

- [ ] **Step 1: Tauri CLI 설치 확인**

Run: `cargo install tauri-cli --version "^2.0" 2>/dev/null; cargo tauri --version`
Expected: `tauri-cli 2.x.x`

- [ ] **Step 2: desktop/package.json 생성**

```json
{
  "name": "hanimo-desktop",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview",
    "tauri": "tauri"
  },
  "dependencies": {
    "@tauri-apps/api": "^2.0.0",
    "@tauri-apps/plugin-shell": "^2.0.0",
    "@tauri-apps/plugin-fs": "^2.0.0",
    "react": "^18.3.0",
    "react-dom": "^18.3.0",
    "zustand": "^5.0.0",
    "@monaco-editor/react": "^4.6.0",
    "tailwindcss": "^4.0.0",
    "@tailwindcss/vite": "^4.0.0"
  },
  "devDependencies": {
    "@types/react": "^18.3.0",
    "@types/react-dom": "^18.3.0",
    "@tauri-apps/cli": "^2.0.0",
    "typescript": "^5.5.0",
    "vite": "^6.0.0",
    "@vitejs/plugin-react": "^4.0.0"
  }
}
```

- [ ] **Step 3: vite.config.ts 생성**

```typescript
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

const host = process.env.TAURI_DEV_HOST;

export default defineConfig({
  plugins: [react(), tailwindcss()],
  clearScreen: false,
  server: {
    port: 1420,
    strictPort: true,
    host: host || false,
    hmr: host ? { protocol: "ws", host, port: 1421 } : undefined,
    watch: { ignored: ["**/src-tauri/**"] },
  },
});
```

- [ ] **Step 4: tsconfig.json 생성**

```json
{
  "compilerOptions": {
    "target": "ES2021",
    "useDefineForClassFields": true,
    "lib": ["ES2021", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "forceConsistentCasingInFileNames": true
  },
  "include": ["src"]
}
```

- [ ] **Step 5: index.html 생성**

```html
<!doctype html>
<html lang="ko">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>hanimo</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
```

- [ ] **Step 6: React entry point 생성**

Create `desktop/src/main.tsx`:

```tsx
import React from "react";
import ReactDOM from "react-dom/client";
import "./styles/globals.css";

function App() {
  return <div className="h-screen bg-neutral-900 text-white flex items-center justify-center">
    <h1 className="text-2xl">hanimo 🐶</h1>
  </div>;
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
```

Create `desktop/src/styles/globals.css`:

```css
@import "tailwindcss";

html, body, #root {
  margin: 0;
  padding: 0;
  height: 100%;
  overflow: hidden;
}
```

- [ ] **Step 7: Tauri Rust 설정 — Cargo.toml**

```toml
[package]
name = "hanimo-desktop"
version = "0.1.0"
edition = "2021"

[lib]
name = "hanimo_desktop_lib"
crate-type = ["staticlib", "cdylib", "rlib"]

[build-dependencies]
tauri-build = { version = "2", features = [] }

[dependencies]
tauri = { version = "2", features = [] }
tauri-plugin-shell = "2"
tauri-plugin-fs = "2"
serde = { version = "1", features = ["derive"] }
serde_json = "1"
```

Create `desktop/src-tauri/build.rs`:

```rust
fn main() {
    tauri_build::build()
}
```

- [ ] **Step 8: Tauri main.rs 생성**

```rust
#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_fs::init())
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
```

Create `desktop/src-tauri/src/main.rs`:

```rust
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

fn main() {
    hanimo_desktop_lib::run()
}
```

Create `desktop/src-tauri/src/lib.rs` (the `run` function lives here):

```rust
#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_fs::init())
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
```

- [ ] **Step 9: tauri.conf.json 생성**

```json
{
  "$schema": "https://raw.githubusercontent.com/nicedavid98/tauri-schema/refs/heads/main/tauri.conf.json.schema.json",
  "productName": "hanimo",
  "version": "0.1.0",
  "identifier": "com.hanimo.app",
  "build": {
    "frontendDist": "../dist",
    "devUrl": "http://localhost:1420",
    "beforeDevCommand": "npm run dev",
    "beforeBuildCommand": "npm run build"
  },
  "app": {
    "title": "hanimo 🐶",
    "windows": [
      {
        "title": "hanimo 🐶",
        "width": 1280,
        "height": 800,
        "minWidth": 900,
        "minHeight": 600,
        "decorations": true,
        "resizable": true
      }
    ],
    "security": {
      "csp": null
    }
  },
  "bundle": {
    "active": true,
    "targets": "all",
    "icon": [
      "icons/32x32.png",
      "icons/128x128.png",
      "icons/128x128@2x.png",
      "icons/icon.icns",
      "icons/icon.ico"
    ]
  }
}
```

Create `desktop/src-tauri/capabilities/default.json`:

```json
{
  "$schema": "https://raw.githubusercontent.com/nicedavid98/tauri-schema/refs/heads/main/capabilities.schema.json",
  "identifier": "default",
  "description": "Default capabilities for hanimo desktop",
  "windows": ["main"],
  "permissions": [
    "core:default",
    "shell:allow-execute",
    "shell:allow-stdin-write",
    "fs:default",
    "fs:allow-read",
    "fs:allow-write"
  ]
}
```

- [ ] **Step 10: 의존성 설치 및 빌드 테스트**

Run:
```bash
cd desktop && npm install
cd src-tauri && cargo build
```
Expected: 빌드 성공, 에러 없음

- [ ] **Step 11: 개발 서버 실행 테스트**

Run: `cd desktop && npm run tauri dev`
Expected: 네이티브 윈도우에 "hanimo 🐶" 텍스트 표시

- [ ] **Step 12: Commit**

```bash
git add desktop/
git commit -m "feat(desktop): initialize Tauri v2 project with React + Vite

Scaffold desktop app with Tauri v2 wrapper, React 18 frontend,
TailwindCSS, and Vite build system. Minimal hello-world window."
```

---

## Task 2: Node.js Sidecar 통합

**Files:**
- Create: `desktop/src-tauri/src/sidecar.rs`
- Create: `desktop/src-tauri/src/commands.rs`
- Modify: `desktop/src-tauri/src/lib.rs`
- Modify: `desktop/src-tauri/tauri.conf.json`
- Create: `desktop/src/lib/ipc.ts`
- Create: `desktop/src/hooks/use-sidecar.ts`

기존 `src/headless.ts`가 JSON stdin/stdout 프로토콜을 이미 제공하므로, sidecar로 이를 그대로 활용한다.

- [ ] **Step 1: sidecar.rs — Node.js 프로세스 관리**

```rust
use std::io::{BufRead, BufReader, Write};
use std::process::{Child, Command, Stdio};
use std::sync::Mutex;

pub struct Sidecar {
    process: Mutex<Option<Child>>,
}

impl Sidecar {
    pub fn new() -> Self {
        Self {
            process: Mutex::new(None),
        }
    }

    pub fn start(&self, node_path: &str, script_path: &str, args: &[&str]) -> Result<(), String> {
        let child = Command::new(node_path)
            .arg(script_path)
            .args(args)
            .stdin(Stdio::piped())
            .stdout(Stdio::piped())
            .stderr(Stdio::piped())
            .spawn()
            .map_err(|e| format!("Failed to start sidecar: {}", e))?;

        let mut proc = self.process.lock().map_err(|e| e.to_string())?;
        *proc = Some(child);
        Ok(())
    }

    pub fn send(&self, message: &str) -> Result<(), String> {
        let mut proc = self.process.lock().map_err(|e| e.to_string())?;
        if let Some(ref mut child) = *proc {
            if let Some(ref mut stdin) = child.stdin {
                writeln!(stdin, "{}", message).map_err(|e| format!("Write failed: {}", e))?;
                stdin.flush().map_err(|e| format!("Flush failed: {}", e))?;
            }
        }
        Ok(())
    }

    pub fn read_line(&self) -> Result<Option<String>, String> {
        let mut proc = self.process.lock().map_err(|e| e.to_string())?;
        if let Some(ref mut child) = *proc {
            if let Some(ref mut stdout) = child.stdout {
                let mut reader = BufReader::new(stdout);
                let mut line = String::new();
                match reader.read_line(&mut line) {
                    Ok(0) => Ok(None),
                    Ok(_) => Ok(Some(line.trim().to_string())),
                    Err(e) => Err(format!("Read failed: {}", e)),
                }
            } else {
                Ok(None)
            }
        } else {
            Err("Sidecar not running".to_string())
        }
    }

    pub fn stop(&self) -> Result<(), String> {
        let mut proc = self.process.lock().map_err(|e| e.to_string())?;
        if let Some(ref mut child) = *proc {
            child.kill().map_err(|e| format!("Kill failed: {}", e))?;
        }
        *proc = None;
        Ok(())
    }
}
```

- [ ] **Step 2: commands.rs — Tauri IPC 커맨드**

```rust
use crate::sidecar::Sidecar;
use tauri::State;

#[tauri::command]
pub async fn start_sidecar(sidecar: State<'_, Sidecar>) -> Result<String, String> {
    // In production, node binary is bundled; in dev, use system node
    let node_path = "node";
    // Path to the hanimo headless entry point (built dist)
    let script_path = "../dist/cli.js";
    let args = &["--headless"];

    sidecar.start(node_path, script_path, args)?;
    Ok("Sidecar started".to_string())
}

#[tauri::command]
pub async fn send_prompt(sidecar: State<'_, Sidecar>, content: String) -> Result<(), String> {
    let message = serde_json::json!({
        "type": "prompt",
        "content": content
    })
    .to_string();
    sidecar.send(&message)
}

#[tauri::command]
pub async fn stop_sidecar(sidecar: State<'_, Sidecar>) -> Result<(), String> {
    sidecar.stop()
}
```

- [ ] **Step 3: lib.rs 업데이트 — sidecar와 commands 등록**

```rust
mod commands;
mod sidecar;

use sidecar::Sidecar;

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_fs::init())
        .manage(Sidecar::new())
        .invoke_handler(tauri::generate_handler![
            commands::start_sidecar,
            commands::send_prompt,
            commands::stop_sidecar,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
```

- [ ] **Step 4: Tauri 이벤트 기반 스트리밍으로 개선 — commands.rs 수정**

sidecar stdout을 비동기로 읽어 Tauri 이벤트로 프론트엔드에 전달:

```rust
use tauri::{AppHandle, Emitter, State};
use crate::sidecar::Sidecar;
use std::io::{BufRead, BufReader};
use std::thread;

#[tauri::command]
pub async fn start_sidecar(
    app: AppHandle,
    sidecar: State<'_, Sidecar>,
) -> Result<String, String> {
    let node_path = "node";
    let script_path = "../dist/cli.js";
    let args = &["--headless"];

    sidecar.start(node_path, script_path, args)?;

    // Spawn thread to read stdout and emit events
    let stdout = sidecar.take_stdout()?;
    thread::spawn(move || {
        let reader = BufReader::new(stdout);
        for line in reader.lines() {
            match line {
                Ok(text) => {
                    let _ = app.emit("sidecar-event", text);
                }
                Err(_) => break,
            }
        }
    });

    Ok("Sidecar started".to_string())
}

#[tauri::command]
pub async fn send_prompt(sidecar: State<'_, Sidecar>, content: String) -> Result<(), String> {
    let message = serde_json::json!({
        "type": "prompt",
        "content": content
    })
    .to_string();
    sidecar.send(&message)
}

#[tauri::command]
pub async fn stop_sidecar(sidecar: State<'_, Sidecar>) -> Result<(), String> {
    sidecar.stop()
}
```

Update `sidecar.rs` to add `take_stdout`:

```rust
use std::io::Write;
use std::process::{Child, ChildStdout, Command, Stdio};
use std::sync::Mutex;

pub struct Sidecar {
    process: Mutex<Option<Child>>,
    stdout: Mutex<Option<ChildStdout>>,
}

impl Sidecar {
    pub fn new() -> Self {
        Self {
            process: Mutex::new(None),
            stdout: Mutex::new(None),
        }
    }

    pub fn start(&self, node_path: &str, script_path: &str, args: &[&str]) -> Result<(), String> {
        let mut child = Command::new(node_path)
            .arg(script_path)
            .args(args)
            .stdin(Stdio::piped())
            .stdout(Stdio::piped())
            .stderr(Stdio::piped())
            .spawn()
            .map_err(|e| format!("Failed to start sidecar: {}", e))?;

        let stdout = child.stdout.take();
        let mut proc = self.process.lock().map_err(|e| e.to_string())?;
        *proc = Some(child);
        let mut out = self.stdout.lock().map_err(|e| e.to_string())?;
        *out = stdout;
        Ok(())
    }

    pub fn take_stdout(&self) -> Result<ChildStdout, String> {
        let mut out = self.stdout.lock().map_err(|e| e.to_string())?;
        out.take().ok_or_else(|| "No stdout available".to_string())
    }

    pub fn send(&self, message: &str) -> Result<(), String> {
        let mut proc = self.process.lock().map_err(|e| e.to_string())?;
        if let Some(ref mut child) = *proc {
            if let Some(ref mut stdin) = child.stdin {
                writeln!(stdin, "{}", message).map_err(|e| format!("Write failed: {}", e))?;
                stdin.flush().map_err(|e| format!("Flush failed: {}", e))?;
            }
        }
        Ok(())
    }

    pub fn stop(&self) -> Result<(), String> {
        let mut proc = self.process.lock().map_err(|e| e.to_string())?;
        if let Some(ref mut child) = *proc {
            child.kill().map_err(|e| format!("Kill failed: {}", e))?;
        }
        *proc = None;
        Ok(())
    }
}
```

- [ ] **Step 5: 프론트엔드 IPC 래퍼 — ipc.ts**

Create `desktop/src/lib/ipc.ts`:

```typescript
import { invoke } from "@tauri-apps/api/core";
import { listen, type UnlistenFn } from "@tauri-apps/api/event";

export interface SidecarEvent {
  type: "text" | "tool-call" | "tool-result" | "done" | "error";
  data: unknown;
  timestamp: string;
}

export async function startSidecar(): Promise<string> {
  return invoke<string>("start_sidecar");
}

export async function sendPrompt(content: string): Promise<void> {
  return invoke("send_prompt", { content });
}

export async function stopSidecar(): Promise<void> {
  return invoke("stop_sidecar");
}

export async function onSidecarEvent(
  callback: (event: SidecarEvent) => void
): Promise<UnlistenFn> {
  return listen<string>("sidecar-event", (e) => {
    try {
      const parsed = JSON.parse(e.payload) as SidecarEvent;
      callback(parsed);
    } catch {
      // ignore malformed JSON
    }
  });
}
```

- [ ] **Step 6: use-sidecar hook 생성**

Create `desktop/src/hooks/use-sidecar.ts`:

```typescript
import { useEffect, useRef, useCallback } from "react";
import {
  startSidecar,
  sendPrompt,
  stopSidecar,
  onSidecarEvent,
  type SidecarEvent,
} from "../lib/ipc";

interface UseSidecarOptions {
  onEvent: (event: SidecarEvent) => void;
}

export function useSidecar({ onEvent }: UseSidecarOptions) {
  const started = useRef(false);

  useEffect(() => {
    let unlisten: (() => void) | undefined;

    async function init() {
      if (started.current) return;
      started.current = true;

      unlisten = await onSidecarEvent(onEvent);
      await startSidecar();
    }

    init();

    return () => {
      unlisten?.();
      stopSidecar();
    };
  }, [onEvent]);

  const send = useCallback(async (content: string) => {
    await sendPrompt(content);
  }, []);

  return { send };
}
```

- [ ] **Step 7: 빌드 테스트**

Run: `cd desktop/src-tauri && cargo check`
Expected: 컴파일 에러 없음

- [ ] **Step 8: Commit**

```bash
git add desktop/
git commit -m "feat(desktop): add Node.js sidecar IPC integration

Rust sidecar manager spawns hanimo headless mode as child process.
Tauri events stream stdout to React frontend via use-sidecar hook."
```

---

## Task 3: 테마 시스템 (다크/라이트)

**Files:**
- Create: `desktop/src/lib/theme.ts`
- Create: `desktop/src/stores/theme-store.ts`

- [ ] **Step 1: 테마 정의 — theme.ts**

```typescript
export interface Theme {
  name: "dark" | "light";
  colors: {
    bg: string;
    bgSecondary: string;
    bgTertiary: string;
    border: string;
    text: string;
    textSecondary: string;
    textMuted: string;
    accent: string;
    accentHover: string;
    userBubble: string;
    assistantBubble: string;
    success: string;
    error: string;
    sidebarBg: string;
    sidebarIcon: string;
    sidebarIconActive: string;
    editorBg: string;
    tabBg: string;
    tabActiveBg: string;
    inputBg: string;
    inputBorder: string;
  };
}

export const darkTheme: Theme = {
  name: "dark",
  colors: {
    bg: "#0a0a0a",
    bgSecondary: "#111111",
    bgTertiary: "#1a1a1a",
    border: "#2a2a2a",
    text: "#e0e0e0",
    textSecondary: "#888888",
    textMuted: "#555555",
    accent: "#6366f1",
    accentHover: "#818cf8",
    userBubble: "#1a1a2e",
    assistantBubble: "#0d1117",
    success: "#7ee787",
    error: "#ff7b72",
    sidebarBg: "#0d0d0d",
    sidebarIcon: "#666666",
    sidebarIconActive: "#e0e0e0",
    editorBg: "#0d1117",
    tabBg: "#111111",
    tabActiveBg: "#1a1a1a",
    inputBg: "#1a1a1a",
    inputBorder: "#333333",
  },
};

export const lightTheme: Theme = {
  name: "light",
  colors: {
    bg: "#ffffff",
    bgSecondary: "#f8f9fa",
    bgTertiary: "#f0f0f0",
    border: "#e0e0e0",
    text: "#333333",
    textSecondary: "#666666",
    textMuted: "#999999",
    accent: "#1a73e8",
    accentHover: "#1557b0",
    userBubble: "#e8f4fd",
    assistantBubble: "#f0f0f0",
    success: "#28a745",
    error: "#dc3545",
    sidebarBg: "#f0f0f0",
    sidebarIcon: "#999999",
    sidebarIconActive: "#333333",
    editorBg: "#ffffff",
    tabBg: "#f8f9fa",
    tabActiveBg: "#ffffff",
    inputBg: "#ffffff",
    inputBorder: "#d0d0d0",
  },
};

export function getTheme(name: "dark" | "light"): Theme {
  return name === "dark" ? darkTheme : lightTheme;
}
```

- [ ] **Step 2: 테마 store — theme-store.ts**

```typescript
import { create } from "zustand";
import { persist } from "zustand/middleware";
import { type Theme, getTheme } from "../lib/theme";

interface ThemeState {
  mode: "dark" | "light";
  theme: Theme;
  toggle: () => void;
  setMode: (mode: "dark" | "light") => void;
}

export const useThemeStore = create<ThemeState>()(
  persist(
    (set) => ({
      mode: "dark",
      theme: getTheme("dark"),
      toggle: () =>
        set((state) => {
          const next = state.mode === "dark" ? "light" : "dark";
          return { mode: next, theme: getTheme(next) };
        }),
      setMode: (mode) => set({ mode, theme: getTheme(mode) }),
    }),
    { name: "hanimo-theme" }
  )
);
```

- [ ] **Step 3: Commit**

```bash
git add desktop/src/lib/theme.ts desktop/src/stores/theme-store.ts
git commit -m "feat(desktop): add dark/light theme system with Zustand persistence"
```

---

## Task 4: VSCode 스타일 레이아웃 셸

**Files:**
- Create: `desktop/src/App.tsx`
- Create: `desktop/src/components/layout/Sidebar.tsx`
- Create: `desktop/src/components/layout/FileTree.tsx`
- Create: `desktop/src/components/layout/EditorArea.tsx`
- Create: `desktop/src/components/layout/ChatPanel.tsx`
- Modify: `desktop/src/main.tsx`

- [ ] **Step 1: App.tsx — 루트 레이아웃**

```tsx
import { Sidebar } from "./components/layout/Sidebar";
import { FileTree } from "./components/layout/FileTree";
import { EditorArea } from "./components/layout/EditorArea";
import { ChatPanel } from "./components/layout/ChatPanel";
import { useThemeStore } from "./stores/theme-store";

export function App() {
  const { theme } = useThemeStore();
  const c = theme.colors;

  return (
    <div
      className="h-screen flex overflow-hidden"
      style={{ background: c.bg, color: c.text }}
    >
      <Sidebar />
      <FileTree />
      <EditorArea />
      <ChatPanel />
    </div>
  );
}
```

- [ ] **Step 2: Sidebar.tsx — 아이콘 사이드바**

```tsx
import { useState } from "react";
import { useThemeStore } from "../../stores/theme-store";

type SidebarTab = "files" | "search" | "chat" | "settings";

export function Sidebar() {
  const [active, setActive] = useState<SidebarTab>("files");
  const { theme, toggle } = useThemeStore();
  const c = theme.colors;

  const items: { id: SidebarTab; icon: string; label: string }[] = [
    { id: "files", icon: "📁", label: "Explorer" },
    { id: "search", icon: "🔍", label: "Search" },
    { id: "chat", icon: "💬", label: "Chat" },
    { id: "settings", icon: "⚙️", label: "Settings" },
  ];

  return (
    <div
      className="flex flex-col items-center py-2 gap-1"
      style={{ width: 48, background: c.sidebarBg, borderRight: `1px solid ${c.border}` }}
    >
      {items.map((item) => (
        <button
          key={item.id}
          onClick={() => setActive(item.id)}
          title={item.label}
          className="w-10 h-10 flex items-center justify-center rounded-md text-lg cursor-pointer transition-colors"
          style={{
            background: active === item.id ? c.bgTertiary : "transparent",
            color: active === item.id ? c.sidebarIconActive : c.sidebarIcon,
          }}
        >
          {item.icon}
        </button>
      ))}
      <div className="flex-1" />
      <button
        onClick={toggle}
        title="Toggle theme"
        className="w-10 h-10 flex items-center justify-center rounded-md text-lg cursor-pointer"
        style={{ color: c.sidebarIcon }}
      >
        {theme.name === "dark" ? "🌙" : "☀️"}
      </button>
    </div>
  );
}
```

- [ ] **Step 3: FileTree.tsx — 파일 탐색기 (플레이스홀더)**

```tsx
import { useThemeStore } from "../../stores/theme-store";

export function FileTree() {
  const { theme } = useThemeStore();
  const c = theme.colors;

  return (
    <div
      className="flex flex-col overflow-y-auto"
      style={{
        width: 220,
        background: c.bgSecondary,
        borderRight: `1px solid ${c.border}`,
      }}
    >
      <div
        className="px-3 py-2 text-xs font-semibold uppercase tracking-wider"
        style={{ color: c.textMuted }}
      >
        Explorer
      </div>
      <div className="px-3 py-1 text-sm" style={{ color: c.textSecondary }}>
        Open a folder to get started
      </div>
    </div>
  );
}
```

- [ ] **Step 4: EditorArea.tsx — 에디터 영역 (플레이스홀더)**

```tsx
import { useThemeStore } from "../../stores/theme-store";

export function EditorArea() {
  const { theme } = useThemeStore();
  const c = theme.colors;

  return (
    <div className="flex-1 flex flex-col" style={{ background: c.editorBg }}>
      <div
        className="flex items-center px-2 gap-1"
        style={{
          height: 36,
          background: c.bgSecondary,
          borderBottom: `1px solid ${c.border}`,
        }}
      >
        <span
          className="text-xs px-3 py-1 rounded-t"
          style={{ background: c.tabActiveBg, color: c.text }}
        >
          Welcome
        </span>
      </div>
      <div className="flex-1 flex items-center justify-center">
        <div className="text-center" style={{ color: c.textMuted }}>
          <div className="text-4xl mb-4">🐶</div>
          <div className="text-lg font-medium" style={{ color: c.text }}>
            hanimo
          </div>
          <div className="text-sm mt-1">Start a conversation to begin</div>
        </div>
      </div>
    </div>
  );
}
```

- [ ] **Step 5: ChatPanel.tsx — 채팅 패널 (플레이스홀더)**

```tsx
import { useThemeStore } from "../../stores/theme-store";

export function ChatPanel() {
  const { theme } = useThemeStore();
  const c = theme.colors;

  return (
    <div
      className="flex flex-col"
      style={{
        width: 340,
        background: c.bgSecondary,
        borderLeft: `1px solid ${c.border}`,
      }}
    >
      <div
        className="flex items-center px-3 gap-2"
        style={{
          height: 36,
          borderBottom: `1px solid ${c.border}`,
        }}
      >
        <span className="text-sm">🐶</span>
        <span className="text-sm font-medium">hanimo chat</span>
      </div>
      <div className="flex-1 overflow-y-auto p-3">
        <div className="text-sm" style={{ color: c.textMuted }}>
          Send a message to start chatting with hanimo
        </div>
      </div>
      <div className="p-3" style={{ borderTop: `1px solid ${c.border}` }}>
        <div
          className="flex items-center rounded-lg px-3 py-2"
          style={{ background: c.inputBg, border: `1px solid ${c.inputBorder}` }}
        >
          <input
            type="text"
            placeholder="Message hanimo..."
            className="flex-1 bg-transparent outline-none text-sm"
            style={{ color: c.text }}
          />
          <button
            className="ml-2 px-2 py-1 rounded text-xs"
            style={{ background: c.accent, color: "#fff" }}
          >
            Send
          </button>
        </div>
      </div>
    </div>
  );
}
```

- [ ] **Step 6: main.tsx 업데이트**

```tsx
import React from "react";
import ReactDOM from "react-dom/client";
import { App } from "./App";
import "./styles/globals.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
```

- [ ] **Step 7: 실행 테스트**

Run: `cd desktop && npm run tauri dev`
Expected: VSCode 스타일 레이아웃이 렌더링됨 — 좌측 사이드바, 파일트리, 에디터 영역, 우측 채팅 패널. 테마 토글 버튼으로 다크↔라이트 전환.

- [ ] **Step 8: Commit**

```bash
git add desktop/src/
git commit -m "feat(desktop): add VSCode-style layout with sidebar, file tree, editor, chat panel

Four-panel layout: icon sidebar (48px), file tree (220px), editor area (flex),
chat panel (340px). Dark/light theme toggle in sidebar."
```

---

## Task 5: 채팅 기능 구현

**Files:**
- Create: `desktop/src/stores/chat-store.ts`
- Create: `desktop/src/components/chat/MessageList.tsx`
- Create: `desktop/src/components/chat/MessageBubble.tsx`
- Create: `desktop/src/components/chat/ChatInput.tsx`
- Modify: `desktop/src/components/layout/ChatPanel.tsx`

- [ ] **Step 1: chat-store.ts — 채팅 상태 관리**

```typescript
import { create } from "zustand";

export interface ChatMessage {
  id: string;
  role: "user" | "assistant" | "tool-call" | "tool-result";
  content: string;
  toolName?: string;
  isError?: boolean;
  timestamp: number;
}

interface ChatState {
  messages: ChatMessage[];
  isStreaming: boolean;
  streamingContent: string;
  addMessage: (msg: Omit<ChatMessage, "id" | "timestamp">) => void;
  setStreaming: (streaming: boolean) => void;
  appendStreamingContent: (content: string) => void;
  clearStreamingContent: () => void;
  clear: () => void;
}

let nextId = 0;

export const useChatStore = create<ChatState>((set) => ({
  messages: [],
  isStreaming: false,
  streamingContent: "",
  addMessage: (msg) =>
    set((state) => ({
      messages: [
        ...state.messages,
        { ...msg, id: String(++nextId), timestamp: Date.now() },
      ],
    })),
  setStreaming: (isStreaming) => set({ isStreaming }),
  appendStreamingContent: (content) =>
    set((state) => ({
      streamingContent: state.streamingContent + content,
    })),
  clearStreamingContent: () => set({ streamingContent: "" }),
  clear: () => set({ messages: [], streamingContent: "" }),
}));
```

- [ ] **Step 2: MessageBubble.tsx — 단일 메시지 렌더링**

```tsx
import { useThemeStore } from "../../stores/theme-store";
import type { ChatMessage } from "../../stores/chat-store";

interface Props {
  message: ChatMessage;
}

export function MessageBubble({ message }: Props) {
  const { theme } = useThemeStore();
  const c = theme.colors;

  const isUser = message.role === "user";
  const isTool = message.role === "tool-call" || message.role === "tool-result";

  if (isTool) {
    return (
      <div
        className="mx-3 my-1 px-3 py-2 rounded-md text-xs font-mono"
        style={{
          background: c.bgTertiary,
          color: message.isError ? c.error : c.success,
          border: `1px solid ${c.border}`,
        }}
      >
        <span style={{ color: c.textMuted }}>
          {message.role === "tool-call" ? "▶" : "✓"} {message.toolName}
        </span>
        <div className="mt-1 whitespace-pre-wrap" style={{ color: c.textSecondary }}>
          {message.content}
        </div>
      </div>
    );
  }

  return (
    <div className={`px-3 py-1 ${isUser ? "flex justify-end" : ""}`}>
      <div
        className="max-w-[85%] rounded-xl px-3 py-2 text-sm whitespace-pre-wrap"
        style={{
          background: isUser ? c.userBubble : c.assistantBubble,
          color: isUser ? c.accent : c.text,
          border: `1px solid ${c.border}`,
        }}
      >
        {message.content}
      </div>
    </div>
  );
}
```

- [ ] **Step 3: MessageList.tsx — 메시지 목록**

```tsx
import { useEffect, useRef } from "react";
import { useChatStore } from "../../stores/chat-store";
import { useThemeStore } from "../../stores/theme-store";
import { MessageBubble } from "./MessageBubble";

export function MessageList() {
  const { messages, isStreaming, streamingContent } = useChatStore();
  const { theme } = useThemeStore();
  const c = theme.colors;
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages, streamingContent]);

  return (
    <div className="flex-1 overflow-y-auto py-2">
      {messages.map((msg) => (
        <MessageBubble key={msg.id} message={msg} />
      ))}
      {isStreaming && streamingContent && (
        <div className="px-3 py-1">
          <div
            className="max-w-[85%] rounded-xl px-3 py-2 text-sm whitespace-pre-wrap"
            style={{
              background: c.assistantBubble,
              color: c.text,
              border: `1px solid ${c.border}`,
            }}
          >
            {streamingContent}
            <span className="animate-pulse">▊</span>
          </div>
        </div>
      )}
      <div ref={bottomRef} />
    </div>
  );
}
```

- [ ] **Step 4: ChatInput.tsx — 메시지 입력**

```tsx
import { useState, useCallback, type KeyboardEvent } from "react";
import { useThemeStore } from "../../stores/theme-store";
import { useChatStore } from "../../stores/chat-store";

interface Props {
  onSend: (content: string) => void;
}

export function ChatInput({ onSend }: Props) {
  const [input, setInput] = useState("");
  const { theme } = useThemeStore();
  const { isStreaming } = useChatStore();
  const c = theme.colors;

  const handleSend = useCallback(() => {
    const trimmed = input.trim();
    if (!trimmed || isStreaming) return;
    onSend(trimmed);
    setInput("");
  }, [input, isStreaming, onSend]);

  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === "Enter" && !e.shiftKey) {
        e.preventDefault();
        handleSend();
      }
    },
    [handleSend]
  );

  return (
    <div className="p-3" style={{ borderTop: `1px solid ${c.border}` }}>
      <div
        className="flex items-end rounded-lg px-3 py-2"
        style={{ background: c.inputBg, border: `1px solid ${c.inputBorder}` }}
      >
        <textarea
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Message hanimo..."
          rows={1}
          className="flex-1 bg-transparent outline-none text-sm resize-none"
          style={{ color: c.text, maxHeight: 120 }}
          disabled={isStreaming}
        />
        <button
          onClick={handleSend}
          disabled={isStreaming || !input.trim()}
          className="ml-2 px-3 py-1 rounded text-xs font-medium transition-colors"
          style={{
            background: input.trim() && !isStreaming ? c.accent : c.bgTertiary,
            color: input.trim() && !isStreaming ? "#fff" : c.textMuted,
            cursor: input.trim() && !isStreaming ? "pointer" : "default",
          }}
        >
          {isStreaming ? "..." : "Send"}
        </button>
      </div>
    </div>
  );
}
```

- [ ] **Step 5: ChatPanel.tsx 업데이트 — 실제 컴포넌트 조합**

```tsx
import { useCallback } from "react";
import { useThemeStore } from "../../stores/theme-store";
import { useChatStore } from "../../stores/chat-store";
import { useSidecar } from "../../hooks/use-sidecar";
import { MessageList } from "../chat/MessageList";
import { ChatInput } from "../chat/ChatInput";
import type { SidecarEvent } from "../../lib/ipc";

export function ChatPanel() {
  const { theme } = useThemeStore();
  const c = theme.colors;
  const {
    addMessage,
    setStreaming,
    appendStreamingContent,
    clearStreamingContent,
  } = useChatStore();

  const handleSidecarEvent = useCallback(
    (event: SidecarEvent) => {
      switch (event.type) {
        case "text":
          appendStreamingContent(event.data as string);
          break;
        case "tool-call": {
          const tc = event.data as { toolName: string; args: unknown };
          addMessage({
            role: "tool-call",
            content: JSON.stringify(tc.args, null, 2),
            toolName: tc.toolName,
          });
          break;
        }
        case "tool-result": {
          const tr = event.data as {
            toolName: string;
            result: string;
            isError: boolean;
          };
          addMessage({
            role: "tool-result",
            content: tr.result,
            toolName: tr.toolName,
            isError: tr.isError,
          });
          break;
        }
        case "done": {
          const done = event.data as { response: string };
          clearStreamingContent();
          addMessage({ role: "assistant", content: done.response });
          setStreaming(false);
          break;
        }
        case "error": {
          const err = event.data as { message: string };
          clearStreamingContent();
          addMessage({
            role: "assistant",
            content: `Error: ${err.message}`,
            isError: true,
          });
          setStreaming(false);
          break;
        }
      }
    },
    [addMessage, setStreaming, appendStreamingContent, clearStreamingContent]
  );

  const { send } = useSidecar({ onEvent: handleSidecarEvent });

  const handleSend = useCallback(
    (content: string) => {
      addMessage({ role: "user", content });
      setStreaming(true);
      clearStreamingContent();
      send(content);
    },
    [addMessage, setStreaming, clearStreamingContent, send]
  );

  return (
    <div
      className="flex flex-col"
      style={{
        width: 340,
        background: c.bgSecondary,
        borderLeft: `1px solid ${c.border}`,
      }}
    >
      <div
        className="flex items-center px-3 gap-2"
        style={{ height: 36, borderBottom: `1px solid ${c.border}` }}
      >
        <span className="text-sm">🐶</span>
        <span className="text-sm font-medium">hanimo chat</span>
      </div>
      <MessageList />
      <ChatInput onSend={handleSend} />
    </div>
  );
}
```

- [ ] **Step 6: 실행 테스트**

Run: `cd desktop && npm run tauri dev`
Expected: 오른쪽 채팅 패널에서 메시지 입력 → sidecar로 전송 → 스트리밍 응답 표시

- [ ] **Step 7: Commit**

```bash
git add desktop/src/
git commit -m "feat(desktop): implement chat panel with sidecar streaming

Chat store (Zustand), message bubbles, streaming indicator, tool call/result
display. Connected to Node.js sidecar via Tauri events."
```

---

## Task 6: 파일 트리 & 에디터 통합

**Files:**
- Create: `desktop/src/stores/editor-store.ts`
- Create: `desktop/src/hooks/use-file-system.ts`
- Create: `desktop/src/components/editor/TabBar.tsx`
- Create: `desktop/src/components/editor/MonacoWrapper.tsx`
- Modify: `desktop/src/components/layout/FileTree.tsx`
- Modify: `desktop/src/components/layout/EditorArea.tsx`

- [ ] **Step 1: editor-store.ts — 에디터 상태**

```typescript
import { create } from "zustand";

export interface OpenFile {
  path: string;
  name: string;
  content: string;
  language: string;
  isDirty: boolean;
}

interface EditorState {
  openFiles: OpenFile[];
  activeFilePath: string | null;
  projectRoot: string | null;
  openFile: (file: OpenFile) => void;
  closeFile: (path: string) => void;
  setActiveFile: (path: string) => void;
  updateFileContent: (path: string, content: string) => void;
  markSaved: (path: string) => void;
  setProjectRoot: (root: string) => void;
}

export const useEditorStore = create<EditorState>((set) => ({
  openFiles: [],
  activeFilePath: null,
  projectRoot: null,
  openFile: (file) =>
    set((state) => {
      const exists = state.openFiles.find((f) => f.path === file.path);
      if (exists) return { activeFilePath: file.path };
      return {
        openFiles: [...state.openFiles, file],
        activeFilePath: file.path,
      };
    }),
  closeFile: (path) =>
    set((state) => {
      const remaining = state.openFiles.filter((f) => f.path !== path);
      const newActive =
        state.activeFilePath === path
          ? remaining[remaining.length - 1]?.path ?? null
          : state.activeFilePath;
      return { openFiles: remaining, activeFilePath: newActive };
    }),
  setActiveFile: (path) => set({ activeFilePath: path }),
  updateFileContent: (path, content) =>
    set((state) => ({
      openFiles: state.openFiles.map((f) =>
        f.path === path ? { ...f, content, isDirty: true } : f
      ),
    })),
  markSaved: (path) =>
    set((state) => ({
      openFiles: state.openFiles.map((f) =>
        f.path === path ? { ...f, isDirty: false } : f
      ),
    })),
  setProjectRoot: (root) => set({ projectRoot: root }),
}));

const EXT_LANG_MAP: Record<string, string> = {
  ts: "typescript",
  tsx: "typescriptreact",
  js: "javascript",
  jsx: "javascriptreact",
  json: "json",
  md: "markdown",
  css: "css",
  html: "html",
  py: "python",
  rs: "rust",
  toml: "toml",
  yaml: "yaml",
  yml: "yaml",
  sh: "shell",
  bash: "shell",
};

export function getLanguageFromPath(path: string): string {
  const ext = path.split(".").pop()?.toLowerCase() ?? "";
  return EXT_LANG_MAP[ext] ?? "plaintext";
}
```

- [ ] **Step 2: use-file-system.ts — 파일 시스템 훅**

```typescript
import { useCallback } from "react";
import { readTextFile, writeTextFile, readDir } from "@tauri-apps/plugin-fs";
import { open } from "@tauri-apps/plugin-dialog";
import { useEditorStore, getLanguageFromPath } from "../stores/editor-store";

export interface FileEntry {
  name: string;
  path: string;
  isDirectory: boolean;
  children?: FileEntry[];
}

export function useFileSystem() {
  const { openFile, markSaved, setProjectRoot } = useEditorStore();

  const openFolder = useCallback(async () => {
    const selected = await open({ directory: true, multiple: false });
    if (typeof selected === "string") {
      setProjectRoot(selected);
      return selected;
    }
    return null;
  }, [setProjectRoot]);

  const listDirectory = useCallback(
    async (dirPath: string): Promise<FileEntry[]> => {
      const entries = await readDir(dirPath);
      const result: FileEntry[] = [];
      for (const entry of entries) {
        if (entry.name.startsWith(".")) continue;
        if (entry.name === "node_modules") continue;
        result.push({
          name: entry.name,
          path: `${dirPath}/${entry.name}`,
          isDirectory: entry.isDirectory,
        });
      }
      return result.sort((a, b) => {
        if (a.isDirectory !== b.isDirectory) return a.isDirectory ? -1 : 1;
        return a.name.localeCompare(b.name);
      });
    },
    []
  );

  const openFileFromPath = useCallback(
    async (filePath: string) => {
      const content = await readTextFile(filePath);
      const name = filePath.split("/").pop() ?? filePath;
      openFile({
        path: filePath,
        name,
        content,
        language: getLanguageFromPath(filePath),
        isDirty: false,
      });
    },
    [openFile]
  );

  const saveFile = useCallback(
    async (filePath: string, content: string) => {
      await writeTextFile(filePath, content);
      markSaved(filePath);
    },
    [markSaved]
  );

  return { openFolder, listDirectory, openFileFromPath, saveFile };
}
```

Add `@tauri-apps/plugin-dialog` dependency:

Run: `cd desktop && npm install @tauri-apps/plugin-dialog`

Add plugin to `desktop/src-tauri/Cargo.toml` dependencies:

```toml
tauri-plugin-dialog = "2"
```

Add to `desktop/src-tauri/src/lib.rs`:

```rust
.plugin(tauri_plugin_dialog::init())
```

Add to `desktop/src-tauri/capabilities/default.json` permissions:

```json
"dialog:allow-open"
```

- [ ] **Step 3: FileTree.tsx — 실제 파일 탐색기**

```tsx
import { useState, useEffect, useCallback } from "react";
import { useThemeStore } from "../../stores/theme-store";
import { useEditorStore } from "../../stores/editor-store";
import { useFileSystem, type FileEntry } from "../../hooks/use-file-system";

function FileItem({
  entry,
  depth,
  onFileClick,
  onToggleDir,
  expandedDirs,
}: {
  entry: FileEntry;
  depth: number;
  onFileClick: (path: string) => void;
  onToggleDir: (path: string) => void;
  expandedDirs: Set<string>;
}) {
  const { theme } = useThemeStore();
  const c = theme.colors;
  const isExpanded = expandedDirs.has(entry.path);

  return (
    <>
      <button
        className="w-full text-left px-2 py-0.5 text-sm flex items-center gap-1 hover:opacity-80"
        style={{
          paddingLeft: `${depth * 16 + 8}px`,
          color: c.text,
          background: "transparent",
          border: "none",
          cursor: "pointer",
        }}
        onClick={() =>
          entry.isDirectory ? onToggleDir(entry.path) : onFileClick(entry.path)
        }
      >
        <span className="text-xs">
          {entry.isDirectory ? (isExpanded ? "📂" : "📁") : "📄"}
        </span>
        <span className="truncate">{entry.name}</span>
      </button>
      {entry.isDirectory && isExpanded && entry.children?.map((child) => (
        <FileItem
          key={child.path}
          entry={child}
          depth={depth + 1}
          onFileClick={onFileClick}
          onToggleDir={onToggleDir}
          expandedDirs={expandedDirs}
        />
      ))}
    </>
  );
}

export function FileTree() {
  const { theme } = useThemeStore();
  const c = theme.colors;
  const { projectRoot } = useEditorStore();
  const { openFolder, listDirectory, openFileFromPath } = useFileSystem();
  const [entries, setEntries] = useState<FileEntry[]>([]);
  const [expandedDirs, setExpandedDirs] = useState<Set<string>>(new Set());

  useEffect(() => {
    if (projectRoot) {
      listDirectory(projectRoot).then(setEntries);
    }
  }, [projectRoot, listDirectory]);

  const handleToggleDir = useCallback(
    async (path: string) => {
      const next = new Set(expandedDirs);
      if (next.has(path)) {
        next.delete(path);
      } else {
        next.add(path);
        const children = await listDirectory(path);
        setEntries((prev) => {
          const update = (items: FileEntry[]): FileEntry[] =>
            items.map((item) =>
              item.path === path
                ? { ...item, children }
                : item.children
                  ? { ...item, children: update(item.children) }
                  : item
            );
          return update(prev);
        });
      }
      setExpandedDirs(next);
    },
    [expandedDirs, listDirectory]
  );

  if (!projectRoot) {
    return (
      <div
        className="flex flex-col items-center justify-center gap-2"
        style={{
          width: 220,
          background: c.bgSecondary,
          borderRight: `1px solid ${c.border}`,
        }}
      >
        <button
          onClick={openFolder}
          className="px-3 py-1.5 rounded text-sm cursor-pointer"
          style={{ background: c.accent, color: "#fff" }}
        >
          Open Folder
        </button>
      </div>
    );
  }

  return (
    <div
      className="flex flex-col overflow-y-auto"
      style={{
        width: 220,
        background: c.bgSecondary,
        borderRight: `1px solid ${c.border}`,
      }}
    >
      <div
        className="px-3 py-2 text-xs font-semibold uppercase tracking-wider"
        style={{ color: c.textMuted }}
      >
        {projectRoot.split("/").pop()}
      </div>
      {entries.map((entry) => (
        <FileItem
          key={entry.path}
          entry={entry}
          depth={0}
          onFileClick={openFileFromPath}
          onToggleDir={handleToggleDir}
          expandedDirs={expandedDirs}
        />
      ))}
    </div>
  );
}
```

- [ ] **Step 4: TabBar.tsx — 에디터 탭 바**

```tsx
import { useThemeStore } from "../../stores/theme-store";
import { useEditorStore } from "../../stores/editor-store";

export function TabBar() {
  const { theme } = useThemeStore();
  const c = theme.colors;
  const { openFiles, activeFilePath, setActiveFile, closeFile } =
    useEditorStore();

  if (openFiles.length === 0) return null;

  return (
    <div
      className="flex items-center overflow-x-auto"
      style={{
        height: 36,
        background: c.bgSecondary,
        borderBottom: `1px solid ${c.border}`,
      }}
    >
      {openFiles.map((file) => (
        <button
          key={file.path}
          onClick={() => setActiveFile(file.path)}
          className="flex items-center gap-1 px-3 text-xs whitespace-nowrap h-full border-r"
          style={{
            background:
              file.path === activeFilePath ? c.tabActiveBg : c.tabBg,
            color: file.path === activeFilePath ? c.text : c.textSecondary,
            borderColor: c.border,
            cursor: "pointer",
            borderBottom:
              file.path === activeFilePath
                ? `2px solid ${c.accent}`
                : "2px solid transparent",
          }}
        >
          {file.isDirty && (
            <span style={{ color: c.accent, fontSize: 10 }}>●</span>
          )}
          {file.name}
          <span
            onClick={(e) => {
              e.stopPropagation();
              closeFile(file.path);
            }}
            className="ml-1 hover:opacity-100 opacity-50"
            style={{ fontSize: 10 }}
          >
            ✕
          </span>
        </button>
      ))}
    </div>
  );
}
```

- [ ] **Step 5: MonacoWrapper.tsx — Monaco 에디터**

```tsx
import Editor from "@monaco-editor/react";
import { useThemeStore } from "../../stores/theme-store";
import { useEditorStore } from "../../stores/editor-store";
import { useFileSystem } from "../../hooks/use-file-system";
import { useCallback, useEffect, useRef } from "react";

export function MonacoWrapper() {
  const { theme } = useThemeStore();
  const { openFiles, activeFilePath, updateFileContent } = useEditorStore();
  const { saveFile } = useFileSystem();
  const editorRef = useRef<unknown>(null);

  const activeFile = openFiles.find((f) => f.path === activeFilePath);

  const handleChange = useCallback(
    (value: string | undefined) => {
      if (activeFilePath && value !== undefined) {
        updateFileContent(activeFilePath, value);
      }
    },
    [activeFilePath, updateFileContent]
  );

  // Ctrl/Cmd+S to save
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === "s") {
        e.preventDefault();
        if (activeFile && activeFile.isDirty) {
          saveFile(activeFile.path, activeFile.content);
        }
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [activeFile, saveFile]);

  if (!activeFile) return null;

  return (
    <Editor
      height="100%"
      language={activeFile.language}
      value={activeFile.content}
      theme={theme.name === "dark" ? "vs-dark" : "light"}
      onChange={handleChange}
      onMount={(editor) => {
        editorRef.current = editor;
      }}
      options={{
        minimap: { enabled: false },
        fontSize: 13,
        lineNumbers: "on",
        scrollBeyondLastLine: false,
        wordWrap: "on",
        automaticLayout: true,
        padding: { top: 8 },
      }}
    />
  );
}
```

- [ ] **Step 6: EditorArea.tsx 업데이트**

```tsx
import { useThemeStore } from "../../stores/theme-store";
import { useEditorStore } from "../../stores/editor-store";
import { TabBar } from "../editor/TabBar";
import { MonacoWrapper } from "../editor/MonacoWrapper";

export function EditorArea() {
  const { theme } = useThemeStore();
  const c = theme.colors;
  const { openFiles } = useEditorStore();

  return (
    <div className="flex-1 flex flex-col" style={{ background: c.editorBg }}>
      <TabBar />
      {openFiles.length > 0 ? (
        <div className="flex-1">
          <MonacoWrapper />
        </div>
      ) : (
        <div className="flex-1 flex items-center justify-center">
          <div className="text-center" style={{ color: c.textMuted }}>
            <div className="text-4xl mb-4">🐶</div>
            <div className="text-lg font-medium" style={{ color: c.text }}>
              hanimo
            </div>
            <div className="text-sm mt-1">
              Open a file or start a conversation
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
```

- [ ] **Step 7: 실행 테스트**

Run: `cd desktop && npm run tauri dev`
Expected: Open Folder 버튼 → 폴더 선택 → 파일 트리 표시 → 파일 클릭 → Monaco 에디터에서 열림 → 탭 전환/닫기 작동 → Cmd+S 저장

- [ ] **Step 8: Commit**

```bash
git add desktop/
git commit -m "feat(desktop): add file tree explorer and Monaco editor with tabs

File tree with lazy directory expansion, Monaco editor with syntax
highlighting, tab management, Cmd+S save, dirty file indicator."
```

---

## Task 7: 온보딩 위저드

**Files:**
- Create: `desktop/src/stores/onboarding-store.ts`
- Create: `desktop/src/components/onboarding/OnboardingWizard.tsx`
- Modify: `desktop/src/App.tsx`

- [ ] **Step 1: onboarding-store.ts**

```typescript
import { create } from "zustand";
import { persist } from "zustand/middleware";

interface OnboardingState {
  completed: boolean;
  provider: string;
  apiKey: string;
  model: string;
  step: number;
  setProvider: (provider: string) => void;
  setApiKey: (key: string) => void;
  setModel: (model: string) => void;
  nextStep: () => void;
  prevStep: () => void;
  complete: () => void;
  reset: () => void;
}

export const useOnboardingStore = create<OnboardingState>()(
  persist(
    (set) => ({
      completed: false,
      provider: "",
      apiKey: "",
      model: "",
      step: 0,
      setProvider: (provider) => set({ provider }),
      setApiKey: (apiKey) => set({ apiKey }),
      setModel: (model) => set({ model }),
      nextStep: () => set((s) => ({ step: s.step + 1 })),
      prevStep: () => set((s) => ({ step: Math.max(0, s.step - 1) })),
      complete: () => set({ completed: true }),
      reset: () =>
        set({ completed: false, provider: "", apiKey: "", model: "", step: 0 }),
    }),
    { name: "hanimo-onboarding" }
  )
);

export const PROVIDERS = [
  { id: "openai", name: "OpenAI", models: ["gpt-4o", "gpt-4o-mini", "gpt-4.1", "gpt-4.1-mini"] },
  { id: "anthropic", name: "Anthropic", models: ["claude-sonnet-4-20250514", "claude-haiku-4-20250414", "claude-opus-4-20250514"] },
  { id: "google", name: "Google", models: ["gemini-2.5-flash", "gemini-2.5-pro"] },
  { id: "deepseek", name: "DeepSeek", models: ["deepseek-chat", "deepseek-coder", "deepseek-reasoner"] },
  { id: "groq", name: "Groq", models: ["qwen-qwq-32b", "llama-3.3-70b-versatile"] },
  { id: "ollama", name: "Ollama (Local)", models: [] },
];
```

- [ ] **Step 2: OnboardingWizard.tsx**

```tsx
import { useThemeStore } from "../../stores/theme-store";
import { useOnboardingStore, PROVIDERS } from "../../stores/onboarding-store";

export function OnboardingWizard() {
  const { theme } = useThemeStore();
  const c = theme.colors;
  const {
    step,
    provider,
    apiKey,
    model,
    setProvider,
    setApiKey,
    setModel,
    nextStep,
    prevStep,
    complete,
  } = useOnboardingStore();

  const selectedProvider = PROVIDERS.find((p) => p.id === provider);

  return (
    <div
      className="h-screen flex items-center justify-center"
      style={{ background: c.bg }}
    >
      <div
        className="w-full max-w-md rounded-2xl p-8"
        style={{
          background: c.bgSecondary,
          border: `1px solid ${c.border}`,
        }}
      >
        <div className="text-center mb-8">
          <div className="text-4xl mb-2">🐶</div>
          <h1 className="text-xl font-bold" style={{ color: c.text }}>
            Welcome to hanimo
          </h1>
          <p className="text-sm mt-1" style={{ color: c.textSecondary }}>
            {step === 0 && "Choose your LLM provider"}
            {step === 1 && "Enter your API key"}
            {step === 2 && "Select a model"}
          </p>
        </div>

        {/* Step 0: Provider */}
        {step === 0 && (
          <div className="flex flex-col gap-2">
            {PROVIDERS.map((p) => (
              <button
                key={p.id}
                onClick={() => {
                  setProvider(p.id);
                  if (p.id === "ollama") {
                    setApiKey("ollama");
                    nextStep();
                  }
                  nextStep();
                }}
                className="text-left px-4 py-3 rounded-lg text-sm transition-colors"
                style={{
                  background: provider === p.id ? c.accent : c.bgTertiary,
                  color: provider === p.id ? "#fff" : c.text,
                  border: `1px solid ${c.border}`,
                  cursor: "pointer",
                }}
              >
                {p.name}
              </button>
            ))}
          </div>
        )}

        {/* Step 1: API Key */}
        {step === 1 && (
          <div className="flex flex-col gap-4">
            <input
              type="password"
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
              placeholder={`Enter ${selectedProvider?.name ?? ""} API key...`}
              className="w-full px-4 py-3 rounded-lg text-sm outline-none"
              style={{
                background: c.inputBg,
                color: c.text,
                border: `1px solid ${c.inputBorder}`,
              }}
            />
            <div className="flex gap-2">
              <button
                onClick={prevStep}
                className="flex-1 py-2 rounded-lg text-sm"
                style={{
                  background: c.bgTertiary,
                  color: c.textSecondary,
                  cursor: "pointer",
                }}
              >
                Back
              </button>
              <button
                onClick={nextStep}
                disabled={!apiKey.trim()}
                className="flex-1 py-2 rounded-lg text-sm"
                style={{
                  background: apiKey.trim() ? c.accent : c.bgTertiary,
                  color: apiKey.trim() ? "#fff" : c.textMuted,
                  cursor: apiKey.trim() ? "pointer" : "default",
                }}
              >
                Next
              </button>
            </div>
          </div>
        )}

        {/* Step 2: Model */}
        {step === 2 && (
          <div className="flex flex-col gap-2">
            {(selectedProvider?.models ?? []).map((m) => (
              <button
                key={m}
                onClick={() => setModel(m)}
                className="text-left px-4 py-3 rounded-lg text-sm font-mono"
                style={{
                  background: model === m ? c.accent : c.bgTertiary,
                  color: model === m ? "#fff" : c.text,
                  border: `1px solid ${c.border}`,
                  cursor: "pointer",
                }}
              >
                {m}
              </button>
            ))}
            <div className="flex gap-2 mt-4">
              <button
                onClick={prevStep}
                className="flex-1 py-2 rounded-lg text-sm"
                style={{
                  background: c.bgTertiary,
                  color: c.textSecondary,
                  cursor: "pointer",
                }}
              >
                Back
              </button>
              <button
                onClick={complete}
                disabled={!model}
                className="flex-1 py-2 rounded-lg text-sm font-medium"
                style={{
                  background: model ? c.accent : c.bgTertiary,
                  color: model ? "#fff" : c.textMuted,
                  cursor: model ? "pointer" : "default",
                }}
              >
                Start hanimo 🐶
              </button>
            </div>
          </div>
        )}

        {/* Progress dots */}
        <div className="flex justify-center gap-2 mt-6">
          {[0, 1, 2].map((i) => (
            <div
              key={i}
              className="w-2 h-2 rounded-full"
              style={{ background: step >= i ? c.accent : c.bgTertiary }}
            />
          ))}
        </div>
      </div>
    </div>
  );
}
```

- [ ] **Step 3: App.tsx — 온보딩 게이트 추가**

```tsx
import { Sidebar } from "./components/layout/Sidebar";
import { FileTree } from "./components/layout/FileTree";
import { EditorArea } from "./components/layout/EditorArea";
import { ChatPanel } from "./components/layout/ChatPanel";
import { OnboardingWizard } from "./components/onboarding/OnboardingWizard";
import { useThemeStore } from "./stores/theme-store";
import { useOnboardingStore } from "./stores/onboarding-store";

export function App() {
  const { theme } = useThemeStore();
  const { completed } = useOnboardingStore();
  const c = theme.colors;

  if (!completed) {
    return <OnboardingWizard />;
  }

  return (
    <div
      className="h-screen flex overflow-hidden"
      style={{ background: c.bg, color: c.text }}
    >
      <Sidebar />
      <FileTree />
      <EditorArea />
      <ChatPanel />
    </div>
  );
}
```

- [ ] **Step 4: 실행 테스트**

Run: `cd desktop && npm run tauri dev`
Expected: 첫 실행 시 온보딩 위저드 표시 → 프로바이더 선택 → API 키 입력 → 모델 선택 → 메인 앱으로 전환. 재실행 시 온보딩 스킵.

- [ ] **Step 5: Commit**

```bash
git add desktop/src/
git commit -m "feat(desktop): add onboarding wizard with provider/key/model setup

Three-step wizard: provider selection, API key input, model selection.
State persisted via Zustand + localStorage. Gates main app until complete."
```

---

## Task 8: 빌드 파이프라인 & 아이콘

**Files:**
- Create: `desktop/src-tauri/icons/` (아이콘 파일들)
- Create: `.github/workflows/desktop-build.yml`
- Modify: `desktop/src-tauri/tauri.conf.json` (최종 빌드 설정)

- [ ] **Step 1: 기본 아이콘 생성**

Run: `cd desktop && npx tauri icon --input ../assets/hanimo-icon.png`

If no icon source exists, create a placeholder:

Run: `cd desktop && npx tauri icon`
Expected: `src-tauri/icons/` 디렉토리에 각 플랫폼용 아이콘 생성

- [ ] **Step 2: GitHub Actions 워크플로우 생성**

Create `.github/workflows/desktop-build.yml`:

```yaml
name: Build Desktop App

on:
  push:
    tags: ["v*"]
  workflow_dispatch:

jobs:
  build:
    strategy:
      fail-fast: false
      matrix:
        include:
          - platform: macos-latest
            args: "--target aarch64-apple-darwin"
          - platform: macos-latest
            args: "--target x86_64-apple-darwin"
          - platform: ubuntu-22.04
            args: ""
          - platform: windows-latest
            args: ""

    runs-on: ${{ matrix.platform }}

    steps:
      - uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: "20"

      - name: Install Rust stable
        uses: dtolnay/rust-toolchain@stable
        with:
          targets: ${{ matrix.platform == 'macos-latest' && 'aarch64-apple-darwin,x86_64-apple-darwin' || '' }}

      - name: Install Linux dependencies
        if: matrix.platform == 'ubuntu-22.04'
        run: |
          sudo apt-get update
          sudo apt-get install -y libwebkit2gtk-4.1-dev libappindicator3-dev librsvg2-dev patchelf

      - name: Install frontend dependencies
        working-directory: desktop
        run: npm ci

      - name: Build hanimo core
        run: npm run build

      - name: Build desktop app
        uses: tauri-apps/tauri-action@v0
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          projectPath: desktop
          tauriScript: npx tauri
          args: ${{ matrix.args }}

      - name: Upload artifacts
        uses: actions/upload-artifact@v4
        with:
          name: hanimo-${{ matrix.platform }}
          path: desktop/src-tauri/target/release/bundle/**/*
```

- [ ] **Step 3: 로컬 빌드 테스트**

Run:
```bash
cd desktop
npm run build
cd src-tauri
cargo tauri build
```
Expected: `target/release/bundle/` 에 플랫폼별 설치 파일 생성

- [ ] **Step 4: Commit**

```bash
git add desktop/src-tauri/icons/ .github/workflows/desktop-build.yml
git commit -m "feat(desktop): add build pipeline and app icons

GitHub Actions workflow for cross-platform builds (macOS arm64/x64,
Linux, Windows). Tauri icon set for all platforms."
```

---

## Task 9: 통합 테스트 & 최종 검증

**Files:**
- Modify: `desktop/package.json` (test script 추가)

- [ ] **Step 1: 프론트엔드 스토어 테스트 설정**

Run: `cd desktop && npm install -D vitest jsdom @testing-library/react @testing-library/jest-dom`

Add to `desktop/package.json` scripts:
```json
"test": "vitest run",
"test:watch": "vitest"
```

Add `desktop/vite.config.ts` test config:

```typescript
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

const host = process.env.TAURI_DEV_HOST;

export default defineConfig({
  plugins: [react(), tailwindcss()],
  clearScreen: false,
  server: {
    port: 1420,
    strictPort: true,
    host: host || false,
    hmr: host ? { protocol: "ws", host, port: 1421 } : undefined,
    watch: { ignored: ["**/src-tauri/**"] },
  },
  test: {
    environment: "jsdom",
    globals: true,
  },
});
```

- [ ] **Step 2: chat-store 테스트**

Create `desktop/src/__tests__/chat-store.test.ts`:

```typescript
import { describe, it, expect, beforeEach } from "vitest";
import { useChatStore } from "../stores/chat-store";

describe("ChatStore", () => {
  beforeEach(() => {
    useChatStore.setState({
      messages: [],
      isStreaming: false,
      streamingContent: "",
    });
  });

  it("adds a user message", () => {
    useChatStore.getState().addMessage({ role: "user", content: "hello" });
    const msgs = useChatStore.getState().messages;
    expect(msgs).toHaveLength(1);
    expect(msgs[0].role).toBe("user");
    expect(msgs[0].content).toBe("hello");
  });

  it("manages streaming state", () => {
    const store = useChatStore.getState();
    store.setStreaming(true);
    store.appendStreamingContent("hello ");
    store.appendStreamingContent("world");
    expect(useChatStore.getState().isStreaming).toBe(true);
    expect(useChatStore.getState().streamingContent).toBe("hello world");
  });

  it("clears messages", () => {
    useChatStore.getState().addMessage({ role: "user", content: "test" });
    useChatStore.getState().clear();
    expect(useChatStore.getState().messages).toHaveLength(0);
  });
});
```

- [ ] **Step 3: editor-store 테스트**

Create `desktop/src/__tests__/editor-store.test.ts`:

```typescript
import { describe, it, expect, beforeEach } from "vitest";
import { useEditorStore, getLanguageFromPath } from "../stores/editor-store";

describe("EditorStore", () => {
  beforeEach(() => {
    useEditorStore.setState({
      openFiles: [],
      activeFilePath: null,
      projectRoot: null,
    });
  });

  it("opens a file and sets it active", () => {
    useEditorStore.getState().openFile({
      path: "/test/foo.ts",
      name: "foo.ts",
      content: "const x = 1;",
      language: "typescript",
      isDirty: false,
    });
    expect(useEditorStore.getState().openFiles).toHaveLength(1);
    expect(useEditorStore.getState().activeFilePath).toBe("/test/foo.ts");
  });

  it("does not duplicate open files", () => {
    const file = {
      path: "/test/foo.ts",
      name: "foo.ts",
      content: "",
      language: "typescript",
      isDirty: false,
    };
    useEditorStore.getState().openFile(file);
    useEditorStore.getState().openFile(file);
    expect(useEditorStore.getState().openFiles).toHaveLength(1);
  });

  it("closes a file and updates active", () => {
    const store = useEditorStore.getState();
    store.openFile({ path: "/a.ts", name: "a.ts", content: "", language: "typescript", isDirty: false });
    store.openFile({ path: "/b.ts", name: "b.ts", content: "", language: "typescript", isDirty: false });
    store.closeFile("/b.ts");
    expect(useEditorStore.getState().openFiles).toHaveLength(1);
    expect(useEditorStore.getState().activeFilePath).toBe("/a.ts");
  });

  it("marks file as dirty on content update", () => {
    useEditorStore.getState().openFile({
      path: "/test.ts",
      name: "test.ts",
      content: "old",
      language: "typescript",
      isDirty: false,
    });
    useEditorStore.getState().updateFileContent("/test.ts", "new");
    expect(useEditorStore.getState().openFiles[0].isDirty).toBe(true);
  });
});

describe("getLanguageFromPath", () => {
  it("maps .ts to typescript", () => {
    expect(getLanguageFromPath("foo.ts")).toBe("typescript");
  });

  it("maps .py to python", () => {
    expect(getLanguageFromPath("script.py")).toBe("python");
  });

  it("returns plaintext for unknown extensions", () => {
    expect(getLanguageFromPath("file.xyz")).toBe("plaintext");
  });
});
```

- [ ] **Step 4: 테스트 실행**

Run: `cd desktop && npm test`
Expected: All tests pass

- [ ] **Step 5: 전체 앱 실행 검증**

Run: `cd desktop && npm run tauri dev`

Verify:
1. 첫 실행 → 온보딩 위저드 표시
2. 프로바이더/키/모델 설정 완료 → 메인 앱 전환
3. 사이드바 아이콘 클릭 동작
4. 테마 토글 (다크 ↔ 라이트)
5. Open Folder → 파일 트리 표시
6. 파일 클릭 → Monaco 에디터에서 열림
7. 탭 전환/닫기
8. Cmd+S 저장
9. 채팅 메시지 전송 (sidecar 연결 시)

- [ ] **Step 6: Final commit**

```bash
git add desktop/
git commit -m "test(desktop): add store unit tests and verify MVP integration

Chat store and editor store tests. Verified: onboarding wizard, theme
toggle, file tree, Monaco editor, tab management, chat panel."
```

---

## Summary

| Task | Description | Estimated Steps |
|------|-------------|----------------|
| 1 | Tauri 프로젝트 초기화 | 12 |
| 2 | Node.js Sidecar 통합 | 8 |
| 3 | 테마 시스템 | 3 |
| 4 | VSCode 레이아웃 셸 | 8 |
| 5 | 채팅 기능 | 7 |
| 6 | 파일 트리 & 에디터 | 8 |
| 7 | 온보딩 위저드 | 5 |
| 8 | 빌드 파이프라인 | 4 |
| 9 | 통합 테스트 & 검증 | 6 |
| **Total** | | **61 steps** |
