package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/arunvm/travail-backend/models"

	"github.com/gin-gonic/gin"
)

func (server *server) signup(c *gin.Context) {
	var user models.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		panic(err)
	}

	log.Println(user)
	err = user.Create(server.db)
	if err != nil {
		panic(err)
	}

	c.Status(http.StatusOK)
	return
}
