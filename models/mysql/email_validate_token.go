package mysql

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/arunvm/travail-backend/email"
	"github.com/arunvm/travail-backend/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func (db *Mysql) CreateEmailValidationToken(user *models.User, emailClient email.Email) error {

	tx := db.Client.Begin()

	token, err := emailValidateToken(tx, user)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "CreateEmailValidationToken",
			"subFunc": "emailValidateToken",
		}).Error(err)
		return err
	}

	err = emailClient.SendValidationEmail(user.Name, user.Email, token)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "CreateEmailValidationToken",
			"subFunc": "emailClient.SendValidationEmail",
			"userID":  user.ID,
		})
		return err
	}

	tx.Commit()
	return nil
}

func emailValidateToken(db *gorm.DB, user *models.User) (string, error) {
	err := db.Table("email_validate_tokens").Where("user_id = ?", user.ID).
		UpdateColumn("invalidate", true).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.WithFields(log.Fields{
				"func":   "InvalidateEmailTokens",
				"info":   "updating all existing tokens as invalidated",
				"userID": user.ID,
			}).Error(err)
			return "", err
		}
	}

	createdAt := time.Now().Unix()

	ev := &models.EmailValidateToken{
		UserID:     user.ID,
		CreatedAt:  createdAt,
		ExpiresAt:  createdAt + int64(24*60*60),
		Token:      randToken(),
		Invalidate: false,
	}

	err = db.Create(ev).Error
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

func (db *Mysql) VerifyEmail(token string) error {
	var ev models.EmailValidateToken

	err := db.Client.Find(&ev, "token = ? AND invalidate = false AND expires_at > ?", token, time.Now().Unix()).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func": "VerifyEmail",
			"info": "error retrieving token from db",
		}).Error(err)
		return err
	}

	err = db.Client.Table("users").Where("id = ?", ev.UserID).UpdateColumn("verified", true).Error
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

func (db *Mysql) InvalidateEmailTokens(userID int) error {
	err := db.Client.Table("email_validate_tokens").Where("user_id = ?", userID).
		UpdateColumn("invalidate", true).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "InvalidateEmailTokens",
			"info":   "updating all existing tokens as invalidated",
			"userID": userID,
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
