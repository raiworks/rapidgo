---
title: "Requests & Validation"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Requests & Validation

## Abstract

This document covers the two validation approaches: the framework's
built-in validator (zero dependencies, fluent API) and struct-based
validation using `go-playground/validator` via Gin's binding system.
It includes API and SSR validation flows.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Built-in Validator](#2-built-in-validator)
3. [Validation Methods](#3-validation-methods)
4. [API Validation Flow](#4-api-validation-flow)
5. [SSR Validation Flow](#5-ssr-validation-flow)
6. [Struct-based Validation](#6-struct-based-validation)
7. [Request Structs](#7-request-structs)
8. [When to Use Which](#8-when-to-use-which)
9. [Security Considerations](#9-security-considerations)
10. [References](#10-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Built-in validator** — The framework's zero-dependency validation
  engine in `core/validation/`.
- **Struct-based validation** — Validation via struct tags using
  `go-playground/validator`, integrated through Gin's `ShouldBind*`
  methods.

## 2. Built-in Validator

Located in `core/validation/`. Provides a fluent, chainable API with
no external dependencies:

```go
package validation

import (
    "fmt"
    "net"
    "net/mail"
    "regexp"
    "strings"
    "unicode/utf8"
)

type Errors map[string][]string

func (e Errors) HasErrors() bool { return len(e) > 0 }

func (e Errors) Add(field, message string) {
    e[field] = append(e[field], message)
}

func (e Errors) First(field string) string {
    if msgs, ok := e[field]; ok && len(msgs) > 0 {
        return msgs[0]
    }
    return ""
}

type Validator struct {
    errors Errors
}

func New() *Validator {
    return &Validator{errors: make(Errors)}
}

func (v *Validator) Errors() Errors { return v.errors }
func (v *Validator) Valid() bool     { return !v.errors.HasErrors() }
```

## 3. Validation Methods

All methods return `*Validator` for chaining:

| Method | Signature | Description |
|--------|-----------|-------------|
| `Required` | `(field, value string)` | Value is not empty |
| `MinLength` | `(field, value string, min int)` | Minimum string length |
| `MaxLength` | `(field, value string, max int)` | Maximum string length |
| `Email` | `(field, value string)` | Valid email format |
| `URL` | `(field, value string)` | Starts with `http://` or `https://` |
| `Matches` | `(field, value, pattern string)` | Matches regex pattern |
| `In` | `(field, value string, allowed []string)` | Value in allowed list |
| `Confirmed` | `(field, value, confirmation string)` | Two values match |
| `IP` | `(field, value string)` | Valid IP address |

### Method Implementations

```go
func (v *Validator) Required(field, value string) *Validator {
    if strings.TrimSpace(value) == "" {
        v.errors.Add(field, field+" is required")
    }
    return v
}

func (v *Validator) MinLength(field, value string, min int) *Validator {
    if utf8.RuneCountInString(value) < min {
        v.errors.Add(field, fmt.Sprintf("%s must be at least %d characters", field, min))
    }
    return v
}

func (v *Validator) MaxLength(field, value string, max int) *Validator {
    if utf8.RuneCountInString(value) > max {
        v.errors.Add(field, fmt.Sprintf("%s must be at most %d characters", field, max))
    }
    return v
}

func (v *Validator) Email(field, value string) *Validator {
    if _, err := mail.ParseAddress(value); err != nil {
        v.errors.Add(field, field+" must be a valid email")
    }
    return v
}

func (v *Validator) URL(field, value string) *Validator {
    if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
        v.errors.Add(field, field+" must be a valid URL")
    }
    return v
}

func (v *Validator) Matches(field, value, pattern string) *Validator {
    if matched, _ := regexp.MatchString(pattern, value); !matched {
        v.errors.Add(field, field+" format is invalid")
    }
    return v
}

func (v *Validator) In(field, value string, allowed []string) *Validator {
    for _, a := range allowed {
        if value == a {
            return v
        }
    }
    v.errors.Add(field, fmt.Sprintf("%s must be one of: %s",
        field, strings.Join(allowed, ", ")))
    return v
}

func (v *Validator) Confirmed(field, value, confirmation string) *Validator {
    if value != confirmation {
        v.errors.Add(field, field+" confirmation does not match")
    }
    return v
}

func (v *Validator) IP(field, value string) *Validator {
    if net.ParseIP(value) == nil {
        v.errors.Add(field, field+" must be a valid IP address")
    }
    return v
}
```

## 4. API Validation Flow

Use the built-in validator with form or JSON data, returning errors
as a JSON envelope:

```go
func CreateUser(c *gin.Context) {
    name := c.PostForm("name")
    email := c.PostForm("email")
    password := c.PostForm("password")
    passwordConfirm := c.PostForm("password_confirmation")

    v := validation.New()
    v.Required("name", name).MinLength("name", name, 2).MaxLength("name", name, 100)
    v.Required("email", email).Email("email", email)
    v.Required("password", password).MinLength("password", password, 8)
    v.Confirmed("password", password, passwordConfirm)

    if !v.Valid() {
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "error":   "validation failed",
            "details": v.Errors(),
        })
        return
    }
    // proceed with validated data
}
```

## 5. SSR Validation Flow

For server-rendered forms, pass errors and old input back to the
template. Optionally use flash messages for post-redirect-get:

```go
if !v.Valid() {
    c.HTML(http.StatusUnprocessableEntity, "users/create.html", gin.H{
        "errors": v.Errors(),
        "old":    gin.H{"name": name, "email": email},
    })
    return
}
```

With flash messages (post-redirect-get pattern):

```go
if !v.Valid() {
    sessionMgr.FlashErrors(data, v.Errors())
    sessionMgr.FlashOldInput(data, map[string]string{
        "name":  c.PostForm("name"),
        "email": c.PostForm("email"),
    })
    c.Set("session", data)
    c.Redirect(http.StatusFound, "/users/create")
    return
}
```

## 6. Struct-based Validation

For API endpoints, use struct tags with Gin's built-in binding
(powered by `go-playground/validator`):

```go
func CreateUserAPI(c *gin.Context) {
    var req requests.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        responses.Error(c, 422, err.Error())
        return
    }
    // req is validated — pass to service
}
```

## 7. Request Structs

Define request structs in `http/requests/`:

```go
package requests

type CreateUserRequest struct {
    Name     string `json:"name" binding:"required,min=2,max=100"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Age      int    `json:"age" binding:"omitempty,gte=1,lte=150"`
}

type UpdateUserRequest struct {
    Name  string `json:"name" binding:"omitempty,min=2,max=100"`
    Email string `json:"email" binding:"omitempty,email"`
}

type PaginationRequest struct {
    Page    int `form:"page" binding:"omitempty,min=1"`
    PerPage int `form:"per_page" binding:"omitempty,min=1,max=100"`
}
```

Common binding tags:

| Tag | Description |
|-----|-------------|
| `required` | Field must be present and non-zero |
| `min=N` | Minimum length (string) or value (number) |
| `max=N` | Maximum length (string) or value (number) |
| `email` | Valid email format |
| `gte=N` | Greater than or equal |
| `lte=N` | Less than or equal |
| `omitempty` | Skip validation if empty |

## 8. When to Use Which

| Approach | Best For | Pros | Cons |
|----------|----------|------|------|
| **Built-in validator** | SSR forms, complex conditional logic | Zero dependencies, fluent API, full control | More verbose |
| **Struct-based** | API endpoints, JSON body validation | Concise, declarative, Gin-integrated | Less flexibility for custom rules |

**Recommendation:** Use struct-based validation for API endpoints and
the built-in validator for SSR forms where you need old input and
flash message support.

## 9. Security Considerations

- All user input **MUST** be validated before processing.
- Validation error messages **MUST NOT** expose internal
  implementation details (e.g., database column names).
- Regex patterns used in `Matches()` **SHOULD** be anchored (`^...$`)
  to prevent partial matches.
- File upload validations (size, type) **MUST** be enforced on the
  server side, even if client-side checks exist.

## 10. References

- [Controllers](controllers.md)
- [Responses](responses.md)
- [Views — Form Templates](views.md#6-form-templates-with-validation)
- [Sessions — Flash Messages](../security/sessions.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
