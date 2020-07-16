package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) getBugs(c *gin.Context) {
	admin, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createBug",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	bugs, err := server.db.GetBugs()
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createBug",
			"subFunc": "admin.GetBugs",
			"adminID": admin.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching bugs")
		return
	}

	c.JSON(http.StatusOK, bugs)
	return
}
