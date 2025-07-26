package service

import (
	"PingMeMaybe/libs/db/models"
	"PingMeMaybe/libs/dto"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"log"
)

type notificationProcessorService struct {
	db *pgx.Conn
}

type INotificationProcessorService interface {
	HandleNotificationQueueItems(ctx context.Context, task *asynq.Task) error
}

func NewNotificationProcessorService(db *pgx.Conn) INotificationProcessorService {
	return &notificationProcessorService{
		db,
	}
}

func (n notificationProcessorService) HandleNotificationQueueItems(ctx context.Context, task *asynq.Task) error {
	var p dto.PostNotificationDTO
	query := `UPDATE notifications SET status = $1 WHERE transaction_id = $2`

	if err := json.Unmarshal(task.Payload(), &p); err != nil {
		_, err := n.db.Exec(ctx, query, models.NotificationStatusFailed, task.ResultWriter().TaskID())
		return err
	}

	// If the payload is successfully extracted, then mark notification as successful (for now)
	_, err := n.db.Exec(ctx, query, models.NotificationStatusSuccess, task.ResultWriter().TaskID())
	if err != nil {
		fmt.Println("Error updating notification status:", err)
		return err
	}
	log.Printf("🔔 Sending notification to user %d: %s \n %s", p.Id, p.Title, p.Description)

	return nil // returning nil means success
}
