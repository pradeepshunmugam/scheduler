/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"scheduler/cli/cmd"
	"scheduler/cli/logger"

	_ "github.com/lib/pq"
)

func main() {
	logger.Init()
	//logger.Log.Info("Initializing main function.")
	cmd.Execute()
}
