/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"scheduler/cli/cmd"
	"scheduler/cli/db"
	"scheduler/cli/logger"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	err := godotenv.Load("../.env")
	if err != nil {
		logger.Log.Error("Unable to load env", zap.Error(err))
	}
	cmd.Execute()
	db := db.GetDB()
	defer db.Close()
}
