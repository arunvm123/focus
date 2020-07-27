package main

import (
	"net/http"

	"github.com/arunvm/focus/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) createTask(c *gin.Context) {
	var args models.CreateTaskArgs

	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createTask",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "createTask",
			"info":   "error decoding request body",
			"userID": user.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	task, err := server.db.CreateTask(&args, user)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createTask",
			"subFunc": "user.CreateTask",
			"userID":  user.ID,
			"args":    args,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when creating task")
		return
	}

	c.JSON(http.StatusOK, task)
	return
}
