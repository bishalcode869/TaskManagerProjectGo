package utils

import (
	"TaskManager/internal/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mocking the JWTSecret for testing purposes
func init() {
	// override the secret key used for sining the JWT during tests
	config.Config = &config.AppConfig{
		JWTSecret: "testsecret",
	}

}

// Test function for GenerateJWT
func TestGenerateJWT(t *testing.T) {
	userID := uint(123)
	expiration := time.Hour * 24

	// Generate JWT token
	token, err := GenerateJWT(userID, expiration)
	// Assert no error
	require.NoError(t, err, "JWT should be generated without error")

	// Assert token is not empty
	require.NotEmpty(t, token, "Token should not be empty")

	// Parse and validate token
	parsedUserID, err := ValidateToken(token)
	require.NoError(t, err, "Token should be valid")

	// Assert the user ID from the token matches the one passed during generation
	assert.Equal(t, userID, parsedUserID, "User ID from token should match the generated one")

}

func TestValidateToken_InvalidToken(t *testing.T) {
	// Using an invalid token string
	invalidToken := "invalid_token_string"

	// Try to validate the invalid token
	_, err := ValidateToken(invalidToken)

	// Assert that an error is returned
	assert.Error(t, err, "Validation of an invalid token should return an error")
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	// Generate a token with a very short expiration
	userID := uint(123)
	expiration := time.Second * 1
	token, err := GenerateJWT(userID, expiration)
	require.NoError(t, err)

	// Wait for the token to expire
	time.Sleep(time.Second * 2)

	// Validate the expired token
	_, err = ValidateToken(token)

	// Assert that an error is returned for expired token
	assert.Error(t, err, "Validation of an expired token should return an error")
}
