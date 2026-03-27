export interface SubTask {
  id: string;
  description: string;
  type: 'analyze' | 'code' | 'review' | 'research' | 'general';
}

export interface SubAgentResult {
  taskId: string;
  status: 'fulfilled' | 'rejected';
  result?: string;
  error?: string;
}

export type OrchestratorEvent =
  | { type: 'decompose-start' }
  | { type: 'decompose-done'; tasks: SubTask[] }
  | { type: 'sub-agent-start'; taskId: string; index: number; total: number }
  | { type: 'sub-agent-done'; taskId: string; index: number; total: number; status: 'fulfilled' | 'rejected' }
  | { type: 'synthesize-start'; results: SubAgentResult[] }
  | { type: 'synthesize-done' };
