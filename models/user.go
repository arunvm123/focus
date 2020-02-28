package models

import (
	"github.com/jinzhu/gorm"
)

// User model
type User struct {
	ID       int    `json:"id" gorm:"primary_key"`
	Email    string `json:"email" gorm:"unique;not null"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Create is a helper function to create a new user
func (user *User) Create(db *gorm.DB) error {
	return db.Create(&user).Error
}

// GetUserFromEmail returns user details from the given email id
func GetUserFromEmail(db *gorm.DB, email string) (*User, error) {
	var user User

	err := db.Find(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
