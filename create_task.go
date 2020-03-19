package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
)

func (server *server) createTask(c *gin.Context) {
	var task *models.Task

	err := json.NewDecoder(c.Request.Body).Decode(&task)
	if err != nil {
		log.Printf("Error when decoding request body\n%v", err)
		c.JSON(http.StatusInternalServerError, "Request body not properly formatted")
		return
	}

	user, err := getUserFromContext(c)
	if err != nil {
		log.Println("Unable to fetch user")
		c.JSON(http.StatusInternalServerError, "Error fetching user")
		return
	}

	err = user.CreateTask(server.db, task)
	if err != nil {
		log.Printf("Error creating task")
		c.JSON(http.StatusInternalServerError, "Error when creating task")
		return
	}

	c.JSON(http.StatusOK, task)
	return
}
