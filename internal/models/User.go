package models

import "gorm.io/gorm"

// User represents a user in the system
type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"uniqueIndex;not null" `
	Username string `json:"username" gorm:"uniqueIndex;not null" `
	Password string `json:"-" `
}
