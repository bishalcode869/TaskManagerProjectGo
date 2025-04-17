package services

import (
	"TaskManager/internal/models"
	"TaskManager/internal/repositories"
	"TaskManager/pkg/utils"
	"errors"
	"fmt"
	"time"
)

// AuthService interface defines the methods for authentication-related operations
type AuthService interface {
	LoginUser(username, email, password string) (*models.User, string, error)
	RegisterUser(username, password, email string) (*models.User, string, error)
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

// Helper function to check if a user exists by either username or email
func (s *AuthServiceImpl) userExists(username, email string) (*models.User, error) {
	if username != "" {
		user, err := s.AuthRepo.GetUserByUsername(username)
		if err == nil && user != nil {
			return nil, fmt.Errorf("Username already taken")
		}
	}
	if email != "" {
		user, err := s.AuthRepo.GetUserByEmail(email)
		if err == nil && user != nil {
			return nil, fmt.Errorf("Email already registered")
		}
	}
	return nil, nil
}

// RegisterUser registers a new user, hashes their password, and returns a JWT token
func (s *AuthServiceImpl) RegisterUser(username, password, email string) (*models.User, string, error) {
	// Check if username or email already exists
	if _, err := s.userExists(username, email); err != nil {
		return nil, "", err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", errors.New("Failed to hash password")
	}

	// Create user
	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	// Save user to DB
	createdUser, err := s.AuthRepo.CreateUser(user)
	if err != nil {
		return nil, "", err
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(createdUser.ID, time.Hour*24)
	if err != nil {
		return nil, "", err
	}

	return createdUser, token, nil
}

// LoginUser authenticates a user and returns a JWT token
func (s *AuthServiceImpl) LoginUser(username, email, password string) (*models.User, string, error) {
	var user *models.User
	var err error

	if email != "" {
		user, err = s.AuthRepo.GetUserByEmail(email)
	} else if username != "" {
		user, err = s.AuthRepo.GetUserByUsername(username)
	}

	if err != nil || user == nil {
		return nil, "", fmt.Errorf("Invalid username/email or password")
	}

	// Compare hashed password
	if err := utils.ComparePasswords(user.Password, password); err != nil {
		return nil, "", fmt.Errorf("Invalid username/email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, time.Hour*24)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to generate token")
	}
	return user, token, nil
}
