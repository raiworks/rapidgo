# 💬 Discussion: Views & Templates

> **Feature**: `17` — Views & Templates
> **Depends on**: #07 (Router), #15 (Controllers)
> **Status**: 🟢 RESOLVED

---

## Context

Features #07 and #15 built the router and controller layers. All current controllers return JSON via `c.JSON()`. The framework needs server-side rendering (SSR) using Go's standard `html/template` package, template function registration (including the existing `route` helper from #07), static file serving, and a sample template to verify the pipeline.

## Blueprint Reference

The blueprint's **View / Template Engine** section (lines 1937–1954) prescribes:
- `LoadHTMLGlob("resources/views/**/*")` on the Gin engine
- Templates stored in `resources/views/`
- Templ mentioned as optional alternative (out of scope)

The framework doc `docs/framework/http/views.md` expands on:
- Template directory structure (`resources/views/`)
- Data passing via `gin.H`
- Template functions (`route` helper from named routes)
- Static file serving: `/static` → `resources/static/`, `/uploads` → `storage/uploads/`, `/favicon.ico`

## Scope Decision

Feature #17 implements the **template engine setup and static file serving** on the Router, plus a sample home template and updated controller. Specifically:

1. **Router methods**: `LoadTemplates(pattern)`, `SetFuncMap(funcMap)`, `Static(path, dir)`, `StaticFile(path, file)` — thin wrappers on the Gin engine
2. **Default template functions**: `DefaultFuncMap()` function exposing `route` helper
3. **RouterProvider update**: Call `SetFuncMap(DefaultFuncMap())` then `LoadTemplates("resources/views/**/*")` and configure static paths
4. **Sample template**: `resources/views/home.html`
5. **Home controller update**: Render `home.html` via `c.HTML()` instead of `c.JSON()`

**Out of scope**: Templ integration, form validation templates, flash messages, CSRF (those belong to later features: Sessions #20, Validation #24+, CSRF #29+).

## Key Decisions

| # | Decision | Rationale |
|---|---|---|
| 1 | Thin wrappers on Router | Consistent with existing Router pattern (Get, Post, etc. all delegate to engine) |
| 2 | `DefaultFuncMap()` function in `core/router/` | Co-locates with `Route()` named routes helper already in `core/router/named.go` |
| 3 | Static serving in RouterProvider | Blueprint shows it in router setup; provider is the right boot-time hook |
| 4 | `c.HTML()` for Home | Proves the template pipeline works end-to-end |
| 5 | `LoadTemplates` wraps `LoadHTMLGlob` | Framework-consistent naming; single call site |
