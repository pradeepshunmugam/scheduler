package db

import (
	"fmt"
	"scheduler/cli/logger"

	"go.uber.org/zap"
)

//Update the URL and respective details in DB - scheduler.

func DBInsert(name, url, cron, sample, email string) {
	db := GetDB()
	insertQuery := `INSERT INTO url_list(name, url, cron, sample, email) VALUES($1, $2, $3, $4, $5)`
	_, err := db.Exec(insertQuery, name, url, cron, sample, email)
	if err != nil {
		logger.Log.Error("Unable to insert", zap.Error(err))
		fmt.Println("Unable to add")
	} else {
		logger.Log.Info("Data inserted", zap.String("name", name), zap.String("url", url), zap.String("cron", cron), zap.String("sample", sample), zap.String("email", email))
		fmt.Println("Added succesfully")
	}
}

//Delete the url which is not needed.

func DeleteRow(name string) {
	db := GetDB()
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
					fmt.Println("Unable to Delete")
				} else {
					logger.Log.Info("Deleted succefully.", zap.String("name", name))
					fmt.Println("Deleted Succesfully")
				}

			} else {
				fmt.Println("Name is not availble")
			}
		}
	}

}

func InsertUser(username, hashedPwdStr, email string) string {
	db := GetDB()
	insertUser := `Insert INTO users(username, password, email) VALUES($1, $2, $3)`
	_, err := db.Exec(insertUser, username, hashedPwdStr, email)
	if err != nil {
		logger.Log.Error("Unable to insert user", zap.Error(err))
		return "unable to create user"
	}
	logger.Log.Info("user inserted sucessfully", zap.String("user", username))
	return "User created"

}
