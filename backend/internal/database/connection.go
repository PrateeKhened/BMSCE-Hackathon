package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

// DB holds our database connection
type DB struct {
	*sql.DB
}

// NewConnection creates a new database connection
// Decision: Using sql.DB directly wrapped in our struct for better control
func NewConnection(driverName, dataSourceName string) (*DB, error) {
	// Decision: Using sql.Open instead of a higher-level ORM for:
	// 1. Better performance and control
	// 2. Simpler debugging
	// 3. No additional learning curve for team
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Decision: Test the connection immediately to fail fast
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Decision: Configure connection pool for performance
	// Max open connections - prevents resource exhaustion
	db.SetMaxOpenConns(25)
	// Max idle connections - keeps connections ready for reuse
	db.SetMaxIdleConns(25)

	log.Printf("Database connection established successfully")

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// GetDB returns the underlying sql.DB for advanced operations
// Decision: Expose underlying DB for migrations and complex queries
func (db *DB) GetDB() *sql.DB {
	return db.DB
}