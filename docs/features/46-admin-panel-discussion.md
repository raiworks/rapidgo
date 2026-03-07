# Feature #46 — Admin Panel Scaffolding: Discussion

> **Status**: ✅ SHIPPED
> **Depends On**: #15 (Controllers), #17 (Views)
> **Branch**: `docs/46-admin-panel-scaffolding`

---

## What Problem Does This Solve?

Every web application eventually needs an admin panel — a protected area where administrators manage data through CRUD interfaces. Building admin controllers, route groups, role-based middleware, and HTML templates from scratch for each resource is repetitive boilerplate that slows development.

Currently, RapidGo has `make:controller` which generates API-style controllers returning JSON. There is no way to quickly scaffold an admin CRUD controller that renders HTML templates, no admin-specific middleware for role checks, and no base admin layout template.

## What Does This Feature Add?

Three framework additions:

1. **`make:admin [Resource]` CLI command** — Generates an admin controller with 7 CRUD methods rendering HTML templates, plus 4 view templates (index, show, create, edit) and a shared admin layout (created once, never overwritten).

2. **`AdminOnly()` middleware** — A Gin middleware that checks the user's role from the request context. If the role is not `"admin"`, the request is aborted with 403 Forbidden. Designed to stack with the existing `AuthMiddleware()`.

3. **Admin layout template** — A base HTML template (`resources/views/admin/layout.html`) with Tailwind CSS, a sidebar navigation placeholder, and a content area. Generated once by the first `make:admin` invocation.

### Key Design Decisions

1. **Extends existing `scaffold()` pattern** — The new `adminScaffold()` helper reuses the same file-safety approach (`os.O_CREATE|os.O_EXCL` to prevent overwrites) but uses `[[ ]]` template delimiters to avoid conflicts with Go/Gin `{{ }}` template syntax in generated HTML files.

2. **ResourceController interface** — Generated admin controllers implement the existing `ResourceController` interface (Index, Create, Store, Show, Edit, Update, Destroy), so they work with `group.Resource()` for route registration. No new interfaces.

3. **Role from context** — `AdminOnly()` reads `c.GetString("role")` from the Gin context. The developer is responsible for setting this value (e.g., by extending `AuthMiddleware` or adding a middleware that loads the user and sets the role). This keeps the middleware decoupled from any specific auth implementation.

4. **Generated code is a starting point** — The generated controller and templates are minimal scaffolds with TODO comments. Developers are expected to customize them — adding model queries, form handling, validation, and styling.

5. **No runtime admin framework** — This is purely a code generation tool. There is no `Admin` struct, no dynamic resource registration, no auto-generated CRUD from model reflection. The generated code is plain Go that developers own and modify.

## What's Out of Scope?

- **Automatic CRUD from models** — No model introspection or reflection-based admin. The developer writes the query logic.
- **Admin authentication flow** — Login/logout pages. The existing `AuthMiddleware()` handles token validation; `AdminOnly()` only checks the role.
- **Permission system / RBAC** — Fine-grained permissions beyond a simple role check. Developers can implement this by extending `AdminOnly()`.
- **Pre-built admin dashboard with charts** — The generated dashboard template is a placeholder. Analytics and charts are application-specific.
- **Asset pipeline / CSS bundling** — Uses Tailwind CSS via CDN link in the layout template.
- **Auto-registration of admin routes** — The developer registers routes manually in `routes/admin.go` or similar, following the existing pattern from `routes/web.go` and `routes/api.go`.

## How Will Developers Use It?

### 1. Generate admin scaffold

```bash
rapidgo make:admin Post
```

Output:
```
Admin controller created: http/controllers/admin/post_controller.go
Admin view created: resources/views/admin/post/index.html
Admin view created: resources/views/admin/post/show.html
Admin view created: resources/views/admin/post/create.html
Admin view created: resources/views/admin/post/edit.html
Admin layout created: resources/views/admin/layout.html
Admin dashboard created: resources/views/admin/dashboard.html
```

### 2. Register admin routes

```go
// routes/admin.go
func RegisterAdmin(r *router.Router) {
    admin := r.Group("/admin", middleware.AuthMiddleware(), middleware.AdminOnly())
    admin.Get("/", adminControllers.Dashboard)
    admin.Resource("/posts", &adminControllers.PostController{})
}
```

### 3. Customize generated code

Edit the controller to add model queries, form binding, and business logic. Edit the templates to build the actual UI.
