package models

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type FCMNotificationToken struct {
	UserID int    `json:"userID" gorm:"primary_key;auto_increment:false"`
	Token  string `json:"token" gorm:"primary_key"`
}

// Create is a helper function to create details of email validation token
func (token *FCMNotificationToken) Create(db *gorm.DB) error {
	return db.Create(&token).Error
}

// Save is a helper function to update details of email validation token
func (token *FCMNotificationToken) Save(db *gorm.DB) error {
	return db.Save(&token).Error
}

type AddNotificationTokenArgs struct {
	Token string `json:"token" binding:"required"`
}

func (user *User) AddNotificationToken(db *gorm.DB, args *AddNotificationTokenArgs) error {
	var fcmToken FCMNotificationToken

	fcmToken.Token = args.Token
	fcmToken.UserID = user.ID

	err := fcmToken.Save(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "AddNotificationToken",
			"info":   "saving notification token",
			"userID": user.ID,
		})
		return err
	}

	return nil
}

func (user *User) GetNotificationTokens(db *gorm.DB) ([]string, error) {
	var tokens []string

	err := db.Table("fcm_notification_tokens").Where("user_id = ?", user.ID).
		Pluck("token", &tokens).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "GetNotificationTokens",
			"info":   "retrieving notification tokens",
			"userID": user.ID,
		}).Error(err)
		return []string{}, err
	}

	return tokens, nil
}
