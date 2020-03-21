package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/arunvm/travail-backend/config"
	"github.com/arunvm/travail-backend/emails"
	"github.com/dgrijalva/jwt-go"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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
		c.JSON(http.StatusConflict, "Email already exists")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error when hashing password\n%v", err)
		c.JSON(http.StatusInternalServerError, "Error hashing password")
		return
	}

	user.Password = string(passwordHash)
	user.Verified = false

	tx := server.db.Begin()
	err = user.Create(tx)
	if err != nil {
		tx.Rollback()
		log.Printf("Error when inserting data to database\n%v", err)
		c.JSON(http.StatusInternalServerError, "Error inserting data into DB")
		return
	}

	token, err := models.CreateEmailValidationToken(tx, &user)
	if err != nil {
		tx.Rollback()
		log.Printf("Error when creating email validate token\n%v", err)
		c.JSON(http.StatusInternalServerError, "Error creating email validate token")
		return
	}

	err = emails.SendValidationEmail(server.email, &user, token)
	if err != nil {
		tx.Rollback()
		log.Printf("Error sending validation email\n%v", err)
		c.JSON(http.StatusInternalServerError, "Error sending email")
		return
	}

	tx.Commit()

	c.Status(http.StatusOK)
	return
}

func (server *server) login(c *gin.Context) {
	var loginData loginRequest

	err := json.NewDecoder(c.Request.Body).Decode(&loginData)
	if err != nil {
		log.Printf("Error when decoding request body\n%v", err)
		c.JSON(http.StatusInternalServerError, "Request body not properly formatted")
		return
	}

	user, err := models.GetUserFromEmail(server.db, loginData.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, "User does not exist, Please sign up")
			return
		}
		log.Printf("Error when looking up user with email\n%v", err)
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}

	if !user.Verified {
		c.JSON(http.StatusUnauthorized, "Email not verified")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		log.Printf("Passwords do not match\n%v", err)
		c.JSON(http.StatusUnauthorized, "Wrong password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	config, err := config.GetConfig()
	if err != nil {
		log.Printf("Failed to read config file\n%v", err)
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}

	signedToken, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		log.Printf("Error when signing token\n%v", err)
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}

	c.SetCookie("Authorization", signedToken, 0, "", "travail.in", false, false)

	c.JSON(http.StatusOK, struct {
		Token string `json:"token"`
		Name  string `json:"name"`
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{
		Token: signedToken,
		Email: user.Email,
		ID:    user.ID,
		Name:  user.Name,
	})
	return
}
