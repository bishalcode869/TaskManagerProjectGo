package models

import "gorm.io/gorm"

// User represents a user in the system
type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Username string `json:"username" gorm:"uniqueIndex;not null" validate:"required,min=3,max=30"`
	Password string `json:"-" validate:"required,min=8"`
}
