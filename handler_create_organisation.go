package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) createOrganisation(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createOrganisation",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var args models.CreateOrganisationArgs
	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createOrganisation",
			"subFunc": "c.ShouldBindJSON",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	tx := server.db.Begin()
	err = user.CreateOrganisation(tx, &args)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "createOrganisation",
			"subFunc": "user.CreateOrganisation",
			"userID":  user.ID,
			"args":    args,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when creating organisation")
		return
	}

	tx.Commit()
	c.Status(http.StatusOK)
	return
}
