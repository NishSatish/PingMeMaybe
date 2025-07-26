package main

import (
	"PingMeMaybe/gateway/server"
	"PingMeMaybe/libs/db"
	"context"
	"log"
)

func main() {
	dbConn, err := db.InitDBConn()
	defer dbConn.Close(context.Background())
	if err != nil {
		log.Fatal("Failed to initialize database connection: " + err.Error())
	}

	server.StartServer()
}
