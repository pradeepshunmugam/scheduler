package main

import (
	"fmt"
	"scheduler/scheduler/db.go"
	"scheduler/scheduler/logger"
	"scheduler/scheduler/scheduler"
	"time"
)

func main() {
	logger.Init()
	db := db.GetDB()
	defer db.Close()
	for {
		//scheduler.UpdateNextRun()
		scheduler.Run()
		fmt.Println("Waiting for a minute")
		time.Sleep(60 * time.Second)

	}

}
