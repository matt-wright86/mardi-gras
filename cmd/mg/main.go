package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matt-wright86/mardi-gras/internal/app"
	"github.com/matt-wright86/mardi-gras/internal/data"
	"github.com/matt-wright86/mardi-gras/internal/tmux"
)

func main() {
	path := flag.String("path", "", "Path to .beads/issues.jsonl file")
	blockTypesFlag := flag.String("block-types", "", "Comma-separated dependency types that count as blockers (default: blocks)")
	statusMode := flag.Bool("status", false, "Output tmux status line and exit")
	flag.Parse()

	// Parse blocking types from flag, env var, or default
	blockingTypes := parseBlockingTypes(*blockTypesFlag)

	pathExplicit := *path != ""

	// Resolve JSONL path
	jsonlPath := *path
	if jsonlPath == "" {
		// Walk up from cwd looking for .beads/issues.jsonl
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
			os.Exit(1)
		}
		jsonlPath = findBeadsFile(cwd)
		if jsonlPath == "" {
			fmt.Fprintf(os.Stderr, "No .beads/issues.jsonl found.\n\n")
			fmt.Fprintf(os.Stderr, "Run mg from inside a project with Beads, or specify a path:\n")
			fmt.Fprintf(os.Stderr, "  mg --path /path/to/.beads/issues.jsonl\n")
			os.Exit(1)
		}
	}

	// Load issues
	issues, err := data.LoadIssues(jsonlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading issues from %s: %v\n", jsonlPath, err)
		os.Exit(1)
	}

	if *statusMode {
		groups := data.GroupByParade(issues, blockingTypes)
		fmt.Print(tmux.StatusLine(groups))
		return
	}

	// Run TUI
	model := app.New(issues, jsonlPath, pathExplicit, blockingTypes)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// parseBlockingTypes builds the blocking types set from flag, env var, or default.
func parseBlockingTypes(flagVal string) map[string]bool {
	raw := flagVal
	if raw == "" {
		raw = os.Getenv("MG_BLOCK_TYPES")
	}
	if raw == "" {
		return data.DefaultBlockingTypes
	}
	types := make(map[string]bool)
	for _, t := range strings.Split(raw, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			types[t] = true
		}
	}
	if len(types) == 0 {
		return data.DefaultBlockingTypes
	}
	return types
}

// findBeadsFile walks up from dir looking for .beads/issues.jsonl.
func findBeadsFile(dir string) string {
	for {
		candidate := filepath.Join(dir, ".beads", "issues.jsonl")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
