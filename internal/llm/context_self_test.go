package llm

import (
	"strings"
	"testing"
)

func TestIsHanimoOwnRepo(t *testing.T) {
	if !IsHanimoOwnRepo() {
		t.Error("expected IsHanimoOwnRepo to return true when running under hanimo's test suite")
	}
}

func TestGatherSystemContextSelfWarning(t *testing.T) {
	ctx := GatherSystemContext()
	if !strings.Contains(ctx, "SELF-REPO WARNING") {
		t.Errorf("expected self-repo warning in system context, got:\n%s", ctx)
	}
}
