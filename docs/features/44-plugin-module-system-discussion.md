# 💬 Discussion: Plugin / Module System

> **Feature**: `44` — Plugin / Module System
> **Status**: 🟢 COMPLETE
> **Date**: 2026-03-07

---

## What Are We Building?

A framework-native plugin system that allows self-contained modules to integrate with all RapidGo subsystems — services, routes, middleware, CLI commands, migrations, events, and queue handlers — through a single unified interface. Plugins are Go packages that implement the `Plugin` interface and are registered in application code.

---

## Why?

- **Modularity**: large applications need self-contained feature bundles (notifications, analytics, admin panel) that can be added or removed cleanly
- **Reusability**: plugins can be shared across projects as Go modules
- **Extension points**: the framework already has registries (middleware, migrations, seeders, queue handlers) — a plugin system unifies them under one interface
- **Laravel parity**: Laravel packages/service providers are a major ecosystem driver

---

## Prior Art

| System | Approach | Notes |
|---|---|---|
| Laravel | Service Providers + Package Discovery | Providers are plugins. Composer auto-discovers providers via `extra.laravel.providers` in composer.json |
| WordPress | Plugin API (hooks/filters) | File-based discovery, global hook system |
| Gin middleware | func(c *gin.Context) | Per-handler only, no lifecycle |
| Go plugin stdlib | `plugin.Open()` | CGO-only, Linux-only, no Windows, abandoned pattern |

---

## Constraints

1. **Compile-time composition** — plugins are Go packages imported and registered in application code, NOT runtime-loaded binaries (Go's `plugin` stdlib is CGO-only, Linux-only, and effectively abandoned)
2. **Provider-based** — plugins implement `container.Provider` (Register + Boot) plus metadata — no new lifecycle, reuse existing
3. **No discovery magic** — developer explicitly registers plugins in `app/plugins.go` (like providers in `root.go`), no directory scanning
4. **Namespace convention** — plugin services use `"pluginname.service"` naming in the container to prevent collisions
5. **Optional interfaces** — routes, CLI commands, migrations, and event listeners are declared via optional interfaces, not forced on every plugin
6. **MVP scope** — Plugin interface, PluginManager, route/CLI/migration hooks, example plugin. NO asset loading, no fluent config, no remote plugin registry

---

## Decision Log

| # | Decision | Rationale |
|---|---|---|
| 1 | Compile-time import, not runtime loading | Go `plugin` stdlib is CGO-only / Linux-only. Compile-time gives type safety, IDE support, and cross-platform compatibility |
| 2 | Extend `container.Provider`, don't replace it | Plugins ARE providers with metadata. No new lifecycle. No breaking changes to existing boot flow |
| 3 | Explicit registration, no directory scanning | Predictable, debuggable, no hidden magic. Developer controls what runs |
| 4 | Optional interfaces for routes/CLI/migrations | Not every plugin needs routes. Pattern: `if p, ok := plugin.(RouteRegistrar); ok { ... }` |
| 5 | PluginManager struct manages ordering and hooks | Central coordinator that handles dependency ordering, duplicate detection, and subsystem wiring |
| 6 | `"pluginname.service"` naming convention | Container is a flat namespace; prefix-based convention prevents collisions without requiring framework changes to Container |
| 7 | Expose `rootCmd` via `cli.RootCmd()` accessor | Plugins need to add CLI commands. Accessor function (not exported var) is safe and testable |

---

## Open Questions

_None — all resolved._

---

## Next

Architecture doc → `44-plugin-module-system-architecture.md`
