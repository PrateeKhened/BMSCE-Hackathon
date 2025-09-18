package tests

import (
	"testing"
	"time"

	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/config"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/database"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/models"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/services"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/pkg/types"
)

// setupAuthTest creates test services and database
func setupAuthTest(t *testing.T) (*services.AuthService, *database.DB) {
	// Decision: Use in-memory database for isolated tests
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Driver: "sqlite3",
			DSN:    ":memory:",
		},
		JWT: config.JWTConfig{
			Secret:     "test-secret-key-for-testing-only",
			Expiration: time.Hour * 24,
		},
	}

	db, err := database.Setup(cfg)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	// Decision: Create tables for testing (in real app, migrations handle this)
	createTestTables(t, db)

	// Decision: Create service instances with test configuration
	passwordService := services.NewPasswordServiceWithCost(4) // Lower cost for faster tests
	jwtService := services.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)
	userRepo := models.NewUserRepository(db.GetDB())
	authService := services.NewAuthService(userRepo, passwordService, jwtService)

	return authService, db
}

// createTestTables creates necessary tables for testing
func createTestTables(t *testing.T, db *database.DB) {
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

	_, err := db.Exec(createUserTable)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}
}

// TestPasswordService tests password hashing functionality
func TestPasswordService(t *testing.T) {
	// Decision: Test password service in isolation
	passwordService := services.NewPasswordService()

	password := "test_password_123"

	// Test password hashing
	hash, err := passwordService.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if len(hash) == 0 {
		t.Fatal("Hash should not be empty")
	}

	// Test password verification - correct password
	if !passwordService.CheckPassword(password, hash) {
		t.Fatal("Password verification should succeed with correct password")
	}

	// Test password verification - incorrect password
	if passwordService.CheckPassword("wrong_password", hash) {
		t.Fatal("Password verification should fail with incorrect password")
	}

	// Test that same password produces different hashes (due to salt)
	hash2, err := passwordService.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password second time: %v", err)
	}

	if hash == hash2 {
		t.Fatal("Same password should produce different hashes due to salt")
	}

	t.Log("Password service test passed")
}

// TestJWTService tests JWT token functionality
func TestJWTService(t *testing.T) {
	jwtService := services.NewJWTService("test-secret", time.Hour)

	userID := 123
	email := "test@example.com"

	// Test token generation
	token, err := jwtService.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if len(token) == 0 {
		t.Fatal("Token should not be empty")
	}

	// Test token validation
	claims, err := jwtService.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Fatalf("Expected user ID %d, got %d", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Fatalf("Expected email %s, got %s", email, claims.Email)
	}

	// Test invalid token
	_, err = jwtService.ValidateToken("invalid-token")
	if err == nil {
		t.Fatal("Should fail to validate invalid token")
	}

	// Test token refresh
	newToken, err := jwtService.RefreshToken(token)
	if err != nil {
		t.Fatalf("Failed to refresh token: %v", err)
	}

	// Decision: Validate that refreshed token is valid and contains same user data
	// (Tokens might be identical if refreshed in same second, which is acceptable)
	newClaims, err := jwtService.ValidateToken(newToken)
	if err != nil {
		t.Fatalf("Refreshed token should be valid: %v", err)
	}

	if newClaims.UserID != userID || newClaims.Email != email {
		t.Fatal("Refreshed token should contain same user data")
	}

	t.Log("JWT service test passed")
}

