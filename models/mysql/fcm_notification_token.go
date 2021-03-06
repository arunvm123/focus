package mysql

import (
	"github.com/arunvm/focus/models"
	log "github.com/sirupsen/logrus"
)

func (db *Mysql) AddNotificationToken(args *models.AddNotificationTokenArgs, user *models.User) error {
	var fcmToken models.FCMNotificationToken

	fcmToken.Token = args.Token
	fcmToken.UserID = user.ID

	err := db.Client.Save(fcmToken).Error
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

func (db *Mysql) GetNotificationTokens(userId int) ([]string, error) {
	var tokens []string

	err := db.Client.Table("fcm_notification_tokens").Where("user_id = ?", userId).
		Pluck("token", &tokens).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "GetNotificationTokens",
			"info":   "retrieving notification tokens",
			"userID": userId,
		}).Error(err)
		return []string{}, err
	}

	return tokens, nil
}
