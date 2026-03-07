# 🏗️ Architecture: Service Mode

> **Feature**: `56` — Service Mode Architecture
> **Discussion**: [`56-service-mode-discussion.md`](56-service-mode-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-07

---

## Overview

Add a service mode system that lets the framework run in selective modes — web-only, API-only, WebSocket-only, or any combination. Controlled by `RAPIDGO_MODE` env var or `--mode` CLI flag. Default `all` preserves full backward compatibility. Three implementation phases: (A) mode infrastructure and parsing, (B) mode-aware provider loading and route registration, (C) multi-port serving with separate Gin engines per mode.

## Scope

This feature covers **Phases 1–3** from the discussion doc:

| Phase | What | This Feature |
|---|---|---|
| Phase 1 | Optional Providers | ✅ Phase A+B |
| Phase 2 | Service Mode Flag | ✅ Phase A+B |
| Phase 3 | Multi-Port Serving | ✅ Phase C |
| Phase 4 | Remove Global State | ❌ Deferred → Feature #57 |
| Phase 5 | Worker/Queue Mode | ❌ Deferred → Feature #42 |

---

## File Structure

```
core/service/
├── mode.go                      # NEW — Mode type, constants, parsing, validation
└── mode_test.go                 # NEW — Tests for mode parsing and operations

core/cli/
├── root.go                      # MODIFY — NewApp(mode) with conditional providers
└── serve.go                     # MODIFY — --mode flag, multi-port serving

core/server/
└── server.go                    # MODIFY — Add ListenAndServeMulti() for multi-port

app/providers/
├── router_provider.go           # MODIFY — Mode-aware route registration
└── middleware_provider.go       # MODIFY — Mode-aware middleware loading

routes/
└── ws.go                        # NEW — RegisterWS() placeholder

.env                             # MODIFY — Add RAPIDGO_MODE, WEB_PORT, API_PORT, WS_PORT
```

**Files changed**: 6 modified, 3 new (including test file)

---

## Data Model

No database changes. No migrations.

---

## Design — Phase A: Service Mode Infrastructure

### Mode Type (`core/service/mode.go`)

```go
package service

import (
	"fmt"
	"strings"
)

// Mode represents which services the application should run.
// Uses bitmask for easy combination.
type Mode uint8

const (
	ModeWeb    Mode = 1 << iota // Web SSR (templates, static files)
	ModeAPI                     // JSON API endpoints
	ModeWS                      // WebSocket service

	ModeAll = ModeWeb | ModeAPI | ModeWS // Monolith — all HTTP services
)

// modeNames maps string identifiers to Mode constants.
var modeNames = map[string]Mode{
	"web": ModeWeb,
	"api": ModeAPI,
	"ws":  ModeWS,
	"all": ModeAll,
}

// ParseMode parses a comma-separated mode string into a Mode bitmask.
// Valid inputs: "all", "web", "api", "ws", "api,ws", "web,api", etc.
// Returns error for empty or invalid mode strings.
func ParseMode(s string) (Mode, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return 0, fmt.Errorf("service mode cannot be empty")
	}

	var m Mode
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		flag, ok := modeNames[part]
		if !ok {
			return 0, fmt.Errorf("invalid service mode: %q (valid: all, web, api, ws)", part)
		}
		m |= flag
	}

	if m == 0 {
		return 0, fmt.Errorf("service mode cannot be empty")
	}
	return m, nil
}

// Has returns true if the mode includes the given flag.
func (m Mode) Has(flag Mode) bool {
	return m&flag != 0
}

// String returns a human-readable representation of the mode.
func (m Mode) String() string {
	if m == ModeAll {
		return "all"
	}
	var parts []string
	if m.Has(ModeWeb) {
		parts = append(parts, "web")
	}
	if m.Has(ModeAPI) {
		parts = append(parts, "api")
	}
	if m.Has(ModeWS) {
		parts = append(parts, "ws")
	}
	if len(parts) == 0 {
		return "none"
	}
	return strings.Join(parts, ",")
}

// Services returns the list of individual modes active in this bitmask.
func (m Mode) Services() []Mode {
	var s []Mode
	if m.Has(ModeWeb) {
		s = append(s, ModeWeb)
	}
	if m.Has(ModeAPI) {
		s = append(s, ModeAPI)
	}
	if m.Has(ModeWS) {
		s = append(s, ModeWS)
	}
	return s
}

// PortEnvKey returns the environment variable name for the port of a single-mode constant.
func (m Mode) PortEnvKey() string {
	switch m {
	case ModeWeb:
		return "WEB_PORT"
	case ModeAPI:
		return "API_PORT"
	case ModeWS:
		return "WS_PORT"
	default:
		return "APP_PORT"
	}
}
```

---

## Design — Phase B: Mode-Aware Bootstrap

### Modified `NewApp()` (`core/cli/root.go`)

```go
// NewApp creates and boots a RapidGo application configured for the given mode.
func NewApp(mode service.Mode) *app.App {
	application := app.New()

	// Always required
	application.Register(&providers.ConfigProvider{})
	application.Register(&providers.LoggerProvider{})

	// DB required for HTTP modes that may access data
	if mode.Has(service.ModeWeb) || mode.Has(service.ModeAPI) || mode.Has(service.ModeWS) {
		application.Register(&providers.DatabaseProvider{})
	}

	// Session only needed for web mode (cookie-based auth)
	if mode.Has(service.ModeWeb) {
		application.Register(&providers.SessionProvider{})
	}

	// Middleware and Router for any HTTP mode
	application.Register(&providers.MiddlewareProvider{Mode: mode})
	application.Register(&providers.RouterProvider{Mode: mode})

	application.Boot()
	return application
}
```

The existing `NewApp()` (no args) is replaced. All callers (serve, migrate, migrate_rollback, migrate_status, seed commands) are updated — non-serve commands pass `service.ModeAll`.

### Modified `RouterProvider` (`app/providers/router_provider.go`)

```go
// RouterProvider creates the router and registers route definitions.
type RouterProvider struct {
	Mode service.Mode
}

// Register creates a new Router and registers it as "router" in the container.
func (p *RouterProvider) Register(c *container.Container) {
	c.Instance("router", router.New())
}

// Boot sets up templates, static serving, and loads route definitions based on mode.
func (p *RouterProvider) Boot(c *container.Container) {
	r := container.MustMake[*router.Router](c, "router")

	// Template engine and static serving — only for web mode
	if p.Mode.Has(service.ModeWeb) {
		r.SetFuncMap(router.DefaultFuncMap())
		viewsDir := filepath.Join("resources", "views")
		if info, err := os.Stat(viewsDir); err == nil && info.IsDir() {
			r.LoadTemplates(viewsDir)
		}
		if info, err := os.Stat("resources/static"); err == nil && info.IsDir() {
			r.Static("/static", "./resources/static")
		}
		if info, err := os.Stat("storage/uploads"); err == nil && info.IsDir() {
			r.Static("/uploads", "./storage/uploads")
		}
	}

	// Route definitions — conditional on mode
	if p.Mode.Has(service.ModeWeb) {
		routes.RegisterWeb(r)
	}
	if p.Mode.Has(service.ModeAPI) {
		routes.RegisterAPI(r)
	}
	if p.Mode.Has(service.ModeWS) {
		routes.RegisterWS(r)
	}

	// Health check — available in any HTTP mode when DB is present
	if c.Has("db") {
		health.Routes(r, func() *gorm.DB {
			return container.MustMake[*gorm.DB](c, "db")
		})
	}
}
```

### Modified `MiddlewareProvider` (`app/providers/middleware_provider.go`)

```go
// MiddlewareProvider registers built-in middleware aliases and groups.
type MiddlewareProvider struct {
	Mode service.Mode
}

// Boot registers middleware aliases relevant to the current mode.
func (p *MiddlewareProvider) Boot(c *container.Container) {
	// Always register — universally useful
	middleware.RegisterAlias("recovery", middleware.Recovery())
	middleware.RegisterAlias("requestid", middleware.RequestID())
	middleware.RegisterAlias("cors", middleware.CORS())
	middleware.RegisterAlias("error_handler", middleware.ErrorHandler())
	middleware.RegisterAlias("ratelimit", middleware.RateLimitMiddleware())

	// Web-only middleware
	if p.Mode.Has(service.ModeWeb) {
		middleware.RegisterAlias("csrf", middleware.CSRFMiddleware())
	}

	// Auth — needed by both web (session) and api (JWT)
	middleware.RegisterAlias("auth", middleware.AuthMiddleware())

	middleware.RegisterGroup("global",
		middleware.Recovery(),
		middleware.RequestID(),
	)
}
```

### Modified Serve Command (`core/cli/serve.go`)

```go
var servePort string
var serveMode string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load .env so RAPIDGO_MODE can be read from .env file
		config.Load()

		// Resolve mode: CLI flag > env var > default "all"
		modeStr := config.Env("RAPIDGO_MODE", "all")
		if serveMode != "" {
			modeStr = serveMode
		}

		mode, err := service.ParseMode(modeStr)
		if err != nil {
			return fmt.Errorf("invalid service mode: %w", err)
		}

		application := NewApp(mode)

		// ... banner and logging (same as today, adding mode) ...

		// Delegate to single-port or multi-port based on active services
		services := mode.Services()
		if len(services) <= 1 || allSamePort(services) {
			// Single server — backward compatible
			return serveSingle(application, mode)
		}
		// Multi-port — one server per service
		return serveMulti(application, mode)
	},
}

func init() {
	serveCmd.Flags().StringVarP(&servePort, "port", "p", "", "port to listen on (overrides APP_PORT)")
	serveCmd.Flags().StringVarP(&serveMode, "mode", "m", "", "service mode: all, web, api, ws, or comma-separated (overrides RAPIDGO_MODE)")
}
```

### New WebSocket Routes (`routes/ws.go`)

```go
package routes

import "github.com/RAiWorks/RapidGo/core/router"

// RegisterWS defines WebSocket routes.
// Currently a placeholder — WebSocket routes will be added as the application grows.
func RegisterWS(r *router.Router) {
	// WebSocket route registration will be added here.
	// Example: r.Get("/ws", controllers.WebSocketHandler)
}
```

---

## Design — Phase C: Multi-Port Serving

### Multi-Server Support (`core/server/server.go`)

Add a new function alongside the existing `ListenAndServe`:

```go
// ServiceConfig identifies a named HTTP service to run on a specific port.
type ServiceConfig struct {
	Name   string       // "web", "api", "ws" — for logging
	Config Config       // Standard server config (addr, handler, timeouts)
}

// ListenAndServeMulti starts multiple HTTP servers on separate ports and
// blocks until SIGINT/SIGTERM. All servers are shut down gracefully.
func ListenAndServeMulti(services []ServiceConfig) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	servers := make([]*http.Server, len(services))
	errCh := make(chan error, len(services))

	for i, svc := range services {
		srv := &http.Server{
			Addr:         svc.Config.Addr,
			Handler:      svc.Config.Handler,
			ReadTimeout:  svc.Config.ReadTimeout,
			WriteTimeout: svc.Config.WriteTimeout,
			IdleTimeout:  svc.Config.IdleTimeout,
		}
		servers[i] = srv

		go func(name string) {
			slog.Info("service starting", "name", name, "addr", srv.Addr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errCh <- fmt.Errorf("service %s: %w", name, err)
			}
		}(svc.Name)
	}

	// Wait for signal or server error
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	slog.Info("shutting down all services…")
	stop()

	shutdownTimeout := services[0].Config.ShutdownTimeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown all servers
	var firstErr error
	for i, srv := range servers {
		if err := srv.Shutdown(shutdownCtx); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("shutdown %s: %w", services[i].Name, err)
		}
	}
	slog.Info("all services stopped")
	return firstErr
}
```

### Multi-Port Serve Logic (`core/cli/serve.go`)

```go
// serveSingle starts one HTTP server on a single port (backward compatible).
func serveSingle(application *app.App, mode service.Mode) error {
	port := resolvePort(mode)
	r := container.MustMake[*router.Router](application.Container, "router")
	return server.ListenAndServe(server.Config{
		Addr:            ":" + port,
		Handler:         r,
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 30 * time.Second,
	})
}

// serveMulti starts separate HTTP servers per service on different ports.
func serveMulti(application *app.App, mode service.Mode) error {
	var services []server.ServiceConfig

	for _, svc := range mode.Services() {
		r := router.New()
		applyRoutesForMode(r, application.Container, svc)
		port := resolvePortForMode(svc)

		services = append(services, server.ServiceConfig{
			Name: svc.String(),
			Config: server.Config{
				Addr:            ":" + port,
				Handler:         r,
				ReadTimeout:     15 * time.Second,
				WriteTimeout:    15 * time.Second,
				IdleTimeout:     60 * time.Second,
				ShutdownTimeout: 30 * time.Second,
			},
		})
	}

	return server.ListenAndServeMulti(services)
}

// resolvePort returns the port for the active mode.
// Single-mode uses mode-specific port env var, else APP_PORT.
func resolvePort(mode service.Mode) string {
	if servePort != "" {
		return servePort
	}
	services := mode.Services()
	if len(services) == 1 {
		return config.Env(services[0].PortEnvKey(), config.Env("APP_PORT", "8080"))
	}
	return config.Env("APP_PORT", "8080")
}

// resolvePortForMode returns the port for a specific service mode.
func resolvePortForMode(m service.Mode) string {
	return config.Env(m.PortEnvKey(), config.Env("APP_PORT", "8080"))
}

// applyRoutesForMode registers routes on a router for a specific mode.
func applyRoutesForMode(r *router.Router, c *container.Container, m service.Mode) {
	if m.Has(service.ModeWeb) {
		r.SetFuncMap(router.DefaultFuncMap())
		// ... template/static setup ...
		routes.RegisterWeb(r)
	}
	if m.Has(service.ModeAPI) {
		routes.RegisterAPI(r)
	}
	if m.Has(service.ModeWS) {
		routes.RegisterWS(r)
	}

	// Health check — each per-service router gets its own health endpoints
	if c.Has("db") {
		health.Routes(r, func() *gorm.DB {
			return container.MustMake[*gorm.DB](c, "db")
		})
	}
}

// allSamePort returns true if all services resolve to the same port.
func allSamePort(services []service.Mode) bool {
	if len(services) <= 1 {
		return true
	}
	port := config.Env(services[0].PortEnvKey(), config.Env("APP_PORT", "8080"))
	for _, s := range services[1:] {
		if config.Env(s.PortEnvKey(), config.Env("APP_PORT", "8080")) != port {
			return false
		}
	}
	return true
}
```

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `RAPIDGO_MODE` | `all` | Service mode: `all`, `web`, `api`, `ws`, or comma-separated |
| `WEB_PORT` | (falls back to `APP_PORT`) | Port for web SSR service |
| `API_PORT` | (falls back to `APP_PORT`) | Port for API service |
| `WS_PORT` | (falls back to `APP_PORT`) | Port for WebSocket service |

**Precedence**: CLI `--mode` flag > `RAPIDGO_MODE` env var > default `"all"`
**Precedence**: CLI `--port` flag > mode-specific port var > `APP_PORT` > `"8080"`

---

## Data Flow

### Default Mode (`all`) — Zero Behavioral Change

```
main → Execute → serve --mode=all
    → ParseMode("all") → ModeWeb|ModeAPI|ModeWS
    → NewApp(ModeAll)
        → Config, Logger, DB, Session, Middleware(All), Router(All)
        → Boot: templates + static + RegisterWeb + RegisterAPI + RegisterWS + health
    → serveSingle: one router on APP_PORT (:8080)
    → ListenAndServe — identical to current behavior
```

### Single Mode (`api`)

```
main → Execute → serve --mode=api
    → ParseMode("api") → ModeAPI
    → NewApp(ModeAPI)
        → Config, Logger, DB (no Session), Middleware(API), Router(API)
        → Boot: RegisterAPI only (no templates, no web routes)
    → serveSingle: one router on API_PORT (or APP_PORT)
    → ListenAndServe
```

### Multi-Port (`api,ws`)

```
main → Execute → serve --mode=api,ws
    → ParseMode("api,ws") → ModeAPI|ModeWS
    → NewApp(ModeAPI|ModeWS)
        → Config, Logger, DB, Middleware(API|WS), Router(API|WS)
        → Boot: RegisterAPI + RegisterWS on main router
    → Ports differ? (API_PORT=8081, WS_PORT=8082)
        → Yes → serveMulti: separate routers, separate servers
        → No  → serveSingle: one router, one port
```

---

## Trade-offs

| Decision | Pro | Con | Rationale |
|---|---|---|---|
| Bitmask mode (not string enum) | Easy combination, `Has()` checks | Slightly less obvious | Standard Go pattern, clean API |
| Shared container for multi-port | One DB pool, simple | Services can interfere | Independent containers deferred to #57 |
| Mode passed via provider struct field | No provider interface changes | Providers now have state | Minimal change, Mode is read-only |
| CSRF only in web mode | API doesn't need CSRF | Could be needed for API forms | APIs use token auth, not cookies |
| Main router orphaned in multi-port | Simple — no special-case logic in Boot() | RouterProvider.Boot() registers routes on a container router that serveMulti() doesn't use | Acceptable waste — route registration is fast; cleaner than adding multi-port awareness to providers |
| Session only in web mode | Lighter API mode | API can't use session-based auth by default | API should use JWT; if session needed, use mode=web,api |
| Runtime mode (no build tags) | Simpler, one binary | Larger binary | Build tags deferred as future optimization |

---

## Security Considerations

- **Mode validation**: Invalid mode strings fail fast with clear error (never start with zero services)
- **No mode bypass**: Excluded routes are never registered — not just hidden, but absent from the Gin engine
- **Session isolation**: API mode has no session provider — no session fixation attack surface
- **CSRF per mode**: CSRF middleware only loaded for web mode — API mode uses token auth

---

## Backward Compatibility

| Scenario | Behavior |
|---|---|
| `RAPIDGO_MODE` not set, no `--mode` flag | Default `"all"` — identical to current behavior |
| No `WEB_PORT`/`API_PORT`/`WS_PORT` set | Falls back to `APP_PORT` (then `8080`) |
| `NewApp()` callers (migrate, seed, scaffold) | Pass `ModeAll` — full provider set, same as before |
| Existing tests | All run with `ModeAll` by default — no changes needed |
