package models

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Notification struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	// This field must only be a marshal of the NotificationPayload struct
	Payload       []byte             `json:"payload"`
	ChannelID     int                `json:"channel_id"`
	TransactionId string             `json:"transaction_id"`
	Status        NotificationStatus `json:"status"`
}

type NotificationPayload struct {
	Link string `json:"link"`
}

type NotificationStatus string

const (
	NotificationStatusProcessing NotificationStatus = "PROCESSING"
	NotificationStatusSuccess    NotificationStatus = "SUCCESS"
	NotificationStatusFailed     NotificationStatus = "FAILED"
)

type NotificationRepo struct {
	DB *pgx.Conn
}

type INotificationRepository interface {
	CreateNotification(ctx context.Context, notification Notification) (int, error)
	GetNotificationByID(ctx context.Context, id int) (*Notification, error)
}

func NewNotificationRepo(db *pgx.Conn) INotificationRepository {
	return &NotificationRepo{
		DB: db,
	}
}

func (r *NotificationRepo) CreateNotification(ctx context.Context, notification Notification) (int, error) {
	var id int
	query := `INSERT INTO notifications (title, description, payload, transaction_id, status) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.DB.QueryRow(ctx,
		query,
		notification.Title,
		notification.Description,
		notification.Payload,
		notification.TransactionId,
		notification.Status).Scan(&id)
	if err != nil {
		fmt.Println("error saving notification", err)
		return 0, err
	}
	return id, nil
}

func (r *NotificationRepo) GetNotificationByID(ctx context.Context, id int) (*Notification, error) {
	query := `SELECT id, title, description, user_id, channel_id, transaction_id 
			  FROM notifications WHERE id = $1`
	row := r.DB.QueryRow(ctx, query, id)

	var notification Notification
	err := row.Scan(
		&notification.ID,
		&notification.Title,
		&notification.Description,
		&notification.ChannelID,
		&notification.TransactionId)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}
