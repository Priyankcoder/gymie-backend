package models

import (
	"time"

	"gorm.io/gorm"
)

// NutritionDay represents a day's nutrition tracking
type NutritionDay struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"userId"`
	Date      time.Time      `gorm:"not null;index:idx_user_date,unique" json:"date"`
	Notes     string         `json:"notes,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User  User   `gorm:"foreignKey:UserID" json:"-"`
	Meals []Meal `gorm:"foreignKey:NutritionDayID;constraint:OnDelete:CASCADE" json:"meals,omitempty"`

	// Calculated fields (not stored in DB)
	TotalCalories int `gorm:"-" json:"totalCalories"`
	TotalProtein  int `gorm:"-" json:"totalProtein"`
	TotalCarbs    int `gorm:"-" json:"totalCarbs"`
	TotalFat      int `gorm:"-" json:"totalFat"`
}

// Meal represents a meal within a nutrition day
type Meal struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	NutritionDayID uint           `gorm:"not null;index" json:"nutritionDayId"`
	Name           string         `gorm:"not null" json:"name"` // "Breakfast", "Lunch", etc.
	Time           *time.Time     `json:"time,omitempty"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	NutritionDay NutritionDay `gorm:"foreignKey:NutritionDayID" json:"-"`
	Foods        []Food       `gorm:"foreignKey:MealID;constraint:OnDelete:CASCADE" json:"foods,omitempty"`

	// Calculated fields (not stored in DB)
	TotalCalories int `gorm:"-" json:"totalCalories"`
	TotalProtein  int `gorm:"-" json:"totalProtein"`
	TotalCarbs    int `gorm:"-" json:"totalCarbs"`
	TotalFat      int `gorm:"-" json:"totalFat"`
}

// Food represents a food item within a meal
type Food struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	MealID    uint           `gorm:"not null;index" json:"mealId"`
	Name      string         `gorm:"not null" json:"name"`
	Calories  int            `gorm:"not null" json:"calories"`
	Protein   int            `json:"protein"` // in grams
	Carbs     int            `json:"carbs"`   // in grams
	Fat       int            `json:"fat"`     // in grams
	Quantity  float64        `gorm:"default:1" json:"quantity"`
	Unit      string         `json:"unit,omitempty"` // "serving", "grams", "oz", etc.
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Meal Meal `gorm:"foreignKey:MealID" json:"-"`
}

// NutritionDayCreateRequest represents the request to create a nutrition day
type NutritionDayCreateRequest struct {
	Date  time.Time           `json:"date" binding:"required"`
	Notes string              `json:"notes,omitempty"`
	Meals []MealCreateRequest `json:"meals,omitempty"`
}

// MealCreateRequest represents the request to create a meal
type MealCreateRequest struct {
	Name  string              `json:"name" binding:"required"`
	Time  *time.Time          `json:"time,omitempty"`
	Foods []FoodCreateRequest `json:"foods,omitempty"`
}

// FoodCreateRequest represents the request to create a food item
type FoodCreateRequest struct {
	Name     string  `json:"name" binding:"required"`
	Calories int     `json:"calories" binding:"required"`
	Protein  int     `json:"protein"`
	Carbs    int     `json:"carbs"`
	Fat      int     `json:"fat"`
	Quantity float64 `json:"quantity" binding:"required"`
	Unit     string  `json:"unit,omitempty"`
}

// NutritionDayUpdateRequest represents the request to update a nutrition day
type NutritionDayUpdateRequest struct {
	Date  *time.Time `json:"date,omitempty"`
	Notes *string    `json:"notes,omitempty"`
}

// MealUpdateRequest represents the request to update a meal
type MealUpdateRequest struct {
	Name *string    `json:"name,omitempty"`
	Time *time.Time `json:"time,omitempty"`
}

// AddMealRequest represents the request to add a meal item to a nutrition day
type AddMealRequest struct {
	MealType string `json:"mealType" binding:"required"` // "breakfast", "lunch", "dinner", "snack"
	Name     string `json:"name" binding:"required"`      // Food item name
	Calories int    `json:"calories" binding:"required"`
	Protein  int    `json:"protein"`
	Carbs    int    `json:"carbs"`
	Fat      int    `json:"fat"`
}

// FoodUpdateRequest represents the request to update a food item
type FoodUpdateRequest struct {
	Name     *string  `json:"name,omitempty"`
	Calories *int     `json:"calories,omitempty"`
	Protein  *int     `json:"protein,omitempty"`
	Carbs    *int     `json:"carbs,omitempty"`
	Fat      *int     `json:"fat,omitempty"`
	Quantity *float64 `json:"quantity,omitempty"`
	Unit     *string  `json:"unit,omitempty"`
}

// NutritionStatsResponse represents nutrition statistics
type NutritionStatsResponse struct {
	TotalDays       int64      `json:"total_days"`
	AverageCalories float64    `json:"average_calories"`
	AverageProtein  float64    `json:"average_protein"`
	AverageCarbs    float64    `json:"average_carbs"`
	AverageFat      float64    `json:"average_fat"`
	LastLoggedDate  *time.Time `json:"last_logged_date,omitempty"`
}

// CalculateTotals calculates the total macros for a nutrition day
func (n *NutritionDay) CalculateTotals() {
	n.TotalCalories = 0
	n.TotalProtein = 0
	n.TotalCarbs = 0
	n.TotalFat = 0

	for _, meal := range n.Meals {
		meal.CalculateTotals()
		n.TotalCalories += meal.TotalCalories
		n.TotalProtein += meal.TotalProtein
		n.TotalCarbs += meal.TotalCarbs
		n.TotalFat += meal.TotalFat
	}
}

// CalculateTotals calculates the total macros for a meal
func (m *Meal) CalculateTotals() {
	m.TotalCalories = 0
	m.TotalProtein = 0
	m.TotalCarbs = 0
	m.TotalFat = 0

	for _, food := range m.Foods {
		m.TotalCalories += int(float64(food.Calories) * food.Quantity)
		m.TotalProtein += int(float64(food.Protein) * food.Quantity)
		m.TotalCarbs += int(float64(food.Carbs) * food.Quantity)
		m.TotalFat += int(float64(food.Fat) * food.Quantity)
	}
}
