
# Gymie API Contract

> Complete API specification based on the React Native frontend implementation

## Overview

This document defines the API contract between the Gymie React Native app and the Go backend. All endpoints, request/response formats, and data models are derived from the current frontend implementation.

---

## Base Configuration

```
Base URL: https://api.gymie.com/v1
Authentication: JWT Bearer Token
Content-Type: application/json
API Version: 1.0.0
```

---

## Data Models

### User

```typescript
interface User {
  id: string;              // UUID
  email: string;
  name?: string;
  age?: number;
  height?: number;         // in cm
  weight?: number;         // in kg
  unit: 'kg' | 'lb';
  createdAt: string;       // ISO 8601
  updatedAt: string;       // ISO 8601
}
```

### Workout

```typescript
interface Workout {
  id: string;              // UUID
  userId: string;          // UUID
  date: string;            // YYYY-MM-DD
  exercises: Exercise[];
  duration?: number;       // minutes
  totalVolume?: number;    // kg (calculated)
  notes?: string;
  createdAt: string;       // ISO 8601
  updatedAt: string;       // ISO 8601
}
```

### Exercise

```typescript
interface Exercise {
  id: string;              // UUID
  workoutId: string;       // UUID
  name: string;
  sets: WorkoutSet[];
  notes?: string;
  muscleGroup?: string;    // chest, back, legs, shoulders, arms, core
  equipment?: string;      // barbell, dumbbell, machine, bodyweight, cable
}
```

### WorkoutSet

```typescript
interface WorkoutSet {
  weight: number;          // in user's preferred unit (kg/lb)
  reps: number;
  completed: boolean;
  restTime?: number;       // seconds
  notes?: string;
}
```

### NutritionDay

```typescript
interface NutritionDay {
  id: string;              // UUID
  userId: string;          // UUID
  date: string;            // YYYY-MM-DD
  meals: Meal[];
  goals: MacroGoals;
  totals: {
    calories: number;
    protein: number;       // grams
    carbs: number;         // grams
    fat: number;           // grams
  };
}
```

### Meal

```typescript
interface Meal {
  id: string;              // UUID
  nutritionDayId: string;  // UUID
  name: string;
  mealType: 'breakfast' | 'lunch' | 'dinner' | 'snack';
  time?: string;           // HH:mm
  foods: FoodItem[];
  totals: {
    calories: number;
    protein: number;
    carbs: number;
    fat: number;
  };
}
```

### FoodItem

```typescript
interface FoodItem {
  id: string;              // UUID
  mealId: string;          // UUID
  name: string;
  servingSize: number;
  servingUnit: string;     // g, oz, cup, piece, etc.
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
}
```

### MacroGoals

```typescript
interface MacroGoals {
  calories: number;
  protein: number;         // grams
  carbs: number;           // grams
  fat: number;             // grams
}
```

### ProgressPhoto

```typescript
interface ProgressPhoto {
  id: string;              // UUID
  userId: string;          // UUID
  url: string;             // Cloudflare R2 URL
  thumbnailUrl?: string;
  date: string;            // YYYY-MM-DD
  weight?: number;         // optional weight at time of photo
  notes?: string;
  createdAt: string;       // ISO 8601
}
```

### WeightEntry

```typescript
interface WeightEntry {
  id: string;              // UUID
  userId: string;          // UUID
  weight: number;          // in user's preferred unit
  date: string;            // YYYY-MM-DD
  createdAt: string;       // ISO 8601
}
```

### WorkoutTemplate

```typescript
interface WorkoutTemplate {
  id: string;              // UUID
  userId?: string;         // UUID (null for system templates)
  name: string;
  description?: string;
  exercises: TemplateExercise[];
  isPublic: boolean;       // system templates are public
  category?: string;
  difficulty?: 'beginner' | 'intermediate' | 'advanced';
  createdAt: string;
  updatedAt: string;
}
```

### TemplateExercise

```typescript
interface TemplateExercise {
  name: string;
  sets: number;
  reps: number | string;   // can be "8-12" or number
  restTime?: number;       // seconds
  notes?: string;
}
```

### ExerciseLibrary (Read-only)

```typescript
interface ExerciseInfo {
  id: string;
  name: string;
  description?: string;
  muscleGroup: string;
  equipment: string;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  instructions?: string[];
  videoUrl?: string;
  imageUrl?: string;
}
```

