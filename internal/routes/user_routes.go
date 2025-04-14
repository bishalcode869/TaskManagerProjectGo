package routes

import (
	"TaskManager/internal/controllers"
	"TaskManager/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes sets up the routes realted to users
func SetupUserRoutes(router *gin.Engine, userController *controllers.UserController) {
	userRoutes := router.Group("/users")
	{

		// applying jwt middleware
		userRoutes.Use(middleware.AuthRequired())
		userRoutes.POST("/", userController.CreateUser)
		userRoutes.GET("/:id", userController.GetUserByID)
		userRoutes.GET("/", userController.GetAllUsers)
		userRoutes.PUT("/:id", userController.UpdateUser)
		userRoutes.DELETE("/:id", userController.DeleteUser)
	}
}
