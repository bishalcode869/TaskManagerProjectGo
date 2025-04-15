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

	app, err := bootstrap.InitializeApp()
	if err != nil {
		log.Fatal("❌ App initialization failed:", err)
	}

	router := gin.Default()
	router.Use(middleware.Errorhandler())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, it is working"})
	})

	routes.SetupUserRoutes(router, app.Controller.User)
	routes.SetupAuthRoutes(router, app.Controller.Auth)

	log.Println("Server is running at http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
