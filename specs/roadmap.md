# Roadmap

Phases are intentionally small. Each phase should leave the project in a
working, committable state.

Each functionality phase (2 onward) layers one new piece of cron syntax onto
the parser **and** the plain-English explainer, with tests for both.

Status legend: `[ ]` not started · `[~]` in progress · `[x]` done

---

## Phase 1 — Project bootstrap `[x]`

- Initialize `go.mod` (`github.com/dibarrola/gocron`, Go 1.26+)
- `main.go`: read `os.Args`, parse `--next` and `--json` flags, print raw
  expression back to confirm wiring works
- Directory layout:
  ```
  main.go
  internal/
    cron/      # parser + explainer live here for now
    schedule/  # next-run calculator lives here (later)
    output/    # formatting (text + JSON) lives here (later)
  ```

## Phase 2 — Simple expressions, CI, and README `[x]`

The minimum end-to-end slice: parse and explain expressions that only use
wildcards (`*`) and fixed integer values across all five fields, and gate
future work behind CI.

- **Parser** (`internal/cron`):
  - Split a 5-field expression: `minute hour dom month dow`
  - Each field is either `*` (wildcard) or a single integer
  - Surface a `ValidationError` for malformed input (bad field count,
    disallowed characters) and a `FieldError` for unparseable tokens
- **Explainer** (`ExplainExpression`):
  - `* * * * *` → "every minute"
  - `5 * * * *` → "at minute 5"
  - `30 6 * * *` → "at 06:30"
  - Combinations with day-of-month, month, day-of-week (named, e.g.
    "Monday", "September")
- **Tests** for parser and explainer covering each shape above
- **CI** — GitHub Actions workflow at `.github/workflows/ci.yml`,
  triggered on `push` and `pull_request`:
  1. `gofmt -l .` (fail on unformatted files)
  2. `go vet ./...`
  3. `go test -race ./...`
  - One job, sequential steps; Go version pinned to match `go.mod`
  - No linter yet — add in a later phase if churn warrants it
- **README** at repo root: one-paragraph description, install
  instructions (`go install ...`), and a usage example showing the
  simple-expression explainer working end-to-end

## Phase 3 — Ranges `[x]`

Extend parser and explainer to handle `a-b` per field.

- Parser: accept `1-5`, validate `a <= b` and both within the field's bounds
- Explainer phrasing, e.g.:
  - `0 9-17 * * *` → "at minute 0 past hours 9 through 17"
  - `* * * * 1-5` → "every minute on Monday through Friday"
- Tests: valid ranges, inverted ranges, out-of-bounds ranges

## Phase 4 — Steps `[ ]`

Extend parser and explainer to handle `*/N` and `a-b/N` per field.

- Parser: accept `*/5`, `1-30/5`; validate step > 0 and that the base
  (wildcard or range) is valid
- Explainer phrasing, e.g.:
  - `*/5 * * * *` → "every 5 minutes"
  - `0 9-17/2 * * *` → "at minute 0 past every 2nd hour from 9 through 17"
- Tests: bare step on wildcard, step on range, zero/negative step

## Phase 5 — Lists `[ ]`

Extend parser and explainer to handle comma-separated combinations of any
previously-supported token.

- Parser: accept `1,2,3`, `0,15,30,45`, `1-5,10,*/15` — each comma element
  is independently any of: single value, range, step
- Explainer phrasing, e.g.:
  - `0,15,30,45 * * * *` → "at minutes 0, 15, 30, and 45"
  - `0 9,12,17 * * 1-5` → "at minute 0 past hours 9, 12, and 17 on Monday through Friday"
- Tests: pure lists, mixed-token lists, duplicate entries, empty elements

## Phase 6 — Next-N run calculator `[ ]`

- Given a parsed expression and a `time.Time` start, return the next N
  matching times
- Strategy: advance minute-by-minute from `start+1m`, check each candidate
  against all five fields
- Handle month/dom edge cases (e.g. Feb 30 never matches)
- Cap iteration to avoid infinite loops on expressions that can never fire
  (e.g. `* * 31 2 *`)
- Unit tests: known expressions with known next-run times

## Phase 7 — Output formatting `[ ]`

- Add `github.com/fatih/color` dependency (deferred from Phase 1; first
  consumed here)
- Plain-text formatter (colored via `fatih/color` when TTY, plain otherwise):
  - Header line with the raw expression
  - Meaning line
  - Numbered list of next N run times
- JSON formatter:
  ```json
  {
    "expression": "*/5 * * * *",
    "meaning": "Every 5 minutes",
    "next_runs": ["2026-05-08T14:00:00Z", "..."]
  }
  ```
- Wire everything together in `main.go`; full end-to-end flow works

## Phase 8 — Extended formats `[ ]`

_Not in scope until Phase 7 is complete and stable._

- **6-field (seconds)**: detect 6 tokens, treat first field as seconds (0–59)
- **Predefined strings**: `@yearly`, `@monthly`, `@weekly`, `@daily`,
  `@hourly`, `@midnight` — expand to their 5-field equivalents before parsing
- **Timezone prefix**: `TZ=America/New_York 0 9 * * *` — parse the prefix,
  use `time.LoadLocation` to compute runs in the specified zone

## Phase 9 — Polish & release prep `[ ]`

- Friendly, structured error messages (bad field, out-of-range value, etc.)
- `--help` output
- Expanded README with examples for every supported syntax
- Goreleaser config for cross-platform GitHub releases
