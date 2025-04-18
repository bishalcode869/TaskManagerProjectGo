package services

import (
	"TaskManager/internal/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepo is a testify/mock for UserRepository
type MockUserRepo struct{ mock.Mock }

func (m *MockUserRepo) CreateUser(u *models.User) (*models.User, error) {
	args := m.Called(u)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserRepo) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserRepo) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserRepo) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}
func (m *MockUserRepo) UpdateUser(u *models.User) (*models.User, error) {
	args := m.Called(u)
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserRepo) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// Writing the Tests
func TestUserService_CreateUser_DuplilcateEmail(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)

	req := &models.User{Email: "e@mail", Username: "u", Password: "pass"}

	// Simulate email already taken
	mockRepo.
		On("GetUserByEmail", "e@mail").
		Return(&models.User{}, nil)

	_, err := svc.CreateUser(req)
	assert.EqualError(t, err, "email already taken")
	mockRepo.AssertNotCalled(t, "CreateUser")
}

func TestUserService_CreateUser_DuplicateUsername(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)

	req := &models.User{Email: "a@mail", Username: "u", Password: "pass"}

	// Email OK, but username taken
	mockRepo.
		On("GetUserByEmail", "a@mail").
		Return(nil, errors.New("not found"))
	mockRepo.
		On("GetUserByUsername", "u").
		Return(&models.User{}, nil)

	_, err := svc.CreateUser(req)
	assert.EqualError(t, err, "username already taken")
	mockRepo.AssertNotCalled(t, "CreateUser")
}

func TestUserService_CreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)

	req := &models.User{Email: "e@mail", Username: "u", Password: "pass"}

	// No duplicates
	mockRepo.
		On("GetUserByEmail", "e@mail").
		Return(nil, errors.New("not found"))
	mockRepo.
		On("GetUserByUsername", "u").
		Return(nil, errors.New("not found"))

	// Capture CreateUser call
	mockRepo.
		On("CreateUser", mock.AnythingOfType("*models.User")).
		Return(func(u *models.User) *models.User { return u }, nil)

	created, err := svc.CreateUser(req)
	assert.NoError(t, err)
	assert.NotEqual(t, "pass", created.Password, "password should be hashed")
	mockRepo.AssertExpectations(t)
}
