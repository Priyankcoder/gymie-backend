
# Backend Development Context

> Everything you need to start building the Gymie Go backend in a new repository

## Quick Start

When starting the backend in a new folder, share these documents with Stitch:

1. **[`BACKEND_TECH_STACK.md`](BACKEND_TECH_STACK.md)** - Complete tech stack recommendations
2. **[`API_CONTRACT.md`](API_CONTRACT.md)** - All API endpoints and data models
3. **This file** - Project context and requirements

---

## Project Overview

**Gymie** is a comprehensive fitness tracking mobile application (React Native + Expo) that helps users:
- Track workouts with exercises and sets
- Log daily nutrition and meals
- Monitor progress with photos and weight tracking
- Use workout templates and exercise library

**Current State:**
- ✅ Frontend: Fully functional React Native app
- ✅ Mock Data: Uses local API simulation
- 🚧 Backend: Need to build Go backend to replace mocks

---

## Business Requirements

### Core Features (MVP)

1. **User Management**
   - Email/password authentication
   - User profiles with goals (calories, macros)
   - Unit preferences (kg vs lb)

2. **Workout Tracking**
   - Create/edit/delete workouts
   - Multiple exercises per workout
   - Track sets (weight, reps, completion status)
   - Calculate workout volume and stats
   - View workout history by date range

3. **Nutrition Tracking**
   - Log meals by type (breakfast, lunch, dinner, snack)
   - Track macros (protein, carbs, fat) and calories
   - Set daily macro goals
   - View daily and historical nutrition data

4. **Progress Monitoring**
   - Upload progress photos with timestamps
   - Track body weight over time
   - View exercise-specific progress (max weight, volume)
   - Generate weekly/monthly statistics

5. **Workout Templates**
   - System-provided workout templates
   - User-created custom templates
   - Use templates to create workouts

6. **Exercise Library**
   - Searchable exercise database
   - Filter by muscle group and equipment
   - Exercise instructions and metadata

### Future Features (Phase 2)

- AI meal estimation from photos
- Social features (share workouts, leaderboards)
- Workout plan creation (multi-week programs)
- Advanced analytics and insights
- Wearable device integration
- Recipe generator

---

## Technical Context

### Frontend Architecture

**Framework:** React Native (Expo SDK 52)
**State Management:** React Context API
**Navigation:** Expo Router (file-based routing)
**Storage:** AsyncStorage
**Key Libraries:**
- `@react-navigation/native` for navigation
- `expo-image-picker` for photo uploads
- `react-native-chart-kit` for visualizations
- `date-fns` for date manipulation

**Project Structure:**
```
src/
├── components/          # Reusable UI components
│   ├── features/       # Feature-specific components
│   └── ui/             # Generic UI components
├── contexts/           # React Context providers
├── hooks/              # Custom React hooks
├── services/           # API services (localApi.ts = mock)
├── types/              # TypeScript type definitions
└── utils/              # Utility functions
```

### Current Data Flow

```
Frontend → localApi.ts (mock) → AsyncStorage → Frontend
```

**After Backend:**
```
Frontend → Real API Client → Go Backend → PostgreSQL → Frontend
```

---

## Data Models

### User Data Structure

```typescript
User {
  id: UUID
  email: string
  name?: string
  age?: number
  height?: number (cm)
  weight?: number (kg/lb)
  unit: 'kg' | 'lb'
  goals: MacroGoals {
    calories: number
    protein: number (g)
    carbs: number (g)
    fat: number (g)
  }
}
```

### Workout Data Structure

```typescript
Workout {
  id: UUID
  userId: UUID
  date: string (YYYY-MM-DD)
  exercises: Exercise[] {
    id: UUID
    name: string
    sets: WorkoutSet[] {
      weight: number
      reps: number
      completed: boolean
      restTime?: number (seconds)
    }
    muscleGroup?: string
    equipment?: string
  }
  totalVolume?: number (calculated)
  duration?: number (minutes)
}
```

### Nutrition Data Structure

```typescript
NutritionDay {
  id: UUID
  userId: UUID
  date: string (YYYY-MM-DD)
  meals: Meal[] {
    id: UUID
    name: string
    mealType: 'breakfast' | 'lunch' | 'dinner' | 'snack'
    foods: FoodItem[] {
      name: string
      calories: number
      protein: number (g)
      carbs: number (g)
      fat: number (g)
    }
  }
  totals: {
    calories: number
    protein: number
    carbs: number
    fat: number
  }
}
```

