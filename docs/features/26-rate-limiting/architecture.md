# Feature #26 — Rate Limiting: Architecture

## Component overview

```
.env                       RATE_LIMIT=60-M
    │
    ▼
core/middleware/ratelimit.go    RateLimitMiddleware() gin.HandlerFunc
    │
    ├── reads RATE_LIMIT env var (default "60-M")
    ├── parses rate via limiter.NewRateFromFormatted()
    ├── creates memory.NewStore()
    ├── creates limiter.New(store, rate)
    └── returns mgin.NewMiddleware(instance)
    │
    ▼
app/providers/middleware_provider.go
    └── RegisterAlias("ratelimit", middleware.RateLimitMiddleware())
```

## New file

| File | Purpose |
|------|---------|
| `core/middleware/ratelimit.go` | `RateLimitMiddleware()` — configurable IP-based rate limiter |

## Modified file

| File | Change |
|------|--------|
| `app/providers/middleware_provider.go` | Add `"ratelimit"` alias |

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/ulule/limiter/v3` | latest | Rate limiting core |
| `github.com/ulule/limiter/v3/drivers/middleware/gin` | latest | Gin middleware adapter |
| `github.com/ulule/limiter/v3/drivers/store/memory` | latest | In-memory token bucket store |

## Behaviour

1. On boot, `RateLimitMiddleware()` reads `RATE_LIMIT` from env (default `"60-M"`).
2. Parses rate format: `"<count>-<period>"` where period is S (second), M (minute), H (hour), D (day).
3. Creates an in-memory store and a limiter instance.
4. Returns the `mgin` Gin middleware handler.
5. Each request is keyed by client IP. If the rate limit is exceeded, the middleware returns `429 Too Many Requests` with standard `X-RateLimit-*` headers.

## Response headers (set by ulule/limiter)

| Header | Description |
|--------|-------------|
| `X-RateLimit-Limit` | Maximum requests in the window |
| `X-RateLimit-Remaining` | Remaining requests in the current window |
| `X-RateLimit-Reset` | Unix timestamp when the window resets |
