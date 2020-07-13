package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func (server *server) forgotPassword(c *gin.Context) {
	var args struct {
		Email string `json:"email" binding:"required,email"`
	}

	err := c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "forgotPassword",
			"subFunc": "c.ShouldBindJSON",
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	user, err := models.GetUserFromEmail(server.db, args.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// If user does not exist, it is not necessary to provide more info
			c.Status(http.StatusOK)
			return
		}
		log.WithFields(log.Fields{
			"func":    "forgotPassword",
			"subFunc": "models.GetUserFromEmail",
			"email":   args.Email,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "error fetching user details")
		return
	}

	tx := server.db.Begin()
	token, err := user.CreateForgotPasswordToken(tx)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "forgotPassword",
			"subFunc": "user.CreateForgotPasswordToken",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "error when creating token")
		return
	}

	err = server.email.SendForgotPasswordEmail(user, token)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "forgotPassword",
			"subFunc": "emails.SendForgotPasswordnEmail",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "error when seniding email")
		return
	}

	tx.Commit()
	c.Status(http.StatusOK)
	return
}
