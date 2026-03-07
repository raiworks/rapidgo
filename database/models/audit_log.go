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
