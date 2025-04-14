package routes

import (
	"TaskManager/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes sets up the routes related to authentication
func SetupAuthRoutes(router *gin.Engine, authController *controllers.AuthController) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
	}
}
