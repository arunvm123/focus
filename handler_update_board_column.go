package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) updateBoardColumn(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updateBoardColumn",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	var args models.UpdateBoardColumnArgs
	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "updateBoardColumn",
			"info":   "error decoding request body",
			"userID": user.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	args.BoardID = c.GetString("boardID")
	err = models.UpdateBoardColumn(server.db, &args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updateBoardColumn",
			"subFuuc": "user.UpdateBoardColumn",
			"args":    args,
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when updating board column")
		return
	}

	c.Status(http.StatusOK)
	return
}
