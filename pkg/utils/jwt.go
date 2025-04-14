package utils

import (
	"TaskManager/internal/config"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateJWT generates a JWT token for a given user ID with configurable expiration
func GenerateJWT(userID uint, expiration time.Duration) (string, error) {
	jwtSecret := []byte(config.Config.JWTSecret)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(expiration).Unix()

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return tokenString, nil
}

// ValidateToken validates the JWT token and returns the user ID
func ValidateToken(tokenString string) (uint, error) {
	jwtSecret := []byte(config.Config.JWTSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return 0, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id not found in token")
	}

	return uint(userID), nil
}
