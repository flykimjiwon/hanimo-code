import type { LanguageModelV1, ToolSet } from 'ai';
import { runAgentLoop } from '../core/agent-loop.js';
import { SubAgent } from './sub-agent.js';
import type { SubTask, SubAgentResult, OrchestratorEvent } from './types.js';
import type { AgentEvent, AgentLoopResult } from '../core/types.js';

interface OrchestratorOptions {
  model: LanguageModelV1;
  systemPrompt: string;
  messages: import('../core/types.js').Message[];
  tools?: ToolSet;
  maxSteps?: number;
  subAgentCount: number;
  subAgentModel: LanguageModelV1;
  onEvent?: (event: AgentEvent) => void;
  onOrchestratorEvent?: (event: OrchestratorEvent) => void;
  abortSignal?: AbortSignal;
}

const DECOMPOSE_PROMPT = `You are a task decomposition engine. Given a user request, break it into independent sub-tasks that can be worked on in parallel by separate AI agents.

Rules:
1. Each sub-task should be self-contained and focused.
2. Each task needs a clear description of what to do.
3. Assign a type from: analyze, code, review, research, general.

Respond with ONLY a valid JSON array (no markdown fences, no extra text):
[
  {
    "id": "task-1",
    "description": "Description of what to do",
    "type": "code"
  }
]

If the request is simple enough for one task, return a single-element array.`;

function parseSubTasks(text: string, fallbackDescription: string): SubTask[] {
  // Try to extract JSON array from the response
  const match = /\[[\s\S]*\]/.exec(text);
  if (!match) {
    return [{ id: 'task-1', description: fallbackDescription, type: 'general' }];
  }

  try {
    const parsed: unknown = JSON.parse(match[0]);
    if (!Array.isArray(parsed) || parsed.length === 0) {
      return [{ id: 'task-1', description: fallbackDescription, type: 'general' }];
    }

    const validTypes = new Set(['analyze', 'code', 'review', 'research', 'general']);
    const tasks: SubTask[] = [];

    for (const item of parsed) {
      if (
        typeof item === 'object' &&
        item !== null &&
        typeof (item as Record<string, unknown>)['id'] === 'string' &&
        typeof (item as Record<string, unknown>)['description'] === 'string'
      ) {
        const raw = item as Record<string, unknown>;
        const rawType = raw['type'];
        const type = typeof rawType === 'string' && validTypes.has(rawType)
          ? (rawType as SubTask['type'])
          : 'general';
        tasks.push({
          id: raw['id'] as string,
          description: raw['description'] as string,
          type,
        });
      }
    }

    return tasks.length > 0
      ? tasks
      : [{ id: 'task-1', description: fallbackDescription, type: 'general' }];
  } catch {
    return [{ id: 'task-1', description: fallbackDescription, type: 'general' }];
  }
}

export class Orchestrator {
  async run(options: OrchestratorOptions): Promise<AgentLoopResult> {
    const {
      model,
      systemPrompt,
      messages,
      tools,
      maxSteps = 25,
      subAgentCount,
      subAgentModel,
      onEvent,
      onOrchestratorEvent,
      abortSignal,
    } = options;

    const emit = (event: OrchestratorEvent): void => {
      onOrchestratorEvent?.(event);
    };

    // Extract the last user message as the primary request
    const userMessage = (() => {
      for (let i = messages.length - 1; i >= 0; i--) {
        const msg = messages[i];
        if (msg && msg.role === 'user') {
          return typeof msg.content === 'string' ? msg.content : JSON.stringify(msg.content);
        }
      }
      return '';
    })();

    // Step 1: Decompose
    emit({ type: 'decompose-start' });

    const decomposeResult = await runAgentLoop({
      model,
      systemPrompt: DECOMPOSE_PROMPT,
      messages: [{ role: 'user' as const, content: userMessage }],
      maxSteps: 1,
      abortSignal,
    });

    const tasks = parseSubTasks(decomposeResult.response, userMessage);
    // Limit to requested sub-agent count
    const cappedTasks = tasks.slice(0, subAgentCount);

    emit({ type: 'decompose-done', tasks: cappedTasks });

    // Step 2: Execute sub-agents in parallel
    const total = cappedTasks.length;
    const subAgentMaxSteps = Math.max(1, Math.floor(maxSteps / 2));

    const taskPromises = cappedTasks.map((task, index) => {
      emit({ type: 'sub-agent-start', taskId: task.id, index, total });
      const agent = new SubAgent(subAgentModel, tools, subAgentMaxSteps);
      return agent.runTask(task, abortSignal).then((result) => {
        emit({
          type: 'sub-agent-done',
          taskId: task.id,
          index,
          total,
          status: result.status,
        });
        return result;
      });
    });

    const settled = await Promise.allSettled(taskPromises);

    const results: SubAgentResult[] = settled.map((outcome, index) => {
      const task = cappedTasks[index];
      const taskId = task?.id ?? `task-${index}`;
      if (outcome.status === 'fulfilled') {
        return outcome.value;
      }
      const errMsg = outcome.reason instanceof Error
        ? outcome.reason.message
        : String(outcome.reason);
      return { taskId, status: 'rejected' as const, error: errMsg };
    });

    // Step 3: Synthesize
    emit({ type: 'synthesize-start', results });

    const resultsSummary = results
      .map((r) => {
        if (r.status === 'fulfilled') {
          return `Task ${r.taskId} (success):\n${r.result ?? ''}`;
        }
        return `Task ${r.taskId} (failed): ${r.error ?? 'unknown error'}`;
      })
      .join('\n\n---\n\n');

    const synthesizeMessages: import('../core/types.js').Message[] = [
      ...messages,
      {
        role: 'user' as const,
        content: `Sub-agent results:\n\n${resultsSummary}\n\nPlease synthesize these results into a single coherent response to the original request.`,
      },
    ];

    const finalResult = await runAgentLoop({
      model,
      systemPrompt,
      messages: synthesizeMessages,
      tools,
      maxSteps,
      onEvent,
      abortSignal,
    });

    emit({ type: 'synthesize-done' });

    return finalResult;
  }
}
