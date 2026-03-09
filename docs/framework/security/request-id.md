---
title: "Request ID / Tracing"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Request ID / Tracing

## Abstract

This document covers the Request ID middleware that attaches a unique
identifier to every request for logging, debugging, and distributed
tracing.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Middleware Implementation](#2-middleware-implementation)
3. [Behavior](#3-behavior)
4. [Using Request IDs](#4-using-request-ids)
5. [Security Considerations](#5-security-considerations)
6. [References](#6-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Request ID** — A unique string identifier assigned to each HTTP
  request for tracing across logs and services.

## 2. Middleware Implementation

```go
func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.GetHeader("X-Request-ID")
        if id == "" {
            b := make([]byte, 16)
            rand.Read(b)
            id = hex.EncodeToString(b)
        }
        c.Set("request_id", id)
        c.Header("X-Request-ID", id)
        c.Next()
    }
}
```

## 3. Behavior

| Scenario | Action |
|----------|--------|
| Client sends `X-Request-ID` header | Reused as-is |
| No header present | Generated (32-char hex from `crypto/rand`) |

The request ID is:
1. Stored in the Gin context as `request_id`
2. Returned in the `X-Request-ID` response header

## 4. Using Request IDs

### In Structured Logging

```go
slog.Info("processing request",
    "request_id", c.GetString("request_id"),
    "path", c.Request.URL.Path,
)
```

### In Error Responses

Include the request ID for client-side debugging:

```go
c.JSON(500, gin.H{
    "error":      "internal server error",
    "request_id": c.GetString("request_id"),
})
```

### In API Responses

The request ID is exposed via the `X-Request-ID` response header,
which is included in `ExposeHeaders` of the CORS configuration.

## 5. Security Considerations

- Request IDs from clients **SHOULD** be validated (length, format)
  in production to prevent log injection.
- Don't include sensitive information in request IDs.

## 6. References

- [Middleware](../http/middleware.md)
- [Logging](../core/logging.md)
- [CORS](cors.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
