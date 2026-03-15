package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gymie-backend/internal/config"
	"github.com/yourusername/gymie-backend/internal/handlers"
	"github.com/yourusername/gymie-backend/internal/middleware"
	"github.com/yourusername/gymie-backend/internal/models"
	"github.com/yourusername/gymie-backend/internal/service"
)

// SetupRouter sets up the Gin router with all routes
func SetupRouter(services *service.Services, cfg *config.Config) *gin.Engine {
	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware(cfg))
	router.Use(middleware.ErrorHandler())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.NewSuccessResponse(
			"Server is running",
			gin.H{
				"status":  "healthy",
				"version": "1.0.0",
			},
		))
	})

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// Initialize handlers
		authHandler := handlers.NewAuthHandler(services.Auth, services.DB, services.Email)
		userHandler := handlers.NewUserHandler(services.User)
		workoutHandler := handlers.NewWorkoutHandler(services.Workout)
		nutritionHandler := handlers.NewNutritionHandler(services.Nutrition)
		progressHandler := handlers.NewProgressHandler(services.Progress)
		workoutPlanHandler := handlers.NewWorkoutPlanHandler(services.WorkoutPlan)
		offlineNutritionHandler := handlers.NewOfflineNutritionHandler(services.OfflineNutrition)

		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/google", authHandler.LoginWithGoogle)

			// Email verification routes (public)
			auth.GET("/verify-email", authHandler.VerifyEmail)
			auth.POST("/resend-verification", authHandler.ResendVerification)
			auth.GET("/verification-status/:email", authHandler.GetVerificationStatus)
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(services.Auth))
		{
			// Auth routes
			authProtected := protected.Group("/auth")
			{
				authProtected.GET("/me", authHandler.Me)
			}

			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", userHandler.GetProfile)
				users.PUT("/profile", userHandler.UpdateProfile)
				users.DELETE("/account", userHandler.DeleteAccount)
			}

			// Workout routes
			workouts := protected.Group("/workouts")
			{
				workouts.POST("", workoutHandler.Create)
				workouts.GET("", workoutHandler.List)
				workouts.GET("/stats", workoutHandler.GetStats)
				workouts.GET("/:id", workoutHandler.GetByID)
				workouts.PUT("/:id", workoutHandler.Update)
				workouts.DELETE("/:id", workoutHandler.Delete)
			}

			// Nutrition routes
			nutrition := protected.Group("/nutrition")
			{
				nutrition.POST("", nutritionHandler.Create)
				nutrition.GET("/date", nutritionHandler.GetByDate)
				nutrition.GET("/range", nutritionHandler.GetByDateRange)
				nutrition.GET("/stats", nutritionHandler.GetStats)
				nutrition.GET("/:id", nutritionHandler.GetByID)
				nutrition.PUT("/:id", nutritionHandler.Update)
				nutrition.DELETE("/:id", nutritionHandler.Delete)
				nutrition.POST("/:id/meals", nutritionHandler.AddMeal)
				nutrition.DELETE("/:id/meals/:meal_id", nutritionHandler.DeleteMeal)
			}

			// Progress routes
			progress := protected.Group("/progress")
			{
				// Progress photos
				progress.POST("/photos", progressHandler.CreateProgressPhoto)
				progress.GET("/photos", progressHandler.GetProgressPhotos)
				progress.PUT("/photos/:id", progressHandler.UpdateProgressPhoto)
				progress.DELETE("/photos/:id", progressHandler.DeleteProgressPhoto)

				// Weight tracking
				progress.POST("/weight", progressHandler.CreateWeightEntry)
				progress.GET("/weight", progressHandler.GetWeightEntries)
				progress.GET("/weight/stats", progressHandler.GetWeightProgress)
				progress.PUT("/weight/:id", progressHandler.UpdateWeightEntry)
				progress.DELETE("/weight/:id", progressHandler.DeleteWeightEntry)

				// Upload URL generation
				progress.GET("/upload-url", progressHandler.GenerateUploadURL)
			}

			// Workout Plan routes
			workoutPlans := protected.Group("/workout-plans")
			{
				workoutPlans.POST("", workoutPlanHandler.CreatePlan)
				workoutPlans.GET("", workoutPlanHandler.ListPlans)
				workoutPlans.GET("/active", workoutPlanHandler.GetActivePlan)
				workoutPlans.GET("/:id", workoutPlanHandler.GetPlan)
				workoutPlans.PUT("/:id", workoutPlanHandler.UpdatePlan)
				workoutPlans.DELETE("/:id", workoutPlanHandler.DeletePlan)
				workoutPlans.PUT("/:id/active", workoutPlanHandler.SetActivePlan)
				workoutPlans.POST("/:id/clone", workoutPlanHandler.ClonePlan)
				workoutPlans.PUT("/:id/recurrence", workoutPlanHandler.SetRecurrence)
			}

			// Scheduled Workout routes
			scheduledWorkouts := protected.Group("/scheduled-workouts")
			{
				scheduledWorkouts.GET("", workoutPlanHandler.ListScheduled)
				scheduledWorkouts.GET("/today", workoutPlanHandler.GetTodaysWorkout)
				scheduledWorkouts.GET("/:id", workoutPlanHandler.GetScheduled)
				scheduledWorkouts.PUT("/:id/status", workoutPlanHandler.UpdateScheduledStatus)
				scheduledWorkouts.DELETE("/:id", workoutPlanHandler.DeleteScheduled)
				scheduledWorkouts.POST("/generate", workoutPlanHandler.GenerateScheduled)
				scheduledWorkouts.DELETE("/plan/:plan_id", workoutPlanHandler.ClearPlanSchedule)
			}

			// Offline Nutrition Sync routes
			sync := protected.Group("/sync")
			{
				// Correction sync
				sync.POST("/corrections", offlineNutritionHandler.SyncCorrections)
				sync.GET("/corrections/stats", offlineNutritionHandler.GetCorrectionStats)

				// Model version management
				sync.GET("/model-versions", offlineNutritionHandler.GetModelVersions)
				sync.GET("/nutrition-db", offlineNutritionHandler.DownloadNutritionDB)

				// Dish search and lookup
				sync.GET("/dishes/search", offlineNutritionHandler.SearchDishes)
				sync.GET("/dishes/:dish_id", offlineNutritionHandler.GetDishByID)
			}
		}
	}

	return router
}
