---
title: "Error Handling"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Error Handling

## Abstract

This document describes the framework's centralized error handling
strategy — the error middleware that catches panics, returns consistent
JSON or HTML error responses, and integrates with structured logging.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Overview](#2-overview)
3. [Error Middleware](#3-error-middleware)
4. [JSON vs HTML Responses](#4-json-vs-html-responses)
5. [Structured Error Logging](#5-structured-error-logging)
6. [Debug Mode Behavior](#6-debug-mode-behavior)
7. [Custom Error Types](#7-custom-error-types)
8. [Security Considerations](#8-security-considerations)
9. [References](#9-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

## 2. Overview

The framework uses a central error middleware registered on the Gin
engine to catch unhandled errors and panics. This ensures:

- **Consistency** — All error responses follow the same format.
- **Safety** — Stack traces are never exposed in production.
- **Observability** — All errors are logged with structured context.

## 3. Error Middleware

```go
package middleware

import (
    "log/slog"
    "net/http"

    "github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            slog.Error("request error", "path", c.Request.URL.Path, "err", err)

            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "internal server error",
            })
        }
    }
}
```

The middleware:

1. Calls `c.Next()` to execute the handler chain.
2. After the handler returns, checks `c.Errors` for any errors added
   via `c.Error(err)`.
3. Logs the error with path context using `slog`.
4. Returns a generic error response to the client.

## 4. JSON vs HTML Responses

The error response format **SHOULD** match the request type:

### API requests (JSON)

When the request accepts JSON (e.g., API routes with
`Accept: application/json`), errors are returned as:

```json
{
    "success": false,
    "error": "internal server error"
}
```

### Web requests (HTML)

For SSR web pages, errors **SHOULD** render an error template:

```go
if c.GetHeader("Accept") == "application/json" ||
   strings.HasPrefix(c.Request.URL.Path, "/api") {
    c.JSON(status, gin.H{"error": message})
} else {
    c.HTML(status, "errors/500.html", gin.H{"error": message})
}
```

## 5. Structured Error Logging

All errors **MUST** be logged with structured context using `slog`:

```go
slog.Error("request error",
    "path", c.Request.URL.Path,
    "method", c.Request.Method,
    "err", err,
    "request_id", c.GetString("request_id"),
)
```

This produces JSON log output in production:

```json
{
    "level": "ERROR",
    "msg": "request error",
    "path": "/api/users",
    "method": "POST",
    "err": "duplicate key value",
    "request_id": "a1b2c3d4e5f6"
}
```

The `request_id` field enables correlating errors with specific
requests across log aggregation systems.

## 6. Debug Mode Behavior

Error detail level is controlled by `APP_DEBUG`:

### Development (`APP_DEBUG=true`)

- Use `gin.Recovery()` — provides detailed stack traces.
- Error responses **MAY** include the actual error message.

```go
if config.IsDebug() {
    r.Use(gin.Recovery())
} else {
    r.Use(middleware.ErrorHandler())
}
```

### Production (`APP_DEBUG=false`)

- Use `ErrorHandler()` — returns generic messages only.
- Actual error details are logged but **MUST NOT** be sent to clients.
- Stack traces **MUST NOT** appear in responses.

## 7. Custom Error Types

For domain-specific errors, define custom error types that carry
status codes and user-facing messages:

```go
package errors

type AppError struct {
    Code    int    // HTTP status code
    Message string // User-facing message
    Err     error  // Internal error (logged, not exposed)
}

func (e *AppError) Error() string {
    return e.Message
}

func NotFound(msg string) *AppError {
    return &AppError{Code: 404, Message: msg}
}

func BadRequest(msg string) *AppError {
    return &AppError{Code: 400, Message: msg}
}

func Internal(err error) *AppError {
    return &AppError{Code: 500, Message: "internal server error", Err: err}
}
```

Usage in services:

```go
func (s *UserService) GetByID(id uint) (*models.User, error) {
    var user models.User
    if err := s.DB.First(&user, id).Error; err != nil {
        return nil, errors.NotFound("user not found")
    }
    return &user, nil
}
```

## 8. Security Considerations

- **MUST NOT** expose stack traces, internal error messages, or
  database error details in production responses.
- **MUST NOT** leak file paths, SQL queries, or configuration values
  in error messages.
- `APP_DEBUG` **MUST** be `false` in production.
- All errors **MUST** be logged for auditability, even when the
  client receives a generic message.
- Validation errors (422) **MAY** include field-specific details
  since these are user-input related, not internal failures.

## 9. References

- [Logging](logging.md)
- [Configuration](configuration.md)
- [Middleware](../http/middleware.md)
- [Responses](../http/responses.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
