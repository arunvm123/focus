package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
)

func (server *server) createTasks(c *gin.Context) {
	var createTasksArgs []models.CreateTaskArgs

	err := json.NewDecoder(c.Request.Body).Decode(&createTasksArgs)
	if err != nil {
		log.Printf("Error when decoding request body\n%v", err)
		c.JSON(http.StatusInternalServerError, "Request body not properly formatted")
		return
	}

	user, ok := c.Keys["user"].(*models.User)
	if !ok {
		log.Println("Unable to fetch user")
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	err = user.CreateTasks(server.db, &createTasksArgs)
	if err != nil {
		log.Printf("Error when creating tasks\n%v", err)
		c.JSON(http.StatusInternalServerError, "Error when creating tasks")
		return
	}

	c.JSON(http.StatusOK, nil)
	return
}
