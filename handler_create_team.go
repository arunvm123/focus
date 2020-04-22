package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) createTeam(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	var args models.CreateTeamArgs
	err = c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "createTeam",
			"info":   "error decoding request body",
			"userID": user.ID,
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	args.OrganisationID = c.Keys["organisationID"].(string)

	tx := server.db.Begin()
	err = user.CreateTeam(tx, &args)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "createTeam",
			"subFunc": "user.CreateTeam",
			"userID":  user.ID,
			"args":    args,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when creating team")
		return
	}

	tx.Commit()
	c.Status(http.StatusOK)
	return
}
