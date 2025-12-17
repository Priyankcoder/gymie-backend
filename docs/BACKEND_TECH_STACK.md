
# Gymie Backend Tech Stack

> Comprehensive guide for building a scalable Go backend for the Gymie fitness tracking application

## Table of Contents

- [Core Technology Decisions](#core-technology-decisions)
- [Recommended Tech Stack](#recommended-tech-stack)
- [Database Design](#database-design)
- [Architecture Patterns](#architecture-patterns)
- [Deployment Strategy](#deployment-strategy)
- [Scaling Roadmap](#scaling-roadmap)
- [Cost Analysis](#cost-analysis)

---

## Core Technology Decisions

### Why Go?

- **Performance**: Native compilation, no runtime overhead, ~40% faster than Node.js
- **Concurrency**: Built-in goroutines perfect for handling multiple workout tracking sessions
- **Type Safety**: Strong typing reduces bugs in production
- **Deployment**: Single binary, no dependencies, easy containerization
- **Scalability**: Designed for cloud-native, microservices architecture
- **Community**: Strong ecosystem for web development

---

## Recommended Tech Stack

### 🎯 Minimal Viable Product (MVP) Stack

```yaml
Framework:         Gin (github.com/gin-gonic/gin)
Database:          PostgreSQL 15+
ORM:               GORM v2 (github.com/go-gorm/gorm)
Cache:             Redis 7+
Authentication:    JWT (github.com/golang-jwt/jwt/v5)
File Storage:      Cloudflare R2 or AWS S3
Background Jobs:   Asynq (github.com/hibiken/asynq)
Logging:           Zap (go.uber.org/zap)
API Docs:          Swagger (github.com/swaggo/swag)
Testing:           Testify (github.com/stretchr/testify)
Hosting:           Fly.io or Railway.app
```

**Estimated Monthly Cost**: $20-40
**Supports**: Up to 10,000 users

---

## Framework Options

### Option 1: Gin (Recommended)

```go
import "github.com/gin-gonic/gin"
```

**Pros:**
- Most popular Go web framework
- Excellent documentation and community
- Built-in validation and middleware
- ~40x faster than Martini
- Mature and battle-tested

**Cons:**
- Slightly heavier than Fiber
- More opinionated structure

**Use when:** You want stability, community support, and proven production use

### Option 2: Fiber

```go
import "github.com/gofiber/fiber/v2"
```

**Pros:**
- Express.js-like syntax (familiar for JavaScript devs)
- Fastest Go framework (built on fasthttp)
- Zero-allocation router
- Great middleware ecosystem

**Cons:**
- Smaller community than Gin
- Some edge cases with fasthttp

**Use when:** You need maximum performance and like Express.js patterns

### Option 3: Echo

```go
import "github.com/labstack/echo/v4"
```

**Pros:**
- Minimalist and lightweight
- Good performance
- Simple learning curve

**Cons:**
- Smaller ecosystem
- Less active maintenance

---

## Database Strategy

### Primary Database: PostgreSQL

```bash
# Version
PostgreSQL 15+

# Go Drivers
github.com/lib/pq        # Standard driver
github.com/jackc/pgx/v5  # Better performance, connection pooling

# ORM
github.com/go-gorm/gorm
github.com/go-gorm/driver/postgres
```

#### Why PostgreSQL?

1. **ACID Compliance**: Critical for workout/nutrition data integrity
2. **JSON Support**: Store flexible workout plans, exercise metadata
3. **Time-Series**: Excellent for progress tracking over time
4. **Advanced Indexing**: Fast queries on user workouts, exercises
5. **Scalability**: Proven to handle millions of users
6. **PostGIS**: Future support for gym location features

#### Schema Design Principles

```sql
-- Users table with proper constraints
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Workouts with time-series optimization
CREATE TABLE workouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_date (user_id, date DESC)
);

-- Exercises with JSONB for flexibility
CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_id UUID REFERENCES workouts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    sets JSONB NOT NULL, -- [{weight: 100, reps: 10, completed: true}]
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Alternative: MongoDB

**Consider MongoDB if:**
- You need extreme schema flexibility
- You're prototyping rapidly
- Your data is mostly document-based

**Skip MongoDB because:**
- Relational data (users → workouts → exercises)
- ACID transactions needed
- Complex queries required

---

## Caching Strategy

### Redis

```go
import "github.com/redis/go-redis/v9"
```

#### Use Cases

1. **Session Management**
   ```go
   // Store user sessions
   redis.Set(ctx, "session:"+sessionID, userData, 24*time.Hour)
   ```

2. **Rate Limiting**
   ```go
   // Limit API requests per user
   redis.Incr(ctx, "rate:"+userID)
   redis.Expire(ctx, "rate:"+userID, 60*time.Second)
   ```

3. **Caching Hot Data**
   ```go
   // Cache exercise library
   redis.Set(ctx, "exercises:all", exerciseList, 1*time.Hour)
   ```

4. **Leaderboards** (future social features)
   ```go
   // Sorted sets for rankings
   redis.ZAdd(ctx, "leaderboard:weekly", score, userID)
   ```

5. **Real-time Data**
   ```go
   // Active workout tracking
   redis.HSet(ctx, "workout:active:"+userID, "current_exercise", data)
   ```

---

## Authentication & Authorization

### JWT-Based Authentication

```go
import "github.com/golang-jwt/jwt/v5"
```

#### Implementation Strategy

```go
type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

// Access token: 15 minutes
// Refresh token: 7 days
```

#### Security Features

1. **Refresh Token Rotation**
   - Issue new refresh token on each refresh
   - Invalidate old tokens
   - Detect token theft

2. **Token Storage**
   - Access token: Memory only (React Native)
   - Refresh token: Secure storage (React Native Keychain)

3. **Social Login** (Phase 2)
   ```go
   import "golang.org/x/oauth2"
   
   // Google OAuth
   // Apple Sign In
   // Facebook Login (optional)
   ```

### Middleware Stack

```go
// Rate limiting
middleware.RateLimiter()

// CORS
middleware.CORS()

// Authentication
middleware.JWTAuth()

// Request logging
middleware.Logger()

// Error recovery
middleware.Recovery()
```

---

## API Design

### RESTful API Structure

```
Base URL: https://api.gymie.com/v1

Authentication:
POST   /auth/register           # Register new user
POST   /auth/login              # Login
POST   /auth/refresh            # Refresh token
POST   /auth/logout             # Logout
POST   /auth/forgot-password    # Password reset
POST   /auth/reset-password     # Confirm reset

Workouts:
GET    /workouts                # List user workouts
POST   /workouts                # Create workout
GET    /workouts/:id            # Get workout details
PUT    /workouts/:id            # Update workout
DELETE /workouts/:id            # Delete workout
GET    /workouts/today          # Today's workout

Exercises:
GET    /exercises               # Exercise library
GET    /exercises/:id           # Exercise details
POST   /exercises/search        # Search exercises
GET    /exercises/categories    # Exercise categories

Nutrition:
GET    /nutrition               # Daily nutrition
POST   /nutrition/meals         # Add meal
GET    /nutrition/history       # Nutrition history
POST   /nutrition/ai-estimate   # AI meal estimation

Progress:
GET    /progress/stats          # Progress statistics
GET    /progress/photos         # Progress photos
POST   /progress/photos         # Upload photo
GET    /progress/weight         # Weight history
POST   /progress/weight         # Log weight

Templates:
GET    /templates               # Workout templates
GET    /templates/:id           # Template details
POST   /templates               # Create template
```

### GraphQL Alternative (Future)

```go
import "github.com/99designs/gqlgen"
```

**Use GraphQL for:**
- Complex dashboard queries
- Flexible data requirements
- Reducing over-fetching

**Example Query:**
```graphql
query DashboardData {
  user {
    profile
    todayWorkout {
      exercises { name sets }
    }
    todayNutrition {
      meals { name calories }
    }
    weeklyProgress {
      weight
      volume
    }
  }
}
```

---

## File Storage

### Cloudflare R2 (Recommended)

```go
import (
    "github.com/aws/aws-sdk-go-v2/service/s3"
)
```

**Advantages:**
- S3-compatible API
- Zero egress fees (massive savings)
- Global CDN included
- 10GB free storage

**Use for:**
- Progress photos
- Profile pictures
- Recipe images
- Export files

**Cost**: ~$0.015/GB/month (vs S3's $0.023/GB)

### Alternative: AWS S3

**Use S3 if:**
- You're already in AWS ecosystem
- Need advanced features (versioning, lifecycle)
- Require compliance certifications

---

## Background Jobs

### Asynq (Redis-backed)

```go
import "github.com/hibiken/asynq"
```

#### Job Types

1. **Image Processing**
   ```go
   // Resize progress photos
   // Generate thumbnails
   // Optimize file sizes
   ```

2. **Notifications**
   ```go
   // Workout reminders
   // Weekly summaries
   // Achievement notifications
   ```

3. **Data Aggregation**
   ```go
   // Daily stats calculation
   // Weekly reports
   // Monthly insights
   ```

4. **AI Tasks**
   ```go
   // Meal estimation processing
   // Recipe generation
   // Workout plan creation
   ```

---

## Search & Analytics

### Exercise Search: Meilisearch

```bash
# Docker
docker run -d -p 7700:7700 getmeili/meilisearch:latest

# Go client
github.com/meilisearch/meilisearch-go
```

**Features:**
- Typo tolerance
- Instant search
- Faceted search (by muscle group, equipment)
- Ranking based on relevance

### Analytics: TimescaleDB

**Extension of PostgreSQL for time-series data**

```sql
-- Convert workouts table to hypertable
SELECT create_hypertable('workouts', 'date');

-- Efficient time-range queries
SELECT AVG(volume) 
FROM workouts 
WHERE user_id = $1 
  AND date >= now() - interval '30 days';
```

---

## Real-time Features

### WebSockets (Optional)

```go
import "github.com/gorilla/websocket"
```

**Use cases:**
- Live workout timer sync
- Real-time rest timer
- Social features (future)
- Live leaderboard updates

**Alternative:** Server-Sent Events (SSE) for simpler one-way updates

---

## Monitoring & Observability

### Logging: Zap

```go
import "go.uber.org/zap"

logger, _ := zap.NewProduction()
logger.Info("Workout created",
    zap.String("user_id", userID),
    zap.String("workout_id", workoutID),
)
```

### Metrics: Prometheus

```go
import "github.com/prometheus/client_golang/prometheus"

// Custom metrics
workoutCreated := prometheus.NewCounter(...)
apiLatency := prometheus.NewHistogram(...)
```

### Error Tracking: Sentry

```go
import "github.com/getsentry/sentry-go"

sentry.CaptureException(err)
```

### APM: Datadog or New Relic

- Request tracing
- Performance monitoring
- Database query analysis
- Real-time alerts

---

## Testing Strategy

### Unit Tests

```go
import (
    "testing"
    "github.com/stretchr/testify"
)

func TestCreateWorkout(t *testing.T) {
    // Arrange
    repo := NewMockRepository()
    service := NewWorkoutService(repo)
    
    // Act
    workout, err := service.Create(context.Background(), data)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, workout)
}
```

### Integration Tests

```go
import "github.com/ory/dockertest"

// Spin up PostgreSQL for tests
pool, _ := dockertest.NewPool("")
resource, _ := pool.Run("postgres", "15", []string{...})
```

### API Tests

```go
import "github.com/gavv/httpexpect"

e := httpexpect.New(t, server.URL)
e.POST("/workouts").
    WithJSON(data).
    Expect().
    Status(http.StatusCreated).
    JSON().Object().
    ContainsKey("id")
```

---

## Project Structure

```
gymie-backend/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
│
├── internal/                        # Private application code
│   ├── api/
│   │   ├── handlers/               # HTTP handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── workout_handler.go
│   │   │   └── nutrition_handler.go
│   │   ├── middleware/             # HTTP middleware
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   └── logger.go
│   │   └── routes/                 # Route definitions
│   │       └── routes.go
│   │
│   ├── models/                     # Data models
│   │   ├── user.go
│   │   ├── workout.go
│   │   └── nutrition.go
│   │
│   ├── repository/                 # Data access layer
│   │   ├── user_repo.go
│   │   ├── workout_repo.go
│   │   └── nutrition_repo.go
│   │
│   ├── service/                    # Business logic
│   │   ├── auth_service.go
│   │   ├── workout_service.go
│   │   └── nutrition_service.go
│   │
│   └── utils/                      # Helper functions
│       ├── jwt.go
│       ├── validation.go
│       └── errors.go
│
├── pkg/                            # Public libraries (reusable)
│   ├── auth/
│   │   └── jwt.go
│   ├── cache/
│   │   └── redis.go
│   ├── storage/
│   │   └── s3.go
│   └── logger/
│       └── logger.go
│
├── migrations/                     # Database migrations
│   ├── 001_create_users.up.sql
│   ├── 001_create_users.down.sql
│   ├── 002_create_workouts.up.sql
│   └── 002_create_workouts.down.sql
│
├── tests/                          # Integration tests
│   ├── api/
│   └── integration/
│
├── scripts/                        # Utility scripts
│   ├── seed.go                    # Seed database
│   └── migrate.go                 # Run migrations
│
├── configs/                        # Configuration files
│   ├── config.yaml
│   └── config.prod.yaml
│
├── docs/                          # Documentation
│   ├── api/                       # API documentation
│   └── architecture/              # Architecture docs
│
├── .env.example                   # Environment variables template
├── .gitignore
├── docker-compose.yml             # Local development setup
├── Dockerfile                     # Production container
├── go.mod                         # Go modules
├── go.sum
├── Makefile                       # Common commands
└── README.md
```

---

## Deployment Strategy

### Development Environment

```yaml
# docker-compose.yml
version: '3.8'
services:
  postgres:
    image: postgres:15-alpine
    ports: ["5432:5432"]
    environment:
      POSTGRES_DB: gymie_dev
      POSTGRES_PASSWORD: dev_password
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports: ["6379:6379"]
    volumes:
      - redis_data:/data

  api:
    build: .
    ports: ["8080:8080"]
    depends_on: [postgres, redis]
    environment:
      DATABASE_URL: postgres://postgres:dev_password@postgres/gymie_dev
      REDIS_URL: redis://redis:6379
```

### Production Deployment

#### Option 1: Fly.io (Recommended for MVP)

```toml
# fly.toml
app = "gymie-api"

[env]
  PORT = "8080"

[[services]]
  http_checks = []
  internal_port = 8080
  protocol = "tcp"
  
  [services.concurrency]
    hard_limit = 250
    soft_limit = 200

[[services.ports]]
  handlers = ["http"]
  port = 80

[[services.ports]]
  handlers = ["tls", "http"]
  port = 443
```

**Cost**: ~$5/month to start

#### Option 2: Railway.app

```json
{
  "build": {
    "builder": "NIXPACKS",
    "buildCommand": "go build -o main cmd/api/main.go"
  },
  "deploy": {
    "startCommand": "./main",
    "restartPolicyType": "ON_FAILURE"
  }
}
```

**Cost**: ~$5/month to start

#### Option 3: AWS (Production-grade)

```yaml
# ECS Task Definition
{
  "family": "gymie-api",
  "containerDefinitions": [{
    "name": "api",
    "image": "gymie/api:latest",
    "memory": 512,
    "cpu": 256,
    "essential": true,
    "portMappings": [{
      "containerPort": 8080
    }]
  }]
}
```

**Cost**: ~$30-50/month

### CI/CD Pipeline

```yaml
# .github/workflows/deploy.yml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      - run: go test ./...
      - run: go build -v ./...

  deploy:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: superfly/flyctl-actions/setup-flyctl@master
      - run: flyctl deploy --remote-only
```

---

## Scaling Roadmap

### Phase 1: MVP (0-10K users)

**Infrastructure:**
- Single server instance
- Managed PostgreSQL (1 core, 1GB RAM)
- Redis (512MB)
- Cloudflare R2 storage

**Cost**: $20-40/month

**Handles:**
- ~100 concurrent users
- ~1000 API requests/minute
- Basic monitoring

### Phase 2: Growth (10K-100K users)

**Upgrades:**
- 2-3 server instances with load balancer
- PostgreSQL read replicas
- Redis cluster
- CDN for static assets

**New Services:**
- APM monitoring (Datadog/New Relic)
- Error tracking (Sentry)
- Rate limiting per user

**Cost**: $200-500/month

**Handles:**
- ~1000 concurrent users
- ~10,000 API requests/minute
- Advanced analytics

### Phase 3: Scale (100K-1M users)

**Architecture Changes:**
- Microservices split:
  - Auth service
  - Workout service
  - Nutrition service
  - Progress service
- Message queue (RabbitMQ/Kafka)
- Database sharding by user_id
- Multi-region deployment

**Cost**: $2,000-5,000/month

**Handles:**
- ~10,000 concurrent users
- ~100,000 API requests/minute
- Global latency <100ms

### Phase 4: Enterprise (1M+ users)

**Full Scale:**
- Kubernetes orchestration
- Service mesh (Istio)
- Multi-region database clusters
- Advanced caching strategies
- Dedicated DevOps team

**Cost**: $10,000+/month

---

## Cost Analysis

### Development Costs

| Service | Provider | Cost |
|---------|----------|------|
| Hosting | Fly.io Free Tier | $0 |
| Database | Supabase Free Tier | $0 |
| Redis | Upstash Free Tier | $0 |
| Storage | Cloudflare R2 Free Tier | $0 |
| **Total** | | **$0/month** |

### MVP Production (10K users)

| Service | Provider | Cost |
|---------|----------|------|
| Hosting | Fly.io (2GB RAM) | $5-10 |
| Database | Supabase Pro | $25 |
| Redis | Upstash Pay-as-you-go | $5 |
| Storage | Cloudflare R2 | $5 |
| Monitoring | Free tiers | $0 |
| **Total** | | **$40-45/month** |

### Growth Production (100K users)

| Service | Provider | Cost |
|---------|----------|------|
| Hosting | Fly.io (4GB RAM, 2 instances) | $30 |
| Database | Supabase Pro + replicas | $75 |
| Redis | Upstash Pro | $20 |
| Storage | Cloudflare R2 (100GB) | $10 |
| Monitoring | Datadog | $15 |
| CDN | Cloudflare | $20 |
| **Total** | | **$170/month** |

---

## Security Considerations

### Essential Security Measures

1. **Environment Variables**
   ```bash
   # Never commit secrets
   DATABASE_URL=postgresql://...
   JWT_SECRET=random_secure_string
   AWS_ACCESS_KEY=...
   ```

2. **HTTPS Only**
   - Force TLS 1.3
   - HSTS headers
   - Certificate pinning (mobile app)

3. **Rate Limiting**
   ```go
   // Per user rate limits
   LOGIN: 5 attempts / 15 minutes
   API: 100 requests / minute
   UPLOADS: 10 files / hour
   ```

4. **Input Validation**
   ```go
   // Use struct tags with validator
   type CreateWorkout struct {
       Date string `json:"date" binding:"required,date"`
       Exercises []Exercise `json:"exercises" binding:"required,dive"`
   }
   ```

5. **SQL Injection Prevention**
   ```go
   // Use parameterized queries with GORM
   db.Where("user_id = ?", userID).Find(&workouts)
   ```

6. **Password Security**
   ```go
   import "golang.org/x/crypto/bcrypt"
   
   // Hash with bcrypt cost 12
   hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
   ```

---

## Performance Optimization

### Database

1. **Indexing Strategy**
   ```sql
   -- Composite indexes for common queries
   CREATE INDEX idx_workouts_user_date ON workouts(user_id, date DESC);
   CREATE INDEX idx_exercises_workout ON exercises(workout_id);
   CREATE INDEX idx_nutrition_user_date ON nutrition(user_id, date DESC);
   ```

2. **Connection Pooling**
   ```go
   db.SetMaxOpenConns(25)
   db.SetMaxIdleConns(5)
   db.SetConnMaxLifetime(5 * time.Minute)
   ```

3. **Query Optimization**
   ```go
   // Preload related data
   db.Preload("Exercises").Find(&workouts)
   
   // Select specific fields
   db.Select("id, name, date").Find(&workouts)
   ```

### API

1. **Response Compression**
   ```go
   router.Use(gzip.Gzip(gzip.DefaultCompression))
   ```

2. **Response Caching**
   ```go
   // Cache exercise library
   if cached, err := redis.Get(ctx, "exercises"); err == nil {
       return cached
   }
   ```

3. **Pagination**
   ```go
   // Cursor-based pagination
   db.Where("id > ?", lastID).Limit(20).Find(&items)
   ```

---

## Monitoring Checklist

### Metrics to Track

- [ ] Request rate (requests/second)
- [ ] Response time (p50, p95, p99)
- [ ] Error rate (4xx, 5xx)
- [ ] Database query time
- [ ] Cache hit rate
- [ ] Active users
- [ ] Queue depth (background jobs)

### Alerts to Configure

- [ ] API response time > 500ms
- [ ] Error rate > 1%
- [ ] Database connections > 80%
- [ ] Disk usage > 85%
- [ ] Memory usage > 90%
- [ ] Failed background jobs

---

## Migration Strategy

### From Mock Data to Real Backend

1. **Phase 1: API Contract** (Week 1)
   - Define OpenAPI spec
   - Update frontend to use environment variables
   - Create API client with Axios/Fetch

2. **Phase 2: Core APIs** (Week 2-3)
   - Auth endpoints
   - Workout CRUD
   - Exercise library

3. **Phase 3: Advanced Features** (Week 4-5)
   - Nutrition tracking
   - Progress photos
   - Background jobs

4. **Phase 4: Testing & Deployment** (Week 6)
   - Integration tests
   - Load testing
   - Production deployment

---

## Resources & Learning

### Go Resources
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Web Examples](https://gowebexamples.com/)

### Architecture
- [Twelve-Factor App](https://12factor.net/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

### Performance
- [High Performance Go](https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html)

---

## Next Steps

1. **Set up development environment**
   ```bash
   # Initialize Go module
   go mod init github.com/yourusername/gymie-backend
   
   # Install dependencies
   go get github.com/gin-gonic/gin
   go get gorm.io/gorm
   go get github.com/redis/go-redis/v9
   ```

2. **Create initial project structure**
   ```bash
   mkdir -p cmd/api internal/{api,models,repository,service} pkg
   ```

3. **Set up local database**
   ```bash
   docker-compose up -d postgres redis
   ```

4. **Start building!**

---

## Questions?

For implementation help, refer to:
- **API Design**: See REST endpoint specifications above
- **Database Schema**: See PostgreSQL schema examples
- **Deployment**: Choose Fly.io for fastest start
- **Scaling**: Follow the roadmap as you grow

Good luck building Gymie backend! 🚀
