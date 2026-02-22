# Contributing to Mardi Gras

Thanks for your interest in making the parade better! This guide covers everything you need to get started.

## Prerequisites

- **Go 1.24+** ([install](https://go.dev/doc/install))
- **Git**
- **golangci-lint** for linting ([install](https://golangci-lint.run/welcome/install/))
- A Beads project, or use the included `testdata/sample.jsonl`

## Getting Started

```bash
git clone https://github.com/matt-wright86/mardi-gras.git
cd mardi-gras
make build
```

Run against the included sample data:

```bash
make dev
```

This builds the `mg` binary and launches it with `testdata/sample.jsonl`.

## Development Commands

```bash
make build        # compile the mg binary
make run          # build + run (auto-detects .beads/issues.jsonl)
make dev          # build + run with sample data
make test         # go test ./...
make lint         # golangci-lint run ./...
make fmt          # go fmt ./...
make tidy         # go mod tidy
make clean        # remove binary and dist/
```

CI runs tests with `-race` and lints with the same `.golangci.yml` config, so run `make test` and `make lint` locally before pushing.

## Project Structure

```
cmd/mg/main.go        Entry point (flags, path resolution, bootstrap)

internal/
  app/                Root BubbleTea model (lifecycle, routing, layout)
  data/               Domain types, JSONL parsing, filtering, file watcher
  views/              Parade (left pane) and Detail (right pane)
  components/         Header, Footer, Help overlay, Float utility
  agent/              Claude Code integration and tmux dispatch
  tmux/               tmux status line widget (--status mode)
  ui/                 Theme colors, lipgloss styles, Unicode symbols

testdata/             Sample JSONL for development
docs/                 Architecture docs and screenshots
```

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for a deeper walkthrough of the data flow, BubbleTea model structure, and dependency graph between packages.

## How to Contribute

### Reporting Bugs

Open a GitHub issue with:

- What you expected vs. what happened
- Steps to reproduce
- Terminal emulator and OS
- Output of `mg --version`

### Suggesting Features

Open an issue describing the feature and why it would be useful. The [README](README.md#possible-future-ideas) lists some ideas we've been thinking about.

### Submitting a Pull Request

1. Fork the repo and create a branch from `main`.
2. Make your changes.
3. Add or update tests if the change affects behavior.
4. Run `make fmt && make lint && make test` and fix any issues.
5. Write clear commit messages that explain the *why*, not just the *what*.
6. Open a PR against `main`.

Keep PRs focused -- one feature or fix per PR makes review faster for everyone.

## Code Conventions

- **Formatting**: `gofmt` (enforced by CI).
- **Linting**: golangci-lint with the config in `.golangci.yml`.
- **Naming**: follow standard Go conventions. Exported names should be clear without a package prefix.
- **Errors**: return errors rather than panicking. Use `fmt.Errorf` with `%w` for wrapping.
- **Tests**: table-driven tests where appropriate. Test files live alongside the code they test.
- **Dependencies**: Mardi Gras intentionally has a small dependency footprint (just the Charmbracelet toolkit). Propose new dependencies in the PR description with a rationale.

## Architecture Notes

Mardi Gras follows the [Elm Architecture](https://guide.elm-lang.org/architecture/) via BubbleTea:

- **Model** holds all state in `app.Model`.
- **Update** routes messages (key presses, file changes, agent events) to handlers.
- **View** composes sub-models (parade, detail, header, footer) into the final screen.

Key design constraints:

- `data` and `ui` have no internal dependencies beyond the standard library and lipgloss -- keep them that way.
- No package imports `app` -- it is the root.
- Single binary, no runtime dependencies. Cross-compiles via GoReleaser.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
