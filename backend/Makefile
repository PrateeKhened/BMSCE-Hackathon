# Medical Report Simplifier Backend Makefile

# Variables
BINARY_NAME=medical-report-server
MAIN_PATH=./cmd/server
MIGRATION_PATH=./migrations
DB_DSN=./medical_reports.db

# Go commands
.PHONY: help build run clean test deps migrate-up migrate-down migrate-status

help: ## Display available commands
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	CGO_ENABLED=1 go build -o $(BINARY_NAME) $(MAIN_PATH)

run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	go run $(MAIN_PATH)

dev: ## Run with hot reload (install air first: go install github.com/cosmtrek/air@latest)
	@echo "Starting development server with hot reload..."
	air

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f $(DB_DSN)

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Database migration commands
migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	goose -dir $(MIGRATION_PATH) sqlite3 $(DB_DSN) up

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	goose -dir $(MIGRATION_PATH) sqlite3 $(DB_DSN) down

migrate-status: ## Check migration status
	@echo "Checking migration status..."
	goose -dir $(MIGRATION_PATH) sqlite3 $(DB_DSN) status

migrate-create: ## Create new migration (usage: make migrate-create NAME=create_users_table)
	@echo "Creating migration: $(NAME)..."
	goose -dir $(MIGRATION_PATH) create $(NAME) sql

# Development helpers
fmt: ## Format code
	go fmt ./...

lint: ## Run golangci-lint
	golangci-lint run

init-db: ## Initialize database with migrations
	@echo "Initializing database..."
	mkdir -p uploads
	make migrate-up

# Docker commands (future)
docker-build: ## Build Docker image
	docker build -t medical-report-backend .

docker-run: ## Run Docker container
	docker run -p 8080:8080 medical-report-backend

# Setup project for first time
setup: deps init-db ## Setup project for first time development
	@echo "Project setup complete!"
	@echo "Run 'make run' to start the server"