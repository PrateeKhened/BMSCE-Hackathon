package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/models"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/services"
)

// UserContextKey is the key for storing user in request context
// Decision: Use custom type to avoid context key collisions
type UserContextKey string

const (
	UserKey UserContextKey = "user"
)

// AuthMiddleware provides JWT authentication middleware
// Decision: Use struct to inject auth service dependency
type AuthMiddleware struct {
	authService *services.AuthService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth is middleware that requires valid JWT authentication
// Decision: Return middleware function for flexible use with different routes
func (am *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Decision: Extract token from Authorization header
		token := extractBearerToken(r)
		if token == "" {
			writeUnauthorizedResponse(w, "Authorization token required")
			return
		}

		// Decision: Validate token and get user information
		user, err := am.authService.GetUserFromToken(token)
		if err != nil {
			writeUnauthorizedResponse(w, "Invalid or expired token")
			return
		}

		// Decision: Check if user account is still active
		if !user.IsActive {
			writeUnauthorizedResponse(w, "Account is deactivated")
			return
		}

		// Decision: Add user to request context for handlers to use
		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth is middleware that extracts user if token is present but doesn't require it
// Decision: Useful for endpoints that behave differently for authenticated users
func (am *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r)
		if token != "" {
			// Decision: Only add user to context if token is valid
			if user, err := am.authService.GetUserFromToken(token); err == nil && user.IsActive {
				ctx := context.WithValue(r.Context(), UserKey, user)
				r = r.WithContext(ctx)
			}
		}

		// Decision: Always continue to next handler regardless of auth status
		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext extracts the authenticated user from request context
// Decision: Utility function for handlers to easily get current user
func GetUserFromContext(r *http.Request) (*models.User, bool) {
	user, ok := r.Context().Value(UserKey).(*models.User)
	return user, ok
}

// extractBearerToken extracts JWT token from Authorization header
// Decision: Support standard "Bearer <token>" format
func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Decision: Split on space and check for "Bearer" prefix
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

// writeUnauthorizedResponse writes a standardized unauthorized response
// Decision: Consistent error format across all auth failures
func writeUnauthorizedResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusUnauthorized)

	response := `{"error": true, "message": "` + message + `", "status": 401}`
	w.Write([]byte(response))
}