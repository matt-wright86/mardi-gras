package gastown

import (
	"testing"

	"github.com/matt-wright86/mardi-gras/internal/data"
)

func TestComputeScorecards(t *testing.T) {
	q85 := float32(0.85)
	q70 := float32(0.70)
	q90 := float32(0.90)
	crystTrue := true
	crystFalse := false

	issues := []data.Issue{
		{
			ID: "a", Status: data.StatusClosed, Assignee: "Toast",
			QualityScore: &q85, Crystallizes: &crystTrue,
		},
		{
			ID: "b", Status: data.StatusClosed, Assignee: "Toast",
			QualityScore: &q90, Crystallizes: &crystTrue,
		},
		{
			ID: "c", Status: data.StatusClosed, Assignee: "Muffin",
			QualityScore: &q70, Crystallizes: &crystFalse,
		},
		// Open issue — should not count
		{
			ID: "d", Status: data.StatusOpen, Assignee: "Toast",
			QualityScore: &q90,
		},
		// Closed but no assignee, has creator
		{
			ID: "e", Status: data.StatusClosed,
			Creator:      &data.EntityRef{Name: "Whiskers"},
			QualityScore: &q85,
		},
		// Closed with no quality score
		{
			ID: "f", Status: data.StatusClosed, Assignee: "Toast",
		},
	}

	cards := ComputeScorecards(issues)

	if len(cards) != 3 {
		t.Fatalf("expected 3 scorecards, got %d", len(cards))
	}

	// Should be sorted by quality descending
	// Toast: (0.85+0.90)/2 = 0.875, Whiskers: 0.85, Muffin: 0.70
	if cards[0].Name != "Toast" {
		t.Fatalf("expected Toast first (highest quality), got %q", cards[0].Name)
	}
	if cards[0].IssuesClosed != 3 {
		t.Fatalf("Toast IssuesClosed = %d, want 3", cards[0].IssuesClosed)
	}
	if cards[0].TotalScored != 2 {
		t.Fatalf("Toast TotalScored = %d, want 2", cards[0].TotalScored)
	}
	if cards[0].AvgQuality < 0.87 || cards[0].AvgQuality > 0.88 {
		t.Fatalf("Toast AvgQuality = %.3f, want ~0.875", cards[0].AvgQuality)
	}
	if cards[0].Crystallizing != 2 {
		t.Fatalf("Toast Crystallizing = %d, want 2", cards[0].Crystallizing)
	}

	if cards[1].Name != "Whiskers" {
		t.Fatalf("expected Whiskers second, got %q", cards[1].Name)
	}

	if cards[2].Name != "Muffin" {
		t.Fatalf("expected Muffin third, got %q", cards[2].Name)
	}
	if cards[2].Ephemeral != 1 {
		t.Fatalf("Muffin Ephemeral = %d, want 1", cards[2].Ephemeral)
	}
}

func TestComputeScorecardsEmpty(t *testing.T) {
	cards := ComputeScorecards(nil)
	if len(cards) != 0 {
		t.Fatalf("expected 0 scorecards, got %d", len(cards))
	}
}

func TestComputeScorecardsNoClosedIssues(t *testing.T) {
	issues := []data.Issue{
		{ID: "a", Status: data.StatusOpen, Assignee: "Toast"},
		{ID: "b", Status: data.StatusInProgress, Assignee: "Muffin"},
	}
	cards := ComputeScorecards(issues)
	if len(cards) != 0 {
		t.Fatalf("expected 0 scorecards for no closed issues, got %d", len(cards))
	}
}
