package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gymie-backend/internal/models"
	"github.com/yourusername/gymie-backend/internal/service"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
				"unauthorized",
				"Authorization header is required",
				nil,
			))
			c.Abort()
			return
		}

		// Check Bearer token format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
				"unauthorized",
				"Invalid authorization header format. Expected: Bearer <token>",
				nil,
			))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
				"unauthorized",
				"Invalid or expired token",
				err.Error(),
			))
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)

		c.Next()
	}
}

// GetUserID retrieves the user ID from the context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	if !ok {
		return 0, false
	}
	return id, true
}

// GetUserEmail retrieves the user email from the context
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("userEmail")
	if !exists {
		return "", false
	}
	e, ok := email.(string)
	if !ok {
		return "", false
	}
	return e, true
}
