package main

import (
	"log"
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
)

func (server *server) getTasks(c *gin.Context) {
	user, ok := c.Keys["user"].(*models.User)
	if !ok {
		log.Println("Unable to fetch user")
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	tasks, err := user.GetTasks(server.db)
	if err != nil {
		log.Printf("Error when fetching tasks\n%v", err)
		c.JSON(http.StatusInternalServerError, "Error fetching tasks")
		return
	}

	c.JSON(http.StatusOK, tasks)
	return
}
