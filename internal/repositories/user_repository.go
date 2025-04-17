// internal/repositories/user_repository.go
package repositories

import (
	"TaskManager/internal/models"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// UserRepository interface defines the methods for user-related DB operations
type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)       // Get user by email
	GetUserByUsername(username string) (*models.User, error) // Get user by username
	GetAllUsers() ([]models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id uint) error
}

// UserRepositoryImpl is the concrete implementation of the UserRepository interface
type UserRepositoryImpl struct {
	DB *gorm.DB
}

// NewUserRepository creates and returns a new UserRepository instance
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}

// CreateUser adds a new user to the database
func (repo *UserRepositoryImpl) CreateUser(user *models.User) (*models.User, error) {
	if err := repo.DB.Create(user).Error; err != nil {
		log.Println("Error creating user:", err)
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user by their ID
func (repo *UserRepositoryImpl) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := repo.DB.First(&user, id).Error; err != nil {
		log.Println("Error fetching user by ID:", err)
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email
func (repo *UserRepositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := repo.DB.Where("email = ?", email).First(&user).Error; err != nil {
		log.Println("Error fetching user by email", err)
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by their username
func (repo *UserRepositoryImpl) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := repo.DB.Where("username = ?", username).First(&user).Error; err != nil {
		log.Println("Error fetching user by username:", err)
		return nil, err
	}
	return &user, nil
}

// GetAllUsers retrieves all users from the database
func (repo *UserRepositoryImpl) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := repo.DB.Find(&users).Error; err != nil {
		log.Println("Error fetching all users:", err)
		return nil, err
	}
	return users, nil
}

// UpdateUser updates an existing user's information
func (repo *UserRepositoryImpl) UpdateUser(user *models.User) (*models.User, error) {
	if err := repo.DB.Save(user).Error; err != nil {
		log.Println("Error updating user:", err)
		return nil, err
	}
	return user, nil
}

// DeleteUser deletes a user by their ID
func (repo *UserRepositoryImpl) DeleteUser(id uint) error {
	var user models.User
	if err := repo.DB.First(&user, id).Error; err != nil {
		return fmt.Errorf("user not found")
	}
	return repo.DB.Delete(&user).Error
}
