package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/matt-wright86/mardi-gras/internal/ui"
)

// FooterBinding is a key-description pair.
type FooterBinding struct {
	Key  string
	Desc string
}

// Footer renders the keybinding help bar.
type Footer struct {
	Width    int
	Bindings []FooterBinding
}

// ParadeBindings are the default keybindings for the parade view.
var ParadeBindings = []FooterBinding{
	{Key: "j/k", Desc: "navigate"},
	{Key: "tab", Desc: "switch pane"},
	{Key: "c", Desc: "toggle closed"},
	{Key: "q", Desc: "quit"},
}

// DetailBindings are keybindings when the detail pane is focused.
var DetailBindings = []FooterBinding{
	{Key: "j/k", Desc: "scroll"},
	{Key: "tab", Desc: "switch pane"},
	{Key: "esc", Desc: "back"},
	{Key: "q", Desc: "quit"},
}

// View renders the footer.
func (f Footer) View() string {
	var parts []string
	for _, b := range f.Bindings {
		key := ui.FooterKey.Render(b.Key)
		desc := ui.FooterDesc.Render(b.Desc)
		parts = append(parts, key+" "+desc)
	}
	content := strings.Join(parts, "  ")
	return ui.FooterStyle.Width(f.Width).Render(content)
}

// NewFooter creates a footer with the given width and pane focus.
func NewFooter(width int, detailFocused bool) Footer {
	bindings := ParadeBindings
	if detailFocused {
		bindings = DetailBindings
	}
	return Footer{Width: width, Bindings: bindings}
}

// Divider returns a full-width horizontal divider line.
func Divider(width int) string {
	return lipgloss.NewStyle().
		Foreground(ui.DimPurple).
		Width(width).
		Render(strings.Repeat(ui.DividerH, width))
}
