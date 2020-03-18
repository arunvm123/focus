package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
)

func (server *server) updateList(c *gin.Context) {
	user, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var args models.UpdateListArgs
	err = json.NewDecoder(c.Request.Body).Decode(&args)
	if err != nil {
		log.Printf("Error when decoding request body\n%v", err)
		c.JSON(http.StatusInternalServerError, "Request body not properly formatted")
		return
	}

	err = user.UpdateList(server.db, &args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
	return
}
