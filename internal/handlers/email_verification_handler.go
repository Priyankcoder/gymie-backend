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

type EmailVerificationHandler struct {
	db           *gorm.DB
	emailService service.EmailServiceInterface
}

func NewEmailVerificationHandler(db *gorm.DB, emailService service.EmailServiceInterface) *EmailVerificationHandler {
	return &EmailVerificationHandler{
		db:           db,
		emailService: emailService,
	}
}

// VerifyEmail handles email verification with token
// @Summary Verify user email
// @Description Verify user email address using verification token
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/verify-email [get]
func (h *EmailVerificationHandler) VerifyEmail(c *gin.Context) {
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

	// Verify the email
	user.EmailVerified = true
	user.VerificationToken = "" // Clear the token after use

	if err := h.db.Save(&user).Error; err != nil {
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

// ResendVerificationEmail resends verification email
// @Summary Resend verification email
// @Description Resend verification email to user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.ResendVerificationRequest true "Email address"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 429 {object} models.ErrorResponse
// @Router /auth/resend-verification [post]
func (h *EmailVerificationHandler) ResendVerificationEmail(c *gin.Context) {
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

	// Check rate limiting - only allow resend if last token was created > 5 minutes ago
	if user.VerificationTokenExpiresAt.Sub(time.Now()) > 23*time.Hour+55*time.Minute {
		c.JSON(http.StatusTooManyRequests, models.NewErrorResponse(
			"rate_limit",
			"Please wait a few minutes before requesting another verification email",
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

	user.VerificationToken = verificationToken
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

// GetVerificationStatus checks if email is verified
// @Summary Check email verification status
// @Description Check if user's email is verified
// @Tags auth
// @Accept json
// @Produce json
// @Param email query string true "Email address"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/verification-status [get]
func (h *EmailVerificationHandler) GetVerificationStatus(c *gin.Context) {
	email := c.Query("email")

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

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Verification status retrieved",
		gin.H{
			"email":    user.Email,
			"verified": user.EmailVerified,
		},
	))
}
