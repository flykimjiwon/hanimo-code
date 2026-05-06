// Bidirectional tool-confirmation channel.
//
// Flow:
//   1. Chat loop sees CheckTool() == "ask", emits SSE event "tool_confirm"
//      with the tool call id and creates a pending channel keyed by that id.
//   2. WebView shows a modal, user clicks Allow/Deny, posts to /confirm.
//   3. Handler delivers the decision into the channel; loop unblocks.
//
// Timeout: 5 minutes. If unanswered, treated as "deny" so we never wedge
// the agent loop indefinitely.
package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type confirmDecision struct {
	Approve bool
	// If !Approve, reason gets surfaced as the tool result.
	Reason string
}

var (
	confirmMu      sync.Mutex
	confirmPending = map[string]chan confirmDecision{}
)

func registerConfirm(callID string) chan confirmDecision {
	ch := make(chan confirmDecision, 1)
	confirmMu.Lock()
	confirmPending[callID] = ch
	confirmMu.Unlock()
	return ch
}

func resolveConfirm(callID string, dec confirmDecision) bool {
	confirmMu.Lock()
	ch, ok := confirmPending[callID]
	if ok {
		delete(confirmPending, callID)
	}
	confirmMu.Unlock()
	if !ok {
		return false
	}
	select {
	case ch <- dec:
	default:
	}
	return true
}

func waitConfirm(callID string, timeout time.Duration) confirmDecision {
	confirmMu.Lock()
	ch := confirmPending[callID]
	confirmMu.Unlock()
	if ch == nil {
		return confirmDecision{Approve: false, Reason: "no pending confirmation"}
	}
	select {
	case d := <-ch:
		return d
	case <-time.After(timeout):
		confirmMu.Lock()
		delete(confirmPending, callID)
		confirmMu.Unlock()
		return confirmDecision{Approve: false, Reason: "user did not respond in time"}
	}
}

// /confirm — POST {id, approve, reason?}
func (s *server) handleConfirm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		ID      string `json:"id"`
		Approve bool   `json:"approve"`
		Reason  string `json:"reason,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if !resolveConfirm(body.ID, confirmDecision{Approve: body.Approve, Reason: body.Reason}) {
		http.Error(w, "no pending confirmation for that id", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
