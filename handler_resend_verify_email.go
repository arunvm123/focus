package main

import (
	"net/http"

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

	user, err := server.db.GetUserFromEmail(args.Email)
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

	err = server.db.CreateEmailValidationToken(user, server.email)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "resendVerifyEmail",
			"subFunc": "models.CreateEmailValidationToken",
			"userID":  user.ID,
		})
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
	return
}
