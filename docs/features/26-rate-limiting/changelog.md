# Feature #26 — Rate Limiting: Changelog

## [Unreleased]

### Added
- `core/middleware/ratelimit.go` — `RateLimitMiddleware()` with configurable rate via `RATE_LIMIT` env var (default `"60-M"`).
- `"ratelimit"` middleware alias in `middleware_provider.go`.
- `github.com/ulule/limiter/v3` dependency (memory store + Gin driver).
- 6 test cases (TC-32 to TC-37) for rate limiting behaviour.

### Changed
- `app/providers/middleware_provider.go` — added `"ratelimit"` alias registration.

### Deviation log
_None expected._
