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
