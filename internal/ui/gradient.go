package ui

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

// ApplyMardiGrasGradient applies a smooth Purple -> Gold -> Green gradient to the text.
func ApplyMardiGrasGradient(text string) string {
	runes := []rune(text)
	width := len(runes)
	if width == 0 {
		return ""
	}

	c1, _ := colorful.Hex(string(Purple))
	c2, _ := colorful.Hex(string(Gold))
	c3, _ := colorful.Hex(string(Green))

	var b strings.Builder
	for i, r := range runes {
		t := 0.0
		if width > 1 {
			t = float64(i) / float64(width-1)
		}
		var c colorful.Color
		if t < 0.5 {
			c = c1.BlendLuv(c2, t*2)
		} else {
			c = c2.BlendLuv(c3, (t-0.5)*2)
		}

		s := lipgloss.NewStyle().Foreground(lipgloss.Color(c.Hex()))
		b.WriteString(s.Render(string(r)))
	}
	return b.String()
}

// ApplyShimmerGradient applies the Mardi Gras gradient with a phase offset that shifts over time.
// The offset (0.0-1.0) rotates the gradient start point, creating a wave effect.
// A sine-based brightness modulation adds sparkle to individual characters.
func ApplyShimmerGradient(text string, offset float64) string {
	runes := []rune(text)
	width := len(runes)
	if width == 0 {
		return ""
	}

	c1, _ := colorful.Hex(string(Purple))
	c2, _ := colorful.Hex(string(Gold))
	c3, _ := colorful.Hex(string(Green))

	var b strings.Builder
	for i, r := range runes {
		// Base position with offset shift (wraps around)
		t := 0.0
		if width > 1 {
			t = float64(i)/float64(width-1) + offset
		}
		// Wrap to [0, 1]
		t = t - math.Floor(t)

		// Three-stop gradient with wrap: Purple → Gold → Green → Purple
		var c colorful.Color
		switch {
		case t < 1.0/3:
			c = c1.BlendLuv(c2, t*3)
		case t < 2.0/3:
			c = c2.BlendLuv(c3, (t-1.0/3)*3)
		default:
			c = c3.BlendLuv(c1, (t-2.0/3)*3)
		}

		// Sine-based brightness sparkle
		sparkle := 0.8 + 0.2*math.Sin(float64(i)*0.7+offset*math.Pi*6)
		h, s, l := c.Hsl()
		c = colorful.Hsl(h, s, l*sparkle)

		s2 := lipgloss.NewStyle().Foreground(lipgloss.Color(c.Hex()))
		b.WriteString(s2.Render(string(r)))
	}
	return b.String()
}

// ApplyPartialMardiGrasGradient applies the gradient as if the text was `totalLength` characters long,
// ensuring a partial progress bar maps to the correct segment of the full color spectrum.
func ApplyPartialMardiGrasGradient(text string, totalLength int) string {
	runes := []rune(text)
	if totalLength == 0 {
		return ""
	}

	c1, _ := colorful.Hex(string(Purple))
	c2, _ := colorful.Hex(string(Gold))
	c3, _ := colorful.Hex(string(Green))

	var b strings.Builder
	for i, r := range runes {
		t := 0.0
		if totalLength > 1 {
			t = float64(i) / float64(totalLength-1)
		}
		var c colorful.Color
		if t < 0.5 {
			c = c1.BlendLuv(c2, t*2)
		} else {
			c = c2.BlendLuv(c3, (t-0.5)*2)
		}

		s := lipgloss.NewStyle().Foreground(lipgloss.Color(c.Hex()))
		b.WriteString(s.Render(string(r)))
	}
	return b.String()
}
