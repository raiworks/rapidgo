# 🏗️ Architecture: OAuth2 / Social Login

> **Feature**: `49` — OAuth2 / Social Login
> **Status**: ✅ SHIPPED
> **Package**: `core/oauth`
> **Date**: 2026-03-07

---

## Component Overview

```
core/oauth/
    │
    └── oauth.go              (NEW)
        ├── Provider struct
        │   ├── Name         string
        │   ├── ClientID     string
        │   ├── ClientSecret string
        │   ├── RedirectURL  string
        │   ├── Scopes       []string
        │   ├── AuthURL      string
        │   ├── TokenURL     string
        │   ├── UserInfoURL  string
        │   └── parseUser    func(map[string]any) UserInfo
        │
        ├── UserInfo struct
        │   ├── ID        string
        │   ├── Email     string
        │   ├── Name      string
        │   ├── AvatarURL string
        │   └── Raw       map[string]any
        │
        ├── NewProvider(name, clientID, clientSecret, redirectURL string, scopes []string, authURL, tokenURL, userInfoURL string, parser func(map[string]any) UserInfo) *Provider
        ├── Google(clientID, clientSecret, redirectURL string) *Provider
        ├── GitHub(clientID, clientSecret, redirectURL string) *Provider
        ├── (*Provider).AuthCodeURL(state string) string
        ├── (*Provider).Exchange(ctx context.Context, code string) (*oauth2.Token, error)
        ├── (*Provider).FetchUser(ctx context.Context, token *oauth2.Token) (UserInfo, error)
        └── GenerateState() string
```

---

## Structs

### Provider

```go
// Provider configures an OAuth2 provider (e.g., Google, GitHub).
type Provider struct {
    Name         string
    ClientID     string
    ClientSecret string
    RedirectURL  string
    Scopes       []string
    AuthURL      string
    TokenURL     string
    UserInfoURL  string
    parseUser    func(map[string]any) UserInfo
}
```

| Field | Type | Visibility | Description |
|-------|------|------------|-------------|
| `Name` | `string` | Exported | Provider name (e.g., "google", "github") |
| `ClientID` | `string` | Exported | OAuth2 client ID |
| `ClientSecret` | `string` | Exported | OAuth2 client secret |
| `RedirectURL` | `string` | Exported | Callback URL (e.g., `http://localhost:8080/auth/google/callback`) |
| `Scopes` | `[]string` | Exported | OAuth2 scopes to request |
| `AuthURL` | `string` | Exported | Provider's authorization endpoint |
| `TokenURL` | `string` | Exported | Provider's token endpoint |
| `UserInfoURL` | `string` | Exported | Provider's user info API endpoint |
| `parseUser` | `func(map[string]any) UserInfo` | Unexported | Maps raw provider JSON to normalized `UserInfo` |

### UserInfo

```go
// UserInfo holds normalized user profile data from an OAuth2 provider.
type UserInfo struct {
    ID        string
    Email     string
    Name      string
    AvatarURL string
    Raw       map[string]any
}
```

| Field | Type | Description |
|-------|------|-------------|
| `ID` | `string` | Provider-specific user ID (always a string for uniformity) |
| `Email` | `string` | User's email (may be empty if provider doesn't return it) |
| `Name` | `string` | User's display name |
| `AvatarURL` | `string` | User's avatar/profile picture URL |
| `Raw` | `map[string]any` | Full raw JSON response from the provider |

---

## Functions

### NewProvider

```go
func NewProvider(name, clientID, clientSecret, redirectURL string, scopes []string, authURL, tokenURL, userInfoURL string, parser func(map[string]any) UserInfo) *Provider
```

Creates a custom OAuth2 provider with all fields. Used by `Google()` and `GitHub()` internally, and by developers for custom providers.

### Google

```go
func Google(clientID, clientSecret, redirectURL string) *Provider
```

Creates a pre-configured Google OAuth2 provider:
- **AuthURL**: `https://accounts.google.com/o/oauth2/v2/auth`
- **TokenURL**: `https://oauth2.googleapis.com/token`
- **UserInfoURL**: `https://www.googleapis.com/oauth2/v2/userinfo`
- **Scopes**: `["openid", "email", "profile"]`
- **parseUser**: Maps `id`, `email`, `name`, `picture`

### GitHub

```go
func GitHub(clientID, clientSecret, redirectURL string) *Provider
```

Creates a pre-configured GitHub OAuth2 provider:
- **AuthURL**: `https://github.com/login/oauth/authorize`
- **TokenURL**: `https://github.com/login/oauth/access_token`
- **UserInfoURL**: `https://api.github.com/user`
- **Scopes**: `["user:email"]`
- **parseUser**: Maps `id` (converted to string), `email`, `name` (falls back to `login`), `avatar_url`