### Progress Data Structure

```typescript
ProgressPhoto {
  id: UUID
  userId: UUID
  url: string (Cloudflare R2)
  date: string
  weight?: number
}

WeightEntry {
  id: UUID
  userId: UUID
  weight: number
  date: string
}
```

---

## API Requirements

### Authentication
- JWT-based authentication
- Access token (15 min expiry)
- Refresh token (7 days)
- Token rotation on refresh

### Endpoints Needed

**Core (MVP):**
```
POST   /auth/register
POST   /auth/login
POST   /auth/refresh
POST   /auth/logout

GET    /workouts
POST   /workouts
GET    /workouts/:id
PUT    /workouts/:id
DELETE /workouts/:id
GET    /workouts/today

GET    /nutrition?date=YYYY-MM-DD
POST   /nutrition/meals
PUT    /nutrition/meals/:id
DELETE /nutrition/meals/:id

GET    /progress/stats
GET    /progress/photos
POST   /progress/photos
GET    /progress/weight
POST   /progress/weight

GET    /exercises (library)
GET    /templates
POST   /templates

GET    /users/profile
PUT    /users/profile
```

**See [`API_CONTRACT.md`](API_CONTRACT.md) for complete specifications**

---

## Database Schema

### Required Tables

1. **users**
   - id, email, password_hash, name, age, height, weight, unit
   - goals (JSONB: calories, protein, carbs, fat)
   - created_at, updated_at

2. **workouts**
   - id, user_id, date, duration, notes
   - created_at, updated_at
   - INDEX: (user_id, date DESC)

3. **exercises**
   - id, workout_id, name, muscle_group, equipment
   - sets (JSONB array)
   - created_at

4. **nutrition_days**
   - id, user_id, date
   - goals (JSONB)
   - created_at, updated_at
   - INDEX: (user_id, date DESC)

5. **meals**
   - id, nutrition_day_id, name, meal_type, time
   - created_at, updated_at

6. **food_items**
   - id, meal_id, name, serving_size, serving_unit
   - calories, protein, carbs, fat

7. **progress_photos**
   - id, user_id, url, thumbnail_url, date, weight, notes
   - created_at
   - INDEX: (user_id, date DESC)

8. **weight_entries**
   - id, user_id, weight, date
   - created_at
   - INDEX: (user_id, date DESC)

9. **templates**
   - id, user_id (nullable for system templates), name, description
   - exercises (JSONB array)
   - is_public, category, difficulty
   - created_at, updated_at

10. **exercise_library**
    - id, name, description, muscle_group, equipment
    - difficulty, instructions (JSONB array)
    - video_url, image_url

---

## Performance Requirements

### Response Times
- Authentication: <500ms
- List endpoints: <300ms
- Create/Update: <500ms
- Image upload: <2s

### Scalability Targets
- **MVP**: 10,000 users
- **Phase 2**: 100,000 users
- **Phase 3**: 1,000,000+ users

### Concurrent Users
- **MVP**: ~100 simultaneous
- **Growth**: ~1,000 simultaneous
- **Scale**: ~10,000 simultaneous

---

## Security Requirements

1. **Authentication**
   - Bcrypt password hashing (cost 12)
   - JWT with RS256 algorithm
   - Refresh token rotation
   - Session invalidation on logout

2. **Authorization**
   - Users can only access their own data
   - Admin role for system templates

3. **Rate Limiting**
   - Login: 5 attempts per 15 min
   - API: 100 req/min per user
   - Uploads: 10 files/hour per user

4. **Data Validation**
   - Input sanitization
   - SQL injection prevention (parameterized queries)
   - XSS prevention
   - CORS configuration

5. **File Upload Security**
   - File type validation
   - Size limits (10MB max)
   - Virus scanning (future)
   - Signed URLs for access

---

## Infrastructure Requirements

### MVP Setup

**Hosting:** Fly.io or Railway
**Database:** PostgreSQL 15+ (managed)
**Cache:** Redis 7+ (managed)
**Storage:** Cloudflare R2
**Monitoring:** Basic logs + Sentry

**Estimated Cost:** $20-40/month

### Services Needed

