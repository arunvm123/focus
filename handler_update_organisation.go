package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) updateOrganisation(c *gin.Context) {
	admin, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updateOrganisation",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var args models.UpdateOrganisationArgs
	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updateOrganisation",
			"subFunc": "c.ShouldBindJSON",
			"adminID": admin.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	args.ID = c.Keys["organisationID"].(string)
	err = server.db.UpdateOrganisation(&args, admin)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "updateOrganisation",
			"subFunc": "admin.UpdateOrganisation",
			"adminID": admin.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when updating organisation")
		return
	}

	c.Status(http.StatusOK)
	return
}
