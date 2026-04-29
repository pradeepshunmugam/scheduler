package db

import (
	"database/sql"
	"fmt"
	"os"
	"scheduler/cli/logger"
	"sync"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var db *sql.DB
var once sync.Once

func initDb() {
	var err error
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	Dbname := os.Getenv("DB_NAME_SCHEDULER")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	connString := fmt.Sprintf(`host = %v port = %v  dbname = %v user = %v password = %v sslmode=disable`, host, port, Dbname, user, password)
	db, err = sql.Open("postgres", connString)
	if err != nil {
		logger.Log.Error("Unable to connect to db", zap.Error(err))
		return
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Unable to ping the DB")
		return
	}
}

func GetDB() *sql.DB {
	once.Do(initDb)
	return db

}
