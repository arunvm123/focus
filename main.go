package main

import (
	"log"
	"net/http"

	"github.com/arunvm/travail-backend/config"

	"github.com/arunvm/travail-backend/models"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type server struct {
	db     *gorm.DB
	routes *gin.Engine
}

func newServer() *server {
	s := server{}
	return &s
}

func main() {
	server := newServer()

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

	server.routes = initialiseRoutes(server)

	http.ListenAndServe(":5000", server.routes)
}
