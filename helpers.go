package main

import (
	"errors"
	"log"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
)

func getUserFromContext(c *gin.Context) (*models.User, error) {
	user, ok := c.Keys["user"].(*models.User)
	if !ok {
		log.Println("Unable to fetch user")
		return nil, errors.New("Error fetching user")
	}

	return user, nil
}
