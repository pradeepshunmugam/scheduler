package dbconn

import (
	"database/sql"
	"fmt"
	"os"
	"scheduler/cli/logger"

	"go.uber.org/zap"
)

//Update the URL and respective details in DB - scheduler.

func DBInsert(name, url, cron, sample, email string) string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	Dbname := os.Getenv("DB_NAME_SCHEDULER")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	connString := fmt.Sprintf(`host = %v port = %v  dbname = %v user = %v password = %v sslmode=disable`, host, port, Dbname, user, password)

	//connString := `host=localhost port=5432 dbname=scheduler user=postgres password=admin sslmode=disable`
	db, err := sql.Open("postgres", connString)
	insertQuery := `INSERT INTO url_list(name, url, cron, sample, email) VALUES($1, $2, $3, $4, $5)`
	if err != nil {
		logger.Log.Error("Unable to connect to db", zap.Error(err))
		return "Unable to add"
	} else {
		_, err := db.Exec(insertQuery, name, url, cron, sample, email)
		if err != nil {
			logger.Log.Error("Unable to insert", zap.Error(err))
			return "Unable to add"
		} else {
			logger.Log.Info("Data inserted", zap.String("name", name), zap.String("url", url), zap.String("cron", cron), zap.String("sample", sample), zap.String("email", email))
			return "Added succesfully"
		}
	}
	return ""

}

func DeleteRow(name string) string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	Dbname := os.Getenv("DB_NAME_SCHEDULER")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	connString := fmt.Sprintf(`host = %v port = %v  dbname = %v user = %v password = %v sslmode=disable`, host, port, Dbname, user, password)

	//connString := `host=localhost port=5432 dbname=scheduler user=postgres password=admin sslmode=disable`
	db, err := sql.Open("postgres", connString)
	if err != nil {
		logger.Log.Error("Unable to connnect to DB to delete the url.", zap.Error(err))
		return "Unable to delete the url."
	}
	defer db.Close()
	validateQuery := `SELECT count(*) FROM url_list where name = $1`
	res, err := db.Query(validateQuery, name)
	if err != nil {
		logger.Log.Error("Unable to connect the DB to check the value availabilty", zap.Error(err))
	}
	var result int
	for res.Next() {
		err := res.Scan(&result)
		if err != nil {
			logger.Log.Error("Result in scanning the output of number of count", zap.Error(err))
		} else {
			if result > 0 {
				deleteQuery := `DELETE FROM url_list where name = $1`
				_, err = db.Exec(deleteQuery, name)
				if err != nil {
					logger.Log.Info("unable to delete the value.", zap.Error(err))
					return "Unable to Delete"
				} else {
					logger.Log.Info("Deleted succefully.", zap.String("name", name))
					return "Deleted Succesfully"
				}

			} else {
				return "Name is not availble"
			}
		}
	}
	return ""

}
