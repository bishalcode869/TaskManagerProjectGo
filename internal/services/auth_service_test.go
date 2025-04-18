package services

import (
	"TaskManager/internal/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_RegisterAndLogin(t *testing.T) {
	mockRepo := new(MockUserRepo)
	authSvc := NewAuthService(mockRepo)

	username, email, pass := "bob", "bob@example.com", "secret123"

	// RegisterUser mocks
	mockRepo.
		On("GetUserByUsername", username).
		Return(nil, errors.New("not found"))
	mockRepo.
		On("GetUserByEmail", email).
		Return(nil, errors.New("not found"))
	mockRepo.
		On("CreateUser", mock.AnythingOfType("*models.User")).
		Return(func(u *models.User) *models.User { return u }, nil)

	user, token, err := authSvc.RegisterUser(username, pass, email)
	assert.NoError(t, err)
	assert.Equal(t, username, user.Username)
	assert.NotEmpty(t, token)

	// LoginUser mocks
	mockRepo.
		On("GetUserByUsername", username).
		Return(user, nil)

	u2, tok2, err := authSvc.LoginUser(username, "", pass)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, u2.ID)
	assert.NotEmpty(t, tok2)
}

func TestAuthService_Login_Invalid(t *testing.T) {
	mockRepo := new(MockUserRepo)
	authSvc := NewAuthService(mockRepo)

	// No user found
	mockRepo.
		On("GetUserByUsername", "foo").
		Return(nil, errors.New("not found"))

	_, _, err := authSvc.LoginUser("foo", "", "p")
	assert.Error(t, err)
}
