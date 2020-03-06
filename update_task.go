package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func (server *server) updateTask(c *gin.Context) {
	var args models.UpdateTaskArgs

	err := json.NewDecoder(c.Request.Body).Decode(&args)
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

	err = user.UpdateTask(server.db, args)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, "No such task")
			return
		}
		log.Printf("Error when updating task\n%v", err)
		c.JSON(http.StatusInternalServerError, "Error when updating task")
		return
	}

	c.JSON(http.StatusOK, nil)
	return
}
