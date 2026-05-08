# Roadmap

Phases are intentionally small. Each phase should leave the project in a
working, committable state.

Status legend: `[ ]` not started · `[~]` in progress · `[x]` done

---

## Phase 1 — Project bootstrap `[ ]`

- Initialize `go.mod` (`github.com/dibarrola/gocron`, Go 1.22+)
- Add `github.com/fatih/color` dependency
- `main.go`: read `os.Args`, parse `--next` and `--json` flags, print raw
  expression back to confirm wiring works
- Directory layout:
  ```
  main.go
  internal/
    cron/      # parser lives here
    explain/   # English explainer lives here
    schedule/  # next-run calculator lives here
    output/    # formatting (text + JSON) lives here
  ```

## Phase 2 — 5-field cron parser `[ ]`

Fields: `minute hour dom month dow`

- Define a `Field` type with its name, allowed range, and parsed value set
- Parse each field token into one of:
  - Wildcard (`*`)
  - Single value (`5`)
  - Range (`1-5`)
  - Step (`*/5`, `1-5/2`)
  - List (`1,2,3`) — any combination of the above separated by commas
- Expand every field into a `[]int` of matching values
- Validate bounds per field (minute 0–59, hour 0–23, dom 1–31, month 1–12,
  dow 0–6)
- Return a `CronExpr` struct; return a descriptive error on bad input
- Unit tests: valid expressions, invalid bounds, malformed tokens

## Phase 2.5 — CI pipeline `[ ]`

_Set up before continuing feature work so every subsequent PR is gated._

- **GitHub Actions workflow** at `.github/workflows/ci.yml`, triggered on
  `push` and `pull_request` to any branch
- Jobs (all run on `ubuntu-latest`, Go version pinned to match `go.mod`):
  1. **`fmt`** — `gofmt -l .` fails if any file is not formatted
  2. **`vet`** — `go vet ./...`
  3. **`lint`** — `golangci-lint run` with a minimal `.golangci.yml`
     (enable `errcheck`, `staticcheck`, `unused` at minimum)
  4. **`test`** — `go test -race ./...`
- Jobs are independent and run in parallel; `test` is the only one that
  produces an artifact (coverage summary printed to stdout, no upload yet)
- Add `golangci-lint` version pin to the workflow (not installed globally)
- No caching in v1; add later if build times become noticeable

## Phase 3 — Plain-English explainer `[ ]`

- Given a `CronExpr`, produce a single human-readable sentence
- Handle the most common patterns explicitly:
  - `* * * * *` → "Every minute"
  - `*/N * * * *` → "Every N minutes"
  - `0 * * * *` → "Every hour, on the hour"
  - `0 H * * *` → "Every day at HH:00"
  - `0 H D M *` → "At HH:00, on day D of month M"
  - etc.
- Fall back to a structured description for patterns that don't match a named
  shortcut: "At minutes [0,15,30,45], every hour, every day"
- Unit tests covering each named pattern and the fallback

## Phase 4 — Next-N run calculator `[ ]`

- Given a `CronExpr` and a `time.Time` start, return the next N matching times
- Strategy: advance minute-by-minute from `start+1m`, check each candidate
  against all five fields
- Handle month/dom edge cases (e.g. Feb 30 never matches)
- Cap iteration to avoid infinite loops on expressions that can never fire
  (e.g. `* * 31 2 *`)
- Unit tests: known expressions with known next-run times

## Phase 5 — Output formatting `[ ]`

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

## Phase 6 — Extended formats `[ ]`

_Not in scope until Phase 5 is complete and stable._

- **6-field (seconds)**: detect 6 tokens, treat first field as seconds (0–59)
- **Predefined strings**: `@yearly`, `@monthly`, `@weekly`, `@daily`,
  `@hourly`, `@midnight` — expand to their 5-field equivalents before parsing
- **Timezone prefix**: `TZ=America/New_York 0 9 * * *` — parse the prefix,
  use `time.LoadLocation` to compute runs in the specified zone

## Phase 7 — Polish & release prep `[ ]`

_After Phase 6, or in parallel once core is stable._

- Friendly, structured error messages (bad field, out-of-range value, etc.)
- `--help` output
- README with install instructions and usage examples
- Goreleaser config for cross-platform GitHub releases
