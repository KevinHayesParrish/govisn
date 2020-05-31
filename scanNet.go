package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

/*
 * TODO:
 	*
*/

// SCANNETVERSION is the file version number
const SCANNETVERSION = "0.0.1"

func scanNet(debugFlag bool, seed string, community string, database *sql.DB) []ScannedRouter {

	fmt.Println("\nfunc scanNet version", SCANNETVERSION, "started.")
	if debugFlag {
		fmt.Println("seed=", seed, "community=", community)
	}

	var scannedRouters []ScannedRouter

	if debugFlag {
		fmt.Println("Returning, scannedRouters=", scannedRouters)
	}
	fmt.Println("func scanNet", SCANNETVERSION, "ended.")
	return scannedRouters
}
