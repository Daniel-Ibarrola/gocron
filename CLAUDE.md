# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```sh
go build -o gocron .      # build the binary
go test ./...             # run all tests
go test -race ./...       # run tests with race detector (used in CI)
go test ./internal/cron/  # run tests for a single package
go vet ./...              # static analysis
gofmt -l .               # check formatting (CI fails on violations)
gofmt -w .               # reformat all files
```

## Architecture

**gocron** is a CLI that explains 5-field cron expressions in plain English. No third-party cron or CLI libraries — stdlib only (intentional practice constraint).

### Data flow

```
main.go  →  cron.ExplainExpression(expr)
               ├── validate(expr)          [validation.go]  character & field-count check
               └── parseExpression(expr)   [parse.go]       tokenise → []fieldSpec
                       └── returns *cronExpression
           →  explainMinutes / explainHours / explainDaysOfMonth / ...  [cron.go]
               └── strings.Join(parts, " ")  →  human string
```

### Key types (`internal/cron/parse.go`)

- `fieldKind` — `fieldWildcard | fieldSingle | fieldRange` (more kinds added per phase)
- `fieldSpec{kind, lo, hi}` — represents one parsed cron field; single values store the same number in both `lo` and `hi`
- `cronExpression{minutes, hours, daysOfMonth, months, daysOfWeek fieldSpec}` — the full parsed expression

### Error types (`internal/cron/`)

- `ValidationError` — bad character or wrong field count (caught before parsing)
- `FieldError` — per-field parse failure (bad integer, out-of-bounds, inverted range)

### Phased development

New cron syntax (steps, lists, etc.) is added one phase at a time — each phase extends **both** the parser and the explainer together with tests for both. The roadmap is in `specs/roadmap.md`. Current branch (`phase-3-ranges`) adds `a-b` range support (Phase 3). Upcoming phases: steps (`*/N`, `a-b/N`), lists, next-run calculator, output formatting.

### Testing conventions

Table-driven tests in `*_test.go` alongside each source file. Tests cover the parser and explainer independently. See `internal/cron/parse_test.go` and `internal/cron/explain_test.go` for existing patterns.
