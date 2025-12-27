
# Offline-First Nutrition System Architecture

## 🎯 Core Philosophy

**The app must work 100% offline. Online connectivity enriches, never blocks.**

---

## 📱 System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        USER DEVICE                               │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  FOOD IMAGE CAPTURE                                        │  │
│  └─────────────────┬──────────────────────────────────────────┘  │
│                    │                                              │
│                    ▼                                              │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  ON-DEVICE ML (MobileNetV3-Small)                         │  │
│  │  • Size: ~5-7 MB                                          │  │
│  │  • Input: 224x224 RGB                                     │  │
│  │  • Output: dish_id + confidence                           │  │
│  │  • Inference: <100ms on low-end devices                   │  │
│  └─────────────────┬──────────────────────────────────────────┘  │
│                    │                                              │
│                    ▼                                              │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  RULE-BASED PORTION ESTIMATION                            │  │
│  │  • No ML, just image area heuristics                      │  │
│  │  • Output: small | medium | large                         │  │
│  └─────────────────┬──────────────────────────────────────────┘  │
│                    │                                              │
│                    ▼                                              │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  LOCAL SQLite NUTRITION DATABASE                          │  │
│  │  • 500+ Indian dishes                                     │  │
│  │  • Deterministic nutrition calculation                    │  │
│  │  • Portion multipliers applied                            │  │
│  └─────────────────┬──────────────────────────────────────────┘  │
│                    │                                              │
│                    ▼                                              │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  INSTANT RESULT DISPLAY                                   │  │
│  │  • No network required                                    │  │
│  │  • User can correct if needed                             │  │
│  └─────────────────┬──────────────────────────────────────────┘  │
│                    │                                              │
│                    ▼ (when online)                               │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  BACKGROUND SYNC                                          │  │
│  │  • Upload corrections                                     │  │
│  │  • Check for DB updates                                   │  │
│  │  • Check for model updates                                │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                                │
                                │ (async, non-blocking)
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      BACKEND (Go)                                │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  CORRECTION AGGREGATION                                   │  │
│  │  • Learn from user corrections                            │  │
│  │  • Improve dish aliases                                   │  │
│  │  • Update nutrition DB                                    │  │
│  └───────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  RECIPE GENERATION (LLM)                                  │  │
│  │  • Generate new recipes on-demand                         │  │
│  │  • Nutrition calculated server-side                       │  │
│  └───────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  MODEL & DB VERSIONING                                    │  │
│  │  • Push updates to devices                                │  │
│  │  • Rollback support                                       │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🔑 Key Design Decisions

### 1. Why MobileNetV3-Small?

| Metric | Value | Why It Matters |
|--------|-------|----------------|
| Model Size | 5-7 MB | Works on 2GB RAM Android devices |
| Inference Time | <100ms | Real-time user experience |
| CPU Only | Yes | No GPU required |
| Accuracy | ~75% top-1 | Good enough with user correction |
| Battery Impact | Minimal | Efficient architecture |

### 2. Why Local SQLite?

- **Instant Results**: No network latency
- **Offline Support**: Works on flights, in basements
- **Deterministic**: Same input = same output
- **Privacy**: No data leaves device unless user syncs
- **Low Memory**: <5MB for 500 dishes

### 3. Why Rule-Based Portion?

- **Fast**: No ML inference needed
- **Predictable**: Users understand the logic
- **Correctable**: Easy to adjust if wrong
- **Universal**: Works across all platforms

---

## 📊 Dish ID Taxonomy

### Design Principles

1. **Uppercase Snake Case**: `DAL_MAKHANI`, not `dal makhani`
2. **Language Agnostic**: Use common transliterations
3. **Hierarchical**: Parent → Child relationships
4. **Stable**: Never change existing IDs, only add
5. **Descriptive**: Clear what the dish is

### Taxonomy Structure

```
CATEGORY_TYPE_VARIANT_PREPARATION
```

Examples:
- `DAL_MAKHANI` (category_type)
- `RICE_BASMATI_PLAIN` (category_type_variant)
- `CHICKEN_BREAST_GRILLED` (category_type_preparation)
- `ROTI_WHOLE_WHEAT` (category_type_variant)

### Taxonomy File (Shared)