// TestAuthServiceSignup tests user registration
func TestAuthServiceSignup(t *testing.T) {
	authService, db := setupAuthTest(t)
	defer db.Close()

	// Test successful signup
	signupReq := &types.SignupRequest{
		Email:    "newuser@example.com",
		Password: "secure_password_123",
		FullName: "New User",
	}

	response, err := authService.SignUp(signupReq)
	if err != nil {
		t.Fatalf("Signup should succeed: %v", err)
	}

	if response == nil {
		t.Fatal("Response should not be nil")
	}

	if response.Token == "" {
		t.Fatal("Token should be provided after signup")
	}

	if response.User.Email != signupReq.Email {
		t.Fatalf("Expected email %s, got %s", signupReq.Email, response.User.Email)
	}

	// Test duplicate email signup
	_, err = authService.SignUp(signupReq)
	if err == nil {
		t.Fatal("Should fail to signup with duplicate email")
	}

	// Test invalid email
	invalidReq := &types.SignupRequest{
		Email:    "invalid-email",
		Password: "secure_password_123",
		FullName: "Test User",
	}

	_, err = authService.SignUp(invalidReq)
	if err == nil {
		t.Fatal("Should fail to signup with invalid email")
	}

	// Test short password
	shortPasswordReq := &types.SignupRequest{
		Email:    "another@example.com",
		Password: "123", // Too short
		FullName: "Test User",
	}

	_, err = authService.SignUp(shortPasswordReq)
	if err == nil {
		t.Fatal("Should fail to signup with short password")
	}

	t.Log("Auth service signup test passed")
}

// TestAuthServiceLogin tests user authentication
func TestAuthServiceLogin(t *testing.T) {
	authService, db := setupAuthTest(t)
	defer db.Close()

	// First create a user
	signupReq := &types.SignupRequest{
		Email:    "loginuser@example.com",
		Password: "test_password_123",
		FullName: "Login User",
	}

	_, err := authService.SignUp(signupReq)
	if err != nil {
		t.Fatalf("Failed to create user for login test: %v", err)
	}

	// Test successful login
	loginReq := &types.LoginRequest{
		Email:    signupReq.Email,
		Password: signupReq.Password,
	}

	response, err := authService.Login(loginReq)
	if err != nil {
		t.Fatalf("Login should succeed: %v", err)
	}

	if response.Token == "" {
		t.Fatal("Token should be provided after login")
	}

	if response.User.Email != signupReq.Email {
		t.Fatalf("Expected email %s, got %s", signupReq.Email, response.User.Email)
	}

	// Test invalid email
	invalidEmailReq := &types.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "test_password_123",
	}

	_, err = authService.Login(invalidEmailReq)
	if err == nil {
		t.Fatal("Should fail to login with non-existent email")
	}

	// Test wrong password
	wrongPasswordReq := &types.LoginRequest{
		Email:    signupReq.Email,
		Password: "wrong_password",
	}

	_, err = authService.Login(wrongPasswordReq)
	if err == nil {
		t.Fatal("Should fail to login with wrong password")
	}

	t.Log("Auth service login test passed")
}

// TestAuthServiceTokenValidation tests token-based user retrieval
func TestAuthServiceTokenValidation(t *testing.T) {
	authService, db := setupAuthTest(t)
	defer db.Close()

	// Create and login user
	signupReq := &types.SignupRequest{
		Email:    "tokenuser@example.com",
		Password: "test_password_123",
		FullName: "Token User",
	}

	loginResponse, err := authService.SignUp(signupReq)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test valid token
	user, err := authService.GetUserFromToken(loginResponse.Token)
	if err != nil {
		t.Fatalf("Should validate valid token: %v", err)
	}

	if user.Email != signupReq.Email {
		t.Fatalf("Expected email %s, got %s", signupReq.Email, user.Email)
	}

	// Test invalid token
	_, err = authService.GetUserFromToken("invalid.token.here")
	if err == nil {
		t.Fatal("Should fail to validate invalid token")
	}

	// Test token refresh
	newToken, err := authService.RefreshToken(loginResponse.Token)
	if err != nil {
		t.Fatalf("Should refresh valid token: %v", err)
	}

	// Decision: Focus on functionality - refreshed token should be valid
	// (May be identical if refreshed immediately, which is acceptable)

	// Validate refreshed token works
	_, err = authService.GetUserFromToken(newToken)
	if err != nil {
		t.Fatalf("Refreshed token should be valid: %v", err)
	}

	t.Log("Auth service token validation test passed")
}