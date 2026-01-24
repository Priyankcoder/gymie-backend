
.PHONY: help run migrate build test clean docker-up docker-down delete-user test-email

# Default target
help:
	@echo "Available commands:"
	@echo "  make run              - Run the application"
	@echo "  make migrate          - Run database migrations"
	@echo "  make build            - Build the application"
	@echo "  make test             - Run tests"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make docker-up        - Start Docker services (PostgreSQL, Redis)"
	@echo "  make docker-down      - Stop Docker services"
	@echo "  make dev              - Start Docker services and run migrations"
	@echo "  make delete-user      - Delete a user by email (usage: make delete-user EMAIL=user@example.com)"
	@echo "  make test-email       - Test SMTP email configuration"

# Run the application
run:
	go run cmd/api/main.go

# Run database migrations
migrate:
	go run migrations/migrate.go

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Start Docker services
docker-up:
	docker-compose up -d

# Stop Docker services
docker-down:
	docker-compose down

# Development: Start Docker and run migrations
dev: docker-up
	@echo "Waiting for database to be ready..."
	@sleep 3
	@make migrate
	@echo "Development environment is ready!"
	@echo "Run 'make run' to start the server"

# Install dependencies
deps:
	go mod download
	go mod tidy

# Delete a user by email
delete-user:
	@if [ -z "$(EMAIL)" ]; then \
		echo "❌ Error: EMAIL parameter is required"; \
		echo "Usage: make delete-user EMAIL=user@example.com"; \
		exit 1; \
	fi
	@./scripts/delete-user.sh $(EMAIL)

# Test email configuration
test-email:
	@echo "🧪 Testing SMTP email configuration..."
	@./scripts/test-email.sh
