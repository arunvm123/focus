package main

import (
	"encoding/json"
	"io/ioutil"
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

	c.SetCookie("Authorization", signedToken, 0, "", "travail.in", false, false)

	c.JSON(http.StatusOK, struct {
		Token      string  `json:"token"`
		Name       string  `json:"name"`
		ID         int     `json:"id"`
		Email      string  `json:"email"`
		ProfilePic *string `json:"profilePic"`
	}{
		Token:      signedToken,
		Email:      user.Email,
		ID:         user.ID,
		Name:       user.Name,
		ProfilePic: user.ProfilePic,
	})
	return
}

func (server *server) loginWithGoogle(c *gin.Context) {
	var args struct {
		AccessToken string `json:"access_token" binding:"required"`
	}

	err := c.ShouldBindJSON(&args)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "loginWithGoogle",
			"info": "decoding request body",
		}).Error(err)
		c.JSON(http.StatusBadRequest, "Request body not properly formatted")
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + args.AccessToken)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "loginWithGoogle",
			"info": "retrieving user info from google",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Internal error")
		return
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "loginWithGoogle",
			"subFunc": "ioutil.ReadAll",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "error reading response")
		return
	}

	var userInfo models.LoginWithGoogleArgs
	err = json.Unmarshal(contents, &userInfo)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "loginWithGoogle",
			"info": "unmarshalling info to struct",
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error unmarshalling data")
		return
	}

	// When fetching user info from google api, if the access token is invalid, it doesnt raise an error
	// This is to ensure that valid info was retrieved.
	if userInfo.Email == "" {
		c.JSON(http.StatusUnauthorized, "Invalid access token")
		return
	}

	var user *models.User
	if models.CheckIfUserExists(server.db, userInfo.Email) == true {
		user, err = models.GetUserFromEmail(server.db, userInfo.Email)
		if err != nil {
			log.WithFields(log.Fields{
				"func":    "loginWithGoogle",
				"subFunc": "models.GetUserFromEmail",
				"email":   userInfo.Email,
			}).Error(err)
			c.JSON(http.StatusInternalServerError, "Error when fetching user details")
			return
		}

		if user.GoogleOauth == false {
			c.JSON(http.StatusUnauthorized, "Login with email and password")
			return
		}
	} else {
		user, err = models.SignUpWithGoogle(server.db, &userInfo)
		if err != nil {
			log.WithFields(log.Fields{
				"func":    "loginWithGoogle",
				"subFunc": "models.SignUpWithGoogle",
				"email":   userInfo.Email,
			}).Error(err)
			c.JSON(http.StatusInternalServerError, "Error when signing up with google")
			return
		}
	}

	signedToken, err := getJWTToken(user.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "loginWithGoogle",
			"subFunc": "getJWTToken",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error while retrieving token")
		return
	}

	c.SetCookie("Authorization", signedToken, 0, "", "travail.in", false, false)

	c.JSON(http.StatusOK, struct {
		Token      string  `json:"token"`
		Name       string  `json:"name"`
		ID         int     `json:"id"`
		Email      string  `json:"email"`
		ProfilePic *string `json:"profilePic"`
	}{
		Token:      signedToken,
		Email:      user.Email,
		ID:         user.ID,
		Name:       user.Name,
		ProfilePic: user.ProfilePic,
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
