package data

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSampleIssues(t *testing.T) {
	// Find sample data relative to project root
	path := filepath.Join("..", "..", "testdata", "sample.jsonl")
	issues, err := LoadIssues(path)
	if err != nil {
		t.Fatalf("LoadIssues: %v", err)
	}

	if len(issues) != 10 {
		t.Fatalf("expected 10 issues, got %d", len(issues))
	}

	// Verify sorting: active issues come first
	for i, issue := range issues {
		if issue.Status == StatusClosed {
			// All remaining should be closed
			for j := i; j < len(issues); j++ {
				if issues[j].Status != StatusClosed {
					t.Errorf("issue %d (%s) is active but comes after closed issue %d", j, issues[j].ID, i)
				}
			}
			break
		}
	}
}

func TestGroupByParade(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "sample.jsonl")
	issues, err := LoadIssues(path)
	if err != nil {
		t.Fatalf("LoadIssues: %v", err)
	}

	groups := GroupByParade(issues)

	rolling := groups[ParadeRolling]
	linedUp := groups[ParadeLinedUp]
	stalled := groups[ParadeStalled]
	passed := groups[ParadePastTheStand]

	if len(rolling) != 2 {
		t.Errorf("expected 2 rolling, got %d", len(rolling))
	}
	// mg-006 depends on mg-001 (in_progress), so it should be stalled
	if len(stalled) != 1 {
		t.Errorf("expected 1 stalled, got %d", len(stalled))
	}
	if len(stalled) > 0 && stalled[0].ID != "mg-006" {
		t.Errorf("expected stalled issue to be mg-006, got %s", stalled[0].ID)
	}
	if len(linedUp) != 4 {
		t.Errorf("expected 4 lined up, got %d", len(linedUp))
	}
	if len(passed) != 3 {
		t.Errorf("expected 3 past the stand, got %d", len(passed))
	}
}

func TestIsBlocked(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "sample.jsonl")
	issues, err := LoadIssues(path)
	if err != nil {
		t.Fatalf("LoadIssues: %v", err)
	}

	issueMap := BuildIssueMap(issues)

	// mg-006 is blocked by mg-001 (in_progress)
	mg006 := issueMap["mg-006"]
	if mg006 == nil {
		t.Fatal("mg-006 not found")
	}
	if !mg006.IsBlocked(issueMap) {
		t.Error("mg-006 should be blocked")
	}

	blockers := mg006.BlockedByIDs(issueMap)
	if len(blockers) != 1 || blockers[0] != "mg-001" {
		t.Errorf("expected blockers [mg-001], got %v", blockers)
	}

	// mg-001 blocks mg-006
	mg001 := issueMap["mg-001"]
	blocks := mg001.BlocksIDs(issues)
	if len(blocks) != 1 || blocks[0] != "mg-006" {
		t.Errorf("expected mg-001 blocks [mg-006], got %v", blocks)
	}
}

func TestLoadRealBeads(t *testing.T) {
	// Test against real Beads data if available
	path := "/Users/matthewwright/Work/voice-vault/.beads/issues.jsonl"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("real beads data not available")
	}

	issues, err := LoadIssues(path)
	if err != nil {
		t.Fatalf("LoadIssues (real data): %v", err)
	}

	if len(issues) == 0 {
		t.Error("expected at least 1 issue from real data")
	}

	groups := GroupByParade(issues)
	t.Logf("Real data: %d total, %d rolling, %d lined up, %d stalled, %d passed",
		len(issues),
		len(groups[ParadeRolling]),
		len(groups[ParadeLinedUp]),
		len(groups[ParadeStalled]),
		len(groups[ParadePastTheStand]),
	)
}
