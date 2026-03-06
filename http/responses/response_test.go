package responses

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestSuccess_Returns200WithData(t *testing.T) {
	c, w := setupTestContext()

	data := map[string]string{"name": "RGo"}
	Success(c, data)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if !resp.Success {
		t.Fatal("expected success to be true")
	}
	if resp.Data == nil {
		t.Fatal("expected data to be present")
	}
}

func TestCreated_Returns201WithData(t *testing.T) {
	c, w := setupTestContext()

	data := map[string]string{"id": "1"}
	Created(c, data)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if !resp.Success {
		t.Fatal("expected success to be true")
	}
	if resp.Data == nil {
		t.Fatal("expected data to be present")
	}
}

func TestError_Returns404WithMessage(t *testing.T) {
	c, w := setupTestContext()

	Error(c, http.StatusNotFound, "not found")

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if resp.Success {
		t.Fatal("expected success to be false")
	}
	if resp.Error != "not found" {
		t.Fatalf("expected error 'not found', got '%s'", resp.Error)
	}
}

func TestError_Returns422WithMessage(t *testing.T) {
	c, w := setupTestContext()

	Error(c, http.StatusUnprocessableEntity, "validation failed")

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if resp.Success {
		t.Fatal("expected success to be false")
	}
	if resp.Error != "validation failed" {
		t.Fatalf("expected error 'validation failed', got '%s'", resp.Error)
	}
}

func TestPaginated_ReturnsDataWithMeta(t *testing.T) {
	c, w := setupTestContext()

	items := []string{"a", "b", "c"}
	Paginated(c, items, 1, 10, 25)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if !resp.Success {
		t.Fatal("expected success to be true")
	}
	if resp.Meta == nil {
		t.Fatal("expected meta to be present")
	}
	if resp.Meta.Page != 1 {
		t.Fatalf("expected page 1, got %d", resp.Meta.Page)
	}
	if resp.Meta.PerPage != 10 {
		t.Fatalf("expected per_page 10, got %d", resp.Meta.PerPage)
	}
	if resp.Meta.Total != 25 {
		t.Fatalf("expected total 25, got %d", resp.Meta.Total)
	}
	if resp.Meta.TotalPages != 3 {
		t.Fatalf("expected total_pages 3, got %d", resp.Meta.TotalPages)
	}
}

func TestPaginated_ExactDivision(t *testing.T) {
	c, w := setupTestContext()

	items := []string{"a", "b"}
	Paginated(c, items, 1, 10, 30)

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if resp.Meta.TotalPages != 3 {
		t.Fatalf("expected total_pages 3, got %d", resp.Meta.TotalPages)
	}
}

func TestSuccess_OmitsErrorField(t *testing.T) {
	c, w := setupTestContext()

	Success(c, map[string]string{"ok": "yes"})

	body := w.Body.String()
	if strings.Contains(body, `"error"`) {
		t.Fatalf("expected no 'error' field in JSON, got: %s", body)
	}
}

func TestError_OmitsDataField(t *testing.T) {
	c, w := setupTestContext()

	Error(c, http.StatusBadRequest, "bad")

	body := w.Body.String()
	if strings.Contains(body, `"data"`) {
		t.Fatalf("expected no 'data' field in JSON, got: %s", body)
	}
}
