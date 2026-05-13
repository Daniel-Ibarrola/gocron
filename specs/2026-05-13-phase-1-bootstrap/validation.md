# Phase 1 — Project bootstrap: Validation

The phase is mergeable when **every** item below passes.

## Build & static checks

- [ ] `go build ./...` exits 0 and produces a `gocron` binary.
- [ ] `go vet ./...` exits 0 with no warnings.
- [ ] `go mod tidy` produces no diff against the committed `go.mod`
      (`/go.sum` should not exist — no dependencies yet).
- [ ] `go.mod` first line reads `module github.com/dibarrola/gocron`.
- [ ] `go.mod` declares `go 1.26`.

## CLI behavior — happy path

Flags must appear **before** the cron expression. The stdlib `flag` package
stops parsing at the first non-flag argument, so anything after the
expression is ignored. Build once, then run:

```
$ ./gocron "*/5 * * * *"
expression: */5 * * * *
next: 3
json: false
```

```
$ ./gocron --next 5 "*/5 * * * *"
expression: */5 * * * *
next: 5
json: false
```

```
$ ./gocron --next 5 --json "*/5 * * * *"
expression: */5 * * * *
next: 5
json: true
```

- [ ] Each invocation above produces the shown output verbatim (no extra
      lines, no trailing whitespace).
- [ ] Exit status is 0 for all three.

## CLI behavior — error path

```
$ ./gocron
gocron: missing cron expression
$ echo $?
1
```

- [ ] No expression given: error message lands on **stderr** (not stdout).
- [ ] Exit status is non-zero.

## Repository state

- [ ] No `internal/` directory exists (deferred to later phases).
- [ ] No `github.com/fatih/color` import or `go.sum` entry exists.
- [ ] `specs/roadmap.md`:
      - Phase 1 is marked `[x]`.
      - Phase 1 Go version reads `1.26+`, not `1.22+`.
      - The `fatih/color` bullet has moved from Phase 1 to Phase 5.
- [ ] `specs/2026-05-13-phase-1-bootstrap/` contains `plan.md`,
      `requirements.md`, and `validation.md`.

## Ready to merge

When all boxes are checked, push the branch and open a PR against `main`.
