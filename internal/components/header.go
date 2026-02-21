package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/matt-wright86/mardi-gras/internal/data"
	"github.com/matt-wright86/mardi-gras/internal/ui"
)

// Header renders the top title bar with bead string and counts.
type Header struct {
	Width  int
	Groups map[data.ParadeStatus][]data.Issue
}

// View renders the header.
func (h Header) View() string {
	rolling := len(h.Groups[data.ParadeRolling])
	linedUp := len(h.Groups[data.ParadeLinedUp])
	stalled := len(h.Groups[data.ParadeStalled])
	total := rolling + linedUp + stalled + len(h.Groups[data.ParadePastTheStand])

	title := ui.HeaderStyle.Render(
		fmt.Sprintf("%s MARDI GRAS %s", ui.FleurDeLis, ui.FleurDeLis),
	)

	counts := ui.HeaderCounts.Render(fmt.Sprintf(
		"%d rolling | %d lined up | %d stalled | %d total",
		rolling, linedUp, stalled, total,
	))

	titleLine := lipgloss.JoinHorizontal(
		lipgloss.Center,
		title,
		"  ",
		counts,
	)

	// Pad to full width
	titleLine = lipgloss.NewStyle().Width(h.Width).Render(titleLine)

	beadStr := h.renderBeadString()

	return lipgloss.JoinVertical(lipgloss.Left, titleLine, beadStr)
}

// renderBeadString creates the decorative bead string separator.
func (h Header) renderBeadString() string {
	colors := []lipgloss.Style{
		ui.BeadStylePurple,
		ui.BeadStyleGold,
		ui.BeadStyleGreen,
	}
	beads := []string{ui.BeadRound, ui.BeadDiamond}

	var parts []string
	ci := 0
	for i := 0; len(strings.Join(parts, "")) < h.Width-2; i++ {
		bead := beads[i%2]
		style := colors[ci%3]
		parts = append(parts, style.Render(bead))
		parts = append(parts, colors[ci%3].Render(ui.BeadDash))
		ci++
	}

	line := strings.Join(parts, "")
	return lipgloss.NewStyle().Width(h.Width).Render(line)
}
