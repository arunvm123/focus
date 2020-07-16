package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) deleteTasks(c *gin.Context) {
	var args models.DeleteTasksArgs

	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "deleteTasks",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "deleteTasks",
			"subFunc": "c.ShouldBindJSON",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	if len(args.TaskIDs) == 0 {
		c.JSON(http.StatusBadRequest, "Please provide a non emoty array")
		return
	}

	err = server.db.DeleteTasks(&args, user)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "deleteTasks",
			"subFunc": "user.DeleteTasks",
			"args":    args,
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when deleting tasks")
		return
	}

	c.Status(http.StatusOK)
	return
}
