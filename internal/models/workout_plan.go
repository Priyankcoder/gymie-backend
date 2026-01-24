package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// WorkoutPlan represents a workout plan/program
type WorkoutPlan struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"userId"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description,omitempty"`
	Type        string         `gorm:"not null" json:"type"` // ppl, push_pull, upper_lower, bro_split, full_body, custom
	IsActive    bool           `gorm:"default:false" json:"isActive"`
	Color       string         `json:"color,omitempty"`
	Recurrence  datatypes.JSON `json:"recurrence,omitempty"` // JSON field for PlanRecurrence
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User      User               `gorm:"foreignKey:UserID" json:"-"`
	Days      []WorkoutPlanDay   `gorm:"foreignKey:PlanID;constraint:OnDelete:CASCADE" json:"days,omitempty"`
	Scheduled []ScheduledWorkout `gorm:"foreignKey:PlanID;constraint:OnDelete:CASCADE" json:"-"`
}

// WorkoutPlanDay represents a day in a workout plan
type WorkoutPlanDay struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	PlanID     uint           `gorm:"not null;index" json:"planId"`
	DayIndex   int            `gorm:"not null" json:"dayIndex"` // 0-6 for days or sequential
	Name       string         `gorm:"not null" json:"name"`
	IsRestDay  bool           `gorm:"default:false" json:"isRestDay"`
	TemplateID *uint          `json:"templateId,omitempty"`
	Exercises  datatypes.JSON `json:"exercises"` // JSON array of TemplateExercise
	Notes      string         `json:"notes,omitempty"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Plan WorkoutPlan `gorm:"foreignKey:PlanID" json:"-"`
}

// ScheduledWorkout represents a scheduled workout instance
type ScheduledWorkout struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"userId"`
	PlanID    uint           `gorm:"not null;index" json:"planId"`
	PlanDayID uint           `gorm:"not null" json:"planDayId"`
	Date      time.Time      `gorm:"not null;index" json:"date"`
	Status    string         `gorm:"not null;default:scheduled" json:"status"` // scheduled, completed, skipped, rescheduled
	WorkoutID *uint          `json:"workoutId,omitempty"`                      // Reference to completed workout
	Notes     string         `json:"notes,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User    User           `gorm:"foreignKey:UserID" json:"-"`
	Plan    WorkoutPlan    `gorm:"foreignKey:PlanID" json:"-"`
	PlanDay WorkoutPlanDay `gorm:"foreignKey:PlanDayID" json:"-"`
	Workout *Workout       `gorm:"foreignKey:WorkoutID" json:"workout,omitempty"`
}

// Request/Response DTOs

// PlanRecurrence represents the recurrence pattern for a plan
type PlanRecurrence struct {
	Type          string   `json:"type"` // weekly, biweekly, monthly, custom
	Interval      int      `json:"interval"`
	StartDate     string   `json:"startDate"`
	EndDate       *string  `json:"endDate,omitempty"`
	RestDays      []int    `json:"restDays"` // Day indices that are rest days
	ExcludedDates []string `json:"excludedDates,omitempty"`
}

// TemplateExercise represents an exercise template
type TemplateExercise struct {
	Name         string      `json:"name"`
	TargetSets   int         `json:"targetSets"`
	TargetReps   interface{} `json:"targetReps"` // Can be int or string like "8-12"
	TargetWeight *float64    `json:"targetWeight,omitempty"`
	RestSeconds  *int        `json:"restSeconds,omitempty"`
	Notes        string      `json:"notes,omitempty"`
	SupersetWith string      `json:"supersetWith,omitempty"`
}

// WorkoutPlanCreateRequest represents the request to create a workout plan
type WorkoutPlanCreateRequest struct {
	Name        string                  `json:"name" binding:"required"`
	Description string                  `json:"description,omitempty"`
	Type        string                  `json:"type" binding:"required"`
	IsActive    bool                    `json:"isActive"`
	Color       string                  `json:"color,omitempty"`
	Days        []WorkoutPlanDayRequest `json:"days" binding:"required"`
}

// WorkoutPlanDayRequest represents a day in the plan creation request
type WorkoutPlanDayRequest struct {
	DayIndex   int                `json:"dayIndex" binding:"required"`
	Name       string             `json:"name" binding:"required"`
	IsRestDay  bool               `json:"isRestDay"`
	TemplateID *uint              `json:"templateId,omitempty"`
	Exercises  []TemplateExercise `json:"exercises"`
	Notes      string             `json:"notes,omitempty"`
}

// WorkoutPlanUpdateRequest represents the request to update a workout plan
type WorkoutPlanUpdateRequest struct {
	Name        *string                 `json:"name,omitempty"`
	Description *string                 `json:"description,omitempty"`
	Type        *string                 `json:"type,omitempty"`
	IsActive    *bool                   `json:"isActive,omitempty"`
	Color       *string                 `json:"color,omitempty"`
	Days        []WorkoutPlanDayRequest `json:"days,omitempty"`
}

// SetRecurrenceRequest represents the request to set plan recurrence
type SetRecurrenceRequest struct {
	Recurrence *PlanRecurrence `json:"recurrence"`
}

// ScheduledWorkoutCreateRequest represents the request to create a scheduled workout
type ScheduledWorkoutCreateRequest struct {
	PlanID    uint      `json:"planId" binding:"required"`
	PlanDayID uint      `json:"planDayId" binding:"required"`
	Date      time.Time `json:"date" binding:"required"`
	Notes     string    `json:"notes,omitempty"`
}

// ScheduledWorkoutUpdateStatusRequest represents the request to update scheduled workout status
type ScheduledWorkoutUpdateStatusRequest struct {
	Status    string `json:"status" binding:"required,oneof=scheduled completed skipped rescheduled"`
	WorkoutID *uint  `json:"workoutId,omitempty"`
	Notes     string `json:"notes,omitempty"`
}

// GenerateScheduleRequest represents the request to generate scheduled workouts from recurrence
type GenerateScheduleRequest struct {
	PlanID    uint   `json:"planId" binding:"required"`
	StartDate string `json:"startDate" binding:"required"`
	EndDate   string `json:"endDate" binding:"required"`
}
