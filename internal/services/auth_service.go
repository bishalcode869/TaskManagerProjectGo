package services

import (
	"errors"
	"fmt"
	"time"

	"TaskManager/internal/models"
	"TaskManager/internal/repositories"
	"TaskManager/pkg/utils"

	"gorm.io/gorm"
)

// AuthService interface defines the methods for authentication-related operations
type AuthService interface {
	RegisterUser(username, password, email string) (*models.User, string, error)
	LoginUser(username, email, password string) (*models.User, string, error)
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

// userExists checks if a user exists by username or email (internal helper)
func (s *AuthServiceImpl) userExists(username, email string) error {
	if username != "" {
		user, err := s.AuthRepo.GetUserByUsername(username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error checking username: %v", err)
		}
		if user != nil {
			return errors.New("username already taken")
		}
	}
	if email != "" {
		user, err := s.AuthRepo.GetUserByEmail(email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error checking email: %v", err)
		}
		if user != nil {
			return errors.New("email already registered")
		}
	}
	return nil
}

// RegisterUser registers a new user, hashes their password, and returns a JWT token
func (s *AuthServiceImpl) RegisterUser(username, password, email string) (*models.User, string, error) {
	// ensure username/email are not already in use
	if err := s.userExists(username, email); err != nil {
		return nil, "", err
	}

	// hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", errors.New("failed to hash password")
	}

	// create and persist user
	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}
	createdUser, err := s.AuthRepo.CreateUser(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %v", err)
	}

	// generate token
	token, err := utils.GenerateJWT(createdUser.ID, time.Hour*24)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}
	return createdUser, token, nil
}

// LoginUser authenticates a user and returns a JWT token
func (s *AuthServiceImpl) LoginUser(username, email, password string) (*models.User, string, error) {
	var (
		user *models.User
		err  error
	)

	// lookup by email or username
	if email != "" {
		user, err = s.AuthRepo.GetUserByEmail(email)
	} else if username != "" {
		user, err = s.AuthRepo.GetUserByUsername(username)
	}

	if err != nil || user == nil {
		return nil, "", errors.New("invalid username/email or password")
	}

	// verify password
	if err := utils.ComparePasswords(user.Password, password); err != nil {
		return nil, "", errors.New("invalid username/email or password")
	}

	// generate token
	token, err := utils.GenerateJWT(user.ID, time.Hour*24)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}
	return user, token, nil
}
