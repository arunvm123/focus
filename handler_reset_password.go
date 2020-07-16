package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) resetPassword(c *gin.Context) {
	var args struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	err := c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "resetPassword",
			"subFunc": "c.ShouldBindJSON",
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	err = server.db.ResetPassword(args.Token, args.Password)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "resetPassword",
			"subFunc": "models.UpdatePassword",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "error while resetting password")
		return
	}

	c.Status(http.StatusOK)
	return
}
