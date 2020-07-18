package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) inviteToOrganisation(c *gin.Context) {
	admin, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "inviteToOrganisation",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var args models.InviteToOrganisationArgs
	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "inviteToOrganisation",
			"subFunc": "c.ShouldBindJSON",
			"adminID": admin.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	if args.Email == admin.Email {
		c.JSON(http.StatusBadRequest, "Provide another email id")
		return
	}

	args.OrganisationID = c.Keys["organisationID"].(string)

	tx := server.db.Begin()
	orgInfo, err := tx.CreateOrganisationInviteToken(&args, admin)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "inviteToOrganisation",
			"subFunc": "admin.CreateOrganisationInviteToken",
			"adminID": admin.ID,
			"args":    args,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when creating invitation token")
		return
	}

	err = server.email.SendOrganisationInvite(admin.Name, orgInfo.Email, orgInfo.Token, orgInfo.OrganisationName)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "inviteToOrganisation",
			"subFunc": "emails.SendInviteToOrganisation",
			"adminID": admin.ID,
			"args":    args,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when sending organisation invite")
		return
	}

	tx.Commit()
	c.Status(http.StatusOK)
	return
}
