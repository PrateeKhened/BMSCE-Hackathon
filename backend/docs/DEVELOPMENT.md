# Development Guide

## Quick Start

1. **Clone and Setup**:
   ```bash
   cd backend
   make setup
   ```

2. **Start Development Server**:
   ```bash
   make run
   # Or with hot reload:
   make dev
   ```

3. **Test the Setup**:
   ```bash
   curl http://localhost:8080/health
   ```

## Development Commands

| Command | Purpose |
|---------|---------|
| `make deps` | Download Go dependencies |
| `make run` | Start the server |
| `make dev` | Start with hot reload (install `air` first) |
| `make build` | Build production binary |
| `make test` | Run all tests |
| `make test-coverage` | Run tests with coverage report |
| `make fmt` | Format code |
| `make lint` | Run linter |
| `make clean` | Clean build artifacts |

## Database Commands

| Command | Purpose |
|---------|---------|
| `make migrate-up` | Apply pending migrations |
| `make migrate-down` | Rollback last migration |
| `make migrate-status` | Check migration status |
| `make migrate-create NAME=migration_name` | Create new migration |
| `make init-db` | Initialize database for first time |

## Project Structure Rules

### File Naming Conventions
- **Files**: `snake_case.go` (e.g., `user_service.go`)
- **Packages**: `lowercase` (e.g., `package auth`)
- **Constants**: `PascalCase` (e.g., `MaxFileSize`)
- **Variables/Functions**: `camelCase` (e.g., `userService`)

### Package Organization Rules
1. **`internal/`**: Private code, not importable by other projects
2. **`pkg/`**: Public packages that can be imported
3. **`cmd/`**: Application entry points
4. **One concept per package**: Each package should have a single responsibility

### Code Style Guidelines
1. **Error Handling**: Always handle errors explicitly
2. **Context**: Pass `context.Context` for cancellation and timeouts
3. **Interfaces**: Define interfaces where they're used, not where they're implemented
4. **Testing**: Test files end with `_test.go`

## Testing Guidelines

### Unit Test Structure
```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   types.SignupRequest
        want    *types.User
        wantErr bool
    }{
        // test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Integration Test Setup
```go
func TestAPI(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    // Create test server
    server := setupTestServer(t, db)

    // Run tests
}
```

## Environment Setup

### Required Environment Variables
```bash
# Copy example environment file
cp .env.example .env

# Edit with your settings
vim .env
```

### Development Dependencies
```bash
# Hot reload tool
go install github.com/cosmtrek/air@latest

# Linter
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Migration tool (if not using make commands)
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## API Testing

### Using curl
```bash
# Health check
curl http://localhost:8080/health

# Sign up
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","full_name":"Test User"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### Using HTTP files (VS Code REST Client)
Create `test.http` file:
```http
### Health Check
GET http://localhost:8080/health

### Sign Up
POST http://localhost:8080/api/auth/signup
Content-Type: application/json

{
  "email": "test@example.com",
  "password": "password123",
  "full_name": "Test User"
}
```

## Database Schema Changes

### Creating Migrations
```bash
# Create a new migration
make migrate-create NAME=add_user_profile_fields

# This creates:
# migrations/20231218120000_add_user_profile_fields.sql
```

### Migration File Structure
```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS users;
```

## Debugging

### Logging
- Use structured logging (add logger in future)
- Log at appropriate levels (DEBUG, INFO, WARN, ERROR)
- Include context in log messages

### Common Issues
1. **Port already in use**: Change PORT in .env or kill existing process
2. **Database locked**: Close any DB browser connections
3. **File permissions**: Ensure uploads/ directory is writable
4. **Missing dependencies**: Run `make deps`

## Code Review Checklist

- [ ] All errors are handled
- [ ] Tests are written for new functionality
- [ ] Code follows Go conventions
- [ ] No sensitive data in logs
- [ ] Database queries use prepared statements
- [ ] Input validation is present
- [ ] Documentation is updated

## Performance Considerations

1. **Database**: Use appropriate indexes
2. **File Upload**: Stream large files, don't load into memory
3. **Caching**: Add Redis for session storage (future)
4. **Connection Pooling**: Configure database connection pool
5. **Graceful Shutdown**: Handle SIGTERM properly