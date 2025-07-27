package server

import (
	"PingMeMaybe/gateway/pkg/service"
	"PingMeMaybe/libs/config"
	"PingMeMaybe/libs/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetRoutes(r *gin.Engine, dbConn *pgxpool.Pool) *gin.Engine {
	asynqClient := config.GetAsynqClient()
	dbService := db.NewDBService(dbConn)

	services := service.InitAppServices(asynqClient, dbService)

	r.POST("/notification", services.NotificationsService().QueueNotification)

	return r
}
