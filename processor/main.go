package main

import (
	"PingMeMaybe/libs/db"
	"PingMeMaybe/processor/server"
)

func main() {
	dbConn, err := db.InitDBPoolConn()
	if err != nil {
		panic("Failed to initialize database connection: " + err.Error())
	}

	server.StartServer(dbConn)
}
