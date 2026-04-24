package scheduler

import (
	"database/sql"
	"fmt"
	"os"
	"scheduler/logger"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Jobs struct {
	url      string
	name     string
	sample   int
	cron     string
	email    string
	last_run string
	next_run string
}

func ReadURL() {
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
	rows, err := dbconn.Query("SELECT name, url, sample, cron, email, last_run, next_run FROM url_list;")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic")
		}
	}()
	if err != nil {
		logger.Log.Error("Unableto read the urls!!, ", zap.Error(err))
		panic(err)
	}
	for rows.Next() {
		var job Jobs
		err := rows.Scan(&job.name, &job.url, &job.sample, &job.cron, &job.email, &job.last_run, &job.next_run)
		defer rows.Close()
		if err != nil {
			logger.Log.Error("Unableto scan the urls!!, ", zap.Error(err))
		}
		log(job)
	}

}
func log(job Jobs) {
	//logger.Log.Info("Read", zap.String("name", job.name), zap.String("url", job.url), zap.String("cron", job.cron), zap.String("email", job.email), zap.String("last_run", job.last_run), zap.String("next_run", job.next_run))

	calculateNextRun(job)
}
