package service

import "github.com/jackc/pgx/v5"

type ProcessorServices struct {
	INotificationProcessorService
}

type IProcessorServices interface {
	INotificationProcessorService
}

func NewProcessorServices(db *pgx.Conn) IProcessorServices {
	return &ProcessorServices{
		INotificationProcessorService: NewNotificationProcessorService(db),
	}
}
