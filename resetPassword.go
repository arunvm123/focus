package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/models"
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

	tx := server.db.Begin()
	err = models.ResetPassword(tx, args.Token, args.Password)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "resetPassword",
			"subFunc": "models.UpdatePassword",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "error while resetting password")
		return
	}

	tx.Commit()
	c.Status(http.StatusOK)
	return
}
