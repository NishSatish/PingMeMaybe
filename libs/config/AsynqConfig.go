package config

import "github.com/hibiken/asynq"

func GetAsynqClient() *asynq.Client {
	LoadEnv(".")

	return asynq.NewClient(asynq.RedisClientOpt{
		Addr:     GetConfig().GetString("REDIS_CLUSTER"),
		Username: GetConfig().GetString("REDIS_USERNAME"),
		Password: GetConfig().GetString("REDIS_PASSWORD"),
	})
}
