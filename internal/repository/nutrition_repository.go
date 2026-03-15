package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yourusername/gymie-backend/internal/models"
	"gorm.io/gorm"
)

// NutritionRepository defines the interface for nutrition data operations
type NutritionRepository interface {
	Create(ctx context.Context, day *models.NutritionDay) error
	GetByID(ctx context.Context, id uint, userID uint) (*models.NutritionDay, error)
	GetByDate(ctx context.Context, userID uint, date time.Time) (*models.NutritionDay, error)
	GetByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]models.NutritionDay, error)
	Update(ctx context.Context, day *models.NutritionDay) error
	Delete(ctx context.Context, id uint, userID uint) error

	// Meal operations
	CreateMeal(ctx context.Context, meal *models.Meal) error
	UpdateMeal(ctx context.Context, meal *models.Meal) error
	DeleteMeal(ctx context.Context, id uint, nutritionDayID uint) error

	// Food operations
	CreateFood(ctx context.Context, food *models.Food) error
	UpdateFood(ctx context.Context, food *models.Food) error
	DeleteFood(ctx context.Context, id uint, mealID uint) error

	// Stats
	GetNutritionStats(ctx context.Context, userID uint) (*models.NutritionStatsResponse, error)
}

type nutritionRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewNutritionRepository creates a new nutrition repository
func NewNutritionRepository(db *gorm.DB, redis *redis.Client) NutritionRepository {
	return &nutritionRepository{
		db:    db,
		redis: redis,
	}
}

// Create creates a new nutrition day with meals and foods
func (r *nutritionRepository) Create(ctx context.Context, day *models.NutritionDay) error {
	return r.db.WithContext(ctx).Create(day).Error
}

// GetByID retrieves a nutrition day by ID with caching
func (r *nutritionRepository) GetByID(ctx context.Context, id uint, userID uint) (*models.NutritionDay, error) {
	cacheKey := fmt.Sprintf("nutrition:%d", id)

	// Try cache first
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var day models.NutritionDay
		if err := json.Unmarshal([]byte(cached), &day); err == nil {
			if day.UserID == userID {
				day.CalculateTotals()
				return &day, nil
			}
		}
	}

	// Fetch from database
	var day models.NutritionDay
	if err := r.db.WithContext(ctx).
		Preload("Meals.Foods").
		Where("id = ? AND user_id = ?", id, userID).
		First(&day).Error; err != nil {
		return nil, err
	}

	// Calculate totals
	day.CalculateTotals()

	// Cache the result
	if data, err := json.Marshal(day); err == nil {
		r.redis.Set(ctx, cacheKey, data, 15*time.Minute)
	}

	return &day, nil
}

// GetByDate retrieves a nutrition day for a specific date
func (r *nutritionRepository) GetByDate(ctx context.Context, userID uint, date time.Time) (*models.NutritionDay, error) {
	// Normalize date to start of day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	var day models.NutritionDay
	if err := r.db.WithContext(ctx).
		Preload("Meals.Foods").
		Where("user_id = ? AND DATE(date) = DATE(?)", userID, startOfDay).
		First(&day).Error; err != nil {
		return nil, err
	}

	// Calculate totals
	day.CalculateTotals()

	return &day, nil
}

// GetByDateRange retrieves nutrition days within a date range
func (r *nutritionRepository) GetByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]models.NutritionDay, error) {
	var days []models.NutritionDay

	if err := r.db.WithContext(ctx).
		Preload("Meals.Foods").
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("date DESC").
		Find(&days).Error; err != nil {
		return nil, err
	}

	// Calculate totals for each day
	for i := range days {
		days[i].CalculateTotals()
	}

	return days, nil
}

// Update updates a nutrition day and invalidates cache
func (r *nutritionRepository) Update(ctx context.Context, day *models.NutritionDay) error {
	if err := r.db.WithContext(ctx).Save(day).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("nutrition:%d", day.ID)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// Delete soft deletes a nutrition day and invalidates cache
func (r *nutritionRepository) Delete(ctx context.Context, id uint, userID uint) error {
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.NutritionDay{}).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("nutrition:%d", id)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// CreateMeal creates a new meal
func (r *nutritionRepository) CreateMeal(ctx context.Context, meal *models.Meal) error {
	if err := r.db.WithContext(ctx).Create(meal).Error; err != nil {
		return err
	}

	// Invalidate nutrition day cache
	cacheKey := fmt.Sprintf("nutrition:%d", meal.NutritionDayID)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// UpdateMeal updates a meal
func (r *nutritionRepository) UpdateMeal(ctx context.Context, meal *models.Meal) error {
	if err := r.db.WithContext(ctx).Save(meal).Error; err != nil {
		return err
	}

	// Invalidate nutrition day cache
	cacheKey := fmt.Sprintf("nutrition:%d", meal.NutritionDayID)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// DeleteMeal deletes a meal
func (r *nutritionRepository) DeleteMeal(ctx context.Context, id uint, nutritionDayID uint) error {
	if err := r.db.WithContext(ctx).
		Where("id = ? AND nutrition_day_id = ?", id, nutritionDayID).
		Delete(&models.Meal{}).Error; err != nil {
		return err
	}

	// Invalidate nutrition day cache
	cacheKey := fmt.Sprintf("nutrition:%d", nutritionDayID)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// CreateFood creates a new food item
func (r *nutritionRepository) CreateFood(ctx context.Context, food *models.Food) error {
	return r.db.WithContext(ctx).Create(food).Error
}

// UpdateFood updates a food item
func (r *nutritionRepository) UpdateFood(ctx context.Context, food *models.Food) error {
	return r.db.WithContext(ctx).Save(food).Error
}

// DeleteFood deletes a food item
func (r *nutritionRepository) DeleteFood(ctx context.Context, id uint, mealID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND meal_id = ?", id, mealID).
		Delete(&models.Food{}).Error
}

// GetNutritionStats retrieves nutrition statistics for a user
func (r *nutritionRepository) GetNutritionStats(ctx context.Context, userID uint) (*models.NutritionStatsResponse, error) {
	var stats models.NutritionStatsResponse

	// Get total days
	if err := r.db.WithContext(ctx).
		Model(&models.NutritionDay{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalDays).Error; err != nil {
		return nil, err
	}

	// Get last logged date
	var lastDay models.NutritionDay
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("date DESC").
		First(&lastDay).Error; err == nil {
		stats.LastLoggedDate = &lastDay.Date
	}

	// Calculate averages (this is simplified - in production you might want to aggregate in DB)
	var days []models.NutritionDay
	if err := r.db.WithContext(ctx).
		Preload("Meals.Foods").
		Where("user_id = ?", userID).
		Find(&days).Error; err == nil && len(days) > 0 {

		var totalCal, totalPro, totalCarb, totalFat float64
		for _, day := range days {
			day.CalculateTotals()
			totalCal += float64(day.TotalCalories)
			totalPro += float64(day.TotalProtein)
			totalCarb += float64(day.TotalCarbs)
			totalFat += float64(day.TotalFat)
		}

		count := float64(len(days))
		stats.AverageCalories = totalCal / count
		stats.AverageProtein = totalPro / count
		stats.AverageCarbs = totalCarb / count
		stats.AverageFat = totalFat / count
	}

	return &stats, nil
}