---

## API Response Format

### Success Response

```typescript
interface ApiResponse<T> {
  success: true;
  data: T;
  message?: string;
}
```

### Error Response

```typescript
interface ApiError {
  success: false;
  error: string;
  code: string;           // ERROR_CODE for client handling
  details?: any;
}
```

### Error Codes

```typescript
// Authentication
AUTH_REQUIRED = "Authentication required"
AUTH_INVALID_TOKEN = "Invalid or expired token"
AUTH_INVALID_CREDENTIALS = "Invalid email or password"

// Validation
VALIDATION_ERROR = "Validation failed"
VALIDATION_EMAIL_EXISTS = "Email already exists"

// Not Found
NOT_FOUND = "Resource not found"
WORKOUT_NOT_FOUND = "Workout not found"
USER_NOT_FOUND = "User not found"

// Permissions
FORBIDDEN = "Access forbidden"
NOT_OWNER = "You don't own this resource"

// Rate Limiting
RATE_LIMIT_EXCEEDED = "Too many requests"

// Server
INTERNAL_ERROR = "Internal server error"
```

---

## Authentication Endpoints

### Register

```http
POST /auth/register
```

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securePassword123",
  "name": "John Doe"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "name": "John Doe",
      "createdAt": "2025-01-01T00:00:00Z"
    },
    "tokens": {
      "accessToken": "jwt_token",
      "refreshToken": "refresh_token",
      "expiresIn": 900
    }
  }
}
```

### Login

```http
POST /auth/login
```

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": { /* User object */ },
    "tokens": {
      "accessToken": "jwt_token",
      "refreshToken": "refresh_token",
      "expiresIn": 900
    }
  }
}
```

### Refresh Token

```http
POST /auth/refresh
```

**Request:**
```json
{
  "refreshToken": "refresh_token_here"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "accessToken": "new_jwt_token",
    "refreshToken": "new_refresh_token",
    "expiresIn": 900
  }
}
```

### Logout

```http
POST /auth/logout
Authorization: Bearer {accessToken}
```

