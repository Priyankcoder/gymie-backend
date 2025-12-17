
package models

import (
	"time"

	"gorm.io/gorm"
)

// ProgressPhoto represents a progress photo
type ProgressPhoto struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Date        time.Time      `gorm:"not null;index" json:"date"`
	ImageURL    string         `gorm:"not null" json:"image_url"`
	Weight      *float64       `json:"weight,omitempty"` // in kg
	Notes       string         `json:"notes,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// WeightEntry represents a weight tracking entry
type WeightEntry struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Date      time.Time      `gorm:"not null;index" json:"date"`
	Weight    float64        `gorm:"not null" json:"weight"` // in kg
	Notes     string         `json:"notes,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// ProgressPhotoCreateRequest represents the request to create a progress photo
type ProgressPhotoCreateRequest struct {
	Date   time.Time `json:"date" binding:"required"`
	Weight *float64  `json:"weight,omitempty"`
	Notes  string    `json:"notes,omitempty"`
	// Image will be uploaded separately via multipart form
}

// ProgressPhotoUpdateRequest represents the request to update a progress photo
type ProgressPhotoUpdateRequest struct {
	Date   *time.Time `json:"date,omitempty"`
	Weight *float64   `json:"weight,omitempty"`
	Notes  *string    `json:"notes,omitempty"`
}

// WeightEntryCreateRequest represents the request to create a weight entry
type WeightEntryCreateRequest struct {
	Date   time.Time `json:"date" binding:"required"`
	Weight float64   `json:"weight" binding:"required"`
	Notes  string    `json:"notes,omitempty"`
}

// WeightEntryUpdateRequest represents the request to update a weight entry
type WeightEntryUpdateRequest struct {
	Date   *time.Time `json:"date,omitempty"`
	Weight *float64   `json:"weight,omitempty"`
	Notes  *string    `json:"notes,omitempty"`
}

// WeightProgressResponse represents weight progress over time
type WeightProgressResponse struct {
	StartWeight   *float64     `json:"start_weight,omitempty"`
	CurrentWeight *float64     `json:"current_weight,omitempty"`
	GoalWeight    *float64     `json:"goal_weight,omitempty"`
	TotalChange   float64      `json:"total_change"`
	AverageChange float64      `json:"average_change"` // per week
	Entries       []WeightEntry `json:"entries"`
}

// UploadURLResponse represents a presigned upload URL response
type UploadURLResponse struct {
	UploadURL string    `json:"upload_url"`
	FileURL   string    `json:"file_url"`
	ExpiresAt time.Time `json:"expires_at"`
}
