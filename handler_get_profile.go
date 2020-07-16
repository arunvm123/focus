package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) getProfile(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getProfile",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	profile, err := server.db.GetProfile(user)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getProfile",
			"subFunc": "user.GetProfile",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching profile")
		return
	}

	c.JSON(http.StatusOK, profile)
	return
}
