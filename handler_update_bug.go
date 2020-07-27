package main

import (
	"net/http"

	"github.com/arunvm/focus/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) updateBug(c *gin.Context) {
	var args models.UpdateBugArgs

	admin, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updateBug",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updateBug",
			"info":    "error decoding request body",
			"adminID": admin.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	err = server.db.UpdateBug(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updateBug",
			"subFunc": "admin.UpdateBug",
			"adminID": admin.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when updating bug")
		return
	}

	c.Status(http.StatusOK)
	return
}
