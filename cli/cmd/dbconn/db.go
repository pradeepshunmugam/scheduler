package dbconn

import (
	"database/sql"
	"fmt"
)

//Update the URL and respective details in DB - scheduler.

func DBInsert(name, url, cron, sample, email string) {
	connString := `host=localhost port=5432 dbname=scheduler user=postgres password=admin sslmode=disable`
	db, err := sql.Open("postgres", connString)
	insertQuery := `INSERT INTO url_list(name, url, cron, sample, email) VALUES($1, $2, $3, $4, $5)`
	if err != nil {
		fmt.Println("Unable to connect to db")
	} else {
		_, err := db.Exec(insertQuery, name, url, cron, sample, email)
		if err != nil {
			fmt.Println("Unable to insert", err)
		} else {
			fmt.Println("Data inserted", name, url, cron, sample, email)
		}
	}

}
