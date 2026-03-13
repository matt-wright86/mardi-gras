# Epic Filter — Design Spec

## Overview

Add the ability to filter the parade list to show only issues belonging to a specific epic. Users can select an epic via a fuzzy-searchable picker (triggered by `e` keybinding or command palette), and the parade narrows to that epic and all its descendants.

## State Changes

### App Model (`internal/app/app.go`)

Add two fields:

- `epicFilter string` — the selected epic's issue ID (empty = no filter active)
- `epicPicking bool` — whether the epic picker overlay is currently open

### Epic Picker

Reuse the existing `Palette` component pattern — a fuzzy-searchable list populated with all issues where `IssueType == TypeEpic` from the current dataset.

Each entry displays: epic ID + title.

When an epic filter is already active, prepend a "Clear epic filter" entry at the top of the list.

### Trigger

- **Keybinding `e`**: When no epic filter active → open picker. When epic filter active → clear the filter (quick toggle).
- **Palette action `ActionFilterByEpic`**: Always opens the picker (with clear option if filter is active).

## Filtering Logic

### New Helper: `FilterByEpic`

Location: `internal/data/filter.go`

```go
func FilterByEpic(issues []Issue, epicID string) []Issue
```

Keep issues where:
- `issue.ID == epicID`, OR
- `strings.HasPrefix(issue.ID, epicID+".")`

This captures the epic itself and all descendants at any depth.

### Integration into `rebuildParade()`

Apply epic filter early in the pipeline, before text filter and focus mode:

```
all issues
  → FilterByEpic (if epicFilter != "")
  → FilterIssuesWithHighlights (text filter)
  → FocusFilter (if focusMode)
  → GroupByParade
  → NewParadeWithData
```

This lets text filter and focus mode compose naturally on top of the epic-scoped set.

## UI Changes

### Header Badge (`internal/components/header.go`)

When `epicFilter != ""`, render a badge after the existing count badges:

```
epic:mg-007
```

Styled with the theme's epic color to distinguish it from other badges.

### Footer Bindings (`internal/components/footer.go`)

Add `e` to the parade keybinding hints (label: "epic").

### Palette Command (`internal/components/palette.go`)

Add `ActionFilterByEpic` to the action enum and a corresponding command:

```
Name: "Filter by epic"
Desc: "Show only issues under a specific epic"
Key:  "e"
Action: ActionFilterByEpic
```

Always included in `buildPaletteCommands()` (epics are a core issue type, not gated on Gas Town).

## Epic Picker Component

The picker is a lightweight overlay similar to the command palette:

- Text input at top for fuzzy search
- Scrollable list of epic entries (ID + title)
- `enter` selects, `esc` cancels
- Fuzzy matching on ID + title (reuse `github.com/sahilm/fuzzy`)

Implementation: Add an `EpicPicker` model in `internal/components/` or reuse `Palette` with a different command set. Since the palette already supports fuzzy search over a list of `{Name, Desc}` entries, the simplest approach is to build epic entries as `PaletteCommand` items with `Action: ActionSelectEpic` (or similar), and reuse the palette in a "epic picker mode."

### Message Flow

1. User presses `e` (or selects palette action) → app sets `epicPicking = true`, builds epic command list, opens palette
2. User selects an epic → palette returns `PaletteResult{Action: ActionSelectEpic}` with the epic ID
3. App sets `epicFilter = selectedEpicID`, calls `rebuildParade()`
4. User presses `e` again → app clears `epicFilter`, calls `rebuildParade()`

## Files Changed

| File | Change |
|------|--------|
| `internal/data/filter.go` | Add `FilterByEpic()` function |
| `internal/data/filter_test.go` | Tests for `FilterByEpic()` |
| `internal/app/app.go` | Add `epicFilter`/`epicPicking` state, `e` key handler, epic picker flow, integrate into `rebuildParade()` |
| `internal/components/palette.go` | Add `ActionFilterByEpic` (and possibly `ActionClearEpicFilter`) to action enum, add command to `buildPaletteCommands()` |
| `internal/components/header.go` | Render epic filter badge when active |
| `internal/components/footer.go` | Add `e` to parade keybinding hints |

## Testing

- `FilterByEpic`: unit tests for exact match, prefix match, no match, empty epicID passthrough
- Integration: verify `rebuildParade()` applies epic filter before text filter and focus mode (existing test patterns in `filter_test.go`)
