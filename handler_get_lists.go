package main

import (
	"net/http"

	"github.com/arunvm/focus/models"
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

	var args models.GetListsArgs
	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getLists",
			"subFunc": "c.ShouldBindJSON",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	lists, err := server.db.GetLists(&args, user)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getLists",
			"subFunc": "user.GetLists",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when retrieving lists")
		return
	}

	if len(*lists) != 0 {
		(*lists)[0].Active = true
	}

	c.JSON(http.StatusOK, lists)
	return
}
