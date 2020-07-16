package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) getNotificationTokens(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getNotificationTokens",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	tokens, err := server.db.GetNotificationTokens(user)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getNotificationTokens",
			"subFunc": "user.GetNotificationTokens",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when retrieving tokens")
		return
	}

	c.JSON(http.StatusOK, tokens)

}
