package main

import (
	"net/http"

	"github.com/rs/cors"

	"github.com/arunvm/travail-backend/config"

	"github.com/arunvm/travail-backend/models"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sendgrid/sendgrid-go"
	log "github.com/sirupsen/logrus"
)

type server struct {
	db     *gorm.DB
	routes *gin.Engine
	email  *sendgrid.Client
}

func newServer() *server {
	s := server{}
	return &s
}

func main() {
	server := newServer()

	log.SetFormatter(&log.JSONFormatter{})

	config, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Failed to read config\n%v", err)
	}

	// <username>:<pw>@tcp(<HOST>:<port>)/<dbname>")
	server.db, err = gorm.Open("mysql", config.Database.User+":"+config.Database.Password+"@tcp("+config.Database.Host+":"+config.Database.Port+")/"+config.Database.DatabaseName+"?parseTime=true")
	if err != nil {
		panic(err)
	}

	server.db.LogMode(true)
	models.MigrateDB(server.db)

	// email client
	server.email = sendgrid.NewSendClient(config.SendgridKey)

	server.routes = initialiseRoutes(server)

	routes := cors.AllowAll().Handler(server.routes)

	http.ListenAndServeTLS(":5000", "./certs/fullchain.pem", "./certs/privkey.pem", routes)
}
