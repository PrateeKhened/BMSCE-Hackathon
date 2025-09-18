# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an AI-powered Medical Report Simplifier with:
- **Frontend**: Web application (managed by frontend developer)
- **Backend**: Go REST API for authentication, file upload, AI processing, and chat functionality

## Development Commands

### Quick Start
```bash
cd backend
make setup          # First-time setup: deps + database initialization
make run            # Start development server
```

### Core Development
```bash
make deps           # Download Go dependencies
make build          # Build production binary
make test           # Run all tests
make test-coverage  # Run tests with HTML coverage report
make fmt            # Format Go code
make lint           # Run golangci-lint
```

### Database Operations
```bash
make migrate-up     # Apply pending migrations
make migrate-down   # Rollback last migration
make migrate-status # Check migration status
make migrate-create NAME=migration_name  # Create new migration
make init-db        # Initialize database (creates uploads/ and runs migrations)
```

### Testing Specific Components
```bash
go test ./tests/ -v                    # Run database/integration tests
go test ./internal/models/ -v          # Test database models
go test -run TestUserModel ./tests/ -v # Run single test function
CGO_ENABLED=1 go test ./tests/ -v      # Ensure SQLite driver works in tests
```

## Architecture Overview

### Repository Pattern Implementation
- **Models** (`internal/models/`): Repository interfaces + SQL implementations
- **User Model**: CRUD operations with email-based auth, soft delete capability
- **Report Model**: File metadata, AI processing status tracking, user associations
- **Chat Message Model**: Conversation history with soft delete

### Database Design
- **SQLite** with Goose migrations in `migrations/`
- **Foreign key constraints** enabled for data integrity
- **Strategic indexing** for common queries (email lookups, user reports, chat history)
- **Soft deletes** for users and chat messages, hard delete for reports

### Configuration
- Environment-based config in `internal/config/config.go`
- Defaults for development, environment variables for production
- Database, JWT, file upload, and server settings centralized

### API Structure (Planned)
```
/api/auth/*     - Authentication endpoints
/api/reports/*  - File upload and report management
/api/reports/{id}/chat/* - Chat functionality per report
/health         - Health check endpoint
```

### Error Handling
- Custom error types in `pkg/errors/` with HTTP status codes
- Repository pattern returns `nil` for not-found (not errors)
- Structured error responses across API

## Development Workflow

### Database Schema Changes
1. Create migration: `make migrate-create NAME=descriptive_name`
2. Edit generated SQL file in `migrations/`
3. Test: `make migrate-up` then `make migrate-down`
4. Update corresponding model in `internal/models/` if needed

### Testing Strategy
- **Unit tests**: Repository operations, business logic
- **Integration tests**: Database connectivity, model operations
- **In-memory SQLite** for test isolation
- Test files in `tests/` directory, can be organized by component

### Implementation Status
- ✅ **Database Foundation**: Complete with migrations, repository pattern, and tests
  - 3 tables: users, reports, chat_messages with proper relationships
  - Foreign key constraints and strategic indexing
  - Repository interfaces with full CRUD operations
  - Comprehensive test coverage with in-memory SQLite
- ✅ **Authentication System**: Complete (JWT + bcrypt)
  - User registration, login, logout, token refresh
  - Secure password hashing with bcrypt
  - JWT middleware for protected endpoints
- ✅ **File Upload System**: Complete with validation and processing
  - Support for PDF, TXT, DOCX files (up to 20MB)
  - Secure file storage with unique naming
  - File type validation and metadata tracking
- ✅ **AI Integration**: Complete Gemini API integration
  - Real-time medical report analysis
  - Health metrics extraction with 0-100 scoring for speedometer display
  - Configurable prompt system in `prompts/medical_analysis_prompt.txt`
  - Comprehensive error handling and JSON parsing
- 📋 **Chat Functionality**: Planned (AI conversation about reports)

### Git Workflow
- **Frequent commits** for each small feature with full tests
- **Backend-only changes** (frontend managed separately)
- Always `git pull` before `git push` to integrate frontend updates
- Commit message format: `feat: descriptive message` with implementation details

## Key Files

- `Makefile`: All development commands and workflows
- `internal/database/setup.go`: Database connection with foreign key constraints
- `internal/database/connection.go`: Database connection pooling and configuration
- `migrations/*.sql`: Goose migration files for schema changes
- `internal/models/*.go`: Repository pattern implementations (user.go, report.go, chat_message.go)
- `internal/config/config.go`: Environment-based configuration management
- `internal/services/ai_service.go`: Gemini AI integration with configurable prompts
- `internal/handlers/report.go`: Complete report management with file upload and AI processing
- `prompts/medical_analysis_prompt.txt`: **Configurable AI system prompt** for medical analysis
- `prompts/README.md`: Documentation for modifying AI prompts
- `pkg/types/*.go`: Shared data structures and request/response types
- `pkg/errors/errors.go`: Custom error types with HTTP status codes
- `tests/database_test.go`: Database connectivity and model operation tests
- `cmd/server/main.go`: Main application entry point with full service integration

## Important Implementation Details

### Database Connection
- SQLite with CGO_ENABLED=1 requirement for compilation
- Connection pooling configured (25 max open/idle connections)
- Foreign key constraints enabled via PRAGMA
- In-memory database (`:memory:`) used for testing

### Repository Pattern Structure
```go
// Each model has interface + SQL implementation
type UserRepository interface { ... }
type SQLUserRepository struct { db *sql.DB }

// Repositories return nil for not-found (not errors)
func (r *SQLUserRepository) GetByEmail(email string) (*User, error)
```

### Migration Management
- Goose tracks applied migrations in `goose_db_version` table
- Migration files follow timestamp naming: `YYYYMMDDHHMMSS_description.sql`
- Both UP and DOWN migrations required for each file