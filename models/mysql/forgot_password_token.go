package mysql

import (
	"time"

	"github.com/arunvm/travail-backend/models"
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (db *Mysql) CreateForgotPasswordToken(user *models.User) (string, error) {
	var token models.ForgotPasswordToken

	err := db.Client.Find(&token, "user_id = ?", user.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			token.UserID = user.ID
			token.Token = xid.New().String()
			token.ExpiresAt = time.Now().Unix() + int64(60*60)

			err = db.Client.Create(token).Error
			if err != nil {
				log.WithFields(log.Fields{
					"func":    "CreateForgotPasswordToken",
					"subFunc": "token.Create",
					"userID":  user.ID,
				}).Error(err)
				return "", err
			}
			return token.Token, nil
		}
		return "", err
	} else {
		token.Token = xid.New().String()
		token.ExpiresAt = time.Now().Unix() + int64(60*60)
		err = db.Client.Save(token).Error
		if err != nil {
			log.WithFields(log.Fields{
				"func":    "CreateForgotPasswordToken",
				"subFunc": "token.Save",
				"userID":  user.ID,
			}).Error(err)
			return "", err
		}
	}

	return token.Token, nil
}

func (db *Mysql) ResetPassword(token, password string) error {
	var user models.User

	tx := db.Client.Begin()

	err := tx.Table("forgot_password_tokens").Joins("JOIN users on forgot_password_tokens.user_id = users.id").
		Select("users.*").
		Find(&user, "token = ? AND expires_at > ?", token, time.Now().Unix()).Error
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func": "ResetPassword",
			"info": "retrieving user info if token exists and has not expired",
		}).Error(err)
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "ResetPassword",
			"subFunc": "bcrypt.GenerateFromPassword",
			"userID":  user.ID,
		}).Error(err)
		return err
	}

	user.Password = string(passwordHash)
	err = tx.Save(user).Error
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "ResetPassword",
			"subFunc": "user.Save",
			"userID":  user.ID,
		}).Error(err)
		return err
	}

	err = tx.Delete(models.ForgotPasswordToken{UserID: user.ID}).Error
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":   "ResetPassword",
			"info":   "deleting token",
			"userID": user.ID,
		}).Error(err)
		return err
	}

	tx.Commit()
	return nil
}
