package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// --- Test Helpers ---

// newTestRouter creates a Router in test mode without Gin debug output.
func newTestRouter() *Router {
	gin.SetMode(gin.TestMode)
	return &Router{engine: gin.New()}
}

// okHandler returns a handler that responds with 200 and the given body.
func okHandler(body string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, body)
	}
}

// doRequest performs an HTTP request against the router and returns the recorder.
func doRequest(r *Router, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	r.ServeHTTP(w, req)
	return w
}

// --- mockResourceController implements ResourceController for testing ---

type mockResourceController struct{}

func (m *mockResourceController) Index(c *gin.Context)   { c.String(http.StatusOK, "index") }
func (m *mockResourceController) Create(c *gin.Context)  { c.String(http.StatusOK, "create") }
func (m *mockResourceController) Store(c *gin.Context)   { c.String(http.StatusOK, "store") }
func (m *mockResourceController) Show(c *gin.Context)    { c.String(http.StatusOK, "show") }
func (m *mockResourceController) Edit(c *gin.Context)    { c.String(http.StatusOK, "edit") }
func (m *mockResourceController) Update(c *gin.Context)  { c.String(http.StatusOK, "update") }
func (m *mockResourceController) Destroy(c *gin.Context) { c.String(http.StatusOK, "destroy") }

// TC-01: Router.New creates a valid router
func TestNew_CreatesValidRouter(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	r := New()
	if r == nil {
		t.Fatal("New() returned nil router")
	}
	if r.Engine() == nil {
		t.Fatal("Engine() returned nil")
	}
}

// TC-02: Gin mode set to ReleaseMode in production
func TestGinMode_Production(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	New()
	if gin.Mode() != gin.ReleaseMode {
		t.Errorf("expected %q, got %q", gin.ReleaseMode, gin.Mode())
	}
}

// TC-03: Gin mode set to TestMode in testing
func TestGinMode_Testing(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	New()
	if gin.Mode() != gin.TestMode {
		t.Errorf("expected %q, got %q", gin.TestMode, gin.Mode())
	}
}

// TC-04: Gin mode defaults to DebugMode for development
func TestGinMode_Development(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	New()
	if gin.Mode() != gin.DebugMode {
		t.Errorf("expected %q, got %q", gin.DebugMode, gin.Mode())
	}
}

// TC-05: GET route responds with 200
func TestGet_Route(t *testing.T) {
	r := newTestRouter()
	r.Get("/ping", okHandler("pong"))
	w := doRequest(r, http.MethodGet, "/ping")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "pong" {
		t.Errorf("expected 'pong', got %q", w.Body.String())
	}
}

