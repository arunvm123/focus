package main

import (
	"encoding/json"
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) updateProfile(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updateProfile",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var args models.UpdateProfileArgs
	err = json.NewDecoder(c.Request.Body).Decode(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "updateProfile",
			"info":   "error decoding request body",
			"userID": user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Request body not properly formatted")
		return
	}

	err = user.UpdateProfile(server.db, args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updateProfile",
			"subFunc": "user.UpdateProfile",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when updating prrofile")
		return
	}

	c.Status(http.StatusOK)
	return
}
