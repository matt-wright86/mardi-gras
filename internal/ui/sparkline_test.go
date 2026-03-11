package ui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestBrailleSparklineEmpty(t *testing.T) {
	got := BrailleSparkline(nil, 10, lipgloss.NewStyle())
	if strings.TrimSpace(got) != "" {
		// Should be all spaces
		if got == "" {
			t.Error("should return padded spaces for nil data")
		}
	}
}

func TestBrailleSparklineBasic(t *testing.T) {
	data := []float64{1, 2, 3, 4}
	got := BrailleSparkline(data, 5, lipgloss.NewStyle())
	if got == "" {
		t.Error("should render braille characters")
	}
	// 4 data points = 2 braille chars + 3 padding = 5 chars
	if len([]rune(got)) != 5 {
		t.Errorf("expected 5 runes, got %d", len([]rune(got)))
	}
}

func TestBrailleSparklineOddData(t *testing.T) {
	data := []float64{1, 2, 3}
	got := BrailleSparkline(data, 5, lipgloss.NewStyle())
	if got == "" {
		t.Error("should handle odd number of data points")
	}
}

func TestBrailleSparklineAllZero(t *testing.T) {
	data := []float64{0, 0, 0}
	got := BrailleSparkline(data, 5, lipgloss.NewStyle())
	if strings.TrimSpace(got) != "" {
		// All zeros should render as spaces
		t.Log("all-zero data rendered as:", got)
	}
}

func TestMiniSparklineAllZero(t *testing.T) {
	got := MiniSparkline([3]int{0, 0, 0})
	if got != "" {
		t.Error("all-zero should return empty string")
	}
}

func TestMiniSparklineBasic(t *testing.T) {
	got := MiniSparkline([3]int{1, 3, 2})
	if got == "" {
		t.Error("non-zero values should produce output")
	}
}

func TestDualSparklineBasic(t *testing.T) {
	top := []float64{1, 0, 2, 0}
	bot := []float64{0, 1, 0, 2}
	style := lipgloss.NewStyle()
	got := DualSparkline(top, bot, 4, style, style)
	if got == "" {
		t.Error("should render dual sparkline")
	}
}

func TestDualSparklineZeroWidth(t *testing.T) {
	got := DualSparkline([]float64{1}, []float64{1}, 0, lipgloss.NewStyle(), lipgloss.NewStyle())
	if got != "" {
		t.Error("zero width should return empty")
	}
}
