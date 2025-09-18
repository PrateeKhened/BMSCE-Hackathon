package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/config"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/database"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/handlers"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/middleware"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/models"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/router"
	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/services"
)

func main() {
	// Decision: Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
		log.Printf("Using system environment variables")
	}

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
	reportRepo := models.NewReportRepository(db.GetDB())

	// Decision: Initialize services (business logic layer)
	passwordService := services.NewPasswordService()
	jwtService := services.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)
	authService := services.NewAuthService(userRepo, passwordService, jwtService)

	// Initialize AI service for Gemini integration
	aiService, err := services.NewAIService(cfg.AI.GeminiAPIKey)
	if err != nil {
		log.Printf("Warning: AI service initialization failed: %v", err)
		log.Printf("Report analysis will not be available")
	}
	defer func() {
		if aiService != nil {
			aiService.Close()
		}
	}()

	// Decision: Initialize handlers (HTTP layer)
	authHandler := handlers.NewAuthHandler(authService)
	reportHandler := handlers.NewReportHandler(reportRepo, authService, aiService, cfg.Upload.UploadPath, cfg.Upload.MaxFileSize)

	// Decision: Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Decision: Setup router with all dependencies
	rt := router.NewRouter(authHandler, reportHandler, authMiddleware)
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
	log.Println("  GET  /health                    - Health check")
	log.Println("  POST /api/auth/signup           - User registration")
	log.Println("  POST /api/auth/login            - User login")
	log.Println("  POST /api/auth/logout           - User logout")
	log.Println("  GET  /api/auth/me               - Get current user (requires auth)")
	log.Println("  POST /api/auth/refresh          - Refresh JWT token (requires auth)")
	log.Println("  GET  /api/reports               - Get user's reports (requires auth)")
	log.Println("  POST /api/reports               - Upload medical report (requires auth)")
	log.Println("  GET  /api/reports/{id}          - Get specific report (requires auth)")
	log.Println("  DELETE /api/reports/{id}        - Delete report (requires auth)")
	log.Println("  GET  /api/reports/{id}/summary  - Get AI analysis summary (requires auth)")
	log.Println("  GET  /api/reports/{id}/metrics  - Get health metrics for speedometer (requires auth)")

	log.Printf("Server ready and listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}