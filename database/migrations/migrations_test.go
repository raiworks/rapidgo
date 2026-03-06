package migrations

import (
	"testing"

	"github.com/RAiWorks/RGo/database/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// setupTestDB opens SQLite :memory: and returns a clean *gorm.DB.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	return db
}

// TC-01: models.All() returns all registered models.
func TestModelsAll(t *testing.T) {
	all := models.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 models, got %d", len(all))
	}
	// Verify types
	if _, ok := all[0].(*models.User); !ok {
		t.Fatal("expected first model to be *User")
	}
	if _, ok := all[1].(*models.Post); !ok {
		t.Fatal("expected second model to be *Post")
	}
}

// TC-02: NewMigrator auto-creates schema_migrations table.
func TestNewMigrator_CreatesTable(t *testing.T) {
	db := setupTestDB(t)
	m, err := NewMigrator(db)
	if err != nil {
		t.Fatalf("NewMigrator failed: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil Migrator")
	}
	if !db.Migrator().HasTable(&SchemaMigration{}) {
		t.Fatal("expected schema_migrations table to exist")
	}
}

// TC-03: Run() applies pending migrations.
func TestMigrator_RunPending(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	Register(Migration{
		Version: "20260306_create_alpha",
		Up: func(db *gorm.DB) error {
			return db.Exec("CREATE TABLE alpha (id INTEGER PRIMARY KEY)").Error
		},
		Down: func(db *gorm.DB) error {
			return db.Exec("DROP TABLE alpha").Error
		},
	})
	Register(Migration{
		Version: "20260306_create_beta",
		Up: func(db *gorm.DB) error {
			return db.Exec("CREATE TABLE beta (id INTEGER PRIMARY KEY)").Error
		},
		Down: func(db *gorm.DB) error {
			return db.Exec("DROP TABLE beta").Error
		},
	})

	db := setupTestDB(t)
	m, _ := NewMigrator(db)

	n, err := m.Run()
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if n != 2 {
		t.Fatalf("expected 2 applied, got %d", n)
	}

	// Tables exist
	if !db.Migrator().HasTable("alpha") {
		t.Fatal("expected alpha table to exist")
	}
	if !db.Migrator().HasTable("beta") {
		t.Fatal("expected beta table to exist")
	}

	// Tracking rows exist
	var count int64
	db.Model(&SchemaMigration{}).Count(&count)
	if count != 2 {
		t.Fatalf("expected 2 tracking rows, got %d", count)
	}

	// Batch should be 1
	var rec SchemaMigration
	db.First(&rec)
	if rec.Batch != 1 {
		t.Fatalf("expected batch 1, got %d", rec.Batch)
	}
}

// TC-04: Run() is idempotent — second call applies nothing.
func TestMigrator_RunIdempotent(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	Register(Migration{
		Version: "20260306_create_gamma",
		Up: func(db *gorm.DB) error {
			return db.Exec("CREATE TABLE gamma (id INTEGER PRIMARY KEY)").Error
		},
		Down: func(db *gorm.DB) error {
			return db.Exec("DROP TABLE gamma").Error
		},
	})

	db := setupTestDB(t)
	m, _ := NewMigrator(db)

	n1, _ := m.Run()
	if n1 != 1 {
		t.Fatalf("expected 1 applied first run, got %d", n1)
	}

	n2, _ := m.Run()
	if n2 != 0 {
		t.Fatalf("expected 0 applied second run, got %d", n2)
	}
}

// TC-05: Migrations run in version-sorted order.
func TestMigrator_RunOrder(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	var order []string

	// Register out of order
	Register(Migration{
		Version: "20260306_b_second",
		Up: func(db *gorm.DB) error {
			order = append(order, "b")
			return nil
		},
		Down: func(db *gorm.DB) error { return nil },
	})
	Register(Migration{
		Version: "20260306_a_first",
		Up: func(db *gorm.DB) error {
			order = append(order, "a")
			return nil
		},
		Down: func(db *gorm.DB) error { return nil },
	})

	db := setupTestDB(t)
	m, _ := NewMigrator(db)
	m.Run()

	if len(order) != 2 || order[0] != "a" || order[1] != "b" {
		t.Fatalf("expected [a, b], got %v", order)
	}
}

