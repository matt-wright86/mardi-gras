package ui

import (
	"strings"
	"testing"
)

func TestRenderStars(t *testing.T) {
	// Full stars for excellent quality
	result := RenderStars(1.0)
	if !strings.Contains(result, SymStar) {
		t.Error("expected filled stars in output")
	}

	// Zero stars for zero quality
	result = RenderStars(0.0)
	if !strings.Contains(result, SymStarEmpty) {
		t.Error("expected empty stars in output")
	}
}

func TestRenderStarsCompact(t *testing.T) {
	result := RenderStarsCompact(0.8)
	if !strings.Contains(result, SymStar) {
		t.Error("expected star symbol in compact output")
	}
	if !strings.Contains(result, "4") {
		t.Error("expected '4' in compact output for 0.8 score")
	}
}

func TestRenderStarsCompactClampsEdges(t *testing.T) {
	// Score 0.0 → ★0
	low := RenderStarsCompact(0.0)
	if !strings.Contains(low, "0") {
		t.Error("expected '0' for zero score")
	}

	// Score 1.0 → ★5
	high := RenderStarsCompact(1.0)
	if !strings.Contains(high, "5") {
		t.Error("expected '5' for max score")
	}
}
