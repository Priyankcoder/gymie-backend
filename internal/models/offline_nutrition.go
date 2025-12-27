
package models

import (
	"time"
)

// DishMaster represents the master dish information
type DishMaster struct {
	DishID      string    `json:"dish_id" gorm:"primaryKey"`
	DisplayName string    `json:"display_name" gorm:"not null"`
	Category    string    `json:"category" gorm:"not null;index"`
	Cuisine     string    `json:"cuisine" gorm:"not null;index"`
	Aliases     []string  `json:"aliases" gorm:"type:text[]"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name for GORM
func (DishMaster) TableName() string {
	return "dish_master"
}

// DishNutritionMaster represents the master nutrition data
type DishNutritionMaster struct {
	DishID           string    `json:"dish_id" gorm:"primaryKey"`
	BaseServingGrams int       `json:"base_serving_grams" gorm:"not null"`
	Calories         float64   `json:"calories" gorm:"not null"`
	Protein          float64   `json:"protein" gorm:"not null"`
	Carbs            float64   `json:"carbs" gorm:"not null"`
	Fat              float64   `json:"fat" gorm:"not null"`
	Fiber            float64   `json:"fiber" gorm:"default:0"`
	Sodium           float64   `json:"sodium" gorm:"default:0"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TableName specifies the table name for GORM
func (DishNutritionMaster) TableName() string {
	return "dish_nutrition_master"
}

// UserCorrection represents a user's correction
type UserCorrection struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	UserID            uint      `json:"user_id" gorm:"not null;index"`
	ImageHash         string    `json:"image_hash" gorm:"not null;index"`
	PredictedDishID   string    `json:"predicted_dish_id" gorm:"not null;index"`
	CorrectedDishID   *string   `json:"corrected_dish_id,omitempty"`
	PredictedPortion  string    `json:"predicted_portion" gorm:"not null"`
	CorrectedPortion  *string   `json:"corrected_portion,omitempty"`
	Confidence        float64   `json:"confidence" gorm:"not null"`
	DeviceType        string    `json:"device_type" gorm:"not null"` // android, ios, web
	AppVersion        string    `json:"app_version"`
	CreatedAt         time.Time `json:"created_at"`
}

// TableName specifies the table name for GORM
func (UserCorrection) TableName() string {
	return "user_corrections"
}

// ModelVersion represents a model version
type ModelVersion struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ModelType   string    `json:"model_type" gorm:"not null;uniqueIndex:idx_model_version"` // vision, nutrition_db
	Version     string    `json:"version" gorm:"not null;uniqueIndex:idx_model_version"`
	Checksum    string    `json:"checksum" gorm:"not null"`
	SizeBytes   int64     `json:"size_bytes" gorm:"not null"`
	DownloadURL string    `json:"download_url" gorm:"not null"`
	IsActive    bool      `json:"is_active" gorm:"default:false;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name for GORM
func (ModelVersion) TableName() string {
	return "model_versions"
}

// =============================================================================
// Request/Response Models
// =============================================================================

// SyncCorrectionsRequest represents the correction sync request
type SyncCorrectionsRequest struct {
	Corrections []CorrectionData `json:"corrections" binding:"required"`
}

// CorrectionData represents a single correction
type CorrectionData struct {
	ImageHash        string  `json:"image_hash" binding:"required"`
	PredictedDishID  string  `json:"predicted_dish_id" binding:"required"`
	CorrectedDishID  *string `json:"corrected_dish_id,omitempty"`
	PredictedPortion string  `json:"predicted_portion" binding:"required"`
	CorrectedPortion *string `json:"corrected_portion,omitempty"`
	Confidence       float64 `json:"confidence" binding:"required"`
	DeviceType       string  `json:"device_type" binding:"required"`
	AppVersion       string  `json:"app_version"`
}

// SyncCorrectionsResponse represents the correction sync response
type SyncCorrectionsResponse struct {
	Synced   int      `json:"synced"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors,omitempty"`
}

// ModelVersionsResponse represents the model versions response
type ModelVersionsResponse struct {
	VisionModel    *ModelVersionInfo `json:"vision_model"`
	NutritionDB    *ModelVersionInfo `json:"nutrition_db"`
	ServerTime     time.Time         `json:"server_time"`
}

// ModelVersionInfo represents a single model version info
type ModelVersionInfo struct {
	Version     string `json:"version"`
	Checksum    string `json:"checksum"`
	SizeBytes   int64  `json:"size_bytes"`
	DownloadURL string `json:"download_url"`
}

// NutritionDBResponse represents the nutrition database download response
type NutritionDBResponse struct {
	Version     string    `json:"version"`
	Checksum    string    `json:"checksum"`
	SizeBytes   int64     `json:"size_bytes"`
	DownloadURL string    `json:"download_url"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DishSearchRequest represents a dish search request
type DishSearchRequest struct {
	Query    string `json:"query" binding:"required"`
	Limit    int    `json:"limit"`
	Category string `json:"category,omitempty"`
}

// DishSearchResponse represents a dish search response
type DishSearchResponse struct {
	Dishes []DishSearchResult `json:"dishes"`
	Total  int                `json:"total"`
}

// DishSearchResult represents a single dish search result
type DishSearchResult struct {
	DishID      string   `json:"dish_id"`
	DisplayName string   `json:"display_name"`
	Category    string   `json:"category"`
	Cuisine     string   `json:"cuisine"`
	Aliases     []string `json:"aliases"`
	Nutrition   *DishNutritionMaster `json:"nutrition,omitempty"`
}
