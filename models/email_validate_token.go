package models

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type EmailValidateToken struct {
	UserID    int    `json:"uesrId" gorm:"primary_key"`
	Token     string `json:"token"`
	CreatedAt int64  `json:"createdAt"`
	ExpiresAt int64  `json:"expiresAt"`
}

// Create is a helper function to create details of email validation token
func (ev *EmailValidateToken) Create(db *gorm.DB) error {
	return db.Create(&ev).Error
}

// Save is a helper function to update details of email validation token
func (ev *EmailValidateToken) Save(db *gorm.DB) error {
	return db.Save(&ev).Error
}

type ValidateEmailArgs struct {
	Token string `json:"token"`
}

func CreateEmailValidationToken(db *gorm.DB, user *User) (string, error) {
	createdAt := time.Now().Unix()

	ev := &EmailValidateToken{
		UserID:    user.ID,
		CreatedAt: createdAt,
		ExpiresAt: createdAt + int64(time.Hour*24*60),
		Token:     randToken(),
	}

	err := ev.Create(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "CreateEmailValidationToken",
			"subFunc": "ev.Create",
			"userID":  user.ID,
		}).Error(err)
		return "", err
	}

	return ev.Token, nil
}

func VerifyEmail(db *gorm.DB, token string) error {
	var ev EmailValidateToken

	err := db.Find(&ev, "token = ?", token).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func": "VerifyEmail",
			"info": "error retrieving token from db",
		}).Error(err)
		return err
	}

	err = db.Table("users").Where("id = ?", ev.UserID).UpdateColumn("verified", true).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "VerifyEmail",
			"info":   "setting email verified as true",
			"userID": ev.UserID,
		}).Error(err)
		return err
	}

	return nil
}

func randToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}