# Modol Level 2 Skill System Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Upgrade Modol's skill system from knowledge injection to workflow-aware Level 2 with triggers, steps, gates, and 8 built-in skills.

**Architecture:** Three new modules (`skill-parser.ts`, `skill-engine.ts`, `skill-prompt.ts`) handle parsing, state management, and prompt generation. `agent-loop.ts` gets ~15 lines for strict gate support. Built-in skills live in `src/skills/*.md` with YAML frontmatter.

**Tech Stack:** TypeScript, Vitest, YAML frontmatter parsing (manual — no dependency)

---

## File Structure

| Action | File | Responsibility |
|--------|------|---------------|
| Create | `src/core/skill-parser.ts` | Parse YAML frontmatter + extract steps from markdown |
| Create | `src/core/skill-engine.ts` | Trigger matching, active skill state, gate management |
| Create | `src/core/skill-prompt.ts` | Generate system prompt section for active skill |
| Modify | `src/core/skills.ts` | Integrate new modules, maintain backward compat |
| Modify | `src/core/types.ts:32-37` | Add `skill:gate` event type to `AgentEvent` union |
| Modify | `src/core/agent-loop.ts:82-215` | Add strict gate check after stream completion |
| Modify | `src/core/system-prompt.ts:65-66` | Use new skill prompt builder |
| Modify | `src/tui/hooks/use-commands.ts:570-588` | Extend `/skill` command with activate/deactivate/next/status |
| Modify | `src/tui/components/status-bar.tsx:78-157` | Show active skill + current step |
| Create | `src/skills/brainstorming.md` | Built-in: brainstorming workflow |
| Create | `src/skills/debugging.md` | Built-in: debugging workflow |
| Create | `src/skills/git-workflow.md` | Built-in: git commit workflow |
| Create | `src/skills/planning.md` | Built-in: planning workflow |
| Create | `src/skills/tdd.md` | Built-in: TDD workflow |
| Create | `src/skills/code-review.md` | Built-in: code review workflow |
| Create | `src/skills/refactoring.md` | Built-in: refactoring workflow |
| Create | `src/skills/security-check.md` | Built-in: security check workflow |
| Create | `tests/skill-parser.test.ts` | Parser unit tests |
| Create | `tests/skill-engine.test.ts` | Engine unit tests |
| Create | `tests/skill-prompt.test.ts` | Prompt generation tests |

---

### Task 1: Skill Parser Module

**Files:**
- Create: `src/core/skill-parser.ts`
- Create: `tests/skill-parser.test.ts`

- [ ] **Step 1: Write failing tests for frontmatter parsing**

```typescript
// tests/skill-parser.test.ts
import { describe, it, expect } from 'vitest';
import { parseSkillFile, parseFrontmatter } from '../src/core/skill-parser.js';
import type { ParsedSkill, SkillMeta } from '../src/core/skill-parser.js';

describe('skill-parser', () => {
  describe('parseFrontmatter', () => {
    it('parses valid YAML frontmatter', () => {
      const raw = `---
name: brainstorming
description: Design exploration
triggers: ["design", "brainstorm"]
roles: ["dev", "plan"]
type: workflow
steps:
  - id: explore
    title: "Explore context"
    gate: none
  - id: propose
    title: "Propose approach"
    gate: strict
---

## Step: explore
Explore the project context.

## Step: propose
Propose 2-3 approaches.`;

      const result = parseSkillFile(raw, '/tmp/brainstorming.md');
      expect(result.meta.name).toBe('brainstorming');
      expect(result.meta.type).toBe('workflow');
      expect(result.meta.triggers).toEqual(['design', 'brainstorm']);
      expect(result.meta.roles).toEqual(['dev', 'plan']);
      expect(result.meta.steps).toHaveLength(2);
      expect(result.meta.steps[0].id).toBe('explore');
      expect(result.meta.steps[0].gate).toBe('none');
      expect(result.meta.steps[1].gate).toBe('strict');
    });

    it('extracts step content from markdown sections', () => {
      const raw = `---
name: test-skill
description: Test
type: workflow
steps:
  - id: first
    title: "First step"
    gate: none
  - id: second
    title: "Second step"
    gate: soft
---

## Step: first
Do the first thing.
With multiple lines.

## Step: second
Do the second thing.`;

      const result = parseSkillFile(raw, '/tmp/test.md');
      expect(result.meta.steps[0].content).toBe('Do the first thing.\nWith multiple lines.');
      expect(result.meta.steps[1].content).toBe('Do the second thing.');
    });

    it('falls back to knowledge type when no frontmatter', () => {
      const raw = `# My Knowledge Skill

Some knowledge content here.`;

      const result = parseSkillFile(raw, '/tmp/knowledge.md');
      expect(result.meta.type).toBe('knowledge');
      expect(result.meta.name).toBe('knowledge');
      expect(result.meta.steps).toEqual([]);
      expect(result.rawContent).toContain('Some knowledge content here.');
    });

    it('falls back to knowledge type on malformed frontmatter', () => {
      const raw = `---
