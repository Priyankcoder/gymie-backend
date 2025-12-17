
package models

import (
	"time"

	"gorm.io/gorm"
)

// Workout represents a workout session
type Workout struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Name      string         `gorm:"not null" json:"name"`
	Date      time.Time      `gorm:"not null;index" json:"date"`
	Duration  *int           `json:"duration,omitempty"` // in minutes
	Notes     string         `json:"notes,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User      User       `gorm:"foreignKey:UserID" json:"-"`
	Exercises []Exercise `gorm:"foreignKey:WorkoutID;constraint:OnDelete:CASCADE" json:"exercises,omitempty"`
}

// Exercise represents an exercise within a workout
type Exercise struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	WorkoutID  uint           `gorm:"not null;index" json:"workout_id"`
	Name       string         `gorm:"not null" json:"name"`
	MuscleGroup string        `json:"muscle_group,omitempty"` // "chest", "back", "legs", etc.
	Order      int            `gorm:"not null" json:"order"`  // Order of exercise in workout
	Notes      string         `json:"notes,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Workout Workout      `gorm:"foreignKey:WorkoutID" json:"-"`
	Sets    []WorkoutSet `gorm:"foreignKey:ExerciseID;constraint:OnDelete:CASCADE" json:"sets,omitempty"`
}

// WorkoutSet represents a set within an exercise
type WorkoutSet struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ExerciseID uint           `gorm:"not null;index" json:"exercise_id"`
	SetNumber  int            `gorm:"not null" json:"set_number"`
	Weight     *float64       `json:"weight,omitempty"` // in kg or lbs
	Reps       *int           `json:"reps,omitempty"`
	Completed  bool           `gorm:"default:false" json:"completed"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Exercise Exercise `gorm:"foreignKey:ExerciseID" json:"-"`
}

// WorkoutCreateRequest represents the request to create a workout
type WorkoutCreateRequest struct {
	Name      string                     `json:"name" binding:"required"`
	Date      time.Time                  `json:"date" binding:"required"`
	Duration  *int                       `json:"duration,omitempty"`
	Notes     string                     `json:"notes,omitempty"`
	Exercises []ExerciseCreateRequest    `json:"exercises,omitempty"`
}

// ExerciseCreateRequest represents the request to create an exercise
type ExerciseCreateRequest struct {
	Name        string             `json:"name" binding:"required"`
	MuscleGroup string             `json:"muscle_group,omitempty"`
	Order       int                `json:"order" binding:"required"`
	Notes       string             `json:"notes,omitempty"`
	Sets        []SetCreateRequest `json:"sets,omitempty"`
}

// SetCreateRequest represents the request to create a set
type SetCreateRequest struct {
	SetNumber int      `json:"set_number" binding:"required"`
	Weight    *float64 `json:"weight,omitempty"`
	Reps      *int     `json:"reps,omitempty"`
	Completed bool     `json:"completed"`
}

// WorkoutUpdateRequest represents the request to update a workout
type WorkoutUpdateRequest struct {
	Name     *string    `json:"name,omitempty"`
	Date     *time.Time `json:"date,omitempty"`
	Duration *int       `json:"duration,omitempty"`
	Notes    *string    `json:"notes,omitempty"`
}

// ExerciseUpdateRequest represents the request to update an exercise
type ExerciseUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	MuscleGroup *string `json:"muscle_group,omitempty"`
	Order       *int    `json:"order,omitempty"`
	Notes       *string `json:"notes,omitempty"`
}

// SetUpdateRequest represents the request to update a set
type SetUpdateRequest struct {
	SetNumber *int     `json:"set_number,omitempty"`
	Weight    *float64 `json:"weight,omitempty"`
	Reps      *int     `json:"reps,omitempty"`
	Completed *bool    `json:"completed,omitempty"`
}

// WorkoutStatsResponse represents workout statistics
type WorkoutStatsResponse struct {
	TotalWorkouts    int64   `json:"total_workouts"`
	TotalExercises   int64   `json:"total_exercises"`
	TotalSets        int64   `json:"total_sets"`
	TotalVolume      float64 `json:"total_volume"` // weight * reps
	AverageDuration  float64 `json:"average_duration"`
	LastWorkoutDate  *time.Time `json:"last_workout_date,omitempty"`
}
