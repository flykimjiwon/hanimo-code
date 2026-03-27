import { describe, it, expect } from 'vitest';
import { todoTool } from '../src/tools/todo.js';

describe('todo tool', () => {
  it('should add a task', async () => {
    const result = await todoTool.execute({ action: 'add', task: 'Test task' }, { toolCallId: '1', messages: [], abortSignal: undefined as unknown as AbortSignal });
    expect(result).toMatchObject({ success: true, action: 'added' });
    expect((result as { item: { id: number } }).item.id).toBeGreaterThan(0);
  });

  it('should list tasks', async () => {
    const result = await todoTool.execute({ action: 'list' }, { toolCallId: '2', messages: [], abortSignal: undefined as unknown as AbortSignal });
    expect(result).toMatchObject({ success: true, action: 'list' });
    expect((result as { items: unknown[] }).items.length).toBeGreaterThan(0);
  });

  it('should update task status', async () => {
    const addResult = await todoTool.execute({ action: 'add', task: 'Update me' }, { toolCallId: '3', messages: [], abortSignal: undefined as unknown as AbortSignal });
    const id = (addResult as { item: { id: number } }).item.id;
    const result = await todoTool.execute({ action: 'update', id, status: 'done' }, { toolCallId: '4', messages: [], abortSignal: undefined as unknown as AbortSignal });
    expect(result).toMatchObject({ success: true, action: 'updated' });
    expect((result as { item: { status: string } }).item.status).toBe('done');
  });

  it('should remove a task', async () => {
    const addResult = await todoTool.execute({ action: 'add', task: 'Remove me' }, { toolCallId: '5', messages: [], abortSignal: undefined as unknown as AbortSignal });
    const id = (addResult as { item: { id: number } }).item.id;
    const result = await todoTool.execute({ action: 'remove', id }, { toolCallId: '6', messages: [], abortSignal: undefined as unknown as AbortSignal });
    expect(result).toMatchObject({ success: true, action: 'removed' });
  });

  it('should clear all tasks', async () => {
    // Add a couple tasks first
    await todoTool.execute({ action: 'add', task: 'clear test 1' }, { toolCallId: '8', messages: [], abortSignal: undefined as unknown as AbortSignal });
    await todoTool.execute({ action: 'add', task: 'clear test 2' }, { toolCallId: '9', messages: [], abortSignal: undefined as unknown as AbortSignal });

    const result = await todoTool.execute({ action: 'clear' }, { toolCallId: '10', messages: [], abortSignal: undefined as unknown as AbortSignal });
    expect(result).toMatchObject({ success: true, action: 'cleared' });

    const listResult = await todoTool.execute({ action: 'list' }, { toolCallId: '11', messages: [], abortSignal: undefined as unknown as AbortSignal });
    expect((listResult as { summary: { total: number } }).summary.total).toBe(0);
  });

  it('should fail without task description on add', async () => {
    const result = await todoTool.execute({ action: 'add' }, { toolCallId: '7', messages: [], abortSignal: undefined as unknown as AbortSignal });
    expect(result).toMatchObject({ success: false });
  });
});
