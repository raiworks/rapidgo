package providers

import (
	"testing"

	"github.com/RAiWorks/RGo/core/app"
	"github.com/RAiWorks/RGo/core/config"
	"github.com/RAiWorks/RGo/core/container"
)

// TC-01: ConfigProvider implements Provider interface (compile-time check)
var _ container.Provider = (*ConfigProvider)(nil)

// TC-02: LoggerProvider implements Provider interface (compile-time check)
var _ container.Provider = (*LoggerProvider)(nil)

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
