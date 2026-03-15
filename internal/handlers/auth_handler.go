package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gymie-backend/internal/models"
	"github.com/yourusername/gymie-backend/internal/service"
	"github.com/yourusername/gymie-backend/internal/utils"
	"gorm.io/gorm"
)

type AuthHandler struct {
	authService  service.AuthService
	db           *gorm.DB
	emailService service.EmailServiceInterface
}

func NewAuthHandler(authService service.AuthService, db *gorm.DB, emailService service.EmailServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		db:           db,
		emailService: emailService,
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

// VerifyEmail verifies user's email address
// @Summary Verify email
// @Description Verify user's email address using verification token
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/verify-email [get]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Verification token is required",
			nil,
		))
		return
	}

	var user models.User
	if err := h.db.Where("verification_token = ?", token).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, models.NewErrorResponse(
				"invalid_token",
				"Invalid or expired verification token",
				nil,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"server_error",
			"Failed to verify email",
			err.Error(),
		))
		return
	}

	// Check if token is expired
	if time.Now().After(user.VerificationTokenExpiresAt) {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"token_expired",
			"Verification token has expired. Please request a new one.",
			nil,
		))
		return
	}

	// Check if already verified
	if user.EmailVerified {
		c.JSON(http.StatusOK, models.NewSuccessResponse(
			"Email already verified",
			gin.H{"verified": true},
		))
		return
	}

	// Verify the email - use Updates with map to set token to NULL (not empty string)
	// to avoid unique constraint violation on verification_token
	if err := h.db.Model(&user).Updates(map[string]interface{}{
		"email_verified":               true,
		"verification_token":           nil,
		"verification_token_expires_at": nil,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"server_error",
			"Failed to verify email",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Email verified successfully! You can now log in.",
		gin.H{
			"verified": true,
			"email":    user.Email,
		},
	))
}

// ResendVerification resends verification email
// @Summary Resend verification email
// @Description Resend verification email to user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{email string} true "Email address"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/resend-verification [post]
func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid email address",
			err.Error(),
		))
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Don't reveal if email exists (security best practice)
		c.JSON(http.StatusOK, models.NewSuccessResponse(
			"If the email exists and is not verified, a verification link has been sent",
			nil,
		))
		return
	}

	// Check if already verified
	if user.EmailVerified {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"already_verified",
			"Email is already verified",
			nil,
		))
		return
	}

	// Check rate limiting - only allow resend if last token was created > 60 seconds ago
	if user.VerificationTokenExpiresAt.Sub(time.Now()) > 23*time.Hour+59*time.Minute {
		c.JSON(http.StatusTooManyRequests, models.NewErrorResponse(
			"rate_limit",
			"Please wait a moment before requesting another verification email",
			nil,
		))
		return
	}

	// Generate new verification token
	verificationToken, err := utils.GenerateVerificationToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"server_error",
			"Failed to generate verification token",
			err.Error(),
		))
		return
	}

	tokenExpiry := time.Now().Add(24 * time.Hour)

	user.VerificationToken = &verificationToken
	user.VerificationTokenExpiresAt = tokenExpiry

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"server_error",
			"Failed to update verification token",
			err.Error(),
		))
		return
	}

	// Send verification email (pass token, not full link)
	if err := h.emailService.SendVerificationEmail(user.Email, user.Name, verificationToken); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"email_error",
			"Failed to send verification email",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Verification email sent successfully",
		gin.H{"email": user.Email},
	))
}

// GetVerificationStatus gets user's email verification status
// @Summary Get verification status
// @Description Check if user's email is verified
// @Tags auth
// @Produce json
// @Param email path string true "Email address"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /auth/verification-status/{email} [get]
func (h *AuthHandler) GetVerificationStatus(c *gin.Context) {
	email := c.Param("email")

	if email == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Email is required",
			nil,
		))
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.NewErrorResponse(
				"not_found",
				"User not found",
				nil,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"server_error",
			"Failed to check verification status",
			err.Error(),
		))
		return
	}

	// Calculate cooldown remaining (60 seconds from token creation)
	var cooldownRemaining int64 = 0
	if !user.EmailVerified && !user.VerificationTokenExpiresAt.IsZero() {
		// Token creation time = expiry time - 24 hours
		tokenCreationTime := user.VerificationTokenExpiresAt.Add(-24 * time.Hour)

		// Time since token was created
		timeSinceCreation := time.Since(tokenCreationTime)

		// Cooldown is 60 seconds from creation
		cooldownDuration := 60 * time.Second

		if timeSinceCreation < cooldownDuration {
			cooldownRemaining = int64((cooldownDuration - timeSinceCreation).Seconds())
		}
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Verification status retrieved",
		gin.H{
			"email":             user.Email,
			"verified":          user.EmailVerified,
			"cooldownRemaining": cooldownRemaining,
		},
	))
}
