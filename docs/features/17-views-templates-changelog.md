# 📝 Changelog: Views & Templates

> **Feature**: `17` — Views & Templates
> **Branch**: `feature/17-views-templates`
> **Started**: 2026-03-06
> **Completed**: 2026-03-06

---

## Log

### Phase A — View Engine & Static Serving
- Created `core/router/view.go` — `DefaultFuncMap()` with `route` helper
- Added `SetFuncMap`, `LoadTemplates`, `Static`, `StaticFile` methods to `core/router/router.go`
- Updated `app/providers/router_provider.go` — template engine + static serving in Boot
- Created `resources/views/home.html` — sample HTML5 template
- Updated `http/controllers/home_controller.go` — `c.JSON()` → `c.HTML()`
- `go build` clean, `go vet` clean

### Phase B — Testing
- Created `core/router/view_test.go` — 5 tests (FuncMap, template rendering, static serving)
- Updated `http/controllers/controllers_test.go` — Home test asserts HTML response with templates loaded
- `go test ./core/router/... -v` — all pass
- `go test ./http/controllers/... -v` — all pass
- `go test ./... -count=1` — 177 total tests, 0 failures

### Phase C — Changelog & Cross-Check
- Code vs architecture: 1 deviation (see below)
- Changelog updated

---

## Deviations from Plan

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|
| Provider guards for template/static paths | Unconditional `LoadTemplates` and `Static` calls | Added `filepath.Glob` check before `LoadTemplates`, `os.Stat` checks before `Static` | `LoadHTMLGlob` panics when pattern matches no files; tests run from package dirs where `resources/views/` doesn't exist |

## Key Decisions Made During Build

| Decision | Context | Date |
|---|---|---|
| Guard LoadTemplates with glob check | `TestFullAppBootstrap_WithRouter` and `TestNewApp_ReturnsBootedApp` panicked because `resources/views/**/*` matched no files when tests run from `app/providers/` or `core/cli/` directories | 2026-03-06 |
| Guard Static with os.Stat check | Same issue — directory doesn't exist when running from test package dirs | 2026-03-06 |
