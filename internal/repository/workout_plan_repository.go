
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

// WorkoutPlanRepository defines the interface for workout plan operations
type WorkoutPlanRepository interface {
	// Plan operations
	Create(ctx context.Context, plan *models.WorkoutPlan) error
	GetByID(ctx context.Context, id uint, userID uint) (*models.WorkoutPlan, error)
	GetByUser(ctx context.Context, userID uint) ([]models.WorkoutPlan, error)
	GetActive(ctx context.Context, userID uint) (*models.WorkoutPlan, error)
	Update(ctx context.Context, plan *models.WorkoutPlan) error
	Delete(ctx context.Context, id uint, userID uint) error
	SetActive(ctx context.Context, id uint, userID uint) error
	
	// Day operations
	CreateDay(ctx context.Context, day *models.WorkoutPlanDay) error
	CreateDays(ctx context.Context, days []models.WorkoutPlanDay) error
	UpdateDay(ctx context.Context, day *models.WorkoutPlanDay) error
	DeleteDay(ctx context.Context, id uint, planID uint) error
	DeleteDays(ctx context.Context, planID uint) error
	
	// Scheduled workout operations
	CreateScheduled(ctx context.Context, scheduled *models.ScheduledWorkout) error
	GetScheduledByID(ctx context.Context, id uint, userID uint) (*models.ScheduledWorkout, error)
	GetScheduledByUser(ctx context.Context, userID uint) ([]models.ScheduledWorkout, error)
	GetScheduledByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]models.ScheduledWorkout, error)
	GetScheduledByDate(ctx context.Context, userID uint, date time.Time) ([]models.ScheduledWorkout, error)
	GetTodaysScheduled(ctx context.Context, userID uint) (*models.ScheduledWorkout, error)
	UpdateScheduled(ctx context.Context, scheduled *models.ScheduledWorkout) error
	DeleteScheduled(ctx context.Context, id uint, userID uint) error
	DeleteScheduledByPlan(ctx context.Context, planID uint, userID uint) error
	BulkCreateScheduled(ctx context.Context, scheduled []models.ScheduledWorkout) error
}

type workoutPlanRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewWorkoutPlanRepository creates a new workout plan repository
func NewWorkoutPlanRepository(db *gorm.DB, redis *redis.Client) WorkoutPlanRepository {
	return &workoutPlanRepository{
		db:    db,
		redis: redis,
	}
}

// Create creates a new workout plan
func (r *workoutPlanRepository) Create(ctx context.Context, plan *models.WorkoutPlan) error {
	return r.db.WithContext(ctx).Create(plan).Error
}

// GetByID retrieves a workout plan by ID
func (r *workoutPlanRepository) GetByID(ctx context.Context, id uint, userID uint) (*models.WorkoutPlan, error) {
	cacheKey := fmt.Sprintf("plan:%d", id)
	
	// Try cache first
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var plan models.WorkoutPlan
		if err := json.Unmarshal([]byte(cached), &plan); err == nil {
			if plan.UserID == userID {
				return &plan, nil
			}
		}
	}

	// Fetch from database
	var plan models.WorkoutPlan
	if err := r.db.WithContext(ctx).
		Preload("Days").
		Where("id = ? AND user_id = ?", id, userID).
		First(&plan).Error; err != nil {
		return nil, err
	}

	// Cache the result
	if data, err := json.Marshal(plan); err == nil {
		r.redis.Set(ctx, cacheKey, data, 10*time.Minute)
	}

	return &plan, nil
}

