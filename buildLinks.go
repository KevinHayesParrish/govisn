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

	/* TODO
	 *
	 * Populate RouterName, DestinationName, DestinationIP, NextHopName and NextHopIP from RouteTable elements.
	 * calculate LinkID using CRC of RouterName/DestinationIP/NextHopIP
	 * Write Links Row to database
	 *
	 */

	fmt.Println("func buildLinks version", buildLinksVersion, "stopped")
	return database
}