```json
{
  "version": "1.0.0",
  "dishes": {
    "DAL_MAKHANI": {
      "display_name": "Dal Makhani",
      "aliases": ["makhani dal", "black dal", "kaali dal"],
      "category": "dal",
      "cuisine": "indian"
    },
    "CHICKEN_BIRYANI": {
      "display_name": "Chicken Biryani",
      "aliases": ["chicken biryani", "murgh biryani"],
      "category": "rice",
      "cuisine": "indian"
    }
  }
}
```

---

## 💾 SQLite Schema

### Database: `nutrition.db`

```sql
-- Dish Master Table
CREATE TABLE dish_master (
    dish_id TEXT PRIMARY KEY,
    display_name TEXT NOT NULL,
    category TEXT NOT NULL,
    cuisine TEXT NOT NULL,
    aliases TEXT, -- JSON array
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- Nutrition Data Table
CREATE TABLE dish_nutrition (
    dish_id TEXT PRIMARY KEY,
    base_serving_grams INTEGER NOT NULL,
    calories REAL NOT NULL,
    protein REAL NOT NULL,
    carbs REAL NOT NULL,
    fat REAL NOT NULL,
    fiber REAL DEFAULT 0,
    sodium REAL DEFAULT 0,
    FOREIGN KEY (dish_id) REFERENCES dish_master(dish_id)
);

-- User Corrections (Local Cache)
CREATE TABLE user_corrections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    image_hash TEXT NOT NULL,
    predicted_dish_id TEXT NOT NULL,
    corrected_dish_id TEXT,
    predicted_portion TEXT NOT NULL,
    corrected_portion TEXT,
    confidence REAL NOT NULL,
    synced INTEGER DEFAULT 0,
    created_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- Model Metadata
CREATE TABLE model_metadata (
    model_type TEXT PRIMARY KEY, -- 'vision' or 'nutrition_db'
    version TEXT NOT NULL,
    checksum TEXT NOT NULL,
    size_bytes INTEGER NOT NULL,
    updated_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- Indexes
CREATE INDEX idx_dish_category ON dish_master(category);
CREATE INDEX idx_corrections_synced ON user_corrections(synced);
CREATE INDEX idx_corrections_dish ON user_corrections(predicted_dish_id);
```

---

## 🎨 Portion Estimation Algorithm

### Input
- Image dimensions (W x H)
- Dish classification confidence
- (Optional) Plate detection bounding box

### Algorithm

```python
def estimate_portion(image_width, image_height, dish_area):
    """
    Rule-based portion estimation
    No ML, just heuristics
    """
    total_pixels = image_width * image_height
    dish_ratio = dish_area / total_pixels
    
    # Thresholds tuned from user studies
    if dish_ratio < 0.15:
        return "small"
    elif dish_ratio < 0.35:
        return "medium"
    else:
        return "large"
```

### Portion Multipliers

```json
{
  "small": 0.75,
  "medium": 1.0,
  "large": 1.3
}
```

### Nutrition Calculation

```python
def calculate_nutrition(dish_id, portion):
    # Query SQLite
    base_nutrition = db.query("SELECT * FROM dish_nutrition WHERE dish_id = ?", dish_id)
    
    # Apply multiplier
    multiplier = PORTION_MULTIPLIERS[portion]
    
    return {
        "calories": base_nutrition.calories * multiplier,
        "protein": base_nutrition.protein * multiplier,
        "carbs": base_nutrition.carbs * multiplier,
        "fat": base_nutrition.fat * multiplier
    }
```

---

## 🔄 Sync Flow

### When Device Comes Online

```
1. Check model versions
   GET /api/sync/model-versions
   
2. If nutrition DB outdated:
   GET /api/sync/nutrition-db
   
3. If vision model outdated:
   GET /api/sync/vision-model
   
4. Upload pending corrections:
   POST /api/sync/corrections
   
5. Mark as synced in local DB
```

### Sync Priorities

1. **Critical**: Nutrition DB updates (security fixes)
2. **High**: Model updates (accuracy improvements)
3. **Normal**: User corrections upload
4. **Low**: Analytics data

---

## 📲 Cross-Platform Implementation

### Android (Kotlin + TensorFlow Lite)

