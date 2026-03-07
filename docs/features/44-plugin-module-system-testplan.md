# 🧪 Test Plan: Plugin / Module System

> **Feature**: `44` — Plugin / Module System
> **Tasks**: [`44-plugin-module-system-tasks.md`](44-plugin-module-system-tasks.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-07

---

## Test File

`core/plugin/plugin_test.go`

---

## Test Helpers

Internal mock plugins used across tests:

```go
// mockPlugin — minimal Plugin (Provider + Name only)
type mockPlugin struct { name string; registerCalled, bootCalled bool }

// mockRoutePlugin — Plugin + RouteRegistrar
type mockRoutePlugin struct { mockPlugin; routesCalled bool }

// mockCommandPlugin — Plugin + CommandRegistrar
type mockCommandPlugin struct { mockPlugin; commandsCalled bool }

// mockEventPlugin — Plugin + EventRegistrar
type mockEventPlugin struct { mockPlugin; eventsCalled bool }

// mockFullPlugin — Plugin + all optional interfaces
type mockFullPlugin struct { mockPlugin; routesCalled, commandsCalled, eventsCalled bool }
```

---

## Unit Tests

### 1. Manager Construction

| # | Test | Expectation |
|---|------|-------------|
| T01 | `TestNewManager` | `NewManager()` returns non-nil `*PluginManager`, `Plugins()` returns empty slice |

### 2. Plugin Registration

| # | Test | Expectation |
|---|------|-------------|
| T02 | `TestAddPlugin` | `Add(plugin)` returns nil error, `Plugins()` has len 1 |
| T03 | `TestAddMultiplePlugins` | Add 3 plugins, `Plugins()` returns 3 in registration order |
| T04 | `TestAddDuplicateNameReturnsError` | Second `Add()` with same `Name()` returns non-nil error |
| T05 | `TestAddDuplicateDoesNotModifyList` | After duplicate error, `Plugins()` still has len 1 |

### 3. Provider Lifecycle

| # | Test | Expectation |
|---|------|-------------|
| T06 | `TestRegisterAll` | `RegisterAll(c)` calls `Register(c)` on each plugin, verified via mock flag |
| T07 | `TestBootAll` | `BootAll(c)` calls `Boot(c)` on each plugin |
| T08 | `TestLifecycleOrder` | Register order: A before B. Boot order: A before B. Verified via ordered append to shared slice |
| T09 | `TestPluginServiceAccessible` | Plugin binds `"test.greeting"` in Register, `c.Make("test.greeting")` returns expected value after RegisterAll |

### 4. Route Registration

| # | Test | Expectation |
|---|------|-------------|
| T10 | `TestRegisterRoutes` | Plugin implementing `RouteRegistrar` has `routesCalled == true` after `RegisterRoutes(r)` |
| T11 | `TestRegisterRoutesSkipsNonRegistrar` | Plugin NOT implementing `RouteRegistrar` is not affected |
| T12 | `TestRegisterRoutesMixed` | Mix of route and non-route plugins — only route plugins called |

### 5. Command Registration

| # | Test | Expectation |
|---|------|-------------|
| T13 | `TestRegisterCommands` | Plugin returns `[]*cobra.Command`, commands added to root |
| T14 | `TestRegisterCommandsSkipsNonRegistrar` | Non-CommandRegistrar plugin is skipped |
| T15 | `TestRegisteredCommandIsUsable` | Command added by plugin is findable via `root.Commands()` |

### 6. Event Registration

| # | Test | Expectation |
|---|------|-------------|
| T16 | `TestRegisterEvents` | Plugin implementing `EventRegistrar` has `eventsCalled == true` |
| T17 | `TestRegisterEventsSkipsNonRegistrar` | Non-EventRegistrar plugin is skipped |

### 7. Full Integration

| # | Test | Expectation |
|---|------|-------------|
| T18 | `TestFullPluginAllInterfaces` | Plugin implementing all optional interfaces: all hooks called |

---

## Coverage Matrix

| Component | Tests | Coverage |
|-----------|-------|----------|
| `NewManager()` | T01 | Construction |
| `Add()` | T02–T05 | Registration, duplicates |
| `RegisterAll()` | T06, T08, T09 | Provider Register phase |
| `BootAll()` | T07, T08 | Provider Boot phase |
| `RegisterRoutes()` | T10–T12 | Route optional interface |
| `RegisterCommands()` | T13–T15 | CLI optional interface |
| `RegisterEvents()` | T16–T17 | Event optional interface |
| Full lifecycle | T18 | All interfaces together |

**Total: 18 tests**

---

## Next

Changelog → `44-plugin-module-system-changelog.md`
