package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *server) getLists(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	lists, err := user.GetLists(server.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, lists)
	return
}
