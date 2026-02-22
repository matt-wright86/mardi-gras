package gastown

import (
	"encoding/json"
	"os"
	"testing"
)

func TestDetectNoEnvVars(t *testing.T) {
	// Clear GT env vars for this test
	os.Unsetenv("GT_ROLE")
	os.Unsetenv("GT_RIG")
	os.Unsetenv("GT_SCOPE")
	os.Unsetenv("GT_POLECAT")
	os.Unsetenv("GT_CREW")
	env := Detect()
	if env.Active {
		t.Error("expected Active=false with no GT_ env vars")
	}
	if env.Role != "" {
		t.Errorf("expected empty Role, got %q", env.Role)
	}
}

func TestDetectWithEnvVars(t *testing.T) {
	t.Setenv("GT_ROLE", "polecat")
	t.Setenv("GT_RIG", "beads")
	t.Setenv("GT_POLECAT", "Toast")
	env := Detect()
	if !env.Active {
		t.Error("expected Active=true")
	}
	if env.Role != "polecat" {
		t.Errorf("expected Role=polecat, got %q", env.Role)
	}
	if env.Worker != "Toast" {
		t.Errorf("expected Worker=Toast, got %q", env.Worker)
	}
}

func TestDetectCrewWorker(t *testing.T) {
	t.Setenv("GT_ROLE", "crew")
	t.Setenv("GT_RIG", "beads")
	t.Setenv("GT_CREW", "Muffin")
	env := Detect()
	if env.Worker != "Muffin" {
		t.Errorf("expected Worker=Muffin, got %q", env.Worker)
	}
}

func TestTownStatusAgentForIssue(t *testing.T) {
	status := &TownStatus{
		Agents: []AgentRuntime{
			{Name: "Toast", Role: "polecat", HookBead: "bd-a1b2", Running: true, HasWork: true},
			{Name: "Muffin", Role: "polecat", HookBead: "bd-c3d4", Running: true, HasWork: true},
			{Name: "Whiskers", Role: "polecat", HookBead: "", Running: true, HasWork: false},
		},
	}
	agent := status.AgentForIssue("bd-a1b2")
	if agent == nil || agent.Name != "Toast" {
		t.Errorf("expected Toast for bd-a1b2, got %v", agent)
	}
	if status.AgentForIssue("bd-nope") != nil {
		t.Error("expected nil for unknown issue")
	}
}

func TestTownStatusAgentForIssueNil(t *testing.T) {
	var status *TownStatus
	if status.AgentForIssue("bd-a1b2") != nil {
		t.Error("expected nil from nil TownStatus")
	}
}

func TestTownStatusActiveAgentMap(t *testing.T) {
	status := &TownStatus{
		Agents: []AgentRuntime{
			{Name: "Toast", HookBead: "bd-a1b2", Running: true},
			{Name: "Muffin", HookBead: "", Running: true},        // no hook
			{Name: "Stale", HookBead: "bd-e5f6", Running: false}, // not running
		},
	}
	m := status.ActiveAgentMap()
	if len(m) != 1 {
		t.Errorf("expected 1 active agent, got %d", len(m))
	}
	if m["bd-a1b2"] != "Toast" {
		t.Errorf("expected Toast for bd-a1b2, got %q", m["bd-a1b2"])
	}
}

func TestTownStatusActiveAgentMapNil(t *testing.T) {
	var status *TownStatus
	m := status.ActiveAgentMap()
	if len(m) != 0 {
		t.Errorf("expected empty map from nil TownStatus, got %d", len(m))
	}
}

func TestTownStatusWorkingCount(t *testing.T) {
	status := &TownStatus{
		Agents: []AgentRuntime{
			{Name: "Toast", HasWork: true, Running: true},
			{Name: "Muffin", HasWork: false, Running: true},
			{Name: "Stale", HasWork: true, Running: false},
		},
	}
	if got := status.WorkingCount(); got != 1 {
		t.Errorf("expected 1 working, got %d", got)
	}

	var nilStatus *TownStatus
	if got := nilStatus.WorkingCount(); got != 0 {
		t.Errorf("expected 0 from nil, got %d", got)
	}
}

func TestTownStatusUnreadMail(t *testing.T) {
	status := &TownStatus{
		Agents: []AgentRuntime{
			{Name: "Toast", Mail: 3},
			{Name: "Muffin", Mail: 0},
			{Name: "Whiskers", Mail: 2},
		},
	}
	if got := status.UnreadMail(); got != 5 {
		t.Errorf("expected 5 unread, got %d", got)
	}

	var nilStatus *TownStatus
	if got := nilStatus.UnreadMail(); got != 0 {
		t.Errorf("expected 0 from nil, got %d", got)
	}
}

func TestTownStatusParsing(t *testing.T) {
	// Verify our structs can parse representative JSON
	raw := `{
		"agents": [
			{"name":"Toast","role":"polecat","rig":"beads","running":true,
			 "has_work":true,"work_title":"Fix login","hook_bead":"bd-a1b2",
			 "state":"working","unread_mail":0}
		],
		"rigs": [{"name":"beads","agents":3}],
		"convoys": [{"id":"hq-x1","title":"Auth fixes","status":"open","done":2,"total":5}]
	}`
	var status TownStatus
	if err := json.Unmarshal([]byte(raw), &status); err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(status.Agents) != 1 || status.Agents[0].Name != "Toast" {
		t.Errorf("unexpected agents: %+v", status.Agents)
	}
	if status.WorkingCount() != 1 {
		t.Errorf("expected 1 working, got %d", status.WorkingCount())
	}
	if len(status.Rigs) != 1 || status.Rigs[0].Name != "beads" {
		t.Errorf("unexpected rigs: %+v", status.Rigs)
	}
	if len(status.Convoys) != 1 || status.Convoys[0].ID != "hq-x1" {
		t.Errorf("unexpected convoys: %+v", status.Convoys)
	}
}
