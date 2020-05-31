package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *server) deleteBoard(c *gin.Context) {
	// TODO: Need to discuss on who has permission to delete board

	c.Status(http.StatusNotImplemented)
	return
}
