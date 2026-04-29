package db

import (
	"database/sql"
	"fmt"
	"os"
	"scheduler/cli/logger"
	"sync"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var db *sql.DB
var once sync.Once

func initDB() {
	err := godotenv.Load("../.env")
	if err != nil {
		logger.Log.Error("Unable to load environment. ", zap.Error(err))
		return
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME_SCHEDULER")
	password := os.Getenv("DB_PASSWORD")
	connString := fmt.Sprintf(`host = %v port = %v  dbname = %v user = %v password = %v sslmode=disable`, host, port, dbname, user, password)

	db, err = sql.Open("postgres", connString)
	if err != nil {
		fmt.Println("Unable to connect to DB", err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Unable to ping DB", err)
	}
}

func GetDB() *sql.DB {
	once.Do(initDB)
	return db
}
