package ui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestNewGradient(t *testing.T) {
	g := NewGradient(BrightGreen, BrightGold, lipgloss.Color("#E74C3C"))
	// Should produce 101 styles (0-100%)
	s0 := g.At(0)
	s100 := g.At(100)
	if s0.GetForeground() == s100.GetForeground() {
		t.Error("gradient endpoints should differ")
	}
}

func TestGradientAtClamping(t *testing.T) {
	g := GradientProgress
	// Should not panic on out-of-range
	_ = g.At(-10)
	_ = g.At(200)
}

func TestGradientBar(t *testing.T) {
	bar := GradientBar(50, 10, GradientProgress)
	if bar == "" {
		t.Error("gradient bar should not be empty")
	}
	// At 50%, should have both filled and empty blocks
	stripped := lipgloss.NewStyle().Render(bar)
	if !strings.Contains(stripped, "█") {
		t.Error("gradient bar should contain filled blocks")
	}
}

func TestGradientBarEdgeCases(t *testing.T) {
	if got := GradientBar(0, 10, GradientProgress); got == "" {
		t.Error("0% bar should still render empty blocks")
	}
	if got := GradientBar(100, 10, GradientProgress); got == "" {
		t.Error("100% bar should render")
	}
	if got := GradientBar(50, 0, GradientProgress); got != "" {
		t.Error("zero-width bar should be empty")
	}
}
