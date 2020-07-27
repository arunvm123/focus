package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) signup(c *gin.Context) {
	var args models.SignUpArgs
	err := c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "signup",
			"info": "decoding request body",
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	if server.db.CheckIfUserExists(args.Email) == true {
		c.JSON(http.StatusConflict, "Email already exists")
		return
	}

	tx := server.tx.Begin()
	user, token, err := tx.UserSignup(&args, false)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "signup",
			"subFunc": "models.UserSignup",
			"email":   args.Email,
		})
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = server.email.SendValidationEmail(user.Name, user.Email, token)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "UserSignup",
			"subFunc": "emailClient.SendValidationEmail",
			"userID":  user.ID,
		})
		c.JSON(http.StatusInternalServerError, "Error while sending email")
		return
	}

	tx.Commit()
	c.Status(http.StatusOK)
	return
}
