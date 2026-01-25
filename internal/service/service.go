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
	Email            EmailServiceInterface // Brevo email service
	DB               *gorm.DB
}

// NewServices creates a new Services instance
func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              EMAIL SERVICE INITIALIZATION                    ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	
	// Initialize Brevo email service (REQUIRED)
	if cfg.BrevoAPIKey == "" {
		fmt.Println("❌ CRITICAL ERROR: BREVO_API_KEY not configured!")
		fmt.Println("   └─ Email service cannot be initialized without Brevo")
		fmt.Println("   └─ Please add BREVO_API_KEY to your environment variables")
		panic("BREVO_API_KEY environment variable is required")
	}
	
	fmt.Println("✅ Brevo API Key detected")
	fmt.Printf("   └─ Using Brevo Email Service\n\n")
	
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
