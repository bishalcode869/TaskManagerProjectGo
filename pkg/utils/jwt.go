package utils

import (
	"TaskManager/internal/config"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// jwtSecret should be set in an environment variable, for example
var jwtSecret = []byte(config.Config.JWTSecret) // Store JWT secret securely, not directly in the code

// GenerateJWT generates a JWT token for a given user ID with configurable expiration
func GenerateJWT(userID uint, expiration time.Duration) (string, error) {
	// Create a new JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(expiration).Unix() // Token expires in the given duration

	// Generate the token string
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return tokenString, nil
}

// ValidateToken validates the JWT token and returns the user ID
func ValidateToken(tokenString string) (uint, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method is valid
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return 0, fmt.Errorf("invalid token: %v", err)
	}

	// Extract the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// Get the user ID from claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id not found in token")
	}

	return uint(userID), nil
}
