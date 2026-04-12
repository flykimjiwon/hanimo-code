package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
)

// Request represents a JSON-RPC 2.0 request.
type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int64       `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// Response represents a JSON-RPC 2.0 response.
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

// RPCError represents a JSON-RPC 2.0 error.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *RPCError) Error() string {
	return fmt.Sprintf("rpc error %d: %s", e.Code, e.Message)
}

// Tool represents an MCP tool definition.
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// Transport defines the interface for MCP communication.
type Transport interface {
	Send(data []byte) error
	Receive() ([]byte, error)
	Close() error
}

// Client is a JSON-RPC 2.0 client for the MCP protocol.
type Client struct {
	transport Transport
	nextID    atomic.Int64
	mu        sync.Mutex
	tools     []Tool
}

// NewClient creates a new MCP client with the given transport.
func NewClient(transport Transport) *Client {
	return &Client{
		transport: transport,
	}
}

// call sends a JSON-RPC request and returns the raw result.
func (c *Client) call(method string, params interface{}) (json.RawMessage, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := c.nextID.Add(1)

	req := Request{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	if err := c.transport.Send(data); err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	respData, err := c.transport.Receive()
	if err != nil {
		return nil, fmt.Errorf("receive response: %w", err)
	}

	var resp Response
	if err := json.Unmarshal(respData, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Result, nil
}

// notify sends a JSON-RPC notification (no ID, no response expected).
func (c *Client) notify(method string, params interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	type Notification struct {
		JSONRPC string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params,omitempty"`
	}

	notif := Notification{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	data, err := json.Marshal(notif)
	if err != nil {
		return fmt.Errorf("marshal notification: %w", err)
	}

	return c.transport.Send(data)
}

// Initialize performs the MCP initialization handshake.
func (c *Client) Initialize(ctx context.Context) error {
	initParams := map[string]interface{}{
		"protocolVersion": "2025-06-18",
		"capabilities":    map[string]interface{}{},
		"clientInfo": map[string]interface{}{
			"name":    "hanimo",
			"version": "0.1.0",
		},
	}

	_, err := c.call("initialize", initParams)
	if err != nil {
		return fmt.Errorf("initialize: %w", err)
	}

	if err := c.notify("notifications/initialized", nil); err != nil {
		return fmt.Errorf("send initialized notification: %w", err)
	}

	return nil
}

// ListTools retrieves the list of available tools from the server.
func (c *Client) ListTools(ctx context.Context) ([]Tool, error) {
	result, err := c.call("tools/list", nil)
	if err != nil {
		return nil, fmt.Errorf("list tools: %w", err)
	}

	var listResult struct {
		Tools []Tool `json:"tools"`
	}
	if err := json.Unmarshal(result, &listResult); err != nil {
		return nil, fmt.Errorf("unmarshal tools: %w", err)
	}

	c.tools = listResult.Tools
	return listResult.Tools, nil
}

// CallTool invokes a tool on the MCP server and returns the text result.
func (c *Client) CallTool(ctx context.Context, name string, args map[string]interface{}) (string, error) {
	params := map[string]interface{}{
		"name":      name,
		"arguments": args,
	}

	result, err := c.call("tools/call", params)
	if err != nil {
		return "", fmt.Errorf("call tool %s: %w", name, err)
	}

	var callResult struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(result, &callResult); err != nil {
		return "", fmt.Errorf("unmarshal tool result: %w", err)
	}

	if len(callResult.Content) == 0 {
		return "", nil
	}

	return callResult.Content[0].Text, nil
}

// Close closes the underlying transport.
func (c *Client) Close() error {
	return c.transport.Close()
}
