package main

import (
	"PingMeMaybe/libs/db"
	"PingMeMaybe/processor/server"
)

func main() {
	dbConn, err := db.InitDBConn()
	if err != nil {
		panic("Failed to initialize database connection: " + err.Error())
	}

	server.StartServer(dbConn)
}
