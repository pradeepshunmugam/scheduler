package main

import (
	"database/sql"
	"fmt"
	"os"
	"scheduler/logger"
	"scheduler/scheduler/scheduler"
	"sync"

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
	var wg sync.WaitGroup
	wg.Add(2)
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME_SCHEDULER")

	connString := fmt.Sprintf(`host = %v port = %v dbname = %v user = %v password = %v sslmode = disable`, dbhost, dbport, dbname, dbuser, dbpassword)
	dbconn, err := sql.Open("postgres", connString)
	if err != nil {
		logger.Log.Error("Unable to connect to DB to read the urls!!, ", zap.Error(err))
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic")
		}
	}()
	if err != nil {
		logger.Log.Error("Unableto read the urls!!, ", zap.Error(err))
		panic(err)
	}
	go scheduler.UpdateNextRun(dbconn, &wg)
	go scheduler.Run(dbconn, &wg)
	wg.Wait()

}
