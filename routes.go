package main

import (
	"github.com/gin-gonic/gin"
)

func initialiseRoutes(server *server) *gin.Engine {
	r := gin.Default()

	public := r.Group("/")
	public.POST("/signup", server.signup)
	public.POST("/login", server.login)

	private := r.Group("/")
	private.Use(server.tokenAuthorisationMiddleware())
	private.POST("/create/tasks", server.createTasks)
	private.POST("/create/task", server.createTask)

	return r
}
