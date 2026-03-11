package database

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// paginationItem is a minimal model for pagination tests.
type paginationItem struct {
	ID   uint   `gorm:"primaryKey"`
	Name string
}

func (paginationItem) TableName() string { return "items" }

func setupPaginationDB(t *testing.T, count int) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&paginationItem{}); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}

	for i := 1; i <= count; i++ {
		db.Create(&paginationItem{ID: uint(i), Name: fmt.Sprintf("item-%d", i)})
	}
	return db
}

// ─── Page-based Pagination ───────────────────────────────────────────────────

// T-016: Correct offset calculation
func TestPaginate_OffsetCalculation(t *testing.T) {
	db := setupPaginationDB(t, 10)

	var items []paginationItem
	res, err := Paginate(db.Model(&paginationItem{}), 2, 3, &items)
	if err != nil {
		t.Fatalf("Paginate error: %v", err)
	}

	if len(items) != 3 {
		t.Fatalf("got %d items, want 3", len(items))
	}
	if items[0].ID != 4 {
		t.Errorf("first item ID = %d, want 4 (offset 3)", items[0].ID)
	}
	if res.Page != 2 {
		t.Errorf("Page = %d, want 2", res.Page)
	}
	if res.Total != 10 {
		t.Errorf("Total = %d, want 10", res.Total)
	}
	if res.TotalPages != 4 { // ceil(10/3) = 4
		t.Errorf("TotalPages = %d, want 4", res.TotalPages)
	}
}

// T-017: Out-of-range page returns empty slice
func TestPaginate_OutOfRange(t *testing.T) {
	db := setupPaginationDB(t, 5)

	var items []paginationItem
	res, err := Paginate(db.Model(&paginationItem{}), 100, 5, &items)
	if err != nil {
		t.Fatalf("Paginate error: %v", err)
	}

	if len(items) != 0 {
		t.Errorf("expected empty slice, got %d items", len(items))
	}
	if res.Total != 5 {
		t.Errorf("Total = %d, want 5", res.Total)
	}
}

// T-018: PerPage clamping (0→15, 200→100, negative→15)
func TestPaginate_PerPageClamping(t *testing.T) {
	db := setupPaginationDB(t, 20)

	tests := []struct {
		input int
		want  int
	}{
		{0, 15},
		{-1, 15},
		{200, 100},
		{5, 5},
	}

	for _, tt := range tests {
		var items []paginationItem
		res, err := Paginate(db.Model(&paginationItem{}), 1, tt.input, &items)
		if err != nil {
			t.Fatalf("Paginate(perPage=%d) error: %v", tt.input, err)
		}
		if res.PerPage != tt.want {
			t.Errorf("Paginate(perPage=%d): PerPage = %d, want %d", tt.input, res.PerPage, tt.want)
		}
	}
}

// T-018b: Page clamping (0→1, negative→1)
func TestPaginate_PageClamping(t *testing.T) {
	db := setupPaginationDB(t, 5)

	var items []paginationItem
	res, err := Paginate(db.Model(&paginationItem{}), -1, 5, &items)
	if err != nil {
		t.Fatalf("Paginate error: %v", err)
	}
	if res.Page != 1 {
		t.Errorf("Page = %d, want 1 (clamped)", res.Page)
	}
	if items[0].ID != 1 {
		t.Errorf("first item ID = %d, want 1", items[0].ID)
	}
}

// ─── Cursor-based Pagination ─────────────────────────────────────────────────

// T-019: First page (no cursor)
func TestCursorPaginate_FirstPage(t *testing.T) {
	db := setupPaginationDB(t, 10)

	var items []paginationItem
	res, err := CursorPaginate(db.Model(&paginationItem{}), "", "id", 3, "next", &items)
	if err != nil {
		t.Fatalf("CursorPaginate error: %v", err)
	}

	if len(items) != 3 {
		t.Fatalf("got %d items, want 3", len(items))
	}
	if items[0].ID != 1 || items[2].ID != 3 {
		t.Errorf("items = [%d..%d], want [1..3]", items[0].ID, items[2].ID)
	}
	if !res.HasMore {
		t.Error("HasMore = false, want true")
	}
	if res.NextCursor == "" {
		t.Error("NextCursor is empty, want non-empty")
	}
}

