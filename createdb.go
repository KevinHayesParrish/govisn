package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

/*
func createdb() {
	database, _ := sql.Open("sqlite3", "./nraboy.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	statement.Exec("Nic", "Raboy")
	rows, _ := database.Query("SELECT id, firstname, lastname FROM people")
	var id int
	var firstname string
	var lastname string
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname)
		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	}
}
*/

func main() {
	database, _ := sql.Open("sqlite3", "./hashed.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, router1 TEXT, router2 TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO people (router1, router2) VALUES (?, ?)")
	statement.Exec("media", "router")
	rows, _ := database.Query("SELECT id, router1, router2 FROM people")
	var id int
	var router1 string
	var router2 string
	for rows.Next() {
		rows.Scan(&id, &router1, &router2)
		fmt.Println(strconv.Itoa(id) + ": " + router1 + " " + router2)
	}
}
