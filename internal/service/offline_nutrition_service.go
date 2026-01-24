package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/yourusername/gymie-backend/internal/models"
	"gorm.io/gorm"
)

// OfflineNutritionService handles offline nutrition system operations
type OfflineNutritionService struct {
	db *gorm.DB
}

// NewOfflineNutritionService creates a new offline nutrition service
func NewOfflineNutritionService(db *gorm.DB) *OfflineNutritionService {
	return &OfflineNutritionService{
		db: db,
	}
}

// =============================================================================
// CORRECTION SYNC
// =============================================================================

// SyncCorrections processes user corrections from devices
func (s *OfflineNutritionService) SyncCorrections(ctx context.Context, userID uint, corrections []models.CorrectionData) (*models.SyncCorrectionsResponse, error) {
	response := &models.SyncCorrectionsResponse{
		Synced: 0,
		Failed: 0,
		Errors: []string{},
	}

	for _, correction := range corrections {
		// Check if already exists (idempotency)
		var existing models.UserCorrection
		err := s.db.WithContext(ctx).
			Where("user_id = ? AND image_hash = ?", userID, correction.ImageHash).
			First(&existing).Error

		if err == nil {
			// Already synced
			response.Synced++
			continue
		}

		if err != gorm.ErrRecordNotFound {
			// Database error
			response.Failed++
			response.Errors = append(response.Errors, fmt.Sprintf("hash %s: db error", correction.ImageHash))
			continue
		}

		// Create new correction
		userCorrection := models.UserCorrection{
			UserID:           userID,
			ImageHash:        correction.ImageHash,
			PredictedDishID:  correction.PredictedDishID,
			CorrectedDishID:  correction.CorrectedDishID,
			PredictedPortion: correction.PredictedPortion,
			CorrectedPortion: correction.CorrectedPortion,
			Confidence:       correction.Confidence,
			DeviceType:       correction.DeviceType,
			AppVersion:       correction.AppVersion,
		}

		if err := s.db.WithContext(ctx).Create(&userCorrection).Error; err != nil {
			response.Failed++
			response.Errors = append(response.Errors, fmt.Sprintf("hash %s: %v", correction.ImageHash, err))
			continue
		}

		response.Synced++
	}

	return response, nil
}

// =============================================================================
// MODEL VERSIONS
// =============================================================================

// GetModelVersions returns the latest model versions
func (s *OfflineNutritionService) GetModelVersions(ctx context.Context) (*models.ModelVersionsResponse, error) {
	response := &models.ModelVersionsResponse{}

	// Get vision model version
	var visionModel models.ModelVersion
	err := s.db.WithContext(ctx).
		Where("model_type = ? AND is_active = ?", "vision", true).
		Order("created_at DESC").
		First(&visionModel).Error

	if err == nil {
		response.VisionModel = &models.ModelVersionInfo{
			Version:     visionModel.Version,
			Checksum:    visionModel.Checksum,
			SizeBytes:   visionModel.SizeBytes,
			DownloadURL: visionModel.DownloadURL,
		}
	}

	// Get nutrition DB version
	var nutritionDB models.ModelVersion
	err = s.db.WithContext(ctx).
		Where("model_type = ? AND is_active = ?", "nutrition_db", true).
		Order("created_at DESC").
		First(&nutritionDB).Error

	if err == nil {
		response.NutritionDB = &models.ModelVersionInfo{
			Version:     nutritionDB.Version,
			Checksum:    nutritionDB.Checksum,
			SizeBytes:   nutritionDB.SizeBytes,
			DownloadURL: nutritionDB.DownloadURL,
		}
	}

	return response, nil
}

// =============================================================================
// DISH SEARCH
// =============================================================================

