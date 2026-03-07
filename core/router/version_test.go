package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// T01
func TestVersion_ReturnsRouteGroup(t *testing.T) {
	r := newTestRouter()
	g := r.Version("v1")
	if g == nil {
		t.Fatal("expected non-nil RouteGroup")
	}
}

// T02
func TestVersion_PrefixesRoutes(t *testing.T) {
	r := newTestRouter()
	v1 := r.Version("v1")
	v1.Get("/users", okHandler("users-v1"))

	w := doRequest(r, http.MethodGet, "/api/v1/users")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "users-v1" {
		t.Fatalf("expected 'users-v1', got %q", w.Body.String())
	}
}

// T03
func TestVersion_MultipleVersions(t *testing.T) {
	r := newTestRouter()
	v1 := r.Version("v1")
	v2 := r.Version("v2")
	v1.Get("/items", okHandler("items-v1"))
	v2.Get("/items", okHandler("items-v2"))

	w1 := doRequest(r, http.MethodGet, "/api/v1/items")
	w2 := doRequest(r, http.MethodGet, "/api/v2/items")

	if w1.Body.String() != "items-v1" {
		t.Fatalf("v1: expected 'items-v1', got %q", w1.Body.String())
	}
	if w2.Body.String() != "items-v2" {
		t.Fatalf("v2: expected 'items-v2', got %q", w2.Body.String())
	}
}

// T04
func TestVersion_SupportsAllMethods(t *testing.T) {
	r := newTestRouter()
	v1 := r.Version("v1")
	v1.Get("/res", okHandler("get"))
	v1.Post("/res", okHandler("post"))
	v1.Put("/res", okHandler("put"))
	v1.Delete("/res", okHandler("delete"))

	tests := []struct {
		method string
		want   string
	}{
		{http.MethodGet, "get"},
		{http.MethodPost, "post"},
		{http.MethodPut, "put"},
		{http.MethodDelete, "delete"},
	}
	for _, tt := range tests {
		w := doRequest(r, tt.method, "/api/v1/res")
		if w.Code != http.StatusOK {
			t.Fatalf("%s: expected 200, got %d", tt.method, w.Code)
		}
		if w.Body.String() != tt.want {
			t.Fatalf("%s: expected %q, got %q", tt.method, tt.want, w.Body.String())
		}
	}
}

// T05
func TestVersion_NoDeprecationHeaders(t *testing.T) {
	r := newTestRouter()
	v2 := r.Version("v2")
	v2.Get("/ping", okHandler("pong"))

	w := doRequest(r, http.MethodGet, "/api/v2/ping")
	if w.Header().Get("Sunset") != "" {
		t.Fatal("non-deprecated version should not have Sunset header")
	}
	if w.Header().Get("X-API-Deprecated") != "" {
		t.Fatal("non-deprecated version should not have X-API-Deprecated header")
	}
}

// T06
func TestDeprecatedVersion_SetsSunsetHeader(t *testing.T) {
	r := newTestRouter()
	sunset := "Sat, 01 Jun 2026 00:00:00 GMT"
	v1 := r.DeprecatedVersion("v1", sunset)
	v1.Get("/data", okHandler("ok"))

	w := doRequest(r, http.MethodGet, "/api/v1/data")
	if w.Header().Get("Sunset") != sunset {
		t.Fatalf("expected Sunset %q, got %q", sunset, w.Header().Get("Sunset"))
	}
}

// T07
func TestDeprecatedVersion_SetsDeprecatedHeader(t *testing.T) {
	r := newTestRouter()
	v1 := r.DeprecatedVersion("v1", "Sat, 01 Jun 2026 00:00:00 GMT")
	v1.Get("/data", okHandler("ok"))

	w := doRequest(r, http.MethodGet, "/api/v1/data")
	if w.Header().Get("X-API-Deprecated") != "true" {
		t.Fatalf("expected X-API-Deprecated 'true', got %q", w.Header().Get("X-API-Deprecated"))
	}
}

// T08
func TestDeprecatedVersion_PrefixesRoutes(t *testing.T) {
	r := newTestRouter()
	v1 := r.DeprecatedVersion("v1", "Sat, 01 Jun 2026 00:00:00 GMT")
	v1.Get("/users", okHandler("users-v1"))

	w := doRequest(r, http.MethodGet, "/api/v1/users")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "users-v1" {
		t.Fatalf("expected 'users-v1', got %q", w.Body.String())
	}
}

// T09
func TestDeprecatedVersion_HeadersOnAllRoutes(t *testing.T) {
	r := newTestRouter()
	sunset := "Sat, 01 Jun 2026 00:00:00 GMT"
	v1 := r.DeprecatedVersion("v1", sunset)
	v1.Get("/a", okHandler("a"))
	v1.Post("/b", okHandler("b"))

	for _, tc := range []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/a"},
		{http.MethodPost, "/api/v1/b"},
	} {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(tc.method, tc.path, nil)
		r.ServeHTTP(w, req)

		if w.Header().Get("Sunset") != sunset {
			t.Fatalf("%s %s: missing Sunset header", tc.method, tc.path)
		}
		if w.Header().Get("X-API-Deprecated") != "true" {
			t.Fatalf("%s %s: missing X-API-Deprecated header", tc.method, tc.path)
		}
	}
}

// T10
func TestVersion_NestedGroup(t *testing.T) {
	r := newTestRouter()
	v1 := r.Version("v1")
	admin := v1.Group("/admin")
	admin.Get("/dashboard", okHandler("admin-dashboard"))

	w := doRequest(r, http.MethodGet, "/api/v1/admin/dashboard")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "admin-dashboard" {
		t.Fatalf("expected 'admin-dashboard', got %q", w.Body.String())
	}
}