// TC-06: Rollback() undoes the last batch.
func TestMigrator_Rollback(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	Register(Migration{
		Version: "20260306_create_delta",
		Up: func(db *gorm.DB) error {
			return db.Exec("CREATE TABLE delta (id INTEGER PRIMARY KEY)").Error
		},
		Down: func(db *gorm.DB) error {
			return db.Exec("DROP TABLE delta").Error
		},
	})
	Register(Migration{
		Version: "20260306_create_epsilon",
		Up: func(db *gorm.DB) error {
			return db.Exec("CREATE TABLE epsilon (id INTEGER PRIMARY KEY)").Error
		},
		Down: func(db *gorm.DB) error {
			return db.Exec("DROP TABLE epsilon").Error
		},
	})

	db := setupTestDB(t)
	m, _ := NewMigrator(db)
	m.Run()

	n, err := m.Rollback()
	if err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}
	if n != 2 {
		t.Fatalf("expected 2 rolled back, got %d", n)
	}

	if db.Migrator().HasTable("delta") {
		t.Fatal("expected delta table to be dropped")
	}
	if db.Migrator().HasTable("epsilon") {
		t.Fatal("expected epsilon table to be dropped")
	}

	var count int64
	db.Model(&SchemaMigration{}).Count(&count)
	if count != 0 {
		t.Fatalf("expected 0 tracking rows, got %d", count)
	}
}

// TC-07: Rollback only undoes the last batch.
func TestMigrator_RollbackBatches(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	Register(Migration{
		Version: "20260306_create_zeta",
		Up: func(db *gorm.DB) error {
			return db.Exec("CREATE TABLE zeta (id INTEGER PRIMARY KEY)").Error
		},
		Down: func(db *gorm.DB) error {
			return db.Exec("DROP TABLE zeta").Error
		},
	})

	db := setupTestDB(t)
	m, _ := NewMigrator(db)

	// Batch 1: zeta
	m.Run()

	// Register another, batch 2
	Register(Migration{
		Version: "20260306_create_eta",
		Up: func(db *gorm.DB) error {
			return db.Exec("CREATE TABLE eta (id INTEGER PRIMARY KEY)").Error
		},
		Down: func(db *gorm.DB) error {
			return db.Exec("DROP TABLE eta").Error
		},
	})
	m.Run()

	// Rollback should only undo batch 2 (eta)
	n, err := m.Rollback()
	if err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 rolled back, got %d", n)
	}

	if !db.Migrator().HasTable("zeta") {
		t.Fatal("expected zeta table to still exist")
	}
	if db.Migrator().HasTable("eta") {
		t.Fatal("expected eta table to be dropped")
	}
}

// TC-08: Status() returns correct applied/pending status.
func TestMigrator_Status(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	Register(Migration{
		Version: "20260306_a_create_theta",
		Up:      func(db *gorm.DB) error { return nil },
		Down:    func(db *gorm.DB) error { return nil },
	})

	db := setupTestDB(t)
	m, _ := NewMigrator(db)

	// Run first migration
	m.Run()

	// Register a second one (not yet run)
	Register(Migration{
		Version: "20260306_b_create_iota",
		Up:      func(db *gorm.DB) error { return nil },
		Down:    func(db *gorm.DB) error { return nil },
	})

	statuses, err := m.Status()
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}
	if len(statuses) != 2 {
		t.Fatalf("expected 2 statuses, got %d", len(statuses))
	}

	// First (a_create_theta) should be applied
	if !statuses[0].Applied || statuses[0].Batch != 1 {
		t.Fatalf("expected theta applied (batch 1), got applied=%v batch=%d", statuses[0].Applied, statuses[0].Batch)
	}
	// Second (b_create_iota) should be pending
	if statuses[1].Applied {
		t.Fatal("expected iota to be pending")
	}
}

// TC-09: Rollback with no applied migrations returns 0.
func TestMigrator_RollbackEmpty(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	db := setupTestDB(t)
	m, _ := NewMigrator(db)

	n, err := m.Rollback()
	if err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}
	if n != 0 {
		t.Fatalf("expected 0 rolled back, got %d", n)
	}
}
