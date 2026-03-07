# 💬 Discussion: OAuth2 / Social Login

> **Feature**: `49` — OAuth2 / Social Login
> **Status**: ✅ SHIPPED
> **Depends on**: Feature #21 (Authentication) — ✅ SHIPPED
> **Date**: 2026-03-07

---

## What Problem Does This Solve?

Feature #21 provides JWT-based authentication with username/password. Many modern applications also need users to sign in via third-party providers (Google, GitHub, etc.) using the OAuth2 authorization code flow. Without framework support, every application rebuilds:

1. **OAuth2 redirect URLs** — building the authorization URL with state, scopes, etc.
2. **Callback handling** — exchanging the authorization code for tokens
3. **User info fetching** — getting the user's profile from the provider
4. **State verification** — CSRF protection via the `state` parameter
5. **Provider configuration** — wiring up client IDs, secrets, endpoints

Feature #49 provides a reusable OAuth2 integration layer in `core/oauth` that handles the protocol mechanics, leaving account creation/linking as an application-level concern.

---

## Scope

### In Scope

| Item | Description |
|------|-------------|
| Provider config | Struct to hold client ID, secret, scopes, endpoints |
| Built-in providers | Google and GitHub pre-configured (endpoints and user-info URLs) |
| Custom providers | Developers can define their own providers |
| Authorization URL | Generate the redirect URL with state and scopes |
| Code exchange | Exchange authorization code for access token |
| User info | Fetch user profile from the provider's API |
| State generation | Cryptographically random state parameter |
| State verification | Compare state from callback to stored state |
| UserInfo struct | Normalized user info (ID, Email, Name, AvatarURL, Raw) |

### Out of Scope

| Item | Rationale |
|------|-----------|
| Account creation/linking | Application-level — the framework provides the user info, the app decides what to do |
| Database models for OAuth accounts | Application-level concern |
| Token refresh | App-level — the framework provides the initial access token |
| Session management for OAuth state | App can use existing session system or any other storage |
| OAuth2 client_credentials / implicit flows | Only authorization code flow — the standard for web apps |
| OpenID Connect (OIDC) | Future enhancement — #49 covers plain OAuth2 |

---

## Design Decisions

| # | Decision | Choice | Rationale |
|---|----------|--------|-----------|
| D1 | Package location | `core/oauth/` | New package — OAuth2 is a distinct domain from JWT auth |
| D2 | OAuth2 library | `golang.org/x/oauth2` | Official Go OAuth2 package, widely used, well-maintained |
| D3 | Provider pattern | `Provider` struct with pre-built Google/GitHub constructors | Keeps it simple — no interface needed, just config structs |
| D4 | User info normalization | `UserInfo` struct with `ID`, `Email`, `Name`, `AvatarURL`, `Raw` | Every provider returns different JSON — normalization simplifies app code |
| D5 | State management | Framework generates state, app stores/verifies it | State is a CSRF token — storage is app-level (session, cookie, etc.) |
| D6 | HTTP client for user info | `oauth2.Client(ctx, token)` from x/oauth2 | Automatically attaches the Bearer token, handles transport |
| D7 | Error handling | Return errors, don't panic | OAuth2 involves external services — errors are expected |
| D8 | No middleware | Functions only — developer wires routes explicitly | OAuth2 routes are few and specific — middleware pattern doesn't fit |

---

## Open Questions

| # | Question | Resolution |
|---|----------|------------|
| Q1 | Should we support PKCE (Proof Key for Code Exchange)? | ✅ No — PKCE is for public clients (SPAs, mobile). Server-side web apps use client_secret. Future enhancement if needed. |
| Q2 | Should we provide route helpers (e.g., `RegisterOAuthRoutes`)? | ✅ No — just provide the functions. Route registration is 2 lines of code and app-specific. |
| Q3 | Should we store the raw provider response? | ✅ Yes — `UserInfo.Raw` holds the full JSON map. Developers can access provider-specific fields. |
| Q4 | How do we handle providers that don't return email? | ✅ `UserInfo.Email` can be empty. The app decides how to handle it (e.g., prompt user). |
| Q5 | Should we support multiple providers simultaneously? | ✅ Yes — each `Provider` is independent. The app can create as many as needed. |

---

## OAuth2 Authorization Code Flow

```
1. App redirects user → Provider's authorization URL
   (with client_id, redirect_uri, state, scopes)

2. User authenticates with provider, grants permissions

3. Provider redirects back to app's callback URL
   (with authorization code + state)

4. App verifies state, exchanges code for access token

5. App uses access token to fetch user info from provider API

6. App creates/links user account based on user info
```

Steps 1-5 are handled by `core/oauth`. Step 6 is application-level.

---

## Prior Art

| Framework | Pattern |
|-----------|---------|
| Laravel Socialite | Provider classes with `redirect()` and `user()` methods |
| Passport.js (Node) | Strategy pattern with serialize/deserialize |
| Django allauth | Full account management with social providers |
| Goth (Go) | Provider interface with `BeginAuth` / `CompleteAuth` |

The RapidGo approach is closest to Laravel Socialite — simple provider config with redirect + callback functions. No complex strategy pattern or account management.
