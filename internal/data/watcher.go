package data

import tea "github.com/charmbracelet/bubbletea"

// FileChangedMsg signals that the issues file was modified on disk.
type FileChangedMsg struct {
	Issues []Issue
}

// WatchFile starts watching a JSONL file for changes.
// Stub for v1 — will use fsnotify in a future version.
func WatchFile(_ string) tea.Cmd {
	return nil
}
