package app

import (
	"context"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/matt-wright86/mardi-gras/internal/components"
	"github.com/matt-wright86/mardi-gras/internal/data"
)

func newReadyModel(t *testing.T) Model {
	t.Helper()

	issues := []data.Issue{
		testIssue("open-1", data.StatusOpen),
		testIssue("open-2", data.StatusOpen),
	}
	m := New(issues, data.Source{}, data.DefaultBlockingTypes)
	m.startedAt = time.Now().Add(-time.Second)
	model, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 20})
	return model.(Model)
}

func TestSyncSelectionClearsDetailWhenSelectionMissing(t *testing.T) {
	m := newReadyModel(t)
	if m.detail.Issue == nil {
		t.Fatal("expected detail issue to be initialized")
	}

	m.parade.SelectedIssue = nil
	m.syncSelection()

	if m.detail.Issue != nil {
		t.Fatalf("expected detail issue to be cleared, got %+v", m.detail.Issue)
	}
}

func TestExecutePaletteActionToggleGasTownClearsProblems(t *testing.T) {
	m := newReadyModel(t)
	m.gtEnv.Available = true
	m.showProblems = true

	model, cmd := m.executePaletteAction(components.ActionToggleGasTown)
	got := model.(Model)

	if !got.showGasTown {
		t.Fatal("expected showGasTown to be true")
	}
	if got.showProblems {
		t.Fatal("expected showProblems to be cleared when Gas Town opens")
	}
	if !got.gasTownTicking {
		t.Fatal("expected Gas Town ticking to start")
	}
	if cmd == nil {
		t.Fatal("expected Gas Town activation command")
	}
}

func TestExecutePaletteActionCycleLayoutActivatesGasTown(t *testing.T) {
	m := newReadyModel(t)
	m.gtEnv.Available = true
	m.showProblems = true

	model, cmd := m.executePaletteAction(components.ActionCycleLayout)
	got := model.(Model)

	if got.layoutPreset != LayoutGasTown {
		t.Fatalf("expected layout preset %v, got %v", LayoutGasTown, got.layoutPreset)
	}
	if !got.showGasTown {
		t.Fatal("expected Gas Town layout to open Gas Town")
	}
	if got.showProblems {
		t.Fatal("expected Problems panel to be hidden when cycling to Gas Town layout")
	}
	if cmd == nil {
		t.Fatal("expected layout cycle to return bootstrap commands")
	}
}

func TestCreateBranchCmdUsesProjectDir(t *testing.T) {
	cmd := createBranchCmd(context.Background(), "/tmp/project", "feat/mg-452-test")
	if cmd.Dir != "/tmp/project" {
		t.Fatalf("expected command dir %q, got %q", "/tmp/project", cmd.Dir)
	}
}
