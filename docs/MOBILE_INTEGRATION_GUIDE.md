
# Mobile Integration Guide: Offline-First Nutrition System

This guide explains how to integrate the offline-first nutrition recognition system into Android, iOS, and React Native applications.

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Setup & Installation](#setup--installation)
4. [Implementation Guide](#implementation-guide)
5. [Model Integration](#model-integration)
6. [Database Integration](#database-integration)
7. [Sync Strategy](#sync-strategy)
8. [Testing](#testing)
9. [Troubleshooting](#troubleshooting)

---

## Overview

### Key Features

- **100% Offline Operation**: Full nutrition recognition without internet
- **On-Device ML**: MobileNetV3-Small for fast, CPU-only inference
- **Local Database**: SQLite with 500+ Indian dishes
- **Background Sync**: Automatic sync when connection available
- **Deterministic Results**: Consistent nutrition data via lookup tables

### Performance Targets

- **Inference Time**: <100ms on mid-range devices
- **Model Size**: <5MB (quantized)
- **Database Size**: <5MB
- **Memory Usage**: <50MB peak
- **Battery Impact**: Minimal (<2% per hour of active use)

---

## Architecture

```
┌─────────────────────────────────────────┐
│           Mobile Application            │
├─────────────────────────────────────────┤
│  ┌──────────────┐  ┌─────────────────┐ │
│  │   Camera     │  │   Nutrition     │ │
│  │   Module     │──│   Recognition   │ │
│  └──────────────┘  │   Service       │ │
│                    └─────────────────┘ │
│         │                  │           │
│         ▼                  ▼           │
│  ┌──────────────┐  ┌─────────────────┐ │
│  │  TFLite/     │  │   SQLite        │ │
│  │  CoreML      │  │   Database      │ │
│  │  Runtime     │  │   (nutrition.db)│ │
│  └──────────────┘  └─────────────────┘ │
│                           │             │
│                           ▼             │
│                    ┌─────────────────┐  │
│                    │  Sync Service   │  │
│                    │  (Background)   │  │
│                    └─────────────────┘  │
│                           │             │
└───────────────────────────┼─────────────┘
                            │ HTTPS
                            ▼
                    ┌─────────────────┐
                    │  Backend API    │
                    │  (Corrections & │
                    │   Updates)      │
                    └─────────────────┘
```

---

## Setup & Installation

### Android (Kotlin/Java)

#### 1. Add Dependencies

**build.gradle (app level):**
```gradle
dependencies {
    // TensorFlow Lite
    implementation 'org.tensorflow:tensorflow-lite:2.14.0'
    implementation 'org.tensorflow:tensorflow-lite-support:0.4.4'
    
    // SQLite
    implementation "androidx.room:room-runtime:2.6.0"
    kapt "androidx.room:room-compiler:2.6.0"
    implementation "androidx.room:room-ktx:2.6.0"
    
    // Coroutines
    implementation 'org.jetbrains.kotlinx:kotlinx-coroutines-android:1.7.3'
    
    // Work Manager (for background sync)
    implementation "androidx.work:work-runtime-ktx:2.9.0"
    
    // Networking
    implementation 'com.squareup.retrofit2:retrofit:2.9.0'
    implementation 'com.squareup.retrofit2:converter-gson:2.9.0'
}
```

#### 2. Add Model Assets

Place these files in `app/src/main/assets/`:
- `mobilenet_v3_small_food.tflite` (vision model)
- `nutrition.db` (nutrition database)

#### 3. Configure Permissions

**AndroidManifest.xml:**
```xml
<uses-permission android:name="android.permission.CAMERA" />
<uses-permission android:name="android.permission.INTERNET" />
<uses-permission android:name="android.permission.ACCESS_NETWORK_STATE" />
```

---

### iOS (Swift)

#### 1. Add Frameworks

**Podfile:**
```ruby
target 'YourApp' do
  use_frameworks!
  
  # Core ML
  # (Built-in, no pod needed)
  
  # SQLite
  pod 'SQLite.swift', '~> 0.14.1'
  
  # Networking
  pod 'Alamofire', '~> 5.8'
end
```

#### 2. Add Model Files

Add to Xcode project:
- `MobileNetV3Food.mlmodel` (vision model)
- `nutrition.db` (nutrition database)

Ensure both are in "Copy Bundle Resources" in Build Phases.

#### 3. Configure Permissions

**Info.plist:**
```xml
<key>NSCameraUsageDescription</key>
<string>We need camera access to recognize food for nutrition tracking</string>
```

---

### React Native (TypeScript)

#### 1. Install Dependencies

```bash
npm install --save \
  @tensorflow/tfjs \
  @tensorflow/tfjs-react-native \
  react-native-fs \
  react-native-sqlite-storage \
  @react-native-async-storage/async-storage

# Link native dependencies
npx pod-install ios
```

#### 2. Setup Assets

Copy model and database files to:
- Android: `android/app/src/main/assets/`
- iOS: Add via Xcode to bundle

---

## Implementation Guide

### Android Implementation

#### 1. Model Wrapper

```kotlin
// NutritionClassifier.kt
class NutritionClassifier(private val context: Context) {
    private lateinit var interpreter: Interpreter
    private lateinit var labels: List<String>
    
    init {
        loadModel()
        loadLabels()
    }
    
    private fun loadModel() {
        val modelBuffer = loadModelFile("mobilenet_v3_small_food.tflite")
        val options = Interpreter.Options().apply {
            setNumThreads(4)
            setUseNNAPI(false) // CPU only for consistency
        }
        interpreter = Interpreter(modelBuffer, options)
    }
    
    private fun loadModelFile(filename: String): ByteBuffer {
        val assetFileDescriptor = context.assets.openFd(filename)
        val inputStream = FileInputStream(assetFileDescriptor.fileDescriptor)
        val fileChannel = inputStream.channel
        val startOffset = assetFileDescriptor.startOffset
        val declaredLength = assetFileDescriptor.declaredLength
        return fileChannel.map(FileChannel.MapMode.READ_ONLY, startOffset, declaredLength)
    }
    
    private fun loadLabels() {
        // Load from dish_taxonomy.json or embedded list
        labels = context.assets.open("dish_labels.txt")
            .bufferedReader()
            .readLines()
    }
    
    fun classify(bitmap: Bitmap): ClassificationResult {
        val startTime = System.currentTimeMillis()
        
        // Preprocess image
        val inputImage = preprocessImage(bitmap)
        
        // Run inference
        val output = Array(1) { FloatArray(labels.size) }
        interpreter.run(inputImage, output)
        
        // Get top prediction
        val probabilities = output[0]
        val maxIndex = probabilities.indices.maxByOrNull { probabilities[it] } ?: 0
        
        val inferenceTime = System.currentTimeMillis() - startTime
        
        return ClassificationResult(
            dishId = labels[maxIndex],
            confidence = probabilities[maxIndex],
            inferenceTimeMs = inferenceTime
        )
    }
    
    private fun preprocessImage(bitmap: Bitmap): ByteBuffer {
        val inputSize = 224
        val byteBuffer = ByteBuffer.allocateDirect(4 * inputSize * inputSize * 3)
        byteBuffer.order(ByteOrder.nativeOrder())
        
        val resizedBitmap = Bitmap.createScaledBitmap(bitmap, inputSize, inputSize, true)
        val intValues = IntArray(inputSize * inputSize)
        resizedBitmap.getPixels(intValues, 0, resizedBitmap.width, 0, 0, 
            resizedBitmap.width, resizedBitmap.height)
        
        var pixel = 0
        for (i in 0 until inputSize) {
            for (j in 0 until inputSize) {
                val value = intValues[pixel++]
                // Normalize to [-1, 1]
                byteBuffer.putFloat(((value shr 16 and 0xFF) - 127.5f) / 127.5f)
                byteBuffer.putFloat(((value shr 8 and 0xFF) - 127.5f) / 127.5f)
                byteBuffer.putFloat(((value and 0xFF) - 127.5f) / 127.5f)
            }
        }
        
        return byteBuffer
    }
    
    fun close() {
        interpreter.close()
    }
}

data class ClassificationResult(
    val dishId: String,
    val confidence: Float,
    val inferenceTimeMs: Long
)
```

#### 2. Database Access

```kotlin
// NutritionDatabase.kt
@Database(entities = [DishEntity::class, NutritionEntity::class], version = 1)
abstract class NutritionDatabase : RoomDatabase() {
    abstract fun dishDao(): DishDao
    
    companion object {
        @Volatile
        private var INSTANCE: NutritionDatabase? = null
        
        fun getDatabase(context: Context): NutritionDatabase {
            return INSTANCE ?: synchronized(this) {
                val instance = Room.databaseBuilder(
                    context.applicationContext,
                    NutritionDatabase::class.java,
                    "nutrition_database"
                )
                .createFromAsset("nutrition.db")
                .build()
                INSTANCE = instance
                instance
            }
        }
    }
}

@Dao
interface DishDao {
    @Query("SELECT * FROM dish_master WHERE dish_id = :dishId")
    suspend fun getDish(dishId: String): DishEntity?
    
    @Query("""
        SELECT d.*, n.* 
        FROM dish_master d 
        LEFT JOIN dish_nutrition_master n ON d.dish_id = n.dish_id 
        WHERE d.dish_id = :dishId
    """)
    suspend fun getDishWithNutrition(dishId: String): DishWithNutrition?
    
    @Query("""
        SELECT * FROM dish_master 
        WHERE LOWER(display_name) LIKE '%' || LOWER(:query) || '%'
        LIMIT :limit
    """)
    suspend fun searchDishes(query: String, limit: Int = 10): List<DishEntity>
}

@Entity(tableName = "dish_master")
data class DishEntity(
    @PrimaryKey val dish_id: String,
    val display_name: String,
    val category: String,
    val cuisine: String
)

@Entity(tableName = "dish_nutrition_master")
data class NutritionEntity(
    @PrimaryKey val dish_id: String,
    val base_serving_grams: Int,
    val calories: Double,
    val protein: Double,
    val carbs: Double,
    val fat: Double,
    val fiber: Double,
    val sodium: Double
)

data class DishWithNutrition(
    @Embedded val dish: DishEntity,
    @Embedded val nutrition: NutritionEntity?
)
```

#### 3. Nutrition Service

```kotlin
// NutritionService.kt
class NutritionService(
    private val classifier: NutritionClassifier,
    private val database: NutritionDatabase,
    private val syncService: SyncService
) {
    suspend fun recognizeFood(bitmap: Bitmap): NutritionResult {
        // 1. Classify image
        val classification = classifier.classify(bitmap)
        
        // 2. Look up nutrition data
        val dishWithNutrition = database.dishDao()
            .getDishWithNutrition(classification.dishId)
        
        if (dishWithNutrition == null) {
            throw IllegalStateException("Dish not found: ${classification.dishId}")
        }
        
        // 3. Estimate portion size (simple heuristic)
        val portionSize = estimatePortion(bitmap)
        val portionMultiplier = when (portionSize) {
            PortionSize.SMALL -> 0.7
            PortionSize.MEDIUM -> 1.0
            PortionSize.LARGE -> 1.5
        }
        
        // 4. Calculate scaled nutrition
        val nutrition = dishWithNutrition.nutrition!!
        val scaledNutrition = ScaledNutrition(
            calories = nutrition.calories * portionMultiplier,
            protein = nutrition.protein * portionMultiplier,
            carbs = nutrition.carbs * portionMultiplier,
            fat = nutrition.fat * portionMultiplier,
            fiber = nutrition.fiber * portionMultiplier,
            sodium = nutrition.sodium * portionMultiplier
        )
        
        return NutritionResult(
            dish = dishWithNutrition.dish,
            nutrition = scaledNutrition,
            portionSize = portionSize,
            confidence = classification.confidence,
            inferenceTimeMs = classification.inferenceTimeMs
        )
    }
    
    private fun estimatePortion(bitmap: Bitmap): PortionSize {
        // Simple heuristic based on image analysis
        // In production, use a more sophisticated approach
        val pixelDensity = calculateFoodPixelDensity(bitmap)
        return when {
            pixelDensity < 0.3 -> PortionSize.SMALL
            pixelDensity > 0.6 -> PortionSize.LARGE
            else -> PortionSize.MEDIUM
        }
    }
    
    suspend fun submitCorrection(
        imageHash: String,
        predictedDishId: String,
        correctedDishId: String?,
        predictedPortion: String,
        correctedPortion: String?,
        confidence: Float
    ) {
        val correction = CorrectionData(
            imageHash = imageHash,
            predictedDishId = predictedDishId,
            correctedDishId = correctedDishId,
            predictedPortion = predictedPortion,
            correctedPortion = correctedPortion,
            confidence = confidence,
            deviceType = "android",
            appVersion = BuildConfig.VERSION_NAME
        )
        
        syncService.queueCorrection(correction)
    }
}

enum class PortionSize {
    SMALL, MEDIUM, LARGE
}

data class NutritionResult(
    val dish: DishEntity,
    val nutrition: ScaledNutrition,
    val portionSize: PortionSize,
    val confidence: Float,
    val inferenceTimeMs: Long
)

data class ScaledNutrition(
    val calories: Double,
    val protein: Double,
    val carbs: Double,
    val fat: Double,
    val fiber: Double,
    val sodium: Double
)
```

#### 4. Background Sync Worker

```kotlin
// SyncWorker.kt
class SyncWorker(
    context: Context,
    params: WorkerParameters
) : CoroutineWorker(context, params) {
    
    override suspend fun doWork(): Result {
        val syncService = SyncService(applicationContext)
        
        return try {
            // Sync corrections
            syncService.syncPendingCorrections()
            
            // Check for model updates
            syncService.checkForUpdates()
            
            Result.success()
        } catch (e: Exception) {
            Log.e("SyncWorker", "Sync failed", e)
            Result.retry()
        }
    }
}

// Schedule periodic sync
fun schedulePeriodicSync(context: Context) {
    val constraints = Constraints.Builder()
        .setRequiredNetworkType(NetworkType.CONNECTED)
        .setRequiresBatteryNotLow(true)
        .build()
    
    val syncWork = PeriodicWorkRequestBuilder<SyncWorker>(
        repeatInterval = 6,
        repeatIntervalTimeUnit = TimeUnit.HOURS
    )
    .setConstraints(constraints)
    .setBackoffCriteria(
        BackoffPolicy.EXPONENTIAL,
        WorkRequest.MIN_BACKOFF_MILLIS,
        TimeUnit.MILLISECONDS
    )
    .build()
    
    WorkManager.getInstance(context)
        .enqueueUniquePeriodicWork(
            "nutrition_sync",
            ExistingPeriodicWorkPolicy.KEEP,
            syncWork
        )
}
```

---

### iOS Implementation

#### 1. Model Wrapper

```swift
// NutritionClassifier.swift
import CoreML
import Vision

class NutritionClassifier {
    private var model: VNCoreMLModel?
    private let labels: [String]
    
    init() throws {
        // Load model
        let configuration = MLModelConfiguration()
        configuration.computeUnits = .cpuOnly // Force CPU for consistency
        
        let mlModel = try MobileNetV3Food(configuration: configuration)
        self.model = try VNCoreMLModel(for: mlModel.model)
        
        // Load labels
        self.labels = try loadLabels()
    }
    
    private func loadLabels() throws -> [String] {
        guard let path = Bundle.main.path(forResource: "dish_labels", ofType: "txt"),
              let content = try? String(contentsOfFile: path) else {
            throw NSError(domain: "NutritionClassifier", code: -1)
        }
        return content.components(separatedBy: .newlines)
    }
    
    func classify(image: UIImage, completion: @escaping (ClassificationResult) -> Void) {
        guard let model = model,
              let ciImage = CIImage(image: image) else {
            return
        }
        
        let startTime = Date()
        
        let request = VNCoreMLRequest(model: model) { [weak self] request, error in
            guard let results = request.results as? [VNClassificationObservation],
                  let topResult = results.first,
                  let labels = self?.labels else {
                return
            }
            
            let inferenceTime = Date().timeIntervalSince(startTime) * 1000
            
            let result = ClassificationResult(
                dishId: topResult.identifier,
                confidence: topResult.confidence,
                inferenceTimeMs: Int64(inferenceTime)
            )
            
            completion(result)
        }
        
        let handler = VNImageRequestHandler(ciImage: ciImage, options: [:])
        try? handler.perform([request])
    }
}

struct ClassificationResult {
    let dishId: String
    let confidence: Float
    let inferenceTimeMs: Int64
}
```

#### 2. Database Access

```swift
// NutritionDatabase.swift
import SQLite

class NutritionDatabase {
    private var db: Connection?
    
    // Tables
    private let dishMaster = Table("dish_master")
    private let nutritionMaster = Table("dish_nutrition_master")
    
    // Columns
    private let dishId = Expression<String>("dish_id")
    private let displayName = Expression<String>("display_name")
    private let category = Expression<String>("category")
    private let cuisine = Expression<String>("cuisine")
    private let calories = Expression<Double>("calories")
    private let protein = Expression<Double>("protein")
    private let carbs = Expression<Double>("carbs")
    private let fat = Expression<Double>("fat")
    
    init() throws {
        guard let path = Bundle.main.path(forResource: "nutrition", ofType: "db") else {
            throw DatabaseError.fileNotFound
        }
        db = try Connection(path, readonly: true)
    }
    
    func getDish(dishId: String) throws -> Dish? {
        guard let db = db else { return nil }
        
        let query = dishMaster.filter(self.dishId == dishId)
        
        guard let row = try db.pluck(query) else {
            return nil
        }
        
        return Dish(
            dishId: try row.get(self.dishId),
            displayName: try row.get(displayName),
            category: try row.get(category),
            cuisine: try row.get(cuisine)
        )
    }
    
    func getNutrition(dishId: String) throws -> Nutrition? {
        guard let db = db else { return nil }
        
        let query = nutritionMaster.filter(self.dishId == dishId)
        
        guard let row = try db.pluck(query) else {
            return nil
        }
        
        return Nutrition(
            calories: try row.get(calories),
            protein: try row.get(protein),
            carbs: try row.get(carbs),
            fat: try row.get(fat)
        )
    }
    
    func searchDishes(query: String, limit: Int = 10) throws -> [Dish] {
        guard let db = db else { return [] }
        
        let search = dishMaster
            .filter(displayName.like("%\(query)%"))
            .limit(limit)
        
        return try db.prepare(search).map { row in
            Dish(
                dishId: try row.get(self.dishId),
                displayName: try row.get(displayName),
                category: try row.get(category),
                cuisine: try row.get(cuisine)
            )
        }
    }
}

struct Dish {
    let dishId: String
    let displayName: String
    let category: String
    let cuisine: String
}

struct Nutrition {
    let calories: Double
    let protein: Double
    let carbs: Double
    let fat: Double
}

enum DatabaseError: Error {
    case fileNotFound
}
```

---

## Sync Strategy

### Correction Queue

Store corrections locally and sync when online:

```kotlin
// Android
class CorrectionQueue(private val context: Context) {
    private val prefs = context.getSharedPreferences("corrections", Context.MODE_PRIVATE)
    private val gson = Gson()
    
    fun add(correction: CorrectionData) {
        val queue = getQueue().toMutableList()
        queue.add(correction)
        saveQueue(queue)
    }
    
    fun getQueue(): List<CorrectionData> {
        val json = prefs.getString("queue", "[]") ?: "[]"
        return gson.fromJson(json, Array<CorrectionData>::class.java).toList()
    }
    
    fun clear() {
        prefs.edit().putString("queue", "[]").apply()
    }
    
    private fun saveQueue(queue: List<CorrectionData>) {
        val json = gson.toJson(queue)
        prefs.edit().putString("queue", json).apply()
    }
}
```

### Model Updates

Check for updates and download in background:

```kotlin
suspend fun checkForUpdates() {
    val api = createApiService()
    val response = api.getModelVersions()
    
    val localVersion = getLocalModelVersion()
    
    if (response.visionModel.version != localVersion) {
        downloadModel(response.visionModel.downloadUrl)
    }
}

private suspend fun downloadModel(url: String) {
    withContext(Dispatchers.IO) {
        // Download model
        val file = File(context.filesDir, "model_update.tflite")
        // ... download logic
        
        // Verify checksum
        val checksum = calculateSHA256(file)
        if (checksum == expectedChecksum) {
            // Replace old model
            file.renameTo(File(context.filesDir, "mobilenet_v3_small_food.tflite"))
        }
    }
}
```

---

## Testing

### Unit Tests

```kotlin
@Test
fun testClassification() = runTest {
    val bitmap = loadTestImage("biryani.jpg")
    val result = classifier.classify(bitmap)
    
    assertEquals("BIRYANI_CHICKEN", result.dishId)
    assertTrue(result.confidence > 0.5)
    assertTrue(result.inferenceTimeMs < 200)
}

@Test
fun testNutritionLookup() = runTest {
    val dish = database.dishDao().getDishWithNutrition("BIRYANI_CHICKEN")
    
    assertNotNull(dish)
    assertEquals("Chicken Biryani", dish?.dish?.display_name)
    assertTrue(dish?.nutrition?.calories ?: 0.0 > 400)
}
```

### Integration Tests

```kotlin
@Test
fun testEndToEndRecognition() = runTest {
    val service = NutritionService(classifier, database, syncService)
    val bitmap = loadTestImage("biryani.jpg")
    
    val result = service.recognizeFood(bitmap)
    
    assertNotNull(result)
    assertEquals("BIRYANI_CHICKEN", result.dish.dish_id)
    assertTrue(result.nutrition.calories > 0)
}
```

---

## Troubleshooting

### Common Issues

**1. Model Not Loading**
```
Error: Failed to load model
Solution: Ensure model file is in assets/ and filename is correct
```

**2. Slow Inference**
```
Issue: Inference taking >500ms
Solutions:
- Verify CPU-only mode is enabled
- Check image preprocessing (should resize to 224x224)
- Profile with Android Profiler
```

**3. Database Not Found**
```
Error: Unable to open database
Solution: Verify nutrition.db is in assets/ and createFromAsset() is called
```

**4. Sync Failures**
```
Issue: Corrections not syncing
Solutions:
- Check network connectivity
- Verify auth token is valid
- Check WorkManager logs
```

### Performance Optimization

1. **Image Preprocessing**: Resize images before classification
2. **Database Queries**: Add indexes on frequently searched columns
3. **Memory Management**: Recycle bitmaps after use
4. **Background Work**: Use WorkManager for sync operations

---

## Next Steps

1. Implement camera integration
2. Add user correction UI
3. Integrate with nutrition tracking
4. Implement offline indicator
5. Add unit and integration tests

For additional support, refer to:
- API Documentation: `OFFLINE_NUTRITION_API.md`
- Architecture Documentation: `OFFLINE_FIRST_NUTRITION_ARCHITECTURE.md`
