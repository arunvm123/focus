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

	user, err := models.GetUserFromID(server.db, int(userID))
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
