# Feature #40 — Testing Infrastructure: Design

## Package

`testing/testutil/testutil.go`

## Public API

```go
// NewTestRouter creates a Router in Gin test mode.
func NewTestRouter(t *testing.T) *router.Router

// NewTestDB opens an in-memory SQLite database for tests.
func NewTestDB(t *testing.T, models ...interface{}) *gorm.DB

// DoRequest performs an HTTP request against an http.Handler and returns the recorder.
func DoRequest(handler http.Handler, method, path string) *httptest.ResponseRecorder

// AssertStatus fails the test if the response code doesn't match.
func AssertStatus(t *testing.T, got, want int)

// AssertJSONKey fails if the JSON response body doesn't contain the expected key/value.
func AssertJSONKey(t *testing.T, body []byte, key, want string)
```

## File Layout

```
testing/testutil/
  testutil.go      — Utility functions
  testutil_test.go — Self-tests
```
