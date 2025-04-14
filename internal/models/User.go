package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"uniqueIndex"`
	Username string `json:"username" gorm:"uniqueIndex"`
	Password string `json:"-"`
}
