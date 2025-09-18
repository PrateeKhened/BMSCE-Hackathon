package tests

import (
	"testing"

	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/config"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/database"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/models"
)

// TestDatabaseConnection tests basic database connectivity
// Decision: Test database connection to catch configuration issues early
func TestDatabaseConnection(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Driver: "sqlite3",
			DSN:    ":memory:", // Use in-memory database for testing
		},
	}

	db, err := database.Setup(cfg)
	if err != nil {
		t.Fatalf("Failed to setup database: %v", err)
	}
	defer db.Close()

	// Test connection with a simple query
	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		t.Fatalf("Database connection test failed: %v", err)
	}

	if result != 1 {
		t.Fatalf("Expected 1, got %d", result)
	}

	t.Log("Database connection test passed")
}

// TestUserModel tests basic user model operations
// Decision: Test user model since it's the foundation for authentication
func TestUserModel(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Driver: "sqlite3",
			DSN:    ":memory:",
		},
	}

	db, err := database.Setup(cfg)
	if err != nil {
		t.Fatalf("Failed to setup database: %v", err)
	}
	defer db.Close()

	// Create tables manually for testing (in real app, migrations handle this)
	createUserTable := `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			full_name TEXT NOT NULL,
			email_verified BOOLEAN DEFAULT FALSE,
			is_active BOOLEAN DEFAULT TRUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`

	_, err = db.Exec(createUserTable)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	// Test user repository
	repo := models.NewUserRepository(db.GetDB())

	// Test user creation
	user := &models.User{
		Email:        "test@example.com",
		PasswordHash: "hashed_password_123",
		FullName:     "Test User",
		EmailVerified: false,
		IsActive:     true,
	}

	err = repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Verify user was created with ID
	if user.ID == 0 {
		t.Fatal("User ID should be set after creation")
	}

	// Test getting user by ID
	retrievedUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}

	if retrievedUser == nil {
		t.Fatal("User should be found")
	}

	if retrievedUser.Email != user.Email {
		t.Fatalf("Expected email %s, got %s", user.Email, retrievedUser.Email)
	}

	// Test getting user by email
	userByEmail, err := repo.GetByEmail(user.Email)
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}

	if userByEmail == nil {
		t.Fatal("User should be found by email")
	}

	if userByEmail.ID != user.ID {
		t.Fatalf("Expected user ID %d, got %d", user.ID, userByEmail.ID)
	}

	t.Log("User model test passed")
}