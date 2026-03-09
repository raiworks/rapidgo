---
title: "CORS"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# CORS

## Abstract

This document covers Cross-Origin Resource Sharing (CORS)
configuration using `gin-contrib/cors`.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Configuration](#2-configuration)
3. [CORS Middleware](#3-cors-middleware)
4. [Allowed Headers](#4-allowed-headers)
5. [Security Considerations](#5-security-considerations)
6. [References](#6-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **CORS** — Cross-Origin Resource Sharing; a mechanism that allows
  restricted resources on a web page to be requested from a different
  origin.
- **Preflight** — An `OPTIONS` request sent by browsers before the
  actual request to check CORS policy.

## 2. Configuration

`.env`:

```env
CORS_ALLOWED_ORIGINS=https://example.com,https://app.example.com
```

If not set or empty, defaults to `*` (allows all origins).

## 3. CORS Middleware

Library: `github.com/gin-contrib/cors`

```go
package middleware

import (
    "os"
    "strings"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
    allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
    if len(allowedOrigins) == 0 || allowedOrigins[0] == "" {
        allowedOrigins = []string{"*"}
    }

    return cors.New(cors.Config{
        AllowOrigins:     allowedOrigins,
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-CSRF-Token"},
        ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}
```

## 4. Allowed Headers

| Header | Purpose |
|--------|---------|
| `Origin` | Identifies the requesting origin |
| `Content-Type` | Request body format |
| `Authorization` | JWT Bearer token |
| `X-CSRF-Token` | CSRF token for AJAX requests |

### Exposed Headers

| Header | Purpose |
|--------|---------|
| `Content-Length` | Response body size |
| `X-Request-ID` | Request tracing ID |

## 5. Security Considerations

- In production, **MUST** specify exact allowed origins rather than
  `*`.
- `AllowCredentials: true` requires specific origins — browsers
  reject `*` with credentials.
- CORS headers are applied to the `api` middleware group by default;
  web routes typically don't need CORS.
- `MaxAge` (12 hours) controls how long browsers cache preflight
  responses — higher values reduce preflight requests.

## 6. References

- [gin-contrib/cors](https://github.com/gin-contrib/cors)
- [Middleware](../http/middleware.md)
- [Rate Limiting](rate-limiting.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
