package providers

import (
	"testing"

	"github.com/RAiWorks/RGo/core/app"
	"github.com/RAiWorks/RGo/core/config"
	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/core/middleware"
	"github.com/RAiWorks/RGo/core/router"
)

// TC-01: ConfigProvider implements Provider interface (compile-time check)
var _ container.Provider = (*ConfigProvider)(nil)

// TC-02: LoggerProvider implements Provider interface (compile-time check)
var _ container.Provider = (*LoggerProvider)(nil)

// TC-09: RouterProvider implements Provider interface (compile-time check)
var _ container.Provider = (*RouterProvider)(nil)

// TC-16: MiddlewareProvider implements Provider interface (compile-time check)
var _ container.Provider = (*MiddlewareProvider)(nil)

// TC-10: DatabaseProvider implements Provider interface (compile-time check)
var _ container.Provider = (*DatabaseProvider)(nil)

// TC-03: ConfigProvider.Register loads config
func TestConfigProvider_RegisterLoadsConfig(t *testing.T) {
	t.Setenv("APP_NAME", "TestApp")

	c := container.New()
	p := &ConfigProvider{}

	p.Register(c)

	appName := config.Env("APP_NAME", "")
	if appName != "TestApp" {
		t.Fatalf("expected 'TestApp', got '%s'", appName)
	}
}

// TC-04: LoggerProvider.Boot sets up logger
func TestLoggerProvider_BootSetsUpLogger(t *testing.T) {
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_OUTPUT", "stdout")

	c := container.New()

	cp := &ConfigProvider{}
	cp.Register(c)

	lp := &LoggerProvider{}
	// Should not panic
	lp.Boot(c)
}

// TC-05: Full App bootstrap with both providers
func TestFullAppBootstrap(t *testing.T) {
	t.Setenv("APP_NAME", "BootstrapTest")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_OUTPUT", "stdout")

	a := app.New()

	a.Register(&ConfigProvider{})
	a.Register(&LoggerProvider{})
	a.Boot()

	appName := config.Env("APP_NAME", "")
	if appName != "BootstrapTest" {
		t.Fatalf("expected 'BootstrapTest', got '%s'", appName)
	}
}

// TC-06: ConfigProvider.Boot is no-op
func TestConfigProvider_BootIsNoOp(t *testing.T) {
	c := container.New()
	p := &ConfigProvider{}
	// Should not panic
	p.Boot(c)
}

// TC-07: LoggerProvider.Register is no-op
func TestLoggerProvider_RegisterIsNoOp(t *testing.T) {
	c := container.New()
	p := &LoggerProvider{}
	// Should not panic
	p.Register(c)
}

// TC-08: Provider registration order — Config before Logger
func TestProviderOrder_ConfigBeforeLogger(t *testing.T) {
	t.Setenv("APP_NAME", "OrderTest")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_OUTPUT", "stdout")

	a := app.New()

	a.Register(&ConfigProvider{})
	a.Register(&LoggerProvider{})
	a.Boot()

	logLevel := config.Env("LOG_LEVEL", "")
	if logLevel != "debug" {
		t.Fatalf("expected 'debug', got '%s'", logLevel)
	}
}

// TC-09: RouterProvider registers router as "router" in container
func TestRouterProvider_RegistersRouter(t *testing.T) {
	t.Setenv("APP_ENV", "testing")

	c := container.New()
	p := &RouterProvider{}
	p.Register(c)

	r := container.MustMake[*router.Router](c, "router")
	if r == nil {
		t.Fatal("expected non-nil *Router from container")
	}
	if r.Engine() == nil {
		t.Fatal("expected non-nil Engine")
	}
}

// TC-10: Full App bootstrap with all three providers
func TestFullAppBootstrap_WithRouter(t *testing.T) {
	t.Setenv("APP_NAME", "RouterBootstrapTest")
	t.Setenv("APP_ENV", "testing")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_OUTPUT", "stdout")

	a := app.New()
	a.Register(&ConfigProvider{})
	a.Register(&LoggerProvider{})
	a.Register(&RouterProvider{})
	a.Boot()

	r := container.MustMake[*router.Router](a.Container, "router")
	if r == nil {
		t.Fatal("expected non-nil *Router after full bootstrap")
	}
}

// TC-11: DatabaseProvider.Register registers "db" binding
func TestDatabaseProvider_RegistersBinding(t *testing.T) {
	c := container.New()
	p := &DatabaseProvider{}
	p.Register(c)

	if !c.Has("db") {
		t.Fatal("expected 'db' binding to be registered")
	}
}

// TC-12: Full bootstrap with DatabaseProvider (SQLite integration)
func TestDatabaseProvider_FullBootstrap(t *testing.T) {
	t.Setenv("DB_DRIVER", "sqlite")
	t.Setenv("DB_NAME", ":memory:")
	t.Setenv("DB_MAX_OPEN_CONNS", "1")
	t.Setenv("DB_MAX_IDLE_CONNS", "1")
	t.Setenv("DB_CONN_MAX_LIFETIME", "1")
	t.Setenv("DB_CONN_MAX_IDLE_TIME", "1")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_OUTPUT", "stdout")

	a := app.New()
	a.Register(&ConfigProvider{})
	a.Register(&LoggerProvider{})
	a.Register(&DatabaseProvider{})
	a.Boot()

	db := a.Container.Make("db")
	if db == nil {
		t.Fatal("expected non-nil db after full bootstrap")
	}
}

// TC-17: MiddlewareProvider Boot registers built-in aliases
func TestMiddlewareProvider_RegistersAliases(t *testing.T) {
	middleware.ResetRegistry()

	c := container.New()
	p := &MiddlewareProvider{}
	p.Boot(c)

	// Verify all 4 built-in aliases resolve without panic
	if middleware.Resolve("recovery") == nil {
		t.Fatal("recovery alias not registered")
	}
	if middleware.Resolve("requestid") == nil {
		t.Fatal("requestid alias not registered")
	}
	if middleware.Resolve("cors") == nil {
		t.Fatal("cors alias not registered")
	}
	if middleware.Resolve("error_handler") == nil {
		t.Fatal("error_handler alias not registered")
	}

	// Verify default group
	global := middleware.ResolveGroup("global")
	if len(global) != 2 {
		t.Fatalf("expected global group length 2, got %d", len(global))
	}
}
