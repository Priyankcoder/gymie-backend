
package main

import (
	"log"

	"github.com/yourusername/gymie-backend/internal/config"
	"github.com/yourusername/gymie-backend/internal/models"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("Running database migrations...")

	// Auto migrate all models
	err = db.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.Workout{},
		&models.Exercise{},
		&models.WorkoutSet{},
		&models.NutritionDay{},
		&models.Meal{},
		&models.Food{},
		&models.ProgressPhoto{},
		&models.WeightEntry{},
	)

	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully!")
}
