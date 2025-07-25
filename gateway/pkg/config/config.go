package config

import "github.com/spf13/viper"

var config *viper.Viper

// Use this fn in the main to load in the env vars
func LoadEnv(path string) {
	// these 3 lines right here to build the path to the .env file
	// for ex if i pass "." as the arg, it will look at "./.env"
	config = viper.New()
	config.AddConfigPath(path)
	config.SetConfigName("local")
	config.SetConfigType("env")

	config.AutomaticEnv()

	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return
}

// GetConfig returns the config
func GetConfig() *viper.Viper {
	return config
}