**Request:**
```json
{
  "refreshToken": "refresh_token_here"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

---

## Workout Endpoints

### Get All Workouts

```http
GET /workouts?startDate=2025-01-01&endDate=2025-01-31&limit=50&offset=0
Authorization: Bearer {accessToken}
```

**Query Parameters:**
- `startDate` (optional): Filter from date (YYYY-MM-DD)
- `endDate` (optional): Filter to date (YYYY-MM-DD)
- `limit` (optional): Number of results (default: 50)
- `offset` (optional): Pagination offset (default: 0)

**Response:**
```json
{
  "success": true,
  "data": {
    "workouts": [
      {
        "id": "uuid",
        "userId": "uuid",
        "date": "2025-01-15",
        "exercises": [/* Exercise objects */],
        "duration": 60,
        "totalVolume": 5000,
        "createdAt": "2025-01-15T10:00:00Z"
      }
    ],
    "total": 100,
    "hasMore": true
  }
}
```

### Get Today's Workout

```http
GET /workouts/today
Authorization: Bearer {accessToken}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "workout": {
      "id": "uuid",
      "date": "2025-01-15",
      "exercises": [/* Exercise objects */]
    }
  }
}
```

### Get Workout by ID

```http
GET /workouts/:id
Authorization: Bearer {accessToken}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "userId": "uuid",
    "date": "2025-01-15",
    "exercises": [
      {
        "id": "uuid",
        "name": "Bench Press",
        "sets": [
          { "weight": 100, "reps": 10, "completed": true },
          { "weight": 100, "reps": 8, "completed": true }
        ]
      }
    ]
  }
}
```

### Create Workout

```http
POST /workouts
Authorization: Bearer {accessToken}
```

**Request:**
```json
{
  "date": "2025-01-15",
  "exercises": [
    {
      "name": "Bench Press",
      "sets": [
        { "weight": 100, "reps": 10, "completed": false }
      ],
      "muscleGroup": "chest"
    }
  ],
  "notes": "Great workout!"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "userId": "uuid",
    "date": "2025-01-15",
    "exercises": [/* Exercise objects with IDs */],
    "createdAt": "2025-01-15T10:00:00Z"
  }
}
```

### Update Workout

```http
PUT /workouts/:id
Authorization: Bearer {accessToken}
```

**Request:**
```json
{
  "exercises": [
    {
      "id": "existing-uuid",
      "name": "Bench Press",
      "sets": [
        { "weight": 100, "reps": 10, "completed": true }
      ]
    }
  ],
  "notes": "Updated notes"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "exercises": [/* Updated exercises */],
    "updatedAt": "2025-01-15T11:00:00Z"
  }
}
```

### Delete Workout

```http
DELETE /workouts/:id
Authorization: Bearer {accessToken}
```

**Response:**
```json
{
  "success": true,
  "message": "Workout deleted successfully"
}
```

---

## Exercise Library Endpoints

### Get Exercise Library

```http
GET /exercises?muscleGroup=chest&equipment=barbell&search=bench&limit=50
Authorization: Bearer {accessToken}
```

**Query Parameters:**
- `muscleGroup` (optional): Filter by muscle group
- `equipment` (optional): Filter by equipment
- `search` (optional): Search by name
- `limit` (optional): Number of results (default: 50)

**Response:**
```json
{
  "success": true,
  "data": {
    "exercises": [
      {
        "id": "uuid",
        "name": "Bench Press",
        "muscleGroup": "chest",
        "equipment": "barbell",
        "difficulty": "intermediate",
        "description": "...",
        "instructions": ["Step 1", "Step 2"]
      }
    ],
    "total": 150
  }
}
```

### Get Exercise by ID

```http
GET /exercises/:id
Authorization: Bearer {accessToken}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Bench Press",
    "description": "...",
    "muscleGroup": "chest",
    "equipment": "barbell",
    "instructions": ["Step 1", "Step 2"],
    "videoUrl": "https://..."
  }
}
```

---

## Nutrition Endpoints

### Get Nutrition by Date

```http
GET /nutrition?date=2025-01-15
Authorization: Bearer {accessToken}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "date": "2025-01-15",
    "meals": [
      {
        "id": "uuid",
        "name": "Breakfast",
        "mealType": "breakfast",
        "time": "08:00",
        "foods": [
          {
            "id": "uuid",
            "name": "Eggs",
            "servingSize": 2,
            "servingUnit": "large",
            "calories": 140,
            "protein": 12,
            "carbs": 1,
            "fat": 10
          }
        ],
        "totals": {
          "calories": 140,
          "protein": 12,
          "carbs": 1,
          "fat": 10
        }
      }
    ],
    "goals": {
      "calories": 2500,
      "protein": 180,
      "carbs": 250,
      "fat": 80
    },
    "totals": {
      "calories": 2300,
      "protein": 175,
      "carbs": 240,
      "fat": 75
    }
  }
}
```

### Add Meal

```http
POST /nutrition/meals
Authorization: Bearer {accessToken}
```

**Request:**
```json
{
  "date": "2025-01-15",
  "name": "Breakfast",
  "mealType": "breakfast",
  "time": "08:00",
  "foods": [
    {
      "name": "Eggs",
      "servingSize": 2,
      "servingUnit": "large",
      "calories": 140,
      "protein": 12,
      "carbs": 1,
      "fat": 10
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Breakfast",
    "foods": [/* Food items with IDs */],
    "totals": { "calories": 140, "protein": 12, "carbs": 1, "fat": 10 }
  }
}
```

### Update Meal

```http
PUT /nutrition/meals/:id
Authorization: Bearer {accessToken}
```

**Request:**
```json
{
  "name": "Updated Breakfast",
  "foods": [/* Updated food items */]
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Updated Breakfast",
    "foods": [/* Updated foods */],
    "updatedAt": "2025-01-15T09:00:00Z"
  }
}
```

### Delete Meal

```http
DELETE /nutrition/meals/:id
Authorization: Bearer {accessToken}
```

**Response:**
```json
{
  "success": true,
  "message": "Meal deleted successfully"
}
```

### AI Meal Estimation (Future)

```http
POST /nutrition/ai-estimate
Authorization: Bearer {accessToken}
Content-Type: multipart/form-data
```

**Request:**
```
photo: <file>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "estimatedFoods": [
      {
        "name": "Chicken Breast",
        "servingSize": 200,
        "servingUnit": "g",
        "calories": 330,
        "protein": 62,
        "carbs": 0,
        "fat": 7,
        "confidence": 0.85
      }
    ],
    "totalCalories": 330,
    "totalProtein": 62
  }
}
```

---

## Progress Endpoints

### Get Progress Stats

```http
GET /progress/stats?startDate=2025-01-01&endDate=2025-01-31
Authorization: Bearer {accessToken}
```

**Query Parameters:**
- `startDate` (required): Start date for range
- `endDate` (required): End date for range

**Response:**
```json
{
  "success": true,
  "data": {
    "workouts": {
      "total": 20,
      "totalVolume": 100000,
      "avgDuration": 55,
      "volumeByDate": [
        { "date": "2025-01-01", "volume": 5000 },
        { "date": "2025-01-02", "volume": 5500 }
      ]
    },
    "exercises": [
      {
        "name": "Bench Press",
        "totalSets": 60,
        "totalReps": 600,
        "maxWeight": 120,
        "avgWeight": 105,
        "progressData": [
          { "date": "2025-01-01", "maxWeight": 100 },
          { "date": "2025-01-08", "maxWeight": 105 }
        ]
      }
    ],
    "bodyWeight": {
      "current": 80,
      "start": 85,
      "change": -5,
      "history": [
        { "date": "2025-01-01", "weight": 85 },
        { "date": "2025-01-15", "weight": 82.5 }
      ]
    }
  }
}
```

### Get Progress Photos

```http
GET /progress/photos?startDate=2025-01-01&endDate=2025-01-31&limit=20
Authorization: Bearer {accessToken}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "photos": [
      {
        "id": "uuid",
        "url": "https://r2.cloudflare.com/...",
        "thumbnailUrl": "https://r2.cloudflare.com/.../thumb",
        "date": "2025-01-15",
        "weight": 80,
        "notes": "Front view",
        "createdAt": "2025-01-15T10:00:00Z"
      }
    ],
    "total": 15
  }
}
```

### Upload Progress Photo

```http
POST /progress/photos
Authorization: Bearer {accessToken}
Content-Type: multipart/form-data
```

**Request:**
```
photo: <file>
date: 2025-01-15
weight: 80
notes: Front view
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "url": "https://r2.cloudflare.com/...",
    "thumbnailUrl": "https://r2.cloudflare.com/.../thumb",
    "date": "2025-01-15",
    "createdAt": "2025-01-15T10:00:00Z"
  }
}
```

### Delete Progress Photo

```http
DELETE /progress/photos/:id
Authorization: Bearer {accessToken}
```

**Response:**
```json
{
  "success": true,
  "message": "Photo deleted successfully"
}
```

### Get Weight History

```http
GET /progress/weight?startDate=2025-01-01&endDate=2025-01-31
Authorization: Bearer {accessToken}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "entries": [
      {
        "id": "uuid",
        "weight": 85,
        "date": "2025-01-01",
        "createdAt": "2025-01-01T08:00:00Z"
      },
      {
        "id": "uuid",
        "weight": 84.5,
        "date": "2025-01-08",
        "createdAt": "2025-01-08T08:00:00Z"
      }
    ],
    "current": 84.5,
    "start": 85,
    "change": -0.5
  }
}
```

### Log Weight

```http
POST /progress/weight
Authorization: Bearer {accessToken}
```

**Request:**
```json
{
  "weight": 84.5,
  "date": "2025-01-15"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "weight": 84.5,
    "date": "2025-01-15",
    "createdAt": "2025-01-15T08:00:00Z"
  }
}
```

---

## Template Endpoints

### Get Templates

```http
GET /templates?category=strength&difficulty=intermediate
Authorization: Bearer {accessToken}
```

**Query Parameters:**
- `category` (optional): Filter by category
- `difficulty` (optional): Filter by difficulty
- `isPublic` (optional): true for system templates, false for user templates

**Response:**
```json
{
  "success": true,
  "data": {
    "templates": [
      {
        "id": "uuid",
        "name": "Push Day",
        "description": "Chest, shoulders, and triceps",
        "exercises": [
          {
            "name": "Bench Press",
            "sets": 4,
            "reps": "8-10",
            "restTime": 120
          }
        ],
        "isPublic": true,
        "difficulty": "intermediate"
      }
    ]
  }
}
```

### Get Template by ID

```http
GET /templates/:id
Authorization: Bearer {accessToken}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Push Day",
    "exercises": [/* Template exercises */],
    "isPublic": true
  }
}
```

### Create Template

```http
POST /templates
Authorization: Bearer {accessToken}
```

**Request:**
```json
{
  "name": "My Custom Workout",
  "description": "Custom description",
  "exercises": [
    {
      "name": "Squat",
      "sets": 5,
      "reps": 5,
      "restTime": 180,
      "notes": "Heavy"
    }
  ],
  "isPublic": false
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "My Custom Workout",
    "userId": "uuid",
    "exercises": [/* Exercises */],
    "createdAt": "2025-01-15T10:00:00Z"
  }
}
```

---

## User Profile Endpoints

### Get Profile

```http
GET /users/profile
Authorization: Bearer {accessToken}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "age": 30,
    "height": 180,
    "weight": 80,
    "unit": "kg",
    "goals": {
      "calories": 2500,
      "protein": 180,
      "carbs": 250,
      "fat": 80
    },
    "stats": {
      "totalWorkouts": 100,
      "totalVolume": 500000,
      "currentStreak": 5
    }
  }
}
```

### Update Profile

```http
PUT /users/profile
Authorization: Bearer {accessToken}
```

**Request:**
```json
{
  "name": "John Updated",
  "age": 31,
  "height": 180,
  "weight": 78,
  "unit": "kg",
  "goals": {
    "calories": 2600,
    "protein": 190,
    "carbs": 260,
    "fat": 85
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "John Updated",
    "age": 31,
    "updatedAt": "2025-01-15T10:00:00Z"
  }
}
```

---

## Rate Limits

```
Authentication:
  - Login: 5 attempts per 15 minutes per IP
  - Register: 3 attempts per hour per IP

