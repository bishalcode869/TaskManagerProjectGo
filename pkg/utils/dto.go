// internal/dto/user.go
package utils

import "time"

// UserResponse defines the response structure for user data
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserCreateRequest defines the request structure for user creation
type UserCreateRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserUpdateRequest defines the request structure for updating user data
type UserUpdateRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}
