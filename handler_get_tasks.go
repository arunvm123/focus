package main

import (
	"net/http"

	"github.com/arunvm/focus/models"
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
	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "getTasks",
			"info":   "Error when decoding request body",
			"userID": user.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	tasks, err := server.db.GetTasks(&args, user)
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
