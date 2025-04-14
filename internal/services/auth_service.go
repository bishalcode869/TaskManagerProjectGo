package services

import (
	"TaskManager/internal/models"
	"TaskManager/internal/repositories"
)

// AuthService interface defines the methods for authentication-related operations
type AuthService interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
}

// AuthServiceImpl is the concrete implementation of the AuthService interface
type AuthServiceImpl struct {
	AuthRepo repositories.UserRepository
}

// NewAuthService creates and returns a new AuthService instance
func NewAuthService(authRepo repositories.UserRepository) AuthService {
	return &AuthServiceImpl{
		AuthRepo: authRepo,
	}
}

// GetUserByEmail retrieves a user by their email using the repository method
func (s *AuthServiceImpl) GetUserByEmail(email string) (*models.User, error) {
	return s.AuthRepo.GetUserByEmail(email)
}

// GetUserByUsername retrieves a user by their username using the repository method
func (s *AuthServiceImpl) GetUserByUsername(username string) (*models.User, error) {
	return s.AuthRepo.GetUserByUsername(username)
}

// CreateUser creates a new user using the repository method
func (s *AuthServiceImpl) CreateUser(user *models.User) (*models.User, error) {
	return s.AuthRepo.CreateUser(user)
}
