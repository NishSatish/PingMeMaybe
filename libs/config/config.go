package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var config *viper.Viper

// Use this fn in the main to load in the env vars
func LoadEnv(path string) {
	// these 3 lines right here to build the path to the .env file
	// for ex if i pass "." as the arg, it will look at "./.env"
	config = viper.New()
	config.AddConfigPath(path)
	config.SetConfigName("local")
	config.SetConfigType("env")

	config.AutomaticEnv() // For using injected vars in the docker container

	err := config.ReadInConfig()
	if err != nil {
		fmt.Println("Cannot find local.env, switching to automatic env mode")
	}
	return
}

// GetConfig returns the config
func GetConfig() *viper.Viper {
	return config
}
