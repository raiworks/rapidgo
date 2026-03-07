# 📝 Changelog: OAuth2 / Social Login

> **Feature**: `49` — OAuth2 / Social Login
> **Status**: ✅ SHIPPED
> **Date**: 2026-03-07

---

## Added

- `core/oauth/oauth.go` — `Provider`, `UserInfo`, `Google()`, `GitHub()`, `NewProvider()`, `AuthCodeURL()`, `Exchange()`, `FetchUser()`, `GenerateState()`
- `core/oauth/oauth_test.go` — 15 unit tests for provider config, auth flow, user info normalization

## Dependencies

| Package | Version | License | Purpose |
|---------|---------|---------|---------|
| `golang.org/x/oauth2` | latest | BSD-3-Clause | OAuth2 authorization code flow |

## Files

| File | Action |
|------|--------|
| `core/oauth/oauth.go` | NEW |
| `core/oauth/oauth_test.go` | NEW |

## Migration Guide

- No migrations required
- No new framework-level environment variables
- No breaking changes — new package, no existing code modified
- App-level env vars: `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GITHUB_CLIENT_ID`, `GITHUB_CLIENT_SECRET` (app reads them, not the framework)
- Use `oauth.Google(...)` or `oauth.GitHub(...)` to create providers
- Use `provider.AuthCodeURL(state)` for redirect, `provider.Exchange(ctx, code)` + `provider.FetchUser(ctx, token)` for callback
