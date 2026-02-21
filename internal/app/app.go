package app

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/matt-wright86/mardi-gras/internal/components"
	"github.com/matt-wright86/mardi-gras/internal/data"
	"github.com/matt-wright86/mardi-gras/internal/views"
)

// Pane tracks which panel is focused.
type Pane int

const (
	PaneParade Pane = iota
	PaneDetail
)

// Model is the root BubbleTea model.
type Model struct {
	issues      []data.Issue
	groups      map[data.ParadeStatus][]data.Issue
	parade      views.Parade
	detail      views.Detail
	header      components.Header
	activPane   Pane
	width       int
	height      int
	ready       bool
}

// New creates a new app model from loaded issues.
func New(issues []data.Issue) Model {
	groups := data.GroupByParade(issues)
	return Model{
		issues:    issues,
		groups:    groups,
		activPane: PaneParade,
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layout()
		m.ready = true
		return m, nil

	case data.FileChangedMsg:
		m.issues = msg.Issues
		m.groups = data.GroupByParade(msg.Issues)
		m.layout()
		return m, nil
	}

	// Forward to detail viewport when focused
	if m.activPane == PaneDetail {
		var cmd tea.Cmd
		m.detail.Viewport, cmd = m.detail.Viewport.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "tab":
		if m.activPane == PaneParade {
			m.activPane = PaneDetail
			m.detail.Focused = true
		} else {
			m.activPane = PaneParade
			m.detail.Focused = false
		}
		return m, nil

	case "esc":
		if m.activPane == PaneDetail {
			m.activPane = PaneParade
			m.detail.Focused = false
		}
		return m, nil

	case "c":
		m.parade.ToggleClosed()
		m.syncSelection()
		return m, nil
	}

	// Navigation keys depend on active pane
	if m.activPane == PaneParade {
		switch msg.String() {
		case "j", "down":
			m.parade.MoveDown()
			m.syncSelection()
		case "k", "up":
			m.parade.MoveUp()
			m.syncSelection()
		case "g":
			// Go to top
			m.parade.Cursor = 0
			m.parade.ScrollOffset = 0
			for i, item := range m.parade.Items {
				if !item.IsHeader {
					m.parade.Cursor = i
					break
				}
			}
			m.syncSelection()
		case "G":
			// Go to bottom
			for i := len(m.parade.Items) - 1; i >= 0; i-- {
				if !m.parade.Items[i].IsHeader {
					m.parade.Cursor = i
					break
				}
			}
			m.syncSelection()
		case "enter":
			m.activPane = PaneDetail
			m.detail.Focused = true
		}
		return m, nil
	}

	// Detail pane navigation
	if m.activPane == PaneDetail {
		var cmd tea.Cmd
		switch msg.String() {
		case "j", "down":
			m.detail.Viewport.LineDown(1)
		case "k", "up":
			m.detail.Viewport.LineUp(1)
		default:
			m.detail.Viewport, cmd = m.detail.Viewport.Update(msg)
		}
		return m, cmd
	}

	return m, nil
}

// syncSelection updates the detail panel with the currently selected issue.
func (m *Model) syncSelection() {
	if m.parade.SelectedIssue != nil {
		m.detail.SetIssue(m.parade.SelectedIssue)
	}
}

// layout recalculates dimensions for all sub-components.
func (m *Model) layout() {
	// Reserve lines for header (2) + divider (1) + footer (1) + divider (1)
	headerH := 2
	footerH := 2 // divider + footer
	bodyH := m.height - headerH - footerH
	if bodyH < 1 {
		bodyH = 1
	}

	// Split width: ~40% parade, ~60% detail
	paradeW := m.width * 2 / 5
	if paradeW < 30 {
		paradeW = 30
	}
	detailW := m.width - paradeW

	m.header = components.Header{
		Width:  m.width,
		Groups: m.groups,
	}

	m.parade.SetSize(paradeW, bodyH)
	m.detail.SetSize(detailW, bodyH)
	m.detail.AllIssues = m.issues
	m.detail.IssueMap = data.BuildIssueMap(m.issues)

	// Initialize parade if first layout
	if len(m.parade.Items) == 0 {
		m.parade = views.NewParade(m.issues, paradeW, bodyH)
		m.syncSelection()
	}

	// Initialize detail viewport
	m.detail.Viewport = viewport.New(detailW-2, bodyH)
	if m.parade.SelectedIssue != nil {
		m.detail.SetIssue(m.parade.SelectedIssue)
	}
}

// View implements tea.Model.
func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	header := m.header.View()

	// Body: parade | detail side by side
	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.parade.View(),
		m.detail.View(),
	)

	footer := components.NewFooter(m.width, m.activPane == PaneDetail)
	divider := components.Divider(m.width)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
		divider,
		footer.View(),
	)
}
