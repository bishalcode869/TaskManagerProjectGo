package services_test

import (
	"TaskManager/internal/models"
	"TaskManager/internal/services"
	"TaskManager/mocks"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	svc := &services.AuthServiceImpl{
		AuthRepo:     mockRepo,
		HashPassword: func(pw string) (string, error) { return "hashedPw", nil },
		GenerateJWT:  func(id uint, ttl time.Duration) (string, error) { return "signedToken", nil },
		TokenTTL:     time.Hour,
	}

	mockRepo.EXPECT().GetUserByUsername("john").Return(nil, gorm.ErrRecordNotFound)
	mockRepo.EXPECT().GetUserByEmail("john@example.com").Return(nil, gorm.ErrRecordNotFound)
	mockRepo.EXPECT().CreateUser(gomock.Any()).DoAndReturn(
		func(u *models.User) (*models.User, error) { return u, nil },
	)

	user, token, err := svc.RegisterUser("john", "pass123", "john@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "signedToken", token)
	assert.Equal(t, "john", user.Username)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, "hashedPw", user.Password)
}

func TestRegisterUser_DuplicateUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	svc := &services.AuthServiceImpl{AuthRepo: mockRepo}

	mockRepo.EXPECT().GetUserByUsername("john").Return(&models.User{Username: "john"}, nil)

	user, token, err := svc.RegisterUser("john", "pass123", "john@example.com")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	assert.Equal(t, "username already taken", err.Error())
}

func TestRegisterUser_DuplicateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	svc := &services.AuthServiceImpl{AuthRepo: mockRepo}

	mockRepo.EXPECT().GetUserByUsername("john").Return(nil, gorm.ErrRecordNotFound)
	mockRepo.EXPECT().GetUserByEmail("john@example.com").Return(&models.User{Email: "john@example.com"}, nil)

	user, token, err := svc.RegisterUser("john", "pass123", "john@example.com")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	assert.Equal(t, "email already registered", err.Error())
}

func TestRegisterUser_HashError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	svc := &services.AuthServiceImpl{
		AuthRepo:     mockRepo,
		HashPassword: func(pw string) (string, error) { return "", errors.New("hash fail") },
	}

	mockRepo.EXPECT().GetUserByUsername("john").Return(nil, gorm.ErrRecordNotFound)
	mockRepo.EXPECT().GetUserByEmail("john@example.com").Return(nil, gorm.ErrRecordNotFound)

	user, token, err := svc.RegisterUser("john", "pass123", "john@example.com")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	assert.Equal(t, "failed to hash password", err.Error())
}

func TestLoginUser_Success_ByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	svc := &services.AuthServiceImpl{
		AuthRepo:        mockRepo,
		ComparePassword: func(hash, pw string) error { return nil },
		GenerateJWT:     func(id uint, ttl time.Duration) (string, error) { return "jwtToken", nil },
		TokenTTL:        time.Hour,
	}

	stored := &models.User{Username: "john", Password: "hashed"}
	mockRepo.EXPECT().GetUserByEmail(gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	mockRepo.EXPECT().GetUserByUsername("john").Return(stored, nil)

	user, token, err := svc.LoginUser("john", "", "pass123")
	assert.NoError(t, err)
	assert.Equal(t, stored, user)
	assert.Equal(t, "jwtToken", token)
}

func TestLoginUser_Fail_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	svc := &services.AuthServiceImpl{AuthRepo: mockRepo}

	mockRepo.EXPECT().GetUserByEmail(gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	mockRepo.EXPECT().GetUserByUsername("john").Return(nil, gorm.ErrRecordNotFound)

	user, token, err := svc.LoginUser("john", "", "pass123")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
}

func TestLoginUser_Fail_BadPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	stored := &models.User{Username: "john", Password: "hashed"}
	svc := &services.AuthServiceImpl{
		AuthRepo:        mockRepo,
		ComparePassword: func(hash, pw string) error { return errors.New("mismatch") },
	}

	mockRepo.EXPECT().GetUserByEmail(gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	mockRepo.EXPECT().GetUserByUsername("john").Return(stored, nil)

	user, token, err := svc.LoginUser("john", "", "pass123")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
}

func TestLoginUser_JWTError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	stored := &models.User{Email: "john@example.com", Password: "hashed"}
	svc := &services.AuthServiceImpl{
		AuthRepo:        mockRepo,
		ComparePassword: func(hash, pw string) error { return nil },
		GenerateJWT:     func(id uint, ttl time.Duration) (string, error) { return "", errors.New("jwt fail") },
	}

	mockRepo.EXPECT().GetUserByEmail("john@example.com").Return(stored, nil)

	user, token, err := svc.LoginUser("", "john@example.com", "pass123")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
}
