// internal/controllers/user_controller.go
package controllers

import (
	"TaskManager/internal/models"
	"TaskManager/internal/services"
	dto "TaskManager/pkg/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController handles HTTP requests related to user operations
type UserController struct {
	UserService services.UserService
}

// NewUserController creates and returns a new UserController instance
func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// CreateUser handles creating a new user
func (u *UserController) CreateUser(c *gin.Context) {
	var userRequest dto.UserCreateRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		log.Println("Error binding user data:", err)
		return
	}

	// Basic field validation
	if userRequest.Email == "" || userRequest.Username == "" || userRequest.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Password is hashed in the service layer
	user := models.User{
		Email:    userRequest.Email,
		Username: userRequest.Username,
		Password: userRequest.Password,
	}

	newUser, err := u.UserService.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("Error creating user:", err)
		return
	}

	userResponse := dto.UserResponse{
		ID:        newUser.ID,
		Email:     newUser.Email,
		Username:  newUser.Username,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	c.JSON(http.StatusCreated, userResponse)
	log.Println("User created successfully:", newUser.Username)
}

// GetUserByID handles retrieving a user by their ID
func (u *UserController) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := u.UserService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userResponse := dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, userResponse)
}

// GetAllUsers handles retrieving all users
func (u *UserController) GetAllUsers(c *gin.Context) {
	users, err := u.UserService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, userResponses)
}

// UpdateUser handles updating an existing user
func (u *UserController) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var userRequest dto.UserCreateRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := u.UserService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if userRequest.Email != "" {
		user.Email = userRequest.Email
	}
	if userRequest.Username != "" {
		user.Username = userRequest.Username
	}

	updatedUser, err := u.UserService.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResponse := dto.UserResponse{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		Username:  updatedUser.Username,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}

	log.Printf("User updated successfully: %s (ID: %d)", updatedUser.Username, updatedUser.ID)
	c.JSON(http.StatusOK, userResponse)
}

// DeleteUser handles deleting a user by their ID
func (u *UserController) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = u.UserService.DeleteUser(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("User with ID %d deleted successfully", id)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
