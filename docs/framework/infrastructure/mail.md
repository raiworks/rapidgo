---
title: "Mail"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Mail

## Abstract

This document covers the email sending system using
`wneessen/go-mail` — configuration, the Mailer struct, and sending
HTML emails.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Configuration](#2-configuration)
3. [Mailer](#3-mailer)
4. [Usage](#4-usage)
5. [Provider Registration](#5-provider-registration)
6. [Security Considerations](#6-security-considerations)
7. [References](#7-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

## 2. Configuration

`.env`:

```env
MAIL_HOST=smtp.example.com
MAIL_PORT=587
MAIL_USERNAME=noreply@example.com
MAIL_PASSWORD=secret
MAIL_FROM_NAME=MyApp
MAIL_FROM_ADDRESS=noreply@example.com
MAIL_ENCRYPTION=tls
```

## 3. Mailer

Library: `github.com/wneessen/go-mail`

```go
package mail

import (
    "os"
    "strconv"

    gomail "github.com/wneessen/go-mail"
)

type Mailer struct {
    Host     string
    Port     int
    Username string
    Password string
    From     string
    FromName string
}

func NewMailer() *Mailer {
    port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
    return &Mailer{
        Host:     os.Getenv("MAIL_HOST"),
        Port:     port,
        Username: os.Getenv("MAIL_USERNAME"),
        Password: os.Getenv("MAIL_PASSWORD"),
        From:     os.Getenv("MAIL_FROM_ADDRESS"),
        FromName: os.Getenv("MAIL_FROM_NAME"),
    }
}

func (m *Mailer) Send(to, subject, htmlBody string) error {
    msg := gomail.NewMsg()
    msg.FromFormat(m.FromName, m.From)
    msg.To(to)
    msg.Subject(subject)
    msg.SetBodyString(gomail.TypeTextHTML, htmlBody)

    client, err := gomail.NewClient(m.Host,
        gomail.WithPort(m.Port),
        gomail.WithSMTPAuth(gomail.SMTPAuthPlain),
        gomail.WithUsername(m.Username),
        gomail.WithPassword(m.Password),
        gomail.WithTLSPolicy(gomail.TLSMandatory),
    )
    if err != nil {
        return err
    }
    return client.DialAndSend(msg)
}
```

## 4. Usage

### Direct Usage

```go
mailer := mail.NewMailer()
mailer.Send("user@example.com", "Welcome!", "<h1>Welcome to MyApp</h1>")
```

### With Events

Send emails asynchronously via the event system:

```go
events.Listen("user.created", func(payload interface{}) {
    user := payload.(*models.User)
    mailer.Send(user.Email, "Welcome!", "<h1>Welcome!</h1>")
})
```

### From Container

```go
mailer := container.MustMake[*mail.Mailer](app.Container, "mail")
mailer.Send(to, subject, htmlBody)
```

## 5. Provider Registration

```go
type MailProvider struct{}

func (p *MailProvider) Register(c *container.Container) {
    c.Singleton("mail", func(c *container.Container) interface{} {
        return mail.NewMailer()
    })
}

func (p *MailProvider) Boot(c *container.Container) {}
```

## 6. Security Considerations

- SMTP credentials **MUST** be stored in `.env` and **MUST NOT** be
  committed to version control.
- `TLSMandatory` ensures encrypted SMTP connections — never use
  plain SMTP in production.
- Validate recipient email addresses before sending to prevent
  mail injection.
- Consider rate limiting outbound emails to prevent abuse.

## 7. References

- [wneessen/go-mail](https://github.com/wneessen/go-mail)
- [Events](events.md)
- [Service Providers](../core/service-providers.md)
- [Configuration](../core/configuration.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
