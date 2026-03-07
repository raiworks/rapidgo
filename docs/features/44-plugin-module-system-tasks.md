# 📋 Tasks: Plugin / Module System

> **Feature**: `44` — Plugin / Module System
> **Architecture**: [`44-plugin-module-system-architecture.md`](44-plugin-module-system-architecture.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-07

---

## Phase A — Core Plugin Package (`core/plugin/plugin.go`)

| # | Task | Detail |
|---|------|--------|
| A1 | Create `core/plugin/plugin.go` | Package declaration, imports |
| A2 | Define `Plugin` interface | Embeds `container.Provider`, adds `Name() string` |
| A3 | Define `RouteRegistrar` interface | `RegisterRoutes(r *router.Router)` |
| A4 | Define `CommandRegistrar` interface | `Commands() []*cobra.Command` |
| A5 | Define `EventRegistrar` interface | `RegisterEvents(d *events.Dispatcher)` |
| A6 | Implement `PluginManager` struct | `plugins []Plugin`, `names map[string]bool` |
| A7 | Implement `NewManager()` | Initialize map and slice |
| A8 | Implement `Add()` | Duplicate name check, append to slice |
| A9 | Implement `Plugins()` | Return registered plugins slice |
| A10 | Implement `RegisterAll()` | Iterate and call `Register(c)` |
| A11 | Implement `BootAll()` | Iterate and call `Boot(c)` |
| A12 | Implement `RegisterRoutes()` | Type-assert `RouteRegistrar`, call `RegisterRoutes(r)` |
| A13 | Implement `RegisterCommands()` | Type-assert `CommandRegistrar`, call `Commands()`, add to root |
| A14 | Implement `RegisterEvents()` | Type-assert `EventRegistrar`, call `RegisterEvents(d)` |

**Exit**: `core/plugin` package compiles.

---

## Phase B — Tests (`core/plugin/plugin_test.go`)

| # | Task | Detail |
|---|------|--------|
| B1 | `TestNewManager` | NewManager returns non-nil with empty Plugins() |
| B2 | `TestAddPlugin` | Add a plugin, Plugins() has len 1, correct Name() |
| B3 | `TestAddMultiplePlugins` | Add 3 plugins, Plugins() returns 3 in order |
| B4 | `TestAddDuplicateNameReturnsError` | Add two plugins with same Name(), error on second |
| B5 | `TestRegisterAll` | RegisterAll() calls Register(c) on each plugin |
| B6 | `TestBootAll` | BootAll() calls Boot(c) on each plugin |
| B7 | `TestRegisterRoutes` | Plugin implementing RouteRegistrar gets called |
| B8 | `TestRegisterRoutesSkipsNonRegistrar` | Plugin NOT implementing RouteRegistrar is skipped |
| B9 | `TestRegisterCommands` | Plugin implementing CommandRegistrar adds commands to root |
| B10 | `TestRegisterCommandsSkipsNonRegistrar` | Plugin NOT implementing CommandRegistrar is skipped |
| B11 | `TestRegisterEvents` | Plugin implementing EventRegistrar gets called |
| B12 | `TestRegisterEventsSkipsNonRegistrar` | Plugin NOT implementing EventRegistrar is skipped |
| B13 | `TestPluginLifecycleOrder` | Register → Boot order verified across multiple plugins |
| B14 | `TestPluginServicesAccessible` | Plugin binds service in Register, resolvable after RegisterAll |

**Exit**: All tests pass.

---

## Phase C — CLI Accessor + Wiring

| # | Task | Detail |
|---|------|--------|
| C1 | Add `RootCmd()` accessor to `core/cli/root.go` | Returns `*cobra.Command` for plugin command registration |
| C2 | Create `app/plugins.go` | `RegisterPlugins(m *plugin.PluginManager)` function |

**Exit**: CLI is extensible, app plugin registration point exists.

---

## Phase D — Example Plugin

| # | Task | Detail |
|---|------|--------|
| D1 | Create `plugins/example/example.go` | ExamplePlugin implementing Plugin, RouteRegistrar, CommandRegistrar |
| D2 | Register in `app/plugins.go` | Import and add ExamplePlugin |

**Exit**: Example plugin compiles and demonstrates all interfaces.

---

## Phase E — Verification

| # | Task | Detail |
|---|------|--------|
| E1 | Run plugin tests | `go test ./core/plugin/...` — all pass |
| E2 | Run full test suite | `go test ./...` — all 33+ packages pass |
| E3 | Build binary | `go build -o bin/rapidgo.exe ./cmd` — compiles cleanly |

**Exit**: Feature is complete and tested.

---

## Summary

| Phase | Files | Tasks |
|-------|-------|-------|
| A — Core | 1 new | 14 |
| B — Tests | 1 new | 14 |
| C — Wiring | 1 mod, 1 new | 2 |
| D — Example | 1 new, 1 mod | 2 |
| E — Verify | — | 3 |
| **Total** | **4 new, 2 mod** | **35** |

---

## Next

Test plan → `44-plugin-module-system-testplan.md`
