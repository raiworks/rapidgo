# 🏗️ Architecture: Plugin / Module System

> **Feature**: `44` — Plugin / Module System
> **Discussion**: [`44-plugin-module-system-discussion.md`](44-plugin-module-system-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-07

---

## Overview

A compile-time plugin system built on the existing Provider pattern. Plugins are Go packages that implement the `Plugin` interface (which embeds `container.Provider`) plus optional interfaces for routes, CLI commands, and event listeners. A `PluginManager` coordinates plugin registration, dependency validation, and subsystem wiring. Developers explicitly register plugins in application code — no directory scanning or runtime loading.

---

## File Structure

```
core/
├── plugin/
│   ├── plugin.go              # Plugin interface, optional interfaces, PluginManager
│   └── plugin_test.go         # Tests for PluginManager

core/
├── cli/
│   ├── root.go                # MODIFY — add RootCmd() accessor for plugin CLI registration

app/
├── plugins.go                 # NEW — RegisterPlugins() function, example plugin registration

plugins/
├── example/
│   └── example.go             # NEW — ExamplePlugin demonstrating all interfaces
```

**Total**: 3 new files, 1 modified file

---

## Component Design

### Plugin Interface

**Package**: `core/plugin`
**File**: `plugin.go`

```go
// Plugin is a self-contained module that integrates with the RapidGo framework.
// It embeds container.Provider for service registration (Register + Boot).
type Plugin interface {
    container.Provider

    // Name returns a unique identifier for the plugin (e.g., "notifications").
    Name() string
}
```

The core interface is deliberately minimal — just a Provider with a name. All other capabilities are declared through optional interfaces.

### Optional Interfaces

```go
// RouteRegistrar is implemented by plugins that register HTTP routes.
type RouteRegistrar interface {
    RegisterRoutes(r *router.Router)
}

// CommandRegistrar is implemented by plugins that add CLI commands.
type CommandRegistrar interface {
    Commands() []*cobra.Command
}

// EventRegistrar is implemented by plugins that listen to framework events.
type EventRegistrar interface {
    RegisterEvents(d *events.Dispatcher)
}
```

Pattern: the `PluginManager` checks each plugin with type assertions:
```go
if p, ok := plugin.(RouteRegistrar); ok {
    p.RegisterRoutes(router)
}
```

### PluginManager

```go
// PluginManager coordinates plugin registration and subsystem wiring.
type PluginManager struct {
    plugins []Plugin
    names   map[string]bool
}

// NewManager creates a PluginManager.
func NewManager() *PluginManager

// Add registers a plugin. Returns an error if:
//   - a plugin with the same Name() is already registered
func (m *PluginManager) Add(p Plugin) error

// Plugins returns all registered plugins.
func (m *PluginManager) Plugins() []Plugin

// RegisterAll calls Register(c) on each plugin in order.
func (m *PluginManager) RegisterAll(c *container.Container)

// BootAll calls Boot(c) on each plugin in order.
func (m *PluginManager) BootAll(c *container.Container)

// RegisterRoutes calls RegisterRoutes(r) on each plugin that implements RouteRegistrar.
func (m *PluginManager) RegisterRoutes(r *router.Router)

// RegisterCommands calls Commands() on each plugin that implements CommandRegistrar
// and adds the returned commands to the root Cobra command.
func (m *PluginManager) RegisterCommands(root *cobra.Command)

// RegisterEvents calls RegisterEvents(d) on each plugin that implements EventRegistrar.
func (m *PluginManager) RegisterEvents(d *events.Dispatcher)
```

### Implementation Details

`Add()`:
- Checks `m.names[p.Name()]` for duplicates, returns error if found
- Stores plugin in ordered slice
- Marks name in set

`RegisterAll()` / `BootAll()`:
- Iterates plugins in registration order, calls `Register(c)` / `Boot(c)`
- Same two-phase lifecycle as `app.App` already uses

`RegisterRoutes()`:
- Iterates plugins, type-asserts `RouteRegistrar`, calls `RegisterRoutes(r)`
- Called from `RouterProvider.Boot()` after framework routes are registered

`RegisterCommands()`:
- Iterates plugins, type-asserts `CommandRegistrar`, calls `Commands()`, adds each to root
- Called from `root.go` after framework commands are registered

`RegisterEvents()`:
- Iterates plugins, type-asserts `EventRegistrar`, calls `RegisterEvents(d)`
- Called during boot if an event dispatcher exists in the container

---

## CLI Root Accessor

**File**: `core/cli/root.go`
**Change**: Add accessor function for the root command.

```go
// RootCmd returns the root Cobra command, allowing plugins to add subcommands.
func RootCmd() *cobra.Command {
    return rootCmd
}
```

This keeps `rootCmd` unexported while providing controlled access.

---

## Integration Flow

### Application Bootstrap

**File**: `app/plugins.go`

```go
package app

import (
    "github.com/RAiWorks/RapidGo/core/plugin"
    exampleplugin "github.com/RAiWorks/RapidGo/plugins/example"
)

// RegisterPlugins registers all application plugins with the manager.
func RegisterPlugins(m *plugin.PluginManager) {
    m.Add(exampleplugin.New())
}
```

### Boot Sequence (in CLI commands that need plugins)

```go
// 1. Create plugin manager.
pm := plugin.NewManager()
app.RegisterPlugins(pm)

// 2. Register plugin services (before app providers that may depend on them).
pm.RegisterAll(application.Container)

// 3. Boot application (providers boot, including RouterProvider which calls pm.RegisterRoutes).
application.Boot()

// 4. Boot plugins (after all providers, so plugins can resolve any service).
pm.BootAll(application.Container)

// 5. Wire optional interfaces.
pm.RegisterCommands(cli.RootCmd())
if application.Container.Has("events") {
    pm.RegisterEvents(container.MustMake[*events.Dispatcher](application.Container, "events"))
}
```

For MVP, the plugin wiring is called from individual CLI commands that need it (serve, work, schedule:run). This keeps the integration explicit without modifying the `App` struct.

---

## Example Plugin

**File**: `plugins/example/example.go`

```go
package example

import (
    "fmt"
    "log/slog"

    "github.com/RAiWorks/RapidGo/core/container"
    "github.com/RAiWorks/RapidGo/core/router"
    "github.com/gin-gonic/gin"
    "github.com/spf13/cobra"
)

// ExamplePlugin demonstrates the plugin interface.
type ExamplePlugin struct{}

func New() *ExamplePlugin { return &ExamplePlugin{} }

// --- Plugin interface ---

func (p *ExamplePlugin) Name() string { return "example" }

func (p *ExamplePlugin) Register(c *container.Container) {
    c.Singleton("example.greeting", func(c *container.Container) interface{} {
        return "Hello from the Example Plugin!"
    })
}

func (p *ExamplePlugin) Boot(c *container.Container) {
    slog.Info("example plugin booted")
}

// --- RouteRegistrar ---

func (p *ExamplePlugin) RegisterRoutes(r *router.Router) {
    g := r.Group("/example")
    g.Get("/hello", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Hello from ExamplePlugin"})
    })
}

// --- CommandRegistrar ---

func (p *ExamplePlugin) Commands() []*cobra.Command {
    return []*cobra.Command{
        {
            Use:   "example:greet",
            Short: "Print a greeting from the example plugin",
            Run: func(cmd *cobra.Command, args []string) {
                fmt.Println("Hello from ExamplePlugin CLI!")
            },
        },
    }
}
```

---

## Data Flow

```
Application startup:
    → pm := plugin.NewManager()
    → app.RegisterPlugins(pm)           — developer registers plugins
    → pm.RegisterAll(container)          — plugin.Register() binds services
    → application.Boot()                 — framework providers boot (DB, Router, etc.)
    → pm.BootAll(container)              — plugin.Boot() runs (can resolve any service)
    → pm.RegisterRoutes(router)          — plugins that implement RouteRegistrar add routes
    → pm.RegisterCommands(rootCmd)       — plugins that implement CommandRegistrar add CLI commands
    → pm.RegisterEvents(dispatcher)      — plugins that implement EventRegistrar hook into events
```

---

## Naming Convention

Plugin services in the container follow `"pluginname.service"` convention:

```go
// Plugin "notifications" registers:
c.Singleton("notifications.mailer", ...)
c.Singleton("notifications.config", ...)

// Framework services (no prefix):
c.Singleton("db", ...)
c.Singleton("router", ...)
```

---

## Security Considerations

- **Compile-time only**: no dynamic code loading, no eval, no remote package fetching at runtime
- **Name uniqueness**: `PluginManager.Add()` rejects duplicate plugin names — prevents silent overwrite
- **Container isolation**: naming convention prevents plugins from accidentally overwriting framework services
- **No new attack surface**: plugins run in the same process with the same permissions as the application

---

## Trade-offs & Alternatives

| Approach | Pros | Cons | Verdict |
|---|---|---|---|
| Provider + optional interfaces | Reuses existing pattern, type-safe, IDE-friendly | Compile-time only | ✅ Selected — Go-idiomatic |
| Go `plugin` stdlib (runtime) | Hot-reload, dynamic loading | CGO-only, Linux-only, effectively abandoned | ❌ Not portable |
| HashiCorp go-plugin (gRPC) | Process isolation, language-agnostic | Massive complexity, IPC overhead, external dep | ❌ Over-engineered for this |
| Config-driven (YAML plugin manifest) | Declarative | Loses type safety, reflection-heavy | ❌ Un-Go-like |
| Directory scanning | Auto-discovery | Hidden magic, harder to debug, import side-effects | ❌ Against explicit philosophy |

---

## Future Iterations (NOT in this feature)

- Plugin configuration files (`plugins/name/config.toml`)
- Plugin asset loading (views, static files)
- Plugin dependency graph with topological sorting
- Plugin enable/disable via config
- Plugin versioning with framework compatibility checks
- `make:plugin` scaffold CLI command
- Plugin marketplace / remote registry

---

## Next

Tasks doc → `44-plugin-module-system-tasks.md`
