/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"scheduler/setup/db"
	"scheduler/setup/logger"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	err := godotenv.Load("../.env")

	if err != nil {
		//logger.Log.Error("Unable to load env", zap.Error(err))
		logger.Log.Error("Unable to load env", zap.Error(err))
	}
	db.CreateDB() //call the createDB function from db package to create DB and schema.
}
