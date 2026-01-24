package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"` // Never send password in JSON
	Name      string         `gorm:"not null" json:"name"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Email verification
	EmailVerified              bool      `gorm:"default:false" json:"emailVerified"`
	VerificationToken          string    `gorm:"index:idx_users_verification_token;unique" json:"-"`
	VerificationTokenExpiresAt time.Time `json:"-"`

	// Profile information
	Profile *UserProfile `gorm:"foreignKey:UserID" json:"profile,omitempty"`

	// Relationships
	Workouts       []Workout       `gorm:"foreignKey:UserID" json:"workouts,omitempty"`
	NutritionDays  []NutritionDay  `gorm:"foreignKey:UserID" json:"nutritionDays,omitempty"`
	ProgressPhotos []ProgressPhoto `gorm:"foreignKey:UserID" json:"progressPhotos,omitempty"`
	WeightEntries  []WeightEntry   `gorm:"foreignKey:UserID" json:"weightEntries,omitempty"`
}

// UserProfile represents additional user profile information
type UserProfile struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"uniqueIndex;not null" json:"userId"`
	DisplayName    *string   `json:"displayName,omitempty"`    // Custom display name (can be different from Name)
	ProfilePicture *string   `json:"profilePicture,omitempty"` // URL to profile picture
	Bio            *string   `json:"bio,omitempty"`            // User bio/description
	Height         *float64  `json:"height,omitempty"`         // in cm
	Weight         *float64  `json:"weight,omitempty"`         // in kg
	Age            *int      `json:"age,omitempty"`
	Gender         string    `json:"gender,omitempty"` // "male", "female", "other"
	Goal           string    `json:"goal,omitempty"`   // "lose_weight", "gain_muscle", "maintain"
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// UserRegisterRequest represents the registration request
type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

// UserLoginRequest represents the login request
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserUpdateRequest represents the user update request
type UserUpdateRequest struct {
	Name           *string  `json:"name,omitempty"`
	DisplayName    *string  `json:"displayName,omitempty"`
	ProfilePicture *string  `json:"profilePicture,omitempty"`
	Bio            *string  `json:"bio,omitempty"`
	Email          *string  `json:"email,omitempty" binding:"omitempty,email"`
	Height         *float64 `json:"height,omitempty"`
	Weight         *float64 `json:"weight,omitempty"`
	Age            *int     `json:"age,omitempty"`
	Gender         *string  `json:"gender,omitempty"`
	Goal           *string  `json:"goal,omitempty"`
}

// UserResponse represents the user response (without sensitive data)
type UserResponse struct {
	ID            uint         `json:"id"`
	Email         string       `json:"email"`
	Name          string       `json:"name"`
	EmailVerified bool         `json:"emailVerified"`
	Profile       *UserProfile `json:"profile,omitempty"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		Name:          u.Name,
		EmailVerified: u.EmailVerified,
		Profile:       u.Profile,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}
