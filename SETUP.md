
# Backend Setup Guide

Complete guide to set up and run the Gymie backend API.

## Prerequisites

Ensure you have the following installed:
- **Go 1.24+**: [Download here](https://golang.org/dl/)
- **Docker & Docker Compose**: [Download here](https://www.docker.com/products/docker-desktop)
- **Make** (optional): For using Makefile commands

## Step-by-Step Setup

### 1. Environment Configuration

Create your environment file:

```bash
cd backend
cp .env.example .env
```

Edit `.env` and update the following critical settings:

```env
# IMPORTANT: Change JWT_SECRET in production!
JWT_SECRET=your-very-secure-secret-key-change-this

# Database (default values work with docker-compose)
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/gymie_dev

# Redis (default values work with docker-compose)
REDIS_URL=redis://localhost:6379
```

### 2. Start Database Services

Start PostgreSQL and Redis using Docker:

```bash
# Start services in background
make docker-up

# Or using docker-compose directly
docker-compose up -d

# Verify services are running
docker ps
```

You should see:
- `gymie-postgres` on port 5432
- `gymie-redis` on port 6379

### 3. Install Go Dependencies

```bash
go mod download
```

### 4. Run Database Migrations

This creates all necessary database tables:

```bash
make migrate

# Or directly
go run migrations/migrate.go
```

Expected output:
```
Running database migrations...
Migrations completed successfully!
```

### 5. Start the API Server

```bash
make run

# Or directly
go run cmd/api/main.go
```

Expected output:
```
Database connection established
Redis connection established
Starting server on port 8080
```

### 6. Test the API

Open a new terminal and test the health endpoint:

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "success": true,
  "message": "Server is running",
  "data": {
    "status": "healthy",
    "version": "1.0.0"
  }
}
```

## Quick Start (All-in-One)

For a complete setup in one command:

```bash
make dev
```

This will:
1. Start Docker services (PostgreSQL & Redis)
2. Wait for services to be ready
3. Run database migrations
4. Display instructions to start the server

Then just run:
```bash
make run
```

## Testing the API

### Register a User

```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

Save the `token` from the response.

### Get User Profile

```bash
curl http://localhost:8080/v1/users/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Create a Workout

```bash
curl -X POST http://localhost:8080/v1/workouts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "name": "Morning Workout",
    "date": "2024-12-18T08:00:00Z",
    "duration": 60,
    "exercises": [
      {
        "name": "Bench Press",
        "muscle_group": "chest",
        "order": 1,
        "sets": [
          {
            "set_number": 1,
            "weight": 80,
            "reps": 10,
            "completed": true
          }
        ]
      }
    ]
  }'
```

## Troubleshooting

### Port Already in Use

If port 8080 is already in use:

1. Find the process:
```bash
lsof -i :8080
```

2. Kill it:
```bash
kill -9 <PID>
```

3. Or change the port in `.env`:
```env
PORT=8081
```

### Database Connection Failed

1. Check if PostgreSQL is running:
```bash
docker ps | grep postgres
```

2. Check logs:
```bash
docker logs gymie-postgres
```

3. Restart PostgreSQL:
```bash
docker restart gymie-postgres
```

4. Verify connection:
```bash
docker exec -it gymie-postgres psql -U postgres -d gymie_dev
```

### Redis Connection Failed

1. Check if Redis is running:
```bash
docker ps | grep redis
```

2. Test connection:
```bash
docker exec -it gymie-redis redis-cli ping
```
Should return: `PONG`

3. Restart Redis:
```bash
docker restart gymie-redis
```

### Migration Errors

If migrations fail:

1. Check database connection
2. Drop and recreate database:
```bash
docker exec -it gymie-postgres psql -U postgres -c "DROP DATABASE gymie_dev;"
docker exec -it gymie-postgres psql -U postgres -c "CREATE DATABASE gymie_dev;"
```

3. Run migrations again:
```bash
make migrate
```

### Build Errors

1. Clean and rebuild:
```bash
make clean
go mod tidy
make build
```

2. Verify Go version:
```bash
go version
```
Should be 1.24 or higher.

## Development Workflow

### Daily Development

1. Start Docker services (if not running):
```bash
make docker-up
```

2. Start the server:
```bash
make run
```

3. Make changes to code

4. Server will need manual restart (no hot reload by default)

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Building for Production

```bash
# Build binary
make build

# Binary will be in bin/api
./bin/api
```

### Stopping Services

```bash
# Stop API server: Ctrl+C

# Stop Docker services
make docker-down
```

## Database Management

### Access PostgreSQL

```bash
docker exec -it gymie-postgres psql -U postgres -d gymie_dev
```

Common commands:
- `\dt` - List all tables
- `\d table_name` - Describe table structure
- `SELECT * FROM users;` - Query users
- `\q` - Quit

### Access Redis

```bash
docker exec -it gymie-redis redis-cli
```

Common commands:
- `KEYS *` - List all keys
- `GET key_name` - Get value
- `FLUSHALL` - Clear all data (caution!)
- `exit` - Quit

### Reset Database

```bash
# Stop services
make docker-down

# Remove volumes (deletes all data!)
docker-compose down -v

# Start fresh
make docker-up
make migrate
```

## Environment Variables Reference

### Required

```env
DATABASE_URL=postgresql://user:pass@host:port/dbname
JWT_SECRET=your-secret-key
REDIS_URL=redis://localhost:6379
```

### Optional

```env
PORT=8080                          # Server port
ENVIRONMENT=development            # development/production
JWT_EXPIRATION=24                  # Token expiration in hours

# Storage (for file uploads)
STORAGE_ENDPOINT=https://...
STORAGE_ACCESS_KEY=...
STORAGE_SECRET_KEY=...
STORAGE_BUCKET=gymie-dev

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60

# CORS
CORS_ALLOWED_ORIGINS=*
```

## Next Steps

1. ✅ Backend is running
2. 📱 Test with frontend app (see `../frontend/README.md`)
3. 📚 Read API documentation (`docs/API_CONTRACT.md`)
4. 🚀 Deploy to production (see `docs/BACKEND_TECH_STACK.md`)

## Support

- **API Documentation**: `docs/API_CONTRACT.md`
- **Architecture Guide**: `docs/BACKEND_CONTEXT.md`
- **Deployment Guide**: `docs/BACKEND_TECH_STACK.md`
- **Frontend Integration**: `../frontend/docs/README.md`

## Common Make Commands

```bash
make help          # Show all available commands
make dev           # Complete development setup
make run           # Start the server
make migrate       # Run database migrations
make test          # Run tests
make build         # Build binary
make docker-up     # Start Docker services
make docker-down   # Stop Docker services
make clean         # Clean build artifacts
```

Happy coding! 🚀
