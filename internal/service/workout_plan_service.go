package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/yourusername/gymie-backend/internal/models"
	"github.com/yourusername/gymie-backend/internal/repository"
)

// WorkoutPlanService defines the interface for workout plan business logic
type WorkoutPlanService interface {
	// Plan operations
	Create(ctx context.Context, userID uint, req *models.WorkoutPlanCreateRequest) (*models.WorkoutPlan, error)
	GetByID(ctx context.Context, id uint, userID uint) (*models.WorkoutPlan, error)
	GetAll(ctx context.Context, userID uint) ([]models.WorkoutPlan, error)
	GetActive(ctx context.Context, userID uint) (*models.WorkoutPlan, error)
	Update(ctx context.Context, id uint, userID uint, req *models.WorkoutPlanUpdateRequest) (*models.WorkoutPlan, error)
	Delete(ctx context.Context, id uint, userID uint) error
	SetActive(ctx context.Context, id uint, userID uint) error
	Clone(ctx context.Context, id uint, userID uint, newName string) (*models.WorkoutPlan, error)
	SetRecurrence(ctx context.Context, id uint, userID uint, req *models.SetRecurrenceRequest) error

	// Day operations
	UpdateDay(ctx context.Context, planID uint, dayID uint, userID uint, exercises []models.TemplateExercise) error
	AddDay(ctx context.Context, planID uint, userID uint, req *models.WorkoutPlanDayRequest) error
	RemoveDay(ctx context.Context, planID uint, dayID uint, userID uint) error

	// Scheduled workout operations
	CreateScheduled(ctx context.Context, userID uint, req *models.ScheduledWorkoutCreateRequest) (*models.ScheduledWorkout, error)
	GetScheduledByID(ctx context.Context, id uint, userID uint) (*models.ScheduledWorkout, error)
	GetAllScheduled(ctx context.Context, userID uint) ([]models.ScheduledWorkout, error)
	GetScheduledByDateRange(ctx context.Context, userID uint, startDate, endDate string) ([]models.ScheduledWorkout, error)
	GetScheduledByDate(ctx context.Context, userID uint, date string) ([]models.ScheduledWorkout, error)
	GetTodaysWorkout(ctx context.Context, userID uint) (map[string]interface{}, error)
	UpdateScheduledStatus(ctx context.Context, id uint, userID uint, req *models.ScheduledWorkoutUpdateStatusRequest) (*models.ScheduledWorkout, error)
	DeleteScheduled(ctx context.Context, id uint, userID uint) error
	GenerateFromRecurrence(ctx context.Context, userID uint, req *models.GenerateScheduleRequest) ([]models.ScheduledWorkout, error)
	ClearPlanSchedule(ctx context.Context, planID uint, userID uint) error
}

type workoutPlanService struct {
	planRepo repository.WorkoutPlanRepository
}

// NewWorkoutPlanService creates a new workout plan service
func NewWorkoutPlanService(planRepo repository.WorkoutPlanRepository) WorkoutPlanService {
	return &workoutPlanService{
		planRepo: planRepo,
	}
}

// Create creates a new workout plan
func (s *workoutPlanService) Create(ctx context.Context, userID uint, req *models.WorkoutPlanCreateRequest) (*models.WorkoutPlan, error) {
	// If setting as active, deactivate other plans
	if req.IsActive {
		if err := s.planRepo.SetActive(ctx, 0, userID); err != nil {
			return nil, err
		}
	}

	plan := &models.WorkoutPlan{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		IsActive:    req.IsActive,
		Color:       req.Color,
	}

	// Convert days
	for _, dayReq := range req.Days {
		exercisesJSON, err := json.Marshal(dayReq.Exercises)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal exercises: %w", err)
		}

		day := models.WorkoutPlanDay{
			DayIndex:   dayReq.DayIndex,
			Name:       dayReq.Name,
			IsRestDay:  dayReq.IsRestDay,
			TemplateID: dayReq.TemplateID,
			Exercises:  exercisesJSON,
			Notes:      dayReq.Notes,
		}
		plan.Days = append(plan.Days, day)
	}

	if err := s.planRepo.Create(ctx, plan); err != nil {
		return nil, err
	}

	return plan, nil
}

// GetByID retrieves a workout plan by ID
func (s *workoutPlanService) GetByID(ctx context.Context, id uint, userID uint) (*models.WorkoutPlan, error) {
	return s.planRepo.GetByID(ctx, id, userID)
}

// GetAll retrieves all workout plans for a user
func (s *workoutPlanService) GetAll(ctx context.Context, userID uint) ([]models.WorkoutPlan, error) {
	return s.planRepo.GetByUser(ctx, userID)
}

