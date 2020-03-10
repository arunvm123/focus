package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *server) getProfile(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	profile, err := user.GetProfile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error fetching profile")
		return
	}

	c.JSON(http.StatusOK, profile)
	return
}
