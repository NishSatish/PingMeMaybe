package cron

import (
	"PingMeMaybe/libs/db"
	"PingMeMaybe/libs/utils"
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"golang.org/x/sync/errgroup"
	"log"
)

type MarkFailuresCron struct {
	db *db.DBService
}

type MarkFailuresCronInterface interface {
	// StartMarkFailuresCron To keep scanning the db and mark transactions that are stuck in processing as failed
	StartMarkFailuresCron()
}

func NewMarkFailuresCron(db *db.DBService) MarkFailuresCronInterface {
	return &MarkFailuresCron{
		db,
	}
}

func (m *MarkFailuresCron) StartMarkFailuresCron() {
	cronJob := cron.New()
	_, err := cronJob.AddFunc("@every 10s", func() {

		notifications, err := m.db.Notifications.GetPendingNotifications(context.Background())
		if err != nil {
			log.Fatal("Could not fecth notifications:", err)
		}
		if len(notifications) == 0 {
			fmt.Println("No notifications found")
			return
		}

		g := new(errgroup.Group)

		for _, n := range notifications {
			notification := n // capture range variable
			if utils.IsOlderThanOneDay(notification.CreatedAt) {
				g.Go(func() error {
					fmt.Printf("Marking notification %d \n", notification.TransactionId)
					err := m.db.Notifications.MarkNotificationAsFailed(context.Background(), notification.TransactionId)
					if err != nil {
						log.Printf("Failed to update notification %d status: %v\n", notification.Title, err)
					} else {
						fmt.Printf("Notification %d marked as failed successfully\n", notification.Title)
					}
					return err
				})
			}
		}

		if err := g.Wait(); err != nil {
			log.Println("some notifications failed to update.")
		} else {
			log.Println("All old notifications marked as failed successfully.")
		}
	})
	if err != nil {
		log.Fatal("Failed to start cron job:", err)
		return
	}
	cronJob.Start()
}
