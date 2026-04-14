/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"scheduler/cmd"
	"scheduler/db"

	_ "github.com/lib/pq"
)

func main() {
	cmd.Execute()
	db.CreateDB()
}
