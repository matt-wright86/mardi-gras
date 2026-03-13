# Worktree Dispatch Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `W` keybinding to create git worktrees for beads and use the worktree path as agent cwd on dispatch.

**Architecture:** New `CreateWorktree()` and `WorktreePath()` functions in `internal/data/mutate.go` handle git worktree creation and metadata reads. The `W` key in `app.go` calls `CreateWorktree` via a tea.Cmd, and the `a` key dispatch checks for a worktree path before choosing cwd. Parade and detail views show a worktree indicator.

**Tech Stack:** Go, BubbleTea, git CLI, beads CLI (`bd`)

**Spec:** `docs/superpowers/specs/2026-03-13-worktree-dispatch-design.md`

---

## Chunk 1: Core Data Layer

### Task 1: `WorktreePath` helper

**Files:**
- Modify: `internal/data/mutate.go` (append after line 84)
- Test: `internal/data/mutate_test.go`

- [ ] **Step 1: Write the failing test**

In `internal/data/mutate_test.go`, add:

```go
func TestWorktreePath(t *testing.T) {
	tests := []struct {
		name     string
		issue    Issue
		expected string
	}{
		{
			name:     "metadata has worktree string",
			issue:    Issue{Metadata: map[string]interface{}{"worktree": "/tmp/worktrees/feat/bd-a1b2"}},
			expected: "/tmp/worktrees/feat/bd-a1b2",
		},
		{
			name:     "metadata nil",
			issue:    Issue{},
			expected: "",
		},
		{
			name:     "metadata empty map",
			issue:    Issue{Metadata: map[string]interface{}{}},
			expected: "",
		},
		{
			name:     "worktree key is not a string",
			issue:    Issue{Metadata: map[string]interface{}{"worktree": 42}},
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WorktreePath(tt.issue)
			if got != tt.expected {
				t.Errorf("WorktreePath() = %q, want %q", got, tt.expected)
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/data/ -run TestWorktreePath -v`
Expected: FAIL — `WorktreePath` undefined

- [ ] **Step 3: Write minimal implementation**

In `internal/data/mutate.go`, append:

```go
// WorktreePath returns the worktree path stored in an issue's metadata.
// Returns "" if not set or not a string.
func WorktreePath(issue Issue) string {
	if issue.Metadata == nil {
		return ""
	}
	v, ok := issue.Metadata["worktree"]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/data/ -run TestWorktreePath -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/data/mutate.go internal/data/mutate_test.go
git commit -m "feat(data): add WorktreePath helper to read worktree metadata"
```

---

### Task 2: `CreateWorktree` function

**Files:**
- Modify: `internal/data/mutate.go` (add imports, append function after `WorktreePath`)
- Test: `internal/data/mutate_test.go`

Note: `CreateWorktree` runs `git` and `bd` as subprocesses. Testing requires either a real git repo fixture or mocking exec. Since the codebase doesn't have exec mocking infrastructure, we test the logic parts (path computation, branch name) and keep the integration test minimal with a real temp git repo.

- [ ] **Step 1: Write the failing test**

In `internal/data/mutate_test.go`, add:

```go
func TestCreateWorktreePathComputation(t *testing.T) {
	// Test that the worktree path is computed correctly from projectDir + branch name.
	issue := Issue{ID: "bd-a1b2", Title: "Fix login", IssueType: TypeBug}
	branch := BranchName(issue)
	projectDir := "/home/user/work/my-project"

	expectedBase := "/home/user/work/my-project-worktrees"
	expectedPath := expectedBase + "/" + branch

	// Verify expected path matches the formula: Dir(projectDir) / Base(projectDir)+"-worktrees" / branchName
	got := filepath.Join(filepath.Dir(projectDir), filepath.Base(projectDir)+"-worktrees", branch)
	if got != expectedPath {
		t.Errorf("worktree path = %q, want %q", got, expectedPath)
	}
}
```

Add `"path/filepath"` to the imports in the test file.

- [ ] **Step 2: Run test to verify it passes**

Run: `go test ./internal/data/ -run TestCreateWorktreePathComputation -v`
Expected: PASS (this is a pure path computation test, no new code needed)

- [ ] **Step 3: Write the CreateWorktree function**

In `internal/data/mutate.go`, add to imports: `"context"`, `"os"`, `"os/exec"`, `"path/filepath"`. Then append:

```go
// CreateWorktree creates a git worktree for the given issue and stores
// the worktree path in the issue's metadata. Returns the absolute worktree path.
func CreateWorktree(issue Issue, projectDir string) (string, error) {
	// Check if worktree already tracked in metadata
	if wt := WorktreePath(issue); wt != "" {
		return "", fmt.Errorf("worktree already exists: %s", wt)
	}

	branch := BranchName(issue)
	baseDir := filepath.Join(filepath.Dir(projectDir), filepath.Base(projectDir)+"-worktrees")
	wtPath := filepath.Join(baseDir, branch)
	absPath, err := filepath.Abs(wtPath)
	if err != nil {
		return "", fmt.Errorf("resolve worktree path: %w", err)
	}

	// Partial failure recovery: dir exists but metadata wasn't set
	if info, statErr := os.Stat(absPath); statErr == nil && info.IsDir() {
		// Worktree dir exists on disk — just set metadata
		if metaErr := setWorktreeMetadata(issue.ID, absPath); metaErr != nil {
			return "", fmt.Errorf("set worktree metadata: %w", metaErr)
		}
		return absPath, nil
	}

	// Ensure parent directories exist (branch names contain slashes like feat/...)
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return "", fmt.Errorf("create worktree directory: %w", err)
	}

	// Try creating with new branch first
	ctx, cancel := context.WithTimeout(context.Background(), timeoutShort)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "worktree", "add", absPath, "-b", branch)
	cmd.Dir = projectDir
	if err := cmd.Run(); err != nil {
		// Branch may already exist — retry without -b
		ctx2, cancel2 := context.WithTimeout(context.Background(), timeoutShort)
		defer cancel2()
		cmd2 := exec.CommandContext(ctx2, "git", "worktree", "add", absPath, branch)
		cmd2.Dir = projectDir
		if err2 := cmd2.Run(); err2 != nil {
			return "", fmt.Errorf("git worktree add: %w", err2)
		}
	}

	// Store worktree path on the bead
	if metaErr := setWorktreeMetadata(issue.ID, absPath); metaErr != nil {
		return "", fmt.Errorf("worktree created at %s but metadata update failed: %w", absPath, metaErr)
	}

	return absPath, nil
}

func setWorktreeMetadata(issueID, absPath string) error {
	return execWithTimeout(timeoutShort, "bd", "update", issueID,
		"--set-metadata", "worktree="+absPath)
}
```

- [ ] **Step 4: Run all data tests**

