/** @jsxImportSource @opentui/react */
/**
 * Hanimo TUI — built with OpenTUI React
 */
import { createCliRenderer, TextAttributes, SyntaxStyle } from '@opentui/core';
import { createRoot } from '@opentui/react';
import { useState, useCallback, useRef } from 'react';
import type { LanguageModelV1, ToolSet } from 'ai';
import { runAgentLoop, estimateCost } from '../core/agent-loop.js';
import type { AgentEvent, Message } from '../core/types.js';
import type { RoleManager } from '../roles/role-manager.js';
import type { RoleDefinition } from '../roles/types.js';
import { COMMAND_LIST } from '../tui/hooks/use-commands.js';

// ── Types ──

interface DisplayMessage {
  id: string;
  role: 'user' | 'assistant' | 'tool-call' | 'tool-result' | 'system';
  content: string;
  toolName?: string;
}

interface AppProps {
  provider: string;
  model: string;
  modelInstance: LanguageModelV1;
  systemPrompt: string;
  tools?: ToolSet;
  maxSteps?: number;
  initialPrompt?: string;
  activeRole?: RoleDefinition;
  streaming?: boolean;
}

// ── Helpers ──

let nextId = 0;
function msgId(): string { return `m-${++nextId}`; }

function formatTokens(n: number): string {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`;
  if (n >= 1000) return `${(n / 1000).toFixed(1)}k`;
  return String(n);
}

function formatCost(cost: number): string {
  if (cost === 0) return '$0';
  if (cost < 0.01) return `$${cost.toFixed(4)}`;
  return `$${cost.toFixed(2)}`;
}

// ── App Component ──

function App({
  provider,
  model,
  modelInstance,
  systemPrompt,
  tools,
  maxSteps = 25,
  activeRole,
  streaming = true,
}: AppProps) {
  const [messages, setMessages] = useState<DisplayMessage[]>([]);
  const [inputValue, setInputValue] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [streamingText, setStreamingText] = useState('');
  const [isStreaming, setIsStreaming] = useState(false);
  const [usage, setUsage] = useState({ promptTokens: 0, completionTokens: 0, totalCost: 0 });
  const [currentTool, setCurrentTool] = useState<string | null>(null);

  const conversationRef = useRef<Message[]>([]);
  const abortRef = useRef<AbortController | null>(null);

  const syntaxStyle = useRef(SyntaxStyle.create()).current;
  const roleIcon = activeRole?.icon ?? '\u26A1';
  const roleName = activeRole?.name ?? 'Hanimo';
  const cwd = process.cwd().replace(process.env['HOME'] ?? '', '~');

  const sendMessage = useCallback(async (text: string) => {
    if (!text.trim() || isLoading) return;

    const userMsg: DisplayMessage = { id: msgId(), role: 'user', content: text };
    setMessages((prev: DisplayMessage[]) => [...prev, userMsg]);
    setInputValue('');
    setIsLoading(true);
    setStreamingText('');
    setIsStreaming(true);

    conversationRef.current.push({ role: 'user', content: text });
    const controller = new AbortController();
    abortRef.current = controller;

    try {
      const result = await runAgentLoop({
        model: modelInstance,
        systemPrompt,
        messages: conversationRef.current,
        tools,
        maxSteps,
        abortSignal: controller.signal,
        streaming,
        onEvent: (event: AgentEvent) => {
          switch (event.type) {
            case 'token':
              setStreamingText((prev: string) => prev + event.content);
              break;
            case 'tool-call':
              setCurrentTool(event.toolName);
              setMessages((prev: DisplayMessage[]) => [...prev, {
                id: msgId(), role: 'tool-call' as const, content: event.toolName, toolName: event.toolName,
              }]);
              break;
            case 'tool-result':
              setCurrentTool(null);
              setMessages((prev: DisplayMessage[]) => [...prev, {
                id: msgId(), role: 'tool-result' as const,
                content: event.isError ? `Error: ${event.result}` : event.result,
                toolName: event.toolName,
              }]);
              break;
            case 'done':
              if (event.usage) {
                const modelId = (modelInstance as unknown as { modelId?: string }).modelId ?? '';
                const cost = estimateCost(modelId, event.usage);
                setUsage((prev: { promptTokens: number; completionTokens: number; totalCost: number }) => ({
                  promptTokens: prev.promptTokens + event.usage.promptTokens,
                  completionTokens: prev.completionTokens + event.usage.completionTokens,
                  totalCost: prev.totalCost + cost.totalCost,
                }));
              }
              break;
          }
        },
      });

      if (result.response) {
        setMessages((prev: DisplayMessage[]) => [...prev, { id: msgId(), role: 'assistant' as const, content: result.response }]);
      }
      conversationRef.current = result.messages;
    } catch (err) {
      if (err instanceof Error && err.name === 'AbortError') {
        setMessages((prev: DisplayMessage[]) => [...prev, { id: msgId(), role: 'system' as const, content: '(\uCDE8\uC18C\uB428)' }]);
      } else {
        const msg = err instanceof Error ? err.message : String(err);
        setMessages((prev: DisplayMessage[]) => [...prev, { id: msgId(), role: 'system' as const, content: `\u274C ${msg}` }]);
      }
    } finally {
      setIsLoading(false);
      setStreamingText('');
      setIsStreaming(false);
      setCurrentTool(null);
      abortRef.current = null;
    }
  }, [isLoading, modelInstance, systemPrompt, tools, maxSteps, streaming]);

  // Command matches
  const commandMatches = inputValue.startsWith('/')
    ? COMMAND_LIST.filter(c => c.name.startsWith(inputValue.slice(1).toLowerCase())).slice(0, 6)
    : [];

  const totalTokens = usage.promptTokens + usage.completionTokens;
  const statusRight = `${cwd}  ${model}  ${formatTokens(totalTokens)} tok  ${formatCost(usage.totalCost)}`;

  const hasMessages = messages.length > 0 || isLoading;

  return (
    <box flexDirection="column" flexGrow={1}>
      {/* Chat area */}
      <scrollbox flexGrow={1}>
        {!hasMessages ? (
          <box justifyContent="center" alignItems="center" flexGrow={1}>
            <box flexDirection="column" alignItems="center" gap={1}>
              <text content={[
                '  \u2584\u2580\u2580\u2580\u2580\u2584    \u2588\u2588   \u2588\u2588  \u2588\u2588\u2588\u2588\u2588  \u2588\u2588\u2588   \u2588\u2588 \u2588\u2588 \u2588\u2588\u2588   \u2588\u2588\u2588  \u2588\u2588\u2588\u2588\u2588\u2588',
                '  \u2588\u25D5\u1D25\u25D5 \u2588    \u2588\u2588   \u2588\u2588 \u2588\u2588   \u2588\u2588 \u2588\u2588\u2588\u2588  \u2588\u2588 \u2588\u2588 \u2588\u2588\u2588\u2588 \u2588\u2588\u2588\u2588 \u2588\u2588    \u2588\u2588',
                '  \u2580\u2584\u2584\u2584\u2584\u2580    \u2588\u2588\u2588\u2588\u2588\u2588\u2588 \u2588\u2588\u2588\u2588\u2588\u2588\u2588 \u2588\u2588 \u2588\u2588 \u2588\u2588 \u2588\u2588 \u2588\u2588 \u2588\u2588\u2588 \u2588\u2588 \u2588\u2588    \u2588\u2588',
                '  \u224B\u2588  \u2588\u224B    \u2588\u2588   \u2588\u2588 \u2588\u2588   \u2588\u2588 \u2588\u2588  \u2588\u2588\u2588\u2588 \u2588\u2588 \u2588\u2588     \u2588\u2588 \u2588\u2588    \u2588\u2588',
                '    \u2580\u2580      \u2588\u2588   \u2588\u2588 \u2588\u2588   \u2588\u2588 \u2588\u2588   \u2588\u2588\u2588 \u2588\u2588 \u2588\u2588     \u2588\u2588  \u2588\u2588\u2588\u2588\u2588\u2588',
              ].join('\n')} />
              <text attributes={TextAttributes.DIM} content={`\u2500\u2500\u2500 ${provider}/${model}  \u00B7  ${roleIcon} ${roleName} \u2500\u2500\u2500`} />
            </box>
          </box>
        ) : (
          <box flexDirection="column" paddingX={1}>
            {messages.map((msg) => {
              switch (msg.role) {
                case 'user':
                  return (
                    <box key={msg.id}>
                      <text attributes={TextAttributes.BOLD} fg="cyan" content={'\u276F '} />
                      <text attributes={TextAttributes.BOLD} content={msg.content} />
                    </box>
                  );
                case 'assistant':
                  return (
                    <box key={msg.id}>
                      <markdown syntaxStyle={syntaxStyle} content={msg.content}  />
                    </box>
                  );
                case 'tool-call':
                  return (
                    <box key={msg.id}>
                      <text fg="yellow" content={`\u26A1 ${msg.toolName}`} />
                    </box>
                  );
                case 'tool-result': {
                  const lines = msg.content.split('\n');
                  const preview = lines.length > 5 ? lines.slice(0, 5).join('\n') + `\n... (${lines.length - 5} more)` : msg.content;
                  return (
                    <box key={msg.id} flexDirection="column">
                      <text attributes={TextAttributes.DIM} content={`\u2514 ${msg.toolName}`} />
                      <box paddingLeft={2}>
                        <text attributes={TextAttributes.DIM} fg={msg.content.startsWith('Error:') ? 'red' : undefined} content={preview} />
                      </box>
                    </box>
                  );
                }
                case 'system':
                  return (
                    <box key={msg.id}>
                      <text attributes={TextAttributes.DIM} content={msg.content} />
                    </box>
                  );
              }
            })}

            {/* Streaming markdown */}
            {streamingText && (
              <box>
                <markdown syntaxStyle={syntaxStyle} content={streamingText} streaming={isStreaming}  />
              </box>
            )}

            {/* Loading */}
            {isLoading && !streamingText && (
              <box>
                <text fg="cyan" content={currentTool ? `\u26A1 ${currentTool}...` : '\u25CF Thinking...'} />
              </box>
            )}
          </box>
        )}
      </scrollbox>

      {/* Command dropdown */}
      {commandMatches.length > 0 && (
        <box flexDirection="column" paddingX={2}>
          {commandMatches.map((cmd) => (
            <box key={cmd.name}>
              <text attributes={TextAttributes.DIM} content="  /" />
              <text attributes={TextAttributes.BOLD} content={cmd.name.padEnd(16)} />
              <text attributes={TextAttributes.DIM} content={cmd.descriptionKo} />
            </box>
          ))}
        </box>
      )}

      {/* Input */}
      <box borderStyle="rounded" paddingX={1} flexDirection="column">
        <box justifyContent="space-between">
          <text attributes={TextAttributes.BOLD} content={`${roleIcon} ${roleName}`} />
          <text attributes={TextAttributes.DIM} content={'Shift+Enter \u21B5  Tab \u21C4 mode'} />
        </box>
        <input
          value={inputValue}
          onChange={(v: string) => setInputValue(v)}
          onSubmit={(v: unknown) => sendMessage(String(v))}
          placeholder={isLoading ? '\uC751\uB2F5 \uB300\uAE30 \uC911...' : '\uBA54\uC2DC\uC9C0\uB97C \uC785\uB825\uD558\uC138\uC694...'}
        />
      </box>

      {/* Bottom hints */}
      <box justifyContent="space-between" paddingX={1}>
        <text attributes={TextAttributes.DIM} content={isLoading ? 'Ctrl+C \uCDE8\uC18C' : 'Esc \uBA54\uB274'} />
        <text attributes={TextAttributes.DIM} content={statusRight} />
      </box>
    </box>
  );
}

// ── Entry Point ──

export interface StartAppOptions {
  provider: string;
  model: string;
  modelInstance: LanguageModelV1;
  systemPrompt: string;
  tools?: ToolSet;
  maxSteps?: number;
  initialPrompt?: string;
  providerConfig?: { apiKey?: string; baseURL?: string };
  roleManager?: RoleManager;
  activeRole?: RoleDefinition;
  networkMode?: string;
  streaming?: boolean;
}

export async function startApp(options: StartAppOptions): Promise<void> {
  const renderer = await createCliRenderer();
  const root = createRoot(renderer);

  root.render(
    <App
      provider={options.provider}
      model={options.model}
      modelInstance={options.modelInstance}
      systemPrompt={options.systemPrompt}
      tools={options.tools}
      maxSteps={options.maxSteps}
      initialPrompt={options.initialPrompt}
      activeRole={options.activeRole}
      streaming={options.streaming}
    />
  );
}
