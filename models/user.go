package models

import (
	"log"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
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

// Save is a helper function to update user
func (user *User) Save(db *gorm.DB) error {
	return db.Save(&user).Error
}

// UserProfile describes users profile
type UserProfile struct {
	ID    int    `json:"id" gorm:"primary_key"`
	Email string `json:"email" gorm:"unique;not null"`
	Name  string `json:"name"`
}

type UpdateProfileArgs struct {
	// Email    *string `json:"email,omitempty"`
	Name     *string `json:"name,omitempty"`
	Password *string `json:"password,omitempty"`
}

// GetProfile organises user data and returns it
func (user *User) GetProfile() (*UserProfile, error) {
	return &UserProfile{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

func (user *User) UpdateProfile(db *gorm.DB, args UpdateProfileArgs) error {
	if args.Name != nil {
		user.Name = *args.Name
	}
	if args.Password != nil {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(*args.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error when hashing password\n%v", err)
			return err
		}

		user.Password = string(passwordHash)
	}

	err := user.Save(db)
	if err != nil {
		log.Printf("Error when updating user\n%v", err)
		return err
	}

	return nil
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

// GetUserFromID returns user details from the given user id
func GetUserFromID(db *gorm.DB, userID int) (*User, error) {
	var user User

	err := db.Find(&user, "id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
