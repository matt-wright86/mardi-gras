package data

import (
	"encoding/json"
	"testing"
	"time"
)

func TestQualityStars(t *testing.T) {
	tests := []struct {
		score float32
		want  int
	}{
		{0.0, 0},
		{0.1, 1},
		{0.3, 2},
		{0.5, 3},
		{0.7, 4},
		{0.9, 5},
		{1.0, 5},
		{0.89, 4}, // 4.45 rounds to 4
		{0.91, 5}, // 4.55 rounds to 5
	}
	for _, tt := range tests {
		got := QualityStars(tt.score)
		if got != tt.want {
			t.Errorf("QualityStars(%.2f) = %d, want %d", tt.score, got, tt.want)
		}
	}
}

func TestQualityStarsClamps(t *testing.T) {
	if QualityStars(-0.5) != 0 {
		t.Fatal("negative score should clamp to 0")
	}
	if QualityStars(2.0) != 5 {
		t.Fatal("score > 1.0 should clamp to 5")
	}
}

func TestIssueHOPFieldsParsing(t *testing.T) {
	jsonl := `{"id":"bd-001","title":"Test HOP","status":"closed","priority":2,"issue_type":"task","created_at":"2026-02-20T10:00:00Z","updated_at":"2026-02-23T10:00:00Z","quality_score":0.85,"crystallizes":true,"creator":{"name":"polecat-alpha","platform":"gastown","uri":"hop://gastown/mardi_gras/polecat-alpha"},"validations":[{"validator":{"name":"witness","platform":"gastown"},"outcome":"accepted","quality_score":0.9,"timestamp":"2026-02-22T14:00:00Z"},{"validator":{"name":"refinery","platform":"gastown"},"outcome":"accepted","quality_score":0.8,"comment":"Clean implementation","timestamp":"2026-02-22T15:00:00Z"}]}`

	var issue Issue
	if err := json.Unmarshal([]byte(jsonl), &issue); err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if issue.QualityScore == nil || *issue.QualityScore != 0.85 {
		t.Fatalf("QualityScore = %v, want 0.85", issue.QualityScore)
	}

	if issue.Crystallizes == nil || !*issue.Crystallizes {
		t.Fatalf("Crystallizes = %v, want true", issue.Crystallizes)
	}

	if issue.Creator == nil {
		t.Fatal("Creator should not be nil")
	}
	if issue.Creator.Name != "polecat-alpha" {
		t.Fatalf("Creator.Name = %q, want 'polecat-alpha'", issue.Creator.Name)
	}
	if issue.Creator.URI != "hop://gastown/mardi_gras/polecat-alpha" {
		t.Fatalf("Creator.URI = %q", issue.Creator.URI)
	}

	if len(issue.Validations) != 2 {
		t.Fatalf("len(Validations) = %d, want 2", len(issue.Validations))
	}

	v := issue.Validations[0]
	if v.Validator.Name != "witness" {
		t.Fatalf("Validations[0].Validator.Name = %q", v.Validator.Name)
	}
	if v.Outcome != OutcomeAccepted {
		t.Fatalf("Validations[0].Outcome = %q, want 'accepted'", v.Outcome)
	}
	if v.QualityScore != 0.9 {
		t.Fatalf("Validations[0].QualityScore = %f, want 0.9", v.QualityScore)
	}

	v1 := issue.Validations[1]
	if v1.Comment != "Clean implementation" {
		t.Fatalf("Validations[1].Comment = %q", v1.Comment)
	}
}

func TestIssueHOPFieldsOmitted(t *testing.T) {
	jsonl := `{"id":"bd-002","title":"No HOP","status":"open","priority":2,"issue_type":"task","created_at":"2026-02-20T10:00:00Z","updated_at":"2026-02-23T10:00:00Z"}`

	var issue Issue
	if err := json.Unmarshal([]byte(jsonl), &issue); err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if issue.QualityScore != nil {
		t.Fatalf("QualityScore should be nil, got %v", *issue.QualityScore)
	}
	if issue.Crystallizes != nil {
		t.Fatal("Crystallizes should be nil")
	}
	if issue.Creator != nil {
		t.Fatal("Creator should be nil")
	}
	if len(issue.Validations) != 0 {
		t.Fatalf("Validations should be empty, got %d", len(issue.Validations))
	}
}

func TestValidationOutcomes(t *testing.T) {
	v := Validation{
		Validator:    EntityRef{Name: "test"},
		Outcome:      OutcomeRevision,
		QualityScore: 0.5,
		Timestamp:    time.Now(),
	}
	if v.Outcome != "revision_requested" {
		t.Fatalf("OutcomeRevision = %q, want 'revision_requested'", v.Outcome)
	}
}

func TestQualityLabel(t *testing.T) {
	tests := []struct {
		score float32
		want  string
	}{
		{0.95, "excellent"},
		{0.75, "good"},
		{0.55, "fair"},
		{0.35, "poor"},
		{0.1, "low"},
	}
	for _, tt := range tests {
		got := QualityLabel(tt.score)
		if got != tt.want {
			t.Errorf("QualityLabel(%.2f) = %q, want %q", tt.score, got, tt.want)
		}
	}
}
