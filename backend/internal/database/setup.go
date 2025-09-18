package database

import (
	"fmt"
	"log"

	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/config"
)

// Setup initializes the database connection and returns a DB instance
// Decision: Centralized database setup function for consistent initialization
func Setup(cfg *config.Config) (*DB, error) {
	// Decision: Log connection attempt for debugging
	log.Printf("Connecting to database: driver=%s, dsn=%s", cfg.Database.Driver, cfg.Database.DSN)

	db, err := NewConnection(cfg.Database.Driver, cfg.Database.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to setup database: %w", err)
	}

	// Decision: Enable foreign key constraints for SQLite
	// This ensures referential integrity between tables
	if cfg.Database.Driver == "sqlite3" {
		if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
			return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
		}
		log.Println("Foreign key constraints enabled")
	}

	log.Println("Database setup completed successfully")
	return db, nil
}