package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gymie-backend/internal/middleware"
	"github.com/yourusername/gymie-backend/internal/models"
	"github.com/yourusername/gymie-backend/internal/service"
)

type WorkoutPlanHandler struct {
	planService service.WorkoutPlanService
}

func NewWorkoutPlanHandler(planService service.WorkoutPlanService) *WorkoutPlanHandler {
	return &WorkoutPlanHandler{
		planService: planService,
	}
}

// CreatePlan creates a new workout plan
// @Summary Create workout plan
// @Description Create a new workout plan with days
// @Tags workout-plans
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.WorkoutPlanCreateRequest true "Plan details"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /workout-plans [post]
func (h *WorkoutPlanHandler) CreatePlan(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req models.WorkoutPlanCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	plan, err := h.planService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"create_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(
		"Workout plan created successfully",
		plan,
	))
}

// GetPlan retrieves a workout plan by ID
// @Summary Get workout plan
// @Description Get a specific workout plan by ID
// @Tags workout-plans
// @Produce json
// @Security Bearer
// @Param id path int true "Plan ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /workout-plans/{id} [get]
func (h *WorkoutPlanHandler) GetPlan(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	plan, err := h.planService.GetByID(c.Request.Context(), uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"not_found",
			"Workout plan not found",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Workout plan retrieved successfully",
		plan,
	))
}

// ListPlans lists all workout plans
// @Summary List workout plans
// @Description Get all workout plans for the authenticated user
// @Tags workout-plans
// @Produce json
// @Security Bearer
// @Success 200 {object} models.SuccessResponse
// @Router /workout-plans [get]
func (h *WorkoutPlanHandler) ListPlans(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	plans, err := h.planService.GetAll(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"fetch_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Workout plans retrieved successfully",
		plans,
	))
}

// GetActivePlan retrieves the active workout plan
// @Summary Get active workout plan
// @Description Get the currently active workout plan
// @Tags workout-plans
// @Produce json
// @Security Bearer
// @Success 200 {object} models.SuccessResponse
// @Router /workout-plans/active [get]
func (h *WorkoutPlanHandler) GetActivePlan(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	plan, err := h.planService.GetActive(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"fetch_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Active workout plan retrieved successfully",
		plan,
	))
}

// UpdatePlan updates a workout plan
// @Summary Update workout plan
// @Description Update a workout plan's details
// @Tags workout-plans
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Plan ID"
// @Param request body models.WorkoutPlanUpdateRequest true "Update details"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /workout-plans/{id} [put]
func (h *WorkoutPlanHandler) UpdatePlan(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req models.WorkoutPlanUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	plan, err := h.planService.Update(c.Request.Context(), uint(id), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"update_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Workout plan updated successfully",
		plan,
	))
}

// DeletePlan deletes a workout plan
// @Summary Delete workout plan
// @Description Delete a workout plan
// @Tags workout-plans
// @Produce json
// @Security Bearer
// @Param id path int true "Plan ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /workout-plans/{id} [delete]
func (h *WorkoutPlanHandler) DeletePlan(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.planService.Delete(c.Request.Context(), uint(id), userID); err != nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"delete_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Workout plan deleted successfully",
		nil,
	))
}

// SetActivePlan sets a plan as active
// @Summary Set active plan
// @Description Set a workout plan as the active plan
// @Tags workout-plans
// @Produce json
// @Security Bearer
// @Param id path int true "Plan ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /workout-plans/{id}/active [put]
func (h *WorkoutPlanHandler) SetActivePlan(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.planService.SetActive(c.Request.Context(), uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"set_active_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Workout plan set as active successfully",
		nil,
	))
}

// ClonePlan clones a workout plan
// @Summary Clone workout plan
// @Description Create a copy of an existing workout plan
// @Tags workout-plans
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Plan ID"
// @Param request body map[string]string true "New name"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /workout-plans/{id}/clone [post]
func (h *WorkoutPlanHandler) ClonePlan(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	plan, err := h.planService.Clone(c.Request.Context(), uint(id), userID, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"clone_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(
		"Workout plan cloned successfully",
		plan,
	))
}

// SetRecurrence sets the recurrence pattern for a plan
// @Summary Set plan recurrence
// @Description Set or update the recurrence pattern for a workout plan
// @Tags workout-plans
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Plan ID"
// @Param request body models.SetRecurrenceRequest true "Recurrence details"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /workout-plans/{id}/recurrence [put]
func (h *WorkoutPlanHandler) SetRecurrence(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req models.SetRecurrenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	if err := h.planService.SetRecurrence(c.Request.Context(), uint(id), userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"set_recurrence_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Recurrence set successfully",
		nil,
	))
}

