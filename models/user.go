package models

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	ID          int     `json:"id" gorm:"primary_key"`
	Email       string  `json:"email" gorm:"unique;not null"`
	Name        string  `json:"name"`
	Password    string  `json:"password"`
	Verified    bool    `json:"verified"`
	ProfilePic  *string `json:"profilePic"`
	GoogleOauth bool    `json:"googleOauth"`
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
	ID         int     `json:"id"`
	Email      string  `json:"email"`
	Name       string  `json:"name"`
	ProfilePic *string `json:"profilePic"`
}

type SignUpArgs struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginArgs struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginWithGoogleArgs struct {
	Email   string  `json:"email" binding:"required,email"`
	Name    string  `json:"name" binding:"required"`
	Picture *string `json:"picture"`
}

type UpdateProfileArgs struct {
	// Email    *string `json:"email,omitempty"`
	Name     *string `json:"name,omitempty"`
	Password *string `json:"password,omitempty"`
}

type UpdatePasswordArgs struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}

// GetProfile organises user data and returns it
func (user *User) GetProfile() (*UserProfile, error) {
	return &UserProfile{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		ProfilePic: user.ProfilePic,
	}, nil
}

func (user *User) UpdateProfile(db *gorm.DB, args UpdateProfileArgs) error {
	if args.Name != nil {
		user.Name = *args.Name
	}
	if args.Password != nil {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(*args.Password), bcrypt.DefaultCost)
		if err != nil {
			log.WithFields(log.Fields{
				"func":    "UpdateProfile",
				"subFunc": "bcrypt.GenerateFromPassword",
				"userID":  user.ID,
			}).Error(err)
			return err
		}

		user.Password = string(passwordHash)
	}

	err := user.Save(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UpdateProfile",
			"subFunc": "user.Save",
			"userID":  user.ID,
		}).Error(err)
		return err
	}

	return nil
}

// GetUserFromEmail returns user details from the given email id
func GetUserFromEmail(db *gorm.DB, email string) (*User, error) {
	var user User

	err := db.Find(&user, "email = ?", email).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":  "GetUserFromEmail",
			"info":  "retrieving user info from email",
			"email": email,
		}).Error(err)
		return nil, err
	}

	return &user, nil
}

// GetUserFromID returns user details from the given user id
func GetUserFromID(db *gorm.DB, userID int) (*User, error) {
	var user User

	err := db.Find(&user, "id = ?", userID).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "GetUserFromID",
			"info":   "retrieving user info from id",
			"userID": userID,
		}).Error(err)
		return nil, err
	}

	return &user, nil
}

func CheckIfUserExists(db *gorm.DB, email string) bool {
	var count int
	err := db.Table("users").Where("email = ?", email).Count(&count).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":  "CheckIfUserExists",
			"info":  "checking if user with specified email exitst",
			"email": email,
		}).Error(err)
		return true
	}

	if count > 0 {
		return true
	}

	return false
}

func UserSignup(db *gorm.DB, args *SignUpArgs) (*User, error) {
	var user User

	user.Email = args.Email
	user.Name = args.Name

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UserSignup",
			"subFunc": "bcrypt.GenerateFromPassword",
			"email":   args.Email,
		}).Error(err)
		return nil, err
	}

	user.Password = string(passwordHash)
	user.Verified = false
	user.GoogleOauth = false

	err = user.Create(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UserSignup",
			"subFunc": "user.Create",
			"email":   args.Email,
		}).Error(err)
		return nil, err
	}

	return &user, nil
}

func SignUpWithGoogle(db *gorm.DB, args *LoginWithGoogleArgs) (*User, error) {
	var user User

	user.Email = args.Email
	user.Name = args.Name
	user.ProfilePic = args.Picture
	user.Verified = true
	user.GoogleOauth = true

	err := user.Create(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "SignUpWithGoogle",
			"subFunc": "user.Create",
			"email":   args.Email,
		}).Error(err)
		return nil, err
	}

	return &user, nil
}

func (user *User) UpdatePassword(db *gorm.DB, args *UpdatePasswordArgs) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(args.CurrentPassword))
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UpdatePassword",
			"subFunc": "bcrypt.CompareHashAndPassword",
			"userID":  user.ID,
		}).Error(err)
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(args.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UpdatePassword",
			"subFunc": "bcrypt.GenerateFromPassword",
			"userID":  user.ID,
		}).Error(err)
		return err
	}

	user.Password = string(hashedPassword)
	err = user.Save(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UpdatePassword",
			"subFunc": "user.Save",
			"userID":  user.ID,
		}).Error(err)
		return err
	}

	return nil
}
