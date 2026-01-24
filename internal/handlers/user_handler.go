package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gymie-backend/internal/middleware"
	"github.com/yourusername/gymie-backend/internal/models"
	"github.com/yourusername/gymie-backend/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile retrieves the current user's profile
// @Summary Get user profile
// @Description Get the authenticated user's profile with details
// @Tags users
// @Produce json
// @Security Bearer
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			"unauthorized",
			"User not authenticated",
			nil,
		))
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"not_found",
			"User not found",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"User profile retrieved successfully",
		user.ToResponse(),
	))
}

// UpdateProfile updates the current user's profile
// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.UserUpdateRequest true "Update details"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			"unauthorized",
			"User not authenticated",
			nil,
		))
		return
	}

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	user, err := h.userService.Update(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"update_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"User profile updated successfully",
		user.ToResponse(),
	))
}

// DeleteAccount deletes the current user's account
// @Summary Delete user account
// @Description Delete the authenticated user's account (soft delete)
// @Tags users
// @Produce json
// @Security Bearer
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /users/account [delete]
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			"unauthorized",
			"User not authenticated",
			nil,
		))
		return
	}

	if err := h.userService.Delete(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"delete_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Account deleted successfully",
		nil,
	))
}
