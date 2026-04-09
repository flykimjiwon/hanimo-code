package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// PaletteItem represents a single command in the palette.
type PaletteItem struct {
	Label       string
	Description string
	Action      string // slash command
}

// PaletteItems is the default list of commands available in the palette.
var PaletteItems = []PaletteItem{
	{Label: "Save Session", Description: "Save current session", Action: "/save"},
	{Label: "Load Session", Description: "Load previous session", Action: "/load"},
	{Label: "Search Sessions", Description: "Search past sessions", Action: "/search"},
	{Label: "Switch Model", Description: "Change LLM model", Action: "/model"},
	{Label: "Switch Provider", Description: "Change provider", Action: "/provider"},
	{Label: "Usage Stats", Description: "Token usage and cost", Action: "/usage"},
	{Label: "Diagnostics", Description: "Run code diagnostics", Action: "/diagnostics"},
	{Label: "Remember", Description: "Save a memory", Action: "/remember"},
	{Label: "Memories", Description: "Show saved memories", Action: "/memories"},
	{Label: "Config", Description: "Show configuration", Action: "/config"},
	{Label: "Switch Theme", Description: "Change color theme", Action: "/theme"},
	{Label: "Clear", Description: "Clear conversation", Action: "/clear"},
	{Label: "Help", Description: "Show help", Action: "/help"},
}

// FuzzyFilter returns items matching the query (case-insensitive substring match).
func FuzzyFilter(query string, items []PaletteItem) []PaletteItem {
	if query == "" {
		return items
	}
	query = strings.ToLower(query)
	var matched []PaletteItem
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Label), query) ||
			strings.Contains(strings.ToLower(item.Description), query) ||
			strings.Contains(item.Action, query) {
			matched = append(matched, item)
		}
	}
	return matched
}

// RenderPalette renders the command palette as a floating overlay.
func RenderPalette(items []PaletteItem, selected int, query string, width int) string {
	paletteWidth := 50
	if width < 60 {
		paletteWidth = width - 10
	}
	if paletteWidth < 30 {
		paletteWidth = 30
	}

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true)

	// Search input
	searchStyle := lipgloss.NewStyle().
		Foreground(ColorText)

	promptStyle := lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true)

	// Build content
	var lines []string
	lines = append(lines, titleStyle.Render("Command Palette"))
	lines = append(lines, "")

	// Search bar
	searchDisplay := query
	if searchDisplay == "" {
		searchDisplay = lipgloss.NewStyle().Foreground(ColorMuted).Render("Type to filter...")
	} else {
		searchDisplay = searchStyle.Render(query)
	}
	lines = append(lines, promptStyle.Render("> ")+searchDisplay)
	lines = append(lines, lipgloss.NewStyle().Foreground(ColorMuted).Render(strings.Repeat("─", paletteWidth-4)))

	// Items (max 10 visible)
	maxVisible := 10
	if len(items) < maxVisible {
		maxVisible = len(items)
	}

	if len(items) == 0 {
		lines = append(lines, lipgloss.NewStyle().Foreground(ColorMuted).Italic(true).Render("  No matching commands"))
	} else {
		// Calculate scroll offset to keep selected item visible
		offset := 0
		if selected >= maxVisible {
			offset = selected - maxVisible + 1
		}

		for i := offset; i < offset+maxVisible && i < len(items); i++ {
			item := items[i]
			labelWidth := paletteWidth - 8
			label := item.Label
			if lipgloss.Width(label) > labelWidth/2 {
				label = label[:labelWidth/2]
			}

			if i == selected {
				sel := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#1E1E2E")).
					Background(ColorPrimary).
					Bold(true)
				desc := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#1E1E2E")).
					Background(ColorPrimary)

				line := fmt.Sprintf(" %s  %s", label, item.Description)
				// Pad to fill width
				pad := paletteWidth - 4 - lipgloss.Width(line)
				if pad > 0 {
					line += strings.Repeat(" ", pad)
				}
				_ = desc // desc style used implicitly via sel
				lines = append(lines, sel.Render(line))
			} else {
				labelStyle := lipgloss.NewStyle().Foreground(ColorText)
				descStyle := lipgloss.NewStyle().Foreground(ColorMuted)
				line := fmt.Sprintf(" %s  %s", labelStyle.Render(label), descStyle.Render(item.Description))
				lines = append(lines, line)
			}
		}
	}

	lines = append(lines, "")
	hintStyle := lipgloss.NewStyle().Foreground(ColorMuted)
	lines = append(lines, hintStyle.Render("↑↓ navigate  Enter select  Esc close"))

	content := strings.Join(lines, "\n")

	// Floating box
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 2).
		Width(paletteWidth).
		Background(lipgloss.Color("#1E1E2E"))

	return box.Render(content)
}
