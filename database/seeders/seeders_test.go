package seeders

import (
	"errors"
	"strings"
	"testing"

	"github.com/RAiWorks/RGo/database/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Post{}); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
	return db
}

// mockSeeder tracks whether Seed was called.
type mockSeeder struct {
	name   string
	called bool
	err    error
}

func (m *mockSeeder) Name() string         { return m.name }
func (m *mockSeeder) Seed(db *gorm.DB) error { m.called = true; return m.err }

// TC-01: Register adds a seeder to the registry.
func TestRegister_AddsSeeder(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	Register(&mockSeeder{name: "test"})

	names := Names()
	if len(names) != 1 {
		t.Fatalf("expected 1 seeder, got %d", len(names))
	}
	if names[0] != "test" {
		t.Fatalf("expected name 'test', got %q", names[0])
	}
}

// TC-02: RunAll executes all registered seeders.
func TestRunAll_ExecutesSeeders(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	a := &mockSeeder{name: "a"}
	b := &mockSeeder{name: "b"}
	Register(a)
	Register(b)

	db := setupTestDB(t)
	if err := RunAll(db); err != nil {
		t.Fatalf("RunAll failed: %v", err)
	}
	if !a.called {
		t.Fatal("expected seeder 'a' to be called")
	}
	if !b.called {
		t.Fatal("expected seeder 'b' to be called")
	}
}

// TC-03: RunAll stops on first error.
func TestRunAll_StopsOnError(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	a := &mockSeeder{name: "a"}
	b := &mockSeeder{name: "b", err: errors.New("seed failed")}
	Register(a)
	Register(b)

	db := setupTestDB(t)
	err := RunAll(db)
	if err == nil {
		t.Fatal("expected error from RunAll")
	}
	if !strings.Contains(err.Error(), "seed failed") {
		t.Fatalf("expected error to contain 'seed failed', got %q", err.Error())
	}
}

// TC-04: RunByName executes the named seeder only.
func TestRunByName_FindsSeeder(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	a := &mockSeeder{name: "first"}
	b := &mockSeeder{name: "second"}
	Register(a)
	Register(b)

	db := setupTestDB(t)
	if err := RunByName(db, "second"); err != nil {
		t.Fatalf("RunByName failed: %v", err)
	}
	if a.called {
		t.Fatal("expected seeder 'first' NOT to be called")
	}
	if !b.called {
		t.Fatal("expected seeder 'second' to be called")
	}
}

// TC-05: RunByName returns error for unknown name.
func TestRunByName_NotFound(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	err := RunByName(setupTestDB(t), "nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown seeder")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected 'not found' in error, got %q", err.Error())
	}
}

// TC-06: UserSeeder creates admin and regular user.
func TestUserSeeder_CreatesUsers(t *testing.T) {
	db := setupTestDB(t)

	seeder := &UserSeeder{}
	if err := seeder.Seed(db); err != nil {
		t.Fatalf("UserSeeder.Seed failed: %v", err)
	}

	var users []models.User
	db.Find(&users)

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}

	var admin models.User
	db.Where("email = ?", "admin@example.com").First(&admin)
	if admin.Role != "admin" {
		t.Fatalf("expected admin role, got %q", admin.Role)
	}

	var user models.User
	db.Where("email = ?", "user@example.com").First(&user)
	if user.Role != "user" {
		t.Fatalf("expected user role, got %q", user.Role)
	}
}

// TC-07: UserSeeder is idempotent — no duplicates on second run.
func TestUserSeeder_Idempotent(t *testing.T) {
	db := setupTestDB(t)

	seeder := &UserSeeder{}
	seeder.Seed(db)
	seeder.Seed(db)

	var count int64
	db.Model(&models.User{}).Count(&count)
	if count != 2 {
		t.Fatalf("expected 2 users after double seed, got %d", count)
	}
}
