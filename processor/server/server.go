package server

import (
	"PingMeMaybe/libs/messagePatterns"
	"PingMeMaybe/processor/pkg/service"
	_ "encoding/json"
	"github.com/hibiken/asynq"
	"log"
)

func StartServer() {
	startMuxServer()
}

func startMuxServer() {
	// This is a background processor, wont be exposed by HTTP
	services := service.NewProcessorServices()

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
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
