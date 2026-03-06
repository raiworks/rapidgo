package database

import (
	"gorm.io/gorm"
)

// TransferCredits atomically moves credits from one user to another.
// Both users must exist and the operation uses gorm.Expr for race-safe SQL.
func TransferCredits(db *gorm.DB, fromID, toID uint, amount int) error {
	return WithTransaction(db, func(tx *gorm.DB) error {
		// Verify source user exists
		if err := tx.Table("users").First(&struct{ ID uint }{}, fromID).Error; err != nil {
			return err
		}
		// Verify destination user exists
		if err := tx.Table("users").First(&struct{ ID uint }{}, toID).Error; err != nil {
			return err
		}
		// Deduct from source
		if err := tx.Table("users").Where("id = ?", fromID).
			Update("credits", gorm.Expr("credits - ?", amount)).Error; err != nil {
			return err
		}
		// Add to destination
		if err := tx.Table("users").Where("id = ?", toID).
			Update("credits", gorm.Expr("credits + ?", amount)).Error; err != nil {
			return err
		}
		return nil
	})
}
