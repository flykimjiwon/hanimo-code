export interface RoleModelConfig {
  preferred?: string;
  fallback?: string;
}

export interface RoleSubAgentConfig {
  enabled: boolean;
  count: 3 | 5 | 10;
  model?: string;
}

export interface RoleDefinition {
  id: string;
  name: string;
  description: string;
  icon: string;
  systemPrompt: string;
  tools: string[];
  maxSteps: number;
  temperature?: number;
  model?: RoleModelConfig;
  subAgents?: RoleSubAgentConfig;
}
