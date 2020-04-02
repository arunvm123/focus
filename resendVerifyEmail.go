package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/emails"
	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func (server *server) resendVerifyEmail(c *gin.Context) {
	var args struct {
		Email string `json:"email" binding:"required,email"`
	}

	err := c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "resendVerifyEmail",
			"subFunc": "c.ShouldBindJSON",
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	user, err := models.GetUserFromEmail(server.db, args.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, "Please sign up")
			return
		}
		log.WithFields(log.Fields{
			"func":    "resendVerifyEmail",
			"subFunc": "models.GetUserFromEmail",
			"email":   args.Email,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "error fetching user details")
		return
	}

	if user.Verified {
		c.JSON(http.StatusBadRequest, "user already verified")
		return
	}

	tx := server.db.Begin()
	err = models.InvalidateEmailTokens(tx, user.ID)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "resendVerifyEmail",
			"subFunc": "models.InvalidateEmailTokens",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "error while invalidating previous tokens")
		return
	}

	token, err := models.CreateEmailValidationToken(tx, user)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "resendVerifyEmail",
			"subFunc": "models.CreateEmailValidationToken",
			"userID":  user.ID,
		})
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = emails.SendValidationEmail(server.email, user, token)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "resendVerifyEmail",
			"subFunc": "emails.SendValidationEmail",
			"userID":  user.ID,
		})
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	tx.Commit()
	c.Status(http.StatusOK)
	return
}
