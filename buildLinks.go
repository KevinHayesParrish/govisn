package main

import (
	"database/sql"
	"fmt"

	//"log"
	//"math"
	_ "github.com/mattn/go-sqlite3"
)

//buildLinksVersion is the file version sequence number
const buildLinksVersion = "0.0.1"

func buildLinks(debugFlag bool, database *sql.DB) *sql.DB {
	fmt.Println("func buildLinks version", buildLinksVersion, "started")

	fmt.Println("func buildLinks version", buildLinksVersion, "stopped")
	return database
}
