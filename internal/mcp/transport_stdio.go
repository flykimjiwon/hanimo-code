package mcp

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

// StdioTransport communicates with an MCP server via stdin/stdout
// using Content-Length framing.
type StdioTransport struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	reader *bufio.Reader
}

// NewStdioTransport starts the given command and returns a transport
// that communicates over its stdin/stdout.
func NewStdioTransport(command string, args []string) (*StdioTransport, error) {
	cmd := exec.Command(command, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("get stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return nil, fmt.Errorf("get stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		stdin.Close()
		return nil, fmt.Errorf("start command: %w", err)
	}

	return &StdioTransport{
		cmd:    cmd,
		stdin:  stdin,
		reader: bufio.NewReader(stdout),
	}, nil
}

// Send writes a Content-Length framed message to the process stdin.
func (t *StdioTransport) Send(data []byte) error {
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
	if _, err := io.WriteString(t.stdin, header); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	if _, err := t.stdin.Write(data); err != nil {
		return fmt.Errorf("write body: %w", err)
	}
	return nil
}

// Receive reads a Content-Length framed message from the process stdout.
func (t *StdioTransport) Receive() ([]byte, error) {
	var contentLength int

	for {
		line, err := t.reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("read header line: %w", err)
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			break
		}

		if strings.HasPrefix(line, "Content-Length: ") {
			val := strings.TrimPrefix(line, "Content-Length: ")
			contentLength, err = strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("parse content length %q: %w", val, err)
			}
		}
	}

	if contentLength == 0 {
		return nil, fmt.Errorf("missing or zero Content-Length header")
	}

	body := make([]byte, contentLength)
	if _, err := io.ReadFull(t.reader, body); err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return body, nil
}

// Close closes stdin and waits for the process to exit.
func (t *StdioTransport) Close() error {
	t.stdin.Close()
	return t.cmd.Wait()
}
