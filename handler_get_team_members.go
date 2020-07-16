package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) getTeamMembers(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getTeamMembers",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	members, err := server.db.GetTeamMembers(c.Keys["teamID"].(string))
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getTeamMembers",
			"subFunc": "models.GetTeamMembers",
			"userID":  user.ID,
			"teamID":  c.Keys["teamID"].(string),
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when retrieving team members")
		return
	}

	c.JSON(http.StatusOK, members)
	return
}
