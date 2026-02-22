package data

import (
	"testing"
	"time"
)

func TestParentID(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{"mg-007", ""},
		{"mg-007.1", "mg-007"},
		{"mg-007.2.1", "mg-007.2"},
		{"bd-a3f8.1.1", "bd-a3f8.1"},
		{"simple", ""},
	}
	for _, tc := range tests {
		t.Run(tc.id, func(t *testing.T) {
			iss := Issue{ID: tc.id}
			if got := iss.ParentID(); got != tc.want {
				t.Errorf("ParentID(%q) = %q, want %q", tc.id, got, tc.want)
			}
		})
	}
}

func TestNestingDepth(t *testing.T) {
	tests := []struct {
		id   string
		want int
	}{
		{"mg-007", 0},
		{"mg-007.1", 1},
		{"mg-007.2.1", 2},
		{"a.b.c.d", 3},
	}
	for _, tc := range tests {
		t.Run(tc.id, func(t *testing.T) {
			iss := Issue{ID: tc.id}
			if got := iss.NestingDepth(); got != tc.want {
				t.Errorf("NestingDepth(%q) = %d, want %d", tc.id, got, tc.want)
			}
		})
	}
}

func TestIsOverdue(t *testing.T) {
	past := time.Now().Add(-48 * time.Hour)
	future := time.Now().Add(48 * time.Hour)

	tests := []struct {
		name   string
		dueAt  *time.Time
		status Status
		want   bool
	}{
		{"nil due", nil, StatusOpen, false},
		{"past due, open", &past, StatusOpen, true},
		{"past due, closed", &past, StatusClosed, false},
		{"future due, open", &future, StatusOpen, false},
		{"past due, in_progress", &past, StatusInProgress, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			iss := Issue{DueAt: tc.dueAt, Status: tc.status}
			if got := iss.IsOverdue(); got != tc.want {
				t.Errorf("IsOverdue() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIsDeferred(t *testing.T) {
	past := time.Now().Add(-24 * time.Hour)
	future := time.Now().Add(5 * 24 * time.Hour)

	tests := []struct {
		name       string
		deferUntil *time.Time
		want       bool
	}{
		{"nil", nil, false},
		{"past", &past, false},
		{"future", &future, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			iss := Issue{DeferUntil: tc.deferUntil}
			if got := iss.IsDeferred(); got != tc.want {
				t.Errorf("IsDeferred() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDueLabel(t *testing.T) {
	tests := []struct {
		name   string
		offset time.Duration
		want   string
	}{
		{"3 days overdue", -3 * 24 * time.Hour, "3d overdue"},
		{"due today (slightly past)", -2 * time.Hour, "due today"},
		{"due today (slightly future)", 6 * time.Hour, "due today"},
		{"1 day left", 36 * time.Hour, "1d left"},
		{"5 days left", 5*24*time.Hour + 12*time.Hour, "5d left"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			due := time.Now().Add(tc.offset)
			iss := Issue{DueAt: &due}
			got := iss.DueLabel()
			if got != tc.want {
				t.Errorf("DueLabel() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDueLabelNil(t *testing.T) {
	iss := Issue{}
	if got := iss.DueLabel(); got != "" {
		t.Errorf("DueLabel() with nil DueAt = %q, want empty", got)
	}
}

func TestDeferLabel(t *testing.T) {
	future := time.Now().Add(5*24*time.Hour + 12*time.Hour)
	iss := Issue{DeferUntil: &future}
	got := iss.DeferLabel()
	if got != "deferred 5d" {
		t.Errorf("DeferLabel() = %q, want %q", got, "deferred 5d")
	}
}

func TestDeferLabelNil(t *testing.T) {
	iss := Issue{}
	if got := iss.DeferLabel(); got != "" {
		t.Errorf("DeferLabel() with nil DeferUntil = %q, want empty", got)
	}
}

func TestDeferLabelPast(t *testing.T) {
	past := time.Now().Add(-24 * time.Hour)
	iss := Issue{DeferUntil: &past}
	if got := iss.DeferLabel(); got != "" {
		t.Errorf("DeferLabel() with past DeferUntil = %q, want empty", got)
	}
}