// TC-06: POST route responds with 200
func TestPost_Route(t *testing.T) {
	r := newTestRouter()
	r.Post("/data", okHandler("created"))
	w := doRequest(r, http.MethodPost, "/data")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// TC-07: PUT route responds with 200
func TestPut_Route(t *testing.T) {
	r := newTestRouter()
	r.Put("/items/:id", okHandler("updated"))
	w := doRequest(r, http.MethodPut, "/items/1")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// TC-08: DELETE route responds with 200
func TestDelete_Route(t *testing.T) {
	r := newTestRouter()
	r.Delete("/items/:id", okHandler("deleted"))
	w := doRequest(r, http.MethodDelete, "/items/1")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// TC-09: PATCH route responds with 200
func TestPatch_Route(t *testing.T) {
	r := newTestRouter()
	r.Patch("/items/:id", okHandler("patched"))
	w := doRequest(r, http.MethodPatch, "/items/1")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// TC-10: OPTIONS route responds with 200
func TestOptions_Route(t *testing.T) {
	r := newTestRouter()
	r.Options("/cors", okHandler("options"))
	w := doRequest(r, http.MethodOptions, "/cors")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// TC-11: Route group adds prefix to paths
func TestGroup_AddsPrefix(t *testing.T) {
	r := newTestRouter()
	api := r.Group("/api")
	api.Get("/users", okHandler("users"))
	w := doRequest(r, http.MethodGet, "/api/users")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// TC-12: Nested route group combines prefixes
func TestGroup_Nested(t *testing.T) {
	r := newTestRouter()
	api := r.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/users", okHandler("v1-users"))
	w := doRequest(r, http.MethodGet, "/api/v1/users")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "v1-users" {
		t.Errorf("expected 'v1-users', got %q", w.Body.String())
	}
}

// TC-13: Resource registers all 7 RESTful routes
func TestResource_RegistersAll7Routes(t *testing.T) {
	r := newTestRouter()
	ctrl := &mockResourceController{}
	r.Resource("/posts", ctrl)

	tests := []struct {
		method string
		path   string
		body   string
	}{
		{http.MethodGet, "/posts", "index"},
		{http.MethodGet, "/posts/create", "create"},
		{http.MethodPost, "/posts", "store"},
		{http.MethodGet, "/posts/1", "show"},
		{http.MethodGet, "/posts/1/edit", "edit"},
		{http.MethodPut, "/posts/1", "update"},
		{http.MethodDelete, "/posts/1", "destroy"},
	}

	for _, tt := range tests {
		w := doRequest(r, tt.method, tt.path)
		if w.Code != http.StatusOK {
			t.Errorf("%s %s: expected 200, got %d", tt.method, tt.path, w.Code)
		}
		if w.Body.String() != tt.body {
			t.Errorf("%s %s: expected %q, got %q", tt.method, tt.path, tt.body, w.Body.String())
		}
	}
}

// TC-14: APIResource registers 5 routes (no Create/Edit)
func TestAPIResource_Registers5Routes(t *testing.T) {
	r := newTestRouter()
	ctrl := &mockResourceController{}
	r.APIResource("/users", ctrl)

	// These 5 should return 200
	okRoutes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/users"},
		{http.MethodPost, "/users"},
		{http.MethodGet, "/users/1"},
		{http.MethodPut, "/users/1"},
		{http.MethodDelete, "/users/1"},
	}
	for _, tt := range okRoutes {
		w := doRequest(r, tt.method, tt.path)
		if w.Code != http.StatusOK {
			t.Errorf("%s %s: expected 200, got %d", tt.method, tt.path, w.Code)
		}
	}

	// Edit route should NOT be registered.
	// Note: GET /users/create matches /users/:id with id="create" (Gin behavior),
	// so we verify no explicit /edit sub-route is registered.
	w := doRequest(r, http.MethodGet, "/users/1/edit")
	if w.Code != http.StatusNotFound {
		t.Errorf("GET /users/1/edit: expected 404, got %d", w.Code)
	}
}

// TC-15: Resource routes on RouteGroup combine prefix
func TestResource_OnGroup(t *testing.T) {
	r := newTestRouter()
	ctrl := &mockResourceController{}
	api := r.Group("/api")
	api.Resource("/posts", ctrl)

	w := doRequest(r, http.MethodGet, "/api/posts")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "index" {
		t.Errorf("expected 'index', got %q", w.Body.String())
	}

	w = doRequest(r, http.MethodGet, "/api/posts/1")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// TC-16: APIResource on RouteGroup combines prefix
func TestAPIResource_OnGroup(t *testing.T) {
	r := newTestRouter()
	ctrl := &mockResourceController{}
	api := r.Group("/api")
	api.APIResource("/users", ctrl)

	w := doRequest(r, http.MethodGet, "/api/users")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	w = doRequest(r, http.MethodPut, "/api/users/5")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// TC-17: Named route — Name and Route generate correct URL
func TestNamedRoute_SingleParam(t *testing.T) {
	ResetNamedRoutes()
	Name("users.show", "/users/:id")
	result := Route("users.show", "42")
	if result != "/users/42" {
		t.Errorf("expected '/users/42', got %q", result)
	}
}

// TC-18: Named route — multiple parameters
func TestNamedRoute_MultipleParams(t *testing.T) {
	ResetNamedRoutes()
	Name("posts.comment", "/posts/:postId/comments/:commentId")
	result := Route("posts.comment", "5", "99")
	if result != "/posts/5/comments/99" {
		t.Errorf("expected '/posts/5/comments/99', got %q", result)
	}
}

// TC-19: Named route — unknown name returns "/"
func TestNamedRoute_UnknownName(t *testing.T) {
	ResetNamedRoutes()
	result := Route("nonexistent")
	if result != "/" {
		t.Errorf("expected '/', got %q", result)
	}
}

// TC-20: Named route — no parameters returns pattern as-is
func TestNamedRoute_NoParams(t *testing.T) {
	ResetNamedRoutes()
	Name("home", "/")
	result := Route("home")
	if result != "/" {
		t.Errorf("expected '/', got %q", result)
	}
}

// TC-21: Router implements http.Handler (ServeHTTP)
func TestServeHTTP_ImplementsHandler(t *testing.T) {
	r := newTestRouter()
	r.Get("/test", okHandler("handler"))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "handler" {
		t.Errorf("expected 'handler', got %q", w.Body.String())
	}

	// Compile-time check that Router satisfies http.Handler
	var _ http.Handler = r
}

// TC-22: Group middleware executes for group routes only
func TestGroup_MiddlewareIsolation(t *testing.T) {
	r := newTestRouter()

	// Root route — no middleware
	r.Get("/public", okHandler("public"))

	// Group with middleware that sets a header
	api := r.Group("/api", func(c *gin.Context) {
		c.Header("X-API", "true")
		c.Next()
	})
	api.Get("/data", okHandler("data"))

	// /api/data should have the X-API header
	w := doRequest(r, http.MethodGet, "/api/data")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Header().Get("X-API") != "true" {
		t.Error("expected X-API header on group route")
	}

	// /public should NOT have the X-API header
	w = doRequest(r, http.MethodGet, "/public")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Header().Get("X-API") != "" {
		t.Error("X-API header should not be present on public route")
	}
}

// TC-23: RouterProvider implements Provider interface — tested in providers_test.go

// TC-24: RouterProvider registers router as "router" in container — tested in providers_test.go

// TC-25: Route parameter extraction via Gin context
func TestRouteParam_Extraction(t *testing.T) {
	r := newTestRouter()
	r.Get("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(http.StatusOK, id)
	})

	w := doRequest(r, http.MethodGet, "/users/42")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "42" {
		t.Errorf("expected '42', got %q", w.Body.String())
	}
}

// TC-26: NoRoute registers a custom 404 handler
func TestNoRoute_CustomHandler(t *testing.T) {
	r := newTestRouter()
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "custom 404")
	})
	r.Get("/exists", okHandler("found"))

	w := doRequest(r, http.MethodGet, "/nonexistent")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
	if w.Body.String() != "custom 404" {
		t.Errorf("expected 'custom 404', got %q", w.Body.String())
	}
}

// TC-27: GlobalHandlers returns registered global middleware
func TestGlobalHandlers_ReturnsMiddleware(t *testing.T) {
	r := newTestRouter()
	if len(r.GlobalHandlers()) != 0 {
		t.Fatalf("expected 0 global handlers on new router, got %d", len(r.GlobalHandlers()))
	}

	r.Use(func(c *gin.Context) { c.Next() })
	r.Use(func(c *gin.Context) { c.Next() })

	if len(r.GlobalHandlers()) != 2 {
		t.Fatalf("expected 2 global handlers, got %d", len(r.GlobalHandlers()))
	}
}
