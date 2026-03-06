# Feature #29 — Mail / Email: Discussion

## What problem does this solve?

Applications need to send transactional emails — password resets, welcome messages, notifications, alerts. A `Mailer` abstraction centralizes SMTP configuration and provides a clean API for sending HTML emails.

## Why now?

Core infrastructure (#02, #05), config (#02), and authentication (#21) are shipped. Auth flows often require email (password reset, verification). Mail is the next foundational service.

## What does the blueprint specify?

- `Mailer` struct with Host, Port, Username, Password, From, FromName fields.
- `NewMailer()` factory reading from env vars: `MAIL_HOST`, `MAIL_PORT`, `MAIL_USERNAME`, `MAIL_PASSWORD`, `MAIL_FROM_ADDRESS`, `MAIL_FROM_NAME`.
- `Send(to, subject, htmlBody)` method using `go-mail` library.
- Uses `gomail.SMTPAuthPlain` and `gomail.TLSMandatory`.
- Library: `github.com/wneessen/go-mail`.

## Design decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Library | `github.com/wneessen/go-mail` | Blueprint-specified; actively maintained, modern Go mail library |
| Package location | `core/mail/` | Existing stub directory, follows `core/*` pattern |
| SMTP auth | Plain auth with TLS | Blueprint default; standard for SMTP relay services |
| Testing | Unit test struct creation + config; no live SMTP tests | Can't send real emails in CI; verify config wiring and message construction |

## What is out of scope?

- Email templates / template rendering (app-level concern).
- Queue / async sending (future enhancement).
- Multiple mail drivers (Mailgun, SES, etc.).
- Attachment support (can be added later via `go-mail` API).
