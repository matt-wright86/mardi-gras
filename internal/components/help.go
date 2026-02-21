package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/matt-wright86/mardi-gras/internal/ui"
)

// Help renders the global ? help modal.
type Help struct {
	Width  int
	Height int
}

// NewHelp creates a new help rendering component.
func NewHelp(width, height int) Help {
	return Help{Width: width, Height: height}
}

// View returns the rendered modal block positioned at the center of the terminal.
func (h Help) View() string {
	var lines []string

	lines = append(lines, "")
	lines = append(lines, "  "+ui.HelpTitle.Render("[ MARDI GRAS HELP ]"))
	lines = append(lines, "")

	// Global Keybindings
	lines = append(lines, "  "+ui.HelpSection.Render("GLOBAL BINDINGS"))
	lines = append(lines, "  "+h.formatBinding("q", "Quit application"))
	lines = append(lines, "  "+h.formatBinding("tab", "Switch active pane"))
	lines = append(lines, "  "+h.formatBinding("?", "Toggle this help menu"))
	lines = append(lines, "")

	// Parade Keybindings
	lines = append(lines, "  "+ui.HelpSection.Render("PARADE PANE"))
	lines = append(lines, "  "+h.formatBinding("j / k", "Navigate up/down"))
	lines = append(lines, "  "+h.formatBinding("g / G", "Jump to top/bottom"))
	lines = append(lines, "  "+h.formatBinding("enter", "Focus detail pane"))
	lines = append(lines, "  "+h.formatBinding("c", "Toggle closed issues"))
	lines = append(lines, "  "+h.formatBinding("/", "Enter filter mode"))
	lines = append(lines, "")

	// Detail Keybindings
	lines = append(lines, "  "+ui.HelpSection.Render("DETAIL PANE"))
	lines = append(lines, "  "+h.formatBinding("j / k", "Scroll up/down"))
	lines = append(lines, "  "+h.formatBinding("esc", "Back to parade pane"))
	lines = append(lines, "  "+h.formatBinding("/", "Enter filter mode"))
	lines = append(lines, "")

	// Filter Mode Keybindings
	lines = append(lines, "  "+ui.HelpSection.Render("FILTER MODE"))
	lines = append(lines, "  "+h.formatBinding("esc", "Clear query and exit filter mode"))
	lines = append(lines, "  "+h.formatBinding("enter", "Apply query and exit filter mode"))
	lines = append(lines, "  "+h.formatBinding("type:bug", "Match issue type"))
	lines = append(lines, "  "+h.formatBinding("p0, p1...", "Match priority level"))
	lines = append(lines, "")

	content := strings.Join(lines, "\n")
	box := ui.HelpOverlayBg.Render(content)

	// Since Lipgloss doesn't have true absolute positioning (z-index), we place
	// the box in a container that represents the full terminal screen, centering it.
	return lipgloss.Place(h.Width, h.Height, lipgloss.Center, lipgloss.Center, box)
}

func (h Help) formatBinding(key, desc string) string {
	// Pad the key to align the descriptions nicely
	keyStr := ui.HelpKey.Render(key)
	padding := 12 - len(key)
	if padding < 1 {
		padding = 1
	}
	return keyStr + strings.Repeat(" ", padding) + ui.HelpDesc.Render(desc)
}
