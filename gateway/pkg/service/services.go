package service

import (
	"PingMeMaybe/gateway/pkg/service/notifications"
	"PingMeMaybe/libs/db"
	"github.com/hibiken/asynq"
)

// Basically accumulate all the services here

// A little complication to add branching, basically so that 10 services dont have all their functions nested under one services object

type AppServices struct {
	Notifications notifications.NotificationsServiceInterface
}

type AppServicesInterface interface {
	NotificationsService() notifications.NotificationsServiceInterface
}

func (a *AppServices) NotificationsService() notifications.NotificationsServiceInterface {
	return a.Notifications
}

func InitAppServices(asynq *asynq.Client, dbService *db.DBService) AppServicesInterface {
	return &AppServices{
		Notifications: notifications.NewNotificationsService(asynq, dbService.NotificationsRepository()),
	}
}
