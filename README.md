
# Gymie Backend

Go backend API for the Gymie fitness tracking mobile application.

## Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin
- **Database**: PostgreSQL 15+
- **ORM**: GORM v2
- **Cache**: Redis 7+
- **Authentication**: JWT
- **Storage**: Cloudflare R2
- **Deployment**: Fly.io

## Documentation

Complete documentation is available in the `/docs` folder:

- **[README.md](docs/README.md)** - Start here for documentation overview
- **[BACKEND_CONTEXT.md](docs/BACKEND_CONTEXT.md)** - Project context and requirements
- **[API_CONTRACT.md](docs/API_CONTRACT.md)** - Complete API specification
- **[BACKEND_TECH_STACK.md](docs/BACKEND_TECH_STACK.md)** - Technology guide and scaling

## Quick Start

### Prerequisites

- Go 1.22 or higher
- PostgreSQL 15+
- Redis 7+
- Docker and Docker Compose (for local development)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd gymie-backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your local settings
   ```

4. **Start dependencies with Docker**
   ```bash
   docker-compose up -d
   ```

5. **Run database migrations**
   ```bash
   make migrate-up
   ```

6. **Start the server**
   ```bash
   go run cmd/api/main.go
   ```

The API will be available at `http://localhost:8080`

## Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/           # HTTP handlers
│   │   ├── middleware/         # HTTP middleware
│   │   └── routes/             # Route definitions
│   ├── models/                 # Data models
│   ├── repository/             # Data access layer
│   ├── service/                # Business logic
│   └── utils/                  # Helper functions
├── pkg/                        # Public libraries
├── migrations/                 # Database migrations
├── tests/                      # Integration tests
├── docs/                       # Documentation
├── docker-compose.yml          # Local development setup
├── Dockerfile                  # Production container
└── go.mod                      # Go modules
```

## API Endpoints

See [API_CONTRACT.md](docs/API_CONTRACT.md) for complete API documentation.

### Core Endpoints

- **Authentication**: `/auth/*`
- **Workouts**: `/workouts/*`
- **Nutrition**: `/nutrition/*`
- **Progress**: `/progress/*`
- **Exercises**: `/exercises/*`
- **Templates**: `/templates/*`
- **User Profile**: `/users/*`

## Development Commands

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Build
go build -o bin/api cmd/api/main.go

# Run
./bin/api

# Format code
go fmt ./...

# Lint
golangci-lint run

# Database migrations
make migrate-up
make migrate-down
make migrate-create NAME=your_migration_name
```

## Environment Variables

Create a `.env` file based on `.env.example`:

```bash
# Server
PORT=8080

# Database
DATABASE_URL=postgresql://user:password@localhost:5432/gymie_dev

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=your-secret-key-here
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h

# Storage (Cloudflare R2)
AWS_ACCESS_KEY_ID=your-key
AWS_SECRET_ACCESS_KEY=your-secret
S3_BUCKET=gymie-dev
S3_REGION=auto
S3_ENDPOINT=https://your-account.r2.cloudflarestorage.com
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Run specific test
go test -run TestCreateWorkout ./internal/api/handlers

# Run with verbose output
go test -v ./...
```

## Deployment

### Fly.io (Recommended for MVP)

```bash
# Install flyctl
curl -L https://fly.io/install.sh | sh

# Login
flyctl auth login

# Create app
flyctl launch

# Deploy
flyctl deploy

# View logs
flyctl logs
```

See [BACKEND_TECH_STACK.md](docs/BACKEND_TECH_STACK.md#deployment-strategy) for detailed deployment instructions.

## Contributing

1. Create a feature branch
2. Make your changes
3. Write tests
4. Run `go fmt ./...` and `go test ./...`
5. Submit a pull request

## License

MIT

## Related Repositories

- **Frontend**: [gymie-frontend](../frontend) - React Native mobile app

## Support

For questions or issues, please refer to the documentation in the `/docs` folder or create an issue on GitHub.
