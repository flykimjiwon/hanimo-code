import { McpClient } from './client.js';
import type { McpServerConfig } from '../config/schema.js';
import type { ToolSet } from 'ai';
import { mergeToolSets } from '../tools/registry.js';

interface ConnectedServer {
  client: McpClient;
  config: McpServerConfig;
  toolSet: ToolSet;
}

export class McpBridge {
  private servers: ConnectedServer[] = [];

  async loadFromConfig(
    servers: Record<string, McpServerConfig>
  ): Promise<void> {
    for (const [key, config] of Object.entries(servers)) {
      if (!config.enabled) continue;

      const client = new McpClient(key);
      try {
        await client.connect(config);
        const toolSet = await client.buildToolSet();
        this.servers.push({ client, config, toolSet });
      } catch (err: unknown) {
        const message = err instanceof Error ? err.message : String(err);
        console.error(`[MCP] Failed to connect to server "${key}": ${message}`);
      }
    }
  }

  getAvailableTools(
    networkMode: 'online' | 'offline'
  ): ToolSet | undefined {
    const sets: Array<ToolSet | undefined> = [];

    for (const server of this.servers) {
      if (networkMode === 'offline' && server.config.onlineOnly === true) {
        continue;
      }
      sets.push(server.toolSet);
    }

    return mergeToolSets(...sets);
  }

  async disconnectAll(): Promise<void> {
    await Promise.all(
      this.servers.map(async (server) => {
        try {
          await server.client.disconnect();
        } catch (err: unknown) {
          const message = err instanceof Error ? err.message : String(err);
          console.error(
            `[MCP] Error disconnecting "${server.config.name}": ${message}`
          );
        }
      })
    );
    this.servers = [];
  }

  getConnectedCount(): number {
    return this.servers.filter((s) => s.client.isConnected()).length;
  }
}
