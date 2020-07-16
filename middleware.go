package main

import (
	"net/http"

	"github.com/arunvm/travail-backend/config"
	"github.com/arunvm/travail-backend/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (server *server) tokenAuthorisationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, "Provide token")
			c.Abort()
			return
		}

		user, err := server.getUserFromToken(token)
		if err != nil {
			log.WithFields(log.Fields{
				"func":    "tokenAuthorisationMiddleware",
				"subFunc": "server.getUserFromToken",
			}).Error(err)
			c.JSON(http.StatusUnauthorized, "Invalid user")
			c.Abort()
			return
		}

		c.Keys = make(map[string]interface{})
		c.Keys["user"] = user
		c.Next()
	}
}

func (server *server) checkIfOrganisationAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := getUserFromContext(c)
		if err != nil {
			log.WithFields(log.Fields{
				"func":    "checkIfOrganisationAdmin",
				"subFunc": "getUserFromContext",
			}).Error(err)
			c.JSON(http.StatusUnauthorized, "Invalid user")
			c.Abort()
			return
		}

		orgID := c.Query("organisationID")
		if orgID == "" {
			c.JSON(http.StatusUnauthorized, "Provide organisation id")
			c.Abort()
			return
		}
		c.Keys["organisationID"] = orgID

		if server.db.CheckIfOrganisationAdmin(orgID, user) == false {
			c.JSON(http.StatusUnauthorized, "Invalid user")
			c.Abort()
			return
		}

		c.Next()
	}
}

func (server *server) checkIfOrganisationMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := getUserFromContext(c)
		if err != nil {
			log.WithFields(log.Fields{
				"func":    "checkIfOrganisationMember",
				"subFunc": "getUserFromContext",
			}).Error(err)
			c.JSON(http.StatusUnauthorized, "Invalid user")
			return
		}

		orgID := c.Query("organisationID")
		if orgID == "" {
			c.JSON(http.StatusUnauthorized, "Provide organisation id")
			c.Abort()
			return
		}
		c.Keys["organisationID"] = orgID

		if server.db.CheckIfOrganisationMember(orgID, user) == false {
			c.JSON(http.StatusUnauthorized, "Invalid user")
			c.Abort()
			return
		}

		c.Next()
	}
}

func (server *server) checkIfTeamAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := getUserFromContext(c)
		if err != nil {
			log.WithFields(log.Fields{
				"func":    "checkIfTeamAdmin",
				"subFunc": "getUserFromContext",
			}).Error(err)
			c.JSON(http.StatusUnauthorized, "Invalid user")
			return
		}

		teamID := c.Query("teamID")
		if teamID == "" {
			c.JSON(http.StatusUnauthorized, "Provide team id")
			c.Abort()
			return
		}
		c.Keys["teamID"] = teamID

		if server.db.CheckIfTeamAdmin(teamID, user) == false {
			c.JSON(http.StatusUnauthorized, "Invalid user")
			c.Abort()
			return
		}

		c.Next()
	}
}

func (server *server) checkIfTeamMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := getUserFromContext(c)
		if err != nil {
			log.WithFields(log.Fields{
				"func":    "checkIfTeamAdmin",
				"subFunc": "getUserFromContext",
			}).Error(err)
			c.JSON(http.StatusUnauthorized, "Invalid user")
			return
		}

		teamID := c.Query("teamID")
		if teamID == "" {
			c.JSON(http.StatusUnauthorized, "Provide team id")
			c.Abort()
			return
		}
		c.Keys["teamID"] = teamID

		if server.db.CheckIfTeamMember(teamID, user) == false {
			c.JSON(http.StatusUnauthorized, "Invalid user")
			c.Abort()
			return
		}

		c.Next()
	}
}

func (server *server) checkIfBoardPartOfTeam() gin.HandlerFunc {
	return func(c *gin.Context) {
		if server.db.CheckIfBoardPartOfTeam(c.Query("boardID"), c.Query("teamID")) == false {
			c.JSON(http.StatusUnauthorized, "Board not part of team")
			c.Abort()
			return
		}

		c.Keys["boardID"] = c.Query("boardID")
		c.Next()
	}
}

func (server *server) checkIfColumnPartOfBoard() gin.HandlerFunc {
	return func(c *gin.Context) {
		if server.db.CheckIfColumnPartOfBoard(c.Query("boardColumnID"), c.Query("boardID")) == false {
			c.JSON(http.StatusUnauthorized, "Board column not of specified board")
			c.Abort()
			return
		}

		c.Keys["boardColumnID"] = c.Query("boardColumnID")
		c.Next()
	}
}

func (server *server) checkIfAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := getUserFromContext(c)
		if err != nil {
			log.WithFields(log.Fields{
				"func":    "CheckIfAdminMiddleware",
				"subFunc": "getUserFromContext",
			}).Error(err)
			c.JSON(http.StatusUnauthorized, "Invalid user")
			c.Abort()
			return
		}

		confguration, err := config.GetConfig()
		if err != nil {
			log.WithFields(log.Fields{
				"func":    "CheckIfAdminMiddleware",
				"subFunc": "config.GetConfig",
			}).Error(err)
			c.JSON(http.StatusInternalServerError, "Error reading config file")
			c.Abort()
			return
		}

		var flag bool
		for i := 0; i < len(confguration.AdminIDs); i++ {
			if user.ID == confguration.AdminIDs[i] {
				flag = true
				break
			}
		}

		if flag == false {
			c.JSON(http.StatusUnauthorized, "User not admin")
			c.Abort()
			return
		}

		c.Next()
	}
}

func (server *server) getUserFromToken(token string) (*models.User, error) {
	config, err := config.GetConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getUserFromToken",
			"subFunc": "config.GetConfig",
		}).Error(err)
		return nil, err
	}

	parsedString, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getUserFromToken",
			"subFunc": "jwt.Parse",
		}).Error(err)
		return nil, err
	}

	userID := parsedString.Claims.(jwt.MapClaims)["id"].(float64)

	user, err := server.db.GetUserFromID(int(userID))
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getUserFromToken",
			"subFunc": "models.GetUserFromID",
			"userID":  int(userID),
		}).Error(err)
		return nil, err
	}

	return user, nil
}
