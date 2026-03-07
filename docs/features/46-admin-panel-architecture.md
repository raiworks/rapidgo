# Feature #46 — Admin Panel Scaffolding: Architecture

> **Status**: 🔨 IN PROGRESS
> **Depends On**: #15 (Controllers), #17 (Views)
> **Branch**: `docs/46-admin-panel-scaffolding`

---

## New File: `core/middleware/admin.go`

### Function

```go
// AdminOnly returns a Gin middleware that restricts access to admin users.
// It reads the "role" value from the Gin context and aborts with 403 if
// the role is not "admin". Must be used after a middleware that sets
// c.Set("role", ...). Note: the built-in AuthMiddleware only sets "user_id";
// you must add your own middleware to load and set the user's role.
func AdminOnly() gin.HandlerFunc
```

### Behavior

1. Reads `c.GetString("role")` from the request context.
2. If the role is `"admin"`, calls `c.Next()`.
3. Otherwise, aborts with `c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})`.

No database access, no JWT parsing — purely checks a context value set by a prior middleware.

---

## New File: `core/cli/make_admin.go`

### Functions

```go
// adminScaffold generates a file from a template using [[ ]] delimiters.
// This avoids conflicts with {{ }} used in Go/Gin HTML templates.
// The full output path (not just the directory) is provided.
// Creates parent directories with os.MkdirAll(filepath.Dir(path), 0750).
// Prevents overwrites with os.O_CREATE|os.O_EXCL|os.O_WRONLY.
func adminScaffold(kind, name, tpl, path string, out io.Writer) error
```

Template data passed to `adminScaffold`:

| Key | Value | Example |
|---|---|---|
| `Name` | Original resource name | `Post` |
| `Resource` | `toSnakeCase(Name)` | `post` |

### Command

```go
var makeAdminCmd = &cobra.Command{
    Use:   "make:admin [resource]",
    Short: "Create an admin controller and views for a resource",
    Args:  cobra.ExactArgs(1),
    RunE:  runMakeAdmin,
}
```

### `runMakeAdmin` Logic

1. Compute `resource = toSnakeCase(name)`.
2. Generate controller: `http/controllers/admin/{resource}_controller.go` using `adminControllerTpl`.
3. Generate 4 view templates in `resources/views/admin/{resource}/`:
   - `index.html` using `adminIndexTpl`
   - `show.html` using `adminShowTpl`
   - `create.html` using `adminCreateTpl`
   - `edit.html` using `adminEditTpl`
4. Generate shared files (skip if exists, no error):
   - `resources/views/admin/layout.html` using `adminLayoutTpl`
   - `resources/views/admin/dashboard.html` using `adminDashboardTpl`

For shared files (layout, dashboard): `runMakeAdmin` calls `adminScaffold()` which uses `os.O_CREATE|os.O_EXCL|os.O_WRONLY`. If the file already exists, `runMakeAdmin` checks `errors.Is(err, os.ErrExist)` and skips silently (no error returned, no message printed). For resource files (controller, views), the error is returned to the caller.

### Generated Controller Template (`adminControllerTpl`)

Package: `admin` (subpackage of controllers).

```go
package admin

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type [[.Name]]Controller struct{}

func (ctrl *[[.Name]]Controller) Index(c *gin.Context) {
    c.HTML(http.StatusOK, "admin/[[.Resource]]/index.html", gin.H{
        "title": "[[.Name]] List",
    })
}

func (ctrl *[[.Name]]Controller) Create(c *gin.Context) {
    c.HTML(http.StatusOK, "admin/[[.Resource]]/create.html", gin.H{
        "title": "Create [[.Name]]",
    })
}

func (ctrl *[[.Name]]Controller) Store(c *gin.Context) {
    // TODO: bind form, validate, save to database
    c.Redirect(http.StatusFound, "/admin/[[.Resource]]")
}

func (ctrl *[[.Name]]Controller) Show(c *gin.Context) {
    id := c.Param("id")
    c.HTML(http.StatusOK, "admin/[[.Resource]]/show.html", gin.H{
        "title": "[[.Name]] Details",
        "id":    id,
    })
}

func (ctrl *[[.Name]]Controller) Edit(c *gin.Context) {
    id := c.Param("id")
    c.HTML(http.StatusOK, "admin/[[.Resource]]/edit.html", gin.H{
        "title": "Edit [[.Name]]",
        "id":    id,
    })
}

func (ctrl *[[.Name]]Controller) Update(c *gin.Context) {
    // TODO: bind form, validate, update database
    c.Redirect(http.StatusFound, "/admin/[[.Resource]]")
}

func (ctrl *[[.Name]]Controller) Destroy(c *gin.Context) {
    // TODO: delete from database
    c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
```

