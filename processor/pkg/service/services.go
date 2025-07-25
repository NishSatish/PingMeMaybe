package service

type ProcessorServices struct {
	INotificationProcessorService
}

type IProcessorServices interface {
	INotificationProcessorService
}

func NewProcessorServices() IProcessorServices {
	return &ProcessorServices{
		INotificationProcessorService: NewNotificationProcessorService(),
	}
}
