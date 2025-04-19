// internal/services/user_service.go
package services

import (
	"TaskManager/internal/models"
	"TaskManager/internal/repositories"
	"TaskManager/pkg/utils"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
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
	HashFunc func(string) (string, error)
}

// NewUserService creates and returns a new UserService instance
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &UserServiceImpl{
		UserRepo: userRepo,
		HashFunc: utils.HashPassword,
	}
}

// CreateUser adds a new user by calling the repository's CreateUser method
func (s *UserServiceImpl) CreateUser(user *models.User) (*models.User, error) {
	// Validate email and username uniqueness
	if _, err := s.UserRepo.GetUserByEmail(user.Email); err == nil {
		return nil, errors.New("email already taken")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Error checking email:", err)
		return nil, fmt.Errorf("unexpected error checking email: %v", err)
	}

	if _, err := s.UserRepo.GetUserByUsername(user.Username); err == nil {
		return nil, errors.New("username already taken")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Error checking username:", err)
		return nil, fmt.Errorf("unexpected error checking username: %v", err)
	}

	// Hash the password before saving
	hashedPassword, err := s.HashFunc(user.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return nil, err
	}
	user.Password = hashedPassword

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
