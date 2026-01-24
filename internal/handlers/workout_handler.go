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

type WorkoutHandler struct {
	workoutService service.WorkoutService
}

func NewWorkoutHandler(workoutService service.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{
		workoutService: workoutService,
	}
}

// Create creates a new workout
// @Summary Create workout
// @Description Create a new workout with exercises and sets
// @Tags workouts
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.WorkoutCreateRequest true "Workout details"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /workouts [post]
func (h *WorkoutHandler) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req models.WorkoutCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	workout, err := h.workoutService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"create_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(
		"Workout created successfully",
		workout,
	))
}

// GetByID retrieves a workout by ID
// @Summary Get workout
// @Description Get a specific workout by ID
// @Tags workouts
// @Produce json
// @Security Bearer
// @Param id path int true "Workout ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /workouts/{id} [get]
func (h *WorkoutHandler) GetByID(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	workout, err := h.workoutService.GetByID(c.Request.Context(), uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"not_found",
			"Workout not found",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Workout retrieved successfully",
		workout,
	))
}

// List lists workouts for the current user
// @Summary List workouts
// @Description Get all workouts for the authenticated user
// @Tags workouts
// @Produce json
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} models.PaginatedResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /workouts [get]
func (h *WorkoutHandler) List(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var query models.ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid query parameters",
			err.Error(),
		))
		return
	}

	var dateQuery models.DateRangeQuery
	c.ShouldBindQuery(&dateQuery)

	// Check if date range is provided
	if dateQuery.StartDate != "" && dateQuery.EndDate != "" {
		startDate, err1 := time.Parse("2006-01-02", dateQuery.StartDate)
		endDate, err2 := time.Parse("2006-01-02", dateQuery.EndDate)

		if err1 == nil && err2 == nil {
			workouts, err := h.workoutService.GetByDateRange(c.Request.Context(), userID, startDate, endDate)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
					"fetch_failed",
					err.Error(),
					nil,
				))
				return
			}

			c.JSON(http.StatusOK, models.NewSuccessResponse(
				"Workouts retrieved successfully",
				workouts,
			))
			return
		}
	}

	// Regular pagination
	workouts, total, err := h.workoutService.GetByUser(c.Request.Context(), userID, &query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"fetch_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewPaginatedResponse(workouts, total, query.Page, query.PageSize))
}

// Update updates a workout
// @Summary Update workout
// @Description Update an existing workout
// @Tags workouts
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Workout ID"
// @Param request body models.WorkoutUpdateRequest true "Update details"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /workouts/{id} [put]
func (h *WorkoutHandler) Update(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req models.WorkoutUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	workout, err := h.workoutService.Update(c.Request.Context(), uint(id), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"update_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Workout updated successfully",
		workout,
	))
}

// Delete deletes a workout
// @Summary Delete workout
// @Description Delete a workout
// @Tags workouts
// @Produce json
// @Security Bearer
// @Param id path int true "Workout ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /workouts/{id} [delete]
func (h *WorkoutHandler) Delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.workoutService.Delete(c.Request.Context(), uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"delete_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Workout deleted successfully",
		nil,
	))
}

// GetStats retrieves workout statistics
// @Summary Get workout stats
// @Description Get workout statistics for the authenticated user
// @Tags workouts
// @Produce json
// @Security Bearer
// @Success 200 {object} models.SuccessResponse
// @Router /workouts/stats [get]
func (h *WorkoutHandler) GetStats(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	stats, err := h.workoutService.GetStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"fetch_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Workout statistics retrieved successfully",
		stats,
	))
}
