# gocron

A standalone Go CLI that explains cron expressions in plain English and
(eventually) prints the next N scheduled run times. Think
[crontab.guru](https://crontab.guru), but offline and pipeable.

> Status: early development. Today the CLI explains expressions that use
> wildcards (`*`) and fixed integer values. Ranges, steps, lists, and the
> next-run calculator are tracked in [`specs/roadmap.md`](specs/roadmap.md).

## Install

Requires Go 1.26 or newer.

```sh
go install github.com/dibarrola/gocron@latest
```

Or build from source:

```sh
git clone https://github.com/dibarrola/gocron.git
cd gocron
go build -o gocron .
```

## Usage

```sh
gocron "<cron expression>"
```

Example:

```sh
$ gocron "30 6 * * *"
expression: 30 6 * * *
next: 3
json: false
at 06:30
```

Supported today (5-field expressions, `minute hour dom month dow`):

| Expression    | Explanation                          |
| ------------- | ------------------------------------ |
| `* * * * *`   | every minute                         |
| `5 * * * *`   | at minute 5                          |
| `30 6 * * *`  | at 06:30                             |
| `0 9 * * 1`   | at 09:00 on Monday                   |
| `0 0 1 1 *`   | at 00:00 on day of month 1 in January |

Flags `--next` and `--json` are wired up but not yet functional — they will
become useful once the next-run calculator and output formatter land.

## Development

```sh
go test ./...
go vet ./...
gofmt -l .
```

CI runs the same three checks on every push and pull request.

## Project layout

```
main.go              # CLI entry point
internal/cron/       # parser + explainer
specs/               # mission, tech stack, and roadmap
```

See [`specs/roadmap.md`](specs/roadmap.md) for what's planned next.
