# 🏗️ Architecture: Service Providers

> **Feature**: `06` — Service Providers
> **Discussion**: [`06-service-providers-discussion.md`](06-service-providers-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Overview

Two concrete service providers (`ConfigProvider`, `LoggerProvider`) bootstrap the existing framework features through the Provider lifecycle. `cmd/main.go` is updated from direct function calls to the `App` bootstrap pattern. This establishes the provider pattern that all future features will follow.

## File Structure

```
app/providers/
├── config_provider.go      # ConfigProvider — loads .env, registers config accessor
└── logger_provider.go      # LoggerProvider — sets up slog in Boot phase

app/providers/
└── providers_test.go       # Tests for both providers

cmd/
└── main.go                 # MODIFIED — App bootstrap replaces direct calls
```

## Component Design

### ConfigProvider (`app/providers/config_provider.go`)

**Responsibility**: Load environment configuration and register it in the container
**Package**: `providers`

```go
package providers

import (
	"github.com/RAiWorks/RGo/core/config"
	"github.com/RAiWorks/RGo/core/container"
)

// ConfigProvider loads environment configuration and registers
// config accessors in the container.
type ConfigProvider struct{}

// Register loads the .env file. Called first — all other providers
// may read config values during their own Register().
func (p *ConfigProvider) Register(c *container.Container) {
	config.Load()
}

// Boot is a no-op. Config is fully available after Load().
func (p *ConfigProvider) Boot(c *container.Container) {}
```

**Design notes**:
- `config.Load()` in `Register()` because other providers may read `config.Env()` during their own `Register()` calls
- No service registered in container — the `config` package uses package-level functions (`config.Env()`, `config.IsDebug()`) which are already globally accessible
- `Boot()` is empty — config doesn't need post-registration setup

### LoggerProvider (`app/providers/logger_provider.go`)

**Responsibility**: Initialize structured logging after config is loaded
**Package**: `providers`

```go
package providers

import (
	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/core/logger"
)

// LoggerProvider sets up structured logging via slog.
type LoggerProvider struct{}

// Register is a no-op. Logger setup requires config values
// which may not be loaded if this isn't the first provider.
func (p *LoggerProvider) Register(c *container.Container) {}

// Boot initializes the logger. Runs after all providers have
// registered, so config values (LOG_LEVEL, LOG_FORMAT, LOG_OUTPUT)
// are guaranteed available.
func (p *LoggerProvider) Boot(c *container.Container) {
	logger.Setup()
}
```

**Design notes**:
- `logger.Setup()` in `Boot()` because it reads `LOG_LEVEL`, `LOG_FORMAT`, `LOG_OUTPUT` from env vars — config must be loaded first
- `Register()` is empty — nothing to bind; `slog` is a global logger
- `Boot()` order: LoggerProvider boots after ConfigProvider (registration order preserved)

### Updated main.go (`cmd/main.go`)

**Responsibility**: Application entrypoint using App bootstrap pattern
**Package**: `main`

```go
package main

import (
	"fmt"
	"log/slog"

	"github.com/RAiWorks/RGo/app/providers"
	"github.com/RAiWorks/RGo/core/app"
	"github.com/RAiWorks/RGo/core/config"
)

func main() {
	application := app.New()

	// Register providers (order matters)
	application.Register(&providers.ConfigProvider{})  // 1. Config first — loads .env
	application.Register(&providers.LoggerProvider{})  // 2. Logger — uses config in Boot

	// Boot all providers
	application.Boot()

	// Banner
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

	slog.Info("server initialized",
		"app", appName,
		"port", appPort,
		"env", appEnv,
	)
}
```

**Changes from current main.go**:
- Removed direct `config.Load()` call — now handled by `ConfigProvider.Register()`
- Removed direct `logger.Setup()` call — now handled by `LoggerProvider.Boot()`
- Added `app.New()` → `Register()` → `Boot()` bootstrap pattern
- Added imports for `app` and `providers` packages
- Banner and log output remain identical

## Data Flow

```
main.go
  → app.New()                              creates Container
  → app.Register(&ConfigProvider{})        → config.Load() — .env loaded
  → app.Register(&LoggerProvider{})        → (no-op in Register)
  → app.Boot()
      → ConfigProvider.Boot()              → (no-op)
      → LoggerProvider.Boot()              → logger.Setup() — slog configured
  → Banner + slog.Info(...)                application running
```

## Configuration

No new environment variables. Uses existing:
- `APP_NAME`, `APP_ENV`, `APP_PORT`, `APP_DEBUG` (from #02)
- `LOG_LEVEL`, `LOG_FORMAT`, `LOG_OUTPUT` (from #03)

## Security Considerations

- Provider factories that access credentials should read from env vars inside closures, not store as struct fields (pattern established for future providers)
- `config.Load()` reads `.env` — file must not contain production secrets in version control (documented in `.env` header)

## Trade-offs & Alternatives

| Approach | Pros | Cons | Verdict |
|---|---|---|---|
| Providers for Config/Logger | Follows blueprint pattern, consistent bootstrap, extensible | Slight indirection vs direct calls | ✅ Selected |
| Keep direct calls in main.go | Simpler, no abstraction | Doesn't establish provider pattern, inconsistent with blueprint | ❌ Deferred |
| Register config/logger as container services | Full DI for all features | Config uses package-level funcs, logger is global slog — wrapping adds no value | ❌ Over-engineered |

## Next

Create tasks doc → `06-service-providers-tasks.md`
