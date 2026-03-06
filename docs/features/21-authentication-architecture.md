# 🏗️ Architecture: Authentication

> **Feature**: `21` — Authentication
> **Status**: FINAL
> **Date**: 2026-03-06

---

## 1. Overview

Add JWT authentication support and auth middleware to the framework. This feature creates the `core/auth/` package for token operations and `core/middleware/auth.go` for route protection. Also adds a GORM `BeforeCreate` hook on the User model for automatic password hashing.

## 2. File Structure

### New Files

| File | Purpose |
|------|---------|
| `core/auth/jwt.go` | `GenerateToken`, `ValidateToken` — JWT helpers |
| `core/middleware/auth.go` | `AuthMiddleware` — Bearer token validation |
| `core/auth/auth_test.go` | Tests for JWT functions |

### Modified Files

| File | Change |
|------|--------|
| `database/models/user.go` | Add `BeforeCreate` hook for password hashing |
| `database/models/models_test.go` | Add test for BeforeCreate hook |
| `go.mod` / `go.sum` | Add `golang-jwt/jwt/v5` dependency |

**Total**: 3 new files, 3 modified files

## 3. Dependencies

### New External Dependency

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/golang-jwt/jwt/v5` | latest | JWT token creation and parsing |

### Internal Dependencies

- `app/helpers` — `HashPassword()` for User `BeforeCreate` hook
- `core/middleware/registry.go` — `RegisterAlias()` for `"auth"` alias (used at route setup, not in this feature's code)

## 4. Detailed Design

### 4.1 JWT Package — `core/auth/jwt.go`

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

// ValidateToken parses and validates a JWT string.
// Returns the claims if the token is valid.
func ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not set")
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

**Key decisions**:
- Validates signing method is HMAC (prevents `none` algorithm attack)
- Reads `JWT_EXPIRY` from env (default 3600s = 1 hour), not hardcoded 24h
- Adds `iat` (issued at) claim for token age tracking
- Returns explicit error when `JWT_SECRET` is not set

### 4.2 Auth Middleware — `core/middleware/auth.go`

```go
package middleware

import (
	"net/http"
	"strings"

	"github.com/RAiWorks/RGo/core/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT Bearer tokens on protected routes.
// On success, sets "user_id" in the Gin context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid Authorization header",
			})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := auth.ValidateToken(tokenStr)
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

### 4.3 User Model Hook — `database/models/user.go`

```go
package models

import (
	"strings"

	"github.com/RAiWorks/RGo/app/helpers"
	"gorm.io/gorm"
)

// User represents an application user.
type User struct {
	BaseModel
	Name     string `gorm:"size:100;not null" json:"name"`
	Email    string `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	Role     string `gorm:"size:50;default:user" json:"role"`
	Active   bool   `gorm:"default:true" json:"active"`
	Posts    []Post `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}

// BeforeCreate hashes the password before inserting into the database.
// Skips if the password is already bcrypt-hashed (starts with "$2a$" or "$2b$").
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" && !strings.HasPrefix(u.Password, "$2a$") && !strings.HasPrefix(u.Password, "$2b$") {
		hashed, err := helpers.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashed
	}
	return nil
}
```

## 5. Data Flow

### JWT Authentication Flow
```
Client → Authorization: Bearer <token>
  → AuthMiddleware
    → Extract token from header
    → auth.ValidateToken(token)
      → jwt.Parse with HMAC check
      → Return claims or error
    → Set c.Set("user_id", claims["user_id"])
    → c.Next() → Handler
```

### Session Authentication Flow (existing — no new code)
```
Client → Cookie: session_id=abc123
  → SessionMiddleware (existing, #20)
    → mgr.Start() → load session data
    → c.Set("session", data) where data["user_id"] exists
    → c.Next() → Handler checks data["user_id"]
```

## 6. Configuration

All env vars already exist in `.env`:

| Variable | Default | Purpose |
|----------|---------|---------|
| `JWT_SECRET` | `change-me-to-a-random-string` | HMAC signing key |
| `JWT_EXPIRY` | `3600` | Token lifetime in seconds |

No new env vars needed.

## 7. Security

1. **Signing method validation**: `ValidateToken` checks `t.Method.(*jwt.SigningMethodHMAC)` — prevents `none` algorithm attack
2. **Empty secret protection**: Both functions return error if `JWT_SECRET` is empty
3. **Password never exposed**: User model has `json:"-"` on Password field
4. **Bcrypt double-hash prevention**: `BeforeCreate` checks for `$2a$`/`$2b$` prefix
5. **Constant-time password comparison**: Uses bcrypt's built-in comparison (via helpers)

## 8. What This Feature Does NOT Include

| Item | Reason |
|------|--------|
| Refresh tokens | Blueprint mentions but provides no code — future enhancement |
| Login/register controllers | App-level concern, not framework infrastructure |
| Auth routes | App-level — framework provides the building blocks |
| `AuthProvider` service provider | Not in blueprint — auth functions are stateless, no container registration needed |
| Session-based auth code | Already works via #20's session system — no new code needed |
| `"admin"` / `"verified"` aliases | User-defined middleware, not framework-provided |