1. **Compute**
   - 1 server instance (1GB RAM, 1 CPU)
   - Auto-restart on failure
   - SSL certificate (automatic)

2. **Database**
   - PostgreSQL with 10GB storage
   - Automated backups (daily)
   - Connection pooling

3. **Cache**
   - Redis with 512MB memory
   - Persistence enabled
   - Used for sessions, rate limiting

4. **Storage**
   - Cloudflare R2 bucket
   - 10GB initial storage
   - CDN for image delivery

5. **Monitoring**
   - Error tracking (Sentry)
   - Logs aggregation
   - Performance metrics

---

## Development Workflow

### Setup Steps

1. **Initialize Go Project**
   ```bash
   mkdir gymie-backend && cd gymie-backend
   go mod init github.com/yourusername/gymie-backend
   ```

2. **Install Dependencies**
   ```bash
   go get github.com/gin-gonic/gin
   go get gorm.io/gorm
   go get github.com/redis/go-redis/v9
   go get github.com/golang-jwt/jwt/v5
   ```

3. **Set Up Database**
   ```bash
   docker-compose up -d postgres redis
   ```

4. **Run Migrations**
   ```bash
   make migrate-up
   ```

5. **Start Server**
   ```bash
   go run cmd/api/main.go
   ```

### Development Checklist

- [ ] Set up project structure
- [ ] Configure environment variables
- [ ] Create database migrations
- [ ] Implement authentication (JWT)
- [ ] Build workout endpoints
- [ ] Build nutrition endpoints
- [ ] Build progress endpoints
- [ ] Implement file uploads
- [ ] Add rate limiting
- [ ] Write tests (unit + integration)
- [ ] Set up CI/CD
- [ ] Deploy to staging
- [ ] Update frontend to use real API
- [ ] Deploy to production

---

## Testing Strategy

### Unit Tests
- Test all business logic
- Mock database calls
- Use testify for assertions
- Target: >80% coverage

### Integration Tests
- Test API endpoints end-to-end
- Use test database
- Test authentication flow
- Test data validation

### Load Tests
- Simulate 1000 concurrent users
- Test response times under load
- Identify bottlenecks

---

## Migration Plan

### Phase 1: Backend Setup (Week 1)
- Set up Go project structure
- Configure database and Redis
- Implement authentication
- Deploy to Fly.io

### Phase 2: Core APIs (Week 2-3)
- Build workout endpoints
- Build nutrition endpoints
- Build progress endpoints
- Write tests

### Phase 3: Integration (Week 4)
- Update React Native app
- Replace localApi.ts with real API
- Test all features
- Fix bugs

### Phase 4: Launch (Week 5)
- Deploy to production
- Monitor performance
- Gather user feedback
- Iterate

---

## Environment Variables

### Development
```bash
PORT=8080
DATABASE_URL=postgresql://postgres:password@localhost:5432/gymie_dev
REDIS_URL=redis://localhost:6379
JWT_SECRET=dev_secret_key
AWS_ACCESS_KEY_ID=dev_key
AWS_SECRET_ACCESS_KEY=dev_secret
S3_BUCKET=gymie-dev
```

### Production
```bash
PORT=8080
DATABASE_URL=postgresql://user:pass@host:5432/gymie_prod
REDIS_URL=redis://host:6379
JWT_SECRET=<random_secure_key>
AWS_ACCESS_KEY_ID=<production_key>
AWS_SECRET_ACCESS_KEY=<production_secret>
S3_BUCKET=gymie-prod
SENTRY_DSN=<sentry_url>
```

---

## Code Style Guide

### Go Conventions
```go
// Package naming
package workout

// Struct naming (PascalCase)
type WorkoutService struct {}

// Function naming (camelCase for private, PascalCase for public)
func (s *WorkoutService) CreateWorkout() {}
func validateInput() {}

// Error handling
if err != nil {
    return nil, fmt.Errorf("failed to create workout: %w", err)
}

// Context usage
func (s *WorkoutService) GetWorkout(ctx context.Context, id string) {}
```

### API Response Format
```go
// Success
{
    "success": true,
    "data": { ... }
}

// Error
{
    "success": false,
    "error": "Error message",
    "code": "ERROR_CODE"
}
```

---

## Common Gotchas

