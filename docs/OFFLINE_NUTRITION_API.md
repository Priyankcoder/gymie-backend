
# Offline Nutrition API Documentation

This document describes the REST API endpoints for the offline-first nutrition system.

## Overview

The Offline Nutrition API provides endpoints for:
- Syncing user corrections from devices
- Fetching model versions and updates
- Searching dishes and nutrition data
- Accessing correction analytics

All endpoints require authentication via JWT token in the `Authorization` header:
```
Authorization: Bearer <jwt_token>
```

---

## Endpoints

### 1. Sync User Corrections

Upload nutrition prediction corrections from the client device.

**Endpoint:** `POST /v1/sync/corrections`

**Request Body:**
```json
{
  "corrections": [
    {
      "image_hash": "sha256:abc123...",
      "predicted_dish_id": "BIRYANI_CHICKEN",
      "corrected_dish_id": "BIRYANI_MUTTON",
      "predicted_portion": "medium",
      "corrected_portion": "large",
      "confidence": 0.85,
      "device_type": "android",
      "app_version": "1.0.0"
    }
  ]
}
```

**Field Descriptions:**
- `image_hash` (required): SHA-256 hash of the food image for idempotency
- `predicted_dish_id` (required): Dish ID predicted by the model
- `corrected_dish_id` (optional): User's corrected dish ID (null if prediction was correct)
- `predicted_portion` (required): Predicted portion size (small/medium/large)
- `corrected_portion` (optional): User's corrected portion size
- `confidence` (required): Model's confidence score (0.0 - 1.0)
- `device_type` (required): Device platform (android/ios/web)
- `app_version` (optional): App version string

**Response:**
```json
{
  "status": "success",
  "message": "Corrections synced successfully",
  "data": {
    "synced": 1,
    "failed": 0,
    "errors": []
  }
}
```

**Notes:**
- Idempotent: Duplicate corrections (same `image_hash`) are ignored
- Supports batch upload of multiple corrections
- Errors are returned for individual corrections that fail

---

### 2. Get Model Versions

Fetch the latest available model versions for update checks.

**Endpoint:** `GET /v1/sync/model-versions`

**Response:**
```json
{
  "status": "success",
  "message": "Model versions retrieved successfully",
  "data": {
    "vision_model": {
      "version": "1.2.0",
      "checksum": "sha256:def456...",
      "size_bytes": 5242880,
      "download_url": "https://cdn.example.com/models/vision-v1.2.0.tflite"
    },
    "nutrition_db": {
      "version": "2.0.1",
      "checksum": "sha256:ghi789...",
      "size_bytes": 4567890,
      "download_url": "https://cdn.example.com/nutrition-db-v2.0.1.db"
    },
    "server_time": "2024-01-15T10:30:00Z"
  }
}
```

**Usage:**
- Client compares `version` with locally stored version
- If newer version available, download from `download_url`
- Verify download integrity using `checksum` (SHA-256)
- `server_time` helps with clock synchronization

---

### 3. Search Dishes

Search for dishes by name, alias, or category.

**Endpoint:** `GET /v1/sync/dishes/search`

**Query Parameters:**
- `query` (required): Search term (e.g., "biryani", "curry")
- `category` (optional): Filter by category (e.g., "rice", "curry", "bread")
- `limit` (optional): Max results (default: 10, max: 50)

**Example Request:**
```
GET /v1/sync/dishes/search?query=biryani&category=rice&limit=5
```

**Response:**
```json
{
  "status": "success",
  "message": "Dishes found",
  "data": {
    "dishes": [
      {
        "dish_id": "BIRYANI_CHICKEN",
        "display_name": "Chicken Biryani",
        "category": "rice",
        "cuisine": "indian",
        "aliases": ["chicken biryani", "murgh biryani", "biryani"],
        "nutrition": {
          "dish_id": "BIRYANI_CHICKEN",
          "base_serving_grams": 250,
          "calories": 450,
          "protein": 25,
          "carbs": 55,
          "fat": 15,
          "fiber": 3,
          "sodium": 800
        }
      }
    ],
    "total": 1
  }
}
```

**Notes:**
- Search is case-insensitive
- Matches against `display_name`, `dish_id`, and `aliases`
- Returns full nutrition data for each dish

---

### 4. Get Dish By ID

Retrieve a specific dish with its nutrition information.

**Endpoint:** `GET /v1/sync/dishes/:dish_id`

**Example Request:**
```
GET /v1/sync/dishes/BIRYANI_CHICKEN
```

**Response:**
```json
{
  "status": "success",
  "message": "Dish retrieved successfully",
  "data": {
    "dish_id": "BIRYANI_CHICKEN",
    "display_name": "Chicken Biryani",
    "category": "rice",
    "cuisine": "indian",
    "aliases": ["chicken biryani", "murgh biryani"],
    "nutrition": {
      "dish_id": "BIRYANI_CHICKEN",
      "base_serving_grams": 250,
      "calories": 450,
      "protein": 25,
      "carbs": 55,
      "fat": 15,
      "fiber": 3,
      "sodium": 800
    }
  }
}
```

**Error Response (404):**
```json
{
  "status": "error",
  "error": "not_found",
  "message": "Dish not found",
  "details": "dish not found: INVALID_ID"
}
```

---

### 5. Get Correction Statistics

Retrieve analytics about user corrections.

**Endpoint:** `GET /v1/sync/corrections/stats`

**Query Parameters:**
- `dish_id` (optional): Filter by specific dish ID

