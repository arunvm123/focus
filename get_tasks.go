package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
)

func (server *server) getTasks(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		log.Println("Unable to fetch user")
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	var args models.GetTasksArgs
	err = json.NewDecoder(c.Request.Body).Decode(&args)
	if err != nil {
		log.Printf("Error when decoding request body\n%v", err)
		c.JSON(http.StatusInternalServerError, "Request body not properly formatted")
		return
	}

	tasks, err := user.GetTasks(server.db, &args)
	if err != nil {
		log.Printf("Error when fetching tasks\n%v", err)
		c.JSON(http.StatusInternalServerError, "Error fetching tasks")
		return
	}

	c.JSON(http.StatusOK, tasks)
	return
}