// SearchDishes searches for dishes by name or alias
func (s *OfflineNutritionService) SearchDishes(ctx context.Context, query string, category string, limit int) (*models.DishSearchResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	query = strings.ToLower(strings.TrimSpace(query))

	// Search in dish master
	var dishes []models.DishMaster
	q := s.db.WithContext(ctx).
		Where("LOWER(display_name) LIKE ? OR LOWER(dish_id) LIKE ?", "%"+query+"%", "%"+query+"%")

	if category != "" {
		q = q.Where("category = ?", category)
	}

	if err := q.Limit(limit).Find(&dishes).Error; err != nil {
		return nil, fmt.Errorf("failed to search dishes: %w", err)
	}

	// Get nutrition data for found dishes
	dishIDs := make([]string, len(dishes))
	for i, dish := range dishes {
		dishIDs[i] = dish.DishID
	}

	var nutritionData []models.DishNutritionMaster
	if len(dishIDs) > 0 {
		s.db.WithContext(ctx).
			Where("dish_id IN ?", dishIDs).
			Find(&nutritionData)
	}

	// Map nutrition to dishes
	nutritionMap := make(map[string]*models.DishNutritionMaster)
	for i := range nutritionData {
		nutritionMap[nutritionData[i].DishID] = &nutritionData[i]
	}

	// Build response
	results := make([]models.DishSearchResult, len(dishes))
	for i, dish := range dishes {
		results[i] = models.DishSearchResult{
			DishID:      dish.DishID,
			DisplayName: dish.DisplayName,
			Category:    dish.Category,
			Cuisine:     dish.Cuisine,
			Aliases:     dish.Aliases,
			Nutrition:   nutritionMap[dish.DishID],
		}
	}

	return &models.DishSearchResponse{
		Dishes: results,
		Total:  len(results),
	}, nil
}

// =============================================================================
// NUTRITION DATABASE MANAGEMENT
// =============================================================================

// GetDishByID retrieves a single dish with nutrition
func (s *OfflineNutritionService) GetDishByID(ctx context.Context, dishID string) (*models.DishSearchResult, error) {
	var dish models.DishMaster
	if err := s.db.WithContext(ctx).Where("dish_id = ?", dishID).First(&dish).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("dish not found: %s", dishID)
		}
		return nil, fmt.Errorf("failed to fetch dish: %w", err)
	}

	var nutrition models.DishNutritionMaster
	_ = s.db.WithContext(ctx).Where("dish_id = ?", dishID).First(&nutrition).Error

	result := &models.DishSearchResult{
		DishID:      dish.DishID,
		DisplayName: dish.DisplayName,
		Category:    dish.Category,
		Cuisine:     dish.Cuisine,
		Aliases:     dish.Aliases,
		Nutrition:   &nutrition,
	}

	return result, nil
}

// =============================================================================
// CORRECTION ANALYTICS
// =============================================================================

// GetCorrectionStats returns correction statistics
func (s *OfflineNutritionService) GetCorrectionStats(ctx context.Context, dishID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total corrections
	var totalCount int64
	query := s.db.WithContext(ctx).Model(&models.UserCorrection{})
	if dishID != "" {
		query = query.Where("predicted_dish_id = ?", dishID)
	}
	query.Count(&totalCount)
	stats["total_corrections"] = totalCount

	// Corrections by dish
	var dishCorrections []struct {
		PredictedDishID string
		Count           int64
	}
	s.db.WithContext(ctx).
		Model(&models.UserCorrection{}).
		Select("predicted_dish_id, COUNT(*) as count").
		Group("predicted_dish_id").
		Order("count DESC").
		Limit(10).
		Scan(&dishCorrections)
	stats["top_corrected_dishes"] = dishCorrections

	// Corrections by device type
	var deviceCorrections []struct {
		DeviceType string
		Count      int64
	}
	s.db.WithContext(ctx).
		Model(&models.UserCorrection{}).
		Select("device_type, COUNT(*) as count").
		Group("device_type").
		Scan(&deviceCorrections)
	stats["by_device"] = deviceCorrections

	return stats, nil
}
