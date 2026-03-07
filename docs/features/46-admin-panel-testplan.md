# Feature #46 — Admin Panel Scaffolding: Test Plan

> **Status**: 🔨 IN PROGRESS
> **Depends On**: #15 (Controllers), #17 (Views)
> **Branch**: `feature/46-admin-panel-scaffolding`

---

## Test Cases

### AdminOnly Middleware — `core/middleware/middleware_test.go`

| ID | Test | Assertion |
|---|---|---|
| T01 | AdminOnly allows request when role is "admin" | Status 200, handler reached |
| T02 | AdminOnly blocks request when role is "user" | Status 403, JSON error "admin access required" |
| T03 | AdminOnly blocks request when role is not set | Status 403, JSON error "admin access required" |

### Admin Scaffold — `core/cli/make_scaffold_test.go`

| ID | Test | Assertion |
|---|---|---|
| T04 | make:admin generates controller file | File exists at `http/controllers/admin/{resource}_controller.go`, contains `type {Name}Controller struct{}`, package is `admin` |
| T05 | make:admin generates 4 view templates | Files exist at `resources/views/admin/{resource}/{index,show,create,edit}.html`, each contains `{{ .title }}` |
| T06 | make:admin generates layout template | `resources/views/admin/layout.html` exists, contains "Admin Panel" |
| T07 | make:admin generates dashboard template | `resources/views/admin/dashboard.html` exists, contains "Dashboard" |
| T08 | make:admin skips layout if already exists | Pre-create layout file, run make:admin, verify layout content unchanged |
| T09 | make:admin prevents controller overwrite | Run make:admin twice with same name, second call returns error |
| T10 | adminScaffold uses custom delimiters correctly | Template with `[[ .Name ]]` produces correct output, `{{ .title }}` passes through literally |

### Coverage Summary

| Area | Tests | Coverage |
|---|---|---|
| AdminOnly middleware | T01–T03 | Allow, deny (wrong role), deny (no role) |
| File generation | T04–T07 | Controller, views (4), layout, dashboard |
| Overwrite protection | T08–T09 | Shared files skip, resource files error |
| Template delimiters | T10 | `[[ ]]` substitution, `{{ }}` passthrough |
| **Total** | **10** | |