// === Scheduled Workouts ===

// ListScheduled lists all scheduled workouts
// @Summary List scheduled workouts
// @Description Get all scheduled workouts for the authenticated user
// @Tags scheduled-workouts
// @Produce json
// @Security Bearer
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} models.SuccessResponse
// @Router /scheduled-workouts [get]
func (h *WorkoutPlanHandler) ListScheduled(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var scheduled []models.ScheduledWorkout
	var err error

	if startDate != "" && endDate != "" {
		scheduled, err = h.planService.GetScheduledByDateRange(c.Request.Context(), userID, startDate, endDate)
	} else {
		scheduled, err = h.planService.GetAllScheduled(c.Request.Context(), userID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"fetch_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Scheduled workouts retrieved successfully",
		scheduled,
	))
}

// GetScheduled retrieves a scheduled workout by ID
// @Summary Get scheduled workout
// @Description Get a specific scheduled workout by ID
// @Tags scheduled-workouts
// @Produce json
// @Security Bearer
// @Param id path int true "Scheduled Workout ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /scheduled-workouts/{id} [get]
func (h *WorkoutPlanHandler) GetScheduled(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	scheduled, err := h.planService.GetScheduledByID(c.Request.Context(), uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"not_found",
			"Scheduled workout not found",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Scheduled workout retrieved successfully",
		scheduled,
	))
}

// GetTodaysWorkout retrieves today's scheduled workout
// @Summary Get today's workout
// @Description Get today's scheduled workout with plan and day details
// @Tags scheduled-workouts
// @Produce json
// @Security Bearer
// @Success 200 {object} models.SuccessResponse
// @Router /scheduled-workouts/today [get]
func (h *WorkoutPlanHandler) GetTodaysWorkout(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	result, err := h.planService.GetTodaysWorkout(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"fetch_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Today's workout retrieved successfully",
		result,
	))
}

// UpdateScheduledStatus updates the status of a scheduled workout
// @Summary Update scheduled workout status
// @Description Update the status of a scheduled workout (completed, skipped, etc.)
// @Tags scheduled-workouts
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Scheduled Workout ID"
// @Param request body models.ScheduledWorkoutUpdateStatusRequest true "Status update"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /scheduled-workouts/{id}/status [put]
func (h *WorkoutPlanHandler) UpdateScheduledStatus(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req models.ScheduledWorkoutUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	scheduled, err := h.planService.UpdateScheduledStatus(c.Request.Context(), uint(id), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"update_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Scheduled workout status updated successfully",
		scheduled,
	))
}

// DeleteScheduled deletes a scheduled workout
// @Summary Delete scheduled workout
// @Description Delete a scheduled workout
// @Tags scheduled-workouts
// @Produce json
// @Security Bearer
// @Param id path int true "Scheduled Workout ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /scheduled-workouts/{id} [delete]
func (h *WorkoutPlanHandler) DeleteScheduled(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.planService.DeleteScheduled(c.Request.Context(), uint(id), userID); err != nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"delete_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Scheduled workout deleted successfully",
		nil,
	))
}

// GenerateScheduled generates scheduled workouts from recurrence
// @Summary Generate scheduled workouts
// @Description Generate scheduled workouts from a plan's recurrence pattern
// @Tags scheduled-workouts
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.GenerateScheduleRequest true "Generation params"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /scheduled-workouts/generate [post]
func (h *WorkoutPlanHandler) GenerateScheduled(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req models.GenerateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	scheduled, err := h.planService.GenerateFromRecurrence(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"generation_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(
		"Scheduled workouts generated successfully",
		scheduled,
	))
}

// ClearPlanSchedule clears all scheduled workouts for a plan
// @Summary Clear plan schedule
// @Description Delete all scheduled workouts for a specific plan
// @Tags scheduled-workouts
// @Produce json
// @Security Bearer
// @Param plan_id path int true "Plan ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /scheduled-workouts/plan/{plan_id} [delete]
func (h *WorkoutPlanHandler) ClearPlanSchedule(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	planID, _ := strconv.ParseUint(c.Param("plan_id"), 10, 32)

	if err := h.planService.ClearPlanSchedule(c.Request.Context(), uint(planID), userID); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"clear_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Plan schedule cleared successfully",
		nil,
	))
}
