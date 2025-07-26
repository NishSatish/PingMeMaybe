package db

import (
	configLib "PingMeMaybe/libs/config"
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

func InitDBConn() (*pgx.Conn, error) {
	configLib.LoadEnv(".")
	conn, err := pgx.Connect(context.Background(), configLib.GetConfig().GetString("DATABASE_SESSION_POOLING_MODE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("DATABASE CONNECTED, version:", version)
	return conn, err
}
