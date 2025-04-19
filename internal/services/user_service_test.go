package services_test

import (
	"TaskManager/internal/models"
	"TaskManager/internal/services"
	"TaskManager/mocks"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocking the UserRepository
	mockRepo := mocks.NewMockUserRepository(ctrl)
	userSvc := services.NewUserService(mockRepo)

	// Test data
	newUser := &models.User{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
	}

	// Set up expectations
	mockRepo.EXPECT().GetUserByEmail(newUser.Email).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.EXPECT().GetUserByUsername(newUser.Username).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.EXPECT().CreateUser(gomock.Any()).Return(newUser, nil)

	// Call the service method
	createdUser, err := userSvc.CreateUser(newUser)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, newUser.Email, createdUser.Email)
	assert.Equal(t, newUser.Username, createdUser.Username)
}

func TestCreateUser_EmailAlreadyTaken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userSvc := services.NewUserService(mockRepo)

	newUser := &models.User{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
	}

	// Mock GetUserByEmail to simulate email already taken
	mockRepo.EXPECT().GetUserByEmail(newUser.Email).Return(&models.User{}, nil)

	// Call the service method
	createdUser, err := userSvc.CreateUser(newUser)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, createdUser)
	assert.Equal(t, "email already taken", err.Error())
}

func TestCreateUser_UsernameAlreadyTaken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userSvc := services.NewUserService(mockRepo)

	newUser := &models.User{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
	}

	// Mock the behavior for both email and username
	mockRepo.EXPECT().GetUserByEmail(newUser.Email).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.EXPECT().GetUserByUsername(newUser.Username).Return(&models.User{}, nil)

	// Call the service method
	createdUser, err := userSvc.CreateUser(newUser)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, createdUser)
	assert.Equal(t, "username already taken", err.Error())
}

func TestCreateUser_HashPasswordError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockHashError := errors.New("hashing error")

	// Create service with mocked hashing function
	serviceImpl := &services.UserServiceImpl{
		UserRepo: mockRepo,
		HashFunc: func(password string) (string, error) {
			return "", mockHashError // simulate hashing error
		},
	}
	userSvc := serviceImpl

	// Test data for new user
	newUser := &models.User{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
	}

	// Mock repository methods for checking email and username
	mockRepo.EXPECT().GetUserByEmail(newUser.Email).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.EXPECT().GetUserByUsername(newUser.Username).Return(nil, gorm.ErrRecordNotFound)

	// Mock CreateUser method (it shouldn't be called since hashing fails)
	mockRepo.EXPECT().CreateUser(gomock.Any()).Times(0) // Ensure CreateUser is not called

	// Call the service method
	createdUser, err := userSvc.CreateUser(newUser)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, createdUser)
	assert.Equal(t, mockHashError.Error(), err.Error()) // Check if the error is the expected hash error
}
