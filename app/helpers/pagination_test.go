package helpers

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type paginationItem struct {
	ID   uint   `gorm:"primarykey"`
	Name string
}

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
		db.Create(&paginationItem{ID: uint(i), Name: "item"})
	}
	return db
}

// TC-01: Page < 1 clamped to 1.
func TestPaginate_PageClamped(t *testing.T) {
	db := setupPaginationDB(t, 5)
	var dest []paginationItem
	result, err := Paginate(db.Model(&paginationItem{}), 0, 10, &dest)
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if result.Page != 1 {
		t.Fatalf("Page = %d, want 1", result.Page)
	}
}

// TC-02: PerPage < 1 clamped to 15.
func TestPaginate_PerPageTooLow(t *testing.T) {
	db := setupPaginationDB(t, 5)
	var dest []paginationItem
	result, err := Paginate(db.Model(&paginationItem{}), 1, 0, &dest)
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if result.PerPage != 15 {
		t.Fatalf("PerPage = %d, want 15", result.PerPage)
	}
}

// TC-03: PerPage > 100 clamped to 15.
func TestPaginate_PerPageTooHigh(t *testing.T) {
	db := setupPaginationDB(t, 5)
	var dest []paginationItem
	result, err := Paginate(db.Model(&paginationItem{}), 1, 200, &dest)
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if result.PerPage != 15 {
		t.Fatalf("PerPage = %d, want 15", result.PerPage)
	}
}

// TC-04: Valid inputs pass through.
func TestPaginate_ValidInputs(t *testing.T) {
	db := setupPaginationDB(t, 25)
	var dest []paginationItem
	result, err := Paginate(db.Model(&paginationItem{}), 2, 10, &dest)
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if result.Page != 2 {
		t.Fatalf("Page = %d, want 2", result.Page)
	}
	if result.PerPage != 10 {
		t.Fatalf("PerPage = %d, want 10", result.PerPage)
	}
}

// TC-05: TotalPages ceiling division (25 / 10 = 3).
func TestPaginate_TotalPagesCeiling(t *testing.T) {
	db := setupPaginationDB(t, 25)
	var dest []paginationItem
	result, err := Paginate(db.Model(&paginationItem{}), 1, 10, &dest)
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if result.Total != 25 {
		t.Fatalf("Total = %d, want 25", result.Total)
	}
	if result.TotalPages != 3 {
		t.Fatalf("TotalPages = %d, want 3", result.TotalPages)
	}
}

// TC-06: TotalPages exact division (30 / 10 = 3).
func TestPaginate_TotalPagesExact(t *testing.T) {
	db := setupPaginationDB(t, 30)
	var dest []paginationItem
	result, err := Paginate(db.Model(&paginationItem{}), 1, 10, &dest)
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if result.TotalPages != 3 {
		t.Fatalf("TotalPages = %d, want 3", result.TotalPages)
	}
}

// TC-07: TotalPages zero when empty.
func TestPaginate_EmptyTable(t *testing.T) {
	db := setupPaginationDB(t, 0)
	var dest []paginationItem
	result, err := Paginate(db.Model(&paginationItem{}), 1, 10, &dest)
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if result.Total != 0 {
		t.Fatalf("Total = %d, want 0", result.Total)
	}
	if result.TotalPages != 0 {
		t.Fatalf("TotalPages = %d, want 0", result.TotalPages)
	}
}

// TC-08: Offset calculation — page 2 returns correct slice.
func TestPaginate_OffsetSlice(t *testing.T) {
	db := setupPaginationDB(t, 25)
	var dest []paginationItem
	result, err := Paginate(db.Model(&paginationItem{}), 2, 10, &dest)
	if err != nil {
		t.Fatalf("Paginate: %v", err)
	}
	if len(dest) != 10 {
		t.Fatalf("len(dest) = %d, want 10", len(dest))
	}
	if result.Total != 25 {
		t.Fatalf("Total = %d, want 25", result.Total)
	}
	// Page 3 should have 5 remaining items.
	var page3 []paginationItem
	r3, _ := Paginate(db.Model(&paginationItem{}), 3, 10, &page3)
	if len(page3) != 5 {
		t.Fatalf("page 3 len = %d, want 5", len(page3))
	}
	if r3.TotalPages != 3 {
		t.Fatalf("TotalPages = %d, want 3", r3.TotalPages)
	}
}

// TC-09: Returns GORM error for invalid table.
func TestPaginate_GormError(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	// Query a table that does not exist.
	type ghost struct {
		ID int `gorm:"primarykey"`
	}
	var dest []ghost
	_, pErr := Paginate(db.Model(&ghost{}), 1, 10, &dest)
	if pErr == nil {
		t.Fatal("expected error for non-existent table")
	}
}
