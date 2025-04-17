package main

import (
	"TaskManager/internal/bootstrap"
	"TaskManager/internal/middleware"
	"TaskManager/internal/routes"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	// Initialize application (config, DB, controllers, etc.)
	app, err := bootstrap.InitializeApp()
	if err != nil {
		log.Fatal("❌ App initialization failed:", err)
	}

	// Initalize Gin router
	router := gin.Default()

	// Global middleware
	router.Use(middleware.Errorhandler())

	// Health check or welcome route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "🚀 Hello, TaskManager API is working!"})
	})

	routes.SetupUserRoutes(router, app.Controller.User)
	routes.SetupAuthRoutes(router, app.Controller.Auth)

	log.Println("Server is running at http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
