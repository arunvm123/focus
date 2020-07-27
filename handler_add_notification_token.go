package main

import (
	"net/http"

	"github.com/arunvm/focus/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) addNotificationToken(c *gin.Context) {
	var args models.AddNotificationTokenArgs

	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "addNotificationToken",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "addNotificationToken",
			"info":   "decoding request body",
			"userId": user.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	err = server.db.AddNotificationToken(&args, user)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "addNotificationToken",
			"subFunc": "models.AddNotificationToken",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when saving token")
		return
	}

	c.Status(http.StatusOK)
	return
}
