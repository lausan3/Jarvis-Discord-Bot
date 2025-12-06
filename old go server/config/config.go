package config

import (
	"main/infra/logger"

	"github.com/spf13/viper"
)

type Configuration struct {
	Bot BotConfiguration
}

func Setup() error {
	var configuration *Configuration

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("Error reading config file, %s", err.Error())
		return err
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		logger.Errorf("Could not decode: %v", err)
		return err
	}

	return nil
}
