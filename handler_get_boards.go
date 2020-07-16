package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) getBoards(c *gin.Context) {
	var args models.GetBoardsArgs
	args.TeamID = c.Keys["teamID"].(string)

	boards, err := server.db.GetBoards(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getBoards",
			"subFunc": "models.GetBoards",
			"args":    args,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when retrieving boards")
		return
	}

	c.JSON(http.StatusOK, &boards)
	return
}
