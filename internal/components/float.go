package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/matt-wright86/mardi-gras/internal/ui"
)

// Float renders a decorative ASCII parade float. Stub for v1.
func Float(title string, width int) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.BrightGold).
		Foreground(ui.BrightPurple).
		Bold(true).
		Width(width-4).
		Align(lipgloss.Center).
		Padding(0, 1)

	return style.Render(ui.FleurDeLis + " " + title + " " + ui.FleurDeLis)
}
