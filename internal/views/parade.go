package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/matt-wright86/mardi-gras/internal/data"
	"github.com/matt-wright86/mardi-gras/internal/ui"
)

// paradeSection defines how each parade group renders.
type paradeSection struct {
	Title  string
	Symbol string
	Style  lipgloss.Style
	Status data.ParadeStatus
}

var sections = []paradeSection{
	{Title: "ROLLING", Symbol: ui.SymRolling, Style: ui.SectionRolling, Status: data.ParadeRolling},
	{Title: "LINED UP", Symbol: ui.SymLinedUp, Style: ui.SectionLinedUp, Status: data.ParadeLinedUp},
	{Title: "STALLED", Symbol: ui.SymStalled, Style: ui.SectionStalled, Status: data.ParadeStalled},
	{Title: "PAST THE STAND", Symbol: ui.SymPassed, Style: ui.SectionPassed, Status: data.ParadePastTheStand},
}

// ParadeItem is a renderable entry — either a section header or an issue.
type ParadeItem struct {
	IsHeader bool
	Section  paradeSection
	Issue    *data.Issue
}

// Parade is the grouped issue list view.
type Parade struct {
	Items         []ParadeItem
	Cursor        int
	ShowClosed    bool
	Width         int
	Height        int
	ScrollOffset  int
	AllIssues     []data.Issue
	Groups        map[data.ParadeStatus][]data.Issue
	issueMap      map[string]*data.Issue
	SelectedIssue *data.Issue
}

// NewParade creates a parade view from a set of issues.
func NewParade(issues []data.Issue, width, height int) Parade {
	groups := data.GroupByParade(issues)
	issueMap := data.BuildIssueMap(issues)
	p := Parade{
		ShowClosed: false,
		Width:      width,
		Height:     height,
		AllIssues:  issues,
		Groups:     groups,
		issueMap:   issueMap,
	}
	p.rebuildItems()
	if len(p.Items) > 0 {
		// Move cursor to first non-header item
		for i, item := range p.Items {
			if !item.IsHeader {
				p.Cursor = i
				p.SelectedIssue = item.Issue
				break
			}
		}
	}
	return p
}

// rebuildItems flattens groups into the renderable item list.
func (p *Parade) rebuildItems() {
	p.Items = nil
	for _, sec := range sections {
		issues := p.Groups[sec.Status]
		if len(issues) == 0 {
			continue
		}

		// Closed section: show collapsed count or expanded list
		if sec.Status == data.ParadePastTheStand {
			p.Items = append(p.Items, ParadeItem{IsHeader: true, Section: sec})
			if p.ShowClosed {
				for i := range issues {
					p.Items = append(p.Items, ParadeItem{Issue: &issues[i]})
				}
			}
			continue
		}

		p.Items = append(p.Items, ParadeItem{IsHeader: true, Section: sec})
		for i := range issues {
			p.Items = append(p.Items, ParadeItem{Issue: &issues[i]})
		}
	}
}

// MoveUp moves the cursor up, skipping headers.
func (p *Parade) MoveUp() {
	for i := p.Cursor - 1; i >= 0; i-- {
		if !p.Items[i].IsHeader {
			p.Cursor = i
			p.SelectedIssue = p.Items[i].Issue
			p.ensureVisible()
			return
		}
	}
}

// MoveDown moves the cursor down, skipping headers.
func (p *Parade) MoveDown() {
	for i := p.Cursor + 1; i < len(p.Items); i++ {
		if !p.Items[i].IsHeader {
			p.Cursor = i
			p.SelectedIssue = p.Items[i].Issue
			p.ensureVisible()
			return
		}
	}
}

