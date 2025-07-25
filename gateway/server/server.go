package server

import (
	_ "encoding/json"
	"github.com/gin-gonic/gin"
)

func StartServer() {
	startHttpServer()
}

func startHttpServer() {
	r := gin.Default()

	serverWithRoutes := SetRoutes(r)

	serverWithRoutes.Run(":8080")
}