// T-020: Second page using cursor
func TestCursorPaginate_WithCursor(t *testing.T) {
	db := setupPaginationDB(t, 10)

	// First page
	var page1 []paginationItem
	res1, _ := CursorPaginate(db.Model(&paginationItem{}), "", "id", 3, "next", &page1)

	// Second page using NextCursor
	var page2 []paginationItem
	res2, err := CursorPaginate(db.Model(&paginationItem{}), res1.NextCursor, "id", 3, "next", &page2)
	if err != nil {
		t.Fatalf("CursorPaginate page 2 error: %v", err)
	}

	if len(page2) != 3 {
		t.Fatalf("page 2 got %d items, want 3", len(page2))
	}
	if page2[0].ID != 4 || page2[2].ID != 6 {
		t.Errorf("page 2 items = [%d..%d], want [4..6]", page2[0].ID, page2[2].ID)
	}
	if !res2.HasMore {
		t.Error("page 2 HasMore = false, want true")
	}
}

// T-021: Last page (has_more=false)
func TestCursorPaginate_LastPage(t *testing.T) {
	db := setupPaginationDB(t, 5)

	// Cursor after ID=3 → should get items 4,5
	cursor := base64.StdEncoding.EncodeToString([]byte("3"))

	var items []paginationItem
	res, err := CursorPaginate(db.Model(&paginationItem{}), cursor, "id", 5, "next", &items)
	if err != nil {
		t.Fatalf("CursorPaginate error: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}
	if res.HasMore {
		t.Error("HasMore = true, want false (last page)")
	}
}

// T-022: Prev direction
func TestCursorPaginate_PrevDirection(t *testing.T) {
	db := setupPaginationDB(t, 10)

	// Cursor at ID=7, direction prev → should get items before 7
	cursor := base64.StdEncoding.EncodeToString([]byte("7"))

	var items []paginationItem
	res, err := CursorPaginate(db.Model(&paginationItem{}), cursor, "id", 3, "prev", &items)
	if err != nil {
		t.Fatalf("CursorPaginate error: %v", err)
	}

	if len(items) != 3 {
		t.Fatalf("got %d items, want 3", len(items))
	}
	// prev direction orders DESC, so items should be 6,5,4
	if items[0].ID != 6 || items[2].ID != 4 {
		t.Errorf("items = [%d..%d], want [6..4] (DESC)", items[0].ID, items[2].ID)
	}
	if !res.HasMore {
		t.Error("HasMore = false, want true")
	}
}

// T-023: Empty table
func TestCursorPaginate_EmptyTable(t *testing.T) {
	db := setupPaginationDB(t, 0)

	var items []paginationItem
	res, err := CursorPaginate(db.Model(&paginationItem{}), "", "id", 10, "next", &items)
	if err != nil {
		t.Fatalf("CursorPaginate error: %v", err)
	}

	if len(items) != 0 {
		t.Errorf("expected empty, got %d items", len(items))
	}
	if res.HasMore {
		t.Error("HasMore = true, want false")
	}
	if res.NextCursor != "" {
		t.Errorf("NextCursor = %q, want empty", res.NextCursor)
	}
}

// T-023b: Invalid cursor returns error
func TestCursorPaginate_InvalidCursor(t *testing.T) {
	db := setupPaginationDB(t, 5)

	var items []paginationItem
	_, err := CursorPaginate(db.Model(&paginationItem{}), "not-valid-base64!!!", "id", 5, "next", &items)
	if err == nil {
		t.Fatal("expected error for invalid cursor, got nil")
	}
}
