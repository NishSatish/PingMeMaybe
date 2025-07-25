package service

import (
	"PingMeMaybe/libs/dto"
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"log"
)

type notificationProcessorService struct {
}

type INotificationProcessorService interface {
	HandleNotificationQueueItems(ctx context.Context, task *asynq.Task) error
}

func NewNotificationProcessorService() INotificationProcessorService {
	return &notificationProcessorService{}
}

func (n notificationProcessorService) HandleNotificationQueueItems(ctx context.Context, task *asynq.Task) error {
	var p dto.PostNotificationDTO
	if err := json.Unmarshal(task.Payload(), &p); err != nil {
		return err
	}

	// ✅ Your custom logic here:
	log.Printf("🔔 Sending notification to user %d: %s \n %s", p.Id, p.Title, p.Description)

	return nil // returning nil means success
}
