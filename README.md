
# Gymie Backend API

Go-based REST API server for the Gymie fitness tracking application.

## Tech Stack

- **Language**: Go 1.24+
- **Framework**: Gin
- **Database**: PostgreSQL with GORM
- **Cache**: Redis
- **Authentication**: JWT
- **Storage**: Cloudflare R2 (S3-compatible)

## Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── models/                  # Data models
│   ├── handlers/                # HTTP handlers (controllers)
│   ├── middleware/              # HTTP middleware
│   ├── repository/              # Database operations
│   ├── service/                 # Business logic
│   ├── routes/                  # Route definitions
│   └── utils/                   # Utility functions
├── migrations/                  # Database migrations
├── docs/                        # Documentation
├── docker-compose.yml           # Docker services
├── Makefile                     # Build commands
└── .env.example                 # Environment variables template
```

## Quick Start

### Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose
- Make (optional, for using Makefile commands)

### 1. Clone and Setup

```bash
cd backend

# Copy environment variables
cp .env.example .env

# Edit .env with your configuration
nano .env
```

### 2. Start Database Services

```bash
# Using Make
make docker-up

# Or using docker-compose directly
docker-compose up -d
```

This starts:
- PostgreSQL on port 5432
- Redis on port 6379

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run Database Migrations

```bash
# Using Make
make migrate

# Or directly
go run migrations/migrate.go
```

### 5. Start the Server

```bash
# Using Make
make run

# Or directly
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

## Development Commands

```bash
# Start Docker services
make docker-up

# Run migrations
make migrate

# Start the server
make run

# Run tests
make test

# Run tests with coverage
make test-coverage

# Build the application
make build

# Stop Docker services
make docker-down

# Complete development setup
make dev
```

## API Endpoints

### Authentication
- `POST /v1/auth/register` - Register a new user
- `POST /v1/auth/login` - Login user
- `GET /v1/auth/me` - Get current user (protected)

### Users
- `GET /v1/users/profile` - Get user profile
- `PUT /v1/users/profile` - Update user profile
- `DELETE /v1/users/account` - Delete user account

### Workouts
- `POST /v1/workouts` - Create workout
- `GET /v1/workouts` - List workouts
- `GET /v1/workouts/:id` - Get workout by ID
- `PUT /v1/workouts/:id` - Update workout
- `DELETE /v1/workouts/:id` - Delete workout
- `GET /v1/workouts/stats` - Get workout statistics

### Nutrition
- `POST /v1/nutrition` - Create nutrition day
- `GET /v1/nutrition/:id` - Get nutrition day by ID
- `GET /v1/nutrition/date?date=YYYY-MM-DD` - Get nutrition by date
- `GET /v1/nutrition/range?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD` - Get nutrition by date range
- `PUT /v1/nutrition/:id` - Update nutrition day
- `DELETE /v1/nutrition/:id` - Delete nutrition day
- `GET /v1/nutrition/stats` - Get nutrition statistics

### Progress
- `POST /v1/progress/photos` - Create progress photo
- `GET /v1/progress/photos` - List progress photos
- `PUT /v1/progress/photos/:id` - Update progress photo
- `DELETE /v1/progress/photos/:id` - Delete progress photo
- `POST /v1/progress/weight` - Create weight entry
- `GET /v1/progress/weight` - List weight entries
- `GET /v1/progress/weight/stats` - Get weight progress
- `PUT /v1/progress/weight/:id` - Update weight entry
- `DELETE /v1/progress/weight/:id` - Delete weight entry
- `GET /v1/progress/upload-url?filename=photo.jpg` - Generate upload URL

### Health Check
- `GET /health` - Server health status

## Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

Get a token by calling `/v1/auth/login` or `/v1/auth/register`.

## Environment Variables

See [`.env.example`](.env.example) for all available configuration options.

### Required Variables

```env
# Database
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/gymie_dev

# JWT Secret (MUST change in production)
JWT_SECRET=your-secret-key-change-this-in-production

# Redis
REDIS_URL=redis://localhost:6379
```

### Optional Variables

- `PORT` - Server port (default: 8080)
- `ENVIRONMENT` - Environment mode (development/production)
- Storage configuration for file uploads
- Rate limiting settings
- CORS settings

## Database Migrations

The application uses GORM's AutoMigrate feature. To run migrations:

```bash
go run migrations/migrate.go
```

This will create all necessary tables:
- users
- user_profiles
- workouts
- exercises
- workout_sets
- nutrition_days
- meals
- foods
- progress_photos
- weight_entries

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests for specific package
go test ./internal/service/...
```

## Building for Production

```bash
# Build binary
go build -o bin/api cmd/api/main.go

# Run binary
./bin/api
```

## Docker Deployment

### Build Docker Image

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o api cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]
```

### Deploy with Docker Compose

```bash
docker-compose up -d
```

## Troubleshooting

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check PostgreSQL logs
docker logs gymie-postgres

# Restart PostgreSQL
docker restart gymie-postgres
```

### Redis Connection Issues

```bash
# Check if Redis is running
docker ps | grep redis

# Test Redis connection
redis-cli ping

# Restart Redis
docker restart gymie-redis
```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

## API Documentation

For detailed API documentation with request/response examples, see:
- [API Contract](docs/API_CONTRACT.md)
- [Backend Context](docs/BACKEND_CONTEXT.md)
- [Tech Stack Guide](docs/BACKEND_TECH_STACK.md)

## Contributing

1. Create a feature branch
2. Make your changes
3. Write tests
4. Run tests and ensure they pass
5. Submit a pull request

## License

MIT

## Support

For issues and questions, please refer to the documentation in the `docs/` directory.
