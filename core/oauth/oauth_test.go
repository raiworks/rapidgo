package oauth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"golang.org/x/oauth2"
)

// T01: Google returns correctly configured Provider
func TestGoogle_ReturnsProvider(t *testing.T) {
	p := Google("cid", "csecret", "http://localhost/cb")
	if p == nil {
		t.Fatal("Google returned nil")
	}
	if p.Name != "google" {
		t.Fatalf("expected name 'google', got %q", p.Name)
	}
	if p.AuthURL != "https://accounts.google.com/o/oauth2/v2/auth" {
		t.Fatalf("unexpected AuthURL: %s", p.AuthURL)
	}
	if p.TokenURL != "https://oauth2.googleapis.com/token" {
		t.Fatalf("unexpected TokenURL: %s", p.TokenURL)
	}
	if p.UserInfoURL != "https://www.googleapis.com/oauth2/v2/userinfo" {
		t.Fatalf("unexpected UserInfoURL: %s", p.UserInfoURL)
	}
	if len(p.Scopes) != 3 {
		t.Fatalf("expected 3 scopes, got %d", len(p.Scopes))
	}
}

// T02: GitHub returns correctly configured Provider
func TestGitHub_ReturnsProvider(t *testing.T) {
	p := GitHub("cid", "csecret", "http://localhost/cb")
	if p == nil {
		t.Fatal("GitHub returned nil")
	}
	if p.Name != "github" {
		t.Fatalf("expected name 'github', got %q", p.Name)
	}
	if p.AuthURL != "https://github.com/login/oauth/authorize" {
		t.Fatalf("unexpected AuthURL: %s", p.AuthURL)
	}
	if p.TokenURL != "https://github.com/login/oauth/access_token" {
		t.Fatalf("unexpected TokenURL: %s", p.TokenURL)
	}
	if p.UserInfoURL != "https://api.github.com/user" {
		t.Fatalf("unexpected UserInfoURL: %s", p.UserInfoURL)
	}
	if len(p.Scopes) != 1 || p.Scopes[0] != "user:email" {
		t.Fatalf("unexpected scopes: %v", p.Scopes)
	}
}

// T03: NewProvider creates a custom provider
func TestNewProvider_CustomProvider(t *testing.T) {
	parser := func(raw map[string]any) UserInfo {
		return UserInfo{ID: str(raw["sub"])}
	}
	p := NewProvider("custom", "c", "s", "http://cb", []string{"read"}, "http://auth", "http://token", "http://user", parser)
	if p.Name != "custom" {
		t.Fatalf("expected name 'custom', got %q", p.Name)
	}
	if p.ClientID != "c" || p.ClientSecret != "s" {
		t.Fatal("client credentials not set")
	}
	if p.RedirectURL != "http://cb" {
		t.Fatalf("unexpected RedirectURL: %s", p.RedirectURL)
	}
	if p.AuthURL != "http://auth" || p.TokenURL != "http://token" || p.UserInfoURL != "http://user" {
		t.Fatal("endpoints not set correctly")
	}
	if p.parseUser == nil {
		t.Fatal("parseUser is nil")
	}
}

// T04: AuthCodeURL contains state parameter
func TestAuthCodeURL_ContainsState(t *testing.T) {
	p := Google("cid", "csecret", "http://localhost/cb")
	u := p.AuthCodeURL("test-state-123")

	parsed, err := url.Parse(u)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}
	if parsed.Query().Get("state") != "test-state-123" {
		t.Fatalf("expected state 'test-state-123', got %q", parsed.Query().Get("state"))
	}
}

// T05: AuthCodeURL contains client_id
func TestAuthCodeURL_ContainsClientID(t *testing.T) {
	p := Google("my-client-id", "csecret", "http://localhost/cb")
	u := p.AuthCodeURL("state")

	parsed, err := url.Parse(u)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}
	if parsed.Query().Get("client_id") != "my-client-id" {
		t.Fatalf("expected client_id 'my-client-id', got %q", parsed.Query().Get("client_id"))
	}
}

// T06: AuthCodeURL contains redirect_uri
func TestAuthCodeURL_ContainsRedirectURI(t *testing.T) {
	p := Google("cid", "csecret", "http://localhost:8080/auth/callback")
	u := p.AuthCodeURL("state")

	parsed, err := url.Parse(u)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}
	if parsed.Query().Get("redirect_uri") != "http://localhost:8080/auth/callback" {
		t.Fatalf("unexpected redirect_uri: %q", parsed.Query().Get("redirect_uri"))
	}
}

// T07: AuthCodeURL contains scope
func TestAuthCodeURL_ContainsScope(t *testing.T) {
	p := Google("cid", "csecret", "http://localhost/cb")
	u := p.AuthCodeURL("state")

	parsed, err := url.Parse(u)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}
	scope := parsed.Query().Get("scope")
	if scope == "" {
		t.Fatal("scope parameter is empty")
	}
	// Google scopes: openid email profile
	for _, s := range []string{"openid", "email", "profile"} {
		if !strings.Contains(scope, s) {
			t.Fatalf("scope missing %q: %s", s, scope)
		}
	}
}

