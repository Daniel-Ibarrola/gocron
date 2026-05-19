# Phase 3 — Ranges: Requirements

## Goal

Extend the parser and explainer to handle `a-b` range tokens in any of the
five cron fields, and refactor the internal field representation so that
Phase 4 (steps) and Phase 5 (lists) can extend it without another rewrite.

## In scope

### Parser

- Accept `a-b` in any field, where `a` and `b` are integers and `a <= b`.
- Validate that both endpoints fall within the field's bounds (see below).
- Reject inverted ranges (`b < a`) with a `FieldError`.
- Reject out-of-bounds endpoints with a `FieldError`.
- Allow degenerate ranges (`5-5`); normalize them to a single value at
  parse time so the explainer never has to say "5 through 5".
- Also validate that single integer values fall within field bounds —
  the Phase 2 parser accepted out-of-range singles (`0 99 * * *`). Fix
  that here, since the bounds table is being introduced anyway.

### Field bounds

| Field         | Min | Max |
|---------------|----:|----:|
| minute        |   0 |  59 |
| hour          |   0 |  23 |
| day of month  |   1 |  31 |
| month         |   1 |  12 |
| day of week   |   0 |   6 |

Day-of-week: Sunday = 0. `7` is **not** accepted as an alias for Sunday in
this phase (classic crontab quirk; explicitly rejected for now).

### Internal data model

Replace the current `cronExpression` struct (five `int` fields, `-1` for
wildcard) with a per-field tagged shape:

```go
type fieldKind int

const (
    fieldWildcard fieldKind = iota
    fieldSingle
    fieldRange
)

type fieldSpec struct {
    kind   fieldKind
    lo, hi int // single: lo==hi; range: lo<=hi; wildcard: zero
}

type cronExpression struct {
    minutes, hours, daysOfMonth, months, daysOfWeek fieldSpec
}
```

A `fieldBounds` table indexed by field position provides the min/max for
validation.

### Explainer

Per-field phrasing, dispatched on `kind`:

| Field        | wildcard      | single                 | range                                     |
|--------------|---------------|------------------------|-------------------------------------------|
| minute       | "every minute" / omitted | "at minute M" / combined into "HH:MM" | "at minutes A through B" |
| hour         | omitted       | "past hour H" / combined into "HH:MM" | "past hours A through B" |
| day of month | omitted       | "on day of month D"    | "on days of month A through B"            |
| month        | omitted       | "in MonthName"         | "in MonthA through MonthB"                |
| day of week  | omitted       | "on DayName"           | "on DayA through DayB"                    |

The "at HH:MM" combination only applies when **both** minute and hour are
single values. Other combinations:

- `5-10 9-17 * * *` → "at minutes 5 through 10 past hours 9 through 17"
- `5 9-17 * * *`   → "at minute 5 past hours 9 through 17"
- `5-10 4 * * *`   → "at minutes 5 through 10 past hour 4"
- `* * * 6-8 *`    → "every minute in June through August"
- `* * * * 1-5`    → "every minute on Monday through Friday"
- `* * 1-15 * *`   → "every minute on days of month 1 through 15"

Whole-field ranges (`0-59 * * * *`) are **not** normalized to wildcard —
the user wrote a range, the explainer renders a range.

### Tests

Extend `parse_test.go` and `explain_test.go`:

- Valid range in each of the five fields.
- Range combined with single values across fields.
- Inverted range (`5-1`) → `FieldError`.
- Out-of-bounds range (`0-99` for minute, `0-32` for dom, `0-7` for dow, etc.).
- Out-of-bounds single (`99 * * * *`) → `FieldError` (new).
- Degenerate range (`5-5`) → parses, behaves like single 5.
- `7` as day-of-week → `FieldError`.
- All existing Phase 2 tests still pass against the new data model.

## Out of scope (deferred)

- Step syntax `*/N` and `a-b/N` — Phase 4.
- List syntax `1,2,3` — Phase 5.
- Whole-range collapsing to wildcard in the explainer.
- Day-of-week name input (e.g. `MON-FRI`) — not on the current roadmap.
- Day-of-week `7` as Sunday alias — explicitly rejected for this phase.

## Decisions

1. **Data model: tagged struct (option A)** over enumerated slices or a
   parallel range struct. It extends cleanly to steps (add a `step` field)
   and lists (wrap into `[]fieldSpec`). Slices lose the "this was a range"
   information the explainer needs for "9 through 17" phrasing.
2. **Inverted ranges rejected**, no wrap-around. Wrap-around is
   surprising and inconsistent with what most cron implementations do.
3. **Degenerate ranges (`5-5`) normalized to single** at parse time. The
   parser is the right place to drop redundant structure; the explainer
   shouldn't have to special-case it.
4. **`7` rejected for day-of-week.** Classic crontab accepts it as
   Sunday, but supporting it requires special-case validation logic for
   one field. Defer until/unless a real user asks for it.
5. **Bounds-check singles in this phase.** Phase 2 accepted
   `0 99 * * *` because no bounds existed yet. Now that we're adding the
   bounds table, applying it everywhere is one extra line — leaving the
   gap would be deliberate inconsistency.
6. **Branch name**: `phase-3-ranges`.
7. **No new package.** Everything stays in `internal/cron`. The roadmap
   layout reserves `internal/schedule` and `internal/output` for later
   phases; nothing in Phase 3 belongs there.

## Context

- `specs/mission.md` — `gocron` is a standalone, offline, scriptable cron
  explainer. Range support is one of the table-stakes pieces of cron syntax.
- `specs/tech-stack.md` — stdlib only; no third-party cron library.
  `strings` + `strconv` are the only tools needed for `a-b` parsing.
- `specs/roadmap.md` Phase 3 — defines the surface area: `a-b` parsing,
  bounds validation, explainer phrasing, and tests for valid/inverted/
  out-of-bounds ranges.
- Current repo state: Phase 2 merged via PR #1. Parser handles wildcards
  and singles. Explainer covers all combinations of those. The data model
  is `cronExpression{minutes, hours, ...: int}` with `-1` sentinel for
  wildcard — this is the shape being replaced.
