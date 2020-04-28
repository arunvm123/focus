package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// The middleware checks if the user is an admin. So if it reaches this handler
// the user is an admin and can straight away send the response
func (server *server) adminCheck(c *gin.Context) {
	c.JSON(http.StatusOK, struct {
		Admin bool `json:"admin"`
	}{
		Admin: true,
	})
	return
}
