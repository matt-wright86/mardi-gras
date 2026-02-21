package data

import (
	"testing"
)

func TestFilterIssues(t *testing.T) {
	issues := []Issue{
		{ID: "vv-001", Title: "Fix login bug", IssueType: TypeBug, Priority: PriorityCritical},
		{ID: "vv-002", Title: "Add search feature", IssueType: TypeFeature, Priority: PriorityHigh},
		{ID: "vv-003", Title: "Update documentation", IssueType: TypeChore, Priority: PriorityLow},
		{ID: "vv-004", Title: "Refactor auth flow", IssueType: TypeTask, Priority: PriorityMedium},
	}

	tests := []struct {
		name     string
		query    string
		expected []string // expected issue IDs
	}{
		{
			name:     "Empty query",
			query:    "",
			expected: []string{"vv-001", "vv-002", "vv-003", "vv-004"},
		},
		{
			name:     "Free text partial word",
			query:    "log",
			expected: []string{"vv-001"},
		},
		{
			name:     "Free text multiple words",
			query:    "fix bug",
			expected: []string{"vv-001"},
		},
		{
			name:     "Free text case insensitive",
			query:    "SEARCH",
			expected: []string{"vv-002"},
		},
		{
			name:     "Free text matches ID",
			query:    "vv-003",
			expected: []string{"vv-003"},
		},
		{
			name:     "Type filter exact",
			query:    "type:bug",
			expected: []string{"vv-001"},
		},
		{
			name:     "Type filter and free text",
			query:    "type:feature search",
			expected: []string{"vv-002"},
		},
		{
			name:     "Type filter no match",
			query:    "type:epic",
			expected: []string{},
		},
		{
			name:     "Priority short label (p0)",
			query:    "p0",
			expected: []string{"vv-001"},
		},
		{
			name:     "Priority short label (p1)",
			query:    "P1",
			expected: []string{"vv-002"},
		},
		{
			name:     "Priority explicit number",
			query:    "priority:3",
			expected: []string{"vv-003"},
		},
		{
			name:     "Priority explicit name",
			query:    "priority:medium",
			expected: []string{"vv-004"},
		},
		{
			name:     "Multiple tokens combined",
			query:    "type:feature p1 add",
			expected: []string{"vv-002"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterIssues(issues, tt.query)
			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d results, got %d", len(tt.expected), len(result))
			}

			// Verify IDs
			resMap := make(map[string]bool)
			for _, r := range result {
				resMap[r.ID] = true
			}
			for _, exp := range tt.expected {
				if !resMap[exp] {
					t.Errorf("expected issue %s to be in result, but it was not", exp)
				}
			}
		})
	}
}
