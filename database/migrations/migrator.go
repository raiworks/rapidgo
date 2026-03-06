// Package migrations provides a file-based migration engine for GORM.
package migrations

import (
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"
)

// SchemaMigration tracks applied migrations in the database.
type SchemaMigration struct {
	ID        uint   `gorm:"primaryKey"`
	Version   string `gorm:"size:255;uniqueIndex;not null"`
	Batch     int    `gorm:"not null"`
	AppliedAt time.Time
}

// MigrationFunc is a function that performs a migration step.
type MigrationFunc func(db *gorm.DB) error

// Migration represents a single migration with up and down functions.
type Migration struct {
	Version string
	Up      MigrationFunc
	Down    MigrationFunc
}

// MigrationStatus represents the status of a single migration.
type MigrationStatus struct {
	Version string
	Applied bool
	Batch   int
}

// registry holds all registered migrations (populated via Register).
var registry []Migration

// Register adds a migration to the global registry.
// Called from migration files (typically in init()).
func Register(m Migration) {
	registry = append(registry, m)
}

// ResetRegistry clears all registered migrations. Used in tests only.
func ResetRegistry() {
	registry = nil
}

// Migrator manages database migrations.
type Migrator struct {
	DB *gorm.DB
}

// NewMigrator creates a Migrator and ensures the schema_migrations table exists.
func NewMigrator(db *gorm.DB) (*Migrator, error) {
	if err := db.AutoMigrate(&SchemaMigration{}); err != nil {
		return nil, fmt.Errorf("failed to create schema_migrations table: %w", err)
	}
	return &Migrator{DB: db}, nil
}

// Run applies all pending migrations in version order.
// Returns the number of migrations applied.
func (m *Migrator) Run() (int, error) {
	applied, err := m.appliedVersions()
	if err != nil {
		return 0, err
	}

	pending := m.pendingMigrations(applied)
	if len(pending) == 0 {
		return 0, nil
	}

	batch, err := m.nextBatch()
	if err != nil {
		return 0, err
	}

	for _, mig := range pending {
		if err := mig.Up(m.DB); err != nil {
			return 0, fmt.Errorf("migration %s failed: %w", mig.Version, err)
		}
		record := SchemaMigration{
			Version:   mig.Version,
			Batch:     batch,
			AppliedAt: time.Now(),
		}
		if err := m.DB.Create(&record).Error; err != nil {
			return 0, fmt.Errorf("failed to record migration %s: %w", mig.Version, err)
		}
	}

	return len(pending), nil
}

// Rollback undoes the last batch of applied migrations.
// Returns the number of migrations rolled back.
func (m *Migrator) Rollback() (int, error) {
	var maxBatch int
	err := m.DB.Model(&SchemaMigration{}).Select("COALESCE(MAX(batch), 0)").Scan(&maxBatch).Error
	if err != nil {
		return 0, fmt.Errorf("failed to find last batch: %w", err)
	}
	if maxBatch == 0 {
		return 0, nil
	}

	var records []SchemaMigration
	if err := m.DB.Where("batch = ?", maxBatch).Order("version DESC").Find(&records).Error; err != nil {
		return 0, fmt.Errorf("failed to load batch %d: %w", maxBatch, err)
	}

	migrationMap := m.registryMap()
	for _, rec := range records {
		mig, ok := migrationMap[rec.Version]
		if !ok {
			return 0, fmt.Errorf("migration %s not found in registry", rec.Version)
		}
		if mig.Down == nil {
			return 0, fmt.Errorf("migration %s has no Down function", rec.Version)
		}
		if err := mig.Down(m.DB); err != nil {
			return 0, fmt.Errorf("rollback %s failed: %w", rec.Version, err)
		}
		if err := m.DB.Where("version = ?", rec.Version).Delete(&SchemaMigration{}).Error; err != nil {
			return 0, fmt.Errorf("failed to delete migration record %s: %w", rec.Version, err)
		}
	}

	return len(records), nil
}

// Status returns the status of all registered migrations.
func (m *Migrator) Status() ([]MigrationStatus, error) {
	applied, err := m.appliedMap()
	if err != nil {
		return nil, err
	}

	sorted := make([]Migration, len(registry))
	copy(sorted, registry)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Version < sorted[j].Version
	})

	statuses := make([]MigrationStatus, 0, len(sorted))
	for _, mig := range sorted {
		s := MigrationStatus{Version: mig.Version}
		if rec, ok := applied[mig.Version]; ok {
			s.Applied = true
			s.Batch = rec.Batch
		}
		statuses = append(statuses, s)
	}

	return statuses, nil
}

// appliedVersions returns a set of already-applied migration versions.
func (m *Migrator) appliedVersions() (map[string]bool, error) {
	var records []SchemaMigration
	if err := m.DB.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to query schema_migrations: %w", err)
	}
	applied := make(map[string]bool, len(records))
	for _, r := range records {
		applied[r.Version] = true
	}
	return applied, nil
}

// appliedMap returns applied migrations keyed by version.
func (m *Migrator) appliedMap() (map[string]SchemaMigration, error) {
	var records []SchemaMigration
	if err := m.DB.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to query schema_migrations: %w", err)
	}
	applied := make(map[string]SchemaMigration, len(records))
	for _, r := range records {
		applied[r.Version] = r
	}
	return applied, nil
}

// pendingMigrations returns migrations not yet applied, sorted by version.
func (m *Migrator) pendingMigrations(applied map[string]bool) []Migration {
	var pending []Migration
	for _, mig := range registry {
		if !applied[mig.Version] {
			pending = append(pending, mig)
		}
	}
	sort.Slice(pending, func(i, j int) bool {
		return pending[i].Version < pending[j].Version
	})
	return pending
}

// nextBatch returns the next batch number.
func (m *Migrator) nextBatch() (int, error) {
	var maxBatch int
	err := m.DB.Model(&SchemaMigration{}).Select("COALESCE(MAX(batch), 0)").Scan(&maxBatch).Error
	if err != nil {
		return 0, fmt.Errorf("failed to determine next batch: %w", err)
	}
	return maxBatch + 1, nil
}

// registryMap returns all registered migrations keyed by version.
func (m *Migrator) registryMap() map[string]Migration {
	rm := make(map[string]Migration, len(registry))
	for _, mig := range registry {
		rm[mig.Version] = mig
	}
	return rm
}
