package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) getLists(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getLists",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	lists, err := user.GetLists(server.db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getLists",
			"subFunc": "user.GetLists",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, lists)
	return
}
