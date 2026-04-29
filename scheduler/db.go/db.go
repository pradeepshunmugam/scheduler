package db

import (
	"database/sql"
	"fmt"
	"os"
	"scheduler/scheduler/logger"
	"sync"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var db *sql.DB
var once sync.Once

func initDb() {
	var err error
	err = godotenv.Load("../.env")
	if err != nil {
		logger.Log.Error("Unable to load environment. ", zap.Error(err))
		return
	}
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME_SCHEDULER")

	connString := fmt.Sprintf(`host = %v port = %v dbname = %v user = %v password = %v sslmode = disable`, dbhost, dbport, dbname, dbuser, dbpassword)
	db, err = sql.Open("postgres", connString)
	if err != nil {
		logger.Log.Error("Unable to connect to DB to read the urls!!, ", zap.Error(err))
	}
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("Recovered from panic")
	// 	}
	// }()
	// if err != nil {
	// 	logger.Log.Error("Unableto read the urls!!, ", zap.Error(err))
	// 	panic(err)
	// }
	err = db.Ping()
	if err != nil {
		logger.Log.Error("Unable to ping the error", zap.Error(err))
	}
}

func GetDB() *sql.DB {
	once.Do(initDb)
	return db
}
