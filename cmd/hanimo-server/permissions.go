// Permission policy for tool execution. Persisted to ~/.hanimo/permissions.yaml
// so both the TUI and the VS Code extension share the same rules.
//
// Model:
//   - Per-tool toggle: "allow" or "deny" (default = allow).
//   - Shell command denylist: regex patterns checked against shell_exec args
//     before execution. Matching → tool result becomes "denied by policy".
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/flykimjiwon/hanimo/internal/config"
)

type Permissions struct {
	// Tools maps tool name → "allow" | "deny" | "ask". Missing keys default to allow.
	Tools     map[string]string `yaml:"tools" json:"tools"`
	ShellDeny []string          `yaml:"shell_deny" json:"shell_deny"`
}

var (
	permsMu  sync.RWMutex
	cachedP  *Permissions
	shellRes []*regexp.Regexp
)

func defaultPermissions() Permissions {
	return Permissions{
		Tools: map[string]string{
			// All tools allowed by default. Listed for visibility in the UI.
			"file_read": "allow", "file_write": "allow", "file_edit": "allow",
			"list_files": "allow", "shell_exec": "allow",
			"grep_search": "allow", "glob_search": "allow",
			"hashline_read": "allow", "hashline_edit": "allow",
			"git_status": "allow", "git_diff": "allow", "git_log": "allow",
			"diagnostics": "allow", "knowledge_search": "allow", "apply_patch": "allow",
		},
		ShellDeny: []string{
			`\brm\s+-rf\s+/`,             // rm -rf /
			`\brm\s+-rf\s+\*`,            // rm -rf *
			`\bsudo\s+rm\b`,              // sudo rm ...
			`\bgit\s+push\s+.*--force\b`, // git push --force
			`\b:\(\)\s*\{.*:\|:`,         // fork bomb
			`\b>\s*/dev/sd[a-z]`,         // overwriting block devices
			`\bdd\s+if=.*of=/dev/`,       // dd to raw device
			`\bmkfs\b`,                   // formatting filesystems
			`\bchmod\s+777\s+/`,          // recursive perm wipe at root
			`\beval\s+\$\(curl`,          // eval curl pipes
			`\bcurl\s+[^|]*\|\s*sh\b`,    // curl | sh
			`\bcurl\s+[^|]*\|\s*bash\b`,  // curl | bash
		},
	}
}

func permissionsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, config.ConfigDirName)
	_ = os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "permissions.yaml"), nil
}

// loadPermissions returns the cached policy, loading from disk on first use
// and creating defaults if the file is missing.
func loadPermissions() *Permissions {
	permsMu.RLock()
	if cachedP != nil {
		defer permsMu.RUnlock()
		return cachedP
	}
	permsMu.RUnlock()

	permsMu.Lock()
	defer permsMu.Unlock()
	if cachedP != nil {
		return cachedP
	}

	path, err := permissionsPath()
	if err != nil {
		def := defaultPermissions()
		cachedP = &def
		compileShellRegex()
		return cachedP
	}
	data, err := os.ReadFile(path)
	if err != nil {
		def := defaultPermissions()
		cachedP = &def
		_ = savePermissionsLocked()
		compileShellRegex()
		return cachedP
	}
	var p Permissions
	if err := yaml.Unmarshal(data, &p); err != nil {
		def := defaultPermissions()
		cachedP = &def
		compileShellRegex()
		return cachedP
	}
	if p.Tools == nil {
		p.Tools = defaultPermissions().Tools
	}
	cachedP = &p
	compileShellRegex()
	return cachedP
}

func savePermissionsLocked() error {
	path, err := permissionsPath()
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(cachedP)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func compileShellRegex() {
	shellRes = shellRes[:0]
	for _, p := range cachedP.ShellDeny {
		re, err := regexp.Compile(p)
		if err == nil {
			shellRes = append(shellRes, re)
		}
	}
}

// CheckTool returns the policy decision for a tool invocation:
//   "allow" — proceed
//   "deny"  — block, with reason
//   "ask"   — pause and request user confirmation
// shell_exec runs the denylist regardless of tool-level policy.
func CheckTool(tool, args string) (decision string, reason string) {
	p := loadPermissions()
	mode := p.Tools[tool]
	if mode == "" {
		mode = "allow"
	}
	if tool == "shell_exec" {
		permsMu.RLock()
		for _, re := range shellRes {
			if re.MatchString(args) {
				permsMu.RUnlock()
				return "deny", fmt.Sprintf("Shell command blocked by policy (pattern: %s).", re.String())
			}
		}
		permsMu.RUnlock()
	}
	switch mode {
	case "deny":
		return "deny", fmt.Sprintf("Tool %q is denied by policy.", tool)
	case "ask":
		return "ask", ""
	default:
		return "allow", ""
	}
}

// /permissions — read/write the policy. Used by the VS Code settings panel.
func (s *server) handlePermissions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, loadPermissions())
	case http.MethodPut, http.MethodPost:
		var p Permissions
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if p.Tools == nil {
			p.Tools = defaultPermissions().Tools
		}
		permsMu.Lock()
		cachedP = &p
		compileShellRegex()
		_ = savePermissionsLocked()
		permsMu.Unlock()
		writeJSON(w, http.StatusOK, p)
	case http.MethodDelete:
		permsMu.Lock()
		def := defaultPermissions()
		cachedP = &def
		compileShellRegex()
		_ = savePermissionsLocked()
		permsMu.Unlock()
		writeJSON(w, http.StatusOK, cachedP)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// shortenArgs is a debug helper for log/error messages.
func shortenArgs(s string) string {
	if len(s) > 80 {
		return s[:77] + "…"
	}
	return strings.ReplaceAll(s, "\n", " ")
}
