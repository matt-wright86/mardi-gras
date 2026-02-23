package data

import "time"

// EntityRef is a structured attribution reference for the Hierarchy of Proof.
// It supports URIs like hop://gastown/steveyegge/polecat-nux.
type EntityRef struct {
	Name     string `json:"name"`
	Platform string `json:"platform,omitempty"` // "gastown", "github", "human"
	Org      string `json:"org,omitempty"`
	ID       string `json:"id,omitempty"`
	URI      string `json:"uri,omitempty"` // hop://gastown/org/agent
}

// ValidationOutcome is the result of a validation review.
type ValidationOutcome string

const (
	OutcomeAccepted ValidationOutcome = "accepted"
	OutcomeRejected ValidationOutcome = "rejected"
	OutcomeRevision ValidationOutcome = "revision_requested"
)

// Validation records a review of an issue's work by a validator (human or agent).
type Validation struct {
	Validator    EntityRef         `json:"validator"`
	Outcome      ValidationOutcome `json:"outcome"`
	QualityScore float32           `json:"quality_score"` // 0.0-1.0
	Comment      string            `json:"comment,omitempty"`
	Timestamp    time.Time         `json:"timestamp"`
}

// QualityStars converts a 0.0-1.0 quality score to a 0-5 star rating.
func QualityStars(score float32) int {
	stars := int(score*5 + 0.5) // round to nearest
	if stars < 0 {
		return 0
	}
	if stars > 5 {
		return 5
	}
	return stars
}

// QualityLabel returns a human-readable label for a quality score.
func QualityLabel(score float32) string {
	switch {
	case score >= 0.9:
		return "excellent"
	case score >= 0.7:
		return "good"
	case score >= 0.5:
		return "fair"
	case score >= 0.3:
		return "poor"
	default:
		return "low"
	}
}
