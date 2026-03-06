package database

import (
	"errors"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// testUser is a minimal model for transaction tests.
// Avoids importing database/models to keep the database package independent.
type testUser struct {
	ID      uint `gorm:"primarykey"`
	Name    string
	Credits int
}

func (testUser) TableName() string { return "users" }

func setupTxTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&testUser{}); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
	return db
}

// TC-01: WithTransaction commits when callback returns nil
func TestWithTransaction_Commit(t *testing.T) {
	db := setupTxTestDB(t)
	db.Create(&testUser{ID: 1, Name: "Alice", Credits: 100})

	err := WithTransaction(db, func(tx *gorm.DB) error {
		return tx.Model(&testUser{}).Where("id = ?", 1).Update("credits", 200).Error
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	var u testUser
	db.First(&u, 1)
	if u.Credits != 200 {
		t.Fatalf("expected credits 200, got %d", u.Credits)
	}
}

// TC-02: WithTransaction rolls back when callback returns error
func TestWithTransaction_Rollback(t *testing.T) {
	db := setupTxTestDB(t)
	db.Create(&testUser{ID: 1, Name: "Alice", Credits: 100})

	err := WithTransaction(db, func(tx *gorm.DB) error {
		tx.Model(&testUser{}).Where("id = ?", 1).Update("credits", 999)
		return errors.New("forced rollback")
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var u testUser
	db.First(&u, 1)
	if u.Credits != 100 {
		t.Fatalf("expected credits 100 after rollback, got %d", u.Credits)
	}
}

// TC-03: WithTransaction rolls back on panic inside callback
func TestWithTransaction_PanicRollback(t *testing.T) {
	db := setupTxTestDB(t)
	db.Create(&testUser{ID: 1, Name: "Alice", Credits: 100})

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic to propagate")
			}
		}()
		_ = WithTransaction(db, func(tx *gorm.DB) error {
			tx.Model(&testUser{}).Where("id = ?", 1).Update("credits", 999)
			panic("test panic")
		})
	}()

	// Give SQLite a moment to release any locks
	time.Sleep(10 * time.Millisecond)

	var u testUser
	db.First(&u, 1)
	if u.Credits != 100 {
		t.Fatalf("expected credits 100 after panic rollback, got %d", u.Credits)
	}
}

// TC-04: WithTransaction propagates the callback error to the caller
func TestWithTransaction_ErrorPropagation(t *testing.T) {
	db := setupTxTestDB(t)

	sentinel := errors.New("specific error")
	err := WithTransaction(db, func(tx *gorm.DB) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got: %v", err)
	}
}

// TC-05: TransferCredits moves credits atomically
func TestTransferCredits_Success(t *testing.T) {
	db := setupTxTestDB(t)
	db.Create(&testUser{ID: 1, Name: "Alice", Credits: 200})
	db.Create(&testUser{ID: 2, Name: "Bob", Credits: 50})

	if err := TransferCredits(db, 1, 2, 100); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	var alice, bob testUser
	db.First(&alice, 1)
	db.First(&bob, 2)

	if alice.Credits != 100 {
		t.Fatalf("expected Alice credits 100, got %d", alice.Credits)
	}
	if bob.Credits != 150 {
		t.Fatalf("expected Bob credits 150, got %d", bob.Credits)
	}
}

// TC-06: TransferCredits fails if source user doesn't exist
func TestTransferCredits_SourceNotFound(t *testing.T) {
	db := setupTxTestDB(t)
	db.Create(&testUser{ID: 2, Name: "Bob", Credits: 50})

	err := TransferCredits(db, 999, 2, 100)
	if err == nil {
		t.Fatal("expected error for missing source user")
	}

	// Bob's credits should be unchanged
	var bob testUser
	db.First(&bob, 2)
	if bob.Credits != 50 {
		t.Fatalf("expected Bob credits 50 unchanged, got %d", bob.Credits)
	}
}

// TC-07: TransferCredits fails if destination user doesn't exist
func TestTransferCredits_DestNotFound(t *testing.T) {
	db := setupTxTestDB(t)
	db.Create(&testUser{ID: 1, Name: "Alice", Credits: 200})

	err := TransferCredits(db, 1, 999, 100)
	if err == nil {
		t.Fatal("expected error for missing destination user")
	}

	// Alice's credits should be unchanged
	var alice testUser
	db.First(&alice, 1)
	if alice.Credits != 200 {
		t.Fatalf("expected Alice credits 200 unchanged, got %d", alice.Credits)
	}
}
