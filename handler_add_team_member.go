package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) addTeamMember(c *gin.Context) {
	teamAdmin, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "addTeamMember",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	var args models.AddTeamMemberArgs
	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":        "addTeamMember",
			"info":        "error decoding request body",
			"teamAdminID": teamAdmin.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	args.TeamID = c.Keys["teamID"].(string)
	err = teamAdmin.AddTeamMember(server.db, args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":        "addTeamMember",
			"subFunc":     "teamAdmin.AddTeamMember",
			"teamAdminID": teamAdmin.ID,
			"args":        args,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when adding to team")
		return
	}

	c.Status(http.StatusOK)
	return
}
