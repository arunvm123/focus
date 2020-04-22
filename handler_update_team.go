package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) updateTeam(c *gin.Context) {
	teamAdmin, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	var args models.UpdateTeamArgs
	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":        "updateTeam",
			"info":        "error decoding request body",
			"teamAdminID": teamAdmin.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	args.TeamID = c.Keys["teamID"].(string)
	err = teamAdmin.UpdateTeam(server.db, &args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":        "updateTeam",
			"subFunc":     "teamAdmin.UpdateTeam",
			"teamAdminID": teamAdmin.ID,
			"args":        args,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when updating team")
		return
	}

	c.Status(http.StatusOK)
	return
}
