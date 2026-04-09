package agents

import "strings"

const (
	MaxAutoIterations  = 20
	AutoCompleteMarker = "[AUTO_COMPLETE]"
	AutoPauseMarker    = "[AUTO_PAUSE]"
)

// AutoPromptSuffix is appended to the system prompt when auto mode is active
const AutoPromptSuffix = `

You are in AUTONOMOUS MODE. Complete the task independently:
- Use tools to read, write, and test code
- Run diagnostics to verify your work
- When the task is fully complete, output [AUTO_COMPLETE]
- If you're blocked and need human input, output [AUTO_PAUSE]
- Do NOT ask questions — make decisions and proceed`

// CheckAutoMarkers checks response content for auto-mode control markers
func CheckAutoMarkers(content string) (complete bool, pause bool) {
	complete = strings.Contains(content, AutoCompleteMarker)
	pause = strings.Contains(content, AutoPauseMarker)
	return
}
