# 🧪 Test Plan: OAuth2 / Social Login

> **Feature**: `49` — OAuth2 / Social Login
> **Status**: ✅ SHIPPED
> **File**: `core/oauth/oauth_test.go`
> **Date**: 2026-03-07

---

## Test Cases

| TC | Function | Description | Approach |
|----|----------|-------------|----------|
| T01 | `TestGoogle_ReturnsProvider` | Google() returns non-nil Provider with correct config | Assert Name, AuthURL, TokenURL, UserInfoURL, Scopes |
| T02 | `TestGitHub_ReturnsProvider` | GitHub() returns non-nil Provider with correct config | Assert Name, AuthURL, TokenURL, UserInfoURL, Scopes |
| T03 | `TestNewProvider_CustomProvider` | NewProvider creates a custom provider | Create with custom endpoints, assert all fields |
| T04 | `TestAuthCodeURL_ContainsState` | AuthCodeURL includes state parameter | Parse returned URL, assert `state` query param matches |
| T05 | `TestAuthCodeURL_ContainsClientID` | AuthCodeURL includes client_id | Parse returned URL, assert `client_id` query param |
| T06 | `TestAuthCodeURL_ContainsRedirectURI` | AuthCodeURL includes redirect_uri | Parse returned URL, assert `redirect_uri` query param |
| T07 | `TestAuthCodeURL_ContainsScope` | AuthCodeURL includes scopes | Parse returned URL, assert `scope` query param |
| T08 | `TestExchange_Success` | Exchange returns token on valid code | Mock token server, exchange code, assert token not nil |
| T09 | `TestExchange_InvalidCode` | Exchange returns error on invalid code | Mock server returning 401, assert error |
| T10 | `TestFetchUser_Google` | FetchUser returns normalized Google user | Mock user info server with Google JSON, assert fields |
| T11 | `TestFetchUser_GitHub` | FetchUser returns normalized GitHub user | Mock user info server with GitHub JSON, assert fields |
| T12 | `TestFetchUser_RawField` | FetchUser includes raw JSON in UserInfo.Raw | Mock server, assert Raw contains original fields |
| T13 | `TestGenerateState_Length` | GenerateState returns 64-char hex string | Call, assert len == 64, all chars in [0-9a-f] |
| T14 | `TestGenerateState_Unique` | GenerateState produces unique values | Call twice, assert different |
| T15 | `TestFetchUser_GitHub_NameFallback` | FetchUser falls back to `login` when `name` is null | Mock server with `name: null`, assert Name == `login` value |

---

## Test Strategy

- T01–T07 are pure unit tests — no network, no mocks
- T08–T09 use `httptest.NewServer` to mock the OAuth2 token endpoint
- T10–T12 use `httptest.NewServer` to mock the user info API endpoint
- T13–T14 are pure unit tests
- For T08–T09: Provider's `TokenURL` is pointed at the mock server
- For T10–T12: Provider's `UserInfoURL` is pointed at the mock server; a pre-constructed `oauth2.Token` is passed to `FetchUser`

---

## Mock Server Patterns

### Token endpoint mock (T08–T09)
Returns JSON `{"access_token": "mock-token", "token_type": "Bearer"}` on success, HTTP 401 on failure.

### User info endpoint mock (T10–T12)
Returns provider-specific JSON:

**Google mock response:**
```json
{
    "id": "123456",
    "email": "user@gmail.com",
    "name": "Test User",
    "picture": "https://example.com/photo.jpg"
}
```

**GitHub mock response:**
```json
{
    "id": 789,
    "email": "user@github.com",
    "name": "Test User",
    "login": "testuser",
    "avatar_url": "https://example.com/avatar.jpg"
}
```

---

## Pass Criteria

- All 15 tests pass: `go test ./core/oauth/ -count=1 -v`
- Full regression: `go test ./... -count=1`
- Binary builds: `go build -o bin/rapidgo.exe ./cmd`
