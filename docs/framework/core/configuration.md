---
title: "Configuration"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Configuration

## Abstract

This document covers the framework's configuration system — loading
environment variables from `.env` files, the `Env()` helper with
fallback defaults, and environment detection helpers that control
behavior across development, production, and testing environments.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Overview](#2-overview)
3. [.env File Structure](#3-env-file-structure)
4. [Loading Configuration](#4-loading-configuration)
5. [Env Helper](#5-env-helper)
6. [Environment Detection](#6-environment-detection)
7. [Conditional Behavior](#7-conditional-behavior)
8. [Viper vs godotenv](#8-viper-vs-godotenv)
9. [Per-Environment Recommendations](#9-per-environment-recommendations)
10. [Security Considerations](#10-security-considerations)
11. [References](#11-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

## 2. Overview

Configuration is managed through environment variables, loaded from a
`.env` file at application startup. The framework uses
**godotenv** (`github.com/joho/godotenv`) for `.env` loading and
provides helper functions for type-safe access with fallback defaults.

**Viper** (`github.com/spf13/viper`) is the recommended alternative
when YAML, TOML, or JSON configuration files are needed.

## 3. .env File Structure

The `.env` file contains all runtime configuration as key-value pairs:

```env
# Application
APP_NAME=MyApp
APP_PORT=8080
APP_ENV=development
APP_DEBUG=true
APP_KEY=base64:32-byte-random-key-here

# Database
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=secret
DB_NAME=myapp

# JWT
JWT_SECRET=change-me-in-production

# Session
SESSION_DRIVER=db
SESSION_LIFETIME=120
SESSION_COOKIE=framework_session
SESSION_PATH=/
SESSION_DOMAIN=
SESSION_SECURE=false
SESSION_HTTPONLY=true
SESSION_SAMESITE=lax

# Redis (if SESSION_DRIVER=redis or CACHE_DRIVER=redis)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Session File (if SESSION_DRIVER=file)
SESSION_FILE_PATH=storage/sessions

# Cache
CACHE_DRIVER=redis
CACHE_PREFIX=app:
CACHE_TTL=3600

# Mail
MAIL_HOST=smtp.example.com
MAIL_PORT=587
MAIL_USERNAME=noreply@example.com
MAIL_PASSWORD=secret
MAIL_FROM_NAME=MyApp
MAIL_FROM_ADDRESS=noreply@example.com
MAIL_ENCRYPTION=tls

# Storage
STORAGE_DRIVER=local
STORAGE_LOCAL_PATH=storage/uploads
AWS_REGION=us-east-1
AWS_BUCKET=my-bucket
AWS_ACCESS_KEY_ID=xxx
AWS_SECRET_ACCESS_KEY=xxx

# CORS
CORS_ALLOWED_ORIGINS=https://example.com,https://app.example.com

# Rate Limiting
RATE_LIMIT=60-M
RATE_LIMIT_AUTH=5-M

# Caddy (optional)
CADDY_ENABLED=false
CADDY_DOMAIN=localhost
```

## 4. Loading Configuration

The `config.Load()` function **MUST** be called first in `main()`,
before any other initialization:

```go
package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

func Load() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment")
    }
}
```

If no `.env` file is found, the framework falls back to system
environment variables. This allows production deployments to inject
configuration via container orchestrators (Docker, Kubernetes) without
a `.env` file.

## 5. Env Helper

The `Env()` function provides type-safe access to environment variables
with a fallback default:

```go
// Env returns an environment variable with a fallback default.
func Env(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}
```

**Usage:**

```go
port := config.Env("APP_PORT", "8080")
appName := config.Env("APP_NAME", "MyApp")
dbDriver := config.Env("DB_DRIVER", "postgres")
```

The fallback ensures the application always has sensible defaults
even when variables are not set.

## 6. Environment Detection

The framework provides helper functions to detect the current
environment:

```go
package config

import "os"

func AppEnv() string     { return Env("APP_ENV", "development") }
func IsProduction() bool { return AppEnv() == "production" }
func IsDevelopment() bool { return AppEnv() == "development" }
func IsTesting() bool    { return AppEnv() == "testing" }
func IsDebug() bool      { return os.Getenv("APP_DEBUG") == "true" }
```

### Supported environments

| `APP_ENV` | Environment | `APP_DEBUG` |
|-----------|------------|-------------|
| `development` | Local development | `true` (default) |
| `production` | Production deployment | `false` |
| `testing` | Test suite execution | `true` |

## 7. Conditional Behavior

Environment detection enables conditional behavior throughout the
framework:

### Error handling

```go
if config.IsDebug() {
    r.Use(gin.Recovery()) // detailed stack traces
} else {
    r.Use(middleware.ErrorHandler()) // generic error messages
}
```

### Logging

```go
if config.IsProduction() {
    // JSON logs to file — structured, parseable
} else {
    // Pretty-print to stdout — readable for developers
}
```

### Middleware

```go
if !config.IsTesting() {
    r.Use(middleware.RateLimitMiddleware()) // skip rate limiting in tests
}
```

### Session/Cache defaults

- **Development:** `SESSION_DRIVER=memory`, `CACHE_DRIVER=memory`
- **Production:** `SESSION_DRIVER=redis`, `CACHE_DRIVER=redis`

## 8. Viper vs godotenv

| Feature | godotenv | Viper |
|---------|----------|-------|
| `.env` files | Yes | Yes |
| YAML/TOML/JSON | No | Yes |
| Watch for changes | No | Yes |
| Remote config | No | Yes (etcd, Consul) |
| Complexity | Minimal | Higher |

**RECOMMENDED:** Use **godotenv** for `.env`-only projects (most web
apps). Use **Viper** when you need YAML configuration files, config
watching, or remote config stores.

## 9. Per-Environment Recommendations

### Development

```env
APP_ENV=development
APP_DEBUG=true
DB_DRIVER=sqlite
SESSION_DRIVER=memory
CACHE_DRIVER=memory
CORS_ALLOWED_ORIGINS=*
```

### Production

```env
APP_ENV=production
APP_DEBUG=false
DB_DRIVER=postgres
SESSION_DRIVER=redis
CACHE_DRIVER=redis
SESSION_SECURE=true
CORS_ALLOWED_ORIGINS=https://yourdomain.com
```

### Testing

```env
APP_ENV=testing
APP_DEBUG=true
DB_DRIVER=sqlite
SESSION_DRIVER=memory
CACHE_DRIVER=memory
```

## 10. Security Considerations

- The `.env` file **MUST NOT** be committed to version control. Add it
  to `.gitignore`.
- `APP_KEY` **MUST** be exactly 32 bytes (for AES-256-GCM encryption).
  Generate with: `openssl rand -base64 32`.
- `JWT_SECRET` **MUST** be a strong, unique random string in production.
- `APP_DEBUG` **MUST** be `false` in production to prevent stack trace
  exposure.
- `SESSION_SECURE` **MUST** be `true` in production (requires HTTPS).
- Database credentials (`DB_PASS`) and mail credentials
  (`MAIL_PASSWORD`) **MUST** be unique per environment.
- `CORS_ALLOWED_ORIGINS` **MUST NOT** be `*` in production.

## 11. References

- [Application Lifecycle](../architecture/application-lifecycle.md)
- [Environment Variables Reference](../reference/env-reference.md)
- [Service Providers](service-providers.md)
- [Design Principles](../architecture/design-principles.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
