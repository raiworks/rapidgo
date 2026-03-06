package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

// TC-01: Home returns welcome message with status 200
func TestHome_ReturnsWelcome(t *testing.T) {
	c, w := newTestContext("GET", "/")
	controllers.Home(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	m := bodyMap(t, w)
	if m["message"] != "Welcome to RGo" {
		t.Fatalf("expected message 'Welcome to RGo', got '%v'", m["message"])
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
	r := router.New()
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
