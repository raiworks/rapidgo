package database

import "gorm.io/gorm"

// TxFunc is the callback signature for transactional operations.
// Return nil to commit, return an error to rollback.
type TxFunc func(tx *gorm.DB) error

// WithTransaction wraps fn in a database transaction.
// If fn returns nil the transaction commits.
// If fn returns an error or panics the transaction rolls back.
func WithTransaction(db *gorm.DB, fn TxFunc) error {
	return db.Transaction(fn)
}
