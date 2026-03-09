---
title: "SSR Forms with Validation & Flash Messages"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# SSR Forms with Validation & Flash Messages

## Abstract

A guide to building server-side rendered forms with CSRF protection,
validation error display, old input preservation, and flash messages.

## Table of Contents

1. [Overview](#1-overview)
2. [Middleware Setup](#2-middleware-setup)
3. [Create Form Template](#3-create-form-template)
4. [Handle POST with Validation](#4-handle-post-with-validation)
5. [Flash Errors and Old Input](#5-flash-errors-and-old-input)
6. [Flash Success Messages](#6-flash-success-messages)
7. [Complete Flow](#7-complete-flow)
8. [References](#8-references)

## 1. Overview

SSR forms follow a redirect-after-POST pattern:

1. **GET** — render form (with old input and errors if redirected)
2. **POST** — validate → on failure, flash errors & redirect back;
   on success, flash success & redirect forward

## 2. Middleware Setup

Apply the `web` middleware group which includes session and CSRF:

```go
web := r.Group("/")
web.Use(middleware.ResolveGroup("web")...)
// web group includes: session, csrf
```

## 3. Create Form Template

`resources/views/contacts/create.html`:

```html
{{template "layouts/app" .}}
{{define "content"}}
<h1>Contact Us</h1>

{{if .success}}
<div class="alert alert-success">{{.success}}</div>
{{end}}

{{if .errors}}
<ul class="alert alert-danger">
    {{range $field, $msg := .errors}}
    <li>{{$field}}: {{$msg}}</li>
    {{end}}
</ul>
{{end}}

<form method="POST" action="/contact">
    <input type="hidden" name="_csrf" value="{{.csrf_token}}">

    <label>Name</label>
    <input type="text" name="name" value="{{.old.name}}">

    <label>Email</label>
    <input type="email" name="email" value="{{.old.email}}">

    <label>Message</label>
    <textarea name="message">{{.old.message}}</textarea>

    <button type="submit">Send</button>
</form>
{{end}}
```

Key elements:
- Hidden `_csrf` field from session
- `{{.old.name}}` restores previous input on validation failure
- `{{.errors}}` displays validation errors
- `{{.success}}` displays success flash

## 4. Handle POST with Validation

```go
func ContactStore(c *gin.Context) {
    v := validation.New()
    v.Required("name", c.PostForm("name"))
    v.MinLength("name", c.PostForm("name"), 2)
    v.Required("email", c.PostForm("email"))
    v.Email("email", c.PostForm("email"))
    v.Required("message", c.PostForm("message"))
    v.MinLength("message", c.PostForm("message"), 10)

    if v.HasErrors() {
        session.FlashErrors(c, v.Errors)
        session.FlashOldInput(c, map[string]string{
            "name":    c.PostForm("name"),
            "email":   c.PostForm("email"),
            "message": c.PostForm("message"),
        })
        c.Redirect(http.StatusFound, "/contact")
        return
    }

    // Process the contact form...
    session.Flash(c, "success", "Message sent successfully!")
    c.Redirect(http.StatusFound, "/contact")
}
```

## 5. Flash Errors and Old Input

### `FlashErrors`

Stores validation errors in the session for the next request:

```go
session.FlashErrors(c, v.Errors)
// v.Errors is map[string]string: {"name": "name is required"}
```

### `FlashOldInput`

Stores submitted form values so they survive the redirect:

```go
session.FlashOldInput(c, map[string]string{
    "name": c.PostForm("name"),
})
```

### GET Handler Retrieves Flash Data

```go
func ContactCreate(c *gin.Context) {
    c.HTML(http.StatusOK, "contacts/create.html", gin.H{
        "csrf_token": session.Get(c, "_csrf"),
        "errors":     session.GetFlash(c, "errors"),
        "old":        session.GetFlash(c, "old"),
        "success":    session.GetFlash(c, "success"),
    })
}
```

Flash data is automatically cleared after being read.

## 6. Flash Success Messages

```go
session.Flash(c, "success", "Your message has been sent!")
c.Redirect(http.StatusFound, "/contact")
```

## 7. Complete Flow

```text
User submits form
       │
       ▼
POST /contact
       │
   Validate ──── fail ───▶ FlashErrors + FlashOldInput
       │                        │
       │                   302 Redirect
       │                        │
   success                      ▼
       │                   GET /contact
       │                   (renders form with errors + old input)
       ▼
Flash("success", "...")
       │
  302 Redirect
       │
       ▼
  GET /contact
  (renders form with success message)
```

## 8. References

- [Sessions — Flash Messages](../security/sessions.md)
- [CSRF Protection](../security/csrf.md)
- [Requests & Validation](../http/requests-validation.md)
- [Views](../http/views.md)
- [Middleware](../http/middleware.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
