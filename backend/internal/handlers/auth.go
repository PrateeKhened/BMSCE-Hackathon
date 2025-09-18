package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/services"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/pkg/errors"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/pkg/types"
)

// AuthHandler handles authentication HTTP requests
// Decision: Use struct to group related handlers and inject dependencies
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// SignupHandler handles user registration requests
// POST /api/auth/signup
func (ah *AuthHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	// Decision: Only allow POST method for signup
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Decision: Parse JSON request body
	var req types.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Decision: Call authentication service for business logic
	response, err := ah.authService.SignUp(&req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Decision: Return 201 Created for successful user creation
	writeJSONResponse(w, http.StatusCreated, response)
}

// LoginHandler handles user authentication requests
// POST /api/auth/login
func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Decision: Only allow POST method for login
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Decision: Parse JSON request body
	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Decision: Call authentication service
	response, err := ah.authService.Login(&req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Decision: Return 200 OK for successful login
	writeJSONResponse(w, http.StatusOK, response)
}

// LogoutHandler handles user logout requests
// POST /api/auth/logout
// Decision: For now, logout is client-side (delete token). In future, could blacklist tokens.
func (ah *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Decision: Return success message for logout
	// Client should delete the token from storage
	response := types.AuthResponse{
		Message: "Logged out successfully",
		Success: true,
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// MeHandler returns current user information from JWT token
// GET /api/auth/me
func (ah *AuthHandler) MeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Decision: Extract token from Authorization header
	token := extractTokenFromHeader(r)
	if token == "" {
		writeErrorResponse(w, http.StatusUnauthorized, "Authorization token required")
		return
	}

	// Decision: Get user from token using auth service
	user, err := ah.authService.GetUserFromToken(token)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Decision: Return user information (password hash excluded by JSON tag)
	writeJSONResponse(w, http.StatusOK, user)
}

// RefreshHandler generates a new JWT token for valid existing token
// POST /api/auth/refresh
func (ah *AuthHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Decision: Extract token from Authorization header
	token := extractTokenFromHeader(r)
	if token == "" {
		writeErrorResponse(w, http.StatusUnauthorized, "Authorization token required")
		return
	}

	// Decision: Generate new token
	newToken, err := ah.authService.RefreshToken(token)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Decision: Return new token in same format as login
	response := map[string]interface{}{
		"token":   newToken,
		"message": "Token refreshed successfully",
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// extractTokenFromHeader extracts JWT token from Authorization header
// Decision: Support "Bearer <token>" format
func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Decision: Check for "Bearer " prefix
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

// handleServiceError converts service errors to HTTP responses
// Decision: Map custom errors to appropriate HTTP status codes
func handleServiceError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		writeErrorResponse(w, appErr.Code, appErr.Message)
		return
	}

	// Decision: Default to internal server error for unknown errors
	writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
}

// writeJSONResponse writes a JSON response
// Decision: Set proper headers for JSON and CORS
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Decision: Allow CORS for frontend
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Decision: Log error and write minimal error response
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

// writeErrorResponse writes an error response
// Decision: Consistent error format across all endpoints
func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	errorResponse := map[string]interface{}{
		"error":   true,
		"message": message,
		"status":  statusCode,
	}

	writeJSONResponse(w, statusCode, errorResponse)
}