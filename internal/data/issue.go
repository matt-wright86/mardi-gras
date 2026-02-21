package data

import (
	"fmt"
	"time"
)

// Status represents the state of a Beads issue.
type Status string

const (
	StatusOpen       Status = "open"
	StatusInProgress Status = "in_progress"
	StatusClosed     Status = "closed"
)

// IssueType represents the category of a Beads issue.
type IssueType string

const (
	TypeTask    IssueType = "task"
	TypeBug     IssueType = "bug"
	TypeFeature IssueType = "feature"
	TypeChore   IssueType = "chore"
	TypeEpic    IssueType = "epic"
)

// Priority ranges from 0 (critical) to 4 (backlog).
type Priority int

const (
	PriorityCritical Priority = 0
	PriorityHigh     Priority = 1
	PriorityMedium   Priority = 2
	PriorityLow      Priority = 3
	PriorityBacklog  Priority = 4
)

// ParadeStatus maps issues to their parade float group.
type ParadeStatus int

const (
	ParadeRolling    ParadeStatus = iota // in_progress
	ParadeLinedUp                        // open, not blocked
	ParadeStalled                        // open, blocked
	ParadePastTheStand                   // closed
)

// Dependency represents a relationship between two issues.
type Dependency struct {
	IssueID     string `json:"issue_id"`
	DependsOnID string `json:"depends_on_id"`
	Type        string `json:"type"`
	CreatedAt   string `json:"created_at"`
	CreatedBy   string `json:"created_by"`
}

// Issue represents a single Beads issue.
type Issue struct {
	ID                 string       `json:"id"`
	Title              string       `json:"title"`
	Description        string       `json:"description,omitempty"`
	Status             Status       `json:"status"`
	Priority           Priority     `json:"priority"`
	IssueType          IssueType    `json:"issue_type"`
	Owner              string       `json:"owner,omitempty"`
	Assignee           string       `json:"assignee,omitempty"`
	CreatedAt          time.Time    `json:"created_at"`
	CreatedBy          string       `json:"created_by,omitempty"`
	UpdatedAt          time.Time    `json:"updated_at"`
	ClosedAt           *time.Time   `json:"closed_at,omitempty"`
	CloseReason        string       `json:"close_reason,omitempty"`
	Dependencies       []Dependency `json:"dependencies,omitempty"`
	Notes              string       `json:"notes,omitempty"`
	Design             string       `json:"design,omitempty"`
	AcceptanceCriteria string       `json:"acceptance_criteria,omitempty"`
}

// IsBlocked returns true if this issue depends on an unclosed blocker.
func (i *Issue) IsBlocked(issueMap map[string]*Issue) bool {
	for _, dep := range i.Dependencies {
		if dep.Type != "blocks" {
			continue
		}
		if blocker, ok := issueMap[dep.DependsOnID]; ok {
			if blocker.Status != StatusClosed {
				return true
			}
		}
	}
	return false
}

// BlockedByIDs returns the IDs of open issues blocking this one.
func (i *Issue) BlockedByIDs(issueMap map[string]*Issue) []string {
	var blockers []string
	for _, dep := range i.Dependencies {
		if dep.Type != "blocks" {
			continue
		}
		if blocker, ok := issueMap[dep.DependsOnID]; ok {
			if blocker.Status != StatusClosed {
				blockers = append(blockers, dep.DependsOnID)
			}
		}
	}
	return blockers
}

// BlocksIDs returns the IDs of issues that this issue blocks.
func (i *Issue) BlocksIDs(allIssues []Issue) []string {
	var blocked []string
	for _, other := range allIssues {
		for _, dep := range other.Dependencies {
			if dep.Type == "blocks" && dep.DependsOnID == i.ID {
				blocked = append(blocked, other.ID)
			}
		}
	}
	return blocked
}

// Age returns the duration since the issue was created.
func (i *Issue) Age() time.Duration {
	return time.Since(i.CreatedAt)
}

// AgeLabel returns a human-readable age string.
func (i *Issue) AgeLabel() string {
	days := int(i.Age().Hours() / 24)
	switch {
	case days == 0:
		hours := int(i.Age().Hours())
		if hours == 0 {
			return "just now"
		}
		return fmt.Sprintf("%dh", hours)
	case days == 1:
		return "1 day"
	case days < 30:
		return fmt.Sprintf("%d days", days)
	default:
		return fmt.Sprintf("%d weeks", days/7)
	}
}

// ParadeGroup determines which parade section this issue belongs to.
func (i *Issue) ParadeGroup(issueMap map[string]*Issue) ParadeStatus {
	switch i.Status {
	case StatusInProgress:
		return ParadeRolling
	case StatusClosed:
		return ParadePastTheStand
	case StatusOpen:
		if i.IsBlocked(issueMap) {
			return ParadeStalled
		}
		return ParadeLinedUp
	default:
		return ParadeLinedUp
	}
}

// PriorityLabel returns "P0" through "P4".
func PriorityLabel(p Priority) string {
	return fmt.Sprintf("P%d", p)
}

// PriorityName returns the full name for a priority level.
func PriorityName(p Priority) string {
	switch p {
	case PriorityCritical:
		return "Critical"
	case PriorityHigh:
		return "High"
	case PriorityMedium:
		return "Medium"
	case PriorityLow:
		return "Low"
	case PriorityBacklog:
		return "Backlog"
	default:
		return "Unknown"
	}
}

// BuildIssueMap creates a lookup map from a slice of issues.
func BuildIssueMap(issues []Issue) map[string]*Issue {
	m := make(map[string]*Issue, len(issues))
	for idx := range issues {
		m[issues[idx].ID] = &issues[idx]
	}
	return m
}
