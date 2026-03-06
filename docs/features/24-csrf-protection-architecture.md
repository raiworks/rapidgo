# 🏗️ Architecture: CSRF Protection

> **Feature**: `24` — CSRF Protection
> **Status**: FINAL
> **Date**: 2026-03-06

---

## 1. Overview

Implement CSRF protection middleware in `core/middleware/csrf.go`. Generates a per-session token on first request, stores it in the session, and validates it on state-changing requests. Registers as the `"csrf"` middleware alias.

## 2. File Structure

### New Files

| File | Purpose |
|------|---------|
| `core/middleware/csrf.go` | CSRFMiddleware function |

### Modified Files

| File | Change |
|------|--------|
| `app/providers/middleware_provider.go` | Register `"csrf"` alias |

**Total**: 1 new file, 1 modified file

## 3. Dependencies

None new — uses `crypto/rand`, `encoding/hex`, `net/http` (stdlib) and `gin-gonic/gin` (already in go.mod).

## 4. Detailed Design

### `core/middleware/csrf.go`

```go
package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CSRFMiddleware generates and validates per-session CSRF tokens.
// Safe methods (GET, HEAD, OPTIONS) are skipped.
// State-changing methods require a valid token in the _csrf_token form
// field or X-CSRF-Token header.
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, _ := c.Get("session")
		data := sess.(map[string]interface{})

		// Generate token if not present
		token, ok := data["_csrf_token"].(string)
		if !ok || token == "" {
			b := make([]byte, 32)
			rand.Read(b)
			token = hex.EncodeToString(b)
			data["_csrf_token"] = token
			c.Set("session", data)
		}

		// Make token available to templates
		c.Set("csrf_token", token)

		// Skip safe methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Validate token on POST/PUT/PATCH/DELETE
		submitted := c.PostForm("_csrf_token")
		if submitted == "" {
			submitted = c.GetHeader("X-CSRF-Token")
		}
		if submitted != token {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "CSRF token mismatch",
			})
			return
		}

		c.Next()
	}
}
```

### Middleware Provider Change

Add to `Boot()` in `app/providers/middleware_provider.go`:

```go
middleware.RegisterAlias("csrf", middleware.CSRFMiddleware())
```

## 5. Request Flow

```
Request → SessionMiddleware (loads session) → CSRFMiddleware
  ├─ GET/HEAD/OPTIONS → set csrf_token in context → Next()
  └─ POST/PUT/PATCH/DELETE
       ├─ Token matches → Next()
       └─ Token mismatch → 403 Forbidden (abort)
```

**Important**: CSRF middleware MUST run after SessionMiddleware since it reads from `c.Get("session")`.

## 6. Template Integration

```html
<form method="POST" action="/users">
    <input type="hidden" name="_csrf_token" value="{{.csrf_token}}">
    <!-- form fields -->
    <button type="submit">Create</button>
</form>
```

For AJAX, set the `X-CSRF-Token` header with the token value.
