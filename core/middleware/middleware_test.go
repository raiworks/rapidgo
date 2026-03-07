package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/RAiWorks/RapidGo/core/auth"
	"github.com/RAiWorks/RapidGo/core/errors"
	"github.com/RAiWorks/RapidGo/core/session"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// newTestEngine creates a Gin engine in test mode with no default middleware.
func newTestEngine() *gin.Engine {
	return gin.New()
}

// doRequest performs a request against a Gin engine and returns the recorder.
func doRequest(e *gin.Engine, method, path string, headers ...http.Header) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	if len(headers) > 0 {
		req.Header = headers[0]
	}
	e.ServeHTTP(w, req)
	return w
}

// --- Registry Tests ---

// TC-01: RegisterAlias and Resolve round-trip
func TestRegisterAlias_AndResolve(t *testing.T) {
	ResetRegistry()
	handler := func(c *gin.Context) {}
	RegisterAlias("test", handler)

	resolved := Resolve("test")
	if resolved == nil {
		t.Fatal("Resolve returned nil for registered alias")
	}
}

// TC-02: Resolve panics on unknown alias
func TestResolve_PanicsOnUnknown(t *testing.T) {
	ResetRegistry()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Resolve did not panic for unknown alias")
		}
		msg, ok := r.(string)
		if !ok || msg != "middleware not found: nonexistent" {
			t.Fatalf("unexpected panic message: %v", r)
		}
	}()

	Resolve("nonexistent")
}

// TC-03: RegisterGroup and ResolveGroup round-trip
func TestRegisterGroup_AndResolveGroup(t *testing.T) {
	ResetRegistry()
	h1 := func(c *gin.Context) {}
	h2 := func(c *gin.Context) {}
	RegisterGroup("web", h1, h2)

	group := ResolveGroup("web")
	if len(group) != 2 {
		t.Fatalf("expected group length 2, got %d", len(group))
	}
}

// TC-04: ResolveGroup returns nil for unknown group
func TestResolveGroup_ReturnsNilOnUnknown(t *testing.T) {
	ResetRegistry()

	group := ResolveGroup("nonexistent")
	if group != nil {
		t.Fatalf("expected nil for unknown group, got %v", group)
	}
}

// TC-05: ResetRegistry clears all entries
func TestResetRegistry_ClearsAll(t *testing.T) {
	handler := func(c *gin.Context) {}
	RegisterAlias("a", handler)
	RegisterGroup("g", handler)

	ResetRegistry()

	if ResolveGroup("g") != nil {
		t.Fatal("group should be nil after reset")
	}

	defer func() {
		if recover() == nil {
			t.Fatal("Resolve should panic after reset")
		}
	}()
	Resolve("a")
}

// --- Recovery Tests ---

// TC-06: Recovery catches panic and returns 500
func TestRecovery_CatchesPanic(t *testing.T) {
	e := newTestEngine()
	e.Use(Recovery())
	e.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	w := doRequest(e, http.MethodGet, "/panic")
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if body["error"] != "internal server error" {
		t.Fatalf("unexpected error message: %v", body["error"])
	}
}

// TC-07: Recovery passes through normal requests
func TestRecovery_PassesNormalRequest(t *testing.T) {
	e := newTestEngine()
	e.Use(Recovery())
	e.GET("/ok", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := doRequest(e, http.MethodGet, "/ok")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

// --- RequestID Tests ---

// TC-08: RequestID generates UUID when no header present
func TestRequestID_GeneratesUUID(t *testing.T) {
	e := newTestEngine()
	e.Use(RequestID())
	e.GET("/id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"id": c.GetString("request_id")})
	})

	w := doRequest(e, http.MethodGet, "/id")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	headerID := w.Header().Get("X-Request-ID")
	if headerID == "" {
		t.Fatal("X-Request-ID header is empty")
	}

	// Verify UUID format
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	if !uuidPattern.MatchString(headerID) {
		t.Fatalf("X-Request-ID is not valid UUID v4: %s", headerID)
	}

	// Verify body contains same ID
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body["id"] != headerID {
		t.Fatalf("body ID %q != header ID %q", body["id"], headerID)
	}
}

// TC-09: RequestID preserves incoming X-Request-ID
func TestRequestID_PreservesExisting(t *testing.T) {
	e := newTestEngine()
	e.Use(RequestID())
	e.GET("/id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"id": c.GetString("request_id")})
	})

	h := http.Header{}
	h.Set("X-Request-ID", "my-trace-id-123")
	w := doRequest(e, http.MethodGet, "/id", h)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	headerID := w.Header().Get("X-Request-ID")
	if headerID != "my-trace-id-123" {
		t.Fatalf("expected X-Request-ID 'my-trace-id-123', got %q", headerID)
	}

	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body["id"] != "my-trace-id-123" {
		t.Fatalf("body ID should be 'my-trace-id-123', got %q", body["id"])
	}
}

