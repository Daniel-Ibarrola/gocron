# Phase 1 — Project bootstrap: Requirements

## Goal

Replace the GoLand placeholder with a real CLI skeleton that can read an
expression and two flags (`--next`, `--json`), echo the parsed values back,
and exit cleanly. This proves the wiring is sound before any cron logic
lands in Phase 2.

## In scope

- Rewrite `go.mod` module path to `github.com/dibarrola/gocron`.
- Replace the GoLand boilerplate in `main.go` with a real CLI entry point.
- Parse one positional argument (the cron expression).
- Parse two flags via stdlib `flag`:
  - `--next int` (default `3`) — number of upcoming run times to show.
  - `--json bool` (default `false`) — emit JSON instead of text.
- On valid invocation, print the three parsed values back to stdout in this
  exact format (one per line, lowercase keys, single space after colon):
  ```
  expression: <raw expression>
  next: <N>
  json: <bool>
  ```
- On missing expression, print a single-line error to **stderr** and exit
  with a non-zero status.
- Update `specs/roadmap.md` to reflect two decisions made during planning:
  - The Go directive in `go.mod` stays at `1.26` (not `1.22+`).
  - The `fatih/color` dependency is deferred from Phase 1 to Phase 5,
    where it is actually used.

## Out of scope (deferred to later phases)

- Any cron parsing, validation, explanation, or scheduling logic.
- The `internal/cron`, `internal/explain`, `internal/schedule`,
  `internal/output` packages — no stubs created in this phase. Each is added
  in the phase that needs it.
- The `fatih/color` dependency — added in Phase 5.
- Color output, JSON output formatting (the `--json` flag is parsed and
  echoed but does not change the output format yet).
- Friendly `--help` text and polished error messages (Phase 7).
- Unit tests — none are required for this phase. The test suite starts in
  Phase 2 alongside the parser.

## Decisions

1. **Go version in `go.mod`** — keep `go 1.26` (matches the local toolchain).
   Roadmap updated from `1.22+` to `1.26+` to match.
2. **Module path** — rewrite from `gocron` to `github.com/dibarrola/gocron`.
3. **Branch name** — `phase-1-bootstrap`.
4. **No `internal/` stub directories** — git won't track empty dirs, and
   adding `doc.go` placeholders is noise. Each package appears in its phase.
5. **Echo format** — three lines, `key: value`, lowercase keys. This is a
   throwaway debug echo; Phase 5 replaces it with the real formatter.
6. **Missing-expression behavior** — one-line stderr error
   (`gocron: missing cron expression`) and `os.Exit(1)`. This satisfies the
   "clear errors" goal in `mission.md` without prematurely building the
   full error framework planned for Phase 7.
7. **`fatih/color` deferred** — adding an unused dependency in Phase 1 just
   to satisfy a roadmap bullet is wasteful. Moved to Phase 5.

## Context

- `specs/mission.md` — `gocron` is a standalone, offline, scriptable cron
  explainer; v1 is read-only (does not execute jobs).
- `specs/tech-stack.md` — stdlib only for core logic; no CLI framework;
  `flag` package handles arguments.
- Current repo state at the start of this phase: GoLand placeholder
  `main.go`, `go.mod` declares `module gocron` and `go 1.26`, no other
  source files.
