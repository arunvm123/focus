package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) createBoardColumn(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createBoardColumn",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	var args models.CreateBoardColumnArgs
	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "createBoardColumn",
			"info":   "error decoding request body",
			"userID": user.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	args.BoardID = c.GetString("boardID")
	err = user.CreateBoardColumn(server.db, &args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createBoardColumn",
			"subFunc": "user.CreateBoardColumn",
			"userID":  user.ID,
			"args":    args,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when creating board column")
		return
	}

	c.Status(http.StatusOK)
	return
}
