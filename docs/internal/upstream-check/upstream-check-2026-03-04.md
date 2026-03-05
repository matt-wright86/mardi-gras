# Upstream Check — 2026-03-04

## TL;DR

New releases on both repos: Beads v0.58.0 (massive — 76 features, 200+ fixes) and Gas Town v0.10.0 (150 commits — reliability, telemetry, mTLS, cursor hooks). One **potential compatibility risk**: Gas Town removed the `hook_bead` agent bead slot (`fa9dc28`), but mg already enriches `HookBead` from hook data, so no immediate break. Several significant feature opportunities: circuit breakers, cursor hooks support, agent memory commands, heartbeat v2.

## Current baseline

- mg version: v0.5.0 (BubbleTea v2 migration just landed on main)
- Beads: v0.57.0 → **v0.58.0** released 2026-03-03, main has ~15 commits ahead
- Gas Town: v0.9.0 → **v0.10.0** released 2026-03-03, main has ~5 non-backup commits ahead
- Previous check: [upstream-check-2026-03-02.md](upstream-check-2026-03-02.md)

## Breaking changes

### Gas Town: `hook_bead` agent bead slot removed (potential risk)

**Commit**: `fa9dc28` — "Remove agent bead hook slot: use direct bead tracking"

The `HookBead` field on agent beads is no longer written by sling/done/unsling. The work bead's `assignee` + `status=hooked` is now authoritative.

**Impact on mg**: mg's `normalizeStatus()` in `status.go:137-139` already enriches `HookBead` from hook data (`h.Molecule`), not just the agent bead. As long as `gt status --json` still includes hooks with molecule references, mg will continue to work. However, if a future gt version removes the hooks array from status output, `HookBead` will be empty everywhere.

**Risk level**: Low (currently works), but should be monitored. The `AgentForIssue()` and `ActiveAgentMap()` functions both depend on `HookBead`.

**Files affected**: `internal/gastown/status.go`, `internal/views/gastown.go:701-702`, `internal/gastown/problems.go:51`

**Action**: No code change needed now. Add a note to check hook data presence after upgrading to gt v0.10.0.

### Beads: No breaking changes for mg

v0.58.0 is a consolidation release. `bd list --json` output format unchanged. All new features are additive.

## Feature opportunities

### 1. Polecat circuit breaker status (Gas Town v0.10.0)

**What**: Per-bead respawn circuit breaker prevents spawn storms. Circuit breaker state is tracked per-agent.

**How mg could use it**: Surface circuit breaker status in the Gas Town panel agent roster — show when a polecat is in backoff/tripped state. Would help diagnose "why isn't my issue being worked on?" questions.

**Effort**: Small. Read circuit breaker state from `gt status --json` if exposed, or from problems detection.

### 2. Merge queue pre-verification (Gas Town v0.10.0)

**What**: New pre-verify step in molecule polecat work formula. `gt done --pre-verified` for fast-path merges.

**How mg could use it**: Show pre-verification status on convoy items in Gas Town panel. Could add a "pre-verified" badge to MRs.

**Effort**: Small. Parse from convoy/MQ status data.

### 3. OTel telemetry data (Gas Town v0.10.0)

**What**: Agent session, mail, and token usage instrumentation via OpenTelemetry.

**How mg could use it**: If OTel data is queryable, mg could show token usage per agent, session duration, and mail metrics in the Gas Town panel.

**Effort**: Large. Would need OTel query integration or a gt command that exposes the data.

### 4. Heartbeat v2 with agent-reported state (Gas Town v0.10.0)

**What**: Agents now self-report state via heartbeat. More accurate than polling-based inference.

**How mg could use it**: Already benefits from this — `gt status --json` agent states will be more accurate. No code change needed.

**Effort**: None (automatic improvement).

### 5. Cursor hooks support (Gas Town post-v0.10.0)

**What**: `feat/cursor-hooks-support` branch merged — Gas Town now supports Cursor IDE agent hooks alongside Claude Code.

**How mg could use it**: When running in Cursor, mg's agent dispatch could detect and use Cursor-specific hooks. Currently only Claude Code dispatch is supported.

**Effort**: Medium. Would need Cursor detection in `internal/agent/` and hook format handling.

### 6. `bd doctor --agent` mode (Beads v0.58.0)

