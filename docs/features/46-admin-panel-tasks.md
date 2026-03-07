# Feature #46 — Admin Panel Scaffolding: Tasks

> **Status**: ✅ SHIPPED
> **Depends On**: #15 (Controllers), #17 (Views)
> **Branch**: `feature/46-admin-panel-scaffolding`

---

## Task List

### T1 — Create `AdminOnly()` middleware

**File**: `core/middleware/admin.go`

- [ ] Create `AdminOnly() gin.HandlerFunc`
- [ ] Read `c.GetString("role")`; if `"admin"` → `c.Next()`; else → abort 403

### T2 — Create `adminScaffold()` helper

**File**: `core/cli/make_admin.go`

- [ ] Create `adminScaffold(kind, name, tpl, path string, out io.Writer) error`
- [ ] Use `text/template` with `[[ ]]` delimiters
- [ ] Pass `{"Name": name, "Resource": toSnakeCase(name)}` to template
- [ ] Use `os.O_CREATE|os.O_EXCL|os.O_WRONLY` to prevent overwrites
- [ ] Print `"{kind} created: {path}"` on success

### T3 — Create `make:admin` command

**File**: `core/cli/make_admin.go`

- [ ] Define `makeAdminCmd` cobra command: `make:admin [resource]`, `cobra.ExactArgs(1)`
- [ ] Implement `runMakeAdmin`:
  - Generate controller at `http/controllers/admin/{resource}_controller.go`
  - Generate 4 views at `resources/views/admin/{resource}/{index,show,create,edit}.html`
  - Generate layout + dashboard (skip if exists)
- [ ] Define all 7 template constants: `adminControllerTpl`, `adminLayoutTpl`, `adminDashboardTpl`, `adminIndexTpl`, `adminShowTpl`, `adminCreateTpl`, `adminEditTpl`

### T4 — Register command in root

**File**: `core/cli/root.go`

- [ ] Add `rootCmd.AddCommand(makeAdminCmd)` in `init()`

### T5 — Write tests

**Files**: `core/cli/make_scaffold_test.go`, `core/middleware/middleware_test.go`

- [ ] Add admin scaffold tests to `core/cli/make_scaffold_test.go`
- [ ] Add AdminOnly middleware tests to `core/middleware/middleware_test.go`