// GetActive retrieves the active workout plan
func (s *workoutPlanService) GetActive(ctx context.Context, userID uint) (*models.WorkoutPlan, error) {
	return s.planRepo.GetActive(ctx, userID)
}

// Update updates a workout plan
func (s *workoutPlanService) Update(ctx context.Context, id uint, userID uint, req *models.WorkoutPlanUpdateRequest) (*models.WorkoutPlan, error) {
	plan, err := s.planRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		plan.Name = *req.Name
	}
	if req.Description != nil {
		plan.Description = *req.Description
	}
	if req.Type != nil {
		plan.Type = *req.Type
	}
	if req.IsActive != nil {
		if *req.IsActive {
			// Deactivate other plans first
			if err := s.planRepo.SetActive(ctx, id, userID); err != nil {
				return nil, err
			}
		}
		plan.IsActive = *req.IsActive
	}
	if req.Color != nil {
		plan.Color = *req.Color
	}

	// Handle days update
	if req.Days != nil && len(req.Days) > 0 {
		// Delete existing days
		if err := s.planRepo.DeleteDays(ctx, id); err != nil {
			return nil, fmt.Errorf("failed to delete existing days: %w", err)
		}

		// Create new days
		var days []models.WorkoutPlanDay
		for _, dayReq := range req.Days {
			exercisesJSON, err := json.Marshal(dayReq.Exercises)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal exercises: %w", err)
			}

			day := models.WorkoutPlanDay{
				PlanID:    id,
				DayIndex:  dayReq.DayIndex,
				Name:      dayReq.Name,
				IsRestDay: dayReq.IsRestDay,
				Exercises: exercisesJSON,
				Notes:     dayReq.Notes,
			}
			if dayReq.TemplateID != nil {
				day.TemplateID = dayReq.TemplateID
			}
			days = append(days, day)
		}

		if err := s.planRepo.CreateDays(ctx, days); err != nil {
			return nil, fmt.Errorf("failed to create days: %w", err)
		}
	}

	if err := s.planRepo.Update(ctx, plan); err != nil {
		return nil, err
	}

	// Reload plan with updated days
	return s.planRepo.GetByID(ctx, id, userID)
}

// Delete deletes a workout plan
func (s *workoutPlanService) Delete(ctx context.Context, id uint, userID uint) error {
	// Also delete all scheduled workouts for this plan
	if err := s.planRepo.DeleteScheduledByPlan(ctx, id, userID); err != nil {
		return err
	}

	return s.planRepo.Delete(ctx, id, userID)
}

// SetActive sets a plan as active
func (s *workoutPlanService) SetActive(ctx context.Context, id uint, userID uint) error {
	// Set the plan as active
	if err := s.planRepo.SetActive(ctx, id, userID); err != nil {
		return err
	}

	// Get the plan to check if it has recurrence settings
	plan, err := s.planRepo.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	// If the plan has recurrence settings, regenerate scheduled workouts
	if plan.Recurrence != nil && len(plan.Recurrence) > 0 {
		var recurrence models.PlanRecurrence
		if err := json.Unmarshal(plan.Recurrence, &recurrence); err == nil {
			// Generate workouts for the next 12 weeks
			startDate := time.Now()
			if recurrence.StartDate != "" {
				if parsedStart, err := time.Parse("2006-01-02", recurrence.StartDate); err == nil {
					startDate = parsedStart
				}
			}

			endDate := startDate.AddDate(0, 0, 84) // 12 weeks
			if recurrence.EndDate != nil && *recurrence.EndDate != "" {
				if parsedEnd, err := time.Parse("2006-01-02", *recurrence.EndDate); err == nil {
					endDate = parsedEnd
				}
			}

			// Generate scheduled workouts
			req := &models.GenerateScheduleRequest{
				PlanID:    id,
				StartDate: startDate.Format("2006-01-02"),
				EndDate:   endDate.Format("2006-01-02"),
			}

			// Ignore errors from generation - plan is already active
			s.GenerateFromRecurrence(ctx, userID, req)
		}
	}

	return nil
}

// Clone clones a workout plan
func (s *workoutPlanService) Clone(ctx context.Context, id uint, userID uint, newName string) (*models.WorkoutPlan, error) {
	original, err := s.planRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	clone := &models.WorkoutPlan{
		UserID:      userID,
		Name:        newName,
		Description: original.Description,
		Type:        original.Type,
		IsActive:    false,
		Color:       original.Color,
		Days:        make([]models.WorkoutPlanDay, len(original.Days)),
	}

	// Clone days
	for i, day := range original.Days {
		clone.Days[i] = models.WorkoutPlanDay{
			DayIndex:   day.DayIndex,
			Name:       day.Name,
			IsRestDay:  day.IsRestDay,
			TemplateID: day.TemplateID,
			Exercises:  day.Exercises,
			Notes:      day.Notes,
		}
	}

	if err := s.planRepo.Create(ctx, clone); err != nil {
		return nil, err
	}

	return clone, nil
}

