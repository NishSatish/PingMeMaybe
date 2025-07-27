package server

import (
	_ "encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StartServer(db *pgxpool.Pool) {
	r := gin.Default()

	serverWithRoutes := SetRoutes(r, db)

	err := serverWithRoutes.Run(":8080")
	if err != nil {
		panic(err)
		return
	}
}
