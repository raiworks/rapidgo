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
