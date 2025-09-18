package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/config"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/database"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/handlers"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/middleware"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/models"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/router"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/services"
)

func main() {
	// Decision: Load configuration from environment
	cfg := config.Load()
	log.Printf("Starting Medical Report Backend on %s:%s", cfg.Server.Host, cfg.Server.Port)

	// Decision: Initialize database connection
	db, err := database.Setup(cfg)
	if err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}
	defer db.Close()

	// Decision: Initialize repositories (data layer)
	userRepo := models.NewUserRepository(db.GetDB())

	// Decision: Initialize services (business logic layer)
	passwordService := services.NewPasswordService()
	jwtService := services.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)
	authService := services.NewAuthService(userRepo, passwordService, jwtService)

	// Decision: Initialize handlers (HTTP layer)
	authHandler := handlers.NewAuthHandler(authService)

	// Decision: Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Decision: Setup router with all dependencies
	rt := router.NewRouter(authHandler, authMiddleware)
	httpRouter := rt.SetupRoutes()

	// Decision: Configure HTTP server with timeouts
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      httpRouter,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Decision: Log available endpoints for development
	log.Println("Available endpoints:")
	log.Println("  GET  /health              - Health check")
	log.Println("  POST /api/auth/signup     - User registration")
	log.Println("  POST /api/auth/login      - User login")
	log.Println("  POST /api/auth/logout     - User logout")
	log.Println("  GET  /api/auth/me         - Get current user (requires auth)")
	log.Println("  POST /api/auth/refresh    - Refresh JWT token (requires auth)")

	log.Printf("Server ready and listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}