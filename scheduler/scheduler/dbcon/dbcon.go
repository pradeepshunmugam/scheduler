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
	currrentTime := time.Now()
	nextformatted := next.Format("2006-01-02 15:04:05")
	if next.Before(currrentTime) {
		_, err = db.Exec(query, next, name)
		if err != nil {
			logger.Log.Error("Unable to insert. ", zap.Error(err))
		}
		logger.Log.Info("next run updated for job ", zap.String("name :", name), zap.String("next run : ", nextformatted))
	} else {
		logger.Log.Info("next run is already ahead of current time.", zap.String("name", name), zap.Time("Time : ", next), zap.Time("current time : ", currrentTime))
	}

}

// next time.Time, name , url string, sample int
func CheckNextRun(job jobs) {
	currrentTime := time.Now()
	current_minute := currrentTime.Minute()
	next_run_minute := next.Minute()
	nextformatted := next.Format("2006-01-02 15:04:05")
	scheduledJobs := []jobs
	if (current_minute - next_run_minute) < 5 {
		logger.Log.Info("Next run time is less than 5 min.", zap.String("name :", job.name), zap.String("next run : ", nextformatted))
		// fmt.Println("current minute  ", currrentTime.Minute())
		// fmt.Println("next run minute ", next.Minute())
	}
}
