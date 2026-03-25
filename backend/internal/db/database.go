package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database represents a database connection with GORM
type Database struct {
	*gorm.DB
}

// Config holds database configuration
type Config struct {
	URL             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// NewConnection creates a new database connection with connection pooling
func NewConnection(url string) (*Database, error) {
	cfg := Config{
		URL:             url,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
	}

	return NewConnectionWithConfig(cfg)
}

// NewConnectionWithConfig creates a new database connection with custom configuration
func NewConnectionWithConfig(cfg Config) (*Database, error) {
	// Create GORM configuration
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Open database connection
	gormDB, err := gorm.Open(postgres.Open(cfg.URL), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Get underlying SQL database
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database connection established (idle: %d, open: %d, lifetime: %v)",
		cfg.MaxIdleConns, cfg.MaxOpenConns, cfg.ConnMaxLifetime)

	return &Database{gormDB}, nil
}

// Close closes the database connection
func (db *Database) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB for closing: %w", err)
	}

	return sqlDB.Close()
}

// Ping checks if database is reachable
func (db *Database) Ping() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB for ping: %w", err)
	}

	return sqlDB.Ping()
}

// Stats returns database connection statistics
func (db *Database) Stats() map[string]interface{} {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return map[string]interface{}{
			"error": "failed to get sql.DB for stats",
		}
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}

// WithTransaction executes a function within a database transaction
func (db *Database) WithTransaction(fn func(tx *gorm.DB) error) error {
	return db.DB.Transaction(fn)
}
