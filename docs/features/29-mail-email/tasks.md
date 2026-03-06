# Feature #29 — Mail / Email: Tasks

## Prerequisites

- [x] Core infrastructure shipped (#02, #05)
- [x] `core/mail/` directory exists (stub with `.gitkeep`)

## Implementation tasks

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 1 | `go get github.com/wneessen/go-mail` | `go.mod`, `go.sum` | ⬜ |
| 2 | Create `Mailer` struct, `NewMailer()`, `Send()` | `core/mail/mail.go` | ⬜ |
| 3 | Remove `.gitkeep` | `core/mail/.gitkeep` | ⬜ |
| 4 | Write tests | `core/mail/mail_test.go` | ⬜ |
| 5 | Full regression + `go vet` | — | ⬜ |
| 6 | Commit, merge, review doc, roadmap update | — | ⬜ |

## Acceptance criteria

- `Mailer` struct has Host, Port, Username, Password, From, FromName fields.
- `NewMailer()` reads all `MAIL_*` env vars.
- `Send(to, subject, htmlBody)` constructs a message and sends via SMTP.
- Default port is 587 when `MAIL_PORT` not set or invalid.
- Tests verify struct creation and env var reading.
- All existing tests pass (regression).
