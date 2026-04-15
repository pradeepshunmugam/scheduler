package dbconn

import (
	"database/sql"
	"fmt"
)

func DBInsert(name, url string) {
	connString := `host=localhost port=5432 dbname=scheduler user=postgres password=admin sslmode=disable`
	db, err := sql.Open("postgres", connString)
	insert := fmt.Sprintf(`INSERT INTO url_list(name, url) VALUES('%v','%v')`, name, url)
	if err != nil {
		fmt.Println("Unable to connect to db")
	} else {
		_, err := db.Exec(insert)
		if err != nil {
			fmt.Println("Unable to inser")
		} else {
			fmt.Println("URL added into the tool.")
		}
	}

}
