package audit

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/RAiWorks/RapidGo/database/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database with the AuditLog table.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&models.AuditLog{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

// ── T01: NewLogger returns a Logger ──────────────────────────────────────────

func TestNewLogger_ReturnsLogger(t *testing.T) {
	db := setupTestDB(t)
	l := NewLogger(db)
	if l == nil {
		t.Fatal("NewLogger() returned nil")
	}
}

// ── T02: Log creates an entry with action "create" ──────────────────────────

func TestLog_CreateAction(t *testing.T) {
	db := setupTestDB(t)
	l := NewLogger(db)

	err := l.Log(Entry{
		UserID:    1,
		Action:    "create",
		ModelType: "User",
		ModelID:   42,
	})
	if err != nil {
		t.Fatalf("Log() error: %v", err)
	}

	var record models.AuditLog
	if err := db.First(&record).Error; err != nil {
		t.Fatalf("failed to find record: %v", err)
	}
	if record.UserID != 1 {
		t.Errorf("UserID = %d, want 1", record.UserID)
	}
	if record.Action != "create" {
		t.Errorf("Action = %q, want 'create'", record.Action)
	}
	if record.ModelType != "User" {
		t.Errorf("ModelType = %q, want 'User'", record.ModelType)
	}
	if record.ModelID != 42 {
		t.Errorf("ModelID = %d, want 42", record.ModelID)
	}
}

// ── T03: Log stores OldValues and NewValues as JSON ─────────────────────────

func TestLog_UpdateWithOldNewValues(t *testing.T) {
	db := setupTestDB(t)
	l := NewLogger(db)

	err := l.Log(Entry{
		UserID:    1,
		Action:    "update",
		ModelType: "Post",
		ModelID:   10,
		OldValues: map[string]interface{}{"title": "Old Title"},
		NewValues: map[string]interface{}{"title": "New Title"},
	})
	if err != nil {
		t.Fatalf("Log() error: %v", err)
	}

	var record models.AuditLog
	db.First(&record)

	var oldVals, newVals map[string]interface{}
	if err := json.Unmarshal([]byte(record.OldValues), &oldVals); err != nil {
		t.Fatalf("OldValues is not valid JSON: %v", err)
	}
	if err := json.Unmarshal([]byte(record.NewValues), &newVals); err != nil {
		t.Fatalf("NewValues is not valid JSON: %v", err)
	}
	if oldVals["title"] != "Old Title" {
		t.Errorf("OldValues[title] = %v, want 'Old Title'", oldVals["title"])
	}
	if newVals["title"] != "New Title" {
		t.Errorf("NewValues[title] = %v, want 'New Title'", newVals["title"])
	}
}

// ── T04: Log stores delete action with OldValues ────────────────────────────

func TestLog_DeleteAction(t *testing.T) {
	db := setupTestDB(t)
	l := NewLogger(db)

	err := l.Log(Entry{
		UserID:    5,
		Action:    "delete",
		ModelType: "Post",
		ModelID:   99,
		OldValues: map[string]interface{}{"title": "Deleted Post"},
	})
	if err != nil {
		t.Fatalf("Log() error: %v", err)
	}

	var record models.AuditLog
	db.First(&record)
	if record.Action != "delete" {
		t.Errorf("Action = %q, want 'delete'", record.Action)
	}
	if record.OldValues == "" {
		t.Error("OldValues is empty, expected JSON")
	}
}

// ── T05: Log stores Metadata as JSON ────────────────────────────────────────

func TestLog_WithMetadata(t *testing.T) {
	db := setupTestDB(t)
	l := NewLogger(db)

	err := l.Log(Entry{
		UserID:    1,
		Action:    "update",
		ModelType: "User",
		ModelID:   7,
		Metadata:  map[string]interface{}{"ip": "192.168.1.1", "reason": "admin override"},
	})
	if err != nil {
		t.Fatalf("Log() error: %v", err)
	}

	var record models.AuditLog
	db.First(&record)

	var meta map[string]interface{}
	if err := json.Unmarshal([]byte(record.Metadata), &meta); err != nil {
		t.Fatalf("Metadata is not valid JSON: %v", err)
	}
	if meta["ip"] != "192.168.1.1" {
		t.Errorf("Metadata[ip] = %v, want '192.168.1.1'", meta["ip"])
	}
	if meta["reason"] != "admin override" {
		t.Errorf("Metadata[reason] = %v, want 'admin override'", meta["reason"])
	}
}

// ── T06: Nil maps store as empty strings ────────────────────────────────────

func TestLog_NilMapsStoreEmpty(t *testing.T) {
	db := setupTestDB(t)
	l := NewLogger(db)

	err := l.Log(Entry{
		UserID:    1,
		Action:    "create",
		ModelType: "User",
		ModelID:   1,
	})
	if err != nil {
		t.Fatalf("Log() error: %v", err)
	}

	var record models.AuditLog
	db.First(&record)
	if record.OldValues != "" {
		t.Errorf("OldValues = %q, want empty string", record.OldValues)
	}
	if record.NewValues != "" {
		t.Errorf("NewValues = %q, want empty string", record.NewValues)
	}
	if record.Metadata != "" {
		t.Errorf("Metadata = %q, want empty string", record.Metadata)
	}
}

// ── T07: UserID 0 (system action) persists ──────────────────────────────────

func TestLog_ZeroUserID(t *testing.T) {
	db := setupTestDB(t)
	l := NewLogger(db)

	err := l.Log(Entry{
		UserID:    0,
		Action:    "create",
		ModelType: "User",
		ModelID:   1,
	})
	if err != nil {
		t.Fatalf("Log() error: %v", err)
	}

	var record models.AuditLog
	db.First(&record)
	if record.UserID != 0 {
		t.Errorf("UserID = %d, want 0", record.UserID)
	}
}

// ── T08: Custom action persists ─────────────────────────────────────────────

func TestLog_CustomAction(t *testing.T) {
	db := setupTestDB(t)
	l := NewLogger(db)

	err := l.Log(Entry{
		UserID:    3,
		Action:    "login",
		ModelType: "User",
		ModelID:   3,
	})
	if err != nil {
		t.Fatalf("Log() error: %v", err)
	}

	var record models.AuditLog
	db.First(&record)
	if record.Action != "login" {
		t.Errorf("Action = %q, want 'login'", record.Action)
	}
}

// ── T09: Find by UserID ─────────────────────────────────────────────────────

func TestFind_ByUserID(t *testing.T) {
	db := setupTestDB(t)
	l := NewLogger(db)

	l.Log(Entry{UserID: 1, Action: "create", ModelType: "User", ModelID: 1})
	l.Log(Entry{UserID: 2, Action: "create", ModelType: "User", ModelID: 2})
	l.Log(Entry{UserID: 1, Action: "update", ModelType: "User", ModelID: 1})

	logs, err := l.Find("user_id = ?", 1)
	if err != nil {
		t.Fatalf("Find() error: %v", err)
	}
	if len(logs) != 2 {
		t.Errorf("got %d entries, want 2", len(logs))
	}
	for _, entry := range logs {
		if entry.UserID != 1 {
			t.Errorf("entry.UserID = %d, want 1", entry.UserID)
		}
	}
}

// ── T10: Find returns results ordered newest first ──────────────────────────

func TestFind_OrderedNewestFirst(t *testing.T) {
	db := setupTestDB(t)
	l := NewLogger(db)

	l.Log(Entry{UserID: 1, Action: "create", ModelType: "User", ModelID: 1})
	time.Sleep(10 * time.Millisecond) // ensure different timestamps
	l.Log(Entry{UserID: 1, Action: "update", ModelType: "User", ModelID: 1})

	logs, err := l.Find("user_id = ?", 1)
	if err != nil {
		t.Fatalf("Find() error: %v", err)
	}
	if len(logs) != 2 {
		t.Fatalf("got %d entries, want 2", len(logs))
	}
	if !logs[0].CreatedAt.After(logs[1].CreatedAt) && !logs[0].CreatedAt.Equal(logs[1].CreatedAt) {
		t.Error("results not ordered newest first")
	}
	if logs[0].Action != "update" {
		t.Errorf("first result Action = %q, want 'update' (newest)", logs[0].Action)
	}
}

// ── T11: Find returns empty slice when no match ─────────────────────────────

func TestFind_NoResults(t *testing.T) {
	db := setupTestDB(t)
	l := NewLogger(db)

	logs, err := l.Find("user_id = ?", 999)
	if err != nil {
		t.Fatalf("Find() error: %v", err)
	}
	if len(logs) != 0 {
		t.Errorf("got %d entries, want 0", len(logs))
	}
}

// ── T12: ForModel returns matching entries ──────────────────────────────────

func TestForModel_ReturnsMatchingEntries(t *testing.T) {
	db := setupTestDB(t)
	l := NewLogger(db)

	l.Log(Entry{UserID: 1, Action: "create", ModelType: "Post", ModelID: 5})
	l.Log(Entry{UserID: 1, Action: "update", ModelType: "Post", ModelID: 5})
	l.Log(Entry{UserID: 1, Action: "create", ModelType: "Post", ModelID: 6})

	logs, err := l.ForModel("Post", 5)
	if err != nil {
		t.Fatalf("ForModel() error: %v", err)
	}
	if len(logs) != 2 {
		t.Errorf("got %d entries, want 2", len(logs))
	}
}

// ── T13: ForModel distinguishes model types ─────────────────────────────────

func TestForModel_DifferentModelTypes(t *testing.T) {
	db := setupTestDB(t)
	l := NewLogger(db)

	l.Log(Entry{UserID: 1, Action: "create", ModelType: "User", ModelID: 1})
	l.Log(Entry{UserID: 1, Action: "create", ModelType: "Post", ModelID: 1})

	logs, err := l.ForModel("User", 1)
	if err != nil {
		t.Fatalf("ForModel() error: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("got %d entries, want 1", len(logs))
	}
	if logs[0].ModelType != "User" {
		t.Errorf("ModelType = %q, want 'User'", logs[0].ModelType)
	}
}

// ── T14: AuditLog has no DeletedAt field ────────────────────────────────────

func TestAuditLog_NoDeletedAt(t *testing.T) {
	typ := reflect.TypeOf(models.AuditLog{})
	if _, found := typ.FieldByName("DeletedAt"); found {
		t.Error("AuditLog has DeletedAt field — audit logs should be immutable")
	}
	if _, found := typ.FieldByName("UpdatedAt"); found {
		t.Error("AuditLog has UpdatedAt field — audit logs should be immutable")
	}
}