// --- CORS Tests ---

// TC-10: CORS sets default headers
func TestCORS_DefaultHeaders(t *testing.T) {
	e := newTestEngine()
	e.Use(CORS())
	e.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := doRequest(e, http.MethodGet, "/test")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "*" {
		t.Fatalf("expected Allow-Origin '*', got %q", origin)
	}

	methods := w.Header().Get("Access-Control-Allow-Methods")
	if methods == "" {
		t.Fatal("Access-Control-Allow-Methods is empty")
	}

	headers := w.Header().Get("Access-Control-Allow-Headers")
	if headers == "" {
		t.Fatal("Access-Control-Allow-Headers is empty")
	}
}

// TC-11: CORS handles preflight OPTIONS with 204
func TestCORS_PreflightOptions(t *testing.T) {
	e := newTestEngine()
	e.Use(CORS())
	e.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := doRequest(e, http.MethodOptions, "/test")
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "*" {
		t.Fatalf("expected Allow-Origin '*', got %q", origin)
	}
}

// TC-12: CORS accepts custom configuration
func TestCORS_CustomConfig(t *testing.T) {
	cfg := CORSConfig{
		AllowOrigins: []string{"https://example.com"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Content-Type"},
		MaxAge:       3600,
	}

	e := newTestEngine()
	e.Use(CORS(cfg))
	e.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := doRequest(e, http.MethodGet, "/test")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "https://example.com" {
		t.Fatalf("expected Allow-Origin 'https://example.com', got %q", origin)
	}

	maxAge := w.Header().Get("Access-Control-Max-Age")
	if maxAge != "3600" {
		t.Fatalf("expected Max-Age '3600', got %q", maxAge)
	}
}

// TC-26: Default AllowHeaders includes X-CSRF-Token
func TestCORS_DefaultIncludesCSRFToken(t *testing.T) {
	e := newTestEngine()
	e.Use(CORS())
	e.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := doRequest(e, "GET", "/test")
	headers := w.Header().Get("Access-Control-Allow-Headers")
	if !strings.Contains(headers, "X-CSRF-Token") {
		t.Fatalf("expected Allow-Headers to contain X-CSRF-Token, got %q", headers)
	}
}

// TC-27: AllowCredentials header set by default
func TestCORS_DefaultCredentials(t *testing.T) {
	e := newTestEngine()
	e.Use(CORS())
	e.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := doRequest(e, "GET", "/test")
	creds := w.Header().Get("Access-Control-Allow-Credentials")
	if creds != "true" {
		t.Fatalf("expected Allow-Credentials 'true', got %q", creds)
	}
}

// TC-28: ExposeHeaders set by default
func TestCORS_DefaultExposeHeaders(t *testing.T) {
	e := newTestEngine()
	e.Use(CORS())
	e.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := doRequest(e, "GET", "/test")
	expose := w.Header().Get("Access-Control-Expose-Headers")
	if !strings.Contains(expose, "Content-Length") || !strings.Contains(expose, "X-Request-ID") {
		t.Fatalf("expected Expose-Headers with Content-Length and X-Request-ID, got %q", expose)
	}
}

// TC-29: CORS_ALLOWED_ORIGINS env overrides default
func TestCORS_EnvOverridesOrigins(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://example.com")

	e := newTestEngine()
	e.Use(CORS())
	e.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := doRequest(e, "GET", "/test")
	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "https://example.com" {
		t.Fatalf("expected 'https://example.com', got %q", origin)
	}
}

// TC-30: Custom config can disable credentials
func TestCORS_CustomDisableCredentials(t *testing.T) {
	cfg := CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: false,
		MaxAge:           3600,
	}

	e := newTestEngine()
	e.Use(CORS(cfg))
	e.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := doRequest(e, "GET", "/test")
	creds := w.Header().Get("Access-Control-Allow-Credentials")
	if creds != "" {
		t.Fatalf("expected no Allow-Credentials header, got %q", creds)
	}
}

