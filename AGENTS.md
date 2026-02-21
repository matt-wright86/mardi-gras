# Repository Guidelines

## Project Structure & Module Organization
- `cmd/mg/main.go`: CLI entrypoint and `--path` flag handling.
- `internal/app`: Bubble Tea root model, keys, and pane orchestration.
- `internal/views`: left parade list + right detail panel rendering.
- `internal/components`: shared header/footer and float blocks.
- `internal/data`: JSONL loading, grouping, dependency/status logic.
- `internal/ui`: symbols, theme palette, and Lipgloss styles.
- `testdata/sample.jsonl`: fixture for tests and local demo runs.
- `docs/`: investigation and product direction notes.

## Beads Data Contract
- Treat `.beads/issues.jsonl` as the source of truth; do not rely on `.beads/.beads.db`.
- Parse JSONL line-by-line and keep reads safe while Beads is running.
- Preserve status semantics: `in_progress` -> Rolling, `open` unblocked -> Lined Up, `open` blocked -> Stalled, `closed` -> Past the Stand.
- Preserve dependency handling for both `blocks` and `parent-child` relationships.
- Optimize for real-world closed-heavy datasets; closed issues should remain collapsible and low-noise by default.

## Build, Test, and Development Commands
- `make build`: build local binary `./mg` from `./cmd/mg`.
- `make run`: build and run using auto-detected `.beads/issues.jsonl`.
- `make run-sample` (or `make dev`): run against `testdata/sample.jsonl`.
- `make test`: execute `go test ./...` across all packages.
- `make fmt`: apply standard Go formatting (`go fmt ./...`).
- `make lint`: run static analysis with `golangci-lint run ./...`.
- `make tidy`: sync module dependencies in `go.mod`/`go.sum`.

## Coding Style & Naming Conventions
- Use idiomatic Go and always format with `make fmt` before committing.
- Keep package boundaries domain-based (`data`, `views`, `components`, `ui`, `app`).
- Exported names use `PascalCase`; unexported helpers use `camelCase`.
- Keep Mardi Gras UI vocabulary and section labels consistent (`ROLLING`, `LINED UP`, `STALLED`, `PAST THE STAND`).
- If keybindings change, update both the in-app footer hints and the README keybinding table in the same PR.

## Testing Guidelines
- Put tests next to implementation as `*_test.go`.
- Name tests `Test<Behavior>` and keep assertions explicit and readable.
- Prefer deterministic tests using fixtures from `testdata/`.
- Run `make test` for all changes; run `make run-sample` to verify `j/k`, `tab`, and `c` behavior.

## Commit & Pull Request Guidelines
- Follow Conventional Commit style used in history (for example `feat:`, `fix:`, `docs:`, `test:`, `chore:`).
- Keep each commit focused on one logical change.
- PRs should include a short summary, validation steps run (commands), and screenshots/GIFs for visible TUI updates.
- Link related issues and call out any follow-up work or known limitations.
