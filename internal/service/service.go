package service

import (
	"fmt"
	
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
	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              EMAIL SERVICE INITIALIZATION                    ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	
	// Initialize email service based on configuration
	var emailServiceInterface EmailServiceInterface
	var legacyEmailService *services.EmailService
	
	if cfg.SendGridAPIKey != "" {
		fmt.Println("✅ SendGrid API Key detected")
		fmt.Printf("   └─ Using SendGrid Email Service (Production)\n\n")
		// Use SendGrid if API key is provided (preferred)
		emailServiceInterface = services.NewSendGridEmailService(cfg)
		// Also initialize legacy service for backward compatibility
		legacyEmailService = services.NewEmailService()
	} else {
		fmt.Println("⚠️  No SendGrid API Key found")
		fmt.Printf("   └─ Falling back to Gmail SMTP Service\n\n")
		// Fall back to SMTP (Gmail) if no SendGrid key
		legacyEmailService = services.NewEmailService()
		emailServiceInterface = legacyEmailService
	}

	return &Services{
		Auth:             NewAuthService(repos.User, cfg, emailServiceInterface),
		User:             NewUserService(repos.User),
		Workout:          NewWorkoutService(repos.Workout),
		Nutrition:        NewNutritionService(repos.Nutrition),
		Progress:         NewProgressService(repos.Progress, cfg),
		WorkoutPlan:      NewWorkoutPlanService(repos.WorkoutPlan),
		OfflineNutrition: NewOfflineNutritionService(repos.DB),
		Email:            legacyEmailService,
		DB:               repos.DB,
	}
}
