# Draft Response: #2372 — Dolt server port cross-project data leakage

**Status**: Draft — do not post without approval
**Issue**: https://github.com/steveyegge/beads/issues/2372

---

## Draft comment

We hit this exact scenario running [mardi-gras](https://github.com/quietpublish/mardi-gras) (a BubbleTea TUI for beads). Two projects — mardi-gras and voice-vault — were both using port 3307, and `bd list --json` in one project silently returned the other project's issues. Took a while to figure out because the data looked plausible.

### What we observed

1. voice-vault had a Dolt server on 3307 (started manually pre-v0.57.0)
2. mardi-gras's config had `DoltServerPort` unset, so it fell through to the Gas Town default (3307)
3. Running `bd list` in mardi-gras connected to voice-vault's server — returned issues with the wrong prefix
4. No error, no warning. The only hint was issue IDs with a prefix we didn't recognize

### What fixed it for us

```bash
# Stop the conflicting server
bd dolt stop

# Re-initialize with an explicit prefix (generates a DerivePort-based port)
bd init --force --prefix mg

# Verify the prefix is set
bd config set issue_prefix mg
```

After re-init, `DerivePort` assigned a hash-based port unique to the absolute path, which resolved the collision. The key was that `bd init --force` on v0.58.0 uses `DerivePort` by default rather than falling back to 3307.

### Confirming the gitignore gap

The issue correctly identifies that `dolt-server.port` was missing from `.beads/.gitignore`. On v0.58.0, our gitignore *does* include it:

```
dolt-server.port
```

So this appears to have been fixed in the v0.57.0→v0.58.0 cycle. If you're on an older version, you can add it manually:

```bash
echo "dolt-server.port" >> .beads/.gitignore
git add .beads/.gitignore && git commit -m "fix: gitignore dolt-server.port"
```

### Supporting the proposed fix

+1 on both layers of the proposed fix:

**Layer 1 (port isolation)**: `DerivePort` is already the right default — the issue is that legacy configs and git-tracked values override it. Deprecating `DoltServerPort` in `metadata.json` with a warning would catch the most common case.

**Layer 2 (database identity verification)**: This would have caught our issue instantly. A project UUID check on connection would turn a silent data leak into a loud, actionable error. Even a simple check like "does the issue prefix in the database match my configured prefix?" would help — that's essentially what tipped us off (wrong prefix in the output).

### Workaround for affected users

If you're seeing cross-project contamination right now:

1. Check which port each project uses: `bd sql "SELECT @@port"` or look in `.beads/metadata.json`
2. If two projects share a port, re-init the newer one: `bd init --force --prefix <your-prefix>`
3. Verify with `bd list` that you see the correct issues
4. Add `dolt-server.port` to `.beads/.gitignore` if it's not already there
5. Check for stale port files: `git ls-files .beads/dolt-server.port` — if tracked, remove it: `git rm --cached .beads/dolt-server.port`

---

## Notes (not for posting)

- Our experience was on macOS, bd v0.56.1→v0.58.0, same as the reporter
- The `bd init --force` path loses existing issues — we lost created beads and had to recreate them. Should mention this risk? Decided not to — the reporter is describing a design issue, not asking for a recovery guide.
- The gitignore fix appears to be in v0.58.0 already (commit unknown, but our fresh init has it)
- @PabloLION's impact data comment is excellent — validates the systemic nature of the problem
- Related to our previous triage drafts for #2098 and #2030
