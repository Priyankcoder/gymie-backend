
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// DishTaxonomy represents the JSON structure
type DishTaxonomy struct {
	Version            string                       `json:"version"`
	LastUpdated        string                       `json:"last_updated"`
	Dishes             map[string]DishInfo          `json:"dishes"`
	PortionMultipliers map[string]float64           `json:"portion_multipliers"`
	ConfidenceThreshold float64                     `json:"confidence_threshold"`
}

type DishInfo struct {	
	DisplayName string   `json:"display_name"`
	Aliases     []string `json:"aliases"`
	Category    string   `json:"category"`
	Cuisine     string   `json:"cuisine"`
	Description string   `json:"description"`
}

// NutritionData represents nutrition information
type NutritionData struct {
	DishID           string
	BaseServingGrams int
	Calories         float64
	Protein          float64
	Carbs            float64
	Fat              float64
	Fiber            float64
	Sodium           float64
}

func main() {
	log.Println("Starting nutrition database generation...")

	// Load dish taxonomy (relative to backend directory)
	taxonomy, err := loadTaxonomy("../../data/dish_taxonomy.json")
	if err != nil {
		log.Fatalf("Failed to load taxonomy: %v", err)
	}

	// Create SQLite database
	db, err := createDatabase("nutrition.db")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create schema
	if err := createSchema(db); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	// Insert dish master data
	if err := insertDishMaster(db, taxonomy); err != nil {
		log.Fatalf("Failed to insert dish master: %v", err)
	}

	// Insert nutrition data
	nutritionData := getIndianDishNutrition()
	if err := insertNutritionData(db, nutritionData); err != nil {
		log.Fatalf("Failed to insert nutrition data: %v", err)
	}

	// Insert model metadata
	if err := insertModelMetadata(db); err != nil {
		log.Fatalf("Failed to insert model metadata: %v", err)
	}

	log.Printf("✅ Database generated successfully!")
	log.Printf("   - %d dishes in master table", len(taxonomy.Dishes))
	log.Printf("   - %d nutrition records", len(nutritionData))
	log.Printf("   - Database file: nutrition.db")
}

func loadTaxonomy(path string) (*DishTaxonomy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var taxonomy DishTaxonomy
	if err := json.Unmarshal(data, &taxonomy); err != nil {
		return nil, err
	}

	return &taxonomy, nil
}

