package main

import (
	"net/http"
	"time"

	"github.com/arunvm/travail-backend/config"
	"github.com/arunvm/travail-backend/emails"
	"github.com/dgrijalva/jwt-go"

	"github.com/arunvm/travail-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (server *server) signup(c *gin.Context) {
	var args models.SignUpArgs
	err := c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "signup",
			"info": "decoding request body",
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	if models.CheckIfUserExists(server.db, args.Email) == true {
		c.JSON(http.StatusConflict, "Email already exists")
		return
	}

	tx := server.db.Begin()
	user, err := models.UserSignup(tx, &args)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "signup",
			"subFunc": "models.UserSignup",
			"email":   args.Email,
		})
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	token, err := models.CreateEmailValidationToken(tx, user)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "signup",
			"subFunc": "models.CreateEmailValidationToken",
			"email":   user.Email,
		})
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = emails.SendValidationEmail(server.email, user, token)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "signup",
			"subFunc": "emails.SendValidationEmail",
			"email":   user.Email,
		})
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	tx.Commit()
	c.Status(http.StatusOK)
	return
}

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
			c.JSON(http.StatusOK, "User does not exist, Please sign up")
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	config, err := config.GetConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "login",
			"subFunc": "config.GetConfig",
			"email":   user.Email,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}

	signedToken, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "login",
			"subFunc": "token.SignedString",
			"email":   user.Email,
		}).Error(err)
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