// TC-31: Preflight with new headers returns 204 with credentials + expose
func TestCORS_PreflightWithNewHeaders(t *testing.T) {
	e := newTestEngine()
	e.Use(CORS())
	e.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := doRequest(e, "OPTIONS", "/test")
	if w.Code != 204 {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Fatal("expected credentials header on preflight")
	}
	if w.Header().Get("Access-Control-Expose-Headers") == "" {
		t.Fatal("expected expose headers on preflight")
	}
}

// --- ErrorHandler Tests ---

// TC-13: ErrorHandler formats AppError as JSON
func TestErrorHandler_FormatsAppError(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_DEBUG", "false")

	e := newTestEngine()
	e.Use(ErrorHandler())
	e.GET("/err", func(c *gin.Context) {
		_ = c.Error(errors.NotFound("user not found"))
		c.Abort()
	})

	w := doRequest(e, http.MethodGet, "/err")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}

	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "user not found" {
		t.Fatalf("unexpected error message: %v", body["error"])
	}
}

// TC-14: ErrorHandler wraps generic error as 500
func TestErrorHandler_WrapsGenericError(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_DEBUG", "false")

	e := newTestEngine()
	e.Use(ErrorHandler())
	e.GET("/err", func(c *gin.Context) {
		_ = c.Error(fmt.Errorf("something broke"))
		c.Abort()
	})

	w := doRequest(e, http.MethodGet, "/err")
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "internal server error" {
		t.Fatalf("unexpected error message: %v", body["error"])
	}
}

// TC-15: ErrorHandler is no-op when no errors
func TestErrorHandler_NoOpWhenNoErrors(t *testing.T) {
	e := newTestEngine()
	e.Use(ErrorHandler())
	e.GET("/ok", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := doRequest(e, http.MethodGet, "/ok")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body["status"] != "ok" {
		t.Fatalf("expected status 'ok', got %v", body["status"])
	}
}

// --- Integration Tests ---

// TC-18: Middleware chain executes in correct order
func TestMiddlewareChain_ExecutionOrder(t *testing.T) {
	e := newTestEngine()
	e.Use(Recovery())
	e.Use(RequestID())
	e.GET("/chain", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"id": c.GetString("request_id")})
	})

	w := doRequest(e, http.MethodGet, "/chain")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	headerID := w.Header().Get("X-Request-ID")
	if headerID == "" {
		t.Fatal("X-Request-ID header missing in chained middleware")
	}
}

// TC-19: UUID format validation
func TestGenerateUUID_Format(t *testing.T) {
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

	ids := make([]string, 10)
	for i := range ids {
		ids[i] = generateUUID()
		if len(ids[i]) != 36 {
			t.Fatalf("UUID length should be 36, got %d: %s", len(ids[i]), ids[i])
		}
		if !uuidPattern.MatchString(ids[i]) {
			t.Fatalf("invalid UUID v4 format: %s", ids[i])
		}
	}

	// Check uniqueness (at least first two)
	if ids[0] == ids[1] {
		t.Fatal("two UUIDs are identical — should be unique")
	}
}

// --- Session Middleware ---

func TestSessionMiddleware_SetsSessionData(t *testing.T) {
	store := session.NewMemoryStore()
	mgr := &session.Manager{
		Store:      store,
		CookieName: "test_session",
		Lifetime:   120 * 60_000_000_000, // 120 minutes in nanoseconds
		Path:       "/",
		HTTPOnly:   true,
		SameSite:   http.SameSiteLaxMode,
	}

	e := newTestEngine()
	e.Use(SessionMiddleware(mgr))
	e.GET("/test", func(c *gin.Context) {
		sid, exists := c.Get("session_id")
		if !exists || sid == "" {
			c.String(500, "no session_id")
			return
		}
		_, exists = c.Get("session")
		if !exists {
			c.String(500, "no session data")
			return
		}
		c.String(200, "ok")
	})

	w := doRequest(e, "GET", "/test")
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if w.Body.String() != "ok" {
		t.Fatalf("expected 'ok', got %q", w.Body.String())
	}
}

