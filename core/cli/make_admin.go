package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

// adminScaffold generates a file from a template using [[ ]] delimiters.
// This avoids conflicts with {{ }} used in Go/Gin HTML templates.
// The full output path (not just the directory) is provided.
// Creates parent directories with os.MkdirAll(filepath.Dir(path), 0750).
// Prevents overwrites with os.O_CREATE|os.O_EXCL|os.O_WRONLY.
func adminScaffold(kind, name, tpl, path string, out io.Writer) error {
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("file already exists or cannot be created: %w", err)
	}
	defer f.Close()

	t := template.Must(template.New(kind).Delims("[[", "]]").Parse(tpl))
	data := map[string]string{
		"Name":     name,
		"Resource": toSnakeCase(name),
	}
	if err := t.Execute(f, data); err != nil {
		return fmt.Errorf("failed to write template: %w", err)
	}

	fmt.Fprintf(out, "%s created: %s\n", kind, path)
	return nil
}

var makeAdminCmd = &cobra.Command{
	Use:   "make:admin [resource]",
	Short: "Create an admin controller and views for a resource",
	Args:  cobra.ExactArgs(1),
	RunE:  runMakeAdmin,
}

func runMakeAdmin(cmd *cobra.Command, args []string) error {
	name := args[0]
	resource := toSnakeCase(name)
	out := cmd.OutOrStdout()

	// Controller
	ctrlPath := filepath.Join("http", "controllers", "admin", resource+"_controller.go")
	if err := adminScaffold("Admin controller", name, adminControllerTpl, ctrlPath, out); err != nil {
		return err
	}

	// Resource views
	views := []struct {
		kind string
		file string
		tpl  string
	}{
		{"Admin view", "index.html", adminIndexTpl},
		{"Admin view", "show.html", adminShowTpl},
		{"Admin view", "create.html", adminCreateTpl},
		{"Admin view", "edit.html", adminEditTpl},
	}
	for _, v := range views {
		path := filepath.Join("resources", "views", "admin", resource, v.file)
		if err := adminScaffold(v.kind, name, v.tpl, path, out); err != nil {
			return err
		}
	}

	// Shared files — skip silently if they already exist
	shared := []struct {
		kind string
		path string
		tpl  string
	}{
		{"Admin layout", filepath.Join("resources", "views", "admin", "layout.html"), adminLayoutTpl},
		{"Admin dashboard", filepath.Join("resources", "views", "admin", "dashboard.html"), adminDashboardTpl},
	}
	for _, s := range shared {
		err := adminScaffold(s.kind, name, s.tpl, s.path, out)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	return nil
}

// --- Admin Templates ---

var adminControllerTpl = `package admin

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
`

var adminLayoutTpl = `<!DOCTYPE html>
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
`

var adminDashboardTpl = `<!DOCTYPE html>
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
`

var adminIndexTpl = `<!DOCTYPE html>
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
`

var adminShowTpl = `<!DOCTYPE html>
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
`

var adminCreateTpl = `<!DOCTYPE html>
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
`

var adminEditTpl = `<!DOCTYPE html>
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
`
