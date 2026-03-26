import { createInterface } from 'node:readline/promises';
import { stdin, stdout } from 'node:process';
import { execFile } from 'node:child_process';
import { promisify } from 'node:util';
import type { LanguageModelV1, ToolSet } from 'ai';
import { runAgentLoop } from './core/agent-loop.js';
import { getModel, clearProviderCache } from './providers/registry.js';
import { loadConfig } from './config/loader.js';
import { PROVIDER_NAMES } from './providers/types.js';
import type { ProviderName } from './providers/types.js';
import type { Message, AgentEvent, TokenUsage } from './core/types.js';

const execFileAsync = promisify(execFile);

const LOCAL_PROVIDERS = new Set(['ollama', 'vllm', 'custom']);

// ANSI helpers
const dim = (s: string): string => `\x1b[2m${s}\x1b[0m`;
const cyan = (s: string): string => `\x1b[36m${s}\x1b[0m`;
const green = (s: string): string => `\x1b[32m${s}\x1b[0m`;
const yellow = (s: string): string => `\x1b[33m${s}\x1b[0m`;
const red = (s: string): string => `\x1b[31m${s}\x1b[0m`;
const bold = (s: string): string => `\x1b[1m${s}\x1b[0m`;
const magenta = (s: string): string => `\x1b[35m${s}\x1b[0m`;

// All slash commands for tab completion
const SLASH_COMMANDS = [
  '/help', '/h',
  '/model', '/m',
  '/models',
  '/provider', '/p',
  '/providers',
  '/tools', '/t',
  '/config',
  '/usage', '/u',
  '/clear',
  '/lang',
  '/system',
  '/setup',
  '/exit', '/quit', '/q',
];

const LANGUAGES: Record<string, string> = {
  ko: '한국어로 응답해주세요.',
  en: 'Respond in English.',
  ja: '日本語で応答してください。',
  zh: '请用中文回答。',
};

interface TextModeOptions {
  provider: string;
  model: string;
  modelInstance: LanguageModelV1;
  systemPrompt: string;
  tools: ToolSet;
  initialPrompt?: string;
}

async function listOllamaModels(): Promise<string[]> {
  try {
    const { stdout: output } = await execFileAsync('ollama', ['list'], { timeout: 10000 });
    const lines = output.trim().split('\n');
    if (lines.length <= 1) return [];
    const models: string[] = [];
    for (let i = 1; i < lines.length; i++) {
      const line = lines[i];
      if (!line) continue;
      const parts = line.trim().split(/\s{2,}/);
      const name = parts[0];
      if (name) models.push(name);
    }
    return models;
  } catch {
    return [];
  }
}

