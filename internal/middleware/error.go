package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gymie-backend/internal/models"
)

// ErrorHandler handles errors and returns appropriate responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Log the error
			log.Printf("Error: %v", err.Err)

			// Determine status code
			statusCode := http.StatusInternalServerError
			if c.Writer.Status() != http.StatusOK {
				statusCode = c.Writer.Status()
			}

			// Return error response
			c.JSON(statusCode, models.NewErrorResponse(
				"error",
				err.Error(),
				nil,
			))
		}
	}
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
					"internal_server_error",
					"An unexpected error occurred",
					nil,
				))
				c.Abort()
			}
		}()
		c.Next()
	}
}
