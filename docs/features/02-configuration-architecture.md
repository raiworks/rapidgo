# 🏗️ Architecture: Configuration System

> **Feature**: `02` — Configuration System
> **Discussion**: [`02-configuration-discussion.md`](02-configuration-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-05

---

## Overview

Create the `core/config` package that loads `.env` files via godotenv, provides typed accessor helpers for environment variables, and exposes environment detection functions. This is the framework's first real package — every subsequent feature depends on it for configuration access.

## File Structure

```
core/config/
├── config.go           # Load() — .env loading via godotenv
├── env.go              # Env(), EnvInt(), EnvBool() — typed accessors
├── environment.go      # AppEnv(), IsProduction(), IsDevelopment(), IsTesting(), IsDebug()
└── config_test.go      # Unit tests for all exported functions

cmd/
└── main.go             # MODIFY — add config.Load() call, display app name/port
```

## Data Model

N/A — no database entities. Configuration is read from environment variables at runtime.

## Component Design

### `core/config/config.go`

**Responsibility**: Load `.env` file into the process environment.
**Package**: `config`

```go
package config

import (
	"log"

	"github.com/joho/godotenv"
)

// Load reads the .env file and sets environment variables.
// If no .env file is found, it logs a message and continues
// (system environment variables are still available).
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}
}
```

### `core/config/env.go`

**Responsibility**: Typed access to environment variables with fallback defaults.
**Package**: `config`

```go
package config

import (
	"os"
	"strconv"
)

// Env returns the value of an environment variable, or the fallback if empty/unset.
func Env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// EnvInt returns the value of an environment variable as an int,
// or the fallback if empty/unset or not a valid integer.
func EnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}

// EnvBool returns the value of an environment variable as a bool.
// Only "true" and "1" are considered truthy. Everything else returns the fallback.
func EnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v == "true" || v == "1"
}
```

### `core/config/environment.go`

**Responsibility**: Environment detection helpers used throughout the framework.
**Package**: `config`

```go
package config

// AppEnv returns the current application environment (development, production, testing).
func AppEnv() string {
	return Env("APP_ENV", "development")
}

// IsProduction returns true if APP_ENV is "production".
func IsProduction() bool {
	return AppEnv() == "production"
}

// IsDevelopment returns true if APP_ENV is "development".
func IsDevelopment() bool {
	return AppEnv() == "development"
}

// IsTesting returns true if APP_ENV is "testing".
func IsTesting() bool {
	return AppEnv() == "testing"
}

// IsDebug returns true if APP_DEBUG is "true" or "1".
func IsDebug() bool {
	return EnvBool("APP_DEBUG", false)
}
```

### `cmd/main.go` (MODIFY)

**Changes**: Add `config.Load()` as the first call, use config to display app name and port.

```go
package main

import (
	"fmt"

	"github.com/RAiWorks/RGo/core/config"
)

func main() {
	config.Load()

	appName := config.Env("APP_NAME", "RGo")
	appPort := config.Env("APP_PORT", "8080")
	appEnv := config.AppEnv()

	fmt.Println("=================================")
	fmt.Printf("  %s Framework\n", appName)
	fmt.Println("  github.com/RAiWorks/RGo")
	fmt.Println("=================================")
	fmt.Printf("  Environment: %s\n", appEnv)
	fmt.Printf("  Port: %s\n", appPort)
	fmt.Printf("  Debug: %v\n", config.IsDebug())
	fmt.Println("=================================")
}
```

## Data Flow

```
Application Start
    │
    ▼
config.Load()
    │
    ├── godotenv.Load() reads .env file
    │   └── Calls os.Setenv() for each key=value pair
    │
    ├── If .env missing → log message, continue
    │
    ▼
Environment variables available via os.Getenv()
    │
    ├── config.Env("KEY", "default")      → string
    ├── config.EnvInt("KEY", 0)           → int
    ├── config.EnvBool("KEY", false)      → bool
    ├── config.AppEnv()                   → "development" | "production" | "testing"
    ├── config.IsProduction()             → bool
    ├── config.IsDevelopment()            → bool
    ├── config.IsTesting()                → bool
    └── config.IsDebug()                  → bool
```

## Configuration

This feature doesn't add new `.env` keys — it reads the keys already defined in Feature #01's `.env` file. The key subset used directly by this package:

| Key | Type | Default | Used By |
|---|---|---|---|
| `APP_ENV` | string | `development` | `AppEnv()`, `IsProduction()`, `IsDevelopment()`, `IsTesting()` |
| `APP_DEBUG` | bool | `false` | `IsDebug()` |
| `APP_NAME` | string | `RGo` | `main.go` banner |
| `APP_PORT` | string | `8080` | `main.go` banner |

## Security Considerations

- **No secrets in code** — all sensitive values come from `.env` or system environment
- **`.env.local` is gitignored** — real credentials never reach version control
- **`IsDebug()` controls error exposure** — stack traces only shown when `APP_DEBUG=true`
- **`IsProduction()` enforces stricter defaults** — used by middleware/session/CORS in future features

## Trade-offs & Alternatives

| Approach | Pros | Cons | Verdict |
|---|---|---|---|
| godotenv only | Minimal, zero transitive deps, `.env` is industry standard | No YAML/TOML, no config watching | ✅ Selected — sufficient for framework needs |
| Viper only | YAML/TOML/JSON, config watching, remote config | Heavier dependency, more complexity | ❌ Over-engineered for `.env`-only |
| godotenv + Viper | Best of both worlds | Two config libraries, potential conflicts | ❌ Add Viper later if needed |
| Custom `.env` parser | No dependencies | Re-inventing the wheel, edge cases | ❌ Unnecessary |

## Next

Create tasks doc → `02-configuration-tasks.md`
