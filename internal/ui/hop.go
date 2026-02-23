package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderStars returns a star rating string like "★★★★☆" for a 0.0-1.0 score.
func RenderStars(score float32) string {
	stars := max(min(int(score*5+0.5), 5), 0)

	color := QualityColor(score)
	filledStyle := lipgloss.NewStyle().Foreground(color)
	emptyStyle := lipgloss.NewStyle().Foreground(Dim)

	var b strings.Builder
	for i := 0; i < stars; i++ {
		b.WriteString(filledStyle.Render(SymStar))
	}
	for i := stars; i < 5; i++ {
		b.WriteString(emptyStyle.Render(SymStarEmpty))
	}
	return b.String()
}

// RenderStarsCompact returns a compact star badge like "★4" for parade rows.
func RenderStarsCompact(score float32) string {
	stars := max(min(int(score*5+0.5), 5), 0)

	color := QualityColor(score)
	return lipgloss.NewStyle().Foreground(color).Render(SymStar + string(rune('0'+stars)))
}
