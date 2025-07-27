package main

import (
	"PingMeMaybe/libs/db"
	"PingMeMaybe/processor/pkg/cron"
	"PingMeMaybe/processor/server"
)

func main() {
	dbConn, err := db.InitDBPoolConn()
	defer dbConn.Close()
	if err != nil {
		panic("Failed to initialize database connection: " + err.Error())
	}

	dbService := db.NewDBService(dbConn)
	crons := cron.GetCrons(dbService)

	// CRONS
	go crons.StartMarkFailuresCron() // v1: every 10 seconds

	// Asynq listener
	server.StartAsynqServer(dbConn)
}
