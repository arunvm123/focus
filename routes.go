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

	private.POST("/create/list", server.createList)
	private.POST("/get/lists", server.getLists)
	private.POST("/update/list", server.updateList)

	private.POST("/create/tasks", server.createTasks)
	private.POST("/create/task", server.createTask)
	private.POST("/get/tasks", server.getTasks)
	private.POST("/update/task", server.updateTask)

	private.GET("/get/profile", server.getProfile)
	private.POST("/update/profile", server.updateProfile)

	return r
}