// SetRecurrence sets the recurrence pattern for a plan
func (s *workoutPlanService) SetRecurrence(ctx context.Context, id uint, userID uint, req *models.SetRecurrenceRequest) error {
	plan, err := s.planRepo.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	if req.Recurrence != nil {
		recurrenceJSON, err := json.Marshal(req.Recurrence)
		if err != nil {
			return fmt.Errorf("failed to marshal recurrence: %w", err)
		}
		plan.Recurrence = recurrenceJSON
	} else {
		plan.Recurrence = nil
	}

	if err := s.planRepo.Update(ctx, plan); err != nil {
		return err
	}

	// If recurrence is set, generate scheduled workouts
	if req.Recurrence != nil {
		endDate := req.Recurrence.StartDate
		if req.Recurrence.EndDate != nil {
			endDate = *req.Recurrence.EndDate
		} else {
			// Generate for 90 days by default
			startTime, _ := time.Parse("2006-01-02", req.Recurrence.StartDate)
			endDate = startTime.AddDate(0, 0, 90).Format("2006-01-02")
		}

		_, err := s.GenerateFromRecurrence(ctx, userID, &models.GenerateScheduleRequest{
			PlanID:    id,
			StartDate: req.Recurrence.StartDate,
			EndDate:   endDate,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateDay updates a day in the plan
func (s *workoutPlanService) UpdateDay(ctx context.Context, planID uint, dayID uint, userID uint, exercises []models.TemplateExercise) error {
	plan, err := s.planRepo.GetByID(ctx, planID, userID)
	if err != nil {
		return err
	}

	// Find the day
	var day *models.WorkoutPlanDay
	for i := range plan.Days {
		if plan.Days[i].ID == dayID {
			day = &plan.Days[i]
			break
		}
	}

	if day == nil {
		return errors.New("day not found")
	}

	exercisesJSON, err := json.Marshal(exercises)
	if err != nil {
		return fmt.Errorf("failed to marshal exercises: %w", err)
	}

	day.Exercises = exercisesJSON

	return s.planRepo.UpdateDay(ctx, day)
}

// AddDay adds a new day to the plan
func (s *workoutPlanService) AddDay(ctx context.Context, planID uint, userID uint, req *models.WorkoutPlanDayRequest) error {
	plan, err := s.planRepo.GetByID(ctx, planID, userID)
	if err != nil {
		return err
	}

	exercisesJSON, err := json.Marshal(req.Exercises)
	if err != nil {
		return fmt.Errorf("failed to marshal exercises: %w", err)
	}

	day := &models.WorkoutPlanDay{
		PlanID:     plan.ID,
		DayIndex:   req.DayIndex,
		Name:       req.Name,
		IsRestDay:  req.IsRestDay,
		TemplateID: req.TemplateID,
		Exercises:  exercisesJSON,
		Notes:      req.Notes,
	}

	return s.planRepo.CreateDay(ctx, day)
}

// RemoveDay removes a day from the plan
func (s *workoutPlanService) RemoveDay(ctx context.Context, planID uint, dayID uint, userID uint) error {
	_, err := s.planRepo.GetByID(ctx, planID, userID)
	if err != nil {
		return err
	}

	return s.planRepo.DeleteDay(ctx, dayID, planID)
}

// CreateScheduled creates a scheduled workout
func (s *workoutPlanService) CreateScheduled(ctx context.Context, userID uint, req *models.ScheduledWorkoutCreateRequest) (*models.ScheduledWorkout, error) {
	scheduled := &models.ScheduledWorkout{
		UserID:    userID,
		PlanID:    req.PlanID,
		PlanDayID: req.PlanDayID,
		Date:      req.Date,
		Status:    "scheduled",
		Notes:     req.Notes,
	}

	if err := s.planRepo.CreateScheduled(ctx, scheduled); err != nil {
		return nil, err
	}

	return scheduled, nil
}

// GetScheduledByID retrieves a scheduled workout by ID
func (s *workoutPlanService) GetScheduledByID(ctx context.Context, id uint, userID uint) (*models.ScheduledWorkout, error) {
	return s.planRepo.GetScheduledByID(ctx, id, userID)
}

// GetAllScheduled retrieves all scheduled workouts
func (s *workoutPlanService) GetAllScheduled(ctx context.Context, userID uint) ([]models.ScheduledWorkout, error) {
	return s.planRepo.GetScheduledByUser(ctx, userID)
}

// GetScheduledByDateRange retrieves scheduled workouts within a date range
func (s *workoutPlanService) GetScheduledByDateRange(ctx context.Context, userID uint, startDate, endDate string) ([]models.ScheduledWorkout, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	return s.planRepo.GetScheduledByDateRange(ctx, userID, start, end)
}

// GetScheduledByDate retrieves scheduled workouts for a specific date
func (s *workoutPlanService) GetScheduledByDate(ctx context.Context, userID uint, date string) ([]models.ScheduledWorkout, error) {
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date: %w", err)
	}

	return s.planRepo.GetScheduledByDate(ctx, userID, dateTime)
}

// GetTodaysWorkout retrieves today's scheduled workout with plan and day details
func (s *workoutPlanService) GetTodaysWorkout(ctx context.Context, userID uint) (map[string]interface{}, error) {
	scheduled, err := s.planRepo.GetTodaysScheduled(ctx, userID)
	if err != nil {
		return nil, err
	}

	if scheduled == nil {
		return map[string]interface{}{
			"scheduled": nil,
			"plan":      nil,
			"day":       nil,
		}, nil
	}

	plan, err := s.planRepo.GetByID(ctx, scheduled.PlanID, userID)
	if err != nil {
		return nil, err
	}

	// Find the day
	var day *models.WorkoutPlanDay
	for i := range plan.Days {
		if plan.Days[i].ID == scheduled.PlanDayID {
			day = &plan.Days[i]
			break
		}
	}

	return map[string]interface{}{
		"scheduled": scheduled,
		"plan":      plan,
		"day":       day,
	}, nil
}

// UpdateScheduledStatus updates the status of a scheduled workout
func (s *workoutPlanService) UpdateScheduledStatus(ctx context.Context, id uint, userID uint, req *models.ScheduledWorkoutUpdateStatusRequest) (*models.ScheduledWorkout, error) {
	scheduled, err := s.planRepo.GetScheduledByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	scheduled.Status = req.Status
	if req.WorkoutID != nil {
		scheduled.WorkoutID = req.WorkoutID
	}
	if req.Notes != "" {
		scheduled.Notes = req.Notes
	}

	if err := s.planRepo.UpdateScheduled(ctx, scheduled); err != nil {
		return nil, err
	}

	return scheduled, nil
}

// DeleteScheduled deletes a scheduled workout
func (s *workoutPlanService) DeleteScheduled(ctx context.Context, id uint, userID uint) error {
	return s.planRepo.DeleteScheduled(ctx, id, userID)
}

// GenerateFromRecurrence generates scheduled workouts from a plan's recurrence
func (s *workoutPlanService) GenerateFromRecurrence(ctx context.Context, userID uint, req *models.GenerateScheduleRequest) ([]models.ScheduledWorkout, error) {
	plan, err := s.planRepo.GetByID(ctx, req.PlanID, userID)
	if err != nil {
		return nil, err
	}

	if plan.Recurrence == nil {
		return nil, errors.New("plan has no recurrence set")
	}

	var recurrence models.PlanRecurrence
	if err := json.Unmarshal([]byte(plan.Recurrence), &recurrence); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recurrence: %w", err)
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	// Clear existing scheduled workouts in this range
	existing, err := s.planRepo.GetScheduledByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	for _, sw := range existing {
		if sw.PlanID == plan.ID {
			s.planRepo.DeleteScheduled(ctx, sw.ID, userID)
		}
	}

	// Generate new scheduled workouts
	var scheduled []models.ScheduledWorkout
	currentDate := startDate
	dayIndex := 0

	for !currentDate.After(endDate) {
		weekday := int(currentDate.Weekday())

		// Check if this is a rest day
		isRestDay := false
		for _, rd := range recurrence.RestDays {
			if rd == weekday {
				isRestDay = true
				break
			}
		}

		if !isRestDay {
			// Find the appropriate day from the plan
			planDay := plan.Days[dayIndex%len(plan.Days)]

			scheduled = append(scheduled, models.ScheduledWorkout{
				UserID:    userID,
				PlanID:    plan.ID,
				PlanDayID: planDay.ID,
				Date:      currentDate,
				Status:    "scheduled",
			})

			dayIndex++
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	if len(scheduled) > 0 {
		if err := s.planRepo.BulkCreateScheduled(ctx, scheduled); err != nil {
			return nil, err
		}
	}

	return scheduled, nil
}

// ClearPlanSchedule clears all scheduled workouts for a plan
func (s *workoutPlanService) ClearPlanSchedule(ctx context.Context, planID uint, userID uint) error {
	return s.planRepo.DeleteScheduledByPlan(ctx, planID, userID)
}
