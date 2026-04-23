package dbcon

import (
	"database/sql"
	"fmt"
	"os"
	"scheduler/logger"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func DBActivity(query string, next time.Time, name string) {
	// DB_HOST=localhost
	// DB_PORT=5432
	// DB_USER=postgres
	// DB_PASSWORD=admin
	// DB_NAME_DEFAULT=postgres
	// DB_NAME_SCHEDULER=scheduler
	err := godotenv.Load("../.env")
	if err != nil {
		logger.Log.Error("Unable to load env file", zap.Error(err))
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	schedulerDB := os.Getenv("DB_NAME_SCHEDULER")
	schedulerDBConn := fmt.Sprintf(`host=%v port=%v user=%v password=%v dbname=%v sslmode=disable`, host, port, user, password, schedulerDB)

	db, err := sql.Open("postgres", schedulerDBConn)
	if err != nil {
		logger.Log.Error("Unable to connect to DB. ", zap.Error(err))
	}
	defer db.Close()
	//update the next run.
	if next.Before(time.Now()) {
		_, err = db.Exec(query, next, name)
		if err != nil {
			logger.Log.Error("Unable to insert. ", zap.Error(err))
		}
		logger.Log.Info("next run updated for job ", zap.String("name :", name), zap.Time("next run : ", next))
	} else {
		logger.Log.Info("next run is already aahead of current time.", zap.String("name", name), zap.Time("Time : ", next))
	}
}
