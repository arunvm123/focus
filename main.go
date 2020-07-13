package main

import (
	"context"
	"flag"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/rs/cors"
	"google.golang.org/api/option"

	"github.com/arunvm/travail-backend/config"
	"github.com/arunvm/travail-backend/email"
	"github.com/arunvm/travail-backend/email/sendgrid"

	"github.com/arunvm/travail-backend/models"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type server struct {
	db         *gorm.DB
	routes     *gin.Engine
	email      email.Email
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

	// Reading file path from flag
	filePath := flag.String("config-path", "config.yaml", "filepath to configuration file")
	flag.Parse()

	// Reading config variables
	config, err := config.Initialise(*filePath)
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
	server.email = sendgrid.New(config.SendgridKey)

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

	http.ListenAndServeTLS(":"+config.Port, "./certs/fullchain.pem", "./certs/privkey.pem", routes)
}
