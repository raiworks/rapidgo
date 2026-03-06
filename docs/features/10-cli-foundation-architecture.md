# 🏗️ Architecture: CLI Foundation

> **Feature**: `10` — CLI Foundation
> **Discussion**: [`10-cli-foundation-discussion.md`](10-cli-foundation-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Overview

The CLI Foundation replaces the monolithic `cmd/main.go` with a Cobra-based command tree. A new `core/cli/` package provides the root command, a `serve` subcommand (migrating current server startup logic), and a `version` subcommand. The `cmd/main.go` becomes a thin shell calling `cli.Execute()`. This establishes the extensible CLI framework that all future commands (make:*, migrate, db:seed) will attach to.

## File Structure

```
core/cli/
├── root.go          # Root command, Execute(), app bootstrap helper
├── serve.go         # Serve command — boots app, starts HTTP server
└── version.go       # Version command — prints framework version

core/cli/
└── cli_test.go      # Tests for CLI commands

cmd/
└── main.go          # Refactored — thin shell calling cli.Execute()
```

### Files Created (3)
| File | Package | Lines (est.) |
|---|---|---|
| `core/cli/root.go` | `cli` | ~50 |
| `core/cli/serve.go` | `cli` | ~55 |
| `core/cli/version.go` | `cli` | ~20 |

### Files Modified (1)
| File | Change |
|---|---|
| `cmd/main.go` | Replace all application logic with single call to `cli.Execute()` |

---

## Component Design

### Root Command (`core/cli/root.go`)

**Responsibility**: Define the root `rgo` command, register all subcommands, provide `Execute()` entry point, and expose `NewApp()` bootstrap helper for subcommands.
**Package**: `cli`

```go
package cli

import (
	"fmt"
	"os"

	"github.com/RAiWorks/RGo/app/providers"
	"github.com/RAiWorks/RGo/core/app"
	"github.com/spf13/cobra"
)

// Version is the current framework version.
const Version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "rgo",
	Short: "RGo — A Go web framework with Laravel-style developer experience",
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)
}

// Execute runs the root command. Called from main().
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// NewApp creates and boots a fully configured RGo application.
// Used by commands that need the application lifecycle (serve, migrate, etc.).
func NewApp() *app.Application {
	application := app.New()
	application.Register(&providers.ConfigProvider{})
	application.Register(&providers.LoggerProvider{})
	application.Register(&providers.DatabaseProvider{})
	application.Register(&providers.MiddlewareProvider{})
	application.Register(&providers.RouterProvider{})
	application.Boot()
	return application
}
```

**Design notes**:
- `rootCmd` is package-level — Cobra convention for simple CLIs
- `init()` registers all subcommands — future features add their commands here
- `Execute()` is the single entry point called from `main()`
- `NewApp()` centralizes the provider chain — avoids duplicating boot logic in every command
- `Version` is a constant — simple, no ldflags injection needed at this stage

### Serve Command (`core/cli/serve.go`)

**Responsibility**: Boot the application and start the HTTP server. Migrates logic from current `cmd/main.go`.
**Package**: `cli`

```go
package cli

import (
	"fmt"
	"log/slog"

	"github.com/RAiWorks/RGo/core/config"
	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/core/router"
	"github.com/spf13/cobra"
)

var servePort string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		application := NewApp()

		port := config.Env("APP_PORT", "8080")
		if servePort != "" {
			port = servePort
		}

		appName := config.Env("APP_NAME", "RGo")
		appEnv := config.AppEnv()

		fmt.Println("=================================")
		fmt.Printf("  %s Framework\n", appName)
		fmt.Println("  github.com/RAiWorks/RGo")
		fmt.Println("=================================")
		fmt.Printf("  Environment: %s\n", appEnv)
		fmt.Printf("  Port: %s\n", port)
		fmt.Printf("  Debug: %v\n", config.IsDebug())
		fmt.Println("=================================")

		slog.Info("server starting",
			"app", appName,
			"port", port,
			"env", appEnv,
		)

		r := container.MustMake[*router.Router](application.Container, "router")
		if err := r.Run(":" + port); err != nil {
			slog.Error("server failed to start", "err", err)
		}
	},
}

func init() {
	serveCmd.Flags().StringVarP(&servePort, "port", "p", "", "port to listen on (overrides APP_PORT)")
}
```

**Design notes**:
- `--port` / `-p` flag overrides `APP_PORT` env var — flag takes precedence only if explicitly provided
- Uses empty string default for `servePort` — allows detecting "not set" vs "set to X"
- Banner and server startup are identical to current `main.go` — backward compatible
- `NewApp()` handles the full provider chain — no duplication
- `serveCmd.Flags()` in `init()` — Cobra convention for persistent flag binding

### Version Command (`core/cli/version.go`)

**Responsibility**: Print the framework version.
**Package**: `cli`

```go
package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the RGo framework version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("RGo Framework v%s\n", Version)
	},
}
```

**Design notes**:
- Reads `Version` constant from `root.go`
- No flags — simple output
- Future: Could include Go version, build date, git commit via ldflags

### Refactored main.go (`cmd/main.go`)

**Responsibility**: Thin entrypoint — delegates to CLI.

```go
package main

import "github.com/RAiWorks/RGo/core/cli"

func main() {
	cli.Execute()
}
```

**Design notes**:
- `main()` is intentionally minimal — all logic lives in `core/cli/`
- Existing `Makefile` targets (`make build`, `make run`) continue to work unchanged
- Binary is still built as `bin/rgo` — no changes needed

---

## Blueprint Adaptations

| Blueprint | Our Implementation | Reason |
|---|---|---|
| Commands defined in `package cmd` | Commands in `package cli` (`core/cli/`) | Go convention: `cmd/` is `package main` only. Library code belongs in named packages for testability. |
| `framework` as binary name | `rgo` as binary name | Our project is called RGo; `rgo` is shorter and consistent with existing Makefile |
| Full `make:*` scaffolding commands | Deferred to later features | Roadmap scopes #10 as "cobra setup, base commands" only. Scaffolding is a separate feature. |
| No `serve` command shown | Added `serve` command | Blueprint lists `framework serve` in the example commands list. Essential for CLI-based workflow. |

---

## Data Flow

```
User runs: rgo serve --port 9090
    │
    ▼
cmd/main.go → cli.Execute()
    │
    ▼
rootCmd.Execute() → Cobra matches "serve" → serveCmd.Run()
    │
    ▼
serveCmd → NewApp() → [Config → Logger → Database → Middleware → Router providers]
    │
    ▼
serveCmd → resolve port (flag > env > default "8080")
    │
    ▼
serveCmd → print banner → start HTTP server on :9090
```

---

## Impact on Existing Code

| Area | Impact |
|---|---|
| `cmd/main.go` | **Rewritten** — all logic moves to `core/cli/serve.go`. Becomes a 5-line file. |
| `Makefile` | **No change** — `go build -o bin/rgo ./cmd/...` and `go run ./cmd/...` still work. |
| `.env` | **No change** — `APP_PORT` is still read; `--port` flag is optional override. |
| Provider chain | **Moved** — same providers in same order, now in `cli.NewApp()` instead of `main()`. |
| Existing tests | **No change** — provider tests don't depend on `main()`. |
