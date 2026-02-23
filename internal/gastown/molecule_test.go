package gastown

import (
	"encoding/json"
	"testing"
)

func TestDAGInfoParsing(t *testing.T) {
	raw := `{
		"root_id": "gs-123",
		"root_title": "Build feature X",
		"total_nodes": 5,
		"tiers": 3,
		"critical_path": ["gs-123.1", "gs-123.3", "gs-123.5"],
		"nodes": {
			"gs-123.1": {"id": "gs-123.1", "title": "Design", "status": "done", "tier": 0},
			"gs-123.2": {"id": "gs-123.2", "title": "Implement A", "status": "done", "tier": 1, "dependencies": ["gs-123.1"]},
			"gs-123.3": {"id": "gs-123.3", "title": "Implement B", "status": "in_progress", "tier": 1, "dependencies": ["gs-123.1"], "parallel": true},
			"gs-123.4": {"id": "gs-123.4", "title": "Test", "status": "blocked", "tier": 2, "dependencies": ["gs-123.2", "gs-123.3"]},
			"gs-123.5": {"id": "gs-123.5", "title": "Submit", "status": "blocked", "tier": 2, "dependencies": ["gs-123.4"]}
		},
		"tier_groups": [["gs-123.1"], ["gs-123.2", "gs-123.3"], ["gs-123.4", "gs-123.5"]]
	}`

	var dag DAGInfo
	if err := json.Unmarshal([]byte(raw), &dag); err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if dag.RootID != "gs-123" {
		t.Errorf("RootID = %q, want %q", dag.RootID, "gs-123")
	}
	if dag.TotalNodes != 5 {
		t.Errorf("TotalNodes = %d, want 5", dag.TotalNodes)
	}
	if dag.Tiers != 3 {
		t.Errorf("Tiers = %d, want 3", dag.Tiers)
	}
	if len(dag.CriticalPath) != 3 {
		t.Errorf("CriticalPath len = %d, want 3", len(dag.CriticalPath))
	}
	if len(dag.Nodes) != 5 {
		t.Errorf("Nodes len = %d, want 5", len(dag.Nodes))
	}
	if len(dag.TierGroups) != 3 {
		t.Errorf("TierGroups len = %d, want 3", len(dag.TierGroups))
	}

	// Check node properties
	node := dag.Nodes["gs-123.3"]
	if node == nil {
		t.Fatal("node gs-123.3 not found")
	}
	if node.Status != "in_progress" {
		t.Errorf("node status = %q, want %q", node.Status, "in_progress")
	}
	if !node.Parallel {
		t.Error("node should be parallel")
	}
	if len(node.Dependencies) != 1 || node.Dependencies[0] != "gs-123.1" {
		t.Errorf("node dependencies = %v, want [gs-123.1]", node.Dependencies)
	}
}

func TestDAGInfoActiveStepID(t *testing.T) {
	tests := []struct {
		name   string
		dag    *DAGInfo
		wantID string
	}{
		{
			name:   "nil DAG",
			dag:    nil,
			wantID: "",
		},
		{
			name: "in_progress step preferred",
			dag: &DAGInfo{
				Nodes: map[string]*DAGNode{
					"s1": {ID: "s1", Status: "done", Tier: 0},
					"s2": {ID: "s2", Status: "in_progress", Tier: 1},
					"s3": {ID: "s3", Status: "ready", Tier: 1},
				},
				TierGroups: [][]string{{"s1"}, {"s2", "s3"}},
			},
			wantID: "s2",
		},
		{
			name: "ready step fallback",
			dag: &DAGInfo{
				Nodes: map[string]*DAGNode{
					"s1": {ID: "s1", Status: "done", Tier: 0},
					"s2": {ID: "s2", Status: "blocked", Tier: 1},
					"s3": {ID: "s3", Status: "ready", Tier: 1},
				},
				TierGroups: [][]string{{"s1"}, {"s2", "s3"}},
			},
			wantID: "s3",
		},
		{
			name: "all done",
			dag: &DAGInfo{
				Nodes: map[string]*DAGNode{
					"s1": {ID: "s1", Status: "done", Tier: 0},
					"s2": {ID: "s2", Status: "done", Tier: 1},
				},
				TierGroups: [][]string{{"s1"}, {"s2"}},
			},
			wantID: "",
		},
		{
			name: "empty nodes",
			dag: &DAGInfo{
				Nodes:      map[string]*DAGNode{},
				TierGroups: [][]string{},
			},
			wantID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.dag.ActiveStepID()
			if got != tt.wantID {
				t.Errorf("ActiveStepID() = %q, want %q", got, tt.wantID)
			}
		})
	}
}

func TestMoleculeProgressParsing(t *testing.T) {
	raw := `{
		"root_id": "gs-123",
		"root_title": "Build feature X",
		"molecule_id": "mol-456",
		"total_steps": 7,
		"done_steps": 3,
		"in_progress_steps": 1,
		"ready_steps": ["gs-123.5"],
		"blocked_steps": ["gs-123.6", "gs-123.7"],
		"percent_complete": 42,
		"complete": false
	}`

	var prog MoleculeProgress
	if err := json.Unmarshal([]byte(raw), &prog); err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if prog.TotalSteps != 7 {
		t.Errorf("TotalSteps = %d, want 7", prog.TotalSteps)
	}
	if prog.DoneSteps != 3 {
		t.Errorf("DoneSteps = %d, want 3", prog.DoneSteps)
	}
	if prog.InProgress != 1 {
		t.Errorf("InProgress = %d, want 1", prog.InProgress)
	}
	if prog.Percent != 42 {
		t.Errorf("Percent = %d, want 42", prog.Percent)
	}
	if prog.Complete {
		t.Error("should not be complete")
	}
	if len(prog.ReadySteps) != 1 {
		t.Errorf("ReadySteps len = %d, want 1", len(prog.ReadySteps))
	}
	if len(prog.BlockedSteps) != 2 {
		t.Errorf("BlockedSteps len = %d, want 2", len(prog.BlockedSteps))
	}
}

func TestStepDoneResultParsing(t *testing.T) {
	raw := `{
		"step_id": "gs-123.3",
		"molecule_id": "mol-456",
		"step_closed": true,
		"next_step_id": "gs-123.4",
		"next_step_title": "Test",
		"parallel_steps": [],
		"complete": false,
		"action": "continue"
	}`

	var result StepDoneResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if result.StepID != "gs-123.3" {
		t.Errorf("StepID = %q, want %q", result.StepID, "gs-123.3")
	}
	if !result.StepClosed {
		t.Error("should be closed")
	}
	if result.NextStepID != "gs-123.4" {
		t.Errorf("NextStepID = %q, want %q", result.NextStepID, "gs-123.4")
	}
	if result.Action != "continue" {
		t.Errorf("Action = %q, want %q", result.Action, "continue")
	}
	if result.Complete {
		t.Error("should not be complete")
	}
}

func TestStepDoneResultComplete(t *testing.T) {
	raw := `{
		"step_id": "gs-123.7",
		"molecule_id": "mol-456",
		"step_closed": true,
		"complete": true,
		"action": "done"
	}`

	var result StepDoneResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if !result.Complete {
		t.Error("should be complete")
	}
	if result.Action != "done" {
		t.Errorf("Action = %q, want %q", result.Action, "done")
	}
}
