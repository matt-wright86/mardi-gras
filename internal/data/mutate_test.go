package data

import "testing"

func TestBranchName(t *testing.T) {
	tests := []struct {
		name     string
		issue    Issue
		expected string
	}{
		{
			name:     "Bug issue",
			issue:    Issue{ID: "bd-a1b2", Title: "Fix login token expiry", IssueType: TypeBug},
			expected: "fix/bd-a1b2-fix-login-token-expiry",
		},
		{
			name:     "Feature issue",
			issue:    Issue{ID: "bd-c3d4", Title: "Add search feature", IssueType: TypeFeature},
			expected: "feat/bd-c3d4-add-search-feature",
		},
		{
			name:     "Task issue",
			issue:    Issue{ID: "bd-e5f6", Title: "Update documentation", IssueType: TypeTask},
			expected: "task/bd-e5f6-update-documentation",
		},
		{
			name:     "Chore issue",
			issue:    Issue{ID: "bd-g7h8", Title: "Clean up CI config", IssueType: TypeChore},
			expected: "chore/bd-g7h8-clean-up-ci-config",
		},
		{
			name:     "Special characters stripped",
			issue:    Issue{ID: "bd-i9j0", Title: "Handle @mentions & #tags (v2)", IssueType: TypeFeature},
			expected: "feat/bd-i9j0-handle-mentions-tags-v2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BranchName(tt.issue)
			if got != tt.expected {
				t.Errorf("BranchName() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"Fix login/auth bug", "fix-login-auth-bug"},
		{"UPPER CASE", "upper-case"},
		{"   spaces   ", "spaces"},
		{"no-change", "no-change"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := slugify(tt.input)
			if got != tt.expected {
				t.Errorf("slugify(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
