package main

import (
	"net/http"

	"github.com/arunvm/focus/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) updateColumnCard(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updateColumnCard",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	var args models.UpdateColumnCardArgs
	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "updateColumnCard",
			"info":   "error decoding request body",
			"userID": user.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	args.ColumnID = c.GetString("boardColumnID")
	err = server.db.UpdateColumnCard(&args, user)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updateColumnCard",
			"subFunc": "user.UpdateColumnCard",
			"args":    args,
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when updating card details")
		return
	}

	c.Status(http.StatusOK)
	return
}
