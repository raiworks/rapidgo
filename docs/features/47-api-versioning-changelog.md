# 📝 Changelog: API Versioning

> **Feature**: `47` — API Versioning
> **Status**: 🟡 IN PROGRESS
> **Date**: 2026-03-07

---

## Added

- `core/router/version.go` — `Version()`, `DeprecatedVersion()`, deprecation headers middleware
- `core/router/version_test.go` — 10 unit tests for versioned routing and deprecation headers

## Files

| File | Action |
|---|---|
| `core/router/version.go` | NEW |
| `core/router/version_test.go` | NEW |

## Migration Guide

- No migrations required
- No new environment variables
- No new dependencies
- No breaking changes — existing routes are unaffected
- Use `r.Version("v1")` instead of `r.Group("/api/v1")` for versioned API groups
- Use `r.DeprecatedVersion("v1", sunsetDate)` to add `Sunset` and `X-API-Deprecated` headers
