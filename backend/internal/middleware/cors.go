package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

// CORSConfig holds CORS configuration
// Decision: Struct for flexible CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         int // Preflight cache time in seconds
}

// DefaultCORSConfig returns a development-friendly CORS configuration
// Decision: Permissive defaults for development, can be restricted in production
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins: []string{"*"}, // Decision: Allow all origins in development
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-Requested-With",
		},
		MaxAge: 86400, // Decision: Cache preflight requests for 24 hours
	}
}

// ProductionCORSConfig returns a production-safe CORS configuration
// Decision: Specific origins for production security
func ProductionCORSConfig(allowedOrigins []string) *CORSConfig {
	config := DefaultCORSConfig()
	config.AllowedOrigins = allowedOrigins // Decision: Specific domains only
	return config
}

// CORS creates a CORS middleware with the given configuration
// Decision: Return middleware function for flexible integration
func CORS(config *CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Decision: Check if origin is allowed and set appropriate header
			if isOriginAllowed(origin, config.AllowedOrigins) {
				// Decision: For wildcard, set "*"; for specific origins, set the origin
				if hasWildcard(config.AllowedOrigins) {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else if origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}

				// Decision: Set CORS headers only for allowed origins
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
				w.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
			}

			// Decision: Handle preflight OPTIONS requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Decision: Continue to next handler for non-preflight requests
			next.ServeHTTP(w, r)
		})
	}
}

// isOriginAllowed checks if an origin is in the allowed list
// Decision: Support wildcard "*" or specific origin matching
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		// Decision: "*" allows all origins (development only)
		if allowed == "*" {
			return true
		}
		// Decision: Exact match for security (requires non-empty origin)
		if origin != "" && allowed == origin {
			return true
		}
	}

	return false
}

// hasWildcard checks if the allowed origins contains "*"
// Decision: Helper function to check for wildcard configuration
func hasWildcard(allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" {
			return true
		}
	}
	return false
}