package session

import (
	"os"
	"testing"
)

func TestIntegration(t *testing.T) {
	dir, err := os.MkdirTemp("", "hanimo-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	if err := InitDB(dir); err != nil {
		t.Fatalf("InitDB: %v", err)
	}
	defer CloseDB()

	// Session CRUD
	sid, err := CreateSession("/tmp/project", "openai", "gpt-4o", "super")
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	if sid == "" {
		t.Fatal("empty session ID")
	}

	if err := SaveMessage(sid, "user", "hello", "", "", 10, 0); err != nil {
		t.Fatalf("SaveMessage: %v", err)
	}
	if err := NameSession(sid, "test session"); err != nil {
		t.Fatalf("NameSession: %v", err)
	}

	s, msgs, err := LoadSession(sid)
	if err != nil {
		t.Fatalf("LoadSession: %v", err)
	}
	if s.Name != "test session" {
		t.Errorf("expected name 'test session', got %q", s.Name)
	}
	if len(msgs) != 1 {
		t.Errorf("expected 1 message, got %d", len(msgs))
	}

	sessions, err := ListSessions(10)
	if err != nil {
		t.Fatalf("ListSessions: %v", err)
	}
	if len(sessions) != 1 {
		t.Errorf("expected 1 session, got %d", len(sessions))
	}

	// Search
	found, err := SearchSessions("hello")
	if err != nil {
		t.Fatalf("SearchSessions: %v", err)
	}
	if len(found) != 1 {
		t.Errorf("expected 1 search result, got %d", len(found))
	}

	// Fork
	forkID, err := ForkSession(sid)
	if err != nil {
		t.Fatalf("ForkSession: %v", err)
	}
	_, forkMsgs, err := LoadSession(forkID)
	if err != nil {
		t.Fatalf("LoadSession fork: %v", err)
	}
	if len(forkMsgs) != 1 {
		t.Errorf("expected 1 forked message, got %d", len(forkMsgs))
	}

	// Memory
	if err := SaveMemory("/tmp/project", "lang", "Go", "user"); err != nil {
		t.Fatalf("SaveMemory: %v", err)
	}
	mems, err := LoadMemories("/tmp/project")
	if err != nil {
		t.Fatalf("LoadMemories: %v", err)
	}
	if len(mems) != 1 {
		t.Errorf("expected 1 memory, got %d", len(mems))
	}

	prompt := FormatMemoriesForPrompt(mems)
	if prompt == "" {
		t.Error("expected non-empty prompt")
	}

	if err := DeleteMemory("/tmp/project", "lang"); err != nil {
		t.Fatalf("DeleteMemory: %v", err)
	}

	// Upsert test
	if err := SaveMemory("/tmp/project", "k", "v1", "auto"); err != nil {
		t.Fatalf("SaveMemory upsert1: %v", err)
	}
	if err := SaveMemory("/tmp/project", "k", "v2", "user"); err != nil {
		t.Fatalf("SaveMemory upsert2: %v", err)
	}
	mems2, _ := LoadMemories("/tmp/project")
	if len(mems2) != 1 || mems2[0].Value != "v2" {
		t.Errorf("upsert failed: got %+v", mems2)
	}

	// Usage
	if err := LogUsage(sid, "openai", "gpt-4o", 500, 200); err != nil {
		t.Fatalf("LogUsage: %v", err)
	}
	u, err := GetSessionUsage(sid)
	if err != nil {
		t.Fatalf("GetSessionUsage: %v", err)
	}
	if u.TokensIn != 500 || u.TokensOut != 200 {
		t.Errorf("unexpected session usage: %+v", u)
	}
	if u.CostUSD <= 0 {
		t.Errorf("expected positive cost, got %f", u.CostUSD)
	}

	total, err := GetTotalUsage()
	if err != nil {
		t.Fatalf("GetTotalUsage: %v", err)
	}
	if total.TokensIn != 500 {
		t.Errorf("unexpected total usage: %+v", total)
	}
}
