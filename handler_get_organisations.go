package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) getOrganisations(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getOrganisations",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	organisations, err := server.db.GetOrganisations(user)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getOrganisations",
			"subFunc": "user.GetOrganisations",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching organisations")
		return
	}

	c.JSON(http.StatusOK, organisations)
	return
}
