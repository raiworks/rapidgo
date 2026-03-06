# Feature #26 — Rate Limiting: Discussion

## What problem does this solve?

Without rate limiting, the application is vulnerable to abuse — brute-force attacks, credential stuffing, API scraping, and denial-of-service. Rate limiting restricts the number of requests a client can make within a time window, protecting server resources and downstream services.

## Why now?

Middleware infrastructure (#08), authentication (#21), and CSRF protection (#24) are shipped. Rate limiting is the natural next security layer — it prevents abuse of the auth and CSRF endpoints and provides a global request throttle.

## What does the blueprint specify?

- `RateLimitMiddleware()` using `ulule/limiter/v3` with the Gin driver.
- In-memory store via `memory.NewStore()` (single-instance default).
- Rate format string from `RATE_LIMIT` env var (default `"60-M"` = 60 requests per minute).
- Blueprint also mentions `RATE_LIMIT_AUTH=5-M` for stricter auth-route limiting.
- Redis store option noted for multi-instance deployments (not implemented now).

## Design decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Library | `ulule/limiter/v3` | Blueprint-specified; mature, well-tested, Gin driver included |
| Store | In-memory | Sufficient for single-instance; Redis can be added later |
| Rate format | `"60-M"` default | Blueprint default; configurable via env |
| Middleware alias | `"ratelimit"` | Consistent with existing alias pattern |
| Scope | Global (IP-based) | `ulule/limiter` defaults to IP-based keying |

## What is out of scope?

- Redis-backed store (future enhancement).
- Per-route rate configuration (can be done manually by applying middleware to specific route groups).
- Custom key extraction (e.g., by user ID).
