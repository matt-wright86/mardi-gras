# ⚜ Mardi Gras

**Your Beads issues deserve a parade, not a spreadsheet.**

`bd list` is fine. You scan the output, count the open ones in your head, try to remember which tasks are blocked, and hope you don't miss anything. It works. But it doesn't _feel_ like anything.

Mardi Gras gives your [Beads](https://github.com/anthropics/beads) issues a parade route. Your in-progress work **rolls**. Your ready tasks are **lined up**. Your blocked work is **stalled** on the route. And everything you've finished? It's already **past the stand** — you threw the beads, the crowd cheered, and it's behind you.

One binary. No config. Just `mg` in any project with a `.beads/` directory.

```
⚜ MARDI GRAS ⚜                 2 rolling | 3 lined up | 1 stalled | 10 total
●─◆─●─◆─●─◆─●─◆─●─◆─●─◆─●─◆─●─◆─●─◆─●─◆─●─◆─●─◆─●─◆─●─◆─●─◆─●─◆─●─◆─●─◆
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 ROLLING ●                      │  Deploy authentication service
  > ● vv-001 Deploy auth    P0  │
    ● vv-002 Fix CI pipeline P1 │  Status:   ● Rolling (in_progress)
                                │  Type:     feature
 LINED UP ♪                     │  Priority: P0 (Critical)
    ♪ vv-003 Add monitoring  P2 │  Owner:    dev@example.com
    ♪ vv-004 Update docs     P3 │  Age:      3 days
    ♪ vv-005 Refactor utils  P4 │
                                │  DEPENDENCIES
 STALLED ⊘                      │  ● blocks → vv-006 (Write tests for...)
    ⊘ vv-006 Write tests     P2 │
                                │
 ▶ PAST THE STAND (4 issues)    │
   [press c to expand]          │
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 / filter  j/k navigate  tab switch pane  c toggle closed  q quit
```

## Why this exists

Beads solves agent context beautifully. But `bd list` wasn't built for humans doing visual triage across dozens of issues. The community knows this — 8+ independent projects have been built trying to solve Beads visibility, from web UIs to desktop apps to other TUIs.

Mardi Gras takes a different angle: instead of building another kanban board, we built a **parade**. Because if you're going to stare at your task list every day, it should at least make you smile.

## Install

```bash
go install github.com/matt-wright86/mardi-gras/cmd/mg@latest
```

Or build from source:

```bash
git clone https://github.com/matt-wright86/mardi-gras.git
cd mardi-gras
make build
```

## Usage

```bash
# Auto-detect .beads/issues.jsonl in current directory
mg

# Point at a specific project
mg --path /path/to/.beads/issues.jsonl
```

Mardi Gras reads your `.beads/issues.jsonl` directly — no daemon, no database, no config file. It polls for changes automatically, so if an agent updates an issue while you're watching, the parade reshuffles in real time.

## Live Updates

- The app watches your JSONL file on a short polling interval.
- External edits (agents, scripts, `bd` commands) are picked up automatically.
- Current view state is preserved on refresh where possible (selection, closed section toggle, active filter query).

## Keybindings

| Key       | Action                                              |
| --------- | --------------------------------------------------- |
| `?`       | Toggle help overlay                                 |
| `j` / `k` | Navigate up/down                                    |
| `/`       | Enter filter mode                                   |
| `esc`     | Clear filter (in filter mode) / Back to parade pane / Close help overlay |
| `enter`   | Apply filter (in filter mode) / Focus detail pane   |
| `tab`     | Switch between parade and detail panes              |
| `c`       | Toggle closed issues                                |
| `g` / `G` | Jump to top / bottom                                |
| `q`       | Quit (or close help overlay if open)                |
| `ctrl+c`  | Quit (global)                                       |

## Help Overlay

Press `?` from anywhere (including while filtering) to open the command reference overlay.

- While help is open, list/detail/filter interactions are paused.
- Press `?`, `esc`, or `q` to close help and return to the previous mode.

## Filtering

Press `/` and the bottom bar becomes a query input.

- `enter`: keep the query applied and return to list navigation.
- `esc`: clear the query and exit filter mode.
- Multiple terms use `AND` semantics (all terms must match).

Supported query forms:

- Free text: `deploy auth` (matches issue ID and title)
- Type token: `type:bug`, `type:feature`, `type:task`, `type:chore`, `type:epic`
- Priority shorthand: `p0` to `p4`
- Priority token: `priority:0` to `priority:4`, or `priority:critical|high|medium|low|backlog`

Examples:

```text
type:feature p1 deploy
priority:high auth
vv-006
```

## The Parade

Every Beads issue maps to a spot on the parade route:

| On the Route         | What It Means                         |
| -------------------- | ------------------------------------- |
| **Rolling** ●        | In progress — the float is moving     |
| **Lined Up** ♪       | Open and unblocked — waiting its turn |
| **Stalled** ⊘        | Blocked by a dependency               |
| **Past the Stand** ✓ | Done — beads have been thrown         |

Closed issues are collapsed by default (because in any real project, 90%+ of your issues are closed). Press `c` to expand them.

## Built with

- [BubbleTea](https://github.com/charmbracelet/bubbletea) — Elm Architecture for the terminal
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — CSS-like styling (the purple, gold, and green)
- [Bubbles](https://github.com/charmbracelet/bubbles) — viewport scrolling

Single binary, no runtime dependencies. Cross-compiles to Linux, macOS, and Windows via [GoReleaser](https://goreleaser.com).

## Contributing

Mardi Gras is early. The parade route is laid, the floats are rolling, but there's plenty of room for more krewes.

## License

MIT
