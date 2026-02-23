package data

import "testing"

func TestParseExternalRef(t *testing.T) {
	tests := []struct {
		input   string
		wantNil bool
		rig     string
		id      string
	}{
		{"external:gastown:gt-c3f2", false, "gastown", "gt-c3f2"},
		{"external:wyvern:wy-e5f6", false, "wyvern", "wy-e5f6"},
		{"external:beads:bd-001", false, "beads", "bd-001"},
		{"bd-001", true, "", ""},        // not external
		{"external:", true, "", ""},       // malformed
		{"external:foo", true, "", ""},    // missing id
		{"blocks", true, "", ""},          // plain dependency type
	}
	for _, tt := range tests {
		ref := ParseExternalRef(tt.input)
		if tt.wantNil {
			if ref != nil {
				t.Errorf("ParseExternalRef(%q) should be nil, got %+v", tt.input, ref)
			}
			continue
		}
		if ref == nil {
			t.Fatalf("ParseExternalRef(%q) should not be nil", tt.input)
		}
		if ref.Rig != tt.rig {
			t.Errorf("ParseExternalRef(%q).Rig = %q, want %q", tt.input, ref.Rig, tt.rig)
		}
		if ref.IssueID != tt.id {
			t.Errorf("ParseExternalRef(%q).IssueID = %q, want %q", tt.input, ref.IssueID, tt.id)
		}
	}
}

func TestCrossRigDeps(t *testing.T) {
	issue := &Issue{
		ID: "bd-001",
		Dependencies: []Dependency{
			{DependsOnID: "bd-002", Type: "blocks"},
			{DependsOnID: "external:gastown:gt-c3f2", Type: "blocks"},
			{DependsOnID: "external:wyvern:wy-e5f6", Type: "related"},
			{DependsOnID: "bd-003", Type: "blocks"},
		},
	}

	refs := CrossRigDeps(issue)
	if len(refs) != 2 {
		t.Fatalf("expected 2 cross-rig deps, got %d", len(refs))
	}
	if refs[0].Rig != "gastown" {
		t.Fatalf("refs[0].Rig = %q, want 'gastown'", refs[0].Rig)
	}
	if refs[1].Rig != "wyvern" {
		t.Fatalf("refs[1].Rig = %q, want 'wyvern'", refs[1].Rig)
	}
}

func TestCrossRigDepsNone(t *testing.T) {
	issue := &Issue{
		ID: "bd-001",
		Dependencies: []Dependency{
			{DependsOnID: "bd-002", Type: "blocks"},
		},
	}

	refs := CrossRigDeps(issue)
	if len(refs) != 0 {
		t.Fatalf("expected 0 cross-rig deps, got %d", len(refs))
	}
}

func TestCrossRigSummary(t *testing.T) {
	issues := []Issue{
		{
			ID: "bd-001",
			Dependencies: []Dependency{
				{DependsOnID: "external:gastown:gt-001", Type: "blocks"},
				{DependsOnID: "external:gastown:gt-002", Type: "blocks"},
			},
		},
		{
			ID: "bd-002",
			Dependencies: []Dependency{
				{DependsOnID: "external:wyvern:wy-001", Type: "blocks"},
			},
		},
		{
			ID: "bd-003",
			Dependencies: []Dependency{
				{DependsOnID: "bd-001", Type: "blocks"},
			},
		},
	}

	summary := CrossRigSummary(issues)
	if summary["gastown"] != 2 {
		t.Fatalf("gastown count = %d, want 2", summary["gastown"])
	}
	if summary["wyvern"] != 1 {
		t.Fatalf("wyvern count = %d, want 1", summary["wyvern"])
	}
	if len(summary) != 2 {
		t.Fatalf("expected 2 rigs in summary, got %d", len(summary))
	}
}

func TestCrossRigSummaryEmpty(t *testing.T) {
	summary := CrossRigSummary(nil)
	if len(summary) != 0 {
		t.Fatalf("expected empty summary, got %d", len(summary))
	}
}