func createDatabase(path string) (*sql.DB, error) {
	// Remove existing database
	os.Remove(path)

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createSchema(db *sql.DB) error {
	schema := `
	-- Dish Master Table
	CREATE TABLE dish_master (
		dish_id TEXT PRIMARY KEY,
		display_name TEXT NOT NULL,
		category TEXT NOT NULL,
		cuisine TEXT NOT NULL,
		aliases TEXT,
		description TEXT,
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
		model_type TEXT PRIMARY KEY,
		version TEXT NOT NULL,
		checksum TEXT NOT NULL,
		size_bytes INTEGER NOT NULL,
		updated_at INTEGER DEFAULT (strftime('%s', 'now'))
	);

	-- Indexes
	CREATE INDEX idx_dish_category ON dish_master(category);
	CREATE INDEX idx_dish_cuisine ON dish_master(cuisine);
	CREATE INDEX idx_corrections_synced ON user_corrections(synced);
	CREATE INDEX idx_corrections_dish ON user_corrections(predicted_dish_id);
	`

	_, err := db.Exec(schema)
	return err
}

func insertDishMaster(db *sql.DB, taxonomy *DishTaxonomy) error {
	stmt, err := db.Prepare(`
		INSERT INTO dish_master (dish_id, display_name, category, cuisine, aliases, description)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for dishID, info := range taxonomy.Dishes {
		aliasesJSON, _ := json.Marshal(info.Aliases)
		_, err := stmt.Exec(
			dishID,
			info.DisplayName,
			info.Category,
			info.Cuisine,
			string(aliasesJSON),
			info.Description,
		)
		if err != nil {
			return fmt.Errorf("failed to insert %s: %w", dishID, err)
		}
	}

	return nil
}

func insertNutritionData(db *sql.DB, data []NutritionData) error {
	stmt, err := db.Prepare(`
		INSERT INTO dish_nutrition (
			dish_id, base_serving_grams, calories, protein, carbs, fat, fiber, sodium
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, nutrition := range data {
		_, err := stmt.Exec(
			nutrition.DishID,
			nutrition.BaseServingGrams,
			nutrition.Calories,
			nutrition.Protein,
			nutrition.Carbs,
			nutrition.Fat,
			nutrition.Fiber,
			nutrition.Sodium,
		)
		if err != nil {
			return fmt.Errorf("failed to insert nutrition for %s: %w", nutrition.DishID, err)
		}
	}

	return nil
}

func insertModelMetadata(db *sql.DB) error {
	stmt, err := db.Prepare(`
		INSERT INTO model_metadata (model_type, version, checksum, size_bytes)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Vision model metadata
	_, err = stmt.Exec("vision", "1.0.0", "sha256:placeholder", 6543210)
	if err != nil {
		return err
	}

	// Nutrition DB metadata
	_, err = stmt.Exec("nutrition_db", "1.0.0", "sha256:placeholder", 4567890)
	return err
}

func getIndianDishNutrition() []NutritionData {
	return []NutritionData{
		// Dal dishes
		{"DAL_MAKHANI", 250, 280, 14, 28, 12, 8, 380},
		{"DAL_TADKA", 250, 180, 12, 25, 4, 7, 320},
		{"DAL_FRY", 250, 190, 12, 26, 5, 7, 340},
		{"SAMBAR", 250, 150, 8, 22, 3, 6, 450},
		{"RAJMA_MASALA", 250, 220, 14, 35, 3, 10, 400},
		{"CHOLE_MASALA", 250, 240, 13, 38, 5, 11, 420},

		// Rice dishes
		{"RICE_BASMATI_PLAIN", 200, 260, 5.4, 56, 0.6, 0.8, 2},
		{"RICE_JEERA", 200, 290, 5.5, 56, 5, 0.8, 240},
		{"BIRYANI_CHICKEN", 350, 480, 22, 62, 14, 3, 680},
		{"BIRYANI_MUTTON", 350, 520, 24, 60, 18, 3, 720},
		{"BIRYANI_VEGETABLE", 300, 380, 8, 65, 9, 4, 520},
		{"PULAO_VEGETABLE", 250, 320, 6, 58, 7, 3, 380},

		// Bread
		{"ROTI_PLAIN", 50, 148, 3.5, 25.5, 3.5, 4, 180},
		{"NAAN_PLAIN", 80, 210, 7, 35, 4, 2.2, 320},
		{"NAAN_BUTTER", 90, 280, 7, 36, 12, 2.2, 380},
		{"NAAN_GARLIC", 90, 290, 7.5, 37, 12, 2.2, 420},
		{"PARATHA_PLAIN", 60, 200, 4, 28, 8, 3, 240},
		{"PARATHA_ALOO", 120, 350, 8, 50, 13, 4, 380},
		{"PURI", 40, 160, 3, 18, 9, 1.5, 180},

		// Chicken dishes
		{"CHICKEN_CURRY", 250, 320, 28, 8, 20, 2, 680},
		{"BUTTER_CHICKEN", 250, 380, 26, 10, 26, 2, 720},
		{"CHICKEN_TIKKA_MASALA", 250, 350, 28, 9, 22, 2, 650},
		{"CHICKEN_KORMA", 250, 360, 26, 12, 24, 3, 580},

		// Paneer dishes
		{"PALAK_PANEER", 200, 280, 15, 10, 20, 4, 480},
		{"PANEER_BUTTER_MASALA", 200, 350, 15, 12, 28, 3, 520},
		{"PANEER_TIKKA_MASALA", 200, 340, 16, 11, 26, 3, 500},
		{"MATAR_PANEER", 200, 280, 14, 15, 18, 5, 450},

		// Vegetable curries
		{"ALOO_GOBI", 200, 150, 3, 22, 6, 4, 380},
		{"BHINDI_MASALA", 200, 120, 4, 15, 5, 5, 320},
		{"BAINGAN_BHARTA", 200, 160, 3, 12, 11, 6, 360},

		// Breakfast
		{"IDLI", 150, 156, 4, 33, 0.5, 2, 280},
		{"DOSA_PLAIN", 150, 168, 4.5, 30, 3, 2, 320},
		{"DOSA_MASALA", 200, 250, 6, 42, 6, 3, 420},
		{"POHA", 200, 250, 6, 40, 7, 3, 380},
		{"UPMA", 200, 200, 5, 35, 5, 3, 420},

		// Snacks
		{"SAMOSA", 100, 262, 5, 33, 13, 3, 340},
		{"PAKORA_VEGETABLE", 100, 280, 6, 28, 16, 3, 380},
		{"VADA", 80, 180, 5, 22, 8, 3, 320},
		{"DHOKLA", 100, 160, 6, 28, 3, 3, 280},
		{"KACHORI", 80, 240, 5, 28, 12, 3, 380},

		// Eggs
		{"EGG_BOILED", 50, 78, 6.3, 0.6, 5.3, 0, 62},
		{"EGG_OMELET", 100, 154, 13, 1.1, 11, 0, 180},
		{"EGG_BHURJI", 150, 220, 15, 4, 16, 1, 380},

		// Desserts
		{"GULAB_JAMUN", 50, 175, 2.5, 28, 6.5, 0.5, 35},
		{"RASGULLA", 60, 140, 4, 25, 3, 0.3, 40},
		{"KHEER", 150, 190, 5, 32, 5, 1, 80},
		{"LADOO", 40, 160, 3, 22, 7, 1.5, 25},
	}
}
