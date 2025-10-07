package main

import (
	"main/bot"
	"main/config"
	"main/infra/logger"
)

func main() {
	if err := config.Setup(); err != nil {
		logger.Fatalf("Error setting up config: %s", err.Error())
	}

	bot.Start()
}
