
package repository

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Repositories holds all repository instances
type Repositories struct {
	User       UserRepository
	Workout    WorkoutRepository
	Nutrition  NutritionRepository
	Progress   ProgressRepository
}

// NewRepositories creates a new Repositories instance
func NewRepositories(db *gorm.DB, redis *redis.Client) *Repositories {
	return &Repositories{
		User:      NewUserRepository(db, redis),
		Workout:   NewWorkoutRepository(db, redis),
		Nutrition: NewNutritionRepository(db, redis),
		Progress:  NewProgressRepository(db, redis),
	}
}