```kotlin
// 1. Load model
val model = Interpreter(loadModelFile("mobilenet_v3.tflite"))

// 2. Preprocess image
val bitmap = resizeBitmap(image, 224, 224)
val input = bitmapToFloatArray(bitmap)

// 3. Run inference
val output = Array(1) { FloatArray(NUM_CLASSES) }
model.run(input, output)

// 4. Get dish_id
val dishId = getDishIdFromIndex(output[0].maxIndex())

// 5. Query SQLite
val nutrition = nutritionDb.query(dishId, portion)

// 6. Display instantly
showResult(dishId, nutrition)
```

### iOS (Swift + CoreML)

```swift
// 1. Load model
let model = try VNCoreMLModel(for: MobileNetV3().model)

// 2. Create request
let request = VNCoreMLRequest(model: model) { request, error in
    guard let results = request.results as? [VNClassificationObservation] else { return }
    let dishId = results.first?.identifier
    
    // 3. Query SQLite
    let nutrition = NutritionDB.shared.query(dishId: dishId, portion: portion)
    
    // 4. Display
    DispatchQueue.main.async {
        self.showResult(dishId: dishId, nutrition: nutrition)
    }
}

// 5. Perform request
let handler = VNImageRequestHandler(cgImage: image.cgImage!)
try handler.perform([request])
```

### Web (ONNX + WebAssembly)

```javascript
// 1. Load ONNX model
const session = await ort.InferenceSession.create('mobilenet_v3.onnx');

// 2. Preprocess image
const tensor = preprocessImage(imageData, 224, 224);

// 3. Run inference
const results = await session.run({ input: tensor });
const dishId = getDishIdFromOutput(results.output);

// 4. Query IndexedDB
const nutrition = await nutritionDB.query(dishId, portion);

// 5. Display instantly
displayResult(dishId, nutrition);
```

---

## 🔐 Security & Privacy

### On Device
- All processing happens locally
- No data sent to server unless user explicitly syncs
- Corrections stored locally, uploaded in batch

### Backend
- Corrections are anonymized
- No images stored on server
- Only image hashes + dish IDs

---

## 📈 Why This Scales

### Low-End Devices (2GB RAM, Low-End CPU)
- Model: 5-7 MB (fits in memory easily)
- Inference: <100ms (acceptable latency)
- SQLite: <5 MB (minimal footprint)
- No GPU required

### Network Resilience
- Works 100% offline
- Sync happens in background
- Failed syncs retry automatically
- No blocking UI

### Battery Life
- CPU-only inference (efficient)
- No continuous camera processing
- Sync batched to minimize radio usage

---

## 🎯 Success Metrics

### Technical
- Model size: <10 MB ✅
- Inference time: <100ms ✅
- Offline functionality: 100% ✅
- Low-end device support: Yes ✅

### User Experience
- Time to result: <1 second
- Correction rate: <20% (with good model)
- Offline usage: 80%+ of interactions
- Battery drain: <2% per 10 predictions

---

## 🚀 Deployment Strategy

### Phase 1: Initial Release
- Bundle model v1.0
- Bundle 500 dishes
- Enable offline mode
- Basic sync

### Phase 2: Intelligence
- User corrections aggregated
- Model v1.1 with better aliases
- Nutrition DB v1.1 with community dishes

### Phase 3: Advanced
- Multi-dish detection
- Ingredient recognition
- Recipe suggestions

---

## 📝 Model Update Protocol

### Version Check
```http
GET /api/sync/model-versions
Response:
{
  "vision_model": {
    "version": "1.0.1",
    "checksum": "sha256:abc123...",
    "size_bytes": 6543210,
    "url": "https://cdn.example.com/models/vision_v1.0.1.tflite"
  },
  "nutrition_db": {
    "version": "1.0.2",
    "checksum": "sha256:def456...",
    "size_bytes": 4567890,
    "url": "https://cdn.example.com/db/nutrition_v1.0.2.db"
  }
}
```

### Update Flow
1. Device checks local version
2. If outdated, downloads in background
3. Verifies checksum
4. Atomically replaces old model/DB
5. Fallback to old version if corrupted

---

## 🧪 Testing Strategy

### On-Device Tests
- Model inference accuracy
- Portion estimation accuracy
- SQLite query performance
- Offline functionality
- Low battery behavior

### Backend Tests
- Sync API reliability
- Model versioning
- Correction aggregation
- Rollback scenarios

---

**This architecture ensures the app works perfectly offline while continuously improving through backend intelligence.**
