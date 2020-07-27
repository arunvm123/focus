package main

import (
	"net/http"

	"github.com/arunvm/focus/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) updateBoard(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updateBoard",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Error fetching user")
		return
	}

	var args models.UpdateBoardArgs
	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "updateBoard",
			"info":   "error decoding request body",
			"userID": user.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	args.TeamID = c.Keys["teamID"].(string)
	err = server.db.UpdateBoard(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updateBoard",
			"subFunc": "user.UpdateBoardArgs",
			"userID":  user.ID,
			"args":    args,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when updating board")
		return
	}

	c.Status(http.StatusOK)
	return
}
