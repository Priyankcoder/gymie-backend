
package models

import (
	"time"

	"gorm.io/gorm"
)

// Recipe represents a recipe with nutrition information
type Recipe struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null;index" json:"title"`
	Description string         `json:"description"`
	ImageURL    string         `json:"imageUrl,omitempty"`
	
	// Nutrition Info
	Calories    int     `gorm:"not null" json:"calories"`
	Protein     float64 `gorm:"not null" json:"protein"`     // in grams
	Carbs       float64 `gorm:"not null" json:"carbs"`       // in grams
	Fat         float64 `gorm:"not null" json:"fat"`         // in grams
	Fiber       float64 `json:"fiber,omitempty"`             // in grams
	
	// Recipe Details
	PrepTime    int      `json:"prepTime"`                   // in minutes
	CookTime    int      `json:"cookTime"`                   // in minutes
	Servings    int      `gorm:"not null;default:1" json:"servings"`
	Difficulty  string   `gorm:"default:'medium'" json:"difficulty"` // easy, medium, hard
	
	// Ingredients and Instructions (JSON arrays)
	Ingredients []string `gorm:"type:json" json:"ingredients"`
	Instructions []string `gorm:"type:json" json:"instructions"`
	
	// Tags and Categories
	Category    string   `gorm:"index" json:"category"`      // breakfast, lunch, dinner, snack
	Tags        []string `gorm:"type:json" json:"tags"`      // high-protein, low-carb, vegan, etc.
	
	// Metadata
	IsAIGenerated bool      `gorm:"default:false" json:"isAiGenerated"`
	Source        string    `json:"source,omitempty"`         // URL or "AI Generated"
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// RecipeCreateRequest represents the request to create a recipe
type RecipeCreateRequest struct {
	Title        string   `json:"title" binding:"required"`
	Description  string   `json:"description"`
	ImageURL     string   `json:"imageUrl"`
	Calories     int      `json:"calories" binding:"required"`
	Protein      float64  `json:"protein" binding:"required"`
	Carbs        float64  `json:"carbs" binding:"required"`
	Fat          float64  `json:"fat" binding:"required"`
	Fiber        float64  `json:"fiber"`
	PrepTime     int      `json:"prepTime"`
	CookTime     int      `json:"cookTime"`
	Servings     int      `json:"servings" binding:"required,min=1"`
	Difficulty   string   `json:"difficulty"`
	Ingredients  []string `json:"ingredients" binding:"required"`
	Instructions []string `json:"instructions" binding:"required"`
	Category     string   `json:"category" binding:"required"`
	Tags         []string `json:"tags"`
	Source       string   `json:"source"`
}

// RecipeGenerateRequest represents the request to AI-generate a recipe
type RecipeGenerateRequest struct {
	Category       string   `json:"category"`                    // breakfast, lunch, dinner, snack
	TargetCalories int      `json:"targetCalories"`              // approximate calories
	TargetProtein  float64  `json:"targetProtein,omitempty"`     // in grams
	DietType       string   `json:"dietType,omitempty"`          // vegan, vegetarian, keto, etc.
	Ingredients    []string `json:"ingredients,omitempty"`       // must-include ingredients
	MaxPrepTime    int      `json:"maxPrepTime,omitempty"`       // in minutes
	Difficulty     string   `json:"difficulty,omitempty"`        // easy, medium, hard
}

// RecipeSearchRequest represents the search parameters
type RecipeSearchRequest struct {
	Query          string   `json:"query" form:"query"`
	Category       string   `json:"category" form:"category"`
	Tags           []string `json:"tags" form:"tags"`
	MaxCalories    int      `json:"maxCalories" form:"maxCalories"`
	MinProtein     float64  `json:"minProtein" form:"minProtein"`
	MaxPrepTime    int      `json:"maxPrepTime" form:"maxPrepTime"`
	Difficulty     string   `json:"difficulty" form:"difficulty"`
	Page           int      `json:"page" form:"page"`
	PageSize       int      `json:"pageSize" form:"pageSize"`
}

// SetDefaults sets default values for pagination
func (r *RecipeSearchRequest) SetDefaults() {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.PageSize < 1 || r.PageSize > 100 {
		r.PageSize = 20
	}
}

// GetOffset calculates the offset for pagination
func (r *RecipeSearchRequest) GetOffset() int {
	return (r.Page - 1) * r.PageSize
}