name: broken
type: [invalid
---

Content here.`;

      const result = parseSkillFile(raw, '/tmp/broken.md');
      expect(result.meta.type).toBe('knowledge');
    });

    it('defaults triggers and roles to empty arrays when omitted', () => {
      const raw = `---
name: minimal
description: Minimal skill
type: knowledge
---

Content.`;

      const result = parseSkillFile(raw, '/tmp/minimal.md');
      expect(result.meta.triggers).toEqual([]);
      expect(result.meta.roles).toEqual([]);
    });

    it('handles steps with missing markdown sections gracefully', () => {
      const raw = `---
name: missing-section
description: Test
type: workflow
steps:
  - id: exists
    title: "Exists"
    gate: none
  - id: missing
    title: "Missing"
    gate: none
---

## Step: exists
This step has content.`;

      const result = parseSkillFile(raw, '/tmp/missing.md');
      expect(result.meta.steps[0].content).toBe('This step has content.');
      expect(result.meta.steps[1].content).toBe('');
    });
  });
});
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `npx vitest run tests/skill-parser.test.ts`
Expected: FAIL — module `../src/core/skill-parser.js` does not exist

- [ ] **Step 3: Implement skill-parser.ts**

```typescript
// src/core/skill-parser.ts

export interface StepDef {
  id: string;
  title: string;
  gate: 'none' | 'soft' | 'strict';
  content: string;
}

export interface SkillMeta {
  name: string;
  description: string;
  triggers: string[];
  roles: string[];
  type: 'workflow' | 'knowledge';
  steps: StepDef[];
}

export interface ParsedSkill {
  meta: SkillMeta;
  rawContent: string;
  path: string;
}

/**
 * Minimal YAML frontmatter parser — handles the subset used by skill files.
 * No external dependency needed.
 */
function parseFrontmatterYaml(yaml: string): Record<string, unknown> {
  const result: Record<string, unknown> = {};
  const lines = yaml.split('\n');
  let currentKey = '';
  let inArray = false;
  let arrayItems: unknown[] = [];
  let inSteps = false;
  let steps: Record<string, unknown>[] = [];
  let currentStep: Record<string, unknown> = {};

  for (const line of lines) {
    // Top-level key: value
    const kvMatch = line.match(/^(\w[\w-]*):\s*(.*)$/);
    if (kvMatch) {
      // Flush previous array/steps
      if (inArray && currentKey) {
        result[currentKey] = arrayItems;
        inArray = false;
        arrayItems = [];
      }
      if (inSteps) {
        if (Object.keys(currentStep).length > 0) steps.push(currentStep);
        result['steps'] = steps;
        inSteps = false;
        steps = [];
        currentStep = {};
      }

      const [, key, value] = kvMatch;
      currentKey = key;

      if (value === '') {
        // Could be start of array or steps
        continue;
      }

      // Inline JSON array: ["a", "b"]
      const arrayMatch = value.match(/^\[(.+)\]$/);
      if (arrayMatch) {
        result[key] = arrayMatch[1]
          .split(',')
          .map(s => s.trim().replace(/^["']|["']$/g, ''));
        continue;
      }

      // Scalar value
      result[key] = value.replace(/^["']|["']$/g, '');
      continue;
    }

    // Array item:   - value
    const arrayItemMatch = line.match(/^\s+-\s+(?!id:)(.+)$/);
    if (arrayItemMatch && !inSteps) {
      inArray = true;
      arrayItems.push(arrayItemMatch[1].replace(/^["']|["']$/g, ''));
      continue;
    }

    // Steps array:   - id: xxx
    const stepStartMatch = line.match(/^\s+-\s+id:\s*(.+)$/);
    if (stepStartMatch) {
      inSteps = true;
      if (Object.keys(currentStep).length > 0) steps.push(currentStep);
      currentStep = { id: stepStartMatch[1].replace(/^["']|["']$/g, '') };
      continue;
    }

    // Step property:     title: xxx  or  gate: xxx
    const stepPropMatch = line.match(/^\s{4,}(\w+):\s*(.+)$/);
    if (stepPropMatch && inSteps) {
      currentStep[stepPropMatch[1]] = stepPropMatch[2].replace(/^["']|["']$/g, '');
      continue;
    }
  }

  // Flush remaining
  if (inArray && currentKey) result[currentKey] = arrayItems;
  if (inSteps) {
    if (Object.keys(currentStep).length > 0) steps.push(currentStep);
    result['steps'] = steps;
  }

  return result;
}

function extractStepContent(markdown: string, stepIds: string[]): Map<string, string> {
  const contents = new Map<string, string>();
  for (const id of stepIds) {
    const pattern = new RegExp(`^## Step:\\s*${id}\\s*$`, 'm');
    const match = pattern.exec(markdown);
    if (!match) {
      contents.set(id, '');
      continue;
    }

    const start = match.index + match[0].length;
    // Find next ## Step: or end of string
    const nextSection = markdown.slice(start).search(/^## Step:/m);
    const end = nextSection >= 0 ? start + nextSection : markdown.length;
    contents.set(id, markdown.slice(start, end).trim());
  }
  return contents;
}

export function parseSkillFile(raw: string, filePath: string): ParsedSkill {
  const filename = filePath.split('/').pop()?.replace('.md', '') ?? 'unknown';

  // Check for frontmatter
  const fmMatch = raw.match(/^---\n([\s\S]*?)\n---\n?([\s\S]*)$/);
  if (!fmMatch) {
    // No frontmatter — legacy knowledge skill
    const headingMatch = raw.match(/^#\s+(.+)$/m);
    return {
      meta: {
        name: filename,
        description: headingMatch?.[1] ?? filename,
        triggers: [],
        roles: [],
        type: 'knowledge',
        steps: [],
      },
      rawContent: raw.trim(),
      path: filePath,
    };
  }

  let parsed: Record<string, unknown>;
  try {
    parsed = parseFrontmatterYaml(fmMatch[1]);
  } catch {
    // Malformed YAML — fallback to knowledge
    return {
      meta: {
        name: filename,
        description: filename,
        triggers: [],
        roles: [],
        type: 'knowledge',
        steps: [],
      },
      rawContent: raw.trim(),
      path: filePath,
    };
  }

  const markdown = fmMatch[2];
  const type = parsed['type'] === 'workflow' ? 'workflow' : 'knowledge';

  const rawSteps = Array.isArray(parsed['steps']) ? parsed['steps'] as Record<string, unknown>[] : [];
  const stepIds = rawSteps.map(s => String(s['id'] ?? ''));
  const stepContents = type === 'workflow' ? extractStepContent(markdown, stepIds) : new Map<string, string>();

  const steps: StepDef[] = rawSteps.map(s => ({
    id: String(s['id'] ?? ''),
    title: String(s['title'] ?? ''),
    gate: (['none', 'soft', 'strict'].includes(String(s['gate'])) ? String(s['gate']) : 'none') as StepDef['gate'],
    content: stepContents.get(String(s['id'] ?? '')) ?? '',
  }));

  const triggers = Array.isArray(parsed['triggers'])
    ? (parsed['triggers'] as unknown[]).map(String)
    : [];
  const roles = Array.isArray(parsed['roles'])
    ? (parsed['roles'] as unknown[]).map(String)
    : [];

  return {
    meta: {
      name: String(parsed['name'] ?? filename),
      description: String(parsed['description'] ?? ''),
      triggers,
      roles,
      type,
      steps,
    },
    rawContent: markdown.trim(),
    path: filePath,
  };
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `npx vitest run tests/skill-parser.test.ts`
Expected: ALL PASS

- [ ] **Step 5: Commit**

```bash
git add src/core/skill-parser.ts tests/skill-parser.test.ts
git commit -m "feat(skills): add skill-parser module with frontmatter + step extraction"
```

---

### Task 2: Skill Engine Module

**Files:**
- Create: `src/core/skill-engine.ts`
- Create: `tests/skill-engine.test.ts`

- [ ] **Step 1: Write failing tests for skill engine**

```typescript
// tests/skill-engine.test.ts
import { describe, it, expect, beforeEach } from 'vitest';
import { SkillEngine } from '../src/core/skill-engine.js';
import type { ParsedSkill } from '../src/core/skill-parser.js';

function makeWorkflowSkill(overrides?: Partial<ParsedSkill['meta']>): ParsedSkill {
  return {
    meta: {
      name: 'test-workflow',
      description: 'Test workflow skill',
      triggers: ['test', '테스트'],
      roles: ['dev'],
      type: 'workflow',
      steps: [
        { id: 'step1', title: 'First', gate: 'none', content: 'Do first thing.' },
        { id: 'step2', title: 'Second', gate: 'strict', content: 'Do second thing.' },
        { id: 'step3', title: 'Third', gate: 'none', content: 'Do third thing.' },
      ],
      ...overrides,
    },
    rawContent: '',
    path: '/tmp/test.md',
  };
}

function makeKnowledgeSkill(): ParsedSkill {
  return {
    meta: {
      name: 'test-knowledge',
      description: 'Test knowledge skill',
      triggers: ['knowledge'],
      roles: [],
      type: 'knowledge',
      steps: [],
    },
    rawContent: 'Some knowledge content.',
    path: '/tmp/knowledge.md',
  };
}

describe('SkillEngine', () => {
  let engine: SkillEngine;

  beforeEach(() => {
    engine = new SkillEngine([makeWorkflowSkill(), makeKnowledgeSkill()]);
  });

  describe('matchTriggers', () => {
    it('matches trigger keyword case-insensitively', () => {
      const match = engine.matchTriggers('please Test this code', 'dev');
      expect(match?.meta.name).toBe('test-workflow');
    });

    it('matches Korean trigger', () => {
      const match = engine.matchTriggers('이거 테스트 해줘', 'dev');
      expect(match?.meta.name).toBe('test-workflow');
    });

    it('returns null when no trigger matches', () => {
      const match = engine.matchTriggers('hello world', 'dev');
      expect(match).toBeNull();
    });

    it('returns null when role is not allowed', () => {
      const match = engine.matchTriggers('test this', 'chat');
      expect(match).toBeNull();
    });

    it('matches all roles when roles array is empty', () => {
      const match = engine.matchTriggers('knowledge stuff', 'chat');
      expect(match?.meta.name).toBe('test-knowledge');
    });

    it('prefers more specific skill (shorter triggers array) on conflict', () => {
      const specific = makeWorkflowSkill({
        name: 'specific',
        triggers: ['test'],
        roles: ['dev'],
      });
      const broad = makeWorkflowSkill({
        name: 'broad',
        triggers: ['test', 'check', 'verify', 'validate'],
        roles: ['dev'],
      });
      const eng = new SkillEngine([broad, specific]);
      const match = eng.matchTriggers('run the test', 'dev');
      expect(match?.meta.name).toBe('specific');
    });
  });

  describe('activate / deactivate', () => {
    it('activates a skill by name', () => {
      const active = engine.activate('test-workflow');
      expect(active.skill.meta.name).toBe('test-workflow');
      expect(active.currentStepIndex).toBe(0);
      expect(active.completedSteps).toEqual([]);
      expect(active.gated).toBe(false);
    });

    it('throws on unknown skill name', () => {
      expect(() => engine.activate('nonexistent')).toThrow();
    });

    it('getActive returns active skill', () => {
      engine.activate('test-workflow');
      expect(engine.getActive()?.skill.meta.name).toBe('test-workflow');
    });

    it('deactivate clears active skill', () => {
      engine.activate('test-workflow');
      engine.deactivate();
      expect(engine.getActive()).toBeNull();
    });

    it('ignores triggers when a skill is already active', () => {
      engine.activate('test-workflow');
      const match = engine.matchTriggers('knowledge stuff', 'dev');
      // matchTriggers should still return the match but activation is blocked externally
      expect(match).not.toBeNull();
      // Engine has active skill
      expect(engine.getActive()?.skill.meta.name).toBe('test-workflow');
    });
  });

  describe('step progression', () => {
    it('advances to next step', () => {
      engine.activate('test-workflow');
      const next = engine.advanceStep();
      expect(next?.id).toBe('step2');
      expect(engine.getActive()?.currentStepIndex).toBe(1);
      expect(engine.getActive()?.completedSteps).toEqual(['step1']);
    });

    it('returns null when all steps complete', () => {
      engine.activate('test-workflow');
      engine.advanceStep(); // -> step2
      engine.approveGate(); // approve strict gate
      engine.advanceStep(); // -> step3
      const next = engine.advanceStep(); // -> done
      expect(next).toBeNull();
      expect(engine.getActive()).toBeNull(); // auto-deactivated
    });

    it('sets gated=true when advancing into a strict step', () => {
      engine.activate('test-workflow');
      engine.advanceStep(); // step1 -> step2 (strict)
      expect(engine.getActive()?.gated).toBe(true);
    });

    it('isGated returns true when at strict step', () => {
      engine.activate('test-workflow');
      engine.advanceStep();
      expect(engine.isGated()).toBe(true);
    });

    it('approveGate clears gated state', () => {
      engine.activate('test-workflow');
      engine.advanceStep();
      expect(engine.isGated()).toBe(true);
      engine.approveGate();
      expect(engine.isGated()).toBe(false);
    });

    it('advanceStep is blocked when gated', () => {
      engine.activate('test-workflow');
      engine.advanceStep(); // -> step2 (strict, gated)
      const next = engine.advanceStep(); // should not advance
      expect(next).toBeNull();
      expect(engine.getActive()?.currentStepIndex).toBe(1); // still at step2
    });
  });

  describe('getPromptData', () => {
    it('returns null when no active skill', () => {
      expect(engine.getPromptData()).toBeNull();
    });

    it('returns prompt data for active skill', () => {
      engine.activate('test-workflow');
      const data = engine.getPromptData();
      expect(data?.skill.meta.name).toBe('test-workflow');
      expect(data?.currentStepIndex).toBe(0);
      expect(data?.completedSteps).toEqual([]);
    });
  });

  describe('getAllSkills', () => {
    it('returns all loaded skills', () => {
      const all = engine.getAllSkills();
      expect(all).toHaveLength(2);
    });
  });
});
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `npx vitest run tests/skill-engine.test.ts`
Expected: FAIL — module `../src/core/skill-engine.js` does not exist

- [ ] **Step 3: Implement skill-engine.ts**

```typescript
// src/core/skill-engine.ts
import type { ParsedSkill, StepDef } from './skill-parser.js';

export interface ActiveSkill {
  skill: ParsedSkill;
  currentStepIndex: number;
  completedSteps: string[];
  startedAt: number;
  gated: boolean;
}

export class SkillEngine {
  private skills: ParsedSkill[];
  private active: ActiveSkill | null = null;

  constructor(skills: ParsedSkill[]) {
    this.skills = skills;
  }

  getAllSkills(): ParsedSkill[] {
    return [...this.skills];
  }

  /**
   * Match user message against trigger keywords.
   * Returns the most specific matching skill (shortest triggers array).
   * Respects role filtering.
   */
  matchTriggers(message: string, role: string): ParsedSkill | null {
    const lower = message.toLowerCase();
    const matches: ParsedSkill[] = [];

    for (const skill of this.skills) {
      if (skill.meta.triggers.length === 0) continue;
      if (skill.meta.roles.length > 0 && !skill.meta.roles.includes(role)) continue;

      const triggered = skill.meta.triggers.some(t => lower.includes(t.toLowerCase()));
      if (triggered) matches.push(skill);
    }

    if (matches.length === 0) return null;

    // Most specific = shortest triggers array
    matches.sort((a, b) => a.meta.triggers.length - b.meta.triggers.length);
    return matches[0];
  }

  activate(skillName: string): ActiveSkill {
    const skill = this.skills.find(s => s.meta.name === skillName);
    if (!skill) throw new Error(`Skill not found: ${skillName}`);

    this.active = {
      skill,
      currentStepIndex: 0,
      completedSteps: [],
      startedAt: Date.now(),
      gated: false,
    };
    return this.active;
  }

  deactivate(): void {
    this.active = null;
  }

  getActive(): ActiveSkill | null {
    return this.active ? { ...this.active } : null;
  }

  /**
   * Advance to the next step. Returns the new step, or null if:
   * - No active skill
   * - Currently gated (must approveGate first)
   * - All steps completed (auto-deactivates)
   */
  advanceStep(): StepDef | null {
    if (!this.active) return null;
    if (this.active.gated) return null;

    const { steps } = this.active.skill.meta;
    const currentStep = steps[this.active.currentStepIndex];

    // Mark current step as completed
    if (currentStep) {
      this.active.completedSteps.push(currentStep.id);
    }

    const nextIndex = this.active.currentStepIndex + 1;

    if (nextIndex >= steps.length) {
      // All steps done — auto-deactivate
      this.active = null;
      return null;
    }

    this.active.currentStepIndex = nextIndex;
    const nextStep = steps[nextIndex];

    // Check if next step is a strict gate
    if (nextStep.gate === 'strict') {
      this.active.gated = true;
    }

    return nextStep;
  }

  isGated(): boolean {
    return this.active?.gated ?? false;
  }

  approveGate(): void {
    if (this.active) {
      this.active.gated = false;
    }
  }

  getPromptData(): { skill: ParsedSkill; currentStepIndex: number; completedSteps: string[] } | null {
    if (!this.active) return null;
    return {
      skill: this.active.skill,
      currentStepIndex: this.active.currentStepIndex,
      completedSteps: [...this.active.completedSteps],
    };
  }
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `npx vitest run tests/skill-engine.test.ts`
Expected: ALL PASS

- [ ] **Step 5: Commit**

```bash
git add src/core/skill-engine.ts tests/skill-engine.test.ts
git commit -m "feat(skills): add skill-engine with trigger matching, step progression, and gates"
```

---

### Task 3: Skill Prompt Builder Module

**Files:**
- Create: `src/core/skill-prompt.ts`
- Create: `tests/skill-prompt.test.ts`

- [ ] **Step 1: Write failing tests for prompt generation**

```typescript
// tests/skill-prompt.test.ts
import { describe, it, expect } from 'vitest';
import { buildSkillPrompt } from '../src/core/skill-prompt.js';
import type { ParsedSkill } from '../src/core/skill-parser.js';

function makeSkill(): ParsedSkill {
  return {
    meta: {
      name: 'brainstorming',
      description: 'Design exploration',
      triggers: [],
      roles: [],
      type: 'workflow',
      steps: [
        { id: 'explore', title: 'Explore context', gate: 'none', content: 'Read project files.' },
        { id: 'clarify', title: 'Clarify requirements', gate: 'none', content: 'Ask questions one at a time.' },
        { id: 'propose', title: 'Propose approach', gate: 'strict', content: 'Suggest 2-3 approaches.' },
        { id: 'design', title: 'Detailed design', gate: 'soft', content: 'Write the design.' },
      ],
    },
    rawContent: '',
    path: '/tmp/brainstorming.md',
  };
}

describe('skill-prompt', () => {
  it('generates workflow prompt with current step highlighted', () => {
    const result = buildSkillPrompt({
      skill: makeSkill(),
      currentStepIndex: 1,
      completedSteps: ['explore'],
    });

    expect(result).toContain('## Active Skill: brainstorming');
    expect(result).toContain('Step 2/4: Clarify requirements');
    expect(result).toContain('Ask questions one at a time.');
    expect(result).toContain('[x] explore');
    expect(result).toContain('[>] clarify');
    expect(result).toContain('[ ] propose');
    expect(result).toContain('[GATE]');
  });

  it('includes STEP_COMPLETE signal instruction', () => {
    const result = buildSkillPrompt({
      skill: makeSkill(),
      currentStepIndex: 0,
      completedSteps: [],
    });
    expect(result).toContain('[STEP_COMPLETE: explore]');
  });

  it('generates knowledge prompt with raw content', () => {
    const knowledgeSkill: ParsedSkill = {
      meta: {
        name: 'typescript-tips',
        description: 'TS knowledge',
        triggers: [],
        roles: [],
        type: 'knowledge',
        steps: [],
      },
      rawContent: 'Use strict mode.\nPrefer interfaces.',
      path: '/tmp/ts.md',
    };

    const result = buildSkillPrompt({
      skill: knowledgeSkill,
      currentStepIndex: 0,
      completedSteps: [],
    });

    expect(result).toContain('## Loaded Skill: typescript-tips');
    expect(result).toContain('Use strict mode.');
    expect(result).not.toContain('Workflow Progress');
  });

  it('returns empty string for null input', () => {
    // The function should handle being called with data for a null skill gracefully
    // In practice, the caller checks for null before calling, but test the prompt content
    const result = buildSkillPrompt({
      skill: makeSkill(),
      currentStepIndex: 3,
      completedSteps: ['explore', 'clarify', 'propose'],
    });
    expect(result).toContain('Step 4/4: Detailed design');
  });
});
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `npx vitest run tests/skill-prompt.test.ts`
Expected: FAIL — module does not exist

- [ ] **Step 3: Implement skill-prompt.ts**

```typescript
// src/core/skill-prompt.ts
import type { ParsedSkill } from './skill-parser.js';

interface PromptData {
  skill: ParsedSkill;
  currentStepIndex: number;
  completedSteps: string[];
}

export function buildSkillPrompt(data: PromptData): string {
  const { skill, currentStepIndex, completedSteps } = data;

  if (skill.meta.type === 'knowledge') {
    return `\n\n## Loaded Skill: ${skill.meta.name}\n${skill.rawContent}`;
  }

  // Workflow prompt
  const steps = skill.meta.steps;
  const currentStep = steps[currentStepIndex];
  if (!currentStep) return '';

  const stepNum = currentStepIndex + 1;
  const totalSteps = steps.length;

  // Progress markers
  const progress = steps.map((s, i) => {
    const gateTag = s.gate === 'strict' ? ' [GATE: 사용자 승인 필요]' : '';
    if (completedSteps.includes(s.id)) {
      return `- [x] ${s.id} — ${s.title}`;
    }
    if (i === currentStepIndex) {
      return `- [>] ${s.id} — ${s.title} (현재)${gateTag}`;
    }
    return `- [ ] ${s.id} — ${s.title}${gateTag}`;
  }).join('\n');

  return `

## Active Skill: ${skill.meta.name} (Step ${stepNum}/${totalSteps}: ${currentStep.title})

### Current Step Instructions
${currentStep.content}

### Workflow Progress
${progress}

### Rules
- 현재 단계의 지시사항을 따르세요.
- 단계가 완료되면 "[STEP_COMPLETE: ${currentStep.id}]"를 응답 끝에 포함하세요.
- [GATE] 표시된 단계 전에는 반드시 사용자 확인을 받으세요.
- 사용자가 명시적으로 요청하지 않는 한 단계를 건너뛰지 마세요.`;
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `npx vitest run tests/skill-prompt.test.ts`
Expected: ALL PASS

- [ ] **Step 5: Commit**

```bash
git add src/core/skill-prompt.ts tests/skill-prompt.test.ts
git commit -m "feat(skills): add skill-prompt builder for workflow and knowledge skills"
```

---

### Task 4: Integrate into Core (skills.ts, system-prompt.ts, types.ts)

**Files:**
- Modify: `src/core/skills.ts`
- Modify: `src/core/system-prompt.ts:5,65-66`
- Modify: `src/core/types.ts:32-37`

- [ ] **Step 1: Update types.ts — add skill:gate event**

In `src/core/types.ts`, add `skill:gate` to the `AgentEvent` union type at line 37:

```typescript
// Add after the existing 'done' event in the union:
  | { type: 'skill:gate'; skillName: string; stepId: string; stepTitle: string };
```

The full `AgentEvent` type becomes:

```typescript
export type AgentEvent =
  | { type: 'token'; content: string }
  | { type: 'tool-call'; toolName: string; args: Record<string, unknown> }
  | { type: 'tool-result'; toolName: string; result: string; isError: boolean }
  | { type: 'error'; error: Error }
  | { type: 'done'; response: string; usage: TokenUsage }
  | { type: 'skill:gate'; skillName: string; stepId: string; stepTitle: string };
```

- [ ] **Step 2: Update skills.ts — add loadAllSkills integrating both built-in and user skills**

Add to `src/core/skills.ts`:

```typescript
import { readFileSync, readdirSync, mkdirSync } from 'node:fs';
import { join, dirname } from 'node:path';
import { homedir } from 'node:os';
import { fileURLToPath } from 'node:url';
import { parseSkillFile } from './skill-parser.js';
import type { ParsedSkill } from './skill-parser.js';

const SKILLS_DIR = join(homedir(), '.modol', 'skills');

// Built-in skills directory (relative to this file's location in dist)
function getBuiltinSkillsDir(): string {
  const thisDir = dirname(fileURLToPath(import.meta.url));
  // In dist: dist/core/skills.js -> dist/skills/
  // In src: src/core/skills.ts -> src/skills/
  return join(thisDir, '..', 'skills');
}

// --- Legacy interface (backward compat) ---
export interface Skill {
  id: string;
  name: string;
  content: string;
  path: string;
}

export function loadSkills(): Skill[] {
  mkdirSync(SKILLS_DIR, { recursive: true });
  let files: string[];
  try {
    files = readdirSync(SKILLS_DIR).filter(f => f.endsWith('.md'));
  } catch { return []; }

  return files.map(f => {
    const path = join(SKILLS_DIR, f);
    const content = readFileSync(path, 'utf-8');
    const id = f.replace('.md', '');
    const headingMatch = content.match(/^#\s+(.+)$/m);
    const name = headingMatch?.[1] ?? id;
    return { id, name, content: content.trim(), path };
  }).filter(s => s.content.length > 0);
}

export function buildSkillsPrompt(skills: Skill[]): string {
  if (skills.length === 0) return '';
  const sections = skills.map(s => `### Skill: ${s.name}\n${s.content}`);
  return `\n\n## Loaded Skills (${skills.length})\n${sections.join('\n\n')}`;
}

export function getSkillsDir(): string {
  return SKILLS_DIR;
}

// --- New Level 2 interface ---
function loadMdFiles(dir: string): ParsedSkill[] {
  try {
    return readdirSync(dir)
      .filter(f => f.endsWith('.md'))
      .map(f => {
        const path = join(dir, f);
        const raw = readFileSync(path, 'utf-8');
        return parseSkillFile(raw, path);
      })
      .filter(s => s.rawContent.length > 0 || s.meta.steps.length > 0);
  } catch {
    return [];
  }
}

/**
 * Load all skills: built-in + user custom.
 * User skills override built-in skills with the same name.
 */
export function loadAllSkills(): ParsedSkill[] {
  mkdirSync(SKILLS_DIR, { recursive: true });

  const builtins = loadMdFiles(getBuiltinSkillsDir());
  const customs = loadMdFiles(SKILLS_DIR);

  // User overrides built-in by name
  const customNames = new Set(customs.map(s => s.meta.name));
  const merged = [
    ...builtins.filter(s => !customNames.has(s.meta.name)),
    ...customs,
  ];

  return merged;
}
```

- [ ] **Step 3: Update system-prompt.ts — use new skill engine when available**

In `src/core/system-prompt.ts`, replace lines 5 and 65-66:

Line 5 — change import:
```typescript
import { loadSkills, buildSkillsPrompt } from './skills.js';
```
No change needed here — the legacy path still works. The new skill engine prompt is injected by the caller (agent-loop), not by `buildSystemPrompt()`.

The `buildSystemPrompt` function continues to use legacy `loadSkills()` + `buildSkillsPrompt()` for knowledge-type skills. Workflow skill prompts are appended by the agent-loop when a skill is active.

No modification needed — the architecture keeps system-prompt.ts unchanged.

- [ ] **Step 4: Run existing tests to verify no regressions**

Run: `npx vitest run tests/skills.test.ts`
Expected: ALL PASS (backward compat maintained)

- [ ] **Step 5: Commit**

```bash
git add src/core/skills.ts src/core/types.ts
git commit -m "feat(skills): integrate Level 2 loading into skills.ts, add skill:gate event type"
```

---

### Task 5: Agent Loop Integration

**Files:**
- Modify: `src/core/agent-loop.ts:82-215`

- [ ] **Step 1: Add SkillEngine to AgentLoopOptions**

In `src/core/types.ts`, add optional `skillEngine` to `AgentLoopOptions`:

```typescript
// Add import at top of types.ts
import type { SkillEngine } from './skill-engine.js';

// Add to AgentLoopOptions interface (after abortSignal):
  skillEngine?: SkillEngine;
```

- [ ] **Step 2: Add skill gate handling to agent-loop.ts**

In `src/core/agent-loop.ts`, after the existing `import` block (line 8), add:

```typescript
import { buildSkillPrompt } from './skill-prompt.js';
```

In the `runAgentLoop` function, after destructuring options (line 93), add `skillEngine`:

```typescript
  const {
    model,
    systemPrompt,
    messages,
    tools,
    maxSteps = 25,
    onEvent,
    abortSignal,
    skillEngine,
  } = options;
```

Before the `try` block at line 124, inject skill prompt into system prompt:

```typescript
  // Inject active skill prompt
  const promptData = skillEngine?.getPromptData() ?? null;
  const skillSection = promptData ? buildSkillPrompt(promptData) : '';
  const fullSystemPrompt = systemPrompt + skillSection;
```

Replace `system: systemPrompt` with `system: fullSystemPrompt` in the `streamText` call at line 128.

After `emit({ type: 'done', ... })` at line 200, add step completion detection:

```typescript
    // Detect skill step completion signal in response
    if (skillEngine) {
      const stepMatch = fullResponse.match(/\[STEP_COMPLETE:\s*(\w[\w-]*)\]/);
      if (stepMatch) {
        skillEngine.advanceStep();
      }

      // Check if now gated
      if (skillEngine.isGated()) {
        const gateData = skillEngine.getPromptData();
        if (gateData) {
          const step = gateData.skill.meta.steps[gateData.currentStepIndex];
          emit({
            type: 'skill:gate',
            skillName: gateData.skill.meta.name,
            stepId: step.id,
            stepTitle: step.title,
          });
        }
      }
    }
```

- [ ] **Step 3: Update AgentLoopResult to include gated state**

In `src/core/types.ts`, add `gated` to `AgentLoopResult`:

```typescript
export interface AgentLoopResult {
  response: string;
  usage: TokenUsage;
  messages: Message[];
  gated?: boolean;
}
```

In `agent-loop.ts`, in the return statement at line 202, add:

```typescript
    return {
      response: fullResponse,
      usage: totalUsage,
      messages: updatedMessages,
      gated: skillEngine?.isGated() ?? false,
    };
```

- [ ] **Step 4: Run full test suite to verify no regressions**

Run: `npx vitest run`
Expected: ALL PASS

- [ ] **Step 5: Commit**

```bash
git add src/core/agent-loop.ts src/core/types.ts
git commit -m "feat(skills): integrate skill engine into agent-loop with strict gate support"
```

---

### Task 6: TUI Integration (Status Bar + Commands)

**Files:**
- Modify: `src/tui/components/status-bar.tsx:78-157`
- Modify: `src/tui/hooks/use-commands.ts:570-588`
- Modify: `src/tui/app.tsx` (pass skill engine to agent hook, handle gate events)

- [ ] **Step 1: Add skill display to status bar**

In `src/tui/components/status-bar.tsx`, add a new prop and display element.

Add to `StatusBarProps` interface (after `verbose`):

```typescript
  activeSkill?: { name: string; stepTitle: string; stepNum: number; totalSteps: number; gated: boolean };
```

Add to the `StatusBar` component JSX, after the `subAgentStatus` section (after line 135):

```typescript
          {props.activeSkill ? (
            <>
              <Text color={colors.dimText}> {'\u2502'} </Text>
              <Text color={props.activeSkill.gated ? colors.warning : colors.statusThinking}>
                {props.activeSkill.gated ? '\u26A0' : '\u2699'} {props.activeSkill.name} ({props.activeSkill.stepNum}/{props.activeSkill.totalSteps}: {props.activeSkill.stepTitle})
              </Text>
            </>
          ) : null}
```

- [ ] **Step 2: Extend /skill command in use-commands.ts**

Replace the existing `skill` command handler (lines 570-588) with:

```typescript
  skill: (args, ctx) => {
    const parts = args.trim().split(/\s+/);
    const sub = parts[0] || 'list';

    // Import skill engine from context (passed via CommandContext)
    const skillEngine = (ctx as unknown as { skillEngine?: import('../../core/skill-engine.js').SkillEngine }).skillEngine;

    if (sub === 'list') {
      if (!skillEngine) {
        ctx.addSystemMessage('Skill engine not initialized.');
        return;
      }
      const skills = skillEngine.getAllSkills();
      if (skills.length === 0) {
        ctx.addSystemMessage('No skills loaded.');
        return;
      }
      const lines = skills.map(s => {
        const typeTag = s.meta.type === 'workflow' ? '[workflow]' : '[knowledge]';
        const triggers = s.meta.triggers.length > 0 ? ` triggers: ${s.meta.triggers.join(', ')}` : '';
        return `  ${s.meta.name.padEnd(20)} ${typeTag}${triggers}`;
      });
      ctx.addSystemMessage(['Skills:', '', ...lines].join('\n'));
      return;
    }

    if (sub === 'activate') {
      const name = parts[1];
      if (!name || !skillEngine) {
        ctx.addSystemMessage('Usage: /skill activate <name>');
        return;
      }
      try {
        skillEngine.activate(name);
        ctx.addSystemMessage(`Skill activated: ${name}`);
      } catch (e) {
        ctx.addSystemMessage(`Error: ${e instanceof Error ? e.message : String(e)}`);
      }
      return;
    }

    if (sub === 'deactivate') {
      if (!skillEngine) return;
      skillEngine.deactivate();
      ctx.addSystemMessage('Skill deactivated.');
      return;
    }

    if (sub === 'next') {
      if (!skillEngine) return;
      if (skillEngine.isGated()) {
        skillEngine.approveGate();
        ctx.addSystemMessage('Gate approved. Continuing workflow.');
      }
      const next = skillEngine.advanceStep();
      if (next) {
        ctx.addSystemMessage(`Advanced to: ${next.title}`);
      } else {
        ctx.addSystemMessage('Skill workflow complete.');
      }
      return;
    }

    if (sub === 'status') {
      if (!skillEngine) return;
      const data = skillEngine.getPromptData();
      if (!data) {
        ctx.addSystemMessage('No active skill.');
        return;
      }
      const step = data.skill.meta.steps[data.currentStepIndex];
      const gated = skillEngine.isGated() ? ' [GATED - /skill next to approve]' : '';
      ctx.addSystemMessage(`Active: ${data.skill.meta.name} — Step ${data.currentStepIndex + 1}/${data.skill.meta.steps.length}: ${step.title}${gated}`);
      return;
    }

    if (sub === 'info') {
      const name = parts[1];
      if (!name || !skillEngine) {
        ctx.addSystemMessage('Usage: /skill info <name>');
        return;
      }
      const skill = skillEngine.getAllSkills().find(s => s.meta.name === name);
      if (!skill) {
        ctx.addSystemMessage(`Skill not found: ${name}`);
        return;
      }
      const lines = [
        `Name: ${skill.meta.name}`,
        `Type: ${skill.meta.type}`,
        `Description: ${skill.meta.description}`,
        `Triggers: ${skill.meta.triggers.join(', ') || 'none'}`,
        `Roles: ${skill.meta.roles.join(', ') || 'all'}`,
        `Path: ${skill.path}`,
      ];
      if (skill.meta.steps.length > 0) {
        lines.push('', 'Steps:');
        skill.meta.steps.forEach((s, i) => {
          const gate = s.gate !== 'none' ? ` [${s.gate}]` : '';
          lines.push(`  ${i + 1}. ${s.title}${gate}`);
        });
      }
      ctx.addSystemMessage(lines.join('\n'));
      return;
    }

    ctx.addSystemMessage('Usage: /skill [list|activate|deactivate|next|status|info]');
  },
```

- [ ] **Step 3: Add skillEngine to CommandContext interface**

In `src/tui/hooks/use-commands.ts`, add to the `CommandContext` interface:

```typescript
  skillEngine?: import('../../core/skill-engine.js').SkillEngine;
```

- [ ] **Step 4: Run test suite**

Run: `npx vitest run`
Expected: ALL PASS

- [ ] **Step 5: Commit**

```bash
git add src/tui/components/status-bar.tsx src/tui/hooks/use-commands.ts
git commit -m "feat(skills): add skill status bar display and /skill subcommands in TUI"
```

---

### Task 7: Wire Everything in app.tsx

**Files:**
- Modify: `src/tui/app.tsx`

- [ ] **Step 1: Initialize SkillEngine in App component**

At the top of `app.tsx`, add imports:

```typescript
import { loadAllSkills } from '../core/skills.js';
import { SkillEngine } from '../core/skill-engine.js';
```

Inside the `App` component, add state initialization (after other `useRef`/`useState` calls):

```typescript
const [skillEngine] = useState(() => new SkillEngine(loadAllSkills()));
const [activeSkillDisplay, setActiveSkillDisplay] = useState<{
  name: string; stepTitle: string; stepNum: number; totalSteps: number; gated: boolean;
} | undefined>(undefined);
```

- [ ] **Step 2: Pass skillEngine to agent loop calls**

Where `runAgentLoop` is called (in `use-agent.ts` hook or directly in app.tsx), pass `skillEngine` in options.

Before calling the agent loop, add trigger matching:

```typescript
// Before sending message to agent loop:
if (!skillEngine.getActive()) {
  const match = skillEngine.matchTriggers(userMessage, currentRole);
  if (match) {
    skillEngine.activate(match.meta.name);
    // Show notification
    addSystemMessage(`\uD83D\uDD27 [${match.meta.name}] 스킬 활성화`);
  }
}
```

- [ ] **Step 3: Handle skill:gate events**

In the `onEvent` callback, add:

```typescript
if (event.type === 'skill:gate') {
  setActiveSkillDisplay({
    name: event.skillName,
    stepTitle: event.stepTitle,
    stepNum: /* current step index + 1 */,
    totalSteps: /* total steps */,
    gated: true,
  });
  addSystemMessage(`[GATE] ${event.skillName} → ${event.stepTitle}. /skill next 로 승인하세요.`);
}
```

Update `activeSkillDisplay` on each agent loop completion:

```typescript
// After agent loop completes:
const data = skillEngine.getPromptData();
if (data) {
  const step = data.skill.meta.steps[data.currentStepIndex];
  setActiveSkillDisplay({
    name: data.skill.meta.name,
    stepTitle: step.title,
    stepNum: data.currentStepIndex + 1,
    totalSteps: data.skill.meta.steps.length,
    gated: skillEngine.isGated(),
  });
} else {
  setActiveSkillDisplay(undefined);
}
```

- [ ] **Step 4: Pass activeSkillDisplay to StatusBar**

```typescript
<StatusBar
  /* ...existing props... */
  activeSkill={activeSkillDisplay}
/>
```

- [ ] **Step 5: Pass skillEngine to CommandContext**

In the command context object, add:

```typescript
skillEngine,
```

- [ ] **Step 6: Run app manually to verify TUI renders**

Run: `npx tsx src/cli.ts --tui`
Expected: App starts without errors. `/skill list` shows loaded skills.

- [ ] **Step 7: Commit**

```bash
git add src/tui/app.tsx
git commit -m "feat(skills): wire skill engine into TUI app with trigger detection and gate UI"
```

---

### Task 8: Built-in Skills — 1차 (brainstorming, debugging, git-workflow)

**Files:**
- Create: `src/skills/brainstorming.md`
- Create: `src/skills/debugging.md`
- Create: `src/skills/git-workflow.md`

- [ ] **Step 1: Create brainstorming.md**

```markdown
---
name: brainstorming
description: 구현 전 설계/요구사항 탐색
triggers: ["설계해", "brainstorm", "디자인해", "기획해", "어떻게 만들"]
roles: ["dev", "plan"]
type: workflow
steps:
  - id: explore
    title: "프로젝트 컨텍스트 파악"
    gate: none
  - id: clarify
    title: "요구사항 명확화"
    gate: none
  - id: propose
    title: "접근 방식 제안"
    gate: strict
  - id: design
    title: "상세 설계"
    gate: soft
---

## Step: explore
현재 프로젝트의 관련 파일, 최근 커밋, 문서를 먼저 읽어 컨텍스트를 파악하세요.

도구를 사용하여:
- 관련 소스 파일 구조 확인 (glob_search)
- 최근 변경사항 확인 (git_log)
- README나 문서 파일 확인 (read_file)

이 단계에서는 코드를 수정하지 마세요.

[STEP_COMPLETE: explore]를 출력하여 다음 단계로 진행하세요.

## Step: clarify
사용자에게 한 번에 하나씩 질문하여 요구사항을 명확히 하세요.

규칙:
- 한 메시지에 질문 하나만
- 가능하면 객관식(A/B/C)으로 제시
- 목적, 제약 조건, 성공 기준에 집중
- 사용자가 "충분해" 또는 구체적 방향을 제시하면 다음 단계로

[STEP_COMPLETE: clarify]를 출력하여 다음 단계로 진행하세요.

## Step: propose
2-3가지 접근 방식을 트레이드오프와 함께 제안하세요.

형식:
1. **추천안**: [방식 이름] — [장점] / [단점]
2. **대안 A**: [방식 이름] — [장점] / [단점]
3. **대안 B**: [방식 이름] — [장점] / [단점]

추천안을 먼저 제시하고 이유를 설명하세요.
사용자가 선택할 때까지 구현을 시작하지 마세요.

[STEP_COMPLETE: propose]를 출력하여 다음 단계로 진행하세요.

## Step: design
사용자가 승인한 방식으로 상세 설계를 작성하세요.

포함할 내용:
- 아키텍처 개요
- 핵심 컴포넌트/모듈 구조
- 데이터 흐름
- 에러 처리 전략
- 테스트 전략

설계가 완료되면 사용자에게 확인을 요청하세요.

[STEP_COMPLETE: design]를 출력하여 스킬을 완료하세요.
```

- [ ] **Step 2: Create debugging.md**

```markdown
---
name: debugging
description: 체계적 버그 분석 및 수정
triggers: ["버그", "에러", "debug", "오류", "안돼", "안됨", "왜 안"]
roles: ["dev"]
type: workflow
steps:
  - id: reproduce
    title: "재현 확인"
    gate: none
  - id: hypothesize
    title: "원인 가설 수립"
    gate: none
  - id: isolate
    title: "원인 격리"
    gate: none
  - id: fix
    title: "수정 적용"
    gate: strict
  - id: verify
    title: "수정 검증"
    gate: none
---

## Step: reproduce
버그를 정확히 재현하세요.

확인할 것:
- 에러 메시지 전문, 스택 트레이스
- 재현 조건 (입력값, 환경, 순서)
- 재현 불가능하면 로그 수집 방법 제안

도구 사용:
- shell_exec: 에러 재현 명령 실행
- read_file: 관련 로그 파일 확인
- grep_search: 에러 메시지로 코드 검색

[STEP_COMPLETE: reproduce]를 출력하여 다음 단계로 진행하세요.

## Step: hypothesize
가능한 원인 2-3개를 가설로 세우세요.

각 가설에 대해:
1. **가설**: [원인 설명]
2. **근거**: [왜 이것이 원인일 수 있는지]
3. **검증 방법**: [어떻게 확인할 수 있는지]

가장 가능성 높은 가설부터 나열하세요.

[STEP_COMPLETE: hypothesize]를 출력하여 다음 단계로 진행하세요.

## Step: isolate
가설을 하나씩 검증하여 근본 원인을 특정하세요.

도구 사용:
- grep_search: 관련 코드 패턴 검색
- read_file: 의심 코드 읽기
- shell_exec: git blame으로 변경 이력 확인
- shell_exec: 테스트 실행으로 가설 검증

가설이 틀리면 다음 가설로. 모두 틀리면 새 가설 수립.
근본 원인을 찾으면 명확히 설명하세요.

[STEP_COMPLETE: isolate]를 출력하여 다음 단계로 진행하세요.

## Step: fix
수정안을 사용자에게 제시하세요.

포함할 내용:
- 근본 원인 요약 (1-2줄)
- 수정 방법 설명
- 영향 범위 (어떤 파일, 어떤 기능)
- 수정 코드 미리보기 (diff 형태)

사용자가 승인하면 코드를 수정하세요.
승인 없이 코드를 변경하지 마세요.

[STEP_COMPLETE: fix]를 출력하여 다음 단계로 진행하세요.

## Step: verify
수정이 올바른지 검증하세요.

확인할 것:
- 원래 버그가 해결되었는지 (재현 시도)
- 기존 테스트 통과 (shell_exec: npm test 등)
- 회귀 없음 (관련 기능 확인)
- 영향 범위 내 다른 코드 문제 없음

[STEP_COMPLETE: verify]를 출력하여 스킬을 완료하세요.
```

- [ ] **Step 3: Create git-workflow.md**

```markdown
---
name: git-workflow
description: 안전한 Git 커밋/브랜치 워크플로우
triggers: ["커밋", "commit", "브랜치", "branch", "PR", "풀리퀘"]
roles: ["dev"]
type: workflow
steps:
  - id: status
    title: "변경사항 확인"
    gate: none
  - id: stage
    title: "스테이징 선별"
    gate: soft
  - id: message
    title: "커밋 메시지 작성"
    gate: soft
  - id: commit
    title: "커밋 실행"
    gate: strict
---

## Step: status
현재 변경사항을 전체 파악하세요.

실행:
- git_status: 변경/추가/삭제된 파일 목록
- git_diff: 실제 변경 내용 확인

주의:
- .env, credentials.json, *.key 등 민감 파일이 있으면 경고
- 바이너리 파일이 포함되어 있으면 알림
- 변경이 너무 많으면 논리적 단위로 분리 제안

[STEP_COMPLETE: status]를 출력하여 다음 단계로 진행하세요.

## Step: stage
변경사항을 논리적 단위로 스테이징하세요.

규칙:
- `git add -A` 대신 파일별로 선택
- 관련 변경끼리 묶기 (기능 단위)
- 민감 파일은 절대 스테이징하지 않음
- 스테이징할 파일 목록을 사용자에게 제시

[STEP_COMPLETE: stage]를 출력하여 다음 단계로 진행하세요.

## Step: message
Conventional Commit 형식으로 커밋 메시지를 작성하세요.

형식: `type(scope): description`

type: feat, fix, refactor, docs, test, chore, style, perf
scope: 변경 범위 (선택)
description: "why" 중심의 간결한 설명

예시:
- feat(auth): add session refresh on token expiry
- fix(cli): prevent crash when config file is missing

메시지를 사용자에게 제안하고 수정 요청을 받으세요.

[STEP_COMPLETE: message]를 출력하여 다음 단계로 진행하세요.

## Step: commit
커밋을 실행하기 전에 최종 확인을 보여주세요.

표시할 내용:
- 스테이지된 파일 목록
- 커밋 메시지 전문
- 브랜치 이름

사용자가 승인하면 git_commit 도구로 커밋하세요.
승인 없이 커밋하지 마세요.

[STEP_COMPLETE: commit]를 출력하여 스킬을 완료하세요.
```

- [ ] **Step 4: Verify built-in skills parse correctly**

Run: `npx tsx -e "import { loadAllSkills } from './src/core/skills.js'; const s = loadAllSkills(); console.log(s.map(x => x.meta.name + ' [' + x.meta.type + '] steps:' + x.meta.steps.length).join('\n'));"`
Expected: All 3 skills listed with correct step counts.

- [ ] **Step 5: Commit**

```bash
git add src/skills/brainstorming.md src/skills/debugging.md src/skills/git-workflow.md
git commit -m "feat(skills): add 1st batch built-in skills (brainstorming, debugging, git-workflow)"
```

---

### Task 9: Built-in Skills — 2차 (planning, tdd, code-review, refactoring, security-check)

**Files:**
- Create: `src/skills/planning.md`
- Create: `src/skills/tdd.md`
- Create: `src/skills/code-review.md`
- Create: `src/skills/refactoring.md`
- Create: `src/skills/security-check.md`

- [ ] **Step 1: Create planning.md**

```markdown
---
name: planning
description: 멀티스텝 구현 계획 수립
triggers: ["계획", "plan", "플랜", "로드맵", "어떤 순서"]
roles: ["dev", "plan"]
type: workflow
steps:
  - id: analyze
    title: "요구사항 분석"
    gate: none
  - id: decompose
    title: "태스크 분해"
    gate: none
  - id: sequence
    title: "실행 순서 결정"
    gate: soft
  - id: write-plan
    title: "계획 문서화"
    gate: strict
---

## Step: analyze
요구사항과 제약 조건을 분석하세요.

확인할 것:
- 기존 코드베이스 구조와의 관계
- 외부 의존성 (라이브러리, API, 서비스)
- 기술적 제약 (성능, 호환성, 보안)
- 비기능 요구사항 (테스트, 문서화)

[STEP_COMPLETE: analyze]를 출력하여 다음 단계로 진행하세요.

## Step: decompose
작업을 독립적이고 테스트 가능한 단위로 분해하세요.

각 단위에 대해:
- 명확한 이름과 설명
- 예상 복잡도 (S/M/L)
- 입력/출력 인터페이스
- 의존하는 다른 단위

[STEP_COMPLETE: decompose]를 출력하여 다음 단계로 진행하세요.

## Step: sequence
의존성 기반으로 실행 순서를 결정하세요.

포함할 것:
- 순서도 (번호 매기기)
- 병렬 실행 가능한 태스크 식별
- 크리티컬 패스 명시
- 각 단계의 검증 방법

[STEP_COMPLETE: sequence]를 출력하여 다음 단계로 진행하세요.

## Step: write-plan
계획을 마크다운 문서로 정리하여 사용자에게 제시하세요.

형식:
- 목표 요약 (1-2문장)
- 태스크 목록 (체크박스)
- 각 태스크별 파일, 변경 내용, 테스트 방법
- 예상 리스크와 대응

사용자 승인 후 파일로 저장하세요.

[STEP_COMPLETE: write-plan]를 출력하여 스킬을 완료하세요.
```

- [ ] **Step 2: Create tdd.md**

```markdown
---
name: tdd
description: 테스트 주도 개발 사이클
triggers: ["tdd", "테스트 먼저", "test first", "테스트부터"]
roles: ["dev"]
type: workflow
steps:
  - id: write-test
    title: "실패하는 테스트 작성"
    gate: none
  - id: run-fail
    title: "테스트 실패 확인"
    gate: none
  - id: implement
    title: "최소 구현"
    gate: strict
  - id: run-pass
    title: "테스트 통과 확인"
    gate: none
  - id: refactor
    title: "리팩토링"
    gate: soft
---

## Step: write-test
구현할 기능의 기대 동작을 테스트로 먼저 작성하세요.

규칙:
- 핵심 동작을 검증하는 테스트 작성
- 엣지 케이스 최소 1개 포함
- 테스트 이름은 기대 동작을 설명 ("should return X when Y")
- 테스트만 작성하고 구현은 아직 하지 않음

[STEP_COMPLETE: write-test]를 출력하여 다음 단계로 진행하세요.

## Step: run-fail
작성한 테스트를 실행하여 올바른 이유로 실패하는지 확인하세요.

확인할 것:
- 테스트가 실패하는가? (실패해야 정상)
- 실패 이유가 "함수/모듈 없음" 또는 "기대값 불일치"인가?
- 엉뚱한 이유로 실패하면 테스트를 수정

shell_exec로 테스트 실행: `npx vitest run [테스트 파일]`

[STEP_COMPLETE: run-fail]를 출력하여 다음 단계로 진행하세요.

## Step: implement
테스트를 통과할 최소한의 구현을 작성하세요.

규칙:
- 테스트가 요구하는 것만 구현
- 불필요한 추가 기능 금지
- 구현 계획을 사용자에게 먼저 제시

사용자 승인 후 구현하세요.

[STEP_COMPLETE: implement]를 출력하여 다음 단계로 진행하세요.

## Step: run-pass
구현 후 테스트를 실행하여 모두 통과하는지 확인하세요.

확인할 것:
- 새 테스트 모두 통과
- 기존 테스트 회귀 없음 (전체 테스트 실행)
- 실패하면 구현 수정 (테스트 수정 아님)

[STEP_COMPLETE: run-pass]를 출력하여 다음 단계로 진행하세요.

## Step: refactor
코드 품질을 개선하세요.

가능한 개선:
- 중복 코드 제거
- 네이밍 개선
- 복잡한 조건문 단순화
- 불필요한 코드 삭제

리팩토링 후 테스트가 계속 통과하는지 확인하세요.

[STEP_COMPLETE: refactor]를 출력하여 스킬을 완료하세요.
```

- [ ] **Step 3: Create code-review.md**

```markdown
---
name: code-review
description: 변경사항 체계적 리뷰
triggers: ["리뷰", "review", "검토", "코드 봐줘", "확인해줘"]
roles: ["dev", "plan"]
type: workflow
steps:
  - id: diff-scan
    title: "변경 범위 파악"
    gate: none
  - id: logic-check
    title: "로직 검증"
    gate: none
  - id: security-check
    title: "보안 체크"
    gate: none
  - id: style-check
    title: "스타일/일관성"
    gate: none
  - id: summary
    title: "리뷰 요약"
    gate: strict
---

## Step: diff-scan
변경 범위를 전체 파악하세요.

실행:
- git_diff: 변경 내용 확인
- git_status: 변경 파일 목록

정리할 것:
- 파일별 변경 규모 (추가/수정/삭제 라인 수)
- 변경의 성격 (기능 추가, 버그 수정, 리팩토링)
- 영향 범위 (어떤 기능에 영향)

[STEP_COMPLETE: diff-scan]를 출력하여 다음 단계로 진행하세요.

## Step: logic-check
비즈니스 로직의 정확성을 검증하세요.

확인할 것:
- 조건문 로직이 올바른가
- 엣지 케이스가 처리되는가
- 오프바이원 에러 없는가
- null/undefined 처리가 적절한가
- 비동기 처리가 올바른가 (race condition, 미처리 Promise)

문제 발견 시 심각도 표시: CRITICAL / WARNING / INFO

[STEP_COMPLETE: logic-check]를 출력하여 다음 단계로 진행하세요.

## Step: security-check
보안 관련 문제를 검사하세요.

확인할 것:
- 사용자 입력이 검증 없이 사용되는가 (인젝션)
- API 키, 비밀이 하드코딩되어 있는가
- 민감 데이터가 로그에 노출되는가
- 파일 경로 조작 가능성 (path traversal)
- 안전하지 않은 정규식 (ReDoS)

[STEP_COMPLETE: security-check]를 출력하여 다음 단계로 진행하세요.

## Step: style-check
코드 스타일과 일관성을 확인하세요.

확인할 것:
- 기존 코드 컨벤션과 일치하는가
- 네이밍이 명확하고 일관적인가
- 파일 구조가 프로젝트 패턴을 따르는가
- 불필요한 주석이나 dead code 없는가
- import 정리가 되어 있는가

[STEP_COMPLETE: style-check]를 출력하여 다음 단계로 진행하세요.

## Step: summary
리뷰 결과를 심각도별로 정리하여 사용자에게 제시하세요.

형식:
### CRITICAL (즉시 수정 필요)
- [파일:라인] 문제 설명 + 수정 제안

### WARNING (수정 권장)
- [파일:라인] 문제 설명 + 수정 제안

### INFO (참고)
- [파일:라인] 개선 사항

문제가 없으면 "리뷰 통과"로 표시.
사용자 확인 후 수정 작업을 진행할 수 있습니다.

[STEP_COMPLETE: summary]를 출력하여 스킬을 완료하세요.
```

- [ ] **Step 4: Create refactoring.md**

```markdown
---
name: refactoring
description: 안전한 코드 리팩토링 가이드
triggers: ["리팩토링", "refactor", "정리", "개선해", "클린업"]
roles: ["dev"]
type: workflow
steps:
  - id: identify
    title: "개선 대상 식별"
    gate: none
  - id: safety-check
    title: "안전성 확인"
    gate: soft
  - id: extract
    title: "리팩토링 실행"
    gate: strict
  - id: verify
    title: "동작 검증"
    gate: none
---

## Step: identify
코드 스멜과 개선 대상을 식별하세요.

확인할 것:
- 긴 함수 (50줄 이상)
- 중복 코드
- 복잡한 조건문 (중첩 3단계 이상)
- 과도한 매개변수 (4개 이상)
- God 클래스/파일 (300줄 이상)
- 불명확한 네이밍

각 문제를 파일:라인과 함께 나열하세요.

[STEP_COMPLETE: identify]를 출력하여 다음 단계로 진행하세요.

## Step: safety-check
리팩토링 전 안전성을 확인하세요.

확인할 것:
- 대상 코드의 테스트 커버리지
- 테스트 없으면 먼저 테스트 추가를 제안
- 영향 범위 분석 (이 코드를 사용하는 곳)
- 리팩토링 후에도 동작이 동일한지 검증 방법

테스트 없이 리팩토링을 진행하지 마세요.

[STEP_COMPLETE: safety-check]를 출력하여 다음 단계로 진행하세요.

## Step: extract
리팩토링 계획을 사용자에게 제시하고 승인 후 실행하세요.

규칙:
- 한 번에 하나의 리팩토링만
- 적용할 패턴 명시 (Extract Function, Rename, Move, Inline 등)
- 변경할 파일과 예상 diff 미리보기
- 동작 변경 없음을 보장

사용자 승인 없이 코드를 변경하지 마세요.

[STEP_COMPLETE: extract]를 출력하여 다음 단계로 진행하세요.

## Step: verify
리팩토링 후 동작이 변경되지 않았는지 검증하세요.

확인할 것:
- 전체 테스트 통과 (shell_exec)
- 리팩토링 전후 동작 동일
- 새로운 경고나 에러 없음

실패하면 리팩토링을 되돌리고 원인 분석.

[STEP_COMPLETE: verify]를 출력하여 스킬을 완료하세요.
```

- [ ] **Step 5: Create security-check.md**

```markdown
---
name: security-check
description: 보안 취약점 스캔
triggers: ["보안", "security", "취약점", "vulnerability", "보안 검사"]
roles: ["dev", "plan"]
type: workflow
steps:
  - id: scan-deps
    title: "의존성 취약점 스캔"
    gate: none
  - id: check-inputs
    title: "입력 검증 확인"
    gate: none
  - id: check-secrets
    title: "비밀 노출 검사"
    gate: none
  - id: report
    title: "보안 리포트"
    gate: strict
---

## Step: scan-deps
의존성의 알려진 취약점을 확인하세요.

실행:
- shell_exec: `npm audit` 또는 `yarn audit`
- read_file: package.json의 의존성 목록 확인
- 오래된 의존성 식별 (마지막 업데이트 1년 이상)

결과를 심각도별로 정리하세요.

[STEP_COMPLETE: scan-deps]를 출력하여 다음 단계로 진행하세요.

## Step: check-inputs
사용자 입력이 들어오는 모든 경로를 추적하세요.

확인할 것:
- CLI 인수, 환경변수, 설정 파일 입력
- 외부 API 응답 처리
- 파일 읽기 (사용자 제공 경로)
- 셸 명령 실행 (사용자 제공 인수)

취약점 패턴:
- 명령 인젝션: 사용자 입력이 shell_exec에 직접 전달
- 경로 탐색: 사용자 입력이 파일 경로에 직접 사용
- XSS: 사용자 입력이 HTML 출력에 직접 포함
- 프롬프트 인젝션: 사용자 입력이 시스템 프롬프트에 직접 주입

[STEP_COMPLETE: check-inputs]를 출력하여 다음 단계로 진행하세요.

## Step: check-secrets
코드와 설정에서 비밀 노출을 검사하세요.

확인할 것:
- grep_search: API 키 패턴 (sk-, pk-, AKIA, ghp_, etc.)
- grep_search: password, secret, token, credential 문자열
- .gitignore에 .env, *.key 등이 포함되어 있는가
- 하드코딩된 URL에 인증 정보가 포함되어 있는가
- 로그 출력에 민감 정보가 포함될 수 있는가

[STEP_COMPLETE: check-secrets]를 출력하여 다음 단계로 진행하세요.

## Step: report
보안 검사 결과를 심각도별로 정리하여 사용자에게 제시하세요.

형식:
### CRITICAL (즉시 수정)
- [파일:라인] 취약점 설명 + 수정 방법

### HIGH (빠른 수정 필요)
- [파일:라인] 취약점 설명 + 수정 방법

### MEDIUM (개선 권장)
- [파일:라인] 취약점 설명 + 수정 방법

### LOW (참고)
- [파일:라인] 개선 사항

발견된 취약점이 없으면 "보안 검사 통과"로 표시.

[STEP_COMPLETE: report]를 출력하여 스킬을 완료하세요.
```

- [ ] **Step 6: Verify all 8 skills parse correctly**

Run: `npx tsx -e "import { loadAllSkills } from './src/core/skills.js'; const s = loadAllSkills(); console.log(s.length + ' skills loaded'); s.forEach(x => console.log('  ' + x.meta.name + ' [' + x.meta.type + '] steps:' + x.meta.steps.length));"`
Expected: 8 skills loaded with correct step counts.

- [ ] **Step 7: Commit**

```bash
git add src/skills/planning.md src/skills/tdd.md src/skills/code-review.md src/skills/refactoring.md src/skills/security-check.md
git commit -m "feat(skills): add 2nd batch built-in skills (planning, tdd, code-review, refactoring, security-check)"
```

---

### Task 10: End-to-End Verification

**Files:**
- No new files

- [ ] **Step 1: Run full test suite**

Run: `npx vitest run`
Expected: ALL PASS (including new skill-parser, skill-engine, skill-prompt tests)

- [ ] **Step 2: Manual TUI test**

Run: `npx tsx src/cli.ts --tui`

Test:
1. Type `/skill list` — should show 8 skills
2. Type `/skill info brainstorming` — should show steps and triggers
3. Type `/skill activate debugging` — should activate
4. Type `/skill status` — should show "debugging Step 1/5: 재현 확인"
5. Type `/skill deactivate` — should deactivate
6. Type "이 코드 설계해줘" — should auto-trigger brainstorming skill

- [ ] **Step 3: Manual text-mode test (if text-mode integration was added)**

Run: `npx tsx src/cli.ts --text`
Type `/skill list`

- [ ] **Step 4: Final commit**

```bash
git add -A
git commit -m "feat(skills): Level 2 skill system — complete implementation with 8 built-in skills

Adds workflow-aware skill system with:
- YAML frontmatter skill parser
- Trigger matching engine with role filtering
- Soft/strict gate support in agent loop
- 8 built-in skills: brainstorming, debugging, git-workflow, planning, tdd, code-review, refactoring, security-check
- /skill commands (list, activate, deactivate, next, status, info)
- TUI status bar integration
- Full backward compatibility with existing knowledge skills"
```