// ToggleClosed shows or hides closed issues.
func (p *Parade) ToggleClosed() {
	p.ShowClosed = !p.ShowClosed
	selectedID := ""
	if p.SelectedIssue != nil {
		selectedID = p.SelectedIssue.ID
	}
	p.rebuildItems()
	// Restore cursor to the same issue if possible
	for i, item := range p.Items {
		if !item.IsHeader && item.Issue.ID == selectedID {
			p.Cursor = i
			p.SelectedIssue = item.Issue
			return
		}
	}
	// Fallback to first selectable item
	for i, item := range p.Items {
		if !item.IsHeader {
			p.Cursor = i
			p.SelectedIssue = item.Issue
			return
		}
	}
}

// ensureVisible adjusts scroll offset so cursor is visible.
func (p *Parade) ensureVisible() {
	if p.Cursor < p.ScrollOffset {
		p.ScrollOffset = p.Cursor
	}
	if p.Cursor >= p.ScrollOffset+p.Height {
		p.ScrollOffset = p.Cursor - p.Height + 1
	}
}

// SetSize updates the available dimensions.
func (p *Parade) SetSize(width, height int) {
	p.Width = width
	p.Height = height
}

// View renders the parade list.
func (p *Parade) View() string {
	var lines []string

	end := p.ScrollOffset + p.Height
	if end > len(p.Items) {
		end = len(p.Items)
	}

	visible := p.Items[p.ScrollOffset:end]

	for idx, item := range visible {
		globalIdx := p.ScrollOffset + idx
		if item.IsHeader {
			lines = append(lines, p.renderSectionHeader(item.Section))
		} else {
			lines = append(lines, p.renderIssue(item.Issue, globalIdx == p.Cursor))
		}
	}

	content := strings.Join(lines, "\n")

	// Pad to fill height
	rendered := strings.Count(content, "\n") + 1
	for rendered < p.Height {
		content += "\n"
		rendered++
	}

	return lipgloss.NewStyle().Width(p.Width).Render(content)
}

func (p *Parade) renderSectionHeader(sec paradeSection) string {
	count := len(p.Groups[sec.Status])
	label := fmt.Sprintf("%s %s", sec.Title, sec.Symbol)

	if sec.Status == data.ParadePastTheStand {
		toggle := ui.Collapsed
		if p.ShowClosed {
			toggle = ui.Expanded
		}
		label = fmt.Sprintf("%s %s (%d issues)", toggle, sec.Title, count)
		if !p.ShowClosed {
			label += "  [press c to expand]"
		}
	}

	return sec.Style.Width(p.Width).Render(label)
}

func (p *Parade) renderIssue(issue *data.Issue, selected bool) string {
	sym := statusSymbol(issue, p.issueMap)
	prio := data.PriorityLabel(issue.Priority)

	// Truncate title to fit
	maxTitle := p.Width - 18
	title := issue.Title
	if len(title) > maxTitle && maxTitle > 3 {
		title = title[:maxTitle-3] + "..."
	}

	prioStyle := ui.BadgePriority.Foreground(ui.PriorityColor(int(issue.Priority)))
	symStyle := lipgloss.NewStyle().Foreground(statusColor(issue, p.issueMap))

	line := fmt.Sprintf("%s %s %s %s",
		symStyle.Render(sym),
		issue.ID,
		title,
		prioStyle.Render(prio),
	)

	if selected {
		cursor := ui.ItemCursor.Render(ui.Cursor + " ")
		return cursor + line
	}

	return ui.ItemNormal.Render(line)
}

func statusSymbol(issue *data.Issue, issueMap map[string]*data.Issue) string {
	switch issue.Status {
	case data.StatusInProgress:
		return ui.SymRolling
	case data.StatusClosed:
		return ui.SymPassed
	default:
		if issue.IsBlocked(issueMap) {
			return ui.SymStalled
		}
		return ui.SymLinedUp
	}
}

func statusColor(issue *data.Issue, issueMap map[string]*data.Issue) lipgloss.Color {
	switch issue.Status {
	case data.StatusInProgress:
		return ui.StatusRolling
	case data.StatusClosed:
		return ui.StatusPassed
	default:
		if issue.IsBlocked(issueMap) {
			return ui.StatusStalled
		}
		return ui.StatusLinedUp
	}
}
