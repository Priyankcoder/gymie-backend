
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gymie-backend/internal/middleware"
	"github.com/yourusername/gymie-backend/internal/models"
	"github.com/yourusername/gymie-backend/internal/service"
)

type NutritionHandler struct {
	nutritionService service.NutritionService
}

func NewNutritionHandler(nutritionService service.NutritionService) *NutritionHandler {
	return &NutritionHandler{
		nutritionService: nutritionService,
	}
}

// Create creates a new nutrition day
// @Summary Create nutrition day
// @Description Create a new nutrition day with meals and foods
// @Tags nutrition
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.NutritionDayCreateRequest true "Nutrition day details"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /nutrition [post]
func (h *NutritionHandler) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req models.NutritionDayCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	day, err := h.nutritionService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"create_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(
		"Nutrition day created successfully",
		day,
	))
}

// GetByID retrieves a nutrition day by ID
// @Summary Get nutrition day
// @Description Get a specific nutrition day by ID
// @Tags nutrition
// @Produce json
// @Security Bearer
// @Param id path int true "Nutrition Day ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /nutrition/{id} [get]
func (h *NutritionHandler) GetByID(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	day, err := h.nutritionService.GetByID(c.Request.Context(), uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"not_found",
			"Nutrition day not found",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Nutrition day retrieved successfully",
		day,
	))
}

// GetByDate retrieves nutrition data for a specific date
// @Summary Get nutrition by date
// @Description Get nutrition data for a specific date
// @Tags nutrition
// @Produce json
// @Security Bearer
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /nutrition/date [get]
func (h *NutritionHandler) GetByDate(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	dateStr := c.Query("date")

	if dateStr == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Date parameter is required",
			nil,
		))
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid date format. Use YYYY-MM-DD",
			err.Error(),
		))
		return
	}

	day, err := h.nutritionService.GetByDate(c.Request.Context(), userID, date)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"not_found",
			"No nutrition data found for this date",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Nutrition day retrieved successfully",
		day,
	))
}

// GetByDateRange retrieves nutrition data within a date range
// @Summary Get nutrition by date range
// @Description Get nutrition data within a date range
// @Tags nutrition
// @Produce json
// @Security Bearer
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /nutrition/range [get]
func (h *NutritionHandler) GetByDateRange(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var dateQuery models.DateRangeQuery
	if err := c.ShouldBindQuery(&dateQuery); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid query parameters",
			err.Error(),
		))
		return
	}

	startDate, err1 := time.Parse("2006-01-02", dateQuery.StartDate)
	endDate, err2 := time.Parse("2006-01-02", dateQuery.EndDate)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid date format. Use YYYY-MM-DD",
			nil,
		))
		return
	}

	days, err := h.nutritionService.GetByDateRange(c.Request.Context(), userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"fetch_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Nutrition days retrieved successfully",
		days,
	))
}

// Update updates a nutrition day
// @Summary Update nutrition day
// @Description Update an existing nutrition day
// @Tags nutrition
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Nutrition Day ID"
// @Param request body models.NutritionDayUpdateRequest true "Update details"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /nutrition/{id} [put]
func (h *NutritionHandler) Update(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req models.NutritionDayUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	day, err := h.nutritionService.Update(c.Request.Context(), uint(id), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"update_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Nutrition day updated successfully",
		day,
	))
}

// Delete deletes a nutrition day
// @Summary Delete nutrition day
// @Description Delete a nutrition day
// @Tags nutrition
// @Produce json
// @Security Bearer
// @Param id path int true "Nutrition Day ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /nutrition/{id} [delete]
func (h *NutritionHandler) Delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.nutritionService.Delete(c.Request.Context(), uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"delete_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Nutrition day deleted successfully",
		nil,
	))
}

// GetStats retrieves nutrition statistics
// @Summary Get nutrition stats
// @Description Get nutrition statistics for the authenticated user
// @Tags nutrition
// @Produce json
// @Security Bearer
// @Success 200 {object} models.SuccessResponse
// @Router /nutrition/stats [get]
func (h *NutritionHandler) GetStats(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	stats, err := h.nutritionService.GetStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"fetch_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Nutrition statistics retrieved successfully",
		stats,
	))
}
