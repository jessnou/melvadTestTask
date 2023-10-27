package main

import (
	"log"
	"melvadTestTask/internal/config"
	"melvadTestTask/internal/db"
	"melvadTestTask/internal/server"
)

func main() {
	dbConn, redisClient, err := db.GetDBConn()
	if err != nil {
		log.Fatalf("Failed to connect %v", err)
	}
	app := config.GetAppConfig(dbConn, redisClient)
	server.StartServer(app)
}
