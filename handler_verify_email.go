package main

import (
	"net/http"

	"github.com/arunvm/focus/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) verifyEmail(c *gin.Context) {
	var token models.ValidateEmailArgs

	err := c.ShouldBindJSON(&token)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "verifyEmail",
			"info": "error decoding request body",
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	err = server.db.VerifyEmail(token.Token)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "verifyEmail",
			"subFunc": "models.VerifyEmail",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error verifying email")
		return
	}

	c.Status(http.StatusOK)
	return
}