### Generated View Templates

All view templates use `{{ }}` for Gin template syntax (passed through literally by `[[ ]]` scaffold delimiters).

**`adminLayoutTpl`** — Base admin layout with Tailwind CSS:
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .title }} — Admin</title>
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
</head>
<body class="bg-gray-100 min-h-screen">
    <nav class="bg-gray-800 text-white p-4">
        <span class="font-bold text-lg">Admin Panel</span>
    </nav>
    <main class="container mx-auto p-6">
        {{ .content }}
    </main>
</body>
</html>
```

**`adminDashboardTpl`** — Simple dashboard page:
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dashboard — Admin</title>
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
</head>
<body class="bg-gray-100 min-h-screen">
    <nav class="bg-gray-800 text-white p-4">
        <span class="font-bold text-lg">Admin Panel</span>
    </nav>
    <main class="container mx-auto p-6">
        <h1 class="text-2xl font-bold mb-4">Dashboard</h1>
        <p class="text-gray-600">Welcome to the admin panel.</p>
    </main>
</body>
</html>
```

**`adminIndexTpl`** — Resource list:
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .title }} — Admin</title>
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
</head>
<body class="bg-gray-100 min-h-screen">
    <nav class="bg-gray-800 text-white p-4">
        <span class="font-bold text-lg">Admin Panel</span>
    </nav>
    <main class="container mx-auto p-6">
        <h1 class="text-2xl font-bold mb-4">{{ .title }}</h1>
        <!-- TODO: list records in a table -->
    </main>
</body>
</html>
```

**`adminShowTpl`** — Single record detail:
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .title }} — Admin</title>
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
</head>
<body class="bg-gray-100 min-h-screen">
    <nav class="bg-gray-800 text-white p-4">
        <span class="font-bold text-lg">Admin Panel</span>
    </nav>
    <main class="container mx-auto p-6">
        <h1 class="text-2xl font-bold mb-4">{{ .title }}</h1>
        <p>ID: {{ .id }}</p>
        <!-- TODO: display record fields -->
    </main>
</body>
</html>
```

**`adminCreateTpl`** — Create form:
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .title }} — Admin</title>
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
</head>
<body class="bg-gray-100 min-h-screen">
    <nav class="bg-gray-800 text-white p-4">
        <span class="font-bold text-lg">Admin Panel</span>
    </nav>
    <main class="container mx-auto p-6">
        <h1 class="text-2xl font-bold mb-4">{{ .title }}</h1>
        <form method="POST" action="/admin/[[.Resource]]">
            <!-- TODO: add form fields -->
            <button type="submit" class="bg-blue-600 text-white px-4 py-2 rounded">Create</button>
        </form>
    </main>
</body>
</html>
```

**`adminEditTpl`** — Edit form:
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .title }} — Admin</title>
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
</head>
<body class="bg-gray-100 min-h-screen">
    <nav class="bg-gray-800 text-white p-4">
        <span class="font-bold text-lg">Admin Panel</span>
    </nav>
    <main class="container mx-auto p-6">
        <h1 class="text-2xl font-bold mb-4">{{ .title }}</h1>
        <form method="POST" action="/admin/[[.Resource]]/{{ .id }}">
            <!-- TODO: add form fields, populate with existing values -->
            <button type="submit" class="bg-blue-600 text-white px-4 py-2 rounded">Update</button>
        </form>
    </main>
</body>
</html>
```

Note: The create and edit templates use `[[ .Resource ]]` for the resource path (substituted at scaffold time) and `{{ .id }}` for Gin template variables (rendered at request time).

---

## Modified File: `core/cli/root.go`

### Change in `init()`

Add `makeAdminCmd` to the root command:

```go
rootCmd.AddCommand(makeAdminCmd)
```

---

## No New Dependencies

This feature uses only the standard library (`text/template`, `os`, `path/filepath`, `fmt`, `io`, `net/http`) and existing framework dependencies (Gin, Cobra).

## No Environment Variables

No env vars. The admin panel is scaffolded on demand via the CLI.

## No Database Changes

No migrations, no new models. The admin controllers interact with existing models.
