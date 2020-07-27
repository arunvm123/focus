package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/arunvm/focus/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

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

	tx := server.tx.Begin()

	var user *models.User
	if tx.CheckIfUserExists(userInfo.Email) == true {
		user, err = tx.GetUserFromEmail(userInfo.Email)
		if err != nil {
			tx.Rollback()
			log.WithFields(log.Fields{
				"func":    "loginWithGoogle",
				"subFunc": "models.GetUserFromEmail",
				"email":   userInfo.Email,
			}).Error(err)
			c.JSON(http.StatusInternalServerError, "Error when fetching user details")
			return
		}

		if user.GoogleOauth == false {
			tx.Rollback()
			c.JSON(http.StatusUnauthorized, "Login with email and password")
			return
		}
	} else {
		var token string

		user, token, err = tx.UserSignup(&models.SignUpArgs{
			Email: userInfo.Email,
			Name:  userInfo.Name,
		}, true)
		if err != nil {
			tx.Rollback()
			log.WithFields(log.Fields{
				"func":    "loginWithGoogle",
				"subFunc": "models.SignUpWithGoogle",
				"email":   userInfo.Email,
			}).Error(err)
			c.JSON(http.StatusInternalServerError, "Error when signing up with google")
			return
		}

		err = server.email.SendValidationEmail(user.Name, user.Email, token)
		if err != nil {
			tx.Rollback()
			log.WithFields(log.Fields{
				"func":    "UserSignup",
				"subFunc": "emailClient.SendValidationEmail",
				"userID":  user.ID,
			})
			c.JSON(http.StatusInternalServerError, "Error while sending email")
			return
		}
	}

	signedToken, err := getJWTToken(user.ID)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "loginWithGoogle",
			"subFunc": "getJWTToken",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error while retrieving token")
		return
	}

	personalTeamID, err := tx.GetPersonalTeamID(user)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "loginWithGoogle",
			"subFunc": "user.GetPersonalTeamID",
			"userID":  user.ID,
		}).Error(err)
		c.JSON(http.StatusInternalServerError, "Error while retrieving personal team id")
		return

	}

	tx.Commit()
	// c.SetCookie("Authorization", signedToken, 0, "", "focus.in", false, false)

	c.JSON(http.StatusOK, struct {
		Token          string  `json:"token"`
		Name           string  `json:"name"`
		ID             int     `json:"id"`
		Email          string  `json:"email"`
		ProfilePic     *string `json:"profilePic"`
		GoogleOauth    bool    `json:"googleOauth"`
		PersonalTeamID string  `json:"personalTeamID"`
	}{
		Token:          signedToken,
		Email:          user.Email,
		ID:             user.ID,
		Name:           user.Name,
		ProfilePic:     user.ProfilePic,
		GoogleOauth:    user.GoogleOauth,
		PersonalTeamID: personalTeamID,
	})
	return
}
