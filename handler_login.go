package main

import (
	"net/http"
	"time"

	"github.com/arunvm/travail-backend/config"
	"github.com/arunvm/travail-backend/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (server *server) login(c *gin.Context) {
	var loginData models.LoginArgs

	err := c.ShouldBindJSON(&loginData)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "login",
			"info": "decoding request body",
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	user, err := models.GetUserFromEmail(server.db, loginData.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, "User does not exist, Please sign up")
			return
		}
		log.WithFields(log.Fields{
			"func":    "login",
			"subFunc": "models.GetUserFromEmail",
			"email":   loginData.Email,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}

	if !user.Verified {
		c.JSON(http.StatusUnauthorized, "Email not verified")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "login",
			"subFunc": "bcrypt.CompareHashAndPassword",
			"email":   user.Email,
		}).Error(err)
		c.JSON(http.StatusUnauthorized, "Wrong password")
		return
	}

	signedToken, err := getJWTToken(user.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getJWTToken",
			"subFunc": "token.SignedString",
			"email":   user.Email,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error while retrieving token")
		return
	}

	personalTeam, err := user.GetPersonalTeam(server.db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "login",
			"subFunc": "user.GetPersonalTeam",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error while retrieving personal team")
		return
	}

	c.SetCookie("Authorization", signedToken, 0, "", "travail.in", false, false)

	c.JSON(http.StatusOK, struct {
		Token        string  `json:"token"`
		Name         string  `json:"name"`
		ID           int     `json:"id"`
		Email        string  `json:"email"`
		ProfilePic   *string `json:"profilePic"`
		GoogleOauth  bool    `json:"googleOauth"`
		PersonalTeam string  `json:"personalTeam"`
	}{
		Token:        signedToken,
		Email:        user.Email,
		ID:           user.ID,
		Name:         user.Name,
		ProfilePic:   user.ProfilePic,
		GoogleOauth:  user.GoogleOauth,
		PersonalTeam: personalTeam,
	})
	return
}

func getJWTToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	config, err := config.GetConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getJWTToken",
			"subFunc": "config.GetConfig",
			"userID":  userID,
		}).Error(err)
		return "", err
	}

	signedToken, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getJWTToken",
			"subFunc": "token.SignedString",
			"userID":  userID,
		}).Error(err)
		return "", err
	}

	return signedToken, nil
}
