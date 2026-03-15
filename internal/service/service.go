package service

import (
	"log"

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
	Email            EmailServiceInterface
	DB               *gorm.DB
}

// NewServices creates a new Services instance
func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	if cfg.BrevoAPIKey == "" {
		log.Fatal("BREVO_API_KEY environment variable is required")
	}

	// Use Brevo as the email service
	emailService := services.NewBrevoEmailService(cfg)

	return &Services{
		Auth:             NewAuthService(repos.User, cfg, emailService),
		User:             NewUserService(repos.User),
		Workout:          NewWorkoutService(repos.Workout),
		Nutrition:        NewNutritionService(repos.Nutrition),
		Progress:         NewProgressService(repos.Progress, cfg),
		WorkoutPlan:      NewWorkoutPlanService(repos.WorkoutPlan),
		OfflineNutrition: NewOfflineNutritionService(repos.DB),
		Email:            emailService, // Brevo email service
		DB:               repos.DB,
	}
}
