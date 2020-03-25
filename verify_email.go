package main

import (
	"encoding/json"
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) verifyEmail(c *gin.Context) {
	var token models.ValidateEmailArgs

	err := json.NewDecoder(c.Request.Body).Decode(&token)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "verifyEmail",
			"info": "error decoding request body",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Request body not properly formatted")
		return
	}

	err = models.VerifyEmail(server.db, token.Token)
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
