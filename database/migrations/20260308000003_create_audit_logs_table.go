package migrations

import (
	"time"

	"gorm.io/gorm"
)

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