// --- Auth Middleware ---

// TC-10: AuthMiddleware rejects request without Authorization header
func TestAuthMiddleware_RejectsMissingHeader(t *testing.T) {
	e := newTestEngine()
	e.Use(AuthMiddleware())
	e.GET("/protected", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := doRequest(e, "GET", "/protected")
	if w.Code != 401 {
		t.Fatalf("expected 401, got %d", w.Code)
	}

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] == nil {
		t.Fatal("expected error message in response body")
	}
}

// TC-11: AuthMiddleware rejects request with invalid token
func TestAuthMiddleware_RejectsInvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")

	e := newTestEngine()
	e.Use(AuthMiddleware())
	e.GET("/protected", func(c *gin.Context) {
		c.String(200, "ok")
	})

	headers := http.Header{}
	headers.Set("Authorization", "Bearer invalid-token-value")
	w := doRequest(e, "GET", "/protected", headers)
	if w.Code != 401 {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

// TC-12: AuthMiddleware sets user_id on valid token
func TestAuthMiddleware_SetsUserID(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-for-auth")
	t.Setenv("JWT_EXPIRY", "3600")

	token, err := auth.GenerateToken(99)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	e := newTestEngine()
	e.Use(AuthMiddleware())
	e.GET("/protected", func(c *gin.Context) {
		uid, exists := c.Get("user_id")
		if !exists {
			c.String(500, "no user_id")
			return
		}
		c.String(200, fmt.Sprintf("user_id=%v", uid))
	})

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)
	w := doRequest(e, "GET", "/protected", headers)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if w.Body.String() != "user_id=99" {
		t.Fatalf("expected 'user_id=99', got %q", w.Body.String())
	}
}

