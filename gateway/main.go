package main

import (
	"PingMeMaybe/gateway/server"
	"PingMeMaybe/libs/db"
	"log"
)

func main() {
	dbConn, err := db.InitDBPoolConn()
	defer dbConn.Close()
	if err != nil {
		log.Fatal("Failed to initialize database connection: " + err.Error())
	}

	server.StartServer(dbConn)
}
