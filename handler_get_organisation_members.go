package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) getOrganisationMembers(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getOrganisationMembers",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	members, err := server.db.GetOrganisationMembers(c.Keys["organisationID"].(string))
	if err != nil {
		log.WithFields(log.Fields{
			"func":           "getOrganisationMembers",
			"subFunc":        "models.GetOrganisationMembers",
			"userID":         user.ID,
			"organisationID": c.Keys["organisationID"].(string),
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching organisation members")
		return
	}

	c.JSON(http.StatusOK, members)
	return
}
