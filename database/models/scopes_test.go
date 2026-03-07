package models

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupScopesTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&User{}, &Post{}); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
	return db
}

// T01: WithTrashed includes soft-deleted records
func TestWithTrashed_IncludesDeletedRecords(t *testing.T) {
	db := setupScopesTestDB(t)
	db.Create(&User{Name: "Alice", Email: "alice@example.com", Password: "pass"})
	db.Create(&User{Name: "Bob", Email: "bob@example.com", Password: "pass"})

	// Soft-delete Alice
	db.Where("email = ?", "alice@example.com").Delete(&User{})

	var users []User
	db.Scopes(WithTrashed).Find(&users)
	if len(users) != 2 {
		t.Fatalf("expected 2 users (including deleted), got %d", len(users))
	}
}

// T02: WithTrashed includes active records
func TestWithTrashed_IncludesActiveRecords(t *testing.T) {
	db := setupScopesTestDB(t)
	db.Create(&User{Name: "Alice", Email: "alice@example.com", Password: "pass"})

	var users []User
	db.Scopes(WithTrashed).Find(&users)
	if len(users) != 1 {
		t.Fatalf("expected 1 active user, got %d", len(users))
	}
	if users[0].Name != "Alice" {
		t.Fatalf("expected 'Alice', got '%s'", users[0].Name)
	}
}

// T03: OnlyTrashed returns only soft-deleted records
func TestOnlyTrashed_ReturnsOnlyDeletedRecords(t *testing.T) {
	db := setupScopesTestDB(t)
	db.Create(&User{Name: "Alice", Email: "alice@example.com", Password: "pass"})
	db.Create(&User{Name: "Bob", Email: "bob@example.com", Password: "pass"})

	// Soft-delete Bob
	db.Where("email = ?", "bob@example.com").Delete(&User{})

	var users []User
	db.Scopes(OnlyTrashed).Find(&users)
	if len(users) != 1 {
		t.Fatalf("expected 1 trashed user, got %d", len(users))
	}
	if users[0].Name != "Bob" {
		t.Fatalf("expected 'Bob', got '%s'", users[0].Name)
	}
}

// T04: OnlyTrashed excludes active records
func TestOnlyTrashed_ExcludesActiveRecords(t *testing.T) {
	db := setupScopesTestDB(t)
	db.Create(&User{Name: "Alice", Email: "alice@example.com", Password: "pass"})
	db.Create(&User{Name: "Bob", Email: "bob@example.com", Password: "pass"})

	// No deletions
	var users []User
	db.Scopes(OnlyTrashed).Find(&users)
	if len(users) != 0 {
		t.Fatalf("expected 0 trashed users, got %d", len(users))
	}
}

// T05: Default query excludes soft-deleted records
func TestDefaultQuery_ExcludesDeletedRecords(t *testing.T) {
	db := setupScopesTestDB(t)
	db.Create(&User{Name: "Alice", Email: "alice@example.com", Password: "pass"})
	db.Create(&User{Name: "Bob", Email: "bob@example.com", Password: "pass"})

	// Soft-delete Alice
	db.Where("email = ?", "alice@example.com").Delete(&User{})

	var users []User
	db.Find(&users)
	if len(users) != 1 {
		t.Fatalf("expected 1 active user, got %d", len(users))
	}
	if users[0].Name != "Bob" {
		t.Fatalf("expected 'Bob', got '%s'", users[0].Name)
	}
}
