package mcp

import "fmt"

// SSETransport is a stub for Server-Sent Events transport.
// Will be implemented in Phase 2.
type SSETransport struct {
	url string
}

// NewSSETransport creates a new SSE transport stub.
func NewSSETransport(url string) (*SSETransport, error) {
	return &SSETransport{url: url}, nil
}

// Send is not yet implemented.
func (t *SSETransport) Send(data []byte) error {
	return fmt.Errorf("SSE transport not yet implemented")
}

// Receive is not yet implemented.
func (t *SSETransport) Receive() ([]byte, error) {
	return nil, fmt.Errorf("SSE transport not yet implemented")
}

// Close is a no-op for the stub.
func (t *SSETransport) Close() error {
	return nil
}
