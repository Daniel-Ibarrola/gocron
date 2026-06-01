# Phase 3 — Ranges: Validation

The phase is mergeable when **every** item below passes.

## Build & static checks

- [ ] `gofmt -l ./...` produces no output.
- [ ] `go vet ./...` exits 0.
- [ ] `go test -race ./...` exits 0; every test in `internal/cron`
      passes, including the new range cases and all retained Phase 2 cases.
- [ ] CI workflow (`.github/workflows/ci.yml`) is green on the branch.

## Parser — valid ranges

All of the following parse without error:

```
1-5 * * * *
0 9-17 * * *
* * 1-15 * *
* * * 6-8 *
* * * * 1-5
5-10 9-17 1-15 6-8 1-5
5-5 * * * *           # degenerate, normalized to single 5
```

## Parser — rejected inputs

Each of these returns a `*FieldError` with the offending token in the
`Field` member:

- [ ] Inverted range: `5-1 * * * *`
- [ ] Out-of-bounds range (minute): `0-99 * * * *`
- [ ] Out-of-bounds range (hour): `* 0-24 * * *`
- [ ] Out-of-bounds range (dom): `* * 0-32 * *`
- [ ] Out-of-bounds range (month): `* * * 0-13 *`
- [ ] Out-of-bounds range (dow): `* * * * 0-7`
- [ ] Out-of-bounds single (minute): `99 * * * *`
- [ ] Day-of-week `7`: `* * * * 7`

## Explainer — range phrasing

Run the binary (`go run . "<expr>"`) and confirm the meaning line:

| Expression          | Expected meaning                                                  |
|---------------------|-------------------------------------------------------------------|
| `0 9-17 * * *`      | at minute 0 past hours 9 through 17                               |
| `5-10 * * * *`      | at minutes 5 through 10                                           |
| `5-10 9-17 * * *`   | at minutes 5 through 10 past hours 9 through 17                   |
| `5-10 4 * * *`      | at minutes 5 through 10 past hour 4                               |
| `* * 1-15 * *`      | every minute on days of month 1 through 15                        |
| `* * * 6-8 *`       | every minute in June through August                               |
| `* * * * 1-5`       | every minute on Monday through Friday                             |
| `5-5 * * * *`       | at minute 5  (degenerate range collapses to single)               |

## Explainer — retained Phase 2 behavior

These must still produce their original Phase 2 output:

| Expression       | Expected meaning                              |
|------------------|-----------------------------------------------|
| `* * * * *`      | every minute                                  |
| `5 * * * *`      | at minute 5                                   |
| `30 6 * * *`    | at 06:30                                      |
| `* 4 * * *`      | every minute past hour 4                      |
| `20 4 4 10 5`    | at 04:20 on day of month 4 on Friday in October |

## Ready to merge

When all boxes are checked, push the branch and open a PR against `main`.
