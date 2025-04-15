package bootstrap

import (
	"TaskManager/internal/config"
	"TaskManager/internal/controllers"
	"TaskManager/internal/models"
	"TaskManager/internal/repositories"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Controller struct {
	User *controllers.UserController
	Auth *controllers.AuthController
}

type AppContainer struct {
	DB         *gorm.DB
	Controller Controller
}

func InitializeApp() (*AppContainer, error) {
	// Load configuration
	log.Println("Loading configuration...")
	config.LoadConfig()

	// Connect to the database
	log.Println("Connecting to  the database...")
	dbService := config.NewDBService()
	db, err := dbService.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("Failed to auto-migrate model: %w", err)
	}

	// Initalize repositories
	log.Println("Initializing repositories...")
	userRepo := repositories.NewUserRepository(db)

	// Initalize controllers
	log.Println("Initializing controllers...")
	userController := controllers.NewUserController(userRepo)
	authController := controllers.NewAuthController(userRepo)

	log.Println("Application initialized successfully.")

	return &AppContainer{
		DB: db,
		Controller: Controller{
			User: userController,
			Auth: authController,
		},
	}, nil
}
