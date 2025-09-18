package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/config"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/database"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/handlers"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/middleware"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/models"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/router"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/services"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/pkg/types"
)

// setupTestServer creates a test HTTP server with all dependencies
func setupTestServer(t *testing.T) *httptest.Server {
	// Decision: Use in-memory database for isolated integration tests
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Driver: "sqlite3",
			DSN:    ":memory:",
		},
		JWT: config.JWTConfig{
			Secret:     "test-secret-key-for-integration-tests",
			Expiration: time.Hour * 24, // 24 hours for testing
		},
	}

	// Decision: Set up complete application stack
	db, err := database.Setup(cfg)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	// Decision: Create all tables for integration testing
	createAllTestTables(t, db)

	// Decision: Initialize all application layers
	userRepo := models.NewUserRepository(db.GetDB())
	reportRepo := models.NewReportRepository(db.GetDB())
	passwordService := services.NewPasswordServiceWithCost(4) // Faster for tests
	jwtService := services.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)
	authService := services.NewAuthService(userRepo, passwordService, jwtService)

	// Initialize AI service (can be nil for auth tests)
	var aiService *services.AIService

	authHandler := handlers.NewAuthHandler(authService)
	reportHandler := handlers.NewReportHandler(reportRepo, authService, aiService, "/tmp/test_uploads", 20971520)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Decision: Create router with all endpoints
	rt := router.NewRouter(authHandler, reportHandler, authMiddleware)
	httpRouter := rt.SetupRoutes()

	// Decision: Return test server for HTTP requests
	return httptest.NewServer(httpRouter)
}

// createAllTestTables creates all necessary tables for integration testing
func createAllTestTables(t *testing.T, db *database.DB) {
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

	createReportTable := `
		CREATE TABLE reports (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			original_filename TEXT NOT NULL,
			file_path TEXT NOT NULL,
			file_type TEXT NOT NULL,
			file_size INTEGER NOT NULL,
			processing_status TEXT DEFAULT 'pending',
			simplified_summary TEXT,
			upload_date DATETIME DEFAULT CURRENT_TIMESTAMP,
			processed_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`

	_, err = db.Exec(createReportTable)
	if err != nil {
		t.Fatalf("Failed to create reports table: %v", err)
	}
}

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// Decision: Test GET /health endpoint
	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to call health endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Decision: Verify response status and content type
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Fatal("Expected Content-Type: application/json")
	}

	// Decision: Parse and validate response body
	var healthResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&healthResponse); err != nil {
		t.Fatalf("Failed to parse health response: %v", err)
	}

	if healthResponse["status"] != "healthy" {
		t.Fatalf("Expected status 'healthy', got %v", healthResponse["status"])
	}

	t.Log("Health endpoint test passed")
}

// TestSignupEndpoint tests the user registration endpoint
func TestSignupEndpoint(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// Decision: Test valid signup request
	signupData := types.SignupRequest{
		Email:    "integration@example.com",
		Password: "integrationtest123",
		FullName: "Integration Test User",
	}

	jsonData, _ := json.Marshal(signupData)
	resp, err := http.Post(server.URL+"/api/auth/signup", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to call signup endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Decision: Verify successful signup (201 Created)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", resp.StatusCode)
	}

	// Decision: Parse response and verify token and user data
	var signupResponse types.LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&signupResponse); err != nil {
		t.Fatalf("Failed to parse signup response: %v", err)
	}

	if signupResponse.Token == "" {
		t.Fatal("Expected token in signup response")
	}

	if signupResponse.User.Email != signupData.Email {
		t.Fatalf("Expected email %s, got %s", signupData.Email, signupResponse.User.Email)
	}

	// Decision: Test duplicate email signup (should fail)
	resp2, err := http.Post(server.URL+"/api/auth/signup", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to call signup endpoint again: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusConflict {
		t.Fatalf("Expected status 409 for duplicate email, got %d", resp2.StatusCode)
	}

	t.Log("Signup endpoint test passed")
}

// TestLoginEndpoint tests the user login endpoint
func TestLoginEndpoint(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// Decision: First create a user via signup
	signupData := types.SignupRequest{
		Email:    "logintest@example.com",
		Password: "logintest123",
		FullName: "Login Test User",
	}

	jsonData, _ := json.Marshal(signupData)
	_, err := http.Post(server.URL+"/api/auth/signup", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create user for login test: %v", err)
	}

	// Decision: Test valid login
	loginData := types.LoginRequest{
		Email:    signupData.Email,
		Password: signupData.Password,
	}

	jsonData, _ = json.Marshal(loginData)
	resp, err := http.Post(server.URL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to call login endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var loginResponse types.LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}

	if loginResponse.Token == "" {
		t.Fatal("Expected token in login response")
	}

	// Decision: Test invalid login (wrong password)
	wrongLoginData := types.LoginRequest{
		Email:    signupData.Email,
		Password: "wrongpassword",
	}

	jsonData, _ = json.Marshal(wrongLoginData)
	resp2, err := http.Post(server.URL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to call login endpoint with wrong password: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected status 401 for wrong password, got %d", resp2.StatusCode)
	}

	t.Log("Login endpoint test passed")
}

// TestProtectedEndpoint tests the /me endpoint with authentication
func TestProtectedEndpoint(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// Decision: Create user and get token
	signupData := types.SignupRequest{
		Email:    "protected@example.com",
		Password: "protectedtest123",
		FullName: "Protected Test User",
	}

	jsonData, _ := json.Marshal(signupData)
	resp, err := http.Post(server.URL+"/api/auth/signup", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	defer resp.Body.Close()

	var signupResponse types.LoginResponse
	json.NewDecoder(resp.Body).Decode(&signupResponse)
	token := signupResponse.Token

	// Decision: Test /me endpoint with valid token
	req, _ := http.NewRequest("GET", server.URL+"/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp2, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to call /me endpoint: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 for /me endpoint, got %d", resp2.StatusCode)
	}

	var userResponse types.User
	if err := json.NewDecoder(resp2.Body).Decode(&userResponse); err != nil {
		t.Fatalf("Failed to parse /me response: %v", err)
	}

	if userResponse.Email != signupData.Email {
		t.Fatalf("Expected email %s, got %s", signupData.Email, userResponse.Email)
	}

	// Decision: Test /me endpoint without token (should fail)
	req3, _ := http.NewRequest("GET", server.URL+"/api/auth/me", nil)
	resp3, err := client.Do(req3)
	if err != nil {
		t.Fatalf("Failed to call /me endpoint without token: %v", err)
	}
	defer resp3.Body.Close()

	if resp3.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected status 401 without token, got %d", resp3.StatusCode)
	}

	t.Log("Protected endpoint test passed")
}

// TestCORSHeaders tests CORS functionality
func TestCORSHeaders(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// Decision: Test preflight OPTIONS request
	req, _ := http.NewRequest("OPTIONS", server.URL+"/api/auth/login", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send OPTIONS request: %v", err)
	}
	defer resp.Body.Close()

	// Decision: Verify CORS headers are present
	allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
	if allowOrigin == "" {
		t.Logf("Headers received: %+v", resp.Header)
		t.Fatal("Expected Access-Control-Allow-Origin header")
	}

	if resp.Header.Get("Access-Control-Allow-Methods") == "" {
		t.Fatal("Expected Access-Control-Allow-Methods header")
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 for OPTIONS request, got %d", resp.StatusCode)
	}

	t.Log("CORS headers test passed")
}