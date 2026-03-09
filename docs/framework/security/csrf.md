---
title: "CSRF Protection"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# CSRF Protection

## Abstract

This document covers Cross-Site Request Forgery (CSRF) protection —
the middleware implementation, token generation, validation, and
usage in HTML forms.

## Table of Contents

1. [Terminology](#1-terminology)
2. [How It Works](#2-how-it-works)
3. [CSRF Middleware](#3-csrf-middleware)
4. [HTML Form Usage](#4-html-form-usage)
5. [AJAX Usage](#5-ajax-usage)
6. [Excluded Methods](#6-excluded-methods)
7. [Security Considerations](#7-security-considerations)
8. [References](#8-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **CSRF** — Cross-Site Request Forgery; an attack that forces
  authenticated users to submit unintended requests.
- **CSRF token** — A random, per-session value that proves the
  request originated from the application's own forms.

## 2. How It Works

1. On the first request, the middleware generates a random 64-character
   hex token and stores it in the session.
2. The token is made available to templates via `csrf_token`.
3. On state-changing requests (POST, PUT, PATCH, DELETE), the
   middleware validates the submitted token against the session token.
4. If the tokens don't match, the request is rejected with 403.

```text
GET /form → Generate token → Store in session → Render in form
POST /form → Read token from form → Compare with session → Accept/Reject
```

## 3. CSRF Middleware

```go
func CSRFMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        sess, _ := c.Get("session")
        data := sess.(map[string]interface{})

        // Generate token if not present
        token, ok := data["_csrf_token"].(string)
        if !ok || token == "" {
            b := make([]byte, 32)
            rand.Read(b)
            token = hex.EncodeToString(b)
            data["_csrf_token"] = token
            c.Set("session", data)
        }

        // Make token available to templates
        c.Set("csrf_token", token)

        // Skip safe methods
        if c.Request.Method == "GET" || c.Request.Method == "HEAD" ||
            c.Request.Method == "OPTIONS" {
            c.Next()
            return
        }

        // Validate token on POST/PUT/PATCH/DELETE
        submitted := c.PostForm("_csrf_token")
        if submitted == "" {
            submitted = c.GetHeader("X-CSRF-Token")
        }
        if submitted != token {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "error": "CSRF token mismatch",
            })
            return
        }

        c.Next()
    }
}
```

## 4. HTML Form Usage

Every form **MUST** include the CSRF token as a hidden field:

```html
<form method="POST" action="/users">
    <input type="hidden" name="_csrf_token" value="{{.csrf_token}}">
    <!-- form fields -->
    <button type="submit">Create</button>
</form>
```

## 5. AJAX Usage

For JavaScript-initiated requests, read the token from a meta tag and
send it as a header:

```html
<meta name="csrf-token" content="{{.csrf_token}}">
```

```javascript
const token = document.querySelector('meta[name="csrf-token"]').content;

fetch('/api/users', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
        'X-CSRF-Token': token
    },
    body: JSON.stringify(data)
});
```

## 6. Excluded Methods

The following HTTP methods are **safe** and skip CSRF validation:

| Method | Validated |
|--------|-----------|
| GET | No |
| HEAD | No |
| OPTIONS | No |
| POST | **Yes** |
| PUT | **Yes** |
| PATCH | **Yes** |
| DELETE | **Yes** |

## 7. Security Considerations

- CSRF tokens **MUST** be generated using `crypto/rand`, not
  `math/rand`.
- Tokens **MUST** be stored in the session, not in cookies
  (double-submit cookie is an alternative but not used here).
- CSRF middleware **MUST** be applied to all web route groups.
- API routes using JWT authentication **MAY** skip CSRF protection
  since JWT tokens already prove request origin.
- Token comparison **SHOULD** use constant-time comparison to
  prevent timing attacks.

## 8. References

- [Sessions](sessions.md)
- [Middleware](../http/middleware.md)
- [Views — Form Templates](../http/views.md#6-form-templates-with-validation)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
