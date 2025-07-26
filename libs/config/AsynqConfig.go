package config

import "github.com/hibiken/asynq"

func GetAsynqClient() *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"})
}
