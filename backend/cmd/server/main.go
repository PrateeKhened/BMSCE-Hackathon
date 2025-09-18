package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prateekkhenedcodes/BMSCE-Hackathon/backend/internal/config"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// TODO: Initialize database connection
	// TODO: Run migrations
	// TODO: Initialize services
	// TODO: Setup routes
	// TODO: Setup middleware

	// Placeholder server setup
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "healthy", "service": "medical-report-backend"}`)
	})

	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}