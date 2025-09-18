# Medical Report Simplifier - Backend Architecture

## Overview

The backend is built with a clean, scalable architecture following Go best practices and domain-driven design principles.

## Directory Structure Explained

### `/cmd` - Application Entry Points
- **`/cmd/server`**: Main HTTP server application
- **`/cmd/migration`**: Database migration runner (future implementation)

### `/internal` - Private Application Code
- **`/internal/auth`**: Authentication and authorization logic
- **`/internal/handlers`**: HTTP request handlers (controllers)
- **`/internal/middleware`**: HTTP middleware (auth, logging, CORS, etc.)
- **`/internal/models`**: Database models and data access layer
- **`/internal/services`**: Business logic and AI integration
- **`/internal/database`**: Database connection and query management
- **`/internal/config`**: Configuration management
- **`/internal/utils`**: Utility functions and helpers

### `/pkg` - Public Packages
- **`/pkg/types`**: Shared data structures and DTOs
- **`/pkg/errors`**: Custom error types and handling

### Other Important Directories
- **`/migrations`**: Database migration files (Goose)
- **`/uploads`**: File storage directory
- **`/tests`**: Unit and integration tests
- **`/docs`**: Project documentation

## Design Patterns Used

### 1. **Repository Pattern**
- Database operations abstracted through interfaces
- Easy to mock for testing
- Database-agnostic implementation

### 2. **Service Layer Pattern**
- Business logic separated from HTTP handlers
- Reusable across different interfaces (REST, GraphQL, CLI)

### 3. **Dependency Injection**
- Services receive dependencies through constructors
- Makes testing and configuration easier

### 4. **Error Handling**
- Custom error types with HTTP status codes
- Consistent error responses across the API

## Configuration Management

The application uses environment-based configuration with sensible defaults:

- **Development**: Uses `.env` file or environment variables
- **Production**: Uses environment variables only
- **Configuration hot-reloading**: Not implemented (add if needed)

## Security Considerations

1. **Password Hashing**: Using bcrypt for password storage
2. **JWT Tokens**: For stateless authentication
3. **File Upload Security**: Type validation and size limits
4. **SQL Injection Prevention**: Using prepared statements
5. **CORS**: Configurable cross-origin policies

## Database Design

### Core Tables
1. **users**: User authentication and profile data
2. **reports**: Uploaded medical reports and metadata
3. **chat_messages**: AI chat history per report
4. **health_metrics**: (Future) Extracted health data points

### Relationships
- Users → Reports (One-to-Many)
- Reports → Chat Messages (One-to-Many)
- Reports → Health Metrics (One-to-Many)

## API Design

### Authentication Endpoints
- `POST /api/auth/signup`: User registration
- `POST /api/auth/login`: User login
- `POST /api/auth/logout`: User logout
- `GET /api/auth/me`: Get current user info

### Report Endpoints
- `POST /api/reports/upload`: Upload medical report
- `GET /api/reports`: List user's reports
- `GET /api/reports/{id}`: Get specific report
- `GET /api/reports/{id}/summary`: Get AI-generated summary

### Chat Endpoints
- `POST /api/reports/{id}/chat`: Send message to AI about report
- `GET /api/reports/{id}/chat`: Get chat history for report

### Health Endpoints
- `GET /health`: Application health check
- `GET /metrics`: Application metrics (future)

## Development Workflow

1. **Setup**: `make setup` - Install dependencies and initialize database
2. **Development**: `make dev` - Run with hot reload (requires Air)
3. **Testing**: `make test` - Run unit and integration tests
4. **Migration**: `make migrate-create NAME=migration_name` - Create new migration
5. **Build**: `make build` - Build production binary

## Testing Strategy

### Unit Tests
- Service layer logic
- Utility functions
- Authentication logic

### Integration Tests
- HTTP endpoints
- Database operations
- File upload functionality

### E2E Tests (Future)
- Complete user workflows
- AI integration testing

## Deployment Considerations

1. **Binary Deployment**: Single binary with embedded assets
2. **Docker Support**: Containerized deployment option
3. **Database**: SQLite for development, PostgreSQL for production
4. **File Storage**: Local filesystem (can be extended to S3/GCS)
5. **Environment Variables**: All configuration via environment
6. **Health Checks**: `/health` endpoint for load balancer

## Next Steps

1. **Phase 1**: Implement authentication and file upload
2. **Phase 2**: Add AI integration for report processing
3. **Phase 3**: Build chat functionality
4. **Phase 4**: Add health metrics dashboard
5. **Phase 5**: Performance optimization and caching

## Technology Stack

- **Language**: Go 1.21+
- **HTTP Router**: Gorilla Mux
- **Database**: SQLite (dev) / PostgreSQL (prod)
- **Authentication**: JWT with bcrypt
- **File Storage**: Local filesystem
- **Testing**: Go's built-in testing + testify
- **Migration**: Goose
- **AI Integration**: OpenAI API (configurable)