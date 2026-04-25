package main

import (
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// emitHashAnchorsForLines flashes anchors at a specific line range. Phase 10
// hashline_edit handler uses this so only the lines actually edited by the
// agent show the 🔒 marker — instead of file_write's broad flash.
func emitHashAnchorsForLines(a *App, newContent string, startLine, count int, clearAfter time.Duration) {
	if a == nil || a.ctx == nil || count <= 0 {
		return
	}
	lines := strings.Split(newContent, "\n")
	if count > len(lines) {
		count = len(lines)
	}
	anchors := make([]hashAnchorPayload, 0, count)
	for i := 0; i < count; i++ {
		ln := lines[i]
		if strings.TrimSpace(ln) == "" {
			continue
		}
		anchors = append(anchors, hashAnchorPayload{Line: startLine + i, Hash: hashLine(ln)})
	}
	runtime.EventsEmit(a.ctx, "hash:anchor", anchors)
	go func() {
		time.Sleep(clearAfter)
		runtime.EventsEmit(a.ctx, "hash:anchor", []hashAnchorPayload{})
	}()
}

// hashAnchorPayload mirrors the frontend HashAnchor type in
// frontend/src/components/hashAnchorGutter.ts.
type hashAnchorPayload struct {
	Line int    `json:"line"`
	Hash string `json:"hash"`
}

// emitHashAnchorsFor publishes a 'hash:anchor' Wails event with the first
// `cap` line anchors of `content`, then auto-clears after `clearAfter`.
// Phase 9 demo: every file_write triggers a brief anchor flash so the
// user sees the brand promise visually. Real per-line edit anchors will
// arrive when hashline_edit tool is ported (Phase 10+).
func emitHashAnchorsFor(a *App, content string, cap int, clearAfter time.Duration) {
	if a == nil || a.ctx == nil {
		return
	}
	lines := strings.Split(content, "\n")
	if len(lines) > cap {
		lines = lines[:cap]
	}
	anchors := make([]hashAnchorPayload, 0, len(lines))
	for i, ln := range lines {
		if strings.TrimSpace(ln) == "" {
			continue
		}
		anchors = append(anchors, hashAnchorPayload{Line: i + 1, Hash: hashLine(ln)})
	}
	runtime.EventsEmit(a.ctx, "hash:anchor", anchors)
	go func() {
		time.Sleep(clearAfter)
		runtime.EventsEmit(a.ctx, "hash:anchor", []hashAnchorPayload{})
	}()
}
