/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"scheduler/setup/db"

	_ "github.com/lib/pq"
)

func main() {
	db.CreateDB() //call the createDB function from db package to create DB and schema.
}
