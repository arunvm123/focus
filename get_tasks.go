package main

import (
	"encoding/json"
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func (server *server) getTasks(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getTasks",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	var args models.GetTasksArgs
	err = json.NewDecoder(c.Request.Body).Decode(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "getTasks",
			"info":   "Error when decoding request body",
			"userID": user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Request body not properly formatted")
		return
	}

	tasks, err := user.GetTasks(server.db, &args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getTasks",
			"subFunc": "user.GetTasks",
			"userID":  user.ID,
			"args":    args,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching tasks")
		return
	}

	c.JSON(http.StatusOK, tasks)
	return
}
