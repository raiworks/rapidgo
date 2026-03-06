// Package seeders provides a registry-based database seeding system.
package seeders

import (
	"fmt"

	"gorm.io/gorm"
)

// Seeder defines the interface for database seeders.
type Seeder interface {
	// Name returns the seeder's unique name (used with --seeder flag).
	Name() string
	// Seed populates the database with data.
	Seed(db *gorm.DB) error
}

// registry holds all registered seeders.
var registry []Seeder

// Register adds a seeder to the global registry.
func Register(s Seeder) {
	registry = append(registry, s)
}

// ResetRegistry clears all registered seeders. Used in tests only.
func ResetRegistry() {
	registry = nil
}

// RunAll executes all registered seeders in registration order.
func RunAll(db *gorm.DB) error {
	for _, s := range registry {
		if err := s.Seed(db); err != nil {
			return fmt.Errorf("seeder %s failed: %w", s.Name(), err)
		}
	}
	return nil
}

// RunByName executes a single seeder by name.
func RunByName(db *gorm.DB, name string) error {
	for _, s := range registry {
		if s.Name() == name {
			return s.Seed(db)
		}
	}
	return fmt.Errorf("seeder %q not found", name)
}

// Names returns the names of all registered seeders.
func Names() []string {
	names := make([]string, len(registry))
	for i, s := range registry {
		names[i] = s.Name()
	}
	return names
}
