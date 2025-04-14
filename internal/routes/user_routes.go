package routes

import (
	"TaskManager/internal/controllers"
	"TaskManager/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, userController *controllers.UserController) {
	userRoutes := router.Group("/users")
	{
		// applying jwt middleware
		userRoutes.Use(middleware.AuthRequired())

		// POST to create a new user
		userRoutes.POST("/", userController.CreateUser)

		// GET a single user by ID
		userRoutes.GET("/:id", userController.GetUserByID)

		// GET all users (make sure to place this route before the :id route)
		userRoutes.GET("/", userController.GetAllUsers)

		// PUT to update a user by ID
		userRoutes.PUT("/:id", userController.UpdateUser)

		// DELETE a user by ID
		userRoutes.DELETE("/:id", userController.DeleteUser)
	}
}
