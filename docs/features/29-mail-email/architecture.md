# Feature #29 — Mail / Email: Architecture

## Component overview

```
.env
    MAIL_HOST=smtp.example.com
    MAIL_PORT=587
    MAIL_USERNAME=noreply@example.com
    MAIL_PASSWORD=secret
    MAIL_FROM_NAME=MyApp
    MAIL_FROM_ADDRESS=noreply@example.com
        │
        ▼
core/mail/mail.go       Mailer struct + NewMailer() + Send()
```

## New file

| File | Purpose |
|------|---------|
| `core/mail/mail.go` | `Mailer` struct, `NewMailer()` factory, `Send()` method |

## Removed

| File | Reason |
|------|--------|
| `core/mail/.gitkeep` | Replaced by real implementation |

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/wneessen/go-mail` | SMTP client with TLS support |

## Mailer struct

```go
type Mailer struct {
    Host     string
    Port     int
    Username string
    Password string
    From     string
    FromName string
}
```

## Methods

| Method | Signature | Behaviour |
|--------|-----------|-----------|
| `NewMailer()` | `func NewMailer() *Mailer` | Reads all `MAIL_*` env vars, returns configured Mailer |
| `Send()` | `func (m *Mailer) Send(to, subject, htmlBody string) error` | Creates message, dials SMTP, sends with TLS |

## Environment variables

| Var | Required | Purpose |
|-----|----------|---------|
| `MAIL_HOST` | Yes | SMTP server hostname |
| `MAIL_PORT` | Yes | SMTP port (typically 587) |
| `MAIL_USERNAME` | Yes | SMTP auth username |
| `MAIL_PASSWORD` | Yes | SMTP auth password |
| `MAIL_FROM_ADDRESS` | Yes | Sender email address |
| `MAIL_FROM_NAME` | Yes | Sender display name |

## Send flow

1. Create `gomail.NewMsg()`.
2. Set From (with display name), To, Subject, HTML body.
3. Create client with host, port, SMTP plain auth, TLS mandatory.
4. `client.DialAndSend(msg)`.