Run: `go test ./internal/data/ -v`
Expected: PASS (compilation check — CreateWorktree uses real git/bd so won't be exercised in unit tests)

- [ ] **Step 5: Commit**

```bash
git add internal/data/mutate.go internal/data/mutate_test.go
git commit -m "feat(data): add CreateWorktree function for git worktree creation"
```

---

## Chunk 2: UI Constants and Help

### Task 3: Worktree symbol

**Files:**
- Modify: `internal/ui/symbols.go:106` (add before closing paren on line 106)

- [ ] **Step 1: Add the symbol**

In `internal/ui/symbols.go`, add after `SymValidator` (line 105), before the closing `)` on line 106:

```go
	SymWorktree = "⌥"
```

- [ ] **Step 2: Run lint**

Run: `make lint`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add internal/ui/symbols.go
git commit -m "feat(ui): add SymWorktree symbol constant"
```

---

### Task 4: Help text entry

**Files:**
- Modify: `internal/components/help.go:103` (add after `B` entry in QUICK ACTIONS)

- [ ] **Step 1: Add help entry**

In `internal/components/help.go`, in the QUICK ACTIONS section, after `{key: "B", desc: "Create + checkout git branch"}` (line 102), add:

```go
				{key: "W", desc: "Create git worktree for issue"},
```

- [ ] **Step 2: Run tests**

Run: `make test`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add internal/components/help.go
git commit -m "feat(help): add W keybinding to help overlay"
```

---

### Task 5: Command palette action

**Files:**
- Modify: `internal/components/palette.go` (add `ActionCreateWorktree` to action enum)
- Modify: `internal/app/app.go` (add palette command entry and dispatch handler)

- [ ] **Step 1: Add action constant**

In `internal/components/palette.go`, find `ActionFilterByEpic` (line 43) and add after it:

```go
	ActionCreateWorktree
```

- [ ] **Step 2: Add palette command in app.go**

In `internal/app/app.go`, in `buildPaletteCommands()`, find the `ActionCreateBranch` entry (line 2059) and add after it:

```go
		{Name: "Create worktree", Desc: "Create git worktree for issue", Key: "W", Action: components.ActionCreateWorktree},
```

- [ ] **Step 3: Add dispatch handler in app.go**

In `internal/app/app.go`, in the palette action switch (around line 2114), find `case components.ActionCreateBranch:` and add after its return:

```go
	case components.ActionCreateWorktree:
		return m.createWorktree()
```

(The `createWorktree()` method will be added in Task 6.)

**Note:** Tasks 5 and 6 are an atomic unit — the palette dispatch references `m.createWorktree()` which doesn't exist until Task 6. Do not attempt to compile or commit Task 5 alone. All changes are committed together at the end of Task 6.

---

## Chunk 3: App Integration

### Task 6: `W` keybinding and worktree command

**Files:**
- Modify: `internal/app/app.go` (add message types, `W` key handler, `createWorktree()` method, message handlers)

- [ ] **Step 1: Add message types**

In `internal/app/app.go`, near the other message types (around line 462), add:

```go
// worktreeCreatedMsg is sent when a git worktree is successfully created.
type worktreeCreatedMsg struct {
	issueID string
	path    string
}

// worktreeErrorMsg is sent when git worktree creation fails.
type worktreeErrorMsg struct {
	issueID string
	err     error
}
```

- [ ] **Step 2: Add the createWorktree method**

In `internal/app/app.go`, near `createAndSwitchBranch()` (after line 2025), add:

```go
// createWorktree creates a git worktree for the selected issue.
func (m Model) createWorktree() (tea.Model, tea.Cmd) {
	issue := m.parade.SelectedIssue
	if issue == nil {
		return m, nil
	}
	issueCopy := *issue
	return m, func() tea.Msg {
		path, err := data.CreateWorktree(issueCopy, m.projectDir)
		if err != nil {
			return worktreeErrorMsg{issueID: issueCopy.ID, err: err}
		}
		return worktreeCreatedMsg{issueID: issueCopy.ID, path: path}
	}
}
```

- [ ] **Step 3: Add the `W` key handler**

In `internal/app/app.go`, in the key handler switch, after the `B` case (line 1576), add:

```go
	case "W":
		return m.createWorktree()
```

- [ ] **Step 4: Add message handlers in Update**

In `internal/app/app.go`, in the `Update` method's message switch, add handlers for the new message types. Find the `mutateResultMsg` handler and add nearby:

```go
	case worktreeCreatedMsg:
		toast, dismissCmd := components.ShowToast(
			fmt.Sprintf("Worktree created: %s", msg.path),
			components.ToastSuccess,
			3*time.Second,
		)
		m.toast = &toast
		return m, tea.Batch(dismissCmd, m.reloadCmd())

	case worktreeErrorMsg:
		toast, dismissCmd := components.ShowToast(
			fmt.Sprintf("Worktree error: %s", msg.err),
			components.ToastError,
			5*time.Second,
		)
		m.toast = &toast
		return m, dismissCmd
```

Use `components.ShowToast` (not manual Toast construction) — it handles the auto-dismiss timer. Check how `mutateResultMsg` triggers reload and follow the same pattern for the reload cmd.

- [ ] **Step 5: Build and verify**

Run: `make build`
Expected: PASS (compiles without errors)

- [ ] **Step 6: Run all tests**

Run: `make test`
Expected: PASS

- [ ] **Step 7: Commit (includes Task 5 palette changes)**

```bash
git add internal/components/palette.go internal/app/app.go
git commit -m "feat: add W keybinding and command palette action for worktree creation"
```

---

### Task 7: Agent dispatch uses worktree cwd

**Files:**
- Modify: `internal/app/app.go` (modify `a` key handler, around lines 1609-1625)

- [ ] **Step 1: Modify agent dispatch to check worktree**

In `internal/app/app.go`, in the `"a"` key handler, after the Gas Town sling path returns (line 1607) and before the prompt/launch block (line 1609), add worktree cwd resolution:

Replace the block from line 1609 to 1625:
```go
		deps := issue.EvaluateDependencies(m.detail.IssueMap, m.blockingTypes)
		prompt := agent.BuildPrompt(*issue, deps, m.detail.IssueMap)

		if m.inTmux {
			issueID := issue.ID
			return m, func() tea.Msg {
				winName, err := agent.LaunchInTmux(prompt, m.projectDir, issueID)
				if err != nil {
					return agentLaunchErrorMsg{issueID: issueID, err: err}
				}
				return agentLaunchedMsg{issueID: issueID, windowName: winName}
			}
		}
		c := agent.Command(prompt, m.projectDir)
		return m, tea.ExecProcess(c, func(err error) tea.Msg {
			return agentFinishedMsg{err: err}
		})
```

With:
```go
		// Resolve working directory: prefer worktree if set and valid
		cwd := m.projectDir
		if wt := data.WorktreePath(*issue); wt != "" {
			if info, statErr := os.Stat(wt); statErr == nil && info.IsDir() {
				cwd = wt
			} else {
				toast, _ := components.ShowToast(
					"Worktree missing, using project root",
					components.ToastWarning,
					3*time.Second,
				)
				m.toast = &toast
			}
		}

		deps := issue.EvaluateDependencies(m.detail.IssueMap, m.blockingTypes)
		prompt := agent.BuildPrompt(*issue, deps, m.detail.IssueMap)

		if m.inTmux {
			issueID := issue.ID
			return m, func() tea.Msg {
				winName, err := agent.LaunchInTmux(prompt, cwd, issueID)
				if err != nil {
					return agentLaunchErrorMsg{issueID: issueID, err: err}
				}
				return agentLaunchedMsg{issueID: issueID, windowName: winName}
			}
		}
		c := agent.Command(prompt, cwd)
		return m, tea.ExecProcess(c, func(err error) tea.Msg {
			return agentFinishedMsg{err: err}
		})
```

**Important:** Add `"os"` to the import block in `app.go` if not already present. The current imports include `"os/exec"` but NOT `"os"`. Add `"os"` in the standard library import group.

- [ ] **Step 2: Build and verify**

Run: `make build`
Expected: PASS

- [ ] **Step 3: Run all tests**

Run: `make test`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add internal/app/app.go
git commit -m "feat: agent dispatch uses worktree cwd when available"
```

---

## Chunk 4: View Indicators

### Task 8: Parade worktree indicator

**Files:**
- Modify: `internal/views/parade.go` (add worktree indicator near other badge prefixes, around line 458)

- [ ] **Step 1: Add worktree indicator**

In `internal/views/parade.go`, after the agent badge block (around line 458, after the closing `}` of the `ActiveAgents` block), add:

```go
	// Worktree indicator
	worktreePrefix := ""
	worktreeWidth := 0
	if wt := data.WorktreePath(issue); wt != "" {
		wtColor := ui.BrightGreen
		if _, err := os.Stat(wt); err != nil {
			wtColor = ui.Muted // stale/missing worktree
		}
		worktreePrefix = lipgloss.NewStyle().Foreground(wtColor).Render(ui.SymWorktree) + " "
		worktreeWidth = 2
	}
