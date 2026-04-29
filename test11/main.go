package main

import (
	"scheduler/test11/db"
	"scheduler/test11/query"
)

func main() {
	db := db.GetDB()
	defer db.Close()
	query.SelectQuery()

}
