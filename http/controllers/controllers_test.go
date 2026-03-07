package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/RAiWorks/RGo/core/router"
	"github.com/RAiWorks/RGo/http/controllers"
	"github.com/RAiWorks/RGo/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// newTestContext creates a Gin context backed by an httptest recorder.
func newTestContext(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, nil)
	return c, w
}

// bodyMap reads the response body as a JSON map.
func bodyMap(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
	return m
}

// testTemplateDir creates a temp dir with home.html for template tests.
func testTemplateDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	tmpl := `<!DOCTYPE html><html><head><title>{{ .title }}</title></head><body><h1>{{ .title }}</h1></body></html>`
	if err := os.WriteFile(filepath.Join(dir, "home.html"), []byte(tmpl), 0644); err != nil {
		t.Fatal(err)
	}
	return dir
}

// TC-01: Home renders HTML template with status 200
func TestHome_RendersHTML(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	dir := testTemplateDir(t)

	r := router.New()
	r.LoadTemplates(dir)
	r.Get("/", controllers.Home)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	ct := w.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Fatalf("expected Content-Type text/html, got '%s'", ct)
	}
	if !strings.Contains(w.Body.String(), "Welcome to RGo") {
		t.Fatalf("expected body to contain 'Welcome to RGo', got: %s", w.Body.String())
	}
}

// TC-02: PostController.Index returns 200 with index message
func TestPostController_Index(t *testing.T) {
	c, w := newTestContext("GET", "/posts")
	ctrl := &controllers.PostController{}
	ctrl.Index(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	m := bodyMap(t, w)
	if m["message"] != "PostController index" {
		t.Fatalf("expected 'PostController index', got '%v'", m["message"])
	}
}

// TC-03: PostController.Store returns 201
func TestPostController_Store(t *testing.T) {
	c, w := newTestContext("POST", "/posts")
	ctrl := &controllers.PostController{}
	ctrl.Store(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
	m := bodyMap(t, w)
	if m["message"] != "PostController store" {
		t.Fatalf("expected 'PostController store', got '%v'", m["message"])
	}
}

// TC-04: PostController.Show returns id from URL param
func TestPostController_Show(t *testing.T) {
	c, w := newTestContext("GET", "/posts/42")
	c.Params = gin.Params{{Key: "id", Value: "42"}}
	ctrl := &controllers.PostController{}
	ctrl.Show(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	m := bodyMap(t, w)
	if m["id"] != "42" {
		t.Fatalf("expected id '42', got '%v'", m["id"])
	}
}

// TC-05: PostController.Update returns 200 with id
func TestPostController_Update(t *testing.T) {
	c, w := newTestContext("PUT", "/posts/7")
	c.Params = gin.Params{{Key: "id", Value: "7"}}
	ctrl := &controllers.PostController{}
	ctrl.Update(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	m := bodyMap(t, w)
	if m["id"] != "7" {
		t.Fatalf("expected id '7', got '%v'", m["id"])
	}
}

// TC-06: PostController.Destroy returns 200
func TestPostController_Destroy(t *testing.T) {
	c, w := newTestContext("DELETE", "/posts/1")
	ctrl := &controllers.PostController{}
	ctrl.Destroy(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	m := bodyMap(t, w)
	if m["message"] != "PostController destroy" {
		t.Fatalf("expected 'PostController destroy', got '%v'", m["message"])
	}
}

// TC-07: PostController implements ResourceController (compile-time check)
func TestPostController_ImplementsResourceController(t *testing.T) {
	var _ router.ResourceController = (*controllers.PostController)(nil)
}

// TC-08: GET / is registered and returns 200
func TestRoutes_HomeRegistered(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	dir := testTemplateDir(t)

	r := router.New()
	r.LoadTemplates(filepath.Join(dir, "*"))
	routes.RegisterWeb(r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

// TC-09: GET /api/posts is registered and returns 200
func TestRoutes_APIPostsRegistered(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	r := router.New()
	routes.RegisterAPI(r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/posts", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
