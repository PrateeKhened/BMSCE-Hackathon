# Medical Report Simplifier - Backend

## Project Structure

```
backend/
├── cmd/                    # Application entry points
│   ├── server/            # Main server application
│   └── migration/         # Migration runner
├── internal/              # Private application code
│   ├── auth/             # Authentication logic
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   ├── models/           # Database models
│   ├── services/         # Business logic
│   ├── database/         # Database connection and queries
│   ├── config/           # Configuration management
│   └── utils/            # Utility functions
├── pkg/                   # Public packages
│   ├── types/            # Shared types and structs
│   └── errors/           # Custom error types
├── migrations/           # Database migrations (Goose)
├── uploads/              # File upload storage
├── tests/                # Test files
│   ├── unit/            # Unit tests
│   └── integration/     # Integration tests
└── docs/                 # Documentation
```

## Features

### Phase 1 - Core Backend
- [x] Project structure setup
- [ ] Database setup with SQLite and Goose migrations
- [ ] User authentication (signup/login/logout)
- [ ] JWT token management
- [ ] File upload endpoint

### Phase 2 - AI Integration
- [ ] AI service integration for report processing
- [ ] Report parsing and simplification
- [ ] Chat functionality with uploaded reports

### Phase 3 - Dashboard & Analytics
- [ ] Health metrics tracking
- [ ] Report history management
- [ ] User dashboard data endpoints

## Technology Stack

- **Language**: Go 1.21+
- **Database**: SQLite
- **Migration**: Goose
- **Authentication**: JWT
- **File Storage**: Local filesystem
- **AI Integration**: TBD (OpenAI/Claude/Local LLM)

## Database Schema

### Users Table
- id (PRIMARY KEY)
- email (UNIQUE)
- password_hash
- full_name
- created_at
- updated_at

### Reports Table
- id (PRIMARY KEY)
- user_id (FOREIGN KEY)
- original_filename
- file_path
- file_type
- simplified_summary
- upload_date
- processed_at

### Chat_Messages Table
- id (PRIMARY KEY)
- report_id (FOREIGN KEY)
- user_message
- ai_response
- created_at

### Health_Metrics Table (Future)
- id (PRIMARY KEY)
- user_id (FOREIGN KEY)
- report_id (FOREIGN KEY)
- metric_type (blood_pressure, diabetes, etc.)
- value
- unit
- date_recorded