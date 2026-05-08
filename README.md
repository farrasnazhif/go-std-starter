# Go User Management API Starter

A minimal Go starter for user authentication, profile management, and email verification. Built with clean architecture principles and dependency injection.

## Tech Stack

- **Go 1.25.4** | **Chi Router** | **PostgreSQL** | **Resend** (email) | **Zap** (logging) | **Swagger** (docs)

## Features

- User registration with email verification
- Account activation via token
- User profiles
- Clean architecture with repository pattern
- Swagger/OpenAPI documentation

## Getting Started

### Prerequisites

- Go 1.25.4+, PostgreSQL 16+, Docker & Docker Compose
- `migrate` CLI: `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

### Quick Start

```bash
# Clone and setup
git clone <repo> && cd go-std-starter

# Start PostgreSQL
docker compose up

# Run migrations
make migrate-up

# Install dependencies and start server
go mod tidy
air  # or: go build -o bin/api ./cmd/api/ && ./bin/api
```

API runs on `http://localhost:8080`

## Environment Variables

Create `.env` file in project root:

```env
ADDR=:8080
ENV=development

# Database
DB_ADDR=postgres://admin:adminpassword@localhost/go-std-starter?sslmode=disable
DB_MAX_OPEN_CONNS=30
DB_MAX_IDLE_CONNS=30
DB_MAX_IDLE_TIME=15m

# URLs
EXTERNAL_URL=localhost:8080
FRONTEND_URL=http://localhost:3000

# Email
MAIL_FROM_EMAIL=onboarding@resend.dev
RESEND_API_KEY=your_key_here
```

Or use `.envrc` with direnv: `direnv allow .`

## API Endpoints

| Method | Endpoint                         | Description      |
| ------ | -------------------------------- | ---------------- |
| GET    | `/api/v1/health`                 | Health check     |
| POST   | `/api/v1/auth/user`              | Register user    |
| GET    | `/api/v1/users/{id}`             | Get user profile |
| PUT    | `/api/v1/users/activate/{token}` | Activate account |
| GET    | `/api/v1/swagger/*`              | API docs         |

## Make Commands

```bash
make test              # Run tests
make migration NAME=x  # Create migration
make migrate-up        # Run migrations
make migrate-down      # Rollback migrations
make gen-docs          # Generate Swagger docs
make seed              # Seed database
```

## Project Structure

```
├── cmd/api/        # HTTP handlers & routes
├── internal/
│   ├── db/         # Database connection
│   ├── store/      # Repository layer
│   ├── mailer/     # Email service (Resend)
│   └── env/        # Config management
├── docs/           # Swagger docs
└── docker-compose.yml
```

## Development

### Hot Reload

```bash
go install github.com/cosmtrek/air@latest
air
```

### Adding Features

1. Create handler in `cmd/api/new_feature.go`
2. Register routes in `cmd/api/api.go`
3. Add store methods in `internal/store/`

### Database Migrations

```bash
make migration NAME=describe_change
# Edit the .up.sql and .down.sql files
make migrate-up
```

### Sending Emails

```go
app.mailer.Send(templateFile, username, email, data, isSandbox)
```

## Notes

- Focused on user management & auth workflows
- No posts, feeds, comments, or social features
- Extend with your own features maintaining the patterns
- Keep business logic in store layer, HTTP in api layer

## License

Unlicense
