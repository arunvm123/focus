package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type ForgotPasswordToken struct {
	UserID    int    `json:"userID" gorm:"primary_key;auto_increment:false"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

func (fpt *ForgotPasswordToken) Create(db *gorm.DB) error {
	return db.Create(&fpt).Error
}

func (fpt *ForgotPasswordToken) Save(db *gorm.DB) error {
	return db.Save(&fpt).Error
}

func (user *User) CreateForgotPasswordToken(db *gorm.DB) (string, error) {
	var token ForgotPasswordToken

	err := db.Find(&token, "user_id = ?", user.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			token.UserID = user.ID
			token.Token = xid.New().String()
			token.ExpiresAt = time.Now().Unix() + int64(60*60)

			err = token.Create(db)
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
	}

	token.Token = xid.New().String()
	token.ExpiresAt = time.Now().Unix() + int64(60*60)
	err = token.Save(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "CreateForgotPasswordToken",
			"subFunc": "token.Save",
			"userID":  user.ID,
		}).Error(err)
		return "", err
	}

	return token.Token, nil
}

func ResetPassword(db *gorm.DB, token, password string) error {
	var user User

	err := db.Table("forgot_password_tokens").Joins("JOIN users on forgot_password_tokens.user_id = users.id").
		Select("users.*").
		Find(&user, "token = ? AND expires_at > ?", token, time.Now().Unix()).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func": "ResetPassword",
			"info": "retrieving user info if token exists and has not expired",
		}).Error(err)
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "ResetPassword",
			"subFunc": "bcrypt.GenerateFromPassword",
			"userID":  user.ID,
		}).Error(err)
		return err
	}

	user.Password = string(passwordHash)
	err = user.Save(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "ResetPassword",
			"subFunc": "user.Save",
			"userID":  user.ID,
		}).Error(err)
		return err
	}

	err = db.Delete(ForgotPasswordToken{UserID: user.ID}).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "ResetPassword",
			"info":   "deleting token",
			"userID": user.ID,
		}).Error(err)
		return err
	}

	return nil
}
