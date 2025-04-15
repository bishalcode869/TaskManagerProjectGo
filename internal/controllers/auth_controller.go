package controllers

import (
	"TaskManager/internal/models"
	"TaskManager/internal/services"
	"TaskManager/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AuthController handles authentication-related HTTP requests
type AuthController struct {
	AuthService services.AuthService
	Validator   *validator.Validate
}

// NewAuthController creates and returns a new AuthController instance
func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
		Validator:   validator.New(),
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

	// Validate input using the validator
	if err := a.Validator.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	// Generate JWT token
	token, err := utils.GenerateJWT(createdUser.ID, time.Hour*24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"token":   token,
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

	// Validate input using the validator
	if err := a.Validator.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Try to find user by username or email
	var user *models.User
	var err error
	if input.Email != "" {
		user, err = a.AuthService.GetUserByEmail(input.Email)
	} else if input.Username != "" {
		user, err = a.AuthService.GetUserByUsername(input.Username)
	}

	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/email or password"})
		return
	}

	// Compare stored hashed password with provided password
	if err := utils.ComparePasswords(user.Password, input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/email or password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, time.Hour*24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
