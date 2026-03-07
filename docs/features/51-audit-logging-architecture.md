# 📐 Architecture: Audit Logging

> **Feature**: `51` — Audit Logging
> **Discussion**: [`51-audit-logging-discussion.md`](51-audit-logging-discussion.md)
> **Status**: 🟡 IN PROGRESS
> **Date**: 2026-03-07

---

## Overview

Feature #51 adds a `core/audit` package that provides a `Logger` for recording model-level changes to an `audit_logs` database table. The package stores who performed an action, what model was affected, the action taken, and optional JSON diffs of old/new values plus freeform metadata. Application code calls `audit.Log()` explicitly — there are no automatic GORM hooks.

---

## File Structure

```
core/
  audit/
    audit.go       ← NEW — Logger, Entry, Log(), Find(), ForModel()
    audit_test.go  ← NEW — unit tests
database/
  models/
    audit_log.go   ← NEW — AuditLog model
    registry.go    ← MODIFIED — add AuditLog to All()
  migrations/
    20260308000003_create_audit_logs_table.go  ← NEW — migration
```

---

## Component Design

### 1. AuditLog Model

**Package**: `database/models`
**File**: `audit_log.go`

```go
package models

import "time"

// AuditLog records a single auditable action on a model.
type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null;default:0" json:"user_id"`
	Action    string    `gorm:"size:50;not null;index" json:"action"`
	ModelType string    `gorm:"size:100;not null;index" json:"model_type"`
	ModelID   uint      `gorm:"not null;index" json:"model_id"`
	OldValues string    `gorm:"type:text" json:"old_values,omitempty"`
	NewValues string    `gorm:"type:text" json:"new_values,omitempty"`
	Metadata  string    `gorm:"type:text" json:"metadata,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
```

**Field details**:

| Field | Type | Tags | Purpose |
|---|---|---|---|
| `ID` | `uint` | `gorm:"primaryKey"` | Auto-increment primary key |
| `UserID` | `uint` | `gorm:"index;not null;default:0"` | ID of the actor; 0 = system/anonymous |
| `Action` | `string` | `gorm:"size:50;not null;index"` | Verb: "create", "update", "delete", or custom |
| `ModelType` | `string` | `gorm:"size:100;not null;index"` | Model name, e.g. "User", "Post" |
| `ModelID` | `uint` | `gorm:"not null;index"` | Primary key of the affected record |
| `OldValues` | `string` | `gorm:"type:text"` | JSON of previous values (for updates/deletes) |
| `NewValues` | `string` | `gorm:"type:text"` | JSON of new values (for creates/updates) |
| `Metadata` | `string` | `gorm:"type:text"` | JSON of extra context (IP, user agent, reason, etc.) |
| `CreatedAt` | `time.Time` | — | Timestamp of the audit entry |

**Design notes**:
- Does NOT embed `BaseModel` — audit logs are immutable, no `UpdatedAt` or `DeletedAt`
- `UserID` defaults to 0 for system-initiated actions (migrations, seeds, cron jobs)
- `OldValues`, `NewValues`, `Metadata` are plain `text` columns with JSON content — works across all GORM-supported drivers

---

### 2. Audit Package

**Package**: `core/audit`
**File**: `audit.go`

```go
package audit

import (
	"encoding/json"

	"github.com/RAiWorks/RapidGo/database/models"
	"gorm.io/gorm"
)

// Logger writes audit log entries to the database.
type Logger struct {
	db *gorm.DB
}

// NewLogger creates a new audit Logger backed by the given database connection.
func NewLogger(db *gorm.DB) *Logger {
	return &Logger{db: db}
}

// Entry holds the data for a single audit log record.
type Entry struct {
	UserID    uint
	Action    string
	ModelType string
	ModelID   uint
	OldValues map[string]interface{}
	NewValues map[string]interface{}
	Metadata  map[string]interface{}
}

// Log persists an audit entry to the database.
func (l *Logger) Log(e Entry) error {
	record := models.AuditLog{
		UserID:    e.UserID,
		Action:    e.Action,
		ModelType: e.ModelType,
		ModelID:   e.ModelID,
	}
	if e.OldValues != nil {
		b, err := json.Marshal(e.OldValues)
		if err != nil {
			return err
		}
		record.OldValues = string(b)
	}
	if e.NewValues != nil {
		b, err := json.Marshal(e.NewValues)
		if err != nil {
			return err
		}
		record.NewValues = string(b)
	}
	if e.Metadata != nil {
		b, err := json.Marshal(e.Metadata)
		if err != nil {
			return err
		}
		record.Metadata = string(b)
	}
	return l.db.Create(&record).Error
}

// Find returns audit log entries matching the given conditions, ordered newest first.
// Conditions use GORM Where syntax: Find("user_id = ?", 42).
func (l *Logger) Find(query string, args ...interface{}) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := l.db.Where(query, args...).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// ForModel returns all audit log entries for a specific model type and ID.
func (l *Logger) ForModel(modelType string, modelID uint) ([]models.AuditLog, error) {
	return l.Find("model_type = ? AND model_id = ?", modelType, modelID)
}
```

