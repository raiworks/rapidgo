# 🏗️ Architecture: CORS Handling

> **Feature**: `25` — CORS Handling
> **Status**: FINAL
> **Date**: 2026-03-06

---

## 1. Overview

Enhance the existing CORS middleware in `core/middleware/cors.go` to add `AllowCredentials`, `ExposeHeaders`, env-based origin configuration, and updated default headers including `X-CSRF-Token`.

## 2. File Structure

### New Files

None.

### Modified Files

| File | Change |
|------|--------|
| `core/middleware/cors.go` | Add fields to CORSConfig, update defaults, add new headers |

**Total**: 0 new files, 1 modified file

## 3. Dependencies

None new — uses `os` (stdlib), already available.

## 4. Detailed Design

### Updated `CORSConfig`

```go
type CORSConfig struct {
	AllowOrigins     []string // Default: from CORS_ALLOWED_ORIGINS env, or ["*"]
	AllowMethods     []string // Default: ["GET","POST","PUT","DELETE","PATCH","OPTIONS"]
	AllowHeaders     []string // Default: ["Origin","Content-Type","Accept","Authorization","X-Request-ID","X-CSRF-Token"]
	ExposeHeaders    []string // Default: ["Content-Length","X-Request-ID"]
	AllowCredentials bool     // Default: true
	MaxAge           int      // Default: 43200 (12 hours)
}
```

### Updated `defaultCORSConfig()`

```go
func defaultCORSConfig() CORSConfig {
	origins := []string{"*"}
	if env := os.Getenv("CORS_ALLOWED_ORIGINS"); env != "" {
		origins = strings.Split(env, ",")
	}

	return CORSConfig{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           43200,
	}
}
```

### Updated `CORS()` handler

Add two new headers in the handler:

```go
if cfg.AllowCredentials {
    c.Header("Access-Control-Allow-Credentials", "true")
}
if len(cfg.ExposeHeaders) > 0 {
    c.Header("Access-Control-Expose-Headers", strings.Join(cfg.ExposeHeaders, ", "))
}
```

## 5. Environment Configuration

```env
CORS_ALLOWED_ORIGINS=https://example.com,https://app.example.com
```

When not set, defaults to `*` (allow all origins).
