package gastown

import (
	"testing"
	"time"

	"github.com/matt-wright86/mardi-gras/internal/data"
)

func TestPredictConvoys(t *testing.T) {
	convoys := []ConvoyDetail{
		{ID: "cv-1", Title: "Auth sprint", Status: "open", Completed: 3, Total: 10},
		{ID: "cv-2", Title: "UI cleanup", Status: "open", Completed: 1, Total: 4},
	}
	velocity := &VelocityMetrics{
		ClosedToday:   2,
		ClosedWeek:    10,
		WorkingAgents: 3,
	}

	preds := PredictConvoys(convoys, velocity)
	if len(preds) != 2 {
		t.Fatalf("expected 2 predictions, got %d", len(preds))
	}

	// cv-1: 7 remaining
	if preds[0].Remaining != 7 {
		t.Fatalf("cv-1 remaining = %d, want 7", preds[0].Remaining)
	}
	if preds[0].ETALabel == "" || preds[0].ETALabel == "unknown" {
		t.Fatalf("cv-1 should have valid ETA, got %q", preds[0].ETALabel)
	}
	if preds[0].Confidence != "high" {
		t.Fatalf("cv-1 confidence = %q, want 'high'", preds[0].Confidence)
	}

	// cv-2: 3 remaining
	if preds[1].Remaining != 3 {
		t.Fatalf("cv-2 remaining = %d, want 3", preds[1].Remaining)
	}
}

func TestPredictConvoysCompletedSkipped(t *testing.T) {
	convoys := []ConvoyDetail{
		{ID: "cv-1", Title: "Done", Status: "open", Completed: 5, Total: 5},
		{ID: "cv-2", Title: "Closed", Status: "closed", Completed: 3, Total: 5},
	}
	velocity := &VelocityMetrics{ClosedWeek: 5}

	preds := PredictConvoys(convoys, velocity)
	if len(preds) != 0 {
		t.Fatalf("expected 0 predictions (completed/closed convoys), got %d", len(preds))
	}
}

func TestPredictConvoysNilVelocity(t *testing.T) {
	convoys := []ConvoyDetail{
		{ID: "cv-1", Title: "Active", Status: "open", Completed: 1, Total: 5},
	}
	preds := PredictConvoys(convoys, nil)
	if preds != nil {
		t.Fatal("expected nil predictions with nil velocity")
	}
}

func TestPredictConvoysZeroRate(t *testing.T) {
	convoys := []ConvoyDetail{
		{ID: "cv-1", Title: "Stalled", Status: "open", Completed: 0, Total: 5},
	}
	velocity := &VelocityMetrics{ClosedToday: 0, ClosedWeek: 0}

	preds := PredictConvoys(convoys, velocity)
	if len(preds) != 1 {
		t.Fatalf("expected 1 prediction, got %d", len(preds))
	}
	if preds[0].ETALabel != "unknown" {
		t.Fatalf("expected 'unknown' ETA, got %q", preds[0].ETALabel)
	}
	if preds[0].Confidence != "low" {
		t.Fatalf("expected 'low' confidence, got %q", preds[0].Confidence)
	}
}

func TestPredictCostBudget(t *testing.T) {
	now := time.Date(2026, 2, 23, 18, 0, 0, 0, time.UTC) // 6pm = 18h elapsed
	costs := &CostsOutput{
		Total:    CostTotal{Cost: 30.0},
		Sessions: 10,
	}

	result := PredictCostBudget(costs, 50.0, now)
	if result == "" {
		t.Fatal("expected non-empty budget prediction")
	}
	if result == "budget exceeded" {
		t.Fatalf("$30 of $50 at 6pm should not be exceeded, got %q", result)
	}
}

func TestPredictCostBudgetExceeded(t *testing.T) {
	now := time.Date(2026, 2, 23, 12, 0, 0, 0, time.UTC)
	costs := &CostsOutput{
		Total:    CostTotal{Cost: 60.0},
		Sessions: 15,
	}

	result := PredictCostBudget(costs, 50.0, now)
	if result != "budget exceeded" {
		t.Fatalf("expected 'budget exceeded', got %q", result)
	}
}

func TestPredictCostBudgetNilCosts(t *testing.T) {
	result := PredictCostBudget(nil, 50.0, time.Now())
	if result != "" {
		t.Fatalf("expected empty, got %q", result)
	}
}

func TestAgentProductivity(t *testing.T) {
	now := time.Now()
	closedAt := now.Add(-1 * time.Hour)
	issues := []data.Issue{
		{ID: "a", Status: data.StatusClosed, ClosedAt: &closedAt},
		{ID: "b", Status: data.StatusClosed, ClosedAt: &closedAt},
		{ID: "c", Status: data.StatusOpen},
	}
	costs := &CostsOutput{Sessions: 4}

	prod := AgentProductivity(issues, costs)
	if prod != 0.5 { // 2 closed / 4 sessions
		t.Fatalf("productivity = %f, want 0.5", prod)
	}
}

func TestAgentProductivityNoSessions(t *testing.T) {
	prod := AgentProductivity(nil, nil)
	if prod != 0 {
		t.Fatalf("expected 0, got %f", prod)
	}
}

func TestFormatETA(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{30 * time.Minute, "30m"},
		{2*time.Hour + 30*time.Minute, "2.5h"},
		{36 * time.Hour, "1.5d"},
		{14 * 24 * time.Hour, "2w"},
	}
	for _, tt := range tests {
		got := formatETA(tt.d)
		if got != tt.want {
			t.Errorf("formatETA(%v) = %q, want %q", tt.d, got, tt.want)
		}
	}
}
