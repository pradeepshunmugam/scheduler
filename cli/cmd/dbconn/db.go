package dbconn

import (
	"database/sql"
	"scheduler/cli/logger"

	"go.uber.org/zap"
)

//Update the URL and respective details in DB - scheduler.

func DBInsert(name, url, cron, sample, email string) {
	connString := `host=localhost port=5432 dbname=scheduler user=postgres password=admin sslmode=disable`
	db, err := sql.Open("postgres", connString)
	insertQuery := `INSERT INTO url_list(name, url, cron, sample, email) VALUES($1, $2, $3, $4, $5)`
	if err != nil {
		logger.Log.Error("Unable to connect to db", zap.Error(err))
	} else {
		_, err := db.Exec(insertQuery, name, url, cron, sample, email)
		if err != nil {
			logger.Log.Error("Unable to insert", zap.Error(err))
		} else {
			logger.Log.Info("Data inserted", zap.String("name", name), zap.String("url", url), zap.String("cron", cron), zap.String("sample", sample), zap.String("email", email))
		}
	}

}
