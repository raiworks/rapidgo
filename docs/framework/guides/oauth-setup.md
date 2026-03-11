---
title: "OAuth2 Provider Setup"
version: "0.1.0"
status: "Final"
date: "2026-03-11"
last_updated: "2026-03-11"
authors:
  - "RAiWorks"
supersedes: ""
---

# OAuth2 Provider Setup

## Abstract

Step-by-step guide to adding Google, GitHub, or custom OAuth2 login to a
RapidGo application using the `core/oauth` package.

## Table of Contents

1. [Overview](#1-overview)
2. [Get Provider Credentials](#2-get-provider-credentials)
3. [Configure the Provider](#3-configure-the-provider)
4. [Redirect Route](#4-redirect-route)
5. [Callback Route](#5-callback-route)
6. [Create or Login User](#6-create-or-login-user)
7. [Custom Providers](#7-custom-providers)
8. [Security Considerations](#8-security-considerations)
9. [Full Example](#9-full-example)
10. [References](#10-references)

---

## 1. Overview

The `core/oauth` package provides:

- **Pre-built providers**: `oauth.Google()`, `oauth.GitHub()`
- **Custom providers**: `oauth.NewProvider()` with your own endpoints
- **UserInfo normalization**: ID, Email, Name, AvatarURL from any provider
- **State generation**: `oauth.GenerateState()` for CSRF protection

Flow: Redirect → Provider Login → Callback → Exchange Code → Fetch User → Login/Register.

## 2. Get Provider Credentials

### Google

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a project (or select existing)
3. Navigate to **APIs & Services → Credentials**
4. Click **Create Credentials → OAuth 2.0 Client ID**
5. Set Application Type to **Web application**
6. Add Authorized Redirect URI: `https://yourapp.com/auth/google/callback`
7. Copy the **Client ID** and **Client Secret**

### GitHub

1. Go to [GitHub Developer Settings](https://github.com/settings/developers)
2. Click **New OAuth App**
3. Set Homepage URL: `https://yourapp.com`
4. Set Authorization Callback URL: `https://yourapp.com/auth/github/callback`
5. Copy the **Client ID** and **Client Secret**

### Environment Variables

```env
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret

GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret

OAUTH_REDIRECT_BASE=https://yourapp.com
```

## 3. Configure the Provider

```go
import "github.com/RAiWorks/RapidGo/v2/core/oauth"

googleProvider := oauth.Google(
    os.Getenv("GOOGLE_CLIENT_ID"),
    os.Getenv("GOOGLE_CLIENT_SECRET"),
    os.Getenv("OAUTH_REDIRECT_BASE")+"/auth/google/callback",
)

githubProvider := oauth.GitHub(
    os.Getenv("GITHUB_CLIENT_ID"),
    os.Getenv("GITHUB_CLIENT_SECRET"),
    os.Getenv("OAUTH_REDIRECT_BASE")+"/auth/github/callback",
)
```

## 4. Redirect Route

Redirect the user to the provider's login page:

```go
r.GET("/auth/google/redirect", func(c *gin.Context) {
    state := oauth.GenerateState()

    // Store state in session for CSRF validation
    session := sessions.Default(c)
    session.Set("oauth_state", state)
    session.Save()

    url := googleProvider.AuthCodeURL(state)
    c.Redirect(http.StatusTemporaryRedirect, url)
})
```

## 5. Callback Route

Handle the provider's callback after the user authenticates:

```go
r.GET("/auth/google/callback", func(c *gin.Context) {
    // Verify state matches (CSRF protection)
    session := sessions.Default(c)
    savedState := session.Get("oauth_state")
    if c.Query("state") != savedState {
        c.AbortWithStatusJSON(400, gin.H{"error": "invalid state"})
        return
    }

    // Exchange code for token
    code := c.Query("code")
    token, err := googleProvider.Exchange(c.Request.Context(), code)
    if err != nil {
        c.AbortWithStatusJSON(500, gin.H{"error": "token exchange failed"})
        return
    }

    // Fetch user profile
    userInfo, err := googleProvider.FetchUser(c.Request.Context(), token)
    if err != nil {
        c.AbortWithStatusJSON(500, gin.H{"error": "failed to fetch user"})
        return
    }

    // userInfo.ID, userInfo.Email, userInfo.Name, userInfo.AvatarURL
    // → create or login user (see next section)
})
```

## 6. Create or Login User

```go
func findOrCreateOAuthUser(db *gorm.DB, info oauth.UserInfo, provider string) (*User, error) {
    var user User
    err := db.Where("oauth_provider = ? AND oauth_id = ?", provider, info.ID).First(&user).Error
    if err == nil {
        return &user, nil // existing user
    }

    // Create new user
    user = User{
        Name:          info.Name,
        Email:         info.Email,
        AvatarURL:     info.AvatarURL,
        OAuthProvider: provider,
        OAuthID:       info.ID,
    }
    return &user, db.Create(&user).Error
}
```

## 7. Custom Providers

For providers not built in (Facebook, Discord, etc.):

```go
facebookProvider := oauth.NewProvider(
    "facebook",
    os.Getenv("FACEBOOK_CLIENT_ID"),
    os.Getenv("FACEBOOK_CLIENT_SECRET"),
    os.Getenv("OAUTH_REDIRECT_BASE")+"/auth/facebook/callback",
    []string{"email", "public_profile"},
    "https://www.facebook.com/v18.0/dialog/oauth",
    "https://graph.facebook.com/v18.0/oauth/access_token",
    "https://graph.facebook.com/v18.0/me?fields=id,name,email,picture.type(large)",
    func(raw map[string]any) oauth.UserInfo {
        avatar := ""
        if pic, ok := raw["picture"].(map[string]any); ok {
            if data, ok := pic["data"].(map[string]any); ok {
                avatar, _ = data["url"].(string)
            }
        }
        return oauth.UserInfo{
            ID:        fmt.Sprintf("%v", raw["id"]),
            Email:     fmt.Sprintf("%v", raw["email"]),
            Name:      fmt.Sprintf("%v", raw["name"]),
            AvatarURL: avatar,
        }
    },
)
```

## 8. Security Considerations

- **Always validate the state parameter** — prevents CSRF attacks.
- **Use HTTPS** in production redirect URLs.
- **Store credentials in environment variables**, never in code.
- **Validate email ownership** — some providers return unverified emails.
- **Rate-limit** the callback endpoint to prevent abuse.

## 9. Full Example

```go
func setupOAuth(r *gin.Engine, db *gorm.DB) {
    google := oauth.Google(
        os.Getenv("GOOGLE_CLIENT_ID"),
        os.Getenv("GOOGLE_CLIENT_SECRET"),
        os.Getenv("OAUTH_REDIRECT_BASE")+"/auth/google/callback",
    )

    r.GET("/auth/google/redirect", func(c *gin.Context) {
        state := oauth.GenerateState()
        session := sessions.Default(c)
        session.Set("oauth_state", state)
        session.Save()
        c.Redirect(http.StatusTemporaryRedirect, google.AuthCodeURL(state))
    })

    r.GET("/auth/google/callback", func(c *gin.Context) {
        session := sessions.Default(c)
        if c.Query("state") != session.Get("oauth_state") {
            c.AbortWithStatus(403)
            return
        }

        token, err := google.Exchange(c.Request.Context(), c.Query("code"))
        if err != nil {
            c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
            return
        }

        info, err := google.FetchUser(c.Request.Context(), token)
        if err != nil {
            c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
            return
        }

        user, err := findOrCreateOAuthUser(db, info, "google")
        if err != nil {
            c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
            return
        }

        // Issue JWT or set session
        jwtToken, _ := auth.GenerateToken(user.ID, user.Email)
        c.JSON(200, gin.H{"token": jwtToken})
    })
}
```

## 10. References

- [OAuth package source](../../core/oauth/oauth.go) — implementation
- [Auth / JWT](../security/authentication.md) — JWT token generation
