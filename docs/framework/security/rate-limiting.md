---
title: "Rate Limiting"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Rate Limiting

## Abstract

This document covers request rate limiting using `ulule/limiter/v3`
with memory and Redis backends.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Configuration](#2-configuration)
3. [Rate Limit Middleware](#3-rate-limit-middleware)
4. [Redis Backend](#4-redis-backend)
5. [Rate Format](#5-rate-format)
6. [Security Considerations](#6-security-considerations)
7. [References](#7-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Rate limiting** — Restricting the number of requests a client can
  make within a time window.

## 2. Configuration

`.env`:

```env
RATE_LIMIT=60-M
RATE_LIMIT_AUTH=5-M
```

## 3. Rate Limit Middleware

Library: `github.com/ulule/limiter/v3`

```go
package middleware

import (
    "os"

    "github.com/gin-gonic/gin"
    limiter "github.com/ulule/limiter/v3"
    mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
    "github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimitMiddleware() gin.HandlerFunc {
    rate := os.Getenv("RATE_LIMIT")
    if rate == "" {
        rate = "60-M" // 60 requests per minute
    }
    r, _ := limiter.NewRateFromFormatted(rate)
    store := memory.NewStore()
    instance := limiter.New(store, r)
    return mgin.NewMiddleware(instance)
}
```

## 4. Redis Backend

For multi-instance deployments, use a shared Redis store:

```go
import sredis "github.com/ulule/limiter/v3/drivers/store/redis"

store, _ := sredis.NewStoreWithOptions(
    redisClient,
    limiter.StoreOptions{Prefix: "rate:"},
)
```

Redis ensures consistent rate limiting across all application
instances.

## 5. Rate Format

| Format | Meaning |
|--------|---------|
| `60-M` | 60 requests per minute |
| `100-H` | 100 requests per hour |
| `5-M` | 5 requests per minute (for auth endpoints) |
| `1000-D` | 1000 requests per day |

### Response Headers

When rate limited, the middleware returns `429 Too Many Requests`
with standard headers:

| Header | Description |
|--------|-------------|
| `X-RateLimit-Limit` | Maximum requests allowed |
| `X-RateLimit-Remaining` | Requests remaining in window |
| `X-RateLimit-Reset` | Unix timestamp when window resets |

## 6. Security Considerations

- Rate limiting **SHOULD** be applied to all API routes.
- Authentication endpoints (login, register) **MUST** have stricter
  limits (e.g., `5-M`) to prevent brute-force attacks.
- Use Redis backend in multi-instance deployments to prevent
  per-instance limits from being bypassed.
- Consider rate limiting by user ID (authenticated) or IP address
  (unauthenticated).

## 7. References

- [ulule/limiter](https://github.com/ulule/limiter)
- [Middleware](../http/middleware.md)
- [CORS](cors.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
