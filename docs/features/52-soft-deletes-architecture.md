# 📐 Architecture: Soft Deletes

> **Feature**: `52` — Soft Deletes
> **Discussion**: [`52-soft-deletes-discussion.md`](52-soft-deletes-discussion.md)
> **Status**: � COMPLETE
> **Date**: 2026-03-07

---

## Overview

Soft deletes add a `DeletedAt` field to `BaseModel` so that all models embedded with `BaseModel` gain automatic GORM soft delete behavior. Queries exclude soft-deleted records by default; `Unscoped()` retrieves or permanently removes them. The feature adds scope helpers (`WithTrashed`, `OnlyTrashed`), service methods (`HardDelete`, `Restore`), a database migration, and comprehensive tests.

---

## File Structure

```
database/
  models/
    base.go          ← MODIFIED — add DeletedAt field
    scopes.go        ← NEW — WithTrashed, OnlyTrashed scope helpers
    scopes_test.go   ← NEW — scope helper tests
  migrations/
    20260308000001_add_soft_deletes.go  ← NEW — migration for existing tables
app/
  services/
    user_service.go       ← MODIFIED — add HardDelete, Restore methods
    user_service_test.go  ← MODIFIED — add soft delete test cases
```

---

## Component Design

### 1. BaseModel Update

**Package**: `database/models`
**File**: `base.go`

```go
package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel provides common fields for all models.
// Embed this in your model structs to get ID, CreatedAt,
// UpdatedAt, and soft delete support via DeletedAt.
type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

**Behavior**:
- GORM automatically sets `deleted_at` to current time on `db.Delete()`
- GORM automatically adds `WHERE deleted_at IS NULL` to all queries
- `db.Unscoped()` bypasses the automatic filter
- `gorm:"index"` ensures the filter is index-backed

---

### 2. Query Scope Helpers

**Package**: `database/models`
**File**: `scopes.go`

```go
package models

import "gorm.io/gorm"

// WithTrashed returns a GORM scope that includes soft-deleted records.
//
//	db.Scopes(models.WithTrashed).Find(&users)
func WithTrashed(db *gorm.DB) *gorm.DB {
	return db.Unscoped()
}

// OnlyTrashed returns a GORM scope that returns only soft-deleted records.
//
//	db.Scopes(models.OnlyTrashed).Find(&users)
func OnlyTrashed(db *gorm.DB) *gorm.DB {
	return db.Unscoped().Where("deleted_at IS NOT NULL")
}
```

**Usage**: Standard GORM scope pattern — `db.Scopes(models.WithTrashed).Find(&users)`

---

### 3. UserService Extensions

**Package**: `app/services`
**File**: `user_service.go`

```go
// HardDelete permanently removes a user from the database.
// This bypasses soft delete and cannot be undone.
func (s *UserService) HardDelete(id uint) error {
	return s.DB.Unscoped().Delete(&models.User{}, id).Error
}

// Restore recovers a soft-deleted user by clearing their deleted_at timestamp.
func (s *UserService) Restore(id uint) error {
	return s.DB.Unscoped().Model(&models.User{}).Where("id = ?", id).Update("deleted_at", nil).Error
}
```

**Existing `Delete()` method** — unchanged code, but behavior changes:
- Before: hard delete (row removed)
- After: soft delete (`deleted_at` set to current time)

---

### 4. Migration

**Package**: `database/migrations`
**File**: `20260308000001_add_soft_deletes.go`

```go
func init() {
	Register(Migration{
		Version: "20260308000001_add_soft_deletes",
		Up: func(db *gorm.DB) error {
			type User struct{ DeletedAt gorm.DeletedAt `gorm:"index"` }
			type Post struct{ DeletedAt gorm.DeletedAt `gorm:"index"` }
			if err := db.Migrator().AddColumn(&User{}, "DeletedAt"); err != nil {
				return err
			}
			if err := db.Migrator().CreateIndex(&User{}, "DeletedAt"); err != nil {
				return err
			}
			if err := db.Migrator().AddColumn(&Post{}, "DeletedAt"); err != nil {
				return err
			}
			return db.Migrator().CreateIndex(&Post{}, "DeletedAt")
		},
		Down: func(db *gorm.DB) error {
			type User struct{ DeletedAt gorm.DeletedAt }
			type Post struct{ DeletedAt gorm.DeletedAt }
			if err := db.Migrator().DropColumn(&Post{}, "DeletedAt"); err != nil {
				return err
			}
			return db.Migrator().DropColumn(&User{}, "DeletedAt")
		},
	})
}
```

**Notes**:
- GORM Migrator API handles database-specific SQL generation for portability
- Column + index for both `users` and `posts` tables
- Down migration drops columns for both tables

---

## Breaking Changes

| Change | Impact | Mitigation |
|---|---|---|
| `Delete()` now soft-deletes | Existing code that calls `Delete()` will no longer permanently remove rows | Use `HardDelete()` for permanent removal |
| Queries exclude soft-deleted rows | `Find()`, `First()`, etc. no longer return deleted records | Use `Scopes(WithTrashed)` to include them |
| `BaseModel` gains `DeletedAt` field | All models get the field; `AutoMigrate` will add the column | Migration provided for existing tables |

---

## Dependencies

- **GORM v1.31.1** — `gorm.DeletedAt` type, `Unscoped()` method (already in go.mod)
- **Feature #11** — `BaseModel`, `User`, `Post` models (already implemented)
- No new external dependencies required

---

## Next

Tasks doc → `52-soft-deletes-tasks.md`
