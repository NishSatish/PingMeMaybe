package server

import (
	"PingMeMaybe/gateway/pkg/service"
	"PingMeMaybe/libs/config"
	"PingMeMaybe/libs/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func SetRoutes(r *gin.Engine, dbConn *pgx.Conn) *gin.Engine {
	asynqClient := config.GetAsynqClient()
	dbService := db.NewDBService(dbConn)

	services := service.InitAppServices(asynqClient, dbService)

	r.POST("/notification", services.NotificationsService().QueueNotification)

	return r
}
