package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/yourusername/gymie-backend/internal/models"
)

// InitDB initializes and returns a database connection
func InitDB(cfg *Config) (*gorm.DB, error) {
	var dsn string

	// Use DATABASE_URL if provided, otherwise construct from individual params
	if cfg.DatabaseURL != "" {
		dsn = cfg.DatabaseURL
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
		)
	}

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if cfg.IsProduction() {
		gormConfig.Logger = logger.Default.LogMode(logger.Error)
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")

	// Auto-migrate database schema
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	return db, nil
}

// autoMigrate automatically migrates database schema to match models
func autoMigrate(db *gorm.DB) error {
	log.Println("Running auto-migration...")

	// List all models that need to be migrated
	models := []interface{}{
		// Core user models
		&models.User{},
		&models.UserProfile{},
		
		// Workout models
		&models.Workout{},
		&models.Exercise{},
		&models.WorkoutSet{},
		
		// Nutrition models
		&models.NutritionDay{},
		&models.Meal{},
		&models.Food{},
		
		// Progress tracking models
		&models.ProgressPhoto{},
		&models.WeightEntry{},
		
		// Workout planning models
		&models.WorkoutPlan{},
		&models.WorkoutPlanDay{},
		&models.ScheduledWorkout{},
		
		// Offline nutrition models
		&models.DishMaster{},
		&models.DishNutritionMaster{},
		&models.UserCorrection{},
		&models.ModelVersion{},
	}

	// Run auto-migration for all models
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("auto-migration failed: %w", err)
	}

	log.Println("✅ Auto-migration completed successfully")
	log.Println("✅ All tables synced including: users (with email_verified), workout_plans, scheduled_workouts, offline_nutrition tables")
	return nil
}
