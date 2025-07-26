package server

import (
	"PingMeMaybe/gateway/pkg/service"
	"PingMeMaybe/libs/config"
	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine) *gin.Engine {
	asynqClient := config.GetAsynqClient()

	services := service.InitAppServices(asynqClient)

	r.POST("/notification", services.NotificationsService().QueueNotification)

	return r
}
