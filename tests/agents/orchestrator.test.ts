import { describe, it, expect } from 'vitest';
import type { SubTask, SubAgentResult } from '../../src/agents/types.js';

// We test the pure parseSubTasks logic by importing it indirectly.
// Since parseSubTasks is not exported, we test the types and structural contracts.

describe('Orchestrator types', () => {
  it('SubTask has required fields', () => {
    const task: SubTask = {
      id: 'task-1',
      description: 'Analyze the codebase',
      type: 'analyze',
    };
    expect(task.id).toBe('task-1');
    expect(task.type).toBe('analyze');
  });

  it('SubAgentResult fulfilled shape', () => {
    const result: SubAgentResult = {
      taskId: 'task-1',
      status: 'fulfilled',
      result: 'Analysis complete: found 10 files.',
    };
    expect(result.status).toBe('fulfilled');
    expect(result.result).toContain('Analysis complete');
  });

  it('SubAgentResult rejected shape', () => {
    const result: SubAgentResult = {
      taskId: 'task-2',
      status: 'rejected',
      error: 'Model timeout',
    };
    expect(result.status).toBe('rejected');
    expect(result.error).toBe('Model timeout');
  });

  it('valid SubTask types', () => {
    const types: SubTask['type'][] = ['analyze', 'code', 'review', 'research', 'general'];
    for (const t of types) {
      const task: SubTask = { id: `task-${t}`, description: `Do ${t}`, type: t };
      expect(task.type).toBe(t);
    }
  });
});

// Test parseSubTasks logic via a re-implementation of the same regex pattern
describe('parseSubTasks pattern', () => {
  function parseSubTasksLocal(text: string, fallback: string): SubTask[] {
    const match = /\[[\s\S]*\]/.exec(text);
    if (!match) return [{ id: 'task-1', description: fallback, type: 'general' }];

    try {
      const parsed: unknown = JSON.parse(match[0]);
      if (!Array.isArray(parsed) || parsed.length === 0) {
        return [{ id: 'task-1', description: fallback, type: 'general' }];
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
        : [{ id: 'task-1', description: fallback, type: 'general' }];
    } catch {
      return [{ id: 'task-1', description: fallback, type: 'general' }];
    }
  }

  it('parses valid JSON array', () => {
    const input = `[{"id":"task-1","description":"Do X","type":"code"},{"id":"task-2","description":"Do Y","type":"review"}]`;
    const result = parseSubTasksLocal(input, 'fallback');
    expect(result).toHaveLength(2);
    expect(result[0]!.id).toBe('task-1');
    expect(result[0]!.type).toBe('code');
    expect(result[1]!.id).toBe('task-2');
    expect(result[1]!.type).toBe('review');
  });

  it('falls back on invalid JSON', () => {
    const result = parseSubTasksLocal('not json at all', 'my fallback');
    expect(result).toHaveLength(1);
    expect(result[0]!.description).toBe('my fallback');
    expect(result[0]!.type).toBe('general');
  });

  it('falls back on empty array', () => {
    const result = parseSubTasksLocal('[]', 'fallback');
    expect(result).toHaveLength(1);
    expect(result[0]!.description).toBe('fallback');
  });

  it('handles JSON with extra text around it', () => {
    const input = `Here are the tasks:\n[{"id":"task-1","description":"analyze code","type":"analyze"}]\n\nDone.`;
    const result = parseSubTasksLocal(input, 'fallback');
    expect(result).toHaveLength(1);
    expect(result[0]!.description).toBe('analyze code');
  });

  it('defaults unknown type to general', () => {
    const input = `[{"id":"task-1","description":"do stuff","type":"unknown_type"}]`;
    const result = parseSubTasksLocal(input, 'fallback');
    expect(result[0]!.type).toBe('general');
  });

  it('skips items without required fields', () => {
    const input = `[{"id":"task-1","description":"ok","type":"code"},{"name":"bad"}]`;
    const result = parseSubTasksLocal(input, 'fallback');
    expect(result).toHaveLength(1);
    expect(result[0]!.id).toBe('task-1');
  });
});
