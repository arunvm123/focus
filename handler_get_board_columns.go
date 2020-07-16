package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) getBoardColumns(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getBoardColumns",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	boardID := c.Query("boardID")
	columns, err := server.db.GetBoardColumns(boardID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getBoardColumns",
			"subFunc": "user.GetBoardColumns",
			"userID":  user.ID,
			"boardID": boardID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching board columns")
		return
	}

	c.JSON(http.StatusOK, &columns)
	return
}
