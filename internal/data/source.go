package data

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	tea "charm.land/bubbletea/v2"
)

// SourceMode indicates how issues are loaded.
type SourceMode int

const (
	SourceJSONL SourceMode = iota // Legacy: read from .beads/issues.jsonl (or --path)
	SourceCLI                     // Preferred: shell out to bd list --json
)

// Source describes how mg loads its issue data.
type Source struct {
	Mode       SourceMode
	Path       string // JSONL file path (SourceJSONL) or empty (SourceCLI)
	ProjectDir string // Project root directory
	Explicit   bool   // True if --path was used
}

// Label returns a display string for the footer.
func (s Source) Label() string {
	if s.Mode == SourceCLI {
		return "bd list"
	}
	if s.Path != "" {
		return filepath.Base(s.Path)
	}
	return "issues.jsonl"
}

// FetchIssuesCLI runs `bd list --json --limit 0 --all` and parses the result.
func FetchIssuesCLI() ([]Issue, error) {
	out, err := runWithTimeout(timeoutMedium, "bd", "list", "--json", "--limit", "0", "--all")
	if err != nil {
		return nil, fmt.Errorf("bd list --json: %w", err)
	}
	var issues []Issue
	if err := json.Unmarshal(out, &issues); err != nil {
		return nil, fmt.Errorf("bd list parse: %w", err)
	}
	SortIssues(issues)
	return issues, nil
}

// FetchCurrentIssueID runs `bd show --current --json` and returns the active issue ID.
// Returns ("", nil) if no current issue exists (bd exits non-zero).
func FetchCurrentIssueID() (string, error) {
	out, err := runWithTimeout(timeoutShort, "bd", "show", "--current", "--json")
	if err != nil {
		return "", nil // bd exits non-zero when no current issue — not an error
	}
	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(out, &result); err != nil {
		return "", fmt.Errorf("bd show --current parse: %w", err)
	}
	return result.ID, nil
}

// DoctorDiagnostic represents a single finding from `bd doctor --agent --json`.
type DoctorDiagnostic struct {
	Name        string   `json:"name"`
	Status      string   `json:"status"`       // "error", "warning", "ok"
	Severity    string   `json:"severity"`     // "blocking", "degraded", etc.
	Category    string   `json:"category"`     // "Core System", "Git Integration", etc.
	Explanation string   `json:"explanation"`  // Human-readable detail
	Observed    string   `json:"observed"`     // What was found
	Expected    string   `json:"expected"`     // What was expected
	Commands    []string `json:"commands"`     // Suggested fix commands
	SourceFiles []string `json:"source_files"` // Upstream source references
}

// DoctorResult holds the full output of `bd doctor --agent --json`.
type DoctorResult struct {
	OK          bool               `json:"overall_ok"`
	Summary     string             `json:"summary"`
	Diagnostics []DoctorDiagnostic `json:"diagnostics"`
}

// FetchDoctorDiagnostics runs `bd doctor --agent --json` and returns findings.
// Only returns error/warning diagnostics (not passed checks).
func FetchDoctorDiagnostics() (*DoctorResult, error) {
	out, err := runWithTimeout(timeoutMedium, "bd", "doctor", "--agent", "--json")
	if err != nil {
		// bd doctor exits non-zero when problems found — still has valid JSON on stdout
		if out == nil {
			return nil, fmt.Errorf("bd doctor: %w", err)
		}
	}
	var result DoctorResult
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("bd doctor parse: %w", err)
	}
	return &result, nil
}

// FetchIssuesNow returns a tea.Cmd that fetches issues via bd CLI immediately
// (no timer delay). Emits FileChangedMsg on success, FileWatchErrorMsg on failure.
func FetchIssuesNow() tea.Cmd {
	return func() tea.Msg {
		issues, err := FetchIssuesCLI()
		if err != nil {
			return FileWatchErrorMsg{Err: err}
		}
		return FileChangedMsg{Issues: issues, LastMod: time.Now()}
	}
}
