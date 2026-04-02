import { describe, it, expect } from 'vitest';
import { RoleManager } from '../src/roles/role-manager.js';
import type { RoleDefinition } from '../src/roles/types.js';

describe('RoleManager', () => {
  it('loads 3 built-in roles on construction', () => {
    const mgr = new RoleManager();
    const roles = mgr.getAllRoles();
    expect(roles.length).toBe(3);

    const ids = roles.map((r) => r.id).sort();
    expect(ids).toEqual(['dev', 'hanimo', 'plan']);
  });

  it('getRole returns correct role by id', () => {
    const mgr = new RoleManager();

    const hanimo = mgr.getRole('hanimo');
    expect(hanimo).toBeDefined();
    expect(hanimo!.name).toBe('Hanimo');
    expect(hanimo!.icon).toBe('⚡');
    expect(hanimo!.tools.length).toBe(16);
    expect(hanimo!.maxSteps).toBe(50);

    const dev = mgr.getRole('dev');
    expect(dev).toBeDefined();
    expect(dev!.name).toBe('Dev');
    expect(dev!.icon).toBe('🔧');
    expect(dev!.tools.length).toBe(16);
    expect(dev!.maxSteps).toBe(25);

    const plan = mgr.getRole('plan');
    expect(plan).toBeDefined();
    expect(plan!.name).toBe('Plan');
    expect(plan!.icon).toBe('📋');
    expect(plan!.tools.length).toBe(10);
    expect(plan!.maxSteps).toBe(15);
  });

  it('getRole returns undefined for unknown id', () => {
    const mgr = new RoleManager();
    expect(mgr.getRole('nonexistent')).toBeUndefined();
  });

  it('createToolSet returns tools for hanimo role', () => {
    const mgr = new RoleManager();
    const hanimo = mgr.getRole('hanimo')!;
    const toolSet = mgr.createToolSet(hanimo);
    expect(toolSet).toBeDefined();
    expect(Object.keys(toolSet!).length).toBe(16);
  });

  it('createToolSet returns tools for dev role', () => {
    const mgr = new RoleManager();
    const dev = mgr.getRole('dev')!;
    const toolSet = mgr.createToolSet(dev);
    expect(toolSet).toBeDefined();
    const keys = Object.keys(toolSet!);
    expect(keys).toContain('read_file');
    expect(keys).toContain('write_file');
    expect(keys).toContain('shell_exec');
    expect(keys).toContain('hashline_read');
    expect(keys).toContain('hashline_edit');
    expect(keys).toContain('webfetch');
    expect(keys).toContain('todo');
    expect(keys).toContain('batch');
    expect(keys).toContain('diagnostics');
    expect(keys.length).toBe(16);
  });

  it('createToolSet returns read-only tools for plan role', () => {
    const mgr = new RoleManager();
    const plan = mgr.getRole('plan')!;
    const toolSet = mgr.createToolSet(plan);
    expect(toolSet).toBeDefined();
    const keys = Object.keys(toolSet!);
    expect(keys).toContain('read_file');
    expect(keys).toContain('hashline_read');
    expect(keys).toContain('glob_search');
    expect(keys).toContain('grep_search');
    expect(keys).toContain('webfetch');
    expect(keys).toContain('batch');
    expect(keys).toContain('diagnostics');
    expect(keys).not.toContain('write_file');
    expect(keys).not.toContain('edit_file');
    expect(keys).not.toContain('shell_exec');
    expect(keys.length).toBe(10);
  });

  it('buildRolePrompt appends role section', () => {
    const mgr = new RoleManager();
    const dev = mgr.getRole('dev')!;
    const base = 'You are a helpful assistant.';
    const result = mgr.buildRolePrompt(dev, base);
    expect(result).toContain(base);
    expect(result).toContain('## Active Role:');
    expect(result).toContain('🔧');
    expect(result).toContain('Dev');
  });

  it('buildRolePrompt returns base prompt if no systemPrompt', () => {
    const mgr = new RoleManager();
    const role: RoleDefinition = {
      id: 'empty',
      name: 'Empty',
      description: 'No prompt',
      icon: '⬜',
      systemPrompt: '',
      tools: [],
      maxSteps: 1,
    };
    const base = 'Base prompt here.';
    const result = mgr.buildRolePrompt(role, base);
    expect(result).toBe(base);
  });

  it('dev role tools list matches all 16 registered tools', () => {
    const mgr = new RoleManager();
    const dev = mgr.getRole('dev')!;
    const expected = [
      'read_file', 'write_file', 'edit_file',
      'hashline_read', 'hashline_edit',
      'git_status', 'git_diff', 'git_commit', 'git_log',
      'shell_exec', 'glob_search', 'grep_search',
      'webfetch', 'todo', 'batch', 'diagnostics',
    ];
    expect(dev.tools.sort()).toEqual(expected.sort());
  });

  it('plan role tools are a subset of dev tools', () => {
    const mgr = new RoleManager();
    const dev = mgr.getRole('dev')!;
    const plan = mgr.getRole('plan')!;
    for (const tool of plan.tools) {
      expect(dev.tools).toContain(tool);
    }
  });
});