**API surface**:

| Function | Signature | Purpose |
|---|---|---|
| `NewLogger` | `NewLogger(db *gorm.DB) *Logger` | Constructor |
| `Log` | `(l *Logger) Log(e Entry) error` | Write an audit entry |
| `Find` | `(l *Logger) Find(query string, args ...interface{}) ([]AuditLog, error)` | Query audit logs with GORM conditions |
| `ForModel` | `(l *Logger) ForModel(modelType string, modelID uint) ([]AuditLog, error)` | Get all entries for a specific record |

---

### 3. Model Registry Update

**File**: `database/models/registry.go`

Add `&AuditLog{}` to the `All()` function so AutoMigrate picks it up:

```go
func All() []interface{} {
	return []interface{}{
		&User{},
		&Post{},
		&AuditLog{},
	}
}
```

---

### 4. Migration

**Package**: `database/migrations`
**File**: `20260308000003_create_audit_logs_table.go`

```go
func init() {
	Register(Migration{
		Version: "20260308000003_create_audit_logs_table",
		Up: func(db *gorm.DB) error {
			type AuditLog struct {
				ID        uint      `gorm:"primaryKey"`
				UserID    uint      `gorm:"index;not null;default:0"`
				Action    string    `gorm:"size:50;not null;index"`
				ModelType string    `gorm:"size:100;not null;index"`
				ModelID   uint      `gorm:"not null;index"`
				OldValues string    `gorm:"type:text"`
				NewValues string    `gorm:"type:text"`
				Metadata  string    `gorm:"type:text"`
				CreatedAt time.Time
			}
			return db.AutoMigrate(&AuditLog{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable("audit_logs")
		},
	})
}
```

---

## Usage Example (Application-Level)

The framework provides building blocks. An application would audit operations like:

```go
// In a service or controller:
auditor := audit.NewLogger(db)

// After creating a user:
auditor.Log(audit.Entry{
	UserID:    currentUserID,
	Action:    "create",
	ModelType: "User",
	ModelID:   newUser.ID,
	NewValues: map[string]interface{}{
		"name":  newUser.Name,
		"email": newUser.Email,
		"role":  newUser.Role,
	},
})

// After updating a post:
auditor.Log(audit.Entry{
	UserID:    currentUserID,
	Action:    "update",
	ModelType: "Post",
	ModelID:   post.ID,
	OldValues: map[string]interface{}{"title": oldTitle},
	NewValues: map[string]interface{}{"title": newTitle},
	Metadata:  map[string]interface{}{"reason": "typo fix"},
})

// Query audit history for a specific record:
logs, _ := auditor.ForModel("User", 42)

// Query by custom conditions:
logs, _ := auditor.Find("user_id = ? AND action = ?", adminID, "delete")
```

---

## Dependencies

- **Existing**: `gorm.io/gorm` — ORM for persistence (already in go.mod)
- **Existing**: `encoding/json` — Go standard library for JSON marshaling
- **Existing**: `database/models` — BaseModel pattern (ID, CreatedAt only — no embed)
- No new external dependencies

---

## Environment Variables

None. Audit logging has no framework-level configuration — it is always available when the `core/audit` package is imported. Application-level decisions (what to audit, retention policies) are left to the developer.

---

## Next

Tasks doc → `51-audit-logging-tasks.md`