// GetByUser retrieves all workout plans for a user
func (r *workoutPlanRepository) GetByUser(ctx context.Context, userID uint) ([]models.WorkoutPlan, error) {
	var plans []models.WorkoutPlan
	if err := r.db.WithContext(ctx).
		Preload("Days").
		Where("user_id = ?", userID).
		Order("is_active DESC, updated_at DESC").
		Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

// GetActive retrieves the active workout plan for a user
func (r *workoutPlanRepository) GetActive(ctx context.Context, userID uint) (*models.WorkoutPlan, error) {
	var plan models.WorkoutPlan
	if err := r.db.WithContext(ctx).
		Preload("Days").
		Where("user_id = ? AND is_active = ?", userID, true).
		First(&plan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &plan, nil
}

// Update updates a workout plan
func (r *workoutPlanRepository) Update(ctx context.Context, plan *models.WorkoutPlan) error {
	// Invalidate cache
	cacheKey := fmt.Sprintf("plan:%d", plan.ID)
	r.redis.Del(ctx, cacheKey)
	
	return r.db.WithContext(ctx).Save(plan).Error
}

// Delete soft deletes a workout plan
func (r *workoutPlanRepository) Delete(ctx context.Context, id uint, userID uint) error {
	// Invalidate cache
	cacheKey := fmt.Sprintf("plan:%d", id)
	r.redis.Del(ctx, cacheKey)
	
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.WorkoutPlan{}).Error
}

// SetActive sets a plan as active and deactivates others
func (r *workoutPlanRepository) SetActive(ctx context.Context, id uint, userID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Deactivate all plans for the user
		if err := tx.Model(&models.WorkoutPlan{}).
			Where("user_id = ?", userID).
			Update("is_active", false).Error; err != nil {
			return err
		}
		
		// Activate the selected plan
		if err := tx.Model(&models.WorkoutPlan{}).
			Where("id = ? AND user_id = ?", id, userID).
			Update("is_active", true).Error; err != nil {
			return err
		}
		
		return nil
	})
}

// CreateDay creates a new day in a workout plan
func (r *workoutPlanRepository) CreateDay(ctx context.Context, day *models.WorkoutPlanDay) error {
	// Invalidate parent plan cache
	cacheKey := fmt.Sprintf("plan:%d", day.PlanID)
	r.redis.Del(ctx, cacheKey)
	
	return r.db.WithContext(ctx).Create(day).Error
}

// UpdateDay updates a day in a workout plan
func (r *workoutPlanRepository) UpdateDay(ctx context.Context, day *models.WorkoutPlanDay) error {
	// Invalidate parent plan cache
	cacheKey := fmt.Sprintf("plan:%d", day.PlanID)
	r.redis.Del(ctx, cacheKey)
	
	return r.db.WithContext(ctx).Save(day).Error
}

// DeleteDay deletes a day from a workout plan
func (r *workoutPlanRepository) DeleteDay(ctx context.Context, id uint, planID uint) error {
	// Invalidate parent plan cache
	cacheKey := fmt.Sprintf("plan:%d", planID)
	r.redis.Del(ctx, cacheKey)
	
	return r.db.WithContext(ctx).
		Where("id = ? AND plan_id = ?", id, planID).
		Delete(&models.WorkoutPlanDay{}).Error
}

// CreateDays creates multiple days in a workout plan
func (r *workoutPlanRepository) CreateDays(ctx context.Context, days []models.WorkoutPlanDay) error {
	if len(days) == 0 {
		return nil
	}
	
	// Invalidate parent plan cache
	if len(days) > 0 {
		cacheKey := fmt.Sprintf("plan:%d", days[0].PlanID)
		r.redis.Del(ctx, cacheKey)
	}
	
	return r.db.WithContext(ctx).Create(&days).Error
}

// DeleteDays deletes all days from a workout plan
func (r *workoutPlanRepository) DeleteDays(ctx context.Context, planID uint) error {
	// Invalidate parent plan cache
	cacheKey := fmt.Sprintf("plan:%d", planID)
	r.redis.Del(ctx, cacheKey)
	
	return r.db.WithContext(ctx).
		Where("plan_id = ?", planID).
		Delete(&models.WorkoutPlanDay{}).Error
}

