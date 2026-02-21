package app

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matt-wright86/mardi-gras/internal/data"
)

func testIssue(id string, status data.Status) data.Issue {
	now := time.Now()
	return data.Issue{
		ID:        id,
		Title:     id,
		Status:    status,
		Priority:  data.PriorityMedium,
		IssueType: data.TypeTask,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestFileChangedMsgPreservesSelectionAndClosedState(t *testing.T) {
	issues := []data.Issue{
		testIssue("open-1", data.StatusOpen),
		testIssue("open-2", data.StatusOpen),
		testIssue("closed-1", data.StatusClosed),
	}

	m := New(issues, "")
	model, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 20})
	got := model.(Model)

	// Move selection to second open issue.
	model, _ = got.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	got = model.(Model)
	if got.parade.SelectedIssue == nil || got.parade.SelectedIssue.ID != "open-2" {
		t.Fatalf("expected selected issue open-2 before refresh, got %+v", got.parade.SelectedIssue)
	}

	// Expand closed section.
	model, _ = got.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	got = model.(Model)
	if !got.parade.ShowClosed {
		t.Fatal("expected closed section expanded before refresh")
	}

	// Simulate file refresh with same issues.
	model, _ = got.Update(data.FileChangedMsg{Issues: issues})
	got = model.(Model)

	if !got.parade.ShowClosed {
		t.Fatal("expected closed section to remain expanded after refresh")
	}
	if got.parade.SelectedIssue == nil || got.parade.SelectedIssue.ID != "open-2" {
		t.Fatalf("expected selected issue open-2 after refresh, got %+v", got.parade.SelectedIssue)
	}
}
