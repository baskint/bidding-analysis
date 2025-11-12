package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool - optimized for connection poolers like Neon
	db.SetMaxOpenConns(10)                 // Reduced from 25
	db.SetMaxIdleConns(2)                  // Reduced from 5
	db.SetConnMaxLifetime(5 * time.Minute) // Changed from 1 hour
	db.SetConnMaxIdleTime(2 * time.Minute) // NEW: Close idle connections faster

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
