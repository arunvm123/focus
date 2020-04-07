package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) updateList(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updateList",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var args models.UpdateListArgs
	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "updateList",
			"info":   "error decoding request body",
			"userID": user.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	err = user.UpdateList(server.db, &args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updateList",
			"subFunc": "user.UpdateList",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
	return
}
