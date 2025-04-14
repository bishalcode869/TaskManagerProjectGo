package middleware

import (
	"TaskManager/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthRequired middleware validates JWT tokens for protected routes
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// The token is passed as "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		userID, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the user ID in the context to use in handlers
		c.Set("user_id", userID)

		// Continue with the request
		c.Next()
	}
}
