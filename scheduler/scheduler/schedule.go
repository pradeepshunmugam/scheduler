package scheduler

import (
	"database/sql"
	"fmt"
	"os"
	"scheduler/logger"

	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

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
	defer dbconn.Close()
	rows, err := dbconn.Query("SELECT name, url, sample FROM url_list;")
	if err != nil {
		logger.Log.Error("Unableto read the urls!!, ", zap.Error(err))
	}
	for rows.Next() {
		var name string
		var url string
		var sample int
		err := rows.Scan(&name, &url, &sample)
		if err != nil {
			logger.Log.Error("Unableto scan the urls!!, ", zap.Error(err))
		}
		logger.Log.Info("Read", zap.String("name", name), zap.String("url", url), zap.Int("sample", sample))
	}
	defer rows.Close()

}
