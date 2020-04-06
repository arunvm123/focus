package main

import (
	"context"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/rs/cors"
	"google.golang.org/api/option"

	"github.com/arunvm/travail-backend/config"

	"github.com/arunvm/travail-backend/models"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sendgrid/sendgrid-go"
	log "github.com/sirupsen/logrus"
)

type server struct {
	db         *gorm.DB
	routes     *gin.Engine
	email      *sendgrid.Client
	pushClient *messaging.Client
}

func newServer() *server {
	s := server{}
	return &s
}

func main() {
	server := newServer()

	// Logging options
	log.SetFormatter(&log.JSONFormatter{})

	// Reading config variables
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

	// FCM push notification
	firebaseApp, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(config.FCMServiceAccountKeyPath))
	if err != nil {
		log.Fatalf("error when initialising firebase app\n%v", err)
	}

	server.pushClient, err = firebaseApp.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error when initialising FCM push notification client\n%v", err)
	}

	err = server.startCronJobs()
	if err != nil {
		log.Fatalf("error starting cron jobs\n%v", err)
	}

	// Setting up routes
	server.routes = initialiseRoutes(server)
	routes := cors.AllowAll().Handler(server.routes)

	http.ListenAndServeTLS(":5000", "./certs/fullchain.pem", "./certs/privkey.pem", routes)
}
