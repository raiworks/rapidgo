package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// T-011: KeyByIP returns client IP
func TestKeyByIP_ReturnsClientIP(t *testing.T) {
	fn := KeyByIP()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.RemoteAddr = "192.168.1.1:12345"

	key := fn(c)
	if key != "192.168.1.1" {
		t.Errorf("KeyByIP() = %q, want %q", key, "192.168.1.1")
	}
}

// T-012: KeyByUserID with context value
func TestKeyByUserID_WithContextValue(t *testing.T) {
	fn := KeyByUserID("userID")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Set("userID", "abc-123")

	key := fn(c)
	if key != "user:abc-123" {
		t.Errorf("KeyByUserID() = %q, want %q", key, "user:abc-123")
	}
}

// T-013: KeyByUserID without context value falls back to IP
func TestKeyByUserID_FallsBackToIP(t *testing.T) {
	fn := KeyByUserID("userID")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.RemoteAddr = "10.0.0.1:9999"

	key := fn(c)
	if key != "10.0.0.1" {
		t.Errorf("KeyByUserID() fallback = %q, want %q", key, "10.0.0.1")
	}
}

// T-014: KeyByHeader with header present
func TestKeyByHeader_WithHeader(t *testing.T) {
	fn := KeyByHeader("X-API-Key")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("X-API-Key", "my-key-123")

	key := fn(c)
	if key != "header:X-API-Key:my-key-123" {
		t.Errorf("KeyByHeader() = %q, want %q", key, "header:X-API-Key:my-key-123")
	}
}

// T-015: KeyByHeader without header falls back to IP
func TestKeyByHeader_FallsBackToIP(t *testing.T) {
	fn := KeyByHeader("X-API-Key")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.RemoteAddr = "172.16.0.1:8080"

	key := fn(c)
	if key != "172.16.0.1" {
		t.Errorf("KeyByHeader() fallback = %q, want %q", key, "172.16.0.1")
	}
}

// T-016: ParseRate valid format
func TestParseRate_ValidFormat(t *testing.T) {
	r, err := ParseRate("100-M")
	if err != nil {
		t.Fatalf("ParseRate(\"100-M\") error: %v", err)
	}
	if r.Limit != 100 {
		t.Errorf("rate.Limit = %d, want 100", r.Limit)
	}
}

// T-017: ParseRate invalid format
func TestParseRate_InvalidFormat(t *testing.T) {
	_, err := ParseRate("invalid")
	if err == nil {
		t.Fatal("ParseRate(\"invalid\") should return error")
	}
}

// KeyByBodyField with JSON body containing field
func TestKeyByBodyField_WithField(t *testing.T) {
	fn := KeyByBodyField("email")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"email":"test@example.com"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.RemoteAddr = "10.0.0.1:1234"

	key := fn(c)
	if key != "body:email:test@example.com" {
		t.Errorf("KeyByBodyField() = %q, want %q", key, "body:email:test@example.com")
	}
}

// KeyByBodyField without field falls back to IP
func TestKeyByBodyField_FallsBackToIP(t *testing.T) {
	fn := KeyByBodyField("email")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"test"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.RemoteAddr = "10.0.0.1:1234"

	key := fn(c)
	if key != "10.0.0.1" {
		t.Errorf("KeyByBodyField() fallback = %q, want %q", key, "10.0.0.1")
	}
}
