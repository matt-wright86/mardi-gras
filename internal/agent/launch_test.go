package agent

import (
	"strings"
	"testing"
	"time"

	"github.com/matt-wright86/mardi-gras/internal/data"
)

func TestAvailable(t *testing.T) {
	// Just verify it runs without panic; result depends on environment.
	_ = Available()
}

func TestBuildPromptFull(t *testing.T) {
	now := time.Now()
	issue := data.Issue{
		ID:                 "mg-001",
		Title:              "Deploy authentication service",
		Description:        "Set up OAuth2 flow for the API gateway.",
		Status:             data.StatusOpen,
		Priority:           data.PriorityCritical,
		IssueType:          data.TypeFeature,
		Owner:              "alice",
		Assignee:           "bob",
		CreatedAt:          now,
		UpdatedAt:          now,
		Notes:              "Needs review from security team.",
		AcceptanceCriteria: "All endpoints require valid JWT.",
		Dependencies: []data.Dependency{
			{IssueID: "mg-001", DependsOnID: "mg-002", Type: "blocks"},
		},
	}

	blocker := data.Issue{
		ID:        "mg-002",
		Title:     "Set up CI pipeline",
		Status:    data.StatusInProgress,
		Priority:  data.PriorityHigh,
		IssueType: data.TypeTask,
		CreatedAt: now,
		UpdatedAt: now,
	}

	issueMap := map[string]*data.Issue{
		"mg-001": &issue,
		"mg-002": &blocker,
	}

	deps := issue.EvaluateDependencies(issueMap, data.DefaultBlockingTypes)
	prompt := BuildPrompt(issue, deps, issueMap)

	for _, want := range []string{
		"mg-001",
		"Deploy authentication service",
		"Set up OAuth2 flow",
		"Owner: alice",
		"Assignee: bob",
		"### Notes",
		"Needs review from security team.",
		"### Acceptance Criteria",
		"All endpoints require valid JWT.",
		"Blocked by: mg-002",
		"Set up CI pipeline",
		"bd update mg-001 --status=in_progress",
		"bd close mg-001",
		"P0",
	} {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing %q\n\nGot:\n%s", want, prompt)
		}
	}
}

func TestBuildPromptMinimal(t *testing.T) {
	now := time.Now()
	issue := data.Issue{
		ID:        "mg-010",
		Title:     "Fix typo in README",
		Status:    data.StatusOpen,
		Priority:  data.PriorityBacklog,
		IssueType: data.TypeChore,
		CreatedAt: now,
		UpdatedAt: now,
	}

	issueMap := map[string]*data.Issue{"mg-010": &issue}
	deps := issue.EvaluateDependencies(issueMap, data.DefaultBlockingTypes)
	prompt := BuildPrompt(issue, deps, issueMap)

	if !strings.Contains(prompt, "mg-010") {
		t.Error("prompt missing issue ID")
	}
	if !strings.Contains(prompt, "Fix typo in README") {
		t.Error("prompt missing title")
	}

	// Optional sections should be absent.
	for _, absent := range []string{
		"### Notes",
		"### Acceptance Criteria",
		"### Dependencies",
		"Owner:",
		"Assignee:",
	} {
		if strings.Contains(prompt, absent) {
			t.Errorf("prompt should not contain %q for minimal issue\n\nGot:\n%s", absent, prompt)
		}
	}
}

func TestBuildPromptDependencies(t *testing.T) {
	now := time.Now()
	issue := data.Issue{
		ID:        "mg-020",
		Title:     "Main task",
		Status:    data.StatusOpen,
		Priority:  data.PriorityMedium,
		IssueType: data.TypeTask,
		CreatedAt: now,
		UpdatedAt: now,
		Dependencies: []data.Dependency{
			{IssueID: "mg-020", DependsOnID: "mg-021", Type: "blocks"},
			{IssueID: "mg-020", DependsOnID: "mg-022", Type: "blocks"},
			{IssueID: "mg-020", DependsOnID: "mg-ghost", Type: "blocks"},
			{IssueID: "mg-020", DependsOnID: "mg-023", Type: "related-to"},
		},
	}

	blocking := data.Issue{
		ID: "mg-021", Title: "Open blocker", Status: data.StatusOpen,
		Priority: data.PriorityMedium, IssueType: data.TypeTask,
		CreatedAt: now, UpdatedAt: now,
	}
	resolved := data.Issue{
		ID: "mg-022", Title: "Done blocker", Status: data.StatusClosed,
		Priority: data.PriorityMedium, IssueType: data.TypeTask,
		CreatedAt: now, UpdatedAt: now,
	}
	related := data.Issue{
		ID: "mg-023", Title: "Related item", Status: data.StatusOpen,
		Priority: data.PriorityLow, IssueType: data.TypeTask,
		CreatedAt: now, UpdatedAt: now,
	}

	issueMap := map[string]*data.Issue{
		"mg-020": &issue,
		"mg-021": &blocking,
		"mg-022": &resolved,
		"mg-023": &related,
	}

	deps := issue.EvaluateDependencies(issueMap, data.DefaultBlockingTypes)
	prompt := BuildPrompt(issue, deps, issueMap)

	for _, want := range []string{
		"Blocked by: mg-021 (Open blocker)",
		"Missing: mg-ghost (not found)",
		"Resolved: mg-022 (Done blocker) -- closed",
		"Related: mg-023 (Related item) -- related-to",
	} {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing %q\n\nGot:\n%s", want, prompt)
		}
	}
}

func TestCommandDir(t *testing.T) {
	cmd := Command("hello world", "/tmp/project")

	if cmd.Dir != "/tmp/project" {
		t.Errorf("expected Dir=%q, got %q", "/tmp/project", cmd.Dir)
	}

	if len(cmd.Args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(cmd.Args), cmd.Args)
	}
	if cmd.Args[0] != "claude" {
		t.Errorf("expected Args[0]=%q, got %q", "claude", cmd.Args[0])
	}
	if cmd.Args[1] != "hello world" {
		t.Errorf("expected Args[1]=%q, got %q", "hello world", cmd.Args[1])
	}
}
