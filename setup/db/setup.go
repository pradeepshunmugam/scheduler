package db

import (
	"database/sql"
	"fmt"
)

func CreateDB() {
	var db string
	connString := `host = localhost port = 5432  dbname = postgres user = postgres password = admin sslmode=disable`
	defaultDb, err := sql.Open("postgres", connString)
	if err != nil {
		fmt.Printf("unable to connect to default db - postgres\n")
		return
	}
	defer defaultDb.Close()
	db = "scheduler"
	var exists bool
	checkDb := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`
	err = defaultDb.QueryRow(checkDb, db).Scan(&exists)
	fmt.Println("Checking db exists or not. ", err)
	if err != nil {
		fmt.Println("unable to check db exist", err)
		return
	}
	if !exists {
		fmt.Printf("DB - %v does not exist.\n", db)
		createQuery := `Create DATABASE `
		defaultDb.Exec(createQuery + db)
		fmt.Printf("DB - %v created.\n", db)
		createSchema_url_list(db)
		createSchema_user(db)
	} else {
		fmt.Println("DB already exists. Checking for schemas.")
		createSchema_url_list(db)
		createSchema_user(db)
	}

}

func createSchema_url_list(db string) {
	var exists bool
	tableName := "url_list"
	connString := fmt.Sprintf("host=localhost port=5432 user=postgres password=admin dbname=%v sslmode=disable", db)
	dbConn, err := sql.Open("postgres", connString)
	if err != nil {
		fmt.Printf("Unable to connect to DB %v\n", err)
	}
	defer dbConn.Close()
	checkUrlList := `SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' and table_name= $1)`
	err = dbConn.QueryRow(checkUrlList, tableName).Scan(&exists)
	if err != nil {
		fmt.Println("Unable to check table.", err)
	}
	if !exists {
		_, err = dbConn.Exec(`Create TABLE IF NOT EXISTS url_list (
		id SERIAL PRIMARY KEY,
		url TEXT NOT NULL,
		name TEXT UNIQUE,
		cron TEXT NOT NULL DEFAULT '*/5 * * * *',
		sample int NOT NULL DEFAULT 1,
        email TEXT ,
		created_at TIMESTAMPTZ DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
		last_run TIMESTAMPTZ,
		next_run TIMESTAMPTZ );`)

		if err != nil {
			fmt.Println("Unable to create table.", err)
		} else {
			fmt.Println("Table created")
		}
	} else {
		fmt.Println("Table already exists.")
	}
}

func createSchema_user(db string) {
	var exists bool
	tableName := "users"
	connString := fmt.Sprintf("host=localhost port=5432 user=postgres password=admin dbname=%v sslmode=disable", db)
	dbConn, err := sql.Open("postgres", connString)
	if err != nil {
		fmt.Printf("Unable to connect to DB %v\n", err)
	}
	defer dbConn.Close()
	checkUser := `SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' and table_name= $1)`
	err = dbConn.QueryRow(checkUser, tableName).Scan(&exists)
	if err != nil {
		fmt.Println("Unable to check table.", err)
	}
	if !exists {
		query := `CREATE TABLE IF NOT EXISTS users (
	id SERIAL,
	username TEXT PRIMARY KEY,
	password TEXT NOT NULL,
	email TEXT NOT NULL)`
		_, err = dbConn.Exec(query)
		if err != nil {
			fmt.Println("Unable to create user table.", err)
		} else {
			fmt.Println("User table created")
		}
	} else {
		fmt.Println("Table exists.")
	}

}