// CreateScheduled creates a scheduled workout
func (r *workoutPlanRepository) CreateScheduled(ctx context.Context, scheduled *models.ScheduledWorkout) error {
	return r.db.WithContext(ctx).Create(scheduled).Error
}

// GetScheduledByID retrieves a scheduled workout by ID
func (r *workoutPlanRepository) GetScheduledByID(ctx context.Context, id uint, userID uint) (*models.ScheduledWorkout, error) {
	var scheduled models.ScheduledWorkout
	if err := r.db.WithContext(ctx).
		Preload("Plan").
		Preload("PlanDay").
		Preload("Workout.Exercises.Sets").
		Where("id = ? AND user_id = ?", id, userID).
		First(&scheduled).Error; err != nil {
		return nil, err
	}
	return &scheduled, nil
}

// GetScheduledByUser retrieves all scheduled workouts for a user
func (r *workoutPlanRepository) GetScheduledByUser(ctx context.Context, userID uint) ([]models.ScheduledWorkout, error) {
	var scheduled []models.ScheduledWorkout
	if err := r.db.WithContext(ctx).
		Preload("Plan").
		Preload("PlanDay").
		Where("user_id = ?", userID).
		Order("date DESC").
		Find(&scheduled).Error; err != nil {
		return nil, err
	}
	return scheduled, nil
}

// GetScheduledByDateRange retrieves scheduled workouts within a date range
func (r *workoutPlanRepository) GetScheduledByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]models.ScheduledWorkout, error) {
	var scheduled []models.ScheduledWorkout
	if err := r.db.WithContext(ctx).
		Preload("Plan").
		Preload("PlanDay").
		Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).
		Order("date ASC").
		Find(&scheduled).Error; err != nil {
		return nil, err
	}
	return scheduled, nil
}

// GetScheduledByDate retrieves scheduled workouts for a specific date
func (r *workoutPlanRepository) GetScheduledByDate(ctx context.Context, userID uint, date time.Time) ([]models.ScheduledWorkout, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	var scheduled []models.ScheduledWorkout
	if err := r.db.WithContext(ctx).
		Preload("Plan").
		Preload("PlanDay").
		Where("user_id = ? AND date >= ? AND date < ?", userID, startOfDay, endOfDay).
		Order("date ASC").
		Find(&scheduled).Error; err != nil {
		return nil, err
	}
	return scheduled, nil
}

// GetTodaysScheduled retrieves today's scheduled workout
func (r *workoutPlanRepository) GetTodaysScheduled(ctx context.Context, userID uint) (*models.ScheduledWorkout, error) {
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	var scheduled models.ScheduledWorkout
	if err := r.db.WithContext(ctx).
		Preload("Plan").
		Preload("PlanDay").
		Where("user_id = ? AND date >= ? AND date < ? AND status = ?", userID, startOfDay, endOfDay, "scheduled").
		First(&scheduled).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &scheduled, nil
}

// UpdateScheduled updates a scheduled workout
func (r *workoutPlanRepository) UpdateScheduled(ctx context.Context, scheduled *models.ScheduledWorkout) error {
	return r.db.WithContext(ctx).Save(scheduled).Error
}

// DeleteScheduled deletes a scheduled workout
func (r *workoutPlanRepository) DeleteScheduled(ctx context.Context, id uint, userID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.ScheduledWorkout{}).Error
}

// DeleteScheduledByPlan deletes all scheduled workouts for a plan
func (r *workoutPlanRepository) DeleteScheduledByPlan(ctx context.Context, planID uint, userID uint) error {
	return r.db.WithContext(ctx).
		Where("plan_id = ? AND user_id = ?", planID, userID).
		Delete(&models.ScheduledWorkout{}).Error
}

// BulkCreateScheduled creates multiple scheduled workouts
func (r *workoutPlanRepository) BulkCreateScheduled(ctx context.Context, scheduled []models.ScheduledWorkout) error {
	if len(scheduled) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&scheduled).Error
}
