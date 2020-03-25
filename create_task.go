package main

import (
	"encoding/json"
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) createTask(c *gin.Context) {
	var task *models.Task

	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createTask",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	err = json.NewDecoder(c.Request.Body).Decode(&task)
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "createTask",
			"info":   "error decoding request body",
			"userID": user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Request body not properly formatted")
		return
	}

	err = user.CreateTask(server.db, task)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createTask",
			"subFunc": "user.CreateTask",
			"userID":  user.ID,
			"args":    task,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when creating task")
		return
	}

	c.JSON(http.StatusOK, task)
	return
}
