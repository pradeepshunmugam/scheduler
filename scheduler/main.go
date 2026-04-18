package main

import (
	"scheduler/logger"
	"scheduler/scheduler/scheduler"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	logger.Init()

	err := godotenv.Load("../.env")
	if err != nil {
		logger.Log.Error("Unable to load environment. ", zap.Error(err))
		return
	}
	scheduler.ReadURL()
}
