package database

import "gorm.io/gorm"

// Resolver holds separate database connections for write and read operations.
type Resolver struct {
	writer *gorm.DB
	reader *gorm.DB
}

// NewResolver creates a Resolver with the given writer and reader connections.
// If reads should go to the same database, pass the same *gorm.DB for both.
func NewResolver(writer, reader *gorm.DB) *Resolver {
	return &Resolver{writer: writer, reader: reader}
}

// Writer returns the write (primary) database connection.
func (r *Resolver) Writer() *gorm.DB {
	return r.writer
}

// Reader returns the read (replica) database connection.
func (r *Resolver) Reader() *gorm.DB {
	return r.reader
}
