
# Nutrition Database Generator

This tool generates a SQLite database from the dish taxonomy JSON file for distribution to mobile clients.

## Usage

```bash
cd backend/cmd/generate-nutrition-db
go run main.go
```

## Output

The tool generates `nutrition.db` in the current directory with the following structure:

### Tables

#### dish_master
- `dish_id` (PRIMARY KEY) - Uppercase snake case dish identifier
- `display_name` - Human-readable dish name
- `category` - Dish category (rice, curry, bread, etc.)
- `cuisine` - Cuisine type (indian, chinese, etc.)
- `aliases` - JSON array of alternative names
- `description` - Dish description
- `created_at`, `updated_at` - Timestamps

#### dish_nutrition
- `dish_id` (PRIMARY KEY) - References dish_master
- `base_serving_grams` - Standard serving size in grams
- `calories` - Calories per serving
- `protein` - Protein in grams
- `carbs` - Carbohydrates in grams
- `fat` - Fat in grams
- `fiber` - Fiber in grams
- `sodium` - Sodium in mg
- `created_at`, `updated_at` - Timestamps

#### model_metadata
- `model_type` - Type of model (vision, nutrition_db)
- `version` - Semantic version
- `checksum` - SHA-256 checksum for integrity
- `size_bytes` - File size in bytes
- `created_at` - Creation timestamp

#### user_corrections (empty, for local queuing)
- Structure for storing corrections before sync

## Verification

```bash
# Check table counts
sqlite3 nutrition.db "SELECT COUNT(*) FROM dish_master;"

# View sample data
sqlite3 nutrition.db "SELECT * FROM dish_master LIMIT 5;"

# Check nutrition data
sqlite3 nutrition.db "SELECT d.display_name, n.calories, n.protein 
FROM dish_master d JOIN dish_nutrition n ON d.dish_id = n.dish_id 
ORDER BY n.calories DESC LIMIT 10;"
```

## Distribution

The generated `nutrition.db` file should be:
1. Copied to mobile app assets (Android: `assets/`, iOS: bundle)
2. Uploaded to CDN for OTA updates
3. Versioned in model_metadata table

## Size

Current database size: ~60KB with 47 dishes
Expected size with 500+ dishes: ~5MB

## Updates

To update the database:
1. Edit `backend/data/dish_taxonomy.json`
2. Re-run this generator
3. Test with sample queries
4. Distribute updated database to clients
