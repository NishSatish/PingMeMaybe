package notifications

import (
	"PingMeMaybe/libs/db/models"
	"PingMeMaybe/libs/dto"
	"PingMeMaybe/libs/messagePatterns"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"log"
	"net/http"
	"time"
)

type notificationsService struct {
	// can be left empty also just for the sake of interface implementation
	asynq                  *asynq.Client
	notificationRepository models.INotificationRepository
}

type NotificationsServiceInterface interface {
	QueueNotification(ctx *gin.Context)
}

// Constructor
func NewNotificationsService(asynq *asynq.Client, notificationsRepository models.INotificationRepository) NotificationsServiceInterface {
	return &notificationsService{
		asynq,
		notificationsRepository,
	}
}

func (n *notificationsService) QueueNotification(ctx *gin.Context) {
	//defer n.asynq.Close()
	var notif dto.PostNotificationDTO
	err := ctx.BindJSON(&notif)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": notif})
		return
	}
	payload, err := json.Marshal(dto.PostNotificationDTO{
		Title:       notif.Title,
		Description: notif.Description,
		Link:        notif.Link,
	})
	task := asynq.NewTask(messagePatterns.DispatchNotification, payload)

	info, err := n.asynq.Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": notif})
		return
	}

	// Save the notification trigger entry in postgres
	notificationPayload, err := json.Marshal(models.NotificationPayload{Link: notif.Link})
	notificationObject := models.Notification{
		Title:         notif.Title,
		Description:   notif.Description,
		Payload:       notificationPayload,
		Status:        models.NotificationStatusProcessing,
		TransactionId: info.ID,
	}
	id, err := n.notificationRepository.CreateNotification(ctx, notificationObject)
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "task_id": info.ID, "queue": info.Queue, "notification_id": id, "payload": payload})
}