### (*Provider).AuthCodeURL

```go
func (p *Provider) AuthCodeURL(state string) string
```

Builds the authorization redirect URL with the configured client ID, redirect URI, scopes, and the provided state parameter. Delegates to `oauth2.Config.AuthCodeURL(state)`.

### (*Provider).Exchange

```go
func (p *Provider) Exchange(ctx context.Context, code string) (*oauth2.Token, error)
```

Exchanges the authorization code for an access token. Delegates to `oauth2.Config.Exchange(ctx, code)`.

### (*Provider).FetchUser

```go
func (p *Provider) FetchUser(ctx context.Context, token *oauth2.Token) (UserInfo, error)
```

1. Creates an HTTP client with the token (`p.oauthConfig().Client(ctx, token)`)
2. GETs `p.UserInfoURL`
3. Decodes the JSON response into `map[string]any`
4. Calls `p.parseUser(raw)` to normalize the data
5. Sets `UserInfo.Raw` to the full response
6. Returns `UserInfo`

### GenerateState

```go
func GenerateState() string
```

Generates a cryptographically random 32-byte hex-encoded string for use as the OAuth2 `state` parameter. Uses `crypto/rand`.

---

## Internal Helper

### (*Provider).oauthConfig

```go
func (p *Provider) oauthConfig() *oauth2.Config
```

Unexported method that constructs an `oauth2.Config` from the Provider's fields. Used internally by `AuthCodeURL`, `Exchange`, and `FetchUser`.

---

## Dependencies

| Package | Purpose | Status |
|---------|---------|--------|
| `golang.org/x/oauth2` | OAuth2 authorization code flow | **NEW** — needs `go get` |
| `crypto/rand` | State parameter generation | Standard library |
| `encoding/hex` | Hex encoding for state | Standard library |
| `encoding/json` | JSON decoding of user info response | Standard library |
| `fmt` | String conversion (GitHub ID int→string) | Standard library |
| `io` | Reading HTTP response body | Standard library |
| `net/http` | HTTP status code check | Standard library |

---

## Files Changed

| File | Action | Description |
|------|--------|-------------|
| `core/oauth/oauth.go` | NEW | Provider, UserInfo, Google, GitHub, auth flow functions |
| `core/oauth/oauth_test.go` | NEW | All test cases |

No existing files modified.

---

## Usage Example

```go
import (
    "github.com/RAiWorks/RapidGo/core/oauth"
)

// Setup
google := oauth.Google(
    os.Getenv("GOOGLE_CLIENT_ID"),
    os.Getenv("GOOGLE_CLIENT_SECRET"),
    "http://localhost:8080/auth/google/callback",
)

// Route: Redirect to Google
r.Get("/auth/google", func(c *gin.Context) {
    state := oauth.GenerateState()
    // Store state in session for verification
    session.Set("oauth_state", state)
    session.Save()
    c.Redirect(http.StatusTemporaryRedirect, google.AuthCodeURL(state))
})

// Route: Handle callback
r.Get("/auth/google/callback", func(c *gin.Context) {
    // Verify state
    if c.Query("state") != session.Get("oauth_state") {
        c.AbortWithStatus(http.StatusForbidden)
        return
    }

    token, err := google.Exchange(c.Request.Context(), c.Query("code"))
    if err != nil { /* handle error */ }

    user, err := google.FetchUser(c.Request.Context(), token)
    if err != nil { /* handle error */ }

    // user.Email, user.Name, user.ID — create/link account
})
```

---

## Google parseUser Mapping

```
Google JSON             → UserInfo
────────────────────────────────────
"id"                    → ID
"email"                 → Email
"name"                  → Name
"picture"               → AvatarURL
(full response)         → Raw
```

## GitHub parseUser Mapping

```
GitHub JSON             → UserInfo
────────────────────────────────────
"id" (number→string)    → ID
"email"                 → Email
"name" (or "login")     → Name
"avatar_url"            → AvatarURL
(full response)         → Raw
```

---

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `GOOGLE_CLIENT_ID` | If using Google | Google OAuth2 client ID |
| `GOOGLE_CLIENT_SECRET` | If using Google | Google OAuth2 client secret |
| `GITHUB_CLIENT_ID` | If using GitHub | GitHub OAuth2 client ID |
| `GITHUB_CLIENT_SECRET` | If using GitHub | GitHub OAuth2 client secret |

These are **not** read by the framework — the app reads them and passes to the provider constructors. No framework env coupling.
