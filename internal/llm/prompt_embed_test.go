package llm

import (
	"strings"
	"testing"
)

func TestSystemPromptEmbedded(t *testing.T) {
	cases := []struct {
		mode    Mode
		mustHave []string
	}{
		{ModeSuper, []string{"PRIMARY DIRECTIVE", "만능 AI 코딩 에이전트", "ASK_USER"}},
		{ModeDev, []string{"PRIMARY DIRECTIVE", "Deep Agent", "TASK_COMPLETE"}},
		{ModePlan, []string{"PRIMARY DIRECTIVE", "Plan — 계획 우선", "단계별"}},
	}
	for _, c := range cases {
		got := SystemPrompt(c.mode)
		if len(got) < 500 {
			t.Errorf("mode=%v: prompt too short (%d bytes) — embed may be broken", c.mode, len(got))
		}
		for _, m := range c.mustHave {
			if !strings.Contains(got, m) {
				t.Errorf("mode=%v: missing %q in prompt", c.mode, m)
			}
		}
	}
}

func TestSystemPromptOrderedSections(t *testing.T) {
	// The three sections (core → body → askuser) must concatenate in order
	// for prompt-cache breakpoints to line up.
	got := SystemPrompt(ModeSuper)
	core := strings.Index(got, "PRIMARY DIRECTIVE")
	body := strings.Index(got, "Smart all-in-one")
	ask := strings.Index(got, "Interactive Questions")
	if core < 0 || body < 0 || ask < 0 {
		t.Fatalf("missing section markers: core=%d body=%d ask=%d", core, body, ask)
	}
	if !(core < body && body < ask) {
		t.Errorf("sections out of order: core=%d body=%d ask=%d", core, body, ask)
	}
}
