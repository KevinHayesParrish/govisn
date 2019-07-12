package main

import (
	"database/sql"
	"fmt"
	"hash/crc32"
	"strconv"

	//	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

//createdsampledbVersion is the file version number
const createsampledbVersion = "0.1.5"

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

func createsampledb() {
	fmt.Println("createdampledb version:", createsampledbVersion)
	database, _ := sql.Open("sqlite3", "./samplenetwork.db")
	/*
	 *	Add Routers table to DB
	 */
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS Routers (RouterID INTEGER NOT NULL PRIMARY KEY, SystemName TEXT, SystemDesc TEXT, UpTime TEXT, Contact TEXT, Location TEXT, GpsLat NUMERIC, GPSLong NUMERIC, GpsAlt NUMERIC)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO Routers (RouterID, SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")

	/*
	 *	Add RouteTable table to DB
	 */
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS RouteTable (RouterID INTEGER, DestAddr TEXT, NextHop TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO RouteTable (RouterID, DestAddr, NextHop) VALUES (?, ?, ?)")

	/*
	 *	Add RouterIP table to DB
	 */
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS RouterIp (RouterID INTEGER, IpAddr TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO RouterIp (RouterID, IpAddr) VALUES (?, ?)")

	/*
	 *	Add RouterMac table to DB
	 */
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS RouterMac (RouterID INTEGER NOT NULL, MacAddr TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO RouterMac (RourterID, MacAddr) VALUES (?, ?)")

	/*
	 *	Add Links table to DB
	 */
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS Links (LinkID INTEGER PRIMARY KEY, FromRouter TEXT, ToRouter TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO Links (LinkID, FromRouter, ToRouter) VALUES (?, ?, ?)")

	/*
	 * Add a set of Links to the database
	 */
	var dest string
	var nextHop string
	dest = "media"
	nextHop = "router"
	destToNextHopLinkStr := dest + nextHop // directional link from dest to nextHop
	nextHopToDestStr := nextHop + dest     // directional link from nextHop to dest

	// add direction link from dest to nextHop to the database
	destToNextHopLinkUint32 := crc32.ChecksumIEEE([]byte(destToNextHopLinkStr))
	statement.Exec(strconv.Itoa(int(destToNextHopLinkUint32)), dest, nextHop)

	// add direction link from nextHop to dest to the database
	nextHopToDestUint32 := crc32.ChecksumIEEE([]byte(nextHopToDestStr))
	statement.Exec(strconv.Itoa(int(nextHopToDestUint32)), nextHop, dest)

	/*
	 * Add another set of Links to the database
	 */
	dest = "router"
	nextHop = "wan router"
	// add direction link from dest to nextHop to the database
	destToNextHopLinkStr = dest + nextHop // directional link from dest to nextHop
	nextHopToDestStr = nextHop + dest     // directional link from nextHop to dest

	destToNextHopLinkUint32 = crc32.ChecksumIEEE([]byte(destToNextHopLinkStr))
	statement.Exec(strconv.Itoa(int(destToNextHopLinkUint32)), dest, nextHop)

	// add direction link from nextHop to dest to the database
	nextHopToDestUint32 = crc32.ChecksumIEEE([]byte(nextHopToDestStr))
	statement.Exec(strconv.Itoa(int(nextHopToDestUint32)), nextHop, dest)

	/*
	 * Add another set of Links to the database
	 */
	dest = "wan router"
	nextHop = "hub"
	// add direction link from dest to nextHop to the database
	destToNextHopLinkStr = dest + nextHop // directional link from dest to nextHop
	nextHopToDestStr = nextHop + dest     // directional link from nextHop to dest

	destToNextHopLinkUint32 = crc32.ChecksumIEEE([]byte(destToNextHopLinkStr))
	statement.Exec(strconv.Itoa(int(destToNextHopLinkUint32)), dest, nextHop)

	// add direction link from nextHop to dest to the database
	nextHopToDestUint32 = crc32.ChecksumIEEE([]byte(nextHopToDestStr))
	statement.Exec(strconv.Itoa(int(nextHopToDestUint32)), nextHop, dest)

	/*
	 * Add another set of Links to the database
	 */
	dest = "wan router"
	nextHop = "old-country-road"
	// add direction link from dest to nextHop to the database
	destToNextHopLinkStr = dest + nextHop // directional link from dest to nextHop
	nextHopToDestStr = nextHop + dest     // directional link from nextHop to dest

	destToNextHopLinkUint32 = crc32.ChecksumIEEE([]byte(destToNextHopLinkStr))
	statement.Exec(strconv.Itoa(int(destToNextHopLinkUint32)), dest, nextHop)

	// add direction link from nextHop to dest to the database
	nextHopToDestUint32 = crc32.ChecksumIEEE([]byte(nextHopToDestStr))
	statement.Exec(strconv.Itoa(int(nextHopToDestUint32)), nextHop, dest)

	/*
	 * Add another set of Links to the database
	 */
	dest = "wan router"
	nextHop = "fukui"
	// add direction link from dest to nextHop to the database
	destToNextHopLinkStr = dest + nextHop // directional link from dest to nextHop
	nextHopToDestStr = nextHop + dest     // directional link from nextHop to dest

	destToNextHopLinkUint32 = crc32.ChecksumIEEE([]byte(destToNextHopLinkStr))
	statement.Exec(strconv.Itoa(int(destToNextHopLinkUint32)), dest, nextHop)

	// add direction link from nextHop to dest to the database
	nextHopToDestUint32 = crc32.ChecksumIEEE([]byte(nextHopToDestStr))
	statement.Exec(strconv.Itoa(int(nextHopToDestUint32)), nextHop, dest)

	/*
	 * Add another set of Links to the database
	 */
	dest = "wan router"
	nextHop = "amsterdam"
	// add direction link from dest to nextHop to the database
	destToNextHopLinkStr = dest + nextHop // directional link from dest to nextHop
	nextHopToDestStr = nextHop + dest     // directional link from nextHop to dest

	destToNextHopLinkUint32 = crc32.ChecksumIEEE([]byte(destToNextHopLinkStr))
	statement.Exec(strconv.Itoa(int(destToNextHopLinkUint32)), dest, nextHop)

	// add direction link from nextHop to dest to the database
	nextHopToDestUint32 = crc32.ChecksumIEEE([]byte(nextHopToDestStr))
	statement.Exec(strconv.Itoa(int(nextHopToDestUint32)), nextHop, dest)

	/*
	* print contents of the db
	 */
	rows, _ := database.Query("SELECT LinkID, FromRouter, ToRouter FROM Links")

	var LinkID int
	var FromRouter string
	var ToRouter string
	for rows.Next() {
		//		fmt.Println(idUint32 + ": " + FromRouter + " " + ToRouter)
		rows.Scan(&LinkID, &FromRouter, &ToRouter)
		fmt.Println(strconv.Itoa(LinkID) + ": " + FromRouter + " " + ToRouter)
	}
}
