package main

import (
	"fmt"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

func testarango() {

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://server:8529"},
	})
	if err != nil {
		// Handle error
	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
	})
	if err != nil {
		// Handle error
	}

	// Open "examples_books" database
	db, err := client.Database(nil, "examples_books")
	if err != nil {
		// Handle error
	}

	// Open "books" collection
	col, err := db.Collection(nil, "books")
	if err != nil {
		// Handle error
	}

	// Create document
	book := Book{
		Title:   "ArangoDB Cookbook",
		NoPages: 257,
	}
	meta, err := col.CreateDocument(nil, book)
	if err != nil {
		// Handle error
	}
	fmt.Printf("Created document in collection '%s' in database '%s'\n", col.Name(), db.Name())

}
