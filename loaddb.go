package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

//loaddbVersion is the file version number
const loadbVersion = "0.1.8"

func loaddb(networkXML string) {
	fmt.Println("loaddb version:", loadbVersion)
	fmt.Println("Loading database from XML document", networkXML) // FOR TESTING ONLY

	// The struc which contains all the Routers in the XML input file.
	type Routers struct {
		XMLName xml.Name `xml:"V15N_Discovered_Network"`
		Routers []Router `xml:"Router"`
	}

	// The Router struct, this contains
	// the router's Name, Descriptions, UpTime
	// Contact, Location and GPS Coordinates.
	// It also contains nested structs for Addresses
	// and Neighbor routers.
	type Router struct {
		XMLName xml.Name `xml:"Router"`
		System  struct {
			XMLName     xml.Name `xml:"System"`
			Name        string   `xml:"Name"`
			Description string   `xml:"Description"`
			UpTime      string   `xml:"Up_Time"`
			Contact     string   `xml:"Contact"`
			Location    string   `xml:"Location"`
			GPS         struct {
				XMLName   xml.Name `xml:"GPS"`
				Latitude  string   `xml:"Latitude"`
				Longitude string   `xml:"Longitude"`
				Altitude  string   `xml:"Altitude"`
			} `xml:"GPS"`
		} `xml:"System"`
		Addresses struct {
			XMLName          xml.Name `xml:"Addresses"`
			NetworkAddresses struct {
				XMLName   xml.Name `xml:"Network_Addresses"`
				IPAddress []string `xml:"IP_Address"`
			} `xml:"Network_Addresses"`
			MediaAddresses struct {
				XMLName      xml.Name `xml:"Media_Addresses"`
				MediaAddress []string `xml:"Media_Address"`
			} `xml:"Media_Addresses"`
		} `xml:"Addresses"`
		Neighbors struct {
			XMLName  xml.Name `xml:"Neighbors"`
			Neighbor []struct {
				XMLName            xml.Name `xml:"Neighbor"`
				DestinationAddress string   `xml:"Destination_Address"`
				NextHop            string   `xml:"Next_Hop"`
			} `xml:"Neighbor"`
		} `xml:"Neighbors"`
	}

	// Open our xmlFile
	xmlFile, err := os.Open(networkXML)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened", networkXML)
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	xmlFileBytes, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("xmlFileBytes=\n" + string(xmlFileBytes)) // TESTING ONLY

	// Initialize the routers array
	var routers Routers

	// Unmarshal our byteArray which contains our discovered network
	err = xml.Unmarshal(xmlFileBytes, &routers)
	if err != nil {
		fmt.Println(err)
		//		return
		panic(err)
	}
	fmt.Println("routers=", routers) //TESTING ONLY

	fmt.Println("routers length=", len(routers.Routers)) // TESTING ONLY

	// Open the database
	databaseName := *DbName + ".db"
	fmt.Println("dabaseName=", databaseName) //TESTING ONLY
	database, _ := sql.Open("sqlite3", databaseName)

	// Create Routers DB table
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS Routers (RouterID INTEGER NOT NULL PRIMARY KEY, SystemName TEXT, SystemDesc TEXT, UpTime TEXT, Contact TEXT, Location TEXT, GpsLat REAL, GPSLong REAL, GpsAlt REAL)")
	statement.Exec()

	// Add Routers to the database
	for i := 0; i < len(routers.Routers); i++ {
		statement, _ = database.Prepare("INSERT INTO Routers (RouterID, SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
		fmt.Println("i=", i)                                               // TESTING ONLY
		fmt.Println("Router Name: " + routers.Routers[i].System.Name)      // TESTING ONLY
		fmt.Println("Description=", routers.Routers[i].System.Description) // TESTING ONLY
		fmt.Println("Up_Time=", routers.Routers[i].System.UpTime)          // TESTING ONLY
		fmt.Println("GPS=", routers.Routers[i].System.GPS)                 // TESTING ONLY
		fmt.Println("Addresses=", routers.Routers[i].Addresses)            // TESTING ONLY
		fmt.Println("Neighbors=", routers.Routers[i].Neighbors)            // TESTING ONLY

		SystemName := routers.Routers[i].System.Name
		RouterIDUint32 := crc32.ChecksumIEEE([]byte(SystemName))
		SystemDesc := routers.Routers[i].System.Description
		UpTime := routers.Routers[i].System.UpTime
		Contact := routers.Routers[i].System.Contact
		Location := routers.Routers[i].System.Location
		GpsLat := routers.Routers[i].System.GPS.Latitude
		GpsLong := routers.Routers[i].System.GPS.Longitude
		GpsAlt := routers.Routers[i].System.GPS.Altitude
		statement.Exec(strconv.Itoa(int(RouterIDUint32)), SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt) // Add router

		// Create RouterIp DB
		statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS RouterIp (RouterID INTEGER, IPAddr TEXT)")
		statement.Exec()
		statement, _ = database.Prepare("INSERT INTO RouterIp (RouterID, IPAddr) VALUES (?, ?)")

		// Add IP Addresses to current router
		for j := 0; j < len(routers.Routers[i].Addresses.NetworkAddresses.IPAddress); j++ {
			IPAddr := routers.Routers[i].Addresses.NetworkAddresses.IPAddress[j]
			statement.Exec(strconv.Itoa(int(RouterIDUint32)), IPAddr) // Add router
		}

		//	Create RouterMac DB table
		statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS RouterMac (RouterID INTEGER NOT NULL, MacAddr TEXT)")
		statement.Exec()
		statement, _ = database.Prepare("INSERT INTO RouterMac (RourterID, MacAddr) VALUES (?, ?)")

		// Add Media Addresses to current router
		for k := 0; k < len(routers.Routers[i].Addresses.MediaAddresses.MediaAddress); k++ {
			MediaAddr := routers.Routers[i].Addresses.MediaAddresses.MediaAddress[k]
			statement.Exec(strconv.Itoa(int(RouterIDUint32)), MediaAddr) // Add router
		}
	}
	//	Create Links DB table
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS Links (LinkID INTEGER PRIMARY KEY, FromRouter TEXT, ToRouter TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO Links (LinkID, FromRouter, ToRouter) VALUES (?, ?, ?)")

	/*
		database, _ := sql.Open("sqlite3", "./networkXML")

		//*	Add Routers table to DB
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

		//*	Add RouteTable table to DB
		statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS RouteTable (RouterID INTEGER, DestAddr TEXT, NextHop TEXT)")
		statement.Exec()
		statement, _ = database.Prepare("INSERT INTO RouteTable (RouterID, DestAddr, NextHop) VALUES (?, ?, ?)")

		//*	Add RouterIP table to DB
		statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS RouterIp (RouterID INTEGER, IpAddr TEXT)")
		statement.Exec()
		statement, _ = database.Prepare("INSERT INTO RouterIp (RouterID, IpAddr) VALUES (?, ?)")

		//*	Add RouterMac table to DB
		statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS RouterMac (RouterID INTEGER NOT NULL, MacAddr TEXT)")
		statement.Exec()
		statement, _ = database.Prepare("INSERT INTO RouterMac (RourterID, MacAddr) VALUES (?, ?)")

		//*	Add Links table to DB
		statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS Links (LinkID INTEGER PRIMARY KEY, FromRouter TEXT, ToRouter TEXT)")
		statement.Exec()
		statement, _ = database.Prepare("INSERT INTO Links (LinkID, FromRouter, ToRouter) VALUES (?, ?, ?)")

		//* Add a set of Links to the database
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

		//* Add another set of Links to the database
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

		//* Add another set of Links to the database
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

		//* Add another set of Links to the database
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

		//* Add another set of Links to the database
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

		//* Add another set of Links to the database
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

		//* Print contents of the db
		rows, _ := database.Query("SELECT LinkID, FromRouter, ToRouter FROM Links")

		var LinkID int
		var FromRouter string
		var ToRouter string
		for rows.Next() {
			//		fmt.Println(idUint32 + ": " + FromRouter + " " + ToRouter)
			rows.Scan(&LinkID, &FromRouter, &ToRouter)
			fmt.Println(strconv.Itoa(LinkID) + ": " + FromRouter + " " + ToRouter)
		}

		//* Add 3DCoordinates table to the database.
		//* This table is used dynamically througout the visualization to hold the 3D
		//* coordinates of the routers.
		statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS Coordinates (RouterID INTEGER, SystemName TEXT, X3D REAL, Y3D REAL, Z3D REAL)")
		if err != nil {
			fmt.Println("Error preparing 3DCoordinates table Create statement. Result=", statement)
			log.Fatal(err)
		}
		statement.Exec()

		statement, err = database.Prepare("INSERT INTO Coordinates (RouterID, SystemName, X3D, Y3D, Z3D) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			fmt.Println("Error preparing Coordinates insert statement. Result=", statement)
			log.Fatal(err)
		}
		RouterIDUint32 = 589093411
		SystemName = "wan-router"
		X3D := "0.0"
		Y3D := "0.0"
		Z3D := "0.0"
		statement.Exec(X3D, Y3D, Z3D) // Add 3D Coordinates
	*/
}
