package gastown

import (
	"testing"

	"github.com/matt-wright86/mardi-gras/internal/data"
)

func TestRecommendFormulasSecurityLabel(t *testing.T) {
	issue := data.Issue{
		ID:        "bd-001",
		Title:     "Fix token validation",
		IssueType: data.TypeBug,
		Labels:    []string{"security"},
	}
	recs := RecommendFormulas(issue)
	if len(recs) == 0 {
		t.Fatal("expected recommendations")
	}
	if recs[0].Formula != "security-audit" {
		t.Fatalf("expected security-audit first, got %q", recs[0].Formula)
	}
}

func TestRecommendFormulasSecurityTitle(t *testing.T) {
	issue := data.Issue{
		ID:        "bd-002",
		Title:     "Add authentication middleware",
		IssueType: data.TypeFeature,
	}
	recs := RecommendFormulas(issue)

	hasSecurityAudit := false
	for _, r := range recs {
		if r.Formula == "security-audit" {
			hasSecurityAudit = true
			break
		}
	}
	if !hasSecurityAudit {
		t.Fatal("expected security-audit recommendation for auth-related issue")
	}
}

func TestRecommendFormulasFeature(t *testing.T) {
	issue := data.Issue{
		ID:        "bd-003",
		Title:     "New dashboard widget",
		IssueType: data.TypeFeature,
		Priority:  data.PriorityMedium,
	}
	recs := RecommendFormulas(issue)

	hasShiny := false
	for _, r := range recs {
		if r.Formula == "shiny" {
			hasShiny = true
			break
		}
	}
	if !hasShiny {
		t.Fatal("expected shiny recommendation for feature type")
	}
}

func TestRecommendFormulasHighPriority(t *testing.T) {
	issue := data.Issue{
		ID:        "bd-004",
		Title:     "Critical production outage",
		IssueType: data.TypeBug,
		Priority:  data.PriorityCritical,
	}
	recs := RecommendFormulas(issue)

	hasRuleOfFive := false
	for _, r := range recs {
		if r.Formula == "rule-of-five" {
			hasRuleOfFive = true
			break
		}
	}
	if !hasRuleOfFive {
		t.Fatal("expected rule-of-five for high-priority issue")
	}
}

func TestRecommendFormulasDefaultFallback(t *testing.T) {
	issue := data.Issue{
		ID:        "bd-005",
		Title:     "Update docs",
		IssueType: data.TypeChore,
		Priority:  data.PriorityBacklog,
	}
	recs := RecommendFormulas(issue)
	if len(recs) == 0 {
		t.Fatal("expected at least one recommendation")
	}

	// Should always have mol-polecat-work somewhere
	hasDefault := false
	for _, r := range recs {
		if r.Formula == "mol-polecat-work" {
			hasDefault = true
			break
		}
	}
	if !hasDefault {
		t.Fatal("expected mol-polecat-work as fallback")
	}
}

func TestRecommendFormulasSortedByScore(t *testing.T) {
	// Security feature with high priority → multiple recommendations
	issue := data.Issue{
		ID:        "bd-006",
		Title:     "Implement security review for auth module",
		IssueType: data.TypeFeature,
		Priority:  data.PriorityHigh,
		Labels:    []string{"security", "review"},
	}
	recs := RecommendFormulas(issue)
	if len(recs) < 2 {
		t.Fatalf("expected multiple recommendations, got %d", len(recs))
	}

	// Verify descending score order
	for i := 1; i < len(recs); i++ {
		if recs[i].Score > recs[i-1].Score {
			t.Fatalf("recommendations not sorted: %d (score %d) after %d (score %d)",
				i, recs[i].Score, i-1, recs[i-1].Score)
		}
	}
}

func TestRecommendFormulasReviewLabel(t *testing.T) {
	issue := data.Issue{
		ID:        "bd-007",
		Title:     "PR needs code review",
		IssueType: data.TypeTask,
	}
	recs := RecommendFormulas(issue)

	hasCodeReview := false
	for _, r := range recs {
		if r.Formula == "code-review" {
			hasCodeReview = true
			break
		}
	}
	if !hasCodeReview {
		t.Fatal("expected code-review for issue with 'review' in title")
	}
}
