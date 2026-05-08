# Mission

## What

`gocron` is a standalone Go CLI that takes a cron expression as input and
produces two things:

1. A plain-English explanation of when the expression fires.
2. The next N scheduled run times (default: 3).

## Why

`crontab.guru` is the de-facto tool for this — but it is web-only, requires a
browser, and cannot be piped or scripted. No comparable standalone CLI written
in Go exists. `gocron` fills that gap: one binary, no network, works anywhere.

## Who

A developer or sysadmin who is writing or debugging a cron job and wants a
quick, offline, scriptable sanity check without leaving the terminal.

## What success looks like

```
$ gocron "*/5 * * * *"

Expression:  */5 * * * *
Meaning:     Every 5 minutes

Next 3 runs:
  1.  Thu May  8 2026  14:00:00
  2.  Thu May  8 2026  14:05:00
  3.  Thu May  8 2026  14:10:00
```

- Colored output in a terminal; degrades to plain text when piped.
- `--json` flag emits machine-readable output.
- Clear, friendly errors for invalid expressions.
- No network access, no config files, no daemons.

## Non-goals (v1)

- Executing or scheduling actual jobs.
- Parsing crontab files (multi-line).
- A TUI or interactive mode.
