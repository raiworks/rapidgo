package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

// UserInfo holds normalized user profile data from an OAuth2 provider.
type UserInfo struct {
	ID        string
	Email     string
	Name      string
	AvatarURL string
	Raw       map[string]any
}

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

// NewProvider creates a custom OAuth2 provider with the given configuration.
func NewProvider(name, clientID, clientSecret, redirectURL string, scopes []string, authURL, tokenURL, userInfoURL string, parser func(map[string]any) UserInfo) *Provider {
	return &Provider{
		Name:         name,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		AuthURL:      authURL,
		TokenURL:     tokenURL,
		UserInfoURL:  userInfoURL,
		parseUser:    parser,
	}
}

// Google creates a pre-configured Google OAuth2 provider.
func Google(clientID, clientSecret, redirectURL string) *Provider {
	return NewProvider(
		"google",
		clientID, clientSecret, redirectURL,
		[]string{"openid", "email", "profile"},
		"https://accounts.google.com/o/oauth2/v2/auth",
		"https://oauth2.googleapis.com/token",
		"https://www.googleapis.com/oauth2/v2/userinfo",
		func(raw map[string]any) UserInfo {
			return UserInfo{
				ID:        str(raw["id"]),
				Email:     str(raw["email"]),
				Name:      str(raw["name"]),
				AvatarURL: str(raw["picture"]),
			}
		},
	)
}

// GitHub creates a pre-configured GitHub OAuth2 provider.
func GitHub(clientID, clientSecret, redirectURL string) *Provider {
	return NewProvider(
		"github",
		clientID, clientSecret, redirectURL,
		[]string{"user:email"},
		"https://github.com/login/oauth/authorize",
		"https://github.com/login/oauth/access_token",
		"https://api.github.com/user",
		func(raw map[string]any) UserInfo {
			name := str(raw["name"])
			if name == "" {
				name = str(raw["login"])
			}
			return UserInfo{
				ID:        fmt.Sprintf("%v", raw["id"]),
				Email:     str(raw["email"]),
				Name:      name,
				AvatarURL: str(raw["avatar_url"]),
			}
		},
	)
}

// oauthConfig builds an oauth2.Config from the Provider's fields.
func (p *Provider) oauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
		RedirectURL:  p.RedirectURL,
		Scopes:       p.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  p.AuthURL,
			TokenURL: p.TokenURL,
		},
	}
}

// AuthCodeURL builds the authorization redirect URL with the provided state.
func (p *Provider) AuthCodeURL(state string) string {
	return p.oauthConfig().AuthCodeURL(state)
}

// Exchange exchanges the authorization code for an access token.
func (p *Provider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.oauthConfig().Exchange(ctx, code)
}

// FetchUser fetches and normalizes the user's profile from the provider API.
func (p *Provider) FetchUser(ctx context.Context, token *oauth2.Token) (UserInfo, error) {
	client := p.oauthConfig().Client(ctx, token)

	resp, err := client.Get(p.UserInfoURL)
	if err != nil {
		return UserInfo{}, fmt.Errorf("oauth: failed to fetch user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UserInfo{}, fmt.Errorf("oauth: user info returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserInfo{}, fmt.Errorf("oauth: failed to read user info body: %w", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return UserInfo{}, fmt.Errorf("oauth: failed to decode user info: %w", err)
	}

	info := p.parseUser(raw)
	info.Raw = raw
	return info, nil
}

// GenerateState returns a cryptographically random 64-character hex string
// for use as the OAuth2 state parameter.
func GenerateState() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("oauth: failed to generate random state: " + err.Error())
	}
	return hex.EncodeToString(b)
}

// str safely converts an interface value to a string.
func str(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}
