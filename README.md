# hanimo

> 오픈소스를 위한 AI 코딩 에이전트

하니모(hanimo)는 터미널에서 동작하는 AI 코딩 에이전트입니다. 단일 바이너리로 배포되며, 14개 이상의 LLM 프로바이더를 지원합니다.

## Features

- **14+ LLM Providers** — Ollama, OpenAI, Anthropic, Google, Novita, OpenRouter, DeepSeek, Groq, Together, Fireworks, Mistral, vLLM, LM Studio, Custom
- **3 Modes** — Super (all-purpose), Dev (coding), Plan (analysis)
- **Hash-Anchored Editing** — MD5 hash per line prevents stale-edit corruption
- **Smart Compaction** — 3-stage context compression for long sessions
- **Autonomous Mode** — `/auto` for self-driving task completion
- **SQLite Sessions** — Save, load, search, fork sessions
- **Project Memory** — Long-term per-project knowledge
- **MCP Support** — Extend tools via Model Context Protocol
- **Built-in Tools** — File R/W/E/D, shell, grep, glob, git, diagnostics
- **Command Palette** — Ctrl+K fuzzy search
- **Single Binary** — No runtime dependencies, cross-platform

## Quick Start

### Binary Download

```bash
# macOS (Apple Silicon)
curl -L -o hanimo https://github.com/flykimjiwon/hanimo/releases/latest/download/hanimo-darwin-arm64
chmod +x hanimo
sudo mv hanimo /usr/local/bin/
```

### Build from Source

```bash
git clone https://github.com/flykimjiwon/hanimo.git
cd hanimo
go build -o hanimo ./cmd/hanimo
./hanimo
```

### Usage

```bash
# With Ollama (default)
ollama pull qwen3:8b
hanimo

# With specific provider
hanimo -p openai -m gpt-4o
hanimo -p anthropic -m claude-sonnet-4
hanimo -p novita -m openai/gpt-oss-120b

# With custom endpoint
hanimo -u http://localhost:8000/v1 -m my-model
```

## Configuration

Config file: `~/.hanimo/config.yaml`

```yaml
default:
  provider: ollama
  model: qwen3:8b

providers:
  novita:
    api_key: "nvt-..."
  openrouter:
    api_key: "sk-or-..."
  anthropic:
    api_key: "sk-ant-..."
```

Environment variables: `HANIMO_API_KEY`, `HANIMO_API_BASE_URL`, `HANIMO_MODEL_SUPER`, `HANIMO_MODEL_DEV`

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| Enter | Send message |
| Shift+Enter | New line |
| Tab | Switch mode (Super → Dev → Plan) |
| Ctrl+K | Command palette |
| Esc | Cancel streaming |
| Ctrl+L | Clear screen |
| Ctrl+C | Exit |

## Slash Commands

| Command | Description |
|---------|-------------|
| `/help` | Show help |
| `/model [name]` | Switch model |
| `/provider [name]` | Switch provider |
| `/auto [task]` | Autonomous mode |
| `/save [name]` | Save session |
| `/load` | Load session |
| `/search [keyword]` | Search sessions |
| `/remember [text]` | Save memory |
| `/memories` | Show memories |
| `/diagnostics [file]` | Code diagnostics |
| `/usage` | Token usage |
| `/config` | Show config |
| `/theme [name]` | Change theme |
| `/clear` | Clear conversation |

## Themes

honey (default), ocean, dracula, nord, forest

## License

MIT
