
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gymie-backend/internal/models"
	"github.com/yourusername/gymie-backend/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.UserRegisterRequest true "Registration details"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	authResponse, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"registration_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(
		"User registered successfully",
		authResponse,
	))
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.UserLoginRequest true "Login credentials"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	authResponse, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			"login_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Login successful",
		authResponse,
	))
}

// LoginWithGoogle handles Google Sign-In authentication
// @Summary Google Sign-In
// @Description Authenticate user with Google ID token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.GoogleSignInRequest true "Google sign-in details"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/google [post]
func (h *AuthHandler) LoginWithGoogle(c *gin.Context) {
	var req models.GoogleSignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	authResponse, err := h.authService.LoginWithGoogle(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			"google_signin_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Google sign-in successful",
		authResponse,
	))
}

// Me returns the current user's profile
// @Summary Get current user
// @Description Get the authenticated user's profile
// @Tags auth
// @Produce json
// @Security Bearer
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			"unauthorized",
			"User not authenticated",
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"User retrieved successfully",
		gin.H{"user_id": userID},
	))
}
