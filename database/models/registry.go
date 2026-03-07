package models

// All returns all model structs for GORM AutoMigrate.
// Add new models here as they are created.
func All() []interface{} {
	return []interface{}{
		&User{},
		&Post{},
		&AuditLog{},
	}
}
