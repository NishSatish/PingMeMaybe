package db

import (
	"PingMeMaybe/libs/db/models"
	"github.com/jackc/pgx/v5"
)

// Service that just interfaces all DB model services from one place

type DBService struct {
	Notifications models.INotificationRepository
}

type DBServiceInterface interface {
	NotificationsRepository() models.INotificationRepository
}

func (this DBService) NotificationsRepository() models.INotificationRepository {
	return this.Notifications
}

func NewDBService(db *pgx.Conn) *DBService {
	return &DBService{
		Notifications: models.NewNotificationRepo(db),
	}
}
