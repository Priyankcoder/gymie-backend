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

// WorkoutRepository defines the interface for workout data operations
type WorkoutRepository interface {
	Create(ctx context.Context, workout *models.Workout) error
	GetByID(ctx context.Context, id uint, userID uint) (*models.Workout, error)
	GetByUser(ctx context.Context, userID uint, query *models.ListQuery) ([]models.Workout, int64, error)
	GetByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]models.Workout, error)
	Update(ctx context.Context, workout *models.Workout) error
	Delete(ctx context.Context, id uint, userID uint) error

	// Exercise operations
	CreateExercise(ctx context.Context, exercise *models.Exercise) error
	UpdateExercise(ctx context.Context, exercise *models.Exercise) error
	DeleteExercise(ctx context.Context, id uint, workoutID uint) error

	// Set operations
	CreateSet(ctx context.Context, set *models.WorkoutSet) error
	UpdateSet(ctx context.Context, set *models.WorkoutSet) error
	DeleteSet(ctx context.Context, id uint, exerciseID uint) error

	// Stats
	GetWorkoutStats(ctx context.Context, userID uint) (*models.WorkoutStatsResponse, error)
}

type workoutRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewWorkoutRepository creates a new workout repository
func NewWorkoutRepository(db *gorm.DB, redis *redis.Client) WorkoutRepository {
	return &workoutRepository{
		db:    db,
		redis: redis,
	}
}

// Create creates a new workout with exercises and sets
func (r *workoutRepository) Create(ctx context.Context, workout *models.Workout) error {
	return r.db.WithContext(ctx).Create(workout).Error
}

// GetByID retrieves a workout by ID with caching
func (r *workoutRepository) GetByID(ctx context.Context, id uint, userID uint) (*models.Workout, error) {
	cacheKey := fmt.Sprintf("workout:%d", id)

	// Try cache first
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var workout models.Workout
		if err := json.Unmarshal([]byte(cached), &workout); err == nil {
			if workout.UserID == userID {
				return &workout, nil
			}
		}
	}

	// Fetch from database
	var workout models.Workout
	if err := r.db.WithContext(ctx).
		Preload("Exercises.Sets").
		Where("id = ? AND user_id = ?", id, userID).
		First(&workout).Error; err != nil {
		return nil, err
	}

	// Cache the result
	if data, err := json.Marshal(workout); err == nil {
		r.redis.Set(ctx, cacheKey, data, 15*time.Minute)
	}

	return &workout, nil
}

// GetByUser retrieves workouts for a user with pagination
func (r *workoutRepository) GetByUser(ctx context.Context, userID uint, query *models.ListQuery) ([]models.Workout, int64, error) {
	query.SetDefaults()

	var workouts []models.Workout
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&models.Workout{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	db := r.db.WithContext(ctx).
		Preload("Exercises.Sets").
		Where("user_id = ?", userID)

	// Apply sorting
	if query.SortBy != "" {
		db = db.Order(fmt.Sprintf("%s %s", query.SortBy, query.Order))
	} else {
		db = db.Order("date DESC")
	}

	if err := db.
		Offset(query.GetOffset()).
		Limit(query.PageSize).
		Find(&workouts).Error; err != nil {
		return nil, 0, err
	}

	return workouts, total, nil
}

// GetByDateRange retrieves workouts within a date range
func (r *workoutRepository) GetByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]models.Workout, error) {
	var workouts []models.Workout

	if err := r.db.WithContext(ctx).
		Preload("Exercises.Sets").
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("date DESC").
		Find(&workouts).Error; err != nil {
		return nil, err
	}

	return workouts, nil
}

// Update updates a workout and invalidates cache
func (r *workoutRepository) Update(ctx context.Context, workout *models.Workout) error {
	if err := r.db.WithContext(ctx).Save(workout).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("workout:%d", workout.ID)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// Delete soft deletes a workout and invalidates cache
func (r *workoutRepository) Delete(ctx context.Context, id uint, userID uint) error {
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Workout{}).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("workout:%d", id)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// CreateExercise creates a new exercise
func (r *workoutRepository) CreateExercise(ctx context.Context, exercise *models.Exercise) error {
	if err := r.db.WithContext(ctx).Create(exercise).Error; err != nil {
		return err
	}

	// Invalidate workout cache
	cacheKey := fmt.Sprintf("workout:%d", exercise.WorkoutID)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// UpdateExercise updates an exercise
func (r *workoutRepository) UpdateExercise(ctx context.Context, exercise *models.Exercise) error {
	if err := r.db.WithContext(ctx).Save(exercise).Error; err != nil {
		return err
	}

	// Invalidate workout cache
	cacheKey := fmt.Sprintf("workout:%d", exercise.WorkoutID)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// DeleteExercise deletes an exercise
func (r *workoutRepository) DeleteExercise(ctx context.Context, id uint, workoutID uint) error {
	if err := r.db.WithContext(ctx).
		Where("id = ? AND workout_id = ?", id, workoutID).
		Delete(&models.Exercise{}).Error; err != nil {
		return err
	}

	// Invalidate workout cache
	cacheKey := fmt.Sprintf("workout:%d", workoutID)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// CreateSet creates a new set
func (r *workoutRepository) CreateSet(ctx context.Context, set *models.WorkoutSet) error {
	return r.db.WithContext(ctx).Create(set).Error
}

// UpdateSet updates a set
func (r *workoutRepository) UpdateSet(ctx context.Context, set *models.WorkoutSet) error {
	return r.db.WithContext(ctx).Save(set).Error
}

// DeleteSet deletes a set
func (r *workoutRepository) DeleteSet(ctx context.Context, id uint, exerciseID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND exercise_id = ?", id, exerciseID).
		Delete(&models.WorkoutSet{}).Error
}

// GetWorkoutStats retrieves workout statistics for a user
func (r *workoutRepository) GetWorkoutStats(ctx context.Context, userID uint) (*models.WorkoutStatsResponse, error) {
	var stats models.WorkoutStatsResponse

	// Get total workouts
	r.db.WithContext(ctx).
		Model(&models.Workout{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalWorkouts)

	// Get total exercises
	r.db.WithContext(ctx).
		Model(&models.Exercise{}).
		Joins("JOIN workouts ON exercises.workout_id = workouts.id").
		Where("workouts.user_id = ?", userID).
		Count(&stats.TotalExercises)

	// Get total sets
	r.db.WithContext(ctx).
		Model(&models.WorkoutSet{}).
		Joins("JOIN exercises ON workout_sets.exercise_id = exercises.id").
		Joins("JOIN workouts ON exercises.workout_id = workouts.id").
		Where("workouts.user_id = ?", userID).
		Count(&stats.TotalSets)

	// Get last workout date
	var lastWorkout models.Workout
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("date DESC").
		First(&lastWorkout).Error; err == nil {
		stats.LastWorkoutDate = &lastWorkout.Date
	}

	// Calculate average duration
	r.db.WithContext(ctx).
		Model(&models.Workout{}).
		Where("user_id = ? AND duration IS NOT NULL", userID).
		Select("AVG(duration)").
		Scan(&stats.AverageDuration)

	return &stats, nil
}
