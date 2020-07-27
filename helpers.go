package main

import (
	"errors"

	"github.com/arunvm/focus/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func getUserFromContext(c *gin.Context) (*models.User, error) {
	user, ok := c.Keys["user"].(*models.User)
	if !ok {
		log.WithFields(log.Fields{
			"func": "getUserFromContext",
			"info": "retrieving user info from context",
		}).Error(errors.New("Error while retrieving user info from context"))
		return nil, errors.New("Error fetching user")
	}

	return user, nil
}
