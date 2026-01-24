
package service

import (
	"github.com/yourusername/gymie-backend/internal/config"
	"github.com/yourusername/gymie-backend/internal/repository"
	"github.com/yourusername/gymie-backend/internal/services"
	"gorm.io/gorm"
)

// Services holds all service instances
type Services struct {
	Auth             AuthService
	User             UserService
	Workout          WorkoutService
	Nutrition        NutritionService
	Progress         ProgressService
	WorkoutPlan      WorkoutPlanService
	OfflineNutrition *OfflineNutritionService
	Email            *services.EmailService
	DB               *gorm.DB
}

// NewServices creates a new Services instance
func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	// Initialize email service
	emailService := services.NewEmailService()
	
	return &Services{
		Auth:             NewAuthService(repos.User, cfg, emailService),
		User:             NewUserService(repos.User),
		Workout:          NewWorkoutService(repos.Workout),
		Nutrition:        NewNutritionService(repos.Nutrition),
		Progress:         NewProgressService(repos.Progress, cfg),
		WorkoutPlan:      NewWorkoutPlanService(repos.WorkoutPlan),
		OfflineNutrition: NewOfflineNutritionService(repos.DB),
		Email:            emailService,
		DB:               repos.DB,
	}
}
