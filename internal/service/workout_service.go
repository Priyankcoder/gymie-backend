
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/gymie-backend/internal/models"
	"github.com/yourusername/gymie-backend/internal/repository"
)

// WorkoutService defines the interface for workout operations
type WorkoutService interface {
	Create(ctx context.Context, userID uint, req *models.WorkoutCreateRequest) (*models.Workout, error)
	GetByID(ctx context.Context, id uint, userID uint) (*models.Workout, error)
	GetByUser(ctx context.Context, userID uint, query *models.ListQuery) ([]models.Workout, int64, error)
	GetByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]models.Workout, error)
	Update(ctx context.Context, id uint, userID uint, req *models.WorkoutUpdateRequest) (*models.Workout, error)
	Delete(ctx context.Context, id uint, userID uint) error
	GetStats(ctx context.Context, userID uint) (*models.WorkoutStatsResponse, error)
}

type workoutService struct {
	workoutRepo repository.WorkoutRepository
}

// NewWorkoutService creates a new workout service
func NewWorkoutService(workoutRepo repository.WorkoutRepository) WorkoutService {
	return &workoutService{
		workoutRepo: workoutRepo,
	}
}

// Create creates a new workout
func (s *workoutService) Create(ctx context.Context, userID uint, req *models.WorkoutCreateRequest) (*models.Workout, error) {
	workout := &models.Workout{
		UserID:   userID,
		Name:     req.Name,
		Date:     req.Date,
		Duration: req.Duration,
		Notes:    req.Notes,
	}

	// Add exercises
	if req.Exercises != nil {
		for _, exReq := range req.Exercises {
			exercise := models.Exercise{
				Name:        exReq.Name,
				MuscleGroup: exReq.MuscleGroup,
				Order:       exReq.Order,
				Notes:       exReq.Notes,
			}

			// Add sets
			if exReq.Sets != nil {
				for _, setReq := range exReq.Sets {
					set := models.WorkoutSet{
						SetNumber: setReq.SetNumber,
						Weight:    setReq.Weight,
						Reps:      setReq.Reps,
						Completed: setReq.Completed,
					}
					exercise.Sets = append(exercise.Sets, set)
				}
			}

			workout.Exercises = append(workout.Exercises, exercise)
		}
	}

	if err := s.workoutRepo.Create(ctx, workout); err != nil {
		return nil, fmt.Errorf("failed to create workout: %w", err)
	}

	return workout, nil
}

// GetByID retrieves a workout by ID
func (s *workoutService) GetByID(ctx context.Context, id uint, userID uint) (*models.Workout, error) {
	workout, err := s.workoutRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workout: %w", err)
	}
	return workout, nil
}

// GetByUser retrieves workouts for a user
func (s *workoutService) GetByUser(ctx context.Context, userID uint, query *models.ListQuery) ([]models.Workout, int64, error) {
	workouts, total, err := s.workoutRepo.GetByUser(ctx, userID, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get workouts: %w", err)
	}
	return workouts, total, nil
}

// GetByDateRange retrieves workouts within a date range
func (s *workoutService) GetByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]models.Workout, error) {
	workouts, err := s.workoutRepo.GetByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get workouts by date range: %w", err)
	}
	return workouts, nil
}

// Update updates a workout
func (s *workoutService) Update(ctx context.Context, id uint, userID uint, req *models.WorkoutUpdateRequest) (*models.Workout, error) {
	// Get existing workout
	workout, err := s.workoutRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workout: %w", err)
	}

	// Update fields
	if req.Name != nil {
		workout.Name = *req.Name
	}
	if req.Date != nil {
		workout.Date = *req.Date
	}
	if req.Duration != nil {
		workout.Duration = req.Duration
	}
	if req.Notes != nil {
		workout.Notes = *req.Notes
	}

	if err := s.workoutRepo.Update(ctx, workout); err != nil {
		return nil, fmt.Errorf("failed to update workout: %w", err)
	}

	return workout, nil
}

// Delete deletes a workout
func (s *workoutService) Delete(ctx context.Context, id uint, userID uint) error {
	if err := s.workoutRepo.Delete(ctx, id, userID); err != nil {
		return fmt.Errorf("failed to delete workout: %w", err)
	}
	return nil
}

// GetStats retrieves workout statistics
func (s *workoutService) GetStats(ctx context.Context, userID uint) (*models.WorkoutStatsResponse, error) {
	stats, err := s.workoutRepo.GetWorkoutStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workout stats: %w", err)
	}
	return stats, nil
}
