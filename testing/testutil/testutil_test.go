package testutil

import (
	"net/http"
	"testing"
)

// TC-01: NewTestRouter returns a usable router.
func TestNewTestRouter(t *testing.T) {
	r := NewTestRouter(t)
	if r == nil {
		t.Fatal("expected non-nil router")
	}
}

// TC-02: NewTestDB returns a working GORM DB and auto-migrates models.
func TestNewTestDB_AutoMigrate(t *testing.T) {
	type Item struct {
		ID   uint
		Name string
	}
	db := NewTestDB(t, &Item{})

	db.Create(&Item{Name: "test"})

	var count int64
	db.Model(&Item{}).Count(&count)
	if count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}
}

// TC-03: DoRequest returns correct response.
func TestDoRequest(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"msg":"hi"}`))
	})

	w := DoRequest(mux, http.MethodGet, "/hello")

	AssertStatus(t, w.Code, http.StatusOK)
	AssertJSONKey(t, w.Body.Bytes(), "msg", "hi")
}

// TC-04: AssertStatus passes on match (implicitly — no failure means pass).
func TestAssertStatus_Pass(t *testing.T) {
	AssertStatus(t, 200, 200)
}

// TC-05: AssertJSONKey validates key/value.
func TestAssertJSONKey_Valid(t *testing.T) {
	body := []byte(`{"status":"ok","count":"5"}`)
	AssertJSONKey(t, body, "status", "ok")
	AssertJSONKey(t, body, "count", "5")
}
