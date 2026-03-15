package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/yourusername/gymie-backend/internal/models"
	"github.com/yourusername/gymie-backend/internal/repository"
)

// NutritionService defines the interface for nutrition operations
type NutritionService interface {
	Create(ctx context.Context, userID uint, req *models.NutritionDayCreateRequest) (*models.NutritionDay, error)
	GetByID(ctx context.Context, id uint, userID uint) (*models.NutritionDay, error)
	GetByDate(ctx context.Context, userID uint, date time.Time) (*models.NutritionDay, error)
	GetByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]models.NutritionDay, error)
	Update(ctx context.Context, id uint, userID uint, req *models.NutritionDayUpdateRequest) (*models.NutritionDay, error)
	Delete(ctx context.Context, id uint, userID uint) error
	GetStats(ctx context.Context, userID uint) (*models.NutritionStatsResponse, error)
	AddMeal(ctx context.Context, nutritionDayID uint, userID uint, req *models.AddMealRequest) (*models.Meal, error)
	DeleteMeal(ctx context.Context, nutritionDayID uint, mealID uint, userID uint) error
}

type nutritionService struct {
	nutritionRepo repository.NutritionRepository
}

// NewNutritionService creates a new nutrition service
func NewNutritionService(nutritionRepo repository.NutritionRepository) NutritionService {
	return &nutritionService{
		nutritionRepo: nutritionRepo,
	}
}

// Create creates a new nutrition day
func (s *nutritionService) Create(ctx context.Context, userID uint, req *models.NutritionDayCreateRequest) (*models.NutritionDay, error) {
	day := &models.NutritionDay{
		UserID: userID,
		Date:   req.Date,
		Notes:  req.Notes,
	}

	// Add meals
	if req.Meals != nil {
		for _, mealReq := range req.Meals {
			meal := models.Meal{
				Name: mealReq.Name,
				Time: mealReq.Time,
			}

			// Add foods
			if mealReq.Foods != nil {
				for _, foodReq := range mealReq.Foods {
					food := models.Food{
						Name:     foodReq.Name,
						Calories: foodReq.Calories,
						Protein:  foodReq.Protein,
						Carbs:    foodReq.Carbs,
						Fat:      foodReq.Fat,
						Quantity: foodReq.Quantity,
						Unit:     foodReq.Unit,
					}
					meal.Foods = append(meal.Foods, food)
				}
			}

			day.Meals = append(day.Meals, meal)
		}
	}

	if err := s.nutritionRepo.Create(ctx, day); err != nil {
		return nil, fmt.Errorf("failed to create nutrition day: %w", err)
	}

	// Calculate totals
	day.CalculateTotals()

	return day, nil
}

// GetByID retrieves a nutrition day by ID
func (s *nutritionService) GetByID(ctx context.Context, id uint, userID uint) (*models.NutritionDay, error) {
	day, err := s.nutritionRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nutrition day: %w", err)
	}
	return day, nil
}

// GetByDate retrieves a nutrition day for a specific date
func (s *nutritionService) GetByDate(ctx context.Context, userID uint, date time.Time) (*models.NutritionDay, error) {
	day, err := s.nutritionRepo.GetByDate(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get nutrition day by date: %w", err)
	}
	return day, nil
}

// GetByDateRange retrieves nutrition days within a date range
func (s *nutritionService) GetByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]models.NutritionDay, error) {
	days, err := s.nutritionRepo.GetByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get nutrition days by date range: %w", err)
	}
	return days, nil
}

// Update updates a nutrition day
func (s *nutritionService) Update(ctx context.Context, id uint, userID uint, req *models.NutritionDayUpdateRequest) (*models.NutritionDay, error) {
	// Get existing nutrition day
	day, err := s.nutritionRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nutrition day: %w", err)
	}

	// Update fields
	if req.Date != nil {
		day.Date = *req.Date
	}
	if req.Notes != nil {
		day.Notes = *req.Notes
	}

	if err := s.nutritionRepo.Update(ctx, day); err != nil {
		return nil, fmt.Errorf("failed to update nutrition day: %w", err)
	}

	// Calculate totals
	day.CalculateTotals()

	return day, nil
}

// Delete deletes a nutrition day
func (s *nutritionService) Delete(ctx context.Context, id uint, userID uint) error {
	if err := s.nutritionRepo.Delete(ctx, id, userID); err != nil {
		return fmt.Errorf("failed to delete nutrition day: %w", err)
	}
	return nil
}

// GetStats retrieves nutrition statistics
func (s *nutritionService) GetStats(ctx context.Context, userID uint) (*models.NutritionStatsResponse, error) {
	stats, err := s.nutritionRepo.GetNutritionStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nutrition stats: %w", err)
	}
	return stats, nil
}

// AddMeal adds a meal item (with one food entry) to a nutrition day
func (s *nutritionService) AddMeal(ctx context.Context, nutritionDayID uint, userID uint, req *models.AddMealRequest) (*models.Meal, error) {
	// Verify nutrition day belongs to user
	if _, err := s.nutritionRepo.GetByID(ctx, nutritionDayID, userID); err != nil {
		return nil, fmt.Errorf("nutrition day not found: %w", err)
	}

	// Capitalize meal type: "breakfast" → "Breakfast"
	mealTypeName := strings.ToUpper(req.MealType[:1]) + strings.ToLower(req.MealType[1:])

	meal := &models.Meal{
		NutritionDayID: nutritionDayID,
		Name:           mealTypeName,
		Foods: []models.Food{
			{
				Name:     req.Name,
				Calories: req.Calories,
				Protein:  req.Protein,
				Carbs:    req.Carbs,
				Fat:      req.Fat,
				Quantity: 1,
			},
		},
	}

	if err := s.nutritionRepo.CreateMeal(ctx, meal); err != nil {
		return nil, fmt.Errorf("failed to add meal: %w", err)
	}

	meal.CalculateTotals()
	return meal, nil
}

// DeleteMeal deletes a meal from a nutrition day
func (s *nutritionService) DeleteMeal(ctx context.Context, nutritionDayID uint, mealID uint, userID uint) error {
	// Verify nutrition day belongs to user
	if _, err := s.nutritionRepo.GetByID(ctx, nutritionDayID, userID); err != nil {
		return fmt.Errorf("nutrition day not found: %w", err)
	}

	if err := s.nutritionRepo.DeleteMeal(ctx, mealID, nutritionDayID); err != nil {
		return fmt.Errorf("failed to delete meal: %w", err)
	}
	return nil
}