export async function startTextMode(options: TextModeOptions): Promise<void> {
  let currentProvider = options.provider;
  let currentModel = options.model;
  let currentModelInstance = options.modelInstance;
  let baseSystemPrompt = options.systemPrompt;
  const { tools, initialPrompt } = options;

  // State
  let toolsEnabled = !LOCAL_PROVIDERS.has(currentProvider);
  let langSuffix = '';
  let cachedOllamaModels: string[] = [];

  const messages: Message[] = [];
  const totalUsage: TokenUsage = { promptTokens: 0, completionTokens: 0, totalTokens: 0 };

  // Pre-cache Ollama models for tab completion
  if (currentProvider === 'ollama') {
    listOllamaModels().then(m => { cachedOllamaModels = m; }).catch(() => {});
  }

  function getSystemPrompt(): string {
    return langSuffix ? `${baseSystemPrompt}\n\n${langSuffix}` : baseSystemPrompt;
  }

  // Tab completion
  function completer(line: string): [string[], string] {
    // Slash command completion
    if (line.startsWith('/')) {
      const parts = line.split(' ');

      // Complete sub-arguments
      if (parts.length >= 2) {
        const cmd = parts[0];
        const partial = parts.slice(1).join(' ');

        if (cmd === '/provider' || cmd === '/p') {
          const hits = PROVIDER_NAMES
            .filter(p => p.startsWith(partial))
            .map(p => `${cmd} ${p}`);
          return [hits.length ? hits : PROVIDER_NAMES.map(p => `${cmd} ${p}`), line];
        }

        if (cmd === '/model' || cmd === '/m') {
          if (cachedOllamaModels.length > 0) {
            const hits = cachedOllamaModels
              .filter(m => m.startsWith(partial))
              .map(m => `${cmd} ${m}`);
            return [hits.length ? hits : cachedOllamaModels.map(m => `${cmd} ${m}`), line];
          }
          return [[], line];
        }

        if (cmd === '/tools' || cmd === '/t') {
          const opts = ['on', 'off'];
          const hits = opts.filter(o => o.startsWith(partial)).map(o => `${cmd} ${o}`);
          return [hits, line];
        }

        if (cmd === '/lang') {
          const langs = Object.keys(LANGUAGES);
          const hits = langs.filter(l => l.startsWith(partial)).map(l => `/lang ${l}`);
          return [hits.length ? hits : langs.map(l => `/lang ${l}`), line];
        }

        return [[], line];
      }

      // Complete command name
      const hits = SLASH_COMMANDS.filter(c => c.startsWith(line));
      return [hits.length ? hits : SLASH_COMMANDS, line];
    }

    return [[], line];
  }

  function printBanner(): void {
    console.log();
    console.log(`  ${bold('devany')} v0.1.0  ${dim('(')}${green(currentProvider)}${dim('/')}${cyan(currentModel)}${dim(')')}`);
    const toolsTag = toolsEnabled ? green('tools:ON') : dim('tools:OFF');
    console.log(`  ${dim('종료: Ctrl+C  |  Tab: 자동완성  |')} ${toolsTag}`);
    console.log(dim('─'.repeat(50)));
    console.log();
  }

  printBanner();

  const rl = createInterface({
    input: stdin,
    output: stdout,
    completer,
  });

  let isRunning = false;
  let abortController: AbortController | null = null;

  // Ctrl+C handling
  process.on('SIGINT', () => {
    if (isRunning && abortController) {
      abortController.abort();
      console.log(`\n  ${dim('(취소됨)')}`);
      isRunning = false;
      abortController = null;
    } else {
      console.log('\n  Bye!');
      rl.close();
      process.exit(0);
    }
  });

  async function switchModel(newModel: string): Promise<boolean> {
    try {
      const config = await loadConfig();
      const providerConfig = config.providers?.[currentProvider] ?? {};
      clearProviderCache();
      currentModelInstance = getModel(currentProvider as ProviderName, newModel, providerConfig);
      currentModel = newModel;
      console.log(`  ${green('✓')} 모델 변경: ${cyan(newModel)}`);
      return true;
    } catch (err) {
      const msg = err instanceof Error ? err.message : String(err);
      console.log(`  ${red('✗')} 모델 변경 실패: ${msg}`);
      return false;
    }
  }

  async function switchProvider(newProvider: string, newModel?: string): Promise<boolean> {
    if (!PROVIDER_NAMES.includes(newProvider as ProviderName)) {
      console.log(`  ${red('✗')} 알 수 없는 프로바이더: ${newProvider}`);
      console.log(`  ${dim('사용 가능:')} ${PROVIDER_NAMES.join(', ')}`);
      return false;
    }
    try {
      const config = await loadConfig();
      const providerConfig = config.providers?.[newProvider] ?? {};
      const modelId = newModel ?? config.model;
      clearProviderCache();
      currentModelInstance = getModel(newProvider as ProviderName, modelId, providerConfig);
      currentProvider = newProvider;
      currentModel = modelId;
      toolsEnabled = !LOCAL_PROVIDERS.has(newProvider);

      // Refresh Ollama model cache
      if (newProvider === 'ollama') {
        listOllamaModels().then(m => { cachedOllamaModels = m; }).catch(() => {});
      }

      console.log(`  ${green('✓')} ${green(newProvider)}/${cyan(modelId)}`);
      console.log(`  ${dim('도구:')} ${toolsEnabled ? green('ON') : red('OFF')}`);
      return true;
    } catch (err) {
      const msg = err instanceof Error ? err.message : String(err);
      console.log(`  ${red('✗')} 프로바이더 변경 실패: ${msg}`);
      return false;
    }
  }

  async function handleInput(input: string): Promise<boolean> {
    const trimmed = input.trim();
    if (!trimmed) return true;

    if (trimmed.startsWith('/')) {
      const parts = trimmed.slice(1).split(' ');
      const cmd = parts[0];
      const args = parts.slice(1);

      switch (cmd) {
        // ── 도움말 ──
        case 'help':
        case 'h':
          console.log();
          console.log(`  ${bold('명령어')} ${dim('(Tab으로 자동완성)')}`);
          console.log();
          console.log(`  ${magenta('모델 & 프로바이더')}`);
          console.log(`    ${cyan('/model')} <name>        ${dim('모델 변경 (즉시 적용)')}`);
          console.log(`    ${cyan('/models')}              ${dim('설치된 모델 목록 (Ollama)')}`);
          console.log(`    ${cyan('/provider')} <name>     ${dim('프로바이더 변경')}`);
          console.log(`    ${cyan('/providers')}           ${dim('프로바이더 목록')}`);
          console.log();
          console.log(`  ${magenta('설정')}`);
          console.log(`    ${cyan('/tools')} [on|off]      ${dim('도구(파일/Git/셸) 토글')}`);
          console.log(`    ${cyan('/lang')} <ko|en|ja|zh>  ${dim('응답 언어 변경')}`);
          console.log(`    ${cyan('/system')} <text>       ${dim('시스템 프롬프트 추가')}`);
          console.log(`    ${cyan('/config')}              ${dim('현재 설정 보기')}`);
          console.log(`    ${cyan('/setup')}               ${dim('초기 설정 다시 실행')}`);
          console.log();
          console.log(`  ${magenta('대화')}`);
          console.log(`    ${cyan('/clear')}               ${dim('대화 초기화')}`);
          console.log(`    ${cyan('/usage')}               ${dim('토큰 사용량')}`);
          console.log();
          console.log(`  ${magenta('종료')}`);
          console.log(`    ${cyan('/exit')}                ${dim('프로그램 종료')}`);
          console.log(`    ${dim('Ctrl+C')}               ${dim('응답 취소 / 종료')}`);
          console.log();
          console.log(`  ${dim('단축: /m=/model /p=/provider /t=/tools /u=/usage /q=/exit')}`);
          console.log();
          return true;

        // ── 모델 ──
        case 'model':
        case 'm':
          if (args.length > 0) {
            await switchModel(args.join(' '));
          } else {
            console.log(`  현재 모델: ${cyan(currentModel)}`);
            console.log(`  ${dim('변경: /model <name> | Tab으로 자동완성')}`);
          }
          return true;

        case 'models': {
          if (currentProvider === 'ollama') {
            console.log(`  ${dim('Ollama 모델 조회 중...')}`);
            const models = await listOllamaModels();
            cachedOllamaModels = models;
            if (models.length > 0) {
              console.log();
              for (const [i, m] of models.entries()) {
                const marker = m === currentModel ? green(' ●') : '  ';
                console.log(`  ${marker} ${dim(String(i + 1) + ')')} ${m}`);
              }
              console.log();
              console.log(`  ${dim('/model <name> 으로 변경 | Tab으로 자동완성')}`);
            } else {
              console.log(`  ${red('모델 없음')} — ${dim('ollama pull <model> 로 설치')}`);
            }
          } else {
            console.log(`  ${dim('Ollama 외 프로바이더는 /model <name> 으로 직접 변경')}`);
          }
          console.log();
          return true;
        }

        // ── 프로바이더 ──
        case 'provider':
        case 'p':
          if (args.length > 0) {
            const newProvider = args[0] ?? '';
            await switchProvider(newProvider, args[1]);
          } else {
            console.log(`  현재: ${green(currentProvider)}/${cyan(currentModel)}`);
            console.log(`  ${dim('변경: /provider <name> [model] | Tab으로 자동완성')}`);
          }
          return true;

        case 'providers':
          console.log();
          for (const p of PROVIDER_NAMES) {
            const marker = p === currentProvider ? green(' ●') : '  ';
            const local = LOCAL_PROVIDERS.has(p) ? dim(' (로컬)') : '';
            console.log(`  ${marker} ${p}${local}`);
          }
          console.log();
          return true;

        // ── 도구 ──
        case 'tools':
        case 't': {
          const arg = args[0]?.toLowerCase();
          if (arg === 'on') {
            toolsEnabled = true;
          } else if (arg === 'off') {
            toolsEnabled = false;
          } else {
            toolsEnabled = !toolsEnabled;
          }
          if (toolsEnabled) {
            console.log(`  ${green('✓')} 도구 ${green('ON')} — 파일 읽기/쓰기, Git, 셸 명령 사용 가능`);
            if (LOCAL_PROVIDERS.has(currentProvider)) {
              console.log(`  ${yellow('⚠')} ${dim('로컬 모델은 tool calling 미지원일 수 있음')}`);
            }
          } else {
            console.log(`  ${green('✓')} 도구 ${red('OFF')} — 순수 대화 모드`);
          }
          return true;
        }

        // ── 언어 ──
        case 'lang': {
          const lang = args[0]?.toLowerCase();
          if (lang && LANGUAGES[lang]) {
            langSuffix = LANGUAGES[lang];
            const names: Record<string, string> = { ko: '한국어', en: 'English', ja: '日本語', zh: '中文' };
            console.log(`  ${green('✓')} 응답 언어: ${names[lang] ?? lang}`);
          } else if (lang === 'off' || lang === 'auto') {
            langSuffix = '';
            console.log(`  ${green('✓')} 언어 제한 해제 (모델 기본)`);
          } else {
            console.log(`  사용법: ${cyan('/lang')} <${Object.keys(LANGUAGES).join('|')}|off>`);
            if (langSuffix) {
              console.log(`  ${dim('현재:')} ${langSuffix}`);
            }
          }
          return true;
        }

        // ── 시스템 프롬프트 ──
        case 'system':
          if (args.length > 0) {
            const extra = args.join(' ');
            baseSystemPrompt = `${options.systemPrompt}\n\n${extra}`;
            console.log(`  ${green('✓')} 시스템 프롬프트 추가됨`);
            console.log(`  ${dim(extra.length > 80 ? extra.slice(0, 80) + '...' : extra)}`);
          } else {
            console.log(`  ${bold('시스템 프롬프트:')}`);
            const preview = getSystemPrompt();
            const lines = preview.split('\n').slice(0, 8);
            for (const line of lines) {
              console.log(`  ${dim(line)}`);
            }
            if (preview.split('\n').length > 8) console.log(`  ${dim('...')}`);
            console.log();
            console.log(`  ${dim('/system <text> 로 추가 | /system reset 으로 초기화')}`);
          }
          return true;

        // ── 설정 보기 ──
        case 'config':
          console.log();
          console.log(`  ${bold('현재 설정:')}`);
          console.log(`    프로바이더  ${green(currentProvider)}`);
          console.log(`    모델        ${cyan(currentModel)}`);
          console.log(`    도구        ${toolsEnabled ? green('ON') : red('OFF')}`);
          console.log(`    언어        ${langSuffix ? dim(langSuffix) : dim('auto')}`);
          console.log(`    대화        ${messages.length}개 메시지`);
          console.log(`    토큰        ${totalUsage.totalTokens.toLocaleString()}`);
          console.log();
          return true;

        // ── 초기 설정 ──
        case 'setup':
          console.log(`  ${dim('devany --setup 으로 재실행하세요.')}`);
          return true;

        // ── 토큰 ──
        case 'usage':
        case 'u':
          console.log();
          console.log(`  ${bold('토큰 사용량:')}`);
          console.log(`    prompt      ${totalUsage.promptTokens.toLocaleString()}`);
          console.log(`    completion  ${totalUsage.completionTokens.toLocaleString()}`);
          console.log(`    total       ${totalUsage.totalTokens.toLocaleString()}`);
          console.log(`    대화        ${messages.length}개 메시지`);
          console.log();
          return true;

        // ── 대화 초기화 ──
        case 'clear':
          messages.length = 0;
          console.log(`  ${green('✓')} 대화 초기화됨.`);
          return true;

        // ── 종료 ──
        case 'exit':
        case 'quit':
        case 'q':
          return false;

        default:
          console.log(`  ${dim('알 수 없는 명령어:')} /${cmd ?? ''} ${dim('— /help 참고')}`);
          return true;
      }
    }

    // ── 일반 메시지 → 에이전트 호출 ──
    messages.push({ role: 'user', content: trimmed });
    isRunning = true;
    abortController = new AbortController();

    process.stdout.write(`\n  ${cyan('▌')} `);

    const onEvent = (event: AgentEvent): void => {
      switch (event.type) {
        case 'token':
          process.stdout.write(event.content);
          break;
        case 'tool-call':
          process.stdout.write(`\n  ${yellow('⚡')} ${dim(event.toolName + '...')}`);
          break;
        case 'tool-result': {
          const preview = event.result.length > 200
            ? event.result.slice(0, 200) + '...'
            : event.result;
          const lines = preview.split('\n').slice(0, 5);
          process.stdout.write(`\n  ${dim('┃')} ${dim(lines.join('\n  ┃ '))}\n  ${cyan('▌')} `);
          break;
        }
        case 'done':
          totalUsage.promptTokens += event.usage.promptTokens;
          totalUsage.completionTokens += event.usage.completionTokens;
          totalUsage.totalTokens += event.usage.totalTokens;
          break;
        case 'error':
          process.stdout.write(`\n  ${red('✗')} ${event.error.message}`);
          break;
      }
    };

    try {
      const result = await runAgentLoop({
        model: currentModelInstance,
        systemPrompt: getSystemPrompt(),
        messages,
        tools: toolsEnabled ? tools : undefined,
        maxSteps: toolsEnabled ? 25 : 1,
        onEvent,
        abortSignal: abortController.signal,
      });

      if (result.response) {
        messages.push({ role: 'assistant', content: result.response });
      } else {
        console.log(`\n  ${yellow('⚠')} 빈 응답`);
        if (currentProvider === 'ollama') {
          console.log(`  ${dim('ollama ps 로 모델 상태 확인 | /models 로 목록 보기')}`);
        }
      }
    } catch (err) {
      if (err instanceof Error && err.name === 'AbortError') {
        // User cancelled
      } else {
        const msg = err instanceof Error ? err.message : String(err);
        console.error(`\n  ${red('✗')} ${msg}`);
        if (msg.includes('ECONNREFUSED')) {
          console.error(`  ${dim('→ ollama serve 로 서버 시작')}`);
        } else if (msg.includes('not found')) {
          console.error(`  ${dim('→ /models 로 사용 가능한 모델 확인')}`);
        }
      }
    } finally {
      isRunning = false;
      abortController = null;
      console.log('\n');
    }

    return true;
  }

  // Handle initial prompt
  if (initialPrompt) {
    const shouldContinue = await handleInput(initialPrompt);
    if (!shouldContinue) {
      rl.close();
      return;
    }
  }

  // Interactive loop
  while (true) {
    let input: string;
    try {
      input = await rl.question(`  ${dim(currentProvider + '/' + currentModel)} ${bold('❯')} `);
    } catch {
      break;
    }

    const shouldContinue = await handleInput(input);
    if (!shouldContinue) break;
  }

  rl.close();
  console.log('  Bye!');
}
