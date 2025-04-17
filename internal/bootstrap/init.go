package bootstrap

import (
	"TaskManager/internal/config"
	"TaskManager/internal/controllers"
	"TaskManager/internal/models"
	"TaskManager/internal/repositories"
	"TaskManager/internal/services"
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
	log.Println("🔧 Loading configuration...")
	config.LoadConfig()

	// Connect to the database
	log.Println("💾 Connecting to the database...")
	dbService := config.NewDBService()
	db, err := dbService.Connect()
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("❌ Failed to auto-migrate models: %w", err)
	}

	// Initalize repositories
	log.Println("📦 Initializing repositories...")
	userRepo := repositories.NewUserRepository(db)

	// Initalize service
	log.Println("🧠 Initializing services...")
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)

	// Initalize controllers
	log.Println("🎮 Initializing controllers...")
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)

	log.Println("✅ Application initialized successfully.")

	// Return everything inside the app container
	return &AppContainer{
		DB: db,
		Controller: Controller{
			User: userController,
			Auth: authController,
		},
	}, nil
}
