---
title: "Authentication"
version: "1.0.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Authentication

## Abstract

This document covers the two authentication strategies: JWT
(stateless, for APIs) and session-based (for SSR web apps), using
`golang-jwt/jwt/v5`.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Strategies](#2-strategies)
3. [JWT Authentication](#3-jwt-authentication)
4. [Auth Middleware](#4-auth-middleware)
5. [Session-based Authentication](#5-session-based-authentication)
6. [Security Considerations](#6-security-considerations)
7. [References](#7-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **JWT** — JSON Web Token; a compact, URL-safe token encoding
  claims.
- **Claims** — Key-value payload embedded in a JWT (e.g., `user_id`,
  `exp`).

## 2. Strategies

| Strategy | Best For | State |
|----------|----------|-------|
| **JWT** | REST APIs, mobile apps, microservices | Stateless (token-based) |
| **Session** | SSR web apps, browser-based UIs | Stateful (server-stored) |

Both can coexist — use JWT for API routes and sessions for web routes.

## 3. JWT Authentication

Library: `github.com/golang-jwt/jwt/v5`

### Token Generation

```go
package auth

import (
    "errors"
    "os"
    "strconv"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a signed JWT for the given user ID.
// Reads JWT_SECRET and JWT_EXPIRY (seconds) from environment.
func GenerateToken(userID uint) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "", errors.New("JWT_SECRET is not set")
    }
    if len(secret) < 32 {
        return "", errors.New("JWT_SECRET must be at least 32 bytes")
    }

    expiry := 3600 // default 1 hour
    if v := os.Getenv("JWT_EXPIRY"); v != "" {
        if parsed, err := strconv.Atoi(v); err == nil {
            expiry = parsed
        }
    }

    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Duration(expiry) * time.Second).Unix(),
        "iat":     time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

### Token Validation

```go
// ValidateToken parses and validates a JWT string.
// Returns the claims if the token is valid.
func ValidateToken(tokenStr string) (jwt.MapClaims, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return nil, errors.New("JWT_SECRET is not set")
    }
    if len(secret) < 32 {
        return nil, errors.New("JWT_SECRET must be at least 32 bytes")
    }

    token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(secret), nil
    })
    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return nil, jwt.ErrTokenUnverifiable
    }

    return claims, nil
}
```

### Login Endpoint

```go
func Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        responses.Error(c, 422, err.Error())
        return
    }

    // Verify credentials
    user, err := userSvc.FindByEmail(req.Email)
    if err != nil || !helpers.CheckPassword(user.Password, req.Password) {
        responses.Error(c, 401, "invalid credentials")
        return
    }

    token, err := auth.GenerateToken(user.ID)
    if err != nil {
        responses.Error(c, 500, "failed to generate token")
        return
    }

    responses.Success(c, gin.H{"token": token})
}
```

## 4. Auth Middleware

Validates the JWT from the `Authorization: Bearer <token>` header:

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        header := c.GetHeader("Authorization")

        if header == "" || !strings.HasPrefix(header, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "missing or invalid Authorization header",
            })
            return
        }

        token := strings.TrimPrefix(header, "Bearer ")
        claims, err := auth.ValidateToken(token)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "invalid or expired token",
            })
            return
        }

        c.Set("user_id", claims["user_id"])
        c.Next()
    }
}
```

## 5. Session-based Authentication

For SSR web apps, store the user ID in the session after login:

```go
func WebLogin(c *gin.Context) {
    email := c.PostForm("email")
    password := c.PostForm("password")

    user, err := userSvc.FindByEmail(email)
    if err != nil || !helpers.CheckPassword(user.Password, password) {
        // Flash error and redirect back
        sessionMgr.Flash(data, "error", "Invalid credentials")
        c.Set("session", data)
        c.Redirect(http.StatusFound, "/login")
        return
    }

    // Store user in session
    data["user_id"] = user.ID
    data["username"] = user.Name
    c.Set("session", data)
    c.Redirect(http.StatusFound, "/dashboard")
}
```

Check authentication in controllers:

```go
func Dashboard(c *gin.Context) {
    sess, _ := c.Get("session")
    data := sess.(map[string]interface{})

    username, _ := data["username"].(string)
    c.HTML(http.StatusOK, "dashboard.html", gin.H{
        "user": username,
    })
}
```

## 6. Security Considerations

- `JWT_SECRET` **MUST** be a strong, random value (minimum 32 bytes)
  and **MUST NOT** be committed to version control.
- JWT tokens **SHOULD** have short lifetimes (default: 1 hour via
  `JWT_EXPIRY` env var). Implement refresh tokens for longer sessions.
- Always use `HS256` or stronger signing methods. Never allow `none`
  algorithm.
- Password comparison **MUST** use `bcrypt.CompareHashAndPassword`
  (constant-time).
- Session-based auth routes **MUST** be protected by CSRF middleware.

## 7. References

- [Sessions](sessions.md)
- [CSRF Protection](csrf.md)
- [Middleware](../http/middleware.md)
- [Crypto Utilities](crypto.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
| 1.0.0 | 2026-03-10 | RAiWorks | Fixed code examples to match source (secret validation, signing method check, correct default expiry) |
