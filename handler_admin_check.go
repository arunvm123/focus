package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *server) adminCheck(c *gin.Context) {
	c.JSON(http.StatusOK, struct {
		Admin bool `json:"admin"`
	}{
		Admin: true,
	})
	return
}
