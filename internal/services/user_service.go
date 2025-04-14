// internal/services/user_service.go
package services

import (
	"TaskManager/internal/models"
	"TaskManager/internal/repositories"
)

// UserService interface defines the methods for user-related business operations
type UserService interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)       // Get user by email
	GetUserByUsername(username string) (*models.User, error) // Get user by username
	GetAllUsers() ([]models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id uint) error
}

// UserServiceImpl is the concrete implementation of the UserService interface
type UserServiceImpl struct {
	UserRepo repositories.UserRepository
}

// NewUserService creates and returns a new UserService instance
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &UserServiceImpl{
		UserRepo: userRepo,
	}
}

// CreateUser adds a new user by calling the repository's CreateUser method
func (s *UserServiceImpl) CreateUser(user *models.User) (*models.User, error) {
	return s.UserRepo.CreateUser(user)
}

// GetUserByID retrieves a user by their ID by calling the repository's GetUserByID method
func (s *UserServiceImpl) GetUserByID(id uint) (*models.User, error) {
	return s.UserRepo.GetUserByID(id)
}

// GetUserByEmail retrieves a user by their email by calling the repository's GetUserByEmail method
func (s *UserServiceImpl) GetUserByEmail(email string) (*models.User, error) {
	return s.UserRepo.GetUserByEmail(email)
}

// GetUserByUsername retrieves a user by their username by calling the repository's GetUserByUsername method
func (s *UserServiceImpl) GetUserByUsername(username string) (*models.User, error) {
	return s.UserRepo.GetUserByUsername(username)
}

// GetAllUsers retrieves all users by calling the repository's GetAllUsers method
func (s *UserServiceImpl) GetAllUsers() ([]models.User, error) {
	return s.UserRepo.GetAllUsers()
}

// UpdateUser updates an existing user's information by calling the repository's UpdateUser method
func (s *UserServiceImpl) UpdateUser(user *models.User) (*models.User, error) {
	return s.UserRepo.UpdateUser(user)
}

// DeleteUser deletes a user by their ID by calling the repository's DeleteUser method
func (s *UserServiceImpl) DeleteUser(id uint) error {
	return s.UserRepo.DeleteUser(id)
}
