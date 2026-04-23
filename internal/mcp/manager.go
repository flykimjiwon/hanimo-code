package mcp

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/flykimjiwon/hanimo/internal/config"
)

// MCPTool is a tool from an MCP server with routing metadata.
type MCPTool struct {
	Name        string
	Description string
	InputSchema map[string]interface{}
	ServerName  string
	OrigName    string
}

// Manager holds multiple MCP clients, one per configured server.
type Manager struct {
	servers []config.MCPServer
	clients []*Client
}

// NewManager creates a Manager for the given server configs.
func NewManager(servers []config.MCPServer) *Manager {
	return &Manager{servers: servers}
}

// Start initializes all clients and collects their tools.
func (m *Manager) Start() error {
	for _, srv := range m.servers {
		client, err := m.createClient(srv)
		if err != nil {
			config.DebugLog("[MCP-MGR] failed to create client server=%s: %v", srv.Name, err)
			m.clients = append(m.clients, nil)
			continue
		}

		ctx := context.Background()
		if err := client.Initialize(ctx); err != nil {
			config.DebugLog("[MCP-MGR] failed to initialize server=%s: %v", srv.Name, err)
			client.Close()
			m.clients = append(m.clients, nil)
			continue
		}

		if _, err := client.ListTools(ctx); err != nil {
			config.DebugLog("[MCP-MGR] failed to list tools server=%s: %v", srv.Name, err)
		}

		m.clients = append(m.clients, client)
		config.DebugLog("[MCP-MGR] started server=%s tools=%d", srv.Name, len(client.tools))
	}
	return nil
}

// createClient creates the appropriate transport and client for a server config.
func (m *Manager) createClient(srv config.MCPServer) (*Client, error) {
	var transport Transport
	var err error

	switch srv.Transport {
	case "stdio":
		if srv.Command == "" {
			return nil, fmt.Errorf("stdio transport requires 'command'")
		}
		// Set environment variables if configured
		for k, v := range srv.Env {
			os.Setenv(k, v)
		}
		transport, err = NewStdioTransport(srv.Command, srv.Args)
		if err != nil {
			return nil, fmt.Errorf("create stdio transport: %w", err)
		}

	case "sse":
		if srv.URL == "" {
			return nil, fmt.Errorf("sse transport requires 'url'")
		}
		transport, err = NewSSETransport(srv.URL)
		if err != nil {
			return nil, fmt.Errorf("create sse transport: %w", err)
		}

	default:
		return nil, fmt.Errorf("unknown transport: %q (supported: stdio, sse)", srv.Transport)
	}

	return NewClient(transport), nil
}

// Stop shuts down all running clients.
func (m *Manager) Stop() {
	for _, c := range m.clients {
		if c != nil {
			c.Close()
		}
	}
}

// AllTools returns the combined tool list from all connected servers,
// with names prefixed as mcp_{servername}_{toolname}.
func (m *Manager) AllTools() []MCPTool {
	var out []MCPTool
	for i, c := range m.clients {
		if c == nil {
			continue
		}
		srv := m.servers[i]
		for _, t := range c.tools {
			out = append(out, MCPTool{
				Name:        mcpToolName(srv.Name, t.Name),
				Description: t.Description,
				InputSchema: t.InputSchema,
				ServerName:  srv.Name,
				OrigName:    t.Name,
			})
		}
	}
	return out
}

// CallTool routes a tool call to the correct server by prefixed name.
func (m *Manager) CallTool(name string, args map[string]interface{}) (string, error) {
	ctx := context.Background()
	for i, c := range m.clients {
		if c == nil {
			continue
		}
		srv := m.servers[i]
		for _, t := range c.tools {
			if mcpToolName(srv.Name, t.Name) == name {
				config.DebugLog("[MCP-MGR] call server=%s tool=%s", srv.Name, t.Name)
				return c.CallTool(ctx, t.Name, args)
			}
		}
	}
	return "", fmt.Errorf("no MCP server has tool: %s", name)
}

// Status returns a human-readable status string for the /mcp command.
func (m *Manager) Status() string {
	if len(m.servers) == 0 {
		return "No MCP servers configured (check config.yaml mcp.servers)"
	}
	var lines []string
	for i, srv := range m.servers {
		c := m.clients[i]
		if c == nil {
			lines = append(lines, fmt.Sprintf("  x %s (%s) — connection failed", srv.Name, srv.Transport))
			continue
		}
		lines = append(lines, fmt.Sprintf("  + %s (%s) — %d tools", srv.Name, srv.Transport, len(c.tools)))
		for _, t := range c.tools {
			lines = append(lines, fmt.Sprintf("      - %s", mcpToolName(srv.Name, t.Name)))
		}
	}
	return strings.Join(lines, "\n")
}

func mcpToolName(serverName, toolName string) string {
	safe := strings.ReplaceAll(serverName, "-", "_")
	safe = strings.ReplaceAll(safe, " ", "_")
	return "mcp_" + safe + "_" + toolName
}
