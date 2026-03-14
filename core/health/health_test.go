package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raiworks/rapidgo/v2/core/router"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupRouter(t *testing.T, db *gorm.DB) *router.Router {
	t.Helper()
	t.Setenv("APP_ENV", "testing")
	r := router.New()
	Routes(r, func() *gorm.DB { return db })
	return r
}

func openDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	return db
}

// TC-01: GET /health returns 200 with {"status":"ok"}
func TestLiveness_ReturnsOK(t *testing.T) {
	r := setupRouter(t, openDB(t))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /health status = %d, want %d", w.Code, http.StatusOK)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status = %q, want %q", body["status"], "ok")
	}
}

// TC-02: GET /health/ready returns 200 with live DB
func TestReadiness_WithLiveDB(t *testing.T) {
	r := setupRouter(t, openDB(t))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /health/ready status = %d, want %d", w.Code, http.StatusOK)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["status"] != "ready" {
		t.Fatalf("status = %q, want %q", body["status"], "ready")
	}
	if body["db"] != "connected" {
		t.Fatalf("db = %q, want %q", body["db"], "connected")
	}
}

// TC-03: GET /health/ready returns 503 when DB is closed
func TestReadiness_WithClosedDB(t *testing.T) {
	db := openDB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB(): %v", err)
	}
	sqlDB.Close()

	r := setupRouter(t, db)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("GET /health/ready status = %d, want %d", w.Code, http.StatusServiceUnavailable)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["status"] != "error" {
		t.Fatalf("status = %q, want %q", body["status"], "error")
	}
	if body["db"] == "" {
		t.Fatal("expected non-empty db error message")
	}
}

// TC-04: GET /health includes version when provided
func TestLiveness_WithVersion(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	r := router.New()
	Routes(r, func() *gorm.DB { return openDB(t) }, "2.4.0")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /health status = %d, want %d", w.Code, http.StatusOK)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status = %q, want %q", body["status"], "ok")
	}
	if body["version"] != "2.4.0" {
		t.Fatalf("version = %q, want %q", body["version"], "2.4.0")
	}
}

// TC-05: GET /health without version has no version field (backward compat)
func TestLiveness_WithoutVersion(t *testing.T) {
	r := setupRouter(t, openDB(t))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /health status = %d, want %d", w.Code, http.StatusOK)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, exists := body["version"]; exists {
		t.Fatal("expected no 'version' field when version not provided")
	}
}

// TC-06: GET /health/ready unaffected by version arg
func TestReadiness_UnaffectedByVersion(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	r := router.New()
	Routes(r, func() *gorm.DB { return openDB(t) }, "2.4.0")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /health/ready status = %d, want %d", w.Code, http.StatusOK)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, exists := body["version"]; exists {
		t.Fatal("readiness endpoint should not include version")
	}
}
