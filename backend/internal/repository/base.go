package repository

import (
	"gorm.io/gorm"
)

// BaseRepository provides common database operations
type BaseRepository struct {
	DB *gorm.DB
}

// NewBaseRepository creates a new BaseRepository
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{DB: db}
}

// WithTransaction executes a function within a database transaction
func (r *BaseRepository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return r.DB.Transaction(fn)
}