1. **Time Zones**
   - Always store UTC in database
   - Convert to user timezone in app
   - Use ISO 8601 format

2. **Units**
   - Store weight in user's preferred unit
   - No automatic conversion
   - Macros always in grams

3. **IDs**
   - Use UUIDs (v4)
   - Never expose sequential IDs
   - Validate UUID format in API

4. **Sets Array**
   - Store as JSONB in PostgreSQL
   - Validate array length (max 20 sets)
   - Allow partial set completion

5. **File Uploads**
   - Generate unique filenames
   - Create thumbnails for photos
   - Clean up failed uploads

---

## Frontend Integration Points

### Update Required Files

1. **`src/services/api.ts`** (create new)
   ```typescript
   import axios from 'axios';
   
   const api = axios.create({
     baseURL: process.env.API_BASE_URL,
     timeout: 30000,
   });
   
   // Add interceptors for auth
   ```

2. **`src/services/localApi.ts`** (remove/replace)
   - Currently contains all mock logic
   - Replace with real API calls

3. **`.env`** (add)
   ```
   API_BASE_URL=https://api.gymie.com/v1
   ```

4. **`src/contexts/AuthContext.tsx`** (update)
   - Add token management
   - Add refresh token logic
   - Add logout cleanup

---

## Monitoring & Observability

### Metrics to Track
- Request rate (req/sec)
- Response time (p50, p95, p99)
- Error rate (%)
- Database query time
- Cache hit rate
- Active users
- Storage usage

### Alerts to Set
- Response time > 1s
- Error rate > 5%
- Database connections > 80%
- Disk usage > 85%
- Memory usage > 90%

### Dashboards Needed
1. **API Health**: Request rate, response times, errors
2. **User Activity**: Active users, signups, logins
3. **Features**: Workouts created, meals logged, photos uploaded
4. **Infrastructure**: CPU, memory, disk, network

---

## Support Documentation

When working on the backend, refer to:

1. **[BACKEND_TECH_STACK.md](BACKEND_TECH_STACK.md)**
   - Complete tech stack details
   - Deployment options
   - Scaling strategies
   - Cost analysis

2. **[API_CONTRACT.md](API_CONTRACT.md)**
   - All API endpoints
   - Request/response formats
   - Data models
   - Error codes

3. **Frontend Source Code**
   - `src/types/index.ts` - TypeScript types
   - `src/services/localApi.ts` - Mock API implementation
   - `src/hooks/` - Business logic in hooks
   - `src/contexts/` - Global state management

---

## Quick Command Reference

```bash
# Development
go run cmd/api/main.go
go test ./...
go test -coverprofile=coverage.out ./...

# Build
go build -o bin/api cmd/api/main.go

# Database
make migrate-up
make migrate-down
make seed

# Docker
docker-compose up -d
docker-compose down
docker-compose logs -f api

# Deploy (Fly.io)
flyctl deploy
flyctl logs
flyctl ssh console
```

---

## Getting Help

### Stitch Context Prompt

When starting backend work in a new folder, share this with Stitch:

```
I'm building a Go backend for the Gymie fitness tracking app. 
Please review these context documents:

1. docs/BACKEND_TECH_STACK.md - Tech stack recommendations
2. docs/API_CONTRACT.md - Complete API specification
3. docs/BACKEND_CONTEXT.md - Project context and requirements

The frontend is a React Native (Expo) app that currently uses mock data. 
I need to build a production-ready Go backend with:
- Authentication (JWT)
- PostgreSQL database
- Redis caching
- Cloudflare R2 for images
- RESTful APIs matching the contract

Let's start with [specific task, e.g., "setting up the project structure"]
```

---

## Success Criteria

The backend is ready when:

- [ ] All API endpoints from contract are implemented
- [ ] Authentication works (register, login, refresh)
- [ ] Users can CRUD workouts, nutrition, progress
- [ ] File uploads work for progress photos
- [ ] Rate limiting is enforced
- [ ] Tests have >70% coverage
- [ ] Response times are <500ms (p95)
- [ ] Error rate is <1%
- [ ] Frontend successfully replaced mock with real API
- [ ] Production deployment is stable
- [ ] Monitoring dashboards are set up

---

This document contains everything needed to build the Gymie backend. Combined with the tech stack and API contract docs, you have complete context to start development in a new repository.

Good luck! 🚀
