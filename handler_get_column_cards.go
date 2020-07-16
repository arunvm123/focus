package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) getColumnCards(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getColumnCards",
			"subFunc": "getUserFromContext",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	cards, err := server.db.GetColumnCards(c.GetString("boardColumnID"))
	if err != nil {
		log.WithFields(log.Fields{
			"func":          "getColumnCards",
			"subFunc":       "models.GetColumnCards",
			"userID":        user.ID,
			"boardColumnID": c.GetString("boardColumnID"),
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error when retrieving column cards")
		return
	}

	c.JSON(http.StatusOK, cards)
	return
}