**What**: Special doctor mode for AI agents — diagnostics without interactive prompts.

**How mg could use it**: Run `bd doctor --agent` in the background and surface warnings in the problems overlay.

**Effort**: Small. Parse doctor output, add to problems detection.

### 7. `bd show --long` for full-detail output (Beads v0.58.0)

**What**: `--long` flag for full-detail issue view.

**How mg could use it**: Use in detail panel to get richer issue data (acceptance criteria, metadata, etc.).

**Effort**: Small. Switch `bd show` call to include `--long` where appropriate.

### 8. `bd backup restore` (Beads v0.58.0)

**What**: Bootstrap from JSONL backup.

**How mg could use it**: Informational — good for recovery scenarios. Could add to help overlay.

**Effort**: None for mg.

### 9. Daemon removal PR (Beads — open PR #2356)

**What**: `feat(bd): remove daemon infrastructure` — removes the bd daemon entirely.

**How mg could use it**: mg already doesn't use the daemon (uses SourceCLI). This validates mg's architecture choice. No action needed.

**Effort**: None.

### 10. Legacy SQLite/JSONL removal PR (Beads — open PR #2350)

**What**: v0.60.0 cleanup removing dead SQLite/JSONL/embedded-Dolt code.

**How mg could use it**: Informational — confirms Dolt-only future. mg's JSONL loader may become fully dead code.

**Effort**: Small (remove dead code when this ships).

## Post-release commits (since v0.58.0 / v0.10.0)

### Beads (post v0.58.0, 2026-03-03/04)

- `b4b586d` feat(doctor): allow suppressing specific warnings via config
- `5ea0eaa` feat: add OpenCode recipe to bd setup
- `d9a719e` fix(doltserver): idle-monitor kills itself via Stop()
- `0e286c0` fix: restore DerivePort as standalone default in DefaultConfig
- `1b921b9` fix(dep): add cross-prefix routing to dep commands
- `f4c575e` refactor: consolidate ExtractPrefix into types package
- Various test refactoring (container-native API, branch-per-test)
- Docs: nvim-beads community tool, QUICKSTART fixes, GIT_INTEGRATION update

### Gas Town (post v0.10.0, 2026-03-03/04)

- `fa9dc28` Remove agent bead hook slot (see breaking changes above)
- `aa7dd7e` Merge cursor-hooks-support (Cursor IDE integration)
- `76ef3fa` refactor: extract shared IsAutonomousRole into hookutil
- `cdb2f04` fix(guard): portable reverse-file for macOS
- `2a6a60f` fix(convoy): add omitempty to strandedConvoyInfo.CreatedAt

## Recommended actions

| # | Action | Priority | Effort | Files |
|---|--------|----------|--------|-------|
| 1 | Monitor `hook_bead` field after gt upgrade to v0.10.0 | high | small | `internal/gastown/status.go` |
| 2 | Upgrade bd to v0.58.0 | medium | small | local binary |
| 3 | Upgrade gt to v0.10.0 | medium | small | local binary |
| 4 | Surface `bd doctor --agent` warnings in problems overlay | medium | small | `internal/views/problems.go`, `internal/gastown/problems.go` |
| 5 | Use `bd show --long` in detail panel | low | small | `internal/data/mutate.go` or detail fetch |
| 6 | Surface circuit breaker state in agent roster | low | small | `internal/views/gastown.go` |
| 7 | Carry forward: render dogs in agent roster | low | medium | `internal/views/gastown.go` |
| 8 | Carry forward: surface `bd show --current` in header/footer | low | medium | `internal/components/header.go` |

## Open PRs to watch

### Beads
- **#2350**: Remove legacy SQLite/JSONL/embedded-Dolt dead code (v0.60.0 cleanup)
- **#2356**: Remove daemon infrastructure
- **#2354**: Detect and offer restore from backup when DB missing
- **#2353**: Remove dead 3-way merge engine remnants
- **#2361**: Skip tombstone entries in `bd init --from-jsonl`

### Gas Town
- **#2358**: Decouple Agent Execution from Presentation (ACP) — significant architecture change
- **#2338**: Reload prefix registry on heartbeat to prevent ghost sessions
- **#2360**: Schema evolution support for `gt wl sync`
