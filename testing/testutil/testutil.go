package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RAiWorks/RGo/core/router"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// NewTestRouter creates a Router in Gin test mode.
func NewTestRouter(t *testing.T) *router.Router {
	t.Helper()
	t.Setenv("APP_ENV", "testing")
	return router.New()
}

// NewTestDB opens an in-memory SQLite database and auto-migrates
// any provided models.
func NewTestDB(t *testing.T, models ...interface{}) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("testutil.NewTestDB: %v", err)
	}
	if len(models) > 0 {
		if err := db.AutoMigrate(models...); err != nil {
			t.Fatalf("testutil.NewTestDB auto-migrate: %v", err)
		}
	}
	return db
}

// DoRequest performs an HTTP request against the handler and returns the recorder.
func DoRequest(handler http.Handler, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	handler.ServeHTTP(w, req)
	return w
}

// AssertStatus fails if got != want.
func AssertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
}

// AssertJSONKey fails if the JSON body doesn't contain key with the expected string value.
func AssertJSONKey(t *testing.T, body []byte, key, want string) {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatalf("AssertJSONKey: unmarshal: %v", err)
	}
	got, ok := m[key]
	if !ok {
		t.Fatalf("AssertJSONKey: key %q not found in JSON", key)
	}
	if s, ok := got.(string); !ok || s != want {
		t.Fatalf("AssertJSONKey: %q = %v, want %q", key, got, want)
	}
}
