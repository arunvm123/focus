package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
)

func (server *server) verifyEmail(c *gin.Context) {
	var token models.ValidateEmailArgs

	err := json.NewDecoder(c.Request.Body).Decode(&token)
	if err != nil {
		log.Printf("Error when decoding request body\n%v", err)
		c.JSON(http.StatusInternalServerError, "Request body not properly formatted")
		return
	}

	err = models.VerifyEmail(server.db, token.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error verifying email")
		return
	}

	c.Status(http.StatusOK)
	return
}
