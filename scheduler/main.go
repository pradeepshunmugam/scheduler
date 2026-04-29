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

	// err := godotenv.Load("../.env")
	// if err != nil {
	// 	logger.Log.Error("Unable to load environment. ", zap.Error(err))
	// 	return
	// }
	// dbhost := os.Getenv("DB_HOST")
	// dbport := os.Getenv("DB_PORT")
	// dbuser := os.Getenv("DB_USER")
	// dbpassword := os.Getenv("DB_PASSWORD")
	// dbname := os.Getenv("DB_NAME_SCHEDULER")

	// connString := fmt.Sprintf(`host = %v port = %v dbname = %v user = %v password = %v sslmode = disable`, dbhost, dbport, dbname, dbuser, dbpassword)
	// dbconn, err := sql.Open("postgres", connString)
	// if err != nil {
	// 	logger.Log.Error("Unable to connect to DB to read the urls!!, ", zap.Error(err))
	// }
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("Recovered from panic")
	// 	}
	// }()
	// if err != nil {
	// 	logger.Log.Error("Unableto read the urls!!, ", zap.Error(err))
	// 	panic(err)
	// }
	// go scheduler.UpdateNextRun(dbconn, &wg)
	// go scheduler.Run(dbconn, &wg)
	// wg.Wait()

	//To make the serive run continuously
	db := db.GetDB()
	defer db.Close()
	for {
		scheduler.UpdateNextRun()
		scheduler.Run()
		fmt.Println("Waiting for a minute")
		time.Sleep(60 * time.Second)

	}

}
