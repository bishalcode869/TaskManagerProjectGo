package controllers

import (
	"TaskManager/internal/models"
	"TaskManager/internal/services"
	"TaskManager/pkg/utils"
	"net/http"
	"time"

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

	// Check if username already exists using AuthService
	existingUser, err := a.AuthService.GetUserByUsername(input.Username)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	// Check if email already exists using AuthService
	existingUserByEmail, err := a.AuthService.GetUserByEmail(input.Email)
	if err == nil && existingUserByEmail != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password using the utility function
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create a new user model instance (password will be hashed)
	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	}

	// Use AuthService to create the user
	createdUser, err := a.AuthService.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Set token expiration (default 24 hours, configurable via query parameter)
	expirationTime := time.Hour * 24
	if expTimeStr := c.DefaultQuery("JWT_EXPIRATION_HOURS", "24"); expTimeStr != "" {
		if expTime, err := time.ParseDuration(expTimeStr + "h"); err == nil {
			expirationTime = expTime
		}
	}

	// Generate JWT token using the created user's ID and the expiration time
	token, err := utils.GenerateJWT(createdUser.ID, expirationTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the JWT token in the response
	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"token":   token,
	})
}

// Login authenticates a user and returns a JWT token
func (a *AuthController) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Bind JSON data to input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Use AuthService to retrieve the user by username
	user, err := a.AuthService.GetUserByUsername(input.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Compare stored hashed password with provided password
	if err := utils.ComparePasswords(user.Password, input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token with default 24-hour expiry
	token, err := utils.GenerateJWT(user.ID, time.Hour*24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the JWT token
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
