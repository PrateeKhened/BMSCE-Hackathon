package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/handlers"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/middleware"
)

// Router holds all router dependencies
// Decision: Struct to organize handlers and middleware
type Router struct {
	authHandler    *handlers.AuthHandler
	authMiddleware *middleware.AuthMiddleware
}

// NewRouter creates a new router with all dependencies
func NewRouter(
	authHandler *handlers.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		authHandler:    authHandler,
		authMiddleware: authMiddleware,
	}
}

// SetupRoutes configures all routes and returns the main router
// Decision: Single function to configure all application routes
func (rt *Router) SetupRoutes() *mux.Router {
	// Decision: Create main router with CORS middleware
	r := mux.NewRouter()

	// Decision: Apply CORS middleware to all routes
	corsMiddleware := middleware.CORS(middleware.DefaultCORSConfig())
	r.Use(corsMiddleware)

	// Decision: Health check endpoint (no auth required)
	r.HandleFunc("/health", rt.healthHandler).Methods("GET", "OPTIONS")

	// Decision: Create API subrouter for versioning
	api := r.PathPrefix("/api").Subrouter()

	// Decision: Setup authentication routes
	rt.setupAuthRoutes(api)

	// Decision: Future route groups will be added here
	// rt.setupReportRoutes(api)
	// rt.setupChatRoutes(api)

	return r
}

// setupAuthRoutes configures authentication endpoints
// Decision: Group auth routes under /api/auth prefix
func (rt *Router) setupAuthRoutes(api *mux.Router) {
	auth := api.PathPrefix("/auth").Subrouter()

	// Decision: Public authentication endpoints (no middleware required)
	auth.HandleFunc("/signup", rt.authHandler.SignupHandler).Methods("POST", "OPTIONS")
	auth.HandleFunc("/login", rt.authHandler.LoginHandler).Methods("POST", "OPTIONS")
	auth.HandleFunc("/logout", rt.authHandler.LogoutHandler).Methods("POST", "OPTIONS")

	// Decision: Protected authentication endpoints (require valid JWT)
	protectedAuth := auth.PathPrefix("").Subrouter()
	protectedAuth.Use(rt.authMiddleware.RequireAuth)
	protectedAuth.HandleFunc("/me", rt.authHandler.MeHandler).Methods("GET", "OPTIONS")
	protectedAuth.HandleFunc("/refresh", rt.authHandler.RefreshHandler).Methods("POST", "OPTIONS")
}

// healthHandler provides application health status
// Decision: Simple health check for load balancers and monitoring
func (rt *Router) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Decision: Include service name and status for identification
	response := `{
		"status": "healthy",
		"service": "medical-report-backend",
		"version": "1.0.0"
	}`

	w.Write([]byte(response))
}

// Future route setup methods will be added here:

// setupReportRoutes will configure report management endpoints
// func (rt *Router) setupReportRoutes(api *mux.Router) {
//     reports := api.PathPrefix("/reports").Subrouter()
//     reports.Use(rt.authMiddleware.RequireAuth) // All report routes require auth
//
//     reports.HandleFunc("", rt.reportHandler.ListReports).Methods("GET")
//     reports.HandleFunc("", rt.reportHandler.UploadReport).Methods("POST")
//     reports.HandleFunc("/{id}", rt.reportHandler.GetReport).Methods("GET")
//     reports.HandleFunc("/{id}", rt.reportHandler.DeleteReport).Methods("DELETE")
//     reports.HandleFunc("/{id}/summary", rt.reportHandler.GetSummary).Methods("GET")
// }

// setupChatRoutes will configure chat endpoints
// func (rt *Router) setupChatRoutes(api *mux.Router) {
//     chat := api.PathPrefix("/reports/{reportId}/chat").Subrouter()
//     chat.Use(rt.authMiddleware.RequireAuth) // All chat routes require auth
//
//     chat.HandleFunc("", rt.chatHandler.SendMessage).Methods("POST")
//     chat.HandleFunc("", rt.chatHandler.GetHistory).Methods("GET")
//     chat.HandleFunc("/{messageId}", rt.chatHandler.DeleteMessage).Methods("DELETE")
// }