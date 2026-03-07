package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TC-01: DefaultFuncMap contains route key
func TestDefaultFuncMap_ContainsRoute(t *testing.T) {
	fm := DefaultFuncMap()
	if fm["route"] == nil {
		t.Fatal("expected FuncMap to contain 'route'")
	}
}

// TC-02: DefaultFuncMap route function resolves named routes
func TestDefaultFuncMap_RouteWorks(t *testing.T) {
	ResetNamedRoutes()
	Name("home", "/")

	fm := DefaultFuncMap()
	routeFn := fm["route"].(func(string, ...string) string)
	result := routeFn("home")
	if result != "/" {
		t.Fatalf("expected '/', got '%s'", result)
	}
}

// TC-03: LoadTemplates + c.HTML renders template with data
func TestLoadTemplates_RenderHTML(t *testing.T) {
	dir := t.TempDir()
	tmpl := `<h1>{{ .title }}</h1>`
	if err := os.WriteFile(filepath.Join(dir, "test.html"), []byte(tmpl), 0644); err != nil {
		t.Fatal(err)
	}

	r := newTestRouter()
	r.LoadTemplates(dir)
	r.Get("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test.html", gin.H{"title": "Hello"})
	})

	w := doRequest(r, "GET", "/")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "<h1>Hello</h1>") {
		t.Fatalf("expected rendered template, got: %s", w.Body.String())
	}
}

// TC-04: Static serves files from a directory
func TestStatic_ServesFiles(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "test.txt"), []byte("static content"), 0644); err != nil {
		t.Fatal(err)
	}

	r := newTestRouter()
	r.Static("/assets", dir)

	w := doRequest(r, "GET", "/assets/test.txt")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "static content") {
		t.Fatalf("expected 'static content', got: %s", w.Body.String())
	}
}

// TC-05: StaticFile serves a single file
func TestStaticFile_ServesSingleFile(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "robots.txt")
	if err := os.WriteFile(filePath, []byte("User-agent: *"), 0644); err != nil {
		t.Fatal(err)
	}

	r := newTestRouter()
	r.StaticFile("/robots.txt", filePath)

	w := doRequest(r, "GET", "/robots.txt")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "User-agent: *") {
		t.Fatalf("expected 'User-agent: *', got: %s", w.Body.String())
	}
}