// TC-13: AuthMiddleware rejects non-Bearer auth scheme
func TestAuthMiddleware_RejectsNonBearerScheme(t *testing.T) {
	e := newTestEngine()
	e.Use(AuthMiddleware())
	e.GET("/protected", func(c *gin.Context) {
		c.String(200, "ok")
	})

	headers := http.Header{}
	headers.Set("Authorization", "Basic abc123")
	w := doRequest(e, "GET", "/protected", headers)
	if w.Code != 401 {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

// TC-14: Auth alias is resolvable after registration
func TestAuthAlias_IsResolvable(t *testing.T) {
	ResetRegistry()
	RegisterAlias("auth", AuthMiddleware())

	resolved := Resolve("auth")
	if resolved == nil {
		t.Fatal("expected auth alias to resolve to a handler")
	}
}

// --- CSRF Middleware ---

// fakeSession injects a session map into the context, simulating SessionMiddleware.
func fakeSession(data map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("session", data)
		c.Next()
	}
}

// TC-15: GET request generates token and sets it in context
func TestCSRF_GETGeneratesToken(t *testing.T) {
	sessData := map[string]interface{}{}
	e := newTestEngine()
	e.Use(fakeSession(sessData))
	e.Use(CSRFMiddleware())

	var gotToken string
	e.GET("/form", func(c *gin.Context) {
		tok, _ := c.Get("csrf_token")
		gotToken = tok.(string)
		c.String(200, "ok")
	})

	w := doRequest(e, "GET", "/form")
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if len(gotToken) != 64 {
		t.Fatalf("expected 64-char hex token, got %d chars", len(gotToken))
	}
}

// TC-16: POST with valid form token passes
func TestCSRF_POSTValidFormToken(t *testing.T) {
	sessData := map[string]interface{}{"_csrf_token": "known-token-value"}
	e := newTestEngine()
	e.Use(fakeSession(sessData))
	e.Use(CSRFMiddleware())
	e.POST("/submit", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	body := strings.NewReader("_csrf_token=known-token-value")
	req := httptest.NewRequest("POST", "/submit", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	e.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// TC-17: POST with valid header token passes
func TestCSRF_POSTValidHeaderToken(t *testing.T) {
	sessData := map[string]interface{}{"_csrf_token": "header-token-value"}
	e := newTestEngine()
	e.Use(fakeSession(sessData))
	e.Use(CSRFMiddleware())
	e.POST("/submit", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/submit", nil)
	req.Header.Set("X-CSRF-Token", "header-token-value")
	e.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// TC-18: POST with missing token returns 403
func TestCSRF_POSTMissingToken(t *testing.T) {
	sessData := map[string]interface{}{"_csrf_token": "some-token"}
	e := newTestEngine()
	e.Use(fakeSession(sessData))
	e.Use(CSRFMiddleware())
	e.POST("/submit", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/submit", nil)
	e.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Fatalf("expected 403, got %d", w.Code)
	}
	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "CSRF token mismatch" {
		t.Fatalf("expected 'CSRF token mismatch', got %v", body["error"])
	}
}

// TC-19: POST with wrong token returns 403
func TestCSRF_POSTWrongToken(t *testing.T) {
	sessData := map[string]interface{}{"_csrf_token": "correct-token"}
	e := newTestEngine()
	e.Use(fakeSession(sessData))
	e.Use(CSRFMiddleware())
	e.POST("/submit", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	body := strings.NewReader("_csrf_token=wrong-token")
	req := httptest.NewRequest("POST", "/submit", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	e.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

// TC-20: HEAD request skips validation
func TestCSRF_HEADSkips(t *testing.T) {
	sessData := map[string]interface{}{"_csrf_token": "some-token"}
	e := newTestEngine()
	e.Use(fakeSession(sessData))
	e.Use(CSRFMiddleware())
	e.HEAD("/ping", func(c *gin.Context) {
		c.Status(200)
	})

	w := doRequest(e, "HEAD", "/ping")
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

// TC-21: OPTIONS request skips validation
func TestCSRF_OPTIONSSkips(t *testing.T) {
	sessData := map[string]interface{}{"_csrf_token": "some-token"}
	e := newTestEngine()
	e.Use(fakeSession(sessData))
	e.Use(CSRFMiddleware())
	e.OPTIONS("/ping", func(c *gin.Context) {
		c.Status(200)
	})

	w := doRequest(e, "OPTIONS", "/ping")
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

// TC-22: PUT with valid token passes
func TestCSRF_PUTValidToken(t *testing.T) {
	sessData := map[string]interface{}{"_csrf_token": "put-token"}
	e := newTestEngine()
	e.Use(fakeSession(sessData))
	e.Use(CSRFMiddleware())
	e.PUT("/update", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/update", nil)
	req.Header.Set("X-CSRF-Token", "put-token")
	e.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// TC-23: DELETE with missing token returns 403
func TestCSRF_DELETEMissingToken(t *testing.T) {
	sessData := map[string]interface{}{"_csrf_token": "del-token"}
	e := newTestEngine()
	e.Use(fakeSession(sessData))
	e.Use(CSRFMiddleware())
	e.DELETE("/remove", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := doRequest(e, "DELETE", "/remove")
	if w.Code != 403 {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

// TC-24: Token persists across requests with same session
func TestCSRF_TokenPersists(t *testing.T) {
	sessData := map[string]interface{}{}
	e := newTestEngine()
	e.Use(fakeSession(sessData))
	e.Use(CSRFMiddleware())

	var firstToken, secondToken string
	e.GET("/a", func(c *gin.Context) {
		tok, _ := c.Get("csrf_token")
		firstToken = tok.(string)
		c.String(200, "ok")
	})
	e.GET("/b", func(c *gin.Context) {
		tok, _ := c.Get("csrf_token")
		secondToken = tok.(string)
		c.String(200, "ok")
	})

	doRequest(e, "GET", "/a")
	doRequest(e, "GET", "/b")

	if firstToken != secondToken {
		t.Fatalf("expected same token, got %q and %q", firstToken, secondToken)
	}
}

// TC-25: "csrf" alias is resolvable
func TestCSRFAlias_IsResolvable(t *testing.T) {
	ResetRegistry()
	RegisterAlias("csrf", CSRFMiddleware())

	resolved := Resolve("csrf")
	if resolved == nil {
		t.Fatal("expected csrf alias to resolve to a handler")
	}
}

// --- Rate Limiting Tests ---

// TC-32: Default rate limit allows requests within limit
func TestRateLimit_AllowsWithinLimit(t *testing.T) {
	t.Setenv("RATE_LIMIT", "5-M")
	e := newTestEngine()
	e.Use(RateLimitMiddleware())
	e.GET("/ok", func(c *gin.Context) { c.String(200, "ok") })

	w := doRequest(e, "GET", "/ok")
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

// TC-33: X-RateLimit-Limit header present
func TestRateLimit_LimitHeaderPresent(t *testing.T) {
	t.Setenv("RATE_LIMIT", "5-M")
	e := newTestEngine()
	e.Use(RateLimitMiddleware())
	e.GET("/ok", func(c *gin.Context) { c.String(200, "ok") })

	w := doRequest(e, "GET", "/ok")
	limit := w.Header().Get("X-RateLimit-Limit")
	if limit != "5" {
		t.Fatalf("expected X-RateLimit-Limit=5, got %q", limit)
	}
}

// TC-34: X-RateLimit-Remaining decrements
func TestRateLimit_RemainingDecrements(t *testing.T) {
	t.Setenv("RATE_LIMIT", "5-M")
	e := newTestEngine()
	e.Use(RateLimitMiddleware())
	e.GET("/ok", func(c *gin.Context) { c.String(200, "ok") })

	w1 := doRequest(e, "GET", "/ok")
	r1 := w1.Header().Get("X-RateLimit-Remaining")

	w2 := doRequest(e, "GET", "/ok")
	r2 := w2.Header().Get("X-RateLimit-Remaining")

	if r1 <= r2 {
		t.Fatalf("expected remaining to decrement: first=%s, second=%s", r1, r2)
	}
}

// TC-35: Requests exceeding limit return 429
func TestRateLimit_ExceedReturns429(t *testing.T) {
	t.Setenv("RATE_LIMIT", "2-M")
	e := newTestEngine()
	e.Use(RateLimitMiddleware())
	e.GET("/ok", func(c *gin.Context) { c.String(200, "ok") })

	doRequest(e, "GET", "/ok") // 1
	doRequest(e, "GET", "/ok") // 2
	w := doRequest(e, "GET", "/ok") // 3 — should be rejected

	if w.Code != 429 {
		t.Fatalf("expected 429, got %d", w.Code)
	}
}

// TC-36: Custom RATE_LIMIT env var is respected
func TestRateLimit_CustomEnv(t *testing.T) {
	t.Setenv("RATE_LIMIT", "10-M")
	e := newTestEngine()
	e.Use(RateLimitMiddleware())
	e.GET("/ok", func(c *gin.Context) { c.String(200, "ok") })

	w := doRequest(e, "GET", "/ok")
	limit := w.Header().Get("X-RateLimit-Limit")
	if limit != "10" {
		t.Fatalf("expected X-RateLimit-Limit=10, got %q", limit)
	}
}

// TC-37: Middleware alias "ratelimit" resolves
func TestRateLimitAlias_IsResolvable(t *testing.T) {
	ResetRegistry()
	RegisterAlias("ratelimit", RateLimitMiddleware())

	resolved := Resolve("ratelimit")
	if resolved == nil {
		t.Fatal("expected ratelimit alias to resolve to a handler")
	}
}

// --- AdminOnly Middleware Tests ---

// TC-38: AdminOnly allows request when role is "admin".
func TestAdminOnly_AllowsAdmin(t *testing.T) {
	e := newTestEngine()
	e.Use(func(c *gin.Context) {
		c.Set("role", "admin")
		c.Next()
	})
	e.Use(AdminOnly())
	e.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := doRequest(e, http.MethodGet, "/admin")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

// TC-39: AdminOnly blocks request when role is "user".
func TestAdminOnly_BlocksUser(t *testing.T) {
	e := newTestEngine()
	e.Use(func(c *gin.Context) {
		c.Set("role", "user")
		c.Next()
	})
	e.Use(AdminOnly())
	e.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := doRequest(e, http.MethodGet, "/admin")
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}

	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "admin access required" {
		t.Fatalf("expected 'admin access required', got %v", body["error"])
	}
}

// TC-40: AdminOnly blocks request when role is not set.
func TestAdminOnly_BlocksNoRole(t *testing.T) {
	e := newTestEngine()
	e.Use(AdminOnly())
	e.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := doRequest(e, http.MethodGet, "/admin")
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}

	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "admin access required" {
		t.Fatalf("expected 'admin access required', got %v", body["error"])
	}
}
