package db

import (
	"PingMeMaybe/libs/db/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service that just interfaces all DB model services from one place

type DBService struct {
	Notifications models.INotificationRepository
	UserCohorts   models.IUserCohortRepository
}

type DBServiceInterface interface {
	NotificationsRepository() models.INotificationRepository
	UserCohortsRepository() models.IUserCohortRepository
}

func (this DBService) NotificationsRepository() models.INotificationRepository {
	return this.Notifications
}

func (this DBService) UserCohortsRepository() models.IUserCohortRepository {
	return this.UserCohorts
}

func NewDBService(db *pgxpool.Pool) *DBService {
	return &DBService{
		Notifications: models.NewNotificationRepo(db),
		UserCohorts:   models.NewUserCohortRepo(db),
	}
}
