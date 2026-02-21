# ⚜ Mardi Gras

A TUI dashboard for the [Beads](https://github.com/anthropics/beads) issue tracker, themed as a Mardi Gras parade.

Your issues are **floats** rolling through the parade route — some rolling, some lined up, some stalled, and some already past the stand.

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

# Specify a path
mg --path /path/to/.beads/issues.jsonl
```

## Keybindings

| Key | Action |
|-----|--------|
| `j` / `k` | Navigate up/down |
| `tab` | Switch between parade and detail panes |
| `enter` | Focus detail pane |
| `esc` | Back to parade pane |
| `c` | Toggle closed issues |
| `g` / `G` | Jump to top / bottom |
| `q` | Quit |

## Parade Vocabulary

| Term | Meaning |
|------|---------|
| **Rolling** ● | In progress — actively being worked on |
| **Lined Up** ♪ | Open and ready — no blockers |
| **Stalled** ⊘ | Open but blocked by dependencies |
| **Past the Stand** ✓ | Closed/completed |
| **Float** | An issue card or group |
| **Bead String** | The decorative separator bar |
| **Krewe** | Your team of contributors |

## Stack

- **Go** — fast, single binary
- **BubbleTea** — Elm Architecture TUI framework
- **Lipgloss** — terminal styling
- **Bubbles** — pre-built TUI components

## License

MIT
