package services

import (
	"testing"

	"github.com/RAiWorks/RapidGo/database/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database with the User table.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
	return db
}

// TC-01: NewUserService returns a valid service
func TestNewUserService_ReturnsService(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)
	if svc == nil {
		t.Fatal("expected non-nil UserService")
	}
	if svc.DB == nil {
		t.Fatal("expected DB to be set")
	}
}

// TC-02: Create inserts a new user
func TestCreate_ReturnsUser(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	user, err := svc.Create("Alice", "alice@example.com", "pass123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID == 0 {
		t.Fatal("expected user ID > 0")
	}
	if user.Name != "Alice" {
		t.Fatalf("expected name 'Alice', got '%s'", user.Name)
	}
	if user.Email != "alice@example.com" {
		t.Fatalf("expected email 'alice@example.com', got '%s'", user.Email)
	}
}

// TC-03: Create rejects duplicate email
func TestCreate_DuplicateEmail_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	_, err := svc.Create("Alice", "alice@example.com", "pass")
	if err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	_, err = svc.Create("Bob", "alice@example.com", "pass")
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
	if err.Error() != "email already exists" {
		t.Fatalf("expected 'email already exists', got '%s'", err.Error())
	}
}

// TC-04: GetByID returns existing user
func TestGetByID_ReturnsUser(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	created, _ := svc.Create("Alice", "alice@example.com", "pass")

	user, err := svc.GetByID(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Name != "Alice" {
		t.Fatalf("expected name 'Alice', got '%s'", user.Name)
	}
	if user.Email != "alice@example.com" {
		t.Fatalf("expected email 'alice@example.com', got '%s'", user.Email)
	}
}

// TC-05: GetByID returns error for non-existent ID
func TestGetByID_NotFound_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	_, err := svc.GetByID(9999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

// TC-06: Update modifies specified fields
func TestUpdate_UpdatesFields(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	created, _ := svc.Create("Alice", "alice@example.com", "pass")

	updated, err := svc.Update(created.ID, map[string]interface{}{"name": "Bob"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "Bob" {
		t.Fatalf("expected name 'Bob', got '%s'", updated.Name)
	}
}

// TC-07: Update returns error for non-existent ID
func TestUpdate_NotFound_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	_, err := svc.Update(9999, map[string]interface{}{"name": "X"})
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

// TC-08 / T12: Delete soft-deletes user (hidden from normal queries, still exists via Unscoped)
func TestDelete_RemovesUser(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	created, _ := svc.Create("Alice", "alice@example.com", "pass")

	err := svc.Delete(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Normal query should not find the user
	_, err = svc.GetByID(created.ID)
	if err == nil {
		t.Fatal("expected error after soft delete, user should not be visible")
	}

	// Unscoped query should still find the user
	var user models.User
	if err := db.Unscoped().First(&user, created.ID).Error; err != nil {
		t.Fatalf("expected soft-deleted user to exist via Unscoped, got: %v", err)
	}
}

// T06: Delete sets deleted_at timestamp
func TestDelete_SoftDeletesUser(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	created, _ := svc.Create("Alice", "alice@example.com", "pass")

	err := svc.Delete(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var user models.User
	db.Unscoped().First(&user, created.ID)
	if !user.DeletedAt.Valid {
		t.Fatal("expected DeletedAt to be set after soft delete")
	}
}

// T07: Delete sets a non-zero deleted_at timestamp
func TestDelete_SetsDeletedAtTimestamp(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	created, _ := svc.Create("Alice", "alice@example.com", "pass")

	if err := svc.Delete(created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var user models.User
	db.Unscoped().First(&user, created.ID)
	if !user.DeletedAt.Valid {
		t.Fatal("expected DeletedAt to be valid")
	}
	if user.DeletedAt.Time.IsZero() {
		t.Fatal("expected DeletedAt timestamp to be non-zero")
	}
}

// T08: HardDelete permanently removes user
func TestHardDelete_PermanentlyRemovesUser(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	created, _ := svc.Create("Alice", "alice@example.com", "pass")

	err := svc.HardDelete(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var user models.User
	err = db.Unscoped().First(&user, created.ID).Error
	if err == nil {
		t.Fatal("expected hard-deleted user to not exist even with Unscoped")
	}
}

// T09: HardDelete on non-existent ID returns no error
func TestHardDelete_NonExistentID_NoError(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	err := svc.HardDelete(9999)
	if err != nil {
		t.Fatalf("expected no error for non-existent ID, got: %v", err)
	}
}

// T10: Restore recovers a soft-deleted user
func TestRestore_RecoversSoftDeletedUser(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	created, _ := svc.Create("Alice", "alice@example.com", "pass")

	// Soft-delete then restore
	svc.Delete(created.ID)
	err := svc.Restore(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should be visible again via normal query
	user, err := svc.GetByID(created.ID)
	if err != nil {
		t.Fatalf("expected restored user to be queryable, got: %v", err)
	}
	if user.DeletedAt.Valid {
		t.Fatal("expected restored user to have nil DeletedAt")
	}
}

// T11: Restore on active user is a no-op
func TestRestore_NonDeletedUser_NoError(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	created, _ := svc.Create("Alice", "alice@example.com", "pass")

	err := svc.Restore(created.ID)
	if err != nil {
		t.Fatalf("expected no error for restoring active user, got: %v", err)
	}

	user, err := svc.GetByID(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Name != "Alice" {
		t.Fatalf("expected 'Alice', got '%s'", user.Name)
	}
}
