# Medical Report Simplifier - Backend

AI-powered medical report simplifier that converts complex medical test reports into patient-friendly explanations.

## Quick Start

```bash
# Clone and setup
cd backend
make setup          # Install dependencies and initialize database

# Development
make run            # Start development server
make test           # Run all tests

# Database operations
make migrate-status # Check current migration status
make migrate-up     # Apply new migrations
```

## Implementation Status

### ✅ Phase 1 - Database Foundation (Complete)
- [x] Project structure with clean architecture
- [x] SQLite database with Goose migrations
- [x] Repository pattern with full CRUD operations
- [x] Database models for users, reports, and chat messages
- [x] Comprehensive test coverage
- [x] Foreign key constraints and strategic indexing

### 🔄 Phase 2 - Authentication & API (In Progress)
- [ ] User authentication (signup/login/logout)
- [ ] JWT token management with bcrypt password hashing
- [ ] HTTP handlers and middleware
- [ ] File upload endpoint with validation

### 📋 Phase 3 - AI Integration (Planned)
- [ ] AI service integration for report processing
- [ ] Medical report parsing and simplification
- [ ] Chat functionality with uploaded reports
- [ ] Processing status tracking

### 🎯 Phase 4 - Dashboard & Analytics (Future)
- [ ] Health metrics extraction and tracking
- [ ] Report history with visual timeline
- [ ] User dashboard with health insights

## Technology Stack

- **Language**: Go 1.21+
- **Database**: SQLite
- **Migration**: Goose
- **Authentication**: JWT
- **File Storage**: Local filesystem
- **AI Integration**: TBD (OpenAI/Claude/Local LLM)

## Database Schema

### Users Table
```sql
id INTEGER PRIMARY KEY AUTOINCREMENT
email TEXT UNIQUE NOT NULL
password_hash TEXT NOT NULL
full_name TEXT NOT NULL
email_verified BOOLEAN DEFAULT FALSE
is_active BOOLEAN DEFAULT TRUE
created_at DATETIME DEFAULT CURRENT_TIMESTAMP
updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
```

### Reports Table
```sql
id INTEGER PRIMARY KEY AUTOINCREMENT
user_id INTEGER NOT NULL (FK → users.id)
original_filename TEXT NOT NULL
file_path TEXT NOT NULL
file_type TEXT NOT NULL
file_size INTEGER NOT NULL
simplified_summary TEXT
processing_status TEXT DEFAULT 'pending' (pending|processing|completed|failed)
upload_date DATETIME DEFAULT CURRENT_TIMESTAMP
processed_at DATETIME
created_at DATETIME DEFAULT CURRENT_TIMESTAMP
updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
```

### Chat_Messages Table
```sql
id INTEGER PRIMARY KEY AUTOINCREMENT
report_id INTEGER NOT NULL (FK → reports.id)
user_message TEXT NOT NULL
ai_response TEXT NOT NULL
created_at DATETIME DEFAULT CURRENT_TIMESTAMP
is_deleted BOOLEAN DEFAULT FALSE
```

## Development Commands

| Command | Purpose |
|---------|---------|
| `make setup` | First-time project setup |
| `make run` | Start development server |
| `make test` | Run all tests |
| `make test-coverage` | Generate HTML coverage report |
| `make migrate-up` | Apply pending migrations |
| `make migrate-down` | Rollback last migration |
| `make migrate-create NAME=name` | Create new migration |
| `make build` | Build production binary |
| `make clean` | Clean build artifacts |

## Testing

```bash
# Run all tests
make test

# Test specific components
go test ./tests/ -v                    # Database integration tests
go test ./internal/models/ -v          # Model unit tests
go test -run TestUserModel ./tests/ -v # Single test function

# Generate coverage report
make test-coverage  # Creates coverage.html
```

## Architecture

### Repository Pattern
- **Database models** implement repository interfaces for type-safe operations
- **Separation of concerns** between data access, business logic, and HTTP handlers
- **Easy testing** with mockable repository interfaces

### Configuration
- **Environment-based** configuration with sensible defaults
- **Development**: Uses `.env` file or direct environment variables
- **Production**: Environment variables only

### Error Handling
- **Custom error types** with HTTP status codes in `pkg/errors/`
- **Consistent API responses** across all endpoints
- **Graceful error handling** with proper logging