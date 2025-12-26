
package service

import (
	"github.com/yourusername/gymie-backend/internal/config"
	"github.com/yourusername/gymie-backend/internal/repository"
)

// Services holds all service instances
type Services struct {
	Auth        AuthService
	User        UserService
	Workout     WorkoutService
	Nutrition   NutritionService
	Progress    ProgressService
	WorkoutPlan WorkoutPlanService
}

// NewServices creates a new Services instance
func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	return &Services{
		Auth:        NewAuthService(repos.User, cfg),
		User:        NewUserService(repos.User),
		Workout:     NewWorkoutService(repos.Workout),
		Nutrition:   NewNutritionService(repos.Nutrition),
		Progress:    NewProgressService(repos.Progress, cfg),
		WorkoutPlan: NewWorkoutPlanService(repos.WorkoutPlan),
	}
}
