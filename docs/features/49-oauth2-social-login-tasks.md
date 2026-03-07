# ✅ Tasks: OAuth2 / Social Login

> **Feature**: `49` — OAuth2 / Social Login
> **Status**: ✅ SHIPPED
> **Date**: 2026-03-07

---

## Implementation Tasks

### T1: Add `golang.org/x/oauth2` dependency

- [ ] Run `go get golang.org/x/oauth2`
- [ ] Verify it appears in `go.mod` and `go.sum`

### T2: Create `core/oauth/oauth.go`

- [ ] Define `UserInfo` struct (`ID`, `Email`, `Name`, `AvatarURL`, `Raw`)
- [ ] Define `Provider` struct (all fields per architecture)
- [ ] Implement `NewProvider(...)` constructor
- [ ] Implement `Google(clientID, clientSecret, redirectURL)` with pre-configured endpoints and parser
- [ ] Implement `GitHub(clientID, clientSecret, redirectURL)` with pre-configured endpoints and parser
- [ ] Implement `(*Provider).oauthConfig()` (unexported) — builds `oauth2.Config`
- [ ] Implement `(*Provider).AuthCodeURL(state string) string`
- [ ] Implement `(*Provider).Exchange(ctx, code) (*oauth2.Token, error)`
- [ ] Implement `(*Provider).FetchUser(ctx, token) (UserInfo, error)`
- [ ] Implement `GenerateState() string`

### T3: Create `core/oauth/oauth_test.go`

- [ ] Write all test cases from the test plan (T01–T15)
- [ ] All tests pass with `go test ./core/oauth/ -count=1`

### T4: Regression

- [ ] `go test ./... -count=1` — all packages pass
- [ ] `go build -o bin/rapidgo.exe ./cmd` — binary builds clean

---

## Acceptance Criteria

1. `Google()` returns a correctly configured Provider
2. `GitHub()` returns a correctly configured Provider
3. `NewProvider()` creates a custom Provider
4. `AuthCodeURL()` returns a valid authorization URL with state
5. `Exchange()` exchanges authorization code for token (tested with mock server)
6. `FetchUser()` fetches and normalizes user info (tested with mock server)
7. `GenerateState()` returns a unique 64-character hex string
8. Google parseUser correctly maps `id`, `email`, `name`, `picture`
9. GitHub parseUser correctly maps `id` (int→string), `email`, `name` (falls back to `login`), `avatar_url`
10. No existing tests broken
