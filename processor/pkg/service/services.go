package service

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProcessorServices struct {
	INotificationProcessorService
}

type IProcessorServices interface {
	INotificationProcessorService
}

func NewProcessorServices(db *pgxpool.Pool) IProcessorServices {
	return &ProcessorServices{
		INotificationProcessorService: NewNotificationProcessorService(db),
	}
}
