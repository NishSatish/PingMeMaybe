package cron

import (
	"PingMeMaybe/libs/db"
)

type Crons struct {
	MarkFailuresCronInterface
}

type CronsInterface interface {
	MarkFailuresCronInterface
}

func GetCrons(dbService *db.DBService) CronsInterface {
	return &Crons{
		NewMarkFailuresCron(dbService),
	}
}
