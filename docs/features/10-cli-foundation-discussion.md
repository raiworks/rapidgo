# 💬 Discussion: CLI Foundation

> **Feature**: `10` — CLI Foundation
> **Status**: 🟢 COMPLETE
> **Branch**: `docs/10-cli-foundation`
> **Depends On**: #01 (Project Setup ✅), #02 (Configuration System ✅)
> **Date Started**: 2026-03-06
> **Date Completed**: 2026-03-06

---

## Summary

Replace the current `cmd/main.go` direct-execution entrypoint with a Cobra-based CLI framework. This introduces a root command (`rgo`), a `serve` subcommand that starts the HTTP server, and a `version` subcommand for build info. The CLI foundation provides the extensible command tree that all future CLI features (make:*, migrate, db:seed) will attach to.

---

## Functional Requirements

- As a **framework developer**, I want a Cobra root command so that all future subcommands (make:*, migrate, serve) attach to a single command tree
- As a **framework developer**, I want a `serve` command that boots the application and starts the HTTP server so that server startup moves from `main()` to a structured CLI subcommand
- As a **framework developer**, I want a `version` command that prints the framework version so that users can verify their installation
- As a **framework user**, I want to run `rgo serve` to start my application so that the CLI interface is consistent and discoverable
- As a **framework user**, I want to run `rgo serve --port 9090` to override the default port so that I can choose a port without editing `.env`
- As a **framework user**, I want to run `rgo version` to see the framework version
- As a **framework user**, I want `rgo` (no args) to display help text listing all available commands

## Current State / Reference

### What Exists
- **`cmd/main.go`**: Direct execution — `func main()` creates the application, registers providers, boots, prints banner, and calls `r.Run()`. No CLI framework.
- **`bin/rgo.exe`**: Built via `go build -o bin/rgo ./cmd/...` — runs main() directly
- **`Makefile`**: `make build`, `make run` targets already work with current `cmd/main.go`
- **Provider order**: Config(1) → Logger(2) → Database(3) → Middleware(4) → Router(5)

### Blueprint Reference
The blueprint (CLI Tools section, lines 480–665) shows:
1. Cobra (`github.com/spf13/cobra`) as the recommended library
2. Example commands: `framework serve`, `framework make:*`, `framework migrate`, `framework db:seed`
3. Code generation scaffolding with templates

The roadmap scopes Feature #10 as **"cobra setup, base commands"** — the foundation only:
- Root command setup
- `serve` command (replaces current main.go logic)
- Command tree ready for future features to add subcommands

### What Works Well
- Application lifecycle (providers, boot) is solid
- Banner printing and server startup logic is correct
- Build pipeline (`make build`, `make run`) works

### What Needs Improvement
- `main()` is monolithic — server startup is hardcoded, no room for alternative commands
- No CLI framework — can't add `migrate`, `make:controller`, or other commands
- No `--port` flag override — port is only readable from `.env`

## Proposed Approach

### Cobra Root Command (`core/cli/root.go`)

Create a new `core/cli/` package with a root command. This is the entry point for all CLI operations. The root command displays help when run without arguments.

### Serve Command (`core/cli/serve.go`)

Move the current server startup logic from `main()` into a `serve` command. The serve command:
1. Creates and boots the application (same provider chain)
2. Prints the startup banner
3. Starts the Gin HTTP server
4. Accepts `--port` flag to override `APP_PORT` from `.env`

### Version Command (`core/cli/version.go`)

A simple command that prints the framework version. The version is defined as a package-level constant.

### Refactored main.go

`cmd/main.go` becomes a thin shell that calls `cli.Execute()`. All application logic moves to the CLI commands.

### Why `core/cli/` not `cmd/`?

The `cmd/` directory is for the binary entrypoint only (`package main`). CLI command definitions belong in `core/cli/` so they can be imported, tested, and extended. This follows Go conventions where `cmd/` is minimal and delegates to library packages.

## Edge Cases & Risks

- [x] `--port` flag vs `APP_PORT` env var — flag takes precedence if provided, env var is the fallback
- [x] Running `rgo` with no args should show help, not start the server
- [x] `serve` command should work identically to current `main.go` when no flags are provided (backward compatible)
- [x] Cobra adds its own `--help` and `completion` commands — these are fine to keep
- [x] Future commands (make:*, migrate) must be able to add themselves without modifying root — achieved via `rootCmd.AddCommand()`

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Feature #01 — Project Setup | Feature | ✅ Done |
| Feature #02 — Configuration System | Feature | ✅ Done |
| `github.com/spf13/cobra` | External | 🔴 Needs install |

## Open Questions

_(All resolved)_

## Decisions Made

| Date | Decision | Rationale |
|---|---|---|
| 2026-03-06 | Use `core/cli/` package, not `cmd/` | `cmd/` is for `package main` only. CLI commands are library code — testable, importable. |
| 2026-03-06 | Scope to root + serve + version only | Roadmap says "cobra setup, base commands". Scaffolding commands belong to later features. |
| 2026-03-06 | `--port` flag on `serve` overrides `APP_PORT` | Provides runtime flexibility without editing `.env`. Flag takes precedence. |
| 2026-03-06 | No `CLIProvider` — CLI bootstraps the app, not the other way around | The CLI creates the app and calls Boot(). Providers don't register commands — commands create the app. |
| 2026-03-06 | Version as constant in `core/cli/` package | Simple, no build-time injection needed. Can be updated to use `ldflags` later. |

## Discussion Complete ✅

**Summary**: Feature #10 replaces the direct-execution `main()` with a Cobra CLI framework. Root command, `serve` subcommand (with `--port` flag), and `version` subcommand. All future CLI features attach to this command tree.
**Completed**: 2026-03-06
**Next**: Create architecture doc → `10-cli-foundation-architecture.md`
