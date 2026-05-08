# Tech Stack

## Language & runtime

- **Go 1.22+**
- Module: `github.com/dibarrola/gocron`
- Single binary, no CGo

## Core logic — stdlib only (intentional)

The cron parser, explainer, and scheduler are implemented from scratch using
only the Go standard library. This is a deliberate practice constraint; no
third-party cron parsing library is used.

Key stdlib packages:
- `strings`, `strconv` — tokenising and parsing field values
- `time` — computing next run times and formatting output
- `encoding/json` — JSON output mode
- `os`, `fmt` — CLI entry point and output

## CLI

No framework. Argument handling is done with `os.Args` and the stdlib `flag`
package. The surface area is small enough that a framework would add more
complexity than it removes.

```
gocron <expression> [flags]

Flags:
  --next  int    Number of upcoming run times to show (default 3)
  --json         Emit JSON instead of colored text
```

## Output & color

- **`github.com/fatih/color`** for terminal colors (ANSI).
  Chosen because it auto-detects whether stdout is a TTY and disables color
  when piped — the right default behavior with minimal code.
- Plain-text and JSON modes never emit ANSI codes.

## Testing

- `testing` (stdlib) for all unit tests.
- Table-driven tests for the parser and scheduler.
- No mocking frameworks; no external test libraries.

## Distribution

- `go build` produces a single static binary.
- GitHub releases via `goreleaser` (added later, not in v1 scope).

## What is explicitly NOT used

- Any third-party cron library (e.g. `robfig/cron`)
- A CLI framework (e.g. `cobra`, `urfave/cli`)
- A logging library
