package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) updatePassword(c *gin.Context) {
	var args models.UpdatePasswordArgs

	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updatePassword",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	if user.GoogleOauth {
		c.JSON(http.StatusUnauthorized, "Cannot change password for oauth accounts")
	}

	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updatePassword",
			"subFunc": "c.ShouldBindJSON",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	err = user.UpdatePassword(server.db, &args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updatePassword",
			"subFunc": "user.UpdatePassword",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error while updating passwerd")
		return
	}

	c.Status(http.StatusOK)
	return
}
