package controllers

import (
	"TaskManager/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthController handles authentication-related HTTP requests
type AuthController struct {
	AuthService services.AuthService
}

// NewAuthController creates and returns a new AuthController instance
func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

// Register registers a new user and returns a JWT token
func (a *AuthController) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	// Bind JSON data to input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call the service layer to register the user
	createdUser, token, err := a.AuthService.RegisterUser(input.Username, input.Password, input.Email)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"token":   token,
		"user":    createdUser,
	})
}

// Login authenticates a user and returns a JWT token
func (a *AuthController) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password" binding:"required"`
	}

	// Bind JSON data to input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Authenticate the user using the service layer
	user, token, err := a.AuthService.LoginUser(input.Username, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/email or"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"message":  "Login successful",
		"token":    token,
	})
}
