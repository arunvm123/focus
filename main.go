package main

import (
	"flag"
	"net/http"

	"github.com/rs/cors"

	"github.com/arunvm/travail-backend/config"
	"github.com/arunvm/travail-backend/email"
	"github.com/arunvm/travail-backend/email/sendgrid"
	"github.com/arunvm/travail-backend/models"
	"github.com/arunvm/travail-backend/models/mysql"
	push "github.com/arunvm/travail-backend/push_notification"
	"github.com/arunvm/travail-backend/push_notification/fcm"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type server struct {
	db     models.DB
	routes http.Handler
	email  email.Email
	push   push.Notification
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
	db, err := mysql.New(config.Database.User + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":" + config.Database.Port + ")/" + config.Database.DatabaseName + "?parseTime=true")
	if err != nil {
		panic(err)
	}

	server.db = db
	mysql.MigrateDB(db.Client)

	// email client
	server.email = sendgrid.New(config.SendgridKey)

	server.push, err = fcm.New(config.FCMServiceAccountKeyPath)
	if err != nil {
		log.Fatalf("error retrieving client for push notification\n%v", err)
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
