package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func (server *server) signup(c *gin.Context) {
	var user models.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		log.Printf("Error when decoding request body\n%v", err)
		c.JSON(http.StatusInternalServerError, "Request body not properly formatted")
		return
	}

	var count int
	err = server.db.Table("users").Where("email = ?", user.Email).Count(&count).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("Error when checking if user exists\n%v", err)
			c.JSON(http.StatusInternalServerError, "Internal error")
			return
		}
	}

	if count > 0 {
		c.JSON(http.StatusOK, "Email already exists")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error when hashing password\n%v", err)
		c.JSON(http.StatusInternalServerError, "Error hashing password")
		return
	}

	user.Password = string(passwordHash)

	err = user.Create(server.db)
	if err != nil {
		log.Printf("Error when inserting data to database\n%v", err)
		c.JSON(http.StatusInternalServerError, "Error inserting data into DB")
		return
	}

	c.Status(http.StatusOK)
	return
}
