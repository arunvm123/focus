package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) createBoard(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createBoard",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var args models.CreateBoardArgs
	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createBoard",
			"subFunc": "c.ShouldBindJSON",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	args.TeamID = c.Keys["teamID"].(string)
	err = server.db.CreateBoard(&args, user)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createBoard",
			"subFunc": "user.CreateBoard",
			"userID":  user.ID,
			"args":    args,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when creating board")
		return
	}

	c.Status(http.StatusOK)
	return
}
