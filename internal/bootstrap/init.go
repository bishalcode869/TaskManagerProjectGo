package bootstrap

import (
	"TaskManager/internal/config"
	"TaskManager/internal/controllers"
	"TaskManager/internal/models"
	"TaskManager/internal/repositories"

	"gorm.io/gorm"
)

type AppContainer struct {
	DB             *gorm.DB
	UserController *controllers.UserController
	AuthController *controllers.AuthController
}

func InitializeApp() (*AppContainer, error) {
	// Load environment variables
	config.LoadConfig()

	dbService := config.NewDBService()
	db, err := dbService.Connect()
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, err
	}

	// Initalize repositories
	userRepo := repositories.NewUserRepository(db)

	// Initalize controllers
	userController := controllers.NewUserController(userRepo)
	authController := controllers.NewAuthController(userRepo)

	return &AppContainer{
		DB:             db,
		UserController: userController,
		AuthController: authController,
	}, nil
}
