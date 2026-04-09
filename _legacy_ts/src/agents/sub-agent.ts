import type { LanguageModelV1, ToolSet } from 'ai';
import { runAgentLoop } from '../core/agent-loop.js';
import type { SubTask, SubAgentResult } from './types.js';

export class SubAgent {
  constructor(
    private model: LanguageModelV1,
    private tools: ToolSet | undefined,
    private maxSteps: number,
  ) {}

  async runTask(task: SubTask, abortSignal?: AbortSignal): Promise<SubAgentResult> {
    const systemPrompt = `You are a sub-agent working on a specific task. Focus only on your assigned task and provide a clear, concise result.\n\nTask: ${task.description}\nTask type: ${task.type}`;

    try {
      const result = await runAgentLoop({
        model: this.model,
        systemPrompt,
        messages: [{ role: 'user' as const, content: task.description }],
        tools: this.tools,
        maxSteps: this.maxSteps,
        abortSignal,
      });
      return { taskId: task.id, status: 'fulfilled', result: result.response };
    } catch (error) {
      const msg = error instanceof Error ? error.message : String(error);
      return { taskId: task.id, status: 'rejected', error: msg };
    }
  }
}
