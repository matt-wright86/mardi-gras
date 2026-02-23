package gastown

import (
	"strings"

	"github.com/matt-wright86/mardi-gras/internal/data"
)

// FormulaRecommendation pairs a formula name with a reason for suggesting it.
type FormulaRecommendation struct {
	Formula string
	Reason  string
	Score   int // higher = stronger match
}

// RecommendFormulas suggests formulas for an issue based on its labels, type,
// and title. Returns recommendations sorted by score (best first).
// This is a heuristic starting point — evolve toward history-based
// recommendation as HOP data becomes available.
func RecommendFormulas(issue data.Issue) []FormulaRecommendation {
	var recs []FormulaRecommendation

	labels := make(map[string]bool)
	for _, l := range issue.Labels {
		labels[strings.ToLower(l)] = true
	}

	titleLower := strings.ToLower(issue.Title)
	descLower := strings.ToLower(issue.Description)

	// Security-related → security-audit
	if labels["security"] || labels["sec"] || labels["vulnerability"] ||
		strings.Contains(titleLower, "security") ||
		strings.Contains(titleLower, "vulnerability") ||
		strings.Contains(titleLower, "auth") ||
		strings.Contains(descLower, "security audit") {
		recs = append(recs, FormulaRecommendation{
			Formula: "security-audit",
			Reason:  "Security-related issue benefits from multi-aspect analysis",
			Score:   90,
		})
	}

	// Code review label → code-review
	if labels["review"] || labels["code-review"] ||
		strings.Contains(titleLower, "review") {
		recs = append(recs, FormulaRecommendation{
			Formula: "code-review",
			Reason:  "Issue tagged for review or mentions review in title",
			Score:   85,
		})
	}

	// Feature type or complex label → shiny (full lifecycle)
	if issue.IssueType == data.TypeFeature || labels["complex"] || labels["epic"] ||
		strings.Contains(titleLower, "implement") ||
		strings.Contains(titleLower, "build") ||
		strings.Contains(titleLower, "add") {
		recs = append(recs, FormulaRecommendation{
			Formula: "shiny",
			Reason:  "Feature work benefits from full design→implement→review→test→submit",
			Score:   80,
		})
	}

	// High-priority or critical → rule-of-five (parallel review)
	if issue.Priority <= data.PriorityHigh || labels["critical"] {
		recs = append(recs, FormulaRecommendation{
			Formula: "rule-of-five",
			Reason:  "High-priority issue benefits from 5x parallel review",
			Score:   75,
		})
	}

	// Bug type → mol-polecat-work (standard lifecycle)
	if issue.IssueType == data.TypeBug {
		recs = append(recs, FormulaRecommendation{
			Formula: "mol-polecat-work",
			Reason:  "Standard polecat lifecycle for bug fixes",
			Score:   70,
		})
	}

	// Task type or chore → mol-polecat-work
	if issue.IssueType == data.TypeTask || issue.IssueType == data.TypeChore {
		recs = append(recs, FormulaRecommendation{
			Formula: "mol-polecat-work",
			Reason:  "Standard polecat lifecycle for tasks",
			Score:   60,
		})
	}

	// Default: always suggest mol-polecat-work as a fallback if not already present
	hasPolecat := false
	for _, r := range recs {
		if r.Formula == "mol-polecat-work" {
			hasPolecat = true
			break
		}
	}
	if !hasPolecat {
		recs = append(recs, FormulaRecommendation{
			Formula: "mol-polecat-work",
			Reason:  "Default polecat workflow",
			Score:   30,
		})
	}

	// Sort by score descending
	sortRecommendations(recs)
	return recs
}

// sortRecommendations sorts by score descending (stable for equal scores).
func sortRecommendations(recs []FormulaRecommendation) {
	for i := 1; i < len(recs); i++ {
		for j := i; j > 0 && recs[j].Score > recs[j-1].Score; j-- {
			recs[j], recs[j-1] = recs[j-1], recs[j]
		}
	}
}
