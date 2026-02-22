package agent

import (
	"os"
	"testing"
)

func TestInTmux(t *testing.T) {
	// Save and restore original value.
	orig := os.Getenv("TMUX")
	defer os.Setenv("TMUX", orig)

	os.Setenv("TMUX", "/tmp/tmux-1000/default,12345,0")
	if !InTmux() {
		t.Error("expected InTmux()=true when TMUX is set")
	}

	os.Unsetenv("TMUX")
	if InTmux() {
		t.Error("expected InTmux()=false when TMUX is unset")
	}
}

func TestWindowName(t *testing.T) {
	tests := []struct {
		issueID string
		want    string
	}{
		{"bd-a1b2", "mg-bd-a1b2"},
		{"mg-001", "mg-mg-001"},
		{"xyz", "mg-xyz"},
	}
	for _, tt := range tests {
		got := WindowName(tt.issueID)
		if got != tt.want {
			t.Errorf("WindowName(%q) = %q, want %q", tt.issueID, got, tt.want)
		}
	}
}

func TestParseAgentPanes(t *testing.T) {
	output := "mg-bd-a1b2\t%5\n\t%0\nmg-bd-c3d4\t%8\n\t%1\n"
	agents := parseAgentPanes(output)

	if len(agents) != 2 {
		t.Fatalf("expected 2 agent panes, got %d: %v", len(agents), agents)
	}

	if agents["bd-a1b2"] != "%5" {
		t.Errorf("missing or wrong entry for bd-a1b2: %v", agents)
	}
	if agents["bd-c3d4"] != "%8" {
		t.Errorf("missing or wrong entry for bd-c3d4: %v", agents)
	}
}

func TestParseAgentPanesEmpty(t *testing.T) {
	agents := parseAgentPanes("")
	if len(agents) != 0 {
		t.Errorf("expected 0 agent panes from empty input, got %d", len(agents))
	}
}

func TestParseAgentPanesNoAgents(t *testing.T) {
	output := "\t%0\n\t%1\n\t%2\n"
	agents := parseAgentPanes(output)
	if len(agents) != 0 {
		t.Errorf("expected 0 agent panes, got %d: %v", len(agents), agents)
	}
}
