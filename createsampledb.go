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
const createsampledbVersion = "0.1.3"

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
	database, _ := sql.Open("sqlite3", "./samplenetwork.db")
	/*
	 *	Add router table to DB
	 */

	/*
	 *	Add RouteTable table to DB
	 */

	/*
	 *	Add RouterIP table to DB
	 */

	/*
	 *	Add RouterMac table to DB
	 */

	/*
	 *	Add links table to DB
	 */
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS links (linkID INTEGER PRIMARY KEY, fromrouter TEXT, torouter TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO links (linkID, fromrouter, torouter) VALUES (?, ?, ?)")

	/*
	 * Add a set of links to the database
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
	 * Add another set of links to the database
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
	 * Add another set of links to the database
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
	 * Add another set of links to the database
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
	 * Add another set of links to the database
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
	 * Add another set of links to the database
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
	rows, _ := database.Query("SELECT linkID, fromrouter, torouter FROM links")

	var linkID int
	var fromrouter string
	var torouter string
	for rows.Next() {
		//		fmt.Println(idUint32 + ": " + fromrouter + " " + torouter)
		rows.Scan(&linkID, &fromrouter, &torouter)
		fmt.Println(strconv.Itoa(linkID) + ": " + fromrouter + " " + torouter)
	}
}