API Endpoints:
  - General: 100 requests per minute per user
  - Uploads: 10 files per hour per user
  - AI Estimation: 20 requests per day per user
```

---

## Pagination

All list endpoints support pagination:

```
?limit=20&offset=0
```

Response includes:
```json
{
  "data": [...],
  "total": 100,
  "limit": 20,
  "offset": 0,
  "hasMore": true
}
```

---

## File Upload

### Supported Formats
- Images: JPEG, PNG, WebP
- Max size: 10MB per file

### Upload Response
```json
{
  "success": true,
  "data": {
    "url": "https://r2.cloudflare.com/gymie/photos/uuid.jpg",
    "thumbnailUrl": "https://r2.cloudflare.com/gymie/photos/uuid_thumb.jpg",
    "size": 2048576,
    "mimeType": "image/jpeg"
  }
}
```

---

## WebSocket Events (Future)

### Connection
```
wss://api.gymie.com/ws?token={jwt_token}
```

### Events
```typescript
// Workout started
{ type: "workout.started", workoutId: "uuid" }

// Set completed
{ type: "set.completed", exerciseId: "uuid", setIndex: 0 }

// Rest timer
{ type: "rest.started", duration: 60 }
{ type: "rest.completed" }
```

---

## Testing

### Test User Credentials

```
Email: test@gymie.com
Password: Test123!
```

### Postman Collection

Available at: `docs/postman/Gymie-API.postman_collection.json`

---

## Implementation Notes

1. **All dates** should be in YYYY-MM-DD format for consistency
2. **All timestamps** should be in ISO 8601 format with UTC timezone
3. **All IDs** should be UUIDs (v4)
4. **Weights** are always stored in user's preferred unit (kg or lb)
5. **Macros** are always in grams
6. **Calories** are always integers
7. **All arrays** can be empty but never null
8. **Optional fields** can be null or omitted

---

## Database Relationships

```
User
 ├── Workouts
 │    └── Exercises
 │         └── Sets (JSONB)
 ├── NutritionDays
 │    └── Meals
 │         └── Foods
 ├── ProgressPhotos
 ├── WeightEntries
 └── Templates
      └── TemplateExercises

ExerciseLibrary (standalone, read-only)
```

---

## Frontend Integration

The React Native app expects this exact API structure. See [`src/services/localApi.ts`](../src/services/localApi.ts:1) for the current mock implementation that should be replaced with real API calls.

### Environment Variables

```typescript
// .env
API_BASE_URL=https://api.gymie.com/v1
API_TIMEOUT=30000
```

### API Client Setup

```typescript
import axios from 'axios';

const api = axios.create({
  baseURL: process.env.API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = getStoredToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
```

---

This API contract is the complete specification needed to build the Go backend. All data models, endpoints, and behaviors are based on the current React Native implementation.
