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
	public.POST("/forgot/password", server.forgotPassword)
	public.POST("/forgot/password/reset", server.resetPassword)

	private := r.Group("/")
	private.Use(server.tokenAuthorisationMiddleware())

	private.POST("/create/list", server.createList)
	private.POST("/get/lists", server.getLists)
	private.POST("/update/list", server.updateList)

	// private.POST("/create/tasks", server.createTasks)
	private.POST("/create/task", server.createTask)
	private.POST("/get/tasks", server.getTasks)
	private.POST("/update/task", server.updateTask)
	private.DELETE("/delete/tasks", server.deleteTasks)

	private.POST("/create/organisation", server.createOrganisation)
	private.GET("/get/organisations", server.getOrganisations)
	private.POST("/accept/organisation/invite", server.acceptOrganisationInvite)

	organisationAdmin := r.Group("/")
	organisationAdmin.Use(server.tokenAuthorisationMiddleware(), server.checkIfOrganisationAdmin())
	organisationAdmin.POST("/update/organisation", server.updateOrganisation)
	organisationAdmin.POST("/organisation/invite", server.inviteToOrganisation)

	organisationMember := r.Group("/")
	organisationMember.Use(server.tokenAuthorisationMiddleware(), server.checkIfOrganisationMember())
	organisationMember.GET("/get/organisation/members", server.getOrganisationMembers)

	organisationMember.POST("/create/team", server.createTeam)

	teamAdmin := r.Group("/")
	teamAdmin.Use(server.tokenAuthorisationMiddleware(), server.checkIfTeamAdmin())
	teamAdmin.POST("/update/team", server.updateTeam)
	teamAdmin.POST("/add/team/member", server.addTeamMember)

	teamMember := r.Group("/")
	teamMember.Use(server.tokenAuthorisationMiddleware(), server.checkIfTeamMember())
	teamMember.GET("/get/team/members", server.getTeamMembers)

	teamMember.POST("/create/board", server.createBoard)
	teamMember.POST("/get/boards", server.getBoards)
	teamMember.POST("/update/board", server.updateBoard)
	teamMember.DELETE("/delete/board", server.deleteBoard) // Implementation Pending

	board := r.Group("/")
	board.Use(server.tokenAuthorisationMiddleware(), server.checkIfTeamMember(), server.checkIfBoardPartOfTeam())
	board.POST("/create/board/column", server.createBoardColumn)
	board.GET("/get/board/columns", server.getBoardColumns)
	board.POST("/update/board/column", server.updateBoardColumn)
	board.DELETE("/delete/board/column", server.deleteBoardColumn) // Implementation Pending

	boardColumn := r.Group("/")
	boardColumn.Use(server.tokenAuthorisationMiddleware(), server.checkIfTeamMember(), server.checkIfBoardPartOfTeam(), server.checkIfColumnPartOfBoard())
	boardColumn.POST("/create/board/column/card", server.createColumnCard)
	boardColumn.GET("/get/board/column/cards", server.getColumnCards)
	boardColumn.POST("/update/board/column/card", server.updateColumnCard)
	boardColumn.DELETE("/delete/board/column/card", server.deleteColumnCard) // Implementation Pending

	private.GET("/get/profile", server.getProfile)
	private.POST("/update/profile", server.updateProfile)
	private.POST("/update/password", server.updatePassword)

	private.POST("/add/notification/token", server.addNotificationToken)
	// private.GET("/get/notification/tokens", server.getNotificationTokens)

	admin := r.Group("/")
	admin.Use(server.tokenAuthorisationMiddleware(), server.checkIfAdminMiddleware())

	private.POST("/create/bug", server.createBug)
	admin.GET("/admin/check", server.adminCheck)
	admin.GET("/get/bugs", server.getBugs)
	admin.POST("/update/bug", server.updateBug)

	return r
}
