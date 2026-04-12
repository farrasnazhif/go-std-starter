# Go User Management API Starter

A clean, structured Go starter project with user authentication, user management, and email invitations. This is a minimal starter focused on user workflows without social features.

## Project Structure

```
├── cmd/
│   ├── api/              # API server
│   │   ├── main.go       # Entry point
│   │   ├── api.go        # Routes and HTTP configuration
│   │   ├── auth.go       # Authentication handlers
│   │   ├── users.go      # User handlers
│   │   ├── health.go     # Health check handler
│   │   ├── errors.go     # Error handling utilities
│   │   └── json.go       # JSON utilities
│   └── migrate/
│       ├── migrations/   # Database migration files (user & invitation tables)
│       └── seed/         # Database seeding scripts
├── internal/
│   ├── db/               # Database connection and setup
│   ├── store/            # Data storage layer (users only)
│   ├── env/              # Environment variable utilities
│   └── mailer/           # Email service integration (SendGrid)
├── docker-compose.yml    # PostgreSQL database setup
├── go.mod & go.sum       # Go dependencies
├── Makefile              # Build and utility commands
└── .air.toml             # Hot reload configuration (for development)
```

## Features

- **User Registration**: Create and register users with email and password
- **User Authentication**: User login and account management
- **User Profiles**: Retrieve user information
- **Email Invitations**: Send invitations via SendGrid with token-based activation
- **Account Activation**: Activate accounts via invitation tokens
- **Database Migrations**: SQL-based schema management
- **API Documentation**: Swagger/OpenAPI integration
- **Clean Architecture**: Separation of concerns with repository pattern

## Tech Stack

- **Go 1.25.4**
- **Chi Router**: Fast, composable HTTP router
- **PostgreSQL**: Primary database
- **SendGrid**: Email service
- **Zap**: Structured logging
- **Swagger**: API documentation

## Getting Started

### Prerequisites

- Go 1.25.4+
- PostgreSQL 16+
- Docker & Docker Compose (for quick database setup)
- migrate CLI tool (for running migrations)

### Installation

1. **Clone or initialize the project**

   ```bash
   git clone <repository-url>
   cd go-std-starter
   ```

2. **Allow direnv** (if using `.envrc` file)

   ```bash
   direnv allow .
   ```

3. **Start PostgreSQL with Docker Compose**

   ```bash
   docker compose --build up
   ```

4. **Install migrate tool** (if not installed)

   ```bash
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ```

5. **Run database migrations**

   ```bash
   make migrate-up
   ```

6. **Install dependencies**

   ```bash
   go mod tidy
   go mod download
   ```

7. **Start the server**

   **Option A: Using air (hot reload for development)**
   
   ```bash
   air
   ```

   **Option B: Build and run manually**
   
   ```bash
   go build -o bin/api ./cmd/api/
   ./bin/api
   ```

The API will be available at `http://localhost:8080`

## Environment Variables

You can set environment variables in two ways:

### Option 1: Using `.env` file
Create a `.env` file in the project root:

```bash
# Server
ADDR=:8080
ENV=development

# Database
DB_ADDR=postgres://admin:adminpassword@localhost/go-std-starter?sslmode=disable
DB_MAX_OPEN_CONNS=30
DB_MAX_IDLE_CONNS=30
DB_MAX_IDLE_TIME=15m

# External URLs
EXTERNAL_URL=localhost:8080
FRONTEND_URL=http://localhost:3000

# Email (SendGrid)
SENDGRID_API_KEY=your_sendgrid_api_key_here
```

### Option 2: Using `.envrc` (with direnv)
If you're using direnv, create a `.envrc` file with your environment variables and run:

```bash
direnv allow .
```

**After modifying `.envrc`**, always run:

```bash
direnv allow .
```

This will reload the environment variables.

## API Endpoints

### Health Check

- `GET /api/v1/health` - Server health status

### Authentication

- `POST /api/v1/auth/user` - Register a new user

### Users

- `GET /api/v1/users/{userID}` - Get user profile
- `PUT /api/v1/users/activate/{token}` - Activate user account

### API Documentation

- `GET /api/v1/swagger/*` - Swagger UI documentation

## Make Commands

```bash
# Run tests
make test

# Create a new migration
make migration NAME=create_table_name

# Run migrations up
make migrate-up

# Run migrations down
make migrate-down

# Generate Swagger documentation
make gen-docs

# Seed database with sample data
make seed
```

## Database Schema

The starter includes migrations for:

- **users table** - User accounts with email, password, and activation status
- **user_invitations table** - Email invitation tokens with expiry for account activation
- **is_active column** - Track account activation status

## Development

### Hot Reload

For development with hot reload, install air:

```bash
go install github.com/cosmtrek/air@latest
```

Then run:

```bash
air
```

### Running Tests

```bash
make test
```

## Adding New Features

### Adding a New Handler

1. Create a handler function in `cmd/api/` (e.g., `cmd/api/new_feature.go`)
2. Register routes in `cmd/api/api.go` within the appropriate router group
3. Add necessary store methods in `internal/store/`

### Adding Database Migrations

```bash
make migration NAME=describe_your_migration
# Edit the generated .up.sql and .down.sql files
make migrate-up
```

### Sending Emails

The mailer is already integrated via SendGrid. Use `app.mailer` in your handlers:

```go
app.mailer.Send(to, subject, htmlBody)
```

## Project Architecture

### Repository Pattern

The `internal/store/` package implements the repository pattern, providing a clean abstraction over database operations.

### Dependency Injection

The `application` struct in `cmd/api/api.go` serves as the service locator, holding references to all dependencies (config, store, logger, mailer).

### Error Handling

Centralized error handling in `cmd/api/errors.go` provides consistent API error responses.

## Notes

- This starter **does not** include posts, feeds, comments, or follower functionality
- It's focused on user management and authentication workflows
- Extend this starter by adding your own features in the appropriate layers
- Keep business logic in the store layer, HTTP handling in the api layer
- Use the logger instance (`app.logger`) for all logging

## License

Unlicense

## Contributing

Feel free to extend and customize this starter for your needs.

---

**Generated from**: go-std-starter repository structure  
**Initial Version**: 0.0.1
