package notifications

import (
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
	asynq *asynq.Client
}

type NotificationsServiceInterface interface {
	QueueNotification(ctx *gin.Context)
}

// Constructor
func NewNotificationsService(asynq *asynq.Client) NotificationsServiceInterface {
	return &notificationsService{
		asynq,
	}
}

func (n *notificationsService) QueueNotification(ctx *gin.Context) {
	fmt.Println("THE INCOMING", n.asynq)
	//defer n.asynq.Close()
	var notif dto.PostNotificationDTO
	err := ctx.BindJSON(&notif)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": notif})
		return
	}
	payload, err := json.Marshal(dto.PostNotificationDTO{
		Id:          notif.Id,
		Title:       notif.Title,
		Description: notif.Description,
		Link:        notif.Link,
	})
	fmt.Println("I GAAT IT! I GAAT IT!", payload)
	task := asynq.NewTask(messagePatterns.DispatchNotification, payload)

	info, err := n.asynq.Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": notif})
		return
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "task_id": info.ID, "queue": info.Queue})
}
