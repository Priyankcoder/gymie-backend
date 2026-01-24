package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gymie-backend/internal/models"
	"github.com/yourusername/gymie-backend/internal/service"
)

// OfflineNutritionHandler handles offline nutrition sync endpoints
type OfflineNutritionHandler struct {
	service *service.OfflineNutritionService
}

// NewOfflineNutritionHandler creates a new offline nutrition handler
func NewOfflineNutritionHandler(service *service.OfflineNutritionService) *OfflineNutritionHandler {
	return &OfflineNutritionHandler{
		service: service,
	}
}

// =============================================================================
// SYNC ENDPOINTS
// =============================================================================

// SyncCorrections handles user correction uploads
// POST /api/sync/corrections
func (h *OfflineNutritionHandler) SyncCorrections(c *gin.Context) {
	// Get user ID from auth context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("auth_required", "Authentication required", nil))
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("auth_error", "Invalid user", nil))
		return
	}

	// Parse request
	var req models.SyncCorrectionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("validation_error", "Invalid request", err.Error()))
		return
	}

	// Process corrections
	response, err := h.service.SyncCorrections(c.Request.Context(), uid, req.Corrections)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("sync_failed", "Failed to sync corrections", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse("Corrections synced successfully", response))
}

// GetModelVersions returns the latest model versions
// GET /api/sync/model-versions
func (h *OfflineNutritionHandler) GetModelVersions(c *gin.Context) {
	response, err := h.service.GetModelVersions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("fetch_failed", "Failed to fetch model versions", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse("Model versions retrieved successfully", response))
}

// SearchDishes searches for dishes by name or alias
// GET /api/sync/dishes/search?query=biryani&category=rice&limit=10
func (h *OfflineNutritionHandler) SearchDishes(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("validation_error", "Query parameter is required", nil))
		return
	}

	category := c.Query("category")
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := parseIntParam(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	response, err := h.service.SearchDishes(c.Request.Context(), query, category, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("search_failed", "Failed to search dishes", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse("Dishes found", response))
}

// GetDishByID retrieves a single dish with nutrition
// GET /api/sync/dishes/:dish_id
func (h *OfflineNutritionHandler) GetDishByID(c *gin.Context) {
	dishID := c.Param("dish_id")
	if dishID == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("validation_error", "Dish ID is required", nil))
		return
	}

	dish, err := h.service.GetDishByID(c.Request.Context(), dishID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse("not_found", "Dish not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse("Dish retrieved successfully", dish))
}

// GetCorrectionStats returns correction analytics
// GET /api/sync/corrections/stats?dish_id=BIRYANI_CHICKEN
func (h *OfflineNutritionHandler) GetCorrectionStats(c *gin.Context) {
	dishID := c.Query("dish_id")

	stats, err := h.service.GetCorrectionStats(c.Request.Context(), dishID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("fetch_failed", "Failed to fetch stats", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse("Stats retrieved successfully", stats))
}

// DownloadNutritionDB provides the nutrition database file
// GET /api/sync/nutrition-db
func (h *OfflineNutritionHandler) DownloadNutritionDB(c *gin.Context) {
	// In production, this would serve a pre-generated SQLite file from CDN
	// For now, return metadata
	response := models.NutritionDBResponse{
		Version:     "1.0.0",
		Checksum:    "sha256:placeholder",
		SizeBytes:   4567890,
		DownloadURL: "/static/nutrition.db",
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse("Nutrition DB info", response))
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func parseIntParam(str string) (int, error) {
	var val int
	_, err := fmt.Sscanf(str, "%d", &val)
	return val, err
}
