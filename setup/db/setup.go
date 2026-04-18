package db

import (
	"database/sql"
	"fmt"
	"os"
	"scheduler/cli/logger"

	"go.uber.org/zap"
)

//Create DB - scheduler while initializing the tool.

func CreateDB() {
	var db string
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	defaultDbname := os.Getenv("DB_NAME_DEFAULT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	connString := fmt.Sprintf(`host = %v port = %v  dbname = %v user = %v password = %v sslmode=disable`, host, port, defaultDbname, user, password)
	defaultDb, err := sql.Open("postgres", connString)
	if err != nil {
		logger.Log.Error("unable to connect to default db - postgres", zap.Error(err))
		return
	}
	defer defaultDb.Close()
	db = "scheduler"
	var exists bool
	checkDb := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`
	err = defaultDb.QueryRow(checkDb, db).Scan(&exists)
	logger.Log.Info("Checking db exists or not. ")
	if err != nil {
		logger.Log.Error("unable to check db exist", zap.Error(err))
		return
	}
	if !exists {
		logger.Log.Info("DB - %v does not exist.\n", zap.String("db", db))
		createQuery := `Create DATABASE `
		defaultDb.Exec(createQuery + db)
		logger.Log.Info("DB - %v created.\n", zap.String("db", db))
		createSchema_url_list(db)
		createSchema_user(db)
	} else {
		logger.Log.Info("DB already exists. Checking for schemas.")
		createSchema_url_list(db)
		createSchema_user(db)
	}

}

//Create schema - "url_list" to store url details.

func createSchema_url_list(db string) {
	var exists bool
	tableName := "url_list"
	connString := fmt.Sprintf("host=localhost port=5432 user=postgres password=admin dbname=%v sslmode=disable", db)
	dbConn, err := sql.Open("postgres", connString)
	if err != nil {
		logger.Log.Error("Unable to connect DB", zap.Error(err))

	}
	defer dbConn.Close()
	checkUrlList := `SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' and table_name= $1)`
	err = dbConn.QueryRow(checkUrlList, tableName).Scan(&exists)
	if err != nil {
		//fmt.Println("Unable to check table.", err)
		logger.Log.Error("Unable to check table", zap.Error(err))

	}
	if !exists {
		_, err = dbConn.Exec(`Create TABLE IF NOT EXISTS url_list (
		id SERIAL PRIMARY KEY,
		url TEXT NOT NULL,
		name TEXT UNIQUE,
		cron TEXT NOT NULL,
		sample TEXT NOT NULL,
        email TEXT ,
		created_at TIMESTAMPTZ DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
		last_run TIMESTAMPTZ,
		next_run TIMESTAMPTZ );`)

		if err != nil {
			//fmt.Println("Unable to create table.", err)
			logger.Log.Error("Unable to create table", zap.Error(err))

		} else {
			logger.Log.Info("Table created")
		}
	} else {
		logger.Log.Info("Table already exists.")
	}
}

//Create schema - "users" for user credential details

func createSchema_user(db string) {
	var exists bool
	tableName := "users"
	connString := fmt.Sprintf("host=localhost port=5432 user=postgres password=admin dbname=%v sslmode=disable", db)
	dbConn, err := sql.Open("postgres", connString)
	if err != nil {
		//fmt.Printf("Unable to connect to DB %v\n", err)
		logger.Log.Error("Unable to connect to DB", zap.Error(err))

	}
	defer dbConn.Close()
	checkUser := `SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' and table_name= $1)`
	err = dbConn.QueryRow(checkUser, tableName).Scan(&exists)
	if err != nil {
		//fmt.Println("Unable to check table.", err)
		logger.Log.Error("Unable to check table", zap.Error(err))

	}
	if !exists {
		query := `CREATE TABLE IF NOT EXISTS users (
	id SERIAL,
	username TEXT PRIMARY KEY,
	password TEXT NOT NULL,
	email TEXT NOT NULL)`
		_, err = dbConn.Exec(query)
		if err != nil {
			//fmt.Println("Unable to create user table.", err)
			logger.Log.Error("Unable to create user table", zap.Error(err))

		} else {
			logger.Log.Info("User table created")
		}
	} else {
		logger.Log.Info("Table exists.")
	}

}