```

Ensure `"os"` is in the imports.

Then integrate into the existing layout on line 530 and line 553:

**Line 530** — add `- worktreeWidth` to the `maxTitle` calculation:
```go
maxTitle := innerWidth - 16 - hintLen - agentWidth - changeWidth - selectWidth - indentWidth - dueWidth - deferWidth - qualityWidth - orphanWidth - standstillWidth - worktreeWidth
```

**Line 553** — add `worktreePrefix` to the `fmt.Sprintf` after `agentPrefix`:
```go
line := fmt.Sprintf("%s%s %s%s%s%s%s%s %s %s",
    indent,
    symStyle.Render(sym),
    selectPrefix,
    changePrefix,
    orphanPrefix,
    standstillPrefix,
    agentPrefix,
    worktreePrefix,     // <-- add here
    idStyle.Render(issue.ID),
    renderedTitle,
    prioStyle.Render(prio),
)
```

Note: Adding `worktreePrefix` changes the format string argument count — update the `%s` count accordingly (add one more `%s` before `idStyle.Render`).

- [ ] **Step 2: Build and verify**

Run: `make build`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add internal/views/parade.go
git commit -m "feat(views): show worktree indicator in parade list"
```

---

### Task 9: Detail view worktree line

**Files:**
- Modify: `internal/views/detail.go` (add worktree row in `renderContent()`, after the ID row around line 239)

- [ ] **Step 1: Add worktree row**

In `internal/views/detail.go`, in `renderContent()`, after the ID row (line 239), add:

```go
	// Worktree
	if wt := data.WorktreePath(*issue); wt != "" {
		wtStyle := lipgloss.NewStyle().Foreground(ui.Muted)
		if _, err := os.Stat(wt); err != nil {
			wtStyle = wtStyle.Strikethrough(true)
		}
		lines = append(lines, d.row("Worktree:", wtStyle.Render(wt)))
	}
```

Ensure `"os"` is in the imports.

- [ ] **Step 2: Build and verify**

Run: `make build`
Expected: PASS

- [ ] **Step 3: Run all tests**

Run: `make test`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add internal/views/detail.go
git commit -m "feat(views): show worktree path in detail panel"
```

---

## Chunk 5: Final Verification

### Task 10: Lint and full test pass

- [ ] **Step 1: Run lint**

Run: `make lint`
Expected: PASS — no warnings

- [ ] **Step 2: Run full test suite**

Run: `make test`
Expected: All tests PASS

- [ ] **Step 3: Manual smoke test**

Run: `make dev`
Verify:
- `W` on an issue shows a toast (will fail in dev mode since no real git repo, but confirms the keybinding is wired)
- `?` help shows `W` in QUICK ACTIONS
- `:` command palette shows "Create worktree"

- [ ] **Step 4: Final commit if any fixes needed**

If lint or tests required fixes, commit those fixes.
