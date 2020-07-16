package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) acceptOrganisationInvite(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "acceptOrganisationInvite",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var args models.AcceptOrganisationInviteArgs
	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "acceptOrganisationInvite",
			"subFunc": "c.ShouldBindJSON",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	err = server.db.AcceptOrganisationInvite(&args, user)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "acceptOrganisationInvite",
			"subFunc": "user.AcceptOrganisationInvite",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when acceptiong invite")
		return
	}

	c.Status(http.StatusOK)
	return
}
