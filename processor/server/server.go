package server

import (
	"PingMeMaybe/libs/config"
	"PingMeMaybe/libs/messagePatterns"
	"PingMeMaybe/processor/pkg/service"
	_ "encoding/json"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func StartAsynqServer(dbConn *pgxpool.Pool) {
	// This is a background processor, wont be exposed by HTTP
	services := service.NewProcessorServices(dbConn)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     config.GetConfig().GetString("REDIS_CLUSTER"),
			Username: config.GetConfig().GetString("REDIS_USERNAME"),
			Password: config.GetConfig().GetString("REDIS_PASSWORD"),
		},
		asynq.Config{
			Concurrency: 10,
			// Priorities
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	// Register handlers with msg patterns
	mux.HandleFunc(messagePatterns.DispatchNotification, services.HandleNotificationQueueItems)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
