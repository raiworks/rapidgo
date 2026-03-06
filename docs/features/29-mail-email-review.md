# Feature #29 — Mail / Email: Review

## Summary

Implemented `Mailer` struct in `core/mail/mail.go` with `NewMailer()` factory and `Send()` method using `github.com/wneessen/go-mail`. Reads SMTP config from `MAIL_*` env vars. Removed `.gitkeep`.

## Files changed

| File | Change |
|------|--------|
| `core/mail/mail.go` | New — `Mailer` struct, `NewMailer()`, `Send(to, subject, htmlBody)` |
| `core/mail/mail_test.go` | New — 5 tests (TC-01 to TC-05) |
| `core/mail/.gitkeep` | Removed |
| `go.mod` / `go.sum` | Added `github.com/wneessen/go-mail` v0.7.2 |

## Test results

| TC | Description | Result |
|----|-------------|--------|
| TC-01 | NewMailer reads all env vars | ✅ PASS |
| TC-02 | Default port 587 when unset | ✅ PASS |
| TC-03 | Parses port from string | ✅ PASS |
| TC-04 | Invalid port falls back to 587 | ✅ PASS |
| TC-05 | Send returns error with invalid host | ✅ PASS |

## Regression

- All 25 packages pass.
- `go vet` clean.

## Deviation log

_None._

## Commit

`e451f21` — merged to `main`.
