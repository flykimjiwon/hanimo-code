// Copyright 2025-2026 Kim Jiwon (김지원). All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/creack/pty"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// termSession is a real PTY-backed terminal session.
//
// 2026-05-03: cmd.StdinPipe / cmd.StdoutPipe 기반의 가짜 PTY에서 진짜 PTY로
// 전환. 이전 구현은 zsh/bash가 non-interactive 모드로 들어가 echo·prompt·
// line editing이 모두 죽었음 (사용자가 입력해도 글자 안 보이는 원인).
type termSession struct {
	cmd    *exec.Cmd
	pty    io.ReadWriteCloser // PTY master
	mu     sync.Mutex
	closed bool
}

func (a *App) GetAvailableShells() []string {
	if runtime.GOOS == "windows" {
		var found []string
		for _, s := range []string{"powershell.exe", "cmd.exe"} {
			if p, err := exec.LookPath(s); err == nil {
				found = append(found, p)
			}
		}
		return found
	}
	var found []string
	for _, s := range []string{"/bin/zsh", "/bin/bash", "/bin/sh"} {
		if _, err := os.Stat(s); err == nil {
			found = append(found, s)
		}
	}
	return found
}

func (a *App) GetCurrentShell() string {
	if a.shellPath != "" {
		return a.shellPath
	}
	if runtime.GOOS == "windows" {
		if p, err := exec.LookPath("powershell.exe"); err == nil {
			return p
		}
		return "cmd.exe"
	}
	if s := os.Getenv("SHELL"); s != "" {
		return s
	}
	return "/bin/bash"
}

func (a *App) SetShell(shell string) error {
	a.shellPath = shell
	a.StopTerminal()
	return a.StartTerminal()
}

func (a *App) StartTerminal() error {
	if a.term != nil && !a.term.closed {
		return nil
	}

	shell := a.GetCurrentShell()
	cmd := exec.Command(shell)
	cmd.Dir = a.cwd

	// 환경 변수 — TERM은 컬러 처리, COLORFGBG는 일부 셸의 컬러 감지 보조
	if runtime.GOOS != "windows" {
		cmd.Env = append(os.Environ(),
			"TERM=xterm-256color",
			"COLORTERM=truecolor",
		)
	}

	// Platform-specific: hide console window on Windows
	hideConsoleWindow(cmd)

	// 진짜 PTY 시작 — pty.Start 가 master pty(io.ReadWriteCloser) 반환
	// cmd.Stdin/Stdout/Stderr 는 자동으로 slave에 연결됨
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	term := &termSession{cmd: cmd, pty: ptmx}
	a.term = term

	// PTY → wails event 스트림 (xterm.js가 그대로 받아서 렌더)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := ptmx.Read(buf)
			if n > 0 {
				term.mu.Lock()
				closed := term.closed
				term.mu.Unlock()
				if closed {
					break
				}
				wailsRuntime.EventsEmit(a.ctx, "term:output", string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
		term.mu.Lock()
		alreadyClosed := term.closed
		term.mu.Unlock()
		if !alreadyClosed {
			wailsRuntime.EventsEmit(a.ctx, "term:output", "\r\n[Process exited]\r\n")
		}
	}()

	return nil
}

// WriteTerminal forwards xterm.js keystrokes to the PTY master. PTY echoes
// back automatically (in canonical mode), so the user sees what they type.
func (a *App) WriteTerminal(input string) {
	if a.term == nil {
		return
	}
	a.term.mu.Lock()
	defer a.term.mu.Unlock()
	if a.term.closed || a.term.pty == nil {
		return
	}
	_, _ = a.term.pty.Write([]byte(input))
}

// ResizeTerminal resizes the PTY window so applications like vim/htop redraw
// at the right dimensions. Without this, line wrapping breaks.
func (a *App) ResizeTerminal(rows, cols int) {
	if a.term == nil || a.term.pty == nil {
		return
	}
	if f, ok := a.term.pty.(*os.File); ok {
		_ = pty.Setsize(f, &pty.Winsize{
			Rows: uint16(rows),
			Cols: uint16(cols),
		})
	}
}

func (a *App) StopTerminal() {
	if a.term == nil {
		return
	}
	a.term.mu.Lock()
	a.term.closed = true
	a.term.mu.Unlock()
	if a.term.pty != nil {
		a.term.pty.Close()
	}
	if a.term.cmd != nil && a.term.cmd.Process != nil {
		_ = a.term.cmd.Process.Kill()
		_ = a.term.cmd.Wait()
	}
	a.term = nil
}
