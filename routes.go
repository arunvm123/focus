package main

import (
	"github.com/gin-gonic/gin"
)

func initialiseRoutes(server *server) *gin.Engine {
	r := gin.Default()

	r.POST("/signup", server.signup)
	r.POST("/login", server.login)

	return r
}