**Example Request:**
```
GET /v1/sync/corrections/stats?dish_id=BIRYANI_CHICKEN
```

**Response:**
```json
{
  "status": "success",
  "message": "Stats retrieved successfully",
  "data": {
    "total_corrections": 1250,
    "top_corrected_dishes": [
      {
        "predicted_dish_id": "BIRYANI_CHICKEN",
        "count": 45
      },
      {
        "predicted_dish_id": "CURRY_CHICKEN",
        "count": 38
      }
    ],
    "by_device": [
      {
        "device_type": "android",
        "count": 800
      },
      {
        "device_type": "ios",
        "count": 450
      }
    ]
  }
}
```

**Notes:**
- Without `dish_id`: Returns global statistics
- With `dish_id`: Returns statistics for that specific dish
- Useful for identifying model weaknesses

---

### 6. Download Nutrition Database

Get metadata and download URL for the nutrition database.

**Endpoint:** `GET /v1/sync/nutrition-db`

**Response:**
```json
{
  "status": "success",
  "message": "Nutrition DB info",
  "data": {
    "version": "1.0.0",
    "checksum": "sha256:placeholder",
    "size_bytes": 4567890,
    "download_url": "/static/nutrition.db"
  }
}
```

**Notes:**
- Returns metadata for the SQLite nutrition database
- Client should download and verify checksum
- Database contains all dishes and nutrition data for offline use

---

## Error Handling

All endpoints follow a consistent error response format:

```json
{
  "status": "error",
  "error": "error_code",
  "message": "Human-readable error message",
  "details": "Additional error details or validation errors"
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `auth_required` | 401 | Authentication token missing |
| `auth_error` | 401 | Invalid authentication token |
| `validation_error` | 400 | Request validation failed |
| `not_found` | 404 | Resource not found |
| `sync_failed` | 500 | Correction sync failed |
| `fetch_failed` | 500 | Failed to fetch data |
| `search_failed` | 500 | Search operation failed |

---

## Rate Limiting

All endpoints are rate-limited to prevent abuse:
- **Correction Sync**: 100 requests/hour
- **Model Versions**: 60 requests/hour
- **Search**: 120 requests/hour
- **Other endpoints**: 60 requests/hour

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1642252800
```

---

## Client Integration Examples

### Android/Kotlin

```kotlin
// Sync corrections
data class CorrectionRequest(
    val corrections: List<CorrectionData>
)

suspend fun syncCorrections(corrections: List<CorrectionData>) {
    val response = apiClient.post("/v1/sync/corrections") {
        setBody(CorrectionRequest(corrections))
    }
    // Handle response
}

// Check for model updates
suspend fun checkModelUpdates() {
    val response = apiClient.get("/v1/sync/model-versions")
    val versions = response.body<ModelVersionsResponse>()
    
    if (versions.visionModel.version != localModelVersion) {
        downloadModel(versions.visionModel.downloadUrl)
    }
}
```

### iOS/Swift

```swift
// Sync corrections
struct CorrectionRequest: Codable {
    let corrections: [CorrectionData]
}

func syncCorrections(_ corrections: [CorrectionData]) async throws {
    let request = CorrectionRequest(corrections: corrections)
    let response = try await apiClient.post("/v1/sync/corrections", body: request)
    // Handle response
}

// Check for model updates
func checkModelUpdates() async throws -> ModelVersionsResponse {
    let response = try await apiClient.get("/v1/sync/model-versions")
    return try JSONDecoder().decode(ModelVersionsResponse.self, from: response)
}
```

### React Native/TypeScript

```typescript
// Sync corrections
interface CorrectionRequest {
  corrections: CorrectionData[];
}

async function syncCorrections(corrections: CorrectionData[]) {
  const response = await fetch('/v1/sync/corrections', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ corrections })
  });
  
  return await response.json();
}

// Check for model updates
async function checkModelUpdates() {
  const response = await fetch('/v1/sync/model-versions', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  return await response.json();
}
```

---

## Best Practices

### 1. Idempotency
Always use stable `image_hash` values to ensure corrections aren't duplicated:
```javascript
const imageHash = await sha256(imageBlob);
```

### 2. Batch Syncing
Sync corrections in batches to reduce network requests:
```javascript
// Queue corrections locally
const queue = [];
queue.push(correction);

// Sync when queue reaches threshold or on app background
if (queue.length >= 10 || isAppBackgrounding) {
  await syncCorrections(queue);
  queue.clear();
}
```

### 3. Model Updates
Check for model updates periodically (e.g., daily) or on app launch:
```javascript
const lastCheck = localStorage.getItem('last_model_check');
const now = Date.now();

if (!lastCheck || (now - parseInt(lastCheck)) > 86400000) {
  await checkModelUpdates();
  localStorage.setItem('last_model_check', now.toString());
}
```

### 4. Offline Queue
Queue corrections when offline and sync when connection is restored:
```javascript
if (navigator.onLine) {
  await syncCorrections(offlineQueue);
  offlineQueue.clear();
} else {
  offlineQueue.push(correction);
}
```

---

## Security Considerations

1. **Authentication**: All endpoints require valid JWT tokens
2. **HTTPS**: Always use HTTPS in production
3. **Checksum Verification**: Verify model downloads using SHA-256 checksums
4. **Rate Limiting**: Respect rate limits to avoid throttling
5. **Data Privacy**: Image hashes are one-way; original images are never uploaded

---

## Support

For questions or issues:
- Email: support@gymie.app
- Documentation: https://docs.gymie.app
- Issue Tracker: https://github.com/yourorg/gymie/issues
