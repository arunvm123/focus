package main

import (
	"github.com/gin-gonic/gin"
)

func initialiseRoutes(server *server) *gin.Engine {
	r := gin.Default()

	public := r.Group("/")
	public.POST("/signup", server.signup)
	public.POST("/login", server.login)
	public.POST("/google/login", server.loginWithGoogle)
	public.POST("/verify/email", server.verifyEmail)
	public.POST("/resend/verify/email", server.resendVerifyEmail)

	private := r.Group("/")
	private.Use(server.tokenAuthorisationMiddleware())

	private.POST("/create/list", server.createList)
	private.POST("/get/lists", server.getLists)
	private.POST("/update/list", server.updateList)

	// private.POST("/create/tasks", server.createTasks)
	private.POST("/create/task", server.createTask)
	private.POST("/get/tasks", server.getTasks)
	private.POST("/update/task", server.updateTask)

	private.GET("/get/profile", server.getProfile)
	private.POST("/update/profile", server.updateProfile)

	private.POST("/add/notification/token", server.addNotificationToken)
	// private.GET("/get/notification/tokens", server.getNotificationTokens)

	return r
}
