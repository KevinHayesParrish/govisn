package main

import (
	"database/sql"
	"fmt"
	"hash/crc32"

	//"log"
	"strconv"

	//	"strconv"

	"github.com/g3n/engine/util/logger"
	_ "github.com/mattn/go-sqlite3"
)

//createdsampledbVersion is the file version number
const createsampledbVersion = "0.1.9"

func createsampledb() {
	var log *logger.Logger
	//	fmt.Println("createdampledb version:", createsampledbVersion)
	log.Debug("createdampledb version %s", createsampledbVersion)

	database, _ := sql.Open("sqlite3", "./samplenetwork.db")
	/*
	 *	Add Routers table to DB
	 */
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS Routers (RouterID INTEGER NOT NULL PRIMARY KEY, SystemName TEXT, SystemDesc TEXT, UpTime TEXT, Contact TEXT, Location TEXT, GpsLat REAL, GPSLong REAL, GpsAlt REAL)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO Routers (RouterID, SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")

	// add router to the database
	SystemName := "media"
	RouterIDUint32 := crc32.ChecksumIEEE([]byte(SystemName))
	SystemDesc := "Hardware: x86 Family 15 Model 2 Stepping 7 AT/AT COMPATIBLE - Software: Windows Version 5.2 (Build 3790 Uniprocessor Free)"
	UpTime := "18h 4m 9s 40"
	Contact := "Kevin Parrish"
	Location := "9218 Faxon Place, Elk Grove, CA 95624 USA"
	GpsLat := "38.419471"
	GpsLong := "-121.357212"
	GpsAlt := "10.668"
	statement.Exec(strconv.Itoa(int(RouterIDUint32)), SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) // Add router

	// add another router to the database
	SystemName = "router"
	RouterIDUint32 = crc32.ChecksumIEEE([]byte(SystemName))
	SystemDesc = "LinkSys WRT54G"
	UpTime = "18h 4m 9s 40"
	Contact = "Kevin Parrish"
	Location = "9218 Faxon Place, Elk Grove, CA 95624 USA"
	GpsLat = "38.419492"
	GpsLong = "-121.357176"
	GpsAlt = "10.668"
	statement.Exec(strconv.Itoa(int(RouterIDUint32)), SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) // Add router

	// add another router to the database
	SystemName = "hub"
	RouterIDUint32 = crc32.ChecksumIEEE([]byte(SystemName))
	SystemDesc = "Cisco 2911"
	UpTime = "18h 4m 9s 40"
	Contact = "Kevin Parrish"
	Location = "35 Hub Drive, Melville, NY 11747 USA"
	GpsLat = "40.758262"
	GpsLong = "-73.436362"
	GpsAlt = "10.668"
	statement.Exec(strconv.Itoa(int(RouterIDUint32)), SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) // Add router

	// add another router to the database
	SystemName = "wan-router"
	RouterIDUint32 = crc32.ChecksumIEEE([]byte(SystemName))
	SystemDesc = "Cisco 2911"
	UpTime = "18h 4m 9s 40"
	Contact = "Kevin Parrish"
	Location = "Wichita, KS USA"
	GpsLat = "37.665342"
	GpsLong = "-97.431564"
	GpsAlt = "396"
	statement.Exec(strconv.Itoa(int(RouterIDUint32)), SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) // Add router

	// add another router to the database
	SystemName = "old-country-road"
	RouterIDUint32 = crc32.ChecksumIEEE([]byte(SystemName))
	SystemDesc = "Cisco 2911"
	UpTime = "18h 4m 9s 40"
	Contact = "Kevin Parrish"
	Location = "201 Old Country Road, Melville, NY 11747 USA"
	GpsLat = "40.795017"
	GpsLong = "-73.41144"
	GpsAlt = "13.716"
	statement.Exec(strconv.Itoa(int(RouterIDUint32)), SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) // Add router

	// add another router to the database
	SystemName = "fukui"
	RouterIDUint32 = crc32.ChecksumIEEE([]byte(SystemName))
	SystemDesc = "Cisco 2911"
	UpTime = "18h 4m 9s 40"
	Contact = "Charlie Dowalo"
	Location = "Fukui, Japan"
	GpsLat = "36.063601"
	GpsLong = "136.220856"
	GpsAlt = "17"
	statement.Exec(strconv.Itoa(int(RouterIDUint32)), SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) // Add router

	// add another router to the database
	SystemName = "amsterdam"
	RouterIDUint32 = crc32.ChecksumIEEE([]byte(SystemName))
	SystemDesc = "Cisco 2911"
	UpTime = "18h 4m 9s 40"
	Contact = "Charlie Dowalo"
	Location = "Amsterdam, Netherlands"
	GpsLat = "52.370359"
	GpsLong = "4.894409"
	GpsAlt = "-2"
	statement.Exec(strconv.Itoa(int(RouterIDUint32)), SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) // Add router

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

	/*
	* Add 3DCoordinates table to the database.
	* This table is used dynamically througout the visualization to hold the 3D
	* coordinates of the routers.
	 */
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS Coordinates (RouterID INTEGER, SystemName TEXT, X3D REAL, Y3D REAL, Z3D REAL)")
	if err != nil {
		//		fmt.Println("Error preparing 3DCoordinates table Create statement. Result=", statement)
		//		log.Fatal(err)
		log.Fatal("Error preparing 3DCoordinates table Create statement.")
	}
	statement.Exec()

	statement, err = database.Prepare("INSERT INTO Coordinates (RouterID, SystemName, X3D, Y3D, Z3D) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		//		fmt.Println("Error preparing Coordinates insert statement. Result=", statement)
		//		log.Fatal(err)
		log.Fatal("Error preparing Coordinates insert statement.")
	}
	RouterIDUint32 = 589093411
	SystemName = "wan-router"
	X3D := "0.0"
	Y3D := "0.0"
	Z3D := "0.0"
	statement.Exec(X3D, Y3D, Z3D) // Add 3D Coordinates
}
