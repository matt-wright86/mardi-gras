# Contributing to Mardi Gras

Thanks for your interest in contributing! Here's how to get started.

## Prerequisites

- Go 1.24+
- A Beads project (or use the included `testdata/sample.jsonl`)

## Setup

```bash
git clone https://github.com/matt-wright86/mardi-gras.git
cd mardi-gras
make build
```

## Development

```bash
make dev          # Build and run with sample data
make test         # Run all tests
make lint         # Run golangci-lint (install: https://golangci-lint.run/welcome/install)
make fmt          # Format code
```

## Making Changes

1. Fork the repo and create a feature branch from `main`.
2. Write tests for new functionality.
3. Run `make test` and `make fmt` before submitting.
4. Keep commits focused — one logical change per commit.
5. Write commit messages that explain *why*, not just *what*.

## Pull Requests

- Keep PRs small and focused.
- Reference any related issues.
- Ensure CI passes (tests + lint).

## Project Structure

```
cmd/mg/            Entry point
internal/
  app/             BubbleTea model, keybindings, update loop
  agent/           Claude Code agent launch
  components/      Header, footer, help modal, dividers
  data/            Issue types, JSONL loader, file watcher, filtering
  tmux/            tmux status line output
  ui/              Styles and color palette
  views/           Parade and detail panel views
testdata/          Sample JSONL for development
```

## Code Style

- Follow existing patterns in the codebase.
- Use `go fmt` formatting.
- Keep dependencies minimal — the goal is a single binary with no runtime deps.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