// T08: Exchange returns token on valid code (mock server)
func TestExchange_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"access_token": "mock-access-token",
			"token_type":   "Bearer",
		})
	}))
	defer ts.Close()

	p := NewProvider("test", "cid", "csecret", "http://cb", nil, "http://auth", ts.URL, "http://user", nil)

	token, err := p.Exchange(context.Background(), "valid-code")
	if err != nil {
		t.Fatalf("Exchange failed: %v", err)
	}
	if token == nil {
		t.Fatal("expected non-nil token")
	}
	if token.AccessToken != "mock-access-token" {
		t.Fatalf("expected 'mock-access-token', got %q", token.AccessToken)
	}
}

// T09: Exchange returns error on invalid code
func TestExchange_InvalidCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"invalid_grant"}`))
	}))
	defer ts.Close()

	p := NewProvider("test", "cid", "csecret", "http://cb", nil, "http://auth", ts.URL, "http://user", nil)

	_, err := p.Exchange(context.Background(), "bad-code")
	if err == nil {
		t.Fatal("expected error for invalid code")
	}
}

// T10: FetchUser returns normalized Google user (mock)
func TestFetchUser_Google(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"id":      "123456",
			"email":   "user@gmail.com",
			"name":    "Test User",
			"picture": "https://example.com/photo.jpg",
		})
	}))
	defer ts.Close()

	p := Google("cid", "csecret", "http://cb")
	p.UserInfoURL = ts.URL

	token := &oauth2.Token{AccessToken: "mock-token", TokenType: "Bearer"}
	info, err := p.FetchUser(context.Background(), token)
	if err != nil {
		t.Fatalf("FetchUser failed: %v", err)
	}
	if info.ID != "123456" {
		t.Fatalf("expected ID '123456', got %q", info.ID)
	}
	if info.Email != "user@gmail.com" {
		t.Fatalf("expected email 'user@gmail.com', got %q", info.Email)
	}
	if info.Name != "Test User" {
		t.Fatalf("expected name 'Test User', got %q", info.Name)
	}
	if info.AvatarURL != "https://example.com/photo.jpg" {
		t.Fatalf("expected avatar URL, got %q", info.AvatarURL)
	}
}

// T11: FetchUser returns normalized GitHub user (mock)
func TestFetchUser_GitHub(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"id":         float64(789),
			"email":      "user@github.com",
			"name":       "Test User",
			"login":      "testuser",
			"avatar_url": "https://example.com/avatar.jpg",
		})
	}))
	defer ts.Close()

	p := GitHub("cid", "csecret", "http://cb")
	p.UserInfoURL = ts.URL

	token := &oauth2.Token{AccessToken: "mock-token", TokenType: "Bearer"}
	info, err := p.FetchUser(context.Background(), token)
	if err != nil {
		t.Fatalf("FetchUser failed: %v", err)
	}
	if info.ID != "789" {
		t.Fatalf("expected ID '789', got %q", info.ID)
	}
	if info.Email != "user@github.com" {
		t.Fatalf("expected email 'user@github.com', got %q", info.Email)
	}
	if info.Name != "Test User" {
		t.Fatalf("expected name 'Test User', got %q", info.Name)
	}
	if info.AvatarURL != "https://example.com/avatar.jpg" {
		t.Fatalf("expected avatar URL, got %q", info.AvatarURL)
	}
}

// T12: FetchUser includes raw JSON in UserInfo.Raw
func TestFetchUser_RawField(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"id":            "123",
			"email":         "user@test.com",
			"name":          "Raw User",
			"custom_field":  "custom_value",
		})
	}))
	defer ts.Close()

	p := Google("cid", "csecret", "http://cb")
	p.UserInfoURL = ts.URL

	token := &oauth2.Token{AccessToken: "mock-token", TokenType: "Bearer"}
	info, err := p.FetchUser(context.Background(), token)
	if err != nil {
		t.Fatalf("FetchUser failed: %v", err)
	}
	if info.Raw == nil {
		t.Fatal("expected Raw to be non-nil")
	}
	if info.Raw["custom_field"] != "custom_value" {
		t.Fatalf("expected custom_field in Raw, got %v", info.Raw["custom_field"])
	}
}

// T13: GenerateState returns 64-character hex string
func TestGenerateState_Length(t *testing.T) {
	state := GenerateState()
	if len(state) != 64 {
		t.Fatalf("expected 64 chars, got %d", len(state))
	}
	for _, c := range state {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Fatalf("non-hex character found: %c", c)
		}
	}
}

// T14: GenerateState produces unique values
func TestGenerateState_Unique(t *testing.T) {
	s1 := GenerateState()
	s2 := GenerateState()
	if s1 == s2 {
		t.Fatal("expected unique states, got identical values")
	}
}

// T15: GitHub FetchUser falls back to login when name is null
func TestFetchUser_GitHub_NameFallback(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// name is null — should fall back to login
		w.Write([]byte(`{"id":42,"email":"user@gh.com","name":null,"login":"fallbackuser","avatar_url":"https://img.com/a.jpg"}`))
	}))
	defer ts.Close()

	p := GitHub("cid", "csecret", "http://cb")
	p.UserInfoURL = ts.URL

	token := &oauth2.Token{AccessToken: "mock-token", TokenType: "Bearer"}
	info, err := p.FetchUser(context.Background(), token)
	if err != nil {
		t.Fatalf("FetchUser failed: %v", err)
	}
	if info.Name != "fallbackuser" {
		t.Fatalf("expected name 'fallbackuser' (from login), got %q", info.Name)
	}
}
